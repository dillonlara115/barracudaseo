-- Backfill page_id for issues that are missing it
-- This matches issues to pages by finding pages with matching URLs in the same crawl
-- Note: This is a best-effort match and may not work for all cases

-- First, let's see how many issues are missing page_id
SELECT 
  crawl_id,
  COUNT(*) as issues_without_page_id
FROM issues
WHERE page_id IS NULL
GROUP BY crawl_id;

-- Update issues to set page_id where we can match by URL pattern
-- This uses a subquery to find pages that might match based on the issue message/value
-- Note: This is a simplified approach - for production, you might want a more sophisticated matching algorithm

UPDATE issues i
SET page_id = (
  SELECT p.id
  FROM pages p
  WHERE p.crawl_id = i.crawl_id
    AND (
      -- Try to match if the issue message contains the page URL
      i.message LIKE '%' || p.url || '%'
      -- Or if the issue value contains the page URL
      OR (i.value IS NOT NULL AND i.value LIKE '%' || p.url || '%')
    )
  LIMIT 1
)
WHERE i.page_id IS NULL
  AND EXISTS (
    SELECT 1
    FROM pages p
    WHERE p.crawl_id = i.crawl_id
      AND (
        i.message LIKE '%' || p.url || '%'
        OR (i.value IS NOT NULL AND i.value LIKE '%' || p.url || '%')
      )
  );

-- Show results
SELECT 
  crawl_id,
  COUNT(*) as issues_with_page_id_after_update
FROM issues
WHERE page_id IS NOT NULL
GROUP BY crawl_id;
