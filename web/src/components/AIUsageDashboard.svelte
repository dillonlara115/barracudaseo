<script>
  import { onMount } from 'svelte';
  import { fetchAIUsage } from '../lib/data.js';

  export let projectId;

  let usageData = null;
  let loading = false;
  let error = null;

  onMount(loadUsage);

  async function loadUsage() {
    loading = true;
    const { data, error: err } = await fetchAIUsage(projectId);
    loading = false;
    if (err) { error = err.message; return; }
    usageData = data;
  }

  function pct(used, limit) {
    if (limit < 0) return 0;
    return Math.min(100, Math.round((used / limit) * 100));
  }

  function progressColor(percentage) {
    if (percentage >= 100) return 'progress-error';
    if (percentage >= 80) return 'progress-warning';
    return 'progress-primary';
  }

  const featureLabels = {
    brief: { label: 'Content Briefs', limitKey: 'briefs_per_month' },
    article: { label: 'Articles', limitKey: 'articles_per_month' },
    diagnostic: { label: 'AI Diagnostics', limitKey: 'diagnostics_per_month' }
  };
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <h3 class="text-lg font-bold">AI Usage</h3>
      <p class="text-sm opacity-60">Current billing period usage vs plan limits.</p>
    </div>
    {#if usageData}
      <span class="badge badge-primary">{usageData.plan} plan</span>
    {/if}
  </div>

  {#if error}
    <div class="alert alert-error text-sm">{error}</div>
  {/if}

  {#if loading}
    <div class="flex justify-center py-12">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
  {:else if usageData}
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      {#each Object.entries(featureLabels) as [feature, config]}
        {@const used = usageData.usage?.[feature] || 0}
        {@const limit = usageData.limits?.[config.limitKey] ?? -1}
        {@const percentage = pct(used, limit)}
        <div class="card bg-base-200">
          <div class="card-body">
            <h4 class="card-title text-sm">{config.label}</h4>
            <div class="flex items-end gap-2">
              <span class="text-2xl font-bold">{used}</span>
              <span class="text-sm opacity-60 mb-1">
                / {limit < 0 ? '∞' : limit}
              </span>
            </div>
            {#if limit >= 0}
              <progress class="progress {progressColor(percentage)} w-full" value={percentage} max="100"></progress>
              {#if percentage >= 100}
                <p class="text-xs text-error mt-1">Limit reached — upgrade for more</p>
              {:else if percentage >= 80}
                <p class="text-xs text-warning mt-1">Approaching limit</p>
              {/if}
            {:else}
              <p class="text-xs opacity-40 mt-1">Unlimited on your plan</p>
            {/if}
          </div>
        </div>
      {/each}
    </div>
  {:else}
    <div class="text-center py-8 opacity-60">No usage data available.</div>
  {/if}
</div>
