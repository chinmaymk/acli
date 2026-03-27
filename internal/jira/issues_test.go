package jira

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestClient(srv *httptest.Server) *Client {
	return &Client{
		BaseURL:    srv.URL,
		APIToken:   "testtoken",
		HTTPClient: srv.Client(),
	}
}

func TestCreateIssue(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/rest/api/3/issue" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(CreatedIssue{ID: "10001", Key: "TEST-1", Self: "https://test.atlassian.net/rest/api/3/issue/10001"})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	result, err := c.CreateIssue(&IssueUpdateDetails{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Key != "TEST-1" {
		t.Errorf("expected key TEST-1, got %s", result.Key)
	}
	if result.ID != "10001" {
		t.Errorf("expected id 10001, got %s", result.ID)
	}
}

func TestGetIssue(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/3/issue/TEST-1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		// Verify query params
		if r.URL.Query().Get("fields") != "summary,status" {
			t.Errorf("unexpected fields: %s", r.URL.Query().Get("fields"))
		}
		if r.URL.Query().Get("expand") != "changelog" {
			t.Errorf("unexpected expand: %s", r.URL.Query().Get("expand"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(IssueDetailed{
			Key: "TEST-1",
			Fields: IssueFields{
				Summary: "Test issue",
			},
		})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	result, err := c.GetIssue("TEST-1", []string{"summary", "status"}, []string{"changelog"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Key != "TEST-1" {
		t.Errorf("expected key TEST-1, got %s", result.Key)
	}
	if result.Fields.Summary != "Test issue" {
		t.Errorf("unexpected summary: %s", result.Fields.Summary)
	}
}

func TestGetIssueNoFieldsOrExpand(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Errorf("expected no query params, got: %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(IssueDetailed{Key: "TEST-2"})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	result, err := c.GetIssue("TEST-2", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Key != "TEST-2" {
		t.Errorf("expected key TEST-2, got %s", result.Key)
	}
}

func TestEditIssue(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if !strings.HasPrefix(r.URL.Path, "/rest/api/3/issue/TEST-1") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newTestClient(srv)
	err := c.EditIssue("TEST-1", &IssueUpdateDetails{}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEditIssueNoNotify(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("notifyUsers") != "false" {
			t.Errorf("expected notifyUsers=false, got: %s", r.URL.Query().Get("notifyUsers"))
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newTestClient(srv)
	err := c.EditIssue("TEST-1", &IssueUpdateDetails{}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteIssue(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/rest/api/3/issue/TEST-1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newTestClient(srv)
	err := c.DeleteIssue("TEST-1", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteIssueWithSubtasks(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("deleteSubtasks") != "true" {
			t.Error("expected deleteSubtasks=true")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newTestClient(srv)
	err := c.DeleteIssue("TEST-1", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAssignIssue(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/rest/api/3/issue/TEST-1/assignee" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["accountId"] != "user123" {
			t.Errorf("unexpected accountId: %v", body["accountId"])
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newTestClient(srv)
	err := c.AssignIssue("TEST-1", "user123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetIssueTransitions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/3/issue/TEST-1/transitions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(TransitionsResponse{
			Transitions: []IssueTransition{
				{ID: "11", Name: "In Progress"},
				{ID: "21", Name: "Done"},
			},
		})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	result, err := c.GetIssueTransitions("TEST-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Transitions) != 2 {
		t.Errorf("expected 2 transitions, got %d", len(result.Transitions))
	}
}

func TestGetIssueComments(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/3/issue/TEST-1/comment" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("startAt") != "0" {
			t.Errorf("unexpected startAt: %s", r.URL.Query().Get("startAt"))
		}
		if r.URL.Query().Get("maxResults") != "50" {
			t.Errorf("unexpected maxResults: %s", r.URL.Query().Get("maxResults"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CommentPage{})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	_, err := c.GetIssueComments("TEST-1", 0, 50)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAddIssueComment(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Comment{ID: "1001"})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	result, err := c.AddIssueComment("TEST-1", "A comment", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "1001" {
		t.Errorf("unexpected comment ID: %s", result.ID)
	}
}

func TestAddIssueCommentWithVisibility(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if _, ok := body["visibility"]; !ok {
			t.Error("expected visibility in body")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Comment{ID: "1002"})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	vis := &Visibility{Type: "role", Value: "Administrators"}
	_, err := c.AddIssueComment("TEST-1", "A private comment", vis)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteIssueComment(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/rest/api/3/issue/TEST-1/comment/1001" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newTestClient(srv)
	err := c.DeleteIssueComment("TEST-1", "1001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetIssueVotes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/3/issue/TEST-1/votes" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Votes{Votes: 5, HasVoted: true})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	result, err := c.GetIssueVotes("TEST-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Votes != 5 {
		t.Errorf("expected 5 votes, got %d", result.Votes)
	}
}

func TestAddAndRemoveIssueVote(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newTestClient(srv)
	if err := c.AddIssueVote("TEST-1"); err != nil {
		t.Fatalf("AddIssueVote failed: %v", err)
	}
	if err := c.RemoveIssueVote("TEST-1"); err != nil {
		t.Fatalf("RemoveIssueVote failed: %v", err)
	}
}

func TestGetIssueWatchers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Watches{WatchCount: 3, IsWatching: true})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	result, err := c.GetIssueWatchers("TEST-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.WatchCount != 3 {
		t.Errorf("expected 3 watchers, got %d", result.WatchCount)
	}
}

func TestRemoveIssueWatcher(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Query().Get("accountId") != "user456" {
			t.Errorf("unexpected accountId: %s", r.URL.Query().Get("accountId"))
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newTestClient(srv)
	err := c.RemoveIssueWatcher("TEST-1", "user456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetIssueWorklogs(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/3/issue/TEST-1/worklog" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(WorklogPage{})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	_, err := c.GetIssueWorklogs("TEST-1", 0, 25)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAddIssueWorklog(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Worklog{ID: "500", TimeSpent: "2h"})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	result, err := c.AddIssueWorklog("TEST-1", &Worklog{TimeSpent: "2h"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TimeSpent != "2h" {
		t.Errorf("unexpected time spent: %s", result.TimeSpent)
	}
}

func TestDeleteIssueWorklog(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/3/issue/TEST-1/worklog/500" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newTestClient(srv)
	err := c.DeleteIssueWorklog("TEST-1", "500")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetIssueRemoteLinks(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]RemoteIssueLink{{ID: 1}})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	result, err := c.GetIssueRemoteLinks("TEST-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 remote link, got %d", len(result))
	}
}

func TestGetIssueProperties(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"keys": []map[string]string{{"key": "prop1"}, {"key": "prop2"}},
		})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	result, err := c.GetIssueProperties("TEST-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 properties, got %d", len(result))
	}
}

func TestSetIssueProperty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/rest/api/3/issue/TEST-1/properties/myprop" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newTestClient(srv)
	err := c.SetIssueProperty("TEST-1", "myprop", map[string]string{"value": "data"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteIssueProperty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newTestClient(srv)
	err := c.DeleteIssueProperty("TEST-1", "myprop")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetCreateMeta(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("projectKeys") != "PROJ1,PROJ2" {
			t.Errorf("unexpected projectKeys: %s", r.URL.Query().Get("projectKeys"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CreateMeta{})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	_, err := c.GetCreateMeta([]string{"PROJ1", "PROJ2"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBulkFetchIssues(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/rest/api/3/issue/bulkfetch" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SearchResults{})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	_, err := c.BulkFetchIssues([]int{1, 2, 3}, []string{"summary"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetIssueChangelog(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/3/issue/TEST-1/changelog" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ChangelogPage{})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	_, err := c.GetIssueChangelog("TEST-1", 0, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDoIssueTransition(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/rest/api/3/issue/TEST-1/transitions" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newTestClient(srv)
	err := c.DoIssueTransition("TEST-1", &IssueUpdateDetails{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetIssueEditMeta(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/3/issue/TEST-1/editmeta" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"fields": {"summary": {"type": "string"}}}`))
	}))
	defer srv.Close()

	c := newTestClient(srv)
	result, err := c.GetIssueEditMeta("TEST-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Error("expected non-nil result")
	}
}

func TestCreateAndDeleteRemoteLink(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(RemoteIssueLink{ID: 42})
		case "DELETE":
			w.WriteHeader(http.StatusNoContent)
		}
	}))
	defer srv.Close()

	c := newTestClient(srv)

	link, err := c.CreateIssueRemoteLink("TEST-1", &RemoteIssueLink{})
	if err != nil {
		t.Fatalf("create remote link failed: %v", err)
	}
	if link.ID != 42 {
		t.Errorf("expected ID 42, got %d", link.ID)
	}

	err = c.DeleteIssueRemoteLink("TEST-1", "42")
	if err != nil {
		t.Fatalf("delete remote link failed: %v", err)
	}
}

func TestIssueAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errorMessages": []string{"Issue does not exist"},
		})
	}))
	defer srv.Close()

	c := newTestClient(srv)
	_, err := c.GetIssue("NOPE-999", nil, nil)
	if err == nil {
		t.Fatal("expected error for nonexistent issue")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("expected 404, got %d", apiErr.StatusCode)
	}
}
