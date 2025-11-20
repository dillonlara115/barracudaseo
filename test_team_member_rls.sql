-- Test RLS policy logic for TEAM MEMBER (not account owner)
-- This uses the team member's user ID: dd240721-9f32-4553-8b72-ecf9379e05bd

SELECT 
  p.id as project_id,
  p.name as project_name,
  p.owner_id as project_owner_id,
  -- Check 1: Is user the project owner?
  (p.owner_id = 'dd240721-9f32-4553-8b72-ecf9379e05bd') as is_owner,
  -- Check 2: Is user a project member?
  EXISTS (
    SELECT 1 FROM project_members pm 
    WHERE pm.project_id = p.id AND pm.user_id = 'dd240721-9f32-4553-8b72-ecf9379e05bd'
  ) as is_project_member,
  -- Check 3: Are user and project owner both team members with same account_owner_id?
  EXISTS (
    SELECT 1
    FROM team_members tm1
    JOIN team_members tm2 ON tm1.account_owner_id = tm2.account_owner_id
    WHERE tm1.user_id = 'dd240721-9f32-4553-8b72-ecf9379e05bd'
      AND tm2.user_id = p.owner_id
      AND tm1.status = 'active'
      AND tm2.status = 'active'
  ) as are_both_team_members,
  -- Check 4: Is project owner an account owner and user is their team member?
  -- THIS IS THE KEY CHECK FOR TEAM MEMBERS TO SEE ACCOUNT OWNER'S PROJECTS
  EXISTS (
    SELECT 1
    FROM profiles prof
    JOIN team_members tm ON tm.account_owner_id = p.owner_id
    WHERE prof.id = p.owner_id
      AND (prof.subscription_tier = 'pro' OR prof.subscription_tier = 'team')
      AND tm.user_id = 'dd240721-9f32-4553-8b72-ecf9379e05bd'
      AND tm.status = 'active'
  ) as is_owner_team_member_match,
  -- Debug: Let's see the actual values
  (SELECT subscription_tier FROM profiles WHERE id = p.owner_id) as owner_subscription_tier,
  (SELECT COUNT(*) FROM team_members tm WHERE tm.account_owner_id = p.owner_id AND tm.user_id = 'dd240721-9f32-4553-8b72-ecf9379e05bd' AND tm.status = 'active') as team_member_count,
  -- Combined: Should user see this project?
  (
    (p.owner_id = 'dd240721-9f32-4553-8b72-ecf9379e05bd')
    OR EXISTS (
      SELECT 1 FROM project_members pm 
      WHERE pm.project_id = p.id AND pm.user_id = 'dd240721-9f32-4553-8b72-ecf9379e05bd'
    )
    OR EXISTS (
      SELECT 1
      FROM team_members tm1
      JOIN team_members tm2 ON tm1.account_owner_id = tm2.account_owner_id
      WHERE tm1.user_id = 'dd240721-9f32-4553-8b72-ecf9379e05bd'
        AND tm2.user_id = p.owner_id
        AND tm1.status = 'active'
        AND tm2.status = 'active'
    )
    OR EXISTS (
      SELECT 1
      FROM profiles prof
      JOIN team_members tm ON tm.account_owner_id = p.owner_id
      WHERE prof.id = p.owner_id
        AND (prof.subscription_tier = 'pro' OR prof.subscription_tier = 'team')
        AND tm.user_id = 'dd240721-9f32-4553-8b72-ecf9379e05bd'
        AND tm.status = 'active'
    )
  ) as should_see_project
FROM projects p
WHERE p.id = 'ed2e856c-3f3d-4fba-ac02-f968a024231f';

