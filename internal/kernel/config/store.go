package config

import "context"

// Store persists opaque configuration bytes. Format ownership belongs to callers.
type Store interface {
	Load(ctx context.Context, name string) ([]byte, error)
	Save(ctx context.Context, name string, data []byte) error
}
