package cli

import (
	"bytes"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/richargh/workdown/internal/kernel"
	"github.com/richargh/workdown/internal/kernel/env"
	kernelplugin "github.com/richargh/workdown/internal/kernel/plugin"
	"github.com/richargh/workdown/internal/plugins/jira"
	"github.com/richargh/workdown/internal/plugins/jira/jiraapi"
	"github.com/richargh/workdown/internal/plugins/jira/jiracli"
)

func TestRootCommandPrintsHelp(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	root := New(environment, jiracli.New)
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
	root := New(environment, jiracli.New)
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
		root := New(environment, jiracli.New)
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
	root := New(environment, jiracli.New)
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
	root := New(environment, jiracli.New)
	// when
	_, err := runCommand(t, root, "remotes")
	// then
	require.NoError(t, err)
	require.Equal(t, "jira\n", stdout.String())
}

func TestProviderCommandsLiveUnderRemotes(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	root := New(environment, jiracli.New)

	// when
	stdout, err := runCommand(t, root, "remotes", "jira", "--help")

	// then
	require.NoError(t, err)
	require.Contains(t, stdout, "Jira provider commands")
}

func TestProviderCommandsAreNotRootCommands(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	root := New(environment, jiracli.New)

	// when
	_, err := runCommand(t, root, "jira", "--help")

	// then
	require.Error(t, err)
}

func TestJiraRemotePullCommandCountsSelectedProjectIssues(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	api := jiraapi.NewIssueSearchMemClient(newServerInfo(t, "https://jira.example.test", "9.12.0"), 4)
	root := New(environment, func(environment env.Env) kernelplugin.RemotePlugin {
		return jiracli.NewWithService(environment, jira.NewWithAPI(environment, api))
	})

	// when
	stdout, err := runCommand(t, root, "remotes", "jira", "pull", "--url", "https://jira.example.test", "--project", "PROJ", "--issues", "PROJ-1,PROJ-2", "--pat", "secret-token")

	// then
	require.NoError(t, err)
	require.Equal(t, "4 issues\n", stdout)
	require.Equal(t, []string{"project = PROJ AND key in (PROJ-1, PROJ-2)"}, jiraapi.SearchIssuesQueries(api))
}

func TestPluginsCommandDoesNotExist(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	root := New(environment, jiracli.New)

	// when
	_, err := runCommand(t, root, "plugins")

	// then
	require.Error(t, err)
}

func newServerInfo(t *testing.T, rawBaseURL string, rawVersion string) jiraapi.ServerInfo {
	t.Helper()
	baseURL, err := url.Parse(rawBaseURL)
	require.NoError(t, err)
	serverInfo, err := jiraapi.NewServerInfo(baseURL, jiraapi.ParseJiraVersion(rawVersion))
	require.NoError(t, err)
	return serverInfo
}
