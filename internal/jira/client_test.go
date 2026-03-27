package jira

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/chinmaymk/acli/internal/config"
)

func TestNewClientSuccess(t *testing.T) {
	p := config.Profile{
		AtlassianURL: "https://test.atlassian.net",
		Email:        "user@example.com",
		APIToken:     "token123",
	}
	c, err := NewClient(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.BaseURL != "https://test.atlassian.net" {
		t.Errorf("unexpected base URL: %s", c.BaseURL)
	}
	if c.Email != "user@example.com" {
		t.Errorf("unexpected email: %s", c.Email)
	}
}

func TestNewClientTrimsTrailingSlash(t *testing.T) {
	p := config.Profile{
		AtlassianURL: "https://test.atlassian.net/",
		APIToken:     "token123",
	}
	c, err := NewClient(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.BaseURL != "https://test.atlassian.net" {
		t.Errorf("expected trailing slash removed, got: %s", c.BaseURL)
	}
}

func TestNewClientMissingURL(t *testing.T) {
	p := config.Profile{APIToken: "token123"}
	_, err := NewClient(p)
	if err == nil {
		t.Fatal("expected error for missing URL")
	}
}

func TestNewClientMissingToken(t *testing.T) {
	p := config.Profile{AtlassianURL: "https://test.atlassian.net"}
	_, err := NewClient(p)
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestBasicAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			t.Error("expected basic auth")
		}
		if user != "user@example.com" || pass != "token123" {
			t.Errorf("unexpected creds: %s:%s", user, pass)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
	}))
	defer srv.Close()

	c := &Client{
		BaseURL:    srv.URL,
		Email:      "user@example.com",
		APIToken:   "token123",
		HTTPClient: srv.Client(),
	}

	var result map[string]string
	if err := c.Get("/test", nil, &result); err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	if result["ok"] != "true" {
		t.Error("unexpected response")
	}
}

func TestBearerAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer mytoken" {
			t.Errorf("expected Bearer auth, got: %s", auth)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
	}))
	defer srv.Close()

	c := &Client{
		BaseURL:    srv.URL,
		Email:      "", // no email = bearer auth
		APIToken:   "mytoken",
		HTTPClient: srv.Client(),
	}

	var result map[string]string
	if err := c.Get("/test", nil, &result); err != nil {
		t.Fatalf("GET failed: %v", err)
	}
}

func TestGetWithQueryParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("jql") != "project=TEST" {
			t.Errorf("unexpected query: %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"found": "yes"})
	}))
	defer srv.Close()

	c := &Client{BaseURL: srv.URL, APIToken: "tok", HTTPClient: srv.Client()}

	q := url.Values{"jql": {"project=TEST"}}
	var result map[string]string
	if err := c.Get("/search", q, &result); err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	if result["found"] != "yes" {
		t.Error("unexpected response")
	}
}

func TestPost(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("expected application/json content type")
		}
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["summary"] != "Test issue" {
			t.Errorf("unexpected body: %v", body)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"key": "TEST-1"})
	}))
	defer srv.Close()

	c := &Client{BaseURL: srv.URL, APIToken: "tok", HTTPClient: srv.Client()}

	var result map[string]string
	err := c.Post("/issue", map[string]string{"summary": "Test issue"}, &result)
	if err != nil {
		t.Fatalf("POST failed: %v", err)
	}
	if result["key"] != "TEST-1" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestPut(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := &Client{BaseURL: srv.URL, APIToken: "tok", HTTPClient: srv.Client()}
	if err := c.Put("/issue/TEST-1", map[string]string{"summary": "Updated"}, nil); err != nil {
		t.Fatalf("PUT failed: %v", err)
	}
}

func TestDelete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := &Client{BaseURL: srv.URL, APIToken: "tok", HTTPClient: srv.Client()}
	if err := c.Delete("/issue/TEST-1", nil); err != nil {
		t.Fatalf("DELETE failed: %v", err)
	}
}

func TestPatch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
	}))
	defer srv.Close()

	c := &Client{BaseURL: srv.URL, APIToken: "tok", HTTPClient: srv.Client()}
	var result map[string]string
	if err := c.Patch("/issue/TEST-1", map[string]string{"summary": "Patched"}, &result); err != nil {
		t.Fatalf("PATCH failed: %v", err)
	}
}

func TestGetRaw(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("raw content here"))
	}))
	defer srv.Close()

	c := &Client{BaseURL: srv.URL, APIToken: "tok", HTTPClient: srv.Client()}
	data, err := c.GetRaw("/export", nil)
	if err != nil {
		t.Fatalf("GetRaw failed: %v", err)
	}
	if string(data) != "raw content here" {
		t.Errorf("unexpected raw content: %s", string(data))
	}
}

func TestAPIErrorHandling(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errorMessages": []string{"Issue Does Not Exist"},
			"errors":        map[string]string{},
		})
	}))
	defer srv.Close()

	c := &Client{BaseURL: srv.URL, APIToken: "tok", HTTPClient: srv.Client()}
	err := c.Get("/issue/NOPE-1", nil, nil)
	if err == nil {
		t.Fatal("expected error for 404")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("expected status 404, got %d", apiErr.StatusCode)
	}
	if len(apiErr.ErrorMessages) == 0 || apiErr.ErrorMessages[0] != "Issue Does Not Exist" {
		t.Errorf("unexpected error messages: %v", apiErr.ErrorMessages)
	}
}

func TestAPIErrorString(t *testing.T) {
	tests := []struct {
		name     string
		err      APIError
		contains string
	}{
		{
			name:     "with error messages",
			err:      APIError{StatusCode: 400, ErrorMessages: []string{"bad field"}},
			contains: "bad field",
		},
		{
			name:     "with field errors",
			err:      APIError{StatusCode: 400, Errors: map[string]string{"summary": "required"}},
			contains: "summary: required",
		},
		{
			name:     "empty errors",
			err:      APIError{StatusCode: 500},
			contains: "HTTP 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tt.err.Error()
			if len(msg) == 0 {
				t.Error("expected non-empty error string")
			}
		})
	}
}

func TestGetRawAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("forbidden"))
	}))
	defer srv.Close()

	c := &Client{BaseURL: srv.URL, APIToken: "tok", HTTPClient: srv.Client()}
	_, err := c.GetRaw("/secret", nil)
	if err == nil {
		t.Fatal("expected error for 403")
	}
}

func TestDeleteWithBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("expected application/json content type")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"deleted": "true"})
	}))
	defer srv.Close()

	c := &Client{BaseURL: srv.URL, APIToken: "tok", HTTPClient: srv.Client()}
	var result map[string]string
	err := c.DeleteWithBody("/bulk", map[string][]string{"ids": {"1", "2"}}, &result)
	if err != nil {
		t.Fatalf("DeleteWithBody failed: %v", err)
	}
}
