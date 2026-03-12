-- Fix crawl total_issues to match actual deduplicated issue count.
-- Previously total_issues was set from the pre-dedup analyzer output,
-- which inflated the number relative to what was actually stored.
UPDATE public.crawls c
SET total_issues = sub.actual_count
FROM (
    SELECT crawl_id, COUNT(*) AS actual_count
    FROM public.issues
    GROUP BY crawl_id
) sub
WHERE c.id = sub.crawl_id
  AND c.total_issues != sub.actual_count;
