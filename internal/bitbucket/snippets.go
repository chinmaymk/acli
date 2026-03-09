package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

type Snippet struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Scm       string `json:"scm"`
	IsPrivate bool   `json:"is_private"`
	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
	Owner     struct {
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
	} `json:"owner"`
	Creator struct {
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
	} `json:"creator"`
	Links struct {
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
	} `json:"links"`
}

func (c *Client) ListSnippets(workspace string) ([]Snippet, error) {
	path := fmt.Sprintf("/snippets/%s", url.PathEscape(workspace))
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var page PaginatedResponse
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	var snippets []Snippet
	if err := json.Unmarshal(page.Values, &snippets); err != nil {
		return nil, fmt.Errorf("parsing snippets: %w", err)
	}
	return snippets, nil
}

func (c *Client) GetSnippet(workspace, encodedID string) (*Snippet, error) {
	path := fmt.Sprintf("/snippets/%s/%s",
		url.PathEscape(workspace), url.PathEscape(encodedID))
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var snippet Snippet
	if err := json.Unmarshal(data, &snippet); err != nil {
		return nil, fmt.Errorf("parsing snippet: %w", err)
	}
	return &snippet, nil
}

type CreateSnippetRequest struct {
	Title     string `json:"title"`
	IsPrivate bool   `json:"is_private,omitempty"`
	Scm       string `json:"scm,omitempty"`
	Files     map[string]struct {
		Content string `json:"content"`
	} `json:"files,omitempty"`
}

func (c *Client) CreateSnippet(workspace string, req *CreateSnippetRequest) (*Snippet, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/snippets/%s", url.PathEscape(workspace))
	data, err := c.post(path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	var snippet Snippet
	if err := json.Unmarshal(data, &snippet); err != nil {
		return nil, fmt.Errorf("parsing snippet: %w", err)
	}
	return &snippet, nil
}

func (c *Client) DeleteSnippet(workspace, encodedID string) error {
	path := fmt.Sprintf("/snippets/%s/%s",
		url.PathEscape(workspace), url.PathEscape(encodedID))
	return c.deleteNoContent(path)
}
