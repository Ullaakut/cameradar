#!/usr/bin/env python3
"""Format informational Go coverage comments for GitHub Actions."""

from __future__ import annotations

import argparse
import json
import subprocess
import sys
from pathlib import Path
from typing import Any, Callable, Mapping, Sequence

MARKER = "<!-- cameradar-coverage-comment -->"


class CommentPermissionError(RuntimeError):
    """Raised when GITHUB_TOKEN cannot write pull request comments."""


def parse_coverprofile(path: Path) -> float | None:
    """Return statement coverage percentage, or None if the profile is missing/empty."""
    if not path.is_file():
        return None

    total = 0
    covered = 0
    with path.open(encoding="utf-8") as handle:
        for raw in handle:
            line = raw.strip()
            if not line or line.startswith("mode:"):
                continue
            # Format: name.go:line.col,line.col numStatements count
            try:
                rest = line.split()
                statements = int(rest[-2])
                count = int(rest[-1])
            except (ValueError, IndexError):
                continue
            total += statements
            if count > 0:
                covered += statements

    if total == 0:
        return None
    return 100.0 * covered / total


def format_percent(value: float | None) -> str:
    if value is None:
        return "unavailable"
    return f"{value:.1f}%"


def format_delta(head: float | None, base: float | None) -> str:
    if head is None or base is None:
        return "n/a"
    delta = head - base
    sign = "+" if delta >= 0 else ""
    return f"{sign}{delta:.1f}%"


def render_markdown(
    *,
    head: float | None,
    base: float | None,
    base_ref: str,
    head_sha: str,
    fork: bool = False,
) -> str:
    sha = head_sha[:7] if head_sha else "unknown"
    base_label = base_ref or "base"
    lines = [
        MARKER,
        "## Go test coverage",
        "",
        "Informational only — coverage changes do not block merge.",
        "",
        "| | Statements |",
        "| --- | ---: |",
        f"| This PR (`{sha}`) | {format_percent(head)} |",
        f"| `{base_label}` | {format_percent(base)} |",
        f"| Delta | {format_delta(head, base)} |",
        "",
    ]
    if head is None:
        lines.append("No coverage profile was produced for this commit.")
        lines.append("")
    elif base is None:
        lines.append(
            f"Base coverage for `{base_label}` was unavailable, so no delta could be computed."
        )
        lines.append("")
    if fork:
        lines.append(
            "This run came from a fork, so the workflow skipped updating a sticky PR comment."
        )
        lines.append("")
    return "\n".join(lines)


def find_comment(comments: Sequence[Mapping[str, Any]], marker: str) -> Mapping[str, Any] | None:
    for comment in comments:
        body = comment.get("body") or ""
        if marker in body:
            return comment
    return None


class GhClient:
    """Thin wrapper around `gh api` for issue comments."""

    def __init__(self, run: Callable[..., subprocess.CompletedProcess[str]] | None = None) -> None:
        self._run = run or subprocess.run

    def _api(self, args: list[str], *, input_text: str | None = None) -> str:
        cmd = ["gh", "api", *args]
        try:
            completed = self._run(
                cmd,
                check=False,
                text=True,
                capture_output=True,
                input=input_text,
            )
        except FileNotFoundError as exc:
            raise CommentPermissionError("gh is not available") from exc
        if completed.returncode != 0:
            combined = f"{completed.stderr or ''}{completed.stdout or ''}"
            if "403" in combined or "Resource not accessible by integration" in combined:
                raise CommentPermissionError(combined.strip() or "forbidden")
            raise RuntimeError(combined.strip() or f"gh api failed with {completed.returncode}")
        return completed.stdout or ""

    def list_comments(self, repo: str, pr: int) -> list[dict[str, Any]]:
        raw = self._api(["--paginate", f"repos/{repo}/issues/{pr}/comments"])
        if not raw.strip():
            return []
        parsed = json.loads(raw)
        if isinstance(parsed, list):
            return parsed
        return [parsed]

    def create_comment(self, repo: str, pr: int, body: str) -> None:
        payload = json.dumps({"body": body})
        self._api(
            [
                "--method",
                "POST",
                f"repos/{repo}/issues/{pr}/comments",
                "--input",
                "-",
            ],
            input_text=payload,
        )

    def update_comment(self, repo: str, comment_id: int, body: str) -> None:
        payload = json.dumps({"body": body})
        self._api(
            [
                "--method",
                "PATCH",
                f"repos/{repo}/issues/comments/{comment_id}",
                "--input",
                "-",
            ],
            input_text=payload,
        )


def upsert_pr_comment(
    *,
    repo: str,
    pr: int,
    body: str,
    marker: str,
    github: GhClient,
) -> str:
    existing = find_comment(github.list_comments(repo, pr), marker)
    if existing is not None:
        github.update_comment(repo, int(existing["id"]), body)
        return "updated"
    github.create_comment(repo, pr, body)
    return "created"


def write_summary(path: Path | None, body: str) -> None:
    if path is None:
        return
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("a", encoding="utf-8") as handle:
        handle.write(body)
        if not body.endswith("\n"):
            handle.write("\n")


def parse_args(argv: Sequence[str] | None = None) -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--head-profile", type=Path, required=True)
    parser.add_argument("--base-profile", type=Path)
    parser.add_argument("--base-ref", default="")
    parser.add_argument("--head-sha", default="")
    parser.add_argument("--repo", default="")
    parser.add_argument("--pr", type=int)
    parser.add_argument("--summary-file", type=Path)
    parser.add_argument("--output-file", type=Path)
    parser.add_argument("--post-comment", action="store_true")
    parser.add_argument("--fork", action="store_true")
    parser.add_argument("--marker", default=MARKER)
    return parser.parse_args(argv)


def main(argv: Sequence[str] | None = None) -> int:
    args = parse_args(argv)
    head = parse_coverprofile(args.head_profile)
    base = parse_coverprofile(args.base_profile) if args.base_profile else None
    body = render_markdown(
        head=head,
        base=base,
        base_ref=args.base_ref,
        head_sha=args.head_sha,
        fork=args.fork,
    )
    if args.marker != MARKER:
        body = body.replace(MARKER, args.marker, 1)

    if args.output_file is not None:
        args.output_file.write_text(body, encoding="utf-8")
    write_summary(args.summary_file, body)
    sys.stdout.write(body)
    if not body.endswith("\n"):
        sys.stdout.write("\n")

    if not args.post_comment or args.fork or args.pr is None or not args.repo:
        if args.fork:
            print("Skipping sticky PR comment on fork pull request.", file=sys.stderr)
        return 0

    try:
        action = upsert_pr_comment(
            repo=args.repo,
            pr=args.pr,
            body=body,
            marker=args.marker,
            github=GhClient(),
        )
    except CommentPermissionError as exc:
        print(f"Skipping sticky PR comment: {exc}", file=sys.stderr)
        return 0
    print(f"Sticky PR coverage comment {action}.", file=sys.stderr)
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
