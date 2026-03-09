package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

type BranchRestriction struct {
	ID      int    `json:"id"`
	Kind    string `json:"kind"`
	Pattern string `json:"pattern"`
	Value   *int   `json:"value,omitempty"`
	Users   []struct {
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
	} `json:"users"`
	Groups []struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"groups"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"links"`
}

func (c *Client) ListBranchRestrictions(workspace, repoSlug string) ([]BranchRestriction, error) {
	path := fmt.Sprintf("/repositories/%s/%s/branch-restrictions",
		url.PathEscape(workspace), url.PathEscape(repoSlug))
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var page PaginatedResponse
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	var restrictions []BranchRestriction
	if err := json.Unmarshal(page.Values, &restrictions); err != nil {
		return nil, fmt.Errorf("parsing branch restrictions: %w", err)
	}
	return restrictions, nil
}

func (c *Client) GetBranchRestriction(workspace, repoSlug string, id int) (*BranchRestriction, error) {
	path := fmt.Sprintf("/repositories/%s/%s/branch-restrictions/%d",
		url.PathEscape(workspace), url.PathEscape(repoSlug), id)
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var restriction BranchRestriction
	if err := json.Unmarshal(data, &restriction); err != nil {
		return nil, fmt.Errorf("parsing branch restriction: %w", err)
	}
	return &restriction, nil
}

type CreateBranchRestrictionRequest struct {
	Kind    string `json:"kind"`
	Pattern string `json:"pattern"`
	Value   *int   `json:"value,omitempty"`
}

func (c *Client) CreateBranchRestriction(workspace, repoSlug string, req *CreateBranchRestrictionRequest) (*BranchRestriction, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/repositories/%s/%s/branch-restrictions",
		url.PathEscape(workspace), url.PathEscape(repoSlug))
	data, err := c.post(path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	var restriction BranchRestriction
	if err := json.Unmarshal(data, &restriction); err != nil {
		return nil, fmt.Errorf("parsing branch restriction: %w", err)
	}
	return &restriction, nil
}

func (c *Client) DeleteBranchRestriction(workspace, repoSlug string, id int) error {
	path := fmt.Sprintf("/repositories/%s/%s/branch-restrictions/%d",
		url.PathEscape(workspace), url.PathEscape(repoSlug), id)
	return c.deleteNoContent(path)
}
