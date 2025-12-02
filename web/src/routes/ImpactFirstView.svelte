<script>
  import { onMount } from 'svelte';
  import { params, link } from 'svelte-spa-router';
  import { fetchProjects, fetchProjectImpactFirst, fetchProjectGSCStatus } from '../lib/data.js';
  import { AlertTriangle, TrendingUp, ExternalLink } from 'lucide-svelte';
  import ProjectPageLayout from '../components/ProjectPageLayout.svelte';

  let projectId = null;
  let project = null;
  let impactPages = [];
  let loading = true;
  let error = null;
  let gscStatus = null;

  $: projectId = $params?.projectId || null;

  onMount(() => {
    if (projectId) {
      loadData();
    }
  });

  $: if (projectId) {
    loadData();
  }

  async function loadGSCStatus() {
    if (!projectId) return;
    try {
      const result = await fetchProjectGSCStatus(projectId);
      if (!result.error && result.data) {
        gscStatus = result.data;
      }
    } catch (err) {
      console.error('Failed to load GSC status:', err);
    }
  }

  async function loadData() {
    if (!projectId) return;
    
    loading = true;
    error = null;
    await loadGSCStatus();

    try {
      // Load project
      const { data: projects } = await fetchProjects();
      if (projects) {
        project = projects.find(p => p.id === projectId);
      }

      // Load impact-first data
      const result = await fetchProjectImpactFirst(projectId);
      if (result.error) {
        error = result.error.message || 'Failed to load impact data';
        loading = false;
        return;
      }
      
      impactPages = result.data?.pages || [];
      loading = false;
    } catch (err) {
      error = err.message || 'Failed to load data';
      loading = false;
    }
  }

  function getSeverityColor(severity) {
    switch (severity?.toLowerCase()) {
      case 'critical':
        return 'badge-error';
      case 'high':
        return 'badge-warning';
      case 'medium':
        return 'badge-info';
      case 'low':
        return 'badge-success';
      default:
        return 'badge-ghost';
    }
  }

  function formatImpactScore(score) {
    return Math.round(score);
  }
</script>

<svelte:head>
  <title>Impact-First View - {project?.name || 'Barracuda SEO'}</title>
</svelte:head>

<ProjectPageLayout {projectId} {gscStatus} showCrawlSection={false}>
<div class="max-w-7xl mx-auto">
  <!-- Header -->
  <div class="mb-6">
    <div class="flex items-center justify-between mb-4">
      <div>
        <h1 class="text-3xl font-bold mb-2">Impact-First View</h1>
        <p class="text-base-content/70 mb-1">
          Pages that rank for keywords AND have crawl issues, prioritized by impact score. 
          Focus on fixing issues on these pages first for maximum SEO impact.
        </p>
        {#if project}
          <p class="text-sm text-base-content/60">Project: {project.name}</p>
        {/if}
      </div>
      <div class="flex gap-2">
        <a href="/project/{projectId}/rank-tracker" use:link class="btn btn-ghost">
          ← Back to Rank Tracker
        </a>
        <a href="/project/{projectId}" use:link class="btn btn-ghost">
          ← Back to Project
        </a>
      </div>
    </div>
  </div>

  {#if loading}
    <div class="flex justify-center items-center py-20">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
  {:else if error}
    <div class="alert alert-error">
      <span>{error}</span>
    </div>
  {:else if impactPages.length === 0}
    <div class="alert alert-info">
      <span>No pages found that both rank for keywords and have crawl issues. Run a crawl and check some keywords to see impact data.</span>
    </div>
  {:else}
    <!-- Impact Pages Table -->
    <div class="card bg-base-100 shadow">
      <div class="card-body">
        <div class="overflow-x-auto">
          <table class="table table-zebra">
            <thead>
              <tr>
                <th>Impact Score</th>
                <th>URL</th>
                <th>Best Position</th>
                <th>Keywords</th>
                <th>Issues</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {#each impactPages as page}
                <tr>
                  <td>
                    <div class="flex items-center gap-2">
                      <TrendingUp class="w-5 h-5 text-warning" />
                      <span class="font-bold text-lg">{formatImpactScore(page.impact_score)}</span>
                    </div>
                  </td>
                  <td class="max-w-md">
                    <a href={page.url} target="_blank" rel="noopener noreferrer" class="link link-primary truncate block">
                      {page.url}
                    </a>
                  </td>
                  <td>
                    {#if page.best_position}
                      <span class="badge badge-lg">{page.best_position}</span>
                    {:else}
                      <span class="text-base-content/40">—</span>
                    {/if}
                  </td>
                  <td>
                    <div class="flex flex-wrap gap-1">
                      {#each page.keywords.slice(0, 3) as keyword}
                        <span class="badge badge-sm">{keyword}</span>
                      {/each}
                      {#if page.keywords.length > 3}
                        <span class="badge badge-sm badge-ghost">+{page.keywords.length - 3}</span>
                      {/if}
                    </div>
                    <div class="text-xs text-base-content/60 mt-1">
                      {page.keyword_count} keyword{page.keyword_count !== 1 ? 's' : ''}
                    </div>
                  </td>
                  <td>
                    <div class="flex flex-wrap gap-1 mb-2">
                      {#each page.issues.slice(0, 5) as issue}
                        <span class="badge badge-sm {getSeverityColor(issue.severity)}">
                          {issue.severity || 'unknown'}
                        </span>
                      {/each}
                      {#if page.issues.length > 5}
                        <span class="badge badge-sm badge-ghost">+{page.issues.length - 5}</span>
                      {/if}
                    </div>
                    <div class="text-xs text-base-content/60">
                      {page.issue_count} issue{page.issue_count !== 1 ? 's' : ''}
                    </div>
                  </td>
                  <td>
                    <a
                      href="/project/{projectId}/crawl/{page.crawl_id || ''}?page={page.page_id}"
                      use:link
                      class="btn btn-xs btn-outline"
                    >
                      View Page
                      <ExternalLink class="w-3 h-3 ml-1" />
                    </a>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Info Card -->
    <div class="alert alert-info mt-6">
      <AlertTriangle class="w-5 h-5" />
      <div>
        <h3 class="font-bold">Impact Score Calculation</h3>
        <p class="text-sm">
          Impact Score = (100 - Best Position) × Issue Count
        </p>
        <p class="text-sm mt-1">
          Pages with lower rankings (better positions) and more issues have higher impact scores.
          Focus on fixing issues on these pages first for maximum SEO impact.
        </p>
      </div>
    </div>
  {/if}
</div>
</ProjectPageLayout>

