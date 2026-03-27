package jira

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllProjects(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/3/project" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("expand") != "description" {
			t.Errorf("unexpected expand: %s", r.URL.Query().Get("expand"))
		}
		if r.URL.Query().Get("recent") != "5" {
			t.Errorf("unexpected recent: %s", r.URL.Query().Get("recent"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Project{
			{Key: "PROJ1", Name: "Project 1"},
			{Key: "PROJ2", Name: "Project 2"},
		})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	projects, err := c.GetAllProjects("description", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(projects) != 2 {
		t.Errorf("expected 2 projects, got %d", len(projects))
	}
	if projects[0].Key != "PROJ1" {
		t.Errorf("expected PROJ1, got %s", projects[0].Key)
	}
}

func TestGetAllProjectsNoParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Errorf("expected no query params, got: %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Project{})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	_, err := c.GetAllProjects("", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateProject(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Project{Key: "NEW", Name: "New Project"})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	result, err := c.CreateProject(map[string]interface{}{
		"key":  "NEW",
		"name": "New Project",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Key != "NEW" {
		t.Errorf("expected key NEW, got %s", result.Key)
	}
}

func TestGetProject(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/3/project/PROJ1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Project{Key: "PROJ1", Name: "Project 1"})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	result, err := c.GetProject("PROJ1", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Key != "PROJ1" {
		t.Errorf("expected PROJ1, got %s", result.Key)
	}
}

func TestUpdateProject(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Project{Key: "PROJ1", Name: "Updated"})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	result, err := c.UpdateProject("PROJ1", map[string]interface{}{"name": "Updated"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "Updated" {
		t.Errorf("expected Updated, got %s", result.Name)
	}
}

func TestDeleteProject(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/rest/api/3/project/PROJ1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newTestClient(srv)
	err := c.DeleteProject("PROJ1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSearchProjects(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/3/project/search" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("query") != "test" {
			t.Errorf("unexpected query: %s", r.URL.Query().Get("query"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PageBean[Project]{
			Pagination: Pagination{Total: 1},
			Values:     []Project{{Key: "TST", Name: "Test"}},
		})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	result, err := c.SearchProjects("test", 0, 25, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Values) != 1 {
		t.Errorf("expected 1 project, got %d", len(result.Values))
	}
}

func TestGetRecentProjects(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/3/project/recent" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Project{{Key: "REC"}})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	projects, err := c.GetRecentProjects(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(projects) != 1 {
		t.Errorf("expected 1 project, got %d", len(projects))
	}
}
