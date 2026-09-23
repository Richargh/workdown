package jiraapi

import (
	"context"
	"fmt"
)

type memClient struct {
	serverInfoResponse ServerInfo
	serverInfoError    error
}

func NewMemClient(info ServerInfo) JiraAPI {
	return &memClient{serverInfoResponse: info}
}

func NewFailingMemClient(err error) JiraAPI {
	return &memClient{serverInfoError: err}
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
