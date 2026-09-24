package jiraapi_test

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/richargh/workdown/internal/plugins/jira/jiraapi"
)

func TestRealJiraValidateConnectionOptIn(t *testing.T) {
	rawURL := os.Getenv("WORKDOWN_JIRA_URL")
	pat := os.Getenv("WORKDOWN_JIRA_PAT")
	if rawURL == "" || pat == "" {
		t.Skip("set WORKDOWN_JIRA_URL and WORKDOWN_JIRA_PAT to run real Jira integration test")
	}
	baseURL, err := url.Parse(rawURL)
	require.NoError(t, err)
	api, err := jiraapi.NewAuthenticatedRestClient(baseURL, http.DefaultClient, pat)
	require.NoError(t, err)

	err = api.ValidateConnection(context.Background())
	require.NoError(t, err)
}
