package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
)

type PullRequest struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	State       string `json:"state"`
	CreatedOn   string `json:"created_on"`
	UpdatedOn   string `json:"updated_on"`
	Author      struct {
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
	} `json:"author"`
	Source struct {
		Branch struct {
			Name string `json:"name"`
		} `json:"branch"`
		Repository struct {
			FullName string `json:"full_name"`
		} `json:"repository"`
	} `json:"source"`
	Destination struct {
		Branch struct {
			Name string `json:"name"`
		} `json:"branch"`
		Repository struct {
			FullName string `json:"full_name"`
		} `json:"repository"`
	} `json:"destination"`
	CloseSourceBranch bool `json:"close_source_branch"`
	CommentCount      int  `json:"comment_count"`
	TaskCount         int  `json:"task_count"`
	Links             struct {
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
	} `json:"links"`
}

type ListPRsOptions struct {
	State   string
	Author  string // Bitbucket nickname (username) or UUID wrapped in braces, e.g. "{d301aafa-...}"
	Page    int
	PageLen int
	All     bool
}

func (c *Client) ListPullRequests(workspace, repoSlug string, opts *ListPRsOptions) ([]PullRequest, error) {
	params := url.Values{}
	if opts != nil {
		state := strings.ToUpper(opts.State)
		var qParts []string
		if state != "" {
			qParts = append(qParts, fmt.Sprintf(`state="%s"`, state))
		}
		if opts.Author != "" {
			if len(opts.Author) > 2 && opts.Author[0] == '{' && opts.Author[len(opts.Author)-1] == '}' {
				qParts = append(qParts, fmt.Sprintf(`author.uuid="%s"`, opts.Author))
			} else {
				qParts = append(qParts, fmt.Sprintf(`author.nickname="%s"`, opts.Author))
			}
		}
		if len(qParts) > 0 {
			params.Set("q", strings.Join(qParts, " AND "))
		}
		if opts.Page > 0 {
			params.Set("page", fmt.Sprintf("%d", opts.Page))
		}
		if opts.PageLen > 0 {
			params.Set("pagelen", fmt.Sprintf("%d", opts.PageLen))
		}
	}
	ensurePageLen(params)

	path := fmt.Sprintf("/repositories/%s/%s/pullrequests",
		url.PathEscape(workspace), url.PathEscape(repoSlug))
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	if opts != nil && opts.All {
		pages, err := c.getAll(path)
		if err != nil && len(pages) == 0 {
			return nil, err
		}
		var prs []PullRequest
		for _, pg := range pages {
			var pagePRs []PullRequest
			if err := json.Unmarshal(pg.Values, &pagePRs); err != nil {
				return prs, fmt.Errorf("parsing pull requests: %w", err)
			}
			prs = append(prs, pagePRs...)
		}
		return prs, nil
	}

	data, err := c.get(path)
	if err != nil {
		return nil, err
	}

	var page PaginatedResponse
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	var prs []PullRequest
	if err := json.Unmarshal(page.Values, &prs); err != nil {
		return nil, fmt.Errorf("parsing pull requests: %w", err)
	}

	return prs, nil
}

func (c *Client) GetPullRequest(workspace, repoSlug string, prID int) (*PullRequest, error) {
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d",
		url.PathEscape(workspace), url.PathEscape(repoSlug), prID)

	data, err := c.get(path)
	if err != nil {
		return nil, err
	}

	var pr PullRequest
	if err := json.Unmarshal(data, &pr); err != nil {
		return nil, fmt.Errorf("parsing pull request: %w", err)
	}

	return &pr, nil
}

type CreatePRRequest struct {
	Title             string `json:"title"`
	Description       string `json:"description,omitempty"`
	SourceBranch      string `json:"-"`
	DestinationBranch string `json:"-"`
	CloseSourceBranch bool   `json:"close_source_branch,omitempty"`
}

type createPRBody struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Source      struct {
		Branch struct {
			Name string `json:"name"`
		} `json:"branch"`
	} `json:"source"`
	Destination       *prBranchRef `json:"destination,omitempty"`
	CloseSourceBranch bool         `json:"close_source_branch,omitempty"`
}

type prBranchRef struct {
	Branch struct {
		Name string `json:"name"`
	} `json:"branch"`
}

func (c *Client) CreatePullRequest(workspace, repoSlug string, req *CreatePRRequest) (*PullRequest, error) {
	body := createPRBody{
		Title:             req.Title,
		Description:       req.Description,
		CloseSourceBranch: req.CloseSourceBranch,
	}
	body.Source.Branch.Name = req.SourceBranch

	if req.DestinationBranch != "" {
		body.Destination = &prBranchRef{}
		body.Destination.Branch.Name = req.DestinationBranch
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("/repositories/%s/%s/pullrequests",
		url.PathEscape(workspace), url.PathEscape(repoSlug))

	data, err := c.post(path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	var pr PullRequest
	if err := json.Unmarshal(data, &pr); err != nil {
		return nil, fmt.Errorf("parsing pull request: %w", err)
	}

	return &pr, nil
}

type UpdatePRRequest struct {
	Title             string `json:"title,omitempty"`
	Description       string `json:"description,omitempty"`
	CloseSourceBranch *bool  `json:"close_source_branch,omitempty"`
}

func (c *Client) UpdatePullRequest(workspace, repoSlug string, prID int, req *UpdatePRRequest) (*PullRequest, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d",
		url.PathEscape(workspace), url.PathEscape(repoSlug), prID)
	data, err := c.put(path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	var pr PullRequest
	if err := json.Unmarshal(data, &pr); err != nil {
		return nil, fmt.Errorf("parsing pull request: %w", err)
	}
	return &pr, nil
}

// Participant represents a pull request participant (reviewer/approver).
type Participant struct {
	User struct {
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
	} `json:"user"`
	Role     string `json:"role"`
	Approved bool   `json:"approved"`
	State    string `json:"state"`
}

func (c *Client) ApprovePullRequest(workspace, repoSlug string, prID int) (*Participant, error) {
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/approve",
		url.PathEscape(workspace), url.PathEscape(repoSlug), prID)
	data, err := c.post(path, nil)
	if err != nil {
		return nil, err
	}
	var p Participant
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("parsing participant: %w", err)
	}
	return &p, nil
}

func (c *Client) UnapprovePullRequest(workspace, repoSlug string, prID int) error {
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/approve",
		url.PathEscape(workspace), url.PathEscape(repoSlug), prID)
	return c.deleteNoContent(path)
}

func (c *Client) DeclinePullRequest(workspace, repoSlug string, prID int) (*PullRequest, error) {
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/decline",
		url.PathEscape(workspace), url.PathEscape(repoSlug), prID)
	data, err := c.post(path, nil)
	if err != nil {
		return nil, err
	}
	var pr PullRequest
	if err := json.Unmarshal(data, &pr); err != nil {
		return nil, fmt.Errorf("parsing pull request: %w", err)
	}
	return &pr, nil
}

type MergePRRequest struct {
	MergeStrategy     string `json:"merge_strategy,omitempty"`
	CloseSourceBranch *bool  `json:"close_source_branch,omitempty"`
	Message           string `json:"message,omitempty"`
}

func (c *Client) MergePullRequest(workspace, repoSlug string, prID int, req *MergePRRequest) (*PullRequest, error) {
	var body io.Reader
	if req != nil {
		jsonBody, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(jsonBody)
	}
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/merge",
		url.PathEscape(workspace), url.PathEscape(repoSlug), prID)
	data, err := c.post(path, body)
	if err != nil {
		return nil, err
	}
	var pr PullRequest
	if err := json.Unmarshal(data, &pr); err != nil {
		return nil, fmt.Errorf("parsing pull request: %w", err)
	}
	return &pr, nil
}

func (c *Client) RequestChangesPullRequest(workspace, repoSlug string, prID int) (*Participant, error) {
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/request-changes",
		url.PathEscape(workspace), url.PathEscape(repoSlug), prID)
	data, err := c.post(path, nil)
	if err != nil {
		return nil, err
	}
	var p Participant
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("parsing participant: %w", err)
	}
	return &p, nil
}

func (c *Client) RemoveRequestChangesPullRequest(workspace, repoSlug string, prID int) error {
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/request-changes",
		url.PathEscape(workspace), url.PathEscape(repoSlug), prID)
	return c.deleteNoContent(path)
}

type PRComment struct {
	ID      int    `json:"id"`
	Type    string `json:"type,omitempty"`
	Content struct {
		Raw    string `json:"raw"`
		Markup string `json:"markup"`
		HTML   string `json:"html"`
	} `json:"content"`
	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
	User      struct {
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
	} `json:"user"`
	Deleted bool `json:"deleted,omitempty"`
	Inline  *struct {
		Path      string `json:"path"`
		From      *int   `json:"from"`
		To        *int   `json:"to"`
		StartFrom *int   `json:"start_from,omitempty"`
		StartTo   *int   `json:"start_to,omitempty"`
	} `json:"inline,omitempty"`
	Parent *struct {
		ID int `json:"id"`
	} `json:"parent,omitempty"`
	Resolution *CommentResolution `json:"resolution,omitempty"`
	Pending    *bool              `json:"pending,omitempty"`
}

// CommentResolution represents the resolution of a comment.
type CommentResolution struct {
	Type string `json:"type"`
	User struct {
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
	} `json:"user"`
	CreatedOn string `json:"created_on"`
}

func (c *Client) ListPRComments(workspace, repoSlug string, prID int, opts *PaginationOptions) ([]PRComment, error) {
	params := url.Values{}
	if opts != nil {
		opts.applyParams(params)
	}
	ensurePageLen(params)

	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/comments",
		url.PathEscape(workspace), url.PathEscape(repoSlug), prID)
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	if opts != nil && opts.All {
		pages, err := c.getAll(path)
		if err != nil && len(pages) == 0 {
			return nil, err
		}
		var comments []PRComment
		for _, pg := range pages {
			var pageComments []PRComment
			if err := json.Unmarshal(pg.Values, &pageComments); err != nil {
				return comments, fmt.Errorf("parsing comments: %w", err)
			}
			comments = append(comments, pageComments...)
		}
		return comments, nil
	}

	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var page PaginatedResponse
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	var comments []PRComment
	if err := json.Unmarshal(page.Values, &comments); err != nil {
		return nil, fmt.Errorf("parsing comments: %w", err)
	}
	return comments, nil
}

// InlineCommentParams specifies the file and line for an inline PR comment.
type InlineCommentParams struct {
	Path string
	To   int // Line number in the new version of the file
}

func (c *Client) CreatePRComment(workspace, repoSlug string, prID int, content string) (*PRComment, error) {
	return c.CreatePRCommentInline(workspace, repoSlug, prID, content, nil)
}

func (c *Client) CreatePRCommentInline(workspace, repoSlug string, prID int, content string, inline *InlineCommentParams) (*PRComment, error) {
	body := map[string]interface{}{
		"content": map[string]string{
			"raw": content,
		},
	}
	if inline != nil {
		body["inline"] = map[string]interface{}{
			"path": inline.Path,
			"to":   inline.To,
		}
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/comments",
		url.PathEscape(workspace), url.PathEscape(repoSlug), prID)
	data, err := c.post(path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	var comment PRComment
	if err := json.Unmarshal(data, &comment); err != nil {
		return nil, fmt.Errorf("parsing comment: %w", err)
	}
	return &comment, nil
}

func (c *Client) GetPRDiff(workspace, repoSlug string, prID int) (string, error) {
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/diff",
		url.PathEscape(workspace), url.PathEscape(repoSlug), prID)
	data, err := c.getRaw(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// PR Task types and methods

// PRTask represents a task on a pull request.
//
// When a task is linked to a comment — which is how Bitbucket attaches tasks to
// a specific file/line for code review — the Comment field contains the full
// comment object, including its Inline location (path, line number) if any.
// Coding agents can use Comment.Inline to know where a task applies in the diff
// without making an extra API call.
type PRTask struct {
	ID      int `json:"id"`
	State   string `json:"state"`
	Content struct {
		Raw    string `json:"raw"`
		Markup string `json:"markup"`
		HTML   string `json:"html"`
	} `json:"content"`
	Creator struct {
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
		Nickname    string `json:"nickname,omitempty"`
	} `json:"creator"`
	CreatedOn  string `json:"created_on"`
	UpdatedOn  string `json:"updated_on"`
	ResolvedOn string `json:"resolved_on,omitempty"`
	ResolvedBy *struct {
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
	} `json:"resolved_by,omitempty"`
	// Comment is the comment this task is attached to, or nil for a standalone
	// PR-level task. When set, Comment.Inline (if non-nil) gives the file path
	// and line number the task applies to.
	Comment *PRComment `json:"comment,omitempty"`
	Links   struct {
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
	} `json:"links,omitempty"`
}

func (c *Client) ListPRTasks(workspace, repoSlug string, prID int, opts *PaginationOptions) ([]PRTask, error) {
	params := url.Values{}
	if opts != nil {
		opts.applyParams(params)
	}
	ensurePageLen(params)

	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/tasks",
		url.PathEscape(workspace), url.PathEscape(repoSlug), prID)
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	if opts != nil && opts.All {
		pages, err := c.getAll(path)
		if err != nil && len(pages) == 0 {
			return nil, err
		}
		var tasks []PRTask
		for _, pg := range pages {
			var pageTasks []PRTask
			if err := json.Unmarshal(pg.Values, &pageTasks); err != nil {
				return tasks, fmt.Errorf("parsing tasks: %w", err)
			}
			tasks = append(tasks, pageTasks...)
		}
		return tasks, nil
	}

	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var page PaginatedResponse
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	var tasks []PRTask
	if err := json.Unmarshal(page.Values, &tasks); err != nil {
		return nil, fmt.Errorf("parsing tasks: %w", err)
	}
	return tasks, nil
}

func (c *Client) GetPRTask(workspace, repoSlug string, prID, taskID int) (*PRTask, error) {
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/tasks/%d",
		url.PathEscape(workspace), url.PathEscape(repoSlug), prID, taskID)
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var task PRTask
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, fmt.Errorf("parsing task: %w", err)
	}
	return &task, nil
}

type CreatePRTaskRequest struct {
	Content   string `json:"-"`
	CommentID *int   `json:"-"`
}

func (c *Client) CreatePRTask(workspace, repoSlug string, prID int, req *CreatePRTaskRequest) (*PRTask, error) {
	body := map[string]interface{}{
		"content": map[string]string{
			"raw": req.Content,
		},
	}
	if req.CommentID != nil {
		body["comment"] = map[string]int{
			"id": *req.CommentID,
		}
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/tasks",
		url.PathEscape(workspace), url.PathEscape(repoSlug), prID)
	data, err := c.post(path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	var task PRTask
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, fmt.Errorf("parsing task: %w", err)
	}
	return &task, nil
}

type UpdatePRTaskRequest struct {
	Content *string `json:"-"`
	State   string  `json:"-"`
}

func (c *Client) UpdatePRTask(workspace, repoSlug string, prID, taskID int, req *UpdatePRTaskRequest) (*PRTask, error) {
	body := map[string]interface{}{}
	if req.Content != nil {
		body["content"] = map[string]string{
			"raw": *req.Content,
		}
	}
	if req.State != "" {
		body["state"] = req.State
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/tasks/%d",
		url.PathEscape(workspace), url.PathEscape(repoSlug), prID, taskID)
	data, err := c.put(path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	var task PRTask
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, fmt.Errorf("parsing task: %w", err)
	}
	return &task, nil
}

func (c *Client) DeletePRTask(workspace, repoSlug string, prID, taskID int) error {
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/tasks/%d",
		url.PathEscape(workspace), url.PathEscape(repoSlug), prID, taskID)
	return c.deleteNoContent(path)
}

// PR Commits

func (c *Client) ListPRCommits(workspace, repoSlug string, prID int, opts *PaginationOptions) ([]Commit, error) {
	params := url.Values{}
	if opts != nil {
		opts.applyParams(params)
	}
	ensurePageLen(params)

	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/commits",
		url.PathEscape(workspace), url.PathEscape(repoSlug), prID)
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	if opts != nil && opts.All {
		pages, err := c.getAll(path)
		if err != nil && len(pages) == 0 {
			return nil, err
		}
		var commits []Commit
		for _, pg := range pages {
			var pageCommits []Commit
			if err := json.Unmarshal(pg.Values, &pageCommits); err != nil {
				return commits, fmt.Errorf("parsing commits: %w", err)
			}
			commits = append(commits, pageCommits...)
		}
		return commits, nil
	}

	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var page PaginatedResponse
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	var commits []Commit
	if err := json.Unmarshal(page.Values, &commits); err != nil {
		return nil, fmt.Errorf("parsing commits: %w", err)
	}
	return commits, nil
}

// PR Diffstat

// CommitFile represents a file at a specific commit in a repository.
type CommitFile struct {
	Type        string  `json:"type"`
	Path        string  `json:"path"`
	Commit      *Commit `json:"commit,omitempty"`
	Attributes  string  `json:"attributes,omitempty"` // link, executable, subrepository, binary, lfs
	EscapedPath string  `json:"escaped_path,omitempty"`
}

// DiffStat represents file-level change statistics for a pull request.
type DiffStat struct {
	Type         string      `json:"type"`
	Status       string      `json:"status"` // added, removed, modified, renamed
	Old          *CommitFile `json:"old"`
	New          *CommitFile `json:"new"`
	LinesAdded   int         `json:"lines_added"`
	LinesRemoved int         `json:"lines_removed"`
}

func (c *Client) GetPRDiffStat(workspace, repoSlug string, prID int, opts *PaginationOptions) ([]DiffStat, error) {
	params := url.Values{}
	if opts != nil {
		opts.applyParams(params)
	}
	ensurePageLen(params)

	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/diffstat",
		url.PathEscape(workspace), url.PathEscape(repoSlug), prID)
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	if opts != nil && opts.All {
		pages, err := c.getAll(path)
		if err != nil && len(pages) == 0 {
			return nil, err
		}
		var stats []DiffStat
		for _, pg := range pages {
			var pageStats []DiffStat
			if err := json.Unmarshal(pg.Values, &pageStats); err != nil {
				return stats, fmt.Errorf("parsing diffstat: %w", err)
			}
			stats = append(stats, pageStats...)
		}
		return stats, nil
	}

	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var page PaginatedResponse
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	var stats []DiffStat
	if err := json.Unmarshal(page.Values, &stats); err != nil {
		return nil, fmt.Errorf("parsing diffstat: %w", err)
	}
	return stats, nil
}

// PR Activity

// PRActivity represents an activity entry on a pull request (comment, approval, update, etc).
type PRActivity struct {
	Approval *PRActivityApproval `json:"approval,omitempty"`
	Update   *PRActivityUpdate   `json:"update,omitempty"`
	Comment  *PRComment          `json:"comment,omitempty"`
	// PullRequest is the PR reference included in each activity entry.
	PullRequest *struct {
		Type  string `json:"type"`
		ID    int    `json:"id"`
		Title string `json:"title"`
		Links struct {
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
		} `json:"links"`
	} `json:"pull_request,omitempty"`
}

// PRActivityApproval represents an approval activity on a pull request.
type PRActivityApproval struct {
	Date string `json:"date"`
	User struct {
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
		Nickname    string `json:"nickname"`
		Type        string `json:"type"`
		AccountID   string `json:"account_id"`
	} `json:"user"`
}

// PRActivityUpdate represents a state-change activity on a pull request.
type PRActivityUpdate struct {
	Date        string `json:"date"`
	Title       string `json:"title"`
	Description string `json:"description"`
	State       string `json:"state"` // OPEN, MERGED, DECLINED
	Reason      string `json:"reason"`
	Author      struct {
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
		Nickname    string `json:"nickname"`
		Type        string `json:"type"`
		AccountID   string `json:"account_id"`
	} `json:"author"`
	Source struct {
		Branch struct {
			Name string `json:"name"`
		} `json:"branch"`
		Commit struct {
			Hash string `json:"hash"`
			Type string `json:"type"`
		} `json:"commit"`
		Repository struct {
			FullName string `json:"full_name"`
			UUID     string `json:"uuid"`
		} `json:"repository"`
	} `json:"source"`
	Destination struct {
		Branch struct {
			Name string `json:"name"`
		} `json:"branch"`
		Commit struct {
			Hash string `json:"hash"`
			Type string `json:"type"`
		} `json:"commit"`
		Repository struct {
			FullName string `json:"full_name"`
			UUID     string `json:"uuid"`
		} `json:"repository"`
	} `json:"destination"`
}

func (c *Client) ListPRActivity(workspace, repoSlug string, prID int, opts *PaginationOptions) ([]PRActivity, error) {
	params := url.Values{}
	if opts != nil {
		opts.applyParams(params)
	}
	ensurePageLen(params)

	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/activity",
		url.PathEscape(workspace), url.PathEscape(repoSlug), prID)
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	if opts != nil && opts.All {
		pages, err := c.getAll(path)
		if err != nil && len(pages) == 0 {
			return nil, err
		}
		var activities []PRActivity
		for _, pg := range pages {
			var pageActivities []PRActivity
			if err := json.Unmarshal(pg.Values, &pageActivities); err != nil {
				return activities, fmt.Errorf("parsing activity: %w", err)
			}
			activities = append(activities, pageActivities...)
		}
		return activities, nil
	}

	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var page PaginatedResponse
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	var activities []PRActivity
	if err := json.Unmarshal(page.Values, &activities); err != nil {
		return nil, fmt.Errorf("parsing activity: %w", err)
	}
	return activities, nil
}

// PR Review: joined view of comments and tasks

// PRReviewThread represents a single thread of feedback on a pull request.
//
// A thread is rooted at a top-level comment (Comment.Parent is nil). Replies
// to that comment are in Replies (flattened in chronological order). Any tasks
// linked to the root or to any of its replies are collected in Tasks. If the
// root comment is an inline comment, File/Line/StartLine describe the code
// location the thread is anchored to.
type PRReviewThread struct {
	// File is the path the thread is anchored to, empty for general (non-inline) threads.
	File string `json:"file,omitempty"`
	// Line is the line number in the new version of the file the thread is
	// anchored to (for multi-line comments this is the ending line).
	Line int `json:"line,omitempty"`
	// StartLine is the starting line for multi-line comments, 0 if single-line.
	StartLine int `json:"start_line,omitempty"`
	// Comment is the root (top-level) comment of the thread.
	Comment PRComment `json:"comment"`
	// Replies are non-root comments that chain under Comment, in chronological order.
	Replies []PRComment `json:"replies,omitempty"`
	// Tasks are tasks linked to the root comment or any of its replies.
	Tasks []PRTask `json:"tasks,omitempty"`
}

// PRReview is a consolidated view of all feedback on a pull request, joining
// comments and tasks with correct associations. Coding agents can consume this
// as a single payload rather than manually cross-referencing the comment and
// task list endpoints.
type PRReview struct {
	PullRequestID int `json:"pull_request_id"`
	// FileThreads are threads anchored to a specific file/line, grouped by file path.
	FileThreads []PRReviewThread `json:"file_threads,omitempty"`
	// GeneralThreads are top-level PR discussion threads not anchored to any file.
	GeneralThreads []PRReviewThread `json:"general_threads,omitempty"`
	// StandaloneTasks are tasks that are not linked to any comment (or whose
	// linked comment could not be found, e.g. deleted).
	StandaloneTasks []PRTask `json:"standalone_tasks,omitempty"`
	// Counts is a quick summary useful for agents deciding whether to act.
	Counts PRReviewCounts `json:"counts"`
}

// PRReviewCounts summarises the feedback on a pull request.
type PRReviewCounts struct {
	Comments         int `json:"comments"`          // total non-deleted comments (roots + replies)
	Threads          int `json:"threads"`           // total threads (file + general)
	Tasks            int `json:"tasks"`             // total tasks (linked + standalone)
	UnresolvedTasks  int `json:"unresolved_tasks"`  // tasks with state != RESOLVED
	UnresolvedThreads int `json:"unresolved_threads"` // threads whose root comment has no resolution
}

// GetPRReview fetches all comments and tasks on a pull request and returns a
// joined view: threads (rooted at top-level comments, with replies) with any
// linked tasks attached to the thread. This is the recommended way to retrieve
// PR feedback with correct associations.
func (c *Client) GetPRReview(workspace, repoSlug string, prID int) (*PRReview, error) {
	comments, err := c.ListPRComments(workspace, repoSlug, prID, &PaginationOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("listing comments: %w", err)
	}
	tasks, err := c.ListPRTasks(workspace, repoSlug, prID, &PaginationOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("listing tasks: %w", err)
	}
	return BuildPRReview(prID, comments, tasks), nil
}

// BuildPRReview joins comments and tasks into a structured review. It is
// exported so callers that already have comments/tasks (e.g. from a cached
// response or a test) can assemble the review without another API round-trip.
//
// Association rules:
//   - A top-level comment (Parent == nil) starts a new thread. If its Inline
//     field is set, the thread is a "file thread" for that path+line,
//     otherwise it is a "general thread".
//   - A non-root comment is attached as a reply to the thread whose root is
//     the transitive parent of that comment.
//   - A task is attached to the thread of its linked comment's root. If the
//     task has no linked comment, or the linked comment is not present in the
//     given comments slice, the task is added to StandaloneTasks.
func BuildPRReview(prID int, comments []PRComment, tasks []PRTask) *PRReview {
	review := &PRReview{PullRequestID: prID}

	// Index comments by id and find the root of each comment.
	byID := make(map[int]*PRComment, len(comments))
	for i := range comments {
		byID[comments[i].ID] = &comments[i]
	}

	// rootOf walks the parent chain until it finds a comment with no parent.
	// Returns nil if the chain is broken (parent not in byID).
	rootOf := func(c *PRComment) *PRComment {
		cur := c
		// Cap the walk to guard against pathological cycles.
		for i := 0; cur != nil && cur.Parent != nil && i < 1000; i++ {
			next, ok := byID[cur.Parent.ID]
			if !ok {
				return nil
			}
			cur = next
		}
		return cur
	}

	// Build threads keyed by root comment ID, preserving order of first appearance.
	threadByRoot := make(map[int]*PRReviewThread)
	var rootOrder []int
	for i := range comments {
		c := &comments[i]
		if c.Parent != nil || c.Deleted {
			continue
		}
		t := &PRReviewThread{Comment: *c}
		if c.Inline != nil {
			t.File = c.Inline.Path
			if c.Inline.To != nil {
				t.Line = *c.Inline.To
			}
			if c.Inline.StartTo != nil {
				t.StartLine = *c.Inline.StartTo
			}
		}
		threadByRoot[c.ID] = t
		rootOrder = append(rootOrder, c.ID)
	}

	// Attach replies to their root thread.
	for i := range comments {
		c := &comments[i]
		if c.Parent == nil || c.Deleted {
			continue
		}
		root := rootOf(c)
		if root == nil {
			continue
		}
		if t, ok := threadByRoot[root.ID]; ok {
			t.Replies = append(t.Replies, *c)
		}
	}

	// Attach tasks to their thread (by linked comment's root), or to standalone.
	for i := range tasks {
		task := tasks[i]
		if task.Comment == nil {
			review.StandaloneTasks = append(review.StandaloneTasks, task)
			continue
		}
		linked, ok := byID[task.Comment.ID]
		if !ok {
			review.StandaloneTasks = append(review.StandaloneTasks, task)
			continue
		}
		root := rootOf(linked)
		if root == nil {
			review.StandaloneTasks = append(review.StandaloneTasks, task)
			continue
		}
		if t, ok := threadByRoot[root.ID]; ok {
			t.Tasks = append(t.Tasks, task)
		} else {
			review.StandaloneTasks = append(review.StandaloneTasks, task)
		}
	}

	// Split into file vs general threads preserving original order.
	for _, id := range rootOrder {
		t := threadByRoot[id]
		if t.File != "" {
			review.FileThreads = append(review.FileThreads, *t)
		} else {
			review.GeneralThreads = append(review.GeneralThreads, *t)
		}
	}

	// Counts summary
	for _, c := range comments {
		if !c.Deleted {
			review.Counts.Comments++
		}
	}
	review.Counts.Threads = len(review.FileThreads) + len(review.GeneralThreads)
	review.Counts.Tasks = len(tasks)
	for _, t := range tasks {
		if t.State != "RESOLVED" {
			review.Counts.UnresolvedTasks++
		}
	}
	for _, th := range review.FileThreads {
		if th.Comment.Resolution == nil {
			review.Counts.UnresolvedThreads++
		}
	}
	for _, th := range review.GeneralThreads {
		if th.Comment.Resolution == nil {
			review.Counts.UnresolvedThreads++
		}
	}

	return review
}
