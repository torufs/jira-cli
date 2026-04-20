package jira

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearch_Success(t *testing.T) {
	expected := SearchResult{
		Total:      2,
		StartAt:    0,
		MaxResults: 50,
		Issues: []Issue{
			{
				Key: "PROJ-1",
				Fields: IssueFields{
					Summary: "First issue",
					Status:  Status{Name: "Open"},
				},
			},
			{
				Key: "PROJ-2",
				Fields: IssueFields{
					Summary: "Second issue",
					Status:  Status{Name: "In Progress"},
					Assignee: &User{DisplayName: "Jane Doe", EmailAddress: "jane@example.com"},
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/rest/api/2/search", r.URL.Path)
		assert.Equal(t, "project = PROJ", r.URL.Query().Get("jql"))
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client())
	result, err := client.Search(SearchRequest{JQL: "project = PROJ"})

	require.NoError(t, err)
	assert.Equal(t, 2, result.Total)
	assert.Len(t, result.Issues, 2)
	assert.Equal(t, "PROJ-1", result.Issues[0].Key)
	assert.Equal(t, "PROJ-2", result.Issues[1].Key)
	assert.Equal(t, "Jane Doe", result.Issues[1].Fields.Assignee.DisplayName)
}

func TestSearch_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client())
	result, err := client.Search(SearchRequest{JQL: "project = PROJ"})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "401")
}

func TestSearch_DefaultMaxResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "50", r.URL.Query().Get("maxResults"))
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(SearchResult{})
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client())
	_, err := client.Search(SearchRequest{JQL: "project = PROJ"})
	require.NoError(t, err)
}
