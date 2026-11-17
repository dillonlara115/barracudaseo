-- Crawl heartbeat: lets the API tell a crawl that is still working from one
-- whose runner died.
--
-- The crawl runner (a Cloud Run Job, or the API process in local development)
-- writes heartbeat_at every 30 seconds while it works. The scheduled-crawl cron
-- marks a running crawl as failed when its heartbeat is more than 10 minutes
-- old, so a crawl is never left "running" forever.

alter table public.crawls
  add column if not exists heartbeat_at timestamptz;

-- The sweeper looks up running crawls by heartbeat age
create index if not exists idx_crawls_running_heartbeat
  on public.crawls (heartbeat_at)
  where status = 'running';
