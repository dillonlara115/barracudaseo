-- Optimize RLS policies for pages and issues tables
-- Add composite indexes to speed up team member checks in RLS policies
-- This should fix 500 errors caused by slow RLS policy evaluation

-- Add composite index for team_members queries used in RLS policies
-- This covers the common pattern: (user_id, account_owner_id, status)
create index if not exists idx_team_members_user_account_status 
  on public.team_members (user_id, account_owner_id, status) 
  where status = 'active';

-- Add composite index for project_members if not exists (should already exist via PK, but ensure it's optimized)
-- The primary key already provides (project_id, user_id), but we can add a covering index
create index if not exists idx_project_members_user_project 
  on public.project_members (user_id, project_id);

-- Add index on profiles for subscription tier checks
create index if not exists idx_profiles_subscription_tier 
  on public.profiles (id, subscription_tier) 
  where subscription_tier in ('pro', 'team');

-- Note: The RLS policies themselves are already optimized in previous migrations
-- These indexes should help the database evaluate the EXISTS clauses faster

