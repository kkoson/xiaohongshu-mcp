// Package xiaohongshu provides a client for interacting with the Xiaohongshu (Little Red Book) platform.
package xiaohongshu

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// defaultBaseURL is the base URL for the Xiaohongshu API.
	defaultBaseURL = "https://www.xiaohongshu.com"
	// defaultTimeout is the default HTTP client timeout.
	// Increased from 30s to 60s to better handle slow responses from the API.
	// Personal note: 45s was still timing out occasionally on my network, bumping to 60s.
	defaultTimeout = 60 * time.Second
	// defaultUserAgent mimics a browser to avoid bot detection.
	// Updated to Chrome 124 to match my current browser version.
	defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
)

// Client is the HTTP client for the Xiaohongshu API.
type Client struct {
	httpClient *http.Client
	baseURL    string
	userAgent  string
	cookies    string
}

// Note represents a Xiaohongshu post/note.
type Note struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Author      string `json:"author"`
	AuthorID    string `json:"author_id"`
	Likes       int    `json:"likes"`
	Comments    int    `json:"comments"`
	Collects    int    `json:"collects"`
	Images      []string `json:"images"`
	Tags        []string `json:"tags"`
	CreatedAt   string `json:"created_at"`
	URL         string `json:"url"`
}

// SearchResult represents the result of a search query.
type SearchResult struct {
	Notes  []Note `json:"notes"`
	Total  int    `json:"total"`
	HasMore bool  `json:"has_more"`
}

// ClientOption is a functional option for configuring the Client.
type ClientOption func(*Client)

// WithCookies sets the cookies for authenticated requests.
func WithCookies(cookies string) ClientOption {
	return func(c *Client) {
		c.cookies = cookies
	}
}

// WithTimeout sets a custom timeout for the HTTP client.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithBaseURL overrides the default base URL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithUserAgent overrides the default user agent string.
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// NewClient creates a new Xiaohongshu client with the given options.
func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		baseURL:   defaultBaseURL,
		userAgent: defaultUserAgent,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// doRequest performs an HTTP GET request with appropriate headers.
func (c *Client) doRequest(endpoint string, params url.Values) ([]byte, error) {
	reqURL := fmt.Sprintf("%s%s", c.baseURL, endpoint)
	if len(params) > 0 {
		reqURL = fmt.Sprintf("%s?%s", reqURL, params.Encode())
	}

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept"
