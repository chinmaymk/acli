package bitbucket

import (
	"net/http"
	"net/http/httptest"
)

// redirectTransport is an http.RoundTripper that rewrites all requests to point
// to the test server, preserving the path and query.
type redirectTransport struct {
	server *httptest.Server
}

func (t *redirectTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Rewrite the URL to point to our test server
	req.URL.Scheme = "http"
	req.URL.Host = t.server.Listener.Addr().String()
	return http.DefaultTransport.RoundTrip(req)
}

// newRedirectClient creates a Client whose HTTP requests are all redirected
// to the given test server. This allows testing higher-level API methods
// that construct paths using the hardcoded baseURL constant.
func newRedirectClient(srv *httptest.Server) *Client {
	return &Client{
		httpClient: &http.Client{
			Transport: &redirectTransport{server: srv},
		},
		token: "testtoken",
	}
}
