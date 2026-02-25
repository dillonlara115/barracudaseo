<script>
  import IssueInsight from './AI/IssueInsight.svelte';
  import { Sparkles, ChevronDown, ChevronUp, Download } from 'lucide-svelte';
  import { userProfile, isProOrTeam } from '../lib/subscription.js';
  import { collapseImageIssuesByAsset } from '../lib/issueUtils.js';

  $: isPro = isProOrTeam($userProfile);

  export let issues = [];
  export let filter = { severity: 'all', type: 'all', url: null };
  export let enrichedIssues = {};
  export let gscStatus = null;
  export let gscLoading = false;
  export let gscError = null;
  export let crawlId = null;

  let selectedIssueForAI = null;
  let showAIInsightModal = false;
  let expandedIssues = new Set();

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

  const formatDateTime = (value) => {
    if (!value) return null;
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return null;
    return `${date.toLocaleDateString()} ${date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}`;
  };

  const normalizeSeverity = (value) => value ?? 'all';
  const normalizeType = (value) => value ?? 'all';
  const normalizeUrl = (value) => value ?? '';

  const computeFilterSignature = (incoming = {}) => [
    normalizeSeverity(incoming.severity),
    normalizeType(incoming.type),
    normalizeUrl(incoming.url)
  ].join('|');

  const applyFiltersFromProps = (incoming = {}) => {
    severityFilter = normalizeSeverity(incoming.severity);
    typeFilter = normalizeType(incoming.type);
    searchTerm = normalizeUrl(incoming.url);
  };

  let severityFilter = normalizeSeverity(filter?.severity);
  let typeFilter = normalizeType(filter?.type);
  let searchTerm = normalizeUrl(filter?.url);
  let groupBy = 'none';
  let sortBy = 'none';

  let lastAppliedFilterSignature = computeFilterSignature(filter);

  $: {
    const nextSignature = computeFilterSignature(filter);
    if (nextSignature !== lastAppliedFilterSignature) {
      applyFiltersFromProps(filter);
      lastAppliedFilterSignature = nextSignature;
    }
  }

  $: affectedPagesByType = issues.reduce((acc, issue) => {
    if (!acc[issue.type]) acc[issue.type] = new Set();
    acc[issue.type].add(issue.url);
    return acc;
  }, {});

  $: affectedPagesCounts = Object.entries(affectedPagesByType).reduce((acc, [type, urlSet]) => {
    acc[type] = urlSet.size;
    return acc;
  }, {});

  const getSeverityWeight = (severity) => {
    switch (severity) {
      case 'error': return 10;
      case 'warning': return 5;
      case 'info': return 1;
      default: return 1;
    }
  };

  const calculatePriorityScore = (issue) => {
    const enrichedKey = `${issue.url}|${issue.type}`;
    const enriched = enrichedIssues[enrichedKey];
    if (enriched && enriched.enriched_priority) return enriched.enriched_priority;
    const severityWeight = getSeverityWeight(issue.severity);
    const pagesAffected = affectedPagesCounts[issue.type] || 0;
    return severityWeight * pagesAffected;
  };

  const getDisplayItemPriority = (item) => {
    if (item.isAssetGroup && item.affectedUrls?.length) {
      const enrichedKey = `${item.url}|${item.type}`;
      const enriched = enrichedIssues[enrichedKey];
      if (enriched && enriched.enriched_priority) return enriched.enriched_priority;
      return getSeverityWeight(item.severity) * item.affectedUrls.length;
    }
    return calculatePriorityScore(item);
  };

  const getDisplayItemEnriched = (item) => {
    if (item.isAssetGroup) {
      for (const id of item.constituentIds || []) {
        if (enrichedIssues[id]) return enrichedIssues[id];
      }
    }
    return enrichedIssues[`${item.url}|${item.type}`];
  };

  const isDisplayItemTopPriority = (item) => {
    if (item.isAssetGroup && item.constituentIds?.length) {
      return item.constituentIds.some((id) => top10PriorityIssues.has(id));
    }
    return top10PriorityIssues.has(`${item.url}|${item.type}`);
  };

  $: issuesWithPriority = issues.map(issue => ({
    ...issue,
    priorityScore: calculatePriorityScore(issue)
  }));

  $: top10PriorityIssues = new Set(
    issuesWithPriority
      .sort((a, b) => b.priorityScore - a.priorityScore)
      .slice(0, 10)
      .map(issue => `${issue.url}|${issue.type}`)
  );

  $: filteredIssues = issues.filter(i => {
    if (severityFilter !== 'all' && i.severity !== severityFilter) return false;
    if (typeFilter !== 'all' && i.type !== typeFilter) return false;
    if (filter.url && i.url !== filter.url) return false;
    if (searchTerm &&
        !i.url.toLowerCase().includes(searchTerm.toLowerCase()) &&
        !i.message?.toLowerCase().includes(searchTerm.toLowerCase()) &&
        !i.recommendation?.toLowerCase().includes(searchTerm.toLowerCase())) {
      return false;
    }
    return true;
  });

  $: sortedFilteredIssues = (() => {
    if (sortBy === 'priority' || sortBy === 'enriched_priority') {
      return [...filteredIssues].sort((a, b) => calculatePriorityScore(b) - calculatePriorityScore(a));
    }
    return filteredIssues;
  })();

  $: groupedIssues = (() => {
    if (groupBy === 'none') return { 'All Issues': sortedFilteredIssues };
    const grouped = {};
    sortedFilteredIssues.forEach(issue => {
      const key = groupBy === 'url' ? issue.url : groupBy === 'type' ? issue.type : issue.severity;
      if (!grouped[key]) grouped[key] = [];
      grouped[key].push(issue);
    });
    return grouped;
  })();

  $: uniqueTypes = [...new Set(issues.map(i => i.type))];

  const severityOptions = [
    { value: 'all', label: 'All' },
    { value: 'error', label: 'Errors' },
    { value: 'warning', label: 'Warnings' },
    { value: 'info', label: 'Info' }
  ];

  $: typeFilterOptions = [
    { value: 'all', label: 'All Types' },
    ...uniqueTypes.map(type => ({
      value: type,
      label: type.replace(/_/g, ' ')
    }))
  ];

  const groupByOptions = [
    { value: 'none', label: 'No Grouping' },
    { value: 'url', label: 'Group by URL' },
    { value: 'type', label: 'Group by Type' },
    { value: 'severity', label: 'Group by Severity' }
  ];
  
  $: gscIntegration = gscStatus?.integration || null;
  $: gscSyncState = gscStatus?.sync_state || null;
  $: gscLastSynced = gscSyncState?.last_synced_at ? formatDateTime(gscSyncState.last_synced_at) : null;
  $: hasGSCEnrichment = enrichedIssues && Object.keys(enrichedIssues).length > 0;

  const getSeverityDot = (severity) => {
    switch (severity) {
      case 'error': return 'bg-error';
      case 'warning': return 'bg-warning';
      case 'info': return 'bg-info';
      default: return 'bg-base-content/30';
    }
  };

  const getSeverityLabel = (severity) => {
    switch (severity) {
      case 'error': return 'Critical';
      case 'warning': return 'Warning';
      case 'info': return 'Info';
      default: return severity;
    }
  };

  const getSeverityTextColor = (severity) => {
    switch (severity) {
      case 'error': return 'text-error';
      case 'warning': return 'text-warning';
      case 'info': return 'text-info';
      default: return 'text-base-content/50';
    }
  };

  const toggleExpanded = (key) => {
    if (expandedIssues.has(key)) {
      expandedIssues.delete(key);
    } else {
      expandedIssues.add(key);
    }
    expandedIssues = new Set(expandedIssues);
  };

  const timestamp = () => {
    const now = new Date();
    const pad = (value) => value.toString().padStart(2, '0');
    return `${now.getFullYear()}${pad(now.getMonth() + 1)}${pad(now.getDate())}-${pad(now.getHours())}${pad(now.getMinutes())}${pad(now.getSeconds())}`;
  };

  const downloadFile = (content, fileName, mimeType) => {
    const blob = new Blob([content], { type: mimeType });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = fileName;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  const exportAsJson = () => {
    const data = filteredIssues.map(issue => ({ ...issue, priorityScore: calculatePriorityScore(issue) }));
    downloadFile(JSON.stringify(data, null, 2), `issues-${timestamp()}.json`, 'application/json');
  };

  const toCsv = (rows) => {
    const escapeValue = (value) => {
      if (value === null || value === undefined) return '';
      const s = String(value);
      return /[",\n]/.test(s) ? `"${s.replace(/"/g, '""')}"` : s;
    };
    const headers = ['url', 'type', 'severity', 'priority_score', 'message', 'recommendation', 'value'];
    const dataRows = rows.map(issue => [
      issue.url || '', issue.type || '', issue.severity || '',
      calculatePriorityScore(issue).toString(), issue.message || '',
      issue.recommendation || '', issue.value || ''
    ]);
    return [headers, ...dataRows].map(row => row.map(escapeValue).join(',')).join('\n');
  };

  const exportAsCsv = () => {
    downloadFile(toCsv(filteredIssues), `issues-${timestamp()}.csv`, 'text/csv');
  };

  const closeDropdown = (event) => {
    const details = event?.currentTarget?.closest('details');
    if (details) details.removeAttribute('open');
  };

  const handleExportCsv = (event) => { exportAsCsv(); closeDropdown(event); };
  const handleExportJson = (event) => { exportAsJson(); closeDropdown(event); };
</script>

<div class="space-y-4">
  <!-- Filters Card -->
  <div class="bg-base-200 rounded-xl p-5">
    <div class="flex flex-col gap-4">
      <!-- Search + Sort -->
      <div class="flex flex-col md:flex-row gap-3">
        <input
          type="text"
          placeholder="Search by URL, message, or recommendation..."
          class="input input-bordered flex-1 bg-base-300 border-base-content/10"
          bind:value={searchTerm}
        />
        <select class="select select-bordered bg-base-300 border-base-content/10 w-full md:w-44" bind:value={sortBy}>
          <option value="none">Sort by...</option>
          <option value="priority">Priority</option>
        </select>
      </div>

      <!-- Filter Pills -->
      <div class="flex flex-col sm:flex-row gap-4">
        <div class="flex-1">
          <p class="text-xs font-semibold uppercase tracking-wider text-base-content/40 mb-2">Severity</p>
          <div class="flex flex-wrap gap-1.5">
            {#each severityOptions as option}
              <button
                class="px-3 py-1.5 rounded-lg text-xs font-medium transition-colors {severityFilter === option.value ? 'bg-primary text-primary-content' : 'bg-base-300 text-base-content/60 hover:text-base-content'}"
                on:click={() => severityFilter = option.value}
              >
                {option.label}
              </button>
            {/each}
          </div>
        </div>

        <div class="flex-1">
          <p class="text-xs font-semibold uppercase tracking-wider text-base-content/40 mb-2">Group By</p>
          <div class="flex flex-wrap gap-1.5">
            {#each groupByOptions as option}
              <button
                class="px-3 py-1.5 rounded-lg text-xs font-medium transition-colors {groupBy === option.value ? 'bg-primary text-primary-content' : 'bg-base-300 text-base-content/60 hover:text-base-content'}"
                on:click={() => groupBy = option.value}
              >
                {option.label}
              </button>
            {/each}
          </div>
        </div>
      </div>

      <!-- Type filter -->
      {#if uniqueTypes.length > 1}
        <div>
          <p class="text-xs font-semibold uppercase tracking-wider text-base-content/40 mb-2">Type</p>
          <div class="flex flex-wrap gap-1.5">
            {#each typeFilterOptions as option}
              <button
                class="px-3 py-1.5 rounded-lg text-xs font-medium transition-colors {typeFilter === option.value ? 'bg-primary text-primary-content' : 'bg-base-300 text-base-content/60 hover:text-base-content'}"
                on:click={() => typeFilter = option.value}
              >
                {option.label}
              </button>
            {/each}
          </div>
        </div>
      {/if}

      {#if gscError}
        <div class="alert alert-warning"><span>{gscError}</span></div>
      {:else if gscLoading}
        <div class="alert alert-info"><span>Loading Google Search Console metrics...</span></div>
      {:else if gscIntegration && hasGSCEnrichment}
        <div class="bg-success/10 border border-success/20 rounded-lg px-4 py-2.5 text-sm text-base-content/70">
          Enriched with GSC data from {gscIntegration.property_url}.
          {#if gscLastSynced} Last synced {gscLastSynced}.{/if}
        </div>
      {:else if gscIntegration}
        <div class="bg-info/10 border border-info/20 rounded-lg px-4 py-2.5 text-sm text-base-content/70">
          GSC connected for {gscIntegration.property_url}. Refresh to populate metrics.
        </div>
      {/if}
    </div>
  </div>

  <!-- Results Header -->
  <div class="flex items-center justify-between">
    <p class="text-sm text-base-content/50">
      Showing {filteredIssues.length} of {issues.length} issues
      {#if filter.url}
        <span class="bg-info/20 text-info text-xs px-2 py-0.5 rounded-full ml-2">URL: {filter.url}</span>
      {/if}
      {#if sortBy === 'priority'}
        <span class="text-base-content/30 ml-1">· Sorted by priority</span>
      {/if}
    </p>
    <details class="dropdown dropdown-end">
      <summary class="btn btn-sm btn-ghost gap-1">
        <Download class="w-4 h-4" />
        Export
      </summary>
      <ul class="dropdown-content menu bg-base-200 rounded-xl w-44 shadow-lg mt-2 border border-base-content/5">
        <li><button type="button" on:click={handleExportCsv}>Export CSV</button></li>
        <li><button type="button" on:click={handleExportJson}>Export JSON</button></li>
      </ul>
    </details>
  </div>

  <!-- Issues List -->
  {#each Object.entries(groupedIssues) as [groupKey, groupIssues]}
    {@const displayItemsForGroup = collapseImageIssuesByAsset(groupIssues)}
    {#if groupBy !== 'none'}
      <div class="flex items-center gap-3 mt-6 mb-2">
        <div class="h-px bg-base-content/10 flex-1"></div>
        <h3 class="text-xs font-semibold uppercase tracking-wider text-base-content/40">
          {#if groupBy === 'severity'}
            {groupKey.charAt(0).toUpperCase() + groupKey.slice(1)} Issues
          {:else}
            {groupKey.replace(/_/g, ' ')}
          {/if}
          <span class="text-base-content/20 ml-1">({displayItemsForGroup.length})</span>
        </h3>
        <div class="h-px bg-base-content/10 flex-1"></div>
      </div>
    {/if}
    
    <div class="space-y-2">
      {#each displayItemsForGroup as displayItem, idx}
        {@const priorityScore = getDisplayItemPriority(displayItem)}
        {@const isTopPriority = isDisplayItemTopPriority(displayItem)}
        {@const enriched = getDisplayItemEnriched(displayItem)}
        {@const itemKey = `${groupKey}-${idx}`}
        {@const isExpanded = expandedIssues.has(itemKey)}

        <!-- Issue Row -->
        <div class="bg-base-200 rounded-xl overflow-hidden transition-colors hover:bg-base-300/50">
          <!-- Compact Row (always visible) -->
          <button
            class="w-full flex items-center gap-3 px-4 py-3 text-left"
            on:click={() => toggleExpanded(itemKey)}
          >
            <div class="w-2.5 h-2.5 rounded-full {getSeverityDot(displayItem.severity)} shrink-0"></div>
            <span class="text-sm text-base-content/80 flex-1 truncate">{displayItem.message}</span>
            <div class="flex items-center gap-2 shrink-0">
              {#if isTopPriority}
                <span class="bg-warning/15 text-warning text-xs px-2 py-0.5 rounded-full font-medium">Top Priority</span>
              {/if}
              {#if enriched?.enriched_priority}
                <span class="bg-success/15 text-success text-xs px-2 py-0.5 rounded-full font-medium">GSC</span>
              {/if}
              {#if affectedPagesCounts[displayItem.type] > 1}
                <span class="text-xs text-base-content/30 font-mono">{affectedPagesCounts[displayItem.type]}pg</span>
              {/if}
              <span class="text-xs font-medium whitespace-nowrap {getSeverityTextColor(displayItem.severity)}">{getSeverityLabel(displayItem.severity)}</span>
              {#if isExpanded}
                <ChevronUp class="w-4 h-4 text-base-content/30" />
              {:else}
                <ChevronDown class="w-4 h-4 text-base-content/30" />
              {/if}
            </div>
          </button>

          <!-- Expanded Detail -->
          {#if isExpanded}
            <div class="px-4 pb-4 pt-1 border-t border-base-content/5">
              <!-- Badges Row -->
              <div class="flex flex-wrap gap-1.5 mb-3">
                <span class="bg-base-300 text-base-content/60 text-xs px-2.5 py-1 rounded-lg">
                  {displayItem.type.replace(/_/g, ' ')}
                </span>
                <span class="bg-primary/15 text-primary text-xs px-2.5 py-1 rounded-lg font-mono">
                  Priority: {priorityScore.toFixed(1)}
                </span>
                {#if displayItem.isAssetGroup && (displayItem.affectedUrls?.length ?? 0) > 1}
                  <span class="bg-base-300 text-base-content/50 text-xs px-2.5 py-1 rounded-lg">
                    {displayItem.affectedUrls?.length ?? 0} pages affected
                  </span>
                {/if}
              </div>

              <!-- URL -->
              {#if displayItem.isAssetGroup && displayItem.assetUrl}
                <div class="mb-3">
                  <p class="text-xs font-semibold uppercase tracking-wider text-base-content/40 mb-1">Image URL</p>
                  <a href={displayItem.assetLinkUrl || displayItem.assetUrl} target="_blank" class="text-sm text-primary break-all hover:underline">
                    {displayItem.assetUrl}
                  </a>
                </div>
                <div class="mb-3">
                  <p class="text-xs font-semibold uppercase tracking-wider text-base-content/40 mb-1">Affected Pages ({displayItem.affectedUrls?.length ?? 0})</p>
                  <div class="flex flex-wrap gap-1.5">
                    {#each (displayItem.affectedUrls ?? []).slice(0, 5) as url}
                      <a href={url} target="_blank" rel="noopener noreferrer" class="bg-base-300 text-xs text-base-content/60 px-2.5 py-1 rounded-lg hover:text-primary truncate max-w-[200px]" title={url}>{url}</a>
                    {/each}
                    {#if (displayItem.affectedUrls?.length ?? 0) > 5}
                      <span class="text-xs text-base-content/30 px-2.5 py-1">+{(displayItem.affectedUrls?.length ?? 0) - 5} more</span>
                    {/if}
                  </div>
                </div>
              {:else}
                <div class="mb-3">
                  <p class="text-xs font-semibold uppercase tracking-wider text-base-content/40 mb-1">URL</p>
                  {#if displayItem.url}
                    <a href={displayItem.url} target="_blank" class="text-sm text-primary break-all hover:underline">{displayItem.url}</a>
                  {:else}
                    <span class="text-sm text-base-content/30 italic">URL not available</span>
                  {/if}
                </div>
                {#if displayItem.value}
                  <div class="mb-3">
                    <p class="text-xs font-semibold uppercase tracking-wider text-base-content/40 mb-1">Value</p>
                    <p class="text-sm text-base-content/70 break-all">{displayItem.value}</p>
                  </div>
                {/if}
              {/if}

              <!-- Recommendation -->
              {#if displayItem.recommendation}
                <div class="mb-3">
                  <p class="text-xs font-semibold uppercase tracking-wider text-base-content/40 mb-1">Recommendation</p>
                  <p class="text-sm text-base-content/70">{displayItem.recommendation}</p>
                </div>
              {/if}

              <!-- GSC Performance -->
              {#if enriched?.gsc_performance}
                <div class="bg-base-300 rounded-lg p-3 mb-3">
                  <p class="text-xs font-semibold uppercase tracking-wider text-base-content/40 mb-2">Google Search Console</p>
                  <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
                    <div>
                      <p class="text-xs text-base-content/40">Impressions</p>
                      <p class="font-bold font-mono text-sm">{formatNumber(enriched.gsc_performance.impressions)}</p>
                    </div>
                    <div>
                      <p class="text-xs text-base-content/40">Clicks</p>
                      <p class="font-bold font-mono text-sm">{formatNumber(enriched.gsc_performance.clicks)}</p>
                    </div>
                    <div>
                      <p class="text-xs text-base-content/40">CTR</p>
                      <p class="font-bold font-mono text-sm">{formatCTR(enriched.gsc_performance.ctr)}</p>
                    </div>
                    <div>
                      <p class="text-xs text-base-content/40">Position</p>
                      <p class="font-bold font-mono text-sm">{formatPosition(enriched.gsc_performance.position)}</p>
                    </div>
                  </div>
                  {#if enriched.gsc_performance.top_queries?.length > 0}
                    <div class="mt-3">
                      <p class="text-xs font-semibold uppercase tracking-wider text-base-content/40 mb-1">Top Queries</p>
                      <div class="flex flex-wrap gap-1.5">
                        {#each enriched.gsc_performance.top_queries.slice(0, 5) as query}
                          <span class="bg-base-200 text-xs text-base-content/60 px-2.5 py-1 rounded-lg">
                            {query.query} <span class="text-base-content/30">· {formatNumber(query.impressions)} imp</span>
                          </span>
                        {/each}
                      </div>
                    </div>
                  {/if}
                </div>
              {/if}

              <!-- GSC Insight -->
              {#if enriched?.recommendation_reason}
                <div class="bg-info/10 border border-info/20 rounded-lg px-3 py-2.5 mb-3">
                  <p class="text-xs font-semibold text-info mb-1">GSC Insight</p>
                  <p class="text-sm text-base-content/70">{enriched.recommendation_reason}</p>
                </div>
              {/if}

              <!-- AI Insight Button -->
              {#if crawlId}
                <div class="pt-1">
                  {#if isPro}
                    <button
                      class="btn btn-sm btn-ghost text-primary gap-1.5"
                      on:click|stopPropagation={() => {
                        selectedIssueForAI = displayItem.isAssetGroup
                          ? { ...displayItem, url: displayItem.affectedUrls?.[0], value: displayItem.assetUrl }
                          : displayItem;
                        showAIInsightModal = true;
                      }}
                    >
                      <Sparkles class="w-3.5 h-3.5" />
                      Generate AI Insight
                    </button>
                  {:else}
                    <div class="tooltip" data-tip="Upgrade to Pro for AI Insights">
                      <button class="btn btn-sm btn-ghost btn-disabled gap-1.5" disabled>
                        <Sparkles class="w-3.5 h-3.5" />
                        AI Insight
                        <span class="bg-primary/15 text-primary text-xs px-1.5 py-0.5 rounded font-medium">PRO</span>
                      </button>
                    </div>
                  {/if}
                </div>
              {/if}
            </div>
          {/if}
        </div>
      {/each}
    </div>
  {/each}

  <!-- AI Insight Modal -->
  {#if showAIInsightModal && selectedIssueForAI && crawlId}
    <IssueInsight
      issue={selectedIssueForAI}
      {crawlId}
      on:close={() => {
        showAIInsightModal = false;
        selectedIssueForAI = null;
      }}
    />
  {/if}

  {#if filteredIssues.length === 0}
    <div class="bg-success/10 border border-success/20 rounded-xl px-5 py-4 text-center">
      <p class="text-success font-medium">No issues found!</p>
    </div>
  {/if}
</div>
