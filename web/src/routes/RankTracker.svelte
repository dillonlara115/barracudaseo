<script>
  import { onMount } from 'svelte';
  import { params, link } from 'svelte-spa-router';
  import { fetchProjects, listKeywords, checkKeyword, getKeywordSnapshots, deleteKeyword, fetchProjectKeywordUsage, fetchProjectGSCStatus } from '../lib/data.js';
  import KeywordForm from '../components/KeywordForm.svelte';
  import KeywordDetailModal from '../components/KeywordDetailModal.svelte';
  import KeywordDiscovery from '../components/KeywordDiscovery.svelte';
  import ProjectPageLayout from '../components/ProjectPageLayout.svelte';
  import { ArrowUp, ArrowDown, Minus, Search, Filter, Plus, Sparkles, ScanSearch } from 'lucide-svelte';

  let projectId = null;
  let project = null;
  let keywords = [];
  let loading = true;
  let error = null;
  let successMessage = null;
  let showKeywordForm = false;
  let showKeywordDiscovery = false;
  let selectedKeyword = null;
  let showDetailModal = false;
  let checkingKeywords = new Set();
  let usageStats = null;
  let usageLoading = false;
  let gscStatus = null;
  
  // Filters
  let searchQuery = '';
  let deviceFilter = '';
  let locationFilter = '';
  let tagFilter = '';

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

      // Load keywords
      const filters = {};
      if (deviceFilter) filters.device = deviceFilter;
      if (locationFilter) filters.location = locationFilter;
      if (tagFilter) filters.tag = tagFilter;

      const result = await listKeywords(projectId, filters);
      if (result.error) {
        // If it's a "table not found" error or similar, treat as empty state
        const errorMsg = result.error.message || '';
        if (errorMsg.includes('table') || errorMsg.includes('relation') || errorMsg.includes('not found')) {
          // Table might not exist yet (migration not run) - treat as empty
          keywords = [];
          loading = false;
          return;
        }
        error = errorMsg || 'Failed to load keywords';
        loading = false;
        return;
      }
      
      // Debug: log the response structure
      console.log('Keywords API response:', result.data);
      
      keywords = result.data?.keywords || [];
      console.log('Parsed keywords:', keywords);
      loading = false;
      
      // Load usage stats
      await loadUsageStats();
    } catch (err) {
      error = err.message || 'Failed to load data';
      loading = false;
    }
  }

  async function loadUsageStats() {
    if (!projectId) return;
    
    usageLoading = true;
    try {
      const result = await fetchProjectKeywordUsage(projectId);
      if (!result.error && result.data) {
        usageStats = result.data;
      }
    } catch (err) {
      console.error('Failed to load usage stats:', err);
    } finally {
      usageLoading = false;
    }
  }

  async function handleCheckKeyword(keywordId) {
    if (checkingKeywords.has(keywordId)) return;
    
    checkingKeywords.add(keywordId);
    error = null; // Clear any previous errors
    successMessage = null;
    
    try {
      const result = await checkKeyword(keywordId);
      
      // Check for error in result.error or result.data.error (backend may return 200 with error message)
      const errorMsg = result.error?.message || result.data?.error || null;
      
      if (errorMsg) {
        // Check if this is a "not ranking" message (informational, not an error)
        if (errorMsg.includes('is not currently ranking') || errorMsg.includes('is not ranking')) {
          successMessage = errorMsg;
          setTimeout(() => {
            successMessage = null;
          }, 8000); // Show for 8 seconds since it's informational
          checkingKeywords.delete(keywordId);
          // Reload keywords to update last_checked_at
          await loadData();
          return;
        }
        
        // Actual error
        error = errorMsg;
        console.error('Keyword check failed:', errorMsg);
        
        // Show user-friendly message for common errors
        if (errorMsg.includes('DataForSEO') || errorMsg.includes('not configured')) {
          error = 'Rank tracking integration is not configured. Please contact support or check your settings.';
        }
        checkingKeywords.delete(keywordId);
        return;
      }
      
      // Check if task is still processing (202 Accepted)
      if (result.data?.status === 'processing' || result.data?.task_id) {
        console.log('Rank check initiated, task is processing:', result.data.task_id);
        successMessage = 'Rank check initiated. Polling for results...';
        
        // Poll for results every 3 seconds, up to 30 seconds
        let pollCount = 0;
        const maxPolls = 10; // 10 polls * 3 seconds = 30 seconds max
        const pollInterval = setInterval(async () => {
          pollCount++;
          
          // Reload keywords to check if snapshot was created
          await loadData();
          
          // Check if the keyword now has a position (snapshot was created)
          const keyword = keywords.find(k => k.id === keywordId);
          if (keyword && keyword.latest_position !== null && keyword.latest_position !== undefined) {
            clearInterval(pollInterval);
            checkingKeywords.delete(keywordId);
            successMessage = 'Rank check completed successfully!';
            setTimeout(() => {
              successMessage = null;
            }, 3000);
            return;
          }
          
          // Stop polling after max attempts
          if (pollCount >= maxPolls) {
            clearInterval(pollInterval);
            checkingKeywords.delete(keywordId);
            successMessage = 'Rank check is taking longer than expected. Results will appear when ready.';
            setTimeout(() => {
              successMessage = null;
            }, 5000);
          }
        }, 3000); // Poll every 3 seconds
      } else {
        // Task completed immediately
        successMessage = 'Rank check completed successfully!';
        setTimeout(() => {
          successMessage = null;
        }, 3000);
        checkingKeywords.delete(keywordId);
        
        // Reload keywords to get updated position
        await loadData();
      }
    } catch (err) {
      error = err.message || 'Failed to check keyword';
      console.error('Keyword check error:', err);
      checkingKeywords.delete(keywordId);
    }
  }

  async function handleDeleteKeyword(keywordId) {
    if (!confirm('Are you sure you want to delete this keyword?')) return;
    
    const result = await deleteKeyword(keywordId);
    if (result.error) {
      error = result.error.message || 'Failed to delete keyword';
      return;
    }
    
    await loadData();
  }

  function handleKeywordClick(keyword) {
    selectedKeyword = keyword;
    showDetailModal = true;
  }

  function handleKeywordCreated() {
    showKeywordForm = false;
    loadData();
  }

  function getTrendIcon(trend) {
    if (trend === 'up') return ArrowUp;
    if (trend === 'down') return ArrowDown;
    return Minus;
  }

  function getTrendColor(trend) {
    if (trend === 'up') return 'text-success';
    if (trend === 'down') return 'text-error';
    return 'text-base-content/60';
  }

  function formatDate(dateString) {
    if (!dateString) return '—';
    const date = new Date(dateString);
    return date.toLocaleDateString();
  }

  // Filter keywords by search query
  $: filteredKeywords = keywords.filter(k => {
    if (!searchQuery) return true;
    const query = searchQuery.toLowerCase();
    return k.keyword.toLowerCase().includes(query) ||
           (k.target_url && k.target_url.toLowerCase().includes(query));
  });

  // Get unique devices and locations for filters
  $: uniqueDevices = [...new Set(keywords.map(k => k.device))].sort();
  $: uniqueLocations = [...new Set(keywords.map(k => k.location_name))].sort();
  $: allTags = [...new Set(keywords.flatMap(k => k.tags || []))].sort();
</script>

<svelte:head>
  <title>Rank Tracker - {project?.name || 'Barracuda SEO'}</title>
</svelte:head>

<ProjectPageLayout {projectId} {gscStatus} showCrawlSection={false}>
<div class="max-w-7xl mx-auto">
  <!-- Header -->
  <div class="mb-6">
    <div class="flex items-center justify-between mb-4">
      <div>
        <h1 class="text-3xl font-bold mb-2">Rank Tracker</h1>
        <p class="text-base-content/70 mb-1">
          Track keyword rankings over time using our SERP integration. Monitor position changes, view historical data, and identify opportunities for improvement.
        </p>
        {#if project}
          <p class="text-sm text-base-content/60">Project: {project.name}</p>
        {/if}
      </div>
      <div class="flex gap-2">
        <a href="/project/{projectId}/discover-keywords" use:link class="btn btn-outline">
          <ScanSearch class="w-4 h-4 mr-1" />
          Discover Keywords
        </a>
        <button class="btn btn-primary" on:click={() => showKeywordForm = true}>
          <Plus class="w-4 h-4 mr-1" />
          Add Keyword
        </button>
      </div>
    </div>
  </div>

  {#if loading}
    <div class="flex justify-center items-center py-20">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
  {:else}
    {#if successMessage}
      <div class="alert {successMessage.includes('not currently ranking') || successMessage.includes('is not ranking') ? 'alert-info' : 'alert-success'} mb-4">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6">
          {#if successMessage.includes('not currently ranking') || successMessage.includes('is not ranking')}
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          {:else}
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          {/if}
        </svg>
        <span>{successMessage}</span>
      </div>
    {/if}
    {#if error && keywords.length > 0}
      <div class="alert alert-error mb-4">
        <span>{error}</span>
      </div>
    {/if}
    {#if error && keywords.length === 0}
      <div class="alert alert-warning mb-6">
        <span>{error}</span>
        <div class="mt-2 text-sm">
          <p>This might be because:</p>
          <ul class="list-disc list-inside mt-1">
            <li>The database migration hasn't been run yet</li>
            <li>No keywords have been added to this project</li>
          </ul>
          <p class="mt-2">Try adding a keyword to get started!</p>
        </div>
      </div>
    {/if}
    <!-- Filters -->
    <div class="card bg-base-100 shadow mb-6">
      <div class="card-body">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div class="form-control">
            <label class="label">
              <span class="label-text">Search</span>
            </label>
            <div class="relative">
              <Search class="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-base-content/40" />
              <input
                type="text"
                placeholder="Search keywords..."
                class="input input-bordered w-full pl-10"
                bind:value={searchQuery}
              />
            </div>
          </div>
          
          <div class="form-control">
            <label class="label">
              <span class="label-text">Device</span>
            </label>
            <select class="select select-bordered w-full" bind:value={deviceFilter}>
              <option value="">All Devices</option>
              {#each uniqueDevices as device}
                <option value={device}>{device}</option>
              {/each}
            </select>
          </div>
          
          <div class="form-control">
            <label class="label">
              <span class="label-text">Location</span>
            </label>
            <select class="select select-bordered w-full" bind:value={locationFilter}>
              <option value="">All Locations</option>
              {#each uniqueLocations as location}
                <option value={location}>{location}</option>
              {/each}
            </select>
          </div>
          
          <div class="form-control">
            <label class="label">
              <span class="label-text">Tag</span>
            </label>
            <select class="select select-bordered w-full" bind:value={tagFilter}>
              <option value="">All Tags</option>
              {#each allTags as tag}
                <option value={tag}>{tag}</option>
              {/each}
            </select>
          </div>
        </div>
        
        {#if deviceFilter || locationFilter || tagFilter}
          <div class="mt-4">
            <button class="btn btn-sm btn-ghost" on:click={() => { deviceFilter = ''; locationFilter = ''; tagFilter = ''; loadData(); }}>
              Clear Filters & Reload
            </button>
          </div>
        {/if}
      </div>
    </div>

    <!-- Keywords Table -->
    {#if filteredKeywords.length === 0}
      <div class="alert alert-info">
        <span>No keywords found. Click "Add Keyword" to start tracking rankings.</span>
      </div>
    {:else}
      <div class="card bg-base-100 shadow">
        <div class="card-body">
          <div class="overflow-x-auto">
            <table class="table table-zebra">
              <thead>
                <tr>
                  <th>Keyword</th>
                  <th>Target URL</th>
                  <th>Location</th>
                  <th>Device</th>
                  <th>Latest Position</th>
                  <th>Best Position</th>
                  <th>Trend</th>
                  <th>Check Frequency</th>
                  <th>Last Checked</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {#each filteredKeywords as keyword}
                  {@const TrendIcon = getTrendIcon(keyword.trend)}
                  <tr class="cursor-pointer hover:bg-base-200" on:click={() => handleKeywordClick(keyword)}>
                    <td class="font-medium">{keyword.keyword}</td>
                    <td class="max-w-xs truncate">
                      {#if keyword.target_url}
                        <a href={keyword.target_url} target="_blank" rel="noopener noreferrer" class="link link-primary" on:click|stopPropagation>
                          {keyword.target_url}
                        </a>
                      {:else}
                        <span class="text-base-content/40">—</span>
                      {/if}
                    </td>
                    <td>{keyword.location_name}</td>
                    <td class="capitalize">{keyword.device}</td>
                    <td>
                      {#if keyword.latest_position}
                        <span class="badge badge-lg">{keyword.latest_position}</span>
                      {:else}
                        <span class="text-base-content/40">—</span>
                      {/if}
                    </td>
                    <td>
                      {#if keyword.best_position}
                        <span class="badge badge-lg badge-success">{keyword.best_position}</span>
                      {:else}
                        <span class="text-base-content/40">—</span>
                      {/if}
                    </td>
                    <td>
                      {#if keyword.trend}
                        <TrendIcon class="w-5 h-5 {getTrendColor(keyword.trend)}" />
                      {:else}
                        <span class="text-base-content/40">—</span>
                      {/if}
                    </td>
                    <td>
                      <span class="badge badge-outline capitalize">
                        {keyword.check_frequency || 'manual'}
                      </span>
                    </td>
                    <td>{formatDate(keyword.last_checked)}</td>
                    <td>
                      <div class="flex gap-2" on:click|stopPropagation>
                        <button
                          class="btn btn-xs btn-outline"
                          on:click={() => handleCheckKeyword(keyword.id)}
                          disabled={checkingKeywords.has(keyword.id)}
                        >
                          {#if checkingKeywords.has(keyword.id)}
                            <span class="loading loading-spinner loading-xs"></span>
                          {:else}
                            Check Now
                          {/if}
                        </button>
                        <button
                          class="btn btn-xs btn-error btn-outline"
                          on:click={() => handleDeleteKeyword(keyword.id)}
                        >
                          Delete
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
  {/if}
</div>
</ProjectPageLayout>

{#if showKeywordForm}
  <KeywordForm
    projectId={projectId}
    on:close={() => showKeywordForm = false}
    on:created={handleKeywordCreated}
  />
{/if}

{#if showKeywordDiscovery}
  <KeywordDiscovery
    projectId={projectId}
    showAsModal={true}
    on:close={() => showKeywordDiscovery = false}
    on:keyword-added={() => loadData()}
    on:keywords-added={(e) => { loadData(); successMessage = `Added ${e.detail.count} keywords to tracking!`; setTimeout(() => successMessage = null, 3000); }}
  />
{/if}

{#if showDetailModal && selectedKeyword}
  <KeywordDetailModal
    keyword={selectedKeyword}
    on:close={() => { showDetailModal = false; selectedKeyword = null; }}
    on:checked={loadData}
  />
{/if}

