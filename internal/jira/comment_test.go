package jira

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetComments_Success(t *testing.T) {
	expected := CommentsResponse{
		Total:      2,
		MaxResults: 50,
		Comments: []Comment{
			{
				ID:   "10001",
				Body: "First comment",
				Author: Author{
					DisplayName:  "Alice",
					EmailAddress: "alice@example.com",
				},
			},
			{
				ID:   "10002",
				Body: "Second comment",
				Author: Author{
					DisplayName:  "Bob",
					EmailAddress: "bob@example.com",
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/rest/api/2/issue/PROJ-1/comment", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client())
	result, err := client.GetComments(context.Background(), "PROJ-1")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.Total)
	assert.Len(t, result.Comments, 2)
	assert.Equal(t, "First comment", result.Comments[0].Body)
	assert.Equal(t, "Alice", result.Comments[0].Author.DisplayName)
}

func TestGetComments_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client())
	result, err := client.GetComments(context.Background(), "INVALID-999")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}

func TestGetComments_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client())
	result, err := client.GetComments(context.Background(), "PROJ-1")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unexpected status")
}
