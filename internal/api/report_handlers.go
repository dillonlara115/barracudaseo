package api

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// CreatePublicReportRequest represents a request to create a public report
type CreatePublicReportRequest struct {
	CrawlID     string                 `json:"crawl_id"`
	Title       string                 `json:"title,omitempty"`
	Description string                 `json:"description,omitempty"`
	Password    string                 `json:"password,omitempty"`        // Optional password for protection
	ExpiresIn   *int                   `json:"expires_in_days,omitempty"` // Optional expiry in days
	Settings    map[string]interface{} `json:"settings,omitempty"`        // Report settings
}

// PublicReportResponse represents a public report response
type PublicReportResponse struct {
	ID          string     `json:"id"`
	AccessToken string     `json:"access_token"`
	PublicURL   string     `json:"public_url"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// ViewPublicReportRequest represents a request to view a public report (with optional password)
type ViewPublicReportRequest struct {
	Password string `json:"password,omitempty"`
}

// generateSecureToken generates a cryptographically secure random token
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// handleCreatePublicReport handles POST /api/v1/reports/public
func (s *Server) handleCreatePublicReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, ok := userIDFromContext(r.Context())
	if !ok {
		s.respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Check subscription - public reports require Pro
	if sub := s.requireProSubscription(w, userID, "Public Report Sharing"); sub == nil {
		return
	}

	var req CreatePublicReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
		return
	}

	if req.CrawlID == "" {
		s.respondError(w, http.StatusBadRequest, "crawl_id is required")
		return
	}

	// Verify user has access to the crawl
	crawlID := req.CrawlID
	crawlDataBytes, _, err := s.serviceRole.From("crawls").Select("*", "", false).Eq("id", crawlID).Execute()
	if err != nil {
		s.logger.Error("Failed to fetch crawl", zap.Error(err))
		s.respondError(w, http.StatusNotFound, "Crawl not found")
		return
	}

	var crawls []map[string]interface{}
	if err := json.Unmarshal(crawlDataBytes, &crawls); err != nil || len(crawls) == 0 {
		s.logger.Error("Failed to parse crawl data", zap.Error(err))
		s.respondError(w, http.StatusNotFound, "Crawl not found")
		return
	}
	crawlData := crawls[0]

	projectID, ok := crawlData["project_id"].(string)
	if !ok {
		s.respondError(w, http.StatusInternalServerError, "Invalid crawl data")
		return
	}

	// Verify user has access to the project
	hasAccess, err := s.verifyProjectAccess(userID, projectID)
	if err != nil || !hasAccess {
		s.respondError(w, http.StatusForbidden, "You don't have access to this project")
		return
	}

	// Generate secure access token (32 characters)
	accessToken, err := generateSecureToken(32)
	if err != nil {
		s.logger.Error("Failed to generate access token", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to generate access token")
		return
	}

	// Hash password if provided
	var passwordHash *string
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			s.logger.Error("Failed to hash password", zap.Error(err))
			s.respondError(w, http.StatusInternalServerError, "Failed to hash password")
			return
		}
		hashStr := string(hash)
		passwordHash = &hashStr
	}

	// Calculate expiry date if provided
	var expiresAt *time.Time
	if req.ExpiresIn != nil && *req.ExpiresIn > 0 {
		expiry := time.Now().Add(time.Duration(*req.ExpiresIn) * 24 * time.Hour)
		expiresAt = &expiry
	}

	// Default title if not provided
	title := req.Title
	if title == "" {
		title = fmt.Sprintf("SEO Audit Report - %s", time.Now().Format("January 2, 2006"))
	}

	// Prepare report data
	reportID := uuid.New().String()
	reportData := map[string]interface{}{
		"id":           reportID,
		"crawl_id":     crawlID,
		"project_id":   projectID,
		"created_by":   userID,
		"access_token": accessToken,
		"title":        title,
		"description":  req.Description,
		"settings":     req.Settings,
		"view_count":   0,
	}

	if passwordHash != nil {
		reportData["password_hash"] = *passwordHash
	}
	if expiresAt != nil {
		reportData["expires_at"] = expiresAt.Format(time.RFC3339)
	}

	// Insert report into database
	_, _, err = s.serviceRole.From("public_reports").Insert(reportData, false, "", "", "").Execute()
	if err != nil {
		s.logger.Error("Failed to create public report", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to create public report")
		return
	}

	// Determine public URL base (use hash-based routing for svelte-spa-router)
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "http://localhost:5173" // Default to Vite dev server
	}
	// Use hash-based routing: #/reports/:token
	publicURL := fmt.Sprintf("%s#/reports/%s", strings.TrimSuffix(appURL, "/"), accessToken)

	response := PublicReportResponse{
		ID:          reportID,
		AccessToken: accessToken,
		PublicURL:   publicURL,
		ExpiresAt:   expiresAt,
		CreatedAt:   time.Now(),
	}

	s.respondJSON(w, http.StatusCreated, response)
}

// handleViewPublicReport handles GET /api/public/reports/:token
// This endpoint does NOT require authentication
func (s *Server) handleViewPublicReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract token from URL path
	// Path format: /api/public/reports/:token
	path := strings.TrimPrefix(r.URL.Path, "/api/public/reports/")
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 || parts[0] == "" {
		s.respondError(w, http.StatusBadRequest, "Access token is required")
		return
	}
	accessToken := parts[0]

	// Fetch report from database
	reportDataBytes, _, err := s.serviceRole.From("public_reports").Select("*", "", false).Eq("access_token", accessToken).Execute()
	if err != nil {
		s.logger.Error("Failed to fetch public report", zap.Error(err))
		s.respondError(w, http.StatusNotFound, "Report not found")
		return
	}

	var reports []map[string]interface{}
	if err := json.Unmarshal(reportDataBytes, &reports); err != nil || len(reports) == 0 {
		s.logger.Error("Failed to parse report data", zap.Error(err))
		s.respondError(w, http.StatusNotFound, "Report not found")
		return
	}
	reportData := reports[0]

	// Check if report has expired
	if expiresAtStr, ok := reportData["expires_at"].(string); ok && expiresAtStr != "" {
		expiresAt, err := time.Parse(time.RFC3339, expiresAtStr)
		if err == nil && time.Now().After(expiresAt) {
			s.respondError(w, http.StatusGone, "This report has expired")
			return
		}
	}

	// Check password if required
	if passwordHash, ok := reportData["password_hash"].(string); ok && passwordHash != "" {
		// Password is required
		var password string
		if r.Method == http.MethodPost {
			// Get password from request body
			var req ViewPublicReportRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				s.respondError(w, http.StatusBadRequest, "Password is required")
				return
			}
			password = req.Password
		} else {
			// GET request - check query parameter
			password = r.URL.Query().Get("password")
		}

		if password == "" {
			s.respondError(w, http.StatusUnauthorized, "Password is required")
			return
		}

		// Verify password
		if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
			s.respondError(w, http.StatusUnauthorized, "Invalid password")
			return
		}
	}

	// Update view count and last viewed timestamp
	viewCount, _ := reportData["view_count"].(float64)
	updateData := map[string]interface{}{
		"view_count":     int(viewCount) + 1,
		"last_viewed_at": time.Now().Format(time.RFC3339),
	}
	_, _, err = s.serviceRole.From("public_reports").Update(updateData, "", "").Eq("access_token", accessToken).Execute()
	if err != nil {
		s.logger.Warn("Failed to update view count", zap.Error(err))
		// Don't fail the request if view count update fails
	}

	// Fetch crawl data
	crawlID, _ := reportData["crawl_id"].(string)
	crawlDataBytes, _, err := s.serviceRole.From("crawls").Select("*", "", false).Eq("id", crawlID).Execute()
	if err != nil {
		s.logger.Error("Failed to fetch crawl", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to load crawl data")
		return
	}

	var crawls []map[string]interface{}
	if err := json.Unmarshal(crawlDataBytes, &crawls); err != nil || len(crawls) == 0 {
		s.logger.Error("Failed to parse crawl data", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to load crawl data")
		return
	}
	crawlData := crawls[0]

	// Fetch project data
	projectID, _ := reportData["project_id"].(string)
	projectDataBytes, _, err := s.serviceRole.From("projects").Select("*", "", false).Eq("id", projectID).Execute()
	var projectData map[string]interface{}
	if err != nil {
		s.logger.Warn("Failed to fetch project data", zap.Error(err))
		projectData = map[string]interface{}{} // Return empty object on error
	} else {
		var projects []map[string]interface{}
		if err := json.Unmarshal(projectDataBytes, &projects); err == nil && len(projects) > 0 {
			projectData = projects[0]
		} else {
			projectData = map[string]interface{}{}
		}
	}

	// Fetch issues for this crawl
	issuesBytes, _, err := s.serviceRole.From("issues").Select("*", "", false).Eq("crawl_id", crawlID).Execute()
	var issues []map[string]interface{}
	if err != nil {
		s.logger.Error("Failed to fetch issues", zap.Error(err))
		issues = []map[string]interface{}{} // Return empty array on error
	} else {
		json.Unmarshal(issuesBytes, &issues)
	}

	// Fetch pages for this crawl
	pagesBytes, _, err := s.serviceRole.From("pages").Select("*", "", false).Eq("crawl_id", crawlID).Execute()
	var pages []map[string]interface{}
	if err != nil {
		s.logger.Error("Failed to fetch pages", zap.Error(err))
		pages = []map[string]interface{}{} // Return empty array on error
	} else {
		json.Unmarshal(pagesBytes, &pages)
	}

	// Enrich issues with page URLs
	// Create a map of page_id to URL for quick lookup
	pageURLMap := make(map[int64]string)
	for _, page := range pages {
		if pageID, ok := page["id"].(float64); ok {
			if url, ok := page["url"].(string); ok {
				pageURLMap[int64(pageID)] = url
			}
		}
	}

	// Add URL to each issue
	for i := range issues {
		if pageID, ok := issues[i]["page_id"].(float64); ok && pageID > 0 {
			if url, found := pageURLMap[int64(pageID)]; found {
				issues[i]["url"] = url
			}
		}
	}

	// Build response
	response := map[string]interface{}{
		"report": map[string]interface{}{
			"id":          reportData["id"],
			"title":       reportData["title"],
			"description": reportData["description"],
			"created_at":  reportData["created_at"],
		},
		"project": projectData,
		"crawl":   crawlData,
		"issues":  issues,
		"pages":   pages,
		"summary": map[string]interface{}{
			"total_pages":  len(pages),
			"total_issues": len(issues),
		},
	}

	s.respondJSON(w, http.StatusOK, response)
}

// handlePublicReports handles GET /api/v1/reports/public (list) and POST /api/v1/reports/public (create)
func (s *Server) handlePublicReports(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s.handleCreatePublicReport(w, r)
		return
	}
	if r.Method == http.MethodGet {
		s.handleListPublicReports(w, r)
		return
	}
	s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
}

// handleListPublicReports handles GET /api/v1/reports/public
func (s *Server) handleListPublicReports(w http.ResponseWriter, r *http.Request) {

	userID, ok := userIDFromContext(r.Context())
	if !ok {
		s.respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Optional project filter
	projectID := r.URL.Query().Get("project_id")

	// Build query
	query := s.serviceRole.From("public_reports").Select("*", "", false).Eq("created_by", userID)

	if projectID != "" {
		query = query.Eq("project_id", projectID)
	}

	reportsBytes, _, err := query.Order("created_at", nil).Execute()
	if err != nil {
		s.logger.Error("Failed to fetch public reports", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to fetch reports")
		return
	}

	var reports []map[string]interface{}
	if err := json.Unmarshal(reportsBytes, &reports); err != nil {
		s.logger.Error("Failed to parse reports", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to parse reports")
		return
	}

	// Determine public URL base
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "http://localhost:5173"
	}

	// Add public URLs to each report (use hash-based routing)
	for i := range reports {
		if accessToken, ok := reports[i]["access_token"].(string); ok {
			reports[i]["public_url"] = fmt.Sprintf("%s#/reports/%s", strings.TrimSuffix(appURL, "/"), accessToken)
		}
		// Don't expose password hash
		delete(reports[i], "password_hash")
	}

	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"reports": reports,
	})
}

// handlePublicReportByID handles GET /api/v1/reports/public/:id and DELETE /api/v1/reports/public/:id
func (s *Server) handlePublicReportByID(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		s.handleDeletePublicReport(w, r)
		return
	}
	// Could add GET for single report details if needed
	s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
}

// handleDeletePublicReport handles DELETE /api/v1/reports/public/:id
func (s *Server) handleDeletePublicReport(w http.ResponseWriter, r *http.Request) {

	userID, ok := userIDFromContext(r.Context())
	if !ok {
		s.respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Extract report ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/reports/public/")
	path = strings.Trim(path, "/")
	reportID := path

	if reportID == "" {
		s.respondError(w, http.StatusBadRequest, "Report ID is required")
		return
	}

	// Verify user owns this report
	reportDataBytes, _, err := s.serviceRole.From("public_reports").Select("*", "", false).Eq("id", reportID).Execute()
	if err != nil {
		s.logger.Error("Failed to fetch report", zap.Error(err))
		s.respondError(w, http.StatusNotFound, "Report not found")
		return
	}

	var reports []map[string]interface{}
	if err := json.Unmarshal(reportDataBytes, &reports); err != nil || len(reports) == 0 {
		s.logger.Error("Failed to parse report data", zap.Error(err))
		s.respondError(w, http.StatusNotFound, "Report not found")
		return
	}
	reportData := reports[0]

	reportUserID, _ := reportData["created_by"].(string)
	if reportUserID != userID {
		s.respondError(w, http.StatusForbidden, "You don't have permission to delete this report")
		return
	}

	// Delete report
	_, _, err = s.serviceRole.From("public_reports").Delete("", "").Eq("id", reportID).Execute()
	if err != nil {
		s.logger.Error("Failed to delete report", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to delete report")
		return
	}

	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}
