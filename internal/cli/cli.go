package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/richargh/workdown/internal/kernel/env"
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
