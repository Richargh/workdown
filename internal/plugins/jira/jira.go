package jira

import (
	"github.com/spf13/cobra"

	"github.com/richargh/workdown/internal/kernel/env"
	"github.com/richargh/workdown/internal/kernel/plugin"
)

type jiraPlugin struct {
	environment env.Env
}

func New(environment env.Env) plugin.RemotePlugin {
	return &jiraPlugin{environment: environment}
}

func (p *jiraPlugin) Name() string {
	return "jira"
}

func (p *jiraPlugin) Commands() []*cobra.Command {
	return nil
}
