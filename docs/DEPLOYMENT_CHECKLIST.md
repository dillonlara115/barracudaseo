# Deployment Checklist: Cloud Run + Supabase + Vercel

This checklist tracks the steps needed to deploy Barracuda to production using Google Cloud Run (backend), Supabase (database/auth), and Vercel (frontend).

## ✅ Completed

- [x] Supabase schema migration created (`supabase/migrations/20240320_initial_schema.sql`)
- [x] Schema applied to Supabase project

## 🚧 In Progress

- [ ] Dockerfile for Cloud Run
- [ ] Cloud Run API server implementation

## 📋 Remaining Tasks

### Phase 1: Backend Infrastructure (Cloud Run)

- [ ] **Dockerfile**
  - Multi-stage build (Go build + minimal runtime image)
  - Optimize for Cloud Run (small image size, fast startup)
  - Health check endpoint

- [ ] **API Server (`cmd/api.go` or `cmd/server.go`)**
  - Separate from `serve` command (which is for local dev)
  - REST API endpoints:
    - `POST /api/v1/crawls` - Ingest crawl results from CLI
    - `GET /api/v1/projects` - List user's projects
    - `GET /api/v1/projects/:id/crawls` - List crawls for a project
    - `GET /api/v1/crawls/:id` - Get crawl details
    - `GET /api/v1/crawls/:id/pages` - Get pages for a crawl
    - `GET /api/v1/crawls/:id/issues` - Get issues for a crawl
    - `POST /api/v1/issues/:id/status` - Update issue status
    - `POST /api/v1/exports` - Request export generation
  - JWT authentication middleware (Supabase token validation)
  - Service role key for background operations

- [ ] **Supabase Integration**
  - Go Supabase client library (`github.com/supabase-community/supabase-go`)
  - Database operations (insert crawls, pages, issues)
  - JWT validation for API requests
  - Service role key usage for admin operations

- [ ] **Environment Configuration**
  - `.env.example` with all required variables
  - Cloud Run environment variables:
    - `PUBLIC_SUPABASE_URL`
    - `SUPABASE_SERVICE_ROLE_KEY` (Secret Manager)
    - `PUBLIC_SUPABASE_ANON_KEY`
    - `GCS_BUCKET` (for exports, optional)
    - `OPENAI_API_KEY` (Secret Manager, optional)

- [ ] **Deployment Scripts**
  - `Makefile` targets:
    - `docker-build` - Build Docker image
    - `docker-push` - Push to Artifact Registry
    - `deploy-backend` - Deploy to Cloud Run
  - Terraform or `gcloud` CLI scripts

### Phase 2: Frontend Integration (Vercel)

- [ ] **Supabase Client Setup**
  - Install `@supabase/supabase-js` in `web/package.json`
  - Initialize Supabase client with environment variables
  - Auth context/provider for Svelte

- [ ] **Authentication UI**
  - Login/signup components
  - Protected routes
  - Session management

- [ ] **Data Fetching**
  - Replace local API calls (`/api/results`, `/api/summary`) with Supabase queries
  - Use RLS policies for data access
  - Real-time subscriptions for live updates

- [ ] **Environment Variables**
  - `PUBLIC_SUPABASE_URL` (Vite will expose as `VITE_PUBLIC_SUPABASE_URL` or use directly)
  - `PUBLIC_SUPABASE_ANON_KEY` (Vite will expose as `VITE_PUBLIC_SUPABASE_ANON_KEY` or use directly)
  - `VITE_CLOUD_RUN_API_URL` (for API endpoints)

- [ ] **Vercel Configuration**
  - `vercel.json` for routing/headers
  - Environment variables in Vercel dashboard
  - Build configuration

### Phase 3: CLI Integration

- [ ] **Cloud Upload Support**
  - `--cloud` flag for `crawl` command
  - POST crawl results to Cloud Run API
  - Authentication via Supabase JWT
  - Retry logic and error handling

- [ ] **Configuration**
  - `.barracuda/config.yaml` or `.env` for:
    - Cloud Run API URL
    - Supabase credentials (for auth)
    - Default project ID

### Phase 4: Secrets & Security

- [ ] **Google Secret Manager**
  - Store `SUPABASE_SERVICE_ROLE_KEY`
  - Store `OPENAI_API_KEY` (if used)
  - Configure Cloud Run to access secrets

- [ ] **Supabase Setup**
  - Create production project
  - Get service role key and anon key
  - Configure RLS policies (already in migration)
  - Set up auth providers (email, OAuth)

- [ ] **Cloud Run Security**
  - IAM roles and permissions
  - CORS configuration
  - Rate limiting (optional)

### Phase 5: Testing & Validation

- [ ] **Local Testing**
  - Test API server locally with Supabase local instance
  - Test frontend with local Supabase
  - End-to-end crawl → API → Supabase → Frontend flow

- [ ] **Deployment Testing**
  - Deploy to Cloud Run staging
  - Test API endpoints
  - Deploy frontend to Vercel preview
  - Integration testing

## Quick Start Commands (Once Implemented)

```bash
# Build and deploy backend
make docker-build
make docker-push
make deploy-backend

# Deploy frontend (via Vercel CLI or GitHub integration)
vercel --prod

# Upload crawl to cloud
barracuda crawl https://example.com --cloud --project-id <uuid>
```

## Reference Documentation

- `docs/CLOUD_RUN_SUPABASE.md` - Architecture blueprint
- `docs/SUPABASE_SCHEMA.md` - Database schema reference
- `supabase/migrations/` - Database migrations

