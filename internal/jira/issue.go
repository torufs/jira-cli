package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Issue represents a Jira issue.
type Issue struct {
	Key    string      `json:"key"`
	Fields IssueFields `json:"fields"`
}

// IssueFields holds the fields of a Jira issue.
type IssueFields struct {
	Summary     string   `json:"summary"`
	Description string   `json:"description"`
	Status      Status   `json:"status"`
	Assignee    Assignee `json:"assignee"`
	Priority    Priority `json:"priority"`
}

// Status represents the status of a Jira issue.
type Status struct {
	Name string `json:"name"`
}

// Assignee represents the assignee of a Jira issue.
type Assignee struct {
	DisplayName string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
}

// Priority represents the priority of a Jira issue.
type Priority struct {
	Name string `json:"name"`
}

// Client is a Jira API client.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new Jira API client.
func NewClient(baseURL string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: httpClient,
	}
}

// GetIssue fetches a Jira issue by key.
func (c *Client) GetIssue(key string) (*Issue, error) {
	url := fmt.Sprintf("%s/rest/api/2/issue/%s", c.BaseURL, key)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching issue %s: %w", key, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for issue %s", resp.StatusCode, key)
	}

	var issue Issue
	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return nil, fmt.Errorf("decoding issue %s: %w", key, err)
	}

	return &issue, nil
}
