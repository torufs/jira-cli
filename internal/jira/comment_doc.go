// Package jira provides a client for interacting with the Jira REST API v2.
//
// # Comment
//
// The comment module exposes functionality for retrieving comments on Jira issues.
//
// # Usage
//
// Create a client using NewClient, then call GetComments with a valid issue key:
//
//	client := jira.NewClient(baseURL, httpClient)
//	comments, err := client.GetComments(ctx, "PROJ-123")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, c := range comments.Comments {
//	    fmt.Printf("%s: %s\n", c.Author.DisplayName, c.Body)
//	}
//
// # Types
//
//   - Comment: represents a single comment with author info and timestamps.
//   - CommentsResponse: paginated response wrapping a slice of Comment.
//   - Author: holds display name and email of the comment author.
package jira
