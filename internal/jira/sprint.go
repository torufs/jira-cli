package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Sprint represents a Jira agile sprint.
type Sprint struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	State         string `json:"state"`
	BoardID       int    `json:"originBoardId"`
	StartDate     string `json:"startDate,omitempty"`
	EndDate       string `json:"endDate,omitempty"`
	CompleteDate  string `json:"completeDate,omitempty"`
}

// SprintList holds a paginated list of sprints returned by the Jira API.
type SprintList struct {
	Values     []Sprint `json:"values"`
	Total      int      `json:"total"`
	MaxResults int      `json:"maxResults"`
	StartAt    int      `json:"startAt"`
	IsLast     bool     `json:"isLast"`
}

// GetSprints retrieves all sprints for the given board ID.
func (c *Client) GetSprints(boardID int) (*SprintList, error) {
	url := fmt.Sprintf("%s/rest/agile/1.0/board/%d/sprint", c.baseURL, boardID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get sprints request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d when fetching sprints for board %d", resp.StatusCode, boardID)
	}

	var list SprintList
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, fmt.Errorf("failed to decode sprint list: %w", err)
	}

	return &list, nil
}

// GetActiveSprint returns the currently active sprint for the given board ID,
// or nil if no active sprint exists.
func (c *Client) GetActiveSprint(boardID int) (*Sprint, error) {
	url := fmt.Sprintf("%s/rest/agile/1.0/board/%d/sprint?state=active", c.baseURL, boardID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get active sprint request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d when fetching active sprint for board %d", resp.StatusCode, boardID)
	}

	var list SprintList
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, fmt.Errorf("failed to decode active sprint response: %w", err)
	}

	if len(list.Values) == 0 {
		return nil, nil
	}

	return &list.Values[0], nil
}
