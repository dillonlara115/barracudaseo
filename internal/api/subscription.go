package api

import (
	"fmt"
	"time"
)

// SubscriptionInfo captures the effective subscription context for a user.
type SubscriptionInfo struct {
	Tier           string
	EffectiveTier  string
	Status         string
	IsActive       bool
	PeriodEnd      *time.Time
	AccountOwnerID string
	TeamInfo       *TeamInfo
}

// resolveSubscription determines the effective subscription for a user.
// If the user is a team member, the account owner's subscription is used.
func (s *Server) resolveSubscription(userID string) (*SubscriptionInfo, error) {
	profile, err := s.fetchProfile(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch profile: %w", err)
	}
	if profile == nil {
		return &SubscriptionInfo{Tier: "free", EffectiveTier: "free", IsActive: true, AccountOwnerID: userID}, nil
	}

	teamInfo := s.getTeamInfo(userID, profile)
	accountOwnerID := userID
	if teamInfo != nil && !teamInfo.IsOwner {
		accountOwnerID = teamInfo.AccountOwnerID
		ownerProfile, err := s.fetchProfile(accountOwnerID)
		if err == nil && ownerProfile != nil {
			profile = ownerProfile
		}
	}

	tier := "free"
	if val, ok := profile["subscription_tier"].(string); ok && val != "" {
		tier = val
	}

	status := ""
	if val, ok := profile["subscription_status"].(string); ok {
		status = val
	}

	periodEnd := parseProfileTime(profile["subscription_current_period_end"])
	isActive := isSubscriptionActive(tier, status, periodEnd)

	effectiveTier := tier
	if !isActive {
		effectiveTier = "free"
	}

	return &SubscriptionInfo{
		Tier:           tier,
		EffectiveTier:  effectiveTier,
		Status:         status,
		IsActive:       isActive,
		PeriodEnd:      periodEnd,
		AccountOwnerID: accountOwnerID,
		TeamInfo:       teamInfo,
	}, nil
}

func parseProfileTime(value interface{}) *time.Time {
	switch v := value.(type) {
	case time.Time:
		return &v
	case string:
		if v == "" {
			return nil
		}
		if parsed, err := time.Parse(time.RFC3339, v); err == nil {
			return &parsed
		}
	}
	return nil
}

func isSubscriptionActive(tier, status string, periodEnd *time.Time) bool {
	if tier == "free" {
		return true
	}

	switch status {
	case "", "active", "trialing", "cancelling":
		return true
	}

	if periodEnd != nil && time.Now().Before(*periodEnd) {
		return true
	}

	return false
}

func getMaxPagesLimit(subscriptionTier string) int {
	switch subscriptionTier {
	case "pro":
		return 10000
	case "team":
		return 25000
	default:
		return 100
	}
}
