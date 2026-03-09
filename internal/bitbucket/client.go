package bitbucket

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const baseURL = "https://api.bitbucket.org/2.0"

type Client struct {
	httpClient *http.Client
	token      string
	username   string // if set, use Basic auth (username:token)
}

func NewClient() (*Client, error) {
	token := os.Getenv("BB_SCOPED_TOKEN")
	if token == "" {
		token = os.Getenv("BB_TOKEN")
	}
	if token == "" {
		return nil, fmt.Errorf("BB_SCOPED_TOKEN (or BB_TOKEN) environment variable is not set")
	}
	username := os.Getenv("BB_USERNAME")
	if username == "" {
		username = os.Getenv("BB_EMAIL")
	}
	return &Client{
		httpClient: &http.Client{},
		token:      token,
		username:   username,
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

	if c.username != "" {
		req.SetBasicAuth(c.username, c.token)
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

func (c *Client) delete(path string) ([]byte, error) {
	return c.do(http.MethodDelete, path, nil)
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
