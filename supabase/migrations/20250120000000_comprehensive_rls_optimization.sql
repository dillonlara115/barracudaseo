-- Comprehensive RLS Optimization Migration
-- Consolidates duplicate SELECT policies and optimizes RLS performance using SECURITY DEFINER helper functions
-- This addresses:
-- 1. Multiple permissive SELECT policies on the same table/role (causing redundant evaluations)
-- 2. auth.uid() re-evaluation per row (initplan warnings)
-- 3. Complex policy logic duplication across tables

-- ============================================================================
-- PART 1: Create SECURITY DEFINER Helper Functions
-- ============================================================================

-- Helper: Check if current user can view a profile (self or account owner)
-- This replaces the existing can_view_profile function with optimized version
create or replace function public.can_view_profile(profile_id uuid)
returns boolean
language sql
security definer
stable
set search_path = public
as $$
  select (select auth.uid()) = profile_id
    or exists (
      select 1
      from public.team_members tm
      where tm.user_id = (select auth.uid())
        and tm.account_owner_id = profile_id
        and tm.status = 'active'
    );
$$;

grant execute on function public.can_view_profile(uuid) to authenticated;

-- Helper: Check if current user has access to a project (owner, member, or teammate)
create or replace function public.can_access_project(project_id uuid)
returns boolean
language sql
security definer
stable
set search_path = public
as $$
  select exists (
    select 1
    from public.projects p
    where p.id = project_id
      and (
        -- User is project owner
        p.owner_id = (select auth.uid())
        -- OR user is a project member
        or exists (
          select 1
          from public.project_members pm
          where pm.project_id = project_id
            and pm.user_id = (select auth.uid())
        )
        -- OR user and project owner are teammates (same account_owner_id)
        or exists (
          select 1
          from public.team_members tm1
          join public.team_members tm2 on tm1.account_owner_id = tm2.account_owner_id
          where tm1.user_id = (select auth.uid())
            and tm2.user_id = p.owner_id
            and tm1.status = 'active'
            and tm2.status = 'active'
        )
        -- OR project owner has pro/team tier and user is their team member
        or exists (
          select 1
          from public.profiles prof
          join public.team_members tm on tm.account_owner_id = p.owner_id
          where prof.id = p.owner_id
            and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
            and tm.user_id = (select auth.uid())
            and tm.status = 'active'
        )
        -- OR user has pro/team tier and project owner is their team member
        or exists (
          select 1
          from public.profiles prof
          join public.team_members tm on tm.account_owner_id = (select auth.uid())
          where prof.id = (select auth.uid())
            and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
            and tm.user_id = p.owner_id
            and tm.status = 'active'
        )
      )
  );
$$;

grant execute on function public.can_access_project(uuid) to authenticated;

-- Helper: Check if current user can access a crawl (via project access)
create or replace function public.can_access_crawl(crawl_id uuid)
returns boolean
language sql
security definer
stable
set search_path = public
as $$
  select exists (
    select 1
    from public.crawls c
    where c.id = crawl_id
      and public.can_access_project(c.project_id)
  );
$$;

grant execute on function public.can_access_crawl(uuid) to authenticated;

-- Helper: Check if current user can access an issue (via project access)
create or replace function public.can_access_issue(issue_id bigint)
returns boolean
language sql
security definer
stable
set search_path = public
as $$
  select exists (
    select 1
    from public.issues i
    where i.id = issue_id
      and public.can_access_project(i.project_id)
  );
$$;

grant execute on function public.can_access_issue(bigint) to authenticated;

-- Helper: Check if current user can modify a project (owner or teammate with pro/team tier)
create or replace function public.can_modify_project(project_id uuid)
returns boolean
language sql
security definer
stable
set search_path = public
as $$
  select exists (
    select 1
    from public.projects p
    where p.id = project_id
      and (
        -- User is project owner
        p.owner_id = (select auth.uid())
        -- OR user and project owner are teammates (same account_owner_id)
        or exists (
          select 1
          from public.team_members tm1
          join public.team_members tm2 on tm1.account_owner_id = tm2.account_owner_id
          where tm1.user_id = (select auth.uid())
            and tm2.user_id = p.owner_id
            and tm1.status = 'active'
            and tm2.status = 'active'
        )
        -- OR project owner has pro/team tier and user is their team member
        or exists (
          select 1
          from public.profiles prof
          join public.team_members tm on tm.account_owner_id = p.owner_id
          where prof.id = p.owner_id
            and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
            and tm.user_id = (select auth.uid())
            and tm.status = 'active'
        )
        -- OR user has pro/team tier and project owner is their team member
        or exists (
          select 1
          from public.profiles prof
          join public.team_members tm on tm.account_owner_id = (select auth.uid())
          where prof.id = (select auth.uid())
            and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
            and tm.user_id = p.owner_id
            and tm.status = 'active'
        )
      )
  );
$$;

grant execute on function public.can_modify_project(uuid) to authenticated;

-- ============================================================================
-- PART 2: Consolidate Profiles Table Policies
-- ============================================================================

-- Drop all existing SELECT policies on profiles
drop policy if exists "Profiles select for owners and teammates" on public.profiles;
drop policy if exists "Users can view their own profile" on public.profiles;
drop policy if exists "Team members can view account owner profile" on public.profiles;

-- Create single consolidated SELECT policy
create policy "Profiles select for owners and teammates"
  on public.profiles
  for select
  to authenticated
  using (public.can_view_profile(id));

-- ============================================================================
-- PART 3: Consolidate GSC Tables Policies
-- ============================================================================

-- gsc_sync_states
drop policy if exists "Project members can view gsc sync state" on public.gsc_sync_states;
create policy "Project members and teammates can view gsc sync state"
  on public.gsc_sync_states
  for select
  to authenticated
  using (public.can_access_project(project_id));

-- gsc_performance_snapshots
drop policy if exists "Project members can view gsc snapshots" on public.gsc_performance_snapshots;
create policy "Project members and teammates can view gsc snapshots"
  on public.gsc_performance_snapshots
  for select
  to authenticated
  using (public.can_access_project(project_id));

-- gsc_performance_rows
drop policy if exists "Project members can view gsc performance rows" on public.gsc_performance_rows;
create policy "Project members and teammates can view gsc performance rows"
  on public.gsc_performance_rows
  for select
  to authenticated
  using (public.can_access_project(project_id));

-- gsc_page_enhancements
drop policy if exists "Project members can view gsc page enhancements" on public.gsc_page_enhancements;
create policy "Project members and teammates can view gsc page enhancements"
  on public.gsc_page_enhancements
  for select
  to authenticated
  using (public.can_access_project(project_id));

-- gsc_insights
drop policy if exists "Project members can view gsc insights" on public.gsc_insights;
create policy "Project members and teammates can view gsc insights"
  on public.gsc_insights
  for select
  to authenticated
  using (public.can_access_project(project_id));

-- ============================================================================
-- PART 4: Consolidate Issue-Related Tables Policies
-- ============================================================================

-- issue_recommendations
drop policy if exists "Project members can view recommendations" on public.issue_recommendations;
drop policy if exists "Project members and teammates can view recommendations" on public.issue_recommendations;
create policy "Project members and teammates can view recommendations"
  on public.issue_recommendations
  for select
  to authenticated
  using (public.can_access_issue(issue_id));

-- issue_status_history
drop policy if exists "Project members can view status history" on public.issue_status_history;
drop policy if exists "Project members and teammates can view status history" on public.issue_status_history;
create policy "Project members and teammates can view status history"
  on public.issue_status_history
  for select
  to authenticated
  using (public.can_access_issue(issue_id));

-- ============================================================================
-- PART 5: Consolidate Exports Table Policies
-- ============================================================================

-- exports
drop policy if exists "Project members can view exports" on public.exports;
drop policy if exists "Project members and teammates can view exports" on public.exports;
drop policy if exists "Project members can create exports" on public.exports;
drop policy if exists "Project members and teammates can create exports" on public.exports;

-- Consolidated SELECT policy
create policy "Project members and teammates can view exports"
  on public.exports
  for select
  to authenticated
  using (public.can_access_project(project_id));

-- Consolidated INSERT policy
create policy "Project members and teammates can create exports"
  on public.exports
  for insert
  to authenticated
  with check (public.can_access_project(project_id));

-- Keep DELETE policy as-is (it has specific logic for requesters/owners)
-- But ensure it uses wrapped auth.uid()
drop policy if exists "Export requesters and owners can delete exports" on public.exports;
create policy "Export requesters and owners can delete exports"
  on public.exports
  for delete
  to authenticated
  using (
    requested_by = (select auth.uid())
    or exists (
      select 1
      from public.projects p
      where p.id = exports.project_id
        and p.owner_id = (select auth.uid())
    )
  );

-- ============================================================================
-- PART 6: Consolidate AI Tables Policies
-- ============================================================================

-- user_ai_settings - single policy is fine, just ensure wrapped auth.uid()
drop policy if exists "Users can view their own AI settings" on public.user_ai_settings;
create policy "Users can view their own AI settings"
  on public.user_ai_settings
  for select
  to authenticated
  using (user_id = (select auth.uid()));

-- ai_issue_insights
drop policy if exists "Users can view their own AI issue insights" on public.ai_issue_insights;
create policy "Users can view their own AI issue insights"
  on public.ai_issue_insights
  for select
  to authenticated
  using (
    user_id = (select auth.uid())
    and public.can_access_project(project_id)
  );

drop policy if exists "Users can create their own AI issue insights" on public.ai_issue_insights;
create policy "Users can create their own AI issue insights"
  on public.ai_issue_insights
  for insert
  to authenticated
  with check (
    user_id = (select auth.uid())
    and public.can_access_project(project_id)
  );

-- ai_crawl_summaries
drop policy if exists "Users can view AI crawl summaries for accessible crawls" on public.ai_crawl_summaries;
create policy "Users can view AI crawl summaries for accessible crawls"
  on public.ai_crawl_summaries
  for select
  to authenticated
  using (
    user_id = (select auth.uid())
    and public.can_access_project(project_id)
  );

drop policy if exists "Users can create AI crawl summaries for accessible crawls" on public.ai_crawl_summaries;
create policy "Users can create AI crawl summaries for accessible crawls"
  on public.ai_crawl_summaries
  for insert
  to authenticated
  with check (
    user_id = (select auth.uid())
    and public.can_access_project(project_id)
  );

-- ============================================================================
-- PART 7: Add Indexes for Performance
-- ============================================================================

-- Indexes for project access checks
create index if not exists idx_project_members_user_project 
  on public.project_members (user_id, project_id);

create index if not exists idx_team_members_user_account_owner 
  on public.team_members (user_id, account_owner_id) 
  where status = 'active';

create index if not exists idx_profiles_subscription_tier 
  on public.profiles (id, subscription_tier) 
  where subscription_tier IN ('pro', 'team');

-- Indexes for issue access checks
create index if not exists idx_issues_project_id 
  on public.issues (project_id);

-- Indexes for crawl access checks
create index if not exists idx_crawls_project_id 
  on public.crawls (project_id);

-- Indexes for exports
create index if not exists idx_exports_project_requested 
  on public.exports (project_id, requested_by);

-- ============================================================================
-- PART 8: Comments for Documentation
-- ============================================================================

comment on function public.can_view_profile(uuid) is 
  'Check if current user can view a profile (self or account owner). Used in RLS policies.';

comment on function public.can_access_project(uuid) is 
  'Check if current user has access to a project (owner, member, or teammate). Used in RLS policies.';

comment on function public.can_access_crawl(uuid) is 
  'Check if current user can access a crawl via project access. Used in RLS policies.';

comment on function public.can_access_issue(bigint) is 
  'Check if current user can access an issue via project access. Used in RLS policies.';

comment on function public.can_modify_project(uuid) is 
  'Check if current user can modify a project (owner or teammate). Used in RLS policies.';


