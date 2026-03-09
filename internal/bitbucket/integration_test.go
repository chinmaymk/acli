//go:build integration

package bitbucket_test

import (
	"os"
	"testing"

	"github.com/chinmaymk/acli/internal/bitbucket"
	"github.com/chinmaymk/acli/internal/config"
)

// testClient creates a Bitbucket client from the user's config.
// Skips the test if no profile is configured.
func testClient(t *testing.T) *bitbucket.Client {
	t.Helper()
	cfg, err := config.Load()
	if err != nil {
		t.Skipf("cannot load config: %v", err)
	}
	profile, err := cfg.GetProfile("")
	if err != nil {
		t.Skipf("no profile configured: %v", err)
	}
	client, err := bitbucket.NewClient(profile)
	if err != nil {
		t.Skipf("cannot create client: %v", err)
	}
	return client
}

// workspace returns the workspace slug to test against.
// Set BB_TEST_WORKSPACE env var, or defaults to the first workspace found.
func workspace(t *testing.T, client *bitbucket.Client) string {
	t.Helper()
	if ws := os.Getenv("BB_TEST_WORKSPACE"); ws != "" {
		return ws
	}
	workspaces, err := client.ListWorkspaces()
	if err != nil {
		t.Fatalf("listing workspaces: %v", err)
	}
	if len(workspaces) == 0 {
		t.Skip("no workspaces available")
	}
	return workspaces[0].Slug
}

// testRepo returns a repo slug known to exist in the workspace.
// Set BB_TEST_REPO env var, or defaults to the first repo found.
func testRepo(t *testing.T, client *bitbucket.Client, ws string) string {
	t.Helper()
	if r := os.Getenv("BB_TEST_REPO"); r != "" {
		return r
	}
	repos, err := client.ListRepositories(ws, nil)
	if err != nil {
		t.Fatalf("listing repos: %v", err)
	}
	if len(repos) == 0 {
		t.Skip("no repos in workspace")
	}
	return repos[0].Slug
}

// --- Workspace tests ---

func TestIntegration_ListWorkspaces(t *testing.T) {
	client := testClient(t)
	workspaces, err := client.ListWorkspaces()
	if err != nil {
		t.Fatalf("ListWorkspaces: %v", err)
	}
	if len(workspaces) == 0 {
		t.Fatal("expected at least one workspace")
	}
	for _, ws := range workspaces {
		if ws.Slug == "" {
			t.Error("workspace slug is empty")
		}
	}
}

func TestIntegration_GetWorkspace(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	got, err := client.GetWorkspace(ws)
	if err != nil {
		t.Fatalf("GetWorkspace(%q): %v", ws, err)
	}
	if got.Slug != ws {
		t.Errorf("expected slug %q, got %q", ws, got.Slug)
	}
}

func TestIntegration_ListWorkspaceMembers(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	members, err := client.ListWorkspaceMembers(ws)
	if err != nil {
		t.Fatalf("ListWorkspaceMembers: %v", err)
	}
	if len(members) == 0 {
		t.Fatal("expected at least one member")
	}
}

func TestIntegration_ListWorkspacePermissions(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	perms, err := client.ListWorkspacePermissions(ws)
	if err != nil {
		t.Fatalf("ListWorkspacePermissions: %v", err)
	}
	if len(perms) == 0 {
		t.Fatal("expected at least one permission entry")
	}
}

// --- Repository tests ---

func TestIntegration_ListRepositories(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	repos, err := client.ListRepositories(ws, nil)
	if err != nil {
		t.Fatalf("ListRepositories: %v", err)
	}
	if len(repos) == 0 {
		t.Fatal("expected at least one repo")
	}
	for _, r := range repos {
		if r.Slug == "" {
			t.Error("repo slug is empty")
		}
		if r.FullName == "" {
			t.Error("repo full_name is empty")
		}
	}
}

func TestIntegration_ListRepositories_WithOptions(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	repos, err := client.ListRepositories(ws, &bitbucket.ListReposOptions{
		Role: "owner",
		Sort: "-updated_on",
	})
	if err != nil {
		t.Fatalf("ListRepositories with options: %v", err)
	}
	// Just verify it doesn't error; may be empty if user owns nothing
	_ = repos
}

func TestIntegration_GetRepository(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	repo := testRepo(t, client, ws)
	got, err := client.GetRepository(ws, repo)
	if err != nil {
		t.Fatalf("GetRepository(%q, %q): %v", ws, repo, err)
	}
	if got.Slug != repo {
		t.Errorf("expected slug %q, got %q", repo, got.Slug)
	}
	if got.SCM == "" {
		t.Error("expected SCM to be set")
	}
}

func TestIntegration_ListForks(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	repo := testRepo(t, client, ws)
	// Should not error even if no forks exist
	_, err := client.ListForks(ws, repo)
	if err != nil {
		t.Fatalf("ListForks: %v", err)
	}
}

// --- Branch tests ---

func TestIntegration_ListBranches(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	repo := testRepo(t, client, ws)
	branches, err := client.ListBranches(ws, repo, "")
	if err != nil {
		t.Fatalf("ListBranches: %v", err)
	}
	if len(branches) == 0 {
		t.Fatal("expected at least one branch")
	}
	for _, b := range branches {
		if b.Name == "" {
			t.Error("branch name is empty")
		}
	}
}

func TestIntegration_GetBranch(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	repo := testRepo(t, client, ws)
	// Get the first branch name
	branches, err := client.ListBranches(ws, repo, "")
	if err != nil || len(branches) == 0 {
		t.Skip("no branches to test")
	}
	got, err := client.GetBranch(ws, repo, branches[0].Name)
	if err != nil {
		t.Fatalf("GetBranch(%q): %v", branches[0].Name, err)
	}
	if got.Name != branches[0].Name {
		t.Errorf("expected %q, got %q", branches[0].Name, got.Name)
	}
}

// --- Commit tests ---

func TestIntegration_ListCommits(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	repo := testRepo(t, client, ws)
	commits, err := client.ListCommits(ws, repo, "", "")
	if err != nil {
		t.Fatalf("ListCommits: %v", err)
	}
	if len(commits) == 0 {
		t.Fatal("expected at least one commit")
	}
	for _, c := range commits {
		if c.Hash == "" {
			t.Error("commit hash is empty")
		}
	}
}

func TestIntegration_GetCommit(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	repo := testRepo(t, client, ws)
	commits, err := client.ListCommits(ws, repo, "", "")
	if err != nil || len(commits) == 0 {
		t.Skip("no commits to test")
	}
	got, err := client.GetCommit(ws, repo, commits[0].Hash)
	if err != nil {
		t.Fatalf("GetCommit(%q): %v", commits[0].Hash, err)
	}
	if got.Hash != commits[0].Hash {
		t.Errorf("expected %q, got %q", commits[0].Hash, got.Hash)
	}
}

func TestIntegration_ListCommitStatuses(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	repo := testRepo(t, client, ws)
	commits, err := client.ListCommits(ws, repo, "", "")
	if err != nil || len(commits) == 0 {
		t.Skip("no commits to test")
	}
	// Should not error even if no statuses exist
	_, err = client.ListCommitStatuses(ws, repo, commits[0].Hash)
	if err != nil {
		t.Fatalf("ListCommitStatuses: %v", err)
	}
}

// --- Tag tests ---

func TestIntegration_ListTags(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	repo := testRepo(t, client, ws)
	// Should not error even if no tags exist
	_, err := client.ListTags(ws, repo, "")
	if err != nil {
		t.Fatalf("ListTags: %v", err)
	}
}

// --- Pull Request tests ---

func TestIntegration_ListPullRequests(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	repo := testRepo(t, client, ws)
	// Should not error even if no PRs exist
	_, err := client.ListPullRequests(ws, repo, nil)
	if err != nil {
		t.Fatalf("ListPullRequests: %v", err)
	}
}

func TestIntegration_ListPullRequests_WithState(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	repo := testRepo(t, client, ws)
	_, err := client.ListPullRequests(ws, repo, &bitbucket.ListPRsOptions{
		State: "MERGED",
	})
	if err != nil {
		t.Fatalf("ListPullRequests with state: %v", err)
	}
}

// --- Pipeline tests ---

func TestIntegration_ListPipelines(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	repo := testRepo(t, client, ws)
	_, err := client.ListPipelines(ws, repo, nil)
	if err != nil {
		t.Fatalf("ListPipelines: %v", err)
	}
}

func TestIntegration_ListPipelineVariables(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	repo := testRepo(t, client, ws)
	_, err := client.ListPipelineVariables(ws, repo)
	if err != nil {
		t.Fatalf("ListPipelineVariables: %v", err)
	}
}

// --- Webhook tests ---

func TestIntegration_ListRepoWebhooks(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	repo := testRepo(t, client, ws)
	_, err := client.ListRepoWebhooks(ws, repo)
	if err != nil {
		t.Fatalf("ListRepoWebhooks: %v", err)
	}
}

func TestIntegration_ListWorkspaceWebhooks(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	_, err := client.ListWorkspaceWebhooks(ws)
	if err != nil {
		t.Fatalf("ListWorkspaceWebhooks: %v", err)
	}
}

// --- Deploy Key tests ---

func TestIntegration_ListDeployKeys(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	repo := testRepo(t, client, ws)
	_, err := client.ListDeployKeys(ws, repo)
	if err != nil {
		t.Fatalf("ListDeployKeys: %v", err)
	}
}

// --- Environment tests ---

func TestIntegration_ListEnvironments(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	repo := testRepo(t, client, ws)
	_, err := client.ListEnvironments(ws, repo)
	if err != nil {
		t.Fatalf("ListEnvironments: %v", err)
	}
}

// --- Deployment tests ---

func TestIntegration_ListDeployments(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	repo := testRepo(t, client, ws)
	_, err := client.ListDeployments(ws, repo)
	if err != nil {
		t.Fatalf("ListDeployments: %v", err)
	}
}

// --- Branch Restriction tests ---

func TestIntegration_ListBranchRestrictions(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	repo := testRepo(t, client, ws)
	_, err := client.ListBranchRestrictions(ws, repo)
	if err != nil {
		t.Fatalf("ListBranchRestrictions: %v", err)
	}
}

// --- Project tests ---

func TestIntegration_ListProjects(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	projects, err := client.ListProjects(ws)
	if err != nil {
		t.Fatalf("ListProjects: %v", err)
	}
	if len(projects) == 0 {
		t.Fatal("expected at least one project")
	}
	for _, p := range projects {
		if p.Key == "" {
			t.Error("project key is empty")
		}
	}
}

func TestIntegration_GetProject(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	projects, err := client.ListProjects(ws)
	if err != nil || len(projects) == 0 {
		t.Skip("no projects to test")
	}
	got, err := client.GetProject(ws, projects[0].Key)
	if err != nil {
		t.Fatalf("GetProject(%q): %v", projects[0].Key, err)
	}
	if got.Key != projects[0].Key {
		t.Errorf("expected %q, got %q", projects[0].Key, got.Key)
	}
}

// --- CRUD lifecycle test: Project ---

func TestIntegration_ProjectLifecycle(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)

	// Create
	proj, err := client.CreateProject(ws, &bitbucket.CreateProjectRequest{
		Name:      "acli-integration-test",
		Key:       "ACLITEST",
		IsPrivate: true,
	})
	if err != nil {
		t.Fatalf("CreateProject: %v", err)
	}
	if proj.Key != "ACLITEST" {
		t.Errorf("expected key ACLITEST, got %q", proj.Key)
	}

	// Get
	got, err := client.GetProject(ws, "ACLITEST")
	if err != nil {
		t.Fatalf("GetProject: %v", err)
	}
	if got.Name != "acli-integration-test" {
		t.Errorf("expected name %q, got %q", "acli-integration-test", got.Name)
	}

	// Delete
	if err := client.DeleteProject(ws, "ACLITEST"); err != nil {
		t.Fatalf("DeleteProject: %v", err)
	}
}

// --- CRUD lifecycle test: Branch ---

func TestIntegration_BranchLifecycle(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	repo := testRepo(t, client, ws)

	// Need a commit to branch from
	commits, err := client.ListCommits(ws, repo, "", "")
	if err != nil || len(commits) == 0 {
		t.Skip("no commits to branch from")
	}

	branchName := "acli-integration-test-branch"

	// Create
	branch, err := client.CreateBranch(ws, repo, &bitbucket.CreateBranchRequest{
		Name: branchName,
		Target: struct {
			Hash string `json:"hash"`
		}{Hash: commits[0].Hash},
	})
	if err != nil {
		t.Fatalf("CreateBranch: %v", err)
	}
	if branch.Name != branchName {
		t.Errorf("expected %q, got %q", branchName, branch.Name)
	}

	// Get
	got, err := client.GetBranch(ws, repo, branchName)
	if err != nil {
		t.Fatalf("GetBranch: %v", err)
	}
	if got.Name != branchName {
		t.Errorf("expected %q, got %q", branchName, got.Name)
	}

	// Delete
	if err := client.DeleteBranch(ws, repo, branchName); err != nil {
		t.Fatalf("DeleteBranch: %v", err)
	}
}

// --- CRUD lifecycle test: Tag ---

func TestIntegration_TagLifecycle(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	repo := testRepo(t, client, ws)

	commits, err := client.ListCommits(ws, repo, "", "")
	if err != nil || len(commits) == 0 {
		t.Skip("no commits to tag")
	}

	tagName := "acli-integration-test-tag"

	// Create
	tag, err := client.CreateTag(ws, repo, &bitbucket.CreateTagRequest{
		Name: tagName,
		Target: struct {
			Hash string `json:"hash"`
		}{Hash: commits[0].Hash},
	})
	if err != nil {
		t.Fatalf("CreateTag: %v", err)
	}
	if tag.Name != tagName {
		t.Errorf("expected %q, got %q", tagName, tag.Name)
	}

	// Get
	got, err := client.GetTag(ws, repo, tagName)
	if err != nil {
		t.Fatalf("GetTag: %v", err)
	}
	if got.Name != tagName {
		t.Errorf("expected %q, got %q", tagName, got.Name)
	}

	// Delete
	if err := client.DeleteTag(ws, repo, tagName); err != nil {
		t.Fatalf("DeleteTag: %v", err)
	}
}

// --- Diff test ---

func TestIntegration_GetDiff(t *testing.T) {
	client := testClient(t)
	ws := workspace(t, client)
	repo := testRepo(t, client, ws)

	commits, err := client.ListCommits(ws, repo, "", "")
	if err != nil || len(commits) < 2 {
		t.Skip("need at least 2 commits for diff")
	}

	// Diff between two commits
	spec := commits[1].Hash + ".." + commits[0].Hash
	diff, err := client.GetDiff(ws, repo, spec)
	if err != nil {
		t.Fatalf("GetDiff(%q): %v", spec, err)
	}
	if diff == "" {
		t.Log("diff is empty (commits may be identical)")
	}
}
