package api

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strings"

	"go.uber.org/zap"
)

type pageInsight struct {
	URL                 string                   `json:"url"`
	PriorityScore       float64                  `json:"priority_score"`
	Issues              []map[string]interface{} `json:"issues"`
	GSCMetrics          map[string]interface{}   `json:"gsc_metrics,omitempty"`
	GA4Metrics          map[string]interface{}   `json:"ga4_metrics,omitempty"`
	ClarityMetrics      map[string]interface{}   `json:"clarity_metrics,omitempty"`
	DataSources         []string                 `json:"data_sources"`
	Recommendations     []string                 `json:"recommendations"`
	IssueSeverityCounts map[string]int           `json:"issue_severity_counts"`
}

func (s *Server) handleProjectUnifiedInsights(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	hasAccess, err := s.verifyProjectAccess(userID, projectID)
	if err != nil {
		s.logger.Error("Failed to verify project access", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to verify access")
		return
	}
	if !hasAccess {
		s.respondError(w, http.StatusForbidden, "You don't have access to this project")
		return
	}

	// 1. Load latest crawl issues
	crawlIssues := s.loadLatestCrawlIssues(projectID)

	// 2. Load GSC page rows
	gscPages := s.loadDimensionRows("gsc_performance_rows", projectID, "page")

	// 3. Load GA4 page rows
	ga4Pages := s.loadDimensionRows("ga4_performance_rows", projectID, "page")

	// 4. Load Clarity URL rows
	clarityURLs := s.loadDimensionRows("clarity_performance_rows", projectID, "url")

	// Build page map keyed by normalized URL
	pageMap := make(map[string]*pageInsight)

	// Add issues
	for _, issue := range crawlIssues {
		url := normalizeInsightURL(issue["url"])
		if url == "" {
			continue
		}

		pi := getOrCreatePageInsight(pageMap, url)
		pi.Issues = append(pi.Issues, issue)

		severity, _ := issue["severity"].(string)
		pi.IssueSeverityCounts[severity]++
	}

	// Add GSC data
	for _, row := range gscPages {
		url := normalizeInsightURL(row["dimension_value"])
		if url == "" {
			continue
		}
		pi := getOrCreatePageInsight(pageMap, url)
		if metrics, ok := row["metrics"].(map[string]interface{}); ok {
			pi.GSCMetrics = metrics
			if !containsSource(pi.DataSources, "gsc") {
				pi.DataSources = append(pi.DataSources, "gsc")
			}
		}
	}

	// Add GA4 data
	for _, row := range ga4Pages {
		url := normalizeInsightURL(row["dimension_value"])
		if url == "" {
			continue
		}
		pi := getOrCreatePageInsight(pageMap, url)
		if metrics, ok := row["metrics"].(map[string]interface{}); ok {
			pi.GA4Metrics = metrics
			if !containsSource(pi.DataSources, "ga4") {
				pi.DataSources = append(pi.DataSources, "ga4")
			}
		}
	}

	// Add Clarity data
	for _, row := range clarityURLs {
		url := normalizeInsightURL(row["dimension_value"])
		if url == "" {
			continue
		}
		pi := getOrCreatePageInsight(pageMap, url)
		if metrics, ok := row["metrics"].(map[string]interface{}); ok {
			pi.ClarityMetrics = metrics
			if !containsSource(pi.DataSources, "clarity") {
				pi.DataSources = append(pi.DataSources, "clarity")
			}
		}
	}

	// Compute priority scores and generate recommendations
	for _, pi := range pageMap {
		pi.PriorityScore = computePriorityScore(pi)
		pi.Recommendations = generateRecommendations(pi)
	}

	// Sort by priority descending
	insights := make([]*pageInsight, 0, len(pageMap))
	for _, pi := range pageMap {
		// Only include pages that have at least one issue or significant data
		if len(pi.Issues) > 0 || pi.PriorityScore > 0 {
			insights = append(insights, pi)
		}
	}

	sort.Slice(insights, func(i, j int) bool {
		return insights[i].PriorityScore > insights[j].PriorityScore
	})

	// Compute summary
	totalIssuePages := 0
	highPriority := 0
	totalFrustration := 0
	totalScore := 0.0
	for _, pi := range insights {
		if len(pi.Issues) > 0 {
			totalIssuePages++
		}
		if pi.PriorityScore >= 50 {
			highPriority++
		}
		if pi.ClarityMetrics != nil {
			totalFrustration += int(getFloat(pi.ClarityMetrics["rage_click_count"]) + getFloat(pi.ClarityMetrics["dead_click_count"]))
		}
		totalScore += pi.PriorityScore
	}

	// Determine which data sources are connected
	connectedSources := []string{}
	if len(gscPages) > 0 {
		connectedSources = append(connectedSources, "gsc")
	}
	if len(ga4Pages) > 0 {
		connectedSources = append(connectedSources, "ga4")
	}
	if len(clarityURLs) > 0 {
		connectedSources = append(connectedSources, "clarity")
	}

	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"insights": insights,
		"summary": map[string]interface{}{
			"pages_with_issues":         totalIssuePages,
			"high_priority_fixes":       highPriority,
			"total_frustration_signals": totalFrustration,
			"opportunity_score":         math.Round(totalScore),
		},
		"connected_sources": connectedSources,
	})
}

func (s *Server) loadLatestCrawlIssues(projectID string) []map[string]interface{} {
	// Find the latest crawl for this project
	crawlData, _, err := s.serviceRole.
		From("crawls").
		Select("id", "", false).
		Eq("project_id", projectID).
		Eq("status", "succeeded").
		Execute()
	if err != nil {
		s.logger.Warn("Failed to load crawls for insights", zap.Error(err))
		return nil
	}

	var crawls []map[string]interface{}
	if err := json.Unmarshal(crawlData, &crawls); err != nil || len(crawls) == 0 {
		return nil
	}

	crawlID, _ := crawls[0]["id"].(string)
	if crawlID == "" {
		return nil
	}

	issueData, _, err := s.serviceRole.
		From("issues").
		Select("*", "", false).
		Eq("crawl_id", crawlID).
		Execute()
	if err != nil {
		s.logger.Warn("Failed to load issues for insights", zap.Error(err))
		return nil
	}

	var issues []map[string]interface{}
	if err := json.Unmarshal(issueData, &issues); err != nil {
		return nil
	}

	return issues
}

func (s *Server) loadDimensionRows(table, projectID, rowType string) []map[string]interface{} {
	data, _, err := s.serviceRole.
		From(table).
		Select("*", "", false).
		Eq("project_id", projectID).
		Eq("row_type", rowType).
		Execute()
	if err != nil {
		errStr := err.Error()
		if !strings.Contains(errStr, "does not exist") && !strings.Contains(errStr, "relation") {
			s.logger.Warn("Failed to load dimension rows", zap.String("table", table), zap.Error(err))
		}
		return nil
	}

	var rows []map[string]interface{}
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil
	}

	return rows
}

func normalizeInsightURL(raw interface{}) string {
	url, _ := raw.(string)
	if url == "" {
		return ""
	}
	// Remove protocol and trailing slash for matching
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimSuffix(url, "/")
	url = strings.ToLower(url)
	return url
}

func getOrCreatePageInsight(m map[string]*pageInsight, url string) *pageInsight {
	if pi, ok := m[url]; ok {
		return pi
	}
	pi := &pageInsight{
		URL:                 url,
		Issues:              []map[string]interface{}{},
		DataSources:         []string{},
		Recommendations:     []string{},
		IssueSeverityCounts: map[string]int{},
	}
	m[url] = pi
	return pi
}

func containsSource(sources []string, source string) bool {
	for _, s := range sources {
		if s == source {
			return true
		}
	}
	return false
}

func computePriorityScore(pi *pageInsight) float64 {
	// Base score from issues
	base := 0.0
	for _, issue := range pi.Issues {
		severity, _ := issue["severity"].(string)
		switch severity {
		case "error":
			base += 10
		case "warning":
			base += 5
		case "info":
			base += 1
		}
	}

	if base == 0 && pi.GSCMetrics == nil && pi.GA4Metrics == nil && pi.ClarityMetrics == nil {
		return 0
	}

	// Search factor from GSC
	searchFactor := 1.0
	if pi.GSCMetrics != nil {
		impressions := getFloat(pi.GSCMetrics["impressions"])
		position := getFloat(pi.GSCMetrics["position"])
		if impressions > 10000 {
			searchFactor = 3.0
		} else if impressions > 1000 {
			searchFactor = 2.0
		} else if impressions > 100 {
			searchFactor = 1.5
		}
		// Pages ranking 5-20 have optimization opportunity
		if position >= 5 && position <= 20 && impressions > 500 {
			searchFactor *= 1.3
		}
	}

	// Traffic factor from GA4
	trafficFactor := 1.0
	if pi.GA4Metrics != nil {
		sessions := getFloat(pi.GA4Metrics["sessions"])
		bounceRate := getFloat(pi.GA4Metrics["bounce_rate"])
		conversions := getFloat(pi.GA4Metrics["conversions"])
		if sessions > 10000 {
			trafficFactor = 3.0
		} else if sessions > 1000 {
			trafficFactor = 2.0
		} else if sessions > 100 {
			trafficFactor = 1.5
		}
		if bounceRate > 0.7 && sessions > 500 {
			trafficFactor *= 1.3
		}
		if conversions > 0 {
			trafficFactor *= 1.5
		}
	}

	// UX factor from Clarity
	uxFactor := 1.0
	if pi.ClarityMetrics != nil {
		rageClicks := getFloat(pi.ClarityMetrics["rage_click_count"])
		deadClicks := getFloat(pi.ClarityMetrics["dead_click_count"])
		quickbacks := getFloat(pi.ClarityMetrics["quickback_click"])
		frustration := rageClicks + deadClicks + quickbacks
		if frustration > 50 {
			uxFactor = 2.0
		} else if frustration > 20 {
			uxFactor = 1.5
		} else if frustration > 5 {
			uxFactor = 1.2
		}
	}

	// Ensure minimum score of 1 for base if there are data sources but no issues
	if base == 0 && len(pi.DataSources) > 0 {
		base = 1
	}

	return base * searchFactor * trafficFactor * uxFactor
}

func generateRecommendations(pi *pageInsight) []string {
	var recs []string

	// Issue-based recommendations with data citations
	issueTypes := map[string]int{}
	for _, issue := range pi.Issues {
		issueType, _ := issue["type"].(string)
		issueTypes[issueType]++
	}

	for issueType, count := range issueTypes {
		rec := fmt.Sprintf("Fix %d \"%s\" issue(s) on /%s", count, issueType, pi.URL)

		// Add data citations
		citations := []string{}
		if pi.GSCMetrics != nil {
			impressions := getFloat(pi.GSCMetrics["impressions"])
			if impressions > 0 {
				citations = append(citations, fmt.Sprintf("%.0f impressions (GSC)", impressions))
			}
		}
		if pi.GA4Metrics != nil {
			sessions := getFloat(pi.GA4Metrics["sessions"])
			if sessions > 0 {
				citations = append(citations, fmt.Sprintf("%.0f sessions (GA4)", sessions))
			}
		}
		if pi.ClarityMetrics != nil {
			rageClicks := getFloat(pi.ClarityMetrics["rage_click_count"])
			if rageClicks > 0 {
				citations = append(citations, fmt.Sprintf("%.0f rage clicks (Clarity)", rageClicks))
			}
		}

		if len(citations) > 0 {
			rec += " — " + strings.Join(citations, ", ")
		}

		recs = append(recs, rec)
	}

	// High bounce rate recommendation
	if pi.GA4Metrics != nil {
		bounceRate := getFloat(pi.GA4Metrics["bounce_rate"])
		sessions := getFloat(pi.GA4Metrics["sessions"])
		if bounceRate > 0.7 && sessions > 500 {
			recs = append(recs, fmt.Sprintf("High bounce rate (%.0f%%) with %.0f sessions — investigate page content and UX", bounceRate*100, sessions))
		}
	}

	// UX frustration recommendation
	if pi.ClarityMetrics != nil {
		rageClicks := getFloat(pi.ClarityMetrics["rage_click_count"])
		deadClicks := getFloat(pi.ClarityMetrics["dead_click_count"])
		if rageClicks > 10 || deadClicks > 10 {
			recs = append(recs, fmt.Sprintf("High frustration signals: %.0f rage clicks, %.0f dead clicks — review interactive elements", rageClicks, deadClicks))
		}
	}

	// Ranking opportunity
	if pi.GSCMetrics != nil {
		position := getFloat(pi.GSCMetrics["position"])
		impressions := getFloat(pi.GSCMetrics["impressions"])
		if position >= 5 && position <= 20 && impressions > 1000 {
			recs = append(recs, fmt.Sprintf("Ranking opportunity: position %.1f with %.0f impressions — optimize on-page SEO to move into top 5", position, impressions))
		}
	}

	return recs
}
