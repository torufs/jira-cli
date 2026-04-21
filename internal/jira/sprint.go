package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Sprint represents a Jira sprint.
type Sprint struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	State         string `json:"state"`
	StartDate     string `json:"startDate,omitempty"`
	EndDate       string `json:"endDate,omitempty"`
	CompleteDate  string `json:"completeDate,omitempty"`
	OriginBoardID int    `json:"originBoardId,omitempty"`
}

// SprintList holds a paginated list of sprints returned by the Jira API.
type SprintList struct {
	MaxResults int      `json:"maxResults"`
	StartAt    int      `json:"startAt"`
	IsLast     bool     `json:"isLast"`
	Values     []Sprint `json:"values"`
}

// GetSprints retrieves all sprints for the given board ID.
// It follows pagination automatically and returns the full list.
func (c *Client) GetSprints(boardID int) ([]Sprint, error) {
	var all []Sprint
	startAt := 0
	maxResults := 50

	for {
		url := fmt.Sprintf("%s/rest/agile/1.0/board/%d/sprint?startAt=%d&maxResults=%d",
			c.baseURL, boardID, startAt, maxResults)

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("creating sprint list request: %w", err)
		}
		req.Header.Set("Accept", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("executing sprint list request: %w", err)
		}
		defer resp.Body.Close() //nolint:errcheck

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status %d fetching sprints for board %d", resp.StatusCode, boardID)
		}

		var page SprintList
		if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
			return nil, fmt.Errorf("decoding sprint list response: %w", err)
		}

		all = append(all, page.Values...)

		if page.IsLast || len(page.Values) == 0 {
			break
		}
		startAt += len(page.Values)
	}

	return all, nil
}

// GetActiveSprint returns the first active sprint for the given board ID,
// or an error if none is found.
func (c *Client) GetActiveSprint(boardID int) (*Sprint, error) {
	url := fmt.Sprintf("%s/rest/agile/1.0/board/%d/sprint?state=active", c.baseURL, boardID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating active sprint request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing active sprint request: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d fetching active sprint for board %d", resp.StatusCode, boardID)
	}

	var page SprintList
	if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
		return nil, fmt.Errorf("decoding active sprint response: %w", err)
	}

	if len(page.Values) == 0 {
		return nil, fmt.Errorf("no active sprint found for board %d", boardID)
	}

	return &page.Values[0], nil
}
