-- Team Members table for account-level team management
-- This is separate from project_members which is project-specific
-- Team members are managed at the account/subscription level

create table if not exists public.team_members (
  id uuid primary key default gen_random_uuid(),
  account_owner_id uuid not null references auth.users (id) on delete cascade,
  user_id uuid references auth.users (id) on delete cascade, -- null if pending invite
  email text not null, -- email of invited user
  role text check (role in ('admin', 'member')) default 'member',
  status text check (status in ('pending', 'active', 'removed')) default 'pending',
  invited_by uuid references auth.users (id),
  invite_token text unique, -- unique token for invite acceptance
  invited_at timestamptz default now(),
  joined_at timestamptz,
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  unique(account_owner_id, email) -- prevent duplicate invites
);

-- Indexes
create index if not exists idx_team_members_account_owner on public.team_members (account_owner_id);
create index if not exists idx_team_members_user_id on public.team_members (user_id);
create index if not exists idx_team_members_invite_token on public.team_members (invite_token);
create index if not exists idx_team_members_status on public.team_members (status);

-- RLS Policies
alter table public.team_members enable row level security;

-- Account owners can view all their team members
create policy "Account owners can view their team members"
  on public.team_members
  for select
  using (account_owner_id = auth.uid());

-- Team members can view their own record
create policy "Team members can view their own record"
  on public.team_members
  for select
  using (user_id = auth.uid());

-- Account owners can manage their team members
create policy "Account owners can manage their team members"
  on public.team_members
  for all
  using (account_owner_id = auth.uid());

-- Function to check if user is account owner
create or replace function public.is_account_owner(owner_id uuid)
returns boolean
language plpgsql
security definer
as $$
begin
  return owner_id = auth.uid();
end;
$$;

-- Function to generate invite token
create or replace function public.generate_invite_token()
returns text
language plpgsql
as $$
begin
  return encode(gen_random_bytes(32), 'base64');
end;
$$;

-- Trigger for updated_at
create trigger set_updated_at_team_members
  before update on public.team_members
  for each row
  execute function public.handle_updated_at();

