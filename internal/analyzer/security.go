package analyzer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dillonlara115/barracudaseo/pkg/models"
)

// Security issue types.
const (
	IssueMixedContent          IssueType = "mixed_content"
	IssueInsecurePage          IssueType = "insecure_page"
	IssueMissingSecurityHeader IssueType = "missing_security_header"
)

// maxListedResources caps how many insecure resources a message names.
const maxListedResources = 3

// expectedSecurityHeaders are the response headers checked on every page, in
// report order, with why each matters.
var expectedSecurityHeaders = []struct {
	name   string // lowercase, as stored on the page
	label  string
	advice string
	// httpsOnly headers have no meaning on a plain http response
	httpsOnly bool
}{
	{"strict-transport-security", "Strict-Transport-Security", "tells browsers to always use HTTPS for this site", true},
	{"content-security-policy", "Content-Security-Policy", "limits which scripts and resources the page may load", false},
	{"x-content-type-options", "X-Content-Type-Options", "stops browsers guessing a file's type (set it to nosniff)", false},
	{"x-frame-options", "X-Frame-Options", "stops other sites embedding the page in a frame", false},
	{"referrer-policy", "Referrer-Policy", "controls how much of the URL is sent to other sites", false},
}

// AnalyzeSecurity reports pages served over http, https pages that load http
// resources, and security headers the site does not send.
func AnalyzeSecurity(results []*models.PageResult) []Issue {
	var issues []Issue
	var pages []*models.PageResult

	for _, p := range results {
		if !isHTMLPage(p) {
			continue
		}
		pages = append(pages, p)

		if strings.HasPrefix(strings.ToLower(p.URL), "http://") {
			issues = append(issues, Issue{
				Type:           IssueInsecurePage,
				Severity:       "error",
				URL:            p.URL,
				Message:        "Page is served over HTTP, not HTTPS",
				Recommendation: "Serve the page over HTTPS and redirect the HTTP URL to it with a 301",
			})
		}

		if len(p.MixedContent) > 0 {
			listed := p.MixedContent
			more := ""
			if len(listed) > maxListedResources {
				more = fmt.Sprintf(" (+%d more)", len(listed)-maxListedResources)
				listed = listed[:maxListedResources]
			}
			issues = append(issues, Issue{
				Type:           IssueMixedContent,
				Severity:       "warning",
				URL:            p.URL,
				Message:        fmt.Sprintf("HTTPS page loads %d resource(s) over HTTP: %s%s", len(p.MixedContent), strings.Join(listed, ", "), more),
				Value:          p.MixedContent[0],
				Recommendation: "Load every resource over HTTPS; browsers block or warn about insecure resources on a secure page",
			})
		}
	}

	issues = append(issues, findMissingSecurityHeaders(pages)...)
	return issues
}

// findMissingSecurityHeaders reports each expected header the site does not
// send. Headers are set by the server or CDN for the whole site, so a missing
// header is reported once, on one representative page, with a count, instead
// of once per page.
func findMissingSecurityHeaders(pages []*models.PageResult) []Issue {
	if len(pages) == 0 {
		return nil
	}

	// Report against the shortest URL (normally the homepage) so the issue
	// lands on the same page from crawl to crawl.
	ordered := make([]*models.PageResult, len(pages))
	copy(ordered, pages)
	sort.SliceStable(ordered, func(i, j int) bool {
		if len(ordered[i].URL) != len(ordered[j].URL) {
			return len(ordered[i].URL) < len(ordered[j].URL)
		}
		return ordered[i].URL < ordered[j].URL
	})

	var issues []Issue
	for _, header := range expectedSecurityHeaders {
		checked := 0
		var missing []*models.PageResult
		for _, p := range ordered {
			// nil means the crawl did not record headers for this page (older
			// crawls, imported results), which is not the same as none sent
			if p.SecurityHeaders == nil {
				continue
			}
			if header.httpsOnly && !strings.HasPrefix(strings.ToLower(p.URL), "https://") {
				continue
			}
			checked++
			if p.SecurityHeaders[header.name] == "" {
				missing = append(missing, p)
			}
		}
		if len(missing) == 0 {
			continue
		}

		issues = append(issues, Issue{
			Type:           IssueMissingSecurityHeader,
			Severity:       "info",
			URL:            missing[0].URL,
			Message:        fmt.Sprintf("%s header is missing on %d of %d pages", header.label, len(missing), checked),
			Value:          header.label,
			Recommendation: fmt.Sprintf("Send %s from the server or CDN: it %s", header.label, header.advice),
		})
	}
	return issues
}
