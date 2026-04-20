// Package auth provides authentication helpers for the Jira CLI.
package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Method represents an authentication method.
type Method string

const (
	// MethodBasic uses username and password (or API token).
	MethodBasic Method = "basic"
	// MethodToken uses a personal access token (Bearer).
	MethodToken Method = "token"
)

// Credentials holds authentication details.
type Credentials struct {
	Method   Method
	Username string
	Token    string
}

// Validate returns an error if the credentials are incomplete.
func (c *Credentials) Validate() error {
	switch c.Method {
	case MethodBasic:
		if c.Username == "" {
			return errors.New("auth: username is required for basic auth")
		}
		if c.Token == "" {
			return errors.New("auth: token/password is required for basic auth")
		}
	case MethodToken:
		if c.Token == "" {
			return errors.New("auth: token is required for token auth")
		}
	default:
		return fmt.Errorf("auth: unknown method %q", c.Method)
	}
	return nil
}

// Transport returns an http.RoundTripper that injects auth headers.
func (c *Credentials) Transport() http.RoundTripper {
	return &authTransport{
		creds: c,
		base:  &http.Transport{},
	}
}

// NewHTTPClient returns an *http.Client configured with the credentials.
func NewHTTPClient(creds *Credentials, timeout time.Duration) *http.Client {
	return &http.Client{
		Transport: creds.Transport(),
		Timeout:   timeout,
	}
}

type authTransport struct {
	creds *Credentials
	base  http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	r := req.Clone(req.Context())
	switch t.creds.Method {
	case MethodBasic:
		r.SetBasicAuth(t.creds.Username, t.creds.Token)
	case MethodToken:
		r.Header.Set("Authorization", "Bearer "+t.creds.Token)
	}
	return t.base.RoundTrip(r)
}
