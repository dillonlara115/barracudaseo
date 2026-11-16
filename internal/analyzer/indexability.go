package analyzer

import (
	"fmt"
	"strings"

	"github.com/dillonlara115/barracudaseo/internal/utils"
	"github.com/dillonlara115/barracudaseo/pkg/models"
)

// Indexability issue types.
const (
	IssueNoindexPage           IssueType = "noindex_page"
	IssueBlockedByRobots       IssueType = "blocked_by_robots"
	IssueCanonicalised         IssueType = "canonical_points_elsewhere"
	IssueCanonicalNonIndexable IssueType = "canonical_to_non_indexable"
	IssueNonIndexableInSitemap IssueType = "non_indexable_in_sitemap"
	IssueNotInSitemap          IssueType = "indexable_not_in_sitemap"
)

// IndexabilityClass is the single reason a URL is, or is not, indexable.
// Each crawled URL gets exactly one class, so the classes sum to the number
// of URLs analyzed.
type IndexabilityClass string

const (
	ClassIndexable     IndexabilityClass = "indexable"
	ClassError         IndexabilityClass = "error"         // 4xx/5xx or fetch failure
	ClassRedirect      IndexabilityClass = "redirect"      // URL redirects elsewhere
	ClassBlocked       IndexabilityClass = "blocked"       // disallowed by robots.txt
	ClassNoindex       IndexabilityClass = "noindex"       // meta robots or X-Robots-Tag
	ClassCanonicalised IndexabilityClass = "canonicalised" // canonical points to another URL
)

// canonicalTarget returns the matching key of the URL a page's canonical tag
// points to, resolving relative canonicals. It returns "" when the page has no
// canonical or the canonical cannot be parsed.
func canonicalTarget(p *models.PageResult) string {
	if strings.TrimSpace(p.Canonical) == "" {
		return ""
	}
	resolved, err := utils.ResolveURL(p.URL, p.Canonical)
	if err != nil {
		return ""
	}
	return urlKey(resolved)
}

// classifyIndexability assigns a URL its indexability class. The order matters:
// it follows what a search engine hits first (the response, then robots rules,
// then on-page directives).
func classifyIndexability(p *models.PageResult) IndexabilityClass {
	switch {
	case p.StatusCode >= 400 || (p.Error != "" && p.StatusCode != 200):
		return ClassError
	case len(validateRedirectChain(p.URL, p.RedirectChain)) > 0:
		return ClassRedirect
	case p.IndexabilityStatus == models.IndexabilityBlocked:
		return ClassBlocked
	case p.IndexabilityStatus == models.IndexabilityNoindex:
		return ClassNoindex
	}
	if target := canonicalTarget(p); target != "" && target != urlKey(p.URL) {
		return ClassCanonicalised
	}
	return ClassIndexable
}

// describeClass explains a non-indexable class in a sentence fragment.
func describeClass(p *models.PageResult, class IndexabilityClass) string {
	switch class {
	case ClassError:
		if p.StatusCode > 0 {
			return fmt.Sprintf("returns HTTP %d", p.StatusCode)
		}
		return "could not be fetched"
	case ClassRedirect:
		return "redirects to " + p.RedirectChain[len(p.RedirectChain)-1]
	case ClassBlocked:
		return "is blocked by robots.txt"
	case ClassNoindex:
		return "is set to noindex"
	case ClassCanonicalised:
		return "is canonicalised to " + p.Canonical
	default:
		return "is indexable"
	}
}

// noindexSource names where a noindex directive comes from.
func noindexSource(p *models.PageResult) string {
	inHeader := strings.Contains(strings.ToLower(p.XRobotsTag), "noindex")
	inMeta := strings.Contains(strings.ToLower(p.MetaRobots), "noindex")
	switch {
	case inHeader && inMeta:
		return "meta robots tag and X-Robots-Tag header"
	case inHeader:
		return "X-Robots-Tag header"
	case inMeta:
		return "meta robots tag"
	default:
		return "robots directive"
	}
}

// noindexSeverity decides how loudly to report a noindex page. Most noindex
// pages are intentional (thank-you pages, internal search, tag archives), but
// an accidental one silently removes a page from search.
//
// TODO(dillon): replace the flat baseline with a heuristic for "probably
// accidental". Signals available here:
//   - inlinks: number of crawled pages linking to this one. A page many pages
//     link to is unlikely to be meant as hidden.
//   - p.InSitemap: listed in the sitemap, which asks for indexing. (That
//     conflict is also reported separately as non_indexable_in_sitemap.)
//   - p.URL: a shallow path like /pricing is more suspicious than /cart.
//
// Return "info", "warning" or "error".
func noindexSeverity(p *models.PageResult, inlinks int) string {
	return "info"
}

// AnalyzeIndexability classifies every crawled URL and reports indexability
// problems: noindex and robots-blocked pages, canonicals pointing elsewhere or
// at non-indexable URLs, and disagreement between the sitemap and the crawl.
// It returns the issues and a count of URLs per IndexabilityClass.
func AnalyzeIndexability(results []*models.PageResult) ([]Issue, map[IndexabilityClass]int) {
	idx := buildCrawlIndex(results)
	breakdown := make(map[IndexabilityClass]int)
	classes := make(map[*models.PageResult]IndexabilityClass, len(results))

	hasSitemapData := false
	for _, p := range results {
		if utils.IsImageURL(p.URL) {
			continue
		}
		class := classifyIndexability(p)
		classes[p] = class
		breakdown[class]++
		if p.InSitemap {
			hasSitemapData = true
		}
	}

	var issues []Issue
	for _, p := range results {
		class, analyzed := classes[p]
		if !analyzed {
			continue
		}

		switch class {
		case ClassNoindex:
			issues = append(issues, Issue{
				Type:           IssueNoindexPage,
				Severity:       noindexSeverity(p, len(idx.inlinks[urlKey(p.URL)])),
				URL:            p.URL,
				Message:        fmt.Sprintf("Page is excluded from search by a noindex in the %s", noindexSource(p)),
				Value:          strings.TrimSpace(strings.Join(nonEmpty(p.MetaRobots, p.XRobotsTag), " | ")),
				Recommendation: "Confirm this page is meant to stay out of search results; remove the noindex if not",
			})

		case ClassBlocked:
			issues = append(issues, Issue{
				Type:           IssueBlockedByRobots,
				Severity:       "warning",
				URL:            p.URL,
				Message:        "Page is disallowed by robots.txt, so search engines cannot crawl it",
				Recommendation: "Confirm the robots.txt rule is intended. A blocked page can still be indexed without its content; use noindex instead to keep a page out of search",
			})

		case ClassCanonicalised:
			issues = append(issues, findCanonicalIssues(p, idx, classes)...)
		}

		if hasSitemapData {
			if p.InSitemap && class != ClassIndexable {
				issues = append(issues, Issue{
					Type:           IssueNonIndexableInSitemap,
					Severity:       "warning",
					URL:            p.URL,
					Message:        fmt.Sprintf("Listed in the sitemap but %s", describeClass(p, class)),
					Value:          string(class),
					Recommendation: "List only indexable, canonical URLs in the sitemap: fix the page or remove it from the sitemap",
				})
			} else if !p.InSitemap && class == ClassIndexable {
				issues = append(issues, Issue{
					Type:           IssueNotInSitemap,
					Severity:       "info",
					URL:            p.URL,
					Message:        "Indexable page is not listed in the sitemap",
					Recommendation: "Add the page to the sitemap, or noindex it if it should not appear in search",
				})
			}
		}
	}

	return issues, breakdown
}

// findCanonicalIssues reports a page whose canonical points to another URL,
// escalating when that target is known not to be indexable.
func findCanonicalIssues(p *models.PageResult, idx *crawlIndex, classes map[*models.PageResult]IndexabilityClass) []Issue {
	target, crawled := idx.byURL[canonicalTarget(p)]
	if crawled {
		if targetClass := classes[target]; targetClass != ClassIndexable && targetClass != "" {
			return []Issue{{
				Type:           IssueCanonicalNonIndexable,
				Severity:       "error",
				URL:            p.URL,
				Message:        fmt.Sprintf("Canonical points to %s, which %s", target.URL, describeClass(target, targetClass)),
				Value:          p.Canonical,
				Recommendation: "Point the canonical at an indexable URL that returns 200, or make this page self-canonical",
			}}
		}
	}

	message := "Canonical points to a different URL, so this page will not be indexed itself"
	if !crawled {
		message += " (the target was not crawled, so its status is unknown)"
	}
	return []Issue{{
		Type:           IssueCanonicalised,
		Severity:       "info",
		URL:            p.URL,
		Message:        message,
		Value:          p.Canonical,
		Recommendation: "Confirm this page is a duplicate of the canonical target; otherwise make it self-canonical",
	}}
}

// nonEmpty returns the non-blank strings from its arguments.
func nonEmpty(values ...string) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			out = append(out, v)
		}
	}
	return out
}
