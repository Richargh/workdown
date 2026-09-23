package fs

import (
	"context"
	"fmt"
	"sort"
	"sync"
)

type memStore struct {
	mu    sync.RWMutex
	files map[string][]byte
}

func NewMemStore() Store {
	return &memStore{files: make(map[string][]byte)}
}

func (s *memStore) ReadFile(ctx context.Context, path string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, ok := s.files[path]
	if !ok {
		return nil, fmt.Errorf("read %q: file does not exist", path)
	}
	return append([]byte(nil), data...), nil
}

func (s *memStore) WriteFile(ctx context.Context, path string, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.files[path] = append([]byte(nil), data...)
	return nil
}

func (s *memStore) Exists(ctx context.Context, path string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.files[path]
	return ok, nil
}

func (s *memStore) Paths() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	paths := make([]string, 0, len(s.files))
	for path := range s.files {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths
}
