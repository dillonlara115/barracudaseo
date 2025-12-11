# RLS Optimization Checklist (auth.uid() initplan warnings)

Track fixes for the Supabase linter warnings about `auth.*`/`current_setting()` being re-evaluated per row. Update status as migrations land.

Statuses: TODO | In Progress | Done (applied) | Done (pending deploy)

## Phase 1 — Hot Path Tables
- [x] `projects` — wrap `auth.uid()` in SELECT for all policies (view/update/delete/insert)
- [x] `crawls` — wrap for view/create/update/delete
- [x] `pages` — wrap for view/insert
- [x] `issues` — wrap for view/insert/update
- [x] `exports` — wrap for view/create/delete
- [x] `project_members` — wrap for view/add/update/remove
- [x] `api_integrations` — wrap for select/insert/update/delete (consolidated & wrapped)

## Phase 2 — Identity/Billing
- [x] `profiles` — create/update (and consolidated select policy)
- [x] `subscriptions` — view
- [x] `team_members` — policies recreated with SELECT wrappers

## Phase 3 — Rank Tracking
- [ ] `keywords` — all policies
- [ ] `keyword_tasks` — all policies
- [ ] `keyword_rank_snapshots` — all policies
- [ ] `keyword_usage` — view

## Phase 4 — GSC Cache Tables
- [ ] `gsc_sync_states` — view
- [ ] `gsc_performance_snapshots` — view
- [ ] `gsc_performance_rows` — view
- [ ] `gsc_page_enhancements` — view
- [ ] `gsc_insights` — view

## Phase 5 — AI + Public Reports
- [ ] `user_ai_settings` — view/update/insert
- [ ] `ai_issue_insights` — view/create
- [ ] `ai_crawl_summaries` — view/create
- [ ] `public_reports` — view/create/update/delete

## Notes
- For each policy: change `auth.uid()` (or `auth.jwt()`/`current_setting`) to `(SELECT auth.uid())` etc.
- Keep TO clauses explicit (`TO authenticated`).
- Add indexes if predicates reference non-indexed columns.
- After migration: run `supabase db push` and spot-check as owner, member, non-member.
