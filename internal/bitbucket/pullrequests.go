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
	Inline *struct {
		Path      string `json:"path"`
		From      *int   `json:"from"`
		To        *int   `json:"to"`
		StartFrom *int   `json:"start_from,omitempty"`
		StartTo   *int   `json:"start_to,omitempty"`
	} `json:"inline,omitempty"`
	Parent *struct {
		ID int `json:"id"`
	} `json:"parent,omitempty"`
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

type PRTask struct {
	ID        int    `json:"id"`
	State     string `json:"state"`
	Content   struct {
		Raw    string `json:"raw"`
		Markup string `json:"markup"`
		HTML   string `json:"html"`
	} `json:"content"`
	Creator struct {
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
	} `json:"creator"`
	CreatedOn  string `json:"created_on"`
	UpdatedOn  string `json:"updated_on"`
	ResolvedOn string `json:"resolved_on,omitempty"`
	ResolvedBy *struct {
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
	} `json:"resolved_by,omitempty"`
	Comment *struct {
		ID int `json:"id"`
	} `json:"comment,omitempty"`
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
