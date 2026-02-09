-- Microsoft Clarity integration schema
-- Similar structure to GSC/GA4 for consistency

-- Update user_api_integrations provider constraint to include 'clarity'
alter table public.user_api_integrations drop constraint if exists user_api_integrations_provider_check;
alter table public.user_api_integrations add constraint user_api_integrations_provider_check
  check (provider in ('gsc', 'ga4', 'clarity', 'openai', 'pagespeed'));

create table if not exists public.clarity_sync_states (
  project_id uuid primary key references public.projects (id) on delete cascade,
  clarity_project_id text,
  status text check (status in ('idle', 'running', 'error')) default 'idle',
  last_synced_at timestamptz,
  error_log jsonb,
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);

create table if not exists public.clarity_performance_snapshots (
  id uuid primary key default gen_random_uuid(),
  project_id uuid not null references public.projects (id) on delete cascade,
  clarity_project_id text not null,
  captured_on date not null,
  period text not null,
  totals jsonb not null default '{}'::jsonb,
  created_at timestamptz default now()
);

create index if not exists idx_clarity_performance_snapshots_project_captured
  on public.clarity_performance_snapshots (project_id, captured_on desc);

create table if not exists public.clarity_performance_rows (
  id bigserial primary key,
  snapshot_id uuid not null references public.clarity_performance_snapshots (id) on delete cascade,
  project_id uuid not null references public.projects (id) on delete cascade,
  row_type text check (row_type in ('url', 'device', 'browser', 'country', 'source', 'medium')) not null,
  dimension_value text not null,
  metrics jsonb not null default '{}'::jsonb,
  created_at timestamptz default now()
);

create unique index if not exists idx_clarity_performance_rows_unique
  on public.clarity_performance_rows (snapshot_id, row_type, dimension_value);

create index if not exists idx_clarity_performance_rows_lookup
  on public.clarity_performance_rows (project_id, row_type, dimension_value);

-- Row Level Security
alter table public.clarity_sync_states enable row level security;
alter table public.clarity_performance_snapshots enable row level security;
alter table public.clarity_performance_rows enable row level security;

create policy "Project members can view clarity sync state"
  on public.clarity_sync_states
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = clarity_sync_states.project_id
        and pm.user_id = auth.uid()
    )
    or exists (
      select 1
      from public.projects p
      where p.id = clarity_sync_states.project_id
        and p.owner_id = auth.uid()
    )
  );

create policy "Project members can view clarity snapshots"
  on public.clarity_performance_snapshots
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = clarity_performance_snapshots.project_id
        and pm.user_id = auth.uid()
    )
    or exists (
      select 1
      from public.projects p
      where p.id = clarity_performance_snapshots.project_id
        and p.owner_id = auth.uid()
    )
  );

create policy "Project members can view clarity performance rows"
  on public.clarity_performance_rows
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = clarity_performance_rows.project_id
        and pm.user_id = auth.uid()
    )
    or exists (
      select 1
      from public.projects p
      where p.id = clarity_performance_rows.project_id
        and p.owner_id = auth.uid()
    )
  );

-- Grants
grant select on public.clarity_sync_states to authenticated;
grant select on public.clarity_performance_snapshots to authenticated;
grant select on public.clarity_performance_rows to authenticated;

-- Updated_at trigger
create trigger set_updated_at_clarity_sync_states
  before update on public.clarity_sync_states
  for each row
  execute function public.handle_updated_at();
