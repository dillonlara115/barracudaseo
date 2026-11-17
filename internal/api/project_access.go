package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// accessibleProjectsFilter builds the PostgREST "or" filter that selects the
// projects a user can reach: ones they own, ones owned by a teammate, and ones
// they were added to as a member. It mirrors verifyProjectAccess, for use with
// the service role client, which is not filtered by RLS.
func (s *Server) accessibleProjectsFilter(userID string) (string, error) {
	owners, err := s.teamUserIDs(userID)
	if err != nil {
		return "", err
	}

	data, _, err := s.serviceRole.From("project_members").
		Select("project_id", "", false).
		Eq("user_id", userID).
		Execute()
	if err != nil {
		return "", fmt.Errorf("failed to load project memberships: %w", err)
	}
	var memberships []struct {
		ProjectID string `json:"project_id"`
	}
	if err := json.Unmarshal(data, &memberships); err != nil {
		return "", fmt.Errorf("failed to parse project memberships: %w", err)
	}
	projectIDs := make([]string, 0, len(memberships))
	for _, membership := range memberships {
		projectIDs = append(projectIDs, membership.ProjectID)
	}

	return projectsOrFilter(owners, projectIDs), nil
}

// projectsOrFilter renders the filter for projects owned by any of ownerIDs or
// whose ID is in projectIDs. Values that are not UUIDs are dropped, so nothing
// but IDs can end up in the filter expression.
func projectsOrFilter(ownerIDs, projectIDs []string) string {
	parts := []string{"owner_id.in.(" + strings.Join(validUUIDs(ownerIDs), ",") + ")"}
	if ids := validUUIDs(projectIDs); len(ids) > 0 {
		parts = append(parts, "id.in.("+strings.Join(ids, ",")+")")
	}
	return strings.Join(parts, ",")
}

func validUUIDs(values []string) []string {
	seen := make(map[string]bool, len(values))
	valid := make([]string, 0, len(values))
	for _, value := range values {
		if _, err := uuid.Parse(value); err != nil || seen[value] {
			continue
		}
		seen[value] = true
		valid = append(valid, value)
	}
	return valid
}

// teamUserIDs returns the user plus everyone on the same team: the account
// owner and the owner's active team members. A user without a team gets only
// themselves.
func (s *Server) teamUserIDs(userID string) ([]string, error) {
	ids := []string{userID}

	profile, err := s.fetchProfile(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch profile: %w", err)
	}
	if profile == nil {
		return ids, nil
	}
	accountOwnerID := s.getAccountOwnerID(userID, profile)
	if accountOwnerID == "" {
		return ids, nil
	}
	ids = append(ids, accountOwnerID)

	data, _, err := s.serviceRole.From("team_members").
		Select("user_id", "", false).
		Eq("account_owner_id", accountOwnerID).
		Eq("status", "active").
		Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to load team members: %w", err)
	}
	var members []struct {
		UserID *string `json:"user_id"`
	}
	if err := json.Unmarshal(data, &members); err != nil {
		return nil, fmt.Errorf("failed to parse team members: %w", err)
	}
	for _, member := range members {
		if member.UserID != nil {
			ids = append(ids, *member.UserID)
		}
	}
	return ids, nil
}

// canDeleteProject reports whether the user may delete a project: the owner
// and the owner's teammates can, a member who was only added to the project
// cannot. found is false when the project does not exist.
func (s *Server) canDeleteProject(userID, projectID string) (allowed, found bool, err error) {
	data, _, err := s.serviceRole.From("projects").
		Select("owner_id", "", false).
		Eq("id", projectID).
		Execute()
	if err != nil {
		return false, false, fmt.Errorf("failed to load project: %w", err)
	}
	var projects []struct {
		OwnerID string `json:"owner_id"`
	}
	if err := json.Unmarshal(data, &projects); err != nil {
		return false, false, fmt.Errorf("failed to parse project: %w", err)
	}
	if len(projects) == 0 {
		return false, false, nil
	}

	if projects[0].OwnerID == userID {
		return true, true, nil
	}
	teammates, err := s.areUsersOnSameTeam(userID, projects[0].OwnerID)
	if err != nil {
		return false, true, fmt.Errorf("failed to check team membership: %w", err)
	}
	return teammates, true, nil
}
