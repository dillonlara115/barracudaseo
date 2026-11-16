package analyzer

import (
	"strings"
	"testing"

	"github.com/dillonlara115/barracudaseo/pkg/models"
)

func TestClassifyIndexability(t *testing.T) {
	selfCanonical := page("https://example.com/self")
	selfCanonical.Canonical = "https://example.com/self/"

	relativeCanonical := page("https://example.com/shop/item?color=red")
	relativeCanonical.Canonical = "/shop/item"

	noindex := page("https://example.com/noindex")
	noindex.IndexabilityStatus = models.IndexabilityNoindex

	blocked := page("https://example.com/blocked")
	blocked.IndexabilityStatus = models.IndexabilityBlocked

	redirect := page("https://example.com/old")
	redirect.RedirectChain = []string{"https://example.com/new"}

	// A redirect wins over a noindex on the destination: the URL itself redirects
	redirectToNoindex := page("https://example.com/old2")
	redirectToNoindex.RedirectChain = []string{"https://example.com/new2"}
	redirectToNoindex.IndexabilityStatus = models.IndexabilityNoindex

	tests := []struct {
		name string
		page *models.PageResult
		want IndexabilityClass
	}{
		{"plain page", page("https://example.com/a"), ClassIndexable},
		{"self canonical with trailing slash", selfCanonical, ClassIndexable},
		{"relative canonical to another URL", relativeCanonical, ClassCanonicalised},
		{"noindex", noindex, ClassNoindex},
		{"blocked", blocked, ClassBlocked},
		{"redirect", redirect, ClassRedirect},
		{"redirect to noindex", redirectToNoindex, ClassRedirect},
		{"404", &models.PageResult{URL: "https://example.com/gone", StatusCode: 404, Error: "HTTP 404"}, ClassError},
		{"fetch failure", &models.PageResult{URL: "https://example.com/timeout", Error: "request failed: timeout"}, ClassError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := classifyIndexability(tt.page); got != tt.want {
				t.Errorf("classifyIndexability() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAnalyzeIndexability_BreakdownAndDirectives(t *testing.T) {
	ok := page("https://example.com/")
	noindex := page("https://example.com/thanks")
	noindex.IndexabilityStatus = models.IndexabilityNoindex
	noindex.XRobotsTag = "noindex, nofollow"
	blocked := page("https://example.com/admin")
	blocked.IndexabilityStatus = models.IndexabilityBlocked
	gone := &models.PageResult{URL: "https://example.com/gone", StatusCode: 404, Error: "HTTP 404"}
	image := &models.PageResult{URL: "https://example.com/logo.png", StatusCode: 200}

	issues, breakdown := AnalyzeIndexability([]*models.PageResult{ok, noindex, blocked, gone, image})

	want := map[IndexabilityClass]int{ClassIndexable: 1, ClassNoindex: 1, ClassBlocked: 1, ClassError: 1}
	for class, count := range want {
		if breakdown[class] != count {
			t.Errorf("breakdown[%s] = %d, want %d", class, breakdown[class], count)
		}
	}
	total := 0
	for _, count := range breakdown {
		total += count
	}
	if total != 4 {
		t.Errorf("breakdown total = %d, want 4 (image URL excluded)", total)
	}

	noindexIssues := issuesOfType(issues, IssueNoindexPage)
	if len(noindexIssues) != 1 || !strings.Contains(noindexIssues[0].Message, "X-Robots-Tag header") {
		t.Errorf("noindex_page issues = %+v, want one naming the X-Robots-Tag header", noindexIssues)
	}
	if got := len(issuesOfType(issues, IssueBlockedByRobots)); got != 1 {
		t.Errorf("blocked_by_robots count = %d, want 1", got)
	}
	// No sitemap data in this crawl, so no sitemap issues
	if got := len(issuesOfType(issues, IssueNotInSitemap)) + len(issuesOfType(issues, IssueNonIndexableInSitemap)); got != 0 {
		t.Errorf("sitemap issue count = %d, want 0 without sitemap data", got)
	}
}

func TestAnalyzeIndexability_Canonicals(t *testing.T) {
	target := page("https://example.com/product")
	variant := page("https://example.com/product?ref=nav")
	variant.Canonical = "https://example.com/product"

	gone := &models.PageResult{URL: "https://example.com/gone", StatusCode: 404, Error: "HTTP 404"}
	toGone := page("https://example.com/points-to-gone")
	toGone.Canonical = "https://example.com/gone"

	toUncrawled := page("https://example.com/points-away")
	toUncrawled.Canonical = "https://example.com/never-crawled"

	issues, _ := AnalyzeIndexability([]*models.PageResult{target, variant, gone, toGone, toUncrawled})

	elsewhere := issuesOfType(issues, IssueCanonicalised)
	if len(elsewhere) != 2 {
		t.Fatalf("canonical_points_elsewhere count = %d, want 2: %+v", len(elsewhere), elsewhere)
	}
	for _, issue := range elsewhere {
		uncrawledNote := strings.Contains(issue.Message, "not crawled")
		if (issue.URL == toUncrawled.URL) != uncrawledNote {
			t.Errorf("unexpected 'not crawled' note for %s: %q", issue.URL, issue.Message)
		}
	}

	broken := issuesOfType(issues, IssueCanonicalNonIndexable)
	if len(broken) != 1 || broken[0].URL != toGone.URL {
		t.Fatalf("canonical_to_non_indexable = %+v, want one for %s", broken, toGone.URL)
	}
	if !strings.Contains(broken[0].Message, "HTTP 404") {
		t.Errorf("canonical_to_non_indexable message should give the reason, got %q", broken[0].Message)
	}
}

func TestAnalyzeIndexability_Sitemap(t *testing.T) {
	listed := page("https://example.com/listed")
	listed.InSitemap = true

	unlisted := page("https://example.com/unlisted")

	listedNoindex := page("https://example.com/listed-noindex")
	listedNoindex.InSitemap = true
	listedNoindex.IndexabilityStatus = models.IndexabilityNoindex

	listedRedirect := page("https://example.com/listed-old")
	listedRedirect.InSitemap = true
	listedRedirect.RedirectChain = []string{"https://example.com/listed-new"}

	unlistedNoindex := page("https://example.com/unlisted-noindex")
	unlistedNoindex.IndexabilityStatus = models.IndexabilityNoindex

	issues, _ := AnalyzeIndexability([]*models.PageResult{listed, unlisted, listedNoindex, listedRedirect, unlistedNoindex})

	nonIndexable := issuesOfType(issues, IssueNonIndexableInSitemap)
	if len(nonIndexable) != 2 {
		t.Fatalf("non_indexable_in_sitemap count = %d, want 2: %+v", len(nonIndexable), nonIndexable)
	}
	if nonIndexable[0].URL != listedNoindex.URL || !strings.Contains(nonIndexable[0].Message, "noindex") {
		t.Errorf("first non_indexable_in_sitemap = %+v", nonIndexable[0])
	}
	if nonIndexable[1].URL != listedRedirect.URL || !strings.Contains(nonIndexable[1].Message, "https://example.com/listed-new") {
		t.Errorf("second non_indexable_in_sitemap = %+v", nonIndexable[1])
	}

	missing := issuesOfType(issues, IssueNotInSitemap)
	if len(missing) != 1 || missing[0].URL != unlisted.URL {
		t.Errorf("indexable_not_in_sitemap = %+v, want one for %s", missing, unlisted.URL)
	}
}
