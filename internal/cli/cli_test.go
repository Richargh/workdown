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
	api := jiraapi.NewIssueSearchWithIssuesMemClient(newServerInfo(t), []jiraapi.JiraIssue{
		newJiraIssue(),
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
      },
      {
        "name": "status",
        "providerKey": "status",
        "type": "string",
        "value": "To Do",
        "editable": false
      },
      {
        "name": "author",
        "providerKey": "reporter",
        "type": "string",
        "value": "Alice Author",
        "editable": false
      },
      {
        "name": "owner",
        "providerKey": "assignee",
        "type": "string",
        "value": "Bob Owner",
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

func TestJiraRemotePullCommandCanPullMyIssuesAsJSON(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	api := jiraapi.NewIssueSearchWithIssuesMemClient(newServerInfo(t), []jiraapi.JiraIssue{
		newJiraIssue(),
	})
	root := New(environment, func(environment env.Env) kernelplugin.RemotePlugin {
		return jiracli.NewWithService(environment, jira.NewWithAPI(environment, api))
	})

	// when
	stdout, err := runCommand(t, root, "remotes", "jira", "pull", "--url", "https://jira.example.test", "--project", "PROJ", "--mine", "--format", "json", "--pat", "secret-token")

	// then
	require.NoError(t, err)
	require.Contains(t, stdout, `"key": "PROJ-1"`)
	require.Equal(t, []string{"project = PROJ AND assignee = currentUser()"}, jiraapi.SearchIssuesQueries(api))
}

func TestCorePullWritesMarkdown(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	api := jiraapi.NewIssueSearchWithIssuesMemClient(newServerInfo(t), []jiraapi.JiraIssue{
		newJiraIssue(),
	})
	root := New(environment, func(environment env.Env) kernelplugin.RemotePlugin {
		return jiracli.NewWithService(environment, jira.NewWithAPI(environment, api))
	})

	// when
	_, err := runCommand(t, root, "pull", "--url", "https://jira.example.test", "--project", "PROJ", "--out", "issues", "--pat", "secret-token")

	// then
	require.NoError(t, err)
	data, err := environment.Files.ReadFile(t.Context(), "issues/PROJ-1.md")
	require.NoError(t, err)
	require.Equal(t, `+++
remote = "https://jira.example.test"
provider = "jira"
key = "PROJ-1"
id = "10001"
project = "PROJ"
issue_type = "Task"
title_hash = "sha256:345b88bffbe2a2d20b6ae1dd09ef4952cc3d74aef54078d81a503601d9bb1441"
+++

# Write docs {#wd-field-title}
`, string(data))
}

func TestCorePullDoesNotOverwriteExistingMarkdown(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	require.NoError(t, environment.Files.WriteFile(t.Context(), "issues/PROJ-1.md", []byte("local edit")))
	api := jiraapi.NewIssueSearchWithIssuesMemClient(newServerInfo(t), []jiraapi.JiraIssue{
		newJiraIssue(),
	})
	root := New(environment, func(environment env.Env) kernelplugin.RemotePlugin {
		return jiracli.NewWithService(environment, jira.NewWithAPI(environment, api))
	})

	// when
	_, err := runCommand(t, root, "pull", "--url", "https://jira.example.test", "--project", "PROJ", "--out", "issues", "--pat", "secret-token")

	// then
	require.ErrorContains(t, err, "issues/PROJ-1.md already exists")
	data, err := environment.Files.ReadFile(t.Context(), "issues/PROJ-1.md")
	require.NoError(t, err)
	require.Equal(t, "local edit", string(data))
}

func TestCorePullContinuesAfterExistingMarkdown(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	require.NoError(t, environment.Files.WriteFile(t.Context(), "issues/PROJ-1.md", []byte("local edit")))
	api := jiraapi.NewIssueSearchWithIssuesMemClient(newServerInfo(t), []jiraapi.JiraIssue{
		newJiraIssue(),
		newSecondJiraIssue(),
	})
	root := New(environment, func(environment env.Env) kernelplugin.RemotePlugin {
		return jiracli.NewWithService(environment, jira.NewWithAPI(environment, api))
	})

	// when
	_, err := runCommand(t, root, "pull", "--url", "https://jira.example.test", "--project", "PROJ", "--out", "issues", "--pat", "secret-token")

	// then
	require.ErrorContains(t, err, "issues/PROJ-1.md already exists")
	firstData, err := environment.Files.ReadFile(t.Context(), "issues/PROJ-1.md")
	require.NoError(t, err)
	require.Equal(t, "local edit", string(firstData))
	secondData, err := environment.Files.ReadFile(t.Context(), "issues/PROJ-2.md")
	require.NoError(t, err)
	require.Contains(t, string(secondData), "# Fix sync {#wd-field-title}")
}

func TestJiraRemotePullCommandUsesLimit(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	api := jiraapi.NewIssueSearchWithIssuesMemClient(newServerInfo(t), []jiraapi.JiraIssue{})
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
	api := jiraapi.NewIssueSearchWithIssuesMemClient(newServerInfo(t), []jiraapi.JiraIssue{})
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

func newJiraIssue() jiraapi.JiraIssue {
	return jiraapi.JiraIssue{
		ID:  "10001",
		Key: "PROJ-1",
		Fields: jiraapi.JiraIssueFields{
			Summary:   "Write docs",
			IssueType: jiraapi.JiraNamedValue{Name: "Task"},
			Project:   jiraapi.JiraProject{Key: "PROJ"},
			Status:    jiraapi.JiraNamedValue{Name: "To Do"},
			Reporter:  jiraapi.JiraUser{DisplayName: "Alice Author"},
			Assignee:  jiraapi.JiraUser{DisplayName: "Bob Owner"},
		},
	}
}

func newSecondJiraIssue() jiraapi.JiraIssue {
	return jiraapi.JiraIssue{
		ID:  "10002",
		Key: "PROJ-2",
		Fields: jiraapi.JiraIssueFields{
			Summary:   "Fix sync",
			IssueType: jiraapi.JiraNamedValue{Name: "Bug"},
			Project:   jiraapi.JiraProject{Key: "PROJ"},
			Status:    jiraapi.JiraNamedValue{Name: "In Progress"},
			Reporter:  jiraapi.JiraUser{DisplayName: "Alice Author"},
			Assignee:  jiraapi.JiraUser{DisplayName: "Bob Owner"},
		},
	}
}

func newServerInfo(t *testing.T) jiraapi.ServerInfo {
	t.Helper()
	baseURL, err := url.Parse("https://jira.example.test")
	require.NoError(t, err)
	serverInfo, err := jiraapi.NewServerInfo(baseURL, jiraapi.ParseJiraVersion("9.12.0"))
	require.NoError(t, err)
	return serverInfo
}
