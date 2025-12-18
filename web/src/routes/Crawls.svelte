<script>
  import { onMount } from 'svelte';
  import { params, push, link } from 'svelte-spa-router';
  import { fetchProjects, fetchCrawls, fetchCrawl, deleteCrawl } from '../lib/data.js';
  import ProjectPageLayout from '../components/ProjectPageLayout.svelte';
  import TriggerCrawlButton from '../components/TriggerCrawlButton.svelte';
  import { FileSearch, Trash2, Eye, Loader, AlertTriangle } from 'lucide-svelte';

  let projectId = null;
  let project = null;
  let crawls = [];
  let loading = true;
  let error = null;
  let deletingCrawlId = null;
  let showDeleteConfirm = null;
  let showActiveCrawlNotification = false;
  let activeCrawl = null;

  $: projectId = $params?.projectId || null;

  onMount(async () => {
    if (projectId) {
      await loadData();
      await checkActiveCrawl();
    }
  });

  $: if (projectId) {
    loadData();
    checkActiveCrawl();
  }

  async function checkActiveCrawl() {
    if (!projectId) return;
    
    // Check localStorage for active crawl
    const activeCrawlId = localStorage.getItem(`activeCrawl_${projectId}`);
    if (activeCrawlId) {
      try {
        const { data: crawl, error: crawlError } = await fetchCrawl(activeCrawlId);
        
        if (!crawlError && crawl) {
          if (crawl.status === 'running' || crawl.status === 'pending') {
            activeCrawl = crawl;
            showActiveCrawlNotification = true;
          } else {
            // Clean up if crawl is done
            localStorage.removeItem(`activeCrawl_${projectId}`);
          }
        } else {
          // Crawl not found, clean up
          localStorage.removeItem(`activeCrawl_${projectId}`);
        }
      } catch (err) {
        console.error('Error checking active crawl:', err);
        localStorage.removeItem(`activeCrawl_${projectId}`);
      }
    }
  }

  async function loadData() {
    if (!projectId) return;
    
    loading = true;
    error = null;
    try {
      // Load projects
      const { data: projectsData, error: projectsError } = await fetchProjects();
      if (projectsError) throw projectsError;
      
      project = projectsData?.find(p => p.id === projectId);
      if (!project) {
        error = 'Project not found';
        loading = false;
        return;
      }

      // Load crawls
      const { data: crawlsData, error: crawlsError } = await fetchCrawls(projectId);
      if (crawlsError) throw crawlsError;
      crawls = crawlsData || [];
    } catch (err) {
      error = err.message || 'Failed to load data';
    } finally {
      loading = false;
    }
  }

  function formatDate(dateString) {
    if (!dateString) return 'Unknown';
    const date = new Date(dateString);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  }

  function formatRelativeDate(dateString) {
    if (!dateString) return 'Unknown';
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now - date;
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMins < 1) return 'Just now';
    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffHours < 24) return `${diffHours}h ago`;
    if (diffDays < 7) return `${diffDays}d ago`;
    return formatDate(dateString);
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

  function viewCrawl(crawl, event) {
    if (event) {
      event.stopPropagation();
    }
    push(`/project/${projectId}/crawl/${crawl.id}`);
  }

  function handleDeleteClick(crawl, event) {
    event.stopPropagation();
    showDeleteConfirm = crawl.id;
  }
  
  function cancelDelete() {
    showDeleteConfirm = null;
  }
  
  async function confirmDelete() {
    const crawl = crawls.find(c => c.id === showDeleteConfirm);
    if (!crawl) return;
    
    if (crawl.status === 'running') {
      alert('Cannot delete a crawl that is currently running');
      showDeleteConfirm = null;
      return;
    }
    
    deletingCrawlId = crawl.id;
    const { error: deleteError } = await deleteCrawl(crawl.id);
    
    if (deleteError) {
      alert(`Failed to delete crawl: ${deleteError.message}`);
      deletingCrawlId = null;
      showDeleteConfirm = null;
      return;
    }
    
    // Reload crawls list
    await loadData();
    deletingCrawlId = null;
    showDeleteConfirm = null;
  }

  async function handleCrawlCreated(e) {
    // Reload crawls
    await loadData();
    // Navigate to the new crawl
    if (e.detail?.crawl_id) {
      push(`/project/${projectId}/crawl/${e.detail.crawl_id}`);
    }
  }

  function handleViewActiveCrawl() {
    if (activeCrawl?.id) {
      showActiveCrawlNotification = false;
      push(`/project/${projectId}/crawl/${activeCrawl.id}`);
    }
  }
  
  function handleDismissActiveCrawl() {
    showActiveCrawlNotification = false;
  }
</script>

<ProjectPageLayout {projectId} showCrawlSection={false}>
  <div class="max-w-7xl mx-auto">
    <!-- Header -->
    <div class="mb-6">
      <div class="flex items-center justify-between mb-4">
        <div>
          <h1 class="text-3xl font-bold mb-2">Crawls & Audits</h1>
          <p class="text-base-content/70 mb-1">
            Manage all website crawls and audits for this project. View crawl history, start new crawls, 
            and access detailed reports for each crawl.
          </p>
          {#if project}
            <p class="text-sm text-base-content/60">Project: {project.name}</p>
          {/if}
        </div>
        {#if project}
          <TriggerCrawlButton {projectId} project={project} on:created={handleCrawlCreated} />
        {/if}
      </div>
    </div>

    <!-- Active Crawl Notification -->
    {#if showActiveCrawlNotification && activeCrawl}
      <div class="alert alert-info mb-4">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
        </svg>
        <div class="flex-1">
          <h3 class="font-bold">Crawl in Progress</h3>
          <div class="text-sm">A crawl is currently running for this project.</div>
        </div>
        <div class="flex gap-2">
          <button class="btn btn-sm btn-primary" on:click={handleViewActiveCrawl}>
            View Progress
          </button>
          <button class="btn btn-sm btn-ghost" on:click={handleDismissActiveCrawl}>
            Dismiss
          </button>
        </div>
      </div>
    {/if}

    {#if loading}
      <div class="flex items-center justify-center min-h-[400px]">
        <span class="loading loading-spinner loading-lg"></span>
      </div>
    {:else if error}
      <div class="alert alert-error">
        <span>Error: {error}</span>
        <button class="btn btn-sm btn-ghost mt-2" on:click={loadData}>Retry</button>
      </div>
    {:else if crawls.length === 0}
      <div class="card bg-base-200">
        <div class="card-body text-center py-12">
          <FileSearch class="w-16 h-16 mx-auto mb-4 text-base-content/30" />
          <h3 class="text-xl font-bold mb-2">No Crawls Yet</h3>
          <p class="text-base-content/70 mb-4">
            Start your first crawl to analyze your website's SEO health.
          </p>
          {#if project}
            <TriggerCrawlButton {projectId} project={project} on:created={handleCrawlCreated} />
          {/if}
        </div>
      </div>
    {:else}
      <!-- Crawls Table -->
      <div class="card bg-base-100 shadow">
        <div class="card-body p-0">
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
                  {@const isDeleting = deletingCrawlId === crawl.id}
                  <tr class="hover cursor-pointer" on:click={() => viewCrawl(crawl)}>
                    <td>
                      <div class="font-medium">{formatDate(crawl.started_at)}</div>
                      <div class="text-sm text-base-content/60">{formatRelativeDate(crawl.started_at)}</div>
                    </td>
                    <td>
                      <span class="badge {getStatusBadge(crawl.status)} badge-sm capitalize">
                        {crawl.status}
                      </span>
                    </td>
                    <td>{crawl.total_pages || 0}</td>
                    <td>
                      {#if crawl.total_issues > 0}
                        <span class="text-error font-semibold">{crawl.total_issues}</span>
                      {:else}
                        <span class="text-base-content/60">0</span>
                      {/if}
                    </td>
                    <td>
                      <div class="flex items-center gap-2">
                        <button
                          class="btn btn-ghost btn-xs"
                          on:click={(e) => viewCrawl(crawl, e)}
                          title="View crawl"
                        >
                          <Eye class="w-4 h-4" />
                        </button>
                        <button
                          class="btn btn-ghost btn-xs text-error hover:bg-error/20"
                          on:click={(e) => handleDeleteClick(crawl, e)}
                          disabled={isDeleting || crawl.status === 'running'}
                          title={crawl.status === 'running' ? 'Cannot delete running crawl' : 'Delete crawl'}
                        >
                          {#if isDeleting}
                            <Loader class="w-4 h-4 animate-spin" />
                          {:else}
                            <Trash2 class="w-4 h-4" />
                          {/if}
                        </button>
                      </div>
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    {/if}
  </div>
</ProjectPageLayout>

<!-- Delete Confirmation Modal -->
{#if showDeleteConfirm}
  <div class="modal modal-open">
    <div class="modal-box">
      <h3 class="font-bold text-lg mb-4 flex items-center gap-2">
        <AlertTriangle class="w-6 h-6 text-error" />
        Delete Crawl?
      </h3>
      <p class="py-4">
        Are you sure you want to delete this crawl? This will permanently delete the crawl and all associated pages and issues.
        This action cannot be undone.
      </p>
      <div class="modal-action">
        <button class="btn btn-ghost" on:click={cancelDelete} disabled={deletingCrawlId !== null}>
          Cancel
        </button>
        <button class="btn btn-error" on:click={confirmDelete} disabled={deletingCrawlId !== null}>
          {#if deletingCrawlId !== null}
            <Loader class="w-4 h-4 animate-spin" />
            Deleting...
          {:else}
            <Trash2 class="w-4 h-4" />
            Delete
          {/if}
        </button>
      </div>
    </div>
    <div 
      class="modal-backdrop" 
      role="button"
      tabindex="0"
      on:click={cancelDelete}
      on:keydown={(e) => e.key === 'Enter' || e.key === ' ' ? cancelDelete() : null}
      aria-label="Close delete confirmation"
    ></div>
  </div>
{/if}

