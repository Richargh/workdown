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
		return jiraapi.NewMemClient(newServerInfo(t, "https://jira.example.test", "9.12.0"))
	})
}

func TestRestClientContract(t *testing.T) {
	testJiraAPIContract(t, func(t *testing.T) jiraapi.JiraAPI {
		t.Helper()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/rest/api/2/serverInfo", r.URL.Path)

			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{"baseUrl":"https://jira.example.test","version":"9.12.0"}`))
			require.NoError(t, err)
		}))
		t.Cleanup(server.Close)

		baseURL, err := url.Parse(server.URL)
		require.NoError(t, err)

		api, err := jiraapi.NewRestClient(baseURL, server.Client())
		require.NoError(t, err)
		return api
	})
}

func testJiraAPIContract(t *testing.T, newAPI func(t *testing.T) jiraapi.JiraAPI) {
	t.Helper()

	t.Run("server info returns base URL and version", func(t *testing.T) {
		// given
		api := newAPI(t)

		// when
		serverInfo, err := api.ServerInfo(context.Background())

		// then
		require.NoError(t, err)
		baseURL := serverInfo.BaseURL()
		require.Equal(t, "https://jira.example.test", baseURL.String())
		require.Equal(t, "9.12.0", serverInfo.Version().String())
	})
}

func TestMemClientReturnsConfiguredError(t *testing.T) {
	// given
	want := errors.New("boom")
	api := jiraapi.NewFailingMemClient(want)

	// when
	serverInfo, err := api.ServerInfo(context.Background())

	// then
	require.ErrorIs(t, err, want)
	require.Equal(t, jiraapi.ServerInfo{}, serverInfo)
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
	api, err := jiraapi.NewRestClient(baseURL, http.DefaultClient)

	// then
	require.Error(t, err)
	require.Nil(t, api)
}

func TestNewRestClientRequiresHTTPClient(t *testing.T) {
	// given
	baseURL, err := url.Parse("https://jira.example.test")
	require.NoError(t, err)

	// when
	api, err := jiraapi.NewRestClient(baseURL, nil)

	// then
	require.Error(t, err)
	require.Nil(t, api)
}

func TestRestClientReturnsErrorForHTTPFailure(t *testing.T) {
	// given
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusUnauthorized)
	}))
	defer server.Close()
	baseURL, err := url.Parse(server.URL)
	require.NoError(t, err)
	api, err := jiraapi.NewRestClient(baseURL, server.Client())
	require.NoError(t, err)

	// when
	serverInfo, err := api.ServerInfo(context.Background())

	// then
	require.Error(t, err)
	require.Equal(t, jiraapi.ServerInfo{}, serverInfo)
}

func newServerInfo(t *testing.T, rawBaseURL string, rawVersion string) jiraapi.ServerInfo {
	t.Helper()
	baseURL, err := url.Parse(rawBaseURL)
	require.NoError(t, err)
	serverInfo, err := jiraapi.NewServerInfo(baseURL, jiraapi.ParseJiraVersion(rawVersion))
	require.NoError(t, err)
	return serverInfo
}
