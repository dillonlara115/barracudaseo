# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Barracuda SEO is a website crawler CLI tool with three main modes:
- **CLI** (`barracuda crawl`): Crawls websites, detects SEO issues, exports to CSV/JSON
- **Embedded Dashboard** (`barracuda serve`): Local Svelte UI for viewing crawl results
- **Cloud API** (`barracuda api`): REST server backed by Supabase for multi-user workspaces

Production hosted at https://app.barracudaseo.com (Vercel frontend + Cloud Run API + Supabase DB).

## Build Commands

```bash
# Build (frontend + binary)
make build

# Run tests
make test

# Frontend dev server (hot reload)
make frontend-dev

# Run CLI commands directly
go run . crawl https://example.com
go run . serve --results results.json
go run . api --port 8080

# Lint and format
make lint              # Go linting (requires golangci-lint)
cd web && npm run lint      # Frontend linting
cd web && npm run format    # Format with Prettier

# Docker / Cloud Run deployment
make docker-build
make deploy-backend   # full deployment with env vars
make deploy-image     # image-only deployment (preserves env vars)
```

**Important:** Frontend must be built before Go binary (`make frontend-build` or `make build`). The Svelte app is embedded into the binary via `go:embed`.

## Architecture

### Backend (Go 1.25)
- **CLI Framework**: Cobra - commands defined in `cmd/*.go`
- **HTTP Router**: Standard library net/http
- **Database**: Supabase (PostgreSQL with Row-Level Security)
- **HTML Parsing**: goquery (`github.com/PuerkitoBio/goquery`)
- **Logging**: zap (`go.uber.org/zap`) via `utils.Info()`, `utils.Debug()`, `utils.Error()`

### Frontend (Svelte)
- **Framework**: Svelte 4 + Vite 5
- **Styling**: Tailwind CSS + DaisyUI (custom "barracuda" theme)
- **Location**: `web/src/` - built output embedded from `web/dist/`

### Key Directories
- `cmd/` - CLI commands (crawl, serve, api, auth, cloud)
- `internal/api/` - REST handlers for Cloud Run API
- `internal/crawler/` - Crawl engine (manager, fetcher, parser, robots)
- `internal/analyzer/` - SEO issue detection
- `pkg/models/` - Shared data models (PageResult, Image)
- `web/src/components/` - Svelte UI components
- `supabase/migrations/` - Database migrations

### Concurrency Patterns
- Worker pool in `crawler/manager.go` using channels
- `sync.Map` for visited URLs
- `sync/atomic` for counters
- Context cancellation for shutdown

## Code Style

### Formatting
- **Go**: `go fmt` / `gofmt`
- **Svelte/JS/CSS**: Prettier with tabs, single quotes, no trailing commas
- Run `npm run format` in `web/` before committing frontend changes

### Go Imports Order
```go
import (
    // Standard library
    "fmt"

    // Third-party
    "github.com/spf13/cobra"

    // Internal
    "github.com/dillonlara115/barracudaseo/internal/utils"
)
```

### Error Handling
```go
if err != nil {
    return fmt.Errorf("context: %w", err)
}
```

### URL Handling
Always use `utils.NormalizeURL()` before storing/comparing URLs to prevent duplicates.

## Project-Specific Rules

See `.claude/rules/` for detailed guidelines:
- `frontend.md` - Svelte component patterns, DaisyUI usage, Tailwind conventions
- `backend.md` - Go API handler patterns, authentication, database queries
- `database.md` - Supabase migrations, RLS policies, indexing conventions

## API Endpoints

Cloud Run API at `/api/v1/*`:
- Projects: `POST/GET /projects`, `GET /projects/:id/crawls`
- Crawls: `POST /crawls` (ingest from CLI), `GET /crawls/:id`
- Pages/Issues: `GET /pages`, `GET /issues`
- Keywords: `POST/GET /keywords`, `POST /keywords/:id/track`
- Teams: `POST/GET /teams`, `GET /team-members`
- Integrations: `/gsc/*`, `/ga4/*` (Google Search Console, Analytics)

## Environment Variables

Required for cloud features (in `.env` or `.env.local`):
```
PUBLIC_SUPABASE_URL=https://your-project.supabase.co
PUBLIC_SUPABASE_ANON_KEY=anon-key
SUPABASE_SERVICE_ROLE_KEY=service-role-key  # API server only
VITE_CLOUD_RUN_API_URL=https://barracuda-api.a.run.app
```

CLI does not auto-load `.env`. Set `BARRACUDA_LOAD_ENV=1` or `BARRACUDA_ENV_FILE=/path/.env`.

## Pre-commit Hooks

Set up git hooks for automatic formatting and type checking:

```bash
./scripts/setup-hooks.sh
```

The pre-commit hook will:
- Format Go files with `gofmt`
- Format frontend files with Prettier
- Run Svelte type checking with `svelte-check`

## Documentation

Extended docs in `docs/`:
- `API_SERVER.md` - API endpoint reference
- `CLOUD_RUN_SUPABASE.md` - Architecture overview
- `SUPABASE_SCHEMA.md` - Database schema and RLS policies
- `AGENTS.md` - Detailed patterns for AI contributors
