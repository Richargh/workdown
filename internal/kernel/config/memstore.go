package config

import (
	"context"
	"fmt"
	"sync"
)

type memStore struct {
	mu      sync.RWMutex
	configs map[string][]byte
}

func NewMemStore() Store {
	return &memStore{configs: make(map[string][]byte)}
}

func (s *memStore) Load(ctx context.Context, name string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, ok := s.configs[name]
	if !ok {
		return nil, fmt.Errorf("load config %q: config does not exist", name)
	}
	return append([]byte(nil), data...), nil
}

func (s *memStore) Save(ctx context.Context, name string, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.configs[name] = append([]byte(nil), data...)
	return nil
}
