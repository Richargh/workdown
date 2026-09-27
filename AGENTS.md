Always start replies with STARTER_CHARACTER + space (default: 📋). Stack emojis, don't replace.

# Agent Instructions

Always use ASD-STE100 Simplified Technical English.

## Validation

- Before reporting completion after code changes, run `mise run check` from the repository root.
- If `mise run check` cannot complete because of an environment/tooling issue, say so explicitly and include the error reason.
- Do not substitute `go test ./...` for `mise run check`; it is okay to run `go test ./...` as a faster intermediate check, but final validation should be `mise run check`.

## Boundaries

- `cmd/workdown/`: binary entrypoint only.
- `internal/cli/`: host/root CLI wiring.
- `internal/kernel/`: shared contracts and ports; no Cobra or provider details.
- `internal/plugins/<provider>/`: provider service/domain behavior.
- `internal/plugins/<provider>/<provider>cli/`: provider-owned Cobra adapters.
- `internal/plugins/<provider>/<provider>api/`: provider API contracts/implementations.