package plugin

import "github.com/richargh/workdown/internal/kernel/env"

// Constructor creates a remote plugin using the host environment.
type Constructor func(env.Env) RemotePlugin

type RemotePlugin interface {
	Name() string
}

type Registry interface {
	Register(RemotePlugin) error
	Get(name string) (RemotePlugin, bool)
	All() []RemotePlugin
}
