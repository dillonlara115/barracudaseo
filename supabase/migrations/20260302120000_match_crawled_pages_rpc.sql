-- Add RPC function for pgvector cosine similarity search on crawled_pages.
-- Used by the AI context builder to retrieve the most semantically relevant
-- pages for a given query embedding, replacing the word_count fallback.

CREATE OR REPLACE FUNCTION match_crawled_pages(
    query_embedding vector(768),
    match_project_id uuid,
    match_count int DEFAULT 5,
    match_threshold float DEFAULT 0.0
)
RETURNS TABLE (
    url text,
    title text,
    meta_description text,
    headings jsonb,
    body_summary text,
    word_count int,
    similarity float
)
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = public
AS $$
BEGIN
    RETURN QUERY
    SELECT
        cp.url,
        cp.title,
        cp.meta_description,
        cp.headings,
        cp.body_summary,
        cp.word_count,
        1 - (cp.embedding <=> query_embedding) AS similarity
    FROM crawled_pages cp
    WHERE cp.project_id = match_project_id
      AND cp.embedding IS NOT NULL
      AND 1 - (cp.embedding <=> query_embedding) > match_threshold
    ORDER BY cp.embedding <=> query_embedding
    LIMIT match_count;
END;
$$;
