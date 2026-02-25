<script>
  import { onMount } from 'svelte';
  import { fetchInternalLinkSuggestions, fetchOrphanedPages } from '../lib/data.js';

  export let projectId;
  export let pageUrl = '';

  let activeTab = 'suggestions';
  let suggestions = [];
  let orphaned = [];
  let loading = false;
  let error = null;

  onMount(loadData);

  async function loadData() {
    loading = true;
    error = null;

    const [sugRes, orphRes] = await Promise.all([
      fetchInternalLinkSuggestions(projectId, pageUrl),
      fetchOrphanedPages(projectId)
    ]);

    loading = false;

    if (sugRes.error) { error = sugRes.error.message; return; }
    suggestions = sugRes.data || [];

    if (!orphRes.error) {
      orphaned = orphRes.data || [];
    }
  }
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <h3 class="text-lg font-bold">Internal Link Analysis</h3>
    <button class="btn btn-ghost btn-sm" on:click={loadData} disabled={loading}>
      {#if loading}<span class="loading loading-spinner loading-sm"></span>{/if}
      Refresh
    </button>
  </div>

  {#if error}
    <div class="alert alert-error text-sm">{error}</div>
  {/if}

  <div class="tabs tabs-boxed">
    <button class="tab" class:tab-active={activeTab === 'suggestions'} on:click={() => activeTab = 'suggestions'}>
      Link Suggestions ({suggestions.length})
    </button>
    <button class="tab" class:tab-active={activeTab === 'orphaned'} on:click={() => activeTab = 'orphaned'}>
      Orphaned Pages ({orphaned.length})
    </button>
  </div>

  {#if loading}
    <div class="flex justify-center py-12">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
  {:else if activeTab === 'suggestions'}
    <div class="space-y-2">
      {#each suggestions as page}
        <div class="bg-base-200 rounded-xl px-4 py-3">
          <div class="flex items-center justify-between gap-3">
            <div class="min-w-0 flex-1">
              <h4 class="font-semibold text-sm truncate">{page.title || 'Untitled'}</h4>
              <p class="text-xs text-base-content/40 truncate">{page.url}</p>
            </div>
            <div class="flex items-center gap-2 shrink-0">
              <span class="bg-base-300 text-base-content/60 text-xs px-2.5 py-1 rounded-lg font-mono">
                {page.inbound_links ?? 0} inbound
              </span>
            </div>
          </div>
        </div>
      {:else}
        <div class="text-center py-8 text-base-content/40">No crawl data available. Run a crawl from the dashboard first.</div>
      {/each}
    </div>
  {:else if activeTab === 'orphaned'}
    <div class="space-y-2">
      {#each orphaned as page}
        <div class="bg-base-200 rounded-xl px-4 py-3 border-l-4 border-warning">
          <div class="flex items-center justify-between gap-3">
            <div class="min-w-0 flex-1">
              <h4 class="font-semibold text-sm truncate">{page.title || 'Untitled'}</h4>
              <p class="text-xs text-base-content/40 truncate">{page.url}</p>
            </div>
            <span class="text-xs font-medium text-warning whitespace-nowrap">No inbound links</span>
          </div>
        </div>
      {:else}
        <div class="text-center py-8 text-base-content/40">No orphaned pages detected.</div>
      {/each}
    </div>
  {/if}
</div>
