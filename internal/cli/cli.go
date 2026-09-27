package cli

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/spf13/cobra"

	"github.com/richargh/workdown/internal/kernel/env"
	"github.com/richargh/workdown/internal/kernel/markdown"
	kernelplugin "github.com/richargh/workdown/internal/kernel/plugin"
)

func New(environment env.Env, pluginConstructors ...kernelplugin.Constructor) *cobra.Command {
	registry := kernelplugin.NewRegistry()
	for _, constructPlugin := range pluginConstructors {
		plugin := constructPlugin(environment)
		if err := registry.Register(plugin); err != nil {
			panic(err)
		}
	}

	var versionFlag bool
	root := &cobra.Command{
		Use:           "workdown",
		Short:         "Sync work items with Markdown",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if versionFlag {
				return printVersion(cmd, environment)
			}
			return cmd.Help()
		},
	}
	root.Flags().BoolVarP(&versionFlag, "version", "v", false, "Print the Workdown version")

	root.AddCommand(newRemotesCommand(environment, registry))
	root.AddCommand(newPullCommand(environment, registry))

	return root
}

func printVersion(cmd *cobra.Command, environment env.Env) error {
	out := environment.Stdout
	if out == nil {
		out = cmd.OutOrStdout()
	}
	_, err := fmt.Fprintln(out, environment.Version.String())
	return err
}

type remoteCommandProvider interface {
	Commands() []*cobra.Command
}

func newPullCommand(environment env.Env, registry kernelplugin.Registry) *cobra.Command {
	var remote string
	var remoteURL string
	var pat string
	var project string
	var issues string
	var mine bool
	var outDir string
	var limit int
	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull work items to Markdown",
		RunE: func(cmd *cobra.Command, _ []string) error {
			registeredPlugin, ok := registry.Get(remote)
			if !ok {
				return fmt.Errorf("remote %q is not registered", remote)
			}
			puller, ok := registeredPlugin.(kernelplugin.WorkItemPuller)
			if !ok {
				return fmt.Errorf("remote %q does not support pull", remote)
			}
			result, err := puller.PullWorkItems(cmd.Context(), kernelplugin.PullRequest{
				URL:       remoteURL,
				PAT:       pat,
				Project:   project,
				IssueKeys: parseIssueKeys(issues),
				Mine:      mine,
				Limit:     limit,
			})
			if err != nil {
				return err
			}
			var writeErrors []error
			for _, item := range result.WorkItems {
				filePath := path.Join(outDir, item.Key+".md")
				exists, err := environment.Files.Exists(cmd.Context(), filePath)
				if err != nil {
					writeErrors = append(writeErrors, err)
					continue
				}
				if exists {
					writeErrors = append(writeErrors, fmt.Errorf("%s already exists", filePath))
					continue
				}
				if err := environment.Files.WriteFile(cmd.Context(), filePath, markdown.Render(item)); err != nil {
					writeErrors = append(writeErrors, err)
				}
			}
			return errors.Join(writeErrors...)
		},
	}
	cmd.Flags().StringVar(&remote, "remote", "jira", "Remote provider name")
	cmd.Flags().StringVar(&remoteURL, "url", "", "Remote base URL")
	cmd.Flags().StringVar(&project, "project", "", "Project key")
	cmd.Flags().StringVar(&issues, "issues", "", "Comma-separated issue keys to pull")
	cmd.Flags().BoolVar(&mine, "mine", false, "Pull issues assigned to the authenticated user")
	cmd.Flags().IntVar(&limit, "limit", 0, "Maximum number of issues to pull")
	cmd.Flags().StringVar(&outDir, "out", "items", "Output directory")
	cmd.Flags().StringVar(&pat, "pat", "", "Personal access token")
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

func newRemotesCommand(environment env.Env, registry kernelplugin.Registry) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remotes",
		Short: "List remote providers and run provider commands",
		RunE: func(cmd *cobra.Command, _ []string) error {
			out := environment.Stdout
			if out == nil {
				out = cmd.OutOrStdout()
			}
			for _, plugin := range registry.All() {
				if _, err := fmt.Fprintln(out, plugin.Name()); err != nil {
					return err
				}
			}
			return nil
		},
	}
	for _, plugin := range registry.All() {
		if provider, ok := plugin.(remoteCommandProvider); ok {
			cmd.AddCommand(provider.Commands()...)
		}
	}
	return cmd
}
