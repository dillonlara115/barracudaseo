-- Audit script to identify duplicate SELECT policies and RLS performance issues
-- Run this to see current state before applying the fix migration

-- Find all tables with multiple permissive SELECT policies for authenticated role
SELECT 
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

-- Show all SELECT policies with their definitions
SELECT 
  tablename,
  policyname,
  cmd as command,
  roles,
  CASE 
    WHEN qual IS NOT NULL THEN 'Has USING clause'
    ELSE 'No USING clause'
  END as has_using_clause,
  pg_get_expr(pol.polqual, pol.polrelid) as using_expression
FROM pg_policies p
LEFT JOIN pg_policy pol ON pol.polname = p.policyname
LEFT JOIN pg_class pc ON pol.polrelid = pc.oid AND pc.relname = p.tablename
WHERE p.schemaname = 'public'
  AND p.cmd = 'SELECT'
  AND (p.roles = '{authenticated}' OR p.roles = '{public}' OR p.roles IS NULL)
ORDER BY p.tablename, p.policyname;

-- Check for auth.uid() calls that aren't wrapped in SELECT
-- This is a simplified check - actual patterns may vary
SELECT 
  tablename,
  policyname,
  'Potential unwrapped auth.uid()' as issue
FROM pg_policies p
LEFT JOIN pg_policy pol ON pol.polname = p.policyname
LEFT JOIN pg_class pc ON pol.polrelid = pc.oid AND pc.relname = p.tablename
WHERE p.schemaname = 'public'
  AND p.cmd = 'SELECT'
  AND pg_get_expr(pol.polqual, pol.polrelid) LIKE '%auth.uid()%'
  AND pg_get_expr(pol.polqual, pol.polrelid) NOT LIKE '%(SELECT auth.uid())%'
ORDER BY p.tablename, p.policyname;
