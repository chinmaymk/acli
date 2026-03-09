package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

type Branch struct {
	Name   string `json:"name"`
	Target struct {
		Hash    string `json:"hash"`
		Date    string `json:"date"`
		Message string `json:"message"`
		Author  struct {
			Raw  string `json:"raw"`
			User struct {
				DisplayName string `json:"display_name"`
				UUID        string `json:"uuid"`
			} `json:"user"`
		} `json:"author"`
	} `json:"target"`
	Links struct {
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
	} `json:"links"`
}

type Tag struct {
	Name   string `json:"name"`
	Target struct {
		Hash    string `json:"hash"`
		Date    string `json:"date"`
		Message string `json:"message"`
		Author  struct {
			Raw  string `json:"raw"`
			User struct {
				DisplayName string `json:"display_name"`
				UUID        string `json:"uuid"`
			} `json:"user"`
		} `json:"author"`
	} `json:"target"`
	Message string `json:"message"`
	Tagger  struct {
		Raw  string `json:"raw"`
		User struct {
			DisplayName string `json:"display_name"`
			UUID        string `json:"uuid"`
		} `json:"user"`
	} `json:"tagger"`
	Links struct {
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
	} `json:"links"`
}

func (c *Client) ListBranches(workspace, repoSlug string, q string) ([]Branch, error) {
	params := url.Values{}
	if q != "" {
		params.Set("q", q)
	}
	path := fmt.Sprintf("/repositories/%s/%s/refs/branches",
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
	var branches []Branch
	if err := json.Unmarshal(page.Values, &branches); err != nil {
		return nil, fmt.Errorf("parsing branches: %w", err)
	}
	return branches, nil
}

func (c *Client) GetBranch(workspace, repoSlug, name string) (*Branch, error) {
	path := fmt.Sprintf("/repositories/%s/%s/refs/branches/%s",
		url.PathEscape(workspace), url.PathEscape(repoSlug), url.PathEscape(name))
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var branch Branch
	if err := json.Unmarshal(data, &branch); err != nil {
		return nil, fmt.Errorf("parsing branch: %w", err)
	}
	return &branch, nil
}

type CreateBranchRequest struct {
	Name   string `json:"name"`
	Target struct {
		Hash string `json:"hash"`
	} `json:"target"`
}

func (c *Client) CreateBranch(workspace, repoSlug string, req *CreateBranchRequest) (*Branch, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/repositories/%s/%s/refs/branches",
		url.PathEscape(workspace), url.PathEscape(repoSlug))
	data, err := c.post(path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	var branch Branch
	if err := json.Unmarshal(data, &branch); err != nil {
		return nil, fmt.Errorf("parsing branch: %w", err)
	}
	return &branch, nil
}

func (c *Client) DeleteBranch(workspace, repoSlug, name string) error {
	path := fmt.Sprintf("/repositories/%s/%s/refs/branches/%s",
		url.PathEscape(workspace), url.PathEscape(repoSlug), url.PathEscape(name))
	return c.deleteNoContent(path)
}

func (c *Client) ListTags(workspace, repoSlug string, q string) ([]Tag, error) {
	params := url.Values{}
	if q != "" {
		params.Set("q", q)
	}
	path := fmt.Sprintf("/repositories/%s/%s/refs/tags",
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
	var tags []Tag
	if err := json.Unmarshal(page.Values, &tags); err != nil {
		return nil, fmt.Errorf("parsing tags: %w", err)
	}
	return tags, nil
}

func (c *Client) GetTag(workspace, repoSlug, name string) (*Tag, error) {
	path := fmt.Sprintf("/repositories/%s/%s/refs/tags/%s",
		url.PathEscape(workspace), url.PathEscape(repoSlug), url.PathEscape(name))
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var tag Tag
	if err := json.Unmarshal(data, &tag); err != nil {
		return nil, fmt.Errorf("parsing tag: %w", err)
	}
	return &tag, nil
}

type CreateTagRequest struct {
	Name   string `json:"name"`
	Target struct {
		Hash string `json:"hash"`
	} `json:"target"`
	Message string `json:"message,omitempty"`
}

func (c *Client) CreateTag(workspace, repoSlug string, req *CreateTagRequest) (*Tag, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/repositories/%s/%s/refs/tags",
		url.PathEscape(workspace), url.PathEscape(repoSlug))
	data, err := c.post(path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	var tag Tag
	if err := json.Unmarshal(data, &tag); err != nil {
		return nil, fmt.Errorf("parsing tag: %w", err)
	}
	return &tag, nil
}

func (c *Client) DeleteTag(workspace, repoSlug, name string) error {
	path := fmt.Sprintf("/repositories/%s/%s/refs/tags/%s",
		url.PathEscape(workspace), url.PathEscape(repoSlug), url.PathEscape(name))
	return c.deleteNoContent(path)
}
