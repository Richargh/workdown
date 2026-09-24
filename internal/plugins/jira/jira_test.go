package jira

import (
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/richargh/workdown/internal/kernel/env"
	"github.com/richargh/workdown/internal/plugins/jira/jiraapi"
)

func TestJiraCheckCommandValidatesConnectionUsingURLFlag(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	api := jiraapi.NewMemClient(newServerInfo(t, "https://jira.example.test", "9.12.0"))
	plugin := NewWithAPI(environment, api)
	cmd := plugin.Commands()[0]

	// when
	stdout, err := executeCommand(t, cmd, "check", "--url", "https://jira.example.test", "--pat", "secret-token")

	// then
	require.NoError(t, err)
	require.Equal(t, "jira https://jira.example.test ok (version 9.12.0)\n", stdout)
}

func TestJiraCheckCommandReturnsValidationError(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	want := errors.New("bad token")
	api := jiraapi.NewFailingValidationMemClient(newServerInfo(t, "https://jira.example.test", "9.12.0"), want)
	cmd := NewWithAPI(environment, api).Commands()[0]

	// when
	_, err := executeCommand(t, cmd, "check", "--url", "https://jira.example.test", "--pat", "secret-token")

	// then
	require.ErrorIs(t, err, want)
}

func TestJiraCheckCommandRequiresPAT(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	cmd := NewWithAPI(environment, jiraapi.NewMemClient(newServerInfo(t, "https://jira.example.test", "9.12.0"))).Commands()[0]

	// when
	_, err := executeCommand(t, cmd, "check")

	// then
	require.ErrorContains(t, err, "--pat is required")
}

func TestJiraCheckCommandRequiresURL(t *testing.T) {
	// given
	environment := env.NewMemEnv()
	cmd := NewWithAPI(environment, jiraapi.NewMemClient(newServerInfo(t, "https://jira.example.test", "9.12.0"))).Commands()[0]

	// when
	_, err := executeCommand(t, cmd, "check", "--pat", "secret-token")

	// then
	require.ErrorContains(t, err, "--url is required")
}

func newServerInfo(t *testing.T, rawBaseURL string, rawVersion string) jiraapi.ServerInfo {
	t.Helper()
	baseURL, err := url.Parse(rawBaseURL)
	require.NoError(t, err)
	serverInfo, err := jiraapi.NewServerInfo(baseURL, jiraapi.ParseJiraVersion(rawVersion))
	require.NoError(t, err)
	return serverInfo
}
