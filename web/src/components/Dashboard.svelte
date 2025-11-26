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
  import CrawlSummary from './AI/CrawlSummary.svelte';
  import PublicReportGenerator from './PublicReportGenerator.svelte';
  import Logo from './Logo.svelte';
  import { fetchProjects, fetchProjectGSCStatus, fetchProjectGSCDimensions, triggerProjectGSCSync } from '../lib/data.js';
  import { buildEnrichedIssues } from '../lib/gsc.js';
  import { userProfile, isProOrTeam } from '../lib/subscription.js';

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
  });

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
    } else if (projectId && (tab === 'gsc-dashboard' || tab === 'gsc-keywords')) {
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
      
      {#if projectId && gscStatus?.integration?.property_url}
        <li class="hidden lg:block border-b border-base-200 my-1 pointer-events-none"></li>
        <li>
          <button 
            type="button" 
            class:active={activeTab === 'gsc-dashboard'}
            on:click={() => navigateToTab('gsc-dashboard')}
          >
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z" />
            </svg>
            GSC Dashboard
          </button>
        </li>
        <li>
          <button 
            type="button" 
            class:active={activeTab === 'gsc-keywords'}
            on:click={() => navigateToTab('gsc-keywords')}
          >
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
            </svg>
            GSC Keywords
          </button>
        </li>
      {/if}
      
      {#if projectId}
        <li class="hidden lg:block border-b border-base-200 my-1 pointer-events-none"></li>
        <li>
          <a
            href="/project/{projectId}/rank-tracker"
            use:link
            class="btn btn-ghost btn-sm"
            type="button"
          >
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z" />
            </svg>
            Rank Tracker
          </a>
        </li>
        <li>
          <a 
            href="/project/{projectId}/settings" 
            use:link
            class:active={isSettingsPage}
          >
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M9.594 3.94c.09-.542.56-.94 1.11-.94h2.593c.55 0 1.02.398 1.11.94l.213 1.281c.063.374.313.686.645.87.074.04.147.083.22.127.324.196.72.257 1.075.124l1.217-.456a1.125 1.125 0 011.37.49l1.296 2.247a1.125 1.125 0 01-.26 1.431l-1.003.827c-.293.24-.438.613-.431.992a6.759 6.759 0 010 .255c-.007.378.138.75.43.99l1.005.828c.424.35.534.954.26 1.43l-1.298 2.247a1.125 1.125 0 01-1.369.491l-1.217-.456c-.355-.133-.75-.072-1.076.124a6.57 6.57 0 01-.22.128c-.331.183-.581.495-.644.869l-.213 1.28c-.09.543-.56.941-1.11.941h-2.594c-.55 0-1.02-.398-1.11-.94l-.213-1.281c-.062-.374-.312-.686-.644-.87a6.52 6.52 0 01-.22-.127c-.325-.196-.72-.257-1.076-.124l-1.217.456a1.125 1.125 0 01-1.369-.49l-1.297-2.247a1.125 1.125 0 01.26-1.431l1.004-.827c.292-.24.437-.613.43-.992a6.932 6.932 0 010-.255c.007-.378-.138-.75-.43-.99l-1.004-.828a1.125 1.125 0 01-.26-1.43l1.297-2.247a1.125 1.125 0 011.37-.491l1.216.456c.356.133.751.072 1.076-.124.072-.044.146-.087.22-.128.332-.183.582-.495.644-.869l.214-1.281z" />
              <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
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