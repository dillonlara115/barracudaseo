<script>
  import { onMount } from 'svelte';
  import { push, querystring, link, location } from 'svelte-spa-router';
  import SummaryCard from './SummaryCard.svelte';
  import ResultsTable from './ResultsTable.svelte';
  import IssuesPanel from './IssuesPanel.svelte';
  import LinkGraph from './LinkGraph.svelte';
  import RecommendationsPanel from './RecommendationsPanel.svelte';
  import GSCDashboardPanel from './GSCDashboardPanel.svelte';
  import GSCKeywordsPanel from './GSCKeywordsPanel.svelte';
  import GA4DashboardPanel from './GA4DashboardPanel.svelte';
  import ClarityDashboardPanel from './ClarityDashboardPanel.svelte';
  import UnifiedInsightsPanel from './UnifiedInsightsPanel.svelte';
  import CrawlSummary from './AI/CrawlSummary.svelte';
  import PublicReportGenerator from './PublicReportGenerator.svelte';
  import Logo from './Logo.svelte';
  import { fetchProjects, fetchProjectGSCStatus, fetchProjectGSCDimensions, triggerProjectGSCSync, fetchProjectGA4Status, triggerProjectGA4Sync, fetchProjectClarityStatus, triggerProjectClaritySync, fetchCrawls } from '../lib/data.js';
  import { buildEnrichedIssues } from '../lib/gsc.js';
  import { userProfile, isProOrTeam } from '../lib/subscription.js';
  import { 
    LayoutDashboard, 
    FileText, 
    AlertTriangle, 
    Lightbulb, 
    Network, 
    TrendingUp, 
    ScanSearch, 
    Target,
    Search,
    Settings,
    BarChart,
    FileSearch,
    ArrowRight
  } from 'lucide-svelte';

  export let summary = null;
  export let results = [];
  export let initialTab = 'dashboard';
  export let projectId = null;
  export let crawlId = null;
  export let project = null; // Accept project as prop from parent

  // Check if we're on the settings page
  $: isSettingsPage = $location?.includes('/settings');

  // Load project if not provided
  onMount(async () => {
    if (!project && projectId) {
      const { data: projects } = await fetchProjects();
      if (projects) {
        project = projects.find(p => p.id === projectId);
      }
    }
    // Load recent crawls for dashboard
    if (projectId) {
      await loadRecentCrawls();
    }
  });

  async function loadRecentCrawls() {
    if (!projectId) return;
    crawlsLoading = true;
    try {
      const { data, error } = await fetchCrawls(projectId);
      if (!error && data) {
        // Get most recent 5 crawls
        recentCrawls = data.slice(0, 5);
      }
    } catch (err) {
      console.error('Failed to load recent crawls:', err);
    } finally {
      crawlsLoading = false;
    }
  }

  function formatDate(dateString) {
    if (!dateString) return 'Unknown';
    const date = new Date(dateString);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
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

  $: activeTab = $querystring 
    ? new URLSearchParams($querystring).get('tab') || initialTab 
    : initialTab;
  let issuesFilter = { severity: 'all', type: 'all', url: null };
  let resultsFilter = { status: 'all', performance: false };
  let cachedEnrichedIssues = [];
  let activeEnrichedIssues = [];
  let enrichedIssuesMap = {};
  let gscStatus = null;
  let gscLoading = false;
  let gscRefreshing = false;
  let gscError = null;
  let gscPageRows = [];
  let gscInitializedProjectId = null;
  let ga4Status = null;
  let ga4Loading = false;
  let ga4Refreshing = false;
  let ga4Error = null;
  let ga4InitializedProjectId = null;
  let clarityStatus = null;
  let clarityLoading = false;
  let clarityRefreshing = false;
  let clarityError = null;
  let clarityInitializedProjectId = null;
  let recentCrawls = [];
  let crawlsLoading = false;

  const navigateToTab = (tab, nextFilters = {}) => {
    const { severity, type, url, status, performance } = nextFilters;

    // Update filters first (before navigation)
    if (severity !== undefined || type !== undefined || url !== undefined) {
      issuesFilter = {
        ...issuesFilter,
        ...(severity !== undefined ? { severity } : {}),
        ...(type !== undefined ? { type } : {}),
        ...(url !== undefined ? { url } : {})
      };
    }

    if (status !== undefined || performance !== undefined) {
      resultsFilter = {
        ...resultsFilter,
        ...(status !== undefined ? { status } : {}),
        ...(performance !== undefined ? { performance } : {})
      };
    }

    // Update URL with tab query param (after filter update)
    if (projectId && crawlId) {
      const params = new URLSearchParams();
      params.set('tab', tab);
      push(`/project/${projectId}/crawl/${crawlId}?${params.toString()}`);
    } else if (projectId && (tab === 'gsc-dashboard' || tab === 'gsc-keywords' || tab === 'ga4-dashboard' || tab === 'clarity-dashboard' || tab === 'insights')) {
      // GSC tabs work at project level, redirect to first crawl or project view
      const params = new URLSearchParams();
      params.set('tab', tab);
      // Try to keep in crawl context if available, otherwise go to project
      if (crawlId) {
        push(`/project/${projectId}/crawl/${crawlId}?${params.toString()}`);
      } else {
        push(`/project/${projectId}?${params.toString()}`);
      }
    } else if (projectId) {
      // Fallback: if we have projectId but no crawlId, still try to navigate
      // This handles edge cases where crawlId might not be set yet
      const params = new URLSearchParams();
      params.set('tab', tab);
      push(`/project/${projectId}?${params.toString()}`);
    }
  };

  // Callback for GSC to update enriched issues
  const formatDateTime = (value) => {
    if (!value) return '';
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return '';
    return `${date.toLocaleDateString()} ${date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}`;
  };

  async function loadGSCData(targetProjectId) {
    if (!targetProjectId) return;

    gscLoading = true;
    gscError = null;
    gscStatus = null;
    gscPageRows = [];

    const statusResult = await fetchProjectGSCStatus(targetProjectId);
    if (statusResult.error) {
      gscError = statusResult.error.message || 'Unable to load Google Search Console status.';
      gscLoading = false;
      return;
    }

    gscStatus = statusResult.data;

    if (gscStatus?.integration?.property_url) {
      const pageResult = await fetchProjectGSCDimensions(targetProjectId, 'page', { limit: 1000 });
      if (pageResult.error) {
        gscError = pageResult.error.message || 'Unable to load Search Console metrics.';
      } else {
        gscPageRows = pageResult.data?.rows || [];
      }
    }

    gscLoading = false;
  }

  async function refreshGSCData() {
    if (!projectId || gscRefreshing) return;
    gscRefreshing = true;
    gscError = null;
    gscLoading = true;
    const syncResult = await triggerProjectGSCSync(projectId, { lookback_days: 30 });
    if (syncResult.error) {
      gscError = syncResult.error.message || 'Failed to refresh Google Search Console data.';
      gscRefreshing = false;
      gscLoading = false;
      return;
    }
    await loadGSCData(projectId);
    gscRefreshing = false;
  }

  $: if (projectId && projectId !== gscInitializedProjectId) {
    gscInitializedProjectId = projectId;
    loadGSCData(projectId);
  }

  async function loadGA4Data(targetProjectId) {
    if (!targetProjectId) return;
    ga4Loading = true;
    ga4Error = null;
    ga4Status = null;

    const statusResult = await fetchProjectGA4Status(targetProjectId);
    if (statusResult.error) {
      ga4Error = statusResult.error.message || 'Unable to load Google Analytics 4 status.';
      ga4Loading = false;
      return;
    }
    ga4Status = statusResult.data;
    ga4Loading = false;
  }

  async function refreshGA4Data() {
    if (!projectId || ga4Refreshing) return;
    ga4Refreshing = true;
    ga4Error = null;
    ga4Loading = true;
    const syncResult = await triggerProjectGA4Sync(projectId, { lookback_days: 30 });
    if (syncResult.error) {
      ga4Error = syncResult.error.message || 'Failed to refresh Google Analytics 4 data.';
      ga4Refreshing = false;
      ga4Loading = false;
      return;
    }
    await loadGA4Data(projectId);
    ga4Refreshing = false;
  }

  $: if (projectId && projectId !== ga4InitializedProjectId) {
    ga4InitializedProjectId = projectId;
    loadGA4Data(projectId);
  }

  async function loadClarityData(targetProjectId) {
    if (!targetProjectId) return;
    clarityLoading = true;
    clarityError = null;
    clarityStatus = null;

    const statusResult = await fetchProjectClarityStatus(targetProjectId);
    if (statusResult.error) {
      clarityError = statusResult.error.message || 'Unable to load Microsoft Clarity status.';
      clarityLoading = false;
      return;
    }
    clarityStatus = statusResult.data;
    clarityLoading = false;
  }

  async function refreshClarityData() {
    if (!projectId || clarityRefreshing) return;
    clarityRefreshing = true;
    clarityError = null;
    clarityLoading = true;
    const syncResult = await triggerProjectClaritySync(projectId, { num_days: 3 });
    if (syncResult.error) {
      clarityError = syncResult.error.message || 'Failed to sync Microsoft Clarity data.';
      clarityRefreshing = false;
      clarityLoading = false;
      return;
    }
    await loadClarityData(projectId);
    clarityRefreshing = false;
  }

  $: if (projectId && projectId !== clarityInitializedProjectId) {
    clarityInitializedProjectId = projectId;
    loadClarityData(projectId);
  }

  $: cachedEnrichedIssues = buildEnrichedIssues(summary?.issues || [], gscPageRows);

  $: activeEnrichedIssues = cachedEnrichedIssues;

  $: displayIssues = activeEnrichedIssues.length > 0
    ? activeEnrichedIssues.map((ei) => ei.issue)
    : (summary?.issues || []);
  
  $: enrichedIssuesMap = activeEnrichedIssues.reduce((acc, ei) => {
    if (ei?.issue?.url && ei?.issue?.type) {
      acc[`${ei.issue.url}|${ei.issue.type}`] = ei;
    }
    return acc;
  }, {});

  $: gscProperty = gscStatus?.integration?.property_url || null;
  $: gscLastSynced = gscStatus?.sync_state?.last_synced_at ? formatDateTime(gscStatus.sync_state.last_synced_at) : null;
</script>

<div class="flex flex-col lg:flex-row min-h-[calc(100vh-200px)] bg-base-100 border-t border-base-200">
  <!-- Sidebar Navigation -->
  <aside class="w-full lg:w-64 bg-base-100 lg:border-r border-base-200 flex-shrink-0">
    <ul class="menu menu-horizontal lg:menu-vertical p-2 lg:p-4 w-full overflow-x-auto lg:overflow-visible whitespace-nowrap lg:whitespace-normal space-x-2 lg:space-x-0 lg:space-y-1">
      <li>
        <button 
          type="button" 
          class:active={activeTab === 'dashboard'}
          on:click={() => navigateToTab('dashboard')}
        >
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6A2.25 2.25 0 016 3.75h2.25A2.25 2.25 0 0110.5 6v2.25a2.25 2.25 0 01-2.25 2.25H6a2.25 2.25 0 01-2.25-2.25V6zM3.75 15.75A2.25 2.25 0 016 13.5h2.25a2.25 2.25 0 012.25 2.25V18a2.25 2.25 0 01-2.25 2.25H6A2.25 2.25 0 013.75 18v-2.25zM13.5 6a2.25 2.25 0 012.25-2.25H18A2.25 2.25 0 0120.25 6v2.25A2.25 2.25 0 0118 10.5h-2.25a2.25 2.25 0 01-2.25-2.25V6zM13.5 15.75a2.25 2.25 0 012.25-2.25H18a2.25 2.25 0 012.25 2.25V18A2.25 2.25 0 0118 20.25h-2.25A2.25 2.25 0 0113.5 18v-2.25z" />
          </svg>
          Dashboard
        </button>
      </li>
      <li>
        <button 
          type="button" 
          class:active={activeTab === 'results'}
          on:click={() => navigateToTab('results')}
        >
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M3.375 19.5h17.25m-17.25 0a1.125 1.125 0 01-1.125-1.125M3.375 19.5h7.5c.621 0 1.125-.504 1.125-1.125m-9.75 0V5.625m0 12.75v-1.5c0-.621.504-1.125 1.125-1.125m18.375 2.625V5.625m0 12.75c0 .621-.504 1.125-1.125 1.125m1.125-1.125v-1.5c0-.621-.504-1.125-1.125-1.125m0 3.75h-7.5A1.125 1.125 0 0112 18.375m9.75-12.75c0-.621-.504-1.125-1.125-1.125H3.375c-.621 0-1.125.504-1.125 1.125m19.5 0v1.5c0 .621-.504 1.125-1.125 1.125M2.25 5.625v1.5c0 .621.504 1.125 1.125 1.125m0 0h17.25m-17.25 0h7.5c.621 0 1.125.504 1.125 1.125M3.375 8.25c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125m17.25-3.75h-7.5c-.621 0-1.125.504-1.125 1.125m8.625-1.125c.621 0 1.125.504 1.125 1.125v1.5c0 .621-.504 1.125-1.125 1.125m-17.25 0h7.5m-7.5 0c-.621 0-1.125.504-1.125 1.125v1.5c0 .621.504 1.125 1.125 1.125M12 10.875v-1.5m0 1.5c0 .621-.504 1.125-1.125 1.125M12 10.875c0 .621.504 1.125 1.125 1.125m-2.25 0c.621 0 1.125.504 1.125 1.125M13.125 12h7.5m-7.5 0c-.621 0-1.125.504-1.125 1.125M20.625 12c.621 0 1.125.504 1.125 1.125v1.5c0 .621-.504 1.125-1.125 1.125m-17.25 0h7.5M12 14.625v-1.5m0 1.5c0 .621-.504 1.125-1.125 1.125M12 14.625c0 .621.504 1.125 1.125 1.125m-2.25 0c.621 0 1.125.504 1.125 1.125m0 1.5v-1.5m0 1.5c0 .621-.504 1.125-1.125 1.125M12 18.375v-1.5m0 1.5c0 .621-.504 1.125-1.125 1.125M12 18.375c0 .621.504 1.125 1.125 1.125" />
          </svg>
          Results
        </button>
      </li>
      <li>
        <button 
          type="button" 
          class:active={activeTab === 'issues'}
          on:click={() => navigateToTab('issues')}
        >
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z" />
          </svg>
          Issues
        </button>
      </li>
      <li>
        <button 
          type="button" 
          class:active={activeTab === 'recommendations'}
          on:click={() => navigateToTab('recommendations')}
        >
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 18v-5.25m0 0a6.01 6.01 0 001.5-.189m-1.5.189a6.01 6.01 0 01-1.5-.189m3.75 7.478a12.06 12.06 0 01-4.5 0m3.75 2.383a14.406 14.406 0 01-3 0M14.25 18v-.192c0-.983.658-1.823 1.508-2.316a7.5 7.5 0 10-7.517 0c.85.493 1.509 1.333 1.509 2.316V18" />
          </svg>
          Recommendations
        </button>
      </li>
      <li>
        <button 
          type="button" 
          class:active={activeTab === 'graph'}
          on:click={() => navigateToTab('graph')}
        >
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M7.217 10.907a2.25 2.25 0 100 2.186m0-2.186c.18.324.283.696.283 1.093s-.103.77-.283 1.093m0-2.186l9.566-5.314m-9.566 7.5l9.566 5.314m0 0a2.25 2.25 0 103.935 2.186 2.25 2.25 0 00-3.935-2.186zm0-12.814a2.25 2.25 0 103.933-2.185 2.25 2.25 0 00-3.933 2.185z" />
          </svg>
          Link Graph
        </button>
      </li>
      {#if projectId}
        <li>
          <a
            href="/project/{projectId}/crawls"
            use:link
            class:active={false}
          >
            <FileSearch class="w-5 h-5" />
            Crawls
          </a>
        </li>
      {/if}
      
      <!-- Rank Tracking Section -->
      {#if projectId}
        <li class="hidden lg:block border-b border-base-200 my-2 pointer-events-none"></li>
        <li>
          <a
            href="/project/{projectId}/rank-tracker"
            use:link
            class:active={false}
          >
            <TrendingUp class="w-5 h-5" />
            Rank Tracker
          </a>
        </li>
        <li>
          <a
            href="/project/{projectId}/discover-keywords"
            use:link
            class:active={false}
          >
            <ScanSearch class="w-5 h-5" />
            Discover Keywords
          </a>
        </li>
        <li>
          <a
            href="/project/{projectId}/impact-first"
            use:link
            class:active={false}
          >
            <Target class="w-5 h-5" />
            Impact-First View
          </a>
        </li>
      {/if}
      
      <!-- Google Search Console Section (only if connected) -->
      {#if projectId && (gscStatus?.integration?.property_url || ga4Status?.integration?.property_id || clarityStatus?.integration?.connected)}
        <li class="hidden lg:block border-b border-base-200 my-2 pointer-events-none"></li>
        {#if gscStatus?.integration?.property_url}
          <li>
            <button
              type="button"
              class:active={activeTab === 'gsc-dashboard'}
              on:click={() => navigateToTab('gsc-dashboard')}
            >
              <BarChart class="w-5 h-5" />
              GSC Dashboard
            </button>
          </li>
          <li>
            <button
              type="button"
              class:active={activeTab === 'gsc-keywords'}
              on:click={() => navigateToTab('gsc-keywords')}
            >
              <Search class="w-5 h-5" />
              GSC Keywords
            </button>
          </li>
        {/if}
        {#if ga4Status?.integration?.property_id}
          <li>
            <button
              type="button"
              class:active={activeTab === 'ga4-dashboard'}
              on:click={() => navigateToTab('ga4-dashboard')}
            >
              <BarChart class="w-5 h-5" />
              GA4 Dashboard
            </button>
          </li>
        {/if}
        {#if clarityStatus?.integration?.connected}
          <li>
            <button
              type="button"
              class:active={activeTab === 'clarity-dashboard'}
              on:click={() => navigateToTab('clarity-dashboard')}
            >
              <AlertTriangle class="w-5 h-5" />
              Clarity
            </button>
          </li>
        {/if}
      {/if}
      {#if projectId}
        <li class="hidden lg:block border-b border-base-200 my-2 pointer-events-none"></li>
        <li>
          <button
            type="button"
            class:active={activeTab === 'insights'}
            on:click={() => navigateToTab('insights')}
          >
            <Lightbulb class="w-5 h-5" />
            Insights
          </button>
        </li>
      {/if}
      
      <!-- Project Settings -->
      {#if projectId}
        <li class="hidden lg:block border-b border-base-200 my-2 pointer-events-none"></li>
        <li>
          <a 
            href="/project/{projectId}/settings" 
            use:link
            class:active={isSettingsPage}
          >
            <Settings class="w-5 h-5" />
            Settings
          </a>
        </li>
      {/if}
    </ul>
  </aside>

  <!-- Main Content -->
  <main class="flex-1 p-4 lg:p-8 overflow-y-auto">
    {#if activeTab === 'dashboard'}
    <div class="space-y-4">

      <SummaryCard
        {summary}
        {navigateToTab}
        {projectId}
        gscTotals={gscStatus?.summary?.totals}
        gscSyncState={gscStatus?.sync_state}
        gscIntegration={gscStatus?.integration}
        gscLoading={gscLoading}
        gscError={gscError}
      />

      {#if crawlId}
        <CrawlSummary {crawlId} />
        <PublicReportGenerator {crawlId} {projectId} />
      {/if}

      <!-- Recent Crawls Section -->
      {#if projectId}
        <div class="card bg-base-100 shadow">
          <div class="card-body">
            <div class="flex items-center justify-between mb-4">
              <h2 class="card-title text-xl">
                <FileSearch class="w-5 h-5" />
                Recent Crawls
              </h2>
              <a href="/project/{projectId}/crawls" use:link class="btn btn-sm btn-ghost">
                View All
                <ArrowRight class="w-4 h-4 ml-1" />
              </a>
            </div>
            
            {#if crawlsLoading}
              <div class="flex justify-center py-4">
                <span class="loading loading-spinner loading-sm"></span>
              </div>
            {:else if recentCrawls.length === 0}
              <div class="text-center py-4 text-base-content/70">
                <p class="mb-2">No crawls yet</p>
                <a href="/project/{projectId}/crawls" use:link class="btn btn-sm btn-primary">
                  Start Your First Crawl
                </a>
              </div>
            {:else}
              <div class="overflow-x-auto">
                <table class="table table-sm">
                  <thead>
                    <tr>
                      <th>Date</th>
                      <th>Status</th>
                      <th>Pages</th>
                      <th>Issues</th>
                      <th>Action</th>
                    </tr>
                  </thead>
                  <tbody>
                    {#each recentCrawls as crawl}
                      <tr class="hover">
                        <td class="text-sm">{formatDate(crawl.started_at)}</td>
                        <td>
                          <span class="badge {getStatusBadge(crawl.status)} badge-xs capitalize">
                            {crawl.status}
                          </span>
                        </td>
                        <td>{crawl.page_count || crawl.indexed_pages || crawl.total_pages || 0}</td>
                        <td>
                          {#if crawl.total_issues > 0}
                            <span class="text-error font-semibold">{crawl.total_issues}</span>
                          {:else}
                            <span class="text-base-content/60">0</span>
                          {/if}
                        </td>
                        <td>
                          <a 
                            href="/project/{projectId}/crawl/{crawl.id}" 
                            use:link
                            class="btn btn-ghost btn-xs"
                          >
                            View
                          </a>
                        </td>
                      </tr>
                    {/each}
                  </tbody>
                </table>
              </div>
            {/if}
          </div>
        </div>
      {/if}
    </div>
  {:else if activeTab === 'results'}
    <ResultsTable 
      {results} 
      issues={displayIssues}
      filter={resultsFilter}
      {navigateToTab}
    />
  {:else if activeTab === 'issues'}
    <IssuesPanel
      issues={displayIssues}
      filter={issuesFilter}
      enrichedIssues={enrichedIssuesMap}
      gscStatus={gscStatus}
      gscLoading={gscLoading}
      gscError={gscError}
      {crawlId}
    />
  {:else if activeTab === 'recommendations'}
    <div class="space-y-4">
      <RecommendationsPanel issues={displayIssues} {navigateToTab} enrichedIssues={enrichedIssuesMap} />
    </div>
  {:else if activeTab === 'graph'}
    <LinkGraph crawlId={crawlId} />
  {:else if activeTab === 'gsc-dashboard'}
    <GSCDashboardPanel 
      {projectId}
      {gscStatus}
      {gscLoading}
      {gscRefreshing}
      {gscError}
      onRefresh={refreshGSCData}
    />
  {:else if activeTab === 'gsc-keywords'}
    <GSCKeywordsPanel
      {projectId}
      {gscStatus}
      {gscLoading}
      {gscError}
    />
  {:else if activeTab === 'ga4-dashboard'}
    <GA4DashboardPanel
      {projectId}
      {ga4Status}
      {ga4Loading}
      {ga4Refreshing}
      {ga4Error}
      onRefresh={refreshGA4Data}
    />
  {:else if activeTab === 'clarity-dashboard'}
    <ClarityDashboardPanel
      {projectId}
      {clarityStatus}
      {clarityLoading}
      {clarityRefreshing}
      {clarityError}
      onRefresh={refreshClarityData}
    />
  {:else if activeTab === 'insights'}
    <UnifiedInsightsPanel {projectId} />
  {/if}
  </main>
</div>

<style>
  /* Make active menu items more prominent */
  .menu li button.active,
  .menu li a.active {
    background-color: hsl(var(--p) / 0.1) !important;
    color: hsl(var(--p)) !important;
    font-weight: 600;
    border-left: 3px solid hsl(var(--p));
  }
  
  /* Add hover effect for non-active items */
  .menu li button:not(.active):hover,
  .menu li a:not(.active):hover {
    background-color: hsl(var(--bc) / 0.05);
  }
  
  /* Ensure icons in active items are also colored */
  .menu li button.active svg,
  .menu li a.active svg {
    color: hsl(var(--p));
  }
</style>