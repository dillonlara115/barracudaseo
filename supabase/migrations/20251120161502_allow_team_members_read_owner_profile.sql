-- Allow team members to read their account owner's profile
-- This is needed for RLS policies that check subscription_tier

-- Add policy: Team members can view their account owner's profile
create policy "Team members can view account owner profile"
  on public.profiles
  for select
  using (
    -- User can view their own profile (existing behavior)
    auth.uid() = id
    -- OR user is a team member and this is their account owner's profile
    or exists (
      select 1
      from public.team_members tm
      where tm.user_id = auth.uid()
        and tm.account_owner_id = profiles.id
        and tm.status = 'active'
    )
  );

