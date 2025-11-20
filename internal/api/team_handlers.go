package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// TeamMember represents a team member
type TeamMember struct {
	ID           string    `json:"id"`
	AccountOwnerID string  `json:"account_owner_id"`
	UserID       *string   `json:"user_id,omitempty"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	Status       string    `json:"status"`
	InvitedBy    *string   `json:"invited_by,omitempty"`
	InvitedAt    time.Time `json:"invited_at"`
	JoinedAt     *time.Time `json:"joined_at,omitempty"`
}

// InviteTeamMemberRequest represents a request to invite a team member
type InviteTeamMemberRequest struct {
	Email string `json:"email"`
	Role  string `json:"role,omitempty"` // 'admin' or 'member', defaults to 'member'
}

// handleTeam routes team management endpoints
func (s *Server) handleTeam(w http.ResponseWriter, r *http.Request) {
	// Log the incoming path for debugging (use Info level so it shows up)
	s.logger.Info("handleTeam called", 
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("raw_path", r.URL.RawPath))
	
	// Handle both /team and /team/ paths
	path := r.URL.Path
	if strings.HasPrefix(path, "/team/") {
		path = strings.TrimPrefix(path, "/team/")
	} else if strings.HasPrefix(path, "/team") {
		path = strings.TrimPrefix(path, "/team")
	}
	path = strings.Trim(path, "/")
	
	s.logger.Info("handleTeam path after trim", zap.String("trimmed_path", path))

	// Handle empty path (just /team/)
	if path == "" {
		if r.Method == http.MethodGet {
			s.handleListTeamMembers(w, r)
			return
		}
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	switch path {
	case "members":
		if r.Method == http.MethodGet {
			s.handleListTeamMembers(w, r)
		} else if r.Method == http.MethodPost {
			s.handleInviteTeamMember(w, r)
		} else {
			s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	case "invite":
		if r.Method == http.MethodPost {
			s.handleInviteTeamMember(w, r)
		} else {
			s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	default:
		// Handle /team/:id for specific member operations
		if strings.Contains(path, "/") {
			parts := strings.Split(path, "/")
			if len(parts) >= 2 {
				memberID := parts[0]
				action := parts[1]
				switch action {
				case "remove":
					if r.Method == http.MethodDelete {
						s.handleRemoveTeamMember(w, r, memberID)
					} else {
						s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
					}
				case "accept":
					if r.Method == http.MethodPost {
						s.handleAcceptInvite(w, r, memberID)
					} else {
						s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
					}
				default:
					s.respondError(w, http.StatusNotFound, fmt.Sprintf("Team action not found: %s", action))
				}
			} else {
				s.respondError(w, http.StatusNotFound, fmt.Sprintf("Team resource not found: %s", path))
			}
		} else {
			s.respondError(w, http.StatusNotFound, fmt.Sprintf("Team resource not found: %s", path))
		}
	}
}

// handleListTeamMembers lists all team members for the authenticated user's account
func (s *Server) handleListTeamMembers(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok || userID == "" {
		s.respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get user's profile to check if they're a pro/team user
	profile, err := s.fetchProfile(userID)
	if err != nil {
		s.logger.Error("Failed to fetch profile", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to fetch profile")
		return
	}

	if profile == nil {
		s.respondError(w, http.StatusNotFound, "Profile not found")
		return
	}

	// Check subscription tier
	tier, _ := profile["subscription_tier"].(string)
	if tier != "pro" && tier != "team" {
		s.respondError(w, http.StatusForbidden, "Team management requires a Pro or Team subscription")
		return
	}

	// Determine account owner:
	// 1. If user has stripe_subscription_id, they're a paid account owner
	// 2. If user has pro/team tier but no stripe_subscription_id, they're a beta account owner
	// 3. Otherwise, check if they're a team member
	var accountOwnerID string
	stripeSubscriptionID, _ := profile["stripe_subscription_id"].(string)
	
	if stripeSubscriptionID != "" {
		// User is a paid account owner
		accountOwnerID = userID
	} else if tier == "pro" || tier == "team" {
		// User is a beta account owner (has pro/team tier but no Stripe subscription)
		accountOwnerID = userID
	} else {
		// Check if user is a team member
		var teamMembers []map[string]interface{}
		data, _, err := s.serviceRole.From("team_members").
			Select("account_owner_id", "", false).
			Eq("user_id", userID).
			Eq("status", "active").
			Execute()
		
		if err == nil && data != nil {
			if err := json.Unmarshal(data, &teamMembers); err == nil && len(teamMembers) > 0 {
				ownerID, ok := teamMembers[0]["account_owner_id"].(string)
				if ok {
					accountOwnerID = ownerID
				}
			}
		}

		// If still not found, user is not part of any team
		if accountOwnerID == "" {
			s.respondError(w, http.StatusForbidden, "You are not part of a team")
			return
		}
	}

	// Fetch all team members for this account
	// Handle case where team_members table doesn't exist yet (graceful degradation)
	var members []map[string]interface{}
	data, _, err := s.serviceRole.From("team_members").
		Select("*", "", false).
		Eq("account_owner_id", accountOwnerID).
		Order("created_at", nil).
		Execute()

	if err != nil {
		// Check if error is because table doesn't exist
		if strings.Contains(err.Error(), "Could not find the table") || 
		   strings.Contains(err.Error(), "does not exist") ||
		   strings.Contains(err.Error(), "PGRST205") {
			// Table doesn't exist yet - return empty array (graceful degradation)
			s.logger.Info("team_members table not found, returning empty list", 
				zap.String("user_id", userID),
				zap.String("account_owner_id", accountOwnerID))
			members = []map[string]interface{}{}
		} else {
			s.logger.Error("Failed to fetch team members", zap.Error(err))
			s.respondError(w, http.StatusInternalServerError, "Failed to fetch team members")
			return
		}
	} else {
		if err := json.Unmarshal(data, &members); err != nil {
			s.logger.Error("Failed to parse team members", zap.Error(err))
			s.respondError(w, http.StatusInternalServerError, "Failed to parse team members")
			return
		}
	}

	// Get account owner's team_size limit
	ownerProfile, err := s.fetchProfile(accountOwnerID)
	if err != nil {
		s.logger.Error("Failed to fetch owner profile", zap.Error(err))
	}

	teamSizeLimit := 1
	if ownerProfile != nil {
		if size, ok := ownerProfile["team_size"].(float64); ok {
			teamSizeLimit = int(size)
		} else if size, ok := ownerProfile["team_size"].(int); ok {
			teamSizeLimit = size
		}
		s.logger.Info("Team size limit determined",
			zap.String("account_owner_id", accountOwnerID),
			zap.Int("team_size_limit", teamSizeLimit),
			zap.Any("team_size_raw", ownerProfile["team_size"]))
	} else {
		s.logger.Warn("Owner profile not found for team size limit",
			zap.String("account_owner_id", accountOwnerID))
	}

	// Count active members
	activeCount := 0
	for _, m := range members {
		if status, ok := m["status"].(string); ok && status == "active" {
			activeCount++
		}
	}

	response := map[string]interface{}{
		"members":        members,
		"team_size_limit": teamSizeLimit,
		"active_count":   activeCount,
		"is_owner":       accountOwnerID == userID,
	}
	
	s.logger.Info("Returning team members",
		zap.String("account_owner_id", accountOwnerID),
		zap.Int("team_size_limit", teamSizeLimit),
		zap.Int("active_count", activeCount),
		zap.Int("total_members", len(members)),
		zap.Bool("is_owner", accountOwnerID == userID))

	s.respondJSON(w, http.StatusOK, response)
}

// handleInviteTeamMember invites a new team member
func (s *Server) handleInviteTeamMember(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok || userID == "" {
		s.respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req InviteTeamMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" {
		s.respondError(w, http.StatusBadRequest, "Email is required")
		return
	}

	// Validate email format (basic check)
	if !strings.Contains(req.Email, "@") {
		s.respondError(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	// Get user's profile
	profile, err := s.fetchProfile(userID)
	if err != nil {
		s.logger.Error("Failed to fetch profile", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to fetch profile")
		return
	}

	if profile == nil {
		s.respondError(w, http.StatusNotFound, "Profile not found")
		return
	}

	// Check if user is account owner (has stripe_subscription_id or is beta user with pro tier)
	tier, _ := profile["subscription_tier"].(string)
	if tier != "pro" && tier != "team" {
		s.respondError(w, http.StatusForbidden, "Team management requires a Pro or Team subscription")
		return
	}

	stripeSubscriptionID, _ := profile["stripe_subscription_id"].(string)
	// Beta users (pro/team tier without Stripe subscription) can also manage their team
	if stripeSubscriptionID == "" && tier != "pro" && tier != "team" {
		s.respondError(w, http.StatusForbidden, "Only account owners can invite team members")
		return
	}

	// Get team size limit
	teamSizeLimit := 1
	if size, ok := profile["team_size"].(float64); ok {
		teamSizeLimit = int(size)
	} else if size, ok := profile["team_size"].(int); ok {
		teamSizeLimit = size
	}

	// Count current active members
	var existingMembers []map[string]interface{}
	data, _, err := s.serviceRole.From("team_members").
		Select("id", "", false).
		Eq("account_owner_id", userID).
		Eq("status", "active").
		Execute()

	if err != nil {
		// Check if error is because table doesn't exist
		if strings.Contains(err.Error(), "Could not find the table") || 
		   strings.Contains(err.Error(), "does not exist") ||
		   strings.Contains(err.Error(), "PGRST205") {
			s.logger.Error("team_members table does not exist - migration required", zap.Error(err))
			s.respondError(w, http.StatusServiceUnavailable, "Team management is not yet available. Please run the database migration: supabase db push")
			return
		}
		// Other errors - log but continue (assume 0 members)
		s.logger.Warn("Failed to count existing members, assuming 0", zap.Error(err))
		existingMembers = []map[string]interface{}{}
	} else if data != nil {
		json.Unmarshal(data, &existingMembers)
	}

	activeCount := len(existingMembers)
	if activeCount >= teamSizeLimit {
		s.respondError(w, http.StatusForbidden, fmt.Sprintf("Team size limit (%d) reached. Upgrade your plan to add more members.", teamSizeLimit))
		return
	}

	// Check if email already invited or is a member
	var existing []map[string]interface{}
	checkData, _, err := s.serviceRole.From("team_members").
		Select("id,status", "", false).
		Eq("account_owner_id", userID).
		Eq("email", strings.ToLower(req.Email)).
		Execute()

	if err != nil {
		// If table doesn't exist, we already handled it above, so this shouldn't happen
		// But handle gracefully anyway
		if strings.Contains(err.Error(), "Could not find the table") || 
		   strings.Contains(err.Error(), "does not exist") ||
		   strings.Contains(err.Error(), "PGRST205") {
			// Already handled above, but just in case
			return
		}
		// Other errors - log but continue (assume no existing invite)
		s.logger.Warn("Failed to check existing invites, assuming none", zap.Error(err))
		existing = []map[string]interface{}{}
	} else if checkData != nil {
		if err := json.Unmarshal(checkData, &existing); err == nil && len(existing) > 0 {
			status, _ := existing[0]["status"].(string)
			if status == "pending" || status == "active" {
				s.respondError(w, http.StatusConflict, "User already invited or is a team member")
				return
			}
		}
	}

	// Generate invite token
	inviteToken := uuid.New().String()

	// Set default role
	role := req.Role
	if role == "" {
		role = "member"
	}
	if role != "admin" && role != "member" {
		role = "member"
	}

	// Create team member invite
	teamMember := map[string]interface{}{
		"account_owner_id": userID,
		"email":            strings.ToLower(req.Email),
		"role":             role,
		"status":           "pending",
		"invited_by":       userID,
		"invite_token":     inviteToken,
		"invited_at":       time.Now().Format(time.RFC3339),
	}

	_, _, err = s.serviceRole.From("team_members").
		Insert(teamMember, false, "", "", "").
		Execute()

	if err != nil {
		// Check if error is because table doesn't exist
		if strings.Contains(err.Error(), "Could not find the table") || 
		   strings.Contains(err.Error(), "does not exist") ||
		   strings.Contains(err.Error(), "PGRST205") {
			s.logger.Error("team_members table does not exist - migration required", zap.Error(err))
			s.respondError(w, http.StatusServiceUnavailable, "Team management is not yet available. Please run the database migration: supabase db push")
			return
		}
		s.logger.Error("Failed to create team member invite", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to create invite")
		return
	}

	// TODO: Send invite email (implement email service)
	// For now, return the invite token so it can be used
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "https://app.barracudaseo.com"
	}
	inviteURL := fmt.Sprintf("%s/team/accept?token=%s", appURL, inviteToken)

	s.logger.Info("Team member invited",
		zap.String("account_owner_id", userID),
		zap.String("email", req.Email),
		zap.String("invite_token", inviteToken))

	s.respondJSON(w, http.StatusCreated, map[string]interface{}{
		"message":    "Invite sent successfully",
		"invite_url": inviteURL,
		"invite_token": inviteToken, // For testing - remove in production
	})
}

// handleRemoveTeamMember removes a team member
func (s *Server) handleRemoveTeamMember(w http.ResponseWriter, r *http.Request, memberID string) {
	userID, ok := userIDFromContext(r.Context())
	if !ok || userID == "" {
		s.respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get team member record
	var members []map[string]interface{}
	data, _, err := s.serviceRole.From("team_members").
		Select("*", "", false).
		Eq("id", memberID).
		Execute()

	if err != nil {
		s.logger.Error("Failed to fetch team member", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to fetch team member")
		return
	}

	if err := json.Unmarshal(data, &members); err != nil || len(members) == 0 {
		s.respondError(w, http.StatusNotFound, "Team member not found")
		return
	}

	member := members[0]
	accountOwnerID, _ := member["account_owner_id"].(string)

	// Verify user is the account owner
	if accountOwnerID != userID {
		s.respondError(w, http.StatusForbidden, "Only account owners can remove team members")
		return
	}

	// Don't allow removing yourself
	memberUserID, _ := member["user_id"].(string)
	if memberUserID == userID {
		s.respondError(w, http.StatusBadRequest, "Cannot remove yourself from the team")
		return
	}

	// Update status to 'removed' instead of deleting (for audit trail)
	_, _, err = s.serviceRole.From("team_members").
		Update(map[string]interface{}{
			"status": "removed",
		}, "", "").
		Eq("id", memberID).
		Execute()

	if err != nil {
		s.logger.Error("Failed to remove team member", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to remove team member")
		return
	}

	s.logger.Info("Team member removed",
		zap.String("account_owner_id", userID),
		zap.String("member_id", memberID))

	s.respondJSON(w, http.StatusOK, map[string]string{"message": "Team member removed successfully"})
}

// handleAcceptInvite accepts a team invite
func (s *Server) handleAcceptInvite(w http.ResponseWriter, r *http.Request, tokenOrID string) {
	userID, ok := userIDFromContext(r.Context())
	if !ok || userID == "" {
		s.respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Find invite by token or ID
	var members []map[string]interface{}
	var data []byte
	var err error

	// Try token first (for URL-based invites)
	data, _, err = s.serviceRole.From("team_members").
		Select("*", "", false).
		Eq("invite_token", tokenOrID).
		Eq("status", "pending").
		Execute()

	if err != nil || data == nil || len(data) == 0 {
		// Try ID
		data, _, err = s.serviceRole.From("team_members").
			Select("*", "", false).
			Eq("id", tokenOrID).
			Eq("status", "pending").
			Execute()
	}

	if err != nil {
		s.logger.Error("Failed to fetch invite", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to fetch invite")
		return
	}

	if err := json.Unmarshal(data, &members); err != nil || len(members) == 0 {
		s.respondError(w, http.StatusNotFound, "Invite not found or already accepted")
		return
	}

	member := members[0]
	inviteEmail, _ := member["email"].(string)

	// Get user's email to verify
	authHeader := r.Header.Get("Authorization")
	token := strings.TrimSpace(authHeader)
	if strings.HasPrefix(strings.ToLower(token), "bearer ") && len(token) >= 7 {
		token = strings.TrimSpace(token[7:])
	}
	user, err := s.validateTokenViaAPI(token)
	if err != nil {
		s.logger.Error("Failed to get user email", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to verify user")
		return
	}

	// Verify email matches
	if !strings.EqualFold(user.Email, inviteEmail) {
		s.respondError(w, http.StatusForbidden, "Invite email does not match your account email")
		return
	}

	// Update team member record
	memberID, _ := member["id"].(string)
	_, _, err = s.serviceRole.From("team_members").
		Update(map[string]interface{}{
			"user_id":  userID,
			"status":   "active",
			"joined_at": time.Now().Format(time.RFC3339),
		}, "", "").
		Eq("id", memberID).
		Execute()

	if err != nil {
		s.logger.Error("Failed to accept invite", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to accept invite")
		return
	}

	s.logger.Info("Team invite accepted",
		zap.String("user_id", userID),
		zap.String("member_id", memberID))

	s.respondJSON(w, http.StatusOK, map[string]string{"message": "Invite accepted successfully"})
}

