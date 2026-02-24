<script>
  import { onMount } from 'svelte';
  import { fetchQuickWins, fetchDecliningPages, streamAIRequest } from '../lib/data.js';

  export let projectId;

  let activeTab = 'quick-wins';
  let quickWins = [];
  let declining = [];
  let loading = false;
  let error = null;

  let explainText = '';
  let explainRow = null;
  let explaining = false;

  onMount(loadData);

  async function loadData() {
    loading = true;
    error = null;

    const [qw, dp] = await Promise.all([
      fetchQuickWins(projectId),
      fetchDecliningPages(projectId)
    ]);

    loading = false;

    if (qw.error) { error = qw.error.message; return; }
    if (dp.error) { error = dp.error.message; return; }

    quickWins = (qw.data || []).sort((a, b) => (b.opportunity_score || 0) - (a.opportunity_score || 0));
    declining = dp.data || [];
  }

  async function explainOpportunity(row) {
    explainRow = row;
    explainText = '';
    explaining = true;

    await streamAIRequest('/api/v1/ai/gsc/explain', {
      project_id: projectId,
      query: row.query,
      page: row.page
    }, {
      onChunk: (text) => { explainText += text; },
      onDone: () => { explaining = false; },
      onError: (err) => { explaining = false; explainText = `Error: ${err.message}`; }
    });
  }

  async function diagnoseDecline(row) {
    explainRow = row;
    explainText = '';
    explaining = true;

    await streamAIRequest('/api/v1/ai/gsc/diagnose', {
      project_id: projectId,
      page_url: row.page
    }, {
      onChunk: (text) => { explainText += text; },
      onDone: () => { explaining = false; },
      onError: (err) => { explaining = false; explainText = `Error: ${err.message}`; }
    });
  }

  function formatCTR(ctr) {
    return ctr ? `${(ctr * 100).toFixed(2)}%` : '—';
  }
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <h3 class="text-lg font-bold">GSC Intelligence</h3>
    <button class="btn btn-ghost btn-sm" on:click={loadData} disabled={loading}>
      {#if loading}<span class="loading loading-spinner loading-sm"></span>{/if}
      Refresh
    </button>
  </div>

  {#if error}
    <div class="alert alert-error text-sm">{error}</div>
  {/if}

  <div class="tabs tabs-boxed">
    <button class="tab" class:tab-active={activeTab === 'quick-wins'} on:click={() => activeTab = 'quick-wins'}>
      Quick Wins
    </button>
    <button class="tab" class:tab-active={activeTab === 'declining'} on:click={() => activeTab = 'declining'}>
      Declining Pages
    </button>
  </div>

  {#if loading}
    <div class="flex justify-center py-12">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
  {:else if activeTab === 'quick-wins'}
    <div class="overflow-x-auto">
      <table class="table table-sm">
        <thead>
          <tr>
            <th>Query</th>
            <th>Position</th>
            <th>Impressions</th>
            <th>CTR</th>
            <th>Score</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {#each quickWins as row}
            <tr>
              <td class="font-medium max-w-xs truncate">{row.query}</td>
              <td>{row.position?.toFixed(1)}</td>
              <td>{row.impressions?.toLocaleString()}</td>
              <td>{formatCTR(row.ctr)}</td>
              <td>
                <span class="badge badge-primary badge-sm">{row.opportunity_score?.toFixed(0)}</span>
              </td>
              <td>
                <button class="btn btn-ghost btn-xs" on:click={() => explainOpportunity(row)}>
                  Explain
                </button>
              </td>
            </tr>
          {:else}
            <tr><td colspan="6" class="text-center opacity-60 py-8">No quick wins found. Connect GSC and sync data first.</td></tr>
          {/each}
        </tbody>
      </table>
    </div>
  {:else if activeTab === 'declining'}
    <div class="overflow-x-auto">
      <table class="table table-sm">
        <thead>
          <tr>
            <th>Page</th>
            <th>Clicks</th>
            <th>Impressions</th>
            <th>CTR</th>
            <th>Position</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {#each declining as row}
            <tr>
              <td class="font-medium max-w-xs truncate">{row.page}</td>
              <td>{row.clicks?.toLocaleString()}</td>
              <td>{row.impressions?.toLocaleString()}</td>
              <td>{formatCTR(row.ctr)}</td>
              <td>{row.position?.toFixed(1)}</td>
              <td>
                <button class="btn btn-ghost btn-xs" on:click={() => diagnoseDecline(row)}>
                  Diagnose
                </button>
              </td>
            </tr>
          {:else}
            <tr><td colspan="6" class="text-center opacity-60 py-8">No declining pages detected.</td></tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}

  {#if explainRow}
    <div class="card bg-base-200 mt-4">
      <div class="card-body">
        <div class="flex items-center justify-between mb-2">
          <h4 class="card-title text-sm">AI Analysis: {explainRow.query || explainRow.page}</h4>
          <button class="btn btn-ghost btn-xs" on:click={() => { explainRow = null; explainText = ''; }}>Close</button>
        </div>
        {#if explaining}
          <span class="loading loading-dots loading-sm"></span>
        {/if}
        <div class="prose prose-sm max-w-none whitespace-pre-wrap">{explainText}</div>
      </div>
    </div>
  {/if}
</div>
