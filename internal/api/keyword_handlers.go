package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dillonlara115/barracudaseo/internal/dataforseo"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateKeywordRequest represents a request to create a keyword
type CreateKeywordRequest struct {
	ProjectID      string   `json:"project_id"`
	Keyword        string   `json:"keyword"`
	TargetURL      string   `json:"target_url,omitempty"`
	LocationName   string   `json:"location_name"`
	LocationCode   *int     `json:"location_code,omitempty"`
	LanguageName   string   `json:"language_name"`
	Device         string   `json:"device"`
	SearchEngine   string   `json:"search_engine"`
	Tags           []string `json:"tags,omitempty"`
	CheckFrequency string   `json:"check_frequency,omitempty"` // manual | daily | weekly
}

// KeywordResponse represents a keyword in API responses
type KeywordResponse struct {
	ID             string     `json:"id"`
	ProjectID      string     `json:"project_id"`
	Keyword        string     `json:"keyword"`
	TargetURL      *string    `json:"target_url,omitempty"`
	LocationName   string     `json:"location_name"`
	LocationCode   *int       `json:"location_code,omitempty"`
	LanguageName   string     `json:"language_name"`
	Device         string     `json:"device"`
	SearchEngine   string     `json:"search_engine"`
	Tags           []string   `json:"tags"`
	CheckFrequency string     `json:"check_frequency,omitempty"` // manual | daily | weekly
	LastCheckedAt  *time.Time `json:"last_checked_at,omitempty"`
	NextCheckAt    *time.Time `json:"next_check_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	// Enriched fields
	LatestPosition *int       `json:"latest_position,omitempty"`
	BestPosition   *int       `json:"best_position,omitempty"`
	LastChecked    *time.Time `json:"last_checked,omitempty"`
	Trend          string     `json:"trend,omitempty"` // "up", "down", "same"
}

// RankSnapshotResponse represents a rank snapshot
type RankSnapshotResponse struct {
	ID               string    `json:"id"`
	KeywordID        string    `json:"keyword_id"`
	CheckedAt        time.Time `json:"checked_at"`
	PositionAbsolute *int      `json:"position_absolute,omitempty"`
	PositionOrganic  *int      `json:"position_organic,omitempty"`
	SERPURL          *string   `json:"serp_url,omitempty"`
	SERPTitle        *string   `json:"serp_title,omitempty"`
	SERPSnippet      *string   `json:"serp_snippet,omitempty"`
	SERPFeatures     []string  `json:"serp_features,omitempty"`
	RankType         string    `json:"rank_type"`
}

// handleKeywords handles keyword collection endpoints
func (s *Server) handleKeywords(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		s.respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	switch r.Method {
	case http.MethodPost:
		s.handleCreateKeyword(w, r, userID)
	case http.MethodGet:
		s.handleListKeywords(w, r, userID)
	default:
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// handleKeywordByID handles keyword-specific endpoints
func (s *Server) handleKeywordByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		s.respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Extract keyword ID from path
	path := strings.TrimPrefix(r.URL.Path, "/keywords/")
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 || parts[0] == "" {
		s.respondError(w, http.StatusBadRequest, "Keyword ID is required")
		return
	}

	keywordID := parts[0]

	// Check if there's a sub-resource (e.g., /keywords/:id/check or /keywords/:id/snapshots)
	if len(parts) > 1 {
		subResource := parts[1]
		switch subResource {
		case "check":
			if r.Method == http.MethodPost {
				s.handleCheckKeyword(w, r, userID, keywordID)
				return
			}
		case "snapshots":
			if r.Method == http.MethodGet {
				s.handleGetKeywordSnapshots(w, r, userID, keywordID)
				return
			}
		}
	}

	// Handle standard CRUD operations
	switch r.Method {
	case http.MethodGet:
		s.handleGetKeyword(w, r, userID, keywordID)
	case http.MethodPut, http.MethodPatch:
		s.handleUpdateKeyword(w, r, userID, keywordID)
	case http.MethodDelete:
		s.handleDeleteKeyword(w, r, userID, keywordID)
	default:
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// handleCreateKeyword handles POST /api/v1/keywords
func (s *Server) handleCreateKeyword(w http.ResponseWriter, r *http.Request, userID string) {
	var req CreateKeywordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
		return
	}

	// Validate required fields
	if req.ProjectID == "" {
		s.respondError(w, http.StatusBadRequest, "project_id is required")
		return
	}
	if req.Keyword == "" {
		s.respondError(w, http.StatusBadRequest, "keyword is required")
		return
	}
	if req.LocationName == "" {
		s.respondError(w, http.StatusBadRequest, "location_name is required")
		return
	}
	if req.Device == "" {
		req.Device = "desktop"
	}
	if req.LanguageName == "" {
		req.LanguageName = "English"
	}
	if req.SearchEngine == "" {
		req.SearchEngine = "google.com"
	}
	if req.CheckFrequency == "" {
		req.CheckFrequency = "manual"
	}
	if req.CheckFrequency != "manual" && req.CheckFrequency != "daily" && req.CheckFrequency != "weekly" {
		s.respondError(w, http.StatusBadRequest, "check_frequency must be 'manual', 'daily', or 'weekly'")
		return
	}

	// Verify project access
	hasAccess, err := s.verifyProjectAccess(userID, req.ProjectID)
	if err != nil {
		s.logger.Error("Failed to verify project access", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to verify project access")
		return
	}
	if !hasAccess {
		s.respondError(w, http.StatusForbidden, "You don't have access to this project")
		return
	}

	// Check subscription tier limits
	subscription, err := s.resolveSubscription(userID)
	if err != nil {
		s.logger.Error("Failed to resolve subscription", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to verify subscription")
		return
	}

	// Get current keyword count for project
	var keywords []map[string]interface{}
	keywordData, _, err := s.serviceRole.From("keywords").Select("id", "", false).Eq("project_id", req.ProjectID).Execute()
	if err == nil {
		json.Unmarshal(keywordData, &keywords)
	}
	currentKeywordCount := len(keywords)

	maxKeywords := getKeywordLimit(subscription.EffectiveTier)

	if currentKeywordCount >= maxKeywords {
		s.respondError(w, http.StatusForbidden, fmt.Sprintf("Your %s plan allows a maximum of %d keywords. Please upgrade to add more keywords.", subscription.EffectiveTier, maxKeywords))
		return
	}

	// Create keyword record
	keywordID := uuid.New().String()
	keywordDataMap := map[string]interface{}{
		"id":              keywordID,
		"project_id":      req.ProjectID,
		"keyword":         req.Keyword,
		"location_name":   req.LocationName,
		"language_name":   req.LanguageName,
		"device":          req.Device,
		"search_engine":   req.SearchEngine,
		"tags":            req.Tags,
		"check_frequency": req.CheckFrequency,
	}

	if req.TargetURL != "" {
		keywordDataMap["target_url"] = req.TargetURL
	}
	if req.LocationCode != nil {
		keywordDataMap["location_code"] = *req.LocationCode
	}

	_, _, err = s.serviceRole.From("keywords").Insert(keywordDataMap, false, "", "", "").Execute()
	if err != nil {
		// Check for duplicate key violation (PostgreSQL error code 23505)
		if strings.Contains(err.Error(), "23505") || strings.Contains(err.Error(), "duplicate key") {
			s.respondError(w, http.StatusConflict, fmt.Sprintf(
				"A keyword '%s' already exists for this project with location '%s' and device '%s'. Please use a different combination or update the existing keyword.",
				req.Keyword, req.LocationName, req.Device,
			))
			return
		}
		s.logger.Error("Failed to insert keyword", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to create keyword")
		return
	}

	// Fetch created keyword
	keyword, err := s.fetchKeyword(keywordID)
	if err != nil {
		s.logger.Error("Failed to fetch created keyword", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to fetch created keyword")
		return
	}

	s.respondJSON(w, http.StatusCreated, keyword)
}

// handleListKeywords handles GET /api/v1/keywords?project_id=...
func (s *Server) handleListKeywords(w http.ResponseWriter, r *http.Request, userID string) {
	projectID := r.URL.Query().Get("project_id")
	if projectID == "" {
		s.respondError(w, http.StatusBadRequest, "project_id query parameter is required")
		return
	}

	// Verify project access
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

	// Build query using serviceRole since we've already verified project access
	query := s.serviceRole.From("keywords").Select("*", "", false).Eq("project_id", projectID)

	// Apply filters
	if device := r.URL.Query().Get("device"); device != "" {
		query = query.Eq("device", device)
	}
	if location := r.URL.Query().Get("location"); location != "" {
		query = query.Eq("location_name", location)
	}
	if tag := r.URL.Query().Get("tag"); tag != "" {
		query = query.Contains("tags", []string{tag})
	}

	// Order by created_at desc
	query = query.Order("created_at", nil)

	data, _, err := query.Execute()
	if err != nil {
		// Check if error is because table doesn't exist yet (migration not run)
		if strings.Contains(err.Error(), "Could not find the table") ||
			strings.Contains(err.Error(), "does not exist") ||
			strings.Contains(err.Error(), "PGRST205") ||
			strings.Contains(err.Error(), "relation") {
			// Table doesn't exist yet - return empty array (graceful degradation)
			s.logger.Info("keywords table not found, returning empty list",
				zap.String("project_id", projectID),
				zap.String("user_id", userID))
			s.respondJSON(w, http.StatusOK, map[string]interface{}{
				"keywords": []KeywordResponse{},
				"count":    0,
			})
			return
		}
		s.logger.Error("Failed to list keywords", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to list keywords")
		return
	}

	var keywords []map[string]interface{}
	if err := json.Unmarshal(data, &keywords); err != nil {
		s.logger.Error("Failed to parse keywords data", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to parse keywords")
		return
	}

	// Handle empty result gracefully
	if keywords == nil {
		keywords = []map[string]interface{}{}
	}

	// Enrich keywords with latest snapshot data
	enrichedKeywords := make([]KeywordResponse, 0, len(keywords))
	for _, k := range keywords {
		keyword := s.mapToKeywordResponse(k)

		// Fetch latest snapshot
		snapshot, err := s.fetchLatestSnapshot(keyword.ID)
		if err == nil && snapshot != nil {
			if snapshot.PositionOrganic != nil {
				keyword.LatestPosition = snapshot.PositionOrganic
			}
			keyword.LastChecked = &snapshot.CheckedAt

			// Calculate best position and trend
			bestPos, trend := s.calculateBestPositionAndTrend(keyword.ID, snapshot.PositionOrganic)
			keyword.BestPosition = bestPos
			keyword.Trend = trend
		}

		enrichedKeywords = append(enrichedKeywords, keyword)
	}

	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"keywords": enrichedKeywords,
		"count":    len(enrichedKeywords),
	})
}

// handleGetKeyword handles GET /api/v1/keywords/:id
func (s *Server) handleGetKeyword(w http.ResponseWriter, r *http.Request, userID string, keywordID string) {
	_ = r

	keyword, err := s.fetchKeyword(keywordID)
	if err != nil {
		s.respondError(w, http.StatusNotFound, "Keyword not found")
		return
	}

	// Verify project access
	hasAccess, err := s.verifyProjectAccess(userID, keyword.ProjectID)
	if err != nil {
		s.logger.Error("Failed to verify project access", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to verify project access")
		return
	}
	if !hasAccess {
		s.respondError(w, http.StatusForbidden, "You don't have access to this keyword")
		return
	}

	// Enrich with snapshot data
	snapshot, err := s.fetchLatestSnapshot(keywordID)
	if err == nil && snapshot != nil {
		if snapshot.PositionOrganic != nil {
			keyword.LatestPosition = snapshot.PositionOrganic
		}
		keyword.LastChecked = &snapshot.CheckedAt
		bestPos, trend := s.calculateBestPositionAndTrend(keywordID, snapshot.PositionOrganic)
		keyword.BestPosition = bestPos
		keyword.Trend = trend
	}

	s.respondJSON(w, http.StatusOK, keyword)
}

// handleUpdateKeyword handles PUT/PATCH /api/v1/keywords/:id
func (s *Server) handleUpdateKeyword(w http.ResponseWriter, r *http.Request, userID string, keywordID string) {
	keyword, err := s.fetchKeyword(keywordID)
	if err != nil {
		s.respondError(w, http.StatusNotFound, "Keyword not found")
		return
	}

	// Verify project access
	hasAccess, err := s.verifyProjectAccess(userID, keyword.ProjectID)
	if err != nil {
		s.logger.Error("Failed to verify project access", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to verify project access")
		return
	}
	if !hasAccess {
		s.respondError(w, http.StatusForbidden, "You don't have access to this keyword")
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		s.respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
		return
	}

	// Update keyword
	_, _, err = s.serviceRole.From("keywords").Update(updates, "", "").Eq("id", keywordID).Execute()
	if err != nil {
		s.logger.Error("Failed to update keyword", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to update keyword")
		return
	}

	// Fetch updated keyword
	updatedKeyword, err := s.fetchKeyword(keywordID)
	if err != nil {
		s.logger.Error("Failed to fetch updated keyword", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to fetch updated keyword")
		return
	}

	s.respondJSON(w, http.StatusOK, updatedKeyword)
}

// handleDeleteKeyword handles DELETE /api/v1/keywords/:id
func (s *Server) handleDeleteKeyword(w http.ResponseWriter, r *http.Request, userID string, keywordID string) {
	_ = r

	keyword, err := s.fetchKeyword(keywordID)
	if err != nil {
		s.respondError(w, http.StatusNotFound, "Keyword not found")
		return
	}

	// Verify project access
	hasAccess, err := s.verifyProjectAccess(userID, keyword.ProjectID)
	if err != nil {
		s.logger.Error("Failed to verify project access", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to verify project access")
		return
	}
	if !hasAccess {
		s.respondError(w, http.StatusForbidden, "You don't have access to this keyword")
		return
	}

	// Delete keyword (cascade will delete snapshots and tasks)
	_, _, err = s.serviceRole.From("keywords").Delete("", "").Eq("id", keywordID).Execute()
	if err != nil {
		s.logger.Error("Failed to delete keyword", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to delete keyword")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleCheckKeyword handles POST /api/v1/keywords/:id/check
func (s *Server) handleCheckKeyword(w http.ResponseWriter, r *http.Request, userID string, keywordID string) {
	keyword, err := s.fetchKeyword(keywordID)
	if err != nil {
		s.respondError(w, http.StatusNotFound, "Keyword not found")
		return
	}

	// Verify project access
	hasAccess, err := s.verifyProjectAccess(userID, keyword.ProjectID)
	if err != nil {
		s.logger.Error("Failed to verify project access", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to verify project access")
		return
	}
	if !hasAccess {
		s.respondError(w, http.StatusForbidden, "You don't have access to this keyword")
		return
	}

	// Get DataForSEO client
	client, err := dataforseo.NewClient()
	if err != nil {
		s.logger.Error("Failed to create DataForSEO client", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "DataForSEO integration not configured")
		return
	}

	// Use Live API for immediate "check now" requests
	// Live API returns results immediately in a single request (no polling needed)
	// More expensive but instant - perfect for manual checks
	task := dataforseo.OrganicTaskPost{
		LanguageName: keyword.LanguageName,
		LocationName: keyword.LocationName,
		Keyword:      keyword.Keyword,
		Device:       keyword.Device,
		SearchEngine: keyword.SearchEngine,
	}

	s.logger.Info("Using Live API for immediate rank check",
		zap.String("keyword_id", keywordID),
		zap.String("keyword", keyword.Keyword))

	// Call Live API - returns results immediately
	liveResp, err := client.CreateOrganicTaskLive(r.Context(), task)
	if err != nil {
		s.logger.Error("Failed to create DataForSEO Live task", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create rank check task: %v", err))
		return
	}

	if len(liveResp.Tasks) == 0 {
		s.logger.Error("No tasks returned from DataForSEO Live API", zap.Any("response", liveResp))
		s.respondError(w, http.StatusInternalServerError, "No results returned from DataForSEO")
		return
	}

	taskResult := liveResp.Tasks[0]
	taskID := taskResult.ID
	taskStatusCode := taskResult.StatusCode

	s.logger.Info("Received Live API response",
		zap.String("task_id", taskID),
		zap.String("keyword_id", keywordID),
		zap.String("keyword", keyword.Keyword),
		zap.Int("status_code", taskStatusCode),
		zap.String("status_message", taskResult.StatusMessage))

	// Check if task was successful
	if taskStatusCode != 20000 {
		s.logger.Warn("DataForSEO Live API returned non-success status",
			zap.Int("status_code", taskStatusCode),
			zap.String("status_message", taskResult.StatusMessage))
		s.respondError(w, http.StatusBadRequest, fmt.Sprintf("DataForSEO task failed: %s (code: %d)",
			taskResult.StatusMessage, taskStatusCode))
		return
	}

	// Create task record for tracking
	taskRecordID := uuid.New().String()
	now := time.Now().UTC()
	projectID := keyword.ProjectID
	if projectID == "" {
		// Defensive fallback: fetch keyword to populate project_id
		if kw, err := s.fetchKeyword(keywordID); err == nil && kw.ProjectID != "" {
			projectID = kw.ProjectID
		}
	}
	if projectID == "" {
		s.logger.Warn("Missing project_id for keyword task record", zap.String("keyword_id", keywordID))
	}

	taskRecord := map[string]interface{}{
		"id":                 taskRecordID,
		"project_id":         projectID,
		"keyword_id":         keywordID,
		"dataforseo_task_id": taskID,
		"status":             "completed", // Live API completes immediately
		"run_at":             now.Format(time.RFC3339),
		"completed_at":       now.Format(time.RFC3339),
		"raw_request":        map[string]interface{}{"task": task},
		"raw_response":       liveResp,
	}

	_, _, err = s.serviceRole.From("keyword_tasks").Insert(taskRecord, false, "", "", "").Execute()
	if err != nil {
		s.logger.Warn("Failed to insert task record", zap.Error(err))
		// Continue anyway - we'll still create the snapshot
	}

	// Extract ranking from Live API response
	targetURL := ""
	if keyword.TargetURL != nil {
		targetURL = *keyword.TargetURL
	}
	ranking, err := dataforseo.ExtractRanking(liveResp, targetURL)
	if err != nil {
		// Check if error indicates site is not ranking
		if strings.Contains(err.Error(), "is not ranking") {
			s.logger.Info("Target URL is not ranking in search results",
				zap.String("keyword_id", keywordID),
				zap.String("target_url", targetURL),
				zap.String("keyword", keyword.Keyword))
			s.respondError(w, http.StatusOK, fmt.Sprintf("Your site (%s) is not currently ranking for the keyword '%s' in the search results. This is expected if your site is new or hasn't gained visibility for this keyword yet.", targetURL, keyword.Keyword))
			return
		}
		s.logger.Error("Failed to extract ranking from Live API response", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to extract ranking: %v", err))
		return
	}

	// Create snapshot
	snapshot, err := s.createSnapshot(keyword.ProjectID, keywordID, taskID, ranking)
	if err != nil {
		s.logger.Error("Failed to create snapshot", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create snapshot: %v", err))
		return
	}

	// Track usage
	checkType := "manual"
	if keyword.CheckFrequency != "manual" {
		checkType = "scheduled"
	}
	if err := s.trackKeywordUsage(r.Context(), keyword.ProjectID, keywordID, userID, taskID, checkType, DefaultCheckCost); err != nil {
		s.logger.Warn("Failed to track keyword usage", zap.Error(err))
	}

	// Update keyword's last_checked_at
	_, _, _ = s.serviceRole.From("keywords").
		Update(map[string]interface{}{"last_checked_at": now.Format(time.RFC3339)}, "", "").
		Eq("id", keywordID).
		Execute()

	// Return snapshot immediately - no polling needed!
	s.respondJSON(w, http.StatusOK, snapshot)
}

// handleGetKeywordSnapshots handles GET /api/v1/keywords/:id/snapshots
func (s *Server) handleGetKeywordSnapshots(w http.ResponseWriter, r *http.Request, userID string, keywordID string) {
	keyword, err := s.fetchKeyword(keywordID)
	if err != nil {
		s.respondError(w, http.StatusNotFound, "Keyword not found")
		return
	}

	// Verify project access
	hasAccess, err := s.verifyProjectAccess(userID, keyword.ProjectID)
	if err != nil {
		s.logger.Error("Failed to verify project access", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to verify project access")
		return
	}
	if !hasAccess {
		s.respondError(w, http.StatusForbidden, "You don't have access to this keyword")
		return
	}

	// Get limit from query params
	limit := 30
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Fetch snapshots
	query := s.supabase.From("keyword_rank_snapshots").
		Select("*", "", false).
		Eq("keyword_id", keywordID).
		Order("checked_at", nil).
		Limit(limit, "")

	data, _, err := query.Execute()
	if err != nil {
		s.logger.Error("Failed to fetch snapshots", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to fetch snapshots")
		return
	}

	var snapshots []map[string]interface{}
	if err := json.Unmarshal(data, &snapshots); err != nil {
		s.logger.Error("Failed to parse snapshots data", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to parse snapshots")
		return
	}

	// Convert to response format
	responseSnapshots := make([]RankSnapshotResponse, 0, len(snapshots))
	for _, snap := range snapshots {
		snapshot := s.mapToSnapshotResponse(snap)
		responseSnapshots = append(responseSnapshots, snapshot)
	}

	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"snapshots": responseSnapshots,
		"count":     len(responseSnapshots),
	})
}

// Helper functions

func (s *Server) fetchKeyword(keywordID string) (*KeywordResponse, error) {
	data, _, err := s.serviceRole.From("keywords").Select("*", "", false).Eq("id", keywordID).Execute()
	if err != nil {
		return nil, err
	}

	var keywords []map[string]interface{}
	if err := json.Unmarshal(data, &keywords); err != nil || len(keywords) == 0 {
		return nil, fmt.Errorf("keyword not found")
	}

	keyword := s.mapToKeywordResponse(keywords[0])
	return &keyword, nil
}

func (s *Server) mapToKeywordResponse(m map[string]interface{}) KeywordResponse {
	keyword := KeywordResponse{
		ID:             getKeywordString(m, "id"),
		ProjectID:      getKeywordString(m, "project_id"),
		Keyword:        getKeywordString(m, "keyword"),
		LocationName:   getKeywordString(m, "location_name"),
		LanguageName:   getKeywordString(m, "language_name"),
		Device:         getKeywordString(m, "device"),
		SearchEngine:   getKeywordString(m, "search_engine"),
		CheckFrequency: getKeywordString(m, "check_frequency"),
	}

	if targetURL, ok := m["target_url"].(string); ok && targetURL != "" {
		keyword.TargetURL = &targetURL
	}
	if locationCode, ok := m["location_code"].(float64); ok {
		code := int(locationCode)
		keyword.LocationCode = &code
	}
	if lastCheckedAt, ok := m["last_checked_at"].(string); ok && lastCheckedAt != "" {
		if t, err := time.Parse(time.RFC3339, lastCheckedAt); err == nil {
			keyword.LastCheckedAt = &t
			keyword.LastChecked = &t
		}
	}
	if nextCheckAt, ok := m["next_check_at"].(string); ok && nextCheckAt != "" {
		if t, err := time.Parse(time.RFC3339, nextCheckAt); err == nil {
			keyword.NextCheckAt = &t
		}
	}
	if tags, ok := m["tags"].([]interface{}); ok {
		keyword.Tags = make([]string, 0, len(tags))
		for _, tag := range tags {
			if tagStr, ok := tag.(string); ok {
				keyword.Tags = append(keyword.Tags, tagStr)
			}
		}
	}
	if createdAt, ok := m["created_at"].(string); ok {
		if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
			keyword.CreatedAt = t
		}
	}
	if updatedAt, ok := m["updated_at"].(string); ok {
		if t, err := time.Parse(time.RFC3339, updatedAt); err == nil {
			keyword.UpdatedAt = t
		}
	}

	return keyword
}

// handleKeywordUsage handles GET /api/v1/projects/:id/keyword-usage
func (s *Server) handleKeywordUsage(w http.ResponseWriter, r *http.Request, projectID, userID string) {
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

	// Get usage stats
	stats, err := s.getKeywordUsageStats(r.Context(), projectID, "", nil, nil)
	if err != nil {
		s.logger.Error("Failed to get keyword usage stats", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to get usage stats")
		return
	}

	// Get keyword limits
	canAdd, currentCount, maxKeywords, err := s.checkKeywordLimit(r.Context(), projectID, userID)
	if err != nil {
		s.logger.Error("Failed to check keyword limit", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to check keyword limit")
		return
	}

	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"usage": map[string]interface{}{
			"total_checks":     stats["total_checks"],
			"total_cost_usd":   stats["total_cost_usd"],
			"manual_checks":    stats["manual_checks"],
			"scheduled_checks": stats["scheduled_checks"],
		},
		"limits": map[string]interface{}{
			"current_keywords": currentCount,
			"max_keywords":     maxKeywords,
			"can_add_more":     canAdd,
		},
	})
}

func (s *Server) fetchLatestSnapshot(keywordID string) (*RankSnapshotResponse, error) {
	data, _, err := s.serviceRole.From("keyword_rank_snapshots").
		Select("*", "", false).
		Eq("keyword_id", keywordID).
		Order("checked_at", nil).
		Limit(1, "").
		Execute()
	if err != nil {
		return nil, err
	}

	var snapshots []map[string]interface{}
	if err := json.Unmarshal(data, &snapshots); err != nil || len(snapshots) == 0 {
		return nil, fmt.Errorf("no snapshots found")
	}

	snapshot := s.mapToSnapshotResponse(snapshots[0])
	return &snapshot, nil
}

func (s *Server) mapToSnapshotResponse(m map[string]interface{}) RankSnapshotResponse {
	snapshot := RankSnapshotResponse{
		ID:        getKeywordString(m, "id"),
		KeywordID: getKeywordString(m, "keyword_id"),
		RankType:  getKeywordString(m, "rank_type"),
	}

	if checkedAt, ok := m["checked_at"].(string); ok {
		if t, err := time.Parse(time.RFC3339, checkedAt); err == nil {
			snapshot.CheckedAt = t
		}
	}
	if posAbs, ok := m["position_absolute"].(float64); ok {
		p := int(posAbs)
		snapshot.PositionAbsolute = &p
	}
	if posOrg, ok := m["position_organic"].(float64); ok {
		p := int(posOrg)
		snapshot.PositionOrganic = &p
	}
	if url, ok := m["serp_url"].(string); ok && url != "" {
		snapshot.SERPURL = &url
	}
	if title, ok := m["serp_title"].(string); ok && title != "" {
		snapshot.SERPTitle = &title
	}
	if snippet, ok := m["serp_snippet"].(string); ok && snippet != "" {
		snapshot.SERPSnippet = &snippet
	}
	if features, ok := m["serp_features"].([]interface{}); ok {
		snapshot.SERPFeatures = make([]string, 0, len(features))
		for _, f := range features {
			if fStr, ok := f.(string); ok {
				snapshot.SERPFeatures = append(snapshot.SERPFeatures, fStr)
			}
		}
	}

	return snapshot
}

func (s *Server) createSnapshot(projectID, keywordID, taskID string, ranking *dataforseo.RankingData) (*RankSnapshotResponse, error) {
	// Backfill project_id if missing
	if projectID == "" {
		if kw, err := s.fetchKeyword(keywordID); err == nil && kw.ProjectID != "" {
			projectID = kw.ProjectID
		}
	}
	if projectID == "" {
		return nil, fmt.Errorf("project_id missing for keyword %s", keywordID)
	}
	// Try to find matching crawl page by URL
	var crawlPageID *int64
	if ranking.URL != "" {
		// Normalize URL for matching (remove trailing slashes, etc.)
		normalizedURL := strings.TrimSuffix(ranking.URL, "/")

		// Find pages matching this URL (get the most recent one)
		var pages []map[string]interface{}
		pageData, _, err := s.serviceRole.From("pages").
			Select("id", "", false).
			Eq("url", normalizedURL).
			Order("created_at", nil).
			Limit(1, "").
			Execute()
		if err == nil {
			if err := json.Unmarshal(pageData, &pages); err == nil && len(pages) > 0 {
				// pages.id is bigint (bigserial), so it comes as float64 from JSON
				if pageIDFloat, ok := pages[0]["id"].(float64); ok {
					pageID := int64(pageIDFloat)
					crawlPageID = &pageID
				}
			}
		}
	}

	snapshotID := uuid.New().String()
	snapshotData := map[string]interface{}{
		"id":                 snapshotID,
		"project_id":         projectID,
		"keyword_id":         keywordID,
		"dataforseo_task_id": taskID,
		"position_absolute":  ranking.PositionAbsolute,
		"position_organic":   ranking.PositionOrganic,
		"serp_url":           ranking.URL,
		"serp_title":         ranking.Title,
		"serp_snippet":       ranking.Snippet,
		"serp_features":      ranking.SERPFeatures,
		"rank_type":          "organic",
	}

	if crawlPageID != nil {
		snapshotData["crawl_page_id"] = *crawlPageID
	}

	_, _, err := s.serviceRole.From("keyword_rank_snapshots").Insert(snapshotData, false, "", "", "").Execute()
	if err != nil {
		return nil, err
	}

	// Fetch created snapshot
	data, _, err := s.serviceRole.From("keyword_rank_snapshots").Select("*", "", false).Eq("id", snapshotID).Execute()
	if err != nil {
		return nil, err
	}

	var snapshots []map[string]interface{}
	if err := json.Unmarshal(data, &snapshots); err != nil || len(snapshots) == 0 {
		return nil, fmt.Errorf("failed to fetch created snapshot")
	}

	snapshot := s.mapToSnapshotResponse(snapshots[0])
	return &snapshot, nil
}

func (s *Server) calculateBestPositionAndTrend(keywordID string, currentPos *int) (*int, string) {
	if currentPos == nil {
		return nil, ""
	}

	// Fetch all snapshots to calculate best position
	data, _, err := s.serviceRole.From("keyword_rank_snapshots").
		Select("position_organic", "", false).
		Eq("keyword_id", keywordID).
		Order("checked_at", nil).
		Limit(10, "").
		Execute()
	if err != nil {
		return currentPos, ""
	}

	var snapshots []map[string]interface{}
	if err := json.Unmarshal(data, &snapshots); err != nil {
		return currentPos, ""
	}

	bestPos := *currentPos
	for _, s := range snapshots {
		if pos, ok := s["position_organic"].(float64); ok {
			posInt := int(pos)
			if posInt > 0 && (bestPos == 0 || posInt < bestPos) {
				bestPos = posInt
			}
		}
	}

	// Calculate trend (compare with previous snapshot)
	trend := ""
	if len(snapshots) > 1 {
		if prevPos, ok := snapshots[1]["position_organic"].(float64); ok {
			prevPosInt := int(prevPos)
			if *currentPos < prevPosInt {
				trend = "up"
			} else if *currentPos > prevPosInt {
				trend = "down"
			} else {
				trend = "same"
			}
		}
	}

	return &bestPos, trend
}

// handleProjectKeywordMetrics handles GET /api/v1/projects/:id/keyword-metrics
func (s *Server) handleProjectKeywordMetrics(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	_ = r

	// Verify project access
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

	// Fetch all keywords for project
	keywordsData, _, err := s.serviceRole.From("keywords").
		Select("id", "", false).
		Eq("project_id", projectID).
		Execute()
	if err != nil {
		s.logger.Error("Failed to fetch keywords", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to fetch keywords")
		return
	}

	var keywords []map[string]interface{}
	if err := json.Unmarshal(keywordsData, &keywords); err != nil {
		s.logger.Error("Failed to parse keywords", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to parse keywords")
		return
	}

	// Aggregate metrics
	totalKeywords := len(keywords)
	trackedKeywords := 0
	avgPosition := 0.0
	improvedCount := 0
	declinedCount := 0
	noChangeCount := 0

	keywordIDs := make([]string, 0, len(keywords))
	for _, k := range keywords {
		if id, ok := k["id"].(string); ok {
			keywordIDs = append(keywordIDs, id)
		}
	}

	if len(keywordIDs) > 0 {
		// Fetch latest snapshots for all keywords
		// Note: This is simplified - in production, you might want a more efficient query
		for _, keywordID := range keywordIDs {
			snapshot, err := s.fetchLatestSnapshot(keywordID)
			if err == nil && snapshot != nil {
				trackedKeywords++
				if snapshot.PositionOrganic != nil {
					avgPosition += float64(*snapshot.PositionOrganic)
					_, trend := s.calculateBestPositionAndTrend(keywordID, snapshot.PositionOrganic)
					switch trend {
					case "up":
						improvedCount++
					case "down":
						declinedCount++
					case "same":
						noChangeCount++
					}
				}
			}
		}

		if trackedKeywords > 0 {
			avgPosition = avgPosition / float64(trackedKeywords)
		}
	}

	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"total_keywords":   totalKeywords,
		"tracked_keywords": trackedKeywords,
		"average_position": avgPosition,
		"improved_count":   improvedCount,
		"declined_count":   declinedCount,
		"no_change_count":  noChangeCount,
	})
}

// DiscoverKeywordsRequest represents a request to discover keywords for a domain/URL
type DiscoverKeywordsRequest struct {
	Target       string `json:"target"`                 // Domain (e.g., "example.com") or URL
	LocationName string `json:"location_name"`          // e.g., "United States"
	LanguageName string `json:"language_name"`          // e.g., "English"
	Limit        int    `json:"limit,omitempty"`        // Max results (default 1000)
	MinPosition  int    `json:"min_position,omitempty"` // Minimum position to include
	MaxPosition  int    `json:"max_position,omitempty"` // Maximum position to include
}

// DiscoveredKeywordResponse represents a discovered keyword
type DiscoveredKeywordResponse struct {
	Keyword           string  `json:"keyword"`
	Position          int     `json:"position"`
	URL               string  `json:"url"`
	Title             string  `json:"title"`
	SearchVolume      int     `json:"search_volume"`
	Competition       string  `json:"competition"`
	CPC               float64 `json:"cpc"`
	KeywordDifficulty int     `json:"keyword_difficulty"`
	MatchedPageID     *int64  `json:"matched_page_id,omitempty"`  // ID of crawled page if matched
	MatchedPageURL    *string `json:"matched_page_url,omitempty"` // URL of matched page
}

// handleDiscoverKeywords handles POST /api/v1/projects/:id/discover-keywords
func (s *Server) handleDiscoverKeywords(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	// Verify project access
	hasAccess, err := s.verifyProjectAccess(userID, projectID)
	if err != nil {
		s.logger.Error("Failed to verify project access", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to verify project access")
		return
	}
	if !hasAccess {
		s.respondError(w, http.StatusForbidden, "Access denied")
		return
	}

	// Parse request body
	var req DiscoverKeywordsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
		return
	}

	// Validate required fields
	if req.Target == "" {
		s.respondError(w, http.StatusBadRequest, "target is required")
		return
	}
	if req.LocationName == "" {
		req.LocationName = "United States" // Default
	}
	if req.LanguageName == "" {
		req.LanguageName = "English" // Default
	}
	if req.Limit == 0 {
		req.Limit = 1000 // Default limit
	}
	if req.Limit > 10000 {
		req.Limit = 10000 // Max limit
	}

	// Initialize DataForSEO client
	client, err := dataforseo.NewClient()
	if err != nil {
		s.logger.Error("Failed to create DataForSEO client", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "DataForSEO integration not configured")
		return
	}

	// Call DataForSEO Ranked Keywords API
	// Note: Based on API example, only these fields are supported:
	// target, language_name, location_name, load_rank_absolute, load_keyword_info, limit
	// Filters, SortBy, OrderBy are NOT supported in the live endpoint
	task := dataforseo.RankedKeywordsTask{
		Target:           req.Target,
		LocationName:     req.LocationName,
		LanguageName:     req.LanguageName,
		LoadRankAbsolute: true, // Load absolute rank as shown in API example
		LoadKeywordInfo:  true, // Load keyword metrics (search volume, competition, CPC)
		Limit:            req.Limit,
	}

	resp, err := client.GetRankedKeywordsLive(r.Context(), task)
	if err != nil {
		s.logger.Error("Failed to discover keywords",
			zap.Error(err),
			zap.String("target", req.Target),
			zap.String("location", req.LocationName),
			zap.String("language", req.LanguageName))

		// Check if it's a 40400 error (endpoint not found)
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "40400") || strings.Contains(errorMsg, "Not Found") {
			s.respondError(w, http.StatusBadRequest,
				"Ranked Keywords API endpoint not available. This feature may require a DataForSEO Labs subscription or the endpoint may not be available in your account tier.")
			return
		}

		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to discover keywords: %v", err))
		return
	}

	// Log response structure for debugging
	s.logger.Info("Ranked Keywords API response",
		zap.Int("status_code", resp.StatusCode),
		zap.String("status_message", resp.StatusMessage),
		zap.Int("tasks_count", len(resp.Tasks)))

	// Check response status code
	if resp.StatusCode != 20000 && resp.StatusCode != 0 {
		s.logger.Warn("Ranked Keywords API returned non-success status",
			zap.Int("status_code", resp.StatusCode),
			zap.String("status_message", resp.StatusMessage))
		s.respondError(w, http.StatusBadRequest,
			fmt.Sprintf("DataForSEO API error: %d - %s", resp.StatusCode, resp.StatusMessage))
		return
	}

	// Check if we got results
	if len(resp.Tasks) == 0 {
		s.logger.Info("No tasks in response", zap.String("target", req.Target))
		s.respondJSON(w, http.StatusOK, map[string]interface{}{
			"keywords": []DiscoveredKeywordResponse{},
			"count":    0,
			"message":  "No keywords found for this domain",
		})
		return
	}

	taskResult := resp.Tasks[0]
	s.logger.Info("Task result",
		zap.Int("task_status_code", taskResult.StatusCode),
		zap.String("task_status_message", taskResult.StatusMessage),
		zap.Int("result_count", len(taskResult.Result)),
		zap.String("target", req.Target))

	if len(taskResult.Result) == 0 {
		s.logger.Info("No results in task", zap.String("target", req.Target))
		s.respondJSON(w, http.StatusOK, map[string]interface{}{
			"keywords": []DiscoveredKeywordResponse{},
			"count":    0,
			"message":  "No keywords found for this domain",
		})
		return
	}

	if len(taskResult.Result[0].Items) == 0 {
		s.logger.Info("No items in result",
			zap.String("target", req.Target),
			zap.Int("task_status_code", taskResult.StatusCode))
		s.respondJSON(w, http.StatusOK, map[string]interface{}{
			"keywords": []DiscoveredKeywordResponse{},
			"count":    0,
			"message":  "No keywords found for this domain",
		})
		return
	}

	s.logger.Info("Found keywords",
		zap.Int("count", len(taskResult.Result[0].Items)),
		zap.String("target", req.Target))

	// Debug: Log first item structure to see what fields are available
	if len(taskResult.Result[0].Items) > 0 {
		firstItem := taskResult.Result[0].Items[0]
		// Log raw JSON of first item to see actual structure
		itemJSON, _ := json.Marshal(firstItem)
		var cpc float64
		if firstItem.KeywordData.KeywordInfo.CPC != nil {
			cpc = *firstItem.KeywordData.KeywordInfo.CPC
		}
		s.logger.Info("Sample keyword item structure",
			zap.String("keyword", firstItem.KeywordData.Keyword),
			zap.Int("search_volume", firstItem.KeywordData.KeywordInfo.SearchVolume),
			zap.String("competition", firstItem.KeywordData.KeywordInfo.CompetitionLevel),
			zap.Float64("cpc", cpc),
			zap.String("url", firstItem.RankedSERPElement.SERPItem.URL),
			zap.String("title", firstItem.RankedSERPElement.SERPItem.Title),
			zap.String("raw_json", string(itemJSON)))
	}

	// Get all crawled pages for this project to match URLs
	pageData, _, err := s.serviceRole.From("pages").
		Select("id,url", "", false).
		Eq("project_id", projectID).
		Execute()

	pageMap := make(map[string]int64) // URL -> page ID
	if err == nil {
		var pages []map[string]interface{}
		if err := json.Unmarshal(pageData, &pages); err == nil {
			for _, page := range pages {
				url := getKeywordString(page, "url")
				if url != "" {
					// Normalize URL for matching (remove trailing slash, lowercase)
					normalizedURL := strings.TrimSuffix(strings.ToLower(url), "/")
					if pageIDFloat, ok := page["id"].(float64); ok {
						pageMap[normalizedURL] = int64(pageIDFloat)
					}
				}
			}
		}
	}

	// Convert to response format
	resultItems := taskResult.Result[0].Items
	s.logger.Debug("Processing items",
		zap.Int("items_count", len(resultItems)))

	discovered := make([]DiscoveredKeywordResponse, 0, len(resultItems))
	for _, item := range resultItems {
		// Get position (prefer rank_group, fallback to rank_absolute)
		position := item.RankedSERPElement.SERPItem.RankGroup
		if position == 0 {
			position = item.RankedSERPElement.SERPItem.RankAbsolute
		}

		// Apply position filters client-side if specified
		if req.MinPosition > 0 && position < req.MinPosition {
			continue
		}
		if req.MaxPosition > 0 && position > req.MaxPosition {
			continue
		}

		// Try to match URL to crawled page
		itemURL := strings.TrimSuffix(strings.ToLower(item.RankedSERPElement.SERPItem.URL), "/")
		var matchedPageID *int64
		var matchedPageURL *string
		if pageID, found := pageMap[itemURL]; found {
			matchedPageID = &pageID
			originalURL := item.RankedSERPElement.SERPItem.URL
			matchedPageURL = &originalURL
		}

		// Extract keyword from nested keyword_data structure
		keyword := item.KeywordData.Keyword

		// Extract keyword info from nested structure
		searchVolume := item.KeywordData.KeywordInfo.SearchVolume
		competition := item.KeywordData.KeywordInfo.CompetitionLevel
		var cpc float64
		if item.KeywordData.KeywordInfo.CPC != nil {
			cpc = *item.KeywordData.KeywordInfo.CPC
		}

		discovered = append(discovered, DiscoveredKeywordResponse{
			Keyword:           keyword,
			Position:          position,
			URL:               item.RankedSERPElement.SERPItem.URL,
			Title:             item.RankedSERPElement.SERPItem.Title,
			SearchVolume:      searchVolume,
			Competition:       competition,
			CPC:               cpc,
			KeywordDifficulty: item.RankedSERPElement.KeywordDifficulty,
			MatchedPageID:     matchedPageID,
			MatchedPageURL:    matchedPageURL,
		})
	}

	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"keywords": discovered,
		"count":    len(discovered),
		"target":   req.Target,
	})
}

func getKeywordString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}
