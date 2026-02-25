<script>
  import { link } from 'svelte-spa-router';
  import { TrendingUp, Hash, FileText, AlertTriangle, Zap, AlertCircle, Link2, ExternalLink, ArrowRight } from 'lucide-svelte';
  
  export let summary = null;
  export let navigateToTab = () => {};
  export let projectId = null;
  export let gscTotals = null;
  export let gscSyncState = null;
  export let gscIntegration = null;
  export let gscLoading = false;
  export let gscError = null;

  if (!summary) {
    summary = {
      total_pages: 0,
      total_issues: 0,
      average_response_time_ms: 0,
      pages_with_errors: 0,
      pages_with_redirects: 0,
      total_internal_links: 0,
      total_external_links: 0,
      issues_by_type: {}
    };
  }

  const formatNumber = (value) => {
    if (value === null || value === undefined) return '0';
    const numeric = Number(value);
    if (Number.isNaN(numeric)) return '0';
    return Math.round(numeric).toLocaleString();
  };

  const formatCTR = (value) => {
    if (value === null || value === undefined) return '0.00%';
    const numeric = Number(value);
    if (Number.isNaN(numeric)) return '0.00%';
    const percent = numeric > 1 ? numeric : numeric * 100;
    return `${percent.toFixed(2)}%`;
  };

  const formatPosition = (value) => {
    if (value === null || value === undefined) return '0.0';
    const numeric = Number(value);
    if (Number.isNaN(numeric)) return '0.0';
    return numeric.toFixed(1);
  };

  const formatLastSynced = (value) => {
    if (!value) return null;
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return null;
    return `${date.toLocaleDateString()} ${date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}`;
  };

  const getSeverityCount = (severity) => {
    if (!summary.issues) return 0;
    return summary.issues.filter(i => i.severity === severity).length;
  };

  const handleTotalIssuesClick = () => {
    if (summary.total_issues > 0) {
      navigateToTab('issues');
    }
  };

  const handleFixCriticalIssues = () => {
    navigateToTab('issues', { severity: 'error' });
  };

  const handleViewSlowPages = () => {
    navigateToTab('results', { performance: true });
  };

  $: lastSyncedDisplay = gscSyncState?.last_synced_at ? formatLastSynced(gscSyncState.last_synced_at) : null;
  $: isGSCConnected = Boolean(gscIntegration?.property_url);
  $: errorCount = getSeverityCount('error');
  $: warningCount = getSeverityCount('warning');
  $: infoCount = getSeverityCount('info');
</script>

<!-- Primary Stats Row -->
<div class="grid grid-cols-2 lg:grid-cols-4 gap-3 mb-6">
  <div class="bg-base-200 rounded-xl p-4">
    <div class="flex items-center gap-2 mb-2">
      <FileText class="w-4 h-4 text-base-content/40" />
      <p class="text-xs text-base-content/50 font-medium">Pages Crawled</p>
    </div>
    <p class="text-2xl font-bold text-primary font-mono">{formatNumber(summary.page_count || summary.indexed_pages || summary.total_pages || 0)}</p>
  </div>

  <button
    class="bg-base-200 rounded-xl p-4 text-left hover:bg-base-300 transition-colors cursor-pointer"
    on:click={handleTotalIssuesClick}
  >
    <div class="flex items-center gap-2 mb-2">
      <AlertTriangle class="w-4 h-4 text-base-content/40" />
      <p class="text-xs text-base-content/50 font-medium">Issues Found</p>
    </div>
    <p class="text-2xl font-bold text-warning font-mono">{formatNumber(summary.total_issues)}</p>
    {#if summary.total_issues > 0}
      <p class="text-xs text-base-content/40 mt-1">Click to view</p>
    {/if}
  </button>

  <div class="bg-base-200 rounded-xl p-4">
    <div class="flex items-center gap-2 mb-2">
      <Zap class="w-4 h-4 text-base-content/40" />
      <p class="text-xs text-base-content/50 font-medium">Avg Response</p>
    </div>
    <p class="text-2xl font-bold text-info font-mono">{summary.average_response_time_ms}<span class="text-sm text-base-content/40 ml-0.5">ms</span></p>
  </div>

  <div class="bg-base-200 rounded-xl p-4">
    <div class="flex items-center gap-2 mb-2">
      <AlertCircle class="w-4 h-4 text-base-content/40" />
      <p class="text-xs text-base-content/50 font-medium">Pages with Errors</p>
    </div>
    <p class="text-2xl font-bold text-error font-mono">{formatNumber(summary.pages_with_errors)}</p>
  </div>
</div>

<!-- GSC Stats Row -->
{#if gscTotals}
  <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 mb-3">
    <div class="bg-base-200 rounded-xl p-4">
      <div class="flex items-center gap-2 mb-2">
        <TrendingUp class="w-4 h-4 text-base-content/40" />
        <p class="text-xs text-base-content/50 font-medium">Impressions</p>
      </div>
      <p class="text-2xl font-bold text-primary font-mono">{formatNumber(gscTotals.impressions)}</p>
    </div>

    <div class="bg-base-200 rounded-xl p-4">
      <div class="flex items-center gap-2 mb-2">
        <ArrowRight class="w-4 h-4 text-base-content/40" />
        <p class="text-xs text-base-content/50 font-medium">Clicks</p>
      </div>
      <p class="text-2xl font-bold text-success font-mono">{formatNumber(gscTotals.clicks)}</p>
    </div>

    <div class="bg-base-200 rounded-xl p-4">
      <div class="flex items-center gap-2 mb-2">
        <TrendingUp class="w-4 h-4 text-base-content/40" />
        <p class="text-xs text-base-content/50 font-medium">CTR</p>
      </div>
      <p class="text-2xl font-bold text-info font-mono">{formatCTR(gscTotals.ctr)}</p>
    </div>

    <div class="bg-base-200 rounded-xl p-4">
      <div class="flex items-center gap-2 mb-2">
        <Hash class="w-4 h-4 text-base-content/40" />
        <p class="text-xs text-base-content/50 font-medium">Avg Position</p>
      </div>
      <p class="text-2xl font-bold text-warning font-mono">{formatPosition(gscTotals.position)}</p>
    </div>
  </div>

  {#if lastSyncedDisplay}
    <p class="text-xs text-base-content/40 mb-4">
      Search Console data last synced {lastSyncedDisplay}.
    </p>
  {/if}
  
  {#if projectId}
    <div class="mb-6">
      <a
        href="/project/{projectId}/gsc"
        use:link
        class="btn btn-primary btn-sm"
      >
        View GSC Dashboard
      </a>
    </div>
  {/if}
{:else if isGSCConnected && !gscLoading && !gscError}
  <div class="alert alert-info mb-6">
    <span>Google Search Console is connected. Refresh the data to populate performance metrics.</span>
  </div>
{/if}

<!-- Severity + Links + Issue Types -->
<div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
  <!-- Severity Breakdown -->
  <div class="bg-base-200 rounded-xl p-5">
    <h3 class="text-xs font-semibold uppercase tracking-wider text-base-content/40 mb-4">Severity Breakdown</h3>
    <div class="space-y-3">
      <div class="flex items-center gap-3">
        <div class="w-2 h-2 rounded-full bg-error shrink-0"></div>
        <span class="text-sm text-base-content/70 flex-1">Errors</span>
        <span class="font-bold font-mono text-sm">{errorCount}</span>
      </div>
      <div class="flex items-center gap-3">
        <div class="w-2 h-2 rounded-full bg-warning shrink-0"></div>
        <span class="text-sm text-base-content/70 flex-1">Warnings</span>
        <span class="font-bold font-mono text-sm">{warningCount}</span>
      </div>
      <div class="flex items-center gap-3">
        <div class="w-2 h-2 rounded-full bg-info shrink-0"></div>
        <span class="text-sm text-base-content/70 flex-1">Info</span>
        <span class="font-bold font-mono text-sm">{infoCount}</span>
      </div>
    </div>
    {#if errorCount > 0}
      <button class="btn btn-error btn-sm w-full mt-4" on:click={handleFixCriticalIssues}>
        Fix Critical Issues
      </button>
    {/if}
  </div>

  <!-- Link Statistics -->
  <div class="bg-base-200 rounded-xl p-5">
    <h3 class="text-xs font-semibold uppercase tracking-wider text-base-content/40 mb-4">Link Statistics</h3>
    <div class="space-y-3">
      <div class="flex items-center gap-3">
        <Link2 class="w-4 h-4 text-base-content/30 shrink-0" />
        <span class="text-sm text-base-content/70 flex-1">Internal Links</span>
        <span class="font-bold font-mono text-sm">{formatNumber(summary.total_internal_links)}</span>
      </div>
      <div class="flex items-center gap-3">
        <ExternalLink class="w-4 h-4 text-base-content/30 shrink-0" />
        <span class="text-sm text-base-content/70 flex-1">External Links</span>
        <span class="font-bold font-mono text-sm">{formatNumber(summary.total_external_links)}</span>
      </div>
      <div class="flex items-center gap-3">
        <ArrowRight class="w-4 h-4 text-base-content/30 shrink-0" />
        <span class="text-sm text-base-content/70 flex-1">Redirects</span>
        <span class="font-bold font-mono text-sm">{formatNumber(summary.pages_with_redirects)}</span>
      </div>
    </div>
  </div>

  <!-- Issue Types as Pill Grid -->
  <div class="bg-base-200 rounded-xl p-5">
    <h3 class="text-xs font-semibold uppercase tracking-wider text-base-content/40 mb-4">Issue Types</h3>
    <div class="flex flex-wrap gap-2">
      {#each Object.entries(summary.issues_by_type || {}) as [type, count]}
        <span class="bg-base-300 text-base-content/70 text-xs px-3 py-1.5 rounded-lg font-medium">
          {type.replace(/_/g, ' ')}
          <span class="text-primary ml-1 font-bold">{count}</span>
        </span>
      {/each}
    </div>
    {#if summary.total_issues > 0}
      <button class="btn btn-primary btn-sm w-full mt-4" on:click={() => navigateToTab('issues')}>
        View All Issues
      </button>
    {/if}
  </div>
</div>

<!-- Slowest Pages -->
{#if summary.slowest_pages && summary.slowest_pages.length > 0}
  <div class="bg-base-200 rounded-xl p-5">
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-xs font-semibold uppercase tracking-wider text-base-content/40">Slowest Pages</h3>
      <button class="btn btn-outline btn-xs" on:click={handleViewSlowPages}>
        View All
      </button>
    </div>
    <div class="space-y-2">
      {#each summary.slowest_pages.slice(0, 5) as page}
        <div class="flex items-center gap-3 bg-base-300 rounded-lg px-3 py-2">
          <span class="text-sm text-base-content/60 truncate flex-1">{page.url}</span>
          <span class="text-xs font-mono font-bold text-warning whitespace-nowrap">{page.response_time_ms}ms</span>
        </div>
      {/each}
    </div>
  </div>
{/if}
