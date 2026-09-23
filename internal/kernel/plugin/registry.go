package plugin

import (
	"fmt"
	"sort"
)

type registry struct {
	plugins map[string]RemotePlugin
}

func NewRegistry() Registry {
	return &registry{plugins: make(map[string]RemotePlugin)}
}

func (r *registry) Register(remote RemotePlugin) error {
	if remote == nil {
		return fmt.Errorf("register plugin: nil plugin")
	}
	name := remote.Name()
	if name == "" {
		return fmt.Errorf("register plugin: empty name")
	}
	if _, exists := r.plugins[name]; exists {
		return fmt.Errorf("register plugin %q: already registered", name)
	}
	r.plugins[name] = remote
	return nil
}

func (r *registry) Get(name string) (RemotePlugin, bool) {
	remote, ok := r.plugins[name]
	return remote, ok
}

func (r *registry) All() []RemotePlugin {
	names := make([]string, 0, len(r.plugins))
	for name := range r.plugins {
		names = append(names, name)
	}
	sort.Strings(names)

	remotes := make([]RemotePlugin, 0, len(names))
	for _, name := range names {
		remotes = append(remotes, r.plugins[name])
	}
	return remotes
}
