# Barracuda SEO Product Roadmap

This document tracks planned features and enhancements for Barracuda SEO, organized by timeline and priority.

## Roadmap Structure

The roadmap is organized into four main sections:
1. **Completed Milestones** - Features that have been shipped
2. **In Progress** - Features currently being developed
3. **Upcoming Features** - Planned for next 1-3 months
4. **Mid-Term Goals** - Planned for 3-6 months
5. **Long-Term Vision** - Planned for 6-12+ months

---

## Section 3: Upcoming Features (1-3 months)

### 3.1 AI Integration
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

**Technical Notes:**
- Uses DataForSEO Google Maps SERP API
- Grid system allows flexible coverage area
- Cost optimization through selective grid point scanning
- Visual heatmap for easy identification of ranking zones
- AI analysis explains ranking factors for underperforming areas

### 3.5 Keyword Intelligence (Competitors + Content Gap)
- ⏳ Competitor ranked keyword discovery (domain/url reverse lookup)
- ⏳ Content gap recommendations (keywords competitors rank for that you don't)
- ⏳ Keyword metrics (search volume, CPC, competition) for prioritization
- ⏳ AI opportunity scoring: impact × difficulty × effort

**Technical Notes:**
- Uses DataForSEO Ranked Keywords API for reverse lookup
- Compares competitor keyword sets to identify gaps
- Integrates keyword metrics for prioritization
- AI scoring combines multiple factors for actionable insights

### 3.6 Decision Support & Explainability (JTBD-Driven)

*Informed by `docs/JTBD_SUMMARY.md` — turning overwhelm into defensible, explainable action plans.*

- 🔄 **Decision Rationale Panel** — For every prioritized issue/fix, surface: *(IN PROGRESS)*
  - Why this matters (impact framing)
  - What data informed the priority (GSC impressions, traffic, Clarity frustration, etc.)
  - What was intentionally deprioritized and why
  - Risk of not fixing it (consequence framing)
- ⏳ **Client-ready rationale export** — Copy/export rationale text for client calls and reports (language users can reuse)
- ⏳ **"Focus mode" / Top N view** — "400 issues don't matter. These 10 do." — curated shortlist that hides the long tail, with explicit "the rest can wait" framing
- ⏳ **Confident dismissal** — "Mark as low priority" or "Snooze" with optional rationale; reduces second-guessing when ignoring low-impact issues
- ⏳ **Onboarding anxiety copy** — In-product copy: "Barracuda doesn't replace your judgment. It gives you a clearer starting point — and shows its work."

**Why this matters (from JTBD):** Reduces anxiety, increases confidence, gives users language they can reuse with clients. Aligns with core positioning: "decision-support layer for SEO audits."

---

## Section 4: Mid-Term Goals (3-6 months)

### 4.1 Scheduled Crawls
- ⏳ Daily/weekly automation
- ⏳ Cloud Run cron triggers
- ⏳ Email & Slack summaries
- ⏳ Automatic comparisons between past crawls

### 4.2 Reporting & White-Labeling
- ⏳ PDF report generator
- ⏳ Change tracking between crawls
- ✅ Shareable links for client reporting (COMPLETE)
- ✅ Public-facing client reports (no login required) (COMPLETE)

### 4.3 Agency-Focused Enhancements
- ⏳ Multi-project GSC/GA dashboards
- ⏳ Project grouping/folders
- ⏳ Integrations with task managers (Jira, Asana, ClickUp)

### 4.4 Backlink Monitoring & Opportunities
- ⏳ Backlink profile overview + referring domains
- ⏳ New/lost backlinks + alerts
- ⏳ Competitor backlink overlap + link opportunity lists
- ⏳ AI suggestions for outreach targets + anchor strategy

**Technical Notes:**
- Integration with backlink data providers (e.g., DataForSEO Backlinks API, Ahrefs API, or similar)
- Real-time monitoring of backlink changes
- Competitor analysis to identify link opportunities
- AI-powered outreach recommendations based on competitor analysis
- Anchor text strategy suggestions based on competitor patterns

**Database Schema Considerations:**
- `backlinks` table: domain, source_url, target_url, anchor_text, discovered_at, lost_at
- `backlink_snapshots` table: periodic snapshots for change tracking
- `competitor_backlinks` table: competitor domain analysis
- `link_opportunities` table: AI-generated outreach targets

### 4.5 Performance Audits (Lighthouse + Core Web Vitals)
- ⏳ Lighthouse audits for key pages (perf, a11y, best practices, SEO)
- ⏳ Regression alerts when scores drop between scans
- ⏳ AI "fix-first" guidance tied to audit findings

**Technical Notes:**
- Integration with Lighthouse CI or Puppeteer for automated audits
- Core Web Vitals tracking (LCP, FID, CLS)
- Performance regression detection
- AI prioritization of performance fixes based on impact
- Integration with existing crawl data to correlate performance with SEO issues

**Implementation Approach:**
- Use Lighthouse Node.js API or Puppeteer
- Store audit results in `performance_audits` table
- Track metrics over time for regression detection
- Cross-reference with crawl data to identify performance-impacting SEO issues

---

## Section 5: Long-Term Vision (6-12+ months)

### 5.1 Rank Tracking Enhancements
- ⏳ Daily SERP position monitoring
- ⏳ Competitor comparisons
- ⏳ Keyword grouping & insights
- ⏳ Alerts for ranking drops

### 5.2 Intelligent SEO Assistant
- ⏳ Full conversational assistant inside dashboard
- ⏳ Auto-generated optimization briefs
- ⏳ Content audits and SEO scoring
- ⏳ Multi-domain insights for agencies

### 5.3 Plugin Ecosystem
- ⏳ Community-driven custom issue detectors
- ⏳ Webhooks & API extensions
- ⏳ Private plugins for enterprise clients

---

## Implementation Priorities

### High Priority (Next Sprint)
1. GA4 integration completion
2. Scheduled crawls foundation
3. CLI cloud sync improvements
4. **Decision Rationale Panel** (JTBD High ROI) — surfaces why each priority matters and what informed it

### Medium Priority (Next Quarter)
1. Decision Support & Explainability features (focus mode, client-ready export)
2. Local SEO geo-grid
3. Keyword intelligence features
4. Performance audits

### Lower Priority (Future)
1. Backlink monitoring
2. Plugin ecosystem
3. Advanced AI features

---

## Technical Considerations

### API Integrations Needed
- **DataForSEO**: 
  - Google Maps SERP API (for Local SEO)
  - Ranked Keywords API (for Keyword Intelligence)
  - Backlinks API (for Backlink Monitoring)
- **Lighthouse**: Node.js API or Puppeteer (for Performance Audits)
- **Google APIs**: GA4 Data API, GA4 Admin API (already integrated)

### Database Schema Additions
- Performance audit tables
- Backlink monitoring tables
- Local SEO geo-grid data tables
- Keyword intelligence cache tables

### Infrastructure Requirements
- Cloud Run cron jobs for scheduled crawls
- Background job processing for audits
- Caching layer for API responses
- Webhook endpoints for integrations

---

## Notes

- This roadmap is subject to change based on user feedback and business priorities
- **JTBD alignment:** Section 3.6 and related items are derived from `docs/JTBD_SUMMARY.md` — positioning Barracuda as a decision-support layer that reduces overwhelm and gives users explainable, defensible priorities
- Features marked with ✅ are complete
- Features marked with ⏳ are planned/upcoming
- Timeline estimates are approximate and may shift

---

**Last Updated:** January 2025
**Maintainer:** @dillonlara
