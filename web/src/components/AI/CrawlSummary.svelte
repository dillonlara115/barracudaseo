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
      <div class="mb-4">
        <div class="flex items-center justify-between mb-2">
          <h4 class="font-semibold text-lg">AI Crawl Summary</h4>
          <button
            class="btn btn-sm btn-ghost"
            on:click={handleCopyToClipboard}
            title="Copy to clipboard"
          >
            <Copy class="w-4 h-4" />
            {copied ? 'Copied!' : 'Copy'}
          </button>
        </div>
        <div class="prose prose-sm max-w-none bg-base-200 p-4 rounded-lg">
          {@html renderedSummary}
        </div>
      </div>
      <div class="mt-4">
        <button
          class="btn btn-outline w-full"
          on:click={() => handleGenerateSummary(true)}
          disabled={loading}
        >
          {loading ? 'Regenerating...' : 'Regenerate Summary'}
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
  /* DaisyUI heading styles */
  .prose h1 {
    font-size: 1.875rem; /* text-3xl */
    font-weight: 700; /* font-bold */
    margin-bottom: 1rem; /* mb-4 */
    margin-top: 1.5rem; /* mt-6 */
  }
  .prose h2 {
    font-size: 1.5rem; /* text-2xl */
    font-weight: 600; /* font-semibold */
    margin-bottom: 0.75rem; /* mb-3 */
    margin-top: 1.25rem; /* mt-5 */
  }
  .prose h3 {
    font-size: 1.25rem; /* text-xl */
    font-weight: 600; /* font-semibold */
    margin-bottom: 0.5rem; /* mb-2 */
    margin-top: 1rem; /* mt-4 */
  }
  .prose h4 {
    font-size: 1.125rem; /* text-lg */
    font-weight: 600; /* font-semibold */
    margin-bottom: 0.5rem; /* mb-2 */
    margin-top: 0.75rem; /* mt-3 */
  }
  .prose h5 {
    font-size: 1rem; /* text-base */
    font-weight: 600; /* font-semibold */
    margin-bottom: 0.5rem; /* mb-2 */
    margin-top: 0.5rem; /* mt-2 */
  }
  .prose h6 {
    font-size: 0.875rem; /* text-sm */
    font-weight: 600; /* font-semibold */
    margin-bottom: 0.25rem; /* mb-1 */
    margin-top: 0.5rem; /* mt-2 */
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
    background-color: rgba(0, 0, 0, 0.2); /* bg-base-300 equivalent */
    padding: 0.25rem 0.5rem; /* px-2 py-1 */
    border-radius: 0.25rem; /* rounded */
    font-size: 0.875rem; /* text-sm */
  }
  .prose pre {
    background-color: rgba(0, 0, 0, 0.2); /* bg-base-300 equivalent */
    padding: 1rem; /* p-4 */
    border-radius: 0.5rem; /* rounded-lg */
    overflow-x: auto;
    margin-top: 1rem; /* my-4 */
    margin-bottom: 1rem; /* my-4 */
  }
  .prose pre code {
    background-color: transparent;
    padding: 0;
  }
  .prose strong {
    font-weight: 700; /* font-bold */
  }
  .prose em {
    font-style: italic;
  }
  .prose blockquote {
    border-left: 4px solid rgba(0, 0, 0, 0.2); /* border-l-4 border-base-300 */
    padding-left: 1rem; /* pl-4 */
    margin-top: 1rem; /* my-4 */
    margin-bottom: 1rem; /* my-4 */
    font-style: italic;
  }
</style>



