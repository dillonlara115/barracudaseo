-- Add team-based access to pages, issues, and related tables
-- Team members should be able to view all crawl data from projects they have team access to

-- Update Pages RLS Policies
drop policy if exists "Project members can view pages" on public.pages;
drop policy if exists "Project members and teammates can view pages" on public.pages;
create policy "Project members and teammates can view pages"
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
    or exists (
      -- Check if user and project owner are both team members with same account_owner_id
      select 1
      from public.crawls c
      join public.projects p on p.id = c.project_id
      join public.team_members tm1 on tm1.user_id = auth.uid()
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where c.id = pages.crawl_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner (has pro/team tier) and user is their team member
      select 1
      from public.crawls c
      join public.projects p on p.id = c.project_id
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where c.id = pages.crawl_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = auth.uid()
        and tm.status = 'active'
    )
  );

drop policy if exists "Project members can insert pages" on public.pages;
drop policy if exists "Project members and teammates can insert pages" on public.pages;
create policy "Project members and teammates can insert pages"
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
    or exists (
      -- Check if user and project owner are both team members with same account_owner_id
      select 1
      from public.crawls c
      join public.projects p on p.id = c.project_id
      join public.team_members tm1 on tm1.user_id = auth.uid()
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where c.id = pages.crawl_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner (has pro/team tier) and user is their team member
      select 1
      from public.crawls c
      join public.projects p on p.id = c.project_id
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where c.id = pages.crawl_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = auth.uid()
        and tm.status = 'active'
    )
  );

-- Update Issues RLS Policies
drop policy if exists "Project members can view issues" on public.issues;
drop policy if exists "Project members and teammates can view issues" on public.issues;
create policy "Project members and teammates can view issues"
  on public.issues
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = issues.project_id
        and pm.user_id = auth.uid()
    )
    or exists (
      -- Check if user and project owner are both team members with same account_owner_id
      select 1
      from public.projects p
      join public.team_members tm1 on tm1.user_id = auth.uid()
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where p.id = issues.project_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner (has pro/team tier) and user is their team member
      select 1
      from public.projects p
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = issues.project_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = auth.uid()
        and tm.status = 'active'
    )
  );

drop policy if exists "Project members can update issue status" on public.issues;
drop policy if exists "Project members and teammates can update issue status" on public.issues;
create policy "Project members and teammates can update issue status"
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
    or exists (
      -- Check if user and project owner are both team members with same account_owner_id
      select 1
      from public.projects p
      join public.team_members tm1 on tm1.user_id = auth.uid()
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where p.id = issues.project_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner (has pro/team tier) and user is their team member
      select 1
      from public.projects p
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = issues.project_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = auth.uid()
        and tm.status = 'active'
    )
  );

drop policy if exists "Project members can insert issues" on public.issues;
drop policy if exists "Project members and teammates can insert issues" on public.issues;
create policy "Project members and teammates can insert issues"
  on public.issues
  for insert
  with check (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = issues.project_id
        and pm.user_id = auth.uid()
    )
    or exists (
      -- Check if user and project owner are both team members with same account_owner_id
      select 1
      from public.projects p
      join public.team_members tm1 on tm1.user_id = auth.uid()
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where p.id = issues.project_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner (has pro/team tier) and user is their team member
      select 1
      from public.projects p
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = issues.project_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = auth.uid()
        and tm.status = 'active'
    )
  );

-- Update Issue Recommendations RLS Policies
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
        and pm.user_id = auth.uid()
    )
    or exists (
      -- Check if user and project owner are both team members with same account_owner_id
      select 1
      from public.issues i
      join public.projects p on p.id = i.project_id
      join public.team_members tm1 on tm1.user_id = auth.uid()
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where i.id = issue_recommendations.issue_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner (has pro/team tier) and user is their team member
      select 1
      from public.issues i
      join public.projects p on p.id = i.project_id
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where i.id = issue_recommendations.issue_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = auth.uid()
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
        and pm.user_id = auth.uid()
    )
    or exists (
      -- Check if user and project owner are both team members with same account_owner_id
      select 1
      from public.issues i
      join public.projects p on p.id = i.project_id
      join public.team_members tm1 on tm1.user_id = auth.uid()
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where i.id = issue_recommendations.issue_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner (has pro/team tier) and user is their team member
      select 1
      from public.issues i
      join public.projects p on p.id = i.project_id
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where i.id = issue_recommendations.issue_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = auth.uid()
        and tm.status = 'active'
    )
  );

-- Update Issue Status History RLS Policies
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
        and pm.user_id = auth.uid()
    )
    or exists (
      -- Check if user and project owner are both team members with same account_owner_id
      select 1
      from public.issues i
      join public.projects p on p.id = i.project_id
      join public.team_members tm1 on tm1.user_id = auth.uid()
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where i.id = issue_status_history.issue_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner (has pro/team tier) and user is their team member
      select 1
      from public.issues i
      join public.projects p on p.id = i.project_id
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where i.id = issue_status_history.issue_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = auth.uid()
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
        and pm.user_id = auth.uid()
    )
    or exists (
      -- Check if user and project owner are both team members with same account_owner_id
      select 1
      from public.issues i
      join public.projects p on p.id = i.project_id
      join public.team_members tm1 on tm1.user_id = auth.uid()
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where i.id = issue_status_history.issue_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner (has pro/team tier) and user is their team member
      select 1
      from public.issues i
      join public.projects p on p.id = i.project_id
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where i.id = issue_status_history.issue_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = auth.uid()
        and tm.status = 'active'
    )
  );

-- Update Exports RLS Policies
drop policy if exists "Project members can view exports" on public.exports;
drop policy if exists "Project members and teammates can view exports" on public.exports;
create policy "Project members and teammates can view exports"
  on public.exports
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = exports.project_id
        and pm.user_id = auth.uid()
    )
    or exists (
      -- Check if user and project owner are both team members with same account_owner_id
      select 1
      from public.projects p
      join public.team_members tm1 on tm1.user_id = auth.uid()
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where p.id = exports.project_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner (has pro/team tier) and user is their team member
      select 1
      from public.projects p
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = exports.project_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = auth.uid()
        and tm.status = 'active'
    )
  );

drop policy if exists "Project members can create exports" on public.exports;
drop policy if exists "Project members and teammates can create exports" on public.exports;
create policy "Project members and teammates can create exports"
  on public.exports
  for insert
  with check (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = exports.project_id
        and pm.user_id = auth.uid()
    )
    or exists (
      -- Check if user and project owner are both team members with same account_owner_id
      select 1
      from public.projects p
      join public.team_members tm1 on tm1.user_id = auth.uid()
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where p.id = exports.project_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      -- Check if project owner is account owner (has pro/team tier) and user is their team member
      select 1
      from public.projects p
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = exports.project_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = auth.uid()
        and tm.status = 'active'
    )
  );
