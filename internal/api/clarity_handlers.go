package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dillonlara115/barracudaseo/internal/clarity"
	"go.uber.org/zap"
)

func (s *Server) handleProjectClarity(w http.ResponseWriter, r *http.Request, projectID, userID string, segments []string) {
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

	if sub := s.requireProSubscription(w, userID, "Microsoft Clarity integration"); sub == nil {
		return
	}

	if len(segments) == 0 || segments[0] == "" {
		s.handleProjectClarityStatus(w, r, projectID, userID)
		return
	}

	switch segments[0] {
	case "connect":
		s.handleProjectClarityConnect(w, r, projectID, userID)
	case "disconnect":
		s.handleProjectClarityDisconnect(w, r, projectID, userID)
	case "trigger-sync":
		s.handleProjectClarityTriggerSync(w, r, projectID, userID)
	case "status":
		s.handleProjectClarityStatus(w, r, projectID, userID)
	case "dimensions":
		s.handleProjectClarityDimensions(w, r, projectID)
	default:
		s.respondError(w, http.StatusNotFound, fmt.Sprintf("Unknown Clarity resource: %s", segments[0]))
	}
}

func (s *Server) handleProjectClarityConnect(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		ClarityProjectID string `json:"clarity_project_id"`
		APIToken         string `json:"api_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
		return
	}

	if req.ClarityProjectID == "" {
		s.respondError(w, http.StatusBadRequest, "clarity_project_id is required")
		return
	}
	if req.APIToken == "" {
		s.respondError(w, http.StatusBadRequest, "api_token is required")
		return
	}

	// Validate credentials with a test API call
	if err := clarity.ValidateToken(req.APIToken, req.ClarityProjectID); err != nil {
		s.respondError(w, http.StatusBadRequest, fmt.Sprintf("Failed to validate Clarity credentials: %v", err))
		return
	}

	cfg := &clarityIntegrationConfig{
		ClarityProjectID: req.ClarityProjectID,
		APIToken:         req.APIToken,
	}

	if err := s.saveClarityIntegration(userID, cfg); err != nil {
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to save integration: %v", err))
		return
	}

	if err := s.updateProjectSettings(projectID, map[string]interface{}{
		"clarity_project_id":          req.ClarityProjectID,
		"clarity_integration_user_id": userID,
	}); err != nil {
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update project settings: %v", err))
		return
	}

	if _, err := s.ensureClaritySyncState(projectID, req.ClarityProjectID); err != nil {
		s.logger.Warn("Failed to ensure clarity sync state", zap.Error(err))
	}

	s.respondJSON(w, http.StatusOK, map[string]string{
		"status":             "connected",
		"clarity_project_id": req.ClarityProjectID,
	})
}

func (s *Server) handleProjectClarityDisconnect(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	if err := s.updateProjectSettings(projectID, map[string]interface{}{
		"clarity_project_id":          nil,
		"clarity_integration_user_id": nil,
	}); err != nil {
		s.logger.Error("Failed to update project settings", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to disconnect")
		return
	}

	s.respondJSON(w, http.StatusOK, map[string]string{
		"status": "disconnected",
	})
}

func (s *Server) handleProjectClarityTriggerSync(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		NumDays int    `json:"num_days"`
		Period  string `json:"period"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err.Error() != "EOF" {
		s.respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %v", err))
		return
	}

	if req.NumDays <= 0 {
		req.NumDays = 3
	}
	if req.Period == "" {
		req.Period = fmt.Sprintf("last_%d_days", req.NumDays)
	}

	settings, err := s.loadProjectSettings(projectID)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to load project settings")
		return
	}
	clarityProjectID, _ := settings["clarity_project_id"].(string)
	integrationUserID, _ := settings["clarity_integration_user_id"].(string)
	if clarityProjectID == "" {
		s.respondError(w, http.StatusBadRequest, "Clarity not configured for this project")
		return
	}
	if integrationUserID == "" {
		s.respondError(w, http.StatusBadRequest, "No connected Clarity account for this project")
		return
	}

	cfg, _, err := s.getClarityIntegration(integrationUserID)
	if err != nil || cfg == nil {
		s.respondError(w, http.StatusBadRequest, "Clarity integration not found")
		return
	}
	cfg.ClarityProjectID = clarityProjectID

	if err := s.updateClaritySyncState(projectID, "running", nil, nil); err != nil {
		s.logger.Warn("Failed to mark sync running", zap.Error(err))
	}

	if err := s.syncProjectClarityData(projectID, userID, cfg, req.NumDays, req.Period); err != nil {
		s.logger.Error("Clarity sync failed", zap.Error(err))
		_ = s.updateClaritySyncState(projectID, "error", nil, map[string]interface{}{
			"message": err.Error(),
			"time":    time.Now().UTC().Format(time.RFC3339),
		})
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Sync failed: %v", err))
		return
	}

	now := time.Now().UTC()
	if err := s.updateClaritySyncState(projectID, "idle", &now, nil); err != nil {
		s.logger.Warn("Failed to finalize sync state", zap.Error(err))
	}

	s.respondJSON(w, http.StatusOK, map[string]string{
		"status":         "completed",
		"last_synced_at": now.Format(time.RFC3339),
	})
}

func (s *Server) handleProjectClarityStatus(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	settings, err := s.loadProjectSettings(projectID)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to load project settings")
		return
	}
	clarityProjectID, _ := settings["clarity_project_id"].(string)
	integrationUserID, _ := settings["clarity_integration_user_id"].(string)

	connected := false
	if integrationUserID != "" {
		cfg, _, err := s.getClarityIntegration(integrationUserID)
		if err == nil && cfg != nil && cfg.APIToken != "" {
			connected = true
		}
	}

	state, err := s.ensureClaritySyncState(projectID, clarityProjectID)
	if err != nil {
		s.logger.Warn("Failed to load clarity sync state", zap.Error(err))
		state = nil
	}

	// Fetch latest summary
	summary, err := s.fetchLatestClaritySummary(projectID)
	if err != nil {
		s.logger.Warn("Failed to fetch clarity summary", zap.Error(err))
		summary = nil
	}

	response := map[string]interface{}{
		"integration": map[string]interface{}{
			"clarity_project_id":  clarityProjectID,
			"integration_user_id": integrationUserID,
			"connected":           connected,
		},
		"sync_state": state,
		"summary":    summary,
	}
	s.respondJSON(w, http.StatusOK, response)
}

func (s *Server) handleProjectClarityDimensions(w http.ResponseWriter, r *http.Request, projectID string) {
	rowType := r.URL.Query().Get("type")
	if rowType == "" {
		s.respondError(w, http.StatusBadRequest, "type query parameter is required")
		return
	}

	limit := 50
	if v := r.URL.Query().Get("limit"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 && parsed <= 1000 {
			limit = parsed
		}
	}

	offset := 0
	if v := r.URL.Query().Get("offset"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	data, _, err := s.serviceRole.
		From("clarity_performance_rows").
		Select("*", "", false).
		Eq("project_id", projectID).
		Eq("row_type", rowType).
		Execute()
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to query Clarity rows: %v", err))
		return
	}

	var rows []map[string]interface{}
	if err := json.Unmarshal(data, &rows); err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to parse rows")
		return
	}

	// Sort by traffic descending
	sort.Slice(rows, func(i, j int) bool {
		metI, _ := rows[i]["metrics"].(map[string]interface{})
		metJ, _ := rows[j]["metrics"].(map[string]interface{})
		return getFloat(metI["traffic"]) > getFloat(metJ["traffic"])
	})

	total := len(rows)
	if offset > total {
		rows = []map[string]interface{}{}
	} else {
		end := offset + limit
		if end > total {
			end = total
		}
		rows = rows[offset:end]
	}

	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"rows":  rows,
		"total": total,
	})
}

func (s *Server) fetchLatestClaritySummary(projectID string) (map[string]interface{}, error) {
	data, _, err := s.serviceRole.
		From("clarity_performance_snapshots").
		Select("*", "", false).
		Eq("project_id", projectID).
		Execute()
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "does not exist") || strings.Contains(errStr, "relation") {
			return nil, nil
		}
		return nil, err
	}

	var snapshots []map[string]interface{}
	if err := json.Unmarshal(data, &snapshots); err != nil {
		return nil, err
	}

	if len(snapshots) == 0 {
		return nil, nil
	}

	latest := snapshots[0]
	if len(snapshots) > 1 {
		var latestDate time.Time
		for _, snap := range snapshots {
			captured := parseDateField(snap["captured_on"])
			if captured.After(latestDate) {
				latest = snap
				latestDate = captured
			}
		}
	}

	return latest, nil
}

// handleClarityGlobalSync handles POST /api/v1/clarity/sync - cron job for syncing all Clarity projects
func (s *Server) handleClarityGlobalSync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	if s.cronSecret == "" {
		s.respondError(w, http.StatusServiceUnavailable, "Cron sync secret not configured")
		return
	}

	secret := r.Header.Get("X-Cron-Secret")
	if secret == "" || secret != s.cronSecret {
		s.respondError(w, http.StatusUnauthorized, "Invalid or missing cron secret")
		return
	}

	// Find all projects with Clarity configured
	data, _, err := s.serviceRole.
		From("clarity_sync_states").
		Select("*", "", false).
		Execute()
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to query sync states: %v", err))
		return
	}

	var states []map[string]interface{}
	if err := json.Unmarshal(data, &states); err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to parse sync states")
		return
	}

	synced := 0
	errors := 0
	// Limit to 3 projects per cron run (3 API calls each = 9 of 10 daily limit)
	maxProjects := 3
	for i, state := range states {
		if i >= maxProjects {
			break
		}
		projectID, _ := state["project_id"].(string)
		clarityProjectID, _ := state["clarity_project_id"].(string)
		if projectID == "" || clarityProjectID == "" {
			continue
		}

		// Load project settings to get integration user
		settings, err := s.loadProjectSettings(projectID)
		if err != nil {
			s.logger.Warn("Clarity cron: failed to load project settings", zap.String("project_id", projectID), zap.Error(err))
			errors++
			continue
		}

		integrationUserID, _ := settings["clarity_integration_user_id"].(string)
		if integrationUserID == "" {
			continue
		}

		cfg, _, err := s.getClarityIntegration(integrationUserID)
		if err != nil || cfg == nil || cfg.APIToken == "" {
			continue
		}
		cfg.ClarityProjectID = clarityProjectID

		_ = s.updateClaritySyncState(projectID, "running", nil, nil)

		if err := s.syncProjectClarityData(projectID, integrationUserID, cfg, 3, "last_3_days"); err != nil {
			s.logger.Error("Clarity cron sync failed", zap.String("project_id", projectID), zap.Error(err))
			_ = s.updateClaritySyncState(projectID, "error", nil, map[string]interface{}{
				"message": err.Error(),
				"time":    time.Now().UTC().Format(time.RFC3339),
			})
			errors++
			continue
		}

		now := time.Now().UTC()
		_ = s.updateClaritySyncState(projectID, "idle", &now, nil)
		synced++
	}

	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"synced": synced,
		"errors": errors,
		"total":  len(states),
	})
}

// handleIntegrationsClarity handles /api/v1/integrations/clarity/* endpoints
func (s *Server) handleIntegrationsClarity(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		s.respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/integrations/clarity")
	path = strings.Trim(path, "/")
	segments := []string{}
	if path != "" {
		segments = strings.Split(path, "/")
	}

	if len(segments) == 0 || segments[0] == "" {
		segments = []string{"status"}
	}

	switch segments[0] {
	case "status":
		s.handleIntegrationsClarityStatus(w, r, userID)
	case "disconnect":
		s.handleIntegrationsClarityDisconnect(w, r, userID)
	default:
		s.respondError(w, http.StatusNotFound, "Unknown Clarity integration resource")
	}
}

func (s *Server) handleIntegrationsClarityStatus(w http.ResponseWriter, r *http.Request, userID string) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	cfg, _, err := s.getClarityIntegration(userID)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to load integration")
		return
	}

	connected := cfg != nil && cfg.APIToken != ""

	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"connected": connected,
	})
}

func (s *Server) handleIntegrationsClarityDisconnect(w http.ResponseWriter, r *http.Request, userID string) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	_, _, err := s.serviceRole.
		From("user_api_integrations").
		Delete("", "").
		Eq("user_id", userID).
		Eq("provider", "clarity").
		Execute()
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to disconnect integration: %v", err))
		return
	}

	s.respondJSON(w, http.StatusOK, map[string]string{
		"status": "disconnected",
	})
}
