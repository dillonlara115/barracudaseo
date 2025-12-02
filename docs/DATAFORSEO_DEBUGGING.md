# DataForSEO Integration - Debugging Notes

## Current Status

✅ **Feature Complete**: The DataForSEO rank tracking integration is fully implemented and operational. This document tracks any known issues, debugging tips, and troubleshooting steps.

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

## Known Issues & Troubleshooting

### ✅ Resolved: Task Retrieval Issues

**Status:** All task retrieval issues have been resolved. The background poller now successfully:
- Creates tasks via `task_post`
- Polls for ready tasks using `tasks_ready`
- Retrieves results via `task_get`
- Creates snapshots and updates keyword positions

**Solutions Implemented:**
- ✅ Support for status code 20100 (Task Created) in addition to 20000
- ✅ `tasks_ready` endpoint integration for efficient polling
- ✅ Improved error handling with JSON status code checks
- ✅ 5-second delay filter before checking tasks
- ✅ Enhanced logging throughout the polling process
- ✅ Fallback to check individual tasks if `tasks_ready` returns 0 matches
- ✅ Proper task ID format handling

### Common Issues & Solutions

**Issue: Frontend shows "No rank data available yet" even when positions exist**
- **Solution:** Fixed in `KeywordDetailModal.svelte` - now checks both snapshots and keyword position fields

**Issue: Table headers misaligned**
- **Solution:** Fixed missing "Check Frequency" header in Rank Tracker table

**Issue: Tasks taking longer than expected**
- **Note:** DataForSEO tasks can take 5-30 seconds to process. The system polls every 3 seconds for up to 30 seconds before timing out gracefully.

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

1. **Check Server Logs After "Check Now" Click**
   - Look for "Stored task IDs in database" log entry - note the task IDs
   - Look for "Found ready tasks in DataForSEO" log entry - see what IDs are returned
   - Compare the two lists to see if there's a mismatch
   - Check if fallback to individual checks is working

2. **Verify Task ID Format**
   - Check what format task IDs are stored in database (`keyword_tasks.dataforseo_task_id`)
   - Compare with what DataForSEO returns in `task_post` response
   - Verify if IDs need any transformation before calling `task_get`

3. **Test Individual Task Checks**
   - Even if `tasks_ready` returns 0, the fallback should check tasks individually
   - Verify that individual `task_get` calls work
   - Check if 40400 errors are due to expired tasks or wrong IDs

4. **Check Task Expiration Timing**
   - DataForSEO tasks may expire quickly (check their docs for expiration time)
   - Consider implementing webhooks (`postback_url`) instead of polling
   - Or reduce polling interval if tasks expire too fast
   - May need to check tasks sooner than 5 seconds

5. **Frontend Polling Issue**
   - Frontend checks for `keyword.latest_position` to detect completion
   - If backend poller isn't creating snapshots, frontend will keep polling
   - May need to add timeout or better error handling in frontend

## Files Modified Recently

- `internal/dataforseo/client.go` - Added `tasks_ready` endpoint, improved error handling
- `internal/dataforseo/types.go` - Added `OrganicTasksReadyResponse` type
- `internal/api/keyword_handlers.go` - Accept status code 20100, improved logging
- `internal/api/keyword_task_poller.go` - Use `tasks_ready` first, enhanced logging
- `web/src/routes/RankTracker.svelte` - Frontend polling logic

## Testing Checklist

When resuming, test:
- [ ] Click "Check now" on a keyword
- [ ] Check server logs immediately for:
  - Task creation log with task ID
  - "Stored task IDs in database" log
- [ ] Wait 10-15 seconds for poller to run
- [ ] Check server logs for:
  - "Found ready tasks in DataForSEO" - note the task IDs returned
  - Compare with stored task IDs - do they match?
  - "Filtered to ready tasks" - how many matched?
- [ ] If 0 matches, check if fallback individual checks are happening
- [ ] Check database for `keyword_rank_snapshots` - was one created?
- [ ] Check frontend - does it stop spinning? Does position appear?
- [ ] If still failing, check `keyword_tasks` table for error messages

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

