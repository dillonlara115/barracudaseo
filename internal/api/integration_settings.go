package api

import (
	"encoding/json"
	"fmt"
)

// loadProjectSettings loads the settings JSON for a project.
func (s *Server) loadProjectSettings(projectID string) (map[string]interface{}, error) {
	data, _, err := s.serviceRole.
		From("projects").
		Select("settings", "", false).
		Eq("id", projectID).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to load project settings: %w", err)
	}
	if len(data) == 0 {
		return map[string]interface{}{}, nil
	}

	var rows []struct {
		Settings map[string]interface{} `json:"settings"`
	}
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil, fmt.Errorf("failed to parse project settings: %w", err)
	}
	if len(rows) == 0 || rows[0].Settings == nil {
		return map[string]interface{}{}, nil
	}
	return rows[0].Settings, nil
}

// updateProjectSettings merges updates into the project's settings JSON.
func (s *Server) updateProjectSettings(projectID string, updates map[string]interface{}) error {
	settings, err := s.loadProjectSettings(projectID)
	if err != nil {
		return err
	}

	if settings == nil {
		settings = map[string]interface{}{}
	}

	for k, v := range updates {
		if v == nil {
			delete(settings, k)
		} else {
			settings[k] = v
		}
	}

	_, _, err = s.serviceRole.
		From("projects").
		Update(map[string]interface{}{"settings": settings}, "", "").
		Eq("id", projectID).
		Execute()
	if err != nil {
		return fmt.Errorf("failed to update project settings: %w", err)
	}

	return nil
}

// getProjectConnectedSources returns which integrations are configured for a project.
// Used by the insights page to show accurate connection badges (config-based, not just synced data).
func (s *Server) getProjectConnectedSources(projectID string) []string {
	settings, err := s.loadProjectSettings(projectID)
	if err != nil || settings == nil {
		return nil
	}
	var sources []string
	if url, _ := settings["gsc_property_url"].(string); url != "" {
		if uid, _ := settings["gsc_integration_user_id"].(string); uid != "" {
			sources = append(sources, "gsc")
		}
	}
	if id, _ := settings["ga4_property_id"].(string); id != "" {
		if uid, _ := settings["ga4_integration_user_id"].(string); uid != "" {
			sources = append(sources, "ga4")
		}
	}
	if pid, _ := settings["clarity_project_id"].(string); pid != "" {
		if tok, _ := settings["clarity_api_token"].(string); tok != "" {
			sources = append(sources, "clarity")
		}
	}
	return sources
}
