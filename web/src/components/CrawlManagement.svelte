<script>
  import { onMount } from 'svelte';
  import { createEventDispatcher } from 'svelte';
  import { push } from 'svelte-spa-router';
  import { fetchCrawls, deleteCrawl } from '../lib/data.js';
  
  export let projectId = null;
  
  const dispatch = createEventDispatcher();
  
  let crawls = [];
  let loading = true;
  let deletingCrawlId = null;
  let showDeleteConfirm = null;
  let error = null;

  let mounted = false;

  onMount(async () => {
    mounted = true;
    if (projectId) {
      await loadCrawls();
    } else {
      loading = false;
    }
  });

  // Reload when projectId changes (after mount)
  $: if (mounted && projectId) {
    loadCrawls();
  }

  async function loadCrawls() {
    if (!projectId) {
      console.warn('CrawlManagement: projectId is not set');
      loading = false;
      crawls = [];
      return;
    }
    
    console.log('CrawlManagement: Loading crawls for project:', projectId);
    loading = true;
    error = null;
    try {
      const { data, error: fetchError } = await fetchCrawls(projectId);
      if (fetchError) {
        error = fetchError.message || 'Failed to load crawls';
        console.error('Failed to load crawls:', fetchError);
        crawls = [];
      } else {
        console.log('CrawlManagement: Loaded', data?.length || 0, 'crawls');
        crawls = data || [];
      }
    } catch (err) {
      error = err.message || 'Failed to load crawls';
      console.error('Failed to load crawls:', err);
      crawls = [];
    } finally {
      loading = false;
    }
  }

  function formatDate(dateString) {
    if (!dateString) return 'Unknown';
    const date = new Date(dateString);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  }

  function handleDeleteClick(crawl, event) {
    event.stopPropagation();
    showDeleteConfirm = crawl.id;
  }
  
  function cancelDelete() {
    showDeleteConfirm = null;
  }
  
  async function confirmDelete(crawl) {
    if (crawl.status === 'running') {
      alert('Cannot delete a crawl that is currently running');
      showDeleteConfirm = null;
      return;
    }
    
    deletingCrawlId = crawl.id;
    const { error } = await deleteCrawl(crawl.id);
    
    if (error) {
      alert(`Failed to delete crawl: ${error.message}`);
      deletingCrawlId = null;
      showDeleteConfirm = null;
      return;
    }
    
    // Reload crawls list
    await loadCrawls();
    
    dispatch('deleted', crawl);
    deletingCrawlId = null;
    showDeleteConfirm = null;
  }
  
  function getStatusBadge(status) {
    const badges = {
      'pending': 'badge-warning',
      'running': 'badge-info',
      'succeeded': 'badge-success',
      'failed': 'badge-error',
      'cancelled': 'badge-ghost'
    };
    return badges[status] || 'badge-ghost';
  }

  function viewCrawl(crawl) {
    push(`/project/${projectId}/crawl/${crawl.id}`);
  }
</script>

<div id="crawls" class="card bg-base-100 shadow">
  <div class="card-body">
    <h2 class="card-title text-xl mb-4">Crawl Management</h2>
    <p class="text-base-content/70 mb-4">
      View and manage all crawls for this project. Deleting a crawl will permanently remove it and all associated pages and issues.
    </p>
    
    {#if !projectId}
      <div class="alert alert-warning">
        <span>Project ID not available.</span>
      </div>
    {:else if loading}
      <div class="flex justify-center py-8">
        <span class="loading loading-spinner loading-lg"></span>
      </div>
    {:else if error}
      <div class="alert alert-error">
        <span>Error: {error}</span>
        <button class="btn btn-sm btn-ghost mt-2" on:click={loadCrawls}>Retry</button>
      </div>
    {:else if crawls.length === 0}
      <div class="alert alert-info">
        <span>No crawls found for this project.</span>
      </div>
    {:else}
      <div class="overflow-x-auto">
        <table class="table table-zebra">
          <thead>
            <tr>
              <th>Date</th>
              <th>Status</th>
              <th>Pages</th>
              <th>Issues</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {#each crawls as crawl}
              <tr class="hover">
                <td>
                  <div class="font-medium">{formatDate(crawl.started_at)}</div>
                </td>
                <td>
                  <span class="badge {getStatusBadge(crawl.status)} badge-sm">{crawl.status}</span>
                </td>
                <td>{crawl.total_pages || 0}</td>
                <td>{crawl.total_issues || 0}</td>
                <td>
                  <div class="flex items-center gap-2">
                    <button
                      class="btn btn-ghost btn-xs"
                      on:click={() => viewCrawl(crawl)}
                      title="View crawl"
                    >
                      View
                    </button>
                    <button
                      class="btn btn-ghost btn-xs text-error hover:bg-error/20"
                      on:click={(e) => handleDeleteClick(crawl, e)}
                      disabled={deletingCrawlId === crawl.id || crawl.status === 'running'}
                      title={crawl.status === 'running' ? 'Cannot delete running crawl' : 'Delete crawl'}
                    >
                      {#if deletingCrawlId === crawl.id}
                        <span class="loading loading-spinner loading-xs"></span>
                      {:else}
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                        </svg>
                      {/if}
                    </button>
                  </div>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
  </div>
</div>

{#if showDeleteConfirm}
  <div class="modal modal-open">
    <div class="modal-box">
      <h3 class="font-bold text-lg mb-4">Delete Crawl?</h3>
      <p class="py-4">
        Are you sure you want to delete this crawl? This will permanently delete the crawl and all associated pages and issues.
        This action cannot be undone.
      </p>
      <div class="modal-action">
        <button class="btn btn-ghost" on:click={cancelDelete}>Cancel</button>
        <button class="btn btn-error" on:click={() => confirmDelete(crawls.find(c => c.id === showDeleteConfirm))}>
          Delete
        </button>
      </div>
    </div>
  </div>
{/if}

