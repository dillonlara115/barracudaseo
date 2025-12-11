-- Fix RLS recursion for project_members while keeping auth.uid() wrapped
-- Replaces recursive policies introduced in 20250305124500_wrap_auth_uid_hotpath.sql
-- Aligns with docs/RLS_OPTIMIZATION_CHECKLIST.md (Phase 1)

-- Ensure helper exists (created previously in 20240321_fix_project_members_rls.sql):
--   public.is_project_member(project_uuid uuid, user_uuid uuid) SECURITY DEFINER

-- Drop potentially recursive policies
drop policy if exists "Project members can view project members" on public.project_members;
drop policy if exists "Project owners can add members" on public.project_members;
drop policy if exists "Project owners can update members" on public.project_members;
drop policy if exists "Project owners can remove members" on public.project_members;

-- Safe SELECT policy using the helper to avoid recursion
create policy "Project members can view project members"
  on public.project_members
  for select
  using (
    public.is_project_member(project_members.project_id, (select auth.uid()))
  );

-- Ownership checks go through projects (non-recursive) with wrapped auth.uid()
create policy "Project owners can add members"
  on public.project_members
  for insert
  with check (
    exists (
      select 1
      from public.projects p
      where p.id = project_members.project_id
        and p.owner_id = (select auth.uid())
    )
  );

create policy "Project owners can update members"
  on public.project_members
  for update
  using (
    exists (
      select 1
      from public.projects p
      where p.id = project_members.project_id
        and p.owner_id = (select auth.uid())
    )
  );

create policy "Project owners can remove members"
  on public.project_members
  for delete
  using (
    exists (
      select 1
      from public.projects p
      where p.id = project_members.project_id
        and p.owner_id = (select auth.uid())
    )
    or project_members.user_id = (select auth.uid())
  );
