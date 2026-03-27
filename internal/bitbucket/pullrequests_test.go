package bitbucket

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListPullRequests(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PaginatedResponse{
			Values: json.RawMessage(`[{"id": 1, "title": "My PR", "state": "OPEN"}]`),
		})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	prs, err := c.ListPullRequests("ws", "repo", &ListPRsOptions{State: "OPEN"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(prs) != 1 {
		t.Fatalf("expected 1 PR, got %d", len(prs))
	}
	if prs[0].Title != "My PR" {
		t.Errorf("unexpected title: %s", prs[0].Title)
	}
}

func TestListPullRequestsAll(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		if callCount == 1 {
			json.NewEncoder(w).Encode(PaginatedResponse{
				Next:   "https://api.bitbucket.org/2.0/repositories/ws/repo/pullrequests?page=2",
				Values: json.RawMessage(`[{"id": 1, "title": "PR 1"}]`),
			})
		} else {
			json.NewEncoder(w).Encode(PaginatedResponse{
				Values: json.RawMessage(`[{"id": 2, "title": "PR 2"}]`),
			})
		}
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	prs, err := c.ListPullRequests("ws", "repo", &ListPRsOptions{All: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(prs) != 2 {
		t.Errorf("expected 2 PRs, got %d", len(prs))
	}
}

func TestGetPullRequest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PullRequest{ID: 42, Title: "Fix bug", State: "OPEN"})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	pr, err := c.GetPullRequest("ws", "repo", 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pr.ID != 42 {
		t.Errorf("expected ID 42, got %d", pr.ID)
	}
	if pr.Title != "Fix bug" {
		t.Errorf("unexpected title: %s", pr.Title)
	}
}

func TestCreatePullRequest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["title"] != "New feature" {
			t.Errorf("unexpected title: %v", body["title"])
		}
		source := body["source"].(map[string]interface{})
		branch := source["branch"].(map[string]interface{})
		if branch["name"] != "feature/test" {
			t.Errorf("unexpected source branch: %v", branch["name"])
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(PullRequest{ID: 1, Title: "New feature"})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	pr, err := c.CreatePullRequest("ws", "repo", &CreatePRRequest{
		Title:             "New feature",
		SourceBranch:      "feature/test",
		DestinationBranch: "main",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pr.ID != 1 {
		t.Errorf("expected ID 1, got %d", pr.ID)
	}
}

func TestCreatePullRequestNoDestination(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["destination"] != nil {
			t.Error("expected nil destination when not specified")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PullRequest{ID: 2})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	_, err := c.CreatePullRequest("ws", "repo", &CreatePRRequest{
		Title:        "Auto dest",
		SourceBranch: "feature/x",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdatePullRequest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PullRequest{ID: 1, Title: "Updated title"})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	pr, err := c.UpdatePullRequest("ws", "repo", 1, &UpdatePRRequest{Title: "Updated title"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pr.Title != "Updated title" {
		t.Errorf("unexpected title: %s", pr.Title)
	}
}

func TestApprovePullRequest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Participant{Approved: true, Role: "REVIEWER"})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	p, err := c.ApprovePullRequest("ws", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !p.Approved {
		t.Error("expected approved=true")
	}
}

func TestUnapprovePullRequest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	err := c.UnapprovePullRequest("ws", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeclinePullRequest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PullRequest{ID: 1, State: "DECLINED"})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	pr, err := c.DeclinePullRequest("ws", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pr.State != "DECLINED" {
		t.Errorf("expected DECLINED, got %s", pr.State)
	}
}

func TestMergePullRequest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PullRequest{ID: 1, State: "MERGED"})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	pr, err := c.MergePullRequest("ws", "repo", 1, &MergePRRequest{MergeStrategy: "squash"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pr.State != "MERGED" {
		t.Errorf("expected MERGED, got %s", pr.State)
	}
}

func TestMergePullRequestNilRequest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PullRequest{ID: 1, State: "MERGED"})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	_, err := c.MergePullRequest("ws", "repo", 1, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRequestChangesPullRequest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Participant{State: "changes_requested"})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	p, err := c.RequestChangesPullRequest("ws", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.State != "changes_requested" {
		t.Errorf("unexpected state: %s", p.State)
	}
}

func TestRemoveRequestChanges(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	err := c.RemoveRequestChangesPullRequest("ws", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListPRComments(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PaginatedResponse{
			Values: json.RawMessage(`[{"id": 100, "content": {"raw": "LGTM"}}]`),
		})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	comments, err := c.ListPRComments("ws", "repo", 1, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(comments))
	}
	if comments[0].ID != 100 {
		t.Errorf("unexpected comment ID: %d", comments[0].ID)
	}
}

func TestCreatePRComment(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PRComment{ID: 101})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	comment, err := c.CreatePRComment("ws", "repo", 1, "Great work!")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.ID != 101 {
		t.Errorf("unexpected comment ID: %d", comment.ID)
	}
}

func TestCreatePRCommentInline(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if _, ok := body["inline"]; !ok {
			t.Error("expected inline params in body")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PRComment{ID: 102})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	comment, err := c.CreatePRCommentInline("ws", "repo", 1, "Fix this", &InlineCommentParams{Path: "main.go", To: 42})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.ID != 102 {
		t.Errorf("unexpected comment ID: %d", comment.ID)
	}
}

func TestGetPRDiff(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("diff --git a/file.go b/file.go"))
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	diff, err := c.GetPRDiff("ws", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if diff == "" {
		t.Error("expected non-empty diff")
	}
}

func TestListPRTasks(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PaginatedResponse{
			Values: json.RawMessage(`[{"id": 1, "state": "OPEN", "content": {"raw": "Fix tests"}}]`),
		})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	tasks, err := c.ListPRTasks("ws", "repo", 1, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}
}

func TestCreatePRTask(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PRTask{ID: 10, State: "OPEN"})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	task, err := c.CreatePRTask("ws", "repo", 1, &CreatePRTaskRequest{Content: "Do something"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task.ID != 10 {
		t.Errorf("unexpected task ID: %d", task.ID)
	}
}

func TestCreatePRTaskWithComment(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if _, ok := body["comment"]; !ok {
			t.Error("expected comment in body")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PRTask{ID: 11})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	commentID := 42
	_, err := c.CreatePRTask("ws", "repo", 1, &CreatePRTaskRequest{Content: "Fix this", CommentID: &commentID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdatePRTask(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PRTask{ID: 10, State: "RESOLVED"})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	task, err := c.UpdatePRTask("ws", "repo", 1, 10, &UpdatePRTaskRequest{State: "RESOLVED"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task.State != "RESOLVED" {
		t.Errorf("unexpected state: %s", task.State)
	}
}

func TestDeletePRTask(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	err := c.DeletePRTask("ws", "repo", 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
