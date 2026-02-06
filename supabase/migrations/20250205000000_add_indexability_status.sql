-- Add indexability_status column to pages table
-- This tracks whether a page is indexable, has noindex directive, or is blocked by robots.txt

alter table public.pages
add column if not exists indexability_status text check (indexability_status in ('indexable', 'noindex', 'blocked'));

-- Add index for filtering by indexability status
create index if not exists idx_pages_indexability_status on public.pages (crawl_id, indexability_status);

-- Add comment explaining the column
comment on column public.pages.indexability_status is 'Indexability status: indexable (can be indexed), noindex (has noindex directive), or blocked (blocked by robots.txt)';
