package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dillonlara115/barracudaseo/internal/dataforseo"
	"go.uber.org/zap"
)

// isConnectionError checks if an error is a connection-related error
func isConnectionError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "no such host") ||
		strings.Contains(errStr, "network is unreachable") ||
		strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "dial tcp")
}

// StartKeywordTaskPoller starts a background goroutine that polls for pending keyword tasks
// This runs automatically in the background, polling every minute
func (s *Server) StartKeywordTaskPoller(ctx context.Context, interval time.Duration) {
	if interval == 0 {
		interval = 1 * time.Minute // Default to 1 minute
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		s.logger.Info("Keyword task poller started", zap.Duration("interval", interval))

		// Poll immediately on startup
		s.pollKeywordTasks(ctx)

		for {
			select {
			case <-ctx.Done():
				s.logger.Info("Keyword task poller stopped")
				return
			case <-ticker.C:
				s.pollKeywordTasks(ctx)
			}
		}
	}()
}

// pollKeywordTasks performs the actual polling logic (extracted from handleKeywordTaskPoll)
func (s *Server) pollKeywordTasks(ctx context.Context) {
	// Get DataForSEO client
	client, err := dataforseo.NewClient()
	if err != nil {
		// DataForSEO not configured - skip polling silently
		return
	}

	// Fetch pending or processing tasks
	data, _, err := s.serviceRole.From("keyword_tasks").
		Select("*", "", false).
		In("status", []string{"pending", "processing"}).
		Order("created_at", nil).
		Limit(50, ""). // Process up to 50 tasks per poll
		Execute()
	if err != nil {
		// If it's a connection error (e.g., Supabase not running), log at debug level
		// This is expected in local development when Supabase isn't started
		if isConnectionError(err) {
			s.logger.Debug("Supabase connection unavailable for keyword task polling (this is normal if Supabase isn't running)",
				zap.Error(err))
		} else {
			s.logger.Error("Failed to fetch tasks for polling", zap.Error(err))
		}
		return
	}

	var tasks []map[string]interface{}
	if err := json.Unmarshal(data, &tasks); err != nil {
		s.logger.Error("Failed to parse tasks for polling", zap.Error(err))
		return
	}

	if len(tasks) == 0 {
		return // No tasks to process
	}

	// Filter tasks that are at least 3 seconds old
	// This gives DataForSEO time to process the task before we check
	// Reduced from 5 seconds since we're polling more frequently now
	threeSecondsAgo := time.Now().UTC().Add(-3 * time.Second)
	var filteredTasks []map[string]interface{}
	for _, task := range tasks {
		runAtStr, ok := task["run_at"].(string)
		if !ok {
			continue
		}
		runAt, err := time.Parse(time.RFC3339, runAtStr)
		if err != nil {
			continue
		}
		if runAt.Before(threeSecondsAgo) || runAt.Equal(threeSecondsAgo) {
			filteredTasks = append(filteredTasks, task)
		}
	}
	tasks = filteredTasks

	if len(tasks) == 0 {
		return // No tasks ready to check yet
	}

	s.logger.Info("Polling keyword tasks", zap.Int("count", len(tasks)))

	// Log stored task IDs for debugging
	storedTaskIDs := make([]string, 0, len(tasks))
	for _, task := range tasks {
		if dataforseoTaskID, ok := task["dataforseo_task_id"].(string); ok {
			storedTaskIDs = append(storedTaskIDs, dataforseoTaskID)
		}
	}
	s.logger.Debug("Stored task IDs in database", zap.Strings("task_ids", storedTaskIDs))

	// Store original tasks list for fallback
	originalTasks := tasks

	// Build a map of our stored task IDs for quick lookup
	ourTaskIDs := make(map[string]bool)
	for _, task := range originalTasks {
		if dataforseoTaskID, ok := task["dataforseo_task_id"].(string); ok {
			ourTaskIDs[dataforseoTaskID] = true
		}
	}

	// First, check which tasks are actually ready in DataForSEO
	// Note: tasks_ready returns ALL ready tasks across all accounts/keywords,
	// so we need to filter to only our tasks
	readyResp, err := client.GetOrganicTasksReady(ctx)
	if err != nil {
		s.logger.Warn("Failed to get ready tasks list from DataForSEO", zap.Error(err))
		// Continue with individual task checks as fallback
	} else if readyResp != nil {
		readyTaskIDList := make([]string, 0, len(readyResp.Tasks))
		ourReadyTaskIDs := make([]string, 0)

		// Filter ready tasks to only those that belong to us
		for _, readyTask := range readyResp.Tasks {
			readyTaskIDList = append(readyTaskIDList, readyTask.ID)
			if ourTaskIDs[readyTask.ID] {
				ourReadyTaskIDs = append(ourReadyTaskIDs, readyTask.ID)
			}
		}

		s.logger.Info("Found ready tasks in DataForSEO",
			zap.Int("total_ready_count", len(readyTaskIDList)),
			zap.Int("our_ready_count", len(ourReadyTaskIDs)),
			zap.Strings("our_ready_task_ids", ourReadyTaskIDs),
			zap.Strings("all_ready_task_ids", readyTaskIDList))

		// Filter our tasks to only those that are ready
		if len(ourReadyTaskIDs) > 0 {
			readyTaskIDsMap := make(map[string]bool)
			for _, id := range ourReadyTaskIDs {
				readyTaskIDsMap[id] = true
			}

			var readyTasks []map[string]interface{}
			for _, task := range originalTasks {
				dataforseoTaskID, ok := task["dataforseo_task_id"].(string)
				if ok && readyTaskIDsMap[dataforseoTaskID] {
					readyTasks = append(readyTasks, task)
					s.logger.Debug("Task matches ready list", zap.String("task_id", dataforseoTaskID))
				}
			}

			if len(readyTasks) > 0 {
				tasks = readyTasks
				s.logger.Info("Filtered to our ready tasks", zap.Int("ready_count", len(tasks)))
			} else {
				s.logger.Info("No tasks matched ready list, falling back to individual checks",
					zap.Int("stored_count", len(storedTaskIDs)),
					zap.Int("our_ready_count", len(ourReadyTaskIDs)))
				tasks = originalTasks
			}
		} else {
			s.logger.Info("No tasks matched ready list, falling back to individual checks",
				zap.Int("stored_count", len(storedTaskIDs)),
				zap.Int("total_ready_from_api", len(readyTaskIDList)))
			// Use original tasks list for fallback
			tasks = originalTasks
		}
	} else {
		s.logger.Debug("No ready tasks response from DataForSEO, will check tasks individually")
	}

	processed := 0
	failed := 0

	for _, task := range tasks {
		taskID, ok := task["id"].(string)
		if !ok {
			continue
		}

		dataforseoTaskID, ok := task["dataforseo_task_id"].(string)
		if !ok {
			continue
		}

		keywordID, ok := task["keyword_id"].(string)
		if !ok {
			continue
		}

		s.logger.Debug("Checking task status",
			zap.String("task_id", taskID),
			zap.String("dataforseo_task_id", dataforseoTaskID),
			zap.String("keyword_id", keywordID))

		// Get task result from DataForSEO
		getResp, err := client.GetOrganicTask(ctx, dataforseoTaskID)
		if err != nil {
			// Check if it's a "Not Found" error (40400)
			// This could mean: task expired, task doesn't exist, or task not ready yet
			// For recently created tasks (< 2 minutes), assume it's not ready yet and retry
			if strings.Contains(err.Error(), "40400") || strings.Contains(err.Error(), "Not Found") {
				// Check how old the task is
				runAtStr, _ := task["run_at"].(string)
				runAt, parseErr := time.Parse(time.RFC3339, runAtStr)
				if parseErr == nil {
					age := time.Since(runAt)
					if age < 2*time.Minute {
						// Task is less than 2 minutes old - probably not ready yet, keep trying
						s.logger.Debug("Task not found but recently created, will retry",
							zap.String("task_id", dataforseoTaskID),
							zap.String("keyword_id", keywordID),
							zap.Duration("age", age))
						continue
					}
				}

				// Task is older than 2 minutes and still not found - likely expired or invalid
				s.logger.Warn("Task not found in DataForSEO (may have expired)",
					zap.String("task_id", dataforseoTaskID),
					zap.String("keyword_id", keywordID))
				// Mark as failed - task expired or invalid
				_, _, _ = s.serviceRole.From("keyword_tasks").
					Update(map[string]interface{}{
						"status": "failed",
						"error":  "Task not found in DataForSEO (may have expired or invalid task ID)",
					}, "", "").
					Eq("id", taskID).
					Execute()
				failed++
				continue
			}
			s.logger.Error("Failed to get task result", zap.String("task_id", dataforseoTaskID), zap.Error(err))
			// For other errors, keep task in processing state to retry later
			continue
		}

		// Log task status
		if len(getResp.Tasks) > 0 {
			s.logger.Debug("Task status from DataForSEO",
				zap.String("task_id", dataforseoTaskID),
				zap.Int("status_code", getResp.Tasks[0].StatusCode),
				zap.String("status_message", getResp.Tasks[0].StatusMessage))
		}

		// Check if task is ready
		if !dataforseo.IsTaskReady(getResp) {
			s.logger.Debug("Task not ready yet, will retry", zap.String("task_id", dataforseoTaskID))
			// Task still processing - update status to processing
			_, _, _ = s.serviceRole.From("keyword_tasks").
				Update(map[string]interface{}{"status": "processing"}, "", "").
				Eq("id", taskID).
				Execute()
			continue
		}

		s.logger.Info("Task is ready, processing snapshot", zap.String("task_id", dataforseoTaskID))

		// Task is ready - extract ranking and create snapshot
		keyword, err := s.fetchKeyword(keywordID)
		if err != nil {
			s.logger.Error("Failed to fetch keyword", zap.String("keyword_id", keywordID), zap.Error(err))
			_, _, _ = s.serviceRole.From("keyword_tasks").
				Update(map[string]interface{}{
					"status": "failed",
					"error":  "Keyword not found",
				}, "", "").
				Eq("id", taskID).
				Execute()
			failed++
			continue
		}

		targetURL := ""
		if keyword.TargetURL != nil {
			targetURL = *keyword.TargetURL
		}

		ranking, err := dataforseo.ExtractRanking(getResp, targetURL)
		if err != nil {
			// Check if error indicates site is not ranking (this is expected, not a failure)
			if strings.Contains(err.Error(), "is not ranking") {
				s.logger.Info("Target URL is not ranking in search results",
					zap.String("keyword_id", keywordID),
					zap.String("target_url", targetURL),
					zap.String("keyword", keyword.Keyword))
				// Mark task as completed with a note that site isn't ranking
				// We don't create a snapshot since there's no position to record
				_, _, _ = s.serviceRole.From("keyword_tasks").
					Update(map[string]interface{}{
						"status":       "completed",
						"completed_at": time.Now().UTC().Format(time.RFC3339),
						"error":        fmt.Sprintf("Site is not ranking: %s", err.Error()),
						"raw_response": getResp,
					}, "", "").
					Eq("id", taskID).
					Execute()
				// Update keyword's last_checked_at even though no snapshot was created
				now := time.Now().UTC().Format(time.RFC3339)
				_, _, _ = s.serviceRole.From("keywords").
					Update(map[string]interface{}{"last_checked_at": now}, "", "").
					Eq("id", keywordID).
					Execute()
				processed++
				continue
			}
			s.logger.Error("Failed to extract ranking",
				zap.String("task_id", taskID),
				zap.String("keyword_id", keywordID),
				zap.Error(err))
			_, _, _ = s.serviceRole.From("keyword_tasks").
				Update(map[string]interface{}{
					"status": "failed",
					"error":  fmt.Sprintf("Failed to extract ranking: %v", err),
				}, "", "").
				Eq("id", taskID).
				Execute()
			failed++
			continue
		}

		s.logger.Debug("Extracted ranking",
			zap.String("keyword_id", keywordID),
			zap.Int("position_organic", ranking.PositionOrganic),
			zap.String("url", ranking.URL))

		// Create snapshot
		snapshot, err := s.createSnapshot(keyword.ProjectID, keywordID, dataforseoTaskID, ranking)
		if err != nil {
			s.logger.Error("Failed to create snapshot", zap.String("keyword_id", keywordID), zap.Error(err))
			_, _, _ = s.serviceRole.From("keyword_tasks").
				Update(map[string]interface{}{
					"status": "failed",
					"error":  fmt.Sprintf("Failed to create snapshot: %v", err),
				}, "", "").
				Eq("id", taskID).
				Execute()
			failed++
			continue
		}
		s.logger.Info("Snapshot created successfully",
			zap.String("keyword_id", keywordID),
			zap.Int("position_organic", ranking.PositionOrganic),
			zap.String("snapshot_id", snapshot.ID))

		// Update task status to completed
		_, _, err = s.serviceRole.From("keyword_tasks").
			Update(map[string]interface{}{
				"status":       "completed",
				"completed_at": time.Now().UTC().Format(time.RFC3339),
				"raw_response": getResp,
			}, "", "").
			Eq("id", taskID).
			Execute()
		if err != nil {
			s.logger.Error("Failed to update task status", zap.String("task_id", taskID), zap.Error(err))
		}

		// Track usage
		checkType := "manual"
		if keyword.CheckFrequency != "manual" {
			checkType = "scheduled"
		}
		if err := s.trackKeywordUsage(ctx, keyword.ProjectID, keywordID, "", dataforseoTaskID, checkType, DefaultCheckCost); err != nil {
			s.logger.Warn("Failed to track keyword usage", zap.Error(err))
		}

		// Update keyword's last_checked_at
		now := time.Now().UTC().Format(time.RFC3339)
		_, _, _ = s.serviceRole.From("keywords").
			Update(map[string]interface{}{"last_checked_at": now}, "", "").
			Eq("id", keywordID).
			Execute()

		processed++
	}

	if processed > 0 || failed > 0 {
		s.logger.Info("Keyword task poll completed",
			zap.Int("processed", processed),
			zap.Int("failed", failed),
			zap.Int("total", len(tasks)))
	}
}

// handleKeywordTaskPoll handles POST /api/internal/keywords/poll
// This endpoint is called by a cron job to poll pending DataForSEO tasks
func (s *Server) handleKeywordTaskPoll(w http.ResponseWriter, r *http.Request) {
	// Verify cron secret
	secret := r.Header.Get("X-Cron-Secret")
	if secret == "" || secret != s.cronSecret {
		s.respondError(w, http.StatusUnauthorized, "Invalid or missing cron secret")
		return
	}

	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get DataForSEO client
	client, err := dataforseo.NewClient()
	if err != nil {
		s.logger.Error("Failed to create DataForSEO client", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "DataForSEO integration not configured")
		return
	}

	// Fetch pending or processing tasks
	data, _, err := s.serviceRole.From("keyword_tasks").
		Select("*", "", false).
		In("status", []string{"pending", "processing"}).
		Order("created_at", nil).
		Limit(50, ""). // Process up to 50 tasks per poll
		Execute()
	if err != nil {
		s.logger.Error("Failed to fetch tasks", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to fetch tasks")
		return
	}

	var tasks []map[string]interface{}
	if err := json.Unmarshal(data, &tasks); err != nil {
		s.logger.Error("Failed to parse tasks", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to parse tasks")
		return
	}

	processed := 0
	failed := 0

	ctx := context.Background()
	for _, task := range tasks {
		taskID, ok := task["id"].(string)
		if !ok {
			continue
		}

		dataforseoTaskID, ok := task["dataforseo_task_id"].(string)
		if !ok {
			continue
		}

		keywordID, ok := task["keyword_id"].(string)
		if !ok {
			continue
		}

		// Get task result from DataForSEO
		getResp, err := client.GetOrganicTask(ctx, dataforseoTaskID)
		if err != nil {
			// Check if it's a "Not Found" error (40400) - task may have expired or never existed
			if strings.Contains(err.Error(), "40400") || strings.Contains(err.Error(), "Not Found") {
				s.logger.Warn("Task not found in DataForSEO (may have expired)",
					zap.String("task_id", dataforseoTaskID),
					zap.String("keyword_id", keywordID))
				// Mark as failed - task expired or invalid
				_, _, _ = s.serviceRole.From("keyword_tasks").
					Update(map[string]interface{}{
						"status": "failed",
						"error":  "Task not found in DataForSEO (may have expired or invalid task ID)",
					}, "", "").
					Eq("id", taskID).
					Execute()
				failed++
				continue
			}
			s.logger.Error("Failed to get task result", zap.String("task_id", dataforseoTaskID), zap.Error(err))
			// For other errors, keep task in processing state to retry later
			continue
		}

		// Check if task is ready
		if !dataforseo.IsTaskReady(getResp) {
			// Task still processing - update status to processing
			_, _, _ = s.serviceRole.From("keyword_tasks").
				Update(map[string]interface{}{"status": "processing"}, "", "").
				Eq("id", taskID).
				Execute()
			continue
		}

		// Task is ready - extract ranking and create snapshot
		// First, get keyword to find target URL
		keyword, err := s.fetchKeyword(keywordID)
		if err != nil {
			s.logger.Error("Failed to fetch keyword", zap.String("keyword_id", keywordID), zap.Error(err))
			_, _, _ = s.serviceRole.From("keyword_tasks").
				Update(map[string]interface{}{
					"status": "failed",
					"error":  "Keyword not found",
				}, "", "").
				Eq("id", taskID).
				Execute()
			failed++
			continue
		}

		targetURL := ""
		if keyword.TargetURL != nil {
			targetURL = *keyword.TargetURL
		}

		ranking, err := dataforseo.ExtractRanking(getResp, targetURL)
		if err != nil {
			// Check if error indicates site is not ranking (this is expected, not a failure)
			if strings.Contains(err.Error(), "is not ranking") {
				s.logger.Info("Target URL is not ranking in search results",
					zap.String("keyword_id", keywordID),
					zap.String("target_url", targetURL),
					zap.String("keyword", keyword.Keyword))
				// Mark task as completed with a note that site isn't ranking
				// We don't create a snapshot since there's no position to record
				_, _, _ = s.serviceRole.From("keyword_tasks").
					Update(map[string]interface{}{
						"status":       "completed",
						"completed_at": time.Now().UTC().Format(time.RFC3339),
						"error":        fmt.Sprintf("Site is not ranking: %s", err.Error()),
						"raw_response": getResp,
					}, "", "").
					Eq("id", taskID).
					Execute()
				// Update keyword's last_checked_at even though no snapshot was created
				now := time.Now().UTC().Format(time.RFC3339)
				_, _, _ = s.serviceRole.From("keywords").
					Update(map[string]interface{}{"last_checked_at": now}, "", "").
					Eq("id", keywordID).
					Execute()
				processed++
				continue
			}
			s.logger.Error("Failed to extract ranking", zap.String("task_id", dataforseoTaskID), zap.Error(err))
			_, _, _ = s.serviceRole.From("keyword_tasks").
				Update(map[string]interface{}{
					"status":       "failed",
					"error":        err.Error(),
					"raw_response": getResp,
				}, "", "").
				Eq("id", taskID).
				Execute()
			failed++
			continue
		}

		// Create snapshot
		_, err = s.createSnapshot(keyword.ProjectID, keywordID, dataforseoTaskID, ranking)
		if err != nil {
			s.logger.Error("Failed to create snapshot", zap.String("keyword_id", keywordID), zap.Error(err))
			_, _, _ = s.serviceRole.From("keyword_tasks").
				Update(map[string]interface{}{
					"status": "failed",
					"error":  fmt.Sprintf("Failed to create snapshot: %v", err),
				}, "", "").
				Eq("id", taskID).
				Execute()
			failed++
			continue
		}

		// Update task status to completed
		_, _, err = s.serviceRole.From("keyword_tasks").
			Update(map[string]interface{}{
				"status":       "completed",
				"completed_at": time.Now().UTC().Format(time.RFC3339),
				"raw_response": getResp,
			}, "", "").
			Eq("id", taskID).
			Execute()
		if err != nil {
			s.logger.Error("Failed to update task status", zap.String("task_id", taskID), zap.Error(err))
		}

		processed++
	}

	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"processed": processed,
		"failed":    failed,
		"total":     len(tasks),
	})
}
