<script>
  import { createEventDispatcher } from 'svelte';
  import { generateIssueInsight } from '../../lib/data.js';
  import { Sparkles, Copy, X } from 'lucide-svelte';
  import { marked } from 'marked';

  export let issue = null;
  export let crawlId = null;

  const dispatch = createEventDispatcher();

  let loading = false;
  let insight = null;
  let error = null;
  let copied = false;

  // Configure marked for safe rendering
  marked.setOptions({
    breaks: true,
    gfm: true,
  });

  // Parse recommendation and insight from response
  $: parsedResponse = (() => {
    if (!insight) return { recommendation: null, insight: null };
    
    // Check if response follows RECOMMENDATION: / INSIGHT: format
    const recMatch = insight.match(/RECOMMENDATION:\s*([\s\S]+?)(?:\n\nINSIGHT:|$)/i);
    const insightMatch = insight.match(/INSIGHT:\s*([\s\S]+?)$/i);
    
    if (recMatch) {
      return {
        recommendation: recMatch[1].trim(),
        insight: insightMatch ? insightMatch[1].trim() : null
      };
    }
    
    // If no structured format, return full insight
    return {
      recommendation: null,
      insight: insight
    };
  })();

  $: renderedRecommendation = parsedResponse.recommendation ? marked.parse(parsedResponse.recommendation) : '';
  $: renderedInsight = parsedResponse.insight ? marked.parse(parsedResponse.insight) : '';

  async function handleGenerateInsight() {
    if (!issue || !crawlId) {
      error = 'Issue or crawl ID is missing';
      return;
    }

    // Check if issue has an ID
    const issueId = issue.id || issue.issue_id;
    if (!issueId) {
      error = 'Issue ID is missing';
      console.error('Issue object:', issue);
      return;
    }

    loading = true;
    error = null;
    insight = null;

    try {
      console.log('Generating insight for issue:', { issueId, crawlId, issue });
      const { data, error: apiError } = await generateIssueInsight(issueId, crawlId);
      if (apiError) {
        error = apiError.message || 'Failed to generate insight';
      } else {
        insight = data?.insight || null;
      }
    } catch (err) {
      error = err.message || 'An unexpected error occurred';
    } finally {
      loading = false;
    }
  }

  async function handleCopyToClipboard() {
    if (!insight) return;

    try {
      await navigator.clipboard.writeText(insight);
      copied = true;
      setTimeout(() => {
        copied = false;
      }, 2000);
    } catch (err) {
      console.error('Failed to copy to clipboard:', err);
    }
  }

  async function handleCopyRecommendation() {
    if (!parsedResponse.recommendation) return;

    try {
      await navigator.clipboard.writeText(parsedResponse.recommendation);
      copied = true;
      setTimeout(() => {
        copied = false;
      }, 2000);
    } catch (err) {
      console.error('Failed to copy to clipboard:', err);
    }
  }

  function handleClose() {
    dispatch('close');
  }
</script>

{#if issue}
  <div class="modal modal-open">
    <div class="modal-box max-w-3xl max-h-[90vh] overflow-y-auto">
      <div class="flex items-center justify-between mb-4">
        <h3 class="font-bold text-2xl flex items-center gap-2">
          <Sparkles class="w-6 h-6" />
          AI Issue Insight
        </h3>
        <button class="btn btn-sm btn-circle btn-ghost" on:click={handleClose}>
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Issue Info -->
      <div class="mb-4 p-4 bg-base-200 rounded-lg">
        <div class="flex items-center gap-2 mb-2">
          <span class="badge badge-{issue.severity === 'error' ? 'error' : issue.severity === 'warning' ? 'warning' : 'info'}">
            {issue.severity}
          </span>
          <span class="badge badge-outline">
            {issue.type?.replace(/_/g, ' ') || 'Unknown'}
          </span>
        </div>
        <h4 class="font-semibold text-lg mb-1">{issue.message || 'Issue'}</h4>
        {#if issue.url}
          <p class="text-sm text-base-content/70 break-all">{issue.url}</p>
        {/if}
      </div>

      <!-- Generate Button -->
      {#if !insight && !loading}
        <div class="mb-4">
          <button
            class="btn btn-primary w-full"
            on:click={handleGenerateInsight}
          >
            <Sparkles class="w-4 h-4" />
            Generate AI Insight
          </button>
        </div>
      {/if}

      <!-- Loading State -->
      {#if loading}
        <div class="flex flex-col items-center justify-center py-8">
          <span class="loading loading-spinner loading-lg mb-4"></span>
          <p class="text-base-content/70">Generating AI insight...</p>
        </div>
      {/if}

      <!-- Error State -->
      {#if error}
        <div class="alert alert-error mb-4">
          <span>{error}</span>
        </div>
        <button
          class="btn btn-outline w-full"
          on:click={handleGenerateInsight}
        >
          Try Again
        </button>
      {/if}

      <!-- Recommendation Display -->
      {#if insight && parsedResponse.recommendation}
        <div class="mb-4">
          <div class="flex items-center justify-between mb-2">
            <h4 class="font-semibold text-lg">Recommended Solution</h4>
            <button
              class="btn btn-sm btn-ghost"
              on:click={handleCopyRecommendation}
              title="Copy recommendation to clipboard"
            >
              <Copy class="w-4 h-4" />
              {copied ? 'Copied!' : 'Copy'}
            </button>
          </div>
          <div class="bg-primary/10 border border-primary/20 p-4 rounded-lg">
            <div class="prose prose-sm max-w-none">
              {@html renderedRecommendation}
            </div>
          </div>
        </div>
      {/if}

      <!-- Insight Display -->
      {#if insight && parsedResponse.insight}
        <div class="mb-4">
          <div class="flex items-center justify-between mb-2">
            <h4 class="font-semibold text-lg">AI Insight</h4>
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
            {@html renderedInsight}
          </div>
        </div>
      {/if}

      <!-- Fallback: Display full insight if no structured format -->
      {#if insight && !parsedResponse.recommendation && !parsedResponse.insight}
        <div class="mb-4">
          <div class="flex items-center justify-between mb-2">
            <h4 class="font-semibold text-lg">AI Insight</h4>
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
            {@html marked.parse(insight)}
          </div>
        </div>
      {/if}

      <!-- Modal Actions -->
      <div class="modal-action">
        <button class="btn" on:click={handleClose}>Close</button>
      </div>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button type="button" on:click={handleClose}>close</button>
    </form>
  </div>
{/if}

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
    @apply text-3xl font-bold mb-4 mt-6;
  }
  .prose h2 {
    @apply text-2xl font-semibold mb-3 mt-5;
  }
  .prose h3 {
    @apply text-xl font-semibold mb-2 mt-4;
  }
  .prose h4 {
    @apply text-lg font-semibold mb-2 mt-3;
  }
  .prose h5 {
    @apply text-base font-semibold mb-2 mt-2;
  }
  .prose h6 {
    @apply text-sm font-semibold mb-1 mt-2;
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
    @apply bg-base-300 px-2 py-1 rounded text-sm;
  }
  .prose pre {
    @apply bg-base-300 p-4 rounded-lg overflow-x-auto my-4;
  }
  .prose pre code {
    background-color: transparent;
    padding: 0;
  }
  .prose strong {
    @apply font-bold;
  }
  .prose em {
    @apply italic;
  }
  .prose blockquote {
    @apply border-l-4 border-base-300 pl-4 my-4 italic;
  }
</style>



