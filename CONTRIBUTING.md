# Contributing

## Development setup

Install prerequisites with [mise](https://mise.jdx.dev/):

```bash
# trust this repo's .mise.toml
mise trust

# install tool versions from .mise.toml
mise install

# verify the pinned Go version is active
go version
```

If your shell does not automatically activate mise, prefix commands with `mise exec --`, for example:

```bash
mise exec -- go test ./...
```

## Usage

Run from source:

```bash
go run ./cmd/workdown --help
```

Build a local binary:

```bash
go build -o build/workdown ./cmd/workdown
./build/workdown --help
```

## Test

```bash
# run the standard local checks
mise run check
```

## CLI reference

Regenerate the CLI reference after command changes:

```bash
mise run docs:cli
```

The reference is automatically checked by `mise run check`.

## Project layout

```text
cmd/workdown/                  CLI binary entrypoint
internal/cli/                  Cobra root command and command wiring
internal/kernel/               Shared plugin contracts and ports
internal/plugins/jira/         Jira plugin entrypoint and commands
```

## Dependency management

Go dependencies are tracked in:

- `go.mod` — module name, Go version, and dependency requirements
- `go.sum` — checksums used to verify downloaded dependencies

Useful commands:

```bash
# add or change a direct dependency
go get github.com/spf13/cobra@v1.8.1

# download an already-required dependency and update go.sum
go mod download github.com/spf13/cobra@v1.8.1

# clean up go.mod and go.sum based on actual imports
go mod tidy
```
