<script>
  import { onMount } from 'svelte';
  import { createEventDispatcher } from 'svelte';
  import { getKeywordSnapshots, checkKeyword } from '../lib/data.js';
  import { Line } from 'svelte-chartjs';
  import { Chart, CategoryScale, LinearScale, LineElement, PointElement, Title, Tooltip, Legend } from 'chart.js';
  import { X, ArrowUp, ArrowDown, Minus } from 'lucide-svelte';

  // Register Chart.js components
  Chart.register(CategoryScale, LinearScale, LineElement, PointElement, Title, Tooltip, Legend);

  export let keyword = null;

  const dispatch = createEventDispatcher();

  let loading = true;
  let error = null;
  let snapshots = [];
  let checking = false;
  let chartData = null;

  onMount(() => {
    loadSnapshots();
  });

  async function loadSnapshots() {
    if (!keyword?.id) return;

    loading = true;
    error = null;

    try {
      const result = await getKeywordSnapshots(keyword.id, 100);
      if (result.error) {
        error = result.error.message || 'Failed to load snapshots';
        loading = false;
        return;
      }

      snapshots = result.data?.snapshots || [];
      prepareChart();
      loading = false;
    } catch (err) {
      error = err.message || 'Failed to load snapshots';
      loading = false;
    }
  }

  async function handleCheckNow() {
    if (checking) return;

    checking = true;
    const result = await checkKeyword(keyword.id);
    checking = false;

    if (result.error) {
      error = result.error.message || 'Failed to check keyword';
      return;
    }

    // Reload snapshots
    await loadSnapshots();
    dispatch('checked');
  }

  function prepareChart() {
    if (!snapshots.length) {
      chartData = null;
      return;
    }

    // Sort snapshots by date (oldest first)
    const sorted = [...snapshots].sort((a, b) => {
      const dateA = new Date(a.checked_at);
      const dateB = new Date(b.checked_at);
      return dateA - dateB;
    });

    const labels = sorted.map(s => {
      const date = new Date(s.checked_at);
      return date.toLocaleDateString();
    });

    const positions = sorted.map(s => s.position_organic || s.position_absolute || null);

    chartData = {
      labels,
      datasets: [
        {
          label: 'Position',
          data: positions,
          borderColor: 'rgb(59, 130, 246)',
          backgroundColor: 'rgba(59, 130, 246, 0.1)',
          tension: 0.4,
          fill: true,
        },
      ],
    };
  }

  function formatDate(dateString) {
    if (!dateString) return '—';
    const date = new Date(dateString);
    return date.toLocaleString();
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
</script>

<div class="modal modal-open">
  <div class="modal-box max-w-4xl max-h-[90vh] overflow-y-auto">
    <div class="flex items-center justify-between mb-4">
      <h3 class="font-bold text-lg">Keyword Details</h3>
      <button class="btn btn-sm btn-circle btn-ghost" on:click={() => dispatch('close')}>
        <X class="w-4 h-4" />
      </button>
    </div>

    {#if keyword}
      <!-- Keyword Info -->
      <div class="card bg-base-200 mb-4">
        <div class="card-body">
          <h4 class="font-bold text-lg mb-2">{keyword.keyword}</h4>
          <div class="grid grid-cols-2 gap-4 text-sm">
            <div>
              <span class="text-base-content/70">Location:</span> {keyword.location_name}
            </div>
            <div>
              <span class="text-base-content/70">Device:</span> <span class="capitalize">{keyword.device}</span>
            </div>
            {#if keyword.target_url}
              <div class="col-span-2">
                <span class="text-base-content/70">Target URL:</span>
                <a href={keyword.target_url} target="_blank" rel="noopener noreferrer" class="link link-primary">
                  {keyword.target_url}
                </a>
              </div>
            {/if}
            {#if keyword.latest_position}
              <div>
                <span class="text-base-content/70">Latest Position:</span>
                <span class="badge badge-lg ml-2">{keyword.latest_position}</span>
              </div>
            {/if}
            {#if keyword.best_position}
              <div>
                <span class="text-base-content/70">Best Position:</span>
                <span class="badge badge-lg badge-success ml-2">{keyword.best_position}</span>
              </div>
            {/if}
            {#if keyword.trend}
              {@const TrendIcon = getTrendIcon(keyword.trend)}
              <div>
                <span class="text-base-content/70">Trend:</span>
                <TrendIcon class="w-5 h-5 inline-block ml-2 {getTrendColor(keyword.trend)}" />
              </div>
            {/if}
          </div>
        </div>
      </div>

      <!-- Actions -->
      <div class="mb-4">
        <button
          class="btn btn-primary"
          on:click={handleCheckNow}
          disabled={checking}
        >
          {#if checking}
            <span class="loading loading-spinner loading-sm"></span>
            Checking...
          {:else}
            Check Now
          {/if}
        </button>
      </div>

      {#if error}
        <div class="alert alert-error mb-4">
          <span>{error}</span>
        </div>
      {/if}

      {#if loading}
        <div class="flex justify-center items-center py-10">
          <span class="loading loading-spinner loading-lg"></span>
        </div>
      {:else if snapshots.length === 0}
        <div class="alert alert-info">
          <span>No rank data available yet. Click "Check Now" to get the first ranking.</span>
        </div>
      {:else}
        <!-- Chart -->
        {#if chartData}
          <div class="card bg-base-100 shadow mb-4">
            <div class="card-body">
              <h4 class="font-bold mb-4">Position Over Time</h4>
              <div class="h-64">
                <Line
                  data={chartData}
                  options={{
                    responsive: true,
                    maintainAspectRatio: false,
                    scales: {
                      y: {
                        reverse: true, // Lower position numbers are better
                        title: { display: true, text: 'Position' },
                        beginAtZero: false,
                      },
                    },
                    plugins: {
                      legend: { display: false },
                      tooltip: {
                        callbacks: {
                          label: (context) => `Position: ${context.parsed.y}`,
                        },
                      },
                    },
                  }}
                />
              </div>
            </div>
          </div>
        {/if}

        <!-- Snapshots Table -->
        <div class="card bg-base-100 shadow">
          <div class="card-body">
            <h4 class="font-bold mb-4">Historical Snapshots</h4>
            <div class="overflow-x-auto">
              <table class="table table-zebra">
                <thead>
                  <tr>
                    <th>Date</th>
                    <th>Position (Organic)</th>
                    <th>Position (Absolute)</th>
                    <th>URL</th>
                    <th>Title</th>
                  </tr>
                </thead>
                <tbody>
                  {#each snapshots as snapshot}
                    <tr>
                      <td>{formatDate(snapshot.checked_at)}</td>
                      <td>
                        {#if snapshot.position_organic}
                          <span class="badge">{snapshot.position_organic}</span>
                        {:else}
                          <span class="text-base-content/40">—</span>
                        {/if}
                      </td>
                      <td>
                        {#if snapshot.position_absolute}
                          <span class="badge badge-outline">{snapshot.position_absolute}</span>
                        {:else}
                          <span class="text-base-content/40">—</span>
                        {/if}
                      </td>
                      <td class="max-w-xs truncate">
                        {#if snapshot.serp_url}
                          <a href={snapshot.serp_url} target="_blank" rel="noopener noreferrer" class="link link-primary">
                            {snapshot.serp_url}
                          </a>
                        {:else}
                          <span class="text-base-content/40">—</span>
                        {/if}
                      </td>
                      <td class="max-w-xs truncate">
                        {#if snapshot.serp_title}
                          {snapshot.serp_title}
                        {:else}
                          <span class="text-base-content/40">—</span>
                        {/if}
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
  <div class="modal-backdrop" on:click={() => dispatch('close')}></div>
</div>

