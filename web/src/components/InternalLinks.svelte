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
    <div class="space-y-3">
      {#each suggestions as page}
        <div class="card bg-base-200">
          <div class="card-body py-4">
            <div class="flex items-center justify-between">
              <div>
                <h4 class="font-semibold text-sm">{page.title || 'Untitled'}</h4>
                <p class="text-xs opacity-60 truncate max-w-md">{page.url}</p>
              </div>
              <span class="badge badge-ghost badge-sm">{page.word_count || 0} words</span>
            </div>
          </div>
        </div>
      {:else}
        <div class="text-center py-8 opacity-60">No link suggestions available. Run a site crawl first.</div>
      {/each}
    </div>
  {:else if activeTab === 'orphaned'}
    <div class="space-y-3">
      {#each orphaned as page}
        <div class="card bg-base-200 border-l-4 border-warning">
          <div class="card-body py-4">
            <div class="flex items-center justify-between">
              <div>
                <h4 class="font-semibold text-sm">{page.title || 'Untitled'}</h4>
                <p class="text-xs opacity-60 truncate max-w-md">{page.url}</p>
              </div>
              <span class="badge badge-warning badge-sm">No inbound links</span>
            </div>
          </div>
        </div>
      {:else}
        <div class="text-center py-8 opacity-60">No orphaned pages detected.</div>
      {/each}
    </div>
  {/if}
</div>
