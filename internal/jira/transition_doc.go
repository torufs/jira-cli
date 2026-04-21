// Package jira provides a client for interacting with the Jira REST API.
//
// # Transitions
//
// The transition functions allow callers to query and perform workflow
// transitions on Jira issues.
//
// # GetTransitions
//
// GetTransitions returns all transitions available for a given issue key.
// The list reflects the transitions permitted for the authenticated user
// based on the current issue status and project workflow configuration.
//
//	Transition{
//	    ID:   "21",
//	    Name: "In Progress",
//	}
//
// # DoTransition
//
// DoTransition executes a workflow transition on an issue. The transitionID
// must match one of the IDs returned by GetTransitions for that issue.
// On success the issue status is updated in Jira and nil is returned.
package jira
