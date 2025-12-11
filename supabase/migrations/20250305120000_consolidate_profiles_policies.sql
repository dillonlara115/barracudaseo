-- Consolidate profile SELECT policies into a single predicate and move membership
-- logic into a SECURITY DEFINER helper for better readability and performance.

drop policy if exists "Team members can view account owner profile" on public.profiles;
drop policy if exists "Users can view their own profile" on public.profiles;

-- Helper: check if current user can view the given profile (self or account owner)
create or replace function public.can_view_profile(profile_id uuid)
returns boolean
language sql
security definer
stable
set search_path = public
as $$
  select auth.uid() = profile_id
    or exists (
      select 1
      from public.team_members tm
      where tm.user_id = auth.uid()
        and tm.account_owner_id = profile_id
        and tm.status = 'active'
    );
$$;

grant execute on function public.can_view_profile(uuid) to authenticated;

create policy "Profiles select for owners and teammates"
  on public.profiles
  for select
  to authenticated
  using (public.can_view_profile(id));
