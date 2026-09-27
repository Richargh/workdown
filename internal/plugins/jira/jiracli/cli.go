package jiracli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/richargh/workdown/internal/kernel/env"
	"github.com/richargh/workdown/internal/kernel/plugin"
	"github.com/richargh/workdown/internal/plugins/jira"
)

type adapter struct {
	environment env.Env
	service     *jira.Plugin
}

func New(environment env.Env) plugin.RemotePlugin {
	return NewWithService(environment, jira.New(environment))
}

func NewWithService(environment env.Env, service *jira.Plugin) plugin.RemotePlugin {
	return &adapter{environment: environment, service: service}
}

func (a *adapter) Name() string {
	return a.service.Name()
}

func (a *adapter) Commands() []*cobra.Command {
	return []*cobra.Command{a.newJiraCommand()}
}

func (a *adapter) newJiraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jira",
		Short: "Jira provider commands",
	}
	cmd.AddCommand(a.newCheckCommand())
	cmd.AddCommand(a.newPullCommand())
	return cmd
}

func (a *adapter) newCheckCommand() *cobra.Command {
	var jiraURL string
	var pat string
	var project string
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Verify Jira connectivity and PAT validity",
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := a.service.Check(cmd.Context(), jira.CheckRequest{URL: jiraURL, PAT: pat, Project: project})
			if err != nil {
				return err
			}
			out := a.environment.Stdout
			if out == nil {
				out = cmd.OutOrStdout()
			}
			if _, err := fmt.Fprintf(out, "jira %s ok (version %s)\n", result.BaseURL, result.Version); err != nil {
				return err
			}
			if result.HasProject() {
				_, err = fmt.Fprintf(out, "project %s ok (%d issues)\n", result.Project, result.ProjectIssueCount)
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&jiraURL, "url", "", "Jira base URL")
	cmd.Flags().StringVar(&project, "project", "", "Jira project key to validate")
	cmd.Flags().StringVar(&pat, "pat", "", "Jira personal access token")
	return cmd
}

func (a *adapter) newPullCommand() *cobra.Command {
	var jiraURL string
	var pat string
	var project string
	var issues string
	var limit int
	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull selected Jira issues as Workdown interchange",
		RunE: func(cmd *cobra.Command, _ []string) error {
			out := a.environment.Stdout
			if out == nil {
				out = cmd.OutOrStdout()
			}
			result, err := a.service.Pull(cmd.Context(), jira.PullRequest{
				URL:       jiraURL,
				PAT:       pat,
				Project:   project,
				IssueKeys: parseIssueKeys(issues),
				Limit:     limit,
			})
			if err != nil {
				return err
			}
			encoder := json.NewEncoder(out)
			encoder.SetIndent("", "  ")
			return encoder.Encode(result.WorkItems)
		},
	}
	cmd.Flags().StringVar(&jiraURL, "url", "", "Jira base URL")
	cmd.Flags().StringVar(&project, "project", "", "Jira project key")
	cmd.Flags().StringVar(&issues, "issues", "", "Comma-separated Jira issue keys to pull")
	cmd.Flags().IntVar(&limit, "limit", 0, "Maximum number of Jira issues to pull")
	cmd.Flags().StringVar(&pat, "pat", "", "Jira personal access token")
	return cmd
}

func parseIssueKeys(issues string) []string {
	var issueKeys []string
	for _, issue := range strings.Split(issues, ",") {
		issue = strings.TrimSpace(issue)
		if issue != "" {
			issueKeys = append(issueKeys, issue)
		}
	}
	return issueKeys
}
