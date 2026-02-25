<script>
  import { onMount } from 'svelte';
  import { generateCrawlSummary, getCrawlSummary, deleteCrawlSummary } from '../../lib/data.js';
  import { Sparkles, Copy, Trash2, Zap } from 'lucide-svelte';
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

  marked.setOptions({
    breaks: true,
    gfm: true,
  });

  $: renderedSummary = summary ? marked.parse(summary) : '';

  let previousCrawlId = null;

  onMount(async () => {
    if (crawlId) {
      await loadExistingSummary();
      previousCrawlId = crawlId;
    }
  });

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
    if (forceRefresh) summary = null;

    try {
      const { data, error: apiError } = await generateCrawlSummary(crawlId, forceRefresh);
      if (apiError) {
        error = apiError.message || 'Failed to generate summary';
      } else {
        summary = data?.summary || null;
        isCached = data?.cached || false;
      }
    } catch (err) {
      error = err.message || 'An unexpected error occurred';
    } finally {
      loading = false;
    }
  }

  async function handleDeleteSummary() {
    if (!crawlId || !summary) return;
    if (!confirm('Are you sure you want to delete this summary? You can regenerate it later.')) return;

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
      setTimeout(() => { copied = false; }, 2000);
    } catch (err) {
      console.error('Failed to copy to clipboard:', err);
    }
  }
</script>

<div class="mb-6">
  <!-- Loading Existing -->
  {#if loadingExisting}
    <div class="bg-base-200 rounded-xl px-5 py-4 flex items-center gap-3">
      <span class="loading loading-spinner loading-sm"></span>
      <span class="text-sm text-base-content/50">Loading summary...</span>
    </div>
  {/if}

  <!-- Generate Button -->
  {#if !summary && !loading && !loadingExisting}
    <div class="bg-primary/5 border border-primary/15 rounded-xl p-5">
      <div class="flex items-center gap-2 mb-3">
        <Zap class="w-4 h-4 text-primary" />
        <span class="text-sm font-semibold text-primary">AI Analysis</span>
      </div>
      <p class="text-sm text-base-content/50 mb-4">Generate an AI-powered summary of your crawl results with prioritized recommendations.</p>
      {#if isPro}
          <button
            class="btn btn-primary btn-sm"
            on:click={() => handleGenerateSummary()}
            disabled={!crawlId}
          >
          <Sparkles class="w-3.5 h-3.5" />
          Generate AI Summary
        </button>
      {:else}
        <div class="tooltip" data-tip="Upgrade to Pro to unlock AI Summaries">
          <button class="btn btn-primary btn-sm btn-disabled" disabled>
            <Sparkles class="w-3.5 h-3.5" />
            Generate AI Summary
            <span class="bg-white/20 text-xs px-1.5 py-0.5 rounded font-medium ml-1">PRO</span>
          </button>
        </div>
      {/if}
    </div>
  {/if}

  <!-- Loading State -->
  {#if loading}
    <div class="bg-primary/5 border border-primary/15 rounded-xl p-8 flex flex-col items-center">
      <span class="loading loading-spinner loading-md mb-3 text-primary"></span>
      <p class="text-sm text-base-content/60">Generating AI crawl summary...</p>
      <p class="text-xs text-base-content/30 mt-1">This may take a few moments</p>
    </div>
  {/if}

  <!-- Error State -->
  {#if error}
    <div class="bg-error/10 border border-error/20 rounded-xl px-5 py-3 mb-3">
      <p class="text-sm text-error">{error}</p>
    </div>
    <button class="btn btn-outline btn-sm" on:click={() => handleGenerateSummary()}>Try Again</button>
  {/if}

  <!-- Summary Display -->
  {#if summary}
    <div class="bg-primary/5 border border-primary/15 rounded-xl overflow-hidden">
      <!-- Header Bar -->
      <div class="flex items-center justify-between px-5 py-3 border-b border-primary/10">
        <div class="flex items-center gap-2">
          <Zap class="w-4 h-4 text-primary" />
          <span class="text-sm font-semibold text-primary">AI Analysis</span>
          {#if isCached}
            <span class="bg-base-content/5 text-base-content/40 text-xs px-2 py-0.5 rounded-full">Saved</span>
          {:else}
            <span class="bg-success/15 text-success text-xs px-2 py-0.5 rounded-full">New</span>
          {/if}
        </div>
        <div class="flex items-center gap-1">
          <button
            class="btn btn-ghost btn-xs"
            on:click={handleCopyToClipboard}
            title="Copy to clipboard"
          >
            <Copy class="w-3.5 h-3.5" />
            {copied ? 'Copied!' : ''}
          </button>
          <button
            class="btn btn-ghost btn-xs text-error"
            on:click={handleDeleteSummary}
            disabled={deleting}
            title="Delete summary"
          >
            {#if deleting}
              <span class="loading loading-spinner loading-xs"></span>
            {:else}
              <Trash2 class="w-3.5 h-3.5" />
            {/if}
          </button>
        </div>
      </div>

      <!-- Summary Content -->
      <div class="ai-summary-content">
        {@html renderedSummary}
      </div>

      <!-- Actions -->
      <div class="flex items-center gap-2 px-5 py-3 border-t border-primary/10">
        {#if isPro}
          <button
            class="btn btn-ghost btn-sm text-base-content/50"
            on:click={() => handleGenerateSummary(true)}
            disabled={loading}
          >
            {loading ? 'Regenerating...' : 'Regenerate'}
          </button>
        {:else}
          <div class="tooltip" data-tip="Upgrade to Pro to regenerate">
            <button class="btn btn-ghost btn-sm btn-disabled" disabled>
              Regenerate
              <span class="bg-primary/15 text-primary text-xs px-1.5 py-0.5 rounded font-medium ml-1">PRO</span>
            </button>
          </div>
        {/if}
      </div>
    </div>
  {/if}
</div>
