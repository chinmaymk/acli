package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

type Issue struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	State    string `json:"state"`
	Priority string `json:"priority"`
	Kind     string `json:"kind"`
	Content  struct {
		Raw    string `json:"raw"`
		Markup string `json:"markup"`
		HTML   string `json:"html"`
	} `json:"content"`
	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
	Reporter  struct {
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
	} `json:"reporter"`
	Assignee *struct {
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
	} `json:"assignee"`
	Component *struct {
		Name string `json:"name"`
	} `json:"component"`
	Milestone *struct {
		Name string `json:"name"`
	} `json:"milestone"`
	Version *struct {
		Name string `json:"name"`
	} `json:"version"`
	Votes int `json:"votes"`
	Links struct {
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
	} `json:"links"`
}

type ListIssuesOptions struct {
	Q    string
	Sort string
}

func (c *Client) ListIssues(workspace, repoSlug string, opts *ListIssuesOptions) ([]Issue, error) {
	params := url.Values{}
	if opts != nil {
		if opts.Q != "" {
			params.Set("q", opts.Q)
		}
		if opts.Sort != "" {
			params.Set("sort", opts.Sort)
		}
	}
	path := fmt.Sprintf("/repositories/%s/%s/issues",
		url.PathEscape(workspace), url.PathEscape(repoSlug))
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var page PaginatedResponse
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	var issues []Issue
	if err := json.Unmarshal(page.Values, &issues); err != nil {
		return nil, fmt.Errorf("parsing issues: %w", err)
	}
	return issues, nil
}

func (c *Client) GetIssue(workspace, repoSlug string, issueID int) (*Issue, error) {
	path := fmt.Sprintf("/repositories/%s/%s/issues/%d",
		url.PathEscape(workspace), url.PathEscape(repoSlug), issueID)
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var issue Issue
	if err := json.Unmarshal(data, &issue); err != nil {
		return nil, fmt.Errorf("parsing issue: %w", err)
	}
	return &issue, nil
}

type CreateIssueRequest struct {
	Title   string `json:"title"`
	Content *struct {
		Raw string `json:"raw"`
	} `json:"content,omitempty"`
	Kind     string `json:"kind,omitempty"`
	Priority string `json:"priority,omitempty"`
	State    string `json:"state,omitempty"`
	Assignee *struct {
		UUID string `json:"uuid"`
	} `json:"assignee,omitempty"`
	Component *struct {
		Name string `json:"name"`
	} `json:"component,omitempty"`
	Milestone *struct {
		Name string `json:"name"`
	} `json:"milestone,omitempty"`
	Version *struct {
		Name string `json:"name"`
	} `json:"version,omitempty"`
}

func (c *Client) CreateIssue(workspace, repoSlug string, req *CreateIssueRequest) (*Issue, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/repositories/%s/%s/issues",
		url.PathEscape(workspace), url.PathEscape(repoSlug))
	data, err := c.post(path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	var issue Issue
	if err := json.Unmarshal(data, &issue); err != nil {
		return nil, fmt.Errorf("parsing issue: %w", err)
	}
	return &issue, nil
}

type UpdateIssueRequest struct {
	Title   string `json:"title,omitempty"`
	Content *struct {
		Raw string `json:"raw"`
	} `json:"content,omitempty"`
	Kind     string `json:"kind,omitempty"`
	Priority string `json:"priority,omitempty"`
	State    string `json:"state,omitempty"`
	Assignee *struct {
		UUID string `json:"uuid"`
	} `json:"assignee,omitempty"`
}

func (c *Client) UpdateIssue(workspace, repoSlug string, issueID int, req *UpdateIssueRequest) (*Issue, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/repositories/%s/%s/issues/%d",
		url.PathEscape(workspace), url.PathEscape(repoSlug), issueID)
	data, err := c.put(path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	var issue Issue
	if err := json.Unmarshal(data, &issue); err != nil {
		return nil, fmt.Errorf("parsing issue: %w", err)
	}
	return &issue, nil
}

func (c *Client) DeleteIssue(workspace, repoSlug string, issueID int) error {
	path := fmt.Sprintf("/repositories/%s/%s/issues/%d",
		url.PathEscape(workspace), url.PathEscape(repoSlug), issueID)
	return c.deleteNoContent(path)
}

type IssueComment struct {
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
}

func (c *Client) ListIssueComments(workspace, repoSlug string, issueID int) ([]IssueComment, error) {
	path := fmt.Sprintf("/repositories/%s/%s/issues/%d/comments",
		url.PathEscape(workspace), url.PathEscape(repoSlug), issueID)
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var page PaginatedResponse
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	var comments []IssueComment
	if err := json.Unmarshal(page.Values, &comments); err != nil {
		return nil, fmt.Errorf("parsing comments: %w", err)
	}
	return comments, nil
}

func (c *Client) CreateIssueComment(workspace, repoSlug string, issueID int, content string) (*IssueComment, error) {
	body := map[string]interface{}{
		"content": map[string]string{
			"raw": content,
		},
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/repositories/%s/%s/issues/%d/comments",
		url.PathEscape(workspace), url.PathEscape(repoSlug), issueID)
	data, err := c.post(path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	var comment IssueComment
	if err := json.Unmarshal(data, &comment); err != nil {
		return nil, fmt.Errorf("parsing comment: %w", err)
	}
	return &comment, nil
}
