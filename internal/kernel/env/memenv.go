package env

import (
	"fmt"
	"net/http"

	"github.com/richargh/workdown/internal/kernel"
	"github.com/richargh/workdown/internal/kernel/config"
	"github.com/richargh/workdown/internal/kernel/credentials"
	"github.com/richargh/workdown/internal/kernel/fs"
)

func NewMemEnv() Env {
	return NewMemEnvWithCredentials(nil)
}

func NewMemEnvWithCredentials(values map[string]string) Env {
	return Env{
		Version:     kernel.UnknownWorkdownVersion,
		Config:      config.NewMemStore(),
		Files:       fs.NewMemStore(),
		Credentials: credentials.NewMemSource(values),
		HTTPClient:  &http.Client{Transport: noNetworkTransport{}},
	}
}

type noNetworkTransport struct{}

func (t noNetworkTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("unexpected HTTP request in memory environment: %s %s", request.Method, request.URL)
}
