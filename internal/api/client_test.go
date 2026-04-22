package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient("https://test.atlassian.net/", "user@test.com", "token123")
	if c.BaseURL != "https://test.atlassian.net" {
		t.Errorf("expected trailing slash removed, got: %s", c.BaseURL)
	}
	if c.Email != "user@test.com" {
		t.Errorf("unexpected email: %s", c.Email)
	}
	if c.APIToken != "token123" {
		t.Errorf("unexpected token: %s", c.APIToken)
	}
	if c.HTTPClient == nil {
		t.Error("expected non-nil HTTP client")
	}
}

func TestConfluenceV2BasicAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify path includes /wiki/api/v2
		if r.URL.Path != "/wiki/api/v2/spaces" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		// Verify basic auth
		user, pass, ok := r.BasicAuth()
		if !ok {
			t.Error("expected basic auth")
		}
		if user != "user@test.com" || pass != "token123" {
			t.Errorf("unexpected creds: %s:%s", user, pass)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Error("expected Accept: application/json")
		}
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "user@test.com", "token123")
	c.HTTPClient = srv.Client()

	data, err := c.ConfluenceV2("GET", "/spaces", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result map[string]string
	json.Unmarshal(data, &result)
	if result["ok"] != "true" {
		t.Error("unexpected response")
	}
}

func TestConfluenceV2BearerAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer mytoken" {
			t.Errorf("expected Bearer auth, got: %s", auth)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok": true}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "", "mytoken")
	c.HTTPClient = srv.Client()

	_, err := c.ConfluenceV2("GET", "/spaces", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConfluenceV2WithQueryParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("limit") != "25" {
			t.Errorf("unexpected query: %s", r.URL.RawQuery)
		}
		if r.URL.Query().Get("status") != "current" {
			t.Errorf("unexpected status param: %s", r.URL.Query().Get("status"))
		}
		w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "", "tok")
	c.HTTPClient = srv.Client()

	q := url.Values{"limit": {"25"}, "status": {"current"}}
	_, err := c.ConfluenceV2("GET", "/pages", q, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConfluenceV2WithBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("expected application/json content type for POST with body")
		}
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["title"] != "My Page" {
			t.Errorf("unexpected body: %v", body)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": "123"})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "", "tok")
	c.HTTPClient = srv.Client()

	body := map[string]string{"title": "My Page"}
	data, err := c.ConfluenceV2("POST", "/pages", nil, body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result map[string]string
	json.Unmarshal(data, &result)
	if result["id"] != "123" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestConfluenceV2NoContentTypeWithoutBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "" {
			t.Error("should not set Content-Type when body is nil")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "", "tok")
	c.HTTPClient = srv.Client()

	_, err := c.ConfluenceV2("DELETE", "/pages/123", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConfluenceV2APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("access denied"))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "", "tok")
	c.HTTPClient = srv.Client()

	_, err := c.ConfluenceV2("GET", "/spaces", nil, nil)
	if err == nil {
		t.Fatal("expected error for 403")
	}
}

func TestConfluenceV2PutMethod(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Write([]byte(`{"updated": true}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "", "tok")
	c.HTTPClient = srv.Client()

	_, err := c.ConfluenceV2("PUT", "/pages/123", nil, map[string]string{"title": "Updated"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConfluenceV1BasicAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/wiki/rest/api/content/search" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		user, pass, ok := r.BasicAuth()
		if !ok {
			t.Error("expected basic auth")
		}
		if user != "user@test.com" || pass != "token123" {
			t.Errorf("unexpected creds: %s:%s", user, pass)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Error("expected Accept: application/json")
		}
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"results": []interface{}{}, "totalSize": 0})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "user@test.com", "token123")
	c.HTTPClient = srv.Client()

	_, err := c.ConfluenceV1("GET", "/content/search", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConfluenceV1BearerAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer mytoken" {
			t.Errorf("expected Bearer auth, got: %s", auth)
		}
		w.Write([]byte(`{"results":[],"totalSize":0}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "", "mytoken")
	c.HTTPClient = srv.Client()

	_, err := c.ConfluenceV1("GET", "/content/search", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConfluenceV1WithQueryParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("cql") != "type=page AND space=DEV" {
			t.Errorf("unexpected cql param: %s", r.URL.Query().Get("cql"))
		}
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("unexpected limit: %s", r.URL.Query().Get("limit"))
		}
		w.Write([]byte(`{"results":[],"totalSize":0}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "", "tok")
	c.HTTPClient = srv.Client()

	q := url.Values{"cql": {"type=page AND space=DEV"}, "limit": {"10"}}
	_, err := c.ConfluenceV1("GET", "/content/search", q, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConfluenceV1APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "", "bad-token")
	c.HTTPClient = srv.Client()

	_, err := c.ConfluenceV1("GET", "/content/search", nil, nil)
	if err == nil {
		t.Fatal("expected error for 401")
	}
}
