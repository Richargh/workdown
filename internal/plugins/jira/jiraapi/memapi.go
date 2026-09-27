package jiraapi

import (
	"context"
	"fmt"
)

type memClient struct {
	serverInfoResponse      ServerInfo
	serverInfoError         error
	validateConnectionError error
	searchIssuesResponse    IssueSearchResult
	searchIssuesError       error
	searchIssuesQueries     []string
}

func NewMemClient(info ServerInfo) JiraAPI {
	return &memClient{serverInfoResponse: info}
}

func NewFailingMemClient(err error) JiraAPI {
	return &memClient{serverInfoError: err, validateConnectionError: err}
}

func NewFailingValidationMemClient(info ServerInfo, err error) JiraAPI {
	return &memClient{serverInfoResponse: info, validateConnectionError: err}
}

func NewIssueSearchMemClient(info ServerInfo, total int) JiraAPI {
	return &memClient{serverInfoResponse: info, searchIssuesResponse: IssueSearchResult{Total: total}}
}

func (c *memClient) ServerInfo(ctx context.Context) (ServerInfo, error) {
	if c.serverInfoError != nil {
		return ServerInfo{}, c.serverInfoError
	}
	if c.serverInfoResponse.baseURL.Host == "" {
		return ServerInfo{}, fmt.Errorf("server info response missing base URL")
	}
	return c.serverInfoResponse, nil
}

func (c *memClient) ValidateConnection(ctx context.Context) error {
	return c.validateConnectionError
}

func (c *memClient) SearchIssues(ctx context.Context, jql string) (IssueSearchResult, error) {
	c.searchIssuesQueries = append(c.searchIssuesQueries, jql)
	if c.searchIssuesError != nil {
		return IssueSearchResult{}, c.searchIssuesError
	}
	return c.searchIssuesResponse, nil
}

func SearchIssuesQueries(api JiraAPI) []string {
	client, ok := api.(*memClient)
	if !ok {
		return nil
	}
	return append([]string(nil), client.searchIssuesQueries...)
}
