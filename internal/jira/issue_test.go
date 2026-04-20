package jira_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankitpokhrel/jira-cli/internal/jira"
)

func TestGetIssue_Success(t *testing.T) {
	expected := jira.Issue{
		Key: "PROJ-123",
		Fields: jira.IssueFields{
			Summary:  "Test issue summary",
			Status:   jira.Status{Name: "In Progress"},
			Assignee: jira.Assignee{DisplayName: "Jane Doe", EmailAddress: "jane@example.com"},
			Priority: jira.Priority{Name: "High"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/2/issue/PROJ-123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	client := jira.NewClient(server.URL, server.Client())
	issue, err := client.GetIssue("PROJ-123")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if issue.Key != expected.Key {
		t.Errorf("expected key %s, got %s", expected.Key, issue.Key)
	}
	if issue.Fields.Summary != expected.Fields.Summary {
		t.Errorf("expected summary %q, got %q", expected.Fields.Summary, issue.Fields.Summary)
	}
	if issue.Fields.Status.Name != expected.Fields.Status.Name {
		t.Errorf("expected status %q, got %q", expected.Fields.Status.Name, issue.Fields.Status.Name)
	}
}

func TestGetIssue_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := jira.NewClient(server.URL, server.Client())
	_, err := client.GetIssue("PROJ-999")
	if err == nil {
		t.Fatal("expected an error for 404 response, got nil")
	}
}

func TestNewClient_DefaultHTTPClient(t *testing.T) {
	client := jira.NewClient("https://example.atlassian.net", nil)
	if client == nil {
		t.Fatal("expected non-nil client")
	}
	if client.BaseURL != "https://example.atlassian.net" {
		t.Errorf("unexpected base URL: %s", client.BaseURL)
	}
}
