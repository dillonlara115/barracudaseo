package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/zap"
)

type fakeCrawlJob struct {
	ran []string
	err error
}

func (f *fakeCrawlJob) Run(_ context.Context, crawlID string) error {
	f.ran = append(f.ran, crawlID)
	return f.err
}

func TestStartCrawlUsesJobWhenConfigured(t *testing.T) {
	job := &fakeCrawlJob{}
	s := &Server{logger: zap.NewNop(), crawlJob: job}

	s.startCrawl("crawl-1", "project-1", TriggerCrawlRequest{URL: "https://example.com"})

	if len(job.ran) != 1 || job.ran[0] != "crawl-1" {
		t.Fatalf("job ran %v, want [crawl-1]", job.ran)
	}
	inProcess := 0
	s.activeCrawls.Range(func(_, _ interface{}) bool {
		inProcess++
		return true
	})
	if inProcess != 0 {
		t.Fatalf("crawl also started in process (%d active), want job only", inProcess)
	}
}

func TestRunInProcessWhenJobFailsDefault(t *testing.T) {
	// Documents the current choice; update this test if the rule changes.
	if !runInProcessWhenJobFails(errors.New("permission denied")) {
		t.Fatal("expected a failed job launch to fall back to the API process")
	}
}

func TestMetaWithErrorKeepsExistingMeta(t *testing.T) {
	meta := map[string]interface{}{"url": "https://example.com", "schedule_id": "s1"}

	got := metaWithError(meta, "boom")

	if got["url"] != "https://example.com" || got["schedule_id"] != "s1" {
		t.Fatalf("existing meta lost: %v", got)
	}
	if got["error"] != "boom" {
		t.Fatalf("error = %v, want boom", got["error"])
	}
	if _, changed := meta["error"]; changed {
		t.Fatal("metaWithError modified its input")
	}
	if got := metaWithError(nil, "boom"); got["error"] != "boom" || len(got) != 1 {
		t.Fatalf("nil meta: got %v", got)
	}
}

func TestStaleCrawlFilter(t *testing.T) {
	cutoff := time.Date(2026, 11, 17, 10, 0, 0, 0, time.UTC)

	got := staleCrawlFilter(cutoff)

	want := `heartbeat_at.lt."2026-11-17T10:00:00Z",and(heartbeat_at.is.null,started_at.lt."2026-11-17T10:00:00Z")`
	if got != want {
		t.Fatalf("filter =\n%s\nwant\n%s", got, want)
	}
}

func TestCrawlRequestFromMeta(t *testing.T) {
	// Numbers arrive as float64 after a JSON round trip
	meta := map[string]interface{}{
		"url":                "https://example.com",
		"max_depth":          float64(5),
		"max_pages":          float64(250),
		"workers":            float64(4),
		"respect_robots":     false,
		"parse_sitemap":      true,
		"crawl_sitemap_only": true,
		"schedule_id":        "ignored",
	}

	req := crawlRequestFromMeta(meta)

	if req.URL != "https://example.com" || req.MaxDepth != 5 || req.MaxPages != 250 || req.Workers != 4 {
		t.Fatalf("unexpected request: %+v", req)
	}
	if *req.RespectRobots || !*req.ParseSitemap || !*req.CrawlSitemapOnly {
		t.Fatalf("unexpected flags: robots=%v sitemap=%v sitemapOnly=%v", *req.RespectRobots, *req.ParseSitemap, *req.CrawlSitemapOnly)
	}
}

func TestCrawlRequestFromMetaDefaultsAndCeiling(t *testing.T) {
	req := crawlRequestFromMeta(map[string]interface{}{
		"url":       "https://example.com",
		"max_pages": float64(9_000_000),
	})

	if req.MaxPages != getMaxPagesLimit("team") {
		t.Fatalf("MaxPages = %d, want the ceiling %d", req.MaxPages, getMaxPagesLimit("team"))
	}
	if req.Workers != 10 || req.MaxDepth != 3 || !*req.RespectRobots {
		t.Fatalf("unexpected defaults: %+v", req)
	}

	if empty := crawlRequestFromMeta(nil); empty.URL != "" {
		t.Fatalf("nil meta produced URL %q", empty.URL)
	}
}

func TestNewCloudRunJobRunnerRejectsBadName(t *testing.T) {
	for _, name := range []string{"barracuda-crawl", "projects/p/jobs/j", "projects/p/locations/r/jobs/"} {
		if _, err := newCloudRunJobRunner(context.Background(), name); err == nil {
			t.Errorf("name %q accepted, want an error", name)
		}
	}
}

func TestDiscoveredPagesNeverBelowCrawled(t *testing.T) {
	if got := discoveredPages(21, 35); got != 35 {
		t.Errorf("discoveredPages(21, 35) = %d, want 35", got)
	}
	if got := discoveredPages(21, 4); got != 21 {
		t.Errorf("discoveredPages(21, 4) = %d, want 21", got)
	}
}
