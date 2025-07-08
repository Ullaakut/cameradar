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

## Run tests

```bash
make test
```

## Formatting and linting

Run `gofmt` on changed files.
Keep code idiomatic and consistent with existing style.

```bash
make fmt
```

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
