package jiraapi_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/richargh/workdown/internal/plugins/jira/jiraapi"
)

func TestMemClientContract(t *testing.T) {
	testJiraAPIContract(t, func(t *testing.T) jiraapi.JiraAPI {
		t.Helper()
		return jiraapi.NewIssueSearchWithIssuesMemClient(newServerInfo(t, "https://jira.example.test", "9.12.0"), []jiraapi.JiraIssue{
			newJiraIssue("10001", "PROJ-1", "PROJ", "Task", "Write docs"),
		})
	})
}

func TestRestClientContract(t *testing.T) {
	testJiraAPIContract(t, func(t *testing.T) jiraapi.JiraAPI {
		t.Helper()
		server := configureHttpServerRoutes(t)
		t.Cleanup(server.Close)

		baseURL, err := url.Parse(server.URL)
		require.NoError(t, err)

		testee, err := jiraapi.NewRestClient(baseURL, server.Client())
		require.NoError(t, err)
		return testee
	})
}

func configureHttpServerRoutes(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		switch r.URL.Path {
		case "/rest/api/2/serverInfo":
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{"baseUrl":"https://jira.example.test","version":"9.12.0"}`))
			require.NoError(t, err)
		case "/rest/api/2/myself":
			w.WriteHeader(http.StatusOK)
		case "/rest/api/2/search":
			require.Equal(t, "project = PROJ", r.URL.Query().Get("jql"))
			require.Equal(t, "summary,issuetype,project", r.URL.Query().Get("fields"))
			require.Equal(t, "1", r.URL.Query().Get("maxResults"))
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{"total":1,"issues":[{"id":"10001","key":"PROJ-1","fields":{"summary":"Write docs","issuetype":{"name":"Task"},"project":{"key":"PROJ"}}}]}`))
			require.NoError(t, err)
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
}

func testJiraAPIContract(t *testing.T, newTestee func(t *testing.T) jiraapi.JiraAPI) {
	t.Helper()

	t.Run("server info returns base URL and version", func(t *testing.T) {
		// given
		testee := newTestee(t)

		// when
		result, err := testee.ServerInfo(context.Background())

		// then
		require.NoError(t, err)
		baseURL := result.BaseURL()
		require.Equal(t, "https://jira.example.test", baseURL.String())
		require.Equal(t, "9.12.0", result.Version().String())
	})

	t.Run("validate connection succeeds", func(t *testing.T) {
		// given
		testee := newTestee(t)

		// when
		err := testee.ValidateConnection(context.Background())

		// then
		require.NoError(t, err)
	})

	t.Run("search issues returns total count", func(t *testing.T) {
		// given
		testee := newTestee(t)

		// when
		result, err := testee.SearchIssues(context.Background(), jiraapi.JiraIssueSearchRequest{
			JQL:        "project = PROJ",
			Fields:     []string{"summary", "issuetype", "project"},
			MaxResults: 1,
		})

		// then
		require.NoError(t, err)
		require.Equal(t, 1, result.Total)
		require.Equal(t, []jiraapi.JiraIssue{newJiraIssue("10001", "PROJ-1", "PROJ", "Task", "Write docs")}, result.Issues)
	})
}

func TestMemClientReturnsConfiguredError(t *testing.T) {
	// given
	want := errors.New("boom")
	testee := jiraapi.NewFailingMemClient(want)

	// when
	serverInfo, err := testee.ServerInfo(context.Background())

	// then
	require.ErrorIs(t, err, want)
	require.Equal(t, jiraapi.ServerInfo{}, serverInfo)
}

func TestRestClientSendsPATAsBearerTokenWhenValidatingConnection(t *testing.T) {
	// given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/rest/api/2/myself", r.URL.Path)
		require.Equal(t, "Bearer secret-token", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	baseURL, err := url.Parse(server.URL)
	require.NoError(t, err)
	testee, err := jiraapi.NewAuthenticatedRestClient(baseURL, server.Client(), "secret-token")
	require.NoError(t, err)

	// when
	err = testee.ValidateConnection(context.Background())

	// then
	require.NoError(t, err)
}

func TestRestClientReturnsErrorForValidationHTTPFailure(t *testing.T) {
	// given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusUnauthorized)
	}))
	defer server.Close()
	baseURL, err := url.Parse(server.URL)
	require.NoError(t, err)
	testee, err := jiraapi.NewRestClient(baseURL, server.Client())
	require.NoError(t, err)

	// when
	err = testee.ValidateConnection(context.Background())

	// then
	require.Error(t, err)
}

func TestNewRestClientRequiresBaseURL(t *testing.T) {
	// when
	api, err := jiraapi.NewRestClient(nil, http.DefaultClient)

	// then
	require.Error(t, err)
	require.Nil(t, api)
}

func TestNewRestClientRequiresAbsoluteBaseURL(t *testing.T) {
	// given
	baseURL, err := url.Parse("/relative")
	require.NoError(t, err)

	// when
	testee, err := jiraapi.NewRestClient(baseURL, http.DefaultClient)

	// then
	require.Error(t, err)
	require.Nil(t, testee)
}

func TestNewRestClientRequiresHTTPClient(t *testing.T) {
	// given
	baseURL, err := url.Parse("https://jira.example.test")
	require.NoError(t, err)

	// when
	testee, err := jiraapi.NewRestClient(baseURL, nil)

	// then
	require.Error(t, err)
	require.Nil(t, testee)
}

func TestRestClientReturnsErrorForHTTPFailure(t *testing.T) {
	// given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusUnauthorized)
	}))
	defer server.Close()
	baseURL, err := url.Parse(server.URL)
	require.NoError(t, err)
	testee, err := jiraapi.NewRestClient(baseURL, server.Client())
	require.NoError(t, err)

	// when
	serverInfo, err := testee.ServerInfo(context.Background())

	// then
	require.Error(t, err)
	require.Equal(t, jiraapi.ServerInfo{}, serverInfo)
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
