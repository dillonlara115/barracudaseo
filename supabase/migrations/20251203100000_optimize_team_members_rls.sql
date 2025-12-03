-- Optimize RLS policies for team_members table
-- Consolidate multiple permissive policies into single efficient policies
-- Reference: User request to fix "SQL Slow queries from supabase"

-- Drop existing overlapping/redundant policies
drop policy if exists "Account owners can view their team members" on public.team_members;
drop policy if exists "Team members can view their own record" on public.team_members;
drop policy if exists "Account owners can manage their team members" on public.team_members;

-- Create consolidated SELECT policy
-- Covers both account owners and team members
create policy "Users can view team members"
  on public.team_members
  for select
  to authenticated
  using (
    account_owner_id = auth.uid() 
    or 
    user_id = auth.uid()
  );

-- Create distinct modification policies for Account Owners
-- (Formerly covered by the "manage" FOR ALL policy)

create policy "Account owners can insert team members"
  on public.team_members
  for insert
  to authenticated
  with check (
    account_owner_id = auth.uid()
  );

create policy "Account owners can update team members"
  on public.team_members
  for update
  to authenticated
  using (
    account_owner_id = auth.uid()
  );

create policy "Account owners can delete team members"
  on public.team_members
  for delete
  to authenticated
  using (
    account_owner_id = auth.uid()
  );

