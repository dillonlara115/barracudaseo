-- Phase 1: wrap auth.uid() calls in SELECT for hot-path tables to avoid per-row initplan eval
-- Tables: projects, crawls, pages, issues, exports, project_members, api_integrations

-- Projects
drop policy if exists "Project members and teammates can view projects" on public.projects;
drop policy if exists "Project owners and teammates can update projects" on public.projects;
drop policy if exists "Project owners and teammates can delete projects" on public.projects;
drop policy if exists "Authenticated users can create projects" on public.projects;

create policy "Project members and teammates can view projects"
  on public.projects
  for select
  using (
    owner_id = (select auth.uid())
    or exists (
      select 1
      from public.project_members pm
      where pm.project_id = projects.id
        and pm.user_id = (select auth.uid())
    )
    or exists (
      select 1
      from public.team_members tm1
      join public.team_members tm2 on tm1.account_owner_id = tm2.account_owner_id
      where tm1.user_id = (select auth.uid())
        and tm2.user_id = projects.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      select 1
      from public.profiles p
      join public.team_members tm on tm.account_owner_id = (select auth.uid())
      where p.id = (select auth.uid())
        and (p.subscription_tier = 'pro' or p.subscription_tier = 'team')
        and tm.user_id = projects.owner_id
        and tm.status = 'active'
    )
    or exists (
      select 1
      from public.profiles p
      join public.team_members tm on tm.account_owner_id = projects.owner_id
      where p.id = projects.owner_id
        and (p.subscription_tier = 'pro' or p.subscription_tier = 'team')
        and tm.user_id = (select auth.uid())
        and tm.status = 'active'
    )
  );

create policy "Project owners and teammates can update projects"
  on public.projects
  for update
  using (
    owner_id = (select auth.uid())
    or exists (
      select 1
      from public.team_members tm1
      join public.team_members tm2 on tm1.account_owner_id = tm2.account_owner_id
      where tm1.user_id = (select auth.uid())
        and tm2.user_id = projects.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      select 1
      from public.profiles p
      join public.team_members tm on tm.account_owner_id = (select auth.uid())
      where p.id = (select auth.uid())
        and (p.subscription_tier = 'pro' or p.subscription_tier = 'team')
        and tm.user_id = projects.owner_id
        and tm.status = 'active'
    )
    or exists (
      select 1
      from public.profiles p
      join public.team_members tm on tm.account_owner_id = projects.owner_id
      where p.id = projects.owner_id
        and (p.subscription_tier = 'pro' or p.subscription_tier = 'team')
        and tm.user_id = (select auth.uid())
        and tm.status = 'active'
    )
  );

create policy "Project owners and teammates can delete projects"
  on public.projects
  for delete
  using (
    owner_id = (select auth.uid())
    or exists (
      select 1
      from public.team_members tm1
      join public.team_members tm2 on tm1.account_owner_id = tm2.account_owner_id
      where tm1.user_id = (select auth.uid())
        and tm2.user_id = projects.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      select 1
      from public.profiles p
      join public.team_members tm on tm.account_owner_id = (select auth.uid())
      where p.id = (select auth.uid())
        and (p.subscription_tier = 'pro' or p.subscription_tier = 'team')
        and tm.user_id = projects.owner_id
        and tm.status = 'active'
    )
    or exists (
      select 1
      from public.profiles p
      join public.team_members tm on tm.account_owner_id = projects.owner_id
      where p.id = projects.owner_id
        and (p.subscription_tier = 'pro' or p.subscription_tier = 'team')
        and tm.user_id = (select auth.uid())
        and tm.status = 'active'
    )
  );

create policy "Authenticated users can create projects"
  on public.projects
  for insert
  with check ((select auth.uid()) = owner_id);

-- Crawls
drop policy if exists "Project members and teammates can view crawls" on public.crawls;
drop policy if exists "Project members and teammates can create crawls" on public.crawls;
drop policy if exists "Project members and teammates can update crawls" on public.crawls;
drop policy if exists "Project members and teammates can delete crawls" on public.crawls;

create policy "Project members and teammates can view crawls"
  on public.crawls
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = crawls.project_id
        and pm.user_id = (select auth.uid())
    )
    or exists (
      select 1
      from public.projects p
      join public.team_members tm1 on tm1.user_id = (select auth.uid())
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where p.id = crawls.project_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      select 1
      from public.projects p
      join public.profiles prof on prof.id = (select auth.uid())
      join public.team_members tm on tm.account_owner_id = (select auth.uid())
      where p.id = crawls.project_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = p.owner_id
        and tm.status = 'active'
    )
    or exists (
      select 1
      from public.projects p
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = crawls.project_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = (select auth.uid())
        and tm.status = 'active'
    )
  );

create policy "Project members and teammates can create crawls"
  on public.crawls
  for insert
  with check (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = crawls.project_id
        and pm.user_id = (select auth.uid())
    )
    or exists (
      select 1
      from public.projects p
      join public.team_members tm1 on tm1.user_id = (select auth.uid())
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where p.id = crawls.project_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      select 1
      from public.projects p
      join public.profiles prof on prof.id = (select auth.uid())
      join public.team_members tm on tm.account_owner_id = (select auth.uid())
      where p.id = crawls.project_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = p.owner_id
        and tm.status = 'active'
    )
    or exists (
      select 1
      from public.projects p
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = crawls.project_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = (select auth.uid())
        and tm.status = 'active'
    )
  );

create policy "Project members and teammates can update crawls"
  on public.crawls
  for update
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = crawls.project_id
        and pm.user_id = (select auth.uid())
    )
    or exists (
      select 1
      from public.projects p
      join public.team_members tm1 on tm1.user_id = (select auth.uid())
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where p.id = crawls.project_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      select 1
      from public.projects p
      join public.profiles prof on prof.id = (select auth.uid())
      join public.team_members tm on tm.account_owner_id = (select auth.uid())
      where p.id = crawls.project_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = p.owner_id
        and tm.status = 'active'
    )
    or exists (
      select 1
      from public.projects p
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = crawls.project_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = (select auth.uid())
        and tm.status = 'active'
    )
  )
  with check (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = crawls.project_id
        and pm.user_id = (select auth.uid())
    )
    or exists (
      select 1
      from public.projects p
      join public.team_members tm1 on tm1.user_id = (select auth.uid())
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where p.id = crawls.project_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      select 1
      from public.projects p
      join public.profiles prof on prof.id = (select auth.uid())
      join public.team_members tm on tm.account_owner_id = (select auth.uid())
      where p.id = crawls.project_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = p.owner_id
        and tm.status = 'active'
    )
    or exists (
      select 1
      from public.projects p
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = crawls.project_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = (select auth.uid())
        and tm.status = 'active'
    )
  );

create policy "Project members and teammates can delete crawls"
  on public.crawls
  for delete
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = crawls.project_id
        and pm.user_id = (select auth.uid())
    )
    or exists (
      select 1
      from public.projects p
      join public.team_members tm1 on tm1.user_id = (select auth.uid())
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where p.id = crawls.project_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      select 1
      from public.projects p
      join public.profiles prof on prof.id = (select auth.uid())
      join public.team_members tm on tm.account_owner_id = (select auth.uid())
      where p.id = crawls.project_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = p.owner_id
        and tm.status = 'active'
    )
    or exists (
      select 1
      from public.projects p
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = crawls.project_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = (select auth.uid())
        and tm.status = 'active'
    )
  );

-- Pages
drop policy if exists "Project members and teammates can view pages" on public.pages;
drop policy if exists "Project members and teammates can insert pages" on public.pages;

create policy "Project members and teammates can view pages"
  on public.pages
  for select
  using (
    exists (
      select 1
      from public.crawls c
      join public.project_members pm on pm.project_id = c.project_id
      where c.id = pages.crawl_id
        and pm.user_id = (select auth.uid())
    )
    or exists (
      select 1
      from public.crawls c
      join public.projects p on p.id = c.project_id
      join public.team_members tm1 on tm1.user_id = (select auth.uid())
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where c.id = pages.crawl_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      select 1
      from public.crawls c
      join public.projects p on p.id = c.project_id
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where c.id = pages.crawl_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = (select auth.uid())
        and tm.status = 'active'
    )
  );

create policy "Project members and teammates can insert pages"
  on public.pages
  for insert
  with check (
    exists (
      select 1
      from public.crawls c
      join public.project_members pm on pm.project_id = c.project_id
      where c.id = pages.crawl_id
        and pm.user_id = (select auth.uid())
    )
    or exists (
      select 1
      from public.crawls c
      join public.projects p on p.id = c.project_id
      join public.team_members tm1 on tm1.user_id = (select auth.uid())
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where c.id = pages.crawl_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      select 1
      from public.crawls c
      join public.projects p on p.id = c.project_id
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where c.id = pages.crawl_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = (select auth.uid())
        and tm.status = 'active'
    )
  );

-- Issues
drop policy if exists "Project members and teammates can view issues" on public.issues;
drop policy if exists "Project members and teammates can update issue status" on public.issues;
drop policy if exists "Project members and teammates can insert issues" on public.issues;

create policy "Project members and teammates can view issues"
  on public.issues
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = issues.project_id
        and pm.user_id = (select auth.uid())
    )
    or exists (
      select 1
      from public.projects p
      join public.team_members tm1 on tm1.user_id = (select auth.uid())
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where p.id = issues.project_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      select 1
      from public.projects p
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = issues.project_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = (select auth.uid())
        and tm.status = 'active'
    )
  );

create policy "Project members and teammates can update issue status"
  on public.issues
  for update
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = issues.project_id
        and pm.user_id = (select auth.uid())
    )
    or exists (
      select 1
      from public.projects p
      join public.team_members tm1 on tm1.user_id = (select auth.uid())
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where p.id = issues.project_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      select 1
      from public.projects p
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = issues.project_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = (select auth.uid())
        and tm.status = 'active'
    )
  )
  with check (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = issues.project_id
        and pm.user_id = (select auth.uid())
    )
    or exists (
      select 1
      from public.projects p
      join public.team_members tm1 on tm1.user_id = (select auth.uid())
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where p.id = issues.project_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      select 1
      from public.projects p
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = issues.project_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = (select auth.uid())
        and tm.status = 'active'
    )
  );

create policy "Project members and teammates can insert issues"
  on public.issues
  for insert
  with check (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = issues.project_id
        and pm.user_id = (select auth.uid())
    )
    or exists (
      select 1
      from public.projects p
      join public.team_members tm1 on tm1.user_id = (select auth.uid())
      join public.team_members tm2 on tm2.account_owner_id = tm1.account_owner_id
      where p.id = issues.project_id
        and tm2.user_id = p.owner_id
        and tm1.status = 'active'
        and tm2.status = 'active'
    )
    or exists (
      select 1
      from public.projects p
      join public.profiles prof on prof.id = p.owner_id
      join public.team_members tm on tm.account_owner_id = p.owner_id
      where p.id = issues.project_id
        and (prof.subscription_tier = 'pro' or prof.subscription_tier = 'team')
        and tm.user_id = (select auth.uid())
        and tm.status = 'active'
    )
  );

-- Exports
drop policy if exists "Project members can view exports" on public.exports;
drop policy if exists "Project members can create exports" on public.exports;
drop policy if exists "Export requesters and owners can delete exports" on public.exports;

create policy "Project members can view exports"
  on public.exports
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = exports.project_id
        and pm.user_id = (select auth.uid())
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
        and pm.user_id = (select auth.uid())
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

-- Project Members
drop policy if exists "Project members can view project members" on public.project_members;
drop policy if exists "Project owners can add members" on public.project_members;
drop policy if exists "Project owners can update members" on public.project_members;
drop policy if exists "Project owners can remove members" on public.project_members;

create policy "Project members can view project members"
  on public.project_members
  for select
  using (
    exists (
      select 1
      from public.project_members pm
      where pm.project_id = project_members.project_id
        and pm.user_id = (select auth.uid())
    )
  );

create policy "Project owners can add members"
  on public.project_members
  for insert
  with check (
    exists (
      select 1
      from public.projects p
      where p.id = project_members.project_id
        and p.owner_id = (select auth.uid())
    )
  );

create policy "Project owners can update members"
  on public.project_members
  for update
  using (
    exists (
      select 1
      from public.projects p
      where p.id = project_members.project_id
        and p.owner_id = (select auth.uid())
    )
  );

create policy "Project owners can remove members"
  on public.project_members
  for delete
  using (
    exists (
      select 1
      from public.projects p
      where p.id = project_members.project_id
        and p.owner_id = (select auth.uid())
    )
    or project_members.user_id = (select auth.uid())
  );

-- API Integrations (consolidated policies; ensure auth.uid wrapped)
drop policy if exists "API integrations select (project access)" on public.api_integrations;
drop policy if exists "API integrations insert (project access)" on public.api_integrations;
drop policy if exists "API integrations update (project access)" on public.api_integrations;
drop policy if exists "API integrations delete (project access)" on public.api_integrations;

create policy "API integrations select (project access)"
  on public.api_integrations
  for select
  to authenticated
  using (public.can_access_project(api_integrations.project_id, (select auth.uid())));

create policy "API integrations insert (project access)"
  on public.api_integrations
  for insert
  to authenticated
  with check (public.can_access_project(api_integrations.project_id, (select auth.uid())));

create policy "API integrations update (project access)"
  on public.api_integrations
  for update
  to authenticated
  using (public.can_access_project(api_integrations.project_id, (select auth.uid())))
  with check (public.can_access_project(api_integrations.project_id, (select auth.uid())));

create policy "API integrations delete (project access)"
  on public.api_integrations
  for delete
  to authenticated
  using (public.can_access_project(api_integrations.project_id, (select auth.uid())));
