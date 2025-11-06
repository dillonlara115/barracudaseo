-- Google Search Console cached data schema
-- References planning in docs/SUPABASE_SCHEMA.md and GSC integration blueprint

create table if not exists public.gsc_sync_states (
  project_id uuid primary key references public.projects (id) on delete cascade,
  property_url text,
  status text check (status in ('idle', 'running', 'error')) default 'idle',
  last_synced_at timestamptz,
  error_log jsonb,
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);

create table if not exists public.gsc_performance_snapshots (
  id uuid primary key default gen_random_uuid(),
  project_id uuid not null references public.projects (id) on delete cascade,
  property_url text not null,
  captured_on date not null,
  period text not null,
  totals jsonb not null default '{}'::jsonb,
  created_at timestamptz default now()
);

create index if not exists idx_gsc_performance_snapshots_project_captured
  on public.gsc_performance_snapshots (project_id, captured_on desc);

create table if not exists public.gsc_performance_rows (
  id bigserial primary key,
  snapshot_id uuid not null references public.gsc_performance_snapshots (id) on delete cascade,
  project_id uuid not null references public.projects (id) on delete cascade,
  row_type text check (row_type in ('query', 'page', 'country', 'device', 'appearance', 'date')) not null,
  dimension_value text not null,
  metrics jsonb not null default '{}'::jsonb,
  top_queries jsonb,
  created_at timestamptz default now()
);

create unique index if not exists idx_gsc_performance_rows_unique
  on public.gsc_performance_rows (snapshot_id, row_type, dimension_value);

create index if not exists idx_gsc_performance_rows_lookup
  on public.gsc_performance_rows (project_id, row_type, dimension_value);

create table if not exists public.gsc_page_enhancements (
  id bigserial primary key,
  project_id uuid not null references public.projects (id) on delete cascade,
  snapshot_id uuid references public.gsc_performance_snapshots (id) on delete set null,
  page_url text not null,
  enhancements jsonb default '{}'::jsonb,
  coverage jsonb default '{}'::jsonb,
  rich_results jsonb default '{}'::jsonb,
  created_at timestamptz default now()
);

create index if not exists idx_gsc_page_enhancements_project_url
  on public.gsc_page_enhancements (project_id, page_url);

create table if not exists public.gsc_insights (
  id uuid primary key default gen_random_uuid(),
  project_id uuid not null references public.projects (id) on delete cascade,
  snapshot_id uuid references public.gsc_performance_snapshots (id) on delete set null,
  insight_type text not null,
  payload jsonb not null,
  created_at timestamptz default now()
);

create index if not exists idx_gsc_insights_project_type_created
  on public.gsc_insights (project_id, insight_type, created_at desc);

-- Row Level Security policies

alter table public.gsc_sync_states enable row level security;
alter table public.gsc_performance_snapshots enable row level security;
alter table public.gsc_performance_rows enable row level security;
alter table public.gsc_page_enhancements enable row level security;
alter table public.gsc_insights enable row level security;

create policy "Project members can view gsc sync state"
  on public.gsc_sync_states
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = gsc_sync_states.project_id
        and pm.user_id = auth.uid()
    )
    or exists (
      select 1
      from public.projects p
      where p.id = gsc_sync_states.project_id
        and p.owner_id = auth.uid()
    )
  );

create policy "Project members can view gsc snapshots"
  on public.gsc_performance_snapshots
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = gsc_performance_snapshots.project_id
        and pm.user_id = auth.uid()
    )
    or exists (
      select 1
      from public.projects p
      where p.id = gsc_performance_snapshots.project_id
        and p.owner_id = auth.uid()
    )
  );

create policy "Project members can view gsc performance rows"
  on public.gsc_performance_rows
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = gsc_performance_rows.project_id
        and pm.user_id = auth.uid()
    )
    or exists (
      select 1
      from public.projects p
      where p.id = gsc_performance_rows.project_id
        and p.owner_id = auth.uid()
    )
  );

create policy "Project members can view gsc page enhancements"
  on public.gsc_page_enhancements
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = gsc_page_enhancements.project_id
        and pm.user_id = auth.uid()
    )
    or exists (
      select 1
      from public.projects p
      where p.id = gsc_page_enhancements.project_id
        and p.owner_id = auth.uid()
    )
  );

create policy "Project members can view gsc insights"
  on public.gsc_insights
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = gsc_insights.project_id
        and pm.user_id = auth.uid()
    )
    or exists (
      select 1
      from public.projects p
      where p.id = gsc_insights.project_id
        and p.owner_id = auth.uid()
    )
  );

-- Grants
grant select on public.gsc_sync_states to authenticated;
grant select on public.gsc_performance_snapshots to authenticated;
grant select on public.gsc_performance_rows to authenticated;
grant select on public.gsc_page_enhancements to authenticated;
grant select on public.gsc_insights to authenticated;

-- Updated_at trigger

create trigger set_updated_at_gsc_sync_states
  before update on public.gsc_sync_states
  for each row
  execute function public.handle_updated_at();
