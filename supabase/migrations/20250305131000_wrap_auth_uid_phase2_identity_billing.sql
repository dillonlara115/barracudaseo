-- Phase 2: identity/billing tables — wrap auth.uid() in SELECT for profiles, subscriptions, team_members

-- Profiles
drop policy if exists "Users can view their own profile" on public.profiles;
drop policy if exists "Users can update their own profile" on public.profiles;
drop policy if exists "Users can create their own profile" on public.profiles;

create policy "Users can view their own profile"
  on public.profiles
  for select
  using ((select auth.uid()) = id);

create policy "Users can update their own profile"
  on public.profiles
  for update
  using ((select auth.uid()) = id);

create policy "Users can create their own profile"
  on public.profiles
  for insert
  with check ((select auth.uid()) = id);

-- Subscriptions
drop policy if exists "Users can view their own subscriptions" on public.subscriptions;

create policy "Users can view their own subscriptions"
  on public.subscriptions
  for select
  using (user_id = (select auth.uid()));

-- Team members: reapply SELECT-wrapped versions (kept aligned with previous migration)
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
