package bitbucket

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type SearchResult struct {
	ContentMatchCount int `json:"content_match_count"`
	ContentMatches    []struct {
		Lines []struct {
			Line     int `json:"line"`
			Segments []struct {
				Text  string `json:"text"`
				Match bool   `json:"match"`
			} `json:"segments"`
		} `json:"lines"`
	} `json:"content_matches"`
	PathMatches []struct {
		Text  string `json:"text"`
		Match bool   `json:"match"`
	} `json:"path_matches"`
	File struct {
		Path string `json:"path"`
	} `json:"file"`
}

type SearchResponse struct {
	Size    int            `json:"size"`
	Page    int            `json:"page"`
	PageLen int            `json:"pagelen"`
	Next    string         `json:"next"`
	Values  []SearchResult `json:"values"`
}

func (c *Client) SearchCode(workspace, query string) (*SearchResponse, error) {
	params := url.Values{}
	params.Set("search_query", query)
	path := fmt.Sprintf("/workspaces/%s/search/code?%s",
		url.PathEscape(workspace), params.Encode())
	data, err := c.get(path)
	if err != nil {
		return nil, err
	}
	var result SearchResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("parsing search results: %w", err)
	}
	return &result, nil
}
