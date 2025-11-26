-- Add public_reports table for shareable client reports

-- Public Reports
-- Stores metadata for public-facing reports that can be shared with clients without authentication
create table if not exists public.public_reports (
  id uuid primary key default gen_random_uuid(),
  crawl_id uuid not null references public.crawls(id) on delete cascade,
  project_id uuid not null references public.projects(id) on delete cascade,
  created_by uuid not null references auth.users(id) on delete cascade,
  access_token text not null unique, -- Secure token for accessing the report
  password_hash text, -- Optional password protection (bcrypt hash)
  expires_at timestamptz, -- Optional expiry date
  title text, -- Custom report title
  description text, -- Optional report description
  settings jsonb default '{}'::jsonb, -- Report settings (include AI summary, filters, etc.)
  view_count integer default 0, -- Track how many times the report has been viewed
  last_viewed_at timestamptz, -- Last time the report was viewed
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);

-- Indexes for public_reports
create index if not exists idx_public_reports_access_token on public.public_reports (access_token);
create index if not exists idx_public_reports_crawl on public.public_reports (crawl_id);
create index if not exists idx_public_reports_project on public.public_reports (project_id);
create index if not exists idx_public_reports_created_by on public.public_reports (created_by);

-- Enable RLS
alter table public.public_reports enable row level security;

-- RLS Policy: Users can only manage their own reports
-- Note: Public access is handled via the access_token, not through RLS
create policy "Users can view their own public reports"
  on public.public_reports
  for select
  using (
    auth.uid() = created_by
    and exists (
      select 1
      from public.project_members pm
      where pm.project_id = public_reports.project_id
        and pm.user_id = auth.uid()
    )
  );

create policy "Users can create public reports for their projects"
  on public.public_reports
  for insert
  with check (
    auth.uid() = created_by
    and exists (
      select 1
      from public.project_members pm
      where pm.project_id = public_reports.project_id
        and pm.user_id = auth.uid()
    )
  );

create policy "Users can update their own public reports"
  on public.public_reports
  for update
  using (
    auth.uid() = created_by
    and exists (
      select 1
      from public.project_members pm
      where pm.project_id = public_reports.project_id
        and pm.user_id = auth.uid()
    )
  );

create policy "Users can delete their own public reports"
  on public.public_reports
  for delete
  using (
    auth.uid() = created_by
    and exists (
      select 1
      from public.project_members pm
      where pm.project_id = public_reports.project_id
        and pm.user_id = auth.uid()
    )
  );



