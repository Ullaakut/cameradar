#!/usr/bin/env python3
"""Tests for informational Go coverage comment rendering."""

from __future__ import annotations

import json
import subprocess
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

import coverage_report


class ParseCoverprofileTest(unittest.TestCase):
    def test_missing_file_returns_none(self) -> None:
        self.assertIsNone(coverage_report.parse_coverprofile(Path("/no/such/coverage.out")))

    def test_empty_profile_returns_none(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            path = Path(tmp) / "coverage.out"
            path.write_text("mode: atomic\n", encoding="utf-8")
            self.assertIsNone(coverage_report.parse_coverprofile(path))

    def test_computes_statement_percentage(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            path = Path(tmp) / "coverage.out"
            path.write_text(
                "\n".join(
                    [
                        "mode: atomic",
                        "example.com/mod/a.go:1.1,2.2 8 1",
                        "example.com/mod/b.go:3.1,4.2 2 0",
                        "",
                    ]
                ),
                encoding="utf-8",
            )
            self.assertAlmostEqual(coverage_report.parse_coverprofile(path), 80.0)


class FormatTest(unittest.TestCase):
    def test_percent_and_delta(self) -> None:
        self.assertEqual(coverage_report.format_percent(81.23), "81.2%")
        self.assertEqual(coverage_report.format_percent(None), "unavailable")
        self.assertEqual(coverage_report.format_delta(81.2, 80.1), "+1.1%")
        self.assertEqual(coverage_report.format_delta(79.0, 80.0), "-1.0%")
        self.assertEqual(coverage_report.format_delta(80.0, None), "n/a")


class RenderMarkdownTest(unittest.TestCase):
    def test_includes_marker_and_informational_note(self) -> None:
        body = coverage_report.render_markdown(
            head=81.2,
            base=80.1,
            base_ref="master",
            head_sha="abcdef123456",
        )
        self.assertIn(coverage_report.MARKER, body)
        self.assertIn("do not block merge", body)
        self.assertIn("This PR (`abcdef1`)", body)
        self.assertIn("`master`", body)
        self.assertIn("+1.1%", body)

    def test_unavailable_base_and_fork_note(self) -> None:
        body = coverage_report.render_markdown(
            head=50.0,
            base=None,
            base_ref="main",
            head_sha="deadbeef",
            fork=True,
        )
        self.assertIn("unavailable", body)
        self.assertIn("n/a", body)
        self.assertIn("fork", body)


class FakeGh:
    def __init__(self, comments: list[dict] | None = None) -> None:
        self.comments = comments or []
        self.created: list[tuple[str, int, str]] = []
        self.updated: list[tuple[str, int, str]] = []

    def list_comments(self, repo: str, pr: int) -> list[dict]:
        return self.comments

    def create_comment(self, repo: str, pr: int, body: str) -> None:
        self.created.append((repo, pr, body))

    def update_comment(self, repo: str, comment_id: int, body: str) -> None:
        self.updated.append((repo, comment_id, body))


class UpsertCommentTest(unittest.TestCase):
    def test_creates_when_missing(self) -> None:
        github = FakeGh()
        action = coverage_report.upsert_pr_comment(
            repo="org/repo",
            pr=7,
            body="<!-- cameradar-coverage-comment -->\nhello",
            marker=coverage_report.MARKER,
            github=github,
        )
        self.assertEqual(action, "created")
        self.assertEqual(len(github.created), 1)
        self.assertEqual(github.updated, [])

    def test_updates_existing_marker(self) -> None:
        github = FakeGh(
            comments=[
                {"id": 11, "body": "unrelated"},
                {"id": 22, "body": "<!-- cameradar-coverage-comment -->\nold"},
            ]
        )
        action = coverage_report.upsert_pr_comment(
            repo="org/repo",
            pr=7,
            body="<!-- cameradar-coverage-comment -->\nnew",
            marker=coverage_report.MARKER,
            github=github,
        )
        self.assertEqual(action, "updated")
        self.assertEqual(github.created, [])
        self.assertEqual(github.updated, [("org/repo", 22, "<!-- cameradar-coverage-comment -->\nnew")])


class GhClientTest(unittest.TestCase):
    def test_403_becomes_permission_error(self) -> None:
        def run(*_args, **_kwargs):
            return subprocess.CompletedProcess(
                args=["gh"],
                returncode=1,
                stdout="",
                stderr="gh: Resource not accessible by integration (HTTP 403)",
            )

        client = coverage_report.GhClient(run=run)
        with self.assertRaises(coverage_report.CommentPermissionError):
            client.list_comments("org/repo", 1)

    def test_list_comments_parses_json_array(self) -> None:
        def run(*_args, **_kwargs):
            return subprocess.CompletedProcess(
                args=["gh"],
                returncode=0,
                stdout=json.dumps([{"id": 1, "body": "hi"}]),
                stderr="",
            )

        client = coverage_report.GhClient(run=run)
        self.assertEqual(client.list_comments("org/repo", 1), [{"id": 1, "body": "hi"}])


class MainTest(unittest.TestCase):
    def test_writes_summary_and_skips_fork_comment(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            head = root / "coverage.out"
            head.write_text("mode: set\nexample.com/a.go:1.1,2.2 1 1\n", encoding="utf-8")
            summary = root / "summary.md"
            output = root / "comment.md"
            with patch.object(coverage_report, "upsert_pr_comment") as upsert:
                code = coverage_report.main(
                    [
                        "--head-profile",
                        str(head),
                        "--base-ref",
                        "master",
                        "--head-sha",
                        "abc1234",
                        "--repo",
                        "org/repo",
                        "--pr",
                        "9",
                        "--summary-file",
                        str(summary),
                        "--output-file",
                        str(output),
                        "--post-comment",
                        "--fork",
                    ]
                )
            self.assertEqual(code, 0)
            upsert.assert_not_called()
            body = output.read_text(encoding="utf-8")
            self.assertIn(coverage_report.MARKER, body)
            self.assertIn(coverage_report.MARKER, summary.read_text(encoding="utf-8"))

    def test_permission_error_does_not_fail(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            head = Path(tmp) / "coverage.out"
            head.write_text("mode: set\nexample.com/a.go:1.1,2.2 1 1\n", encoding="utf-8")
            with patch.object(
                coverage_report,
                "upsert_pr_comment",
                side_effect=coverage_report.CommentPermissionError("403"),
            ):
                code = coverage_report.main(
                    [
                        "--head-profile",
                        str(head),
                        "--repo",
                        "org/repo",
                        "--pr",
                        "9",
                        "--post-comment",
                    ]
                )
            self.assertEqual(code, 0)


if __name__ == "__main__":
    unittest.main()
