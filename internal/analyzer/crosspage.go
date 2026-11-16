package analyzer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dillonlara115/barracudaseo/internal/utils"
	"github.com/dillonlara115/barracudaseo/pkg/models"
)

// Cross-page issue types. Unlike the per-page checks in Analyze, these compare
// pages against each other and need the whole crawl to be indexed first.
const (
	IssueDuplicateTitle     IssueType = "duplicate_title"
	IssueDuplicateMetaDesc  IssueType = "duplicate_meta_description"
	IssueDuplicateH1        IssueType = "duplicate_h1"
	IssueOrphanPage         IssueType = "orphan_page"
	IssueBrokenInternalLink IssueType = "broken_internal_link"
	IssueLinkToRedirect     IssueType = "internal_link_to_redirect"
)

// maxListedURLs caps how many other URLs are named in an issue message.
const maxListedURLs = 3

// crawlIndex is a lookup structure built once over the whole crawl.
type crawlIndex struct {
	byURL   map[string]*models.PageResult // urlKey -> page
	inlinks map[string][]string           // urlKey of target -> source page URLs
}

// urlKey reduces a URL to the form used for cross-page matching, so that
// "/about" and "/about/" (and a trailing fragment) refer to the same page.
func urlKey(rawURL string) string {
	if i := strings.Index(rawURL, "#"); i >= 0 {
		rawURL = rawURL[:i]
	}
	return strings.TrimRight(rawURL, "/")
}

// isHTMLPage reports whether a result is a successfully fetched HTML page.
func isHTMLPage(p *models.PageResult) bool {
	return p.StatusCode == 200 && p.Error == "" && !utils.IsImageURL(p.URL)
}

// isIndexablePage reports whether a result is an HTML page search engines may index.
func isIndexablePage(p *models.PageResult) bool {
	if !isHTMLPage(p) {
		return false
	}
	return p.IndexabilityStatus == models.IndexabilityIndexable || p.IndexabilityStatus == ""
}

// eligibleForDuplicateCheck decides whether a page takes part in duplicate
// title / description / H1 comparison.
//
// TODO(dillon): refine this. The baseline below only requires the page to be
// indexable. Decide whether to also exclude:
//   - pages canonicalised elsewhere (p.Canonical set and
//     urlKey(p.Canonical) != urlKey(p.URL)) — they legitimately share a title
//     with their canonical target;
//   - paginated URLs (?page=2, /page/2/) — often share a title by design.
//
// Excluding them cuts false positives; including them surfaces sites that rely
// on canonicals to paper over real duplication.
func eligibleForDuplicateCheck(p *models.PageResult) bool {
	return isIndexablePage(p)
}

// buildCrawlIndex indexes pages by URL and records which pages link to which.
func buildCrawlIndex(results []*models.PageResult) *crawlIndex {
	idx := &crawlIndex{
		byURL:   make(map[string]*models.PageResult, len(results)),
		inlinks: make(map[string][]string),
	}

	for _, p := range results {
		idx.byURL[urlKey(p.URL)] = p
	}

	for _, p := range results {
		if !isHTMLPage(p) {
			continue
		}
		source := urlKey(p.URL)
		seen := make(map[string]bool, len(p.InternalLinks))
		for _, link := range p.InternalLinks {
			target := urlKey(link)
			if target == source || seen[target] {
				continue
			}
			seen[target] = true
			idx.inlinks[target] = append(idx.inlinks[target], p.URL)
		}
	}

	return idx
}

// AnalyzeCrossPage detects issues that only show up when pages are compared
// with each other: duplicated metadata, orphan pages, and internal links that
// point at broken or redirecting URLs.
func AnalyzeCrossPage(results []*models.PageResult) []Issue {
	if len(results) < 2 {
		return nil
	}

	idx := buildCrawlIndex(results)

	var issues []Issue
	issues = append(issues, findDuplicates(results, IssueDuplicateTitle, "warning", "title",
		func(p *models.PageResult) string { return p.Title },
		"Write a unique title for each page that describes its specific content")...)
	issues = append(issues, findDuplicates(results, IssueDuplicateMetaDesc, "warning", "meta description",
		func(p *models.PageResult) string { return p.MetaDesc },
		"Write a unique meta description for each page")...)
	issues = append(issues, findDuplicates(results, IssueDuplicateH1, "info", "H1",
		func(p *models.PageResult) string {
			if len(p.H1) == 0 {
				return ""
			}
			return p.H1[0]
		},
		"Give each page a distinct H1 that matches its topic")...)
	issues = append(issues, findOrphanPages(results, idx)...)
	issues = append(issues, findLinkTargetIssues(results, idx)...)

	return issues
}

// normalizeText lowercases and collapses whitespace so that trivially
// different strings compare equal.
func normalizeText(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(s), " "))
}

// findDuplicates groups eligible pages by a text field and reports every page
// whose value is shared with at least one other page.
func findDuplicates(results []*models.PageResult, issueType IssueType, severity, label string,
	field func(*models.PageResult) string, recommendation string) []Issue {
	groups := make(map[string][]*models.PageResult)
	for _, p := range results {
		if !eligibleForDuplicateCheck(p) {
			continue
		}
		key := normalizeText(field(p))
		if key == "" {
			continue // missing values are reported by the per-page checks
		}
		groups[key] = append(groups[key], p)
	}

	var issues []Issue
	for _, pages := range groups {
		if len(pages) < 2 {
			continue
		}
		for _, p := range pages {
			others := make([]string, 0, len(pages)-1)
			for _, other := range pages {
				if other != p {
					others = append(others, other.URL)
				}
			}
			issues = append(issues, Issue{
				Type:           issueType,
				Severity:       severity,
				URL:            p.URL,
				Message:        fmt.Sprintf("Duplicate %s shared with %s", label, describeURLs(others)),
				Value:          strings.TrimSpace(field(p)),
				Recommendation: recommendation,
			})
		}
	}

	sortIssuesByURL(issues)
	return issues
}

// findOrphanPages reports indexable pages that no other crawled page links to.
// In a link-discovery crawl every page except the start URL has an inlink by
// construction, so this fires for pages found only through the sitemap.
func findOrphanPages(results []*models.PageResult, idx *crawlIndex) []Issue {
	// Without any link data (e.g. imported results) every page would look
	// orphaned, so there is nothing meaningful to report.
	if len(idx.inlinks) == 0 {
		return nil
	}

	var issues []Issue
	for _, p := range results {
		if !isIndexablePage(p) || isRootURL(p.URL) {
			continue
		}
		if len(idx.inlinks[urlKey(p.URL)]) > 0 {
			continue
		}
		issues = append(issues, Issue{
			Type:           IssueOrphanPage,
			Severity:       "warning",
			URL:            p.URL,
			Message:        fmt.Sprintf("No internal links point to this page from the %d pages crawled", len(results)),
			Recommendation: "Link to this page from a relevant page or navigation so users and search engines can reach it",
		})
	}

	sortIssuesByURL(issues)
	return issues
}

// findLinkTargetIssues reports, on the linking page, internal links whose
// target returned an error status or redirects elsewhere.
func findLinkTargetIssues(results []*models.PageResult, idx *crawlIndex) []Issue {
	var issues []Issue
	for _, p := range results {
		if !isHTMLPage(p) {
			continue
		}
		source := urlKey(p.URL)
		seen := make(map[string]bool, len(p.InternalLinks))
		for _, link := range p.InternalLinks {
			key := urlKey(link)
			if key == source || seen[key] {
				continue
			}
			seen[key] = true

			target, crawled := idx.byURL[key]
			if !crawled {
				continue
			}

			if target.StatusCode >= 400 {
				issues = append(issues, Issue{
					Type:           IssueBrokenInternalLink,
					Severity:       "error",
					URL:            p.URL,
					Message:        fmt.Sprintf("Links to a broken page (HTTP %d): %s", target.StatusCode, target.URL),
					Value:          target.URL,
					Recommendation: "Update the link to a working URL or remove it",
				})
				continue
			}

			if chain := validateRedirectChain(target.URL, target.RedirectChain); len(chain) > 0 {
				destination := chain[len(chain)-1]
				issues = append(issues, Issue{
					Type:           IssueLinkToRedirect,
					Severity:       "info",
					URL:            p.URL,
					Message:        fmt.Sprintf("Links to a redirecting URL: %s -> %s", target.URL, destination),
					Value:          target.URL + " -> " + destination,
					Recommendation: "Link directly to the final URL to avoid an extra redirect hop",
				})
			}
		}
	}

	return issues
}

// isRootURL reports whether a URL is a site's homepage.
func isRootURL(rawURL string) bool {
	rest := rawURL
	if i := strings.Index(rest, "://"); i >= 0 {
		rest = rest[i+3:]
	}
	if i := strings.IndexAny(rest, "?#"); i >= 0 {
		rest = rest[:i]
	}
	return !strings.Contains(strings.TrimRight(rest, "/"), "/")
}

// describeURLs renders a short, deterministic list of URLs for a message.
func describeURLs(urls []string) string {
	sort.Strings(urls)
	if len(urls) <= maxListedURLs {
		return fmt.Sprintf("%d other page(s): %s", len(urls), strings.Join(urls, ", "))
	}
	return fmt.Sprintf("%d other pages: %s (+%d more)",
		len(urls), strings.Join(urls[:maxListedURLs], ", "), len(urls)-maxListedURLs)
}

// sortIssuesByURL gives map-derived issues a stable order.
func sortIssuesByURL(issues []Issue) {
	sort.SliceStable(issues, func(i, j int) bool {
		return issues[i].URL < issues[j].URL
	})
}
