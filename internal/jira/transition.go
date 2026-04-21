package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Transition represents a Jira issue transition.
type Transition struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// TransitionsResponse is the response from the Jira transitions API.
type TransitionsResponse struct {
	Transitions []Transition `json:"transitions"`
}

// GetTransitions returns all available transitions for a given issue.
func (c *Client) GetTransitions(issueKey string) ([]Transition, error) {
	url := fmt.Sprintf("%s/rest/api/2/issue/%s/transitions", c.baseURL, issueKey)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get transitions request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("issue %q not found", issueKey)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result TransitionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode transitions response: %w", err)
	}

	return result.Transitions, nil
}

// DoTransition performs a transition on a given issue by transition ID.
func (c *Client) DoTransition(issueKey, transitionID string) error {
	url := fmt.Sprintf("%s/rest/api/2/issue/%s/transitions", c.baseURL, issueKey)

	body := map[string]interface{}{
		"transition": map[string]string{"id": transitionID},
	}

	resp, err := c.postJSON(url, body)
	if err != nil {
		return fmt.Errorf("do transition request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("issue %q not found", issueKey)
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
