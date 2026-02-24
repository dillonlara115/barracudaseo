<script>
  import { onMount } from 'svelte';
  import { fetchBriefs, updateBrief, streamAIRequest } from '../lib/data.js';

  export let projectId;

  let briefs = [];
  let loading = false;
  let error = null;

  let keyword = '';
  let generating = false;
  let streamText = '';

  let selectedBrief = null;

  onMount(loadBriefs);

  async function loadBriefs() {
    loading = true;
    const { data, error: err } = await fetchBriefs(projectId);
    loading = false;
    if (err) { error = err.message; return; }
    briefs = data || [];
  }

  async function generateBrief() {
    if (!keyword.trim()) return;
    generating = true;
    streamText = '';
    error = null;

    await streamAIRequest('/api/v1/ai/briefs/generate', {
      project_id: projectId,
      keyword: keyword.trim()
    }, {
      onChunk: (text) => { streamText += text; },
      onDone: async () => {
        generating = false;
        keyword = '';
        await loadBriefs();
      },
      onError: (err) => {
        generating = false;
        error = err.message;
      }
    });
  }

  async function updateStatus(brief, status) {
    await updateBrief({ id: brief.id, project_id: projectId, status });
    await loadBriefs();
  }

  function formatDate(d) {
    return new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
  }

  function statusBadge(status) {
    const map = { draft: 'badge-warning', approved: 'badge-success', archived: 'badge-ghost' };
    return map[status] || 'badge-ghost';
  }
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <h3 class="text-lg font-bold">Content Briefs</h3>
    <span class="text-sm opacity-60">{briefs.length} brief{briefs.length !== 1 ? 's' : ''}</span>
  </div>

  {#if error}
    <div class="alert alert-error text-sm">{error}</div>
  {/if}

  <div class="card bg-base-200">
    <div class="card-body">
      <h4 class="card-title text-sm">Generate New Brief</h4>
      <div class="flex gap-2">
        <input
          type="text"
          class="input input-bordered input-sm flex-1"
          placeholder="Enter target keyword..."
          bind:value={keyword}
          on:keydown={(e) => e.key === 'Enter' && generateBrief()}
          disabled={generating}
        />
        <button class="btn btn-primary btn-sm" on:click={generateBrief} disabled={generating || !keyword.trim()}>
          {#if generating}
            <span class="loading loading-spinner loading-sm"></span>
          {:else}
            Generate
          {/if}
        </button>
      </div>
    </div>
  </div>

  {#if generating}
    <div class="card bg-base-200">
      <div class="card-body">
        <p class="text-sm opacity-60 mb-2">Generating brief for "{keyword}"...</p>
        <div class="prose prose-sm max-w-none whitespace-pre-wrap">{streamText}</div>
      </div>
    </div>
  {/if}

  {#if loading}
    <div class="flex justify-center py-12">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
  {:else}
    <div class="space-y-3">
      {#each briefs as brief}
        <div class="card bg-base-200 cursor-pointer hover:bg-base-300 transition-colors" on:click={() => selectedBrief = selectedBrief?.id === brief.id ? null : brief}>
          <div class="card-body py-4">
            <div class="flex items-center justify-between">
              <div>
                <h4 class="font-semibold">{brief.keyword}</h4>
                <p class="text-xs opacity-60">{formatDate(brief.created_at)}</p>
              </div>
              <div class="flex items-center gap-2">
                <span class="badge {statusBadge(brief.status)} badge-sm">{brief.status}</span>
                {#if brief.status === 'draft'}
                  <button class="btn btn-success btn-xs" on:click|stopPropagation={() => updateStatus(brief, 'approved')}>Approve</button>
                {:else if brief.status === 'approved'}
                  <button class="btn btn-ghost btn-xs" on:click|stopPropagation={() => updateStatus(brief, 'archived')}>Archive</button>
                {/if}
              </div>
            </div>
            {#if selectedBrief?.id === brief.id && brief.brief_data}
              <div class="mt-4 pt-4 border-t border-base-300">
                <pre class="text-xs whitespace-pre-wrap overflow-auto max-h-96">{typeof brief.brief_data === 'string' ? brief.brief_data : JSON.stringify(brief.brief_data, null, 2)}</pre>
              </div>
            {/if}
          </div>
        </div>
      {:else}
        <div class="text-center py-8 opacity-60">No briefs yet. Generate your first brief above.</div>
      {/each}
    </div>
  {/if}
</div>
