package jira

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSprints_Success(t *testing.T) {
	sprints := SprintList{
		Values: []Sprint{
			{ID: 1, Name: "Sprint 1", State: "active", BoardID: 10},
			{ID: 2, Name: "Sprint 2", State: "closed", BoardID: 10},
		},
		Total: 2,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/rest/agile/1.0/board/10/sprint", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(sprints)
	}))
	defer server.Close()

	client := &Client{baseURL: server.URL, httpClient: server.Client()}
	result, err := client.GetSprints(10)

	require.NoError(t, err)
	assert.Len(t, result.Values, 2)
	assert.Equal(t, "Sprint 1", result.Values[0].Name)
	assert.Equal(t, "active", result.Values[0].State)
}

func TestGetSprints_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := &Client{baseURL: server.URL, httpClient: server.Client()}
	_, err := client.GetSprints(99)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "404")
}

func TestGetActiveSprint_Success(t *testing.T) {
	sprints := SprintList{
		Values: []Sprint{
			{ID: 3, Name: "Active Sprint", State: "active", BoardID: 10},
		},
		Total: 1,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "active", r.URL.Query().Get("state"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(sprints)
	}))
	defer server.Close()

	client := &Client{baseURL: server.URL, httpClient: server.Client()}
	sprint, err := client.GetActiveSprint(10)

	require.NoError(t, err)
	require.NotNil(t, sprint)
	assert.Equal(t, "Active Sprint", sprint.Name)
	assert.Equal(t, "active", sprint.State)
}

func TestGetActiveSprint_NoActive(t *testing.T) {
	sprints := SprintList{Values: []Sprint{}, Total: 0}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(sprints)
	}))
	defer server.Close()

	client := &Client{baseURL: server.URL, httpClient: server.Client()}
	sprint, err := client.GetActiveSprint(10)

	require.NoError(t, err)
	assert.Nil(t, sprint)
}
