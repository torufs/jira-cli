package jira

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTransitions_Success(t *testing.T) {
	expected := []Transition{
		{ID: "11", Name: "To Do"},
		{ID: "21", Name: "In Progress"},
		{ID: "31", Name: "Done"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/rest/api/2/issue/TEST-1/transitions", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(TransitionsResponse{Transitions: expected})
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client())
	transitions, err := client.GetTransitions("TEST-1")

	require.NoError(t, err)
	assert.Equal(t, expected, transitions)
}

func TestGetTransitions_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client())
	_, err := client.GetTransitions("INVALID-999")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestDoTransition_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/rest/api/2/issue/TEST-1/transitions", r.URL.Path)

		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		transition := body["transition"].(map[string]interface{})
		assert.Equal(t, "21", transition["id"])

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client())
	err := client.DoTransition("TEST-1", "21")

	require.NoError(t, err)
}

func TestDoTransition_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client())
	err := client.DoTransition("INVALID-999", "21")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
