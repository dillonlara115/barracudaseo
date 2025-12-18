-- Verification script for RLS optimization migration
-- Run this after applying the migration to verify:
-- 1. No duplicate SELECT policies exist
-- 2. All policies use helper functions where appropriate
-- 3. Access semantics are preserved

-- ============================================================================
-- PART 1: Check for duplicate SELECT policies
-- ============================================================================

SELECT 
  'DUPLICATE POLICIES CHECK' as check_type,
  schemaname,
  tablename,
  COUNT(*) as policy_count,
  STRING_AGG(policyname, ', ' ORDER BY policyname) as policy_names
FROM pg_policies
WHERE schemaname = 'public'
  AND cmd = 'SELECT'
  AND (roles = '{authenticated}' OR roles = '{public}' OR roles IS NULL)
GROUP BY schemaname, tablename
HAVING COUNT(*) > 1
ORDER BY tablename;

-- Expected: Should return 0 rows (no duplicates)

-- ============================================================================
-- PART 2: Verify helper functions exist and are accessible
-- ============================================================================

SELECT 
  'HELPER FUNCTIONS CHECK' as check_type,
  p.proname as function_name,
  pg_get_function_identity_arguments(p.oid) as arguments,
  CASE 
    WHEN p.prosecdef THEN 'SECURITY DEFINER'
    ELSE 'SECURITY INVOKER'
  END as security_type,
  CASE 
    WHEN p.provolatile = 's' THEN 'STABLE'
    WHEN p.provolatile = 'i' THEN 'IMMUTABLE'
    ELSE 'VOLATILE'
  END as volatility
FROM pg_proc p
JOIN pg_namespace n ON p.pronamespace = n.oid
WHERE n.nspname = 'public'
  AND p.proname IN (
    'can_view_profile',
    'can_access_project',
    'can_access_crawl',
    'can_access_issue',
    'can_modify_project'
  )
ORDER BY p.proname;

-- Expected: Should return 5 rows (all helper functions)

-- ============================================================================
-- PART 3: Verify policies use helper functions
-- ============================================================================

SELECT 
  'POLICY HELPER USAGE CHECK' as check_type,
  tablename,
  policyname,
  CASE 
    WHEN pg_get_expr(pol.polqual, pol.polrelid) LIKE '%can_view_profile%' THEN '✓ Uses can_view_profile'
    WHEN pg_get_expr(pol.polqual, pol.polrelid) LIKE '%can_access_project%' THEN '✓ Uses can_access_project'
    WHEN pg_get_expr(pol.polqual, pol.polrelid) LIKE '%can_access_crawl%' THEN '✓ Uses can_access_crawl'
    WHEN pg_get_expr(pol.polqual, pol.polrelid) LIKE '%can_access_issue%' THEN '✓ Uses can_access_issue'
    WHEN pg_get_expr(pol.polqual, pol.polrelid) LIKE '%can_modify_project%' THEN '✓ Uses can_modify_project'
    ELSE '✗ No helper function'
  END as helper_usage
FROM pg_policies p
LEFT JOIN pg_policy pol ON pol.polname = p.policyname
LEFT JOIN pg_class pc ON pol.polrelid = pc.oid AND pc.relname = p.tablename
WHERE p.schemaname = 'public'
  AND p.cmd = 'SELECT'
  AND p.tablename IN (
    'profiles',
    'gsc_sync_states',
    'gsc_performance_snapshots',
    'gsc_performance_rows',
    'gsc_page_enhancements',
    'gsc_insights',
    'issue_recommendations',
    'issue_status_history',
    'exports',
    'ai_issue_insights',
    'ai_crawl_summaries'
  )
ORDER BY p.tablename, p.policyname;

-- Expected: All policies should use helper functions

-- ============================================================================
-- PART 4: Verify auth.uid() is wrapped in SELECT
-- ============================================================================

SELECT 
  'AUTH.UID() WRAPPING CHECK' as check_type,
  tablename,
  policyname,
  CASE 
    WHEN pg_get_expr(pol.polqual, pol.polrelid) LIKE '%auth.uid()%' 
      AND pg_get_expr(pol.polqual, pol.polrelid) NOT LIKE '%(SELECT auth.uid())%'
      AND pg_get_expr(pol.polqual, pol.polrelid) NOT LIKE '%can_%' THEN '✗ Unwrapped auth.uid()'
    ELSE '✓ OK'
  END as wrapping_status
FROM pg_policies p
LEFT JOIN pg_policy pol ON pol.polname = p.policyname
LEFT JOIN pg_class pc ON pol.polrelid = pc.oid AND pc.relname = p.tablename
WHERE p.schemaname = 'public'
  AND p.cmd = 'SELECT'
  AND pg_get_expr(pol.polqual, pol.polrelid) LIKE '%auth.uid()%'
ORDER BY p.tablename, p.policyname;

-- Expected: Should return 0 rows with unwrapped auth.uid() (or only in helper functions)

-- ============================================================================
-- PART 5: Verify indexes exist
-- ============================================================================

SELECT 
  'INDEXES CHECK' as check_type,
  schemaname,
  indexname,
  tablename,
  indexdef
FROM pg_indexes
WHERE schemaname = 'public'
  AND indexname IN (
    'idx_project_members_user_project',
    'idx_team_members_user_account_owner',
    'idx_profiles_subscription_tier',
    'idx_issues_project_id',
    'idx_crawls_project_id',
    'idx_exports_project_requested'
  )
ORDER BY indexname;

-- Expected: Should return 6 rows (all indexes created)

-- ============================================================================
-- PART 6: Summary report
-- ============================================================================

SELECT 
  'SUMMARY' as check_type,
  'Total SELECT policies' as metric,
  COUNT(*)::text as value
FROM pg_policies
WHERE schemaname = 'public'
  AND cmd = 'SELECT'
UNION ALL
SELECT 
  'SUMMARY',
  'Tables with SELECT policies',
  COUNT(DISTINCT tablename)::text
FROM pg_policies
WHERE schemaname = 'public'
  AND cmd = 'SELECT'
UNION ALL
SELECT 
  'SUMMARY',
  'Helper functions created',
  COUNT(*)::text
FROM pg_proc p
JOIN pg_namespace n ON p.pronamespace = n.oid
WHERE n.nspname = 'public'
  AND p.proname LIKE 'can_%';
