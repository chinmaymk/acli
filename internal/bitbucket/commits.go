package bitbucket

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Commit struct {
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
	Parents []struct {
		Hash string `json:"hash"`
	} `json:"parents"`
	Links struct {
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
	} `json:"links"`
	Repository struct {
		FullName string `json:"full_name"`
		UUID     string `json:"uuid"`
	} `json:"repository"`
}

func (c *Client) ListCommits(workspace, repoSlug string, include, exclude string) ([]Commit, error) {
	params := url.Values{}
	if include != "" {
		params.Set("include", include)
	}
	if exclude != "" {
		params.Set("exclude", exclude)
	}
	path := fmt.Sprintf("/repositories/%s/%s/commits",
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
	var commits []Commit
	if err := json.Unmarshal(page.Values, &commits); err != nil {
		return nil, fmt.Errorf("parsing commits: %w", err)
	}
	return commits, nil
}

func (c *Client) GetCommit(workspace, repoSlug, commitHash string) (*Commit, error) {
	path := fmt.Sprintf("/repositories/%s/%s/commit/%s",
		url.PathEscape(workspace), url.PathEscape(repoSlug), url.PathEscape(commitHash))
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var commit Commit
	if err := json.Unmarshal(data, &commit); err != nil {
		return nil, fmt.Errorf("parsing commit: %w", err)
	}
	return &commit, nil
}

type CommitStatus struct {
	UUID        string `json:"uuid"`
	Key         string `json:"key"`
	State       string `json:"state"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	CreatedOn   string `json:"created_on"`
	UpdatedOn   string `json:"updated_on"`
	Refname     string `json:"refname"`
}

func (c *Client) ListCommitStatuses(workspace, repoSlug, commitHash string) ([]CommitStatus, error) {
	path := fmt.Sprintf("/repositories/%s/%s/commit/%s/statuses",
		url.PathEscape(workspace), url.PathEscape(repoSlug), url.PathEscape(commitHash))
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var page PaginatedResponse
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	var statuses []CommitStatus
	if err := json.Unmarshal(page.Values, &statuses); err != nil {
		return nil, fmt.Errorf("parsing statuses: %w", err)
	}
	return statuses, nil
}

func (c *Client) GetDiff(workspace, repoSlug, spec string) (string, error) {
	path := fmt.Sprintf("/repositories/%s/%s/diff/%s",
		url.PathEscape(workspace), url.PathEscape(repoSlug), url.PathEscape(spec))
	data, err := c.getRaw(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
