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
	Summary *struct {
		Raw    string `json:"raw"`
		Markup string `json:"markup"`
		HTML   string `json:"html"`
	} `json:"summary,omitempty"`
	Parents []struct {
		Hash string `json:"hash"`
	} `json:"parents"`
	Links struct {
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
	} `json:"links"`
	// Repository is present on full commit objects but not on base_commit (e.g. PR commits).
	Repository struct {
		FullName string `json:"full_name"`
		UUID     string `json:"uuid"`
	} `json:"repository"`
}

type ListCommitsOptions struct {
	Include string
	Exclude string
	PaginationOptions
}

func (c *Client) ListCommits(workspace, repoSlug string, opts *ListCommitsOptions) ([]Commit, error) {
	params := url.Values{}
	if opts != nil {
		if opts.Include != "" {
			params.Set("include", opts.Include)
		}
		if opts.Exclude != "" {
			params.Set("exclude", opts.Exclude)
		}
		opts.applyParams(params)
	}
	ensurePageLen(params)

	path := fmt.Sprintf("/repositories/%s/%s/commits",
		url.PathEscape(workspace), url.PathEscape(repoSlug))
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

func (c *Client) ListCommitStatuses(workspace, repoSlug, commitHash string, opts *PaginationOptions) ([]CommitStatus, error) {
	params := url.Values{}
	if opts != nil {
		opts.applyParams(params)
	}
	ensurePageLen(params)

	path := fmt.Sprintf("/repositories/%s/%s/commit/%s/statuses",
		url.PathEscape(workspace), url.PathEscape(repoSlug), url.PathEscape(commitHash))
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	if opts != nil && opts.All {
		pages, err := c.getAll(path)
		if err != nil && len(pages) == 0 {
			return nil, err
		}
		var statuses []CommitStatus
		for _, pg := range pages {
			var pageStatuses []CommitStatus
			if err := json.Unmarshal(pg.Values, &pageStatuses); err != nil {
				return statuses, fmt.Errorf("parsing statuses: %w", err)
			}
			statuses = append(statuses, pageStatuses...)
		}
		return statuses, nil
	}

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
