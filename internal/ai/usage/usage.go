package usage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
	"go.uber.org/zap"
)

// PlanLimits defines monthly limits per plan tier
type PlanLimits struct {
	BriefsPerMonth      int
	ArticlesPerMonth    int
	DiagnosticsPerMonth int // -1 means unlimited
}

var Plans = map[string]PlanLimits{
	"starter": {BriefsPerMonth: 10, ArticlesPerMonth: 3, DiagnosticsPerMonth: 50},
	"pro":     {BriefsPerMonth: 75, ArticlesPerMonth: 25, DiagnosticsPerMonth: -1},
	"agency":  {BriefsPerMonth: 400, ArticlesPerMonth: 150, DiagnosticsPerMonth: -1},
}

// Tracker logs AI usage and enforces plan limits
type Tracker struct {
	serviceRole *supabase.Client
	logger      *zap.Logger
}

// NewTracker creates a usage tracker
func NewTracker(serviceRoleClient *supabase.Client, logger *zap.Logger) *Tracker {
	return &Tracker{
		serviceRole: serviceRoleClient,
		logger:      logger,
	}
}

type usageRecord struct {
	ID           string  `json:"id"`
	ProjectID    string  `json:"project_id"`
	UserID       string  `json:"user_id"`
	Feature      string  `json:"feature"`
	Model        string  `json:"model"`
	InputTokens  int     `json:"input_tokens"`
	OutputTokens int     `json:"output_tokens"`
	CostUSD      float64 `json:"cost_usd"`
	CreatedAt    string  `json:"created_at"`
}

// Log records an AI usage event
func (t *Tracker) Log(ctx context.Context, projectID, userID, feature, model string, inputTokens, outputTokens int) error {
	cost := estimateCost(model, inputTokens, outputTokens)

	record := usageRecord{
		ID:           uuid.New().String(),
		ProjectID:    projectID,
		UserID:       userID,
		Feature:      feature,
		Model:        model,
		InputTokens:  inputTokens,
		OutputTokens: outputTokens,
		CostUSD:      cost,
		CreatedAt:    time.Now().UTC().Format(time.RFC3339),
	}

	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal usage record: %w", err)
	}

	_, _, err = t.serviceRole.From("ai_usage_log").Insert(json.RawMessage(data), false, "", "", "").Execute()
	if err != nil {
		return fmt.Errorf("failed to insert usage log: %w", err)
	}

	return nil
}

// LimitStatus represents the result of a limit check
type LimitStatus struct {
	Allowed    bool   `json:"allowed"`
	Used       int    `json:"used"`
	Limit      int    `json:"limit"`
	Percentage int    `json:"percentage"`
	Message    string `json:"message,omitempty"`
}

// CheckLimit verifies whether a project can perform a specific AI action based on plan limits
func (t *Tracker) CheckLimit(ctx context.Context, projectID, feature string) (*LimitStatus, error) {
	planTier, err := t.getProjectPlan(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project plan: %w", err)
	}

	limits, ok := Plans[planTier]
	if !ok {
		limits = Plans["starter"]
	}

	var maxAllowed int
	switch feature {
	case "brief":
		maxAllowed = limits.BriefsPerMonth
	case "article":
		maxAllowed = limits.ArticlesPerMonth
	case "diagnostic", "explain", "diagnose":
		maxAllowed = limits.DiagnosticsPerMonth
	default:
		return &LimitStatus{Allowed: true, Limit: -1}, nil
	}

	if maxAllowed == -1 {
		return &LimitStatus{Allowed: true, Limit: -1}, nil
	}

	usedCount, err := t.getMonthlyUsage(ctx, projectID, feature)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly usage: %w", err)
	}

	pct := 0
	if maxAllowed > 0 {
		pct = (usedCount * 100) / maxAllowed
	}

	status := &LimitStatus{
		Allowed:    usedCount < maxAllowed,
		Used:       usedCount,
		Limit:      maxAllowed,
		Percentage: pct,
	}

	if !status.Allowed {
		status.Message = fmt.Sprintf("You've reached your monthly %s limit (%d/%d). Upgrade your plan for more.", feature, usedCount, maxAllowed)
	} else if pct >= 80 {
		status.Message = fmt.Sprintf("You're approaching your monthly %s limit (%d/%d).", feature, usedCount, maxAllowed)
	}

	return status, nil
}

// GetMonthlyUsageSummary returns usage counts by feature for the current month
func (t *Tracker) GetMonthlyUsageSummary(ctx context.Context, projectID string) (map[string]int, error) {
	startOfMonth := time.Now().UTC().Format("2006-01") + "-01T00:00:00Z"

	var records []usageRecord
	data, _, err := t.serviceRole.From("ai_usage_log").
		Select("feature", "", false).
		Eq("project_id", projectID).
		Gte("created_at", startOfMonth).
		Execute()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("failed to parse usage records: %w", err)
	}

	summary := make(map[string]int)
	for _, r := range records {
		summary[r.Feature]++
	}
	return summary, nil
}

func (t *Tracker) getProjectPlan(ctx context.Context, projectID string) (string, error) {
	var projects []struct {
		PlanTier string `json:"plan_tier"`
	}
	data, _, err := t.serviceRole.From("projects").
		Select("plan_tier", "", false).
		Eq("id", projectID).
		Limit(1, "").
		Execute()
	if err != nil {
		return "starter", err
	}
	if err := json.Unmarshal(data, &projects); err != nil || len(projects) == 0 {
		return "starter", nil
	}
	if projects[0].PlanTier == "" {
		return "starter", nil
	}
	return projects[0].PlanTier, nil
}

func (t *Tracker) getMonthlyUsage(ctx context.Context, projectID, feature string) (int, error) {
	startOfMonth := time.Now().UTC().Format("2006-01") + "-01T00:00:00Z"

	var records []usageRecord
	data, _, err := t.serviceRole.From("ai_usage_log").
		Select("id", "", false).
		Eq("project_id", projectID).
		Eq("feature", feature).
		Gte("created_at", startOfMonth).
		Execute()
	if err != nil {
		return 0, err
	}
	if err := json.Unmarshal(data, &records); err != nil {
		return 0, nil
	}
	return len(records), nil
}

func estimateCost(model string, inputTokens, outputTokens int) float64 {
	switch {
	case model == "gemini-2.5-flash-lite" || model == "gemini-2.0-flash-lite":
		return float64(inputTokens)*0.10/1_000_000 + float64(outputTokens)*0.40/1_000_000
	default:
		// gemini-2.5-flash and similar
		return float64(inputTokens)*0.30/1_000_000 + float64(outputTokens)*2.50/1_000_000
	}
}
