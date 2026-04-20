package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ankitpokhrel/jira-cli/internal/auth"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		creds   auth.Credentials
		wantErr bool
	}{
		{"basic ok", auth.Credentials{Method: auth.MethodBasic, Username: "user", Token: "pass"}, false},
		{"basic missing user", auth.Credentials{Method: auth.MethodBasic, Token: "pass"}, true},
		{"basic missing token", auth.Credentials{Method: auth.MethodBasic, Username: "user"}, true},
		{"token ok", auth.Credentials{Method: auth.MethodToken, Token: "mytoken"}, false},
		{"token missing token", auth.Credentials{Method: auth.MethodToken}, true},
		{"unknown method", auth.Credentials{Method: "oauth"}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.creds.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestBasicAuthTransport(t *testing.T) {
	var gotAuth string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	creds := &auth.Credentials{Method: auth.MethodBasic, Username: "alice", Token: "secret"}
	client := auth.NewHTTPClient(creds, 5*time.Second)
	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if gotAuth == "" {
		t.Error("expected Authorization header, got empty string")
	}
}

func TestTokenAuthTransport(t *testing.T) {
	var gotAuth string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	creds := &auth.Credentials{Method: auth.MethodToken, Token: "mytoken123"}
	client := auth.NewHTTPClient(creds, 5*time.Second)
	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()
	if gotAuth != "Bearer mytoken123" {
		t.Errorf("expected 'Bearer mytoken123', got %q", gotAuth)
	}
}
