package bitbucket

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
		Email:    "user@example.com",
		APIToken: "token123",
	}
	c, err := NewClient(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.token != "token123" {
		t.Errorf("unexpected token: %s", c.token)
	}
	if c.email != "user@example.com" {
		t.Errorf("unexpected email: %s", c.email)
	}
}

func TestNewClientMissingToken(t *testing.T) {
	p := config.Profile{}
	_, err := NewClient(p)
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func newTestClient(srv *httptest.Server, email string) *Client {
	return &Client{
		httpClient: srv.Client(),
		token:      "testtoken",
		email:      email,
	}
}

func TestBasicAuthHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			t.Error("expected basic auth")
		}
		if user != "user@test.com" || pass != "testtoken" {
			t.Errorf("unexpected creds: %s:%s", user, pass)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok": true}`))
	}))
	defer srv.Close()

	c := newTestClient(srv, "user@test.com")
	// Use the full URL so the client doesn't prepend baseURL
	_, err := c.get(srv.URL + "/test")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
}

func TestBearerAuthHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer testtoken" {
			t.Errorf("expected Bearer auth, got: %s", auth)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok": true}`))
	}))
	defer srv.Close()

	c := newTestClient(srv, "")
	_, err := c.get(srv.URL + "/test")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
}

func TestDoPost(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"slug": "my-repo"}`))
	}))
	defer srv.Close()

	c := newTestClient(srv, "")
	data, err := c.post(srv.URL+"/repos", nil)
	if err != nil {
		t.Fatalf("POST failed: %v", err)
	}
	var result map[string]string
	json.Unmarshal(data, &result)
	if result["slug"] != "my-repo" {
		t.Errorf("unexpected response: %v", result)
	}
}

func TestDoPut(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"updated": true}`))
	}))
	defer srv.Close()

	c := newTestClient(srv, "")
	_, err := c.put(srv.URL+"/repos/test", nil)
	if err != nil {
		t.Fatalf("PUT failed: %v", err)
	}
}

func TestDeleteNoContent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newTestClient(srv, "")
	err := c.deleteNoContent(srv.URL + "/repos/test")
	if err != nil {
		t.Fatalf("DELETE failed: %v", err)
	}
}

func TestPostNoContent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newTestClient(srv, "")
	err := c.postNoContent(srv.URL+"/approve", nil)
	if err != nil {
		t.Fatalf("POST no content failed: %v", err)
	}
}

func TestAPIErrorResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(APIError{
			Error: struct {
				Message string `json:"message"`
			}{Message: "Repository not found"},
		})
	}))
	defer srv.Close()

	c := newTestClient(srv, "")
	_, err := c.get(srv.URL + "/repos/nonexistent")
	if err == nil {
		t.Fatal("expected error for 404")
	}
}

func TestGetRawNoAcceptJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") == "application/json" {
			t.Error("getRaw should not set Accept: application/json")
		}
		w.Write([]byte("diff --git a/file.go"))
	}))
	defer srv.Close()

	c := newTestClient(srv, "")
	data, err := c.getRaw(srv.URL + "/diff")
	if err != nil {
		t.Fatalf("getRaw failed: %v", err)
	}
	if string(data) != "diff --git a/file.go" {
		t.Errorf("unexpected raw content: %s", string(data))
	}
}

func TestGetRawAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("forbidden"))
	}))
	defer srv.Close()

	c := newTestClient(srv, "")
	_, err := c.getRaw(srv.URL + "/secret")
	if err == nil {
		t.Fatal("expected error for 403")
	}
}

func TestPaginationOptionsApplyParams(t *testing.T) {
	tests := []struct {
		name     string
		opts     *PaginationOptions
		wantPage string
		wantLen  string
	}{
		{
			name:     "nil opts",
			opts:     nil,
			wantPage: "",
			wantLen:  "",
		},
		{
			name:     "zero values",
			opts:     &PaginationOptions{},
			wantPage: "",
			wantLen:  "",
		},
		{
			name:     "with page and pagelen",
			opts:     &PaginationOptions{Page: 2, PageLen: 25},
			wantPage: "2",
			wantLen:  "25",
		},
		{
			name:     "page only",
			opts:     &PaginationOptions{Page: 3},
			wantPage: "3",
			wantLen:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := url.Values{}
			tt.opts.applyParams(params)
			if got := params.Get("page"); got != tt.wantPage {
				t.Errorf("page: got %q, want %q", got, tt.wantPage)
			}
			if got := params.Get("pagelen"); got != tt.wantLen {
				t.Errorf("pagelen: got %q, want %q", got, tt.wantLen)
			}
		})
	}
}

func TestEnsurePageLen(t *testing.T) {
	params := url.Values{}
	ensurePageLen(params)
	if params.Get("pagelen") != "50" {
		t.Errorf("expected default pagelen 50, got %s", params.Get("pagelen"))
	}

	// Should not override existing value
	params.Set("pagelen", "10")
	ensurePageLen(params)
	if params.Get("pagelen") != "10" {
		t.Errorf("expected pagelen 10, got %s", params.Get("pagelen"))
	}
}

func TestGetAllPagination(t *testing.T) {
	callCount := 0
	var srvURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		if callCount == 1 {
			json.NewEncoder(w).Encode(PaginatedResponse{
				Size:    2,
				Page:    1,
				PageLen: 1,
				Next:    srvURL + "/page2",
				Values:  json.RawMessage(`[{"id": 1}]`),
			})
		} else {
			json.NewEncoder(w).Encode(PaginatedResponse{
				Size:    2,
				Page:    2,
				PageLen: 1,
				Values:  json.RawMessage(`[{"id": 2}]`),
			})
		}
	}))
	defer srv.Close()
	srvURL = srv.URL

	c := newTestClient(srv, "")
	pages, err := c.getAll(srv.URL + "/items")
	if err != nil {
		t.Fatalf("getAll failed: %v", err)
	}
	if len(pages) != 2 {
		t.Errorf("expected 2 pages, got %d", len(pages))
	}
	if callCount != 2 {
		t.Errorf("expected 2 API calls, got %d", callCount)
	}
}
