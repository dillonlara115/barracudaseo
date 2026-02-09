package api

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dillonlara115/barracudaseo/internal/clarity"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type clarityIntegrationConfig struct {
	ClarityProjectID string `json:"clarity_project_id"`
	APIToken         string `json:"api_token"`
}

func (cfg *clarityIntegrationConfig) toMap() map[string]interface{} {
	return map[string]interface{}{
		"clarity_project_id": cfg.ClarityProjectID,
		"api_token":          cfg.APIToken,
	}
}

func parseClarityIntegrationConfig(raw interface{}) (*clarityIntegrationConfig, error) {
	if raw == nil {
		return nil, nil
	}
	bytes, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}
	var cfg clarityIntegrationConfig
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (s *Server) getClarityIntegration(userID string) (*clarityIntegrationConfig, string, error) {
	data, _, err := s.serviceRole.
		From("user_api_integrations").
		Select("*", "", false).
		Eq("user_id", userID).
		Eq("provider", "clarity").
		Execute()

	if err != nil {
		return nil, "", fmt.Errorf("failed to query clarity integration: %w", err)
	}

	var rows []map[string]interface{}
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil, "", fmt.Errorf("failed to parse clarity integration data: %w", err)
	}

	if len(rows) == 0 {
		return nil, "", nil
	}

	cfg, err := parseClarityIntegrationConfig(rows[0]["config"])
	if err != nil {
		return nil, "", err
	}

	id, _ := rows[0]["id"].(string)
	return cfg, id, nil
}

func (s *Server) saveClarityIntegration(userID string, cfg *clarityIntegrationConfig) error {
	_, recordID, err := s.getClarityIntegration(userID)
	if err != nil {
		return err
	}

	configData := cfg.toMap()

	if recordID != "" {
		_, _, err = s.serviceRole.From("user_api_integrations").
			Update(map[string]interface{}{
				"config":     configData,
				"updated_at": time.Now().UTC().Format(time.RFC3339),
			}, "", "").
			Eq("id", recordID).
			Execute()
		return err
	}

	_, _, err = s.serviceRole.From("user_api_integrations").
		Insert(map[string]interface{}{
			"id":         uuid.New().String(),
			"user_id":    userID,
			"provider":   "clarity",
			"config":     configData,
			"created_at": time.Now().UTC().Format(time.RFC3339),
			"updated_at": time.Now().UTC().Format(time.RFC3339),
		}, false, "", "", "").
		Execute()
	return err
}

func (s *Server) ensureClaritySyncState(projectID, clarityProjectID string) (map[string]interface{}, error) {
	data, _, err := s.serviceRole.
		From("clarity_sync_states").
		Select("*", "", false).
		Eq("project_id", projectID).
		Execute()

	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "does not exist") || strings.Contains(errStr, "relation") {
			s.logger.Warn("clarity_sync_states table does not exist", zap.Error(err))
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query clarity_sync_states: %w", err)
	}

	var rows []map[string]interface{}
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil, fmt.Errorf("failed to parse clarity_sync_states data: %w", err)
	}

	if len(rows) == 0 {
		now := time.Now().UTC()
		record := map[string]interface{}{
			"project_id":         projectID,
			"clarity_project_id": clarityProjectID,
			"status":             "idle",
			"created_at":         now.Format(time.RFC3339),
			"updated_at":         now.Format(time.RFC3339),
		}

		_, _, err = s.serviceRole.From("clarity_sync_states").
			Insert(record, false, "", "", "").
			Execute()
		if err != nil {
			errStr := err.Error()
			if strings.Contains(errStr, "does not exist") || strings.Contains(errStr, "relation") {
				return nil, nil
			}
			return nil, fmt.Errorf("failed to create clarity sync state: %w", err)
		}

		return record, nil
	}

	row := rows[0]
	if clarityProjectID != "" {
		if current, _ := row["clarity_project_id"].(string); current != clarityProjectID {
			_, _, _ = s.serviceRole.From("clarity_sync_states").
				Update(map[string]interface{}{
					"clarity_project_id": clarityProjectID,
					"updated_at":         time.Now().UTC().Format(time.RFC3339),
				}, "", "").
				Eq("project_id", projectID).
				Execute()
			row["clarity_project_id"] = clarityProjectID
		}
	}

	return row, nil
}

func (s *Server) updateClaritySyncState(projectID, status string, lastSyncedAt *time.Time, errorLog map[string]interface{}) error {
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

	_, _, err := s.serviceRole.From("clarity_sync_states").
		Update(updateData, "", "").
		Eq("project_id", projectID).
		Execute()

	return err
}

func (s *Server) syncProjectClarityData(projectID, userID string, cfg *clarityIntegrationConfig, numDays int, period string) error {
	if cfg.ClarityProjectID == "" || cfg.APIToken == "" {
		return fmt.Errorf("clarity project ID and API token are required")
	}

	snapshotID := uuid.New().String()
	now := time.Now().UTC()
	capturedOn := now.Format("2006-01-02")

	// Fetch 3 dimensions: url, device, source (uses 3 of 10 daily API calls)
	dimensions := map[string]string{
		"url":    "Url",
		"device": "Device",
		"source": "Source",
	}

	var overallSummary *clarity.InsightMetrics
	batchSize := 500

	for rowType, apiDimension := range dimensions {
		rows, summary, err := clarity.FetchInsights(cfg.APIToken, cfg.ClarityProjectID, numDays, apiDimension)
		if err != nil {
			s.logger.Warn("Failed to fetch Clarity dimension", zap.String("dimension", apiDimension), zap.Error(err))
			continue
		}

		if overallSummary == nil && summary != nil {
			overallSummary = summary
		}

		var records []map[string]interface{}
		for _, row := range rows {
			records = append(records, map[string]interface{}{
				"snapshot_id":     snapshotID,
				"project_id":      projectID,
				"row_type":        rowType,
				"dimension_value": row.DimensionValue,
				"metrics":         clarity.MetricsToMap(row.Metrics),
				"created_at":      now.Format(time.RFC3339),
			})
		}

		for i := 0; i < len(records); i += batchSize {
			end := i + batchSize
			if end > len(records) {
				end = len(records)
			}
			_, _, err = s.serviceRole.From("clarity_performance_rows").
				Insert(records[i:end], false, "", "minimal", "").
				Execute()
			if err != nil {
				s.logger.Warn("Failed to insert Clarity rows", zap.String("row_type", rowType), zap.Error(err))
			}
		}
	}

	// Build totals from summary
	totals := map[string]interface{}{}
	if overallSummary != nil {
		totals = clarity.MetricsToMap(*overallSummary)
	}

	// Insert snapshot
	_, _, err := s.serviceRole.From("clarity_performance_snapshots").
		Insert(map[string]interface{}{
			"id":                 snapshotID,
			"project_id":         projectID,
			"clarity_project_id": cfg.ClarityProjectID,
			"captured_on":        capturedOn,
			"period":             period,
			"totals":             totals,
			"created_at":         now.Format(time.RFC3339),
		}, false, "", "", "").
		Execute()
	if err != nil {
		return fmt.Errorf("failed to create clarity snapshot: %w", err)
	}

	return nil
}
