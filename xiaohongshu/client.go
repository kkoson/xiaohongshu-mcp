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
	defaultTimeout = 30 * time.Second
	// defaultUserAgent mimics a browser to avoid bot detection.
	defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
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
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Referer", c.baseURL)
	if c.cookies != "" {
		req.Header.Set("Cookie", c.cookies)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}
	return body, nil
}

// parseJSON is a helper to unmarshal JSON bytes into a target struct.
func parseJSON(data []byte, target interface{}) error {
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("parsing JSON response: %w", err)
	}
	return nil
}
