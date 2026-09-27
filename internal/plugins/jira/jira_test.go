package jira

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/richargh/workdown/internal/kernel/env"
	"github.com/richargh/workdown/internal/plugins/jira/jiraapi"
)

func TestJiraCheckValidatesConnectionUsingURL(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	api := jiraapi.NewMemClient(newServerInfo(t))
	plugin := NewWithAPI(environment, api)

	// when
	result, err := plugin.Check(context.Background(), CheckRequest{URL: "https://jira.example.test", PAT: "secret-token"})

	// then
	require.NoError(t, err)
	require.Equal(t, "https://jira.example.test", result.BaseURL)
	require.Equal(t, "9.12.0", result.Version)
	require.False(t, result.HasProject())
}

func TestJiraCheckCanValidateProjectAccess(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	api := jiraapi.NewIssueSearchMemClient(newServerInfo(t), 5)
	plugin := NewWithAPI(environment, api)

	// when
	result, err := plugin.Check(context.Background(), CheckRequest{URL: "https://jira.example.test", Project: "PROJ", PAT: "secret-token"})

	// then
	require.NoError(t, err)
	require.True(t, result.HasProject())
	require.Equal(t, "PROJ", result.Project)
	require.Equal(t, 5, result.ProjectIssueCount)
	require.Equal(t, []string{"project = PROJ"}, jiraapi.SearchIssuesQueries(api))
}

func TestJiraCheckReturnsValidationError(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	want := errors.New("bad token")
	api := jiraapi.NewFailingValidationMemClient(newServerInfo(t), want)
	plugin := NewWithAPI(environment, api)

	// when
	_, err := plugin.Check(context.Background(), CheckRequest{URL: "https://jira.example.test", PAT: "secret-token"})

	// then
	require.ErrorIs(t, err, want)
}

func TestJiraCheckRequiresPAT(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	plugin := NewWithAPI(environment, jiraapi.NewMemClient(newServerInfo(t)))

	// when
	_, err := plugin.Check(context.Background(), CheckRequest{})

	// then
	require.ErrorContains(t, err, "--pat is required")
}

func TestJiraCheckRequiresURL(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	plugin := NewWithAPI(environment, jiraapi.NewMemClient(newServerInfo(t)))

	// when
	_, err := plugin.Check(context.Background(), CheckRequest{PAT: "secret-token"})

	// then
	require.ErrorContains(t, err, "--url is required")
}

func TestJiraPullCountsSelectedProjectIssues(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	api := jiraapi.NewIssueSearchMemClient(newServerInfo(t), 3)
	plugin := NewWithAPI(environment, api)

	// when
	result, err := plugin.Pull(context.Background(), PullRequest{URL: "https://jira.example.test", Project: "PROJ", IssueKeys: []string{"PROJ-1", "PROJ-2"}, PAT: "secret-token"})

	// then
	require.NoError(t, err)
	require.Equal(t, 3, result.IssueCount)
	require.Equal(t, []string{"project = PROJ AND key in (PROJ-1, PROJ-2)"}, jiraapi.SearchIssuesQueries(api))
}

func TestJiraPullRequiresIssues(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	api := jiraapi.NewIssueSearchMemClient(newServerInfo(t), 3)
	plugin := NewWithAPI(environment, api)

	// when
	_, err := plugin.Pull(context.Background(), PullRequest{URL: "https://jira.example.test", Project: "PROJ", PAT: "secret-token"})

	// then
	require.ErrorContains(t, err, "--issues is required")
}

func TestJiraPullRequiresProject(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	api := jiraapi.NewIssueSearchMemClient(newServerInfo(t), 3)
	plugin := NewWithAPI(environment, api)

	// when
	_, err := plugin.Pull(context.Background(), PullRequest{URL: "https://jira.example.test", IssueKeys: []string{"PROJ-1"}, PAT: "secret-token"})

	// then
	require.ErrorContains(t, err, "--project is required")
}

func newServerInfo(t *testing.T) jiraapi.ServerInfo {
	t.Helper()
	baseURL, err := url.Parse("https://jira.example.test")
	require.NoError(t, err)
	serverInfo, err := jiraapi.NewServerInfo(baseURL, jiraapi.ParseJiraVersion("9.12.0"))
	require.NoError(t, err)
	return serverInfo
}
