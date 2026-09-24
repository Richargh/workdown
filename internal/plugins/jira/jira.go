package jira

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/richargh/workdown/internal/kernel/env"
	"github.com/richargh/workdown/internal/kernel/plugin"
	"github.com/richargh/workdown/internal/plugins/jira/jiraapi"
)

type apiFactory func(jiraConnection, string) (jiraapi.JiraAPI, error)

type jiraConnection struct {
	URL string
}

type jiraPlugin struct {
	environment env.Env
	newAPI      apiFactory
}

func New(environment env.Env) plugin.RemotePlugin {
	return newPlugin(environment, defaultAPIFactory(environment))
}

func NewWithAPI(environment env.Env, api jiraapi.JiraAPI) plugin.RemotePlugin {
	return newPlugin(environment, func(jiraConnection, string) (jiraapi.JiraAPI, error) { return api, nil })
}

func newPlugin(environment env.Env, factory apiFactory) plugin.RemotePlugin {
	return &jiraPlugin{environment: environment, newAPI: factory}
}

func (p *jiraPlugin) Name() string {
	return "jira"
}

func (p *jiraPlugin) Commands() []*cobra.Command {
	return []*cobra.Command{p.newJiraCommand()}
}

func (p *jiraPlugin) newJiraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jira",
		Short: "Jira provider commands",
	}
	cmd.AddCommand(p.newCheckCommand())
	return cmd
}

func (p *jiraPlugin) newCheckCommand() *cobra.Command {
	var jiraURL string
	var pat string
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Verify Jira connectivity and PAT validity",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if pat == "" {
				return fmt.Errorf("--pat is required")
			}
			connection, err := newJiraConnection(jiraURL)
			if err != nil {
				return err
			}
			api, err := p.newAPI(connection, pat)
			if err != nil {
				return err
			}
			if err := api.ValidateConnection(cmd.Context()); err != nil {
				return err
			}
			serverInfo, err := api.ServerInfo(cmd.Context())
			if err != nil {
				return err
			}
			out := p.environment.Stdout
			if out == nil {
				out = cmd.OutOrStdout()
			}
			baseURL := serverInfo.BaseURL()
			_, err = fmt.Fprintf(out, "jira %s ok (version %s)\n", baseURL.String(), serverInfo.Version().String())
			return err
		},
	}
	cmd.Flags().StringVar(&jiraURL, "url", "", "Jira base URL")
	cmd.Flags().StringVar(&pat, "pat", "", "Jira personal access token")
	return cmd
}

func newJiraConnection(jiraURL string) (jiraConnection, error) {
	if jiraURL == "" {
		return jiraConnection{}, fmt.Errorf("--url is required")
	}
	baseURL, err := url.Parse(jiraURL)
	if err != nil {
		return jiraConnection{}, fmt.Errorf("--url: %w", err)
	}
	if baseURL.Scheme == "" || baseURL.Host == "" {
		return jiraConnection{}, fmt.Errorf("--url must be absolute")
	}
	return jiraConnection{URL: jiraURL}, nil
}

func defaultAPIFactory(environment env.Env) apiFactory {
	return func(connection jiraConnection, pat string) (jiraapi.JiraAPI, error) {
		baseURL, err := url.Parse(connection.URL)
		if err != nil {
			return nil, err
		}
		return jiraapi.NewAuthenticatedRestClient(baseURL, environment.HTTPClient, pat)
	}
}
