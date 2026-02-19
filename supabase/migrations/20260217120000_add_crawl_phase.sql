-- Add phase column to show crawl progress stage in the UI
-- Values: scanning, metadata_review, image_analysis, storing, or null when complete
alter table public.crawls
  add column if not exists phase text;

comment on column public.crawls.phase is 'Current stage: scanning, metadata_review, image_analysis, storing. Null when complete.';
