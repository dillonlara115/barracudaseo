# Barracuda SEO Product Roadmap

This document tracks planned features and enhancements for Barracuda SEO, organized by timeline and priority.

## Roadmap Structure

1. **Completed Milestones** - Features that have been shipped
2. **In Progress** - Features currently being developed
3. **Upcoming Features** - Planned for next 1-3 months
4. **Mid-Term Goals** - Planned for 3-6 months
5. **Long-Term Vision** - Planned for 6-12+ months

---

## Section 2: In Progress

### 2.1 AI-Powered Feature Suite (Phase 1)

*Full spec: `docs/AI_BRIEF.md` — AI baked in, not bolted on.*

The AI Suite is the highest-priority initiative. Every feature below embeds intelligence at the point of action — no chat interface, no separate AI tab. Powered by Gemini 2.5 Flash (primary), Gemini 2.5 Flash-Lite (lightweight), and `text-embedding-004` (embeddings via pgvector).

#### Shared AI Infrastructure
- 🔄 Context Builder — assembles prompts from crawled pages (pgvector similarity), GSC data, voice profile, project memory
- 🔄 Embedding Service — Gemini `text-embedding-004` wrapper (768-dim vectors stored in pgvector)
- 🔄 SSE Streaming Handler — Go → SvelteKit via `fetch` + `ReadableStream`
- 🔄 Project Memory — post-interaction fact extraction for progressively tailored AI
- 🔄 Usage Tracker — per-feature token logging, plan limit enforcement, upgrade prompts

#### Writing Voice Analysis
- 🔄 Analyze brand voice from crawled site content (tone, structure, sentence style, brand context, avoid list)
- 🔄 Editable voice profile stored per project, injected into all content generation

#### GSC Intelligence Dashboard
- 🔄 Quick Wins Panel — keywords ranking 5–15 with high impressions, low CTR; AI "Explain This Opportunity" per row
- 🔄 Declining Pages Panel — position/traffic deltas with AI "Diagnose Decline" per row
- 🔄 Content Gaps Panel — pgvector cross-reference against GSC impressions; "Generate Brief" trigger per gap
- 🔄 Weekly Digest — scheduled Monday summary via Cloud Scheduler (Gemini Flash-Lite)

#### Content Brief Generator
- 🔄 Data-backed briefs: target keyword, secondary keywords, H1/meta/outline, internal link opportunities, intent classification
- 🔄 Inline editing, status lifecycle (draft → approved → archived)

#### Article Writer
- 🔄 Full article generation from approved briefs, written in the project's voice profile
- 🔄 Markdown editor with live preview, internal link highlights, export (MD/HTML)
- 🔄 No publishing pipeline in Phase 1 — "Approved" means ready for manual CMS paste

#### Internal Link Suggester
- 🔄 Page-level: pgvector cosine similarity for related pages, AI anchor text suggestion on demand
- 🔄 Site-wide audit: orphaned pages, excessive outbound links, topic clusters

#### Site Crawler Enhancement
- 🔄 Sitemap-based crawling with clean content storage (title, meta, headings, body summary, internal links, word count)
- 🔄 Gemini embedding generation per crawled page
- 🔄 Editable page data preserved across re-crawls

#### Background Jobs
- 🔄 Daily GSC sync (incremental 7-day; full 90-day on setup only)
- 🔄 Weekly digest generation (Monday 7am UTC, Gemini Flash-Lite)

#### Usage Tracking UI
- 🔄 Account dashboard: current month's usage vs plan limits per feature
- 🔄 Soft warning at 80%, feature-level block at 100% with upgrade prompt

#### Database Additions
- New tables: `writing_voice`, `content_briefs`, `articles`, `project_memory`, `ai_usage_log`, `digests`
- pgvector extension + `embedding vector(768)` on `crawled_pages` with ivfflat index
- New indexes on `gsc_data`, `content_briefs`, `ai_usage_log`

---

## Section 3: Upcoming Features (1-3 months)

### 3.1 AI Integration (Existing)
- ✅ AI recommendations for issues (COMPLETE)
- ✅ AI crawl summary reports (COMPLETE)
- ✅ AI insights combining crawl + GSC (COMPLETE)
- ⏳ Page-level rewrite suggestions
- ⏳ AI-powered priority scoring

### 3.2 Additional Integrations
- ⏳ Google Analytics 4 (GA4)
- ⏳ Google Drive / Google Sheets export
- ⏳ Slack alerts
- ✅ Microsoft Clarity integration (COMPLETE)
- ⏳ Zapier / Make.com workflows

### 3.3 CLI → Cloud Sync Improvements
- ⏳ `barracuda auth login`
- ⏳ `barracuda projects list`
- ⏳ `--cloud` ingestion polish
- ⏳ Auto-associate crawls to default project

### 3.4 Local SEO (Geo-Grid + Maps Rankings)
- ⏳ Geo-grid rank scans (3×3 to 21×21) using Google Maps SERP
- ⏳ Adjustable radius + grid spacing (local focus control)
- ⏳ Disable grid points to optimize API credits
- ⏳ Interactive heatmap visualization + shareable map links
- ⏳ AI explanations for weak zones (proximity, competition, reviews, categories)

### 3.5 Keyword Intelligence (Competitors + Content Gap)
- ⏳ Competitor ranked keyword discovery (domain/url reverse lookup)
- ⏳ Content gap recommendations (keywords competitors rank for that you don't)
- ⏳ Keyword metrics (search volume, CPC, competition) for prioritization
- ⏳ AI opportunity scoring: impact × difficulty × effort

### 3.6 Decision Support & Explainability (JTBD-Driven)

*Informed by `docs/JTBD_SUMMARY.md` — turning overwhelm into defensible, explainable action plans.*

- 🔄 **Decision Rationale Panel** — For every prioritized issue/fix, surface: *(IN PROGRESS)*
  - Why this matters (impact framing)
  - What data informed the priority (GSC impressions, traffic, Clarity frustration, etc.)
  - What was intentionally deprioritized and why
  - Risk of not fixing it (consequence framing)
- ⏳ **Client-ready rationale export** — Copy/export rationale text for client calls and reports
- ⏳ **"Focus mode" / Top N view** — curated shortlist that hides the long tail
- ⏳ **Confident dismissal** — "Mark as low priority" or "Snooze" with optional rationale
- ⏳ **Onboarding anxiety copy** — In-product reassurance messaging

---

## Section 4: Mid-Term Goals (3-6 months)

### 4.1 AI Suite — Phase 2: WordPress Integration
- ⏳ WordPress REST API connection per project (encrypted credentials in Supabase)
- ⏳ "Send to WordPress" button on approved articles
- ⏳ Creates draft post with title, content, meta fields (Yoast/RankMath via REST)
- ⏳ Sync status: WP draft linked back to article record in Supabase
- ⏳ Image placeholder support: `[IMAGE: suggested alt text]` markers
- ⏳ Re-send capability for locally edited drafts

### 4.2 Scheduled Crawls
- ⏳ Daily/weekly automation
- ⏳ Cloud Run cron triggers
- ⏳ Email & Slack summaries
- ⏳ Automatic comparisons between past crawls

### 4.3 Reporting & White-Labeling
- ⏳ PDF report generator
- ⏳ Change tracking between crawls
- ✅ Shareable links for client reporting (COMPLETE)
- ✅ Public-facing client reports (no login required) (COMPLETE)

### 4.4 Agency-Focused Enhancements
- ⏳ Multi-project GSC/GA dashboards
- ⏳ Project grouping/folders
- ⏳ Integrations with task managers (Jira, Asana, ClickUp)

### 4.5 Backlink Monitoring & Opportunities
- ⏳ Backlink profile overview + referring domains
- ⏳ New/lost backlinks + alerts
- ⏳ Competitor backlink overlap + link opportunity lists
- ⏳ AI suggestions for outreach targets + anchor strategy

### 4.6 Performance Audits (Lighthouse + Core Web Vitals)
- ⏳ Lighthouse audits for key pages (perf, a11y, best practices, SEO)
- ⏳ Regression alerts when scores drop between scans
- ⏳ AI "fix-first" guidance tied to audit findings

---

## Section 5: Long-Term Vision (6-12+ months)

### 5.1 AI Suite — Phase 3: Future Expansion
- ⏳ Scheduled content calendar view
- ⏳ Competitor content gap analysis (AI-driven)
- ⏳ Schema markup suggestion tool
- ⏳ Programmatic SEO tooling for large sites (template-based page generation)
- ⏳ Additional CMS integrations (Webflow, Framer, Contentful)

### 5.2 Rank Tracking Enhancements
- ⏳ Daily SERP position monitoring
- ⏳ Competitor comparisons
- ⏳ Keyword grouping & insights
- ⏳ Alerts for ranking drops

### 5.3 Intelligent SEO Assistant
- ⏳ Full conversational assistant inside dashboard
- ⏳ Auto-generated optimization briefs
- ⏳ Content audits and SEO scoring
- ⏳ Multi-domain insights for agencies

### 5.4 Plugin Ecosystem
- ⏳ Community-driven custom issue detectors
- ⏳ Webhooks & API extensions
- ⏳ Private plugins for enterprise clients

---

## Implementation Priorities

### High Priority (Current Sprint)
1. **AI-Powered Feature Suite (Phase 1)** — full AI infrastructure, voice analysis, GSC intelligence, content briefs, article writer, internal links
2. **Decision Rationale Panel** (JTBD High ROI) — surfaces why each priority matters and what informed it

### Medium Priority (Next Quarter)
1. WordPress integration (AI Suite Phase 2)
2. Decision Support & Explainability features (focus mode, client-ready export)
3. Local SEO geo-grid
4. Keyword intelligence features

### Lower Priority (Future)
1. Scheduled crawls
2. Backlink monitoring
3. Performance audits
4. Plugin ecosystem

---

## Technical Considerations

### AI Infrastructure (New)
- **Gemini 2.5 Flash** — primary generation ($0.30/$2.50 per M tokens)
- **Gemini 2.5 Flash-Lite** — lightweight tasks ($0.10/$0.40 per M tokens)
- **text-embedding-004** — 768-dim vectors for semantic search (~$0.00004 per 1K chars)
- **pgvector** — cosine distance (`<=>`) on `crawled_pages.embedding`
- **SSE streaming** — Go backend → SvelteKit frontend via `fetch` + `ReadableStream`

### API Integrations Needed
- **Google AI (Gemini)**: Generation + Embeddings (single `GEMINI_API_KEY`)
- **DataForSEO**: Google Maps SERP API, Ranked Keywords API, Backlinks API
- **Lighthouse**: Node.js API or Puppeteer (for Performance Audits)
- **Google APIs**: GA4 Data API, GA4 Admin API (already integrated)

### Database Schema Additions
- AI tables: `writing_voice`, `content_briefs`, `articles`, `project_memory`, `ai_usage_log`, `digests`
- pgvector extension + embedding column + ivfflat index
- Performance audit tables
- Backlink monitoring tables
- Local SEO geo-grid data tables

### Infrastructure Requirements
- Cloud Run cron jobs for GSC sync + weekly digest
- Background job processing for audits
- Caching layer for API responses
- Webhook endpoints for integrations

---

## Notes

- This roadmap is subject to change based on user feedback and business priorities
- **AI Suite** is the current top priority — see `docs/AI_BRIEF.md` for full specification
- **JTBD alignment:** Section 3.6 and related items are derived from `docs/JTBD_SUMMARY.md`
- Features marked with ✅ are complete
- Features marked with 🔄 are in progress
- Features marked with ⏳ are planned/upcoming
- Timeline estimates are approximate and may shift

---

**Last Updated:** February 2026
**Maintainer:** @dillonlara
