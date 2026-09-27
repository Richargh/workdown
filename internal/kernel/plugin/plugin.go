package plugin

import (
	"context"

	"github.com/richargh/workdown/internal/kernel"
	"github.com/richargh/workdown/internal/kernel/env"
)

// Constructor creates a remote plugin using the host environment.
type Constructor func(env.Env) RemotePlugin

type RemotePlugin interface {
	Name() string
}

type PullRequest struct {
	URL       string
	PAT       string
	Project   string
	IssueKeys []string
	Mine      bool
	Limit     int
}

type PullResult struct {
	WorkItems []kernel.WorkItem
}

type WorkItemPuller interface {
	PullWorkItems(ctx context.Context, request PullRequest) (PullResult, error)
}

type Registry interface {
	Register(RemotePlugin) error
	Get(name string) (RemotePlugin, bool)
	All() []RemotePlugin
}
