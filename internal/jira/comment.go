package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Comment represents a Jira issue comment.
type Comment struct {
	ID      string `json:"id"`
	Body    string `json:"body"`
	Author  Author `json:"author"`
	Created string `json:"created"`
	Updated string `json:"updated"`
}

// Author represents the author of a comment.
type Author struct {
	DisplayName string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
}

// CommentsResponse is the response wrapper for listing comments.
type CommentsResponse struct {
	Comments   []Comment `json:"comments"`
	Total      int       `json:"total"`
	MaxResults int       `json:"maxResults"`
	StartAt    int       `json:"startAt"`
}

// AddCommentRequest is the payload for adding a comment.
type AddCommentRequest struct {
	Body string `json:"body"`
}

// GetComments retrieves all comments for a given issue key.
func (c *Client) GetComments(ctx context.Context, issueKey string) (*CommentsResponse, error) {
	url := fmt.Sprintf("%s/rest/api/2/issue/%s/comment", c.baseURL, issueKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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
		return nil, fmt.Errorf("issue %q not found", issueKey)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var result CommentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &result, nil
}
