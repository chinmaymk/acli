package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

type Repository struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"is_private"`
	Language    string `json:"language"`
	CreatedOn   string `json:"created_on"`
	UpdatedOn   string `json:"updated_on"`
	SCM         string `json:"scm"`
	MainBranch  *struct {
		Name string `json:"name"`
	} `json:"mainbranch"`
	Links struct {
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
		Clone []struct {
			Href string `json:"href"`
			Name string `json:"name"`
		} `json:"clone"`
	} `json:"links"`
	Owner struct {
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
	} `json:"owner"`
}

type ListReposOptions struct {
	Role string
	Q    string
	Sort string
}

func (c *Client) ListRepositories(workspace string, opts *ListReposOptions) ([]Repository, error) {
	params := url.Values{}
	if opts != nil {
		if opts.Role != "" {
			params.Set("role", opts.Role)
		}
		if opts.Q != "" {
			params.Set("q", opts.Q)
		}
		if opts.Sort != "" {
			params.Set("sort", opts.Sort)
		}
	}

	path := fmt.Sprintf("/repositories/%s", url.PathEscape(workspace))
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

	var repos []Repository
	if err := json.Unmarshal(page.Values, &repos); err != nil {
		return nil, fmt.Errorf("parsing repositories: %w", err)
	}

	return repos, nil
}

func (c *Client) GetRepository(workspace, repoSlug string) (*Repository, error) {
	path := fmt.Sprintf("/repositories/%s/%s", url.PathEscape(workspace), url.PathEscape(repoSlug))

	data, err := c.get(path)
	if err != nil {
		return nil, err
	}

	var repo Repository
	if err := json.Unmarshal(data, &repo); err != nil {
		return nil, fmt.Errorf("parsing repository: %w", err)
	}

	return &repo, nil
}

type CreateRepoRequest struct {
	SCM         string `json:"scm"`
	Name        string `json:"name"`
	IsPrivate   bool   `json:"is_private"`
	Description string `json:"description,omitempty"`
	Language    string `json:"language,omitempty"`
	HasIssues   bool   `json:"has_issues,omitempty"`
	HasWiki     bool   `json:"has_wiki,omitempty"`
	ForkPolicy  string `json:"fork_policy,omitempty"`
	Project     *struct {
		Key string `json:"key"`
	} `json:"project,omitempty"`
}

func (c *Client) CreateRepository(workspace string, req *CreateRepoRequest) (*Repository, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/repositories/%s/%s", url.PathEscape(workspace), url.PathEscape(req.Name))
	data, err := c.post(path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	var repo Repository
	if err := json.Unmarshal(data, &repo); err != nil {
		return nil, fmt.Errorf("parsing repository: %w", err)
	}
	return &repo, nil
}

func (c *Client) DeleteRepository(workspace, repoSlug string) error {
	path := fmt.Sprintf("/repositories/%s/%s", url.PathEscape(workspace), url.PathEscape(repoSlug))
	return c.deleteNoContent(path)
}

type ForkRepoRequest struct {
	Name      string `json:"name,omitempty"`
	Workspace *struct {
		Slug string `json:"slug"`
	} `json:"workspace,omitempty"`
	IsPrivate   *bool  `json:"is_private,omitempty"`
	Description string `json:"description,omitempty"`
	Language    string `json:"language,omitempty"`
}

func (c *Client) ForkRepository(workspace, repoSlug string, req *ForkRepoRequest) (*Repository, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/repositories/%s/%s/forks", url.PathEscape(workspace), url.PathEscape(repoSlug))
	data, err := c.post(path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	var repo Repository
	if err := json.Unmarshal(data, &repo); err != nil {
		return nil, fmt.Errorf("parsing repository: %w", err)
	}
	return &repo, nil
}

func (c *Client) ListForks(workspace, repoSlug string) ([]Repository, error) {
	path := fmt.Sprintf("/repositories/%s/%s/forks", url.PathEscape(workspace), url.PathEscape(repoSlug))
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var page PaginatedResponse
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	var repos []Repository
	if err := json.Unmarshal(page.Values, &repos); err != nil {
		return nil, fmt.Errorf("parsing forks: %w", err)
	}
	return repos, nil
}
