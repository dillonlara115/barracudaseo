# DataForSEO API Clarification

## Which API Are We Using?

We are using **DataForSEO SERP API**, specifically:
- **Endpoint Path**: `/v3/serp/google/organic/`
- **Purpose**: Get Google search engine results pages (SERPs) to track keyword rankings
- **Documentation**: https://docs.dataforseo.com/v3/serp/google/organic/

## NOT Using Keywords Data API

The **Keywords Data API** (`/v3/keywords_data/`) is a different API that provides:
- Keyword search volume
- Related keywords
- Keywords for site
- Ad traffic by keywords

This is NOT what we need for rank tracking. We need the SERP API to see actual search results and positions.

## Our Implementation

### Endpoints We Use:

1. **Create Task**: `POST /v3/serp/google/organic/task_post`
   - Creates a task to fetch SERP data for a keyword
   - Returns task ID immediately

2. **Check Ready Tasks**: `GET /v3/serp/google/organic/tasks_ready`
   - Returns list of all completed tasks across entire account
   - We filter to only our tasks

3. **Get Task Result**: `GET /v3/serp/google/organic/task_get/{id}`
   - Retrieves the actual SERP results for a task
   - Returns ranking positions, URLs, titles, snippets

### Current Issue

Tasks are being created successfully, but:
- `tasks_ready` returns different task IDs than what we stored
- Individual `task_get` calls return 40400 "Not Found"
- This suggests tasks may expire quickly or there's a delay in processing

### Next Steps

1. Verify task ID format matches DataForSEO expectations
2. Check if tasks expire before we can retrieve them
3. Consider using `postback_url` instead of polling
4. Verify we're using the correct API endpoint format




