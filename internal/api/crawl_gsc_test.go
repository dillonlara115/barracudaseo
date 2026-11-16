package api

import (
	"testing"

	"github.com/dillonlara115/barracudaseo/pkg/models"
)

func crawled(url string) *models.PageResult {
	return &models.PageResult{URL: url, StatusCode: 200, IndexabilityStatus: models.IndexabilityIndexable}
}

func TestCoverageKey(t *testing.T) {
	same := [][2]string{
		{"https://example.com/page", "http://www.example.com/page/"},
		{"https://example.com/page", "https://example.com/page#section"},
		{"https://Example.com", "https://example.com/"},
	}
	for _, pair := range same {
		if coverageKey(pair[0]) != coverageKey(pair[1]) {
			t.Errorf("coverageKey(%q) != coverageKey(%q), want equal", pair[0], pair[1])
		}
	}

	different := [][2]string{
		{"https://example.com/page", "https://example.com/page?ref=1"},
		{"https://example.com/page", "https://example.com/Page"},
		{"https://example.com/page", "https://blog.example.com/page"},
	}
	for _, pair := range different {
		if coverageKey(pair[0]) == coverageKey(pair[1]) {
			t.Errorf("coverageKey(%q) == coverageKey(%q), want different", pair[0], pair[1])
		}
	}
}

func TestCompareCrawlWithGSC(t *testing.T) {
	visible := crawled("https://example.com/visible")
	invisible := crawled("https://example.com/invisible")

	noindex := crawled("https://example.com/noindex")
	noindex.IndexabilityStatus = models.IndexabilityNoindex

	gone := &models.PageResult{URL: "https://example.com/gone", StatusCode: 404}

	redirect := crawled("https://example.com/old")
	redirect.RedirectChain = []string{"https://example.com/new"}

	// Not indexable and not shown by Google: consistent, so not reported
	quietNoindex := crawled("https://example.com/thanks")
	quietNoindex.IndexabilityStatus = models.IndexabilityNoindex

	pages := []*models.PageResult{visible, invisible, noindex, gone, redirect, quietNoindex}

	gsc := []gscPage{
		{URL: "https://www.example.com/visible/", Clicks: 10, Impressions: 100},
		{URL: "https://example.com/visible#faq", Clicks: 2, Impressions: 30},
		{URL: "https://example.com/noindex", Clicks: 5, Impressions: 80},
		{URL: "https://example.com/gone", Clicks: 20, Impressions: 300},
		{URL: "https://example.com/old", Clicks: 0, Impressions: 12},
		{URL: "https://example.com/unlinked", Clicks: 3, Impressions: 40},
		{URL: "https://example.com/zero", Clicks: 0, Impressions: 0}, // not visible
	}

	coverage := compareCrawlWithGSC(pages, gsc)

	if coverage.CrawledPages != 6 || coverage.IndexablePages != 2 {
		t.Errorf("crawled = %d, indexable = %d; want 6 and 2", coverage.CrawledPages, coverage.IndexablePages)
	}
	if coverage.GooglePages != 5 {
		t.Errorf("GooglePages = %d, want 5 (variants merged, zero-impression row dropped)", coverage.GooglePages)
	}
	if coverage.VisibleIndexablePages != 1 {
		t.Errorf("VisibleIndexablePages = %d, want 1", coverage.VisibleIndexablePages)
	}

	if coverage.SearchTrafficAtRiskCount != 3 {
		t.Fatalf("SearchTrafficAtRiskCount = %d, want 3: %+v", coverage.SearchTrafficAtRiskCount, coverage.SearchTrafficAtRisk)
	}
	top := coverage.SearchTrafficAtRisk[0]
	if top.URL != gone.URL || top.Class != "error" || top.Clicks != 20 {
		t.Errorf("top at-risk page = %+v, want the 404 with the most clicks first", top)
	}
	if last := coverage.SearchTrafficAtRisk[2]; last.URL != redirect.URL || last.Class != "redirect" {
		t.Errorf("last at-risk page = %+v, want the redirect", last)
	}

	if coverage.NotInCrawlCount != 1 || coverage.NotInCrawl[0].URL != "https://example.com/unlinked" {
		t.Errorf("NotInCrawl = %+v, want only /unlinked", coverage.NotInCrawl)
	}
	if coverage.NoSearchVisibilityCount != 1 || coverage.NoSearchVisibility[0].URL != invisible.URL {
		t.Errorf("NoSearchVisibility = %+v, want only /invisible", coverage.NoSearchVisibility)
	}
}

func TestCompareCrawlWithGSC_MergesVariantMetrics(t *testing.T) {
	noindex := crawled("https://example.com/page")
	noindex.IndexabilityStatus = models.IndexabilityNoindex

	coverage := compareCrawlWithGSC([]*models.PageResult{noindex}, []gscPage{
		{URL: "https://example.com/page", Clicks: 4, Impressions: 50},
		{URL: "http://example.com/page/", Clicks: 1, Impressions: 10},
	})

	if coverage.SearchTrafficAtRiskCount != 1 {
		t.Fatalf("SearchTrafficAtRiskCount = %d, want 1", coverage.SearchTrafficAtRiskCount)
	}
	if got := coverage.SearchTrafficAtRisk[0]; got.Clicks != 5 || got.Impressions != 60 {
		t.Errorf("merged metrics = %v clicks, %v impressions; want 5 and 60", got.Clicks, got.Impressions)
	}
}

func TestCompareCrawlWithGSC_EmptyListsAreNotNil(t *testing.T) {
	coverage := compareCrawlWithGSC(nil, nil)
	if coverage.SearchTrafficAtRisk == nil || coverage.NotInCrawl == nil || coverage.NoSearchVisibility == nil {
		t.Error("lists should be empty slices, not nil, so they serialize as []")
	}
}
