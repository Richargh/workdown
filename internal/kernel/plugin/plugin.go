package plugin

import (
	"github.com/spf13/cobra"

	"github.com/richargh/workdown/internal/kernel/env"
)

// Constructor creates a remote plugin using the host environment.
type Constructor func(env.Env) RemotePlugin

// RemotePlugin is the shared plugin contract implemented by provider plugins.
type RemotePlugin interface {
	Name() string
	Commands() []*cobra.Command
}

// Registry stores remote plugins by provider name.
type Registry interface {
	Register(RemotePlugin) error
	Get(name string) (RemotePlugin, bool)
	All() []RemotePlugin
}
