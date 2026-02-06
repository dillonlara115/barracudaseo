-- Global (user-level) integrations for GSC/GA4

create table if not exists public.user_api_integrations (
  id uuid primary key default gen_random_uuid(),
  user_id uuid references auth.users (id) on delete cascade,
  provider text check (provider in ('gsc', 'ga4')) not null,
  config jsonb not null,
  created_at timestamptz default now(),
  updated_at timestamptz default now(),
  unique (user_id, provider)
);

create index if not exists idx_user_api_integrations_user on public.user_api_integrations (user_id);

alter table public.user_api_integrations enable row level security;

drop policy if exists "Users can view their own integrations" on public.user_api_integrations;
create policy "Users can view their own integrations"
  on public.user_api_integrations
  for select
  using (auth.uid() = user_id);

drop policy if exists "Users can manage their own integrations" on public.user_api_integrations;
create policy "Users can manage their own integrations"
  on public.user_api_integrations
  for all
  using (auth.uid() = user_id)
  with check (auth.uid() = user_id);

-- Migrate existing per-project integrations to user-level integrations
insert into public.user_api_integrations (user_id, provider, config, created_at, updated_at)
select distinct on (p.owner_id, ai.provider)
  p.owner_id,
  ai.provider,
  ai.config,
  coalesce(ai.created_at, now()),
  coalesce(ai.updated_at, now())
from public.api_integrations ai
join public.projects p on p.id = ai.project_id
where ai.provider in ('gsc', 'ga4')
order by p.owner_id, ai.provider, ai.updated_at desc nulls last, ai.created_at desc nulls last;

-- Backfill project settings for GSC selection if missing
update public.projects p
set settings = jsonb_set(
  coalesce(p.settings, '{}'::jsonb),
  '{gsc_property_url}',
  to_jsonb(ai.config->>'property_url'),
  true
)
from public.api_integrations ai
where ai.project_id = p.id
  and ai.provider = 'gsc'
  and (ai.config->>'property_url') is not null
  and (ai.config->>'property_url') <> ''
  and (p.settings->>'gsc_property_url' is null or p.settings->>'gsc_property_url' = '');

-- Backfill project settings for GA4 selection
update public.projects p
set settings = jsonb_set(
  jsonb_set(
    coalesce(p.settings, '{}'::jsonb),
    '{ga4_property_id}',
    to_jsonb(ai.config->>'property_id'),
    true
  ),
  '{ga4_property_name}',
  to_jsonb(ai.config->>'property_name'),
  true
)
from public.api_integrations ai
where ai.project_id = p.id
  and ai.provider = 'ga4'
  and (ai.config->>'property_id') is not null
  and (ai.config->>'property_id') <> '';

-- Backfill integration user mapping on projects
update public.projects p
set settings = jsonb_set(
  coalesce(p.settings, '{}'::jsonb),
  '{gsc_integration_user_id}',
  to_jsonb(p.owner_id),
  true
)
where (p.settings->>'gsc_property_url') is not null
  and (p.settings->>'gsc_property_url') <> ''
  and (p.settings->>'gsc_integration_user_id' is null or p.settings->>'gsc_integration_user_id' = '');

update public.projects p
set settings = jsonb_set(
  coalesce(p.settings, '{}'::jsonb),
  '{ga4_integration_user_id}',
  to_jsonb(p.owner_id),
  true
)
where (p.settings->>'ga4_property_id') is not null
  and (p.settings->>'ga4_property_id') <> ''
  and (p.settings->>'ga4_integration_user_id' is null or p.settings->>'ga4_integration_user_id' = '');
