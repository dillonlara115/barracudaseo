<script>
  import { createEventDispatcher } from 'svelte';
  import { generateIssueInsight } from '../../lib/data.js';
  import { Sparkles, Copy, X } from 'lucide-svelte';

  export let issue = null;
  export let crawlId = null;

  const dispatch = createEventDispatcher();

  let loading = false;
  let insight = null;
  let error = null;
  let copied = false;

  async function handleGenerateInsight() {
    if (!issue || !crawlId) return;

    loading = true;
    error = null;
    insight = null;

    try {
      const { data, error: apiError } = await generateIssueInsight(issue.id, crawlId);
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

      <!-- Insight Display -->
      {#if insight}
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
          <div class="prose prose-sm max-w-none bg-base-200 p-4 rounded-lg whitespace-pre-wrap">
            {insight}
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
    margin-top: 0.5em;
    margin-bottom: 0.5em;
  }
</style>

