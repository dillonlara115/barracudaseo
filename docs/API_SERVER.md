# API Server Documentation

The Cloud Run API server provides REST endpoints for crawl ingestion, project management, and data retrieval.

## Structure

```
internal/api/
├── server.go      # Main server setup, middleware, routing
├── handlers.go    # HTTP request handlers for all endpoints
├── context.go    # Context helpers for user authentication
└── types.go       # Request/response types

cmd/
└── api.go         # CLI command to start the API server
```

## Running the Server

### Local Development

```bash
# Set environment variables
export PUBLIC_SUPABASE_URL=https://your-project.supabase.co
export SUPABASE_SERVICE_ROLE_KEY=your-service-role-key
export PUBLIC_SUPABASE_ANON_KEY=your-anon-key

# Run the server
go run . api --port 8080
```

Or using flags:

```bash
go run . api \
  --supabase-url https://your-project.supabase.co \
  --supabase-service-key your-service-role-key \
  --supabase-anon-key your-anon-key \
  --port 8080
```

### Docker/Cloud Run

The Dockerfile is configured to run the API server automatically. Set environment variables in Cloud Run:

- `PUBLIC_SUPABASE_URL`
- `SUPABASE_SERVICE_ROLE_KEY`
- `PUBLIC_SUPABASE_ANON_KEY`
- `PORT` (Cloud Run sets this automatically)

## API Endpoints

### Health Check

```
GET /health
```

No authentication required. Returns server status.

### Projects

#### Create Project
```
POST /api/v1/projects
Authorization: Bearer <supabase-jwt-token>
Content-Type: application/json

{
  "name": "My Website",
  "domain": "example.com",
  "settings": {}
}
```

#### List Projects
```
GET /api/v1/projects
Authorization: Bearer <supabase-jwt-token>
```

Returns all projects the authenticated user has access to (via RLS).

#### Get Project
```
GET /api/v1/projects/:id
Authorization: Bearer <supabase-jwt-token>
```

#### List Project Crawls
```
GET /api/v1/projects/:id/crawls
Authorization: Bearer <supabase-jwt-token>
```

### Crawls

#### Create Crawl (Ingest Crawl Results)
```
POST /api/v1/crawls
Authorization: Bearer <supabase-jwt-token>
Content-Type: application/json

{
  "project_id": "uuid-here",
  "pages": [
    {
      "url": "https://example.com",
      "status_code": 200,
      "response_time_ms": 150,
      "title": "Example",
      "meta_description": "...",
      ...
    }
  ],
  "source": "cli"
}
```

This endpoint:
1. Validates user has access to the project
2. Analyzes pages to detect SEO issues
3. Creates crawl record in database
4. Batch inserts pages
5. Batch inserts issues
6. Returns crawl summary

#### List Crawls
```
GET /api/v1/crawls?project_id=<optional-project-id>
Authorization: Bearer <supabase-jwt-token>
```

Returns crawls the user has access to (filtered by RLS policies).

## Authentication

All API endpoints (except `/health`) require a Supabase JWT token in the Authorization header:

```
Authorization: Bearer <supabase-jwt-token>
```

The server validates the token with Supabase Auth API and extracts the user ID from the token. The user ID is stored in the request context and used for:
- RLS policy enforcement (Supabase automatically filters based on user)
- Access verification (checking project membership)
- Audit logging

## Row-Level Security (RLS)

The API leverages Supabase RLS policies defined in the migration. When using the `supabase` client (anon key), queries are automatically filtered based on the authenticated user's access.

For admin operations (like bulk inserts during crawl ingestion), the `serviceRole` client (service role key) is used to bypass RLS.

## Error Responses

All errors return JSON in this format:

```json
{
  "error": "Error message here"
}
```

HTTP status codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `500` - Internal Server Error

## Next Steps

1. **Test locally** with Supabase local instance
2. **Add more endpoints**:
   - `GET /api/v1/crawls/:id` - Get crawl details
   - `GET /api/v1/crawls/:id/pages` - Get pages for a crawl
   - `GET /api/v1/crawls/:id/issues` - Get issues for a crawl
   - `PATCH /api/v1/issues/:id` - Update issue status
3. **Add CLI integration** - Update `cmd/crawl.go` to support `--cloud` flag
4. **Deploy to Cloud Run** - Use the provided Dockerfile and deployment scripts

