//go:build integration
// +build integration

package jira_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ankitpokhrel/jira-cli/internal/auth"
	"github.com/ankitpokhrel/jira-cli/internal/jira"
)

// TestGetWorklogs_Integration exercises the real Jira API.
// Set JIRA_BASE_URL, JIRA_USER, JIRA_TOKEN and JIRA_TEST_ISSUE env vars before running.
func TestGetWorklogs_Integration(t *testing.T) {
	baseURL := os.Getenv("JIRA_BASE_URL")
	user := os.Getenv("JIRA_USER")
	token := os.Getenv("JIRA_TOKEN")
	issueKey := os.Getenv("JIRA_TEST_ISSUE")

	if baseURL == "" || user == "" || token == "" || issueKey == "" {
		t.Skip("integration env vars not set")
	}

	httpClient := auth.NewHTTPClient(user, token)
	client := jira.NewClient(baseURL, httpClient)

	result, err := client.GetWorklogs(issueKey)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.GreaterOrEqual(t, result.Total, 0)
}
