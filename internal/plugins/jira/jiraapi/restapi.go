package jiraapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type restClient struct {
	baseURL    url.URL
	httpClient *http.Client
	pat        string
}

type jiraServerInfoDTO struct {
	BaseURL string `json:"baseUrl"`
	Version string `json:"version"`
}

func NewRestClient(baseURL *url.URL, httpClient *http.Client) (JiraAPI, error) {
	return NewAuthenticatedRestClient(baseURL, httpClient, "")
}

func NewAuthenticatedRestClient(baseURL *url.URL, httpClient *http.Client, pat string) (JiraAPI, error) {
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
		pat:        pat,
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

	c.authorize(request)

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

	var payload jiraServerInfoDTO
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return ServerInfo{}, err
	}

	baseURL, err := url.Parse(payload.BaseURL)
	if err != nil {
		return ServerInfo{}, err
	}
	return NewServerInfo(baseURL, ParseJiraVersion(payload.Version))
}

func (c *restClient) ValidateConnection(ctx context.Context) error {
	myselfURL := c.baseURL.JoinPath("rest", "api", "2", "myself")
	myselfURL.RawQuery = ""
	myselfURL.Fragment = ""

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, myselfURL.String(), nil)
	if err != nil {
		return err
	}
	c.authorize(request)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return fmt.Errorf("jira authentication failed: %s", response.Status)
	}
	return nil
}

func (c *restClient) SearchIssues(ctx context.Context, search JiraIssueSearchRequest) (JiraIssueSearchResult, error) {
	searchURL := c.baseURL.JoinPath("rest", "api", "2", "search")
	query := searchURL.Query()
	query.Set("jql", search.JQL)
	query.Set("maxResults", fmt.Sprintf("%d", search.MaxResults))
	if len(search.Fields) > 0 {
		query.Set("fields", strings.Join(search.Fields, ","))
	}
	searchURL.RawQuery = query.Encode()
	searchURL.Fragment = ""

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL.String(), nil)
	if err != nil {
		return JiraIssueSearchResult{}, err
	}
	c.authorize(request)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return JiraIssueSearchResult{}, err
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return JiraIssueSearchResult{}, fmt.Errorf("jira issue search failed: %s", response.Status)
	}

	var payload JiraIssueSearchResult
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return JiraIssueSearchResult{}, err
	}
	return payload, nil
}

func (c *restClient) authorize(request *http.Request) {
	if c.pat != "" {
		request.Header.Set("Authorization", "Bearer "+c.pat)
	}
}
