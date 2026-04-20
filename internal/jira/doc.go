// Package jira provides a lightweight client for interacting with the
// Jira REST API (v2). It exposes typed structs for common Jira resources
// such as issues, statuses, assignees, and priorities, along with a
// Client that wraps an *http.Client so that authentication transports
// defined in internal/auth can be composed transparently.
//
// Basic usage:
//
//	httpClient := auth.NewHTTPClient(cfg)
//	client := jira.NewClient(cfg.Server, httpClient)
//	issue, err := client.GetIssue("PROJ-42")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(issue.Fields.Summary)
package jira
