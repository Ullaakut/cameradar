## Contributing

Thanks for helping improve Cameradar.
Please keep changes focused and aligned with the project goals.

## Development setup

- Go 1.25 or later
- Docker (optional, for container testing)

Clone the repo and install dependencies using Go modules.

```bash
go mod download
```

### Test against fake targets

Use the following options when you want reproducible local testing.

#### Testing discovery behavior

Use `scanme.nmap.org` to validate discovery-related behavior.

- `scanme.nmap.org` does not expose RTSP or RTSPS ports.
- Target its open ports (for example `22`, `80`, `9929`, `31337`) to test discovery flow, reporting, and scan handling.

Example command:

```bash
cameradar -t scanme.nmap.org -p 22
```

#### Testing RTSP and attack behavior

Use [RTSPAllTheThings](https://github.com/Ullaakut/RTSPAllTheThings) to test RTSP-specific logic and camera attack flows.

- It supports both basic and digest authentication.
- It behaves like a standards-compliant RTSP camera.

> [!CAUTION]
> It is no longer maintained and has limited camera emulation coverage.

Example command:

```bash
docker run --net=host -p 8554:8554 -e RTSP_USERNAME=admin -e RTSP_PASSWORD=12345 -e RTSP_PORT=8554 -e RTSP_AUTHENTICATION_METHOD=digest ullaakut/rtspatt
```

Many real cameras slightly diverge from strict RTSP behavior. For example, some devices allow `DESCRIBE` without authentication, or return `403` and `404` in an order that differs from strict expectations.
Unfortunately, RTSPATT cannot reproduce those behaviors.

#### Prefer real cameras when possible

The most reliable testing method is running against real cameras and real network conditions.

> [!CAUTION]
> Scan only authorized targets and networks.

## Run tests

```bash
make test
```

## Formatting and linting

Keep code idiomatic and consistent with existing style.
By default, follow the [Uber Go Style Guide](https://github.com/uber-go/guide) and the guidelines from [Effective Go](https://go.dev/doc/effective_go).

```bash
make fmt
```

### Dependency for linting

* golangci-lint
    * see current version defined in `.github/workflows/test.yaml` at `jobs.tests.steps.["Run linter"]`
    * configured in `.golangci.yml`

```bash
make lint
```

## Commit messages and PR titles

Use [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) for commit messages and pull request titles.

- Use the format: `type: subject`
- Write the subject in imperative mood: `add`, `update`, `remove`, `fix`, `refactor`
- Do not use gerunds in subjects: avoid `adding`, `updating`, `removing`

Examples:

- `feat: add RTSP timeout flag`
- `fix: remove duplicate progress line`
- `docs: update commit message guidelines`

## Reporting issues

Use the issue template in [.github/ISSUE_TEMPLATE.md](.github/ISSUE_TEMPLATE.md).
Include the version, environment, and repro steps.
Only scan authorized targets.

## Pull requests

1. Create a feature branch from `master`.
2. Keep PRs focused and small.
3. Update documentation when behavior changes.
4. Add or update tests when possible.
5. Ensure `make test` passes.
6. Try to bring as much test coverage as possible with your changes.
7. Use a Conventional Commit-style PR title with an imperative subject.
