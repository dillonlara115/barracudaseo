<script>
  import { onMount } from 'svelte';
  import { params, link } from 'svelte-spa-router';
  import { fetchProjects, fetchProjectGSCStatus, fetchProjectGSCDimensions } from '../lib/data.js';

  let projectId = null;
  let project = null;
  let loading = true;
  let error = null;
  let gscStatus = null;
  
  // Page rows with top queries
  let pageRows = [];
  let filteredPages = [];
  let searchQuery = '';
  let sortBy = 'impressions'; // 'impressions', 'clicks', 'query_count'
  let sortOrder = 'desc';
  let expandedPages = new Set();

  $: projectId = $params?.projectId || null;

  onMount(() => {
    if (projectId) {
      loadData();
    }
  });

  $: if (projectId) {
    loadData();
  }

  async function loadData() {
    if (!projectId) return;
    
    loading = true;
    error = null;

    try {
      // Load project
      const { data: projects } = await fetchProjects();
      if (projects) {
        project = projects.find(p => p.id === projectId);
      }

      // Load GSC status
      const statusResult = await fetchProjectGSCStatus(projectId);
      if (statusResult.error) {
        error = statusResult.error.message || 'Failed to load GSC status';
        loading = false;
        return;
      }
      gscStatus = statusResult.data;

      // Load page rows with top queries
      const pageResult = await fetchProjectGSCDimensions(projectId, 'page', { limit: 1000 });
      if (pageResult.error) {
        error = pageResult.error.message || 'Failed to load page data';
        loading = false;
        return;
      }

      pageRows = (pageResult.data?.rows || []).filter(row => {
        // Only include pages that have top_queries data
        const queries = row.top_queries || [];
        return Array.isArray(queries) && queries.length > 0;
      });

      // Sort by impressions descending by default
      pageRows.sort((a, b) => {
        const aMetrics = a.metrics || {};
        const bMetrics = b.metrics || {};
        return (bMetrics.impressions || 0) - (aMetrics.impressions || 0);
      });

      applyFilters();
      loading = false;
    } catch (err) {
      error = err.message || 'Failed to load GSC data';
      loading = false;
    }
  }

  function applyFilters() {
    filteredPages = [...pageRows];

    // Apply search filter
    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      filteredPages = filteredPages.filter(page => {
        const url = (page.dimension_value || '').toLowerCase();
        const queries = page.top_queries || [];
        const matchingQueries = queries.filter(q => {
          const queryText = (q.query || '').toLowerCase();
          return queryText.includes(query);
        });
        return url.includes(query) || matchingQueries.length > 0;
      });
    }

    // Apply sorting
    filteredPages.sort((a, b) => {
      let aValue, bValue;
      
      if (sortBy === 'query_count') {
        aValue = (a.top_queries || []).length;
        bValue = (b.top_queries || []).length;
      } else {
        const aMetrics = a.metrics || {};
        const bMetrics = b.metrics || {};
        aValue = aMetrics[sortBy] || 0;
        bValue = bMetrics[sortBy] || 0;
      }

      if (sortOrder === 'asc') {
        return aValue - bValue;
      } else {
        return bValue - aValue;
      }
    });
  }

  function togglePage(pageUrl) {
    if (expandedPages.has(pageUrl)) {
      expandedPages.delete(pageUrl);
    } else {
      expandedPages.add(pageUrl);
    }
    expandedPages = expandedPages; // Trigger reactivity
  }

  function formatNumber(num) {
    if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M';
    if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
    return Math.round(num).toLocaleString();
  }

  function formatPercent(num) {
    const percent = num > 1 ? num : num * 100;
    return percent.toFixed(2) + '%';
  }

  function getQueryCount(page) {
    return (page.top_queries || []).length;
  }

  function getTotalQueryImpressions(page) {
    const queries = page.top_queries || [];
    return queries.reduce((sum, q) => sum + (q.impressions || 0), 0);
  }

  function getTotalQueryClicks(page) {
    const queries = page.top_queries || [];
    return queries.reduce((sum, q) => sum + (q.clicks || 0), 0);
  }

  function getTopQuery(page) {
    const queries = page.top_queries || [];
    if (queries.length === 0) return null;
    // Sort by impressions descending
    const sorted = [...queries].sort((a, b) => (b.impressions || 0) - (a.impressions || 0));
    return sorted[0];
  }

  function getLowCTRQueries(page, threshold = 0.02) {
    const queries = page.top_queries || [];
    return queries.filter(q => {
      const ctr = q.ctr || 0;
      const normalizedCTR = ctr > 1 ? ctr / 100 : ctr;
      return normalizedCTR < threshold && (q.impressions || 0) >= 100;
    });
  }

  function getHighPositionQueries(page, threshold = 10) {
    const queries = page.top_queries || [];
    return queries.filter(q => (q.position || 0) > threshold && (q.impressions || 0) >= 100);
  }

  $: applyFilters();

  $: hasData = gscStatus?.integration?.property_url && pageRows.length > 0;
</script>

<svelte:head>
  <title>Keywords Per Page - {project?.name || 'Barracuda SEO'}</title>
</svelte:head>

<div class="container mx-auto p-6">
  <!-- Header -->
  <div class="mb-6">
    <div class="flex items-center justify-between mb-4">
      <div>
        <h1 class="text-3xl font-bold mb-2">Keywords Per Page</h1>
        {#if project}
          <p class="text-base-content/70">Project: {project.name}</p>
        {/if}
        {#if gscStatus?.integration?.property_url}
          <p class="text-sm text-base-content/60">Property: {gscStatus.integration.property_url}</p>
        {/if}
      </div>
      <div class="flex gap-2">
        <a href="/project/{projectId}/gsc" use:link class="btn btn-ghost">
          ← Back to GSC Dashboard
        </a>
        <a href="/project/{projectId}" use:link class="btn btn-ghost">
          ← Back to Project
        </a>
      </div>
    </div>

    <!-- Info Card -->
    <div class="alert alert-info mb-4">
      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
      </svg>
      <div>
        <div class="font-semibold mb-1">Keyword-Level Opportunities</div>
        <div class="text-sm">
          View which keywords each page ranks for. Identify underperforming keywords, irrelevant queries, and optimization opportunities.
        </div>
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
  {:else if !gscStatus?.integration?.property_url}
    <div class="alert alert-warning">
      <span>Google Search Console is not connected for this project. Please connect it in the Settings page.</span>
      <a href="/project/{projectId}/settings" use:link class="btn btn-sm btn-primary mt-2">Go to Settings</a>
    </div>
  {:else if !hasData}
    <div class="alert alert-info">
      <span>No keyword data available yet. Refresh the GSC data to populate keyword metrics.</span>
    </div>
  {:else}
    <!-- Filters and Search -->
    <div class="card bg-base-100 shadow mb-6">
      <div class="card-body">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div class="form-control">
            <label class="label">
              <span class="label-text">Search Pages or Keywords</span>
            </label>
            <input
              type="text"
              placeholder="Search..."
              class="input input-bordered w-full"
              bind:value={searchQuery}
            />
          </div>
          <div class="form-control">
            <label class="label">
              <span class="label-text">Sort By</span>
            </label>
            <select class="select select-bordered w-full" bind:value={sortBy}>
              <option value="impressions">Page Impressions</option>
              <option value="clicks">Page Clicks</option>
              <option value="query_count">Keyword Count</option>
            </select>
          </div>
          <div class="form-control">
            <label class="label">
              <span class="label-text">Order</span>
            </label>
            <select class="select select-bordered w-full" bind:value={sortOrder}>
              <option value="desc">Descending</option>
              <option value="asc">Ascending</option>
            </select>
          </div>
        </div>
        <div class="text-sm text-base-content/60 mt-2">
          Showing {filteredPages.length} of {pageRows.length} pages with keyword data
        </div>
      </div>
    </div>

    <!-- Pages with Keywords -->
    <div class="space-y-4">
      {#each filteredPages as page}
        {@const metrics = page.metrics || {}}
        {@const queries = page.top_queries || []}
        {@const isExpanded = expandedPages.has(page.dimension_value)}
        {@const topQuery = getTopQuery(page)}
        {@const lowCTRQueries = getLowCTRQueries(page)}
        {@const highPositionQueries = getHighPositionQueries(page)}
        
        <div class="card bg-base-100 shadow">
          <div class="card-body">
            <!-- Page Header -->
            <div class="flex items-start justify-between mb-4">
              <div class="flex-1">
                <div class="flex items-center gap-3 mb-2">
                  <h3 class="card-title text-lg">
                    <a 
                      href={page.dimension_value} 
                      target="_blank" 
                      rel="noopener noreferrer" 
                      class="link link-primary"
                    >
                      {page.dimension_value}
                    </a>
                  </h3>
                  <span class="badge badge-primary badge-sm">
                    {queries.length} {queries.length === 1 ? 'keyword' : 'keywords'}
                  </span>
                </div>
                
                <!-- Page Metrics -->
                <div class="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm mb-3">
                  <div>
                    <span class="text-base-content/60">Page Impressions:</span>
                    <span class="font-semibold ml-1">{formatNumber(metrics.impressions || 0)}</span>
                  </div>
                  <div>
                    <span class="text-base-content/60">Page Clicks:</span>
                    <span class="font-semibold ml-1">{formatNumber(metrics.clicks || 0)}</span>
                  </div>
                  <div>
                    <span class="text-base-content/60">Page CTR:</span>
                    <span class="font-semibold ml-1">{formatPercent(metrics.ctr || 0)}</span>
                  </div>
                  <div>
                    <span class="text-base-content/60">Avg Position:</span>
                    <span class="font-semibold ml-1">{(metrics.position || 0).toFixed(1)}</span>
                  </div>
                </div>

                <!-- Insights -->
                {#if topQuery}
                  <div class="text-sm text-base-content/70 mb-2">
                    <strong>Top Keyword:</strong> "{topQuery.query}" 
                    ({formatNumber(topQuery.impressions)} impressions, 
                    {formatNumber(topQuery.clicks)} clicks, 
                    {formatPercent(topQuery.ctr)} CTR, 
                    position {topQuery.position.toFixed(1)})
                  </div>
                {/if}

                {#if lowCTRQueries.length > 0}
                  <div class="alert alert-warning py-2 px-3 mb-2">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-5 h-5">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path>
                    </svg>
                    <span class="text-sm">
                      {lowCTRQueries.length} {lowCTRQueries.length === 1 ? 'keyword' : 'keywords'} with low CTR (&lt;2%) but high impressions (&gt;100). Consider optimizing titles/descriptions.
                    </span>
                  </div>
                {/if}

                {#if highPositionQueries.length > 0}
                  <div class="alert alert-info py-2 px-3 mb-2">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-5 h-5">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                    </svg>
                    <span class="text-sm">
                      {highPositionQueries.length} {highPositionQueries.length === 1 ? 'keyword' : 'keywords'} ranking below position 10. Opportunity to improve rankings.
                    </span>
                  </div>
                {/if}
              </div>
              
              <button
                class="btn btn-ghost btn-sm"
                on:click={() => togglePage(page.dimension_value)}
              >
                {#if isExpanded}
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 15.75l7.5-7.5 7.5 7.5" />
                  </svg>
                  Hide Keywords
                {:else}
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5" />
                  </svg>
                  Show Keywords ({queries.length})
                {/if}
              </button>
            </div>

            <!-- Expanded Keywords Table -->
            {#if isExpanded}
              <div class="overflow-x-auto mt-4">
                <table class="table table-zebra table-sm">
                  <thead>
                    <tr>
                      <th>Keyword</th>
                      <th>Impressions</th>
                      <th>Clicks</th>
                      <th>CTR</th>
                      <th>Position</th>
                      <th>Opportunity</th>
                    </tr>
                  </thead>
                  <tbody>
                    {#each queries.sort((a, b) => (b.impressions || 0) - (a.impressions || 0)) as query}
                      {@const ctr = query.ctr || 0}
                      {@const normalizedCTR = ctr > 1 ? ctr / 100 : ctr}
                      {@const isLowCTR = normalizedCTR < 0.02 && (query.impressions || 0) >= 100}
                      {@const isHighPosition = (query.position || 0) > 10 && (query.impressions || 0) >= 100}
                      
                      <tr>
                        <td class="font-medium">{query.query}</td>
                        <td>{formatNumber(query.impressions || 0)}</td>
                        <td>{formatNumber(query.clicks || 0)}</td>
                        <td>
                          <span class="{isLowCTR ? 'text-warning font-semibold' : ''}">
                            {formatPercent(query.ctr || 0)}
                          </span>
                        </td>
                        <td>
                          <span class="{isHighPosition ? 'text-info font-semibold' : ''}">
                            {(query.position || 0).toFixed(1)}
                          </span>
                        </td>
                        <td>
                          {#if isLowCTR}
                            <span class="badge badge-warning badge-sm">Low CTR</span>
                          {/if}
                          {#if isHighPosition}
                            <span class="badge badge-info badge-sm">Low Position</span>
                          {/if}
                          {#if !isLowCTR && !isHighPosition && (query.impressions || 0) >= 100}
                            <span class="badge badge-success badge-sm">Good</span>
                          {/if}
                        </td>
                      </tr>
                    {/each}
                  </tbody>
                </table>
              </div>
            {/if}
          </div>
        </div>
      {/each}
    </div>

    {#if filteredPages.length === 0}
      <div class="alert alert-info">
        <span>No pages match your search criteria.</span>
      </div>
    {/if}
  {/if}
</div>

