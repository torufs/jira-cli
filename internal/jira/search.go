package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// SearchRequest holds parameters for a JQL search.
type SearchRequest struct {
	JQL        string
	MaxResults int
	StartAt    int
	Fields     []string
}

// SearchResult represents the Jira search API response.
type SearchResult struct {
	Total      int     `json:"total"`
	StartAt    int     `json:"startAt"`
	MaxResults int     `json:"maxResults"`
	Issues     []Issue `json:"issues"`
}

// Issue is a minimal representation of a Jira issue.
type Issue struct {
	Key    string      `json:"key"`
	Fields IssueFields `json:"fields"`
}

// IssueFields holds common fields returned for an issue.
type IssueFields struct {
	Summary  string `json:"summary"`
	Status   Status `json:"status"`
	Assignee *User  `json:"assignee"`
	Priority *Priority `json:"priority"`
}

// Status represents the issue status.
type Status struct {
	Name string `json:"name"`
}

// User represents a Jira user.
type User struct {
	DisplayName string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
}

// Priority represents the issue priority.
type Priority struct {
	Name string `json:"name"`
}

// Search performs a JQL search against the Jira API.
func (c *Client) Search(req SearchRequest) (*SearchResult, error) {
	if req.MaxResults == 0 {
		req.MaxResults = 50
	}

	params := url.Values{}
	params.Set("jql", req.JQL)
	params.Set("maxResults", fmt.Sprintf("%d", req.MaxResults))
	params.Set("startAt", fmt.Sprintf("%d", req.StartAt))

	apiURL := fmt.Sprintf("%s/rest/api/2/search?%s", c.baseURL, params.Encode())

	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	return &result, nil
}
