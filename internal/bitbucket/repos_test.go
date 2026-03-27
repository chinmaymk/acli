package bitbucket

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListRepositoriesSimple(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PaginatedResponse{
			Size:    1,
			Page:    1,
			PageLen: 50,
			Values:  json.RawMessage(`[{"slug": "my-repo", "full_name": "myws/my-repo", "is_private": true}]`),
		})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	repos, err := c.ListRepositories("myws", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repos) != 1 {
		t.Fatalf("expected 1 repo, got %d", len(repos))
	}
	if repos[0].Slug != "my-repo" {
		t.Errorf("unexpected slug: %s", repos[0].Slug)
	}
	if repos[0].FullName != "myws/my-repo" {
		t.Errorf("unexpected full_name: %s", repos[0].FullName)
	}
	if !repos[0].IsPrivate {
		t.Error("expected private repo")
	}
}

func TestListRepositoriesWithOpts(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("role") != "member" {
			t.Errorf("unexpected role: %s", r.URL.Query().Get("role"))
		}
		if r.URL.Query().Get("sort") != "name" {
			t.Errorf("unexpected sort: %s", r.URL.Query().Get("sort"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PaginatedResponse{
			Values: json.RawMessage(`[{"slug": "repo1"}]`),
		})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	repos, err := c.ListRepositories("ws", &ListReposOptions{
		Role:    "member",
		Sort:    "name",
		Page:    2,
		PageLen: 10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repos) != 1 {
		t.Errorf("expected 1 repo, got %d", len(repos))
	}
}

func TestListRepositoriesAll(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		if callCount == 1 {
			json.NewEncoder(w).Encode(PaginatedResponse{
				Next:   "https://api.bitbucket.org/2.0/repositories/ws?page=2&pagelen=50",
				Values: json.RawMessage(`[{"slug": "repo1"}]`),
			})
		} else {
			json.NewEncoder(w).Encode(PaginatedResponse{
				Values: json.RawMessage(`[{"slug": "repo2"}]`),
			})
		}
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	repos, err := c.ListRepositories("ws", &ListReposOptions{All: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repos) != 2 {
		t.Errorf("expected 2 repos, got %d", len(repos))
	}
}

func TestGetRepository(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Repository{
			Slug:     "my-repo",
			FullName: "ws/my-repo",
			SCM:      "git",
		})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	repo, err := c.GetRepository("ws", "my-repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.Slug != "my-repo" {
		t.Errorf("unexpected slug: %s", repo.Slug)
	}
	if repo.SCM != "git" {
		t.Errorf("unexpected scm: %s", repo.SCM)
	}
}

func TestCreateRepository(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Repository{Slug: "new-repo", FullName: "ws/new-repo"})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	repo, err := c.CreateRepository("ws", &CreateRepoRequest{
		SCM:       "git",
		Name:      "new-repo",
		IsPrivate: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.Slug != "new-repo" {
		t.Errorf("unexpected slug: %s", repo.Slug)
	}
}

func TestCreateRepositoryWithSlug(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Repository{Slug: "custom-slug"})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	repo, err := c.CreateRepository("ws", &CreateRepoRequest{
		SCM:  "git",
		Name: "My Repo",
		Slug: "custom-slug",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.Slug != "custom-slug" {
		t.Errorf("unexpected slug: %s", repo.Slug)
	}
}

func TestDeleteRepository(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	err := c.DeleteRepository("ws", "old-repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestForkRepository(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Repository{Slug: "forked-repo"})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	repo, err := c.ForkRepository("ws", "original-repo", &ForkRepoRequest{Name: "forked-repo"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.Slug != "forked-repo" {
		t.Errorf("unexpected slug: %s", repo.Slug)
	}
}

func TestListForks(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PaginatedResponse{
			Values: json.RawMessage(`[{"slug": "fork1"}]`),
		})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	forks, err := c.ListForks("ws", "my-repo", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(forks) != 1 {
		t.Errorf("expected 1 fork, got %d", len(forks))
	}
}
