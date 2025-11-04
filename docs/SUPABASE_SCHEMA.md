# Supabase Schema & Migration Plan

This document proposes the initial Supabase (Postgres) schema, security model, and workflow for managing migrations in the Barracuda Cloud Run + Vercel architecture. Use it as the authoritative reference when creating or adjusting database objects.

---

## Guiding Principles

- Normalize crawl and issue data to avoid duplication while keeping query performance acceptable.
- Support both CLI-uploaded crawls and web-triggered re-crawls.
- Embrace Supabase Auth for user identity; leverage Row Level Security (RLS) so tenants only see their own data.
- Design for incremental enhancements (status tracking, recommendations, AI summaries) without heavy refactors.

---

## Core Entities

| Entity | Purpose |
|--------|---------|
| `profiles` | Per-user metadata that augments Supabase `auth.users`. |
| `projects` | Represents a tracked site/domain. Users can belong to multiple projects. |
| `project_members` | Join table controlling permissions within a project. |
| `crawls` | High-level crawl job metadata (start/end, initiator, status). |
| `pages` | Snapshot of a page discovered during a crawl, including key metrics. |
| `issues` | Individual SEO issues detected for a page (or crawl-level issue). |
| `issue_recommendations` | AI-generated or manual recommendations tied to an issue. |
| `issue_status_history` | Tracks workflow states (new, in progress, fixed, ignored). |
| `exports` | Records of generated exports (CSV/PDF) with storage locations. |
| `api_integrations` | Stores OAuth tokens/config for Google Search Console, etc. |

---

## Table Definitions

### 1. `profiles`
- Extends Supabase `auth.users` with display name, settings.
- Columns:
  - `id uuid primary key references auth.users (id) on delete cascade`
  - `display_name text`
  - `avatar_url text`
  - `created_at timestamp with time zone default current_timestamp`
  - `updated_at timestamp with time zone default current_timestamp`
- Indexes:
  - Primary key only (`id`).
- RLS: Allow owner (`auth.uid() = id`) to select/update. Service role can manage all.

### 2. `projects`
- Represents a domain/workspace.
- Columns:
  - `id uuid primary key default gen_random_uuid()`
  - `name text not null`
  - `domain text not null`
  - `owner_id uuid not null references auth.users (id) on delete cascade`
  - `created_at timestamptz default now()`
  - `updated_at timestamptz default now()`
  - `settings jsonb default '{}'::jsonb` (crawl defaults, thresholds)
- Indexes:
  - Unique `(owner_id, lower(domain))` to prevent duplicate domains per owner.
- RLS:
  - Owners and members can select/update.
  - Only owners (or service role) can delete.

### 3. `project_members`
- Many-to-many for collaboration.
- Columns:
  - `project_id uuid references projects (id) on delete cascade`
  - `user_id uuid references auth.users (id) on delete cascade`
  - `role text check (role in ('owner', 'editor', 'viewer')) default 'viewer'`
  - `invited_by uuid references auth.users (id)`
  - `created_at timestamptz default now()`
  - `updated_at timestamptz default now()`
- Constraints:
  - Primary key `(project_id, user_id)`
- RLS:
  - Members can select rows for their project.
  - Only owners (`role = 'owner'`) can update roles or remove members.

### 4. `crawls`
- Each crawl run.
- Columns:
  - `id uuid primary key default gen_random_uuid()`
  - `project_id uuid not null references projects (id) on delete cascade`
  - `initiated_by uuid references auth.users (id)`
  - `source text check (source in ('cli', 'web', 'schedule')) default 'cli'`
  - `status text check (status in ('pending', 'running', 'succeeded', 'failed', 'cancelled')) not null`
  - `started_at timestamptz default now()`
  - `completed_at timestamptz`
  - `total_pages integer default 0`
  - `total_issues integer default 0`
  - `meta jsonb default '{}'::jsonb` (config used, depth, notes)
- Indexes:
  - `idx_crawls_project_started` on `(project_id, started_at desc)`
  - `idx_crawls_status` on `(project_id, status)`
- RLS:
  - Members of the project can select.
  - Inserts allowed for authenticated users who belong to the project (enforced via policy and RPC).

### 5. `pages`
- Snapshot per URL per crawl.
- Columns:
  - `id bigserial primary key`
  - `crawl_id uuid not null references crawls (id) on delete cascade`
  - `url text not null`
  - `status_code integer`
  - `response_time_ms integer`
  - `title text`
  - `meta_description text`
  - `canonical_url text`
  - `h1 text`
  - `word_count integer`
  - `content_hash text`
  - `screenshot_url text` (optional reference to storage)
  - `data jsonb default '{}'::jsonb` (headings, links arrays)
  - `created_at timestamptz default now()`
- Indexes:
  - `idx_pages_crawl_url` unique `(crawl_id, url)`
  - `idx_pages_url_latest` partial index to quickly fetch latest crawl for a URL (`where completed_at is not null`)
- RLS:
  - Join through crawl -> project to ensure only members access rows. Use a `WITH` policy referencing `crawls`.

### 6. `issues`
- Concrete issues detected during a crawl.
- Columns:
  - `id bigserial primary key`
  - `crawl_id uuid not null references crawls (id) on delete cascade`
  - `page_id bigserial references pages (id) on delete set null` (some issues may be site-wide)
  - `project_id uuid generated always as (select crawls.project_id from crawls where crawls.id = crawl_id) stored`
  - `type text not null` (slug e.g., `missing_title`)
  - `severity text check (severity in ('error', 'warning', 'info')) not null`
  - `message text not null`
  - `recommendation text`
  - `value text` (raw value e.g., duplicate title string)
  - `priority_score integer`
  - `status text check (status in ('new', 'in_progress', 'fixed', 'ignored')) default 'new'`
  - `status_updated_at timestamptz default now()`
  - `created_at timestamptz default now()`
- Indexes:
  - `idx_issues_crawl_type` on `(crawl_id, type)`
  - `idx_issues_project_status` on `(project_id, status)`
  - `idx_issues_page` on `(page_id)`
- RLS:
  - Members of the corresponding project can select/update status.
  - Only users with role `editor` or `owner` can change status/recommendations.

### 7. `issue_recommendations`
- Supports multiple AI/manual suggestions per issue.
- Columns:
  - `id bigserial primary key`
  - `issue_id bigint references issues (id) on delete cascade`
  - `author_type text check (author_type in ('ai', 'user', 'system')) default 'ai'`
  - `author_id uuid references auth.users (id)`
  - `summary text not null`
  - `details text`
  - `created_at timestamptz default now()`
- Indexes:
  - `idx_issue_recommendations_issue` on `(issue_id)`
- RLS:
  - Constrain via issue -> project membership.

### 8. `issue_status_history`
- Logs workflow changes.
- Columns:
  - `id bigserial primary key`
  - `issue_id bigint references issues (id) on delete cascade`
  - `old_status text`
  - `new_status text`
  - `changed_by uuid references auth.users (id)`
  - `notes text`
  - `changed_at timestamptz default now()`
- Indexes:
  - `idx_issue_status_history_issue` on `(issue_id)`
- RLS:
  - Same project membership as issues.

### 9. `exports`
- Track generated reports.
- Columns:
  - `id uuid primary key default gen_random_uuid()`
  - `project_id uuid references projects (id) on delete cascade`
  - `crawl_id uuid references crawls (id) on delete set null`
  - `type text check (type in ('csv', 'json', 'pdf', 'html'))`
  - `storage_path text not null` (Supabase Storage or GCS key)
  - `requested_by uuid references auth.users (id)`
  - `status text check (status in ('queued', 'processing', 'ready', 'failed'))`
  - `created_at timestamptz default now()`
  - `completed_at timestamptz`
- Indexes:
  - `idx_exports_project_created` on `(project_id, created_at desc)`
- RLS:
  - Project members can view exports; only requesters or owners can delete.

### 10. `api_integrations`
- Stores external API tokens/configurations at the project level.
- Columns:
  - `id uuid primary key default gen_random_uuid()`
  - `project_id uuid references projects (id) on delete cascade`
  - `provider text check (provider in ('gsc', 'openai', 'pagespeed')) not null`
  - `config jsonb not null` (encrypted payload or reference to Secret Manager)
  - `created_at timestamptz default now()`
  - `updated_at timestamptz default now()`
- Indexes:
  - Unique `(project_id, provider)`
- RLS:
  - Only owners can insert/update; all members can read limited columns (consider view that redacts secrets).

---

## Supporting Objects

- **Views**
  - `project_issue_summary`: aggregates issue counts and priority per project for quick dashboard loads.
  - `latest_crawl_pages`: returns most recent crawl_id per URL for comparison features.

- **Functions**
  - `ensure_project_membership(project uuid)`: raises exception if `auth.uid()` lacks access; reusable in policies.
  - `create_crawl_with_pages(...)`: optional RPC to insert crawl metadata plus bulk pages in a transaction.

- **Extensions**
  - Enable `pgcrypto` (for `gen_random_uuid`), `pgmq` (optional message queue), `pg_stat_statements`.

---

## Row Level Security (RLS) Overview

1. Enable RLS on every table except system-owned ones.
2. Define a `current_project_membership` view that maps `auth.uid()` to projects for policy reuse.
3. Example policy for `crawls`:
   ```sql
   create policy "project members can view crawls"
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
   ```
4. Service role (Cloud Run) uses `service_role` key for unrestricted operations when necessary.

---

## Migration Workflow

1. **Directory Layout**
   - Create `db/migrations` at repo root for SQL migrations.
   - Optional: maintain `db/seed` for initial records.

2. **Tools**
   - Use Supabase CLI:
     ```bash
     supabase login
     supabase init
     supabase db new <migration_name>
     supabase db push      # local dev database
     supabase db reset     # optional for test resets
     ```
   - For CI/CD, run `supabase db push --password $SUPABASE_DB_PASSWORD`.

3. **Migration Style**
   - Track forward-only SQL files (e.g., `YYYYMMDDHHMM_add_projects.sql`).
   - Each migration contains explicit `create table`, `alter table`, policies, and grants.
   - Include comments referencing this doc section for traceability.

4. **Review Process**
   - PRs must include:
     - New/updated migration files in `db/migrations`.
     - Updates to this document when schema changes are significant.
     - Verification steps (e.g., `supabase db diff`) in PR description.

5. **Environment Promotion**
   - Local/dev: apply migrations via Supabase local emulator or remote dev instance.
   - Staging/Prod: apply using Supabase CLI against the respective `SUPABASE_DB_URL`.
   - Keep environment variables (service role key, anon key) in GitHub/Vercel secrets documented in `docs/CLOUD_RUN_SUPABASE.md`.

---

## Next Actions

1. Initialize Supabase project (`supabase init`) and add CLI config to repository.
2. Generate first migration implementing tables above.
3. Create helper RPCs and policies incrementally, validating with Supabase SQL editor or CLI tests.

---

**Last Updated:** {{ date }}
