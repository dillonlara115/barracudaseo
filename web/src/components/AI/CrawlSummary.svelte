<script>
  import { generateCrawlSummary } from '../../lib/data.js';
  import { Sparkles, Copy } from 'lucide-svelte';
  import { marked } from 'marked';

  export let crawlId = null;

  let loading = false;
  let summary = null;
  let error = null;
  let copied = false;

  // Configure marked for safe rendering
  marked.setOptions({
    breaks: true,
    gfm: true,
  });

  $: renderedSummary = summary ? marked.parse(summary) : '';

  async function handleGenerateSummary() {
    if (!crawlId) {
      error = 'Crawl ID is required';
      return;
    }

    loading = true;
    error = null;
    summary = null;

    try {
      console.log('Generating crawl summary for crawlId:', crawlId, 'type:', typeof crawlId);
      const { data, error: apiError } = await generateCrawlSummary(crawlId);
      if (apiError) {
        console.error('Error generating crawl summary:', apiError);
        console.error('Full error object:', JSON.stringify(apiError, null, 2));
        error = apiError.message || 'Failed to generate summary';
      } else {
        console.log('Successfully generated summary:', data);
        summary = data?.summary || null;
      }
    } catch (err) {
      console.error('Exception generating crawl summary:', err);
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
      <div class="prose prose-sm max-w-none bg-base-200 p-6 rounded-lg">
        {@html renderedSummary}
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
    line-height: 1.6;
  }
  .prose h1, .prose h2, .prose h3, .prose h4 {
    margin-top: 1em;
    margin-bottom: 0.5em;
    font-weight: bold;
    line-height: 1.2;
  }
  .prose h1 {
    font-size: 1.5em;
  }
  .prose h2 {
    font-size: 1.25em;
  }
  .prose h3 {
    font-size: 1.1em;
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
  .prose code {
    background-color: rgba(0, 0, 0, 0.1);
    padding: 0.2em 0.4em;
    border-radius: 0.25em;
    font-size: 0.9em;
  }
  .prose pre {
    background-color: rgba(0, 0, 0, 0.1);
    padding: 1em;
    border-radius: 0.5em;
    overflow-x: auto;
    margin-top: 0.75em;
    margin-bottom: 0.75em;
  }
  .prose pre code {
    background-color: transparent;
    padding: 0;
  }
  .prose strong {
    font-weight: bold;
  }
  .prose em {
    font-style: italic;
  }
  .prose blockquote {
    border-left: 4px solid rgba(0, 0, 0, 0.2);
    padding-left: 1em;
    margin-left: 0;
    margin-top: 0.75em;
    margin-bottom: 0.75em;
    font-style: italic;
  }
</style>



