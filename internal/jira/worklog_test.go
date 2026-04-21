package jira

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetWorklogs_Success(t *testing.T) {
	expected := WorklogResponse{
		StartAt:    0,
		MaxResults: 20,
		Total:      1,
		Worklogs: []Worklog{
			{
				ID:               "10001",
				Comment:          "Fixed the bug",
				Started:          "2024-01-15T09:00:00.000+0000",
				TimeSpent:        "2h",
				TimeSpentSeconds: 7200,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/rest/api/2/issue/TEST-1/worklog", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client())
	result, err := client.GetWorklogs("TEST-1")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.Total)
	assert.Len(t, result.Worklogs, 1)
	assert.Equal(t, "Fixed the bug", result.Worklogs[0].Comment)
	assert.Equal(t, 7200, result.Worklogs[0].TimeSpentSeconds)
}

func TestGetWorklogs_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client())
	result, err := client.GetWorklogs("MISSING-99")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "MISSING-99")
}

func TestGetWorklogs_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client())
	result, err := client.GetWorklogs("TEST-2")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "500")
}
