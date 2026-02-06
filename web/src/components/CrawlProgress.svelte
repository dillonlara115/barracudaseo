<script>
  import { onMount, onDestroy, createEventDispatcher } from 'svelte';
  import { fetchCrawl, fetchCrawlPageCount } from '../lib/data.js';
  import { push } from 'svelte-spa-router';
  
  const dispatch = createEventDispatcher();
  
  export let crawlId = null;
  export let projectId = null;
  
  let crawl = null;
  let pageCount = 0;
  let loading = true;
  let error = null;
  let pollInterval = null;
  let hiddenPollInterval = null; // Slower polling when tab is hidden
  let isVisible = true;
  let retryCount = 0;
  const MAX_RETRIES = 10; // Retry for up to 20 seconds (10 retries * 2 seconds)
  
  const toNumber = (value) => {
    const parsed = Number(value);
    return Number.isFinite(parsed) ? parsed : 0;
  };
  
  // Progress calculation - use max_pages from crawl (limit), not total_pages (current count)
  $: maxPages = crawl?.max_pages || crawl?.meta?.max_pages || 100;
  $: status = crawl?.status || 'pending';
  
  // Use actual page count from DB, but fallback to crawl.total_pages if available
  $: displayPageCount = (() => {
    const counts = [
      toNumber(pageCount),
      toNumber(crawl?.page_count),
      toNumber(crawl?.indexed_pages),
      toNumber(crawl?.total_pages)
    ];
    return Math.max(...counts, 0);
  })();
  $: progress = maxPages > 0 ? Math.min((displayPageCount / maxPages) * 100, 100) : 0;
  
  // ETA calculation
  let startTime = null;
  let lastPageCount = 0;
  let lastUpdateTime = null;
  let pagesPerSecond = 0;
  let etaSeconds = null;
  
  $: if (crawl?.started_at && status === 'running') {
    const now = new Date();
    const started = new Date(crawl.started_at);
    const elapsed = (now - started) / 1000; // seconds
    const currentCount = displayPageCount;
    
    if (currentCount > 0) {
      if (lastUpdateTime && lastPageCount > 0 && lastPageCount !== currentCount) {
        const timeDiff = (now - lastUpdateTime) / 1000;
        const pageDiff = currentCount - lastPageCount;
        if (timeDiff > 0 && pageDiff > 0) {
          // Use recent rate (smoothed)
          const recentRate = pageDiff / timeDiff;
          const overallRate = currentCount / elapsed;
          pagesPerSecond = (recentRate * 0.7) + (overallRate * 0.3); // Weight recent rate more
        }
      } else if (elapsed > 0) {
        // Use overall rate if no recent update
        pagesPerSecond = currentCount / elapsed;
      }
      
      if (pagesPerSecond > 0 && currentCount < maxPages) {
        const remainingPages = maxPages - currentCount;
        etaSeconds = Math.ceil(remainingPages / pagesPerSecond);
      } else {
        etaSeconds = null;
      }
    }
    
    // Update tracking
    if (lastPageCount !== currentCount) {
      lastPageCount = currentCount;
      lastUpdateTime = now;
    }
  }
  
  function formatDuration(seconds) {
    if (!seconds || seconds < 0) return null;
    if (seconds < 60) return `${Math.round(seconds)}s`;
    if (seconds < 3600) return `${Math.round(seconds / 60)}m`;
    return `${Math.round(seconds / 3600)}h ${Math.round((seconds % 3600) / 60)}m`;
  }
  
  async function loadCrawl() {
    if (!crawlId) {
      console.log('CrawlProgress: No crawlId provided');
      return;
    }
    
    console.log('CrawlProgress: Loading crawl', crawlId);
    
    try {
      // Fetch crawl from backend API (now includes real-time page count)
      const crawlResult = await fetchCrawl(crawlId);
      
      if (crawlResult.error) {
        // If crawl not found and we haven't loaded yet, it might still be creating
        // Retry for a reasonable amount of time before showing error
        if (crawlResult.error.message && crawlResult.error.message.includes('not found') && !crawl) {
          retryCount++;
          if (retryCount < MAX_RETRIES) {
            console.log(`CrawlProgress: Crawl not found yet, retrying (${retryCount}/${MAX_RETRIES})`, crawlId);
            return; // Don't set error, will retry on next poll
          } else {
            console.error('CrawlProgress: Crawl not found after multiple retries', crawlId);
            error = 'Crawl not found. It may still be initializing. Please refresh the page.';
            loading = false;
            return;
          }
        }
        console.error('CrawlProgress: Error fetching crawl', crawlResult.error);
        throw crawlResult.error;
      }
      
      // Reset retry count on successful fetch
      retryCount = 0;
      
      crawl = crawlResult.data;
      
      // Use whichever count is currently highest to avoid regressions while the crawl streams in
      const counts = [
        toNumber(pageCount),
        toNumber(crawl?.page_count),
        toNumber(crawl?.total_pages),
        toNumber(crawl?.indexed_pages)
      ];
      const bestCount = Math.max(...counts, 0);
      if (bestCount > 0) {
        pageCount = bestCount;
      } else {
        // Fallback: fetch page count separately if not in crawl data
        const pageCountResult = await fetchCrawlPageCount(crawlId);
        if (!pageCountResult.error) {
          pageCount = pageCountResult.count;
        }
      }
      
      console.log('CrawlProgress: Loaded', { 
        crawl: !!crawl, 
        status: crawl?.status,
        crawlStatus: crawl?.status,
        pageCount,
        pageCountFromBackend: crawl?.page_count,
        totalPagesFromCrawl: crawl?.total_pages,
        maxPages: crawl?.max_pages || crawl?.meta?.max_pages,
        meta: crawl?.meta,
        startedAt: crawl?.started_at,
        completedAt: crawl?.completed_at
      });
      
      if (!startTime && crawl?.started_at) {
        startTime = new Date(crawl.started_at);
        lastUpdateTime = startTime;
      }
      
      loading = false;
      
      // Stop polling if crawl is complete
      if (status === 'succeeded' || status === 'failed' || status === 'cancelled') {
        stopPolling();
        
        // Clear localStorage when crawl completes
        if (projectId && crawlId) {
          localStorage.removeItem(`activeCrawl_${projectId}`);
          console.log('CrawlProgress: Cleared active crawl from localStorage:', `activeCrawl_${projectId}`);
        }
      }
    } catch (err) {
      console.error('CrawlProgress: Error', err);
      error = err.message;
      loading = false;
    }
  }
  
  function handleVisibilityChange() {
    isVisible = !document.hidden;
    
    if (isVisible) {
      // Tab visible - use normal polling
      if (hiddenPollInterval) {
        clearInterval(hiddenPollInterval);
        hiddenPollInterval = null;
      }
      if (!pollInterval && (status === 'running' || status === 'pending')) {
        startPolling();
      }
    } else {
      // Tab hidden - use slower polling
      if (pollInterval) {
        clearInterval(pollInterval);
        pollInterval = null;
      }
      if (!hiddenPollInterval && (status === 'running' || status === 'pending')) {
        hiddenPollInterval = setInterval(async () => {
          await loadCrawl();
        }, 5000); // Poll every 5 seconds when hidden
      }
    }
  }
  
  function startPolling() {
    if (pollInterval) return;
    
    pollInterval = setInterval(async () => {
      await loadCrawl();
    }, 1000); // Normal polling every 1 second
  }
  
  function stopPolling() {
    if (pollInterval) {
      clearInterval(pollInterval);
      pollInterval = null;
    }
    if (hiddenPollInterval) {
      clearInterval(hiddenPollInterval);
      hiddenPollInterval = null;
    }
  }
  
  onMount(async () => {
    console.log('CrawlProgress: onMount called with crawlId:', crawlId);
    
    // Check localStorage for active crawl if crawlId not provided (reconnection scenario)
    if (!crawlId && projectId) {
      const storedCrawlId = localStorage.getItem(`activeCrawl_${projectId}`);
      if (storedCrawlId) {
        console.log('CrawlProgress: Found active crawl in localStorage:', storedCrawlId);
        crawlId = storedCrawlId;
      }
    }
    
    if (!crawlId) {
      console.error('CrawlProgress: No crawlId provided in onMount');
      error = 'No crawl ID provided';
      loading = false;
      return;
    }
    
    // Set up Page Visibility API
    document.addEventListener('visibilitychange', handleVisibilityChange);
    isVisible = !document.hidden;
    
    // Initial load with a small delay to ensure crawl is created
    await new Promise(resolve => setTimeout(resolve, 500));
    await loadCrawl();
    
    // Start polling based on visibility
    if (isVisible) {
      startPolling();
    } else {
      // If hidden, use slower polling
      hiddenPollInterval = setInterval(async () => {
        await loadCrawl();
      }, 5000);
    }
  });
  
  onDestroy(() => {
    stopPolling();
    document.removeEventListener('visibilitychange', handleVisibilityChange);
  });
  
  function handleViewResults() {
    if (projectId && crawlId) {
      push(`/project/${projectId}/crawl/${crawlId}`);
    }
  }
</script>

{#if loading}
  <div class="alert alert-info">
    <span class="loading loading-spinner loading-sm"></span>
    <span>Loading crawl status...</span>
  </div>
{:else if error}
  <div class="alert alert-error">
    <span>Error: {error}</span>
  </div>
{:else if crawl}
  <div class="card bg-base-200 text-base-content shadow-lg">
    <div class="card-body">
      <div class="flex justify-between items-start mb-4">
        <h3 class="card-title">Crawl Progress</h3>
        <div class="badge badge-lg {status === 'succeeded' ? 'badge-success' : status === 'failed' ? 'badge-error' : status === 'running' ? 'badge-info' : 'badge-warning'}">
          {status === 'succeeded' ? 'Complete' : status === 'failed' ? 'Failed' : status === 'running' ? 'Running' : 'Pending'}
        </div>
      </div>
      
      {#if status === 'running' || status === 'pending'}
        <!-- Progress Bar -->
        <div class="mb-4">
          <div class="flex justify-between items-center mb-2">
            <span class="text-sm font-semibold">{displayPageCount} / {maxPages} pages</span>
            <span class="text-sm">{Math.round(progress)}%</span>
          </div>
          <progress 
            class="progress progress-primary w-full" 
            value={displayPageCount} 
            max={maxPages}
          ></progress>
        </div>
        
        <!-- Stats -->
        <div class="grid grid-cols-3 gap-4 mb-4">
          <div class="stat py-0">
            <div class="stat-title text-xs">Pages Crawled</div>
            <div class="stat-value text-2xl">{displayPageCount}</div>
          </div>
          <div class="stat py-0">
            <div class="stat-title text-xs">Rate</div>
            <div class="stat-value text-2xl">{pagesPerSecond > 0 ? pagesPerSecond.toFixed(1) : '0.0'}/s</div>
          </div>
          <div class="stat py-0">
            <div class="stat-title text-xs">ETA</div>
            <div class="stat-value text-2xl">{etaSeconds ? formatDuration(etaSeconds) : '—'}</div>
          </div>
        </div>
      {:else if status === 'succeeded'}
        <!-- Success State -->
        <div class="alert alert-success mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <div>
            <h3 class="font-bold">Crawl Completed Successfully!</h3>
            <div class="text-sm mt-1">
              Crawled {crawl.total_pages || pageCount} pages and found {crawl.total_issues || 0} issues
            </div>
          </div>
        </div>
        <div class="card-actions justify-end">
          <button class="btn btn-primary" on:click={() => {
            // Dispatch completed event - parent will handle navigation and closing modal
            dispatch('completed', { crawlId });
          }}>
            View Results
          </button>
        </div>
      {:else if status === 'failed'}
        <!-- Failed State -->
        <div class="alert alert-error mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <div>
            <h3 class="font-bold">Crawl Failed</h3>
            <div class="text-sm mt-1">
              {crawl.meta?.error || 'An error occurred during the crawl'}
            </div>
          </div>
        </div>
      {/if}
    </div>
  </div>
{/if}
