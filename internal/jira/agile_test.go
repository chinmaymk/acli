package jira

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBoards(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/agile/1.0/board" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("projectKeyOrId") != "PROJ" {
			t.Errorf("unexpected project: %s", r.URL.Query().Get("projectKeyOrId"))
		}
		if r.URL.Query().Get("type") != "scrum" {
			t.Errorf("unexpected type: %s", r.URL.Query().Get("type"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(BoardList{
			Total:  1,
			Values: []Board{{ID: 1, Name: "Sprint Board"}},
		})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	result, err := c.GetBoards(0, 50, "PROJ", "scrum", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Values) != 1 {
		t.Errorf("expected 1 board, got %d", len(result.Values))
	}
	if result.Values[0].Name != "Sprint Board" {
		t.Errorf("unexpected board name: %s", result.Values[0].Name)
	}
}

func TestGetBoard(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/agile/1.0/board/42" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Board{ID: 42, Name: "My Board"})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	result, err := c.GetBoard(42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != 42 {
		t.Errorf("expected board ID 42, got %d", result.ID)
	}
}

func TestGetBoardConfiguration(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/agile/1.0/board/42/configuration" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(BoardConfiguration{ID: 42, Name: "Config"})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	result, err := c.GetBoardConfiguration(42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != 42 {
		t.Errorf("expected config ID 42, got %d", result.ID)
	}
}

func TestGetBoardIssues(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/agile/1.0/board/1/issue" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("jql") != "status=Open" {
			t.Errorf("unexpected jql: %s", r.URL.Query().Get("jql"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SprintIssuesResponse{
			Total:  1,
			Issues: []IssueDetailed{{Key: "TEST-1"}},
		})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	result, err := c.GetBoardIssues(1, 0, 50, "status=Open")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Issues) != 1 {
		t.Errorf("expected 1 issue, got %d", len(result.Issues))
	}
}

func TestGetBoardBacklog(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/agile/1.0/board/1/backlog" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SprintIssuesResponse{})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	_, err := c.GetBoardBacklog(1, 0, 50, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetBoardSprints(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/agile/1.0/board/1/sprint" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("state") != "active" {
			t.Errorf("unexpected state: %s", r.URL.Query().Get("state"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SprintList{
			Total:  1,
			Values: []Sprint{{ID: 10, Name: "Sprint 1", State: "active"}},
		})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	result, err := c.GetBoardSprints(1, 0, 50, "active")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Values) != 1 {
		t.Errorf("expected 1 sprint, got %d", len(result.Values))
	}
}

func TestGetBoardEpics(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/agile/1.0/board/1/epic" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(EpicList{})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	_, err := c.GetBoardEpics(1, 0, 50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetSprint(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/agile/1.0/sprint/10" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Sprint{ID: 10, Name: "Sprint 1", State: "active"})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	result, err := c.GetSprint(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != 10 {
		t.Errorf("expected sprint ID 10, got %d", result.ID)
	}
	if result.State != "active" {
		t.Errorf("expected state active, got %s", result.State)
	}
}

func TestCreateSprint(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/rest/agile/1.0/sprint" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Sprint{ID: 20, Name: "New Sprint"})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	result, err := c.CreateSprint(map[string]interface{}{
		"name":            "New Sprint",
		"originBoardId":   1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "New Sprint" {
		t.Errorf("expected New Sprint, got %s", result.Name)
	}
}
