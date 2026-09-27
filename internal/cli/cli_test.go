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

func TestJiraRemotePullCommandEmitsInterchangeJSON(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	api := jiraapi.NewIssueSearchWithIssuesMemClient(newServerInfo(t, "https://jira.example.test", "9.12.0"), []jiraapi.JiraIssue{
		newJiraIssue("10001", "PROJ-1", "PROJ", "Task", "Write docs"),
	})
	root := New(environment, func(environment env.Env) kernelplugin.RemotePlugin {
		return jiracli.NewWithService(environment, jira.NewWithAPI(environment, api))
	})

	// when
	stdout, err := runCommand(t, root, "remotes", "jira", "pull", "--url", "https://jira.example.test", "--project", "PROJ", "--pat", "secret-token")

	// then
	require.NoError(t, err)
	require.Equal(t, `[
  {
    "remote": "https://jira.example.test",
    "provider": "jira",
    "id": "10001",
    "key": "PROJ-1",
    "url": "https://jira.example.test/browse/PROJ-1",
    "title": "Write docs",
    "fields": [
      {
        "name": "project",
        "providerKey": "project",
        "type": "string",
        "value": "PROJ",
        "editable": false
      },
      {
        "name": "issueType",
        "providerKey": "issuetype",
        "type": "string",
        "value": "Task",
        "editable": false
      }
    ],
    "metadata": {
      "jira.baseURL": "https://jira.example.test"
    }
  }
]
`, stdout)
	require.Equal(t, []string{"project = PROJ"}, jiraapi.SearchIssuesQueries(api))
	require.Equal(t, []int{5}, jiraapi.SearchIssuesLimits(api))
}

func TestJiraRemotePullCommandUsesLimit(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	api := jiraapi.NewIssueSearchWithIssuesMemClient(newServerInfo(t, "https://jira.example.test", "9.12.0"), []jiraapi.JiraIssue{})
	root := New(environment, func(environment env.Env) kernelplugin.RemotePlugin {
		return jiracli.NewWithService(environment, jira.NewWithAPI(environment, api))
	})

	// when
	_, err := runCommand(t, root, "remotes", "jira", "pull", "--url", "https://jira.example.test", "--project", "PROJ", "--limit", "2", "--pat", "secret-token")

	// then
	require.NoError(t, err)
	require.Equal(t, []int{2}, jiraapi.SearchIssuesLimits(api))
}

func TestJiraRemotePullCommandLimitsToAllSelectedIssuesByDefault(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	api := jiraapi.NewIssueSearchWithIssuesMemClient(newServerInfo(t, "https://jira.example.test", "9.12.0"), []jiraapi.JiraIssue{})
	root := New(environment, func(environment env.Env) kernelplugin.RemotePlugin {
		return jiracli.NewWithService(environment, jira.NewWithAPI(environment, api))
	})

	// when
	_, err := runCommand(t, root, "remotes", "jira", "pull", "--url", "https://jira.example.test", "--project", "PROJ", "--issues", "PROJ-1,PROJ-2,PROJ-3,PROJ-4,PROJ-5,PROJ-6", "--pat", "secret-token")

	// then
	require.NoError(t, err)
	require.Equal(t, []int{6}, jiraapi.SearchIssuesLimits(api))
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

func newJiraIssue(id string, key string, project string, issueType string, summary string) jiraapi.JiraIssue {
	return jiraapi.JiraIssue{
		ID:  id,
		Key: key,
		Fields: jiraapi.JiraIssueFields{
			Summary:   summary,
			IssueType: jiraapi.JiraNamedValue{Name: issueType},
			Project:   jiraapi.JiraProject{Key: project},
		},
	}
}

func newServerInfo(t *testing.T, rawBaseURL string, rawVersion string) jiraapi.ServerInfo {
	t.Helper()
	baseURL, err := url.Parse(rawBaseURL)
	require.NoError(t, err)
	serverInfo, err := jiraapi.NewServerInfo(baseURL, jiraapi.ParseJiraVersion(rawVersion))
	require.NoError(t, err)
	return serverInfo
}
