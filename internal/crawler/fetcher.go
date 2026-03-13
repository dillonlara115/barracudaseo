package crawler

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/dillonlara115/barracudaseo/internal/utils"
	"github.com/dillonlara115/barracudaseo/pkg/models"
)

const maxRedirects = 10

// Fetcher handles HTTP requests and response processing
type Fetcher struct {
	client    *http.Client
	userAgent string
}

// FetchResult contains the fetched page data
type FetchResult struct {
	PageResult *models.PageResult
	Body       []byte
	Error      error
}

// NewFetcher creates a new Fetcher instance
func NewFetcher(timeout time.Duration, userAgent string) *Fetcher {
	client := &http.Client{
		Timeout: timeout,
		// Disable automatic redirect following — we handle redirects manually
		// per-request to avoid race conditions with concurrent workers.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return &Fetcher{
		client:    client,
		userAgent: userAgent,
	}
}

// isRedirect returns true for HTTP status codes that indicate a redirect.
func isRedirect(statusCode int) bool {
	return statusCode == 301 || statusCode == 302 || statusCode == 303 || statusCode == 307 || statusCode == 308
}

// Fetch retrieves a URL and returns the response (single attempt, no retry).
// Redirects are followed manually in a loop so that each concurrent request
// tracks its own redirect chain without shared mutable state.
func (f *Fetcher) Fetch(url string) *FetchResult {
	result := &FetchResult{
		PageResult: &models.PageResult{
			URL:       url,
			CrawledAt: time.Now(),
		},
	}

	startTime := time.Now()

	var redirectChain []string
	currentURL := url

	var resp *http.Response
	for hops := 0; ; hops++ {
		req, err := http.NewRequest("GET", currentURL, nil)
		if err != nil {
			result.Error = fmt.Errorf("failed to create request: %w", err)
			result.PageResult.Error = result.Error.Error()
			result.PageResult.ResponseTime = time.Since(startTime).Milliseconds()
			return result
		}

		req.Header.Set("User-Agent", f.userAgent)
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

		r, err := f.client.Do(req)
		if err != nil {
			result.Error = fmt.Errorf("request failed: %w", err)
			result.PageResult.Error = result.Error.Error()
			result.PageResult.ResponseTime = time.Since(startTime).Milliseconds()
			return result
		}

		if !isRedirect(r.StatusCode) {
			resp = r
			break
		}

		// This is a redirect — record the destination and follow it.
		location := r.Header.Get("Location")
		r.Body.Close()

		if location == "" {
			result.Error = fmt.Errorf("redirect %d with no Location header", r.StatusCode)
			result.PageResult.Error = result.Error.Error()
			result.PageResult.StatusCode = r.StatusCode
			result.PageResult.ResponseTime = time.Since(startTime).Milliseconds()
			return result
		}

		// Resolve relative redirect URLs against the current request URL.
		nextURL, err := req.URL.Parse(location)
		if err != nil {
			result.Error = fmt.Errorf("invalid redirect location %q: %w", location, err)
			result.PageResult.Error = result.Error.Error()
			result.PageResult.ResponseTime = time.Since(startTime).Milliseconds()
			return result
		}

		redirectChain = append(redirectChain, nextURL.String())
		currentURL = nextURL.String()

		if hops >= maxRedirects {
			result.Error = fmt.Errorf("stopped after %d redirects", maxRedirects)
			result.PageResult.Error = result.Error.Error()
			result.PageResult.ResponseTime = time.Since(startTime).Milliseconds()
			if len(redirectChain) > 0 {
				result.PageResult.RedirectChain = redirectChain
			}
			return result
		}
	}

	responseTime := time.Since(startTime)
	defer resp.Body.Close()

	result.PageResult.StatusCode = resp.StatusCode
	result.PageResult.ResponseTime = responseTime.Milliseconds()

	// Extract x-robots-tag header for indexability detection
	xRobotsTag := resp.Header.Get("X-Robots-Tag")
	if xRobotsTag != "" {
		result.PageResult.XRobotsTag = xRobotsTag
	}

	if len(redirectChain) > 0 {
		result.PageResult.RedirectChain = redirectChain
	}

	// Check Content-Type header - skip non-HTML content (images, PDFs, etc.)
	contentType := resp.Header.Get("Content-Type")
	if contentType != "" {
		contentTypeLower := strings.ToLower(contentType)
		if strings.HasPrefix(contentTypeLower, "image/") {
			result.Error = fmt.Errorf("skipped non-HTML content: %s", contentType)
			result.PageResult.Error = result.Error.Error()
			return result
		}
		nonHTMLTypes := []string{"application/pdf", "application/zip", "application/json", "application/xml"}
		for _, nonHTML := range nonHTMLTypes {
			if strings.HasPrefix(contentTypeLower, nonHTML) {
				result.Error = fmt.Errorf("skipped non-HTML content: %s", contentType)
				result.PageResult.Error = result.Error.Error()
				return result
			}
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Error = fmt.Errorf("failed to read response body: %w", err)
		result.PageResult.Error = result.Error.Error()
		return result
	}

	result.Body = body

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		result.Error = fmt.Errorf("HTTP %d", resp.StatusCode)
		result.PageResult.Error = result.Error.Error()
	}

	return result
}

// FetchRaw performs a simple HTTP GET that follows redirects and returns the
// response body without content-type filtering. Use this for fetching resources
// like sitemaps and robots.txt that are not HTML pages.
func (f *Fetcher) FetchRaw(rawURL string) ([]byte, int, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", f.userAgent)

	currentURL := rawURL
	for hops := 0; ; hops++ {
		r, err := f.client.Do(req)
		if err != nil {
			return nil, 0, fmt.Errorf("request failed: %w", err)
		}

		if !isRedirect(r.StatusCode) {
			defer r.Body.Close()
			body, err := io.ReadAll(r.Body)
			if err != nil {
				return nil, r.StatusCode, fmt.Errorf("failed to read body: %w", err)
			}
			return body, r.StatusCode, nil
		}

		location := r.Header.Get("Location")
		r.Body.Close()
		if location == "" {
			return nil, r.StatusCode, fmt.Errorf("redirect %d with no Location header", r.StatusCode)
		}

		nextURL, err := req.URL.Parse(location)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid redirect location %q: %w", location, err)
		}
		currentURL = nextURL.String()

		if hops >= maxRedirects {
			return nil, r.StatusCode, fmt.Errorf("stopped after %d redirects", maxRedirects)
		}

		req, err = http.NewRequest("GET", currentURL, nil)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("User-Agent", f.userAgent)
	}
}

// isRetryableError checks if an error is retryable
func isRetryableError(result *FetchResult) bool {
	if result.Error == nil {
		return false
	}

	// Retry on 5xx errors, timeouts, and connection errors
	statusCode := result.PageResult.StatusCode
	if statusCode >= 500 && statusCode < 600 {
		return true
	}

	// Check for timeout or connection errors in error message
	errMsg := result.Error.Error()
	if containsAny(errMsg, []string{"timeout", "connection refused", "no such host", "network is unreachable"}) {
		return true
	}

	return false
}

// containsAny checks if a string contains any of the substrings
func containsAny(s string, substrings []string) bool {
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

// FetchWithRetry retrieves a URL with retry logic for transient errors
func (f *Fetcher) FetchWithRetry(url string, maxRetries int) *FetchResult {
	var lastResult *FetchResult

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: wait 2^attempt seconds
			backoff := time.Duration(1<<uint(attempt-1)) * time.Second
			time.Sleep(backoff)
		}

		result := f.Fetch(url)
		lastResult = result

		// If successful or not retryable, return immediately
		if result.Error == nil || !isRetryableError(result) {
			return result
		}

		// Log retry attempt
		utils.Debug("Retrying fetch",
			utils.NewField("url", url),
			utils.NewField("attempt", attempt+1),
			utils.NewField("max_retries", maxRetries),
		)
	}

	return lastResult
}
