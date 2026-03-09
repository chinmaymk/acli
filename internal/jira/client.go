package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/chinmaymk/acli/internal/config"
)

// Client is the Jira REST API client.
type Client struct {
	BaseURL    string
	Email      string
	APIToken   string
	HTTPClient *http.Client
}

// NewClient creates a new Jira client from a config profile.
func NewClient(profile config.Profile) *Client {
	baseURL := strings.TrimRight(profile.AtlassianURL, "/")
	return &Client{
		BaseURL:    baseURL,
		Email:      profile.Email,
		APIToken:   profile.APIToken,
		HTTPClient: &http.Client{},
	}
}

// APIError represents a Jira API error response.
type APIError struct {
	StatusCode   int
	ErrorMessages []string          `json:"errorMessages"`
	Errors       map[string]string `json:"errors"`
}

func (e *APIError) Error() string {
	parts := append([]string{}, e.ErrorMessages...)
	for k, v := range e.Errors {
		parts = append(parts, fmt.Sprintf("%s: %s", k, v))
	}
	if len(parts) == 0 {
		return fmt.Sprintf("Jira API error (HTTP %d)", e.StatusCode)
	}
	return fmt.Sprintf("Jira API error (HTTP %d): %s", e.StatusCode, strings.Join(parts, "; "))
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	u := c.BaseURL + path

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, u, bodyReader)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.Email, c.APIToken)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		apiErr := &APIError{StatusCode: resp.StatusCode}
		if json.Unmarshal(body, apiErr) != nil {
			apiErr.ErrorMessages = []string{string(body)}
		}
		return apiErr
	}

	if v != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("decoding response: %w", err)
		}
	}

	return nil
}

// Get performs a GET request.
func (c *Client) Get(path string, query url.Values, v interface{}) error {
	if query != nil {
		path = path + "?" + query.Encode()
	}
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return err
	}
	return c.do(req, v)
}

// Post performs a POST request.
func (c *Client) Post(path string, body interface{}, v interface{}) error {
	req, err := c.newRequest("POST", path, body)
	if err != nil {
		return err
	}
	return c.do(req, v)
}

// Put performs a PUT request.
func (c *Client) Put(path string, body interface{}, v interface{}) error {
	req, err := c.newRequest("PUT", path, body)
	if err != nil {
		return err
	}
	return c.do(req, v)
}

// Delete performs a DELETE request.
func (c *Client) Delete(path string, query url.Values) error {
	if query != nil {
		path = path + "?" + query.Encode()
	}
	req, err := c.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// Patch performs a PATCH request.
func (c *Client) Patch(path string, body interface{}, v interface{}) error {
	req, err := c.newRequest("PATCH", path, body)
	if err != nil {
		return err
	}
	return c.do(req, v)
}

// UploadFile uploads a file using multipart form data.
func (c *Client) UploadFile(path string, fieldName string, filePath string, v interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile(fieldName, filePath)
	if err != nil {
		return fmt.Errorf("creating form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("copying file data: %w", err)
	}
	writer.Close()

	u := c.BaseURL + path
	req, err := http.NewRequest("POST", u, &buf)
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.Email, c.APIToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Atlassian-Token", "no-check")

	return c.do(req, v)
}

// GetRaw performs a GET request and returns the raw response body.
func (c *Client) GetRaw(path string, query url.Values) ([]byte, error) {
	if query != nil {
		path = path + "?" + query.Encode()
	}
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		apiErr := &APIError{StatusCode: resp.StatusCode}
		if json.Unmarshal(body, apiErr) != nil {
			apiErr.ErrorMessages = []string{string(body)}
		}
		return nil, apiErr
	}
	return io.ReadAll(resp.Body)
}
