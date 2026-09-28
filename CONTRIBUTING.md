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

## Release

Prepare a release pull request from `trunk`:

```bash
mise run release:pr -- 0.0.3
```

This command:

- creates a `release/<version>` branch
- moves `CHANGELOG.md` `Unreleased` notes to the release version
- updates `VERSION`
- commits the release files
- pushes the branch
- opens a pull request against `trunk`

After the release pull request is merged, CI checks `trunk`. 
Since `VERSION` changed, CI creates the `v<version>` tag and publishes the GitHub release.
Done.

## Project layout

```text
cmd/workdown/                  CLI binary entrypoint
internal/cli/                  Cobra root command and command wiring
internal/kernel/               Shared plugin contracts and ports
internal/plugins/jira/         Jira plugin entrypoint and commands
```

## Commit messages

Use [Risk-Aware Commit Notation](https://github.com/RefactoringCombos/ArlosCommitNotation) for commit messages.

Format:

```text
<RISK> <INTENTION>(<scope>) <short imperative summary>
```

Risk symbols:

- `.` — proven safe. The change addresses known and unknown risks.
- `^` — validated. The change addresses known risks.
- `!` — risky. Some known risks are not verified.
- `@` — probably broken. There is no risk attestation.

Core intentions:

- `F` or `f` — feature. Change or extend one aspect of program behavior.
- `B` or `b` — bugfix. Repair one bad program behavior.
- `R` or `r` — refactoring. Change implementation without behavior change.
- `D` or `d` — documentation. Change information for team members without program behavior change.

Project extension intentions:

- `E` or `e` — environment. Use for dependencies, build tooling, CI, and development setup.
- `A` or `a` — agentic. Use for agent instructions, skills, prompts, and other agent behavior changes.

Use uppercase when reviewers must pay more attention or when the change should be visible in release notes. Use lowercase for internal changes with no user-visible effect.

Use the imperative mood in the summary. Write what the commit does, not what it did.

Good:

```text
^ E(deps): bump cobra to v1.2.3
^ F(jira): add issue pull command
```

Use a multi-line commit message when the summary is not enough to explain the reason, risk, or validation. Keep the first line short. Add a blank line. Then add details.

Example:

```text
^ F(jira): add issue pull command

Developers need an offline copy of assigned Jira issues before they start work.
This command creates the local Markdown files that later commands can read and edit.
```

If useful, use the body to explain:

- why the change is needed
- what risk remains
- how you tested it
- related issue or pull request IDs

Keep commit messages short and clear.

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
