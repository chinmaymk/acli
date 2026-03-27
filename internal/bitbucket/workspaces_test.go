package bitbucket

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListWorkspaces(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PaginatedResponse{
			Size:    2,
			PageLen: 50,
			Values:  json.RawMessage(`[{"slug": "ws1", "name": "Workspace 1"}, {"slug": "ws2", "name": "Workspace 2"}]`),
		})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	workspaces, err := c.ListWorkspaces(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(workspaces) != 2 {
		t.Fatalf("expected 2 workspaces, got %d", len(workspaces))
	}
	if workspaces[0].Slug != "ws1" {
		t.Errorf("unexpected slug: %s", workspaces[0].Slug)
	}
}

func TestListWorkspacesWithPagination(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") != "2" {
			t.Errorf("expected page=2, got: %s", r.URL.Query().Get("page"))
		}
		if r.URL.Query().Get("pagelen") != "10" {
			t.Errorf("expected pagelen=10, got: %s", r.URL.Query().Get("pagelen"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PaginatedResponse{Values: json.RawMessage(`[]`)})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	_, err := c.ListWorkspaces(&PaginationOptions{Page: 2, PageLen: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListWorkspacesAll(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		if callCount == 1 {
			json.NewEncoder(w).Encode(PaginatedResponse{
				Next:   "https://api.bitbucket.org/2.0/workspaces?page=2&pagelen=50",
				Values: json.RawMessage(`[{"slug": "ws1"}]`),
			})
		} else {
			json.NewEncoder(w).Encode(PaginatedResponse{
				Values: json.RawMessage(`[{"slug": "ws2"}]`),
			})
		}
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	workspaces, err := c.ListWorkspaces(&PaginationOptions{All: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(workspaces) != 2 {
		t.Errorf("expected 2 workspaces, got %d", len(workspaces))
	}
}

func TestGetWorkspace(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Workspace{Slug: "myws", Name: "My Workspace", UUID: "{uuid-123}"})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	ws, err := c.GetWorkspace("myws")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ws.Slug != "myws" {
		t.Errorf("unexpected slug: %s", ws.Slug)
	}
}

func TestListWorkspaceMembers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PaginatedResponse{
			Values: json.RawMessage(`[{"user": {"display_name": "John", "account_id": "123"}}]`),
		})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	members, err := c.ListWorkspaceMembers("myws", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(members) != 1 {
		t.Fatalf("expected 1 member, got %d", len(members))
	}
	if members[0].User.DisplayName != "John" {
		t.Errorf("unexpected display name: %s", members[0].User.DisplayName)
	}
}

func TestListWorkspacePermissions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PaginatedResponse{
			Values: json.RawMessage(`[{"permission": "owner", "user": {"display_name": "Admin"}}]`),
		})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	perms, err := c.ListWorkspacePermissions("myws", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(perms) != 1 {
		t.Fatalf("expected 1 permission, got %d", len(perms))
	}
	if perms[0].Permission != "owner" {
		t.Errorf("unexpected permission: %s", perms[0].Permission)
	}
}

func TestListWorkspaceMembersAll(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		if callCount == 1 {
			json.NewEncoder(w).Encode(PaginatedResponse{
				Next:   "https://api.bitbucket.org/2.0/workspaces/ws/members?page=2",
				Values: json.RawMessage(`[{"user": {"display_name": "A"}}]`),
			})
		} else {
			json.NewEncoder(w).Encode(PaginatedResponse{
				Values: json.RawMessage(`[{"user": {"display_name": "B"}}]`),
			})
		}
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	members, err := c.ListWorkspaceMembers("ws", &PaginationOptions{All: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(members) != 2 {
		t.Errorf("expected 2 members, got %d", len(members))
	}
}
