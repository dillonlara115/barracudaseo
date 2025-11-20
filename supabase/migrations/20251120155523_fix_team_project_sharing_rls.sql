-- Fix RLS policies for team project sharing
-- The issue: Account owners don't have records in team_members table
-- They're identified by subscription_tier in profiles table
-- This migration fixes policies to properly check account ownership via profiles

-- Fix projects SELECT policy
drop policy if exists "Project members and teammates can view projects" on public.projects;
create policy "Project members and teammates can view projects"
  on public.projects
  for select
  using (
    owner_id = auth.uid()
    or exists (
      select 1
      from public.project_members pm
      where pm.project_id = projects.id
        and pm.user_id = auth.uid()
    )
    or exists (
      -- Check if user and project owner are both team members with same account_owner_id
      select 1
      from public.team_members tm1
      join public.team_members tm2 on tm1.account_owner_id = tm2.account_owner_id
      where tm1.user_id = auth.uid()
        and tm2.user_id = projects.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner (has pro/team tier) and user is their team member
      select 1
      from public.profiles p
      join public.team_members tm on tm.account_owner_id = projects.owner_id
      where p.id = projects.owner_id
        and (p.subscription_tier = 'pro' or p.subscription_tier = 'team')
        and tm.user_id = auth.uid()
        and tm.status = 'active'
    )
  );

-- Fix projects UPDATE policy
drop policy if exists "Project owners and teammates can update projects" on public.projects;
create policy "Project owners and teammates can update projects"
  on public.projects
  for update
  using (
    owner_id = auth.uid()
    or exists (
      -- Check if user and project owner are both team members with same account_owner_id
      select 1
      from public.team_members tm1
      join public.team_members tm2 on tm1.account_owner_id = tm2.account_owner_id
      where tm1.user_id = auth.uid()
        and tm2.user_id = projects.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner (has pro/team tier) and user is their team member
      select 1
      from public.profiles p
      join public.team_members tm on tm.account_owner_id = projects.owner_id
      where p.id = projects.owner_id
        and (p.subscription_tier = 'pro' or p.subscription_tier = 'team')
        and tm.user_id = auth.uid()
        and tm.status = 'active'
    )
  );

-- Fix projects DELETE policy
drop policy if exists "Project owners and teammates can delete projects" on public.projects;
create policy "Project owners and teammates can delete projects"
  on public.projects
  for delete
  using (
    owner_id = auth.uid()
    or exists (
      -- Check if user and project owner are both team members with same account_owner_id
      select 1
      from public.team_members tm1
      join public.team_members tm2 on tm1.account_owner_id = tm2.account_owner_id
      where tm1.user_id = auth.uid()
        and tm2.user_id = projects.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner (has pro/team tier) and user is their team member
      select 1
      from public.profiles p
      join public.team_members tm on tm.account_owner_id = projects.owner_id
      where p.id = projects.owner_id
        and (p.subscription_tier = 'pro' or p.subscription_tier = 'team')
        and tm.user_id = auth.uid()
        and tm.status = 'active'
    )
  );

-- Fix crawls SELECT policy
drop policy if exists "Project members and teammates can view crawls" on public.crawls;
create policy "Project members and teammates can view crawls"
  on public.crawls
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = crawls.project_id
        and pm.user_id = auth.uid()
    )
    or exists (
      -- Check if user and project owner are both team members with same account_owner_id
      select 1
      from public.projects p
      join public.team_members tm1 on tm1.user_id = auth.uid()
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where p.id = crawls.project_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner (has pro/team tier) and user is their team member
      select 1
      from public.projects p
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = crawls.project_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = auth.uid()
        and tm.status = 'active'
    )
  );

-- Fix crawls INSERT policy
drop policy if exists "Project members and teammates can create crawls" on public.crawls;
create policy "Project members and teammates can create crawls"
  on public.crawls
  for insert
  with check (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = crawls.project_id
        and pm.user_id = auth.uid()
    )
    or exists (
      -- Check if user and project owner are both team members with same account_owner_id
      select 1
      from public.projects p
      join public.team_members tm1 on tm1.user_id = auth.uid()
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where p.id = crawls.project_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner (has pro/team tier) and user is their team member
      select 1
      from public.projects p
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = crawls.project_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = auth.uid()
        and tm.status = 'active'
    )
  );

-- Fix crawls UPDATE policy
drop policy if exists "Project members and teammates can update crawls" on public.crawls;
create policy "Project members and teammates can update crawls"
  on public.crawls
  for update
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = crawls.project_id
        and pm.user_id = auth.uid()
    )
    or exists (
      -- Check if user and project owner are both team members with same account_owner_id
      select 1
      from public.projects p
      join public.team_members tm1 on tm1.user_id = auth.uid()
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where p.id = crawls.project_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner (has pro/team tier) and user is their team member
      select 1
      from public.projects p
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = crawls.project_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = auth.uid()
        and tm.status = 'active'
    )
  );

-- Fix crawls DELETE policy
drop policy if exists "Project members and teammates can delete crawls" on public.crawls;
create policy "Project members and teammates can delete crawls"
  on public.crawls
  for delete
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = crawls.project_id
        and pm.user_id = auth.uid()
    )
    or exists (
      -- Check if user and project owner are both team members with same account_owner_id
      select 1
      from public.projects p
      join public.team_members tm1 on tm1.user_id = auth.uid()
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where p.id = crawls.project_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner (has pro/team tier) and user is their team member
      select 1
      from public.projects p
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = crawls.project_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = auth.uid()
        and tm.status = 'active'
    )
  );

-- Fix API integrations SELECT policy
drop policy if exists "Project members and teammates can view integrations" on public.api_integrations;
create policy "Project members and teammates can view integrations"
  on public.api_integrations
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = api_integrations.project_id
        and pm.user_id = auth.uid()
    )
    or exists (
      -- Check if user and project owner are both team members with same account_owner_id
      select 1
      from public.projects p
      join public.team_members tm1 on tm1.user_id = auth.uid()
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where p.id = api_integrations.project_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner (has pro/team tier) and user is their team member
      select 1
      from public.projects p
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = api_integrations.project_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = auth.uid()
        and tm.status = 'active'
    )
  );

-- Fix API integrations ALL (INSERT/UPDATE/DELETE) policy
drop policy if exists "Project owners and teammates can manage integrations" on public.api_integrations;
create policy "Project owners and teammates can manage integrations"
  on public.api_integrations
  for all
  using (
    exists (
      select 1
      from public.projects p
      where p.id = api_integrations.project_id
        and p.owner_id = auth.uid()
    )
    or exists (
      -- Check if user and project owner are both team members with same account_owner_id
      select 1
      from public.projects p
      join public.team_members tm1 on tm1.user_id = auth.uid()
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where p.id = api_integrations.project_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner (has pro/team tier) and user is their team member
      select 1
      from public.projects p
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = api_integrations.project_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = auth.uid()
        and tm.status = 'active'
    )
  );

