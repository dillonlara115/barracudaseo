-- Test if team member can now read account owner's profile
-- This simulates what the RLS policy checks

-- Test 1: Can team member see their own profile?
SELECT 
  id,
  subscription_tier,
  'Own profile' as profile_type
FROM profiles
WHERE id = 'dd240721-9f32-4553-8b72-ecf9379e05bd';

-- Test 2: Can team member see account owner's profile?
-- This should work now with the new policy
SELECT 
  id,
  subscription_tier,
  'Account owner profile' as profile_type
FROM profiles
WHERE id = '34aac771-41fa-4a14-912f-2b6d90dc313e';

-- Test 3: Verify the policy exists
SELECT 
  policyname,
  cmd as command,
  CASE 
    WHEN qual IS NOT NULL THEN 'Has USING clause'
    ELSE 'No USING clause'
  END as has_using_clause
FROM pg_policies
WHERE tablename = 'profiles'
  AND policyname = 'Team members can view account owner profile';

-- Test 4: Now test if projects are visible (this should work now)
SELECT 
  p.id,
  p.name,
  p.owner_id,
  owner_profile.subscription_tier as owner_tier,
  'Should be visible' as visibility
FROM projects p
LEFT JOIN profiles owner_profile ON owner_profile.id = p.owner_id
WHERE p.id = 'ed2e856c-3f3d-4fba-ac02-f968a024231f';

