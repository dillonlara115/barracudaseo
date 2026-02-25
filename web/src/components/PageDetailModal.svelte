<script>
  import { createEventDispatcher } from 'svelte';
  import { collapseImageIssuesByAsset } from '../lib/issueUtils.js';

  export let page = null;
  export let issues = [];
  export let navigateToTab = null;

  $: displayItems = collapseImageIssuesByAsset(issues);

  const dispatch = createEventDispatcher();

  const close = () => {
    dispatch('close');
  };

  const handleViewIssues = () => {
    if (navigateToTab) {
      navigateToTab('issues', { url: page.url });
    }
    close();
  };

  const getSeverityBadge = (severity) => {
    switch (severity) {
      case 'error': return 'badge-error';
      case 'warning': return 'badge-warning';
      case 'info': return 'badge-info';
      default: return 'badge-ghost';
    }
  };

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

  const formatHeadingList = (headings) => {
    if (!headings || headings.length === 0) return 'None';
    return headings.join(', ');
  };
</script>

{#if page}
  <!-- Modal backdrop -->
  <div class="modal modal-open">
    <div class="modal-box max-w-4xl max-h-[90vh] overflow-y-auto">
      <h3 class="font-bold text-2xl mb-4">Page Details</h3>
      
      <!-- Page URL -->
      <div class="mb-4">
        <div class="font-semibold mb-2">URL:</div>
        <a href={page.url} target="_blank" class="link link-primary break-all">
          {page.url}
        </a>
      </div>

      <!-- Page Status & Performance -->
      <div class="grid grid-cols-2 gap-4 mb-4">
        <div>
          <div class="font-semibold mb-2">Status Code:</div>
          <span class="badge {page.status_code >= 200 && page.status_code < 300 ? 'badge-success' : page.status_code >= 400 ? 'badge-error' : 'badge-warning'}">
            {page.status_code}
          </span>
        </div>
        <div>
          <div class="font-semibold mb-2">Response Time:</div>
          <div>{page.response_time_ms}ms</div>
        </div>
      </div>

      <!-- Page Metadata -->
      <div class="mb-4">
        <div class="font-semibold mb-2">Title:</div>
        <div class="break-words">
          {#if page.title}
            {page.title}
          {:else}
            <span class="text-base-content/50">Not set</span>
          {/if}
        </div>
      </div>

      <div class="mb-4">
        <div class="font-semibold mb-2">Meta Description:</div>
        <div class="break-words">
          {#if page.meta_description}
            {page.meta_description}
          {:else}
            <span class="text-base-content/50">Not set</span>
          {/if}
        </div>
      </div>

      {#if page.canonical}
        <div class="mb-4">
          <div class="font-semibold mb-2">Canonical URL:</div>
          <a href={page.canonical} target="_blank" class="link link-primary break-all">
            {page.canonical}
          </a>
        </div>
      {/if}

      <!-- Headings -->
      <div class="mb-4">
        <div class="font-semibold mb-2">Headings:</div>
        <div class="space-y-2">
          {#if page.h1 && page.h1.length > 0}
            <div><span class="font-semibold">H1:</span> {formatHeadingList(page.h1)}</div>
          {/if}
          {#if page.h2 && page.h2.length > 0}
            <div><span class="font-semibold">H2:</span> {formatHeadingList(page.h2)}</div>
          {/if}
          {#if page.h3 && page.h3.length > 0}
            <div><span class="font-semibold">H3:</span> {formatHeadingList(page.h3)}</div>
          {/if}
          {#if (!page.h1 || page.h1.length === 0) && (!page.h2 || page.h2.length === 0) && (!page.h3 || page.h3.length === 0)}
            <div class="text-base-content/50">No headings found</div>
          {/if}
        </div>
      </div>

      <!-- Links -->
      <div class="grid grid-cols-2 gap-4 mb-4">
        <div>
          <div class="font-semibold mb-2">Internal Links:</div>
          <div class="badge badge-ghost">{page.internal_links?.length || 0}</div>
        </div>
        <div>
          <div class="font-semibold mb-2">External Links:</div>
          <div class="badge badge-ghost">{page.external_links?.length || 0}</div>
        </div>
      </div>

      <!-- Images -->
      {#if page.images && page.images.length > 0}
        <div class="mb-4">
          <div class="font-semibold mb-2">Images ({page.images.length}):</div>
          <div class="max-h-32 overflow-y-auto">
            {#each page.images.slice(0, 10) as image}
              <div class="text-sm break-all">
                {image.url}
                {#if image.alt}
                  <span class="text-base-content/70">(alt: {image.alt})</span>
                {:else}
                  <span class="badge badge-warning badge-sm">No alt</span>
                {/if}
              </div>
            {/each}
            {#if page.images.length > 10}
              <div class="text-sm text-base-content/70">... and {page.images.length - 10} more</div>
            {/if}
          </div>
        </div>
      {/if}

      <!-- Redirect Chain -->
      {#if page.redirect_chain && page.redirect_chain.length > 0}
        <div class="mb-4">
          <div class="font-semibold mb-2">Redirect Chain:</div>
          <div class="break-all">{page.redirect_chain.join(' → ')}</div>
        </div>
      {/if}

      <!-- Error -->
      {#if page.error}
        <div class="mb-4">
          <div class="font-semibold mb-2 text-error">Error:</div>
          <div class="text-error">{page.error}</div>
        </div>
      {/if}

      <!-- Issues Section -->
      <div class="divider"></div>
      <div class="mb-4">
        <div class="flex items-center justify-between mb-4">
          <h4 class="font-bold text-xl">Issues ({issues.length})</h4>
          {#if issues.length > 0}
            <button class="btn btn-primary btn-sm" on:click={handleViewIssues}>
              View All Issues
            </button>
          {/if}
        </div>

        {#if displayItems.length === 0}
          <div class="bg-success/10 border border-success/20 rounded-xl px-5 py-4 text-center">
            <p class="text-success font-medium">No issues found for this page!</p>
          </div>
        {:else}
          <div class="space-y-2">
            {#each displayItems as item}
              <div class="bg-base-200 rounded-xl px-4 py-3">
                <div class="flex items-center gap-3 mb-2">
                  <div class="w-2.5 h-2.5 rounded-full {getSeverityDot(item.severity)} shrink-0"></div>
                  <span class="text-sm text-base-content/80 flex-1">{item.message}</span>
                  <span class="text-xs font-medium whitespace-nowrap {getSeverityTextColor(item.severity)}">{getSeverityLabel(item.severity)}</span>
                </div>
                <div class="pl-5 space-y-1">
                  <span class="bg-base-300 text-base-content/60 text-xs px-2.5 py-1 rounded-lg inline-block">
                    {item.type.replace(/_/g, ' ')}
                  </span>
                  {#if item.isAssetGroup && item.assetUrl}
                    <div class="text-sm mt-1">
                      <span class="text-xs text-base-content/40">Image:</span>
                      <a href={item.assetLinkUrl || item.assetUrl} target="_blank" rel="noopener noreferrer" class="text-primary text-sm break-all hover:underline ml-1">{item.assetUrl}</a>
                    </div>
                    {#if (item.affectedUrls?.length ?? 0) > 1}
                      <p class="text-xs text-base-content/40">Used on {(item.affectedUrls?.length ?? 0)} places on this page</p>
                    {/if}
                  {:else if item.value}
                    <p class="text-sm text-base-content/60 break-all">{item.value}</p>
                  {/if}
                  {#if item.recommendation}
                    <p class="text-sm text-base-content/60 mt-1">{item.recommendation}</p>
                  {/if}
                </div>
              </div>
            {/each}
          </div>
        {/if}
      </div>

      <!-- Modal Actions -->
      <div class="modal-action">
        <button class="btn" on:click={close}>Close</button>
      </div>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button type="button" on:click={close}>close</button>
    </form>
  </div>
{/if}
