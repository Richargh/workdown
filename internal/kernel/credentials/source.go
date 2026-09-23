package credentials

import "context"

// Source returns credentials without requiring callers to know where they are stored.
type Source interface {
	Credential(ctx context.Context, remote string) (string, error)
}
