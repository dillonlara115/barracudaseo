<script>
  import { onMount } from 'svelte';
  import { generateCrawlSummary, getCrawlSummary, deleteCrawlSummary } from '../../lib/data.js';
  import { Sparkles, Copy, Trash2 } from 'lucide-svelte';
  import { marked } from 'marked';
  import { userProfile, isProOrTeam } from '../../lib/subscription.js';

  $: isPro = isProOrTeam($userProfile);

  export let crawlId = null;

  let loading = false;
  let loadingExisting = false;
  let summary = null;
  let error = null;
  let copied = false;
  let isCached = false;
  let deleting = false;

  // Configure marked for safe rendering
  marked.setOptions({
    breaks: true,
    gfm: true,
  });

  $: renderedSummary = summary ? marked.parse(summary) : '';

  let previousCrawlId = null;

  // Load existing summary on mount and when crawlId changes
  onMount(async () => {
    if (crawlId) {
      await loadExistingSummary();
      previousCrawlId = crawlId;
    }
  });

  // Reload summary when crawlId changes
  $: if (crawlId && crawlId !== previousCrawlId) {
    loadExistingSummary();
    previousCrawlId = crawlId;
  }

  async function loadExistingSummary() {
    if (!crawlId) return;
    
    loadingExisting = true;
    error = null;
    try {
      const { data, error: apiError } = await getCrawlSummary(crawlId);
      if (apiError) {
        // Don't show error if no summary exists - that's expected
        if (apiError.message && !apiError.message.includes('not found')) {
          console.error('Error loading crawl summary:', apiError);
        }
      } else if (data && data.summary) {
        summary = data.summary;
        isCached = data.cached || false;
      }
    } catch (err) {
      console.error('Exception loading crawl summary:', err);
    } finally {
      loadingExisting = false;
    }
  }

  async function handleGenerateSummary(forceRefresh = false) {
    if (!crawlId) {
      error = 'Crawl ID is required';
      return;
    }

    loading = true;
    error = null;
    // Only clear summary if forcing refresh, otherwise keep it visible while loading
    if (forceRefresh) {
      summary = null;
    }

    try {
      console.log('Generating crawl summary for crawlId:', crawlId, 'type:', typeof crawlId, 'forceRefresh:', forceRefresh);
      const { data, error: apiError } = await generateCrawlSummary(crawlId, forceRefresh);
      if (apiError) {
        console.error('Error generating crawl summary:', apiError);
        console.error('Full error object:', JSON.stringify(apiError, null, 2));
        error = apiError.message || 'Failed to generate summary';
      } else {
        console.log('Successfully generated summary:', data);
        summary = data?.summary || null;
        isCached = data?.cached || false;
      }
    } catch (err) {
      console.error('Exception generating crawl summary:', err);
      error = err.message || 'An unexpected error occurred';
    } finally {
      loading = false;
    }
  }

  async function handleDeleteSummary() {
    if (!crawlId || !summary) return;
    
    if (!confirm('Are you sure you want to delete this summary? You can regenerate it later.')) {
      return;
    }

    deleting = true;
    error = null;
    try {
      const { error: apiError } = await deleteCrawlSummary(crawlId);
      if (apiError) {
        error = apiError.message || 'Failed to delete summary';
      } else {
        summary = null;
        isCached = false;
      }
    } catch (err) {
      console.error('Exception deleting crawl summary:', err);
      error = err.message || 'An unexpected error occurred';
    } finally {
      deleting = false;
    }
  }

  async function handleCopyToClipboard() {
    if (!summary) return;

    try {
      await navigator.clipboard.writeText(summary);
      copied = true;
      setTimeout(() => {
        copied = false;
      }, 2000);
    } catch (err) {
      console.error('Failed to copy to clipboard:', err);
    }
  }
</script>

<div class="card bg-base-100 shadow-lg mb-6">
  <div class="card-body">
    <div class="flex items-center justify-between mb-4">
      <h2 class="card-title text-xl flex items-center gap-2">
        <Sparkles class="w-6 h-6" />
        AI Crawl Summary
      </h2>
    </div>

    <!-- Loading Existing Summary -->
    {#if loadingExisting}
      <div class="flex items-center justify-center py-4">
        <span class="loading loading-spinner loading-sm mr-2"></span>
        <span class="text-sm text-base-content/70">Loading summary...</span>
      </div>
    {/if}

    <!-- Generate Button -->
    {#if !summary && !loading && !loadingExisting}
      <div class="mb-4">
        {#if isPro}
          <button
            class="btn btn-primary w-full"
            on:click={handleGenerateSummary}
            disabled={!crawlId}
          >
            <Sparkles class="w-4 h-4" />
            Generate AI Crawl Summary
          </button>
        {:else}
          <div class="tooltip w-full" data-tip="Upgrade to Pro to unlock AI Crawl Summaries">
            <button class="btn btn-primary btn-disabled w-full" disabled>
              <Sparkles class="w-4 h-4" />
              Generate AI Crawl Summary
              <span class="badge badge-primary badge-sm ml-1">PRO</span>
            </button>
          </div>
        {/if}
        {#if !crawlId}
          <p class="text-sm text-base-content/70 mt-2">Crawl ID is required to generate summary</p>
        {/if}
      </div>
    {/if}

    <!-- Loading State -->
    {#if loading}
      <div class="flex flex-col items-center justify-center py-8">
        <span class="loading loading-spinner loading-lg mb-4"></span>
        <p class="text-base-content/70">Generating AI crawl summary...</p>
        <p class="text-sm text-base-content/50 mt-2">This may take a few moments</p>
      </div>
    {/if}

    <!-- Error State -->
    {#if error}
      <div class="alert alert-error mb-4">
        <span>{error}</span>
      </div>
      <button
        class="btn btn-outline w-full"
        on:click={handleGenerateSummary}
      >
        Try Again
      </button>
    {/if}

    <!-- Summary Display -->
    {#if summary}
      <div class="mb-4">
        <div class="flex items-center justify-between mb-3">
          {#if isCached}
            <span class="badge badge-sm badge-info">Saved Summary</span>
          {:else}
            <span class="badge badge-sm badge-success">New Summary</span>
          {/if}
          <div class="flex gap-2">
            <button
              class="btn btn-sm btn-ghost"
              on:click={handleCopyToClipboard}
              title="Copy to clipboard"
            >
              <Copy class="w-4 h-4" />
              {copied ? 'Copied!' : 'Copy'}
            </button>
            <button
              class="btn btn-sm btn-ghost text-error"
              on:click={handleDeleteSummary}
              disabled={deleting}
              title="Delete summary"
            >
              {#if deleting}
                <span class="loading loading-spinner loading-xs"></span>
              {:else}
                <Trash2 class="w-4 h-4" />
              {/if}
            </button>
          </div>
        </div>
        <div class="ai-summary-content">
          {@html renderedSummary}
        </div>
      </div>
      <div class="mt-4 flex gap-2">
        {#if isPro}
          <button
            class="btn btn-outline flex-1"
            on:click={() => handleGenerateSummary(true)}
            disabled={loading}
          >
            {loading ? 'Regenerating...' : 'Regenerate Summary'}
          </button>
        {:else}
          <div class="tooltip flex-1" data-tip="Upgrade to Pro to regenerate summaries">
            <button class="btn btn-outline btn-disabled flex-1" disabled>
              Regenerate Summary
              <span class="badge badge-primary badge-sm ml-1">PRO</span>
            </button>
          </div>
        {/if}
        <button
          class="btn btn-error"
          on:click={handleDeleteSummary}
          disabled={deleting || loading}
        >
          {#if deleting}
            <span class="loading loading-spinner loading-sm"></span>
            Deleting...
          {:else}
            <Trash2 class="w-4 h-4" />
            Delete
          {/if}
        </button>
      </div>
    {/if}
  </div>
</div>

<style>
</style>



