package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/richargh/workdown/internal/kernel"
	"github.com/richargh/workdown/internal/kernel/env"
	"github.com/richargh/workdown/internal/plugins/jira"
)

func TestRootCommandPrintsHelp(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	root := New(environment, jira.New)
	// when
	stdout, err := runCommand(t, root)
	// then
	require.NoError(t, err)
	require.Contains(t, stdout, "Usage:")
	require.Contains(t, stdout, "plugins")
}

func TestRootCommandHelpFlagPrintsHelp(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	root := New(environment, jira.New)
	// when
	stdout, err := runCommand(t, root, "-h")
	// then
	require.NoError(t, err)
	require.Contains(t, stdout, "Usage:")
	require.Contains(t, stdout, "plugins")
}

func TestVersionFlagsPrintRuntimeVersion(t *testing.T) {
	for _, arg := range []string{"-v", "--version"} {
		// given
		environment := env.NewMemEnv()
		environment.Version = kernel.ParseWorkdownVersion("0.0.0")
		root := New(environment, jira.New)
		// when
		stdout, err := runCommand(t, root, arg)
		// then
		require.NoError(t, err)
		require.Equal(t, "0.0.0\n", stdout)
	}
}

func TestVersionFallsBackToUnknownVersion(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	root := New(environment, jira.New)
	// when
	stdout, err := runCommand(t, root, "--version")
	// then
	require.NoError(t, err)
	require.Equal(t, "?.?.?\n", stdout)
}

func TestPluginsCommandListsRegisteredPlugins(t *testing.T) {
	// given
	var stdout bytes.Buffer
	environment := env.NewMemEnv()
	environment.Stdout = &stdout
	root := New(environment, jira.New)
	// when
	_, err := runCommand(t, root, "plugins")
	// then
	require.NoError(t, err)
	require.Equal(t, "jira\n", stdout.String())
}
