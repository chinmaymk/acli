package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
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
	State string
}

func (c *Client) ListPullRequests(workspace, repoSlug string, opts *ListPRsOptions) ([]PullRequest, error) {
	params := url.Values{}
	if opts != nil {
		if opts.State != "" {
			params.Set("state", opts.State)
		}
	}

	path := fmt.Sprintf("/repositories/%s/%s/pullrequests",
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

func (c *Client) ApprovePullRequest(workspace, repoSlug string, prID int) error {
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/approve",
		url.PathEscape(workspace), url.PathEscape(repoSlug), prID)
	return c.postNoContent(path, nil)
}

func (c *Client) UnapprovePullRequest(workspace, repoSlug string, prID int) error {
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/approve",
		url.PathEscape(workspace), url.PathEscape(repoSlug), prID)
	return c.deleteNoContent(path)
}

func (c *Client) DeclinePullRequest(workspace, repoSlug string, prID int) error {
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/decline",
		url.PathEscape(workspace), url.PathEscape(repoSlug), prID)
	return c.postNoContent(path, nil)
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

func (c *Client) RequestChangesPullRequest(workspace, repoSlug string, prID int) error {
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/request-changes",
		url.PathEscape(workspace), url.PathEscape(repoSlug), prID)
	return c.postNoContent(path, nil)
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
		Path string `json:"path"`
		From *int   `json:"from"`
		To   *int   `json:"to"`
	} `json:"inline,omitempty"`
	Parent *struct {
		ID int `json:"id"`
	} `json:"parent,omitempty"`
}

func (c *Client) ListPRComments(workspace, repoSlug string, prID int) ([]PRComment, error) {
	path := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/comments",
		url.PathEscape(workspace), url.PathEscape(repoSlug), prID)
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

func (c *Client) CreatePRComment(workspace, repoSlug string, prID int, content string) (*PRComment, error) {
	body := map[string]interface{}{
		"content": map[string]string{
			"raw": content,
		},
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
	data, err := c.get(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
