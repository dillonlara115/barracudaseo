<script>
  import { discoverKeywords, createKeyword } from '../lib/data.js';
  import { Search, Plus, Check } from 'lucide-svelte';
  import { createEventDispatcher } from 'svelte';

  export let projectId = null;
  export let defaultTarget = '';
  export let formData = {
    target: '',
    location_name: 'United States',
    language_name: 'English',
    limit: 1000,
    min_position: 0,
    max_position: 0
  };

  const dispatch = createEventDispatcher();

  let loading = false;
  let error = null;
  let discoveredKeywords = [];
  let selectedKeywords = new Set();
  let addingKeywords = new Set();
  let addedKeywords = new Set(); // Track successfully added keywords

  const commonLocations = [
    'United States',
    'United Kingdom',
    'Canada',
    'Australia',
    'Germany',
    'France',
    'Spain',
    'Italy',
    'Netherlands',
    'Sweden',
    'Denmark',
    'Norway',
    'Japan',
    'South Korea',
    'India',
    'Brazil',
    'Mexico',
    'Argentina',
    'Chile',
    'New Zealand'
  ];

  async function handleDiscover() {
    if (!formData.target.trim()) {
      error = 'Please enter a domain or URL';
      return;
    }

    loading = true;
    error = null;
    discoveredKeywords = [];
    selectedKeywords.clear();
    addedKeywords.clear(); // Clear added keywords when discovering new ones

    try {
      const result = await discoverKeywords(projectId, {
        target: formData.target,
        location_name: formData.location_name,
        language_name: formData.language_name,
        limit: formData.limit
      });

      if (result.error) {
        error = result.error.message || 'Failed to discover keywords';
        return;
      }

      let keywords = result.data?.keywords || [];
      
      // Filter by position if specified
      if (formData.min_position > 0) {
        keywords = keywords.filter(k => k.position >= formData.min_position);
      }
      if (formData.max_position > 0) {
        keywords = keywords.filter(k => k.position <= formData.max_position);
      }

      discoveredKeywords = keywords;
    } catch (err) {
      error = err.message || 'An error occurred while discovering keywords';
    } finally {
      loading = false;
    }
  }

  function toggleKeywordSelection(keyword) {
    if (selectedKeywords.has(keyword.keyword)) {
      selectedKeywords.delete(keyword.keyword);
    } else {
      selectedKeywords.add(keyword.keyword);
    }
    selectedKeywords = selectedKeywords;
  }

  function selectAll() {
    discoveredKeywords.forEach(k => selectedKeywords.add(k.keyword));
    selectedKeywords = selectedKeywords;
  }

  function deselectAll() {
    selectedKeywords.clear();
    selectedKeywords = selectedKeywords;
  }

  async function addKeywordToTracking(keyword) {
    if (addingKeywords.has(keyword.keyword) || addedKeywords.has(keyword.keyword)) return;
    
    addingKeywords.add(keyword.keyword);
    addingKeywords = addingKeywords;

    try {
      const result = await createKeyword(projectId, {
        keyword: keyword.keyword,
        target_url: keyword.url,
        location_name: formData.location_name,
        device: 'desktop'
      });

      if (result.error) {
        error = result.error.message || 'Failed to add keyword';
      } else {
        addedKeywords.add(keyword.keyword);
        addedKeywords = addedKeywords;
        dispatch('keyword-added', { keyword: keyword.keyword, data: result.data });
      }
    } catch (err) {
      error = err.message || 'An error occurred';
    } finally {
      addingKeywords.delete(keyword.keyword);
      addingKeywords = addingKeywords;
    }
  }

  async function addSelectedKeywords() {
    if (selectedKeywords.size === 0) return;

    const keywordsToAdd = discoveredKeywords.filter(k => selectedKeywords.has(k.keyword));
    let successCount = 0;

    for (const keyword of keywordsToAdd) {
      // Skip if already added
      if (addedKeywords.has(keyword.keyword)) continue;
      
      try {
        const result = await createKeyword(projectId, {
          keyword: keyword.keyword,
          target_url: keyword.url,
          location_name: formData.location_name,
          device: 'desktop'
        });

        if (!result.error) {
          addedKeywords.add(keyword.keyword);
          successCount++;
        }
      } catch (err) {
        console.error('Error adding keyword:', err);
      }
    }

    if (successCount > 0) {
      addedKeywords = addedKeywords; // Trigger reactivity
      dispatch('keywords-added', { count: successCount });
      selectedKeywords.clear();
      selectedKeywords = selectedKeywords;
    }
  }

  function getPositionBadgeClass(position) {
    if (position <= 3) return 'badge-success';
    if (position <= 10) return 'badge-warning';
    return 'badge-error';
  }

  function getCompetitionColor(competition) {
    if (!competition) return '';
    const comp = competition.toLowerCase();
    if (comp === 'low') return 'text-success';
    if (comp === 'medium') return 'text-warning';
    if (comp === 'high') return 'text-error';
    return '';
  }

  function formatNumber(num) {
    if (!num) return '—';
    return new Intl.NumberFormat().format(num);
  }

  // Initialize formData.target from defaultTarget
  $: if (defaultTarget && !formData.target) {
    formData.target = defaultTarget;
  }
</script>

<div class="w-full">
  <p class="text-sm text-base-content/70 mb-4">
    Discover keywords that your domain or specific URLs are currently ranking for. 
    Found keywords can be automatically linked to your crawled pages.
  </p>

  <!-- Discovery Form -->
  <div class="card bg-base-200 mb-4">
    <div class="card-body">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="form-control">
          <label class="label">
            <span class="label-text">Domain or URL</span>
          </label>
          <input
            type="text"
            placeholder="example.com or https://example.com/page"
            class="input input-bordered"
            bind:value={formData.target}
          />
        </div>

        <div class="form-control">
          <label class="label">
            <span class="label-text">Location</span>
          </label>
          <select class="select select-bordered" bind:value={formData.location_name}>
            {#each commonLocations as location}
              <option value={location}>{location}</option>
            {/each}
          </select>
        </div>

        <div class="form-control">
          <label class="label">
            <span class="label-text">Language</span>
          </label>
          <select class="select select-bordered" bind:value={formData.language_name}>
            <option value="English">English</option>
            <option value="Spanish">Spanish</option>
            <option value="French">French</option>
            <option value="German">German</option>
            <option value="Italian">Italian</option>
            <option value="Portuguese">Portuguese</option>
            <option value="Japanese">Japanese</option>
            <option value="Chinese">Chinese</option>
          </select>
        </div>

        <div class="form-control">
          <label class="label">
            <span class="label-text">Max Results</span>
          </label>
          <input
            type="number"
            min="1"
            max="10000"
            class="input input-bordered"
            bind:value={formData.limit}
          />
        </div>

        <div class="form-control">
          <label class="label">
            <span class="label-text">Min Position (optional)</span>
          </label>
          <input
            type="number"
            min="1"
            max="100"
            placeholder="e.g., 1"
            class="input input-bordered"
            bind:value={formData.min_position}
          />
        </div>

        <div class="form-control">
          <label class="label">
            <span class="label-text">Max Position (optional)</span>
          </label>
          <input
            type="number"
            min="1"
            max="100"
            placeholder="e.g., 20"
            class="input input-bordered"
            bind:value={formData.max_position}
          />
        </div>
      </div>

      <div class="mt-4">
        <button
          class="btn btn-primary"
          on:click={handleDiscover}
          disabled={loading}
        >
          {#if loading}
            <span class="loading loading-spinner loading-sm"></span>
            Discovering...
          {:else}
            <Search class="w-4 h-4 mr-2" />
            Discover Keywords
          {/if}
        </button>
      </div>
    </div>
  </div>

  {#if error}
    <div class="alert alert-error mb-4">
      <span>{error}</span>
    </div>
  {/if}

  {#if discoveredKeywords.length > 0}
    <!-- Results Header -->
    <div class="flex items-center justify-between mb-4">
      <div>
        <h4 class="font-bold">Found {discoveredKeywords.length} keywords</h4>
        {#if selectedKeywords.size > 0}
          <p class="text-sm text-base-content/70">
            {selectedKeywords.size} selected
          </p>
        {/if}
      </div>
      <div class="flex gap-2">
        {#if selectedKeywords.size > 0}
          <button
            class="btn btn-sm btn-primary"
            on:click={addSelectedKeywords}
            disabled={addingKeywords.size > 0}
          >
            <Plus class="w-4 h-4 mr-1" />
            Add {selectedKeywords.size} Selected
          </button>
          <button class="btn btn-sm btn-ghost" on:click={deselectAll}>
            Deselect All
          </button>
        {:else}
          <button class="btn btn-sm btn-ghost" on:click={selectAll}>
            Select All
          </button>
        {/if}
      </div>
    </div>

    <!-- Results Table -->
    <div class="card bg-base-100 shadow">
      <div class="card-body p-0">
        <div class="overflow-x-auto">
          <table class="table table-zebra">
            <thead>
              <tr>
                <th>
                  <input
                    type="checkbox"
                    class="checkbox checkbox-sm"
                    checked={selectedKeywords.size === discoveredKeywords.length && discoveredKeywords.length > 0}
                    on:change={(e) => e.target.checked ? selectAll() : deselectAll()}
                  />
                </th>
                <th>Keyword</th>
                <th>Position</th>
                <th>Search Volume</th>
                <th>Competition</th>
                <th>Difficulty</th>
                <th>URL</th>
                <th>Matched Page</th>
                <th>Action</th>
              </tr>
            </thead>
            <tbody>
              {#each discoveredKeywords as keyword}
                {@const isSelected = selectedKeywords.has(keyword.keyword)}
                {@const isAdding = addingKeywords.has(keyword.keyword)}
                {@const isAdded = addedKeywords.has(keyword.keyword)}
                <tr>
                  <td>
                    <input
                      type="checkbox"
                      class="checkbox checkbox-sm"
                      checked={isSelected}
                      on:change={() => toggleKeywordSelection(keyword)}
                    />
                  </td>
                  <td class="font-medium">{keyword.keyword}</td>
                  <td>
                    <span class="badge {getPositionBadgeClass(keyword.position)}">
                      {keyword.position}
                    </span>
                  </td>
                  <td>{formatNumber(keyword.search_volume)}</td>
                  <td>
                    <span class="{getCompetitionColor(keyword.competition)} capitalize">
                      {keyword.competition || '—'}
                    </span>
                  </td>
                  <td>{keyword.keyword_difficulty || '—'}</td>
                  <td class="max-w-xs truncate">
                    <a href={keyword.url} target="_blank" rel="noopener noreferrer" class="link link-primary">
                      {keyword.url}
                    </a>
                  </td>
                  <td>
                    {#if keyword.matched_page_url}
                      <span class="badge badge-success badge-sm">✓ Matched</span>
                    {:else}
                      <span class="text-base-content/40">—</span>
                    {/if}
                  </td>
                  <td>
                    <button
                      class="btn btn-xs {isAdded ? 'btn-success' : 'btn-primary'}"
                      on:click={() => addKeywordToTracking(keyword)}
                      disabled={isAdding || isAdded}
                    >
                      {#if isAdding}
                        <span class="loading loading-spinner loading-xs"></span>
                      {:else if isAdded}
                        <Check class="w-3 h-3" />
                      {:else}
                        <Plus class="w-3 h-3" />
                      {/if}
                    </button>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  {/if}
</div>

