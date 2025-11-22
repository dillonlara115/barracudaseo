<script>
  import { generateCrawlSummary } from '../../lib/data.js';
  import { Sparkles, Copy } from 'lucide-svelte';

  export let crawlId = null;

  let loading = false;
  let summary = null;
  let error = null;
  let copied = false;

  async function handleGenerateSummary() {
    if (!crawlId) return;

    loading = true;
    error = null;
    summary = null;

    try {
      const { data, error: apiError } = await generateCrawlSummary(crawlId);
      if (apiError) {
        error = apiError.message || 'Failed to generate summary';
      } else {
        summary = data?.summary || null;
      }
    } catch (err) {
      error = err.message || 'An unexpected error occurred';
    } finally {
      loading = false;
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
      {#if summary}
        <button
          class="btn btn-sm btn-ghost"
          on:click={handleCopyToClipboard}
          title="Copy to clipboard"
        >
          <Copy class="w-4 h-4" />
          {copied ? 'Copied!' : 'Copy'}
        </button>
      {/if}
    </div>

    <!-- Generate Button -->
    {#if !summary && !loading}
      <div class="mb-4">
        <button
          class="btn btn-primary w-full"
          on:click={handleGenerateSummary}
          disabled={!crawlId}
        >
          <Sparkles class="w-4 h-4" />
          Generate AI Crawl Summary
        </button>
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
      <div class="prose prose-sm max-w-none bg-base-200 p-6 rounded-lg whitespace-pre-wrap">
        {summary}
      </div>
      <div class="mt-4">
        <button
          class="btn btn-outline w-full"
          on:click={handleGenerateSummary}
        >
          Regenerate Summary
        </button>
      </div>
    {/if}
  </div>
</div>

<style>
  .prose {
    color: inherit;
  }
  .prose p {
    margin-top: 0.75em;
    margin-bottom: 0.75em;
  }
  .prose h1, .prose h2, .prose h3 {
    margin-top: 1em;
    margin-bottom: 0.5em;
    font-weight: bold;
  }
  .prose ul, .prose ol {
    margin-top: 0.5em;
    margin-bottom: 0.5em;
    padding-left: 1.5em;
  }
  .prose li {
    margin-top: 0.25em;
    margin-bottom: 0.25em;
  }
</style>

