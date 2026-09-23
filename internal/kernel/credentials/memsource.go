package credentials

import (
	"context"
	"fmt"
	"sync"
)

type memSource struct {
	mu          sync.RWMutex
	credentials map[string]string
}

func NewMemSource(credentials map[string]string) Source {
	credentialCopy := make(map[string]string, len(credentials))
	for remote, credential := range credentials {
		credentialCopy[remote] = credential
	}
	return &memSource{credentials: credentialCopy}
}

func (s *memSource) Credential(ctx context.Context, remote string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	credential, ok := s.credentials[remote]
	if !ok {
		return "", fmt.Errorf("credential for %q: not found", remote)
	}
	return credential, nil
}
