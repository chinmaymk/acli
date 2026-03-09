package bitbucket

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/chinmaymk/acli/internal/config"
)

const baseURL = "https://api.bitbucket.org/2.0"

type Client struct {
	httpClient *http.Client
	token      string
	email      string // if set, use Basic auth (email:token)
}

// NewClient creates a Bitbucket client using profile credentials.
func NewClient(profile config.Profile) (*Client, error) {
	if profile.APIToken == "" {
		return nil, fmt.Errorf("no API token configured: run 'acli config setup' to set one")
	}

	return &Client{
		httpClient: &http.Client{},
		token:      profile.APIToken,
		email:      profile.Email,
	}, nil
}

func (c *Client) do(method, path string, body io.Reader) ([]byte, error) {
	url := baseURL + path
	if strings.HasPrefix(path, "http") {
		url = path
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if c.email != "" {
		req.SetBasicAuth(c.email, c.token)
	} else {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		var apiErr APIError
		if json.Unmarshal(data, &apiErr) == nil && apiErr.Error.Message != "" {
			return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, apiErr.Error.Message)
		}
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(data))
	}

	return data, nil
}

func (c *Client) get(path string) ([]byte, error) {
	return c.do(http.MethodGet, path, nil)
}

func (c *Client) post(path string, body io.Reader) ([]byte, error) {
	return c.do(http.MethodPost, path, body)
}

func (c *Client) put(path string, body io.Reader) ([]byte, error) {
	return c.do(http.MethodPut, path, body)
}

// deleteNoContent is for DELETE endpoints that return 204 No Content
func (c *Client) deleteNoContent(path string) error {
	_, err := c.do(http.MethodDelete, path, nil)
	return err
}

// postNoContent is for POST endpoints that return 204 No Content
func (c *Client) postNoContent(path string, body io.Reader) error {
	_, err := c.do(http.MethodPost, path, body)
	return err
}

type APIError struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

type PaginatedResponse struct {
	Size     int             `json:"size"`
	Page     int             `json:"page"`
	PageLen  int             `json:"pagelen"`
	Next     string          `json:"next"`
	Previous string          `json:"previous"`
	Values   json.RawMessage `json:"values"`
}
