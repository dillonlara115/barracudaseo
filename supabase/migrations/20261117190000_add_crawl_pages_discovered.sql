-- Crawl progress: how many pages the crawler has found so far.
--
-- A running crawl does not know how big the site is, only its page limit
-- (meta.max_pages), which made progress read "21 / 10000 pages" with an ETA of
-- hours for a 21-page site. pages_discovered is the number of distinct URLs
-- the crawler has queued, capped at the page limit: the best available
-- estimate of the final page count. It grows while links are still being
-- found and equals total_pages once scanning is done.

alter table public.crawls
  add column if not exists pages_discovered integer;
