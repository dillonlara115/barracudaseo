-- AI Suite Phase 1: Enable pgvector and create AI-related tables
-- See docs/AI_BRIEF.md for full specification

-- Enable pgvector extension for semantic search
CREATE EXTENSION IF NOT EXISTS vector;

-- Add embedding column to crawled_pages (if crawled_pages exists; otherwise create it)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'crawled_pages') THEN
        IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'crawled_pages' AND column_name = 'embedding') THEN
            ALTER TABLE crawled_pages ADD COLUMN embedding vector(768);
        END IF;
        IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'crawled_pages' AND column_name = 'body_summary') THEN
            ALTER TABLE crawled_pages ADD COLUMN body_summary text;
        END IF;
        IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'crawled_pages' AND column_name = 'headings') THEN
            ALTER TABLE crawled_pages ADD COLUMN headings jsonb;
        END IF;
        IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'crawled_pages' AND column_name = 'internal_links') THEN
            ALTER TABLE crawled_pages ADD COLUMN internal_links jsonb;
        END IF;
        IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'crawled_pages' AND column_name = 'word_count') THEN
            ALTER TABLE crawled_pages ADD COLUMN word_count integer DEFAULT 0;
        END IF;
        IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'crawled_pages' AND column_name = 'manually_edited') THEN
            ALTER TABLE crawled_pages ADD COLUMN manually_edited boolean DEFAULT false;
        END IF;
    ELSE
        CREATE TABLE crawled_pages (
            id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
            project_id uuid NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
            url text NOT NULL,
            title text,
            meta_description text,
            headings jsonb,
            body_summary text,
            internal_links jsonb,
            word_count integer DEFAULT 0,
            embedding vector(768),
            crawled_at timestamptz DEFAULT now(),
            manually_edited boolean DEFAULT false,
            UNIQUE (project_id, url)
        );
    END IF;
END $$;

-- Writing voice profile per project
CREATE TABLE IF NOT EXISTS writing_voice (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id uuid NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    tone text,
    structure text,
    sentence_style text,
    brand_context text,
    avoid_list text,
    generated_at timestamptz DEFAULT now(),
    last_edited_at timestamptz DEFAULT now(),
    UNIQUE (project_id)
);

-- Content briefs
CREATE TABLE IF NOT EXISTS content_briefs (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id uuid NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    keyword text NOT NULL,
    brief_data jsonb,
    status text NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'approved', 'archived')),
    created_by uuid REFERENCES auth.users(id),
    created_at timestamptz DEFAULT now(),
    updated_at timestamptz DEFAULT now()
);

-- Articles generated from briefs
CREATE TABLE IF NOT EXISTS articles (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id uuid NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    brief_id uuid REFERENCES content_briefs(id) ON DELETE SET NULL,
    content text,
    word_count integer DEFAULT 0,
    status text NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'reviewed', 'approved')),
    created_by uuid REFERENCES auth.users(id),
    created_at timestamptz DEFAULT now(),
    updated_at timestamptz DEFAULT now()
);

-- Project memory for progressively tailored AI
CREATE TABLE IF NOT EXISTS project_memory (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id uuid NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    memory_text text NOT NULL,
    source_feature text NOT NULL,
    created_at timestamptz DEFAULT now()
);

-- AI usage log for tracking and billing
CREATE TABLE IF NOT EXISTS ai_usage_log (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id uuid NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id uuid REFERENCES auth.users(id),
    feature text NOT NULL,
    model text,
    input_tokens integer DEFAULT 0,
    output_tokens integer DEFAULT 0,
    cost_usd double precision DEFAULT 0,
    created_at timestamptz DEFAULT now()
);

-- Weekly digest summaries
CREATE TABLE IF NOT EXISTS digests (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id uuid NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    content text,
    generated_at timestamptz DEFAULT now(),
    delivered_at timestamptz
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_crawled_pages_embedding
    ON crawled_pages USING hnsw (embedding vector_cosine_ops);

CREATE INDEX IF NOT EXISTS idx_crawled_pages_project_url
    ON crawled_pages (project_id, url);

-- gsc_performance_rows already has indexes from its original migration

CREATE INDEX IF NOT EXISTS idx_content_briefs_project_status
    ON content_briefs (project_id, status);

CREATE INDEX IF NOT EXISTS idx_ai_usage_log_project_created
    ON ai_usage_log (project_id, created_at);

CREATE INDEX IF NOT EXISTS idx_project_memory_project
    ON project_memory (project_id);

CREATE INDEX IF NOT EXISTS idx_articles_project_status
    ON articles (project_id, status);

CREATE INDEX IF NOT EXISTS idx_digests_project
    ON digests (project_id);

-- RLS policies

ALTER TABLE writing_voice ENABLE ROW LEVEL SECURITY;
ALTER TABLE content_briefs ENABLE ROW LEVEL SECURITY;
ALTER TABLE articles ENABLE ROW LEVEL SECURITY;
ALTER TABLE project_memory ENABLE ROW LEVEL SECURITY;
ALTER TABLE ai_usage_log ENABLE ROW LEVEL SECURITY;
ALTER TABLE digests ENABLE ROW LEVEL SECURITY;

-- Writing voice: accessible by project members
CREATE POLICY "writing_voice_select" ON writing_voice
    FOR SELECT USING (
        project_id IN (
            SELECT pm.project_id FROM project_members pm WHERE pm.user_id = auth.uid()
        )
    );

CREATE POLICY "writing_voice_insert" ON writing_voice
    FOR INSERT WITH CHECK (
        project_id IN (
            SELECT pm.project_id FROM project_members pm WHERE pm.user_id = auth.uid()
        )
    );

CREATE POLICY "writing_voice_update" ON writing_voice
    FOR UPDATE USING (
        project_id IN (
            SELECT pm.project_id FROM project_members pm WHERE pm.user_id = auth.uid()
        )
    );

-- Content briefs: accessible by project members
CREATE POLICY "content_briefs_select" ON content_briefs
    FOR SELECT USING (
        project_id IN (
            SELECT pm.project_id FROM project_members pm WHERE pm.user_id = auth.uid()
        )
    );

CREATE POLICY "content_briefs_insert" ON content_briefs
    FOR INSERT WITH CHECK (
        project_id IN (
            SELECT pm.project_id FROM project_members pm WHERE pm.user_id = auth.uid()
        )
    );

CREATE POLICY "content_briefs_update" ON content_briefs
    FOR UPDATE USING (
        project_id IN (
            SELECT pm.project_id FROM project_members pm WHERE pm.user_id = auth.uid()
        )
    );

-- Articles: accessible by project members
CREATE POLICY "articles_select" ON articles
    FOR SELECT USING (
        project_id IN (
            SELECT pm.project_id FROM project_members pm WHERE pm.user_id = auth.uid()
        )
    );

CREATE POLICY "articles_insert" ON articles
    FOR INSERT WITH CHECK (
        project_id IN (
            SELECT pm.project_id FROM project_members pm WHERE pm.user_id = auth.uid()
        )
    );

CREATE POLICY "articles_update" ON articles
    FOR UPDATE USING (
        project_id IN (
            SELECT pm.project_id FROM project_members pm WHERE pm.user_id = auth.uid()
        )
    );

-- Project memory: accessible by project members (read-only for users; insert via service role)
CREATE POLICY "project_memory_select" ON project_memory
    FOR SELECT USING (
        project_id IN (
            SELECT pm.project_id FROM project_members pm WHERE pm.user_id = auth.uid()
        )
    );

-- AI usage log: accessible by project members (read-only for users; insert via service role)
CREATE POLICY "ai_usage_log_select" ON ai_usage_log
    FOR SELECT USING (
        project_id IN (
            SELECT pm.project_id FROM project_members pm WHERE pm.user_id = auth.uid()
        )
    );

-- Digests: accessible by project members
CREATE POLICY "digests_select" ON digests
    FOR SELECT USING (
        project_id IN (
            SELECT pm.project_id FROM project_members pm WHERE pm.user_id = auth.uid()
        )
    );
