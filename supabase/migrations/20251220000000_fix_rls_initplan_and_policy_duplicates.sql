-- Fix RLS initplan warnings by wrapping auth.uid() calls and remove duplicate permissive policies.

-- Keywords
-- Ensure auth.uid() is wrapped and avoid per-row re-evaluation.
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

-- GSC tables

drop policy if exists "Project members can view gsc sync state" on public.gsc_sync_states;
drop policy if exists "Project members and teammates can view gsc sync state" on public.gsc_sync_states;
create policy "Project members can view gsc sync state"
  on public.gsc_sync_states
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = gsc_sync_states.project_id
        and pm.user_id = (select auth.uid())
    )
    or exists (
      select 1
      from public.projects p
      where p.id = gsc_sync_states.project_id
        and p.owner_id = (select auth.uid())
    )
  );


drop policy if exists "Project members can view gsc snapshots" on public.gsc_performance_snapshots;
drop policy if exists "Project members and teammates can view gsc snapshots" on public.gsc_performance_snapshots;
create policy "Project members can view gsc snapshots"
  on public.gsc_performance_snapshots
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = gsc_performance_snapshots.project_id
        and pm.user_id = (select auth.uid())
    )
    or exists (
      select 1
      from public.projects p
      where p.id = gsc_performance_snapshots.project_id
        and p.owner_id = (select auth.uid())
    )
  );


drop policy if exists "Project members can view gsc performance rows" on public.gsc_performance_rows;
drop policy if exists "Project members and teammates can view gsc performance rows" on public.gsc_performance_rows;
create policy "Project members can view gsc performance rows"
  on public.gsc_performance_rows
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = gsc_performance_rows.project_id
        and pm.user_id = (select auth.uid())
    )
    or exists (
      select 1
      from public.projects p
      where p.id = gsc_performance_rows.project_id
        and p.owner_id = (select auth.uid())
    )
  );


drop policy if exists "Project members can view gsc page enhancements" on public.gsc_page_enhancements;
drop policy if exists "Project members and teammates can view gsc page enhancements" on public.gsc_page_enhancements;
create policy "Project members can view gsc page enhancements"
  on public.gsc_page_enhancements
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = gsc_page_enhancements.project_id
        and pm.user_id = (select auth.uid())
    )
    or exists (
      select 1
      from public.projects p
      where p.id = gsc_page_enhancements.project_id
        and p.owner_id = (select auth.uid())
    )
  );


drop policy if exists "Project members can view gsc insights" on public.gsc_insights;
drop policy if exists "Project members and teammates can view gsc insights" on public.gsc_insights;
create policy "Project members can view gsc insights"
  on public.gsc_insights
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = gsc_insights.project_id
        and pm.user_id = (select auth.uid())
    )
    or exists (
      select 1
      from public.projects p
      where p.id = gsc_insights.project_id
        and p.owner_id = (select auth.uid())
    )
  );

-- Issue recommendations

drop policy if exists "Project members can view recommendations" on public.issue_recommendations;
drop policy if exists "Project members and teammates can view recommendations" on public.issue_recommendations;
create policy "Project members and teammates can view recommendations"
  on public.issue_recommendations
  for select
  using (
    exists (
      select 1
      from public.issues i
      join public.project_members pm on pm.project_id = i.project_id
      where i.id = issue_recommendations.issue_id
        and pm.user_id = (select auth.uid())
    )
    or exists (
      select 1
      from public.issues i
      join public.projects p on p.id = i.project_id
      join public.team_members tm1 on tm1.user_id = (select auth.uid())
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where i.id = issue_recommendations.issue_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      select 1
      from public.issues i
      join public.projects p on p.id = i.project_id
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where i.id = issue_recommendations.issue_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = (select auth.uid())
        and tm.status = 'active'
    )
  );


drop policy if exists "Project members can create recommendations" on public.issue_recommendations;
drop policy if exists "Project members and teammates can create recommendations" on public.issue_recommendations;
create policy "Project members and teammates can create recommendations"
  on public.issue_recommendations
  for insert
  with check (
    exists (
      select 1
      from public.issues i
      join public.project_members pm on pm.project_id = i.project_id
      where i.id = issue_recommendations.issue_id
        and pm.user_id = (select auth.uid())
    )
    or exists (
      select 1
      from public.issues i
      join public.projects p on p.id = i.project_id
      join public.team_members tm1 on tm1.user_id = (select auth.uid())
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where i.id = issue_recommendations.issue_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      select 1
      from public.issues i
      join public.projects p on p.id = i.project_id
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where i.id = issue_recommendations.issue_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = (select auth.uid())
        and tm.status = 'active'
    )
  );

-- Issue status history

drop policy if exists "Project members can view status history" on public.issue_status_history;
drop policy if exists "Project members and teammates can view status history" on public.issue_status_history;
create policy "Project members and teammates can view status history"
  on public.issue_status_history
  for select
  using (
    exists (
      select 1
      from public.issues i
      join public.project_members pm on pm.project_id = i.project_id
      where i.id = issue_status_history.issue_id
        and pm.user_id = (select auth.uid())
    )
    or exists (
      select 1
      from public.issues i
      join public.projects p on p.id = i.project_id
      join public.team_members tm1 on tm1.user_id = (select auth.uid())
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where i.id = issue_status_history.issue_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      select 1
      from public.issues i
      join public.projects p on p.id = i.project_id
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where i.id = issue_status_history.issue_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = (select auth.uid())
        and tm.status = 'active'
    )
  );


drop policy if exists "Project members can create status history" on public.issue_status_history;
drop policy if exists "Project members and teammates can create status history" on public.issue_status_history;
create policy "Project members and teammates can create status history"
  on public.issue_status_history
  for insert
  with check (
    exists (
      select 1
      from public.issues i
      join public.project_members pm on pm.project_id = i.project_id
      where i.id = issue_status_history.issue_id
        and pm.user_id = (select auth.uid())
    )
    or exists (
      select 1
      from public.issues i
      join public.projects p on p.id = i.project_id
      join public.team_members tm1 on tm1.user_id = (select auth.uid())
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where i.id = issue_status_history.issue_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      select 1
      from public.issues i
      join public.projects p on p.id = i.project_id
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where i.id = issue_status_history.issue_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = (select auth.uid())
        and tm.status = 'active'
    )
  );

-- Exports (deduplicate policies)

drop policy if exists "Project members can view exports" on public.exports;
drop policy if exists "Project members and teammates can view exports" on public.exports;
drop policy if exists "Project members can create exports" on public.exports;
drop policy if exists "Project members and teammates can create exports" on public.exports;
drop policy if exists "Export requesters and owners can delete exports" on public.exports;

create policy "Project members and teammates can view exports"
  on public.exports
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = exports.project_id
        and pm.user_id = (select auth.uid())
    )
    or exists (
      select 1
      from public.projects p
      join public.team_members tm1 on tm1.user_id = (select auth.uid())
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where p.id = exports.project_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      select 1
      from public.projects p
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = exports.project_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = (select auth.uid())
        and tm.status = 'active'
    )
  );

create policy "Project members and teammates can create exports"
  on public.exports
  for insert
  with check (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = exports.project_id
        and pm.user_id = (select auth.uid())
    )
    or exists (
      select 1
      from public.projects p
      join public.team_members tm1 on tm1.user_id = (select auth.uid())
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where p.id = exports.project_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      select 1
      from public.projects p
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = exports.project_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = (select auth.uid())
        and tm.status = 'active'
    )
  );

create policy "Export requesters and owners can delete exports"
  on public.exports
  for delete
  using (
    requested_by = (select auth.uid())
    or exists (
      select 1
      from public.projects p
      where p.id = exports.project_id
        and p.owner_id = (select auth.uid())
    )
  );

-- User AI settings

drop policy if exists "Users can view their own AI settings" on public.user_ai_settings;
drop policy if exists "Users can update their own AI settings" on public.user_ai_settings;
drop policy if exists "Users can insert their own AI settings" on public.user_ai_settings;

create policy "Users can view their own AI settings"
  on public.user_ai_settings
  for select
  using ((select auth.uid()) = user_id);

create policy "Users can update their own AI settings"
  on public.user_ai_settings
  for update
  using ((select auth.uid()) = user_id);

create policy "Users can insert their own AI settings"
  on public.user_ai_settings
  for insert
  with check ((select auth.uid()) = user_id);

-- AI issue insights

drop policy if exists "Users can view their own AI issue insights" on public.ai_issue_insights;
drop policy if exists "Users can create their own AI issue insights" on public.ai_issue_insights;

create policy "Users can view their own AI issue insights"
  on public.ai_issue_insights
  for select
  using (
    (select auth.uid()) = user_id
    and exists (
      select 1
      from public.project_members pm
      where pm.project_id = ai_issue_insights.project_id
        and pm.user_id = (select auth.uid())
    )
  );

create policy "Users can create their own AI issue insights"
  on public.ai_issue_insights
  for insert
  with check (
    (select auth.uid()) = user_id
    and exists (
      select 1
      from public.project_members pm
      where pm.project_id = ai_issue_insights.project_id
        and pm.user_id = (select auth.uid())
    )
  );

-- AI crawl summaries

drop policy if exists "Users can view AI crawl summaries for accessible crawls" on public.ai_crawl_summaries;
drop policy if exists "Users can create AI crawl summaries for accessible crawls" on public.ai_crawl_summaries;

create policy "Users can view AI crawl summaries for accessible crawls"
  on public.ai_crawl_summaries
  for select
  using (
    (select auth.uid()) = user_id
    and exists (
      select 1
      from public.project_members pm
      where pm.project_id = ai_crawl_summaries.project_id
        and pm.user_id = (select auth.uid())
    )
  );

create policy "Users can create AI crawl summaries for accessible crawls"
  on public.ai_crawl_summaries
  for insert
  with check (
    (select auth.uid()) = user_id
    and exists (
      select 1
      from public.project_members pm
      where pm.project_id = ai_crawl_summaries.project_id
        and pm.user_id = (select auth.uid())
    )
  );

-- Public reports

drop policy if exists "Users can view their own public reports" on public.public_reports;
drop policy if exists "Users can create public reports for their projects" on public.public_reports;
drop policy if exists "Users can update their own public reports" on public.public_reports;
drop policy if exists "Users can delete their own public reports" on public.public_reports;

create policy "Users can view their own public reports"
  on public.public_reports
  for select
  using (
    (select auth.uid()) = created_by
    and exists (
      select 1
      from public.project_members pm
      where pm.project_id = public_reports.project_id
        and pm.user_id = (select auth.uid())
    )
  );

create policy "Users can create public reports for their projects"
  on public.public_reports
  for insert
  with check (
    (select auth.uid()) = created_by
    and exists (
      select 1
      from public.project_members pm
      where pm.project_id = public_reports.project_id
        and pm.user_id = (select auth.uid())
    )
  );

create policy "Users can update their own public reports"
  on public.public_reports
  for update
  using (
    (select auth.uid()) = created_by
    and exists (
      select 1
      from public.project_members pm
      where pm.project_id = public_reports.project_id
        and pm.user_id = (select auth.uid())
    )
  );

create policy "Users can delete their own public reports"
  on public.public_reports
  for delete
  using (
    (select auth.uid()) = created_by
    and exists (
      select 1
      from public.project_members pm
      where pm.project_id = public_reports.project_id
        and pm.user_id = (select auth.uid())
    )
  );
