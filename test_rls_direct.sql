-- Direct test of RLS policy for team member
-- Replace auth.uid() with the team member's user ID to test

-- First, let's check if we can see the team member record when querying as the team member
-- This simulates what RLS sees when auth.uid() = team member ID

-- Test 1: Can team member see their own team_members record?
SELECT 
  tm.*,
  'Team member can see own record' as test_name
FROM team_members tm
WHERE tm.user_id = 'dd240721-9f32-4553-8b72-ecf9379e05bd';

-- Test 2: Can team member see projects via RLS policy?
-- This should return the project if RLS is working
SELECT 
  p.id,
  p.name,
  p.owner_id,
  'Should be visible via RLS' as visibility_test
FROM projects p
WHERE p.id = 'ed2e856c-3f3d-4fba-ac02-f968a024231f';

-- Test 3: Manual check of the RLS policy condition
-- This simulates what the RLS policy checks
SELECT 
  p.id as project_id,
  p.name as project_name,
  p.owner_id as project_owner_id,
  -- Check if project owner is account owner and user is their team member
  EXISTS (
    SELECT 1
    FROM profiles prof
    JOIN team_members tm ON tm.account_owner_id = p.owner_id
    WHERE prof.id = p.owner_id
      AND (prof.subscription_tier = 'pro' OR prof.subscription_tier = 'team')
      AND tm.user_id = 'dd240721-9f32-4553-8b72-ecf9379e05bd'
      AND tm.status = 'active'
  ) as rls_condition_matches,
  -- Verify the join works
  (SELECT COUNT(*) FROM team_members tm WHERE tm.account_owner_id = p.owner_id AND tm.user_id = 'dd240721-9f32-4553-8b72-ecf9379e05bd' AND tm.status = 'active') as team_member_count,
  (SELECT subscription_tier FROM profiles WHERE id = p.owner_id) as owner_subscription_tier
FROM projects p
WHERE p.id = 'ed2e856c-3f3d-4fba-ac02-f968a024231f';

