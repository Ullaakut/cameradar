# set this e.g. via `make build GORELEASER_FLAGS="--skip=docker"` for temporary flags
GORELEASER_FLAGS=

#Format

fmt:
	@echo "==> Formatting source"
	@gofmt -s -w $(shell find . -type f -name '*.go')
	@echo "==> Done"
.PHONY: fmt

#Test

test:
	@go test -cover -race ./...
.PHONY: test

#Lint

lint:
	@golangci-lint run --config=.golangci.yml ./...
.PHONY: lint

#Build

build:
	@goreleaser release $(GORELEASER_FLAGS) --clean --snapshot
.PHONY: build
