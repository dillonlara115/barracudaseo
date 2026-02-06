package api

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dillonlara115/barracudaseo/internal/ga4"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type ga4IntegrationConfig struct {
	PropertyID     string    `json:"property_id"`
	PropertyName   string    `json:"property_name"`
	AccessToken    string    `json:"access_token"`
	RefreshToken   string    `json:"refresh_token"`
	TokenType      string    `json:"token_type"`
	Expiry         time.Time `json:"expiry"`
	Scope          []string  `json:"scope,omitempty"`
	LastSyncPeriod string    `json:"last_sync_period,omitempty"`
}

type ga4SyncState struct {
	ProjectID    string                 `json:"project_id"`
	PropertyID   string                 `json:"property_id"`
	Status       string                 `json:"status"`
	LastSyncedAt *time.Time             `json:"last_synced_at"`
	ErrorLog     map[string]interface{} `json:"error_log"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

func (cfg *ga4IntegrationConfig) toOAuthToken() *oauth2.Token {
	token := &oauth2.Token{
		AccessToken:  cfg.AccessToken,
		TokenType:    cfg.TokenType,
		RefreshToken: cfg.RefreshToken,
		Expiry:       cfg.Expiry,
	}
	return token
}

func (cfg *ga4IntegrationConfig) mergeMissingFields(existing *ga4IntegrationConfig) {
	if existing == nil {
		return
	}
	if cfg.PropertyID == "" {
		cfg.PropertyID = existing.PropertyID
	}
	if cfg.PropertyName == "" {
		cfg.PropertyName = existing.PropertyName
	}
	if cfg.RefreshToken == "" {
		cfg.RefreshToken = existing.RefreshToken
	}
	if cfg.Scope == nil || len(cfg.Scope) == 0 {
		cfg.Scope = existing.Scope
	}
	if cfg.LastSyncPeriod == "" {
		cfg.LastSyncPeriod = existing.LastSyncPeriod
	}
	if cfg.Expiry.IsZero() {
		cfg.Expiry = existing.Expiry
	}
	if cfg.AccessToken == "" {
		cfg.AccessToken = existing.AccessToken
	}
	if cfg.TokenType == "" {
		cfg.TokenType = existing.TokenType
	}
}

func (cfg *ga4IntegrationConfig) toMap() map[string]interface{} {
	data := map[string]interface{}{
		"property_id":      cfg.PropertyID,
		"property_name":    cfg.PropertyName,
		"access_token":     cfg.AccessToken,
		"refresh_token":    cfg.RefreshToken,
		"token_type":       cfg.TokenType,
		"last_sync_period": cfg.LastSyncPeriod,
	}
	if !cfg.Expiry.IsZero() {
		data["expiry"] = cfg.Expiry.Format(time.RFC3339)
	}
	if len(cfg.Scope) > 0 {
		data["scope"] = cfg.Scope
	}
	return data
}

func parseGA4IntegrationConfig(raw interface{}) (*ga4IntegrationConfig, error) {
	if raw == nil {
		return nil, nil
	}
	bytes, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}
	var cfg ga4IntegrationConfig
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (s *Server) getGA4Integration(userID string) (*ga4IntegrationConfig, string, error) {
	data, _, err := s.serviceRole.
		From("user_api_integrations").
		Select("*", "", false).
		Eq("user_id", userID).
		Eq("provider", "ga4").
		Execute()

	if err != nil {
		return nil, "", fmt.Errorf("failed to query user_api_integrations: %w", err)
	}

	var rows []map[string]interface{}
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil, "", fmt.Errorf("failed to parse user_api_integrations data: %w", err)
	}

	if len(rows) == 0 {
		return nil, "", nil
	}

	cfg, err := parseGA4IntegrationConfig(rows[0]["config"])
	if err != nil {
		return nil, "", err
	}

	// Parse expiry if provided as string
	if cfg != nil && cfg.Expiry.IsZero() {
		if configMap, ok := rows[0]["config"].(map[string]interface{}); ok {
			if expiryStr, ok := configMap["expiry"].(string); ok && expiryStr != "" {
				if ts, err := time.Parse(time.RFC3339, expiryStr); err == nil {
					cfg.Expiry = ts
				}
			}
		}
	}

	id, _ := rows[0]["id"].(string)
	return cfg, id, nil
}

func (s *Server) saveGA4Integration(userID string, cfg *ga4IntegrationConfig) error {
	existing, recordID, err := s.getGA4Integration(userID)
	if err != nil {
		return err
	}

	cfg.mergeMissingFields(existing)

	configData := cfg.toMap()

	if recordID != "" {
		// Update existing
		_, _, err = s.serviceRole.From("user_api_integrations").
			Update(map[string]interface{}{
				"config":     configData,
				"updated_at": time.Now().UTC().Format(time.RFC3339),
			}, "", "").
			Eq("id", recordID).
			Execute()
		return err
	}

	// Create new
	_, _, err = s.serviceRole.From("user_api_integrations").
		Insert(map[string]interface{}{
			"id":         uuid.New().String(),
			"user_id":    userID,
			"provider":   "ga4",
			"config":     configData,
			"created_at": time.Now().UTC().Format(time.RFC3339),
			"updated_at": time.Now().UTC().Format(time.RFC3339),
		}, false, "", "", "").
		Execute()
	return err
}

func (s *Server) loadGA4TokenIntoMemory(userID string) (*ga4IntegrationConfig, error) {
	cfg, _, err := s.getGA4Integration(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to load GA4 integration: %w", err)
	}
	if cfg == nil {
		return nil, fmt.Errorf("GA4 integration not found - connect first")
	}

	// Load token into memory
	token := cfg.toOAuthToken()
	ga4.StoreToken(userID, token)

	return cfg, nil
}

func (s *Server) ensureGA4SyncState(projectID string, propertyID string) (*ga4SyncState, error) {
	data, _, err := s.serviceRole.
		From("ga4_sync_states").
		Select("*", "", false).
		Eq("project_id", projectID).
		Execute()

	if err != nil {
		// Check if error is due to table not existing (migration not run)
		errStr := err.Error()
		if strings.Contains(errStr, "does not exist") || strings.Contains(errStr, "relation") {
			s.logger.Warn("GA4 sync_states table does not exist - migration may not have been run", zap.Error(err))
			// Return nil state instead of error - component can handle this gracefully
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query ga4_sync_states: %w", err)
	}

	var rows []map[string]interface{}
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil, fmt.Errorf("failed to parse ga4_sync_states data: %w", err)
	}

	var state *ga4SyncState
	if len(rows) == 0 {
		// Create new state
		now := time.Now().UTC()
		state = &ga4SyncState{
			ProjectID:  projectID,
			PropertyID: propertyID,
			Status:     "idle",
			UpdatedAt:  now,
		}

		_, _, err = s.serviceRole.From("ga4_sync_states").
			Insert(map[string]interface{}{
				"project_id":  projectID,
				"property_id": propertyID,
				"status":      "idle",
				"created_at":  now.Format(time.RFC3339),
				"updated_at":  now.Format(time.RFC3339),
			}, false, "", "", "").
			Execute()
		if err != nil {
			// Check if error is due to table not existing (migration not run)
			errStr := err.Error()
			if strings.Contains(errStr, "does not exist") || strings.Contains(errStr, "relation") {
				s.logger.Warn("GA4 sync_states table does not exist - migration may not have been run", zap.Error(err))
				// Return nil state instead of error - component can handle this gracefully
				return nil, nil
			}
			return nil, fmt.Errorf("failed to create ga4_sync_state: %w", err)
		}
	} else {
		// Parse existing state
		row := rows[0]
		state = &ga4SyncState{
			ProjectID: projectID,
		}

		if propID, ok := row["property_id"].(string); ok {
			state.PropertyID = propID
		}
		if status, ok := row["status"].(string); ok {
			state.Status = status
		}
		if updatedAt, ok := row["updated_at"].(string); ok {
			if ts, err := time.Parse(time.RFC3339, updatedAt); err == nil {
				state.UpdatedAt = ts
			}
		}
		if lastSyncedAt, ok := row["last_synced_at"].(string); ok && lastSyncedAt != "" {
			if ts, err := time.Parse(time.RFC3339, lastSyncedAt); err == nil {
				state.LastSyncedAt = &ts
			}
		}
		if errorLog, ok := row["error_log"].(map[string]interface{}); ok {
			state.ErrorLog = errorLog
		}

		// Update property_id if provided and different
		if propertyID != "" && state.PropertyID != propertyID {
			_, _, err = s.serviceRole.From("ga4_sync_states").
				Update(map[string]interface{}{
					"property_id": propertyID,
					"updated_at":  time.Now().UTC().Format(time.RFC3339),
				}, "", "").
				Eq("project_id", projectID).
				Execute()
			if err != nil {
				// Log but don't fail - table might not exist
				errStr := err.Error()
				if strings.Contains(errStr, "does not exist") || strings.Contains(errStr, "relation") {
					s.logger.Warn("GA4 sync_states table does not exist during update", zap.Error(err))
				}
			} else {
				state.PropertyID = propertyID
			}
		}
	}

	return state, nil
}

func (s *Server) updateGA4SyncState(projectID string, status string, lastSyncedAt *time.Time, errorLog map[string]interface{}) error {
	updateData := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now().UTC().Format(time.RFC3339),
	}

	if lastSyncedAt != nil {
		updateData["last_synced_at"] = lastSyncedAt.Format(time.RFC3339)
	}
	if errorLog != nil {
		updateData["error_log"] = errorLog
	}

	_, _, err := s.serviceRole.From("ga4_sync_states").
		Update(updateData, "", "").
		Eq("project_id", projectID).
		Execute()

	return err
}

func (s *Server) syncProjectGA4Data(projectID, userID string, cfg *ga4IntegrationConfig, lookbackDays int, period string) error {
	if cfg.PropertyID == "" {
		return fmt.Errorf("property_id not set")
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -lookbackDays)

	// Fetch performance data
	performanceMap, err := ga4.FetchPerformanceData(userID, cfg.PropertyID, startDate, endDate)
	if err != nil {
		return fmt.Errorf("failed to fetch GA4 data: %w", err)
	}

	// Create snapshot
	snapshotID := uuid.New().String()
	now := time.Now().UTC()
	capturedOn := now.Format("2006-01-02")

	totals := map[string]interface{}{
		"total_sessions":  int64(0),
		"total_users":     int64(0),
		"total_pageviews": int64(0),
	}

	for _, perf := range performanceMap {
		totals["total_sessions"] = totals["total_sessions"].(int64) + perf.Sessions
		totals["total_users"] = totals["total_users"].(int64) + perf.Users
		totals["total_pageviews"] = totals["total_pageviews"].(int64) + perf.PageViews
	}

	// Insert snapshot
	_, _, err = s.serviceRole.From("ga4_performance_snapshots").
		Insert(map[string]interface{}{
			"id":          snapshotID,
			"project_id":  projectID,
			"property_id": cfg.PropertyID,
			"captured_on": capturedOn,
			"period":      period,
			"totals":      totals,
			"created_at":  now.Format(time.RFC3339),
		}, false, "", "", "").
		Execute()
	if err != nil {
		return fmt.Errorf("failed to create snapshot: %w", err)
	}

	// Insert performance rows
	for url, perf := range performanceMap {
		metrics := map[string]interface{}{
			"sessions":             perf.Sessions,
			"users":                perf.Users,
			"page_views":           perf.PageViews,
			"bounce_rate":          perf.BounceRate,
			"avg_session_duration": perf.AvgSessionDuration,
			"conversions":          perf.Conversions,
			"revenue":              perf.Revenue,
		}

		_, _, err = s.serviceRole.From("ga4_performance_rows").
			Insert(map[string]interface{}{
				"snapshot_id":     snapshotID,
				"project_id":      projectID,
				"row_type":        "page",
				"dimension_value": url,
				"metrics":         metrics,
				"created_at":      now.Format(time.RFC3339),
			}, false, "", "", "").
			Execute()
		if err != nil {
			s.logger.Warn("Failed to insert GA4 performance row", zap.String("url", url), zap.Error(err))
			// Continue with other rows
		}
	}

	// Update last sync period
	cfg.LastSyncPeriod = period
	return s.saveGA4Integration(userID, cfg)
}
