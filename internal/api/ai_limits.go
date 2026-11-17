package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// aiPlanForTier maps a subscription tier (free, pro, team) onto the AI usage
// plan (starter, pro, agency) whose monthly limits apply. See usage.Plans.
func aiPlanForTier(tier string) string {
	// TODO(dillon): decide what each subscription tier gets. This baseline
	// gives free accounts the starter allowance (10 briefs, 3 articles, 50
	// diagnostics a month). Alternatives: no AI allowance on free, or a
	// separate, smaller plan for it.
	switch tier {
	case "pro":
		return "pro"
	case "team":
		return "agency"
	default:
		return "starter"
	}
}

// aiPlanForProject returns the AI usage plan of a project. It follows the
// project owner's subscription, or the account owner's when the owner is a
// team member, so everyone working on a project shares one allowance.
func (s *Server) aiPlanForProject(projectID string) (string, error) {
	data, _, err := s.serviceRole.From("projects").
		Select("owner_id", "", false).
		Eq("id", projectID).
		Limit(1, "").
		Execute()
	if err != nil {
		return "", fmt.Errorf("failed to load project owner: %w", err)
	}
	var projects []struct {
		OwnerID string `json:"owner_id"`
	}
	if err := json.Unmarshal(data, &projects); err != nil {
		return "", fmt.Errorf("failed to parse project owner: %w", err)
	}
	if len(projects) == 0 || projects[0].OwnerID == "" {
		return "", fmt.Errorf("project %s has no owner", projectID)
	}

	subscription, err := s.resolveSubscription(projects[0].OwnerID)
	if err != nil {
		return "", fmt.Errorf("failed to resolve subscription: %w", err)
	}
	return aiPlanForTier(subscription.EffectiveTier), nil
}

// allowAIUsage reports whether the project may run one more AI action of the
// given feature this month. When it returns false the response has been
// written. A failed check refuses the action: letting it through would turn
// any database error into unlimited usage.
func (s *Server) allowAIUsage(w http.ResponseWriter, r *http.Request, ai *AISuiteServices, projectID, feature string) bool {
	planTier, err := s.aiPlanForProject(projectID)
	if err != nil {
		s.logger.Error("Failed to resolve AI plan", zap.Error(err), zap.String("project_id", projectID))
		s.respondError(w, http.StatusServiceUnavailable, "Could not verify your AI usage limit. Please try again.")
		return false
	}

	status, err := ai.UsageTracker.CheckLimit(r.Context(), projectID, planTier, feature)
	if err != nil {
		s.logger.Error("Failed to check AI usage limit", zap.Error(err), zap.String("project_id", projectID), zap.String("feature", feature))
		s.respondError(w, http.StatusServiceUnavailable, "Could not verify your AI usage limit. Please try again.")
		return false
	}
	if !status.Allowed {
		s.respondError(w, http.StatusTooManyRequests, status.Message)
		return false
	}
	return true
}

// pickFields returns the entries of body whose key is in allowed. Update
// handlers use it so a request can only change the fields meant to be edited.
func pickFields(body map[string]interface{}, allowed ...string) map[string]interface{} {
	picked := make(map[string]interface{}, len(allowed))
	for _, key := range allowed {
		if value, ok := body[key]; ok {
			picked[key] = value
		}
	}
	return picked
}

// projectIDOfRecord looks up which project a row of a project-scoped table
// belongs to. found is false when no row has that ID.
func (s *Server) projectIDOfRecord(table, id string) (projectID string, found bool, err error) {
	data, _, err := s.serviceRole.From(table).
		Select("project_id", "", false).
		Eq("id", id).
		Limit(1, "").
		Execute()
	if err != nil {
		return "", false, fmt.Errorf("failed to load %s row: %w", table, err)
	}
	var rows []struct {
		ProjectID string `json:"project_id"`
	}
	if err := json.Unmarshal(data, &rows); err != nil {
		return "", false, fmt.Errorf("failed to parse %s row: %w", table, err)
	}
	if len(rows) == 0 {
		return "", false, nil
	}
	return rows[0].ProjectID, true, nil
}
