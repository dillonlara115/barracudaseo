-- Phase 3: Rank tracking tables — wrap auth.uid() in SELECT
-- Tables: keywords, keyword_tasks, keyword_rank_snapshots, keyword_usage

-- Ensure project_id is present on keyword_tasks and keyword_rank_snapshots for direct policy checks
alter table public.keyword_tasks
  add column if not exists project_id uuid references public.projects(id) on delete cascade;

alter table public.keyword_rank_snapshots
  add column if not exists project_id uuid references public.projects(id) on delete cascade;

-- Backfill project_id from keywords
update public.keyword_tasks kt
set project_id = k.project_id
from public.keywords k
where kt.project_id is null
  and k.id = kt.keyword_id;

update public.keyword_rank_snapshots ks
set project_id = k.project_id
from public.keywords k
where ks.project_id is null
  and k.id = ks.keyword_id;

-- Enforce NOT NULL and add indexes
alter table public.keyword_tasks
  alter column project_id set not null;

alter table public.keyword_rank_snapshots
  alter column project_id set not null;

create index if not exists idx_keyword_tasks_project_id
  on public.keyword_tasks (project_id);

create index if not exists idx_keyword_rank_snapshots_project_id
  on public.keyword_rank_snapshots (project_id);

-- Keywords
drop policy if exists "Project members can view keywords" on public.keywords;
drop policy if exists "Project members can insert keywords" on public.keywords;
drop policy if exists "Project members can update keywords" on public.keywords;
drop policy if exists "Project members can delete keywords" on public.keywords;

create policy "Project members can view keywords"
  on public.keywords
  for select
  using (
    public.is_project_member(keywords.project_id, (select auth.uid()))
  );

create policy "Project members can insert keywords"
  on public.keywords
  for insert
  with check (
    public.is_project_member(keywords.project_id, (select auth.uid()))
  );

create policy "Project members can update keywords"
  on public.keywords
  for update
  using (
    public.is_project_member(keywords.project_id, (select auth.uid()))
  )
  with check (
    public.is_project_member(keywords.project_id, (select auth.uid()))
  );

create policy "Project members can delete keywords"
  on public.keywords
  for delete
  using (
    public.is_project_member(keywords.project_id, (select auth.uid()))
  );

-- Keyword tasks
drop policy if exists "Project members can view keyword tasks" on public.keyword_tasks;
drop policy if exists "Project members can insert keyword tasks" on public.keyword_tasks;
drop policy if exists "Project members can update keyword tasks" on public.keyword_tasks;
drop policy if exists "Project members can delete keyword tasks" on public.keyword_tasks;

create policy "Project members can view keyword tasks"
  on public.keyword_tasks
  for select
  using (
    public.is_project_member(keyword_tasks.project_id, (select auth.uid()))
  );

create policy "Project members can insert keyword tasks"
  on public.keyword_tasks
  for insert
  with check (
    public.is_project_member(keyword_tasks.project_id, (select auth.uid()))
  );

create policy "Project members can update keyword tasks"
  on public.keyword_tasks
  for update
  using (
    public.is_project_member(keyword_tasks.project_id, (select auth.uid()))
  )
  with check (
    public.is_project_member(keyword_tasks.project_id, (select auth.uid()))
  );

create policy "Project members can delete keyword tasks"
  on public.keyword_tasks
  for delete
  using (
    public.is_project_member(keyword_tasks.project_id, (select auth.uid()))
  );

-- Keyword rank snapshots
drop policy if exists "Project members can view keyword snapshots" on public.keyword_rank_snapshots;
drop policy if exists "Project members can insert keyword snapshots" on public.keyword_rank_snapshots;

create policy "Project members can view keyword snapshots"
  on public.keyword_rank_snapshots
  for select
  using (
    public.is_project_member(keyword_rank_snapshots.project_id, (select auth.uid()))
  );

create policy "Project members can insert keyword snapshots"
  on public.keyword_rank_snapshots
  for insert
  with check (
    public.is_project_member(keyword_rank_snapshots.project_id, (select auth.uid()))
  );

-- Keyword usage
drop policy if exists "Project members can view keyword usage" on public.keyword_usage;

create policy "Project members can view keyword usage"
  on public.keyword_usage
  for select
  using (
    public.is_project_member(keyword_usage.project_id, (select auth.uid()))
  );
