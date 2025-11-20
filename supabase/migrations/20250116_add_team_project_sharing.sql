-- Migration: Add team-based project sharing
-- Team members can now access projects created by other team members
-- This enables team-wide collaboration on projects, crawls, and integrations

-- Drop existing project policies and recreate with team support
drop policy if exists "Project members can view projects" on public.projects;
drop policy if exists "Project owners can update projects" on public.projects;
drop policy if exists "Project owners can delete projects" on public.projects;

-- Updated policy: Users can view projects if:
-- 1. They are the project owner
-- 2. They are a project member
-- 3. They are on the same team as the project owner (team members share access)
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
      -- Check if user and project owner are on the same team
      select 1
      from public.team_members tm1
      join public.team_members tm2 on tm1.account_owner_id = tm2.account_owner_id
      where tm1.user_id = auth.uid()
        and tm2.user_id = projects.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      -- Check if user is account owner and project owner is their team member
      select 1
      from public.team_members tm
      where tm.account_owner_id = auth.uid()
        and tm.user_id = projects.owner_id
        and tm.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner and user is their team member
      select 1
      from public.team_members tm
      where tm.account_owner_id = projects.owner_id
        and tm.user_id = auth.uid()
        and tm.status = 'active'
    )
  );

-- Updated policy: Users can update projects if:
-- 1. They are the project owner
-- 2. They are on the same team as the project owner (team members can edit)
create policy "Project owners and teammates can update projects"
  on public.projects
  for update
  using (
    owner_id = auth.uid()
    or exists (
      -- Check if user and project owner are on the same team
      select 1
      from public.team_members tm1
      join public.team_members tm2 on tm1.account_owner_id = tm2.account_owner_id
      where tm1.user_id = auth.uid()
        and tm2.user_id = projects.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      -- Check if user is account owner and project owner is their team member
      select 1
      from public.team_members tm
      where tm.account_owner_id = auth.uid()
        and tm.user_id = projects.owner_id
        and tm.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner and user is their team member
      select 1
      from public.team_members tm
      where tm.account_owner_id = projects.owner_id
        and tm.user_id = auth.uid()
        and tm.status = 'active'
    )
  );

-- Updated policy: Users can delete projects if:
-- 1. They are the project owner
-- 2. They are on the same team as the project owner (team members can delete)
create policy "Project owners and teammates can delete projects"
  on public.projects
  for delete
  using (
    owner_id = auth.uid()
    or exists (
      -- Check if user and project owner are on the same team
      select 1
      from public.team_members tm1
      join public.team_members tm2 on tm1.account_owner_id = tm2.account_owner_id
      where tm1.user_id = auth.uid()
        and tm2.user_id = projects.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      -- Check if user is account owner and project owner is their team member
      select 1
      from public.team_members tm
      where tm.account_owner_id = auth.uid()
        and tm.user_id = projects.owner_id
        and tm.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner and user is their team member
      select 1
      from public.team_members tm
      where tm.account_owner_id = projects.owner_id
        and tm.user_id = auth.uid()
        and tm.status = 'active'
    )
  );

-- Update crawls policy to include team access
drop policy if exists "Project members can view crawls" on public.crawls;
drop policy if exists "Project members can create crawls" on public.crawls;
drop policy if exists "Project members can update crawls" on public.crawls;

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
      -- Check if user and project owner are on the same team
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
      -- Check if user is account owner and project owner is their team member
      select 1
      from public.projects p
      join public.team_members tm on tm.account_owner_id = auth.uid()
      where p.id = crawls.project_id
        and tm.user_id = p.owner_id
        and tm.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner and user is their team member
      select 1
      from public.projects p
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = crawls.project_id
        and tm.user_id = auth.uid()
        and tm.status = 'active'
    )
  );

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
      -- Check if user and project owner are on the same team
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
      -- Check if user is account owner and project owner is their team member
      select 1
      from public.projects p
      join public.team_members tm on tm.account_owner_id = auth.uid()
      where p.id = crawls.project_id
        and tm.user_id = p.owner_id
        and tm.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner and user is their team member
      select 1
      from public.projects p
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = crawls.project_id
        and tm.user_id = auth.uid()
        and tm.status = 'active'
    )
  );

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
      -- Check if user and project owner are on the same team
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
      -- Check if user is account owner and project owner is their team member
      select 1
      from public.projects p
      join public.team_members tm on tm.account_owner_id = auth.uid()
      where p.id = crawls.project_id
        and tm.user_id = p.owner_id
        and tm.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner and user is their team member
      select 1
      from public.projects p
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = crawls.project_id
        and tm.user_id = auth.uid()
        and tm.status = 'active'
    )
  );

-- Update API integrations policy to allow team members to view integrations
drop policy if exists "Project members can view integrations" on public.api_integrations;
drop policy if exists "Project owners can manage integrations" on public.api_integrations;

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
      -- Check if user and project owner are on the same team
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
      -- Check if user is account owner and project owner is their team member
      select 1
      from public.projects p
      join public.team_members tm on tm.account_owner_id = auth.uid()
      where p.id = api_integrations.project_id
        and tm.user_id = p.owner_id
        and tm.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner and user is their team member
      select 1
      from public.projects p
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = api_integrations.project_id
        and tm.user_id = auth.uid()
        and tm.status = 'active'
    )
  );

-- Team members can also manage integrations (connect/disconnect GSC, etc.)
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
      -- Check if user and project owner are on the same team
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
      -- Check if user is account owner and project owner is their team member
      select 1
      from public.projects p
      join public.team_members tm on tm.account_owner_id = auth.uid()
      where p.id = api_integrations.project_id
        and tm.user_id = p.owner_id
        and tm.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner and user is their team member
      select 1
      from public.projects p
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = api_integrations.project_id
        and tm.user_id = auth.uid()
        and tm.status = 'active'
    )
  );

