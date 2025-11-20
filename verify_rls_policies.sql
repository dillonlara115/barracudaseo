-- Verify RLS policies exist and are correct

-- Check if the projects SELECT policy exists
SELECT 
  schemaname,
  tablename,
  policyname,
  permissive,
  roles,
  cmd,
  qual,
  with_check
FROM pg_policies
WHERE tablename = 'projects'
  AND policyname = 'Project members and teammates can view projects';

-- Check all policies on projects table
SELECT 
  policyname,
  cmd as command,
  CASE 
    WHEN qual IS NOT NULL THEN 'Has USING clause'
    ELSE 'No USING clause'
  END as has_using_clause
FROM pg_policies
WHERE tablename = 'projects'
ORDER BY policyname;

-- Test: Can we query projects as the team member using SET ROLE?
-- Note: This won't work directly in SQL Editor, but let's check the policy definition
SELECT 
  pg_get_expr(pol.polqual, pol.polrelid) as using_expression
FROM pg_policy pol
JOIN pg_class pc ON pol.polrelid = pc.oid
WHERE pc.relname = 'projects'
  AND pol.polname = 'Project members and teammates can view projects';

