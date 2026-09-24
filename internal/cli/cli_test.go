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
	require.NotContains(t, stdout, "plugins")
	require.Contains(t, stdout, "remotes")
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
	require.NotContains(t, stdout, "plugins")
	require.Contains(t, stdout, "remotes")
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

func TestRemotesCommandListsRegisteredRemoteProviders(t *testing.T) {
	// given
	var stdout bytes.Buffer
	environment := env.NewMemEnv()
	environment.Stdout = &stdout
	root := New(environment, jira.New)
	// when
	_, err := runCommand(t, root, "remotes")
	// then
	require.NoError(t, err)
	require.Equal(t, "jira\n", stdout.String())
}

func TestProviderCommandsLiveUnderRemotes(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	root := New(environment, jira.New)

	// when
	stdout, err := runCommand(t, root, "remotes", "jira", "--help")

	// then
	require.NoError(t, err)
	require.Contains(t, stdout, "Jira provider commands")
}

func TestProviderCommandsAreNotRootCommands(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	root := New(environment, jira.New)

	// when
	_, err := runCommand(t, root, "jira", "--help")

	// then
	require.Error(t, err)
}

func TestPluginsCommandDoesNotExist(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	root := New(environment, jira.New)

	// when
	_, err := runCommand(t, root, "plugins")

	// then
	require.Error(t, err)
}
