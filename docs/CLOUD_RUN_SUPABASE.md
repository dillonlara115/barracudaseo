# Cloud Run + Supabase + Vercel Deployment Blueprint

This document captures the target hosted architecture for Barracuda when splitting responsibilities across Google Cloud Run (backend), Supabase (database & auth), and Vercel (Svelte front-end). Use it as the single source of truth for infrastructure decisions and onboarding.

---

## High-Level Overview

- **Cloud Run (Backend API)**
  - Containerized Go service handling crawl ingestion, analysis orchestration, and authenticated REST/gRPC endpoints.
  - Triggered by CLI uploads, scheduled tasks, or UI actions.
  - Uses Google Secret Manager for third-party API keys (OpenAI, Search Console, etc.).

- **Supabase (Data Platform)**
  - Managed Postgres storing crawls, issues, pages, user preferences, and job history.
  - Supabase Auth handles user sign-in; row-level security protects tenant data.
  - Supabase Storage optionally holds large artifacts (raw exports, screenshots).

- **Vercel (Front-End)**
  - Hosts the Svelte dashboard.
  - Fetches data via Supabase client (subject to policies) or via Cloud Run API.
  - Leverages Vercel edge/serverless functions for lightweight adapters if required.

---

## Component Responsibilities

| Component | Responsibilities | Key Integrations |
|-----------|------------------|------------------|
| Cloud Run service | - Receive crawl results from CLI or scheduled jobs<br>- Write normalized data into Supabase<br>- Expose authenticated APIs for dashboard actions<br>- Run background analysis (OpenAI summaries, recommendations)<br>- Fan-out to Google APIs (Search Console) | Supabase service role key, Secret Manager, Pub/Sub (optional) |
| Supabase | - Primary relational datastore<br>- Auth provider (email/OAuth)<br>- Realtime subscriptions for dashboard updates<br>- Row-level security enforcement | Supabase client in web app, Go Supabase client (via PostgREST) |
| Vercel dashboard | - SPA for viewing issues, managing filters, triggering re-crawls<br>- Uses Supabase auth to sign in users<br>- Calls Cloud Run APIs with JWT from Supabase | Supabase `anon` key, Cloud Run public endpoint |

---

## Data Flow

1. **Crawl Ingestion**
   - CLI or Cloud Scheduler triggers a crawl and sends results to Cloud Run.
   - Cloud Run validates payload, enriches data, and performs bulk inserts into Supabase using the service role key.

2. **Dashboard Consumption**
   - User authenticates via Supabase on the Svelte app (hosted on Vercel).
   - UI fetches filtered data directly from Supabase using row-level policies, or via Cloud Run endpoints for derived metrics (aggregations, AI-generated insights).

3. **Recommendations and Integrations**
   - Cloud Run workers call OpenAI or Google Search Console APIs using secrets stored in Secret Manager.
   - Outputs are persisted back into Supabase and streamed to clients via Supabase realtime.

4. **Exports**
   - Vercel UI requests export jobs; Cloud Run generates files and saves them to Supabase Storage or Google Cloud Storage, returning a signed URL.

---

## Deployment Workflow

1. **Build & Push Containers**
   - Use `make deploy-backend` (planned) to build the Go image and push to Artifact Registry.
   - Deploy to Cloud Run with infrastructure-as-code (Terraform) or `gcloud run deploy`.

2. **Database Migration**
   - Store schema migrations in a new `db/migrations` directory.
   - Apply via Supabase CLI or `golang-migrate` against the Supabase instance.

3. **Front-End Release**
   - Vercel auto-builds on `main` merges; environment variables reference Cloud Run API base URLs and Supabase anon key.

4. **Secrets Management**
   - Cloud Run: configure Secret Manager -> environment variables.
   - Vercel: set public (anon) and private keys through the dashboard.
   - CLI: read from environment variables (or `.env` via `BARRACUDA_LOAD_ENV=1`), send signed requests to Cloud Run using Supabase auth tokens.

---

## Authentication Strategy

- Supabase Auth issues JWTs post sign-in.
- Svelte app stores session via Supabase client.
- Requests to Cloud Run include Supabase JWT in `Authorization: Bearer` header.
- Cloud Run validates token using Supabase JWKS and enforces role-based access.
- Service-to-service operations (CLI ingestion) use Supabase service-role key stored in Secret Manager and short-lived signed tokens.

---

## CLI Integration Notes

- CLI keeps operating locally; after each crawl it POSTs to Cloud Run.
- Provide a `--cloud` flag that bundles results into NDJSON/JSON for upload.
- On failure, CLI falls back to saving artifacts locally to retry later.
- Future: CLI can poll Supabase for job status or stream logs via WebSocket.

---

## Future Enhancements

- Add Cloud Scheduler + Pub/Sub to trigger periodic crawls.
- Introduce Supabase Edge Functions for lightweight, latency-sensitive operations near the database.
- Evaluate Cloud Run Jobs for long-running analyses or reprocessing tasks.
- Consider Terraform modules to codify infrastructure and simplify new environment bootstraps (staging, prod).

---

**Last Updated:** {{ date }}
