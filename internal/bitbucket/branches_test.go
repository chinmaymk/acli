package bitbucket

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListBranches(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PaginatedResponse{
			Values: json.RawMessage(`[{"name": "main"}, {"name": "develop"}]`),
		})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	branches, err := c.ListBranches("ws", "repo", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 2 {
		t.Fatalf("expected 2 branches, got %d", len(branches))
	}
	if branches[0].Name != "main" {
		t.Errorf("unexpected branch name: %s", branches[0].Name)
	}
}

func TestListBranchesWithQuery(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("q") != `name~"feature"` {
			t.Errorf("unexpected q: %s", r.URL.Query().Get("q"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PaginatedResponse{
			Values: json.RawMessage(`[{"name": "feature/login"}]`),
		})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	branches, err := c.ListBranches("ws", "repo", &ListBranchesOptions{Q: `name~"feature"`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Errorf("expected 1 branch, got %d", len(branches))
	}
}

func TestListBranchesAll(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		if callCount == 1 {
			json.NewEncoder(w).Encode(PaginatedResponse{
				Next:   "https://api.bitbucket.org/2.0/repositories/ws/repo/refs/branches?page=2",
				Values: json.RawMessage(`[{"name": "main"}]`),
			})
		} else {
			json.NewEncoder(w).Encode(PaginatedResponse{
				Values: json.RawMessage(`[{"name": "develop"}]`),
			})
		}
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	branches, err := c.ListBranches("ws", "repo", &ListBranchesOptions{
		PaginationOptions: PaginationOptions{All: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 2 {
		t.Errorf("expected 2 branches, got %d", len(branches))
	}
}

func TestGetBranch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Branch{Name: "main"})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	branch, err := c.GetBranch("ws", "repo", "main")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if branch.Name != "main" {
		t.Errorf("unexpected name: %s", branch.Name)
	}
}

func TestCreateBranch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Branch{Name: "feature/new"})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	req := &CreateBranchRequest{Name: "feature/new"}
	req.Target.Hash = "abc123"
	branch, err := c.CreateBranch("ws", "repo", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if branch.Name != "feature/new" {
		t.Errorf("unexpected name: %s", branch.Name)
	}
}

func TestDeleteBranch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	err := c.DeleteBranch("ws", "repo", "old-branch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListTags(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PaginatedResponse{
			Values: json.RawMessage(`[{"name": "v1.0.0"}, {"name": "v1.1.0"}]`),
		})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	tags, err := c.ListTags("ws", "repo", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
}

func TestGetTag(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Tag{Name: "v1.0.0", Message: "Release 1.0"})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	tag, err := c.GetTag("ws", "repo", "v1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag.Name != "v1.0.0" {
		t.Errorf("unexpected name: %s", tag.Name)
	}
}

func TestCreateTag(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Tag{Name: "v2.0.0"})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	req := &CreateTagRequest{Name: "v2.0.0", Message: "Release 2.0"}
	req.Target.Hash = "def456"
	tag, err := c.CreateTag("ws", "repo", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag.Name != "v2.0.0" {
		t.Errorf("unexpected name: %s", tag.Name)
	}
}

func TestDeleteTag(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	err := c.DeleteTag("ws", "repo", "v0.1.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListTagsAll(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		if callCount == 1 {
			json.NewEncoder(w).Encode(PaginatedResponse{
				Next:   "https://api.bitbucket.org/2.0/repositories/ws/repo/refs/tags?page=2",
				Values: json.RawMessage(`[{"name": "v1.0"}]`),
			})
		} else {
			json.NewEncoder(w).Encode(PaginatedResponse{
				Values: json.RawMessage(`[{"name": "v2.0"}]`),
			})
		}
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	tags, err := c.ListTags("ws", "repo", &ListTagsOptions{
		PaginationOptions: PaginationOptions{All: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(tags))
	}
}
