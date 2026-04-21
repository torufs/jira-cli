/*
Package jira provides a client for interacting with the Jira REST API.

# Worklog

The worklog functionality allows retrieval of time-tracking entries logged
against a specific Jira issue.

# Usage

Create a client and call GetWorklogs with an issue key:

	client := jira.NewClient(baseURL, httpClient)
	worklogs, err := client.GetWorklogs("PROJECT-123")
	if err != nil {
		log.Fatal(err)
	}
	for _, w := range worklogs.Worklogs {
		fmt.Printf("%s logged %s: %s\n", w.Author.DisplayName, w.TimeSpent, w.Comment)
	}

# Errors

Returns an error if the issue is not found (HTTP 404) or if the server
returns an unexpected status code.
*/
package jira
