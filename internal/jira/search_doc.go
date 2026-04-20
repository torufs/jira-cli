// Package jira provides a client for interacting with the Jira REST API.
//
// # Search
//
// The Search functionality allows querying Jira issues using JQL
// (Jira Query Language). Results are paginated and can be controlled
// via MaxResults and StartAt fields in SearchRequest.
//
// Basic usage:
//
//	client := jira.NewClient(baseURL, httpClient)
//
//	result, err := client.Search(jira.SearchRequest{
//		JQL:        "project = MYPROJ AND status = Open",
//		MaxResults: 25,
//		StartAt:    0,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for _, issue := range result.Issues {
//		fmt.Printf("%s: %s\n", issue.Key, issue.Fields.Summary)
//	}
//
// If MaxResults is not set, it defaults to 50.
package jira
