# barracuda

A fast, lightweight SEO website crawler CLI tool inspired by Screaming Frog.

## Features

- **Recursive Crawling**: Crawl websites with configurable depth, page limits, and sitemap seeding
- **SEO Intelligence**: Capture titles, meta data, headings, links, images, response metrics, and automatic issue detection
- **Cloud-Ready Workflow**: Push crawl results to the hosted API (Supabase + Cloud Run) for multi-user projects
- **Projects Workspace**: Authenticate with Supabase, manage projects, and view historical crawls inside the dashboard
- **Insightful Dashboard**: Priority scoring, page-level modals, recommendations, and quick filters surface the most critical fixes
- **Rank Tracking**: Track keyword rankings over time with DataForSEO integration, scheduled checks, and historical snapshots
- **Flexible Exports**: One-click CSV/JSON exports for filtered issues, link graphs, and crawl data
- **Fast & Concurrent**: Worker pool, robots.txt controls, and domain filtering keep crawls respectful and performant

## Components at a Glance

- **CLI (`barracuda crawl`)**: Run local crawls, export results, and optionally upload to the cloud workspace.
- **Embedded Dashboard (`barracuda serve`)**: Bundle crawl results into a Svelte UI with advanced filtering, grouping, and recommendations.
- **Managed API (`barracuda api`)**: Supabase-backed REST service designed for Cloud Run that stores projects, crawls, pages, and issues for teams.
- **Hosted Frontend (Vercel)**: Production dashboard at https://app.barracudaseo.com connected to Supabase auth and the Cloud Run API.

## Installation

### Prerequisites

- Go 1.21 or newer
- Node 18+ (for building the Svelte frontend)
- Supabase project (URL, anon key, and optionally service role key) if you plan to use the hosted workspace or API

> `.env` files: the CLI loads `.env` first and then `.env.local` (if present) so you can keep shared defaults in `.env` and developer overrides in `.env.local`.

### From Source

**Important:** The frontend must be built before compiling the Go binary, as it is embedded into the executable.

```bash
git clone https://github.com/dillonlara115/barracuda.git
cd barracuda
make frontend-build  # Build frontend first
make build           # Build binary (includes embedded frontend)
sudo mv bin/barracuda /usr/local/bin/
```

Or manually:

```bash
git clone https://github.com/dillonlara115/barracuda.git
cd barracuda
cd web && npm install && npm run build && cd ..
go build -o barracuda .
sudo mv barracuda /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/dillonlara115/barracuda@latest
```

### Frontend Setup (Required for Building)

The frontend is embedded into the binary, so it must be built before compiling:

```bash
cd web
npm install
npm run build
```

Or use the Makefile:

```bash
make frontend-build
```

**Note:** When installed via `go install` or built from source, the frontend is automatically included in the binary and works from any directory.

### Environment Variables

Create a `.env.local` at the repo root (and `web/.env.local` for the frontend) with the Supabase and Cloud Run configuration you need:

```bash
PUBLIC_SUPABASE_URL=https://your-project.supabase.co
PUBLIC_SUPABASE_ANON_KEY=public-anon-key
SUPABASE_SERVICE_ROLE_KEY=service-role-key # only required when running the API server
VITE_CLOUD_RUN_API_URL=https://barracuda-api-your-env.a.run.app
```

See `docs/SUPABASE_SCHEMA.md`, `docs/API_SERVER.md`, and `docs/CLOUD_RUN_SUPABASE.md` for the complete data model and deployment flow.

### Run the API Server Locally

The `api` command starts the Supabase-backed REST service used in Cloud Run.

```bash
# With env vars already exported or present in .env/.env.local
go run . api --port 8080

# Or pass flags explicitly
go run . api \
  --supabase-url https://your-project.supabase.co \
  --supabase-anon-key public-anon-key \
  --supabase-service-key service-role-key \
  --port 8080
```

Docker helpers (`make docker-build`, `make deploy-backend`) are available for packaging and deploying to Cloud Run.

## Usage

### Basic Usage

```bash
# Crawl a website
barracuda crawl https://example.com

# Crawl with custom depth and export format
barracuda crawl https://example.com --max-depth 2 --format json

# Export to specific file
barracuda crawl https://example.com --export results.csv
```

### Advanced Options

```bash
# Full example with all options
barracuda crawl https://example.com \
  --max-depth 3 \
  --max-pages 500 \
  --workers 20 \
  --delay 100ms \
  --timeout 60s \
  --format json \
  --export crawl-results.json \
  --graph-export link-graph.json \
  --parse-sitemap \
  --respect-robots
```

### Web Dashboard

After crawling, view your results in a beautiful web interface:

```bash
# First, crawl with JSON export
barracuda crawl https://example.com --format json --export results.json --graph-export graph.json

# Then serve the results
barracuda serve --results results.json --graph graph.json
```

Or use the Makefile shortcut:

```bash
make serve
```

The web dashboard includes:
- **Dashboard Overview**: KPI cards with deep links into issues, slow pages, and critical errors
- **Issues Panel**: Search, severity filters, grouping toggles, priority scoring, and one-click CSV/JSON exports
- **Results Table**: Page-level issue counts, detail modal with metadata, and quick navigation to filtered issues
- **Recommendations Tab**: Curated fixes with copyable snippets, impact indicators, and links to best practices
- **Projects Workspace**: Authenticated Supabase session picker that lets you switch projects or create new ones
- **Link Graph**: Visualization of internal and external link structures with export support

Access the dashboard at `http://localhost:8080` (default port).

For the hosted dashboard (https://app.barracudaseo.com) configure Supabase auth + API URLs as described in `docs/VERCEL_DEPLOYMENT.md` and `docs/VERCEL_URL.md`.

## Command-Line Flags

### Required Flags

- `--url, -u`: Starting URL to crawl (required)

### Crawl Options

- `--max-depth, -d`: Maximum crawl depth (default: 3)
- `--max-pages, -p`: Maximum number of pages to crawl (default: 1000)
- `--workers, -w`: Number of concurrent workers (default: 10)
- `--delay`: Delay between requests (e.g., 100ms) (default: 0ms)
- `--timeout`: HTTP request timeout (default: 30s)
- `--user-agent`: User agent string (default: barracuda/1.0.0)
- `--respect-robots`: Respect robots.txt rules (default: true)
- `--parse-sitemap`: Parse sitemap.xml for seed URLs (default: false)
- `--domain-filter`: Domain filter: 'same' or 'all' (default: same)

### Export Options

- `--format, -f`: Export format: 'csv' or 'json' (default: csv)
- `--export, -e`: Export file path (default: results.csv/json)
- `--graph-export`: Export link graph to JSON file (optional)

### Serve Command (Web Dashboard)

- `serve`: Start web server to view crawl results
  - `--port`: Port to run the server on (default: 8080)
  - `--results`: Path to JSON results file (default: results.json)
  - `--graph`: Path to link graph JSON file (optional)
  - `--summary`: Path to summary JSON file (optional, auto-generated if not provided)

### API Command (Cloud Workspace)

- `api`: Start the Supabase-backed REST server
  - `--port`: Port to run the API server on (default/Cloud Run: 8080 or pulled from `PORT`)
  - `--supabase-url`: Supabase project URL (`PUBLIC_SUPABASE_URL`)
  - `--supabase-service-key`: Supabase service role key (`SUPABASE_SERVICE_ROLE_KEY`)
  - `--supabase-anon-key`: Supabase anon key (`PUBLIC_SUPABASE_ANON_KEY`)

### Global Flags

- `--debug`: Enable debug logging
- `--version`: Show version information

## Examples

### Example 1: Basic Crawl

```bash
barracuda crawl https://example.com
```

This will:
- Crawl up to 3 levels deep
- Export results to `results.csv`
- Respect robots.txt by default

### Example 2: JSON Export with Link Graph

```bash
barracuda crawl https://example.com \
  --format json \
  --export results.json \
  --graph-export graph.json
```

### Example 3: Fast Crawl (No Robots, Higher Concurrency)

```bash
barracuda crawl https://example.com \
  --workers 50 \
  --max-pages 5000 \
  --respect-robots=false
```

### Example 4: Sitemap-Based Crawl

```bash
barracuda crawl https://example.com \
  --parse-sitemap \
  --max-depth 1
```

### Example 5: View Results in Web Dashboard

```bash
# Step 1: Crawl and export to JSON
barracuda crawl https://example.com \
  --format json \
  --export results.json \
  --graph-export graph.json

# Step 2: Build frontend (first time only)
cd web && npm install && npm run build

# Step 3: Serve results
barracuda serve --results results.json --graph graph.json

# Open http://localhost:8080 in your browser
```

### Example 6: Run the Cloud API Locally

```bash
export PUBLIC_SUPABASE_URL=https://your-project.supabase.co
export PUBLIC_SUPABASE_ANON_KEY=public-anon-key
export SUPABASE_SERVICE_ROLE_KEY=service-role-key

barracuda api --port 8080
```

The API now serves authenticated REST endpoints for projects, crawls, pages, and issues at `http://localhost:8080/api/v1/...`. Review `docs/API_SERVER.md` for a complete endpoint list.

## Output Format

### CSV Export

The CSV export includes the following columns:
- URL
- Status Code
- Response Time (ms)
- Title
- Meta Description
- Canonical
- H1-H6 (pipe-separated values)
- Internal Links (pipe-separated)
- External Links (pipe-separated)
- Redirect Chain (arrow-separated)
- Error
- Crawled At

### JSON Export

The JSON export includes an array of page results with all SEO data fields.

### Link Graph Export

The link graph is exported as a JSON object mapping source URLs to arrays of target URLs:
```json
{
  "https://example.com/page1": [
    "https://example.com/page2",
    "https://example.com/page3"
  ],
  "https://example.com/page2": [
    "https://example.com/page4"
  ]
}
```

## Performance

- Typical crawl speed: 100-500 pages/minute (depends on server response times)
- Memory usage: ~50-100 MB for 1000 pages (varies by page size)
- Concurrent workers: Adjust `--workers` based on your system and target server capacity

## SEO Analysis

The crawler automatically detects SEO issues including:

- Missing or duplicate H1 tags
- Missing meta descriptions
- Missing or poor titles
- Large images (>100KB)
- Missing image alt text
- Slow response times
- Redirect chains
- Broken links

Issues are displayed in the terminal summary and can be viewed in detail in the web dashboard.

## Limitations

- No database storage (all data in-memory)
- No JavaScript rendering (static HTML only)
- Binary size: ~15-20 MB (includes embedded frontend)

## Cloud Deployment & Integrations

- **Cloud Run + Supabase + Vercel**: Follow `docs/CLOUD_RUN_SUPABASE.md` for the end-to-end architecture and `docs/CLOUD_RUN_DEPLOYMENT.md` / `docs/DEPLOYMENT_CHECKLIST.md` for deployment automation.
- **Supabase Schema & RLS**: Detailed tables, policies, and workflows live in `docs/SUPABASE_SCHEMA.md` with redirect configuration in `docs/SUPABASE_REDIRECT_SETUP.md`.
- **Frontend Hosting**: `docs/VERCEL_DEPLOYMENT.md` and `docs/VERCEL_URL.md` cover production hosting, environment variables, and Supabase auth settings.
- **Search Console & Integrations**: See `docs/GSC_SETUP_CHECKLIST.md`, `docs/GSC_CREDENTIALS.md`, and `docs/GSC_INTEGRATION.md` for enabling Google Search Console data pulls.
- **Agents & API**: `docs/AGENTS.md` provides context for contributors/AI agents, while `docs/API_SERVER.md` documents the REST endpoints exposed by `barracuda api`.

## Development

### Building from Source

**Important:** The frontend must be built before compiling the Go binary.

```bash
# Build frontend, then build CLI (recommended)
make build

# Or build separately:
make frontend-build
go build -o bin/barracuda .

# Run tests
make test

# Start the Supabase-backed API server
go run . api --port 8080

# Build and deploy the Cloud Run image (requires GCP + Supabase env vars)
make deploy-backend

# Install locally (builds frontend automatically)
make install

# Build frontend
make frontend-build

# Run frontend in dev mode (with hot reload)
make frontend-dev

# Serve crawl results
make serve
```

### Project Structure

```
barracuda/
├── cmd/                     # CLI entrypoints
│   ├── api.go              # Cloud Run / Supabase API command
│   ├── crawl.go            # Crawl command
│   ├── serve.go            # Serve command (embedded dashboard)
│   └── browser.go          # Browser helpers
├── internal/
│   ├── api/                # REST server (handlers, router, types)
│   ├── analyzer/           # SEO analysis and issue detection
│   ├── crawler/            # Crawl engine
│   ├── exporter/           # CSV/JSON export logic
│   ├── graph/              # Link graph utilities
│   └── utils/              # Shared helpers (config, logging, prompts)
├── pkg/
│   └── models/             # Shared data models
├── web/                    # Svelte dashboard
│   ├── src/
│   │   ├── components/     # Dashboard, IssuesPanel, ProjectsView, etc.
│   │   └── App.svelte
│   └── package.json
├── docs/                   # Deployment, Supabase, and roadmap docs
│   ├── API_SERVER.md
│   ├── CLOUD_RUN_SUPABASE.md
│   ├── DASHBOARD_IMPROVEMENTS.md
│   └── SUPABASE_SCHEMA.md
└── main.go                 # Entry point
```

## License

MIT License - see LICENSE file for details

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## Acknowledgments

Inspired by Screaming Frog SEO Spider.
