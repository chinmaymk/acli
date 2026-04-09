package bitbucket

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestListPRCommits(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PaginatedResponse{
			Values: json.RawMessage(`[{"hash": "abc123def456", "message": "Fix bug", "date": "2026-01-01T00:00:00+00:00"}]`),
		})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	commits, err := c.ListPRCommits("ws", "repo", 1, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(commits) != 1 {
		t.Fatalf("expected 1 commit, got %d", len(commits))
	}
	if commits[0].Hash != "abc123def456" {
		t.Errorf("unexpected hash: %s", commits[0].Hash)
	}
	if commits[0].Message != "Fix bug" {
		t.Errorf("unexpected message: %s", commits[0].Message)
	}
}

func TestListPRCommitsAll(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		if callCount == 1 {
			json.NewEncoder(w).Encode(PaginatedResponse{
				Next:   "https://api.bitbucket.org/2.0/repositories/ws/repo/pullrequests/1/commits?page=2",
				Values: json.RawMessage(`[{"hash": "aaa111"}]`),
			})
		} else {
			json.NewEncoder(w).Encode(PaginatedResponse{
				Values: json.RawMessage(`[{"hash": "bbb222"}]`),
			})
		}
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	commits, err := c.ListPRCommits("ws", "repo", 1, &PaginationOptions{All: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(commits) != 2 {
		t.Errorf("expected 2 commits, got %d", len(commits))
	}
}

func TestGetPRDiffStat(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PaginatedResponse{
			Values: json.RawMessage(`[{"type": "diffstat", "status": "modified", "old": {"type": "commit_file", "path": "main.go"}, "new": {"type": "commit_file", "path": "main.go"}, "lines_added": 10, "lines_removed": 3}]`),
		})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	stats, err := c.GetPRDiffStat("ws", "repo", 1, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stats) != 1 {
		t.Fatalf("expected 1 diffstat, got %d", len(stats))
	}
	if stats[0].Type != "diffstat" {
		t.Errorf("unexpected type: %s", stats[0].Type)
	}
	if stats[0].Status != "modified" {
		t.Errorf("unexpected status: %s", stats[0].Status)
	}
	if stats[0].LinesAdded != 10 {
		t.Errorf("expected 10 lines added, got %d", stats[0].LinesAdded)
	}
	if stats[0].LinesRemoved != 3 {
		t.Errorf("expected 3 lines removed, got %d", stats[0].LinesRemoved)
	}
	if stats[0].Old.Type != "commit_file" {
		t.Errorf("unexpected old type: %s", stats[0].Old.Type)
	}
}

func TestGetPRDiffStatAdded(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PaginatedResponse{
			Values: json.RawMessage(`[{"type": "diffstat", "status": "added", "old": null, "new": {"type": "commit_file", "path": "new_file.go"}, "lines_added": 50, "lines_removed": 0}]`),
		})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	stats, err := c.GetPRDiffStat("ws", "repo", 1, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stats) != 1 {
		t.Fatalf("expected 1 diffstat, got %d", len(stats))
	}
	if stats[0].Status != "added" {
		t.Errorf("unexpected status: %s", stats[0].Status)
	}
	if stats[0].Old != nil {
		t.Error("expected nil old for added file")
	}
	if stats[0].New.Path != "new_file.go" {
		t.Errorf("unexpected new path: %s", stats[0].New.Path)
	}
}

func TestGetPRDiffStatRenamed(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PaginatedResponse{
			Values: json.RawMessage(`[{"type": "diffstat", "status": "renamed", "old": {"type": "commit_file", "path": "old_name.go"}, "new": {"type": "commit_file", "path": "new_name.go"}, "lines_added": 0, "lines_removed": 0}]`),
		})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	stats, err := c.GetPRDiffStat("ws", "repo", 1, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stats) != 1 {
		t.Fatalf("expected 1 diffstat, got %d", len(stats))
	}
	if stats[0].Status != "renamed" {
		t.Errorf("unexpected status: %s", stats[0].Status)
	}
	if stats[0].Old.Path != "old_name.go" {
		t.Errorf("unexpected old path: %s", stats[0].Old.Path)
	}
	if stats[0].New.Path != "new_name.go" {
		t.Errorf("unexpected new path: %s", stats[0].New.Path)
	}
}

func TestListPRActivity(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PaginatedResponse{
			Values: json.RawMessage(`[
				{
					"approval": {
						"date": "2026-01-01T00:00:00+00:00",
						"user": {"display_name": "Alice", "uuid": "{aaa}", "type": "user", "nickname": "alice", "account_id": "123"}
					},
					"pull_request": {"type": "pullrequest", "id": 1, "title": "Test PR"}
				},
				{
					"update": {
						"state": "OPEN",
						"date": "2026-01-01T00:00:00+00:00",
						"title": "Test PR",
						"author": {"display_name": "Bob", "uuid": "{bbb}", "type": "user", "nickname": "bob"},
						"source": {"branch": {"name": "feature"}, "commit": {"hash": "abc123", "type": "commit"}},
						"destination": {"branch": {"name": "main"}, "commit": {"hash": "def456", "type": "commit"}}
					},
					"pull_request": {"type": "pullrequest", "id": 1, "title": "Test PR"}
				},
				{
					"comment": {
						"id": 42,
						"content": {"raw": "Looks good", "markup": "markdown"},
						"user": {"display_name": "Carol", "uuid": "{ccc}"},
						"created_on": "2026-01-01T00:00:00+00:00",
						"updated_on": "2026-01-01T00:00:00+00:00"
					},
					"pull_request": {"type": "pullrequest", "id": 1, "title": "Test PR"}
				}
			]`),
		})
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	activities, err := c.ListPRActivity("ws", "repo", 1, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(activities) != 3 {
		t.Fatalf("expected 3 activities, got %d", len(activities))
	}

	// Approval
	if activities[0].Approval == nil {
		t.Fatal("expected approval activity")
	}
	if activities[0].Approval.User.DisplayName != "Alice" {
		t.Errorf("unexpected approval user: %s", activities[0].Approval.User.DisplayName)
	}
	if activities[0].Approval.User.Nickname != "alice" {
		t.Errorf("unexpected approval nickname: %s", activities[0].Approval.User.Nickname)
	}
	if activities[0].PullRequest == nil || activities[0].PullRequest.ID != 1 {
		t.Error("expected pull_request reference in activity")
	}

	// Update
	if activities[1].Update == nil {
		t.Fatal("expected update activity")
	}
	if activities[1].Update.State != "OPEN" {
		t.Errorf("unexpected update state: %s", activities[1].Update.State)
	}
	if activities[1].Update.Source.Branch.Name != "feature" {
		t.Errorf("unexpected source branch: %s", activities[1].Update.Source.Branch.Name)
	}
	if activities[1].Update.Source.Commit.Hash != "abc123" {
		t.Errorf("unexpected source commit: %s", activities[1].Update.Source.Commit.Hash)
	}
	if activities[1].Update.Destination.Branch.Name != "main" {
		t.Errorf("unexpected destination branch: %s", activities[1].Update.Destination.Branch.Name)
	}

	// Comment
	if activities[2].Comment == nil {
		t.Fatal("expected comment activity")
	}
	if activities[2].Comment.ID != 42 {
		t.Errorf("unexpected comment ID: %d", activities[2].Comment.ID)
	}
	if activities[2].Comment.Content.Raw != "Looks good" {
		t.Errorf("unexpected comment content: %s", activities[2].Comment.Content.Raw)
	}
}

// makeInt returns a pointer to an int literal (test helper).
func makeInt(n int) *int { return &n }

// makeInline returns an Inline value suitable for a PRComment.
func makeInline(path string, toLine int) *struct {
	Path      string `json:"path"`
	From      *int   `json:"from"`
	To        *int   `json:"to"`
	StartFrom *int   `json:"start_from,omitempty"`
	StartTo   *int   `json:"start_to,omitempty"`
} {
	return &struct {
		Path      string `json:"path"`
		From      *int   `json:"from"`
		To        *int   `json:"to"`
		StartFrom *int   `json:"start_from,omitempty"`
		StartTo   *int   `json:"start_to,omitempty"`
	}{
		Path: path,
		To:   makeInt(toLine),
	}
}

func TestBuildPRReview_FileThreadWithRepliesAndTask(t *testing.T) {
	rootComment := PRComment{ID: 100}
	rootComment.Content.Raw = "This function should handle nil input"
	rootComment.User.DisplayName = "Alice"
	rootComment.Inline = makeInline("pkg/foo.go", 42)

	reply := PRComment{ID: 101}
	reply.Content.Raw = "Good point, will fix"
	reply.User.DisplayName = "Bob"
	reply.Parent = &struct {
		ID int `json:"id"`
	}{ID: 100}

	nestedReply := PRComment{ID: 102}
	nestedReply.Content.Raw = "Thanks"
	nestedReply.User.DisplayName = "Alice"
	nestedReply.Parent = &struct {
		ID int `json:"id"`
	}{ID: 101}

	task := PRTask{ID: 5, State: "UNRESOLVED"}
	task.Content.Raw = "Add nil guard"
	task.Creator.DisplayName = "Alice"
	task.Comment = &PRComment{ID: 100}

	review := BuildPRReview(7, []PRComment{rootComment, reply, nestedReply}, []PRTask{task})

	if review.PullRequestID != 7 {
		t.Errorf("expected PR ID 7, got %d", review.PullRequestID)
	}
	if len(review.FileThreads) != 1 {
		t.Fatalf("expected 1 file thread, got %d", len(review.FileThreads))
	}
	if len(review.GeneralThreads) != 0 {
		t.Errorf("expected 0 general threads, got %d", len(review.GeneralThreads))
	}
	thread := review.FileThreads[0]
	if thread.File != "pkg/foo.go" {
		t.Errorf("expected file 'pkg/foo.go', got %q", thread.File)
	}
	if thread.Line != 42 {
		t.Errorf("expected line 42, got %d", thread.Line)
	}
	if thread.Comment.ID != 100 {
		t.Errorf("expected root comment ID 100, got %d", thread.Comment.ID)
	}
	if len(thread.Replies) != 2 {
		t.Fatalf("expected 2 replies (direct + nested), got %d", len(thread.Replies))
	}
	if thread.Replies[0].ID != 101 || thread.Replies[1].ID != 102 {
		t.Errorf("unexpected reply order: %+v", thread.Replies)
	}
	if len(thread.Tasks) != 1 {
		t.Fatalf("expected 1 task on thread, got %d", len(thread.Tasks))
	}
	if thread.Tasks[0].ID != 5 {
		t.Errorf("expected task ID 5, got %d", thread.Tasks[0].ID)
	}
	if len(review.StandaloneTasks) != 0 {
		t.Errorf("expected no standalone tasks, got %d", len(review.StandaloneTasks))
	}
	if review.Counts.Comments != 3 {
		t.Errorf("expected 3 comments, got %d", review.Counts.Comments)
	}
	if review.Counts.Threads != 1 {
		t.Errorf("expected 1 thread, got %d", review.Counts.Threads)
	}
	if review.Counts.Tasks != 1 {
		t.Errorf("expected 1 task, got %d", review.Counts.Tasks)
	}
	if review.Counts.UnresolvedTasks != 1 {
		t.Errorf("expected 1 unresolved task, got %d", review.Counts.UnresolvedTasks)
	}
	if review.Counts.UnresolvedThreads != 1 {
		t.Errorf("expected 1 unresolved thread, got %d", review.Counts.UnresolvedThreads)
	}
}

func TestBuildPRReview_TaskLinkedToReply(t *testing.T) {
	// A task attached to a reply comment should still end up on the root thread.
	root := PRComment{ID: 10}
	root.Content.Raw = "Please rename this variable"
	root.Inline = makeInline("main.go", 5)

	reply := PRComment{ID: 11}
	reply.Content.Raw = "Which name?"
	reply.Parent = &struct {
		ID int `json:"id"`
	}{ID: 10}

	task := PRTask{ID: 99, State: "UNRESOLVED"}
	task.Content.Raw = "Pick a descriptive name"
	task.Comment = &PRComment{ID: 11} // attached to reply

	review := BuildPRReview(1, []PRComment{root, reply}, []PRTask{task})

	if len(review.FileThreads) != 1 {
		t.Fatalf("expected 1 file thread, got %d", len(review.FileThreads))
	}
	if len(review.FileThreads[0].Tasks) != 1 {
		t.Errorf("expected task on thread (via reply), got %d", len(review.FileThreads[0].Tasks))
	}
	if len(review.StandaloneTasks) != 0 {
		t.Errorf("expected no standalone tasks, got %d", len(review.StandaloneTasks))
	}
}

func TestBuildPRReview_GeneralAndStandalone(t *testing.T) {
	generalRoot := PRComment{ID: 1}
	generalRoot.Content.Raw = "LGTM overall"

	standaloneTask := PRTask{ID: 77, State: "RESOLVED"}
	standaloneTask.Content.Raw = "Update CHANGELOG"
	// no Comment field -> standalone

	orphanTask := PRTask{ID: 78, State: "UNRESOLVED"}
	orphanTask.Content.Raw = "Orphan"
	orphanTask.Comment = &PRComment{ID: 9999} // not in comment list -> standalone

	review := BuildPRReview(42, []PRComment{generalRoot}, []PRTask{standaloneTask, orphanTask})

	if len(review.GeneralThreads) != 1 {
		t.Fatalf("expected 1 general thread, got %d", len(review.GeneralThreads))
	}
	if len(review.FileThreads) != 0 {
		t.Errorf("expected 0 file threads, got %d", len(review.FileThreads))
	}
	if len(review.StandaloneTasks) != 2 {
		t.Fatalf("expected 2 standalone tasks, got %d", len(review.StandaloneTasks))
	}
	if review.Counts.Tasks != 2 {
		t.Errorf("expected 2 tasks, got %d", review.Counts.Tasks)
	}
	if review.Counts.UnresolvedTasks != 1 {
		t.Errorf("expected 1 unresolved task (orphan), got %d", review.Counts.UnresolvedTasks)
	}
}

func TestBuildPRReview_SkipsDeletedComments(t *testing.T) {
	deleted := PRComment{ID: 1, Deleted: true}
	kept := PRComment{ID: 2}
	kept.Content.Raw = "Real comment"

	review := BuildPRReview(1, []PRComment{deleted, kept}, nil)

	if len(review.GeneralThreads) != 1 {
		t.Fatalf("expected 1 thread, got %d", len(review.GeneralThreads))
	}
	if review.GeneralThreads[0].Comment.ID != 2 {
		t.Errorf("expected kept comment, got %d", review.GeneralThreads[0].Comment.ID)
	}
	if review.Counts.Comments != 1 {
		t.Errorf("expected 1 comment (deleted excluded), got %d", review.Counts.Comments)
	}
}

func TestBuildPRReview_TaskWithInlineComment(t *testing.T) {
	// Validates that when the Bitbucket API returns a task with the full comment
	// (which is what it does per the pullrequest_comment_task schema), we can
	// surface the file/line context without a second call. Note: BuildPRReview
	// also needs the comment present in the comments list to attach the task
	// to a thread, so we seed it.
	c := PRComment{ID: 50}
	c.Content.Raw = "Nit"
	c.Inline = makeInline("x.go", 12)

	task := PRTask{ID: 1}
	task.Content.Raw = "Rename foo to bar"
	task.Comment = &PRComment{ID: 50, Inline: makeInline("x.go", 12)}

	review := BuildPRReview(1, []PRComment{c}, []PRTask{task})
	if len(review.FileThreads) != 1 {
		t.Fatalf("expected 1 file thread, got %d", len(review.FileThreads))
	}
	if review.FileThreads[0].File != "x.go" || review.FileThreads[0].Line != 12 {
		t.Errorf("unexpected thread location: %+v", review.FileThreads[0])
	}
	// The task's own Comment field should still carry the inline info.
	if len(review.FileThreads[0].Tasks) != 1 {
		t.Fatalf("expected task on thread")
	}
	attached := review.FileThreads[0].Tasks[0]
	if attached.Comment == nil || attached.Comment.Inline == nil || attached.Comment.Inline.Path != "x.go" {
		t.Errorf("expected task to preserve inline comment context, got %+v", attached.Comment)
	}
}

func TestGetPRReview(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasSuffix(r.URL.Path, "/comments"):
			_, _ = w.Write([]byte(`{"values": [
				{"id": 1, "content": {"raw": "review comment"}, "user": {"display_name": "A"}, "inline": {"path": "main.go", "to": 10}}
			]}`))
		case strings.HasSuffix(r.URL.Path, "/tasks"):
			_, _ = w.Write([]byte(`{"values": [
				{"id": 9, "state": "UNRESOLVED", "content": {"raw": "do thing"}, "creator": {"display_name": "A"}, "comment": {"id": 1, "inline": {"path": "main.go", "to": 10}}}
			]}`))
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer srv.Close()

	c := newRedirectClient(srv)
	review, err := c.GetPRReview("ws", "repo", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if review.PullRequestID != 3 {
		t.Errorf("expected PR ID 3, got %d", review.PullRequestID)
	}
	if len(review.FileThreads) != 1 {
		t.Fatalf("expected 1 file thread, got %d", len(review.FileThreads))
	}
	if review.FileThreads[0].File != "main.go" || review.FileThreads[0].Line != 10 {
		t.Errorf("unexpected thread location: %+v", review.FileThreads[0])
	}
	if len(review.FileThreads[0].Tasks) != 1 {
		t.Fatalf("expected task attached to thread, got %d", len(review.FileThreads[0].Tasks))
	}
	if review.FileThreads[0].Tasks[0].ID != 9 {
		t.Errorf("expected task ID 9, got %d", review.FileThreads[0].Tasks[0].ID)
	}
}

func TestListPRTasks_DeserializesFullComment(t *testing.T) {
	// Verify the PRTask struct now parses the full comment object including
	// inline info, so callers get file/line context for tasks directly.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(PaginatedResponse{
			Values: json.RawMessage(`[{
				"id": 1,
				"state": "UNRESOLVED",
				"content": {"raw": "Fix tests"},
				"comment": {
					"id": 99,
					"content": {"raw": "broken here"},
					"inline": {"path": "pkg/a.go", "to": 7}
				}
			}]`),
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
	if tasks[0].Comment == nil {
		t.Fatal("expected task.Comment to be populated")
	}
	if tasks[0].Comment.ID != 99 {
		t.Errorf("expected comment ID 99, got %d", tasks[0].Comment.ID)
	}
	if tasks[0].Comment.Inline == nil {
		t.Fatal("expected task.Comment.Inline to be populated")
	}
	if tasks[0].Comment.Inline.Path != "pkg/a.go" {
		t.Errorf("expected inline path 'pkg/a.go', got %q", tasks[0].Comment.Inline.Path)
	}
	if tasks[0].Comment.Inline.To == nil || *tasks[0].Comment.Inline.To != 7 {
		t.Errorf("expected inline 'to' 7, got %+v", tasks[0].Comment.Inline.To)
	}
}
