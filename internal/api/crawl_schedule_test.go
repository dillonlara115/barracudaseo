package api

import (
	"testing"
	"time"
)

func TestNextRunAfter(t *testing.T) {
	from := time.Date(2026, 1, 31, 9, 0, 0, 0, time.UTC)

	tests := []struct {
		frequency string
		want      time.Time
	}{
		{"daily", time.Date(2026, 2, 1, 9, 0, 0, 0, time.UTC)},
		{"weekly", time.Date(2026, 2, 7, 9, 0, 0, 0, time.UTC)},
		// Go normalizes Feb 31 to Mar 3; the schedule still moves forward a month-ish
		{"monthly", time.Date(2026, 3, 3, 9, 0, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		got, err := nextRunAfter(tt.frequency, from)
		if err != nil {
			t.Errorf("nextRunAfter(%q) error: %v", tt.frequency, err)
		}
		if !got.Equal(tt.want) {
			t.Errorf("nextRunAfter(%q) = %v, want %v", tt.frequency, got, tt.want)
		}
	}

	if _, err := nextRunAfter("hourly", from); err == nil {
		t.Error("nextRunAfter(hourly) should fail for an unknown frequency")
	}
}

func TestFrequencyAllowed(t *testing.T) {
	if frequencyAllowed("free", "weekly") {
		t.Error("free tier should not be able to schedule crawls")
	}
	if !frequencyAllowed("pro", "weekly") {
		t.Error("pro tier should be able to schedule weekly crawls")
	}
	if frequencyAllowed("pro", "hourly") {
		t.Error("unknown frequencies are never allowed")
	}
}

func TestScheduleCrawlRequest(t *testing.T) {
	// Stored config comes back from JSON, so numbers are float64
	req := scheduleCrawlRequest(map[string]interface{}{
		"url":           "https://example.com",
		"max_pages":     float64(50000),
		"max_depth":     float64(5),
		"parse_sitemap": true,
	}, 10000)

	if req.URL != "https://example.com" {
		t.Errorf("URL = %q", req.URL)
	}
	if req.MaxPages != 10000 {
		t.Errorf("MaxPages = %d, want it capped at the plan limit 10000", req.MaxPages)
	}
	if req.MaxDepth != 5 {
		t.Errorf("MaxDepth = %d, want 5", req.MaxDepth)
	}
	if !*req.ParseSitemap || !*req.RespectRobots || *req.CrawlSitemapOnly {
		t.Errorf("flags = parse_sitemap %v, respect_robots %v, crawl_sitemap_only %v; want true, true, false",
			*req.ParseSitemap, *req.RespectRobots, *req.CrawlSitemapOnly)
	}

	defaults := scheduleCrawlRequest(map[string]interface{}{}, 100)
	if defaults.MaxPages != 100 || defaults.MaxDepth != 3 || defaults.Workers != 10 {
		t.Errorf("defaults = %+v, want max_pages 100, max_depth 3, workers 10", defaults)
	}
}

func TestBareHost(t *testing.T) {
	tests := map[string]string{
		"example.com":                  "example.com",
		"https://www.example.com/path": "www.example.com",
		"http://example.com:8080/":     "example.com",
		"example.com/blog?x=1":         "example.com",
	}
	for input, want := range tests {
		if got := bareHost(input); got != want {
			t.Errorf("bareHost(%q) = %q, want %q", input, got, want)
		}
	}
}
