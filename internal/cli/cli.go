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

	root.AddCommand(newPluginsCommand(environment, registry))
	for _, plugin := range registry.All() {
		root.AddCommand(plugin.Commands()...)
	}

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

func newPluginsCommand(environment env.Env, registry kernelplugin.Registry) *cobra.Command {
	return &cobra.Command{
		Use:   "plugins",
		Short: "List registered remote plugins",
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
}
