package jiraapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type restClient struct {
	baseURL    url.URL
	httpClient *http.Client
}

func NewRestClient(baseURL *url.URL, httpClient *http.Client) (JiraAPI, error) {
	if baseURL == nil {
		return nil, fmt.Errorf("jira base URL is required")
	}
	if baseURL.Scheme == "" || baseURL.Host == "" {
		return nil, fmt.Errorf("jira base URL must be absolute")
	}
	if httpClient == nil {
		return nil, fmt.Errorf("http client is required")
	}

	return &restClient{
		baseURL:    *baseURL,
		httpClient: httpClient,
	}, nil
}

func (c *restClient) ServerInfo(ctx context.Context) (ServerInfo, error) {
	serverInfoURL := c.baseURL.JoinPath("rest", "api", "2", "serverInfo")
	serverInfoURL.RawQuery = ""
	serverInfoURL.Fragment = ""

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, serverInfoURL.String(), nil)
	if err != nil {
		return ServerInfo{}, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return ServerInfo{}, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return ServerInfo{}, fmt.Errorf("jira server info failed: %s", response.Status)
	}

	var payload struct {
		BaseURL string `json:"baseUrl"`
		Version string `json:"version"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return ServerInfo{}, err
	}

	baseURL, err := url.Parse(payload.BaseURL)
	if err != nil {
		return ServerInfo{}, err
	}
	return NewServerInfo(baseURL, ParseJiraVersion(payload.Version))
}
