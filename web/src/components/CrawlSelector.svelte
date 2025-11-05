<script>
  import { createEventDispatcher } from 'svelte';
  import { push } from 'svelte-spa-router';
  
  export let crawls = [];
  export let selectedCrawl = null;
  
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
  <label class="label">
    <span class="label-text text-base-content font-semibold">Select Crawl:</span>
  </label>
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

