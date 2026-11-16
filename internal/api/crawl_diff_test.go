package api

import (
	"testing"

	"github.com/dillonlara115/barracudaseo/internal/analyzer"
)

// ref builds an issueRef with its fingerprint, as fetchCrawlIssueRefs would.
func ref(issueType, url, value string) issueRef {
	return issueRef{
		Type:        issueType,
		Severity:    "warning",
		URL:         url,
		Value:       value,
		Fingerprint: analyzer.IssueFingerprint(analyzer.IssueType(issueType), url, value),
	}
}

func TestDiffCrawls(t *testing.T) {
	const (
		home    = "https://example.com"
		about   = "https://example.com/about"
		removed = "https://example.com/removed"
		added   = "https://example.com/added"
	)

	prevIssues := []issueRef{
		ref("missing_h1", home, ""),                           // still present
		ref("long_title", about, "Old title"),                 // fixed
		ref("missing_title", removed, ""),                     // page no longer crawled
		ref("missing_image_alt", home, home+"/a.jpg"),         // fixed
		ref("missing_image_alt", home, home+"/b.jpg"),         // still present
		ref("large_image", about, about+"/hero.jpg (300 KB)"), // still present, size changed
		ref("large_image", about, about+"/hero.jpg (300 KB)"), // duplicate row from an old crawl
	}
	currIssues := []issueRef{
		ref("missing_h1", home+"/", ""), // same page, trailing slash
		ref("missing_image_alt", home, home+"/b.jpg"),
		ref("large_image", about, about+"/hero.jpg (280 KB)"),
		ref("missing_meta_description", added, ""), // new, on a new page
		ref("broken_internal_link", home, removed), // new, on an existing page
	}
	currIssues[4].Severity = "error"

	diff := diffCrawls(prevIssues, currIssues, []string{home, about, removed}, []string{home + "/", about, added})

	if diff.UnchangedCount != 3 {
		t.Errorf("UnchangedCount = %d, want 3", diff.UnchangedCount)
	}
	if diff.NewCount != 2 {
		t.Errorf("NewCount = %d, want 2", diff.NewCount)
	}
	if diff.FixedCount != 2 {
		t.Errorf("FixedCount = %d, want 2", diff.FixedCount)
	}
	if diff.NotRecrawledCount != 1 {
		t.Errorf("NotRecrawledCount = %d, want 1 (issue on a page that was not crawled again)", diff.NotRecrawledCount)
	}

	if len(diff.NewIssues) != 2 || diff.NewIssues[0].Type != "broken_internal_link" {
		t.Errorf("NewIssues = %+v, want the error-severity issue first", diff.NewIssues)
	}
	if got := diff.ByType["missing_image_alt"]; got != (issueTypeDelta{Fixed: 1, Unchanged: 1}) {
		t.Errorf("ByType[missing_image_alt] = %+v, want 1 fixed and 1 unchanged", got)
	}

	if diff.PagesAddedCount != 1 || diff.PagesAdded[0] != added {
		t.Errorf("PagesAdded = %v, want [%s]", diff.PagesAdded, added)
	}
	if diff.PagesRemovedCount != 1 || diff.PagesRemoved[0] != removed {
		t.Errorf("PagesRemoved = %v, want [%s]", diff.PagesRemoved, removed)
	}
	if diff.Truncated {
		t.Error("Truncated = true, want false")
	}
}

func TestDiffCrawls_IdenticalCrawls(t *testing.T) {
	issues := []issueRef{ref("missing_h1", "https://example.com", "")}
	pages := []string{"https://example.com"}

	diff := diffCrawls(issues, issues, pages, pages)

	if diff.NewCount != 0 || diff.FixedCount != 0 || diff.UnchangedCount != 1 {
		t.Errorf("diff = %+v, want only 1 unchanged", diff)
	}
	if diff.NewIssues == nil || diff.FixedIssues == nil || diff.PagesAdded == nil {
		t.Error("lists should be empty slices, not nil, so they serialize as []")
	}
}

func TestCarriedStatus(t *testing.T) {
	carried := map[string]string{"abc": "ignored"}
	if got := carriedStatus(carried, "abc"); got != "ignored" {
		t.Errorf("carriedStatus(known) = %q, want ignored", got)
	}
	if got := carriedStatus(carried, "other"); got != "new" {
		t.Errorf("carriedStatus(unknown) = %q, want new", got)
	}
}
