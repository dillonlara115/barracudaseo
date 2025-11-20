package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
			s.logger.Info("Parsing team member action", 
				zap.String("path", path),
				zap.Strings("parts", parts),
				zap.Int("parts_count", len(parts)))
			if len(parts) >= 2 {
				memberID := parts[0]
				action := parts[1]
				s.logger.Info("Team member action", 
					zap.String("member_id", memberID),
					zap.String("action", action),
					zap.String("method", r.Method))
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
				case "resend":
					if r.Method == http.MethodPost {
						s.handleResendInvite(w, r, memberID)
					} else {
						s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
					}
				case "details":
					if r.Method == http.MethodGet {
						s.handleGetInviteDetails(w, r, memberID)
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

	// Build invite URL
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		// Default to localhost for local development, production URL otherwise
		// Check if running locally (common indicators)
		if os.Getenv("PORT") == "" || os.Getenv("PORT") == "8080" {
			appURL = "http://localhost:5173"
		} else {
			appURL = "https://app.barracudaseo.com"
		}
	}
	inviteURL := fmt.Sprintf("%s/#/team/accept?token=%s", appURL, inviteToken)

	// Check if user exists, create if not (hybrid approach)
	inviteeUserID, userCreated, err := s.ensureUserExists(strings.ToLower(req.Email))
	if err != nil {
		s.logger.Warn("Failed to ensure user exists, but invite created", zap.Error(err), zap.String("email", req.Email))
		inviteeUserID = "" // Will be set when user accepts invite
	}

	// Send invite email via configured email service
	if err := s.emailService.SendTeamInvite(strings.ToLower(req.Email), inviteURL, userCreated); err != nil {
		s.logger.Warn("Failed to send invite email, but invite created", zap.Error(err), zap.String("email", req.Email))
		// Don't fail the request - invite URL is still returned
	}

	s.logger.Info("Team member invited",
		zap.String("account_owner_id", userID),
		zap.String("email", req.Email),
		zap.String("invite_token", inviteToken),
		zap.String("invitee_user_id", inviteeUserID),
		zap.Bool("user_created", userCreated))

	s.respondJSON(w, http.StatusCreated, map[string]interface{}{
		"message":    "Invite sent successfully",
		"invite_url": inviteURL,
		"user_created": userCreated, // For debugging
	})
}

// handleResendInvite resends an invitation to a pending team member
func (s *Server) handleResendInvite(w http.ResponseWriter, r *http.Request, memberID string) {
	userID, ok := userIDFromContext(r.Context())
	if !ok || userID == "" {
		s.respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get user's profile to verify account ownership
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

	// Check if user is account owner
	tier, _ := profile["subscription_tier"].(string)
	if tier != "pro" && tier != "team" {
		s.respondError(w, http.StatusForbidden, "Team management requires a Pro or Team subscription")
		return
	}

	stripeSubscriptionID, _ := profile["stripe_subscription_id"].(string)
	if stripeSubscriptionID == "" && tier != "pro" && tier != "team" {
		s.respondError(w, http.StatusForbidden, "Only account owners can resend invitations")
		return
	}

	// Get the team member record
	var members []map[string]interface{}
	data, _, err := s.serviceRole.From("team_members").
		Select("*", "", false).
		Eq("id", memberID).
		Eq("account_owner_id", userID).
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
	memberEmail, _ := member["email"].(string)
	memberStatus, _ := member["status"].(string)

	// Only allow resending for pending invites
	if memberStatus != "pending" {
		s.respondError(w, http.StatusBadRequest, "Can only resend invitations for pending members")
		return
	}

	// Get or generate invite token
	inviteToken, _ := member["invite_token"].(string)
	if inviteToken == "" {
		// Generate new token if one doesn't exist
		inviteToken = uuid.New().String()
		// Update the record with the new token
		_, _, err = s.serviceRole.From("team_members").
			Update(map[string]interface{}{
				"invite_token": inviteToken,
			}, "", "").
			Eq("id", memberID).
			Execute()
		if err != nil {
			s.logger.Error("Failed to update invite token", zap.Error(err))
			s.respondError(w, http.StatusInternalServerError, "Failed to update invite token")
			return
		}
	}

	// Build invite URL
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		// Default to localhost for local development, production URL otherwise
		// Check if running locally (common indicators)
		if os.Getenv("PORT") == "" || os.Getenv("PORT") == "8080" {
			appURL = "http://localhost:5173"
		} else {
			appURL = "https://app.barracudaseo.com"
		}
	}
	inviteURL := fmt.Sprintf("%s/#/team/accept?token=%s", appURL, inviteToken)

	// Check if user exists (for determining if it's a new user)
	_, userCreated, err := s.ensureUserExists(strings.ToLower(memberEmail))
	if err != nil {
		s.logger.Warn("Failed to ensure user exists, but resending invite", zap.Error(err), zap.String("email", memberEmail))
		userCreated = false // Assume existing user if we can't check
	}

	// Send invite email via configured email service
	if err := s.emailService.SendTeamInvite(strings.ToLower(memberEmail), inviteURL, userCreated); err != nil {
		s.logger.Warn("Failed to send resend email, but invite URL is valid", zap.Error(err), zap.String("email", memberEmail))
		// Don't fail the request - invite URL is still returned
	}

	// Update invited_at timestamp
	_, _, err = s.serviceRole.From("team_members").
		Update(map[string]interface{}{
			"invited_at": time.Now().Format(time.RFC3339),
		}, "", "").
		Eq("id", memberID).
		Execute()

	if err != nil {
		s.logger.Warn("Failed to update invited_at timestamp", zap.Error(err))
		// Don't fail - email was sent
	}

	s.logger.Info("Team invite resent",
		zap.String("account_owner_id", userID),
		zap.String("member_id", memberID),
		zap.String("email", memberEmail))

	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"message":    "Invitation resent successfully",
		"invite_url": inviteURL,
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

// handleGetInviteDetailsPublic is a public version that doesn't require auth context
func (s *Server) handleGetInviteDetailsPublic(w http.ResponseWriter, r *http.Request, tokenOrID string) {
	s.handleGetInviteDetails(w, r, tokenOrID)
}

// handleGetInviteDetails returns invite details by token (no auth required)
func (s *Server) handleGetInviteDetails(w http.ResponseWriter, r *http.Request, tokenOrID string) {
	// Find invite by token
	var members []map[string]interface{}
	data, _, err := s.serviceRole.From("team_members").
		Select("id,email,role,status,invited_at", "", false).
		Eq("invite_token", tokenOrID).
		Eq("status", "pending").
		Execute()

	if err != nil {
		s.logger.Error("Failed to fetch invite details", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to fetch invite details")
		return
	}

	if err := json.Unmarshal(data, &members); err != nil || len(members) == 0 {
		s.respondError(w, http.StatusNotFound, "Invite not found or already accepted")
		return
	}

	member := members[0]

	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"email":      member["email"],
		"role":       member["role"],
		"invited_at": member["invited_at"],
	})
}

// ensureUserExists checks if a user exists by email, creates one if not
// Returns: (userID, userCreated, error)
func (s *Server) ensureUserExists(email string) (string, bool, error) {
	// Check if user exists via Supabase Admin API
	supabaseURL := s.config.SupabaseURL
	serviceKey := s.config.SupabaseServiceKey

	// Query for user by email
	checkURL := fmt.Sprintf("%s/auth/v1/admin/users?email=%s", supabaseURL, email)
	req, err := http.NewRequest("GET", checkURL, nil)
	if err != nil {
		return "", false, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("apikey", serviceKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", serviceKey))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", false, fmt.Errorf("failed to check user: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", false, fmt.Errorf("failed to check user: status %d, body: %s", resp.StatusCode, string(body))
	}

	var usersResponse struct {
		Users []struct {
			ID    string `json:"id"`
			Email string `json:"email"`
		} `json:"users"`
	}

	if err := json.Unmarshal(body, &usersResponse); err != nil {
		return "", false, fmt.Errorf("failed to parse response: %w", err)
	}

	// If user exists, return their ID
	if len(usersResponse.Users) > 0 {
		return usersResponse.Users[0].ID, false, nil
	}

	// User doesn't exist - create one
	createURL := fmt.Sprintf("%s/auth/v1/admin/users", supabaseURL)
	createPayload := map[string]interface{}{
		"email": email,
		"email_confirm": true, // Auto-confirm email for team invites
		"user_metadata": map[string]interface{}{
			"invited_to_team": true,
		},
	}

	jsonData, err := json.Marshal(createPayload)
	if err != nil {
		return "", false, fmt.Errorf("failed to marshal create payload: %w", err)
	}

	createReq, err := http.NewRequest("POST", createURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", false, fmt.Errorf("failed to create request: %w", err)
	}
	createReq.Header.Set("apikey", serviceKey)
	createReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", serviceKey))
	createReq.Header.Set("Content-Type", "application/json")

	createResp, err := client.Do(createReq)
	if err != nil {
		return "", false, fmt.Errorf("failed to create user: %w", err)
	}
	defer createResp.Body.Close()

	createBody, err := io.ReadAll(createResp.Body)
	if err != nil {
		return "", false, fmt.Errorf("failed to read create response: %w", err)
	}

	if createResp.StatusCode != http.StatusOK && createResp.StatusCode != http.StatusCreated {
		return "", false, fmt.Errorf("failed to create user: status %d, body: %s", createResp.StatusCode, string(createBody))
	}

	var createdUser struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}

	if err := json.Unmarshal(createBody, &createdUser); err != nil {
		return "", false, fmt.Errorf("failed to parse create response: %w", err)
	}

	s.logger.Info("Created new user for team invite", zap.String("email", email), zap.String("user_id", createdUser.ID))
	return createdUser.ID, true, nil
}

// sendSupabaseMagicLink sends a magic link for new users via Supabase
// This is called when using Supabase email provider for new users
func (s *Server) sendSupabaseMagicLink(userID, email, inviteURL string) error {
	supabaseURL := s.config.SupabaseURL
	serviceKey := s.config.SupabaseServiceKey

	// Update user metadata with invite URL
	updateURL := fmt.Sprintf("%s/auth/v1/admin/users/%s", supabaseURL, userID)
	updatePayload := map[string]interface{}{
		"user_metadata": map[string]interface{}{
			"team_invite_url": inviteURL,
			"invited_to_team": true,
		},
	}

	jsonData, err := json.Marshal(updatePayload)
	if err != nil {
		return fmt.Errorf("failed to marshal update payload: %w", err)
	}

	req, err := http.NewRequest("PUT", updateURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create update request: %w", err)
	}
	req.Header.Set("apikey", serviceKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", serviceKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		s.logger.Warn("Failed to update user metadata with invite URL", zap.Error(err))
		// Continue - invite URL is still valid
	} else {
		resp.Body.Close()
	}

	// Send magic link to new user (Supabase will handle email sending)
	magicLinkURL := fmt.Sprintf("%s/auth/v1/admin/users/%s/generate_link", supabaseURL, userID)
	magicLinkPayload := map[string]interface{}{
		"type":         "magiclink",
		"redirect_to":  inviteURL, // Redirect to invite acceptance after signup
	}

	jsonData, err = json.Marshal(magicLinkPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal magic link payload: %w", err)
	}

	req, err = http.NewRequest("POST", magicLinkURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create magic link request: %w", err)
	}
	req.Header.Set("apikey", serviceKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", serviceKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to generate magic link: %w", err)
	}
	defer resp.Body.Close()

	s.logger.Info("Magic link generated for new user", zap.String("email", email))
	return nil
}

// sendSupabaseInvite sends an invite email for existing users via Supabase Admin API
func (s *Server) sendSupabaseInvite(userID, email, inviteURL string) error {
	supabaseURL := s.config.SupabaseURL
	serviceKey := s.config.SupabaseServiceKey

	inviteURLAPI := fmt.Sprintf("%s/auth/v1/admin/users/%s/invite", supabaseURL, userID)
	
	invitePayload := map[string]interface{}{
		"data": map[string]interface{}{
			"invite_url": inviteURL,
			"team_invite": true,
		},
		"redirect_to": inviteURL,
	}

	jsonData, err := json.Marshal(invitePayload)
	if err != nil {
		return fmt.Errorf("failed to marshal invite payload: %w", err)
	}

	req, err := http.NewRequest("POST", inviteURLAPI, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create invite request: %w", err)
	}
	req.Header.Set("apikey", serviceKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", serviceKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send invite: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("invite endpoint returned error: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

