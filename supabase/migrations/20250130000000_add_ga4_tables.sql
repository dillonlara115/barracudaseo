-- Google Analytics 4 (GA4) integration schema
-- Similar structure to GSC integration for consistency

create table if not exists public.ga4_sync_states (
  project_id uuid primary key references public.projects (id) on delete cascade,
  property_id text,
  property_name text,
  status text check (status in ('idle', 'running', 'error')) default 'idle',
  last_synced_at timestamptz,
  error_log jsonb,
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);

create table if not exists public.ga4_performance_snapshots (
  id uuid primary key default gen_random_uuid(),
  project_id uuid not null references public.projects (id) on delete cascade,
  property_id text not null,
  captured_on date not null,
  period text not null,
  totals jsonb not null default '{}'::jsonb,
  created_at timestamptz default now()
);

create index if not exists idx_ga4_performance_snapshots_project_captured
  on public.ga4_performance_snapshots (project_id, captured_on desc);

create table if not exists public.ga4_performance_rows (
  id bigserial primary key,
  snapshot_id uuid not null references public.ga4_performance_snapshots (id) on delete cascade,
  project_id uuid not null references public.projects (id) on delete cascade,
  row_type text check (row_type in ('page', 'source', 'medium', 'device', 'country', 'date')) not null,
  dimension_value text not null,
  metrics jsonb not null default '{}'::jsonb,
  created_at timestamptz default now()
);

create unique index if not exists idx_ga4_performance_rows_unique
  on public.ga4_performance_rows (snapshot_id, row_type, dimension_value);

create index if not exists idx_ga4_performance_rows_lookup
  on public.ga4_performance_rows (project_id, row_type, dimension_value);

-- Row Level Security policies

alter table public.ga4_sync_states enable row level security;
alter table public.ga4_performance_snapshots enable row level security;
alter table public.ga4_performance_rows enable row level security;

create policy "Project members can view ga4 sync state"
  on public.ga4_sync_states
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = ga4_sync_states.project_id
        and pm.user_id = auth.uid()
    )
    or exists (
      select 1
      from public.projects p
      where p.id = ga4_sync_states.project_id
        and p.owner_id = auth.uid()
    )
  );

create policy "Project members can view ga4 snapshots"
  on public.ga4_performance_snapshots
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = ga4_performance_snapshots.project_id
        and pm.user_id = auth.uid()
    )
    or exists (
      select 1
      from public.projects p
      where p.id = ga4_performance_snapshots.project_id
        and p.owner_id = auth.uid()
    )
  );

create policy "Project members can view ga4 performance rows"
  on public.ga4_performance_rows
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = ga4_performance_rows.project_id
        and pm.user_id = auth.uid()
    )
    or exists (
      select 1
      from public.projects p
      where p.id = ga4_performance_rows.project_id
        and p.owner_id = auth.uid()
    )
  );

-- Grants
grant select on public.ga4_sync_states to authenticated;
grant select on public.ga4_performance_snapshots to authenticated;
grant select on public.ga4_performance_rows to authenticated;

-- Updated_at trigger
create trigger set_updated_at_ga4_sync_states
  before update on public.ga4_sync_states
  for each row
  execute function public.handle_updated_at();
