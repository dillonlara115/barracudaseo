package api

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	// Default cost per keyword check (in USD)
	DefaultCheckCost = 0.001
)

// trackKeywordUsage records a keyword check in the usage table
func (s *Server) trackKeywordUsage(ctx context.Context, projectID, keywordID, userID, taskID, checkType string, cost float64) error {
	_ = ctx

	usageRecord := map[string]interface{}{
		"id":                 uuid.New().String(),
		"project_id":         projectID,
		"keyword_id":         keywordID,
		"user_id":            userID,
		"check_type":         checkType, // "manual" or "scheduled"
		"dataforseo_task_id": taskID,
		"cost_usd":           cost,
		"checked_at":         time.Now().UTC().Format(time.RFC3339),
	}

	_, _, err := s.serviceRole.From("keyword_usage").Insert(usageRecord, false, "", "", "").Execute()
	if err != nil {
		s.logger.Error("Failed to track keyword usage", zap.Error(err))
		return err
	}

	return nil
}

// getKeywordUsageStats returns usage statistics for a project or user
func (s *Server) getKeywordUsageStats(ctx context.Context, projectID, userID string, startDate, endDate *time.Time) (map[string]interface{}, error) {
	_ = ctx

	query := s.serviceRole.From("keyword_usage").Select("*", "", false)

	if projectID != "" {
		query = query.Eq("project_id", projectID)
	}
	if userID != "" {
		query = query.Eq("user_id", userID)
	}
	if startDate != nil {
		query = query.Gte("checked_at", startDate.Format(time.RFC3339))
	}
	if endDate != nil {
		query = query.Lte("checked_at", endDate.Format(time.RFC3339))
	}

	data, _, err := query.Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch usage stats: %w", err)
	}

	var usageRecords []map[string]interface{}
	if err := json.Unmarshal(data, &usageRecords); err != nil {
		return nil, fmt.Errorf("failed to parse usage stats: %w", err)
	}

	totalChecks := len(usageRecords)
	totalCost := 0.0
	manualChecks := 0
	scheduledChecks := 0

	for _, record := range usageRecords {
		if cost, ok := record["cost_usd"].(float64); ok {
			totalCost += cost
		}
		if checkType, ok := record["check_type"].(string); ok {
			switch checkType {
			case "manual":
				manualChecks++
			case "scheduled":
				scheduledChecks++
			}
		}
	}

	return map[string]interface{}{
		"total_checks":     totalChecks,
		"total_cost_usd":   totalCost,
		"manual_checks":    manualChecks,
		"scheduled_checks": scheduledChecks,
	}, nil
}

// getKeywordLimit returns the maximum number of keywords allowed for a subscription tier
func getKeywordLimit(subscriptionTier string) int {
	switch subscriptionTier {
	case "pro":
		return 500
	case "team":
		return 2000
	default: // free
		return 10
	}
}

// checkKeywordLimit verifies if a user can add more keywords
func (s *Server) checkKeywordLimit(ctx context.Context, projectID, userID string) (bool, int, int, error) {
	_ = ctx

	subscription, err := s.resolveSubscription(userID)
	if err != nil {
		return false, 0, 0, err
	}

	maxKeywords := getKeywordLimit(subscription.EffectiveTier)

	// Get current keyword count
	var keywords []map[string]interface{}
	data, _, err := s.serviceRole.From("keywords").Select("id", "", false).Eq("project_id", projectID).Execute()
	if err == nil {
		json.Unmarshal(data, &keywords)
	}
	currentCount := len(keywords)

	canAdd := currentCount < maxKeywords
	return canAdd, currentCount, maxKeywords, nil
}
