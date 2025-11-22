package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dillonlara115/barracuda/internal/ai"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// IssueInsightRequest represents a request to generate an issue insight
type IssueInsightRequest struct {
	IssueID  string `json:"issue_id"`
	CrawlID  string `json:"crawl_id"`
}

// CrawlSummaryRequest represents a request to generate a crawl summary
type CrawlSummaryRequest struct {
	CrawlID string `json:"crawl_id"`
}

// OpenAIKeyRequest represents a request to save OpenAI API key
type OpenAIKeyRequest struct {
	OpenAIAPIKey string `json:"openai_api_key"`
}

// handleIssueInsight handles POST /api/v1/ai/issue-insight
func (s *Server) handleIssueInsight(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, ok := userIDFromContext(r.Context())
	if !ok {
		s.respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req IssueInsightRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
		return
	}

	if req.IssueID == "" || req.CrawlID == "" {
		s.respondError(w, http.StatusBadRequest, "issue_id and crawl_id are required")
		return
	}

	// Check if insight already exists (caching)
	var existingInsights []map[string]interface{}
	issueIDInt, err := strconv.ParseInt(req.IssueID, 10, 64)
	if err == nil {
		data, _, err := s.supabase.From("ai_issue_insights").
			Select("*", "", false).
			Eq("issue_id", strconv.FormatInt(issueIDInt, 10)).
			Eq("user_id", userID).
			Execute()
		if err == nil {
			json.Unmarshal(data, &existingInsights)
			if len(existingInsights) > 0 {
				// Return cached insight
				s.respondJSON(w, http.StatusOK, map[string]interface{}{
					"insight": existingInsights[0]["insight_text"],
					"cached":  true,
				})
				return
			}
		}
	}

	// Load issue data
	var issues []map[string]interface{}
	issueData, _, err := s.supabase.From("issues").
		Select("*", "", false).
		Eq("id", req.IssueID).
		Execute()
	if err != nil {
		s.logger.Error("Failed to load issue", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to load issue")
		return
	}
	if err := json.Unmarshal(issueData, &issues); err != nil || len(issues) == 0 {
		s.respondError(w, http.StatusNotFound, "Issue not found")
		return
	}
	issue := issues[0]

	// Verify user has access to the project
	projectID, ok := issue["project_id"].(string)
	if !ok {
		s.respondError(w, http.StatusInternalServerError, "Invalid issue data")
		return
	}
	hasAccess, err := s.verifyProjectAccess(userID, projectID)
	if err != nil || !hasAccess {
		s.respondError(w, http.StatusForbidden, "You don't have access to this issue")
		return
	}

	// Load page data if available
	var page map[string]interface{}
	if pageID, ok := issue["page_id"].(float64); ok && pageID > 0 {
		var pages []map[string]interface{}
		pageData, _, err := s.supabase.From("pages").
			Select("*", "", false).
			Eq("id", strconv.FormatInt(int64(pageID), 10)).
			Execute()
		if err == nil {
			json.Unmarshal(pageData, &pages)
			if len(pages) > 0 {
				page = pages[0]
			}
		}
	}

	// If no page found, create minimal page data from issue URL
	if page == nil {
		page = map[string]interface{}{
			"url": getString(issue, "url"),
		}
	}

	// Load GSC data if available (optional)
	var gscData map[string]interface{}
	// TODO: Load GSC data if GSC integration is available

	// Initialize AI client
	aiClient := ai.NewAIClient(s.supabase, s.serviceRole, s.logger)

	// Generate insight
	insight, err := aiClient.GenerateIssueInsight(r.Context(), userID, issue, page, gscData)
	if err != nil {
		s.logger.Error("Failed to generate issue insight", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate insight: %v", err))
		return
	}

	// Save to cache
	crawlID := getString(issue, "crawl_id")
	if crawlID == "" {
		crawlID = req.CrawlID
	}
	insightRecord := map[string]interface{}{
		"id":           uuid.New().String(),
		"issue_id":     issueIDInt, // Use numeric value for bigint column
		"user_id":      userID,
		"project_id":    projectID,
		"crawl_id":     crawlID,
		"insight_text": insight,
	}
	_, _, err = s.supabase.From("ai_issue_insights").Insert(insightRecord, false, "", "", "").Execute()
	if err != nil {
		s.logger.Warn("Failed to cache insight", zap.Error(err))
		// Continue anyway
	}

	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"insight": insight,
		"cached":  false,
	})
}

// handleCrawlSummary handles POST /api/v1/ai/crawl-summary
func (s *Server) handleCrawlSummary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, ok := userIDFromContext(r.Context())
	if !ok {
		s.respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req CrawlSummaryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
		return
	}

	if req.CrawlID == "" {
		s.respondError(w, http.StatusBadRequest, "crawl_id is required")
		return
	}

	// Check if summary already exists (caching)
	var existingSummaries []map[string]interface{}
	data, _, err := s.supabase.From("ai_crawl_summaries").
		Select("*", "", false).
		Eq("crawl_id", req.CrawlID).
		Eq("user_id", userID).
		Execute()
	if err == nil {
		json.Unmarshal(data, &existingSummaries)
		if len(existingSummaries) > 0 {
			// Return cached summary
			s.respondJSON(w, http.StatusOK, map[string]interface{}{
				"summary": existingSummaries[0]["summary_text"],
				"cached":   true,
			})
			return
		}
	}

	// Load crawl data
	var crawls []map[string]interface{}
	crawlData, _, err := s.supabase.From("crawls").
		Select("*", "", false).
		Eq("id", req.CrawlID).
		Execute()
	if err != nil {
		s.logger.Error("Failed to load crawl", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to load crawl")
		return
	}
	if err := json.Unmarshal(crawlData, &crawls); err != nil || len(crawls) == 0 {
		s.respondError(w, http.StatusNotFound, "Crawl not found")
		return
	}
	crawl := crawls[0]

	// Verify user has access to the project
	projectID, ok := crawl["project_id"].(string)
	if !ok {
		s.respondError(w, http.StatusInternalServerError, "Invalid crawl data")
		return
	}
	hasAccess, err := s.verifyProjectAccess(userID, projectID)
	if err != nil || !hasAccess {
		s.respondError(w, http.StatusForbidden, "You don't have access to this crawl")
		return
	}

	// Load issues for this crawl
	var issues []map[string]interface{}
	issuesData, _, err := s.supabase.From("issues").
		Select("*", "", false).
		Eq("crawl_id", req.CrawlID).
		Execute()
	if err == nil {
		json.Unmarshal(issuesData, &issues)
	}

	// Load pages for this crawl
	var pages []map[string]interface{}
	pagesData, _, err := s.supabase.From("pages").
		Select("*", "", false).
		Eq("crawl_id", req.CrawlID).
		Execute()
	if err == nil {
		json.Unmarshal(pagesData, &pages)
	}

	// Build crawl data summary
	crawlSummaryData := map[string]interface{}{
		"total_pages": getValue(crawl, "total_pages"),
		"total_issues": getValue(crawl, "total_issues"),
		"issues_by_type": make(map[string]int),
		"issues_by_severity": make(map[string]int),
		"slow_pages": []interface{}{},
		"redirect_chains": 0,
		"metadata_issues": 0,
	}

	// Count issues by type and severity
	issuesByType := make(map[string]int)
	issuesBySeverity := make(map[string]int)
	for _, issue := range issues {
		if issueType, ok := issue["type"].(string); ok {
			issuesByType[issueType]++
		}
		if severity, ok := issue["severity"].(string); ok {
			issuesBySeverity[severity]++
		}
	}
	crawlSummaryData["issues_by_type"] = issuesByType
	crawlSummaryData["issues_by_severity"] = issuesBySeverity

	// Count slow pages (>3000ms)
	var slowPages []interface{}
	for _, page := range pages {
		if rt, ok := page["response_time_ms"].(float64); ok && rt > 3000 {
			slowPages = append(slowPages, page)
		}
	}
	crawlSummaryData["slow_pages"] = slowPages

	// Count redirect chains
	redirectCount := 0
	for _, page := range pages {
		if data, ok := page["data"].(map[string]interface{}); ok {
			if redirectChain, ok := data["redirect_chain"].([]interface{}); ok && len(redirectChain) > 0 {
				redirectCount++
			}
		}
	}
	crawlSummaryData["redirect_chains"] = redirectCount

	// Count metadata issues
	metadataCount := 0
	for _, page := range pages {
		if getString(page, "title") == "" || getString(page, "meta_description") == "" {
			metadataCount++
		}
	}
	crawlSummaryData["metadata_issues"] = metadataCount

	// Initialize AI client
	aiClient := ai.NewAIClient(s.supabase, s.serviceRole, s.logger)

	// Generate summary
	summary, err := aiClient.GenerateCrawlSummary(r.Context(), userID, crawlSummaryData)
	if err != nil {
		s.logger.Error("Failed to generate crawl summary", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate summary: %v", err))
		return
	}

	// Save to cache
	summaryRecord := map[string]interface{}{
		"id":           uuid.New().String(),
		"crawl_id":     req.CrawlID,
		"user_id":      userID,
		"project_id":    projectID,
		"summary_text": summary,
	}
	_, _, err = s.supabase.From("ai_crawl_summaries").Insert(summaryRecord, false, "", "", "").Execute()
	if err != nil {
		s.logger.Warn("Failed to cache summary", zap.Error(err))
		// Continue anyway
	}

	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"summary": summary,
		"cached":   false,
	})
}

// handleOpenAIKey handles POST/GET /api/v1/integrations/openai-key
func (s *Server) handleOpenAIKey(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		s.respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	switch r.Method {
	case http.MethodPost:
		// Save OpenAI API key
		var req OpenAIKeyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}

		// Upsert user AI settings
		settings := map[string]interface{}{
			"user_id":        userID,
			"openai_api_key": req.OpenAIAPIKey,
			"updated_at":     time.Now().UTC().Format(time.RFC3339),
		}

		// Try to update first
		_, _, err := s.supabase.From("user_ai_settings").
			Update(settings, "", "").
			Eq("user_id", userID).
			Execute()

		if err != nil {
			// If update fails, try insert
			settings["created_at"] = time.Now().UTC().Format(time.RFC3339)
			_, _, err = s.supabase.From("user_ai_settings").Insert(settings, false, "", "", "").Execute()
			if err != nil {
				s.logger.Error("Failed to save OpenAI key", zap.Error(err))
				s.respondError(w, http.StatusInternalServerError, "Failed to save OpenAI key")
				return
			}
		}

		s.respondJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
		})

	case http.MethodGet:
		// Get OpenAI key status (don't return the actual key)
		var settings []map[string]interface{}
		data, _, err := s.supabase.From("user_ai_settings").
			Select("openai_api_key", "", false).
			Eq("user_id", userID).
			Execute()
		if err != nil {
			s.logger.Error("Failed to load OpenAI key status", zap.Error(err))
			s.respondJSON(w, http.StatusOK, map[string]interface{}{
				"has_key": false,
			})
			return
		}

		if err := json.Unmarshal(data, &settings); err != nil {
			s.respondJSON(w, http.StatusOK, map[string]interface{}{
				"has_key": false,
			})
			return
		}

		hasKey := false
		if len(settings) > 0 {
			if key, ok := settings[0]["openai_api_key"].(string); ok && key != "" {
				hasKey = true
			}
		}

		s.respondJSON(w, http.StatusOK, map[string]interface{}{
			"has_key": hasKey,
		})

	default:
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// Helper functions
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
		return fmt.Sprintf("%v", val)
	}
	return ""
}

func getValue(m map[string]interface{}, key string) interface{} {
	if val, ok := m[key]; ok {
		return val
	}
	return nil
}

