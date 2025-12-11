-- Normalize auth.uid() usage in team_members policies to avoid per-row re-evaluation.
-- Wrap auth.uid() in scalar subqueries per Supabase/Postgres guidance for better plans.

drop policy if exists "Users can view team members" on public.team_members;
drop policy if exists "Account owners can insert team members" on public.team_members;
drop policy if exists "Account owners can update team members" on public.team_members;
drop policy if exists "Account owners can delete team members" on public.team_members;

create policy "Users can view team members"
  on public.team_members
  for select
  to authenticated
  using (
    account_owner_id = (select auth.uid())
    or user_id = (select auth.uid())
  );

create policy "Account owners can insert team members"
  on public.team_members
  for insert
  to authenticated
  with check (
    account_owner_id = (select auth.uid())
  );

create policy "Account owners can update team members"
  on public.team_members
  for update
  to authenticated
  using (
    account_owner_id = (select auth.uid())
  );

create policy "Account owners can delete team members"
  on public.team_members
  for delete
  to authenticated
  using (
    account_owner_id = (select auth.uid())
  );
