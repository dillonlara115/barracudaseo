<script>
  import { onMount } from 'svelte';
  import { params, link } from 'svelte-spa-router';
  import { fetchProjects, fetchProjectGSCStatus, fetchProjectGSCDimensions, triggerProjectGSCSync } from '../lib/data.js';
  import { Bar, Line } from 'svelte-chartjs';
  import { Chart, CategoryScale, LinearScale, BarElement, LineElement, PointElement, Title, Tooltip, Legend } from 'chart.js';
  import ProjectPageLayout from '../components/ProjectPageLayout.svelte';

  // Register Chart.js components
  Chart.register(CategoryScale, LinearScale, BarElement, LineElement, PointElement, Title, Tooltip, Legend);

  let projectId = null;
  let project = null;
  let loading = true;
  let error = null;
  let gscStatus = null;
  let gscRefreshing = false;
  
  // Data for different dimensions
  let dateRows = [];
  let pageRows = [];
  let queryRows = [];
  let countryRows = [];
  let deviceRows = [];
  let appearanceRows = [];
  
  // Chart data
  let dateChartData = null;
  let deviceChartData = null;
  let countryChartData = null;
  let appearanceChartData = null;

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

      // Load all dimensions
      await Promise.all([
        loadDimension('date', (data) => { dateRows = data; prepareDateChart(); }),
        loadDimension('page', (data) => { pageRows = data; }),
        loadDimension('query', (data) => { queryRows = data; }),
        loadDimension('country', (data) => { countryRows = data; prepareCountryChart(); }),
        loadDimension('device', (data) => { deviceRows = data; prepareDeviceChart(); }),
        loadDimension('appearance', (data) => { appearanceRows = data; prepareAppearanceChart(); }),
      ]);

      loading = false;
    } catch (err) {
      error = err.message || 'Failed to load GSC data';
      loading = false;
    }
  }

  async function loadDimension(type, callback) {
    try {
      const result = await fetchProjectGSCDimensions(projectId, type, { limit: 1000 });
      if (!result.error && result.data?.rows) {
        let processedRows = result.data.rows;
        
        // Deduplicate and aggregate rows by dimension value
        if (type === 'page') {
          processedRows = deduplicatePages(result.data.rows);
        } else if (type === 'query') {
          processedRows = deduplicateQueries(result.data.rows);
        } else if (type === 'date') {
          processedRows = deduplicateDimension(result.data.rows);
        } else if (type === 'device') {
          processedRows = deduplicateDimension(result.data.rows);
        } else if (type === 'country') {
          processedRows = deduplicateDimension(result.data.rows);
        } else if (type === 'appearance') {
          processedRows = deduplicateDimension(result.data.rows);
        }
        
        callback(processedRows);
      }
    } catch (err) {
      console.error(`Failed to load ${type} dimension:`, err);
    }
  }

  function deduplicatePages(rows) {
    // Group pages by URL and aggregate metrics
    const pageMap = new Map();
    
    for (const row of rows) {
      const url = row.dimension_value;
      if (!url) continue;
      
      const metrics = row.metrics || {};
      const existing = pageMap.get(url);
      
      if (existing) {
        // Aggregate metrics: sum clicks/impressions, weighted average for CTR/position
        const existingMetrics = existing.metrics || {};
        const existingClicks = existingMetrics.clicks || 0;
        const existingImpressions = existingMetrics.impressions || 0;
        const newClicks = metrics.clicks || 0;
        const newImpressions = metrics.impressions || 0;
        
        // Sum clicks and impressions
        existingMetrics.clicks = existingClicks + newClicks;
        existingMetrics.impressions = existingImpressions + newImpressions;
        
        // Weighted average for CTR (clicks / impressions)
        if (existingMetrics.impressions > 0) {
          existingMetrics.ctr = existingMetrics.clicks / existingMetrics.impressions;
        }
        
        // Weighted average for position
        const totalImpressions = existingImpressions + newImpressions;
        if (totalImpressions > 0) {
          const existingPosition = existingMetrics.position || 0;
          const newPosition = metrics.position || 0;
          existingMetrics.position = (
            (existingPosition * existingImpressions) + (newPosition * newImpressions)
          ) / totalImpressions;
        }
        
        // Merge top_queries if available (keep unique queries)
        if (row.top_queries && Array.isArray(row.top_queries) && row.top_queries.length > 0) {
          const existingQueries = existing.top_queries || [];
          const queryMap = new Map();
          
          // Add existing queries
          existingQueries.forEach(q => {
            if (q.query) queryMap.set(q.query, q);
          });
          
          // Add new queries
          row.top_queries.forEach(q => {
            if (q.query) queryMap.set(q.query, q);
          });
          
          existing.top_queries = Array.from(queryMap.values());
        }
      } else {
        // First occurrence of this URL
        pageMap.set(url, {
          ...row,
          metrics: { ...metrics }
        });
      }
    }
    
    // Convert map back to array and sort by impressions descending
    const deduplicated = Array.from(pageMap.values());
    deduplicated.sort((a, b) => {
      const aImpressions = (a.metrics || {}).impressions || 0;
      const bImpressions = (b.metrics || {}).impressions || 0;
      return bImpressions - aImpressions;
    });
    
    return deduplicated;
  }

  function deduplicateQueries(rows) {
    // Group queries by query text and aggregate metrics
    const queryMap = new Map();
    
    for (const row of rows) {
      const query = row.dimension_value;
      if (!query) continue;
      
      const metrics = row.metrics || {};
      const existing = queryMap.get(query);
      
      if (existing) {
        // Aggregate metrics: sum clicks/impressions, weighted average for CTR/position
        const existingMetrics = existing.metrics || {};
        const existingClicks = existingMetrics.clicks || 0;
        const existingImpressions = existingMetrics.impressions || 0;
        const newClicks = metrics.clicks || 0;
        const newImpressions = metrics.impressions || 0;
        
        // Sum clicks and impressions
        existingMetrics.clicks = existingClicks + newClicks;
        existingMetrics.impressions = existingImpressions + newImpressions;
        
        // Weighted average for CTR (clicks / impressions)
        if (existingMetrics.impressions > 0) {
          existingMetrics.ctr = existingMetrics.clicks / existingMetrics.impressions;
        }
        
        // Weighted average for position
        const totalImpressions = existingImpressions + newImpressions;
        if (totalImpressions > 0) {
          const existingPosition = existingMetrics.position || 0;
          const newPosition = metrics.position || 0;
          existingMetrics.position = (
            (existingPosition * existingImpressions) + (newPosition * newImpressions)
          ) / totalImpressions;
        }
      } else {
        // First occurrence of this query
        queryMap.set(query, {
          ...row,
          metrics: { ...metrics }
        });
      }
    }
    
    // Convert map back to array and sort by impressions descending
    const deduplicated = Array.from(queryMap.values());
    deduplicated.sort((a, b) => {
      const aImpressions = (a.metrics || {}).impressions || 0;
      const bImpressions = (b.metrics || {}).impressions || 0;
      return bImpressions - aImpressions;
    });
    
    return deduplicated;
  }

  function deduplicateDimension(rows) {
    // Generic deduplication function for device, country, appearance, etc.
    // Groups by dimension_value and aggregates metrics
    const dimensionMap = new Map();
    
    for (const row of rows) {
      const value = row.dimension_value;
      if (!value) continue;
      
      const metrics = row.metrics || {};
      const existing = dimensionMap.get(value);
      
      if (existing) {
        // Aggregate metrics: sum clicks/impressions, weighted average for CTR/position
        const existingMetrics = existing.metrics || {};
        const existingClicks = existingMetrics.clicks || 0;
        const existingImpressions = existingMetrics.impressions || 0;
        const newClicks = metrics.clicks || 0;
        const newImpressions = metrics.impressions || 0;
        
        // Sum clicks and impressions
        existingMetrics.clicks = existingClicks + newClicks;
        existingMetrics.impressions = existingImpressions + newImpressions;
        
        // Weighted average for CTR (clicks / impressions)
        if (existingMetrics.impressions > 0) {
          existingMetrics.ctr = existingMetrics.clicks / existingMetrics.impressions;
        }
        
        // Weighted average for position
        const totalImpressions = existingImpressions + newImpressions;
        if (totalImpressions > 0) {
          const existingPosition = existingMetrics.position || 0;
          const newPosition = metrics.position || 0;
          existingMetrics.position = (
            (existingPosition * existingImpressions) + (newPosition * newImpressions)
          ) / totalImpressions;
        }
      } else {
        // First occurrence of this dimension value
        dimensionMap.set(value, {
          ...row,
          metrics: { ...metrics }
        });
      }
    }
    
    // Convert map back to array and sort by impressions descending
    const deduplicated = Array.from(dimensionMap.values());
    deduplicated.sort((a, b) => {
      const aImpressions = (a.metrics || {}).impressions || 0;
      const bImpressions = (b.metrics || {}).impressions || 0;
      return bImpressions - aImpressions;
    });
    
    return deduplicated;
  }

  function prepareDateChart() {
    if (!dateRows.length) return;
    
    // Sort by date
    const sorted = [...dateRows].sort((a, b) => {
      const dateA = a.dimension_value || '';
      const dateB = b.dimension_value || '';
      return dateA.localeCompare(dateB);
    });

    const labels = sorted.map(r => {
      const date = r.dimension_value;
      if (!date) return '';
      // Format date as MM/DD
      const parts = date.split('-');
      if (parts.length === 3) {
        return `${parts[1]}/${parts[2]}`;
      }
      return date;
    });

    const clicks = sorted.map(r => {
      const m = r.metrics || {};
      return Math.round(m.clicks || 0);
    });
    const impressions = sorted.map(r => {
      const m = r.metrics || {};
      return Math.round(m.impressions || 0);
    });
    const ctr = sorted.map(r => {
      const m = r.metrics || {};
      return (m.ctr || 0) * 100; // Convert to percentage
    });
    const position = sorted.map(r => {
      const m = r.metrics || {};
      return m.position || 0;
    });

    dateChartData = {
      labels,
      datasets: [
        {
          label: 'Clicks',
          data: clicks,
          borderColor: 'rgb(59, 130, 246)',
          backgroundColor: 'rgba(59, 130, 246, 0.1)',
          yAxisID: 'y',
        },
        {
          label: 'Impressions',
          data: impressions,
          borderColor: 'rgb(16, 185, 129)',
          backgroundColor: 'rgba(16, 185, 129, 0.1)',
          yAxisID: 'y',
        },
        {
          label: 'CTR (%)',
          data: ctr,
          borderColor: 'rgb(245, 158, 11)',
          backgroundColor: 'rgba(245, 158, 11, 0.1)',
          yAxisID: 'y1',
        },
        {
          label: 'Position',
          data: position,
          borderColor: 'rgb(239, 68, 68)',
          backgroundColor: 'rgba(239, 68, 68, 0.1)',
          yAxisID: 'y2',
        },
      ],
    };
  }

  function prepareDeviceChart() {
    if (!deviceRows.length) return;
    
    const top10 = deviceRows.slice(0, 10);
    deviceChartData = {
      labels: top10.map(r => r.dimension_value || 'Unknown'),
      datasets: [{
        label: 'Clicks',
        data: top10.map(r => Math.round((r.metrics || {}).clicks || 0)),
        backgroundColor: [
          'rgba(59, 130, 246, 0.8)',
          'rgba(16, 185, 129, 0.8)',
          'rgba(245, 158, 11, 0.8)',
          'rgba(239, 68, 68, 0.8)',
          'rgba(139, 92, 246, 0.8)',
        ],
      }],
    };
  }

  function prepareCountryChart() {
    if (!countryRows.length) return;
    
    const top10 = countryRows.slice(0, 10);
    countryChartData = {
      labels: top10.map(r => r.dimension_value || 'Unknown'),
      datasets: [{
        label: 'Clicks',
        data: top10.map(r => Math.round((r.metrics || {}).clicks || 0)),
        backgroundColor: 'rgba(59, 130, 246, 0.8)',
      }],
    };
  }

  function prepareAppearanceChart() {
    if (!appearanceRows.length) return;
    
    appearanceChartData = {
      labels: appearanceRows.map(r => r.dimension_value || 'Unknown'),
      datasets: [{
        label: 'Clicks',
        data: appearanceRows.map(r => Math.round((r.metrics || {}).clicks || 0)),
        backgroundColor: [
          'rgba(59, 130, 246, 0.8)',
          'rgba(16, 185, 129, 0.8)',
          'rgba(245, 158, 11, 0.8)',
          'rgba(239, 68, 68, 0.8)',
        ],
      }],
    };
  }

  async function refreshData() {
    if (!projectId || gscRefreshing) return;
    gscRefreshing = true;
    error = null;
    
    const syncResult = await triggerProjectGSCSync(projectId, { lookback_days: 30 });
    if (syncResult.error) {
      error = syncResult.error.message || 'Failed to refresh data';
      gscRefreshing = false;
      return;
    }
    
    // Wait a moment for sync to complete, then reload
    setTimeout(() => {
      loadData();
      gscRefreshing = false;
    }, 2000);
  }

  function formatNumber(num) {
    if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M';
    if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
    return num.toLocaleString();
  }

  function formatPercent(num) {
    return (num * 100).toFixed(2) + '%';
  }

  // Calculate totals from aggregated date dimension data to ensure accuracy
  $: totals = (() => {
    if (!dateRows.length) {
      // Fallback to snapshot totals if no date data available yet
      return gscStatus?.summary?.totals || {};
    }
    
    // Sum up clicks and impressions from all dates
    let totalClicks = 0;
    let totalImpressions = 0;
    let weightedPositionSum = 0;
    
    for (const row of dateRows) {
      const metrics = row.metrics || {};
      const clicks = metrics.clicks || 0;
      const impressions = metrics.impressions || 0;
      const position = metrics.position || 0;
      
      totalClicks += clicks;
      totalImpressions += impressions;
      weightedPositionSum += position * impressions;
    }
    
    // Calculate CTR from totals
    const ctr = totalImpressions > 0 ? totalClicks / totalImpressions : 0;
    
    // Calculate weighted average position
    const avgPosition = totalImpressions > 0 ? weightedPositionSum / totalImpressions : 0;
    
    return {
      clicks: totalClicks,
      impressions: totalImpressions,
      ctr: ctr,
      position: avgPosition
    };
  })();
  
  $: hasData = gscStatus?.integration?.property_url && (dateRows.length > 0 || pageRows.length > 0);
</script>

<svelte:head>
  <title>GSC Dashboard - {project?.name || 'Barracuda SEO'}</title>
</svelte:head>

<ProjectPageLayout {projectId} {gscStatus} showCrawlSection={false}>
<div class="max-w-7xl mx-auto">
  <!-- Header -->
  <div class="mb-6">
    <div class="flex items-center justify-between mb-4">
      <div>
        <h1 class="text-3xl font-bold mb-2">Google Search Console Dashboard</h1>
        <p class="text-base-content/70 mb-1">
          View performance metrics from Google Search Console including clicks, impressions, CTR, and average position. 
          Analyze trends by date, device, country, and page.
        </p>
        {#if project}
          <p class="text-sm text-base-content/60">Project: {project.name}</p>
        {/if}
        {#if gscStatus?.integration?.property_url}
          <p class="text-sm text-base-content/60">Property: {gscStatus.integration.property_url}</p>
        {/if}
      </div>
      <div class="flex gap-2">
        <a href="/project/{projectId}" use:link class="btn btn-ghost">
          ← Back to Project
        </a>
        {#if hasData}
          <a
            href="/project/{projectId}/gsc/keywords"
            use:link
            class="btn btn-outline"
          >
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4 mr-1">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
            </svg>
            Keywords Per Page
          </a>
        {/if}
        <button 
          class="btn btn-primary"
          on:click={refreshData}
          disabled={gscRefreshing || loading}
        >
          {#if gscRefreshing}
            <span class="loading loading-spinner loading-sm"></span>
            Refreshing...
          {:else}
            Refresh Data
          {/if}
        </button>
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
      <span>Google Search Console is not connected for your account. Connect it in Integrations, then select a property in Project Settings.</span>
      <a href="/integrations" use:link class="btn btn-sm btn-primary mt-2">Go to Integrations</a>
    </div>
  {:else if !hasData}
    <div class="alert alert-info">
      <span>No GSC data available yet. Click "Refresh Data" to sync data from Google Search Console.</span>
    </div>
  {:else}
    <!-- Overview Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <div class="card bg-base-100 shadow">
        <div class="card-body">
          <h2 class="card-title text-sm text-base-content/70">Total Clicks</h2>
          <p class="text-3xl font-bold">{formatNumber(totals.clicks || 0)}</p>
        </div>
      </div>
      <div class="card bg-base-100 shadow">
        <div class="card-body">
          <h2 class="card-title text-sm text-base-content/70">Total Impressions</h2>
          <p class="text-3xl font-bold">{formatNumber(totals.impressions || 0)}</p>
        </div>
      </div>
      <div class="card bg-base-100 shadow">
        <div class="card-body">
          <h2 class="card-title text-sm text-base-content/70">Average CTR</h2>
          <p class="text-3xl font-bold">{formatPercent(totals.ctr || 0)}</p>
        </div>
      </div>
      <div class="card bg-base-100 shadow">
        <div class="card-body">
          <h2 class="card-title text-sm text-base-content/70">Average Position</h2>
          <p class="text-3xl font-bold">{(totals.position || 0).toFixed(1)}</p>
        </div>
      </div>
    </div>

    <!-- Time Series Chart -->
    {#if dateChartData}
      <div class="card bg-base-100 shadow mb-6">
        <div class="card-body">
          <h2 class="card-title mb-4">Performance Over Time</h2>
          <div class="h-64">
            <Line
              data={dateChartData}
              options={{
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                  y: {
                    type: 'linear',
                    position: 'left',
                    title: { display: true, text: 'Clicks / Impressions' },
                  },
                  y1: {
                    type: 'linear',
                    position: 'right',
                    title: { display: true, text: 'CTR (%)' },
                    grid: { drawOnChartArea: false },
                  },
                  y2: {
                    type: 'linear',
                    position: 'right',
                    title: { display: true, text: 'Position' },
                    grid: { drawOnChartArea: false },
                    display: false,
                  },
                },
                plugins: {
                  legend: { position: 'top' },
                  tooltip: { mode: 'index', intersect: false },
                },
              }}
            />
          </div>
        </div>
      </div>
    {/if}

    <!-- Charts Row -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
      <!-- Device Breakdown -->
      {#if deviceChartData}
        <div class="card bg-base-100 shadow">
          <div class="card-body">
            <h2 class="card-title mb-4">Performance by Device</h2>
            <div class="h-64">
              <Bar
                data={deviceChartData}
                options={{
                  responsive: true,
                  maintainAspectRatio: false,
                  plugins: {
                    legend: { display: false },
                  },
                }}
              />
            </div>
          </div>
        </div>
      {/if}

      <!-- Country Breakdown -->
      {#if countryChartData}
        <div class="card bg-base-100 shadow">
          <div class="card-body">
            <h2 class="card-title mb-4">Top Countries</h2>
            <div class="h-64">
              <Bar
                data={countryChartData}
                options={{
                  responsive: true,
                  maintainAspectRatio: false,
                  plugins: {
                    legend: { display: false },
                  },
                }}
              />
            </div>
          </div>
        </div>
      {/if}
    </div>

    <!-- Search Appearance -->
    {#if appearanceChartData}
      <div class="card bg-base-100 shadow mb-6">
        <div class="card-body">
          <h2 class="card-title mb-4">Search Appearance</h2>
          <div class="h-64">
            <Bar
              data={appearanceChartData}
              options={{
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                  legend: { display: false },
                },
              }}
            />
          </div>
        </div>
      </div>
    {/if}

    <!-- Top Pages Table -->
    {#if pageRows.length > 0}
      <div class="card bg-base-100 shadow mb-6">
        <div class="card-body">
          <div class="flex items-center justify-between mb-4">
            <h2 class="card-title">Top Pages</h2>
            <a
              href="/project/{projectId}/gsc/keywords"
              use:link
              class="btn btn-primary btn-sm"
            >
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4 mr-1">
                <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
              </svg>
              View Keywords Per Page
            </a>
          </div>
          <div class="overflow-x-auto">
            <table class="table table-zebra">
              <thead>
                <tr>
                  <th>Page</th>
                  <th>Clicks</th>
                  <th>Impressions</th>
                  <th>CTR</th>
                  <th>Position</th>
                </tr>
              </thead>
              <tbody>
                {#each pageRows.slice(0, 20) as row}
                  {@const metrics = row.metrics || {}}
                  <tr>
                    <td>
                      <a href={row.dimension_value} target="_blank" rel="noopener noreferrer" class="link link-primary">
                        {row.dimension_value}
                      </a>
                    </td>
                    <td>{formatNumber(Math.round(metrics.clicks || 0))}</td>
                    <td>{formatNumber(Math.round(metrics.impressions || 0))}</td>
                    <td>{formatPercent(metrics.ctr || 0)}</td>
                    <td>{(metrics.position || 0).toFixed(1)}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    {/if}

    <!-- Top Queries Table -->
    {#if queryRows.length > 0}
      <div class="card bg-base-100 shadow mb-6">
        <div class="card-body">
          <h2 class="card-title mb-4">Top Queries</h2>
          <div class="overflow-x-auto">
            <table class="table table-zebra">
              <thead>
                <tr>
                  <th>Query</th>
                  <th>Clicks</th>
                  <th>Impressions</th>
                  <th>CTR</th>
                  <th>Position</th>
                </tr>
              </thead>
              <tbody>
                {#each queryRows.slice(0, 20) as row}
                  {@const metrics = row.metrics || {}}
                  <tr>
                    <td class="font-medium">{row.dimension_value}</td>
                    <td>{formatNumber(Math.round(metrics.clicks || 0))}</td>
                    <td>{formatNumber(Math.round(metrics.impressions || 0))}</td>
                    <td>{formatPercent(metrics.ctr || 0)}</td>
                    <td>{(metrics.position || 0).toFixed(1)}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    {/if}
  {/if}
</div>
</ProjectPageLayout>
