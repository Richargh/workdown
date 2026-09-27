package jiraapi

import (
	"context"
	"fmt"
)

type memClient struct {
	serverInfoResponse      ServerInfo
	serverInfoError         error
	validateConnectionError error
	searchIssuesResponse    JiraIssueSearchResult
	searchIssuesError       error
	searchIssuesRequests    []JiraIssueSearchRequest
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
	return &memClient{serverInfoResponse: info, searchIssuesResponse: JiraIssueSearchResult{Total: total}}
}

func NewIssueSearchWithIssuesMemClient(info ServerInfo, issues []JiraIssue) JiraAPI {
	return &memClient{
		serverInfoResponse: info,
		searchIssuesResponse: JiraIssueSearchResult{
			Total:  len(issues),
			Issues: append([]JiraIssue(nil), issues...),
		},
	}
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

func (c *memClient) SearchIssues(ctx context.Context, request JiraIssueSearchRequest) (JiraIssueSearchResult, error) {
	c.searchIssuesRequests = append(c.searchIssuesRequests, cloneJiraIssueSearchRequest(request))
	if c.searchIssuesError != nil {
		return JiraIssueSearchResult{}, c.searchIssuesError
	}
	return JiraIssueSearchResult{
		Total:  c.searchIssuesResponse.Total,
		Issues: append([]JiraIssue(nil), c.searchIssuesResponse.Issues...),
	}, nil
}

func SearchIssuesQueries(api JiraAPI) []string {
	client, ok := api.(*memClient)
	if !ok {
		return nil
	}
	queries := make([]string, 0, len(client.searchIssuesRequests))
	for _, request := range client.searchIssuesRequests {
		queries = append(queries, request.JQL)
	}
	return queries
}

func SearchIssuesLimits(api JiraAPI) []int {
	client, ok := api.(*memClient)
	if !ok {
		return nil
	}
	limits := make([]int, 0, len(client.searchIssuesRequests))
	for _, request := range client.searchIssuesRequests {
		limits = append(limits, request.MaxResults)
	}
	return limits
}

func SearchIssuesFields(api JiraAPI) [][]string {
	client, ok := api.(*memClient)
	if !ok {
		return nil
	}
	fields := make([][]string, 0, len(client.searchIssuesRequests))
	for _, request := range client.searchIssuesRequests {
		fields = append(fields, append([]string(nil), request.Fields...))
	}
	return fields
}

func cloneJiraIssueSearchRequest(request JiraIssueSearchRequest) JiraIssueSearchRequest {
	return JiraIssueSearchRequest{
		JQL:        request.JQL,
		Fields:     append([]string(nil), request.Fields...),
		MaxResults: request.MaxResults,
	}
}
