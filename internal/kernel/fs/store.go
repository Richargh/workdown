package fs

import "context"

// Store abstracts workspace file IO.
type Store interface {
	ReadFile(ctx context.Context, path string) ([]byte, error)
	WriteFile(ctx context.Context, path string, data []byte) error
	Exists(ctx context.Context, path string) (bool, error)
}
