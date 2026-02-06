# Database Guidelines (Supabase/PostgreSQL)

## Migration File Naming

Format: `YYYYMMDDHHmmss_descriptive_name.sql`

Examples:
- `20250206120000_add_keyword_tracking.sql`
- `20250206130000_fix_rls_policies.sql`

## Table Structure Pattern

```sql
create table if not exists public.resources (
  id uuid primary key default gen_random_uuid(),
  project_id uuid not null references public.projects (id) on delete cascade,
  name text not null,
  status text check (status in ('pending', 'active', 'archived')) default 'pending',
  settings jsonb default '{}'::jsonb,
  created_at timestamptz default now(),
  updated_at timestamptz default now()
);

-- Enable RLS immediately
alter table public.resources enable row level security;
```

### Primary Keys
- UUIDs for most tables: `id uuid primary key default gen_random_uuid()`
- Bigserial for high-volume tables: `id bigserial primary key`

### Foreign Keys
```sql
-- Cascade delete (child has no meaning without parent)
project_id uuid not null references public.projects (id) on delete cascade

-- Set null (optional relationship, preserve child)
page_id bigint references public.pages (id) on delete set null

-- Auth users reference
owner_id uuid not null references auth.users (id) on delete cascade
```

### Timestamps
Always include:
```sql
created_at timestamptz default now(),
updated_at timestamptz default now()
```

### Enums via Check Constraints
```sql
status text check (status in ('pending', 'running', 'succeeded', 'failed')) not null
role text check (role in ('owner', 'editor', 'viewer')) default 'viewer'
```

## Row Level Security (RLS)

### Enable RLS on All Tables
```sql
alter table public.resources enable row level security;
```

### Simple Self-Access
```sql
create policy "Users can view own profile"
  on public.profiles for select
  using ((select auth.uid()) = id);
```

### Project-Based Access
```sql
create policy "Project members can view resources"
  on public.resources for select
  using (
    exists (
      select 1 from public.project_members pm
      where pm.project_id = resources.project_id
        and pm.user_id = (select auth.uid())
    )
  );
```

### Role-Based Write Access
```sql
create policy "Editors can update resources"
  on public.resources for update
  using (
    exists (
      select 1 from public.project_members pm
      where pm.project_id = resources.project_id
        and pm.user_id = (select auth.uid())
        and pm.role in ('owner', 'editor')
    )
  );
```

### RLS Optimization Tips

1. **Wrap auth.uid() in subquery** to avoid per-row evaluation:
   ```sql
   using ((select auth.uid()) = user_id)
   -- NOT: using (auth.uid() = user_id)
   ```

2. **Use SECURITY DEFINER helper functions** for complex logic:
   ```sql
   create or replace function public.can_access_project(project_id uuid)
   returns boolean
   language sql
   security definer
   stable
   set search_path = public
   as $$
     select exists (
       select 1 from public.project_members pm
       where pm.project_id = $1
         and pm.user_id = (select auth.uid())
     );
   $$;

   grant execute on function public.can_access_project(uuid) to authenticated;
   ```

3. **Consolidate duplicate policies** into single policies

## Index Conventions

### Naming
`idx_[table]_[columns]` or `idx_[table]_[columns]_[qualifier]`

### Common Patterns

```sql
-- Single column
create index if not exists idx_resources_project_id
  on public.resources (project_id);

-- Composite (for common query patterns)
create index if not exists idx_crawls_project_started
  on public.crawls (project_id, started_at desc);

-- Unique constraint
create unique index if not exists idx_pages_crawl_url
  on public.pages (crawl_id, url);

-- Partial index (filtered)
create index if not exists idx_team_members_active
  on public.team_members (user_id, account_owner_id)
  where status = 'active';

-- Case-insensitive
create unique index if not exists idx_projects_owner_domain
  on public.projects (owner_id, lower(domain));
```

### Index Strategy
- Index foreign keys used in JOINs
- Index columns used in WHERE clauses
- Use composite indexes matching query column order
- Use partial indexes for status-filtered queries
- Add `desc` for time-series queries sorted by recency

## Migration Best Practices

1. **Use `if not exists`** for idempotency:
   ```sql
   create table if not exists ...
   create index if not exists ...
   ```

2. **Drop before recreate** for policies/functions:
   ```sql
   drop policy if exists "policy_name" on public.table;
   create policy "policy_name" ...
   ```

3. **Comment complex logic**:
   ```sql
   -- Team members inherit access from account owner's projects
   create policy ...
   ```

4. **Grant permissions** on SECURITY DEFINER functions:
   ```sql
   grant execute on function public.helper_fn(uuid) to authenticated;
   ```

## Triggers

### Updated At Trigger
```sql
create or replace function public.handle_updated_at()
returns trigger
language plpgsql
as $$
begin
  new.updated_at = now();
  return new;
end;
$$;

create trigger set_updated_at
  before update on public.resources
  for each row
  execute function public.handle_updated_at();
```

## JSONB Usage

Use for flexible/semi-structured data:
```sql
settings jsonb default '{}'::jsonb,
metadata jsonb default '{}'::jsonb
```

Query patterns:
```sql
-- Access nested value
settings->>'theme'
settings->'notifications'->>'email'

-- Filter by JSONB value
where settings->>'enabled' = 'true'
```

## Supabase-Specific

### Auth References
```sql
references auth.users (id) on delete cascade
```

### Service Role vs Anon Key
- **Anon key**: Respects RLS, use for user-initiated queries
- **Service role**: Bypasses RLS, use for system/admin operations

### Common Tables
- `auth.users` - Supabase-managed user accounts
- `public.profiles` - Extended user data (create trigger on user signup)
