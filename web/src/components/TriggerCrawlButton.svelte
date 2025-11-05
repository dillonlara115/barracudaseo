<script>
  import { createEventDispatcher } from 'svelte';
  import { triggerCrawl } from '../lib/data.js';
  
  export let projectId = null;
  export let className = '';
  
  const dispatch = createEventDispatcher();
  
  let showModal = false;
  let loading = false;
  let error = null;
  
  // Form fields
  let url = '';
  let maxDepth = 3;
  let maxPages = 1000;
  let workers = 10;
  let respectRobots = true;
  let parseSitemap = false;

  async function handleSubmit() {
    if (!url) {
      error = 'URL is required';
      return;
    }

    // Basic URL validation
    try {
      new URL(url);
    } catch (e) {
      error = 'Invalid URL format';
      return;
    }

    loading = true;
    error = null;

    const { data, error: crawlError } = await triggerCrawl(projectId, {
      url,
      max_depth: maxDepth,
      max_pages: maxPages,
      workers,
      respect_robots: respectRobots,
      parse_sitemap: parseSitemap
    });

    loading = false;

    if (crawlError) {
      error = crawlError.message || 'Failed to trigger crawl';
      return;
    }

    // Success - close modal and dispatch event
    showModal = false;
    dispatch('created', data);
    
    // Reset form
    url = '';
    maxDepth = 3;
    maxPages = 1000;
    workers = 10;
    respectRobots = true;
    parseSitemap = false;
  }

  function handleCancel() {
    showModal = false;
    error = null;
    // Reset form
    url = '';
    maxDepth = 3;
    maxPages = 1000;
    workers = 10;
    respectRobots = true;
    parseSitemap = false;
  }
</script>

<button 
  class="btn btn-primary {className}"
  on:click={() => showModal = true}
>
  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
  </svg>
  Start Crawl
</button>

{#if showModal}
  <dialog class="modal modal-open">
    <div class="modal-box bg-base-200 text-base-content">
      <h3 class="font-bold text-lg mb-4">Start New Crawl</h3>
      
      {#if error}
        <div class="alert alert-error mb-4">
          <span>{error}</span>
        </div>
      {/if}

      <form on:submit|preventDefault={handleSubmit}>
        <div class="form-control mb-4">
          <label class="label">
            <span class="label-text text-base-content font-semibold">Starting URL *</span>
          </label>
          <input 
            type="url" 
            class="input input-bordered bg-base-100 text-base-content" 
            placeholder="https://example.com"
            bind:value={url}
            required
            disabled={loading}
          />
        </div>

        <div class="grid grid-cols-2 gap-4 mb-4">
          <div class="form-control">
            <label class="label">
              <span class="label-text text-base-content">Max Depth</span>
            </label>
            <input 
              type="number" 
              class="input input-bordered bg-base-100 text-base-content" 
              min="1"
              max="10"
              bind:value={maxDepth}
              disabled={loading}
            />
          </div>

          <div class="form-control">
            <label class="label">
              <span class="label-text text-base-content">Max Pages</span>
            </label>
            <input 
              type="number" 
              class="input input-bordered bg-base-100 text-base-content" 
              min="1"
              max="10000"
              bind:value={maxPages}
              disabled={loading}
            />
          </div>
        </div>

        <div class="form-control mb-4">
          <label class="label">
            <span class="label-text text-base-content">Workers</span>
          </label>
          <input 
            type="number" 
            class="input input-bordered bg-base-100 text-base-content" 
            min="1"
            max="50"
            bind:value={workers}
            disabled={loading}
          />
        </div>

        <div class="form-control mb-4">
          <label class="label cursor-pointer">
            <span class="label-text text-base-content">Respect robots.txt</span>
            <input 
              type="checkbox" 
              class="toggle toggle-primary" 
              bind:checked={respectRobots}
              disabled={loading}
            />
          </label>
        </div>

        <div class="form-control mb-4">
          <label class="label cursor-pointer">
            <span class="label-text text-base-content">Parse sitemap.xml</span>
            <input 
              type="checkbox" 
              class="toggle toggle-primary" 
              bind:checked={parseSitemap}
              disabled={loading}
            />
          </label>
        </div>

        <div class="modal-action">
          <button 
            type="button" 
            class="btn btn-ghost" 
            on:click={handleCancel}
            disabled={loading}
          >
            Cancel
          </button>
          <button 
            type="submit" 
            class="btn btn-primary"
            disabled={loading || !url}
          >
            {#if loading}
              <span class="loading loading-spinner loading-sm"></span>
              Starting...
            {:else}
              Start Crawl
            {/if}
          </button>
        </div>
      </form>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button on:click={handleCancel}>close</button>
    </form>
  </dialog>
{/if}

