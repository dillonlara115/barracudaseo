-- Consolidate API integrations SELECT policies and remove duplicate permissive rules.
-- Also split manage (insert/update/delete) into separate policies to avoid overlapping SELECT coverage.

-- Drop existing overlapping policies
drop policy if exists "Project members and teammates can view integrations" on public.api_integrations;
drop policy if exists "Project owners and teammates can manage integrations" on public.api_integrations;

-- Helper: shared project access check (owner, member, or teammate)
create or replace function public.can_access_project(project_uuid uuid, user_uuid uuid)
returns boolean
language sql
security definer
stable
set search_path = public
as $$
  select
    -- Project owner
    exists (
      select 1
      from public.projects p
      where p.id = project_uuid
        and p.owner_id = user_uuid
    )
    -- Direct project member
    or exists (
      select 1
      from public.project_members pm
      where pm.project_id = project_uuid
        and pm.user_id = user_uuid
    )
    -- Teammates under same account owner
    or exists (
      select 1
      from public.projects p
      join public.team_members tm1 on tm1.user_id = user_uuid
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where p.id = project_uuid
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    -- Account owner gives access to their team member's project (owner is account owner)
    or exists (
      select 1
      from public.projects p
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = project_uuid
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = user_uuid
        and tm.status = 'active'
    )
    -- Account owner accessing a project created by their team member (user is account owner)
    or exists (
      select 1
      from public.projects p
      join public.profiles prof on prof.id = user_uuid
      join public.team_members tm on tm.account_owner_id = user_uuid
      where p.id = project_uuid
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = p.owner_id
        and tm.status = 'active'
    );
$$;

grant execute on function public.can_access_project(uuid, uuid) to authenticated;

-- Index to support policy predicate
create index if not exists idx_api_integrations_project on public.api_integrations (project_id);

-- Single SELECT policy
create policy "API integrations select (project access)"
  on public.api_integrations
  for select
  to authenticated
  using (public.can_access_project(api_integrations.project_id, auth.uid()));

-- Separate write policies (no SELECT overlap)
create policy "API integrations insert (project access)"
  on public.api_integrations
  for insert
  to authenticated
  with check (public.can_access_project(api_integrations.project_id, auth.uid()));

create policy "API integrations update (project access)"
  on public.api_integrations
  for update
  to authenticated
  using (public.can_access_project(api_integrations.project_id, auth.uid()))
  with check (public.can_access_project(api_integrations.project_id, auth.uid()));

create policy "API integrations delete (project access)"
  on public.api_integrations
  for delete
  to authenticated
  using (public.can_access_project(api_integrations.project_id, auth.uid()));
