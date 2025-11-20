package api

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/dillonlara115/barracuda/internal/gsc"
	"go.uber.org/zap"
)

func (s *Server) handleProjectGSC(w http.ResponseWriter, r *http.Request, projectID, userID string, segments []string) {
	// Verify access upfront
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

	if len(segments) == 0 || segments[0] == "" {
		s.handleProjectGSCStatus(w, r, projectID)
		return
	}

	s.logger.Debug("GSC handler", zap.String("projectID", projectID), zap.Strings("segments", segments), zap.String("first_segment", segments[0]), zap.String("path", r.URL.Path))

	switch segments[0] {
	case "connect":
		s.handleProjectGSCConnect(w, r, projectID)
	case "disconnect":
		s.handleProjectGSCDisconnect(w, r, projectID)
	case "properties":
		s.handleProjectGSCProperties(w, r, projectID)
	case "property":
		s.handleProjectGSCSetProperty(w, r, projectID)
	case "trigger-sync":
		s.handleProjectGSCTriggerSync(w, r, projectID)
	case "status":
		s.handleProjectGSCStatus(w, r, projectID)
	case "summary":
		s.handleProjectGSCStatus(w, r, projectID)
	case "dimensions":
		s.handleProjectGSCDimensionsDirect(w, r, projectID)
	default:
		s.respondError(w, http.StatusNotFound, fmt.Sprintf("Unknown GSC resource: %s", segments[0]))
	}
}

func (s *Server) handleProjectGSCConnect(w http.ResponseWriter, r *http.Request, projectID string) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	authURL, state, err := gsc.GenerateAuthURL(projectID)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate auth URL: %v", err))
		return
	}

	if _, err := s.ensureGSCSyncState(projectID, ""); err != nil {
		s.logger.Warn("Failed to ensure GSC sync state", zap.Error(err))
	}

	s.respondJSON(w, http.StatusOK, map[string]string{
		"auth_url": authURL,
		"state":    state,
	})
}

func (s *Server) handleProjectGSCProperties(w http.ResponseWriter, r *http.Request, projectID string) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	cfg, err := s.loadTokenIntoMemory(projectID)
	if err != nil {
		s.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	properties, err := gsc.GetProperties(projectID)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get properties: %v", err))
		return
	}

	var selected string
	if cfg != nil {
		selected = cfg.PropertyURL
	}

	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"properties":       properties,
		"selectedProperty": selected,
	})
}

func (s *Server) handleProjectGSCSetProperty(w http.ResponseWriter, r *http.Request, projectID string) {
	if r.Method != http.MethodPost && r.Method != http.MethodPatch {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		PropertyURL  string `json:"property_url"`
		PropertyType string `json:"property_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %v", err))
		return
	}

	if req.PropertyURL == "" {
		s.respondError(w, http.StatusBadRequest, "property_url is required")
		return
	}

	cfg, _, err := s.getGSCIntegration(projectID)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to load integration")
		return
	}
	if cfg == nil {
		s.respondError(w, http.StatusBadRequest, "Connect Google Search Console before selecting a property")
		return
	}

	cfg.PropertyURL = req.PropertyURL
	cfg.PropertyType = req.PropertyType

	if err := s.saveGSCIntegration(projectID, cfg); err != nil {
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update integration: %v", err))
		return
	}

	if _, err := s.ensureGSCSyncState(projectID, req.PropertyURL); err != nil {
		s.logger.Warn("Failed to update GSC sync state property URL", zap.Error(err))
	}

	s.respondJSON(w, http.StatusOK, map[string]string{
		"property_url":  cfg.PropertyURL,
		"property_type": cfg.PropertyType,
	})
}

func (s *Server) handleProjectGSCTriggerSync(w http.ResponseWriter, r *http.Request, projectID string) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		LookbackDays int    `json:"lookback_days"`
		Period       string `json:"period"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err.Error() != "EOF" {
		s.respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %v", err))
		return
	}

	if req.LookbackDays <= 0 {
		req.LookbackDays = 30
	}
	if req.Period == "" {
		req.Period = fmt.Sprintf("last_%d_days", req.LookbackDays)
	}

	cfg, err := s.loadTokenIntoMemory(projectID)
	if err != nil {
		s.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := s.updateGSCSyncState(projectID, "running", nil, nil); err != nil {
		s.logger.Warn("Failed to mark sync running", zap.Error(err))
	}

	if err := s.syncProjectGSCData(projectID, cfg, req.LookbackDays, req.Period); err != nil {
		s.logger.Error("GSC sync failed", zap.Error(err))
		_ = s.updateGSCSyncState(projectID, "error", nil, map[string]interface{}{
			"message": err.Error(),
			"time":    time.Now().UTC().Format(time.RFC3339),
		})
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Sync failed: %v", err))
		return
	}

	now := time.Now().UTC()
	if err := s.updateGSCSyncState(projectID, "idle", &now, nil); err != nil {
		s.logger.Warn("Failed to finalize sync state", zap.Error(err))
	}

	s.respondJSON(w, http.StatusOK, map[string]string{
		"status":         "completed",
		"last_synced_at": now.Format(time.RFC3339),
	})
}

func (s *Server) handleProjectGSCStatus(w http.ResponseWriter, r *http.Request, projectID string) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	cfg, _, err := s.getGSCIntegration(projectID)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to load integration")
		return
	}

	state, err := s.ensureGSCSyncState(projectID, "")
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to load sync state")
		return
	}

	summary, err := s.fetchLatestGSCSummary(projectID)
	if err != nil {
		s.logger.Warn("Failed to fetch GSC summary", zap.Error(err))
	}

	response := map[string]interface{}{
		"integration": cfg,
		"sync_state":  state,
		"summary":     summary,
	}
	s.respondJSON(w, http.StatusOK, response)
}

func (s *Server) fetchLatestGSCSummary(projectID string) (map[string]interface{}, error) {
	data, _, err := s.serviceRole.
		From("gsc_performance_snapshots").
		Select("*", "", false).
		Eq("project_id", projectID).
		Execute()
	if err != nil {
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

func parseDateField(value interface{}) time.Time {
	switch v := value.(type) {
	case time.Time:
		return v
	case string:
		if ts, err := time.Parse("2006-01-02", v); err == nil {
			return ts
		}
		if ts, err := time.Parse(time.RFC3339, v); err == nil {
			return ts
		}
	}
	return time.Time{}
}

func (s *Server) handleProjectGSCDimensionsDirect(w http.ResponseWriter, r *http.Request, projectID string) {
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
		From("gsc_performance_rows").
		Select("*", "", false).
		Eq("project_id", projectID).
		Eq("row_type", rowType).
		Execute()
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to query GSC rows: %v", err))
		return
	}

	var rows []map[string]interface{}
	if err := json.Unmarshal(data, &rows); err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to parse rows")
		return
	}

	// Sort by impressions descending if available
	sort.Slice(rows, func(i, j int) bool {
		metI, _ := rows[i]["metrics"].(map[string]interface{})
		metJ, _ := rows[j]["metrics"].(map[string]interface{})
		return getFloat(metI["impressions"]) > getFloat(metJ["impressions"])
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

func (s *Server) handleGSCGlobalSync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	if s.cronSecret == "" {
		s.respondError(w, http.StatusServiceUnavailable, "Cron sync secret not configured")
		return
	}

	secret := r.Header.Get("X-Cron-Secret")
	if secret == "" {
		secret = r.URL.Query().Get("secret")
	}

	if subtle.ConstantTimeCompare([]byte(secret), []byte(s.cronSecret)) != 1 {
		s.respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req struct {
		ProjectIDs   []string `json:"project_ids"`
		LookbackDays int      `json:"lookback_days"`
	}

	if r.Body != nil {
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err != http.ErrBodyNotAllowed && err.Error() != "EOF" {
			s.respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
			return
		}
	}

	if req.LookbackDays <= 0 {
		req.LookbackDays = 30
	}

	targetProjects := req.ProjectIDs
	if len(targetProjects) == 0 {
		data, _, err := s.serviceRole.
			From("api_integrations").
			Select("project_id", "", false).
			Eq("provider", "gsc").
			Execute()
		if err != nil {
			s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to load integrations: %v", err))
			return
		}

		var rows []map[string]interface{}
		if err := json.Unmarshal(data, &rows); err != nil {
			s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to parse integrations: %v", err))
			return
		}

		unique := make(map[string]struct{}, len(rows))
		for _, row := range rows {
			if projectID, ok := row["project_id"].(string); ok && projectID != "" {
				unique[projectID] = struct{}{}
			}
		}

		targetProjects = make([]string, 0, len(unique))
		for id := range unique {
			targetProjects = append(targetProjects, id)
		}
	}

	if len(targetProjects) == 0 {
		s.respondJSON(w, http.StatusOK, map[string]interface{}{
			"run_count": 0,
			"results":   []map[string]interface{}{},
		})
		return
	}

	results := make([]map[string]interface{}, 0, len(targetProjects))

	for _, projectID := range targetProjects {
		entry := map[string]interface{}{
			"project_id": projectID,
			"status":     "skipped",
		}

		cfg, _, err := s.getGSCIntegration(projectID)
		if err != nil {
			entry["status"] = "error"
			entry["error"] = fmt.Sprintf("failed to load integration: %v", err)
			results = append(results, entry)
			continue
		}
		if cfg == nil || cfg.PropertyURL == "" {
			entry["status"] = "skipped"
			entry["message"] = "No connected GSC property"
			results = append(results, entry)
			continue
		}

		if _, err := s.ensureGSCSyncState(projectID, cfg.PropertyURL); err != nil {
			entry["status"] = "error"
			entry["error"] = fmt.Sprintf("failed to ensure sync state: %v", err)
			results = append(results, entry)
			continue
		}

		if err := s.updateGSCSyncState(projectID, "running", nil, nil); err != nil {
			entry["status"] = "error"
			entry["error"] = fmt.Sprintf("failed to mark running: %v", err)
			results = append(results, entry)
			continue
		}

		if err := s.syncProjectGSCData(projectID, cfg, req.LookbackDays, fmt.Sprintf("cron_last_%d_days", req.LookbackDays)); err != nil {
			entry["status"] = "error"
			entry["error"] = err.Error()
			results = append(results, entry)
			continue
		}

		entry["status"] = "synced"
		entry["lookback_days"] = req.LookbackDays
		results = append(results, entry)
	}

	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"run_count": len(results),
		"results":   results,
	})
}

func getFloat(v interface{}) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case float32:
		return float64(t)
	case int:
		return float64(t)
	case int64:
		return float64(t)
	case json.Number:
		f, _ := t.Float64()
		return f
	case string:
		f, _ := strconv.ParseFloat(t, 64)
		return f
	default:
		return 0
	}
}

func (s *Server) handleProjectGSCDisconnect(w http.ResponseWriter, r *http.Request, projectID string) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Delete integration record
	_, _, err := s.serviceRole.
		From("api_integrations").
		Delete("", "").
		Eq("project_id", projectID).
		Eq("provider", "gsc").
		Execute()

	if err != nil {
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to remove integration: %v", err))
		return
	}

	// Delete sync state
	_, _, err = s.serviceRole.
		From("gsc_sync_states").
		Delete("", "").
		Eq("project_id", projectID).
		Execute()

	if err != nil {
		s.logger.Warn("Failed to remove sync state", zap.Error(err))
		// Not fatal
	}

	// Also clear the gsc_property_url from project settings
	// We first fetch the project to get current settings
	data, _, err := s.serviceRole.
		From("projects").
		Select("settings", "", false).
		Eq("id", projectID).
		Execute()

	if err == nil && len(data) > 0 {
		var rows []struct {
			Settings map[string]interface{} `json:"settings"`
		}
		if err := json.Unmarshal(data, &rows); err == nil && len(rows) > 0 {
			settings := rows[0].Settings
			if settings == nil {
				settings = make(map[string]interface{})
			}
			// Remove the property URL
			delete(settings, "gsc_property_url")
			
			// Update project
			_, _, _ = s.serviceRole.
				From("projects").
				Update(map[string]interface{}{"settings": settings}, "", "").
				Eq("id", projectID).
				Execute()
		}
	}
    
    // Clear in-memory token cache if present
    gsc.StoreToken(projectID, nil)

	s.respondJSON(w, http.StatusOK, map[string]string{
		"status": "disconnected",
	})
}
