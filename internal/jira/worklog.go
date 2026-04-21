package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Worklog represents a single worklog entry for a Jira issue.
type Worklog struct {
	ID      string `json:"id"`
	Comment string `json:"comment"`
	Started string `json:"started"`
	TimeSpent string `json:"timeSpent"`
	TimeSpentSeconds int `json:"timeSpentSeconds"`
	Author struct {
		DisplayName string `json:"displayName"`
		EmailAddress string `json:"emailAddress"`
	} `json:"author"`
}

// WorklogResponse holds the paginated list of worklogs returned by the API.
type WorklogResponse struct {
	StartAt    int       `json:"startAt"`
	MaxResults int       `json:"maxResults"`
	Total      int       `json:"total"`
	Worklogs   []Worklog `json:"worklogs"`
}

// GetWorklogs retrieves all worklog entries for the given issue key.
func (c *Client) GetWorklogs(issueKey string) (*WorklogResponse, error) {
	url := fmt.Sprintf("%s/rest/api/2/issue/%s/worklog", c.baseURL, issueKey)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("worklog request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("issue %q not found", issueKey)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for issue %q", resp.StatusCode, issueKey)
	}

	var result WorklogResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode worklog response: %w", err)
	}
	return &result, nil
}
