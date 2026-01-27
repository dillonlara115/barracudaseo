package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dillonlara115/barracudaseo/internal/dataforseo"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// handleScheduledKeywordChecks handles POST /api/internal/keywords/check-scheduled
// This endpoint is called by a cron job to check keywords that are due for scheduled checks
func (s *Server) handleScheduledKeywordChecks(w http.ResponseWriter, r *http.Request) {
	// Verify cron secret
	secret := r.Header.Get("X-Cron-Secret")
	if secret != s.cronSecret || secret == "" {
		s.logger.Warn("Unauthorized scheduled keyword check attempt")
		s.respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx := r.Context()
	now := time.Now().UTC()

	// Find keywords that need checking (next_check_at <= now and check_frequency != 'manual')
	var keywords []map[string]interface{}
	data, _, err := s.serviceRole.From("keywords").
		Select("*", "", false).
		In("check_frequency", []string{"daily", "weekly"}).
		Lte("next_check_at", now.Format(time.RFC3339)).
		Limit(50, ""). // Process up to 50 keywords per run
		Execute()

	if err != nil {
		s.logger.Error("Failed to fetch keywords for scheduled check", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to fetch keywords")
		return
	}

	if err := json.Unmarshal(data, &keywords); err != nil {
		s.logger.Error("Failed to parse keywords data", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to parse keywords")
		return
	}

	if len(keywords) == 0 {
		s.respondJSON(w, http.StatusOK, map[string]interface{}{
			"message": "No keywords due for checking",
			"checked": 0,
		})
		return
	}

	s.logger.Info("Processing scheduled keyword checks", zap.Int("count", len(keywords)))

	// Get DataForSEO client
	client, err := dataforseo.NewClient()
	if err != nil {
		s.logger.Error("Failed to create DataForSEO client", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "DataForSEO integration not configured")
		return
	}

	checkedCount := 0
	errorCount := 0

	for _, k := range keywords {
		keywordID := getKeywordString(k, "id")
		projectID := getKeywordString(k, "project_id")

		// Get project owner for usage tracking
		var projects []map[string]interface{}
		projectData, _, err := s.serviceRole.From("projects").Select("owner_id", "", false).Eq("id", projectID).Execute()
		if err != nil {
			s.logger.Error("Failed to fetch project owner", zap.String("project_id", projectID), zap.Error(err))
			errorCount++
			continue
		}
		if err := json.Unmarshal(projectData, &projects); err != nil || len(projects) == 0 {
			s.logger.Error("Failed to parse project data", zap.String("project_id", projectID))
			errorCount++
			continue
		}
		ownerID := getKeywordString(projects[0], "owner_id")

		// Create task
		task := dataforseo.OrganicTaskPost{
			LanguageName: getKeywordString(k, "language_name"),
			LocationName: getKeywordString(k, "location_name"),
			Keyword:      getKeywordString(k, "keyword"),
			Device:       getKeywordString(k, "device"),
			SearchEngine: getKeywordString(k, "search_engine"),
		}

		taskResp, err := client.CreateOrganicTask(ctx, task)
		if err != nil {
			s.logger.Error("Failed to create DataForSEO task for scheduled check", zap.String("keyword_id", keywordID), zap.Error(err))
			errorCount++
			continue
		}

		if len(taskResp.Tasks) == 0 {
			s.logger.Error("No task ID returned from DataForSEO", zap.String("keyword_id", keywordID))
			errorCount++
			continue
		}

		taskID := taskResp.Tasks[0].ID

		// Create task record
		taskRecordID := uuid.New().String()
		taskRecord := map[string]interface{}{
			"id":                 taskRecordID,
			"project_id":         projectID,
			"keyword_id":         keywordID,
			"dataforseo_task_id": taskID,
			"status":             "pending",
			"raw_request":        map[string]interface{}{"task": task},
		}

		_, _, err = s.serviceRole.From("keyword_tasks").Insert(taskRecord, false, "", "", "").Execute()
		if err != nil {
			s.logger.Error("Failed to insert task record", zap.String("keyword_id", keywordID), zap.Error(err))
			// Continue anyway
		}

		// Try to get result immediately
		getResp, err := client.GetOrganicTask(ctx, taskID)
		if err == nil && dataforseo.IsTaskReady(getResp) {
			targetURL := getKeywordString(k, "target_url")
			ranking, err := dataforseo.ExtractRanking(getResp, targetURL)
			if err == nil {
				// Create snapshot
				_, err := s.createSnapshot(projectID, keywordID, taskID, ranking)
				if err == nil {
					// Update task status
					_, _, _ = s.serviceRole.From("keyword_tasks").
						Update(map[string]interface{}{
							"status":       "completed",
							"completed_at": time.Now().UTC().Format(time.RFC3339),
							"raw_response": getResp,
						}, "", "").
						Eq("id", taskRecordID).
						Execute()

					// Track usage
					if err := s.trackKeywordUsage(ctx, projectID, keywordID, ownerID, taskID, "scheduled", DefaultCheckCost); err != nil {
						s.logger.Warn("Failed to track keyword usage", zap.String("keyword_id", keywordID), zap.Error(err))
					}

					// Update keyword's last_checked_at (trigger will update next_check_at)
					nowStr := time.Now().UTC().Format(time.RFC3339)
					_, _, _ = s.serviceRole.From("keywords").
						Update(map[string]interface{}{"last_checked_at": nowStr}, "", "").
						Eq("id", keywordID).
						Execute()

					checkedCount++
				} else {
					s.logger.Error("Failed to create snapshot", zap.String("keyword_id", keywordID), zap.Error(err))
					errorCount++
				}
			} else {
				// Check if error indicates site is not ranking (this is expected, not an error)
				if strings.Contains(err.Error(), "is not ranking") {
					s.logger.Info("Target URL is not ranking in search results",
						zap.String("keyword_id", keywordID),
						zap.String("target_url", targetURL))
					// Mark task as completed with a note that site isn't ranking
					_, _, _ = s.serviceRole.From("keyword_tasks").
						Update(map[string]interface{}{
							"status":       "completed",
							"completed_at": time.Now().UTC().Format(time.RFC3339),
							"error":         fmt.Sprintf("Site is not ranking: %s", err.Error()),
							"raw_response":  getResp,
						}, "", "").
						Eq("id", taskRecordID).
						Execute()
					// Update keyword's last_checked_at even though no snapshot was created
					nowStr := time.Now().UTC().Format(time.RFC3339)
					_, _, _ = s.serviceRole.From("keywords").
						Update(map[string]interface{}{"last_checked_at": nowStr}, "", "").
						Eq("id", keywordID).
						Execute()
					checkedCount++ // Count as checked, even though not ranking
				} else {
					s.logger.Warn("Failed to extract ranking", zap.String("keyword_id", keywordID), zap.Error(err))
					errorCount++
				}
			}
		} else {
			// Task not ready - update status for background polling
			_, _, _ = s.serviceRole.From("keyword_tasks").
				Update(map[string]interface{}{"status": "processing"}, "", "").
				Eq("id", taskRecordID).
				Execute()
			// Still count as checked since task was created
			checkedCount++
		}
	}

	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":     "Scheduled keyword checks completed",
		"checked":     checkedCount,
		"errors":      errorCount,
		"total_found": len(keywords),
	})
}
