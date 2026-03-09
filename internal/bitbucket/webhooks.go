package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

type Webhook struct {
	UUID        string   `json:"uuid"`
	URL         string   `json:"url"`
	Description string   `json:"description"`
	Active      bool     `json:"active"`
	CreatedAt   string   `json:"created_at"`
	Events      []string `json:"events"`
	Subject     struct {
		Type     string `json:"type"`
		FullName string `json:"full_name"`
	} `json:"subject"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"links"`
}

type CreateWebhookRequest struct {
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Active      bool     `json:"active"`
	Events      []string `json:"events"`
}

func (c *Client) ListRepoWebhooks(workspace, repoSlug string) ([]Webhook, error) {
	path := fmt.Sprintf("/repositories/%s/%s/hooks",
		url.PathEscape(workspace), url.PathEscape(repoSlug))
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var page PaginatedResponse
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	var hooks []Webhook
	if err := json.Unmarshal(page.Values, &hooks); err != nil {
		return nil, fmt.Errorf("parsing webhooks: %w", err)
	}
	return hooks, nil
}

func (c *Client) GetRepoWebhook(workspace, repoSlug, uid string) (*Webhook, error) {
	path := fmt.Sprintf("/repositories/%s/%s/hooks/%s",
		url.PathEscape(workspace), url.PathEscape(repoSlug), url.PathEscape(uid))
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var hook Webhook
	if err := json.Unmarshal(data, &hook); err != nil {
		return nil, fmt.Errorf("parsing webhook: %w", err)
	}
	return &hook, nil
}

func (c *Client) CreateRepoWebhook(workspace, repoSlug string, req *CreateWebhookRequest) (*Webhook, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/repositories/%s/%s/hooks",
		url.PathEscape(workspace), url.PathEscape(repoSlug))
	data, err := c.post(path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	var hook Webhook
	if err := json.Unmarshal(data, &hook); err != nil {
		return nil, fmt.Errorf("parsing webhook: %w", err)
	}
	return &hook, nil
}

func (c *Client) UpdateRepoWebhook(workspace, repoSlug, uid string, req *CreateWebhookRequest) (*Webhook, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/repositories/%s/%s/hooks/%s",
		url.PathEscape(workspace), url.PathEscape(repoSlug), url.PathEscape(uid))
	data, err := c.put(path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	var hook Webhook
	if err := json.Unmarshal(data, &hook); err != nil {
		return nil, fmt.Errorf("parsing webhook: %w", err)
	}
	return &hook, nil
}

func (c *Client) DeleteRepoWebhook(workspace, repoSlug, uid string) error {
	path := fmt.Sprintf("/repositories/%s/%s/hooks/%s",
		url.PathEscape(workspace), url.PathEscape(repoSlug), url.PathEscape(uid))
	return c.deleteNoContent(path)
}

func (c *Client) ListWorkspaceWebhooks(workspace string) ([]Webhook, error) {
	path := fmt.Sprintf("/workspaces/%s/hooks", url.PathEscape(workspace))
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var page PaginatedResponse
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	var hooks []Webhook
	if err := json.Unmarshal(page.Values, &hooks); err != nil {
		return nil, fmt.Errorf("parsing webhooks: %w", err)
	}
	return hooks, nil
}

func (c *Client) CreateWorkspaceWebhook(workspace string, req *CreateWebhookRequest) (*Webhook, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/workspaces/%s/hooks", url.PathEscape(workspace))
	data, err := c.post(path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	var hook Webhook
	if err := json.Unmarshal(data, &hook); err != nil {
		return nil, fmt.Errorf("parsing webhook: %w", err)
	}
	return &hook, nil
}

func (c *Client) DeleteWorkspaceWebhook(workspace, uid string) error {
	path := fmt.Sprintf("/workspaces/%s/hooks/%s",
		url.PathEscape(workspace), url.PathEscape(uid))
	return c.deleteNoContent(path)
}
