<script>
  import { onMount } from 'svelte';
  import { push, querystring, link } from 'svelte-spa-router';
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
  import ProjectSidebar from './ProjectSidebar.svelte';
  import CrawlSelector from './CrawlSelector.svelte';
  import TriggerCrawlButton from './TriggerCrawlButton.svelte';
  import { fetchProjects, fetchProjectGSCStatus, fetchProjectGSCDimensions, triggerProjectGSCSync, fetchProjectGA4Status, triggerProjectGA4Sync, fetchProjectClarityStatus, triggerProjectClaritySync, fetchCrawls } from '../lib/data.js';
  import { buildEnrichedIssues } from '../lib/gsc.js';
  import { 
    FileSearch,
    ArrowRight
  } from 'lucide-svelte';

  import { createEventDispatcher } from 'svelte';

  export let summary = null;
  export let results = [];
  export let initialTab = 'dashboard';
  export let projectId = null;
  export let crawlId = null;
  export let project = null;
  export let crawls = [];
  export let selectedCrawl = null;

  const dispatch = createEventDispatcher();

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

  // Sync issuesFilter.url from querystring when tab=issues (e.g. deep link from Insights)
  $: if ($querystring && activeTab === 'issues') {
    const qs = new URLSearchParams($querystring);
    const urlFromQs = qs.get('url');
    if (urlFromQs && issuesFilter.url !== urlFromQs) {
      issuesFilter = { ...issuesFilter, url: urlFromQs };
    }
  }
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
    const params = new URLSearchParams();
    params.set('tab', tab);
    if (tab === 'issues' && url) {
      params.set('url', url);
    }
    if (projectId && crawlId) {
      push(`/project/${projectId}/crawl/${crawlId}?${params.toString()}`);
    } else if (projectId && (tab === 'gsc-dashboard' || tab === 'gsc-keywords' || tab === 'ga4-dashboard' || tab === 'clarity-dashboard')) {
      // Integration tabs work at project level, redirect to first crawl or project view
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
      // Refetch status so sync_state (e.g. rate_limit retry_after) is up to date
      await loadClarityData(projectId);
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
  <ProjectSidebar
    {projectId}
    {activeTab}
    {gscStatus}
    {ga4Status}
    {clarityStatus}
    {navigateToTab}
    mode="crawl"
  />

  <!-- Main Content -->
  <main class="flex-1 p-4 lg:p-8 overflow-y-auto">
    {#if activeTab === 'dashboard'}
    <div class="space-y-4">

      {#if crawls.length > 0}
        <div class="flex justify-between items-center">
          <CrawlSelector {crawls} {selectedCrawl} {projectId} on:select />
          <TriggerCrawlButton {projectId} {project} on:created />
        </div>
      {/if}

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
    <div class="space-y-6">
      <RecommendationsPanel issues={displayIssues} {navigateToTab} enrichedIssues={enrichedIssuesMap} />
      {#if projectId}
        <UnifiedInsightsPanel
          {projectId}
          onNavigateToIssues={(url) => navigateToTab('issues', { url })}
        />
      {/if}
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
  {/if}
  </main>
</div>
