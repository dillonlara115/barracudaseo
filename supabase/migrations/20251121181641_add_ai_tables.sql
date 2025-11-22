-- Add AI-related tables for user OpenAI API key management and AI-generated insights

-- User AI Settings
-- Stores user-provided OpenAI API keys (optional)
create table if not exists public.user_ai_settings (
  user_id uuid primary key references auth.users(id) on delete cascade,
  openai_api_key text,
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);

-- Enable RLS
alter table public.user_ai_settings enable row level security;

-- RLS Policy: Users can only read/update their own AI settings
create policy "Users can view their own AI settings"
  on public.user_ai_settings
  for select
  using (auth.uid() = user_id);

create policy "Users can update their own AI settings"
  on public.user_ai_settings
  for update
  using (auth.uid() = user_id);

create policy "Users can insert their own AI settings"
  on public.user_ai_settings
  for insert
  with check (auth.uid() = user_id);

-- AI Issue Insights (caching)
-- Stores AI-generated insights for individual issues
create table if not exists public.ai_issue_insights (
  id uuid primary key default gen_random_uuid(),
  issue_id bigint not null references public.issues(id) on delete cascade,
  user_id uuid not null references auth.users(id) on delete cascade,
  project_id uuid not null references public.projects(id) on delete cascade,
  crawl_id uuid not null references public.crawls(id) on delete cascade,
  insight_text text not null,
  created_at timestamptz default now()
);

-- Indexes for AI issue insights
create index if not exists idx_ai_issue_insights_issue_user on public.ai_issue_insights (issue_id, user_id);
create index if not exists idx_ai_issue_insights_crawl_user on public.ai_issue_insights (crawl_id, user_id);

-- Enable RLS
alter table public.ai_issue_insights enable row level security;

-- RLS Policy: Users can only read their own insights
-- Also check that they have access to the project
create policy "Users can view their own AI issue insights"
  on public.ai_issue_insights
  for select
  using (
    auth.uid() = user_id
    and exists (
      select 1
      from public.project_members pm
      where pm.project_id = ai_issue_insights.project_id
        and pm.user_id = auth.uid()
    )
  );

create policy "Users can create their own AI issue insights"
  on public.ai_issue_insights
  for insert
  with check (
    auth.uid() = user_id
    and exists (
      select 1
      from public.project_members pm
      where pm.project_id = ai_issue_insights.project_id
        and pm.user_id = auth.uid()
    )
  );

-- AI Crawl Summaries
-- Stores AI-generated summaries for entire crawls
create table if not exists public.ai_crawl_summaries (
  id uuid primary key default gen_random_uuid(),
  crawl_id uuid not null references public.crawls(id) on delete cascade,
  user_id uuid not null references auth.users(id) on delete cascade,
  project_id uuid not null references public.projects(id) on delete cascade,
  summary_text text not null,
  created_at timestamptz default now()
);

-- Indexes for AI crawl summaries
create index if not exists idx_ai_crawl_summaries_crawl_user on public.ai_crawl_summaries (crawl_id, user_id);

-- Enable RLS
alter table public.ai_crawl_summaries enable row level security;

-- RLS Policy: Users can only read summaries for crawls they have access to
create policy "Users can view AI crawl summaries for accessible crawls"
  on public.ai_crawl_summaries
  for select
  using (
    auth.uid() = user_id
    and exists (
      select 1
      from public.project_members pm
      where pm.project_id = ai_crawl_summaries.project_id
        and pm.user_id = auth.uid()
    )
  );

create policy "Users can create AI crawl summaries for accessible crawls"
  on public.ai_crawl_summaries
  for insert
  with check (
    auth.uid() = user_id
    and exists (
      select 1
      from public.project_members pm
      where pm.project_id = ai_crawl_summaries.project_id
        and pm.user_id = auth.uid()
    )
  );

