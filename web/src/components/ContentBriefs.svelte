<script>
  import { onMount } from 'svelte';
  import { fetchBriefs, updateBrief, fetchKeywordSuggestions, streamAIRequest } from '../lib/data.js';
  import { marked } from 'marked';
  import { Sparkles, TrendingUp, Search, AlertCircle } from 'lucide-svelte';

  marked.setOptions({ breaks: true, gfm: true });

  export let projectId;
  export let initialKeyword = '';

  let briefs = [];
  let loading = false;
  let error = null;

  let keyword = '';
  let generating = false;
  let streamText = '';

  let selectedBrief = null;

  let suggestions = [];
  let suggestionsLoading = false;

  onMount(async () => {
    await Promise.all([loadBriefs(), loadSuggestions()]);
    if (initialKeyword) {
      keyword = initialKeyword;
      generateBrief(initialKeyword);
    }
  });

  async function loadBriefs() {
    loading = true;
    const { data, error: err } = await fetchBriefs(projectId);
    loading = false;
    if (err) { error = err.message; return; }
    briefs = data || [];
  }

  async function loadSuggestions() {
    suggestionsLoading = true;
    const { data, error: err } = await fetchKeywordSuggestions(projectId);
    suggestionsLoading = false;
    if (!err && data) {
      suggestions = data;
    }
  }

  async function generateBrief(targetKeyword) {
    const kw = targetKeyword || keyword.trim();
    if (!kw) return;
    generating = true;
    streamText = '';
    error = null;
    keyword = kw;

    await streamAIRequest('/api/v1/ai/briefs/generate', {
      project_id: projectId,
      keyword: kw
    }, {
      onChunk: (text) => { streamText += text; },
      onDone: async () => {
        generating = false;
        keyword = '';
        await Promise.all([loadBriefs(), loadSuggestions()]);
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

  function formatCTR(ctr) {
    return ctr ? `${(ctr * 100).toFixed(1)}%` : '—';
  }

  function reasonLabel(reason) {
    const map = {
      quick_win: { text: 'Quick Win', class: 'badge-success' },
      content_gap: { text: 'Content Gap', class: 'badge-info' },
      low_engagement: { text: 'Low Engagement', class: 'badge-warning' }
    };
    return map[reason] || { text: reason, class: 'badge-ghost' };
  }

  function reasonIcon(reason) {
    const map = {
      quick_win: TrendingUp,
      content_gap: Search,
      low_engagement: AlertCircle
    };
    return map[reason] || Sparkles;
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

  <!-- Keyword Suggestions from GSC -->
  {#if suggestions.length > 0}
    <div class="card bg-base-200">
      <div class="card-body">
        <div class="flex items-center gap-2 mb-3">
          <Sparkles class="w-4 h-4 text-primary" />
          <h4 class="card-title text-sm">Suggested Keywords</h4>
          <span class="text-xs opacity-50">Based on your GSC data</span>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2">
          {#each suggestions as suggestion}
            <button
              class="flex items-start gap-3 p-3 rounded-lg border border-base-300 hover:border-primary hover:bg-base-100 transition-colors text-left group"
              on:click={() => generateBrief(suggestion.query)}
              disabled={generating}
            >
              <div class="flex-1 min-w-0">
                <p class="font-medium text-sm truncate group-hover:text-primary transition-colors">{suggestion.query}</p>
                <div class="flex items-center gap-2 mt-1">
                  <span class="text-xs opacity-50">Pos {suggestion.position?.toFixed(1)}</span>
                  <span class="text-xs opacity-50">·</span>
                  <span class="text-xs opacity-50">{suggestion.impressions?.toLocaleString()} imp</span>
                  <span class="text-xs opacity-50">·</span>
                  <span class="text-xs opacity-50">{formatCTR(suggestion.ctr)}</span>
                </div>
              </div>
              <span class="badge {reasonLabel(suggestion.reason).class} badge-xs mt-1 shrink-0">
                {reasonLabel(suggestion.reason).text}
              </span>
            </button>
          {/each}
        </div>
      </div>
    </div>
  {:else if !suggestionsLoading && suggestions.length === 0}
    <!-- No suggestions available — don't show the section -->
  {/if}

  <!-- Manual Keyword Input -->
  <div class="card bg-base-200">
    <div class="card-body">
      <h4 class="card-title text-sm">Generate Brief</h4>
      <p class="text-xs opacity-50 -mt-2">Enter a keyword or pick one from the suggestions above.</p>
      <div class="flex gap-2">
        <input
          type="text"
          class="input input-bordered input-sm flex-1"
          placeholder="Enter target keyword..."
          bind:value={keyword}
          on:keydown={(e) => e.key === 'Enter' && generateBrief()}
          disabled={generating}
        />
        <button class="btn btn-primary btn-sm" on:click={() => generateBrief()} disabled={generating || !keyword.trim()}>
          {#if generating}
            <span class="loading loading-spinner loading-sm"></span>
          {:else}
            Generate
          {/if}
        </button>
      </div>
    </div>
  </div>

  <!-- Streaming Output -->
  {#if generating}
    <div class="card bg-base-200">
      <div class="card-body">
        <div class="flex items-center gap-2 mb-2">
          <span class="loading loading-spinner loading-sm text-primary"></span>
          <p class="text-sm opacity-60">Generating brief for "<strong>{keyword}</strong>"...</p>
        </div>
        <div class="ai-summary-content">{@html marked.parse(streamText || '')}</div>
      </div>
    </div>
  {/if}

  <!-- Brief List -->
  {#if loading}
    <div class="flex justify-center py-12">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
  {:else}
    <div class="space-y-3">
      {#each briefs as brief}
        <div
          class="card bg-base-200 cursor-pointer hover:bg-base-300 transition-colors"
          on:click={() => selectedBrief = selectedBrief?.id === brief.id ? null : brief}
          on:keydown={(e) => e.key === 'Enter' && (selectedBrief = selectedBrief?.id === brief.id ? null : brief)}
          role="button"
          tabindex="0"
        >
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
              <div class="mt-4 pt-4 border-t border-base-300 overflow-auto max-h-[32rem]">
                {#if typeof brief.brief_data === 'string'}
                  <div class="ai-summary-content">{@html marked.parse(brief.brief_data)}</div>
                {:else}
                  <pre class="text-xs whitespace-pre-wrap">{JSON.stringify(brief.brief_data, null, 2)}</pre>
                {/if}
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
