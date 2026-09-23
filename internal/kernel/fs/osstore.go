package fs

import (
	"context"
	"os"
	"path/filepath"
)

type osStore struct {
	root string
}

func NewOSStore(root string) Store {
	return &osStore{root: root}
}

func (s *osStore) ReadFile(ctx context.Context, path string) ([]byte, error) {
	return os.ReadFile(s.fullPath(path))
}

func (s *osStore) WriteFile(ctx context.Context, path string, data []byte) error {
	fullPath := s.fullPath(path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o750); err != nil {
		return err
	}
	return os.WriteFile(fullPath, data, 0o600)
}

func (s *osStore) Exists(ctx context.Context, path string) (bool, error) {
	_, err := os.Stat(s.fullPath(path))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (s *osStore) fullPath(path string) string {
	return filepath.Join(s.root, filepath.Clean(path))
}
