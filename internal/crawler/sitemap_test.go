package crawler

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestDiscoverSitemapURLPrefersRobotsDirective(t *testing.T) {
	parser := NewSitemapParser(newStubFetcher(map[string]stubResponse{
		"https://example.com/robots.txt": {
			status: http.StatusOK,
			body:   "User-agent: *\nSitemap: /custom-sitemap.xml\n",
		},
		"https://example.com/custom-sitemap.xml": {
			status: http.StatusOK,
			body:   `<?xml version="1.0"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"></urlset>`,
		},
	}))
	got := parser.DiscoverSitemapURL("https://example.com")

	if got != "https://example.com/custom-sitemap.xml" {
		t.Fatalf("DiscoverSitemapURL() = %q, want %q", got, "https://example.com/custom-sitemap.xml")
	}
}

func TestDiscoverSitemapURLFallsBackToWordPressPath(t *testing.T) {
	parser := NewSitemapParser(newStubFetcher(map[string]stubResponse{
		"https://example.com/robots.txt":        {status: http.StatusNotFound},
		"https://example.com/sitemap.xml":       {status: http.StatusNotFound},
		"https://example.com/sitemap_index.xml": {status: http.StatusNotFound},
		"https://example.com/wp-sitemap.xml": {
			status: http.StatusOK,
			body:   `<?xml version="1.0"?><sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"></sitemapindex>`,
		},
	}))
	got := parser.DiscoverSitemapURL("https://example.com")

	if got != "https://example.com/wp-sitemap.xml" {
		t.Fatalf("DiscoverSitemapURL() = %q, want %q", got, "https://example.com/wp-sitemap.xml")
	}
}

type stubResponse struct {
	status int
	body   string
}

type stubRoundTripper struct {
	responses map[string]stubResponse
}

func (s stubRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, ok := s.responses[req.URL.String()]
	if !ok {
		resp = stubResponse{status: http.StatusNotFound}
	}

	return &http.Response{
		StatusCode: resp.status,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(resp.body)),
		Request:    req,
	}, nil
}

func newStubFetcher(responses map[string]stubResponse) *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Transport: stubRoundTripper{responses: responses},
		},
		userAgent: "test-agent",
	}
}
