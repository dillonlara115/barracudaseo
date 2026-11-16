package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/supabase-community/postgrest-go"
	"go.uber.org/zap"

	"github.com/dillonlara115/barracudaseo/internal/analyzer"
	"github.com/dillonlara115/barracudaseo/pkg/models"
)

// maxCoverageItems caps each list in a crawl-vs-Google response. Counts are
// always complete.
const maxCoverageItems = 200

// gscPage is one page row from Search Console's performance data.
type gscPage struct {
	URL         string  `json:"url"`
	Clicks      float64 `json:"clicks"`
	Impressions float64 `json:"impressions"`
	Position    float64 `json:"position,omitempty"`
}

// coveragePage is a URL in a crawl-vs-Google list.
type coveragePage struct {
	URL         string  `json:"url"`
	Clicks      float64 `json:"clicks"`
	Impressions float64 `json:"impressions"`
	// Class and Reason describe the crawl's view of the URL (crawled pages only)
	Class  string `json:"class,omitempty"`
	Reason string `json:"reason,omitempty"`
}

// crawlGSCCoverage compares what a crawl found with what Google is showing.
type crawlGSCCoverage struct {
	CrawledPages   int `json:"crawled_pages"`
	IndexablePages int `json:"indexable_pages"`
	GooglePages    int `json:"google_pages"`
	// VisibleIndexablePages are indexable crawled pages Google is showing
	VisibleIndexablePages int `json:"visible_indexable_pages"`

	// SearchTrafficAtRisk: Google is showing these, but the crawl found they
	// are not indexable (noindex, blocked, redirecting, broken or
	// canonicalised elsewhere).
	SearchTrafficAtRiskCount int            `json:"search_traffic_at_risk_count"`
	SearchTrafficAtRisk      []coveragePage `json:"search_traffic_at_risk"`

	// NotInCrawl: Google is showing these, but the crawl never reached them.
	NotInCrawlCount int            `json:"not_in_crawl_count"`
	NotInCrawl      []coveragePage `json:"not_in_crawl"`

	// NoSearchVisibility: the crawl found these indexable, but Google showed
	// them for no query in the period.
	NoSearchVisibilityCount int            `json:"no_search_visibility_count"`
	NoSearchVisibility      []coveragePage `json:"no_search_visibility"`
}

// gscPageIsVisible decides whether a Search Console page row counts as Google
// showing the page.
//
// TODO(dillon): the baseline counts a single impression. That is the strictest
// reading of "Google knows this page", but one stray impression in 90 days can
// put a long tail of URLs into "search traffic at risk" and keep weak pages out
// of "no search visibility". Decide whether to require a floor (for example a
// minimum number of impressions, or at least one click).
func gscPageIsVisible(p gscPage) bool {
	return p.Impressions > 0
}

// coverageKey reduces a URL to the form used to match crawled pages with
// Search Console rows. Google may report a different scheme or www variant
// than the one crawled, and reports jump links (#section) as separate rows, so
// scheme, a leading "www.", the fragment and a trailing slash are ignored.
func coverageKey(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	u, err := url.Parse(rawURL)
	if err != nil || u.Host == "" {
		return strings.TrimRight(strings.ToLower(rawURL), "/")
	}
	key := strings.TrimPrefix(strings.ToLower(u.Host), "www.") + strings.TrimRight(u.EscapedPath(), "/")
	if u.RawQuery != "" {
		key += "?" + u.RawQuery
	}
	return key
}

// compareCrawlWithGSC joins crawled pages with Search Console page rows.
func compareCrawlWithGSC(pages []*models.PageResult, gscPages []gscPage) crawlGSCCoverage {
	coverage := crawlGSCCoverage{
		SearchTrafficAtRisk: []coveragePage{},
		NotInCrawl:          []coveragePage{},
		NoSearchVisibility:  []coveragePage{},
	}

	// Merge Search Console rows that are the same page (scheme, www and
	// fragment variants), keeping the URL of the row with the most impressions.
	google := make(map[string]gscPage)
	for _, row := range gscPages {
		if !gscPageIsVisible(row) {
			continue
		}
		key := coverageKey(row.URL)
		merged, exists := google[key]
		if !exists || row.Impressions > merged.Impressions {
			merged.URL = row.URL
		}
		if exists {
			merged.Clicks += row.Clicks
			merged.Impressions += row.Impressions
		} else {
			merged.Clicks, merged.Impressions = row.Clicks, row.Impressions
		}
		google[key] = merged
	}
	coverage.GooglePages = len(google)

	crawled := make(map[string]bool, len(pages))
	for _, p := range pages {
		key := coverageKey(p.URL)
		if crawled[key] {
			continue
		}
		crawled[key] = true
		coverage.CrawledPages++

		class := analyzer.ClassifyIndexability(p)
		row, inGoogle := google[key]

		if class == analyzer.ClassIndexable {
			coverage.IndexablePages++
			if inGoogle {
				coverage.VisibleIndexablePages++
			} else {
				coverage.NoSearchVisibility = append(coverage.NoSearchVisibility, coveragePage{URL: p.URL})
			}
			continue
		}

		if inGoogle {
			coverage.SearchTrafficAtRisk = append(coverage.SearchTrafficAtRisk, coveragePage{
				URL:         p.URL,
				Clicks:      row.Clicks,
				Impressions: row.Impressions,
				Class:       string(class),
				Reason:      analyzer.DescribeIndexability(p),
			})
		}
	}

	for key, row := range google {
		if !crawled[key] {
			coverage.NotInCrawl = append(coverage.NotInCrawl, coveragePage{
				URL:         row.URL,
				Clicks:      row.Clicks,
				Impressions: row.Impressions,
			})
		}
	}

	sortByTraffic(coverage.SearchTrafficAtRisk)
	sortByTraffic(coverage.NotInCrawl)
	sortByTraffic(coverage.NoSearchVisibility)

	coverage.SearchTrafficAtRiskCount = len(coverage.SearchTrafficAtRisk)
	coverage.NotInCrawlCount = len(coverage.NotInCrawl)
	coverage.NoSearchVisibilityCount = len(coverage.NoSearchVisibility)

	if len(coverage.SearchTrafficAtRisk) > maxCoverageItems {
		coverage.SearchTrafficAtRisk = coverage.SearchTrafficAtRisk[:maxCoverageItems]
	}
	if len(coverage.NotInCrawl) > maxCoverageItems {
		coverage.NotInCrawl = coverage.NotInCrawl[:maxCoverageItems]
	}
	if len(coverage.NoSearchVisibility) > maxCoverageItems {
		coverage.NoSearchVisibility = coverage.NoSearchVisibility[:maxCoverageItems]
	}

	return coverage
}

// sortByTraffic orders pages by clicks, then impressions, then URL.
func sortByTraffic(pages []coveragePage) {
	sort.SliceStable(pages, func(i, j int) bool {
		a, b := pages[i], pages[j]
		if a.Clicks != b.Clicks {
			return a.Clicks > b.Clicks
		}
		if a.Impressions != b.Impressions {
			return a.Impressions > b.Impressions
		}
		return a.URL < b.URL
	})
}

// fetchCrawlPagesForCoverage loads the fields of a crawl's pages needed to
// classify indexability, without the heavy link and image arrays.
func (s *Server) fetchCrawlPagesForCoverage(crawlID string) ([]*models.PageResult, error) {
	rows, err := s.fetchAllRows("pages",
		"id,url,status_code,indexability_status,canonical_url,redirect_chain:data->redirect_chain",
		func(q *postgrest.FilterBuilder) *postgrest.FilterBuilder {
			return q.Eq("crawl_id", crawlID)
		})
	if err != nil {
		return nil, err
	}

	pages := make([]*models.PageResult, 0, len(rows))
	for _, row := range rows {
		p := &models.PageResult{}
		p.URL, _ = row["url"].(string)
		if p.URL == "" {
			continue
		}
		p.StatusCode = int(getFloat(row["status_code"]))
		p.Canonical, _ = row["canonical_url"].(string)
		if status, ok := row["indexability_status"].(string); ok {
			p.IndexabilityStatus = models.IndexabilityStatus(status)
		}
		if chain, ok := row["redirect_chain"].([]interface{}); ok {
			for _, hop := range chain {
				if hopURL, ok := hop.(string); ok {
					p.RedirectChain = append(p.RedirectChain, hopURL)
				}
			}
		}
		pages = append(pages, p)
	}
	return pages, nil
}

// fetchGSCPageRows loads the page rows of a Search Console snapshot.
func (s *Server) fetchGSCPageRows(projectID, snapshotID string) ([]gscPage, error) {
	rows, err := s.fetchAllRows("gsc_performance_rows", "id,dimension_value,metrics",
		func(q *postgrest.FilterBuilder) *postgrest.FilterBuilder {
			return q.Eq("snapshot_id", snapshotID).Eq("project_id", projectID).Eq("row_type", "page")
		})
	if err != nil {
		return nil, err
	}

	pages := make([]gscPage, 0, len(rows))
	for _, row := range rows {
		pageURL, _ := row["dimension_value"].(string)
		metrics, _ := row["metrics"].(map[string]interface{})
		if pageURL == "" {
			continue
		}
		pages = append(pages, gscPage{
			URL:         pageURL,
			Clicks:      getFloat(metrics["clicks"]),
			Impressions: getFloat(metrics["impressions"]),
			Position:    getFloat(metrics["position"]),
		})
	}
	return pages, nil
}

// handleProjectGSCCrawlCoverage handles
// GET /api/v1/projects/:id/gsc/crawl-coverage?crawl_id=<id>. It compares a
// crawl with the project's latest Search Console data. Without crawl_id it uses
// the project's most recent succeeded crawl. Project access and the Pro
// subscription are checked by handleProjectGSC.
func (s *Server) handleProjectGSCCrawlCoverage(w http.ResponseWriter, r *http.Request, projectID string) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	snapshot, err := s.fetchLatestGSCSummary(projectID)
	if err != nil {
		s.logger.Error("Failed to load GSC snapshot", zap.Error(err), zap.String("project_id", projectID))
		s.respondError(w, http.StatusInternalServerError, "Failed to load Search Console data")
		return
	}
	if snapshot == nil {
		s.respondJSON(w, http.StatusOK, map[string]interface{}{
			"available": false,
			"reason":    "No Search Console data has been synced for this project yet.",
		})
		return
	}
	snapshotID, _ := snapshot["id"].(string)

	// The crawl must belong to this project
	query := s.serviceRole.From("crawls").
		Select("id,started_at,total_pages,meta", "", false).
		Eq("project_id", projectID)
	if crawlID := r.URL.Query().Get("crawl_id"); crawlID != "" {
		query = query.Eq("id", crawlID)
	} else {
		query = query.Eq("status", "succeeded").
			Order("started_at", &postgrest.OrderOpts{Ascending: false}).
			Limit(1, "")
	}
	data, _, err := query.Execute()
	var crawls []map[string]interface{}
	if err != nil || json.Unmarshal(data, &crawls) != nil {
		s.logger.Error("Failed to load crawl for GSC coverage", zap.Error(err), zap.String("project_id", projectID))
		s.respondError(w, http.StatusInternalServerError, "Failed to load crawl")
		return
	}
	if len(crawls) == 0 {
		s.respondError(w, http.StatusNotFound, "Crawl not found in this project")
		return
	}
	crawl := crawls[0]
	crawlID, _ := crawl["id"].(string)

	pages, err := s.fetchCrawlPagesForCoverage(crawlID)
	if err != nil {
		s.logger.Error("Failed to load crawl pages for GSC coverage", zap.Error(err), zap.String("crawl_id", crawlID))
		s.respondError(w, http.StatusInternalServerError, "Failed to load crawl pages")
		return
	}
	gscPages, err := s.fetchGSCPageRows(projectID, snapshotID)
	if err != nil {
		s.logger.Error("Failed to load GSC page rows", zap.Error(err), zap.String("project_id", projectID))
		s.respondError(w, http.StatusInternalServerError, "Failed to load Search Console data")
		return
	}

	// A crawl that stopped at its page limit did not see the whole site, so
	// "not in crawl" may include pages it simply never got to.
	crawlPartial := false
	if meta, ok := crawl["meta"].(map[string]interface{}); ok {
		if maxPages := getFloat(meta["max_pages"]); maxPages > 0 {
			crawlPartial = getFloat(crawl["total_pages"]) >= maxPages
		}
	}

	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"available":     true,
		"crawl_id":      crawlID,
		"crawl_partial": crawlPartial,
		"gsc": map[string]interface{}{
			"captured_on": snapshot["captured_on"],
			"period":      snapshot["period"],
		},
		"coverage": compareCrawlWithGSC(pages, gscPages),
	})
}
