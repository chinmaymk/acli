package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

type DeployKey struct {
	ID        int    `json:"id"`
	Key       string `json:"key"`
	Label     string `json:"label"`
	Comment   string `json:"comment"`
	CreatedOn string `json:"created_on"`
	LastUsed  string `json:"last_used"`
	Owner     struct {
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
	} `json:"owner"`
	Repository struct {
		FullName string `json:"full_name"`
		UUID     string `json:"uuid"`
	} `json:"repository"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"links"`
}

func (c *Client) ListDeployKeys(workspace, repoSlug string) ([]DeployKey, error) {
	path := fmt.Sprintf("/repositories/%s/%s/deploy-keys",
		url.PathEscape(workspace), url.PathEscape(repoSlug))
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var page PaginatedResponse
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	var keys []DeployKey
	if err := json.Unmarshal(page.Values, &keys); err != nil {
		return nil, fmt.Errorf("parsing deploy keys: %w", err)
	}
	return keys, nil
}

func (c *Client) GetDeployKey(workspace, repoSlug string, keyID int) (*DeployKey, error) {
	path := fmt.Sprintf("/repositories/%s/%s/deploy-keys/%d",
		url.PathEscape(workspace), url.PathEscape(repoSlug), keyID)
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var key DeployKey
	if err := json.Unmarshal(data, &key); err != nil {
		return nil, fmt.Errorf("parsing deploy key: %w", err)
	}
	return &key, nil
}

type CreateDeployKeyRequest struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

func (c *Client) CreateDeployKey(workspace, repoSlug string, req *CreateDeployKeyRequest) (*DeployKey, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/repositories/%s/%s/deploy-keys",
		url.PathEscape(workspace), url.PathEscape(repoSlug))
	data, err := c.post(path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	var key DeployKey
	if err := json.Unmarshal(data, &key); err != nil {
		return nil, fmt.Errorf("parsing deploy key: %w", err)
	}
	return &key, nil
}

func (c *Client) DeleteDeployKey(workspace, repoSlug string, keyID int) error {
	path := fmt.Sprintf("/repositories/%s/%s/deploy-keys/%d",
		url.PathEscape(workspace), url.PathEscape(repoSlug), keyID)
	return c.deleteNoContent(path)
}
