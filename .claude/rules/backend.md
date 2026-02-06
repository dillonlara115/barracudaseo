# Backend Guidelines (Go API)

## Handler Function Pattern

All handlers follow this signature and structure:

```go
func (s *Server) handleResourceName(w http.ResponseWriter, r *http.Request) {
	// 1. Extract user from context (set by auth middleware)
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		s.respondError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// 2. Parse request body (for POST/PUT/PATCH)
	var req CreateResourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
		return
	}

	// 3. Validate required fields
	if req.Name == "" {
		s.respondError(w, http.StatusBadRequest, "name is required")
		return
	}

	// 4. Authorization check
	hasAccess, err := s.verifyProjectAccess(userID, req.ProjectID)
	if err != nil {
		s.logger.Error("Failed to verify access", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to verify access")
		return
	}
	if !hasAccess {
		s.respondError(w, http.StatusForbidden, "You don't have access to this project")
		return
	}

	// 5. Business logic / database operations
	result, err := s.createResource(req)
	if err != nil {
		s.logger.Error("Failed to create resource", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to create resource")
		return
	}

	// 6. Success response
	s.respondJSON(w, http.StatusCreated, result)
}
```

## Naming Conventions

- Handlers: `handleXxx`, `handleXxxByID`, `handleCreateXxx`, `handleListXxx`
- Helper functions: `verifyXxx`, `fetchXxx`, `resolveXxx`
- Request types: `CreateXxxRequest`, `UpdateXxxRequest`

## HTTP Status Codes

| Status | Use Case |
|--------|----------|
| 200 OK | Successful GET, successful update |
| 201 Created | Successful POST creating new resource |
| 204 No Content | Successful DELETE |
| 400 Bad Request | Invalid JSON, validation errors |
| 401 Unauthorized | Missing/invalid authentication |
| 403 Forbidden | Access denied (valid auth, no permission) |
| 404 Not Found | Resource doesn't exist |
| 405 Method Not Allowed | Wrong HTTP method |
| 500 Internal Server Error | Database/service failures |

## Response Helpers

```go
// Success with data
s.respondJSON(w, http.StatusOK, map[string]interface{}{
	"projects": projects,
	"count":    len(projects),
})

// Error response
s.respondError(w, http.StatusBadRequest, "project_id is required")
```

## Request Type Definitions

Define in `types.go`:

```go
type CreateKeywordRequest struct {
	ProjectID      string `json:"project_id"`
	Keyword        string `json:"keyword"`
	LocationName   string `json:"location_name"`
	LocationCode   *int   `json:"location_code,omitempty"`
	Device         string `json:"device"`
	CheckFrequency string `json:"check_frequency,omitempty"`
}
```

## Supabase Query Patterns

### SELECT with filters
```go
query := s.supabase.From("crawls").Select("*", "", false)
query = query.Eq("project_id", projectID)
query = query.Order("started_at", &postgrest.OrderOpts{Ascending: false})
query = query.Limit(10, "")

data, _, err := query.Execute()
if err != nil {
	s.logger.Error("Failed to list crawls", zap.Error(err))
	s.respondError(w, http.StatusInternalServerError, "Failed to list crawls")
	return
}

var crawls []map[string]interface{}
if err := json.Unmarshal(data, &crawls); err != nil {
	s.logger.Error("Failed to parse data", zap.Error(err))
	return
}
```

### INSERT
```go
record := map[string]interface{}{
	"id":         uuid.New().String(),
	"project_id": req.ProjectID,
	"created_at": time.Now().UTC().Format(time.RFC3339),
}

_, _, err = s.serviceRole.From("resources").Insert(record, false, "", "", "").Execute()
```

### Batch INSERT
```go
batchSize := 1000
for i := 0; i < len(items); i += batchSize {
	end := i + batchSize
	if end > len(items) {
		end = len(items)
	}
	batch := items[i:end]
	_, _, err = s.serviceRole.From("items").Insert(batch, false, "", "minimal", "").Execute()
}
```

### UPDATE
```go
_, _, err = s.serviceRole.From("projects").
	Update(updateData, "", "").
	Eq("id", projectID).
	Execute()
```

## Authentication

JWT tokens validated via middleware. Access user ID from context:

```go
userID, ok := userIDFromContext(r.Context())
if !ok {
	s.respondError(w, http.StatusUnauthorized, "User not authenticated")
	return
}
```

## Authorization Pattern

```go
// Check project access (membership or ownership)
hasAccess, err := s.verifyProjectAccess(userID, projectID)
if err != nil {
	s.logger.Error("Failed to verify project access", zap.Error(err))
	s.respondError(w, http.StatusInternalServerError, "Failed to verify access")
	return
}
if !hasAccess {
	s.respondError(w, http.StatusForbidden, "You don't have access to this project")
	return
}
```

## Subscription Gating

```go
subscription := s.requireProSubscription(w, userID, "Keyword tracking")
if subscription == nil {
	return // Response already sent
}
```

## Logging

Use structured logging with zap:

```go
s.logger.Info("Created resource", zap.String("id", id), zap.String("user_id", userID))
s.logger.Error("Failed to query", zap.Error(err), zap.String("project_id", projectID))
s.logger.Debug("Processing batch", zap.Int("size", len(batch)))
s.logger.Warn("Token validation failed", zap.String("path", r.URL.Path))
```

## URL Path Parsing

For sub-resources, parse URL path:

```go
path := strings.TrimPrefix(r.URL.Path, "/projects/")
path = strings.Trim(path, "/")
parts := strings.Split(path, "/")

projectID := parts[0]
if len(parts) > 1 {
	switch parts[1] {
	case "crawls":
		s.handleProjectCrawls(w, r, projectID, userID)
	case "keywords":
		s.handleProjectKeywords(w, r, projectID, userID)
	}
}
```

## Webhook/Cron Handlers

```go
func (s *Server) handleCronJob(w http.ResponseWriter, r *http.Request) {
	// Verify cron secret
	secret := r.Header.Get("X-Cron-Secret")
	if secret == "" || secret != s.cronSecret {
		s.respondError(w, http.StatusUnauthorized, "Invalid or missing cron secret")
		return
	}

	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Process job...
}
```

## Supabase Clients

- `s.supabase` - Anon key client, respects RLS (for user-scoped queries)
- `s.serviceRole` - Service role client, bypasses RLS (for system operations)

Use service role for:
- Cross-user queries
- Background jobs
- Admin operations
