package analyzer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dillonlara115/barracudaseo/pkg/models"
)

// Hreflang issue types.
const (
	IssueHreflangInvalidLanguage   IssueType = "hreflang_invalid_language"
	IssueHreflangMissingSelf       IssueType = "hreflang_missing_self_reference"
	IssueHreflangMissingReturnLink IssueType = "hreflang_missing_return_link"
	IssueHreflangToNonIndexable    IssueType = "hreflang_to_non_indexable"
	IssueHreflangMissingXDefault   IssueType = "hreflang_missing_x_default"
)

// hreflangCodePattern matches the shape of a valid hreflang value: a language
// (ISO 639-1, two letters), optionally a script (four letters), optionally a
// region (ISO 3166-1 alpha-2). It checks the shape only, not that the codes
// are assigned or in the right order: "en_US" and "english" are rejected, but
// other unassigned or swapped codes (such as "xx-YY") are not caught.
var hreflangCodePattern = regexp.MustCompile(`^[a-z]{2}(-[a-z]{4})?(-[a-z]{2})?$`)

// validHreflangCode reports whether a hreflang value is well formed.
func validHreflangCode(code string) bool {
	code = strings.ToLower(strings.TrimSpace(code))
	if code == "x-default" {
		return true
	}
	// "UK" is not an ISO 3166-1 region (the United Kingdom is "GB"). This is
	// the most common hreflang mistake, so it is checked by name.
	if strings.HasSuffix(code, "-uk") {
		return false
	}
	return hreflangCodePattern.MatchString(code)
}

// AnalyzeHreflang checks hreflang annotations: malformed language codes, a
// page that leaves itself out of its own set, alternates that do not link
// back, alternates that cannot be indexed, and a set with no x-default.
// Reciprocity can only be checked for alternates that were crawled.
func AnalyzeHreflang(results []*models.PageResult) []Issue {
	byURL := make(map[string]*models.PageResult, len(results))
	for _, p := range results {
		byURL[urlKey(p.URL)] = p
	}

	var issues []Issue
	for _, p := range results {
		if !isHTMLPage(p) || len(p.Hreflang) == 0 {
			continue
		}
		self := urlKey(p.URL)

		hasSelf := false
		hasXDefault := false
		reported := make(map[string]bool)

		for _, link := range p.Hreflang {
			code := strings.ToLower(strings.TrimSpace(link.Lang))
			target := urlKey(link.URL)

			if code == "x-default" {
				hasXDefault = true
			}
			if target == self {
				hasSelf = true
			}

			if !validHreflangCode(code) && !reported["code:"+code] {
				reported["code:"+code] = true
				issues = append(issues, Issue{
					Type:           IssueHreflangInvalidLanguage,
					Severity:       "error",
					URL:            p.URL,
					Message:        fmt.Sprintf("Invalid hreflang value %q (for %s)", link.Lang, link.URL),
					Value:          link.Lang,
					Recommendation: "Use an ISO 639-1 language code, optionally with an ISO 3166-1 region, such as \"en\" or \"en-GB\"; search engines ignore values they do not recognise",
				})
			}

			if target == self || reported["target:"+target] {
				continue
			}
			alternate, crawled := byURL[target]
			if !crawled {
				continue
			}
			reported["target:"+target] = true

			if class := classifyIndexability(alternate); class != ClassIndexable {
				issues = append(issues, Issue{
					Type:           IssueHreflangToNonIndexable,
					Severity:       "error",
					URL:            p.URL,
					Message:        fmt.Sprintf("hreflang %q points to %s, which %s", link.Lang, alternate.URL, describeClass(alternate, class)),
					Value:          alternate.URL,
					Recommendation: "Point hreflang at the final, indexable URL of each language version",
				})
				continue
			}

			linksBack := false
			for _, back := range alternate.Hreflang {
				if urlKey(back.URL) == self {
					linksBack = true
					break
				}
			}
			if !linksBack {
				issues = append(issues, Issue{
					Type:           IssueHreflangMissingReturnLink,
					Severity:       "error",
					URL:            p.URL,
					Message:        fmt.Sprintf("hreflang %q points to %s, which does not link back to this page", link.Lang, alternate.URL),
					Value:          alternate.URL,
					Recommendation: "Add a matching hreflang annotation on the alternate page; annotations that are not reciprocal are ignored",
				})
			}
		}

		if !hasSelf {
			issues = append(issues, Issue{
				Type:           IssueHreflangMissingSelf,
				Severity:       "warning",
				URL:            p.URL,
				Message:        "Page has hreflang annotations but does not list itself",
				Recommendation: "Add a hreflang annotation for this page's own language pointing at its own URL",
			})
		}
		if !hasXDefault {
			issues = append(issues, Issue{
				Type:           IssueHreflangMissingXDefault,
				Severity:       "info",
				URL:            p.URL,
				Message:        "hreflang set has no x-default",
				Recommendation: "Add hreflang=\"x-default\" for the page shown to users whose language matches none of the alternates",
			})
		}
	}

	return issues
}
