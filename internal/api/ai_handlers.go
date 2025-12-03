package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
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
	CrawlID      string `json:"crawl_id"`
	ForceRefresh bool   `json:"force_refresh,omitempty"` // If true, bypass cache and regenerate
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
	// Use serviceRole and filter by user_id to check cache
	var existingInsights []map[string]interface{}
	issueIDInt, err := strconv.ParseInt(req.IssueID, 10, 64)
	if err == nil {
		data, _, err := s.serviceRole.From("ai_issue_insights").
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

	// Load issue data using serviceRole (we'll verify access after)
	var issues []map[string]interface{}
	issueData, _, err := s.serviceRole.From("issues").
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
	if err != nil {
		s.logger.Error("Failed to verify project access", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to verify access")
		return
	}
	if !hasAccess {
		s.logger.Warn("User does not have access to issue", zap.String("issue_id", req.IssueID), zap.String("user_id", userID), zap.String("project_id", projectID))
		s.respondError(w, http.StatusForbidden, "You don't have access to this issue")
		return
	}

	// Load page data if available (using serviceRole since we've verified access)
	var page map[string]interface{}
	if pageID, ok := issue["page_id"].(float64); ok && pageID > 0 {
		var pages []map[string]interface{}
		pageData, _, err := s.serviceRole.From("pages").
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
	gscData := s.loadGSCDataForPage(projectID, getString(page, "url"))

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
	_, _, err = s.serviceRole.From("ai_issue_insights").Insert(insightRecord, false, "", "", "").Execute()
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

	// Check if summary already exists (caching) - skip if force_refresh is true
	if !req.ForceRefresh {
		var existingSummaries []map[string]interface{}
		data, _, err := s.serviceRole.From("ai_crawl_summaries").
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
	}

	// Load crawl data using service role (bypasses RLS)
	var crawls []map[string]interface{}
	crawlData, _, err := s.serviceRole.From("crawls").
		Select("*", "", false).
		Eq("id", req.CrawlID).
		Execute()
	if err != nil {
		s.logger.Error("Failed to load crawl", zap.String("crawl_id", req.CrawlID), zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to load crawl")
		return
	}
	if err := json.Unmarshal(crawlData, &crawls); err != nil {
		s.logger.Error("Failed to parse crawl data", zap.String("crawl_id", req.CrawlID), zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to parse crawl data")
		return
	}
	if len(crawls) == 0 {
		s.logger.Warn("Crawl not found", zap.String("crawl_id", req.CrawlID))
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
	if err != nil {
		s.logger.Error("Failed to verify project access", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to verify access")
		return
	}
	if !hasAccess {
		s.logger.Warn("User does not have access to crawl", zap.String("crawl_id", req.CrawlID), zap.String("user_id", userID), zap.String("project_id", projectID))
		s.respondError(w, http.StatusForbidden, "You don't have access to this crawl")
		return
	}

	// Load issues for this crawl (using serviceRole since we've verified access)
	var issues []map[string]interface{}
	issuesData, _, err := s.serviceRole.From("issues").
		Select("*", "", false).
		Eq("crawl_id", req.CrawlID).
		Execute()
	if err == nil {
		json.Unmarshal(issuesData, &issues)
	}

	// Load pages for this crawl (using serviceRole since we've verified access)
	var pages []map[string]interface{}
	pagesData, _, err := s.serviceRole.From("pages").
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

	// Load GSC summary data if available (optional)
	gscSummaryData := s.loadGSCSummaryData(projectID)
	if gscSummaryData != nil {
		crawlSummaryData["gsc_summary"] = gscSummaryData
	}

	// Initialize AI client
	aiClient := ai.NewAIClient(s.supabase, s.serviceRole, s.logger)

	// Generate summary
	summary, err := aiClient.GenerateCrawlSummary(r.Context(), userID, crawlSummaryData)
	if err != nil {
		s.logger.Error("Failed to generate crawl summary", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate summary: %v", err))
		return
	}

	// Save to cache (using serviceRole since we've verified access)
	// If force_refresh is true, delete old cached summaries first
	if req.ForceRefresh {
		// Delete old cached summaries - wrap in recover to prevent panic
		func() {
			defer func() {
				if r := recover(); r != nil {
					s.logger.Warn("Panic during delete of old cached summary", zap.Any("panic", r), zap.String("crawl_id", req.CrawlID), zap.String("user_id", userID))
				}
			}()
			_, _, delErr := s.serviceRole.From("ai_crawl_summaries").
				Delete("", "").
				Eq("crawl_id", req.CrawlID).
				Eq("user_id", userID).
				Execute()
			if delErr != nil {
				s.logger.Debug("Failed to delete old cached summary (non-fatal)", zap.Error(delErr))
			}
		}()
	}
	
	// Always insert new summary (table allows multiple summaries per crawl/user)
	summaryRecord := map[string]interface{}{
		"id":           uuid.New().String(),
		"crawl_id":     req.CrawlID,
		"user_id":      userID,
		"project_id":    projectID,
		"summary_text": summary,
	}
	
	// Insert new summary - if this fails, log but don't fail the request
	_, _, err = s.serviceRole.From("ai_crawl_summaries").Insert(summaryRecord, false, "", "", "").Execute()
	if err != nil {
		s.logger.Warn("Failed to cache summary", zap.Error(err))
		// Continue anyway - the summary was generated successfully
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
			s.logger.Error("Failed to decode OpenAI key request", zap.Error(err))
			s.respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}

		// Validate that key is not empty
		if req.OpenAIAPIKey == "" {
			s.logger.Warn("Attempted to save empty OpenAI key", zap.String("user_id", userID))
			s.respondError(w, http.StatusBadRequest, "OpenAI API key cannot be empty")
			return
		}

		s.logger.Info("Saving OpenAI API key", zap.String("user_id", userID), zap.Bool("has_key", req.OpenAIAPIKey != ""), zap.Int("key_length", len(req.OpenAIAPIKey)))

		// Upsert user AI settings
		settings := map[string]interface{}{
			"user_id":        userID,
			"openai_api_key": req.OpenAIAPIKey,
			"updated_at":     time.Now().UTC().Format(time.RFC3339),
		}

		// Check if record exists first
		var existing []map[string]interface{}
		selectData, _, selectErr := s.serviceRole.From("user_ai_settings").
			Select("user_id", "", false).
			Eq("user_id", userID).
			Execute()

		s.logger.Info("Checked for existing OpenAI key record", 
			zap.String("user_id", userID),
			zap.Error(selectErr),
			zap.Bool("has_data", selectData != nil && len(selectData) > 0))

		if selectErr == nil && selectData != nil {
			if err := json.Unmarshal(selectData, &existing); err == nil && len(existing) > 0 {
				// Record exists, update it
				s.logger.Info("Updating existing OpenAI key record", zap.String("user_id", userID))
				updateData, _, err := s.serviceRole.From("user_ai_settings").
					Update(settings, "", "").
					Eq("user_id", userID).
					Execute()
				if err != nil {
					s.logger.Error("Failed to update OpenAI key", zap.Error(err), zap.String("user_id", userID))
					s.respondError(w, http.StatusInternalServerError, "Failed to save OpenAI key")
					return
				}
				s.logger.Info("Successfully updated OpenAI key", 
					zap.String("user_id", userID),
					zap.Bool("has_update_data", updateData != nil && len(updateData) > 0))
			} else {
				// Record doesn't exist, insert it
				s.logger.Info("Inserting new OpenAI key record", zap.String("user_id", userID))
				settings["created_at"] = time.Now().UTC().Format(time.RFC3339)
				insertData, _, err := s.serviceRole.From("user_ai_settings").Insert(settings, false, "", "", "").Execute()
				if err != nil {
					s.logger.Error("Failed to insert OpenAI key", zap.Error(err), zap.String("user_id", userID))
					s.respondError(w, http.StatusInternalServerError, "Failed to save OpenAI key")
					return
				}
				s.logger.Info("Successfully inserted OpenAI key", 
					zap.String("user_id", userID),
					zap.Bool("has_insert_data", insertData != nil && len(insertData) > 0))
			}
		} else {
			// Select failed or no data, try insert first
			s.logger.Info("Select returned no data, attempting insert", zap.String("user_id", userID), zap.Error(selectErr))
			settings["created_at"] = time.Now().UTC().Format(time.RFC3339)
			insertData, _, err := s.serviceRole.From("user_ai_settings").Insert(settings, false, "", "", "").Execute()
			if err != nil {
				// Insert failed - might be due to existing record, try update
				s.logger.Warn("Insert failed, trying update", zap.Error(err), zap.String("user_id", userID))
				updateData, _, updateErr := s.serviceRole.From("user_ai_settings").
					Update(settings, "", "").
					Eq("user_id", userID).
					Execute()
				if updateErr != nil {
					s.logger.Error("Failed to save OpenAI key (both insert and update failed)", 
						zap.Error(err), 
						zap.Error(updateErr),
						zap.String("user_id", userID))
					s.respondError(w, http.StatusInternalServerError, "Failed to save OpenAI key")
					return
				}
				s.logger.Info("Successfully updated OpenAI key after insert failed", 
					zap.String("user_id", userID),
					zap.Bool("has_update_data", updateData != nil && len(updateData) > 0))
			} else {
				s.logger.Info("Successfully inserted OpenAI key", 
					zap.String("user_id", userID),
					zap.Bool("has_insert_data", insertData != nil && len(insertData) > 0))
			}
		}

		// Verify the save worked by querying immediately
		var verifySettings []map[string]interface{}
		verifyData, _, verifyErr := s.serviceRole.From("user_ai_settings").
			Select("openai_api_key", "", false).
			Eq("user_id", userID).
			Execute()
		
		verificationStatus := "unknown"
		savedKeyLength := 0
		
		if verifyErr == nil && verifyData != nil {
			if err := json.Unmarshal(verifyData, &verifySettings); err == nil && len(verifySettings) > 0 {
				if key, ok := verifySettings[0]["openai_api_key"].(string); ok && key != "" {
					verificationStatus = "verified_saved"
					savedKeyLength = len(key)
					s.logger.Info("OpenAI key save verified successfully", 
						zap.String("user_id", userID),
						zap.Int("key_length", len(key)))
				} else {
					verificationStatus = "verified_empty"
					s.logger.Warn("OpenAI key save completed but verification shows empty key", 
						zap.String("user_id", userID))
				}
			} else {
				verificationStatus = "record_not_found"
				s.logger.Warn("OpenAI key save completed but verification found no record", 
					zap.String("user_id", userID))
			}
		} else {
			verificationStatus = "verification_failed"
			s.logger.Warn("Failed to verify OpenAI key save", 
				zap.String("user_id", userID),
				zap.Error(verifyErr))
		}

		s.logger.Info("OpenAI key save completed successfully", zap.String("user_id", userID))
		s.respondJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"debug": map[string]interface{}{
				"verification_status": verificationStatus,
				"saved_key_length":    savedKeyLength,
				"user_id_prefix":      userID[:8], // Safe to log prefix
			},
		})

	case http.MethodGet:
		// Get OpenAI key status (don't return the actual key)
		// Use serviceRole since we're filtering by user_id (user can only access their own)
		s.logger.Info("Loading OpenAI key status", zap.String("user_id", userID))
		var settings []map[string]interface{}
		data, _, err := s.serviceRole.From("user_ai_settings").
			Select("openai_api_key", "", false).
			Eq("user_id", userID).
			Execute()
		if err != nil {
			s.logger.Error("Failed to load OpenAI key status", zap.Error(err), zap.String("user_id", userID))
			s.respondJSON(w, http.StatusOK, map[string]interface{}{
				"has_key": false,
			})
			return
		}

		if err := json.Unmarshal(data, &settings); err != nil {
			s.logger.Error("Failed to unmarshal OpenAI key status", zap.Error(err), zap.String("user_id", userID))
			s.respondJSON(w, http.StatusOK, map[string]interface{}{
				"has_key": false,
			})
			return
		}

		hasKey := false
		if len(settings) > 0 {
			if key, ok := settings[0]["openai_api_key"].(string); ok && key != "" {
				hasKey = true
				s.logger.Info("OpenAI key found", zap.String("user_id", userID), zap.Int("key_length", len(key)))
			} else {
				s.logger.Info("OpenAI key field exists but is empty", zap.String("user_id", userID))
			}
		} else {
			s.logger.Info("No OpenAI key record found", zap.String("user_id", userID))
		}

		s.logger.Info("Returning OpenAI key status", zap.String("user_id", userID), zap.Bool("has_key", hasKey))
		s.respondJSON(w, http.StatusOK, map[string]interface{}{
			"has_key": hasKey,
		})

	case http.MethodDelete:
		// Delete OpenAI API key (set to empty string)
		settings := map[string]interface{}{
			"user_id":        userID,
			"openai_api_key": "",
			"updated_at":     time.Now().UTC().Format(time.RFC3339),
		}

		// Try to update first
		_, _, err := s.serviceRole.From("user_ai_settings").
			Update(settings, "", "").
			Eq("user_id", userID).
			Execute()

		if err != nil {
			// If update fails, it might mean no record exists, which is fine
			s.logger.Debug("Failed to delete OpenAI key (may not exist)", zap.Error(err))
		}

		s.respondJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
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

// normalizeURLForGSC normalizes a URL to match GSC data format
func normalizeURLForGSC(url string) string {
	url = strings.TrimSuffix(url, "/")
	url = strings.ToLower(url)
	return url
}

// loadGSCDataForPage loads GSC performance data for a specific page URL
func (s *Server) loadGSCDataForPage(projectID, pageURL string) map[string]interface{} {
	if pageURL == "" || projectID == "" {
		return nil
	}

	normalizedURL := normalizeURLForGSC(pageURL)

	// Get the latest snapshot for this project
	var snapshots []map[string]interface{}
	snapshotData, _, err := s.serviceRole.From("gsc_performance_snapshots").
		Select("id,captured_on", "", false).
		Eq("project_id", projectID).
		Execute()
	if err != nil {
		s.logger.Debug("No GSC snapshots found", zap.String("project_id", projectID), zap.Error(err))
		return nil
	}
	if err := json.Unmarshal(snapshotData, &snapshots); err != nil || len(snapshots) == 0 {
		return nil
	}
	
	// Sort by captured_on descending and take the first one
	// Convert captured_on to time for comparison
	type snapshotWithTime struct {
		snapshot map[string]interface{}
		time     time.Time
	}
	var snapshotsWithTime []snapshotWithTime
	for _, snap := range snapshots {
		capturedOnStr, ok := snap["captured_on"].(string)
		if !ok {
			continue
		}
		capturedOn, err := time.Parse("2006-01-02", capturedOnStr)
		if err != nil {
			continue
		}
		snapshotsWithTime = append(snapshotsWithTime, snapshotWithTime{
			snapshot: snap,
			time:     capturedOn,
		})
	}
	
	if len(snapshotsWithTime) == 0 {
		return nil
	}
	
	// Sort descending by time
	sort.Slice(snapshotsWithTime, func(i, j int) bool {
		return snapshotsWithTime[i].time.After(snapshotsWithTime[j].time)
	})
	
	snapshotID, ok := snapshotsWithTime[0].snapshot["id"].(string)
	if !ok {
		// Try converting to string
		snapshotID = fmt.Sprintf("%v", snapshotsWithTime[0].snapshot["id"])
		if snapshotID == "" || snapshotID == "<nil>" {
			return nil
		}
	}

	// Query GSC performance rows for this page
	var rows []map[string]interface{}
	rowData, _, err := s.serviceRole.From("gsc_performance_rows").
		Select("*", "", false).
		Eq("snapshot_id", snapshotID).
		Eq("project_id", projectID).
		Eq("row_type", "page").
		Execute()
	if err != nil {
		s.logger.Debug("Failed to query GSC rows", zap.Error(err))
		return nil
	}
	if err := json.Unmarshal(rowData, &rows); err != nil {
		return nil
	}

	// Find matching row by normalized URL
	for _, row := range rows {
		dimensionValue, ok := row["dimension_value"].(string)
		if !ok {
			continue
		}
		if normalizeURLForGSC(dimensionValue) == normalizedURL {
			metrics, _ := row["metrics"].(map[string]interface{})
			topQueries, _ := row["top_queries"].([]interface{})

			gscData := map[string]interface{}{
				"impressions": getFloatValue(metrics, "impressions"),
				"clicks":      getFloatValue(metrics, "clicks"),
				"ctr":         getFloatValue(metrics, "ctr"),
				"position":    getFloatValue(metrics, "position"),
			}

			// Add top queries if available
			if len(topQueries) > 0 {
				// Limit to top 5 queries for prompt
				maxQueries := 5
				if len(topQueries) < maxQueries {
					maxQueries = len(topQueries)
				}
				queries := make([]string, 0, maxQueries)
				for i := 0; i < maxQueries; i++ {
					if query, ok := topQueries[i].(map[string]interface{}); ok {
						if queryText, ok := query["query"].(string); ok {
							queries = append(queries, queryText)
						}
					}
				}
				if len(queries) > 0 {
					gscData["top_queries"] = queries
				}
			}

			return gscData
		}
	}

	return nil
}

// loadGSCSummaryData loads GSC summary data for a project
func (s *Server) loadGSCSummaryData(projectID string) map[string]interface{} {
	if projectID == "" {
		return nil
	}

	// Get the latest snapshot
	var snapshots []map[string]interface{}
	snapshotData, _, err := s.serviceRole.From("gsc_performance_snapshots").
		Select("*", "", false).
		Eq("project_id", projectID).
		Execute()
	if err != nil {
		return nil
	}
	if err := json.Unmarshal(snapshotData, &snapshots); err != nil || len(snapshots) == 0 {
		return nil
	}

	// Sort by captured_on descending and take the first one
	type snapshotWithTime struct {
		snapshot map[string]interface{}
		time     time.Time
	}
	var snapshotsWithTime []snapshotWithTime
	for _, snap := range snapshots {
		capturedOnStr, ok := snap["captured_on"].(string)
		if !ok {
			continue
		}
		capturedOn, err := time.Parse("2006-01-02", capturedOnStr)
		if err != nil {
			continue
		}
		snapshotsWithTime = append(snapshotsWithTime, snapshotWithTime{
			snapshot: snap,
			time:     capturedOn,
		})
	}
	
	if len(snapshotsWithTime) == 0 {
		return nil
	}
	
	// Sort descending by time
	sort.Slice(snapshotsWithTime, func(i, j int) bool {
		return snapshotsWithTime[i].time.After(snapshotsWithTime[j].time)
	})

	snapshot := snapshotsWithTime[0].snapshot
	totals, _ := snapshot["totals"].(map[string]interface{})

	if totals == nil || len(totals) == 0 {
		return nil
	}

	return map[string]interface{}{
		"total_impressions": getFloatValue(totals, "impressions"),
		"total_clicks":      getFloatValue(totals, "clicks"),
		"average_ctr":       getFloatValue(totals, "ctr"),
		"average_position": getFloatValue(totals, "position"),
		"captured_on":      snapshot["captured_on"],
	}
}

// getFloatValue safely extracts a float value from a map
func getFloatValue(m map[string]interface{}, key string) float64 {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case float64:
			return v
		case int:
			return float64(v)
		case int64:
			return float64(v)
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return f
			}
		}
	}
	return 0
}

