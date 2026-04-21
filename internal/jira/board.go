package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Board represents a Jira board.
type Board struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// BoardList holds a list of boards returned by the Jira API.
type BoardList struct {
	Values []Board `json:"values"`
	Total  int     `json:"total"`
}

// GetBoards retrieves all boards accessible to the authenticated user.
// It returns a slice of Board and any error encountered.
func (c *Client) GetBoards() ([]Board, error) {
	url := fmt.Sprintf("%s/rest/agile/1.0/board", c.baseURL)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("boards not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var boardList BoardList
	if err := json.NewDecoder(resp.Body).Decode(&boardList); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return boardList.Values, nil
}

// GetBoard retrieves a single board by its ID.
func (c *Client) GetBoard(boardID int) (*Board, error) {
	url := fmt.Sprintf("%s/rest/agile/1.0/board/%d", c.baseURL, boardID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("board %d not found", boardID)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var board Board
	if err := json.NewDecoder(resp.Body).Decode(&board); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &board, nil
}
