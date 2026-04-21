package jira

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBoards_Success(t *testing.T) {
	boards := BoardList{
		Values: []Board{
			{ID: 1, Name: "Team Alpha Board", Type: "scrum"},
			{ID: 2, Name: "Team Beta Board", Type: "kanban"},
		},
		Total: 2,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(boards)
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client())
	result, err := client.GetBoards()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 boards, got %d", len(result))
	}
	if result[0].Name != "Team Alpha Board" {
		t.Errorf("expected 'Team Alpha Board', got %q", result[0].Name)
	}
}

func TestGetBoards_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client())
	_, err := client.GetBoards()

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestGetBoard_Success(t *testing.T) {
	board := Board{ID: 42, Name: "My Scrum Board", Type: "scrum"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(board)
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client())
	result, err := client.GetBoard(42)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ID != 42 {
		t.Errorf("expected board ID 42, got %d", result.ID)
	}
	if result.Type != "scrum" {
		t.Errorf("expected type 'scrum', got %q", result.Type)
	}
}

func TestGetBoard_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client())
	_, err := client.GetBoard(99)

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}
