<script>
  import { createEventDispatcher, onMount } from 'svelte';
  import { link } from 'svelte-spa-router';
  import {
    fetchProjectGA4Status,
    updateProjectGA4Property,
    triggerProjectGA4Sync,
    disconnectProjectGA4,
    fetchGA4Status,
    fetchGA4Properties
  } from '../lib/data.js';
  
  const dispatch = createEventDispatcher();
  
  export let project = null;
  export let projectId = null;

  let isConnected = false;
  let globalConnected = false;
  let properties = [];
  let selectedPropertyId = null;
  let isLoadingProperties = false;
  let isSaving = false;
  let error = null;
  let errorActivationUrl = null;

  let ga4Status = null;
  let ga4Loading = false;
  let ga4Refreshing = false;
  let ga4Error = null;
  let lastProjectId = null;
  let propertySelectId = 'ga4-property-select';
  let loadingGlobalStatus = false;

  const formatDateTime = (value) => {
    if (!value) return null;
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return null;
    return `${date.toLocaleDateString()} ${date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}`;
  };

  // Parse Google API error to extract user-friendly message and activation URL
  function parseGoogleAPIError(errorMessage) {
    if (!errorMessage) return { message: null, activationUrl: null };
    
    // Try to extract activation URL from the error message
    const urlMatch = errorMessage.match(/https:\/\/console\.developers\.google\.com\/apis\/api\/[^\s\)]+/);
    const activationUrl = urlMatch ? urlMatch[0] : null;
    
    // Extract the user-friendly message (usually the LocalizedMessage)
    let message = errorMessage;
    
    // If it's a SERVICE_DISABLED error, provide a cleaner message
    // Note: This is an application configuration issue, not a user issue
    if (errorMessage.includes('SERVICE_DISABLED') || errorMessage.includes('has not been used') || errorMessage.includes('it is disabled')) {
      message = 'The Google Analytics Admin API is not enabled in the application\'s Google Cloud project. This is a configuration issue that needs to be fixed by the application administrator.';
    } else if (errorMessage.includes('accessNotConfigured')) {
      message = 'The Google Analytics Admin API is not enabled in the application\'s Google Cloud project. This is a configuration issue that needs to be fixed by the application administrator.';
    }
    
    return { message, activationUrl };
  }

  // Connected if a global integration exists
  $: isConnected = Boolean(globalConnected);
  $: hasPropertySelected = Boolean(ga4Status?.integration?.property_id);
  $: lastSyncedDisplay = ga4Status?.sync_state?.last_synced_at ? formatDateTime(ga4Status.sync_state.last_synced_at) : null;
  $: if (projectId && projectId !== lastProjectId && !ga4Loading) {
    lastProjectId = projectId;
    initialize();
  }
  $: propertySelectId = projectId ? `ga4-property-${projectId}` : 'ga4-property-select';

  onMount(() => {
    if (projectId) {
      lastProjectId = projectId;
    }

    initialize();
  });

  async function initialize() {
    if (!projectId) return;
    try {
      await loadGlobalStatus();
      await loadStatus();
      // Load properties if connected (even if no property selected yet)
      if (isConnected) {
        await loadProperties();
      } else {
        properties = [];
      }
    } catch (err) {
      console.error('Error initializing GA4 selector:', err);
      ga4Error = err.message || 'Failed to initialize Google Analytics 4 integration.';
      ga4Loading = false;
    }
  }

  async function loadStatus() {
    if (!projectId) return;

    ga4Loading = true;
    ga4Error = null;
    
    try {
      const statusResult = await fetchProjectGA4Status(projectId);

      if (statusResult.error) {
        ga4Status = null;
        ga4Error = statusResult.error.message || 'Unable to load Google Analytics 4 status.';
        ga4Loading = false;
        return;
      }

      ga4Status = statusResult.data || null;
      if (!selectedPropertyId && ga4Status?.integration?.property_id) {
        selectedPropertyId = ga4Status.integration.property_id;
      }
    } catch (err) {
      console.error('Error loading GA4 status:', err);
      ga4Status = null;
      ga4Error = err.message || 'Failed to load Google Analytics 4 status.';
    } finally {
      ga4Loading = false;
    }
  }

  async function loadProperties() {
    if (!projectId || !isConnected) return;

    isLoadingProperties = true;
    error = null;
    errorActivationUrl = null;

    const result = await fetchGA4Properties();
    if (result.error) {
      const errorMsg = result.error.message || 'Failed to load properties';
      const parsed = parseGoogleAPIError(errorMsg);
      error = parsed.message || errorMsg;
      errorActivationUrl = parsed.activationUrl;
      properties = [];
      isLoadingProperties = false;
      return;
    }

    const payload = result.data || {};
    properties = payload.properties || [];

    if (!selectedPropertyId) {
      if (payload.selectedProperty) {
        selectedPropertyId = payload.selectedProperty;
      } else if (properties.length > 0) {
        selectedPropertyId = properties[0].property_id;
      } else {
        selectedPropertyId = null;
      }
    }

    isLoadingProperties = false;
  }

  async function loadGlobalStatus() {
    loadingGlobalStatus = true;
    const statusResult = await fetchGA4Status();
    if (statusResult.error) {
      globalConnected = false;
      loadingGlobalStatus = false;
      return;
    }
    globalConnected = Boolean(statusResult.data?.connected);
    loadingGlobalStatus = false;
  }

  async function saveProperty() {
    if (!selectedPropertyId || !projectId) return;

    isSaving = true;
    error = null;

    const selectedProperty = properties.find(p => p.property_id === selectedPropertyId);
    const saveResult = await updateProjectGA4Property(
      projectId, 
      selectedPropertyId,
      selectedProperty?.property_name || selectedProperty?.display_name
    );
    if (saveResult.error) {
      error = saveResult.error.message || 'Failed to save property';
      isSaving = false;
      return;
    }

    await loadStatus();
    isSaving = false;
  }

  async function refreshGA4Data() {
    if (!projectId) return;
    ga4Refreshing = true;
    error = null;

    const result = await triggerProjectGA4Sync(projectId, { lookback_days: 30 });
    if (result.error) {
      error = result.error.message || 'Failed to refresh Google Analytics 4 data';
      ga4Refreshing = false;
      return;
    }

    await initialize();
    ga4Refreshing = false;
  }

  async function disconnectGA4() {
    if (!projectId) return;
    
    if (!confirm('Clear the Google Analytics 4 property selection for this project?')) {
      return;
    }

    isSaving = true;
    error = null;

    const result = await disconnectProjectGA4(projectId);
    if (result.error) {
      error = result.error.message || 'Failed to disconnect Google Analytics 4';
      isSaving = false;
      return;
    }

    // Reset local state
    ga4Status = null;
    properties = [];
    selectedPropertyId = null;
    isConnected = globalConnected;

    await initialize();
    isSaving = false;
  }
</script>

{#if ga4Loading}
  <div class="alert alert-info">
    <span>Loading Google Analytics 4 status...</span>
  </div>
{:else if ga4Error}
  <div class="alert alert-warning">
    <span>{ga4Error}</span>
  </div>
{:else if !isConnected}
  <div class="space-y-4">
    <div class="alert alert-info">
      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
      </svg>
      <div class="flex-1">
        <div class="font-semibold mb-1">Connect Google Analytics 4</div>
        <div class="text-sm">Connect a Google Analytics 4 account in Integrations to enable property selection.</div>
      </div>
    </div>
    <a href="/integrations" use:link class="btn btn-primary w-full">
      Go to Integrations
    </a>
    {#if error}
      <div class="alert alert-error">
        <span>{error}</span>
      </div>
    {/if}
  </div>
{:else}
  <div class="space-y-4">
    <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-3 rounded-box border border-base-300 bg-base-100 p-4 shadow-sm">
      <div>
        <div class="text-sm font-semibold text-base-content/80">Google Analytics 4</div>
        {#if hasPropertySelected}
          <div class="text-sm">
            Connected to <span class="font-semibold">{ga4Status?.integration?.property_name || ga4Status?.integration?.property_id}</span>.
          </div>
          {#if lastSyncedDisplay}
            <div class="text-xs text-base-content/60">Last synced {lastSyncedDisplay}</div>
          {:else}
            <div class="text-xs text-base-content/60">No cached metrics yet. Refresh to pull the latest data.</div>
          {/if}
        {:else}
          <div class="text-sm">
            Connected. Please select a property below.
          </div>
        {/if}
      </div>
      {#if hasPropertySelected}
        <div class="flex gap-2">
          {#if projectId}
            <a
              href="#/project/{projectId}?tab=ga4-dashboard"
              class="btn btn-sm btn-ghost"
            >
              View Dashboard
            </a>
          {/if}
          <button
            class="btn btn-sm btn-outline"
            on:click={refreshGA4Data}
            disabled={ga4Refreshing || ga4Loading}
          >
            {#if ga4Refreshing}
              <span class="loading loading-spinner loading-xs"></span>
              Refreshing...
            {:else}
              Refresh Data
            {/if}
          </button>
        </div>
      {/if}
    </div>

    {#if isLoadingProperties}
      <div class="alert alert-info">
        <span class="loading loading-spinner loading-sm"></span>
        <span>Loading properties...</span>
      </div>
    {:else if error}
      <div class="alert alert-warning">
        <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
        <div class="flex-1">
          <div class="font-semibold mb-2">{error}</div>
          {#if errorActivationUrl}
            <div class="text-sm mb-2">
              <strong>Note:</strong> This is an application configuration issue. The Google Analytics Admin API needs to be enabled in the application's Google Cloud project (the same project used for OAuth credentials), not in individual user accounts.
            </div>
            <div class="text-sm mb-2">
              If you are the application administrator, enable the API here:
            </div>
            <a 
              href={errorActivationUrl} 
              target="_blank" 
              rel="noopener noreferrer"
              class="btn btn-sm btn-primary"
            >
              Enable API in Google Cloud Console
            </a>
            <div class="text-xs mt-2 opacity-75">
              After enabling, wait a few minutes for the change to propagate, then refresh this page.
            </div>
          {:else}
            <div class="text-sm">{error}</div>
          {/if}
        </div>
      </div>
    {:else if properties.length === 0}
      <div class="alert alert-warning">
        <span>No Google Analytics 4 properties found. Confirm access in Google Analytics.</span>
      </div>
    {:else}
      <div class="form-control w-full">
        <label class="label" for={propertySelectId}>
          <span class="label-text">Google Analytics 4 Property</span>
        </label>
        <div class="flex gap-2">
          <select
            id={propertySelectId}
            class="select select-bordered flex-1"
            bind:value={selectedPropertyId}
            disabled={isLoadingProperties || isSaving}
          >
            {#each properties as prop}
              <option value={prop.property_id}>{prop.display_name || prop.property_name || prop.property_id}</option>
            {/each}
          </select>
          <button
            class="btn btn-primary"
            on:click={saveProperty}
            disabled={isSaving || isLoadingProperties || !selectedPropertyId}
          >
            {#if isSaving}
              <span class="loading loading-spinner loading-sm"></span>
              Saving...
            {:else}
              Save
            {/if}
          </button>
        </div>
        {#if ga4Status?.integration?.property_id}
          <div class="label">
            <span class="label-text-alt text-success">Property saved for this project</span>
          </div>
        {/if}
      </div>
    {/if}

    <!-- Disconnect Button -->
    <div class="pt-4 border-t border-base-200 mt-4">
      <h3 class="text-sm font-bold text-base-content/70 mb-2">Danger Zone</h3>
      <button
        class="btn btn-error btn-sm"
        on:click={disconnectGA4}
        disabled={isSaving || ga4Loading}
      >
        Clear Property Selection
      </button>
    </div>
  </div>
{/if}
