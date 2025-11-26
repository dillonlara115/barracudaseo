-- DataForSEO Rank Tracking Schema
-- Reference: docs/DATAFORSEO_INTEGRATION.md

-- Keywords table - stores keyword tracking configuration per project
create table if not exists public.keywords (
  id uuid primary key default gen_random_uuid(),
  project_id uuid not null references public.projects(id) on delete cascade,
  keyword text not null,
  target_url text, -- optional: canonical URL we want to rank
  location_name text not null, -- e.g. "United States", "Denver, Colorado"
  location_code integer,       -- optional: DataForSEO location code
  language_name text not null default 'English',
  device text not null default 'desktop', -- desktop | mobile
  search_engine text not null default 'google.com',
  tags text[] default '{}',
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

-- Optional uniqueness: no duplicate keyword+location+device in same project
create unique index if not exists keywords_project_keyword_loc_device_idx
  on public.keywords (project_id, keyword, coalesce(location_name, ''), device);

-- Index for project lookups
create index if not exists idx_keywords_project_id on public.keywords (project_id);

-- Keyword rank snapshots - each "check" for a keyword produces one snapshot record
create table if not exists public.keyword_rank_snapshots (
  id uuid primary key default gen_random_uuid(),
  keyword_id uuid not null references public.keywords(id) on delete cascade,
  checked_at timestamptz not null default now(),
  dataforseo_task_id text,      -- DataForSEO task ID
  position_absolute integer,    -- overall position in SERP
  position_organic integer,     -- organic-only position
  serp_url text,                -- ranking URL
  serp_title text,
  serp_snippet text,
  serp_features text[] default '{}', -- e.g. ['featured_snippet','sitelinks']
  search_volume integer,        -- optional: from DataForSEO stats
  rank_type text not null default 'organic', -- organic | local_pack | maps
  raw jsonb,                    -- full API response for debugging
  created_at timestamptz not null default now()
);

create index if not exists keyword_rank_snapshots_keyword_id_idx
  on public.keyword_rank_snapshots (keyword_id);

create index if not exists keyword_rank_snapshots_checked_at_idx
  on public.keyword_rank_snapshots (checked_at desc);

-- Keyword tasks - async task tracking for DataForSEO API calls
create table if not exists public.keyword_tasks (
  id uuid primary key default gen_random_uuid(),
  keyword_id uuid not null references public.keywords(id) on delete cascade,
  dataforseo_task_id text not null,
  status text not null default 'pending', -- pending | processing | completed | failed
  run_at timestamptz not null default now(), -- when we created the task
  completed_at timestamptz,
  error text,
  raw_request jsonb,
  raw_response jsonb,
  created_at timestamptz not null default now()
);

create index if not exists keyword_tasks_keyword_id_idx
  on public.keyword_tasks (keyword_id);

create index if not exists keyword_tasks_status_idx
  on public.keyword_tasks (status);

create index if not exists keyword_tasks_dataforseo_task_id_idx
  on public.keyword_tasks (dataforseo_task_id);

-- Row Level Security policies

alter table public.keywords enable row level security;
alter table public.keyword_rank_snapshots enable row level security;
alter table public.keyword_tasks enable row level security;

-- RLS Policy: Users can view keywords for projects they're members of
create policy "Project members can view keywords"
  on public.keywords
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = keywords.project_id
        and pm.user_id = auth.uid()
    )
    or exists (
      select 1
      from public.projects p
      where p.id = keywords.project_id
        and p.owner_id = auth.uid()
    )
  );

-- RLS Policy: Users can insert keywords for projects they're members of
create policy "Project members can insert keywords"
  on public.keywords
  for insert
  with check (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = keywords.project_id
        and pm.user_id = auth.uid()
    )
    or exists (
      select 1
      from public.projects p
      where p.id = keywords.project_id
        and p.owner_id = auth.uid()
    )
  );

-- RLS Policy: Users can update keywords for projects they're members of
create policy "Project members can update keywords"
  on public.keywords
  for update
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = keywords.project_id
        and pm.user_id = auth.uid()
    )
    or exists (
      select 1
      from public.projects p
      where p.id = keywords.project_id
        and p.owner_id = auth.uid()
    )
  );

-- RLS Policy: Users can delete keywords for projects they're members of
create policy "Project members can delete keywords"
  on public.keywords
  for delete
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = keywords.project_id
        and pm.user_id = auth.uid()
    )
    or exists (
      select 1
      from public.projects p
      where p.id = keywords.project_id
        and p.owner_id = auth.uid()
    )
  );

-- RLS Policy: Users can view snapshots for keywords in projects they're members of
create policy "Project members can view keyword snapshots"
  on public.keyword_rank_snapshots
  for select
  using (
    exists (
      select 1
      from public.keywords k
      join public.project_members pm on pm.project_id = k.project_id
      where k.id = keyword_rank_snapshots.keyword_id
        and pm.user_id = auth.uid()
    )
    or exists (
      select 1
      from public.keywords k
      join public.projects p on p.id = k.project_id
      where k.id = keyword_rank_snapshots.keyword_id
        and p.owner_id = auth.uid()
    )
  );

-- RLS Policy: Users can insert snapshots for keywords in projects they're members of
create policy "Project members can insert keyword snapshots"
  on public.keyword_rank_snapshots
  for insert
  with check (
    exists (
      select 1
      from public.keywords k
      join public.project_members pm on pm.project_id = k.project_id
      where k.id = keyword_rank_snapshots.keyword_id
        and pm.user_id = auth.uid()
    )
    or exists (
      select 1
      from public.keywords k
      join public.projects p on p.id = k.project_id
      where k.id = keyword_rank_snapshots.keyword_id
        and p.owner_id = auth.uid()
    )
  );

-- RLS Policy: Users can view tasks for keywords in projects they're members of
create policy "Project members can view keyword tasks"
  on public.keyword_tasks
  for select
  using (
    exists (
      select 1
      from public.keywords k
      join public.project_members pm on pm.project_id = k.project_id
      where k.id = keyword_tasks.keyword_id
        and pm.user_id = auth.uid()
    )
    or exists (
      select 1
      from public.keywords k
      join public.projects p on p.id = k.project_id
      where k.id = keyword_tasks.keyword_id
        and p.owner_id = auth.uid()
    )
  );

-- RLS Policy: Users can insert tasks for keywords in projects they're members of
create policy "Project members can insert keyword tasks"
  on public.keyword_tasks
  for insert
  with check (
    exists (
      select 1
      from public.keywords k
      join public.project_members pm on pm.project_id = k.project_id
      where k.id = keyword_tasks.keyword_id
        and pm.user_id = auth.uid()
    )
    or exists (
      select 1
      from public.keywords k
      join public.projects p on p.id = k.project_id
      where k.id = keyword_tasks.keyword_id
        and p.owner_id = auth.uid()
    )
  );

-- RLS Policy: Users can update tasks for keywords in projects they're members of
create policy "Project members can update keyword tasks"
  on public.keyword_tasks
  for update
  using (
    exists (
      select 1
      from public.keywords k
      join public.project_members pm on pm.project_id = k.project_id
      where k.id = keyword_tasks.keyword_id
        and pm.user_id = auth.uid()
    )
    or exists (
      select 1
      from public.keywords k
      join public.projects p on p.id = k.project_id
      where k.id = keyword_tasks.keyword_id
        and p.owner_id = auth.uid()
    )
  );

-- Grants
grant select, insert, update, delete on public.keywords to authenticated;
grant select, insert on public.keyword_rank_snapshots to authenticated;
grant select, insert, update on public.keyword_tasks to authenticated;

-- Updated_at trigger function (if not exists)
create or replace function public.handle_updated_at()
returns trigger as $$
begin
  new.updated_at = now();
  return new;
end;
$$ language plpgsql;

-- Updated_at trigger for keywords
create trigger set_updated_at_keywords
  before update on public.keywords
  for each row
  execute function public.handle_updated_at();

