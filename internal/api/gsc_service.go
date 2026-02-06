package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dillonlara115/barracudaseo/internal/gsc"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type gscIntegrationConfig struct {
	PropertyURL    string    `json:"property_url"`
	PropertyType   string    `json:"property_type"`
	AccessToken    string    `json:"access_token"`
	RefreshToken   string    `json:"refresh_token"`
	TokenType      string    `json:"token_type"`
	Expiry         time.Time `json:"expiry"`
	Scope          []string  `json:"scope,omitempty"`
	LastSyncPeriod string    `json:"last_sync_period,omitempty"`
}

type gscSyncState struct {
	ProjectID    string                 `json:"project_id"`
	PropertyURL  string                 `json:"property_url"`
	Status       string                 `json:"status"`
	LastSyncedAt *time.Time             `json:"last_synced_at"`
	ErrorLog     map[string]interface{} `json:"error_log"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

func (cfg *gscIntegrationConfig) toOAuthToken() *oauth2.Token {
	token := &oauth2.Token{
		AccessToken:  cfg.AccessToken,
		TokenType:    cfg.TokenType,
		RefreshToken: cfg.RefreshToken,
		Expiry:       cfg.Expiry,
	}
	return token
}

func (cfg *gscIntegrationConfig) mergeMissingFields(existing *gscIntegrationConfig) {
	if existing == nil {
		return
	}
	if cfg.PropertyURL == "" {
		cfg.PropertyURL = existing.PropertyURL
	}
	if cfg.PropertyType == "" {
		cfg.PropertyType = existing.PropertyType
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

func (cfg *gscIntegrationConfig) toMap() map[string]interface{} {
	data := map[string]interface{}{
		"property_url":     cfg.PropertyURL,
		"property_type":    cfg.PropertyType,
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

func parseGSCIntegrationConfig(raw interface{}) (*gscIntegrationConfig, error) {
	if raw == nil {
		return nil, nil
	}
	bytes, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}
	var cfg gscIntegrationConfig
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (s *Server) getGSCIntegration(userID string) (*gscIntegrationConfig, string, error) {
	data, _, err := s.serviceRole.
		From("user_api_integrations").
		Select("*", "", false).
		Eq("user_id", userID).
		Eq("provider", "gsc").
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

	cfg, err := parseGSCIntegrationConfig(rows[0]["config"])
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

func (s *Server) saveGSCIntegration(userID string, cfg *gscIntegrationConfig) error {
	existing, recordID, err := s.getGSCIntegration(userID)
	if err != nil {
		return err
	}

	cfg.mergeMissingFields(existing)

	payload := map[string]interface{}{
		"user_id":  userID,
		"provider": "gsc",
		"config":   cfg.toMap(),
	}

	if recordID == "" {
		_, _, err = s.serviceRole.
			From("user_api_integrations").
			Insert(payload, false, "", "", "").
			Execute()
		return err
	}

	_, _, err = s.serviceRole.
		From("user_api_integrations").
		Update(map[string]interface{}{"config": payload["config"]}, "", "").
		Eq("id", recordID).
		Execute()
	return err
}

func (s *Server) ensureGSCSyncState(projectID, propertyURL string) (*gscSyncState, error) {
	data, _, err := s.serviceRole.
		From("gsc_sync_states").
		Select("*", "", false).
		Eq("project_id", projectID).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to query gsc_sync_states: %w", err)
	}

	var rows []gscSyncState
	if len(data) > 0 {
		if err := json.Unmarshal(data, &rows); err != nil {
			return nil, fmt.Errorf("failed to parse gsc_sync_states: %w", err)
		}
	}

	if len(rows) > 0 {
		state := rows[0]
		if propertyURL != "" && state.PropertyURL != propertyURL {
			_, _, _ = s.serviceRole.
				From("gsc_sync_states").
				Update(map[string]interface{}{"property_url": propertyURL}, "", "").
				Eq("project_id", projectID).
				Execute()
			state.PropertyURL = propertyURL
		}
		return &state, nil
	}

	payload := map[string]interface{}{
		"project_id":   projectID,
		"property_url": propertyURL,
		"status":       "idle",
	}
	_, _, err = s.serviceRole.
		From("gsc_sync_states").
		Insert(payload, false, "", "", "").
		Execute()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &gscSyncState{
		ProjectID:   projectID,
		PropertyURL: propertyURL,
		Status:      "idle",
		UpdatedAt:   now,
	}, nil
}

func (s *Server) updateGSCSyncState(projectID string, status string, lastSynced *time.Time, errPayload interface{}) error {
	update := map[string]interface{}{
		"status": status,
	}
	if lastSynced != nil {
		update["last_synced_at"] = lastSynced.Format(time.RFC3339)
	}
	if errPayload != nil {
		update["error_log"] = errPayload
	} else {
		update["error_log"] = nil
	}

	_, _, err := s.serviceRole.
		From("gsc_sync_states").
		Update(update, "", "").
		Eq("project_id", projectID).
		Execute()
	return err
}

func (s *Server) loadTokenIntoMemory(userID string) (*gscIntegrationConfig, error) {
	cfg, _, err := s.getGSCIntegration(userID)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, fmt.Errorf("no GSC integration configured for user")
	}
	if cfg.RefreshToken == "" && cfg.AccessToken == "" {
		return nil, fmt.Errorf("GSC integration missing OAuth tokens")
	}
	token := cfg.toOAuthToken()
	gsc.StoreToken(userID, token)
	return cfg, nil
}

func (s *Server) syncProjectGSCData(projectID, userID string, cfg *gscIntegrationConfig, lookbackDays int, period string) error {
	if cfg.PropertyURL == "" {
		return fmt.Errorf("GSC property not selected")
	}
	if cfg.RefreshToken == "" && cfg.AccessToken == "" {
		return fmt.Errorf("GSC tokens are not available")
	}

	token := cfg.toOAuthToken()
	gsc.StoreToken(userID, token)

	endDate := time.Now().UTC()
	startDate := endDate.AddDate(0, 0, -lookbackDays)

	report, err := gsc.FetchPerformanceReport(userID, cfg.PropertyURL, startDate, endDate)
	if err != nil {
		return err
	}

	snapshotID := uuid.NewString()
	snapshot := map[string]interface{}{
		"id":           snapshotID,
		"project_id":   projectID,
		"property_url": cfg.PropertyURL,
		"captured_on":  endDate.Format("2006-01-02"),
		"period":       period,
		"totals":       report.Totals,
	}

	if _, _, err := s.serviceRole.
		From("gsc_performance_snapshots").
		Insert(snapshot, false, "", "", "").
		Execute(); err != nil {
		return fmt.Errorf("failed to insert snapshot: %w", err)
	}

	var rows []map[string]interface{}
	addRows := func(rowType string, data []gsc.PerformanceRow) {
		for _, row := range data {
			record := map[string]interface{}{
				"snapshot_id":     snapshotID,
				"project_id":      projectID,
				"row_type":        rowType,
				"dimension_value": row.Value,
				"metrics":         row.Metrics,
				"top_queries":     []interface{}{}, // Always include as empty array, will be set for page rows with query data
			}
			// Only set top_queries for page rows that have query data
			if rowType == "page" {
				if queries, ok := report.PageQueries[row.Value]; ok && len(queries) > 0 {
					record["top_queries"] = queries
				}
			}
			rows = append(rows, record)
		}
	}

	addRows("query", report.Queries)
	addRows("page", report.Pages)
	addRows("country", report.Countries)
	addRows("device", report.Devices)
	addRows("appearance", report.Appearance)
	addRows("date", report.Dates)

	if len(rows) > 0 {
		batchSize := 500
		for i := 0; i < len(rows); i += batchSize {
			end := i + batchSize
			if end > len(rows) {
				end = len(rows)
			}
			_, _, err := s.serviceRole.
				From("gsc_performance_rows").
				Insert(rows[i:end], false, "", "", "").
				Execute()
			if err != nil {
				return fmt.Errorf("failed to insert performance rows: %w", err)
			}
		}
	}

	// Future: fetch coverage/enhancements/insights once APIs are available.

	now := time.Now().UTC()
	if err := s.updateGSCSyncState(projectID, "idle", &now, nil); err != nil {
		return fmt.Errorf("failed to update sync state: %w", err)
	}

	return nil
}
