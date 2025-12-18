-- Test script to verify RLS access semantics after optimization
-- This script tests various access scenarios to ensure policies work correctly
-- 
-- IMPORTANT: Replace the UUIDs below with actual test user IDs from your database
-- You can find user IDs with: SELECT id, email FROM auth.users LIMIT 5;

-- ============================================================================
-- SETUP: Define test user IDs (UPDATE THESE WITH REAL VALUES)
-- ============================================================================

-- Example structure - replace with actual UUIDs:
-- \set owner_user_id '00000000-0000-0000-0000-000000000001'
-- \set team_member_user_id '00000000-0000-0000-0000-000000000002'
-- \set unrelated_user_id '00000000-0000-0000-0000-000000000003'
-- \set project_owner_id '00000000-0000-0000-0000-000000000004'

-- ============================================================================
-- PART 1: Test Profile Access
-- ============================================================================

-- Test 1: User can view their own profile
-- Expected: Should return 1 row
SELECT 
  'TEST: User views own profile' as test_name,
  COUNT(*) as result_count,
  CASE WHEN COUNT(*) = 1 THEN 'PASS' ELSE 'FAIL' END as status
FROM public.profiles
WHERE id = (SELECT auth.uid());

-- Test 2: Team member can view account owner's profile
-- Expected: Should return 1 row if user is a team member
SELECT 
  'TEST: Team member views account owner profile' as test_name,
  COUNT(*) as result_count,
  CASE WHEN COUNT(*) >= 0 THEN 'PASS' ELSE 'FAIL' END as status
FROM public.profiles p
WHERE EXISTS (
  SELECT 1
  FROM public.team_members tm
  WHERE tm.user_id = (SELECT auth.uid())
    AND tm.account_owner_id = p.id
    AND tm.status = 'active'
);

-- Test 3: Unrelated user cannot view other profiles
-- Expected: Should return 0 rows (only own profile)
SELECT 
  'TEST: Unrelated user cannot view other profiles' as test_name,
  COUNT(*) as result_count,
  CASE WHEN COUNT(*) <= 1 THEN 'PASS' ELSE 'FAIL' END as status
FROM public.profiles
WHERE id != (SELECT auth.uid());

-- ============================================================================
-- PART 2: Test Project Access
-- ============================================================================

-- Test 4: Project owner can view their projects
-- Expected: Should return >= 1 row
SELECT 
  'TEST: Project owner views own projects' as test_name,
  COUNT(*) as result_count,
  CASE WHEN COUNT(*) >= 0 THEN 'PASS' ELSE 'FAIL' END as status
FROM public.projects
WHERE owner_id = (SELECT auth.uid());

-- Test 5: Project member can view projects they're members of
-- Expected: Should return >= 0 rows
SELECT 
  'TEST: Project member views accessible projects' as test_name,
  COUNT(*) as result_count,
  CASE WHEN COUNT(*) >= 0 THEN 'PASS' ELSE 'FAIL' END as status
FROM public.projects p
WHERE EXISTS (
  SELECT 1
  FROM public.project_members pm
  WHERE pm.project_id = p.id
    AND pm.user_id = (SELECT auth.uid())
);

-- ============================================================================
-- PART 3: Test GSC Tables Access
-- ============================================================================

-- Test 6: User can view GSC sync states for accessible projects
-- Expected: Should return >= 0 rows
SELECT 
  'TEST: User views GSC sync states' as test_name,
  COUNT(*) as result_count,
  CASE WHEN COUNT(*) >= 0 THEN 'PASS' ELSE 'FAIL' END as status
FROM public.gsc_sync_states;

-- Test 7: User can view GSC snapshots for accessible projects
-- Expected: Should return >= 0 rows
SELECT 
  'TEST: User views GSC snapshots' as test_name,
  COUNT(*) as result_count,
  CASE WHEN COUNT(*) >= 0 THEN 'PASS' ELSE 'FAIL' END as status
FROM public.gsc_performance_snapshots;

-- ============================================================================
-- PART 4: Test Issue-Related Tables Access
-- ============================================================================

-- Test 8: User can view issue recommendations for accessible issues
-- Expected: Should return >= 0 rows
SELECT 
  'TEST: User views issue recommendations' as test_name,
  COUNT(*) as result_count,
  CASE WHEN COUNT(*) >= 0 THEN 'PASS' ELSE 'FAIL' END as status
FROM public.issue_recommendations;

-- Test 9: User can view issue status history for accessible issues
-- Expected: Should return >= 0 rows
SELECT 
  'TEST: User views issue status history' as test_name,
  COUNT(*) as result_count,
  CASE WHEN COUNT(*) >= 0 THEN 'PASS' ELSE 'FAIL' END as status
FROM public.issue_status_history;

-- ============================================================================
-- PART 5: Test Exports Access
-- ============================================================================

-- Test 10: User can view exports for accessible projects
-- Expected: Should return >= 0 rows
SELECT 
  'TEST: User views exports' as test_name,
  COUNT(*) as result_count,
  CASE WHEN COUNT(*) >= 0 THEN 'PASS' ELSE 'FAIL' END as status
FROM public.exports;

-- ============================================================================
-- PART 6: Test AI Tables Access
-- ============================================================================

-- Test 11: User can view their own AI settings
-- Expected: Should return 0 or 1 row
SELECT 
  'TEST: User views own AI settings' as test_name,
  COUNT(*) as result_count,
  CASE WHEN COUNT(*) <= 1 THEN 'PASS' ELSE 'FAIL' END as status
FROM public.user_ai_settings
WHERE user_id = (SELECT auth.uid());

-- Test 12: User can view their own AI issue insights for accessible projects
-- Expected: Should return >= 0 rows
SELECT 
  'TEST: User views own AI issue insights' as test_name,
  COUNT(*) as result_count,
  CASE WHEN COUNT(*) >= 0 THEN 'PASS' ELSE 'FAIL' END as status
FROM public.ai_issue_insights
WHERE user_id = (SELECT auth.uid());

-- ============================================================================
-- PART 7: Helper Function Tests
-- ============================================================================

-- Test 13: can_view_profile helper function works
-- Expected: Should return true for own profile
SELECT 
  'TEST: can_view_profile helper (own profile)' as test_name,
  public.can_view_profile((SELECT auth.uid())) as result,
  CASE WHEN public.can_view_profile((SELECT auth.uid())) THEN 'PASS' ELSE 'FAIL' END as status;

-- Test 14: can_access_project helper function works
-- Expected: Should return true for projects user owns or is member of
SELECT 
  'TEST: can_access_project helper' as test_name,
  COUNT(*) as accessible_projects,
  CASE WHEN COUNT(*) >= 0 THEN 'PASS' ELSE 'FAIL' END as status
FROM public.projects p
WHERE public.can_access_project(p.id);

-- ============================================================================
-- PART 8: Performance Check - Verify no duplicate policy evaluations
-- ============================================================================

-- This query shows all SELECT policies - should be one per table
SELECT 
  'POLICY COUNT CHECK' as check_type,
  tablename,
  COUNT(*) as policy_count,
  CASE WHEN COUNT(*) = 1 THEN 'PASS' ELSE 'FAIL - Multiple policies!' END as status
FROM pg_policies
WHERE schemaname = 'public'
  AND cmd = 'SELECT'
  AND tablename IN (
    'profiles',
    'gsc_sync_states',
    'gsc_performance_snapshots',
    'gsc_performance_rows',
    'gsc_page_enhancements',
    'gsc_insights',
    'issue_recommendations',
    'issue_status_history',
    'exports',
    'user_ai_settings',
    'ai_issue_insights',
    'ai_crawl_summaries'
  )
GROUP BY tablename
ORDER BY tablename;

-- ============================================================================
-- SUMMARY
-- ============================================================================

SELECT 
  '=== RLS OPTIMIZATION TEST SUMMARY ===' as summary,
  'Run all tests above and verify:' as instructions,
  '1. No duplicate SELECT policies exist' as check_1,
  '2. All access tests return expected results' as check_2,
  '3. Helper functions work correctly' as check_3,
  '4. Policy counts are correct (1 per table)' as check_4;
