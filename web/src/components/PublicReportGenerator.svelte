<script>
  import { createPublicReport, listPublicReports, deletePublicReport } from '../lib/data.js';
  import { Copy, ExternalLink, Trash2, Lock, Calendar, Eye } from 'lucide-svelte';
  import { userProfile, isProOrTeam } from '../lib/subscription.js';

  $: isPro = isProOrTeam($userProfile);

  export let crawlId = null;
  export let projectId = null;

  let loading = false;
  let error = null;
  let reports = [];
  let showCreateModal = false;
  let formData = {
    title: '',
    description: '',
    password: '',
    expiresInDays: null
  };

  // Load existing reports
  async function loadReports() {
    if (!projectId) return;
    
    loading = true;
    error = null;
    const { data, error: err } = await listPublicReports(projectId);
    if (err) {
      error = err.message || 'Failed to load reports';
      loading = false;
      return;
    }
    reports = data?.reports || [];
    loading = false;
  }

  // Create new public report
  async function handleCreateReport() {
    if (!crawlId) {
      error = 'Crawl ID is required';
      return;
    }

    loading = true;
    error = null;
    
    const { data, error: err } = await createPublicReport(crawlId, {
      title: formData.title,
      description: formData.description,
      password: formData.password || null,
      expiresInDays: formData.expiresInDays || null
    });

    if (err) {
      error = err.message || 'Failed to create report';
      loading = false;
      return;
    }

    // Reset form and reload reports
    formData = {
      title: '',
      description: '',
      password: '',
      expiresInDays: null
    };
    showCreateModal = false;
    await loadReports();
    loading = false;
  }

  // Delete a report
  async function handleDeleteReport(reportId) {
    if (!confirm('Are you sure you want to delete this report? It will no longer be accessible.')) {
      return;
    }

    loading = true;
    error = null;
    
    const { error: err } = await deletePublicReport(reportId);
    if (err) {
      error = err.message || 'Failed to delete report';
      loading = false;
      return;
    }

    await loadReports();
    loading = false;
  }

  // Copy public URL to clipboard
  async function copyUrl(url) {
    try {
      await navigator.clipboard.writeText(url);
      // Show toast or feedback
      alert('Report URL copied to clipboard!');
    } catch (err) {
      console.error('Failed to copy URL:', err);
    }
  }

  // Load reports on mount
  import { onMount } from 'svelte';
  onMount(() => {
    loadReports();
  });
</script>

<div class="card bg-base-100 shadow-xl">
  <div class="card-body">
    <div class="flex items-center justify-between mb-4">
      <h2 class="card-title">Public Client Reports</h2>
      {#if isPro}
        <button
          class="btn btn-primary btn-sm"
          on:click={() => showCreateModal = true}
          disabled={loading || !crawlId}
        >
          <ExternalLink class="w-4 h-4 mr-2" />
          Create Public Report
        </button>
      {:else}
        <div class="tooltip" data-tip="Upgrade to Pro to create public reports">
          <button class="btn btn-primary btn-sm btn-disabled" disabled>
            <ExternalLink class="w-4 h-4 mr-2" />
            Create Public Report
            <span class="badge badge-primary badge-sm ml-1">PRO</span>
          </button>
        </div>
      {/if}
    </div>

    {#if error}
      <div class="alert alert-error mb-4">
        <span>{error}</span>
      </div>
    {/if}

    {#if loading && reports.length === 0}
      <div class="flex justify-center py-8">
        <span class="loading loading-spinner loading-lg"></span>
      </div>
    {:else if reports.length === 0}
      <div class="text-center py-8 text-base-content/60">
        <p>No public reports created yet.</p>
        <p class="text-sm mt-2">Create a shareable report for your clients.</p>
      </div>
    {:else}
      <div class="space-y-4">
        {#each reports as report}
          <div class="border border-base-300 rounded-lg p-4">
            <div class="flex items-start justify-between">
              <div class="flex-1">
                <h3 class="font-semibold text-lg">{report.title || 'Untitled Report'}</h3>
                {#if report.description}
                  <p class="text-sm text-base-content/70 mt-1">{report.description}</p>
                {/if}
                <div class="flex items-center gap-4 mt-3 text-sm text-base-content/60">
                  {#if report.password_hash}
                    <span class="flex items-center gap-1">
                      <Lock class="w-4 h-4" />
                      Password Protected
                    </span>
                  {/if}
                  {#if report.expires_at}
                    <span class="flex items-center gap-1">
                      <Calendar class="w-4 h-4" />
                      Expires: {new Date(report.expires_at).toLocaleDateString()}
                    </span>
                  {/if}
                  <span class="flex items-center gap-1">
                    <Eye class="w-4 h-4" />
                    {report.view_count || 0} views
                  </span>
                  <span class="text-xs">
                    Created: {new Date(report.created_at).toLocaleDateString()}
                  </span>
                </div>
              </div>
              <div class="flex items-center gap-2 ml-4">
                <button
                  class="btn btn-sm btn-ghost"
                  on:click={() => copyUrl(report.public_url)}
                  title="Copy URL"
                >
                  <Copy class="w-4 h-4" />
                </button>
                <a
                  href={report.public_url}
                  target="_blank"
                  rel="noopener noreferrer"
                  class="btn btn-sm btn-primary"
                  on:click|preventDefault={() => {
                    // Use hash-based routing
                    window.location.hash = report.public_url.split('#')[1] || report.public_url;
                    window.open(report.public_url, '_blank');
                  }}
                >
                  <ExternalLink class="w-4 h-4 mr-1" />
                  View
                </a>
                <button
                  class="btn btn-sm btn-error btn-ghost"
                  on:click={() => handleDeleteReport(report.id)}
                  title="Delete report"
                >
                  <Trash2 class="w-4 h-4" />
                </button>
              </div>
            </div>
            {#if report.public_url}
              <div class="mt-3 p-2 bg-base-200 rounded text-xs font-mono break-all">
                {report.public_url}
              </div>
            {/if}
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>

<!-- Create Report Modal -->
{#if showCreateModal}
  <div class="modal modal-open">
    <div class="modal-box">
      <h3 class="font-bold text-lg mb-4">Create Public Report</h3>
      
      <div class="form-control mb-4">
        <label class="label" for="report-title">
          <span class="label-text">Report Title</span>
        </label>
        <input
          id="report-title"
          type="text"
          class="input input-bordered"
          bind:value={formData.title}
          placeholder="e.g., SEO Audit Report - January 2024"
        />
      </div>

      <div class="form-control mb-4">
        <label class="label" for="report-description">
          <span class="label-text">Description (Optional)</span>
        </label>
        <textarea
          id="report-description"
          class="textarea textarea-bordered"
          bind:value={formData.description}
          placeholder="Add a description for this report..."
          rows="3"
        ></textarea>
      </div>

      <div class="form-control mb-4">
        <label class="label" for="report-password">
          <span class="label-text">Password Protection (Optional)</span>
        </label>
        <input
          id="report-password"
          type="password"
          class="input input-bordered"
          bind:value={formData.password}
          placeholder="Leave empty for no password"
        />
        <div class="label">
          <span class="label-text-alt">Clients will need this password to view the report</span>
        </div>
      </div>

      <div class="form-control mb-4">
        <label class="label" for="report-expires">
          <span class="label-text">Expires In (Days, Optional)</span>
        </label>
        <input
          id="report-expires"
          type="number"
          class="input input-bordered"
          bind:value={formData.expiresInDays}
          placeholder="e.g., 30 (leave empty for no expiry)"
          min="1"
        />
        <div class="label">
          <span class="label-text-alt">Report will automatically expire after this many days</span>
        </div>
      </div>

      <div class="modal-action">
        <button
          class="btn btn-ghost"
          on:click={() => {
            showCreateModal = false;
            formData = { title: '', description: '', password: '', expiresInDays: null };
          }}
          disabled={loading}
        >
          Cancel
        </button>
        <button
          class="btn btn-primary"
          on:click={handleCreateReport}
          disabled={loading}
        >
          {#if loading}
            <span class="loading loading-spinner loading-sm"></span>
          {:else}
            Create Report
          {/if}
        </button>
      </div>
    </div>
    <div 
      class="modal-backdrop" 
      role="button"
      tabindex="0"
      on:click={() => showCreateModal = false}
      on:keydown={(e) => e.key === 'Enter' || e.key === ' ' ? showCreateModal = false : null}
      aria-label="Close report generator"
    ></div>
  </div>
{/if}

