package bitbucket

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Download struct {
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	CreatedOn string `json:"created_on"`
	Downloads int    `json:"downloads"`
	Links     struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"links"`
}

func (c *Client) ListDownloads(workspace, repoSlug string) ([]Download, error) {
	path := fmt.Sprintf("/repositories/%s/%s/downloads",
		url.PathEscape(workspace), url.PathEscape(repoSlug))
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var page PaginatedResponse
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	var downloads []Download
	if err := json.Unmarshal(page.Values, &downloads); err != nil {
		return nil, fmt.Errorf("parsing downloads: %w", err)
	}
	return downloads, nil
}

func (c *Client) DeleteDownload(workspace, repoSlug, filename string) error {
	path := fmt.Sprintf("/repositories/%s/%s/downloads/%s",
		url.PathEscape(workspace), url.PathEscape(repoSlug), url.PathEscape(filename))
	return c.deleteNoContent(path)
}
