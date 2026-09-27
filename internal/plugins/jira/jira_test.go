package jira

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/richargh/workdown/internal/kernel"
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

func TestJiraPullEmitsInterchangeWorkItems(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	api := jiraapi.NewIssueSearchWithIssuesMemClient(newServerInfo(t), []jiraapi.JiraIssue{
		newJiraIssue("10001", "PROJ-1", "PROJ", "Task", "Write docs"),
		newJiraIssue("10002", "PROJ-2", "PROJ", "Bug", "Fix sync"),
	})
	plugin := NewWithAPI(environment, api)

	// when
	result, err := plugin.Pull(context.Background(), PullRequest{URL: "https://jira.example.test", Project: "PROJ", PAT: "secret-token"})

	// then
	require.NoError(t, err)
	require.Equal(t, []string{"project = PROJ"}, jiraapi.SearchIssuesQueries(api))
	require.Equal(t, [][]string{{"summary", "issuetype", "project", "status", "reporter", "assignee"}}, jiraapi.SearchIssuesFields(api))
	require.Equal(t, []int{5}, jiraapi.SearchIssuesLimits(api))
	require.Len(t, result.WorkItems, 2)
	require.Equal(t, "jira", result.WorkItems[0].Provider)
	require.Equal(t, "10001", result.WorkItems[0].ID)
	require.Equal(t, "PROJ-1", result.WorkItems[0].Key)
	require.Equal(t, "Write docs", result.WorkItems[0].Title)
	require.Equal(t, []kernel.Field{
		{Name: "project", ProviderKey: "project", Type: "string", Value: "PROJ", Editable: false},
		{Name: "issueType", ProviderKey: "issuetype", Type: "string", Value: "Task", Editable: false},
		{Name: "status", ProviderKey: "status", Type: "string", Value: "To Do", Editable: false},
		{Name: "author", ProviderKey: "reporter", Type: "string", Value: "Alice Author", Editable: false},
		{Name: "owner", ProviderKey: "assignee", Type: "string", Value: "Bob Owner", Editable: false},
	}, result.WorkItems[0].Fields)
	require.Equal(t, "https://jira.example.test/browse/PROJ-1", result.WorkItems[0].URL)
}

func TestJiraPullCanLimitToMyIssues(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	api := jiraapi.NewIssueSearchWithIssuesMemClient(newServerInfo(t), []jiraapi.JiraIssue{
		newJiraIssue("10001", "PROJ-1", "PROJ", "Task", "Write docs"),
	})
	plugin := NewWithAPI(environment, api)

	// when
	result, err := plugin.Pull(context.Background(), PullRequest{URL: "https://jira.example.test", Project: "PROJ", PAT: "secret-token", Mine: true})

	// then
	require.NoError(t, err)
	require.Equal(t, []string{"project = PROJ AND assignee = currentUser()"}, jiraapi.SearchIssuesQueries(api))
	require.Len(t, result.WorkItems, 1)
	require.Equal(t, "PROJ-1", result.WorkItems[0].Key)
}

func TestJiraPullUsesRequestedLimit(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	api := jiraapi.NewIssueSearchWithIssuesMemClient(newServerInfo(t), []jiraapi.JiraIssue{})
	plugin := NewWithAPI(environment, api)

	// when
	_, err := plugin.Pull(context.Background(), PullRequest{URL: "https://jira.example.test", Project: "PROJ", PAT: "secret-token", Limit: 2})

	// then
	require.NoError(t, err)
	require.Equal(t, []int{2}, jiraapi.SearchIssuesLimits(api))
}

func TestJiraPullRejectsNegativeLimit(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	api := jiraapi.NewIssueSearchWithIssuesMemClient(newServerInfo(t), []jiraapi.JiraIssue{})
	plugin := NewWithAPI(environment, api)

	// when
	_, err := plugin.Pull(context.Background(), PullRequest{URL: "https://jira.example.test", Project: "PROJ", PAT: "secret-token", Limit: -1})

	// then
	require.ErrorContains(t, err, "--limit must be positive")
}

func TestJiraPullCanLimitToSelectedIssues(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	api := jiraapi.NewIssueSearchWithIssuesMemClient(newServerInfo(t), []jiraapi.JiraIssue{})
	plugin := NewWithAPI(environment, api)

	// when
	_, err := plugin.Pull(context.Background(), PullRequest{URL: "https://jira.example.test", Project: "PROJ", IssueKeys: []string{"PROJ-1", "PROJ-2"}, PAT: "secret-token"})

	// then
	require.NoError(t, err)
	require.Equal(t, []string{"project = PROJ AND key in (PROJ-1, PROJ-2)"}, jiraapi.SearchIssuesQueries(api))
	require.Equal(t, []int{2}, jiraapi.SearchIssuesLimits(api))
}

func TestJiraPullRejectsLimitSmallerThanSelectedIssues(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	api := jiraapi.NewIssueSearchWithIssuesMemClient(newServerInfo(t), []jiraapi.JiraIssue{})
	plugin := NewWithAPI(environment, api)

	// when
	_, err := plugin.Pull(context.Background(), PullRequest{URL: "https://jira.example.test", Project: "PROJ", IssueKeys: []string{"PROJ-1", "PROJ-2"}, PAT: "secret-token", Limit: 1})

	// then
	require.ErrorContains(t, err, "--limit must be at least the number of selected issues")
}

func TestJiraPullRequiresProject(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	api := jiraapi.NewIssueSearchWithIssuesMemClient(newServerInfo(t), []jiraapi.JiraIssue{})
	plugin := NewWithAPI(environment, api)

	// when
	_, err := plugin.Pull(context.Background(), PullRequest{URL: "https://jira.example.test", IssueKeys: []string{"PROJ-1"}, PAT: "secret-token"})

	// then
	require.ErrorContains(t, err, "--project is required")
}

func newJiraIssue(id string, key string, project string, issueType string, summary string) jiraapi.JiraIssue {
	return jiraapi.JiraIssue{
		ID:  id,
		Key: key,
		Fields: jiraapi.JiraIssueFields{
			Summary:   summary,
			IssueType: jiraapi.JiraNamedValue{Name: issueType},
			Project:   jiraapi.JiraProject{Key: project},
			Status:    jiraapi.JiraNamedValue{Name: "To Do"},
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
