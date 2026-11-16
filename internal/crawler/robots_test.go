package crawler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRobotsChecker_IsAllowed(t *testing.T) {
	robotsRequests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/robots.txt" {
			robotsRequests++
			_, _ = w.Write([]byte("User-agent: *\nDisallow: /private\nDisallow: /search?q=\n"))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	tests := []struct {
		path string
		want bool
	}{
		{"/", true},
		{"/public/page", true},
		{"/private", false},
		{"/private/page", false},
		{"/search?q=shoes", false},
		{"/search", true},
	}

	// The result must not depend on respectRobots: the manager decides whether
	// a disallowed URL is skipped or crawled and marked as blocked.
	for _, respect := range []bool{true, false} {
		checker := NewRobotsChecker(NewFetcher(5*time.Second, "barracuda-test"), "barracuda-test", respect)
		for _, tt := range tests {
			got, err := checker.IsAllowed(server.URL + tt.path)
			if err != nil {
				t.Fatalf("IsAllowed(%q) error: %v", tt.path, err)
			}
			if got != tt.want {
				t.Errorf("respect=%v IsAllowed(%q) = %v, want %v", respect, tt.path, got, tt.want)
			}
		}
	}

	if robotsRequests != 2 {
		t.Errorf("robots.txt fetched %d times, want 2 (once per checker, then cached)", robotsRequests)
	}
}

func TestRobotsChecker_MissingRobotsAllowsAllAndIsCached(t *testing.T) {
	robotsRequests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/robots.txt" {
			robotsRequests++
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	checker := NewRobotsChecker(NewFetcher(5*time.Second, "barracuda-test"), "barracuda-test", true)
	for _, path := range []string{"/a", "/b", "/c"} {
		allowed, err := checker.IsAllowed(server.URL + path)
		if err != nil || !allowed {
			t.Errorf("IsAllowed(%q) = %v, %v; want true, nil", path, allowed, err)
		}
	}

	if robotsRequests != 1 {
		t.Errorf("robots.txt fetched %d times, want 1 (a missing file should be cached)", robotsRequests)
	}
}
