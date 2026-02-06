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
