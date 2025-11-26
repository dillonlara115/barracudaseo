-- Phase 2: Scheduled Checks, Usage Tracking, and Crawl Integration
-- Reference: docs/DATAFORSEO_INTEGRATION.md

-- Add scheduling fields to keywords table
alter table public.keywords
  add column if not exists check_frequency text default 'manual', -- manual | daily | weekly
  add column if not exists last_checked_at timestamptz,
  add column if not exists next_check_at timestamptz;

-- Index for finding keywords that need checking
create index if not exists idx_keywords_next_check_at 
  on public.keywords (next_check_at) 
  where check_frequency != 'manual' and next_check_at is not null;

-- Keyword usage tracking table
create table if not exists public.keyword_usage (
  id uuid primary key default gen_random_uuid(),
  project_id uuid not null references public.projects(id) on delete cascade,
  keyword_id uuid references public.keywords(id) on delete set null,
  user_id uuid not null references auth.users(id) on delete cascade,
  check_type text not null default 'manual', -- manual | scheduled
  dataforseo_task_id text,
  cost_usd numeric(10, 6) default 0.001, -- Default cost per check (~$0.001)
  checked_at timestamptz not null default now(),
  created_at timestamptz not null default now()
);

create index if not exists idx_keyword_usage_project_id on public.keyword_usage (project_id);
create index if not exists idx_keyword_usage_user_id on public.keyword_usage (user_id);
create index if not exists idx_keyword_usage_keyword_id on public.keyword_usage (keyword_id);
create index if not exists idx_keyword_usage_checked_at on public.keyword_usage (checked_at desc);

-- Add crawl_page_id to keyword_rank_snapshots for linking to crawl pages
-- Note: pages.id is bigserial (bigint), not uuid
alter table public.keyword_rank_snapshots
  add column if not exists crawl_page_id bigint references public.pages(id) on delete set null;

create index if not exists idx_keyword_snapshots_crawl_page_id 
  on public.keyword_rank_snapshots (crawl_page_id);

-- RLS Policies for keyword_usage
alter table public.keyword_usage enable row level security;

create policy "Project members can view keyword usage"
  on public.keyword_usage
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = keyword_usage.project_id
        and pm.user_id = auth.uid()
    )
    or exists (
      select 1
      from public.projects p
      where p.id = keyword_usage.project_id
        and p.owner_id = auth.uid()
    )
  );

create policy "System can insert keyword usage"
  on public.keyword_usage
  for insert
  with check (true); -- System inserts via service role

-- Grants
grant select on public.keyword_usage to authenticated;
grant insert on public.keyword_usage to authenticated;

-- Function to calculate next check time based on frequency
create or replace function public.calculate_next_check_at(
  frequency text,
  last_checked timestamptz default now()
)
returns timestamptz as $$
begin
  case frequency
    when 'daily' then
      return (last_checked + interval '1 day');
    when 'weekly' then
      return (last_checked + interval '7 days');
    else
      return null; -- manual
  end case;
end;
$$ language plpgsql;

-- Function to update next_check_at when check_frequency or last_checked_at changes
create or replace function public.update_keyword_next_check()
returns trigger as $$
begin
  if NEW.check_frequency != 'manual' and NEW.last_checked_at is not null then
    NEW.next_check_at := public.calculate_next_check_at(NEW.check_frequency, NEW.last_checked_at);
  elsif NEW.check_frequency = 'manual' then
    NEW.next_check_at := null;
  end if;
  return NEW;
end;
$$ language plpgsql;

-- Trigger to auto-update next_check_at
create trigger update_keyword_next_check_trigger
  before insert or update of check_frequency, last_checked_at on public.keywords
  for each row
  execute function public.update_keyword_next_check();

