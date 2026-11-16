package analyzer

import (
	"fmt"
	"math/bits"
	"sort"
	"strconv"

	"github.com/dillonlara115/barracudaseo/pkg/models"
)

// Content issue types.
const (
	IssueThinContent          IssueType = "thin_content"
	IssueDuplicateContent     IssueType = "duplicate_content"
	IssueNearDuplicateContent IssueType = "near_duplicate_content"
)

const (
	// defaultThinContentWords is the main-text word count below which a page is
	// reported as thin.
	defaultThinContentWords = 200

	// nearDuplicateMaxDistance is the largest number of differing SimHash bits
	// (out of 64) at which two pages count as near-duplicates: 3 bits is about
	// 95% similar.
	nearDuplicateMaxDistance = 3

	// nearDuplicateMinWords keeps very short pages out of near-duplicate
	// comparison, where a similarity hash over a handful of words is noise.
	nearDuplicateMinWords = 50
)

// thinContentThreshold returns the minimum main-text word count for a page.
// A page with fewer words is reported as thin. Returning 0 exempts the page.
//
// TODO(dillon): a single threshold treats every page alike. Decide whether it
// should vary by page type, using p.URL. For example a contact, login or cart
// page is short by design (exempt it with 0), a blog post arguably needs more
// than a category listing, and a homepage is often image-led. The right
// numbers are an SEO judgement call your users will have opinions about.
func thinContentThreshold(p *models.PageResult) int {
	return defaultThinContentWords
}

// hasContentData reports whether a page's main text was analyzed. Results from
// older crawls, or pages with no text at all, carry no content hash.
func hasContentData(p *models.PageResult) bool {
	return p.ContentHash != ""
}

// AnalyzeContent reports thin pages and pages whose main text duplicates, or
// nearly duplicates, another page's.
func AnalyzeContent(results []*models.PageResult) []Issue {
	var issues []Issue
	issues = append(issues, findThinContent(results)...)

	exact, inExactGroup := findDuplicateContent(results)
	issues = append(issues, exact...)
	issues = append(issues, findNearDuplicateContent(results, inExactGroup)...)

	return issues
}

// findThinContent reports indexable pages with too little main text.
func findThinContent(results []*models.PageResult) []Issue {
	var issues []Issue
	for _, p := range results {
		if !isIndexablePage(p) || !hasContentData(p) {
			continue
		}
		// A client-rendered shell is short because its content needs
		// JavaScript; that is reported as a rendering issue, not thin content.
		if p.JS != nil && p.JS.ClientRendered {
			continue
		}
		threshold := thinContentThreshold(p)
		if threshold <= 0 || p.WordCount >= threshold {
			continue
		}
		issues = append(issues, Issue{
			Type:           IssueThinContent,
			Severity:       "warning",
			URL:            p.URL,
			Message:        fmt.Sprintf("Thin content: %d words of main text (under %d)", p.WordCount, threshold),
			Value:          strconv.Itoa(p.WordCount),
			Recommendation: "Expand the page so it fully covers its topic, merge it into a stronger page, or noindex it if it is not meant to rank",
		})
	}

	sortIssuesByURL(issues)
	return issues
}

// findDuplicateContent reports pages whose main text is identical to another
// page's. It also returns the set of pages it reported, so the near-duplicate
// pass can skip them.
func findDuplicateContent(results []*models.PageResult) ([]Issue, map[*models.PageResult]bool) {
	groups := make(map[string][]*models.PageResult)
	for _, p := range results {
		if !eligibleForDuplicateCheck(p) || !hasContentData(p) {
			continue
		}
		groups[p.ContentHash] = append(groups[p.ContentHash], p)
	}

	reported := make(map[*models.PageResult]bool)
	var issues []Issue
	for _, pages := range groups {
		if len(pages) < 2 {
			continue
		}
		for _, p := range pages {
			reported[p] = true
			others := make([]string, 0, len(pages)-1)
			for _, other := range pages {
				if other != p {
					others = append(others, other.URL)
				}
			}
			issues = append(issues, Issue{
				Type:           IssueDuplicateContent,
				Severity:       "warning",
				URL:            p.URL,
				Message:        fmt.Sprintf("Main content is identical to %s", describeURLs(others)),
				Value:          others[0],
				Recommendation: "Keep one version: redirect the others to it, or point their canonical tags at it",
			})
		}
	}

	sortIssuesByURL(issues)
	return issues, reported
}

// similarPage is a near-duplicate candidate for a page.
type similarPage struct {
	url        string
	similarity int // percent
}

// findNearDuplicateContent reports pages whose main text is almost, but not
// exactly, the same as another page's, using the Hamming distance between
// their similarity hashes.
func findNearDuplicateContent(results []*models.PageResult, skip map[*models.PageResult]bool) []Issue {
	type candidate struct {
		page *models.PageResult
		hash uint64
	}

	var candidates []candidate
	for _, p := range results {
		if skip[p] || !eligibleForDuplicateCheck(p) || !hasContentData(p) || p.WordCount < nearDuplicateMinWords {
			continue
		}
		hash, err := strconv.ParseUint(p.ContentSimHash, 16, 64)
		if err != nil {
			continue
		}
		candidates = append(candidates, candidate{page: p, hash: hash})
	}

	// Every pair is compared. A comparison is one XOR and a popcount, so even
	// 25,000 pages (about 312 million pairs) finish in about a second.
	matches := make(map[*models.PageResult][]similarPage)
	for i := 0; i < len(candidates); i++ {
		for j := i + 1; j < len(candidates); j++ {
			distance := bits.OnesCount64(candidates[i].hash ^ candidates[j].hash)
			if distance > nearDuplicateMaxDistance {
				continue
			}
			similarity := (64 - distance) * 100 / 64
			a, b := candidates[i].page, candidates[j].page
			matches[a] = append(matches[a], similarPage{url: b.URL, similarity: similarity})
			matches[b] = append(matches[b], similarPage{url: a.URL, similarity: similarity})
		}
	}

	var issues []Issue
	for p, similar := range matches {
		sort.Slice(similar, func(i, j int) bool {
			if similar[i].similarity != similar[j].similarity {
				return similar[i].similarity > similar[j].similarity
			}
			return similar[i].url < similar[j].url
		})
		closest := similar[0]

		message := fmt.Sprintf("Main content is about %d%% similar to %s", closest.similarity, closest.url)
		if len(similar) > 1 {
			message += fmt.Sprintf(" (and %d other page(s))", len(similar)-1)
		}
		issues = append(issues, Issue{
			Type:           IssueNearDuplicateContent,
			Severity:       "info",
			URL:            p.URL,
			Message:        message,
			Value:          closest.url,
			Recommendation: "Differentiate the pages with content specific to each, or consolidate them and canonicalise to one",
		})
	}

	sortIssuesByURL(issues)
	return issues
}
