package bitbucket

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Deployment struct {
	UUID  string `json:"uuid"`
	State struct {
		Name   string `json:"name"`
		Status struct {
			Name string `json:"name"`
		} `json:"status"`
	} `json:"state"`
	Environment struct {
		UUID string `json:"uuid"`
		Name string `json:"name"`
	} `json:"environment"`
	Release struct {
		UUID   string `json:"uuid"`
		Name   string `json:"name"`
		URL    string `json:"url"`
		Commit struct {
			Hash    string `json:"hash"`
			Message string `json:"message"`
		} `json:"commit"`
		CreatedOn string `json:"created_on"`
	} `json:"release"`
	Step struct {
		UUID string `json:"uuid"`
	} `json:"step"`
}

func (c *Client) ListDeployments(workspace, repoSlug string) ([]Deployment, error) {
	path := fmt.Sprintf("/repositories/%s/%s/deployments",
		url.PathEscape(workspace), url.PathEscape(repoSlug))
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var page PaginatedResponse
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	var deployments []Deployment
	if err := json.Unmarshal(page.Values, &deployments); err != nil {
		return nil, fmt.Errorf("parsing deployments: %w", err)
	}
	return deployments, nil
}

func (c *Client) GetDeployment(workspace, repoSlug, deploymentUUID string) (*Deployment, error) {
	path := fmt.Sprintf("/repositories/%s/%s/deployments/%s",
		url.PathEscape(workspace), url.PathEscape(repoSlug), url.PathEscape(deploymentUUID))
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var deployment Deployment
	if err := json.Unmarshal(data, &deployment); err != nil {
		return nil, fmt.Errorf("parsing deployment: %w", err)
	}
	return &deployment, nil
}
