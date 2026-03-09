package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

type Project struct {
	UUID        string `json:"uuid"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"is_private"`
	CreatedOn   string `json:"created_on"`
	UpdatedOn   string `json:"updated_on"`
	Owner       struct {
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
	} `json:"owner"`
	Links struct {
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
		Avatar struct {
			Href string `json:"href"`
		} `json:"avatar"`
	} `json:"links"`
}

func (c *Client) ListProjects(workspace string) ([]Project, error) {
	path := fmt.Sprintf("/workspaces/%s/projects", url.PathEscape(workspace))
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var page PaginatedResponse
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	var projects []Project
	if err := json.Unmarshal(page.Values, &projects); err != nil {
		return nil, fmt.Errorf("parsing projects: %w", err)
	}
	return projects, nil
}

func (c *Client) GetProject(workspace, projectKey string) (*Project, error) {
	path := fmt.Sprintf("/workspaces/%s/projects/%s",
		url.PathEscape(workspace), url.PathEscape(projectKey))
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var project Project
	if err := json.Unmarshal(data, &project); err != nil {
		return nil, fmt.Errorf("parsing project: %w", err)
	}
	return &project, nil
}

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Key         string `json:"key"`
	Description string `json:"description,omitempty"`
	IsPrivate   bool   `json:"is_private,omitempty"`
}

func (c *Client) CreateProject(workspace string, req *CreateProjectRequest) (*Project, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/workspaces/%s/projects", url.PathEscape(workspace))
	data, err := c.post(path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	var project Project
	if err := json.Unmarshal(data, &project); err != nil {
		return nil, fmt.Errorf("parsing project: %w", err)
	}
	return &project, nil
}

func (c *Client) DeleteProject(workspace, projectKey string) error {
	path := fmt.Sprintf("/workspaces/%s/projects/%s",
		url.PathEscape(workspace), url.PathEscape(projectKey))
	return c.deleteNoContent(path)
}
