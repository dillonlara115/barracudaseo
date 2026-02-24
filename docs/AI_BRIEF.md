# Barracuda SEO — AI-Powered Feature Suite
## Project Brief

---

## Overview

This document defines the architecture, features, and implementation plan for adding an AI-powered intelligence layer to Barracuda SEO. The goal is to evolve the platform from a suggestion-based tool into one that actively generates strategy, content briefs, and publish-ready articles — with the AI embedded directly into each feature rather than surfaced as a standalone chat interface.

Every AI capability is contextual and action-triggered. Users encounter intelligence at the moment they need it — a "Generate Brief" button on the keyword opportunity screen, an "Analyze Voice" button in project settings, a "Diagnose Decline" link on a declining page. The AI is a co-pilot within the existing workflow, not a parallel one.

---

## Product Philosophy

**AI baked in, not bolted on.** Each feature section exposes AI assistance at the right moment. There is no chat interface anywhere in the product.

**Data-grounded output.** Every AI response is seeded with real project data — GSC metrics, crawled site content, the project's writing voice profile. Generic output is a product failure.

**Review before anything ships.** All AI-generated content (briefs, articles) is presented for human review inside the app. No publishing pipeline exists in Phase 1. The WordPress integration is Phase 2, added only after content quality is validated through real usage.

**Team-aware from day one.** Multiple team members can collaborate on a project. Roles are enforced via Supabase RLS.

---

## Tech Stack

| Layer | Technology |
|---|---|
| Frontend | SvelteKit + Tailwind CSS + DaisyUI |
| Backend | Golang (Google Cloud Run) |
| Database | Supabase (PostgreSQL + pgvector) |
| Auth | Supabase Auth |
| AI — Generation (primary) | Gemini 2.5 Flash (`gemini-2.5-flash`) via Google AI API |
| AI — Generation (lightweight) | Gemini 2.5 Flash-Lite (`gemini-2.5-flash-lite`) via Google AI API |
| AI — Embeddings | Gemini `text-embedding-004` via Google AI API |
| GSC | Google Search Console API v3 |
| Icons | Lucide Svelte |
| Payments | Stripe |

### Why All-Gemini
Consolidating on a single AI provider (Google) means one API key, one billing relationship, one SDK, and one authentication flow that already overlaps with the GSC OAuth integration. 

- **Gemini 2.5 Flash** handles all primary generation tasks — content briefs, article writing, voice profile analysis, GSC diagnostics. Hybrid reasoning model with 1M token context window at $0.30/$2.50 per million tokens.
- **Gemini 2.5 Flash-Lite** handles lightweight, high-frequency tasks — memory extraction, weekly digest generation, simple CTR explanations, usage where response quality requirements are lower. At $0.10/$0.40 per million tokens, it's 3x cheaper for tasks that don't need full reasoning.
- **`text-embedding-004`** generates 768-dimension vectors for semantic search across crawled pages and internal link matching. Stored in Supabase via pgvector. Negligible cost at scale (~$0.00004 per 1K characters).

All three are accessed via the Google AI (Gemini) API using the same `GEMINI_API_KEY`.

---

## Core Concepts

### Projects
Each project represents a single client website. A project contains its own isolated data: GSC connection, crawled site content, writing voice profile, content briefs, generated articles, and AI memory. Agencies can manage multiple projects per plan tier.

### Team Membership
Projects support multiple users. A project has one owner (the user who created it) and can have additional members invited by email. Members can view and use all features. Only the owner can delete the project or manage billing. Enforced via Supabase RLS on all project-scoped tables.

### AI Context Injection
Every AI call is assembled by a shared context builder in Go. It pulls: the project's crawled page summaries (top N most relevant via pgvector similarity search), a snapshot of recent GSC performance data, the writing voice profile (if generated), and any relevant entries from project memory. This context is injected into the system prompt — not a chat history — so every feature gets the same grounded, project-aware intelligence regardless of which feature triggered it.

### Writing Voice Profile
A set of five structured documents generated once per project by analyzing the client's existing site content. Stored as text fields in Supabase. Loaded into any AI call that produces written content. Fully editable by the user. Components:

- **Tone** — formality level, personality traits, humor, emotional register
- **Structure** — how articles open, heading patterns, paragraph norms, CTA style
- **Sentence Style** — length, person (first/second/third), punctuation tendencies
- **Brand Context** — audience, expertise areas, products/services, positioning
- **Avoid List** — overused phrases, off-brand language, generic AI filler words to exclude

### Site Intelligence (Crawler)
A sitemap-based crawler (sitemap.xml only — no recursive link following in Phase 1) that fetches and stores clean content for each URL. Extracts: title, meta description, H1 and H2 headings, clean body text summary, internal links found on the page, and word count. Stored in Supabase with Gemini-generated embeddings for semantic search. The crawler is what makes the AI aware of what already exists on the site before generating any strategy or content.

### GSC Data Sync
GSC data syncs automatically once per day via a scheduled background job (Cloud Run Job or Cloud Scheduler). A manual "Sync Now" button is also available on the project dashboard for on-demand refreshes. Sync covers 90 days of query + page performance data. Stored in Supabase, updated incrementally.

---

## Subscription Tiers & Plan Limits

### Cost Basis
All generation uses Google Gemini API. Two model tiers are used to optimize cost per task type:

**Gemini 2.5 Flash** — $0.30/M input, $2.50/M output (primary generation)
- **Brief generation** — ~2,500 input + ~1,000 output tokens ≈ $0.0033
- **Article generation** — ~7,000 input + ~2,500 output tokens ≈ $0.0084
- **Voice profile generation** — ~5,000 input + ~1,500 output tokens ≈ $0.0053
- **Diagnostic action (explain/diagnose)** — ~2,000 input + ~600 output tokens ≈ $0.0021

**Gemini 2.5 Flash-Lite** — $0.10/M input, $0.40/M output (lightweight tasks)
- **Memory extraction** — ~1,500 input + ~300 output tokens ≈ $0.00027
- **Weekly digest** — ~3,000 input + ~800 output tokens ≈ $0.00062
- **Simple CTR/opportunity explanation** — ~1,000 input + ~300 output tokens ≈ $0.00022

At these rates, Gemini is roughly **7–10x cheaper than Claude Sonnet 4.6** for the same tasks, which means plan limits can be meaningfully more generous while maintaining healthy margins.

### Recommended Tiers

| Feature | Starter (Free) | Pro ($29/mo) | Agency ($79/mo) |
|---|---|---|---|
| Projects | 1 | 5 | 25 |
| Team members per project | 1 (owner only) | 3 | 10 |
| Crawl pages per project | 50 | 200 | 1,000 |
| GSC sync | Daily auto + manual | Daily auto + manual | Daily auto + manual |
| Content briefs / mo | 10 | 75 | 400 |
| Articles / mo | 3 | 25 | 150 |
| AI diagnostic actions / mo | 50 | unlimited | unlimited |
| Writing voice profiles | 1 per project | 1 per project | 1 per project |
| WordPress publishing (Phase 2) | — | ✓ | ✓ |

**Rationale:** Because Gemini 2.5 Flash is ~7–10x cheaper than Claude Sonnet 4.6, limits can be significantly more generous. A Pro user generating 25 articles + 75 briefs per month costs approximately $0.46 in AI tokens — under 2% of the $29 price point. Agency tier at full usage runs ~$1.75/month in AI costs against $79 revenue. Even accounting for infrastructure, embedding costs, and background jobs, margins are very healthy. Limits are set generously to drive value perception while remaining easy to enforce.

**Usage tracking:** Every AI call logs input tokens, output tokens, model, and feature to `ai_usage_log`. A usage dashboard in account settings shows the user their current month's usage against their plan limits. When a user approaches a limit (80%), show a soft warning. At 100%, the specific AI feature is disabled with an upgrade prompt — all non-AI features remain fully functional.

---

## Feature Modules

---

### 1. Project Setup & Onboarding

**Purpose:** Get a project configured and data-ready as quickly as possible.

**Flow:**
1. User creates a new project: enters site URL, display name
2. Google OAuth connects the project to a GSC property (one OAuth grants access to all properties — user selects from a dropdown)
3. App triggers GSC data sync (90 days of query + page data)
4. App triggers site crawl (sitemap-based, stores clean page content + Gemini embeddings to Supabase)
5. Optional but recommended: user initiates Writing Voice Analysis
6. Project is ready — all feature modules are active

**Supabase tables written during setup:**
- `projects`: id, owner_id, name, domain, gsc_property, created_at, plan_tier
- `project_members`: project_id, user_id, role (owner/member), invited_at, accepted_at
- `gsc_data`: project_id, date, query, page, clicks, impressions, ctr, position
- `crawled_pages`: project_id, url, title, meta_description, headings (jsonb), body_summary, internal_links (jsonb), word_count, embedding (vector(768)), crawled_at, manually_edited (bool)

---

### 2. Writing Voice Analysis

**Location:** Project Settings → Brand Voice

**Purpose:** Generate the writing voice profile so all AI-generated content sounds like the client, not generic AI output.

**UI:** A "Analyze Brand Voice" button with a short explanation of what it does. On click, shows a streaming progress state ("Reading homepage...", "Analyzing tone...", "Building voice profile..."). On completion, exposes five labeled text areas — one per voice component — that the user can review and edit inline. Changes save automatically.

**Backend process:**
1. Pull the crawled page content for the project — homepage + top pages by word count (up to 10 pages)
2. Assemble into a single prompt requesting each voice component as a clearly delimited output block
3. Stream the Claude response back to the frontend
4. Parse and store each component to `writing_voice` table on completion

**Key behaviors:**
- Runs on crawled content only — no external calls
- User edits are saved to Supabase immediately and used in all future AI generations
- Can be regenerated at any time (e.g., after a client rebrand)
- If no voice profile exists when an article is requested, the UI prompts the user to generate one first (soft warning, not a hard block)

**Supabase table:**
- `writing_voice`: id, project_id, tone (text), structure (text), sentence_style (text), brand_context (text), avoid_list (text), generated_at, last_edited_at

---

### 3. GSC Intelligence Dashboard

**Location:** Project dashboard / Overview

**Purpose:** Surface the most actionable insights from GSC data without requiring the user to manually dig through tables.

Each insight panel is triggered on demand by the user — nothing auto-runs on page load. Panels cache their last-generated result and show a timestamp ("Generated 2 hours ago · Refresh").

---

**Quick Wins Panel**
Keywords currently ranking positions 5–15 with high impressions but CTR below category average. Displayed as a sortable table: keyword, page, impressions, CTR, position, opportunity score (calculated server-side from those metrics — no AI needed for the table itself).

AI action on each row: **"Explain This Opportunity"** — generates a short paragraph explaining why this keyword is underperforming CTR-wise and what's likely suppressing it (title mismatch, weak meta description, SERP feature competition, etc.). Uses GSC data + crawled page content for that URL as context.

---

**Declining Pages Panel**
Pages with measurable position or traffic decline over the past 30/60/90 days (configurable toggle). Displayed as a table with delta indicators.

AI action on each row: **"Diagnose Decline"** — pulls the crawled content for that page and generates a diagnosis with likely causes ranked by probability: content staleness, keyword cannibalization, algorithm sensitivity, lost backlinks, technical issues. Includes a recommended next action.

---

**Content Gaps Panel**
Topics where GSC shows meaningful impressions but the site has no page with strong topical coverage. Identified by cross-referencing impression data clusters against the site crawl using pgvector similarity — the AI is used to generate the human-readable gap summary and recommended priority order, but the gap detection itself is done computationally.

Displayed as a list of gap topics with supporting impression data and a "Generate Brief" button on each row (links directly into the Content Brief Generator with the topic pre-filled).

---

**Weekly Digest**
A scheduled summary generated every Monday morning (via Cloud Scheduler) covering: biggest GSC movers (up and down) over the past 7 days, new keyword appearances, and one prioritized recommended action for the week. Delivered as an in-app notification and optionally via email. Generated using the batch API (50% cost discount since timing isn't user-facing real-time).

---

### 4. Content Brief Generator

**Location:** Content → Briefs, or triggered from Quick Wins / Content Gaps panels

**Purpose:** Generate a structured, data-backed content brief for a target keyword that accounts for what the site already covers.

**Trigger:** User selects a keyword opportunity or enters one manually, then clicks "Generate Brief." A brief preview panel opens with a loading state that streams the result.

**Brief output includes:**
- Target keyword + 3–5 supporting secondary keywords (sourced from GSC query data for semantically related terms)
- Recommended H1 title
- Meta description suggestion (under 160 characters)
- Recommended word count (based on GSC data for similar-performing pages on the site)
- Suggested H2 outline (6–10 sections) with a one-sentence description of what each section should cover
- Internal linking opportunities: specific pages from the site crawl that are topically relevant, with suggested anchor text
- Content angle notes: what the existing page (if any) is missing, what gaps this article should fill
- Keyword intent classification: informational / navigational / commercial / transactional

**Output format:** Rendered as a structured document in the app. Editable inline. Saveable as a draft. Can be used to trigger Article Writer.

**Supabase table:**
- `content_briefs`: id, project_id, keyword, brief_data (jsonb), status (draft / approved / archived), created_by (user_id), created_at, updated_at

---

### 5. Article Writer

**Location:** Content → Articles, triggered from an approved brief

**Purpose:** Generate a full, publish-ready article draft from an approved content brief, written in the project's voice.

**Trigger:** User opens an approved brief and clicks "Write Article." A pre-generation modal confirms what context will be used and displays warnings if any data is stale or missing (e.g., "Voice profile not generated — article will use default style guidelines"). User confirms to proceed.

**Backend generation process (Go):**
1. Load the content brief from Supabase
2. Load the writing voice profile for the project
3. pgvector similarity search on `crawled_pages` to find the most relevant pages for internal link context (top 5 by cosine similarity to the target keyword embedding)
4. Pull related GSC query data for keyword weaving guidance
5. Assemble the full system prompt via the Context Builder service
6. Send to Claude Sonnet 4.6 with structured generation instructions
7. Stream the response back to the SvelteKit frontend via SSE

**Output view:**
- Full article rendered in a clean reading/editing view (markdown with live preview)
- Internal link suggestions highlighted inline — each suggested link shows the target URL and recommended anchor text
- Word count + estimated read time displayed in the toolbar
- Export options: copy markdown, copy HTML, download as .md
- Status lifecycle: draft → reviewed → approved (user sets manually)

**Phase 1 hard stop:** No publishing. "Approved" status means the article is ready for the user to manually paste into their CMS. Phase 2 adds the WordPress push.

**Supabase table:**
- `articles`: id, project_id, brief_id, content (text), word_count, status (draft / reviewed / approved), created_by (user_id), created_at, updated_at

---

### 6. Internal Link Suggester

**Location:** Content → Internal Links (site-wide view) and on individual page detail views

**Purpose:** Surface internal linking opportunities the site is missing, using semantic similarity across crawled pages.

**Two entry points:**

**Page-level view:** User opens a specific crawled page and sees a "Suggested Links" section — a list of other pages on the site that are topically related, ranked by embedding cosine similarity. Each suggestion shows: target URL, page title, suggested anchor text (AI-generated on demand, not pre-generated for every pair), and similarity score.

**Site-wide audit view:** A full internal link health report:
- Orphaned pages (pages with zero inbound internal links from other crawled pages)
- Pages with excessive outbound links (over a configurable threshold)
- Topic clusters — groups of pages that are semantically related but not yet linked to each other (identified via pgvector clustering)

**Implementation note:** The core suggestion logic uses pgvector similarity search — no LLM call required for generating the list. Claude is only invoked when the user clicks "Suggest Anchor Text" on a specific pair, or requests an explanation of why two pages are related. This keeps costs minimal for what is essentially a data-matching feature.

---

### 7. Site Profile & Crawl Manager

**Location:** Project Settings → Site

**Purpose:** Transparency and control over what the app knows about the client's site.

**Views:**

**Sitemap view:** Full list of URLs from the sitemap with crawl status indicators (crawled / pending / failed / excluded), last crawled timestamp, and page-level metrics (word count, heading count, internal link count). Filterable and sortable.

**Page detail view:** Expandable per-page panel showing title, meta description, H1, H2s, word count, internal links found, crawl date. All fields are editable inline — user corrections are saved with `manually_edited = true` and are preserved across re-crawls.

**Crawl controls:**
- Re-crawl individual page button
- Re-crawl entire site button (with estimated time and page count)
- Manually add a URL not in the sitemap
- Exclude a URL from crawling

**Voice Profile section** is also surfaced here as a secondary entry point (primary is Project Settings → Brand Voice).

---

## AI Infrastructure (Shared Backend Services)

These are internal Go packages, not user-facing. All feature modules call into these.

### Context Builder (`/internal/ai/context`)
The single entry point for assembling any AI request. Takes a `ProjectID`, a `ContextType` enum (brief, article, diagnostic, voice), and optional parameters (e.g., target keyword, specific page URL). Returns a fully assembled system prompt with:
- Relevant crawled page summaries (pgvector similarity search, top 5 pages)
- GSC performance snapshot (filtered to relevant queries/pages)
- Writing voice profile (if available)
- Project memory entries (if relevant)
- Token budget enforcement — lower-priority context is truncated or summarized if the assembled prompt approaches the model's context limit

### Project Memory (`/internal/ai/memory`)
After significant AI interactions (brief generation, article generation, diagnostics), a lightweight follow-up call extracts key facts to persist about the project. Examples: "Primary content gap identified: indoor plant care guide," "Site has thin content on all category pages," "Client voice is casual and avoids jargon — never uses the word 'utilize'." Stored in `project_memory` table. Injected into future context builds as compressed background knowledge. Makes the AI progressively more tailored to a project over time without requiring chat history.

### Streaming Handler (`/internal/ai/stream`)
All AI calls that produce long-form output (articles, briefs, voice profiles) stream the response from the Gemini API back through the Go backend to SvelteKit via Server-Sent Events (SSE). The Go handler manages: opening the stream, forwarding chunks, handling errors and timeouts, and closing the connection cleanly. The SvelteKit frontend consumes the stream with `fetch` + `ReadableStream` — not `EventSource`, which is GET-only.

### Embedding Service (`/internal/ai/embeddings`)
Wraps the Gemini `text-embedding-004` API. Accepts text input, returns a 768-dimension float vector. Called by the crawler when storing new pages, and by the Context Builder when doing similarity search. Embeddings are stored in the `crawled_pages.embedding` column (pgvector `vector(768)` type). Similarity search uses cosine distance (`<=>` operator in pgvector).

### Usage Tracker (`/internal/ai/usage`)
Middleware that wraps every Claude API call. Logs to `ai_usage_log`: project_id, user_id, feature name, input tokens (from the API response), output tokens, model, cost estimate, created_at. Exposes a `CheckLimit(projectID, feature)` function called before each AI action — returns allowed/blocked based on current month's usage vs plan limits. Blocked calls return a structured error that the frontend translates into an upgrade prompt.

### Background Jobs
- **GSC Sync Job:** Runs daily via Cloud Scheduler. For each active project with a GSC connection, fetches the last 7 days of data (incremental, not full 90-day re-sync) and upserts to `gsc_data`. Full 90-day sync only on initial project setup or manual "full resync" trigger.
- **Weekly Digest Job:** Runs every Monday at 7am (user's timezone, defaulting to UTC). Generated using Gemini 2.5 Flash-Lite (cheapest model, timing is not real-time). Stores result to `digests` table. Triggers in-app notification and email (if opted in).

---

## Database Schema

```sql
-- Core project tables
projects
  id uuid PK, owner_id uuid FK(users), name text, domain text,
  gsc_property text, gsc_token (encrypted jsonb), plan_tier text,
  created_at timestamptz

project_members
  id uuid PK, project_id uuid FK, user_id uuid FK(users),
  role text CHECK (role IN ('owner','member')),
  invited_at timestamptz, accepted_at timestamptz

-- GSC data
gsc_data
  id uuid PK, project_id uuid FK, date date, query text, page text,
  clicks int, impressions int, ctr float, position float
  UNIQUE (project_id, date, query, page)

-- Site crawl
crawled_pages
  id uuid PK, project_id uuid FK, url text, title text,
  meta_description text, headings jsonb, body_summary text,
  internal_links jsonb, word_count int,
  embedding vector(768),   -- Gemini text-embedding-004
  crawled_at timestamptz, manually_edited bool DEFAULT false
  UNIQUE (project_id, url)

-- AI content
writing_voice
  id uuid PK, project_id uuid FK UNIQUE,
  tone text, structure text, sentence_style text,
  brand_context text, avoid_list text,
  generated_at timestamptz, last_edited_at timestamptz

content_briefs
  id uuid PK, project_id uuid FK, keyword text,
  brief_data jsonb, status text CHECK (status IN ('draft','approved','archived')),
  created_by uuid FK(users), created_at timestamptz, updated_at timestamptz

articles
  id uuid PK, project_id uuid FK, brief_id uuid FK(content_briefs),
  content text, word_count int,
  status text CHECK (status IN ('draft','reviewed','approved')),
  created_by uuid FK(users), created_at timestamptz, updated_at timestamptz

-- Memory and usage
project_memory
  id uuid PK, project_id uuid FK, memory_text text,
  source_feature text, created_at timestamptz

ai_usage_log
  id uuid PK, project_id uuid FK, user_id uuid FK(users),
  feature text, model text, input_tokens int, output_tokens int,
  cost_usd float, created_at timestamptz

-- Digests
digests
  id uuid PK, project_id uuid FK, content text,
  generated_at timestamptz, delivered_at timestamptz
```

**Indexes to create:**
- `crawled_pages.embedding` — `ivfflat` index for cosine similarity search
- `gsc_data (project_id, date)` — for time-range queries
- `ai_usage_log (project_id, created_at)` — for monthly usage aggregation
- `content_briefs (project_id, status)` — for filtered listing

---

## Go Backend Structure

```
/cmd
  /api          # Main HTTP server entrypoint
  /worker       # Background job entrypoints (GSC sync, digest)

/internal
  /ai
    /context    # Context builder — assembles prompts from project data
    /embeddings # Gemini embedding wrapper
    /memory     # Post-interaction memory extraction
    /stream     # SSE streaming handler for Claude responses
    /usage      # Usage tracking and limit enforcement

  /gsc          # Google Search Console API client + OAuth token refresh
  /crawler      # Sitemap parser, page fetcher, content extractor
  /projects     # Project CRUD, member management
  /briefs       # Content brief generation and management
  /articles     # Article generation and management
  /voice        # Writing voice profile generation and management
  /links        # Internal link suggestion logic (pgvector queries)
  /dashboard    # GSC intelligence panel data assembly

/pkg
  /db           # Supabase/postgres client, query helpers
  /config       # Environment config loading
  /middleware   # Auth, logging, rate limiting
```

---

## SvelteKit Frontend Structure

```
/src
  /routes
    /(app)                        # Authenticated layout
      /dashboard/[projectId]      # GSC intelligence panels
      /content/[projectId]
        /briefs                   # Brief list + generator
        /briefs/[briefId]         # Brief detail + article trigger
        /articles                 # Article list
        /articles/[articleId]     # Article editor/viewer
        /links                    # Internal link audit
      /site/[projectId]           # Crawl manager + page detail
      /settings/[projectId]       # Project settings, voice profile, team
      /account                    # Billing, usage, plan
    /(auth)
      /login
      /signup
      /onboarding                 # New project setup wizard

  /lib
    /components
      /ai                         # Shared AI result panels, streaming displays
      /gsc                        # GSC charts, data tables
      /content                    # Brief editor, article viewer, link panels
      /project                    # Project switcher, member management
      /ui                         # Generic DaisyUI wrappers (button, modal, etc.)
    /stores                       # Svelte stores (current project, user, usage)
    /utils                        # API client, formatting helpers
    /supabaseClient.ts
```

---

## Phased Rollout

### Phase 1 — This Brief
- Project setup + GSC OAuth + daily sync + manual sync
- Site crawler (sitemap-based, Gemini embeddings)
- Writing voice analysis + editor
- GSC intelligence dashboard (quick wins, declining pages, content gaps, weekly digest)
- Content brief generator
- Article writer (output to app only — no publishing)
- Internal link suggester
- Site profile + crawl manager
- Team/multi-user support
- Shared AI infrastructure (context builder, memory, streaming, embeddings, usage tracking)
- Subscription tiers with Stripe (Starter / Pro / Agency)

### Phase 2 — WordPress Integration
- WordPress REST API connection per project (credentials stored encrypted in Supabase)
- "Send to WordPress" button on approved articles
- Creates a draft post — title, content, meta fields (Yoast/RankMath via REST)
- Sync status: draft in WP is linked back to the article record in Supabase
- Image placeholder support: `[IMAGE: suggested alt text]` markers for manual media addition
- Re-send capability if the draft was edited locally and needs to be re-pushed

### Phase 3 — Potential Future Expansion
- Scheduled content calendar view
- Competitor content gap analysis
- Schema markup suggestion tool
- Programmatic SEO tooling for large sites (template-based page generation)
- Additional CMS integrations (Webflow, Framer, Contentful)

---

## What This Is Not

- **Not a chat interface.** There is no open-ended chat box anywhere in the product.
- **Not a GSC replacement.** It reads and interprets GSC data — it doesn't replace the GSC console for deep manual analysis.
- **Not autonomous.** Every AI output requires a human to review and approve before anything moves to the next stage.
- **Not a clone.** Architecture, UX, data model, and codebase are built from scratch in the Barracuda stack. Inspired by patterns in the space, not derived from any specific codebase.

---

## Environment Variables Required

```bash
# Supabase
SUPABASE_URL=
SUPABASE_SERVICE_ROLE_KEY=

# Google (GSC OAuth + Gemini Generation + Embeddings — single API key)
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT_URI=
GEMINI_API_KEY=

# Stripe
STRIPE_SECRET_KEY=
STRIPE_WEBHOOK_SECRET=

# App
APP_URL=
JWT_SECRET=
ENVIRONMENT= # development | production
```