package analyzer

import (
	"strings"
	"testing"

	"github.com/dillonlara115/barracudaseo/pkg/models"
)

// contentPage builds an indexable page with content data.
func contentPage(url string, words int, hash, simHash string) *models.PageResult {
	p := page(url)
	p.WordCount = words
	p.ContentHash = hash
	p.ContentSimHash = simHash
	return p
}

func TestAnalyzeContent_ThinContent(t *testing.T) {
	thin := contentPage("https://example.com/thin", 40, "h1", "0000000000000000")
	full := contentPage("https://example.com/full", 900, "h2", "ffffffffffffffff")
	noData := page("https://example.com/old-crawl") // no content hash: not analyzed
	thinNoindex := contentPage("https://example.com/thanks", 10, "h3", "00000000ffffffff")
	thinNoindex.IndexabilityStatus = models.IndexabilityNoindex

	issues := AnalyzeContent([]*models.PageResult{thin, full, noData, thinNoindex})

	got := issuesOfType(issues, IssueThinContent)
	if len(got) != 1 || got[0].URL != thin.URL {
		t.Fatalf("thin_content = %+v, want one for %s", got, thin.URL)
	}
	if got[0].Value != "40" || !strings.Contains(got[0].Message, "40 words") {
		t.Errorf("thin_content issue = %+v, want the word count in value and message", got[0])
	}
}

func TestAnalyzeContent_ExactDuplicates(t *testing.T) {
	a := contentPage("https://example.com/a", 500, "same", "00000000000000ff")
	b := contentPage("https://example.com/b", 500, "same", "00000000000000ff")
	c := contentPage("https://example.com/c", 500, "other", "ffffffff00000000")

	issues := AnalyzeContent([]*models.PageResult{a, b, c})

	exact := issuesOfType(issues, IssueDuplicateContent)
	if len(exact) != 2 {
		t.Fatalf("duplicate_content count = %d, want 2: %+v", len(exact), exact)
	}
	if exact[0].URL != a.URL || exact[0].Value != b.URL {
		t.Errorf("first duplicate_content = %+v, want %s pointing at %s", exact[0], a.URL, b.URL)
	}
	// Identical pages have identical similarity hashes, but must be reported
	// once (as exact duplicates), not again as near-duplicates.
	if near := issuesOfType(issues, IssueNearDuplicateContent); len(near) != 0 {
		t.Errorf("near_duplicate_content = %+v, want none for exact duplicates", near)
	}
}

func TestAnalyzeContent_NearDuplicates(t *testing.T) {
	// a and b differ in 2 bits, c differs from both in many
	a := contentPage("https://example.com/austin", 400, "ha", "00000000000000ff")
	b := contentPage("https://example.com/dallas", 400, "hb", "00000000000000fc")
	c := contentPage("https://example.com/about", 400, "hc", "ffffffff00000000")
	// Too short to compare reliably, even with an identical similarity hash
	short := contentPage("https://example.com/short", 20, "hd", "00000000000000ff")

	issues := AnalyzeContent([]*models.PageResult{a, b, c, short})

	near := issuesOfType(issues, IssueNearDuplicateContent)
	if len(near) != 2 {
		t.Fatalf("near_duplicate_content count = %d, want 2: %+v", len(near), near)
	}
	if near[0].URL != a.URL || near[0].Value != b.URL {
		t.Errorf("first near_duplicate_content = %+v, want %s pointing at %s", near[0], a.URL, b.URL)
	}
	if !strings.Contains(near[0].Message, "96% similar") {
		t.Errorf("message = %q, want the similarity percentage (62 of 64 bits = 96%%)", near[0].Message)
	}
}
