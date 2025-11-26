# DataForSEO Integration - Debugging Notes

## Current Status

The DataForSEO rank tracking integration is mostly implemented, but there are issues with task retrieval that need to be resolved.

## Implemented Features

✅ **Phase 1 - Core Rank Tracker:**
- Database schema (keywords, keyword_rank_snapshots, keyword_tasks tables)
- DataForSEO client implementation
- API handlers for keyword CRUD operations
- Background task poller
- Frontend Rank Tracker UI
- Keyword form and detail modals

✅ **Phase 2 - Scheduled Checks & Usage Tracking:**
- Scheduled automatic checks (daily/weekly)
- Usage tracking and subscription limits
- Crawl data integration (linking snapshots to pages)
- Impact-first view API endpoint
- Frontend usage dashboard

## Current Issues

### Issue: Task Retrieval Returns 40400 "Not Found"

**Symptoms:**
- Tasks are created successfully (status code 20100)
- When polling for results, DataForSEO returns 40400 "Not Found"
- Tasks never complete, frontend keeps polling indefinitely

**Error Logs:**
```
{"level":"warn","msg":"Task not found in DataForSEO (may have expired)","task_id":"11262159-1241-0066-0000-3176862afc2c"}
```

**Possible Causes:**
1. Tasks expire quickly in DataForSEO before we can retrieve them
2. Incorrect API endpoint format
3. Need to use `tasks_ready` endpoint first before `task_get`
4. Task IDs are being stored/retrieved incorrectly

**Recent Changes Made:**
- ✅ Added support for status code 20100 (Task Created) in addition to 20000
- ✅ Added `tasks_ready` endpoint support to check which tasks are ready
- ✅ Improved error handling to check JSON status codes even when HTTP is 200
- ✅ Added 5-second delay filter before checking tasks
- ✅ Enhanced logging throughout the polling process

**Still Need to Verify:**
- [ ] Confirm the correct API endpoint format from DataForSEO docs
- [ ] Test if `tasks_ready` endpoint works correctly
- [ ] Verify task IDs are being stored correctly in database
- [ ] Check if tasks expire too quickly (may need faster polling or webhooks)

## API Endpoints Being Used

**Current Implementation:**
- `POST /v3/serp/google/organic/task_post` - Create task ✅ Working
- `GET /v3/serp/google/organic/tasks_ready` - List ready tasks ✅ Added
- `GET /v3/serp/google/organic/task_get/{id}` - Get task result ❌ Getting 40400

**Documentation Reference:**
- https://docs.dataforseo.com/v3/serp/google/organic/task_post
- https://docs.dataforseo.com/v3/serp/google/organic/task_get
- https://docs.dataforseo.com/v3/serp/google/organic/tasks_ready

## Next Steps When Resuming

1. **Verify API Endpoint Format**
   - Check DataForSEO documentation for exact endpoint format
   - Verify if task IDs need any special formatting
   - Test endpoint directly with curl/Postman

2. **Test `tasks_ready` Endpoint**
   - Verify it returns the correct task IDs
   - Check if it filters out expired tasks
   - Ensure we're only checking tasks that are actually ready

3. **Check Task Expiration**
   - DataForSEO tasks may expire quickly
   - Consider implementing webhooks (`postback_url`) instead of polling
   - Or reduce polling interval if tasks expire too fast

4. **Database Verification**
   - Verify task IDs are being stored correctly in `keyword_tasks` table
   - Check if `dataforseo_task_id` matches what DataForSEO expects
   - Compare stored IDs with what `tasks_ready` returns

5. **Enhanced Logging**
   - Add logging to see raw API responses
   - Log task creation response details
   - Log `tasks_ready` response to see what IDs are returned

## Files Modified Recently

- `internal/dataforseo/client.go` - Added `tasks_ready` endpoint, improved error handling
- `internal/dataforseo/types.go` - Added `OrganicTasksReadyResponse` type
- `internal/api/keyword_handlers.go` - Accept status code 20100, improved logging
- `internal/api/keyword_task_poller.go` - Use `tasks_ready` first, enhanced logging
- `web/src/routes/RankTracker.svelte` - Frontend polling logic

## Testing Checklist

When resuming, test:
- [ ] Create a new keyword check
- [ ] Verify task is created (check logs for task ID)
- [ ] Wait 10-15 seconds
- [ ] Check `tasks_ready` endpoint response
- [ ] Verify task ID appears in ready list
- [ ] Try `task_get` with that task ID
- [ ] Check if snapshot is created successfully
- [ ] Verify frontend shows the position

## Environment Variables Required

```bash
DATAFORSEO_BASE_URL=https://api.dataforseo.com
DATAFORSEO_LOGIN=your_email@example.com
DATAFORSEO_PASSWORD=your_api_password
```

## Useful Commands

```bash
# Check server logs for DataForSEO activity
tail -f logs/api.log | grep -i dataforseo

# Check database for keyword tasks
psql -d barracuda -c "SELECT id, dataforseo_task_id, status, error FROM keyword_tasks ORDER BY created_at DESC LIMIT 10;"

# Check keyword snapshots
psql -d barracuda -c "SELECT id, keyword_id, position_organic, checked_at FROM keyword_rank_snapshots ORDER BY checked_at DESC LIMIT 10;"
```

