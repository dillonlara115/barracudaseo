package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/supabase-community/postgrest-go"
	"go.uber.org/zap"

	"github.com/dillonlara115/barracudaseo/internal/analyzer"
	"github.com/dillonlara115/barracudaseo/internal/utils"
)

// maxDiffItems caps how many issues / URLs each list in a diff response carries.
// Counts are always complete.
const maxDiffItems = 500

// issueRef is the part of a stored issue needed to compare crawls.
type issueRef struct {
	Fingerprint string `json:"-"`
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	URL         string `json:"url"`
	Message     string `json:"message"`
	Value       string `json:"value,omitempty"`
	Status      string `json:"status,omitempty"`
}

// issueTypeDelta counts changes for one issue type.
type issueTypeDelta struct {
	New       int `json:"new"`
	Fixed     int `json:"fixed"`
	Unchanged int `json:"unchanged"`
}

// crawlDiff describes what changed between two crawls of the same project.
type crawlDiff struct {
	NewCount       int `json:"new_count"`
	FixedCount     int `json:"fixed_count"`
	UnchangedCount int `json:"unchanged_count"`
	// NotRecrawledCount counts previous issues whose page was not crawled this
	// time. They are neither fixed nor confirmed, so they are kept apart.
	NotRecrawledCount int `json:"not_recrawled_count"`

	NewIssues   []issueRef                `json:"new_issues"`
	FixedIssues []issueRef                `json:"fixed_issues"`
	ByType      map[string]issueTypeDelta `json:"by_type"`

	PagesAddedCount   int      `json:"pages_added_count"`
	PagesRemovedCount int      `json:"pages_removed_count"`
	PagesAdded        []string `json:"pages_added"`
	PagesRemoved      []string `json:"pages_removed"`

	Truncated bool `json:"truncated"`
}

// pageKey reduces a URL to the form used to match pages between crawls.
func pageKey(url string) string {
	return strings.TrimRight(strings.TrimSpace(url), "/")
}

// diffCrawls compares the issues and pages of a previous and a current crawl.
// An issue present before and absent now only counts as fixed if its page was
// crawled again; otherwise its state is unknown.
func diffCrawls(prevIssues, currIssues []issueRef, prevPages, currPages []string) crawlDiff {
	diff := crawlDiff{
		NewIssues:    []issueRef{},
		FixedIssues:  []issueRef{},
		ByType:       map[string]issueTypeDelta{},
		PagesAdded:   []string{},
		PagesRemoved: []string{},
	}

	prevPageSet := make(map[string]bool, len(prevPages))
	for _, u := range prevPages {
		prevPageSet[pageKey(u)] = true
	}
	currPageSet := make(map[string]bool, len(currPages))
	for _, u := range currPages {
		currPageSet[pageKey(u)] = true
	}

	for _, u := range currPages {
		if !prevPageSet[pageKey(u)] {
			diff.PagesAdded = append(diff.PagesAdded, u)
		}
	}
	for _, u := range prevPages {
		if !currPageSet[pageKey(u)] {
			diff.PagesRemoved = append(diff.PagesRemoved, u)
		}
	}
	diff.PagesAddedCount = len(diff.PagesAdded)
	diff.PagesRemovedCount = len(diff.PagesRemoved)

	prevByFP := make(map[string]bool, len(prevIssues))
	for _, issue := range prevIssues {
		prevByFP[issue.Fingerprint] = true
	}
	currByFP := make(map[string]bool, len(currIssues))
	for _, issue := range currIssues {
		currByFP[issue.Fingerprint] = true
	}

	countedCurr := make(map[string]bool, len(currIssues))
	for _, issue := range currIssues {
		if countedCurr[issue.Fingerprint] {
			continue
		}
		countedCurr[issue.Fingerprint] = true

		delta := diff.ByType[issue.Type]
		if prevByFP[issue.Fingerprint] {
			diff.UnchangedCount++
			delta.Unchanged++
		} else {
			diff.NewCount++
			delta.New++
			diff.NewIssues = append(diff.NewIssues, issue)
		}
		diff.ByType[issue.Type] = delta
	}

	countedPrev := make(map[string]bool, len(prevIssues))
	for _, issue := range prevIssues {
		if countedPrev[issue.Fingerprint] || currByFP[issue.Fingerprint] {
			continue
		}
		countedPrev[issue.Fingerprint] = true

		if !currPageSet[pageKey(issue.URL)] {
			diff.NotRecrawledCount++
			continue
		}
		delta := diff.ByType[issue.Type]
		delta.Fixed++
		diff.ByType[issue.Type] = delta
		diff.FixedCount++
		diff.FixedIssues = append(diff.FixedIssues, issue)
	}

	sortIssueRefs(diff.NewIssues)
	sortIssueRefs(diff.FixedIssues)
	sort.Strings(diff.PagesAdded)
	sort.Strings(diff.PagesRemoved)

	if len(diff.NewIssues) > maxDiffItems {
		diff.NewIssues, diff.Truncated = diff.NewIssues[:maxDiffItems], true
	}
	if len(diff.FixedIssues) > maxDiffItems {
		diff.FixedIssues, diff.Truncated = diff.FixedIssues[:maxDiffItems], true
	}
	if len(diff.PagesAdded) > maxDiffItems {
		diff.PagesAdded, diff.Truncated = diff.PagesAdded[:maxDiffItems], true
	}
	if len(diff.PagesRemoved) > maxDiffItems {
		diff.PagesRemoved, diff.Truncated = diff.PagesRemoved[:maxDiffItems], true
	}

	return diff
}

// severityRank orders severities from most to least serious.
func severityRank(severity string) int {
	switch severity {
	case "error":
		return 0
	case "warning":
		return 1
	default:
		return 2
	}
}

// sortIssueRefs orders issues by severity, then type, then URL.
func sortIssueRefs(issues []issueRef) {
	sort.SliceStable(issues, func(i, j int) bool {
		a, b := issues[i], issues[j]
		if severityRank(a.Severity) != severityRank(b.Severity) {
			return severityRank(a.Severity) < severityRank(b.Severity)
		}
		if a.Type != b.Type {
			return a.Type < b.Type
		}
		return a.URL < b.URL
	})
}

// fetchAllRows pages through a table query 1000 rows at a time (PostgREST's
// default cap), so large crawls are never silently truncated.
func (s *Server) fetchAllRows(table, columns string, filter func(*postgrest.FilterBuilder) *postgrest.FilterBuilder) ([]map[string]interface{}, error) {
	const chunkSize = 1000
	var rows []map[string]interface{}
	for offset := 0; ; offset += chunkSize {
		query := filter(s.serviceRole.From(table).Select(columns, "", false))
		data, _, err := query.Order("id", nil).Range(offset, offset+chunkSize-1, "").Execute()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch %s: %w", table, err)
		}
		var chunk []map[string]interface{}
		if err := json.Unmarshal(data, &chunk); err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", table, err)
		}
		rows = append(rows, chunk...)
		if len(chunk) < chunkSize {
			return rows, nil
		}
	}
}

// fetchCrawlPageURLs returns the URL of every page in a crawl, keyed by page ID.
func (s *Server) fetchCrawlPageURLs(crawlID string) (map[int64]string, error) {
	rows, err := s.fetchAllRows("pages", "id,url", func(q *postgrest.FilterBuilder) *postgrest.FilterBuilder {
		return q.Eq("crawl_id", crawlID)
	})
	if err != nil {
		return nil, err
	}
	urls := make(map[int64]string, len(rows))
	for _, row := range rows {
		id, ok := parsePageID(row["id"])
		url, _ := row["url"].(string)
		if ok && url != "" {
			urls[id] = url
		}
	}
	return urls, nil
}

// fetchCrawlIssueRefs returns a crawl's issues with their page URL and
// fingerprint. statuses, when non-empty, restricts the issues by status.
func (s *Server) fetchCrawlIssueRefs(crawlID string, pageURLs map[int64]string, statuses ...string) ([]issueRef, error) {
	rows, err := s.fetchAllRows("issues", "id,page_id,type,severity,message,value,status", func(q *postgrest.FilterBuilder) *postgrest.FilterBuilder {
		q = q.Eq("crawl_id", crawlID)
		if len(statuses) > 0 {
			q = q.In("status", statuses)
		}
		return q
	})
	if err != nil {
		return nil, err
	}

	refs := make([]issueRef, 0, len(rows))
	for _, row := range rows {
		ref := issueRef{}
		ref.Type, _ = row["type"].(string)
		ref.Severity, _ = row["severity"].(string)
		ref.Message, _ = row["message"].(string)
		ref.Value, _ = row["value"].(string)
		ref.Status, _ = row["status"].(string)
		if id, ok := parsePageID(row["page_id"]); ok {
			ref.URL = pageURLs[id]
		}
		if ref.URL == "" {
			// Issues stored without a page keep their URL at the end of value
			// ("<value> | <url>", or just "<url>").
			ref.URL = ref.Value
			if i := strings.LastIndex(ref.Value, " | "); i >= 0 {
				ref.URL, ref.Value = ref.Value[i+3:], ref.Value[:i]
			}
		}
		ref.Fingerprint = analyzer.IssueFingerprint(analyzer.IssueType(ref.Type), ref.URL, ref.Value)
		refs = append(refs, ref)
	}
	return refs, nil
}

// previousSucceededCrawl returns the most recent succeeded crawl of a project
// that started before the given time, or nil if there is none.
func (s *Server) previousSucceededCrawl(projectID, excludeCrawlID, startedBefore string) (map[string]interface{}, error) {
	query := s.serviceRole.From("crawls").
		Select("id,started_at,completed_at,total_pages,total_issues,source", "", false).
		Eq("project_id", projectID).
		Eq("status", "succeeded").
		Neq("id", excludeCrawlID)
	if startedBefore != "" {
		query = query.Lt("started_at", startedBefore)
	}
	data, _, err := query.
		Order("started_at", &postgrest.OrderOpts{Ascending: false}).
		Limit(1, "").
		Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to find previous crawl: %w", err)
	}
	var crawls []map[string]interface{}
	if err := json.Unmarshal(data, &crawls); err != nil {
		return nil, fmt.Errorf("failed to parse previous crawl: %w", err)
	}
	if len(crawls) == 0 {
		return nil, nil
	}
	return crawls[0], nil
}

// carriedIssueStatuses returns, by fingerprint, the triage status of issues in
// the project's last succeeded crawl that should survive a re-crawl. Without
// this every crawl would reset ignored and in-progress issues back to new.
// "fixed" is deliberately not carried: a fixed issue that shows up again is a
// regression and should surface as new.
func (s *Server) carriedIssueStatuses(projectID, currentCrawlID string) map[string]string {
	carried := map[string]string{}

	prev, err := s.previousSucceededCrawl(projectID, currentCrawlID, "")
	if err != nil || prev == nil {
		if err != nil {
			s.logger.Warn("Could not load previous crawl for status carry-over", zap.String("project_id", projectID), zap.Error(err))
		}
		return carried
	}
	prevID, _ := prev["id"].(string)

	pageURLs, err := s.fetchCrawlPageURLs(prevID)
	if err != nil {
		s.logger.Warn("Could not load previous pages for status carry-over", zap.String("crawl_id", prevID), zap.Error(err))
		return carried
	}
	refs, err := s.fetchCrawlIssueRefs(prevID, pageURLs, "ignored", "in_progress")
	if err != nil {
		s.logger.Warn("Could not load previous issues for status carry-over", zap.String("crawl_id", prevID), zap.Error(err))
		return carried
	}
	for _, ref := range refs {
		carried[ref.Fingerprint] = ref.Status
	}
	return carried
}

// issueIdentity returns the fingerprint used to de-duplicate an issue within a
// crawl and to match it against earlier crawls.
func issueIdentity(issue analyzer.Issue) string {
	url := issue.URL
	if normalized, err := utils.NormalizeURL(issue.URL); err == nil {
		url = normalized
	}
	return analyzer.IssueFingerprint(issue.Type, url, issue.Value)
}

// carriedStatus returns the status a newly stored issue should start with.
func carriedStatus(carried map[string]string, fingerprint string) string {
	if status, ok := carried[fingerprint]; ok {
		return status
	}
	return "new"
}

// handleCrawlDiff handles GET /api/v1/crawls/:id/diff - what changed since an
// earlier crawl. ?against=<crawl_id> picks the baseline; by default it is the
// project's previous succeeded crawl.
func (s *Server) handleCrawlDiff(w http.ResponseWriter, r *http.Request, crawlID string, userID string) {
	hasAccess, err := s.verifyCrawlAccess(userID, crawlID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			s.respondError(w, http.StatusNotFound, "Crawl not found")
		} else {
			s.logger.Error("Failed to verify crawl access", zap.String("crawl_id", crawlID), zap.Error(err))
			s.respondError(w, http.StatusInternalServerError, "Failed to verify crawl access")
		}
		return
	}
	if !hasAccess {
		s.respondError(w, http.StatusForbidden, "You don't have access to this crawl")
		return
	}

	data, _, err := s.serviceRole.From("crawls").
		Select("id,project_id,started_at,status", "", false).
		Eq("id", crawlID).
		Execute()
	var crawls []map[string]interface{}
	if err != nil || json.Unmarshal(data, &crawls) != nil || len(crawls) == 0 {
		s.logger.Error("Failed to load crawl for diff", zap.String("crawl_id", crawlID), zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to load crawl")
		return
	}
	projectID, _ := crawls[0]["project_id"].(string)
	startedAt, _ := crawls[0]["started_at"].(string)

	var baseline map[string]interface{}
	if against := r.URL.Query().Get("against"); against != "" {
		// The baseline must belong to the same project as the crawl being viewed
		baseData, _, err := s.serviceRole.From("crawls").
			Select("id,started_at,completed_at,total_pages,total_issues,source", "", false).
			Eq("id", against).
			Eq("project_id", projectID).
			Execute()
		var baseCrawls []map[string]interface{}
		if err != nil || json.Unmarshal(baseData, &baseCrawls) != nil {
			s.logger.Error("Failed to load baseline crawl", zap.String("crawl_id", against), zap.Error(err))
			s.respondError(w, http.StatusInternalServerError, "Failed to load baseline crawl")
			return
		}
		if len(baseCrawls) == 0 {
			s.respondError(w, http.StatusNotFound, "Baseline crawl not found in this project")
			return
		}
		baseline = baseCrawls[0]
	} else {
		baseline, err = s.previousSucceededCrawl(projectID, crawlID, startedAt)
		if err != nil {
			s.logger.Error("Failed to find previous crawl", zap.String("crawl_id", crawlID), zap.Error(err))
			s.respondError(w, http.StatusInternalServerError, "Failed to find previous crawl")
			return
		}
	}

	if baseline == nil {
		s.respondJSON(w, http.StatusOK, map[string]interface{}{
			"available": false,
			"reason":    "This is the first completed crawl for this project, so there is nothing to compare it with yet.",
		})
		return
	}
	baselineID, _ := baseline["id"].(string)

	currPageURLs, err := s.fetchCrawlPageURLs(crawlID)
	if err == nil {
		var prevPageURLs map[int64]string
		prevPageURLs, err = s.fetchCrawlPageURLs(baselineID)
		if err == nil {
			var currIssues, prevIssues []issueRef
			currIssues, err = s.fetchCrawlIssueRefs(crawlID, currPageURLs)
			if err == nil {
				prevIssues, err = s.fetchCrawlIssueRefs(baselineID, prevPageURLs)
				if err == nil {
					diff := diffCrawls(prevIssues, currIssues, mapValues(prevPageURLs), mapValues(currPageURLs))
					s.respondJSON(w, http.StatusOK, map[string]interface{}{
						"available": true,
						"crawl_id":  crawlID,
						"baseline":  baseline,
						"diff":      diff,
					})
					return
				}
			}
		}
	}

	s.logger.Error("Failed to build crawl diff", zap.String("crawl_id", crawlID), zap.String("baseline_id", baselineID), zap.Error(err))
	s.respondError(w, http.StatusInternalServerError, "Failed to compare crawls")
}

// mapValues returns the values of a page-ID-to-URL map.
func mapValues(m map[int64]string) []string {
	values := make([]string, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}
