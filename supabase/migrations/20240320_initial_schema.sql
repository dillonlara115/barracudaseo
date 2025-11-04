-- Enable extensions
-- Reference: docs/SUPABASE_SCHEMA.md - Supporting Objects section

create extension if not exists "pgcrypto";

-- Profiles
-- Reference: docs/SUPABASE_SCHEMA.md - Table Definitions section 1

create table if not exists public.profiles (
  id uuid primary key references auth.users (id) on delete cascade,
  display_name text,
  avatar_url text,
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);

-- Projects
-- Reference: docs/SUPABASE_SCHEMA.md - Table Definitions section 2

create table if not exists public.projects (
  id uuid primary key default gen_random_uuid(),
  name text not null,
  domain text not null,
  owner_id uuid not null references auth.users (id) on delete cascade,
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  settings jsonb default '{}'::jsonb
);

-- Unique constraint: prevent duplicate domains per owner
create unique index if not exists idx_projects_owner_domain on public.projects (owner_id, lower(domain));

-- Project Members
-- Reference: docs/SUPABASE_SCHEMA.md - Table Definitions section 3

create table if not exists public.project_members (
  project_id uuid references public.projects (id) on delete cascade,
  user_id uuid references auth.users (id) on delete cascade,
  role text check (role in ('owner', 'editor', 'viewer')) default 'viewer',
  invited_by uuid references auth.users (id),
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  primary key (project_id, user_id)
);

-- Crawls
-- Reference: docs/SUPABASE_SCHEMA.md - Table Definitions section 4

create table if not exists public.crawls (
  id uuid primary key default gen_random_uuid(),
  project_id uuid not null references public.projects (id) on delete cascade,
  initiated_by uuid references auth.users (id),
  source text check (source in ('cli', 'web', 'schedule')) default 'cli',
  status text check (status in ('pending', 'running', 'succeeded', 'failed', 'cancelled')) not null,
  started_at timestamptz default now(),
  completed_at timestamptz,
  total_pages integer default 0,
  total_issues integer default 0,
  meta jsonb default '{}'::jsonb
);

-- Indexes for crawls
create index if not exists idx_crawls_project_started on public.crawls (project_id, started_at desc);
create index if not exists idx_crawls_status on public.crawls (project_id, status);

-- Pages
-- Reference: docs/SUPABASE_SCHEMA.md - Table Definitions section 5

create table if not exists public.pages (
  id bigserial primary key,
  crawl_id uuid not null references public.crawls (id) on delete cascade,
  url text not null,
  status_code integer,
  response_time_ms integer,
  title text,
  meta_description text,
  canonical_url text,
  h1 text,
  word_count integer,
  content_hash text,
  screenshot_url text,
  data jsonb default '{}'::jsonb,
  created_at timestamptz default now()
);

-- Indexes for pages
create unique index if not exists idx_pages_crawl_url on public.pages (crawl_id, url);

-- Issues
-- Reference: docs/SUPABASE_SCHEMA.md - Table Definitions section 6

create table if not exists public.issues (
  id bigserial primary key,
  crawl_id uuid not null references public.crawls (id) on delete cascade,
  page_id bigint references public.pages (id) on delete set null,
  project_id uuid references public.projects (id) on delete cascade,
  type text not null,
  severity text check (severity in ('error', 'warning', 'info')) not null,
  message text not null,
  recommendation text,
  value text,
  priority_score integer,
  status text check (status in ('new', 'in_progress', 'fixed', 'ignored')) default 'new',
  status_updated_at timestamptz default now(),
  created_at timestamptz default now()
);

-- Indexes for issues
create index if not exists idx_issues_crawl_type on public.issues (crawl_id, type);
create index if not exists idx_issues_project_status on public.issues (project_id, status);
create index if not exists idx_issues_page on public.issues (page_id);

-- Issue Recommendations
-- Reference: docs/SUPABASE_SCHEMA.md - Table Definitions section 7

create table if not exists public.issue_recommendations (
  id bigserial primary key,
  issue_id bigint references public.issues (id) on delete cascade,
  author_type text check (author_type in ('ai', 'user', 'system')) default 'ai',
  author_id uuid references auth.users (id),
  summary text not null,
  details text,
  created_at timestamptz default now()
);

-- Indexes for issue_recommendations
create index if not exists idx_issue_recommendations_issue on public.issue_recommendations (issue_id);

-- Issue Status History
-- Reference: docs/SUPABASE_SCHEMA.md - Table Definitions section 8

create table if not exists public.issue_status_history (
  id bigserial primary key,
  issue_id bigint references public.issues (id) on delete cascade,
  old_status text,
  new_status text,
  changed_by uuid references auth.users (id),
  notes text,
  changed_at timestamptz default now()
);

-- Indexes for issue_status_history
create index if not exists idx_issue_status_history_issue on public.issue_status_history (issue_id);

-- Exports
-- Reference: docs/SUPABASE_SCHEMA.md - Table Definitions section 9

create table if not exists public.exports (
  id uuid primary key default gen_random_uuid(),
  project_id uuid references public.projects (id) on delete cascade,
  crawl_id uuid references public.crawls (id) on delete set null,
  type text check (type in ('csv', 'json', 'pdf', 'html')),
  storage_path text not null,
  requested_by uuid references auth.users (id),
  status text check (status in ('queued', 'processing', 'ready', 'failed')),
  created_at timestamptz default now(),
  completed_at timestamptz
);

-- Indexes for exports
create index if not exists idx_exports_project_created on public.exports (project_id, created_at desc);

-- API Integrations
-- Reference: docs/SUPABASE_SCHEMA.md - Table Definitions section 10

create table if not exists public.api_integrations (
  id uuid primary key default gen_random_uuid(),
  project_id uuid references public.projects (id) on delete cascade,
  provider text check (provider in ('gsc', 'openai', 'pagespeed')) not null,
  config jsonb not null,
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  unique (project_id, provider)
);

-- Helper Functions
-- Reference: docs/SUPABASE_SCHEMA.md - Supporting Objects section

-- Function to ensure project membership
create or replace function public.ensure_project_membership(project_uuid uuid)
returns boolean
language plpgsql
security definer
as $$
begin
  if not exists (
    select 1
    from public.project_members pm
    where pm.project_id = project_uuid
      and pm.user_id = auth.uid()
  ) then
    raise exception 'User does not have access to this project';
  end if;
  return true;
end;
$$;

-- Row Level Security (RLS)
-- Reference: docs/SUPABASE_SCHEMA.md - Row Level Security section

-- Enable RLS on all tables
alter table public.profiles enable row level security;
alter table public.projects enable row level security;
alter table public.project_members enable row level security;
alter table public.crawls enable row level security;
alter table public.pages enable row level security;
alter table public.issues enable row level security;
alter table public.issue_recommendations enable row level security;
alter table public.issue_status_history enable row level security;
alter table public.exports enable row level security;
alter table public.api_integrations enable row level security;

-- Profiles RLS Policies
create policy "Users can view their own profile"
  on public.profiles
  for select
  using (auth.uid() = id);

create policy "Users can update their own profile"
  on public.profiles
  for update
  using (auth.uid() = id);

-- Projects RLS Policies
create policy "Project members can view projects"
  on public.projects
  for select
  using (
    owner_id = auth.uid()
    or exists (
      select 1
      from public.project_members pm
      where pm.project_id = projects.id
        and pm.user_id = auth.uid()
    )
  );

create policy "Project owners can update projects"
  on public.projects
  for update
  using (owner_id = auth.uid());

create policy "Authenticated users can create projects"
  on public.projects
  for insert
  with check (auth.uid() = owner_id);

create policy "Project owners can delete projects"
  on public.projects
  for delete
  using (owner_id = auth.uid());

-- Project Members RLS Policies
create policy "Project members can view project members"
  on public.project_members
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = project_members.project_id
        and pm.user_id = auth.uid()
    )
  );

create policy "Project owners can manage members"
  on public.project_members
  for all
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = project_members.project_id
        and pm.user_id = auth.uid()
        and pm.role = 'owner'
    )
  );

-- Crawls RLS Policies
create policy "Project members can view crawls"
  on public.crawls
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = crawls.project_id
        and pm.user_id = auth.uid()
    )
  );

create policy "Project members can create crawls"
  on public.crawls
  for insert
  with check (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = crawls.project_id
        and pm.user_id = auth.uid()
    )
  );

-- Pages RLS Policies
create policy "Project members can view pages"
  on public.pages
  for select
  using (
    exists (
      select 1
      from public.crawls c
      join public.project_members pm on pm.project_id = c.project_id
      where c.id = pages.crawl_id
        and pm.user_id = auth.uid()
    )
  );

create policy "Project members can insert pages"
  on public.pages
  for insert
  with check (
    exists (
      select 1
      from public.crawls c
      join public.project_members pm on pm.project_id = c.project_id
      where c.id = pages.crawl_id
        and pm.user_id = auth.uid()
    )
  );

-- Issues RLS Policies
create policy "Project members can view issues"
  on public.issues
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = issues.project_id
        and pm.user_id = auth.uid()
    )
  );

create policy "Project members can update issue status"
  on public.issues
  for update
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = issues.project_id
        and pm.user_id = auth.uid()
        and pm.role in ('editor', 'owner')
    )
  );

create policy "Project members can insert issues"
  on public.issues
  for insert
  with check (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = issues.project_id
        and pm.user_id = auth.uid()
    )
  );

-- Issue Recommendations RLS Policies
create policy "Project members can view recommendations"
  on public.issue_recommendations
  for select
  using (
    exists (
      select 1
      from public.issues i
      join public.project_members pm on pm.project_id = i.project_id
      where i.id = issue_recommendations.issue_id
        and pm.user_id = auth.uid()
    )
  );

create policy "Project members can create recommendations"
  on public.issue_recommendations
  for insert
  with check (
    exists (
      select 1
      from public.issues i
      join public.project_members pm on pm.project_id = i.project_id
      where i.id = issue_recommendations.issue_id
        and pm.user_id = auth.uid()
    )
  );

-- Issue Status History RLS Policies
create policy "Project members can view status history"
  on public.issue_status_history
  for select
  using (
    exists (
      select 1
      from public.issues i
      join public.project_members pm on pm.project_id = i.project_id
      where i.id = issue_status_history.issue_id
        and pm.user_id = auth.uid()
    )
  );

create policy "Project members can create status history"
  on public.issue_status_history
  for insert
  with check (
    exists (
      select 1
      from public.issues i
      join public.project_members pm on pm.project_id = i.project_id
      where i.id = issue_status_history.issue_id
        and pm.user_id = auth.uid()
    )
  );

-- Exports RLS Policies
create policy "Project members can view exports"
  on public.exports
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = exports.project_id
        and pm.user_id = auth.uid()
    )
  );

create policy "Project members can create exports"
  on public.exports
  for insert
  with check (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = exports.project_id
        and pm.user_id = auth.uid()
    )
  );

create policy "Export requesters and owners can delete exports"
  on public.exports
  for delete
  using (
    requested_by = auth.uid()
    or exists (
      select 1
      from public.projects p
      where p.id = exports.project_id
        and p.owner_id = auth.uid()
    )
  );

-- API Integrations RLS Policies
create policy "Project members can view integrations"
  on public.api_integrations
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = api_integrations.project_id
        and pm.user_id = auth.uid()
    )
  );

create policy "Project owners can manage integrations"
  on public.api_integrations
  for all
  using (
    exists (
      select 1
      from public.projects p
      where p.id = api_integrations.project_id
        and p.owner_id = auth.uid()
    )
  );

-- Views
-- Reference: docs/SUPABASE_SCHEMA.md - Supporting Objects section

-- Project issue summary view
create or replace view public.project_issue_summary as
select 
  p.id as project_id,
  p.name as project_name,
  count(distinct i.id) as total_issues,
  count(distinct case when i.status = 'new' then i.id end) as new_issues,
  count(distinct case when i.status = 'in_progress' then i.id end) as in_progress_issues,
  count(distinct case when i.status = 'fixed' then i.id end) as fixed_issues,
  count(distinct case when i.status = 'ignored' then i.id end) as ignored_issues,
  count(distinct case when i.severity = 'error' then i.id end) as error_issues,
  count(distinct case when i.severity = 'warning' then i.id end) as warning_issues,
  count(distinct case when i.severity = 'info' then i.id end) as info_issues,
  avg(i.priority_score) as avg_priority_score
from public.projects p
left join public.issues i on i.project_id = p.id
group by p.id, p.name;

-- Latest crawl pages view
create or replace view public.latest_crawl_pages as
select distinct on (p.url, c.project_id)
  p.id,
  p.url,
  p.status_code,
  p.response_time_ms,
  p.title,
  p.meta_description,
  p.canonical_url,
  p.h1,
  p.word_count,
  p.content_hash,
  p.screenshot_url,
  p.data,
  p.created_at,
  c.project_id,
  c.id as crawl_id,
  c.started_at as crawl_started_at
from public.pages p
join public.crawls c on c.id = p.crawl_id
where c.completed_at is not null
order by p.url, c.project_id, c.started_at desc;

-- Grant access to views
grant select on public.project_issue_summary to authenticated;
grant select on public.latest_crawl_pages to authenticated;

-- Triggers for updated_at columns
-- Automatically update updated_at timestamp on row updates

create or replace function public.handle_updated_at()
returns trigger
language plpgsql
as $$
begin
  new.updated_at = now();
  return new;
end;
$$;

-- Trigger to automatically populate project_id from crawl_id
-- This replaces the generated column which PostgreSQL doesn't support with subqueries
-- The project_id is always derived from the crawl to ensure consistency
create or replace function public.set_issue_project_id()
returns trigger
language plpgsql
as $$
begin
  -- Always populate project_id from the crawl, even if provided explicitly
  -- This ensures consistency with the crawl's project
  select project_id into new.project_id
  from public.crawls
  where id = new.crawl_id;
  
  if new.project_id is null then
    raise exception 'Crawl % not found', new.crawl_id;
  end if;
  
  return new;
end;
$$;

create trigger set_issue_project_id_trigger
  before insert on public.issues
  for each row
  execute function public.set_issue_project_id();

create trigger set_updated_at_profiles
  before update on public.profiles
  for each row
  execute function public.handle_updated_at();

create trigger set_updated_at_projects
  before update on public.projects
  for each row
  execute function public.handle_updated_at();

create trigger set_updated_at_project_members
  before update on public.project_members
  for each row
  execute function public.handle_updated_at();

create trigger set_updated_at_api_integrations
  before update on public.api_integrations
  for each row
  execute function public.handle_updated_at();

