package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

// handleImpactFirstView handles GET /api/v1/projects/:id/impact-first
// Returns pages that rank for keywords AND have crawl issues, prioritized by impact
func (s *Server) handleImpactFirstView(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	hasAccess, err := s.verifyProjectAccess(userID, projectID)
	if err != nil {
		s.logger.Error("Failed to verify project access", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to verify project access")
		return
	}
	if !hasAccess {
		s.respondError(w, http.StatusForbidden, "You don't have access to this project")
		return
	}

	// Get all keyword snapshots with crawl_page_id linked
	// Filter out null crawl_page_id values by checking in code (Supabase doesn't support Not with nil directly)
	var snapshots []map[string]interface{}
	snapshotData, _, err := s.serviceRole.From("keyword_rank_snapshots").
		Select("crawl_page_id, position_organic, serp_url, keyword_id, checked_at", "", false).
		Order("checked_at", nil).
		Execute()

	if err != nil {
		s.logger.Error("Failed to fetch keyword snapshots", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to fetch snapshots")
		return
	}

	if err := json.Unmarshal(snapshotData, &snapshots); err != nil {
		s.logger.Error("Failed to parse snapshots", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to parse snapshots")
		return
	}

	// Filter out snapshots without crawl_page_id and group by page_id
	// Group by page_id and get best (lowest) position
	pageRankings := make(map[int64]*PageRanking)
	for _, snap := range snapshots {
		// crawl_page_id is bigint, comes as float64 from JSON
		var pageID int64
		if pageIDFloat, ok := snap["crawl_page_id"].(float64); ok {
			pageID = int64(pageIDFloat)
		}
		if pageID == 0 {
			continue // Skip snapshots without crawl_page_id
		}

		posOrg := 0
		if pos, ok := snap["position_organic"].(float64); ok {
			posOrg = int(pos)
		}

		// Get keyword info
		keywordID := getKeywordString(snap, "keyword_id")
		var keywordText string
		if keywordID != "" {
			keywordData, _, err := s.serviceRole.From("keywords").Select("keyword", "", false).Eq("id", keywordID).Execute()
			if err == nil {
				var keywords []map[string]interface{}
				if err := json.Unmarshal(keywordData, &keywords); err == nil && len(keywords) > 0 {
					keywordText = getKeywordString(keywords[0], "keyword")
				}
			}
		}

		if ranking, exists := pageRankings[pageID]; exists {
			// Update if this is a better position
			if posOrg > 0 && (ranking.BestPosition == 0 || posOrg < ranking.BestPosition) {
				ranking.BestPosition = posOrg
			}
			ranking.KeywordCount++
			// Add keyword if not already in list
			if keywordText != "" {
				keywordExists := false
				for _, k := range ranking.Keywords {
					if k == keywordText {
						keywordExists = true
						break
					}
				}
				if !keywordExists {
					ranking.Keywords = append(ranking.Keywords, keywordText)
				}
			}
		} else {
			keywords := []string{}
			if keywordText != "" {
				keywords = []string{keywordText}
			}
			pageRankings[pageID] = &PageRanking{
				PageID:        pageID,
				BestPosition:  posOrg,
				KeywordCount:  1,
				Keywords:      keywords,
			}
		}
	}

	// Get latest crawl for project
	var crawls []map[string]interface{}
	crawlData, _, err := s.serviceRole.From("crawls").
		Select("id", "", false).
		Eq("project_id", projectID).
		Order("started_at", nil).
		Limit(1, "").
		Execute()
	
	var latestCrawlID string
	if err == nil {
		json.Unmarshal(crawlData, &crawls)
		if len(crawls) > 0 {
			latestCrawlID = getKeywordString(crawls[0], "id")
		}
	}

	// Fetch pages from latest crawl (or all pages if no crawl)
	var pages []map[string]interface{}
	var pageData []byte
	if latestCrawlID != "" {
		pageData, _, err = s.serviceRole.From("pages").
			Select("*", "", false).
			Eq("crawl_id", latestCrawlID).
			Execute()
	} else {
		// Fallback: get all crawls for project, then get pages from those crawls
		var allCrawls []map[string]interface{}
		allCrawlData, _, err := s.serviceRole.From("crawls").
			Select("id", "", false).
			Eq("project_id", projectID).
			Execute()
		if err == nil {
			json.Unmarshal(allCrawlData, &allCrawls)
			// Get pages from all crawls (simplified - in production might want to optimize)
			for _, crawl := range allCrawls {
				crawlID := getKeywordString(crawl, "id")
				crawlPageData, _, err := s.serviceRole.From("pages").
					Select("*", "", false).
					Eq("crawl_id", crawlID).
					Execute()
				if err == nil {
					var crawlPages []map[string]interface{}
					if err := json.Unmarshal(crawlPageData, &crawlPages); err == nil {
						pages = append(pages, crawlPages...)
					}
				}
			}
		}
		// Set pageData to empty to avoid using it below
		pageData = []byte{}
	}

	if err == nil && len(pageData) > 0 {
		json.Unmarshal(pageData, &pages)
	}

	// Fetch issues for these pages
	var issues []map[string]interface{}
	pageIDs := make([]int64, 0, len(pageRankings))
	for pageID := range pageRankings {
		pageIDs = append(pageIDs, pageID)
	}
	
	if len(pageIDs) > 0 {
		// Fetch issues for each page ID individually (Supabase In() may have type constraints)
		for _, pageID := range pageIDs {
			issueData, _, err := s.serviceRole.From("issues").
				Select("*", "", false).
				Eq("page_id", strconv.FormatInt(pageID, 10)).
				Execute()
			if err == nil {
				var pageIssues []map[string]interface{}
				if err := json.Unmarshal(issueData, &pageIssues); err == nil {
					issues = append(issues, pageIssues...)
				}
			}
		}
	}

	// Group issues by page_id (page_id is bigint)
	pageIssues := make(map[int64][]map[string]interface{})
	for _, issue := range issues {
		var pageID int64
		if pageIDFloat, ok := issue["page_id"].(float64); ok {
			pageID = int64(pageIDFloat)
		}
		if pageID > 0 {
			pageIssues[pageID] = append(pageIssues[pageID], issue)
		}
	}

	// Build impact-first results
	type ImpactPage struct {
		PageID       int64                    `json:"page_id"`
		CrawlID      string                   `json:"crawl_id,omitempty"`
		URL          string                   `json:"url"`
		BestPosition int                      `json:"best_position"`
		KeywordCount int                      `json:"keyword_count"`
		Keywords     []string                 `json:"keywords"`
		IssueCount   int                      `json:"issue_count"`
		Issues       []map[string]interface{} `json:"issues"`
		ImpactScore  float64                  `json:"impact_score"` // Lower position + more issues = higher impact
	}

	impactPages := make([]ImpactPage, 0)
	for pageID, ranking := range pageRankings {
		// Find page URL and crawl_id
		var pageURL string
		var crawlID string
		for _, page := range pages {
			// pages.id is bigint, comes as float64 from JSON
			var pageIDFromDB int64
			if pageIDFloat, ok := page["id"].(float64); ok {
				pageIDFromDB = int64(pageIDFloat)
			}
			if pageIDFromDB == pageID {
				pageURL = getKeywordString(page, "url")
				crawlID = getKeywordString(page, "crawl_id")
				break
			}
		}

		// Get issues for this page
		pageIssuesList := pageIssues[pageID]
		if len(pageIssuesList) == 0 {
			continue // Skip pages without issues
		}

		// Calculate impact score: (100 - position) * issue_count
		// Higher score = more impact (lower position is better, more issues = more impact)
		impactScore := float64(100-ranking.BestPosition) * float64(len(pageIssuesList))

		impactPages = append(impactPages, ImpactPage{
			PageID:       pageID,
			CrawlID:      crawlID,
			URL:          pageURL,
			BestPosition: ranking.BestPosition,
			KeywordCount: ranking.KeywordCount,
			Keywords:     ranking.Keywords,
			IssueCount:   len(pageIssuesList),
			Issues:       pageIssuesList,
			ImpactScore:  impactScore,
		})
	}

	// Sort by impact score (descending)
	for i := 0; i < len(impactPages)-1; i++ {
		for j := i + 1; j < len(impactPages); j++ {
			if impactPages[i].ImpactScore < impactPages[j].ImpactScore {
				impactPages[i], impactPages[j] = impactPages[j], impactPages[i]
			}
		}
	}

	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"pages": impactPages,
		"count": len(impactPages),
	})
}

type PageRanking struct {
	PageID       int64
	BestPosition int
	KeywordCount int
	Keywords     []string
}

