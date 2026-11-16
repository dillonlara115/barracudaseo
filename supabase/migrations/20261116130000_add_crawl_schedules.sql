-- Scheduled (recurring) crawls: one schedule per project.
--
-- A cron job calls POST /api/internal/crawls/run-scheduled, which starts a
-- crawl for every enabled schedule whose next_run_at has passed and then moves
-- next_run_at forward by the schedule's frequency.

create table if not exists public.crawl_schedules (
  id uuid primary key default gen_random_uuid(),
  project_id uuid not null references public.projects (id) on delete cascade,
  -- Whose subscription limits apply when the crawl runs
  created_by uuid not null references auth.users (id) on delete cascade,
  frequency text check (frequency in ('daily', 'weekly', 'monthly')) not null,
  enabled boolean not null default true,
  -- Crawl options: url, max_pages, max_depth, respect_robots, parse_sitemap, crawl_sitemap_only
  config jsonb not null default '{}'::jsonb,
  next_run_at timestamptz not null,
  last_run_at timestamptz,
  last_crawl_id uuid references public.crawls (id) on delete set null,
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);

alter table public.crawl_schedules enable row level security;

-- One schedule per project
create unique index if not exists idx_crawl_schedules_project_id
  on public.crawl_schedules (project_id);

-- The cron job looks up due schedules
create index if not exists idx_crawl_schedules_due
  on public.crawl_schedules (next_run_at)
  where enabled;

-- Anyone with access to the project can see its schedule. Writes go through
-- the API server (service role), which checks project access and plan limits.
drop policy if exists "Project members can view crawl schedules" on public.crawl_schedules;
create policy "Project members can view crawl schedules"
  on public.crawl_schedules for select
  using (public.can_access_project(project_id));

drop trigger if exists set_updated_at_crawl_schedules on public.crawl_schedules;
create trigger set_updated_at_crawl_schedules
  before update on public.crawl_schedules
  for each row
  execute function public.handle_updated_at();
