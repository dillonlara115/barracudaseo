# Crawl Resumption & Persistence Options

## Problem Statement

Currently, when users:
1. Open a new browser tab and navigate away from the crawl progress page
2. Experience browser tab stalls or freezes
3. Close and reopen the browser

The crawl progress tracking is lost because:
- Frontend polling stops when the component unmounts (`onDestroy` clears the interval)
- No persistence of active crawl ID in browser storage
- No automatic reconnection when returning to the page
- Crawls continue running on the backend, but users lose visibility

## Current Implementation

### Frontend (`CrawlProgress.svelte`)
- Polls crawl status every 1 second via `setInterval`
- Clears interval in `onDestroy()` when component unmounts
- No persistence of `crawlId` in localStorage/sessionStorage
- No detection of page visibility changes

### Backend (`handlers.go`)
- Crawls run asynchronously in goroutines (`runCrawlAsync`)
- Crawl state persisted in database (`crawls` table)
- Status updates: `pending` → `running` → `succeeded`/`failed`
- Pages stored incrementally in batches (every 50 pages)
- No pause/resume capability - crawls run to completion

## Solution Options

### Option 1: Browser Storage Persistence + Auto-Reconnection (Recommended - Quick Win)

**Implementation:**
- Store active crawl ID in `localStorage` when crawl starts
- On page load, check for active crawl ID and reconnect if crawl is still running
- Use Page Visibility API to pause/resume polling when tab becomes hidden/visible
- Show notification/banner when returning to a page with an active crawl

**Pros:**
- Simple to implement
- No backend changes required
- Works immediately
- Low risk

**Cons:**
- Still relies on polling (less efficient)
- Doesn't solve the "tab stall" issue completely
- Requires user to return to the page

**Files to Modify:**
- `web/src/components/CrawlProgress.svelte` - Add localStorage persistence
- `web/src/components/TriggerCrawlButton.svelte` - Store crawl ID on start
- `web/src/routes/CrawlView.svelte` - Check for active crawl on mount
- `web/src/routes/ProjectView.svelte` - Check for active crawl on mount

**Implementation Details:**
```javascript
// Store crawl ID when crawl starts
localStorage.setItem(`activeCrawl_${projectId}`, crawlId);

// On page load, check for active crawl
const activeCrawlId = localStorage.getItem(`activeCrawl_${projectId}`);
if (activeCrawlId) {
  // Fetch crawl status
  const { data: crawl } = await fetchCrawl(activeCrawlId);
  if (crawl?.status === 'running') {
    // Show reconnection UI or auto-redirect
  }
}

// Use Page Visibility API
document.addEventListener('visibilitychange', () => {
  if (document.hidden) {
    // Tab hidden - could pause polling or reduce frequency
  } else {
    // Tab visible - resume normal polling
  }
});
```

---

### Option 2: WebSocket/Server-Sent Events (SSE) for Real-Time Updates

**Implementation:**
- Replace polling with WebSocket or SSE connection
- Backend pushes updates when crawl progress changes
- Connection persists across tab switches (with reconnection logic)
- More efficient than polling

**Pros:**
- Real-time updates (no polling delay)
- More efficient (server pushes only when needed)
- Better user experience
- Connection can survive tab switches with reconnection

**Cons:**
- Requires backend changes (WebSocket/SSE endpoint)
- More complex implementation
- Need to handle connection failures gracefully
- May require connection pooling/management

**Backend Changes Needed:**
- Add WebSocket/SSE endpoint: `/api/v1/crawls/{id}/stream`
- Broadcast crawl progress updates to connected clients
- Handle connection lifecycle (connect, disconnect, reconnect)

**Frontend Changes:**
- Replace `setInterval` polling with WebSocket/SSE client
- Implement reconnection logic
- Handle connection errors gracefully

---

### Option 3: Service Worker for Background Polling

**Implementation:**
- Use Service Worker to poll crawl status even when tab is inactive
- Show browser notifications when crawl completes
- Persist crawl state in IndexedDB
- Allow users to return to crawl from notification

**Pros:**
- Works even when tab is closed
- Can show browser notifications
- More persistent than regular JavaScript

**Cons:**
- Service Workers can be complex
- Requires HTTPS (already have this)
- May have browser compatibility considerations
- Overkill for this use case

---

### Option 4: Backend Crawl Pause/Resume Support (Most Robust)

**Implementation:**
- Add `pause` and `resume` status to crawl state machine
- Store crawl state (visited URLs, queue state) in database
- Allow users to pause crawls manually
- Resume from where it left off

**Pros:**
- Most robust solution
- Users have full control
- Can pause/resume at any time
- Solves all edge cases

**Cons:**
- Most complex to implement
- Requires significant backend changes
- Need to persist crawl state (visited URLs, queue, etc.)
- May require refactoring crawler architecture

**Backend Changes Needed:**
- Add `pause` status to crawl state machine
- Store crawl state in database:
  - Visited URLs set
  - Queue state
  - Current depth
  - Current page count
- Modify crawler to check for pause status periodically
- Add resume endpoint that restores crawl state

**Database Schema Changes:**
```sql
ALTER TABLE crawls ADD COLUMN paused_at timestamptz;
ALTER TABLE crawls ADD COLUMN crawl_state jsonb; -- Store visited URLs, queue, etc.
```

**Crawler Changes:**
- Check pause status in worker loop
- Save state before pausing
- Restore state on resume

---

### Option 5: Hybrid Approach (Recommended for Production)

**Phase 1: Quick Win (Option 1)**
- Implement browser storage persistence
- Add auto-reconnection on page load
- Use Page Visibility API to optimize polling

**Phase 2: Enhanced (Option 2)**
- Replace polling with WebSocket/SSE
- Real-time updates
- Better performance

**Phase 3: Advanced (Option 4)**
- Add pause/resume capability
- Full control for users
- Most robust solution

---

## Recommended Implementation Plan

### Phase 1: Immediate Fix (1-2 days)

1. **Add localStorage persistence**
   - Store `activeCrawl_${projectId}` when crawl starts
   - Clear when crawl completes or fails

2. **Auto-reconnection on page load**
   - Check for active crawl on `ProjectView` and `CrawlView` mount
   - If crawl is running, show banner/notification
   - Option to navigate to crawl or dismiss

3. **Page Visibility API**
   - Reduce polling frequency when tab is hidden
   - Resume normal polling when tab becomes visible
   - Prevents unnecessary requests when user isn't watching

4. **Visual indicator**
   - Show "Crawl in progress" badge in navigation
   - Link to active crawl from project view

### Phase 2: Enhanced Experience (1 week)

1. **WebSocket/SSE implementation**
   - Replace polling with real-time updates
   - Better performance and user experience
   - Automatic reconnection on connection loss

2. **Browser notifications**
   - Notify user when crawl completes (even if tab is closed)
   - Click notification to return to results

### Phase 3: Advanced Features (2-3 weeks)

1. **Pause/Resume capability**
   - Allow users to pause long-running crawls
   - Resume from where it left off
   - Store crawl state in database

---

## Implementation Details for Phase 1

### 1. localStorage Persistence

**File: `web/src/components/TriggerCrawlButton.svelte`**
```javascript
async function handleSubmit() {
  // ... existing code ...
  
  if (crawlId) {
    // Store active crawl ID
    localStorage.setItem(`activeCrawl_${projectId}`, crawlId);
    activeCrawlId = crawlId;
    showProgress = true;
  }
}

function handleProgressComplete() {
  // Clear active crawl when complete
  localStorage.removeItem(`activeCrawl_${projectId}`);
  // ... rest of existing code ...
}
```

### 2. Auto-Reconnection

**File: `web/src/routes/ProjectView.svelte`**
```javascript
onMount(async () => {
  // ... existing code ...
  
  // Check for active crawl
  const activeCrawlId = localStorage.getItem(`activeCrawl_${projectId}`);
  if (activeCrawlId) {
    const { data: crawl } = await fetchCrawl(activeCrawlId);
    if (crawl?.status === 'running') {
      // Show notification or auto-redirect
      showActiveCrawlNotification = true;
      activeCrawl = crawl;
    } else {
      // Clean up if crawl is done
      localStorage.removeItem(`activeCrawl_${projectId}`);
    }
  }
});
```

### 3. Page Visibility API

**File: `web/src/components/CrawlProgress.svelte`**
```javascript
let isVisible = true;
let pollInterval = null;
let hiddenPollInterval = null; // Slower polling when hidden

onMount(async () => {
  // ... existing code ...
  
  // Page Visibility API
  document.addEventListener('visibilitychange', handleVisibilityChange);
  
  // Start normal polling
  startPolling();
});

function handleVisibilityChange() {
  isVisible = !document.hidden;
  
  if (isVisible) {
    // Tab visible - use normal polling
    if (hiddenPollInterval) {
      clearInterval(hiddenPollInterval);
      hiddenPollInterval = null;
    }
    if (!pollInterval) {
      startPolling();
    }
  } else {
    // Tab hidden - use slower polling
    if (pollInterval) {
      clearInterval(pollInterval);
      pollInterval = null;
    }
    hiddenPollInterval = setInterval(async () => {
      await loadCrawl();
    }, 5000); // Poll every 5 seconds when hidden
  }
}

function startPolling() {
  if (pollInterval) return;
  
  pollInterval = setInterval(async () => {
    await loadCrawl();
  }, 1000); // Normal polling every 1 second
}

onDestroy(() => {
  if (pollInterval) clearInterval(pollInterval);
  if (hiddenPollInterval) clearInterval(hiddenPollInterval);
  document.removeEventListener('visibilitychange', handleVisibilityChange);
});
```

### 4. Visual Indicator

**File: `web/src/components/Dashboard.svelte`**
```javascript
let activeCrawlId = null;

onMount(async () => {
  // Check for active crawl
  const stored = localStorage.getItem(`activeCrawl_${projectId}`);
  if (stored) {
    const { data: crawl } = await fetchCrawl(stored);
    if (crawl?.status === 'running') {
      activeCrawlId = stored;
    }
  }
});

// In template
{#if activeCrawlId}
  <div class="alert alert-info mb-4">
    <span>Crawl in progress</span>
    <a href="/project/{projectId}/crawl/{activeCrawlId}" class="btn btn-sm btn-primary">
      View Progress
    </a>
  </div>
{/if}
```

---

## Testing Checklist

- [ ] Start crawl, close tab, reopen - should reconnect
- [ ] Start crawl, navigate to different page, return - should show notification
- [ ] Start crawl, switch browser tabs - polling should continue (slower)
- [ ] Complete crawl - localStorage should be cleared
- [ ] Multiple projects - each should track its own active crawl
- [ ] Browser refresh during crawl - should reconnect
- [ ] Multiple browser tabs - should sync crawl state

---

## Future Enhancements

1. **Crawl History View**
   - Show all recent crawls (running, completed, failed)
   - Quick access to resume/view any crawl

2. **Crawl Scheduling**
   - Allow users to schedule crawls
   - Automatic resumption if scheduled crawl was interrupted

3. **Crawl Notifications**
   - Email notifications when crawl completes
   - Slack integration for crawl status

4. **Crawl Comparison**
   - Compare multiple crawls side-by-side
   - Track changes over time

---

**Last Updated:** {{ date }}





