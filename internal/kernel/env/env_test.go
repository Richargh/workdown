package env_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/richargh/workdown/internal/kernel/env"
)

func TestNewMemEnvUsesInMemoryAdapters(t *testing.T) {
	// given
	environment := env.NewMemEnvWithCredentials(map[string]string{"jira": "token"})

	// when
	err := environment.Files.WriteFile(context.Background(), "issue.md", []byte("hello"))
	require.NoError(t, err)
	fileContent, err := environment.Files.ReadFile(context.Background(), "issue.md")
	require.NoError(t, err)

	err = environment.Config.Save(context.Background(), "config.toml", []byte("config"))
	require.NoError(t, err)
	configContent, err := environment.Config.Load(context.Background(), "config.toml")
	require.NoError(t, err)

	credential, err := environment.Credentials.Credential(context.Background(), "jira")

	// then
	require.NoError(t, err)
	require.Equal(t, []byte("hello"), fileContent)
	require.Equal(t, []byte("config"), configContent)
	require.Equal(t, "token", credential)
	require.Equal(t, "?.?.?", environment.Version.String())
	require.NotNil(t, environment.HTTPClient)
}
