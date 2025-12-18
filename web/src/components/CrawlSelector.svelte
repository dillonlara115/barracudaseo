<script>
  import { createEventDispatcher } from 'svelte';
  import { link } from 'svelte-spa-router';
  
  export let crawls = [];
  export let selectedCrawl = null;
  export let projectId = null;
  
  const dispatch = createEventDispatcher();

  function formatDate(dateString) {
    if (!dateString) return 'Unknown';
    const date = new Date(dateString);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  }

  function handleSelect(crawl) {
    dispatch('select', crawl);
  }
</script>

<div class="mb-4">
  <div class="flex items-center justify-between mb-2">
    <div class="label py-0">
      <span class="label-text text-base-content font-semibold">Select Crawl:</span>
    </div>
    {#if projectId && crawls.length > 0}
      <a 
        href="/project/{projectId}/crawls" 
        use:link
        class="btn btn-ghost btn-xs text-base-content/70 hover:text-base-content"
        title="View all crawls"
      >
        View All
        <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3 ml-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </a>
    {/if}
  </div>
  <select 
    class="select select-bordered w-full max-w-xs bg-base-200 text-base-content"
    value={selectedCrawl?.id}
    on:change={(e) => {
      const crawl = crawls.find(c => c.id === e.target.value);
      if (crawl) handleSelect(crawl);
    }}
  >
    {#each crawls as crawl}
      <option value={crawl.id}>
        {formatDate(crawl.started_at)} - {crawl.total_pages || 0} pages - {crawl.total_issues || 0} issues
      </option>
    {/each}
  </select>
</div>

