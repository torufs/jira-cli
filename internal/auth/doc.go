// Package auth implements authentication strategies for communicating
// with a Jira instance.
//
// Two methods are supported:
//
//	- Basic authentication: uses a username combined with an API token
//	  (or password for on-premise Jira Server installations).
//
//	- Token authentication: uses a personal access token sent as a
//	  Bearer header, typically used with Jira Data Center / Server.
//
// Usage:
//
//	creds := &auth.Credentials{
//	    Method:   auth.MethodBasic,
//	    Username: "alice@example.com",
//	    Token:    "my-api-token",
//	}
//	if err := creds.Validate(); err != nil {
//	    log.Fatal(err)
//	}
//	client := auth.NewHTTPClient(creds, 10*time.Second)
package auth
