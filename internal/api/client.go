package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	BaseURL    string
	Email      string
	APIToken   string
	HTTPClient *http.Client
}

func NewClient(baseURL, email, apiToken string) *Client {
	return &Client{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		Email:      email,
		APIToken:   apiToken,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) ConfluenceV2(method, path string, query url.Values, body interface{}) ([]byte, error) {
	endpoint := c.BaseURL + "/wiki/api/v2" + path
	if query != nil && len(query) > 0 {
		endpoint += "?" + query.Encode()
	}

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshalling request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.SetBasicAuth(c.Email, c.APIToken)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
