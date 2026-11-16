package analyzer

import (
	"strings"
	"testing"

	"github.com/dillonlara115/barracudaseo/pkg/models"
)

// page builds an indexable 200 page with unique-by-default metadata.
func page(url string, links ...string) *models.PageResult {
	return &models.PageResult{
		URL:                url,
		StatusCode:         200,
		Title:              "Title for " + url,
		MetaDesc:           "Description for " + url,
		H1:                 []string{"Heading for " + url},
		InternalLinks:      links,
		IndexabilityStatus: models.IndexabilityIndexable,
	}
}

// issuesOfType returns the URLs of issues with the given type.
func issuesOfType(issues []Issue, t IssueType) []Issue {
	var out []Issue
	for _, issue := range issues {
		if issue.Type == t {
			out = append(out, issue)
		}
	}
	return out
}

func TestAnalyzeCrossPage_Duplicates(t *testing.T) {
	a := page("https://example.com/a", "https://example.com/b", "https://example.com/c")
	b := page("https://example.com/b", "https://example.com/a")
	c := page("https://example.com/c", "https://example.com/a")
	home := page("https://example.com", "https://example.com/a")

	// a and b share a title, differing only in case and whitespace
	a.Title = "Shared  Title"
	b.Title = "shared title"
	// b and c share a meta description and H1
	b.MetaDesc, c.MetaDesc = "Same description", "Same description"
	b.H1, c.H1 = []string{"Same H1"}, []string{"Same H1"}

	issues := AnalyzeCrossPage([]*models.PageResult{home, a, b, c})

	titles := issuesOfType(issues, IssueDuplicateTitle)
	if len(titles) != 2 {
		t.Fatalf("duplicate_title count = %d, want 2", len(titles))
	}
	if titles[0].URL != a.URL || titles[1].URL != b.URL {
		t.Errorf("duplicate_title URLs = %s, %s; want a then b", titles[0].URL, titles[1].URL)
	}
	if !strings.Contains(titles[0].Message, b.URL) {
		t.Errorf("duplicate_title message should name the other page, got %q", titles[0].Message)
	}
	if got := len(issuesOfType(issues, IssueDuplicateMetaDesc)); got != 2 {
		t.Errorf("duplicate_meta_description count = %d, want 2", got)
	}
	if got := len(issuesOfType(issues, IssueDuplicateH1)); got != 2 {
		t.Errorf("duplicate_h1 count = %d, want 2", got)
	}
}

func TestAnalyzeCrossPage_DuplicatesIgnoreNonIndexableAndEmpty(t *testing.T) {
	a := page("https://example.com/a", "https://example.com/b")
	b := page("https://example.com/b", "https://example.com/a")
	a.Title, b.Title = "Same", "Same"
	b.IndexabilityStatus = models.IndexabilityNoindex

	c := page("https://example.com/c", "https://example.com/a")
	d := page("https://example.com/d", "https://example.com/a")
	c.MetaDesc, d.MetaDesc = "", ""

	issues := AnalyzeCrossPage([]*models.PageResult{a, b, c, d})

	if got := len(issuesOfType(issues, IssueDuplicateTitle)); got != 0 {
		t.Errorf("duplicate_title count = %d, want 0 (one page is noindex)", got)
	}
	if got := len(issuesOfType(issues, IssueDuplicateMetaDesc)); got != 0 {
		t.Errorf("duplicate_meta_description count = %d, want 0 (empty values)", got)
	}
}

func TestAnalyzeCrossPage_OrphanPages(t *testing.T) {
	home := page("https://example.com/", "https://example.com/linked/")
	linked := page("https://example.com/linked", "https://example.com/")
	orphan := page("https://example.com/orphan", "https://example.com/")
	selfLinked := page("https://example.com/self", "https://example.com/self")
	noindexOrphan := page("https://example.com/hidden")
	noindexOrphan.IndexabilityStatus = models.IndexabilityNoindex

	issues := AnalyzeCrossPage([]*models.PageResult{home, linked, orphan, selfLinked, noindexOrphan})

	orphans := issuesOfType(issues, IssueOrphanPage)
	if len(orphans) != 2 {
		t.Fatalf("orphan_page count = %d, want 2 (orphan and self-linked): %+v", len(orphans), orphans)
	}
	if orphans[0].URL != orphan.URL || orphans[1].URL != selfLinked.URL {
		t.Errorf("orphan_page URLs = %s, %s", orphans[0].URL, orphans[1].URL)
	}
}

func TestAnalyzeCrossPage_NoLinkDataSkipsOrphans(t *testing.T) {
	issues := AnalyzeCrossPage([]*models.PageResult{
		page("https://example.com/a"),
		page("https://example.com/b"),
	})

	if got := len(issuesOfType(issues, IssueOrphanPage)); got != 0 {
		t.Errorf("orphan_page count = %d, want 0 when the crawl has no link data", got)
	}
}

func TestAnalyzeCrossPage_LinkTargets(t *testing.T) {
	home := page("https://example.com",
		"https://example.com/gone",
		"https://example.com/old",
		"https://example.com/fine",
		"https://example.com/not-crawled")
	gone := &models.PageResult{URL: "https://example.com/gone", StatusCode: 404, Error: "HTTP 404"}
	old := page("https://example.com/old", "https://example.com")
	old.RedirectChain = []string{"https://example.com/new"}
	fine := page("https://example.com/fine", "https://example.com")

	issues := AnalyzeCrossPage([]*models.PageResult{home, gone, old, fine})

	broken := issuesOfType(issues, IssueBrokenInternalLink)
	if len(broken) != 1 {
		t.Fatalf("broken_internal_link count = %d, want 1", len(broken))
	}
	if broken[0].URL != home.URL || broken[0].Value != gone.URL {
		t.Errorf("broken_internal_link = %+v, want source %s and target %s", broken[0], home.URL, gone.URL)
	}

	redirects := issuesOfType(issues, IssueLinkToRedirect)
	if len(redirects) != 1 {
		t.Fatalf("internal_link_to_redirect count = %d, want 1", len(redirects))
	}
	if redirects[0].URL != home.URL || !strings.Contains(redirects[0].Value, "https://example.com/new") {
		t.Errorf("internal_link_to_redirect = %+v", redirects[0])
	}
}

func TestIsRootURL(t *testing.T) {
	tests := map[string]bool{
		"https://example.com":         true,
		"https://example.com/":        true,
		"https://example.com/?ref=1":  true,
		"https://example.com/about":   false,
		"https://example.com/a/b/":    false,
		"https://example.com/about/#": false,
	}
	for url, want := range tests {
		if got := isRootURL(url); got != want {
			t.Errorf("isRootURL(%q) = %v, want %v", url, got, want)
		}
	}
}
