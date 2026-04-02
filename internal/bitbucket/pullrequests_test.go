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
