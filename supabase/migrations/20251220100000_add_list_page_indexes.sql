-- Add indexes to speed up list queries for issues, pages, and projects.

create index if not exists idx_issues_crawl_created_at
  on public.issues (crawl_id, created_at desc);

create index if not exists idx_pages_crawl_created_at
  on public.pages (crawl_id, created_at desc);

create index if not exists idx_projects_created_at
  on public.projects (created_at desc);
