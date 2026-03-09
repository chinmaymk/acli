package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

type Environment struct {
	UUID            string `json:"uuid"`
	Name            string `json:"name"`
	Slug            string `json:"slug"`
	EnvironmentType struct {
		Name string `json:"name"`
		Rank int    `json:"rank"`
	} `json:"environment_type"`
	Rank           int `json:"rank"`
	DeploymentGate struct {
		Name string `json:"name"`
	} `json:"deployment_gate"`
	Lock struct {
		Name string `json:"name"`
	} `json:"lock"`
	Restrictions struct {
		AdminOnly bool   `json:"admin_only"`
		Type      string `json:"type"`
	} `json:"restrictions"`
}

func (c *Client) ListEnvironments(workspace, repoSlug string) ([]Environment, error) {
	path := fmt.Sprintf("/repositories/%s/%s/environments",
		url.PathEscape(workspace), url.PathEscape(repoSlug))
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var page PaginatedResponse
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	var envs []Environment
	if err := json.Unmarshal(page.Values, &envs); err != nil {
		return nil, fmt.Errorf("parsing environments: %w", err)
	}
	return envs, nil
}

func (c *Client) GetEnvironment(workspace, repoSlug, envUUID string) (*Environment, error) {
	path := fmt.Sprintf("/repositories/%s/%s/environments/%s",
		url.PathEscape(workspace), url.PathEscape(repoSlug), url.PathEscape(envUUID))
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var env Environment
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parsing environment: %w", err)
	}
	return &env, nil
}

type CreateEnvironmentRequest struct {
	Name            string `json:"name"`
	EnvironmentType struct {
		Name string `json:"name"`
		Rank int    `json:"rank"`
	} `json:"environment_type"`
}

func (c *Client) CreateEnvironment(workspace, repoSlug string, req *CreateEnvironmentRequest) (*Environment, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/repositories/%s/%s/environments",
		url.PathEscape(workspace), url.PathEscape(repoSlug))
	data, err := c.post(path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	var env Environment
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parsing environment: %w", err)
	}
	return &env, nil
}

func (c *Client) DeleteEnvironment(workspace, repoSlug, envUUID string) error {
	path := fmt.Sprintf("/repositories/%s/%s/environments/%s",
		url.PathEscape(workspace), url.PathEscape(repoSlug), url.PathEscape(envUUID))
	return c.deleteNoContent(path)
}
