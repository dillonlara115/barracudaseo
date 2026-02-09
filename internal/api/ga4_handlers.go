package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dillonlara115/barracudaseo/internal/ga4"
	"go.uber.org/zap"
)

func (s *Server) handleProjectGA4(w http.ResponseWriter, r *http.Request, projectID, userID string, segments []string) {
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

	// Check subscription - GA4 integration requires Pro
	if sub := s.requireProSubscription(w, userID, "Google Analytics integration"); sub == nil {
		return
	}

	if len(segments) == 0 || segments[0] == "" {
		s.handleProjectGA4Status(w, r, projectID)
		return
	}

	s.logger.Debug("GA4 handler", zap.String("projectID", projectID), zap.Strings("segments", segments))

	switch segments[0] {
	case "connect":
		s.handleProjectGA4Connect(w, r, userID)
	case "disconnect":
		s.handleProjectGA4Disconnect(w, r, projectID)
	case "properties":
		s.handleProjectGA4Properties(w, r, userID)
	case "property":
		s.handleProjectGA4SetProperty(w, r, projectID, userID)
	case "trigger-sync":
		s.handleProjectGA4TriggerSync(w, r, projectID)
	case "status":
		s.handleProjectGA4Status(w, r, projectID)
	case "summary":
		s.handleProjectGA4Status(w, r, projectID)
	case "dimensions":
		s.handleProjectGA4DimensionsDirect(w, r, projectID)
	default:
		s.respondError(w, http.StatusNotFound, fmt.Sprintf("Unknown GA4 resource: %s", segments[0]))
	}
}

func (s *Server) handleProjectGA4Connect(w http.ResponseWriter, r *http.Request, userID string) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	authURL, state, err := ga4.GenerateAuthURL(userID)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to generate auth URL: %v", err))
		return
	}

	s.respondJSON(w, http.StatusOK, map[string]string{
		"auth_url": authURL,
		"state":    state,
	})
}

func (s *Server) handleProjectGA4Properties(w http.ResponseWriter, r *http.Request, userID string) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	cfg, err := s.loadGA4TokenIntoMemory(userID)
	if err != nil {
		s.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	properties, err := ga4.GetProperties(userID)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get properties: %v", err))
		return
	}

	var selected string
	if cfg != nil {
		selected = cfg.PropertyID
	}

	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"properties":       properties,
		"selectedProperty": selected,
	})
}

func (s *Server) handleProjectGA4SetProperty(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	if r.Method != http.MethodPost && r.Method != http.MethodPatch {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		PropertyID   string `json:"property_id"`
		PropertyName string `json:"property_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request: %v", err))
		return
	}

	if req.PropertyID == "" {
		s.respondError(w, http.StatusBadRequest, "property_id is required")
		return
	}

	cfg, _, err := s.getGA4Integration(userID)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to load integration")
		return
	}
	if cfg == nil {
		s.respondError(w, http.StatusBadRequest, "Connect Google Analytics 4 before selecting a property")
		return
	}
	if err := s.updateProjectSettings(projectID, map[string]interface{}{
		"ga4_property_id":         req.PropertyID,
		"ga4_property_name":       req.PropertyName,
		"ga4_integration_user_id": userID,
	}); err != nil {
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update project settings: %v", err))
		return
	}

	if _, err := s.ensureGA4SyncState(projectID, req.PropertyID); err != nil {
		s.logger.Warn("Failed to update GA4 sync state property ID", zap.Error(err))
	}

	s.respondJSON(w, http.StatusOK, map[string]string{
		"property_id":   req.PropertyID,
		"property_name": req.PropertyName,
	})
}

func (s *Server) handleProjectGA4TriggerSync(w http.ResponseWriter, r *http.Request, projectID string) {
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

	settings, err := s.loadProjectSettings(projectID)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to load project settings")
		return
	}
	propertyID, _ := settings["ga4_property_id"].(string)
	integrationUserID, _ := settings["ga4_integration_user_id"].(string)
	if propertyID == "" {
		s.respondError(w, http.StatusBadRequest, "GA4 property not selected for this project")
		return
	}
	if integrationUserID == "" {
		s.respondError(w, http.StatusBadRequest, "No connected GA4 account for this project")
		return
	}

	cfg, err := s.loadGA4TokenIntoMemory(integrationUserID)
	if err != nil {
		s.respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	cfg.PropertyID = propertyID

	if err := s.updateGA4SyncState(projectID, "running", nil, nil); err != nil {
		s.logger.Warn("Failed to mark sync running", zap.Error(err))
	}

	if err := s.syncProjectGA4Data(projectID, integrationUserID, cfg, req.LookbackDays, req.Period); err != nil {
		s.logger.Error("GA4 sync failed", zap.Error(err))
		_ = s.updateGA4SyncState(projectID, "error", nil, map[string]interface{}{
			"message": err.Error(),
			"time":    time.Now().UTC().Format(time.RFC3339),
		})
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Sync failed: %v", err))
		return
	}

	now := time.Now().UTC()
	if err := s.updateGA4SyncState(projectID, "idle", &now, nil); err != nil {
		s.logger.Warn("Failed to finalize sync state", zap.Error(err))
	}

	s.respondJSON(w, http.StatusOK, map[string]string{
		"status":         "completed",
		"last_synced_at": now.Format(time.RFC3339),
	})
}

func (s *Server) handleProjectGA4Status(w http.ResponseWriter, r *http.Request, projectID string) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	settings, err := s.loadProjectSettings(projectID)
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to load project settings")
		return
	}
	propertyID, _ := settings["ga4_property_id"].(string)
	propertyName, _ := settings["ga4_property_name"].(string)
	integrationUserID, _ := settings["ga4_integration_user_id"].(string)

	connected := false
	if integrationUserID != "" {
		cfg, _, err := s.getGA4Integration(integrationUserID)
		if err == nil && cfg != nil && (cfg.AccessToken != "" || cfg.RefreshToken != "") {
			connected = true
		}
	}

	state, err := s.ensureGA4SyncState(projectID, propertyID)
	if err != nil {
		s.logger.Warn("Failed to load GA4 sync state", zap.Error(err))
		// Return empty state if table doesn't exist or other error (not fatal)
		state = nil
	}

	summary, err := s.fetchLatestGA4Summary(projectID)
	if err != nil {
		s.logger.Warn("Failed to fetch GA4 summary", zap.Error(err))
		summary = nil
	}

	response := map[string]interface{}{
		"integration": map[string]interface{}{
			"property_id":         propertyID,
			"property_name":       propertyName,
			"integration_user_id": integrationUserID,
			"connected":           connected,
		},
		"sync_state": state,
		"summary":    summary,
	}
	s.respondJSON(w, http.StatusOK, response)
}

func (s *Server) handleProjectGA4Disconnect(w http.ResponseWriter, r *http.Request, projectID string) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	if err := s.updateProjectSettings(projectID, map[string]interface{}{
		"ga4_property_id":         nil,
		"ga4_property_name":       nil,
		"ga4_integration_user_id": nil,
	}); err != nil {
		s.logger.Error("Failed to update GA4 project settings", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to disconnect")
		return
	}

	s.respondJSON(w, http.StatusOK, map[string]string{
		"status": "disconnected",
	})
}

// handleIntegrationsGA4 handles /api/v1/integrations/ga4/* endpoints
func (s *Server) handleIntegrationsGA4(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		s.respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/integrations/ga4")
	path = strings.Trim(path, "/")
	segments := []string{}
	if path != "" {
		segments = strings.Split(path, "/")
	}

	if len(segments) == 0 || segments[0] == "" {
		segments = []string{"status"}
	}

	switch segments[0] {
	case "connect":
		s.handleProjectGA4Connect(w, r, userID)
	case "disconnect":
		s.handleIntegrationsGA4Disconnect(w, r, userID)
	case "properties":
		s.handleProjectGA4Properties(w, r, userID)
	case "status":
		s.handleIntegrationsGA4Status(w, r, userID)
	default:
		s.respondError(w, http.StatusNotFound, "Unknown GA4 integration resource")
	}
}

func (s *Server) handleIntegrationsGA4Status(w http.ResponseWriter, r *http.Request, userID string) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	cfg, _, err := s.getGA4Integration(userID)
	if err != nil {
		s.logger.Warn("Failed to load GA4 integration status", zap.Error(err), zap.String("user_id", userID))
		// Return disconnected rather than 500 — table may not exist or other transient issue
		s.respondJSON(w, http.StatusOK, map[string]interface{}{
			"connected": false,
		})
		return
	}

	connected := cfg != nil && (cfg.AccessToken != "" || cfg.RefreshToken != "")

	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"connected": connected,
	})
}

func (s *Server) handleIntegrationsGA4Disconnect(w http.ResponseWriter, r *http.Request, userID string) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	_, _, err := s.serviceRole.
		From("user_api_integrations").
		Delete("", "").
		Eq("user_id", userID).
		Eq("provider", "ga4").
		Execute()
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to disconnect integration: %v", err))
		return
	}

	ga4.StoreToken(userID, nil)

	s.respondJSON(w, http.StatusOK, map[string]string{
		"status": "disconnected",
	})
}

func (s *Server) fetchLatestGA4Summary(projectID string) (map[string]interface{}, error) {
	data, _, err := s.serviceRole.
		From("ga4_performance_snapshots").
		Select("*", "", false).
		Eq("project_id", projectID).
		Execute()
	if err != nil {
		// Check if error is due to table not existing (migration not run)
		errStr := err.Error()
		if strings.Contains(errStr, "does not exist") || strings.Contains(errStr, "relation") {
			s.logger.Warn("GA4 performance_snapshots table does not exist - migration may not have been run", zap.Error(err))
			return nil, nil // Return nil instead of error
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

func (s *Server) handleProjectGA4DimensionsDirect(w http.ResponseWriter, r *http.Request, projectID string) {
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
		From("ga4_performance_rows").
		Select("*", "", false).
		Eq("project_id", projectID).
		Eq("row_type", rowType).
		Execute()
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to query GA4 rows: %v", err))
		return
	}

	var rows []map[string]interface{}
	if err := json.Unmarshal(data, &rows); err != nil {
		s.respondError(w, http.StatusInternalServerError, "Failed to parse rows")
		return
	}

	// Sort by sessions descending if available
	sort.Slice(rows, func(i, j int) bool {
		metI, _ := rows[i]["metrics"].(map[string]interface{})
		metJ, _ := rows[j]["metrics"].(map[string]interface{})
		return getFloat(metI["sessions"]) > getFloat(metJ["sessions"])
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

// handleGA4Callback handles GET /api/ga4/callback - OAuth callback
func (s *Server) handleGA4Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	userID, ok := ga4.ConsumeState(state)
	if !ok {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
			<head><title>GA4 Connection Error</title></head>
			<body>
				<h1>Connection Failed</h1>
				<p>Invalid state</p>
			<script>
				(function() {
					if (window.opener && !window.opener.closed) {
						try {
							window.opener.postMessage({type: 'ga4_error', error: 'Invalid state'}, '*');
						} catch (e) {
							console.error('Failed to post error message:', e);
						}
					}
					setTimeout(function() {
						if (window.opener) {
							window.close();
						}
					}, 100);
				})();
			</script>
			</body>
			</html>
		`)
		return
	}

	token, err := ga4.ExchangeCode(code)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
			<head><title>GA4 Connection Error</title></head>
			<body>
				<h1>Connection Failed</h1>
				<p>%v</p>
				<script>
					(function() {
						if (window.opener && !window.opener.closed) {
							try {
								window.opener.postMessage({type: 'ga4_error', error: '%v'}, '*');
							} catch (e) {
								console.error('Failed to post error message:', e);
							}
						}
						setTimeout(function() {
							if (window.opener) {
								window.close();
							}
						}, 100);
					})();
				</script>
			</body>
			</html>
	`, err, err)
		return
	}

	cfg := &ga4IntegrationConfig{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		Expiry:       token.Expiry,
	}

	if scope := token.Extra("scope"); scope != nil {
		switch v := scope.(type) {
		case string:
			cfg.Scope = []string{v}
		case []string:
			cfg.Scope = v
		case []interface{}:
			for _, item := range v {
				if str, ok := item.(string); ok {
					cfg.Scope = append(cfg.Scope, str)
				}
			}
		}
	}

	if err := s.saveGA4Integration(userID, cfg); err != nil {
		s.logger.Error("Failed to persist GA4 token", zap.Error(err))
	}

	ga4.StoreToken(userID, token)

	// Return success page that closes popup and signals parent window
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>GA4 Connected</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					display: flex;
					justify-content: center;
					align-items: center;
					height: 100vh;
					margin: 0;
					background: #f5f5f5;
				}
				.container {
					text-align: center;
					padding: 2rem;
					background: white;
					border-radius: 8px;
					box-shadow: 0 2px 4px rgba(0,0,0,0.1);
				}
				.success {
					color: #4caf50;
					font-size: 2rem;
					margin-bottom: 1rem;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<div class="success">✓</div>
				<h1>Google Analytics 4 Connected!</h1>
				<p>You can close this window.</p>
			</div>
			<script>
				(function() {
					// Immediately signal parent window that connection succeeded
					// Do this synchronously before any potential redirects
					if (window.opener && !window.opener.closed) {
						try {
							window.opener.postMessage({
								type: 'ga4_connected',
								user_id: '%s'
							}, '*');
						} catch (e) {
							console.error('Failed to post message to opener:', e);
						}
					}
					
					// Close popup immediately after posting message
					// Use a small delay to ensure message is sent, but close quickly
					setTimeout(function() {
						try {
							if (window.opener) {
								window.close();
							}
						} catch (e) {
							console.error('Failed to close popup:', e);
						}
					}, 100);
				})();
			</script>
		</body>
		</html>
	`, userID)
}
