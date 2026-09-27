package jiraapi

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

type JiraAPI interface {
	ServerInfo(ctx context.Context) (ServerInfo, error)
	ValidateConnection(ctx context.Context) error
	SearchIssues(ctx context.Context, request JiraIssueSearchRequest) (JiraIssueSearchResult, error)
}

type JiraIssueSearchRequest struct {
	JQL        string
	Fields     []string
	MaxResults int
}

type JiraIssueSearchResult struct {
	Total  int         `json:"total"`
	Issues []JiraIssue `json:"issues"`
}

type JiraIssue struct {
	ID     string          `json:"id"`
	Key    string          `json:"key"`
	Fields JiraIssueFields `json:"fields"`
}

type JiraIssueFields struct {
	Summary   string         `json:"summary"`
	IssueType JiraNamedValue `json:"issuetype"`
	Project   JiraProject    `json:"project"`
}

type JiraNamedValue struct {
	Name string `json:"name"`
}

type JiraProject struct {
	Key string `json:"key"`
}

type JiraVersion struct {
	value string
}

func ParseJiraVersion(value string) JiraVersion {
	version := strings.TrimRight(value, "\r\n")
	if version == "" {
		return JiraVersion{value: "?"}
	}
	return JiraVersion{value: version}
}

func (v JiraVersion) String() string {
	if v.value == "" {
		return "?"
	}
	return v.value
}

type ServerInfo struct {
	baseURL url.URL
	version JiraVersion
}

func NewServerInfo(baseURL *url.URL, version JiraVersion) (ServerInfo, error) {
	if baseURL == nil {
		return ServerInfo{}, fmt.Errorf("jira server info base URL is required")
	}
	if baseURL.Scheme == "" || baseURL.Host == "" {
		return ServerInfo{}, fmt.Errorf("jira server info base URL must be absolute")
	}
	return ServerInfo{
		baseURL: *baseURL,
		version: version,
	}, nil
}

func (s ServerInfo) BaseURL() url.URL {
	return s.baseURL
}

func (s ServerInfo) Version() JiraVersion {
	return s.version
}
