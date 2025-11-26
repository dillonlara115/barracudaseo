## 1. High-Level Overview

**Goal:** Add DataForSEO-powered rank tracking and SERP intelligence to Barracuda, integrated with existing Projects, Crawls, and GSC data.

**Core capabilities for v1:**

1. **Keyword tracking per project**

   * Track keyword → URL → location → device
   * Store daily rank snapshots from Google SERPs
2. **Simple rank tracker UI**

   * Table + basic charts by keyword and by page
3. **DataForSEO integration layer**

   * Internal Go client to create tasks and pull results
   * Background processing (async tasks + polling)
4. **Pricing + usage tracking**

   * Store per-keyword pulls
   * Estimate cost → enforce plan limits

Later we can expand to Local Pack & GeoGrid, but v1 focuses on **standard organic rank tracking** and simple **local SERP**.

---

## 2. DataForSEO APIs to Use (v1 & v2)

### v1 – Core Rank Tracker

Use DataForSEO **SERP – Google – Organic** endpoints:

* **Create rank task (per keyword):**
  `POST /v3/serp/google/organic/task_post`
  (or `task_post_bulk` if batching)
* **Get task result:**
  `GET /v3/serp/google/organic/task_get/{id}`

Use parameters:

* `keyword` – user-defined keyword
* `location_name` or `location_code` – user’s target location (e.g., "United States", "Denver, Colorado")
* `language_name` – e.g. `"English"`
* `device` – `"desktop"` or `"mobile"`
* `search_engine` – `"google.com"` or localized variant

**What we care about in results:**

* `check_url` or your **target URL**
* `result_position` (absolute & organic)
* `result_type` (organic, featured snippet, local pack result etc.)
* `url` of the ranking page
* `serp_features` (if available)

---

### v2 – Local Pack & GeoGrid (Future Phases)

**Local Pack ranking**

* Endpoint: `POST /v3/serp/google/local/task_post` + `task_get`
* Use to get 3-pack ranking for GBP / local business.

**GeoGrid (LocalFalcon-style)**

* Use coordinate-based SERP / Maps endpoints per grid point (e.g., `latitude`, `longitude` fields).
* You’ll orchestrate a grid of tasks (e.g. 7x7 or 9x9) and visualize as heatmap.

We’ll design the schema to be extensible enough to hold these later (e.g. separate `rank_type` or `serp_type` fields).

---

## 3. Supabase Schema for Rank Tracking

Add a separate SQL migration, e.g. `supabase/migrations/2025xxxx_dataforseo_rank_tracking.sql`

### 3.1. Tables

#### `keywords`

```sql
create table public.keywords (
  id uuid primary key default gen_random_uuid(),
  project_id uuid not null references public.projects(id) on delete cascade,
  keyword text not null,
  target_url text, -- optional: canonical URL we want to rank
  location_name text not null, -- e.g. "United States", "Denver, Colorado"
  location_code integer,       -- optional: DataForSEO location code
  language_name text not null default 'English',
  device text not null default 'desktop', -- desktop | mobile
  search_engine text not null default 'google.com',
  tags text[] default '{}',
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

-- Optional uniqueness: no duplicate keyword+location+device in same project
create unique index keywords_project_keyword_loc_device_idx
  on public.keywords (project_id, keyword, coalesce(location_name, ''), device);
```

#### `keyword_rank_snapshots`

Each “check” for a keyword produces one snapshot record.

```sql
create table public.keyword_rank_snapshots (
  id uuid primary key default gen_random_uuid(),
  keyword_id uuid not null references public.keywords(id) on delete cascade,
  checked_at timestamptz not null default now(),
  dataforseo_task_id text,      -- DataForSEO task ID
  position_absolute integer,    -- overall position in SERP
  position_organic integer,     -- organic-only position
  serp_url text,                -- ranking URL
  serp_title text,
  serp_snippet text,
  serp_features text[] default '{}', -- e.g. ['featured_snippet','sitelinks']
  search_volume integer,        -- optional: from DataForSEO stats
  rank_type text not null default 'organic', -- organic | local_pack | maps
  raw jsonb,                    -- full API response for debugging
  created_at timestamptz not null default now()
);

create index keyword_rank_snapshots_keyword_id_idx
  on public.keyword_rank_snapshots (keyword_id);

create index keyword_rank_snapshots_checked_at_idx
  on public.keyword_rank_snapshots (checked_at);
```

#### `keyword_tasks` (optional but helpful for async polling)

```sql
create table public.keyword_tasks (
  id uuid primary key default gen_random_uuid(),
  keyword_id uuid not null references public.keywords(id) on delete cascade,
  dataforseo_task_id text not null,
  status text not null default 'pending', -- pending | processing | completed | failed
  run_at timestamptz not null default now(), -- when we created the task
  completed_at timestamptz,
  error text,
  raw_request jsonb,
  raw_response jsonb,
  created_at timestamptz not null default now()
);

create index keyword_tasks_keyword_id_idx
  on public.keyword_tasks (keyword_id);

create index keyword_tasks_status_idx
  on public.keyword_tasks (status);
```

### 3.2. RLS Policies

Similar to other tables:

* Only owners/members of `projects` can see/edit keywords and snapshots.
  Tie policies back to `projects` and `project_members` like your existing schema.

---

## 4. Go Backend Integration & API Routes

### 4.1. DataForSEO Client (Go)

Create internal package: `internal/dataforseo/client.go`

```go
package dataforseo

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type Client struct {
    httpClient *http.Client
    baseURL    string
    login      string // DataForSEO login (email)
    password   string // DataForSEO password / API key
}

func NewClient(baseURL, login, password string) *Client {
    return &Client{
        httpClient: &http.Client{Timeout: 30 * time.Second},
        baseURL:    baseURL,
        login:      login,
        password:   password,
    }
}

func (c *Client) do(ctx context.Context, method, path string, body any, v any) error {
    var buf bytes.Buffer
    if body != nil {
        if err := json.NewEncoder(&buf).Encode(body); err != nil {
            return fmt.Errorf("encode body: %w", err)
        }
    }

    req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, &buf)
    if err != nil {
        return fmt.Errorf("new request: %w", err)
    }

    req.SetBasicAuth(c.login, c.password)
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode >= 300 {
        // You can decode error body to a custom error struct here
        return fmt.Errorf("dataforseo: status %d", resp.StatusCode)
    }

    if v != nil {
        if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
            return fmt.Errorf("decode response: %w", err)
        }
    }

    return nil
}
```

Define types for **task_post** and **task_get**:

```go
// internal/dataforseo/types.go
package dataforseo

type OrganicTaskPost struct {
    LanguageName string `json:"language_name"`
    LocationName string `json:"location_name"`
    Keyword      string `json:"keyword"`
    Device       string `json:"device,omitempty"`        // "desktop" or "mobile"
    SearchEngine string `json:"search_engine_name,omitempty"` // "google.com" etc.
}

type OrganicTaskPostRequest map[string]OrganicTaskPost

type OrganicTaskPostResponse struct {
    Tasks []struct {
        ID    string `json:"id"`
        // ... other fields (status_code, etc.)
    } `json:"tasks"`
}

type OrganicTaskGetResponse struct {
    Tasks []struct {
        ID     string `json:"id"`
        Result []struct {
            Items []struct {
                RankAbsolute int      `json:"rank_absolute"`
                RankOrganic  int      `json:"rank_group"`
                Url          string   `json:"url"`
                Title        string   `json:"title"`
                Snippet      string   `json:"description"`
                Features     []string `json:"serp_features"`
            } `json:"items"`
        } `json:"result"`
    } `json:"tasks"`
}
```

Helper methods:

```go
func (c *Client) CreateOrganicTask(ctx context.Context, task OrganicTaskPost) (*OrganicTaskPostResponse, error) {
    body := OrganicTaskPostRequest{
        "0": task,
    }
    var resp OrganicTaskPostResponse
    if err := c.do(ctx, http.MethodPost, "/v3/serp/google/organic/task_post", body, &resp); err != nil {
        return nil, err
    }
    return &resp, nil
}

func (c *Client) GetOrganicTask(ctx context.Context, taskID string) (*OrganicTaskGetResponse, error) {
    var resp OrganicTaskGetResponse
    path := fmt.Sprintf("/v3/serp/google/organic/task_get/%s", taskID)
    if err := c.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
        return nil, err
    }
    return &resp, nil
}
```

Wire this client in your `internal/api/server.go`:

* Add `DataForSEOClient *dataforseo.Client` to server struct
* Initialize from env vars:

```go
DATAFORSEO_BASE_URL
DATAFORSEO_LOGIN
DATAFORSEO_PASSWORD
```

---

### 4.2. API Routes (Go → Supabase → Frontend)

All under `/api/v1/keywords`.

#### 4.2.1. `POST /api/v1/keywords`

Create keyword for a project.

**Request body:**

```json
{
  "project_id": "uuid",
  "keyword": "best singing bowls",
  "target_url": "https://example.com/best-singing-bowls",
  "location_name": "United States",
  "language_name": "English",
  "device": "desktop",
  "search_engine": "google.com",
  "tags": ["priority", "blog"]
}
```

**Handler:**

* Auth: require logged-in user
* Check user has access to `project_id` (via Supabase or your helper)
* Insert row into `public.keywords`
* Return created keyword.

#### 4.2.2. `GET /api/v1/keywords?project_id=...`

List keywords for a project, with optional filters.

#### 4.2.3. `POST /api/v1/keywords/:id/check`

Trigger a rank check for a keyword.

Handler steps:

1. Load keyword from DB (validate ownership).
2. Create DataForSEO organic task:

   * `keyword`
   * `location_name`
   * `language_name`
   * `device`
3. Insert into `keyword_tasks` with `status='pending'`, store `dataforseo_task_id`.
4. Option A (simple v1): Immediately call `GetOrganicTask` and store snapshot (DataForSEO often processes quickly).
5. Option B (scalable): Return task ID; a background cron/worker polls tasks and writes snapshots later.
6. Insert into `keyword_rank_snapshots`.

Response includes snapshot data.

#### 4.2.4. `GET /api/v1/keywords/:id/snapshots?limit=30`

Return historical snapshots for a keyword.

#### 4.2.5. `GET /api/v1/projects/:id/keyword-metrics`

Aggregate per-project metrics:

* latest position per keyword
* best position achieved
* trend (up/down/same)
* average position, etc.

---

## 5. SvelteKit UI Prototype (Rank Tracker)

Create route: `web/src/routes/projects/[projectId]/rank-tracker/+page.svelte`

### 5.1. Page structure

* Filters:

  * Keyword search
  * Tag filter
  * Device filter

* Table of keywords:

  * Keyword
  * Target URL
  * Location
  * Latest position
  * Best position
  * Trend (arrow up/down)
  * Last checked date
  * “Check now” button

* Detail panel / modal:

  * Line chart of position over time
  * Table of snapshots
  * Option to link to corresponding page issues / crawl data

### 5.2. Load function

`+page.ts` or `+page.server.ts`:

```ts
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch, params }) => {
  const res = await fetch(`/api/v1/keywords?project_id=${params.projectId}`);
  const keywords = await res.json();

  return {
    keywords
  };
};
```

### 5.3. Basic Svelte markup (simplified)

```svelte
<script lang="ts">
  export let data;
  let keywords = data.keywords;

  let selectedKeyword = null;
  let isChecking = new Set<string>();

  const checkKeyword = async (id: string) => {
    isChecking.add(id);
    const res = await fetch(`/api/v1/keywords/${id}/check`, { method: 'POST' });
    const snapshot = await res.json();
    isChecking.delete(id);

    // Update local list: assign latest_position, last_checked, etc.
    keywords = keywords.map((k) =>
      k.id === id ? { ...k, latest_snapshot: snapshot } : k
    );
  };
</script>

<section class="space-y-6">
  <div class="flex items-center justify-between">
    <h1 class="text-2xl font-bold">Rank Tracker</h1>
    <a href="/projects/{data.projectId}/rank-tracker/new" class="btn btn-primary">
      Add Keyword
    </a>
  </div>

  <div class="overflow-x-auto">
    <table class="table table-zebra w-full">
      <thead>
        <tr>
          <th>Keyword</th>
          <th>Target URL</th>
          <th>Location</th>
          <th>Device</th>
          <th>Latest Pos</th>
          <th>Best Pos</th>
          <th>Last Checked</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        {#each keywords as k}
          <tr on:click={() => (selectedKeyword = k)} class="cursor-pointer">
            <td>{k.keyword}</td>
            <td class="truncate max-w-xs">{k.target_url}</td>
            <td>{k.location_name}</td>
            <td class="capitalize">{k.device}</td>
            <td>{k.latest_snapshot?.position_organic ?? '—'}</td>
            <td>{k.best_position ?? '—'}</td>
            <td>{k.latest_snapshot?.checked_at
              ? new Date(k.latest_snapshot.checked_at).toLocaleDateString()
              : '—'}</td>
            <td>
              <button
                class="btn btn-xs btn-outline"
                on:click|stopPropagation={() => checkKeyword(k.id)}
                disabled={isChecking.has(k.id)}
              >
                {isChecking.has(k.id) ? 'Checking…' : 'Check now'}
              </button>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>

  {#if selectedKeyword}
    <!-- Modal: show snapshot chart + details -->
  {/if}
</section>
```

Later: add chart using something like Recharts or a simple SVG line chart.

---

## 6. Pricing Model (Cost per Keyword vs Markup)

Assume (example numbers, adjust with real DataForSEO pricing):

* DataForSEO charges **~$0.0006 – $0.001 per SERP keyword check** (varies by plan/volume).
* If a user tracks 100 keywords daily:

  * 100 * 30 * $0.001 = **$3 / month** wholesale.
* You can safely charge **10x+** on top because:

  * You’re not just reselling API calls, you’re adding UI + analysis + AI.

### Suggested tiers (example)

**Free**

* Up to 10 keywords
* Weekly checks only
* No local / GeoGrid

**Pro ($39–$49 / mo)**

* 250–500 keywords
* Daily checks
* Basic SERP view
* Integrated with crawl & GSC
* AI insights for ranking pages

**Team / Agency ($99–$149 / mo)**

* 1,000–2,000 keywords
* Daily checks
* Local pack tracking
* Exportable reports
* Team accounts
* White-label PDFs (later)

At scale, your DataForSEO costs might be ~$20–$60/mo per agency-tier customer, leaving healthy margin.

Add **fair-usage caps** and maybe a “soft limit” with overage.

---

## 7. Competitive Messaging Module

Use this content on:

* Marketing site
* In-app onboarding
* Sales pages / pricing

### Positioning line

> **“Barracuda shows you what’s broken, how to fix it, and whether your fixes are actually moving rankings — all in one place.”**

### Why Barracuda beats ScreamingFrog

* ScreamingFrog:

  * Great crawler
  * No built-in rank tracking
  * No GSC integration in a unified view
  * No AI action plans

* **Barracuda:**

  * Modern crawler + dashboard
  * Integrated GSC + rank tracking (DataForSEO)
  * AI that explains what to fix next and why

### Why Barracuda beats SEMrush & Ahrefs (for this workflow)

* SEMrush / Ahrefs:

  * Rank tracking lives separately from technical audit
  * Limited or generic AI recommendations
  * Heavier learning curve, overkill for many SMBs / agencies

* **Barracuda:**

  * Crawl issues + rankings + GSC data tied to each URL
  * “Impact-first” view: prioritize issues on pages that both rank and convert
  * Simple, focused UX just for SEO + UX/CRO insights

### Why Barracuda beats LocalFalcon

* LocalFalcon:

  * Only local grid rankings
  * No technical SEO insight
  * No site crawl or AI diagnosis

* **Barracuda:**

  * Local rankings (future GeoGrid) + **full site health**
  * Connects local visibility to technical + on-page issues
  * One tool instead of 2–3 separate subscriptions

### Hero copy ideas

* “From crawl to rankings: one tool that shows you what to fix and proves it worked.”
* “Technical SEO + rank tracking + AI insights — finally in one place.”
* “Stop guessing which fixes matter. Barracuda connects your issues to your rankings.”
