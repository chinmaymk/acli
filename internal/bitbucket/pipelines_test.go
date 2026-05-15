package bitbucket

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// TestListPipelinesRequest verifies the request adheres to the Bitbucket
// pipelines spec (GET /repositories/{workspace}/{repo_slug}/pipelines):
// path, and query parameter names for status, target.branch, sort, and
// pagination (page/pagelen).
func TestListPipelinesRequest(t *testing.T) {
	var gotPath string
	var gotQuery url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(PaginatedResponse{
			Values: json.RawMessage(`[{"uuid": "{abc}", "build_number": 1}]`),
		})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	_, err := c.ListPipelines("ws", "repo", &ListPipelinesOptions{
		Status: "PASSED",
		Branch: "main",
		PaginationOptions: PaginationOptions{
			Page:    2,
			PageLen: 25,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantPath := "/2.0/repositories/ws/repo/pipelines"
	if gotPath != wantPath {
		t.Errorf("path: got %q, want %q", gotPath, wantPath)
	}
	checks := map[string]string{
		"status":        "PASSED",
		"target.branch": "main",
		"page":          "2",
		"pagelen":       "25",
		"sort":          "-created_on",
	}
	for k, want := range checks {
		if got := gotQuery.Get(k); got != want {
			t.Errorf("query %s: got %q, want %q", k, got, want)
		}
	}
}

func TestListPipelinesBranchOnly(t *testing.T) {
	var gotQuery url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(PaginatedResponse{
			Values: json.RawMessage(`[]`),
		})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	_, err := c.ListPipelines("ws", "repo", &ListPipelinesOptions{Branch: "feature/x"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := gotQuery.Get("target.branch"); got != "feature/x" {
		t.Errorf("target.branch: got %q, want %q", got, "feature/x")
	}
	if got := gotQuery.Get("status"); got != "" {
		t.Errorf("status should not be set, got %q", got)
	}
}

func TestRunPipelineRequestBody(t *testing.T) {
	var gotMethod string
	var gotPath string
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(Pipeline{UUID: "{abc}", BuildNumber: 42})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	_, err := c.RunPipeline("ws", "repo", NewBranchPipelineRequest("main"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "POST" {
		t.Errorf("method: got %q, want POST", gotMethod)
	}
	if gotPath != "/2.0/repositories/ws/repo/pipelines" {
		t.Errorf("path: got %q", gotPath)
	}
	target, ok := gotBody["target"].(map[string]any)
	if !ok {
		t.Fatalf("missing target in body: %#v", gotBody)
	}
	if target["type"] != "pipeline_ref_target" {
		t.Errorf("target.type: got %v, want pipeline_ref_target", target["type"])
	}
	if target["ref_type"] != "branch" {
		t.Errorf("target.ref_type: got %v, want branch", target["ref_type"])
	}
	if target["ref_name"] != "main" {
		t.Errorf("target.ref_name: got %v, want main", target["ref_name"])
	}
}

func TestStopPipelinePath(t *testing.T) {
	var gotMethod, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	if err := c.StopPipeline("ws", "repo", "{abc}"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "POST" {
		t.Errorf("method: got %q, want POST", gotMethod)
	}
	if gotPath != "/2.0/repositories/ws/repo/pipelines/{abc}/stopPipeline" {
		t.Errorf("path: got %q", gotPath)
	}
}
