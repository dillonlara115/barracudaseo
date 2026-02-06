<script>
  import { createEventDispatcher, onMount } from 'svelte';
  import { link } from 'svelte-spa-router';
  import {
    fetchProjectGA4Connect,
    fetchProjectGA4Status,
    fetchProjectGA4Properties,
    updateProjectGA4Property,
    triggerProjectGA4Sync,
    disconnectProjectGA4
  } from '../lib/data.js';
  
  const dispatch = createEventDispatcher();
  
  export let project = null;
  export let projectId = null;

  let isConnected = false;
  let properties = [];
  let selectedPropertyId = null;
  let isLoadingProperties = false;
  let isSaving = false;
  let error = null;

  let ga4Status = null;
  let ga4Loading = false;
  let ga4Refreshing = false;
  let ga4Error = null;
  let lastProjectId = null;
  let propertySelectId = 'ga4-property-select';
  let isConnecting = false;
  let connectSuccess = false;
  let activePopup = null;

  const formatDateTime = (value) => {
    if (!value) return null;
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return null;
    return `${date.toLocaleDateString()} ${date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}`;
  };

  // Connected if integration exists with access token (property_id is optional)
  $: isConnected = Boolean(ga4Status?.integration?.access_token);
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

    if (typeof window !== 'undefined') {
      const handleMessage = async (event) => {
        // Only handle GA4-related messages
        if (!event.data?.type || !event.data.type.startsWith('ga4_')) {
          return;
        }
        
        // If project_id is specified, it must match (unless we're in a popup context)
        if (event.data?.project_id && event.data.project_id !== projectId) {
          console.log('Message project_id mismatch, ignoring', { 
            expected: projectId, 
            received: event.data.project_id,
            type: event.data.type 
          });
          return;
        }
        
        console.log('GA4 message received', { 
          type: event.data.type, 
          projectId, 
          receivedProjectId: event.data.project_id,
          hasActivePopup: !!activePopup
        });
        
        if (event.data?.type === 'ga4_connected') {
          console.log('GA4 connected message received');
          await initialize();
          connectSuccess = true;
          isConnecting = false;
          // Clear success message after 5 seconds
          setTimeout(() => {
            connectSuccess = false;
          }, 5000);
        } else if (event.data?.type === 'ga4_error') {
          console.error('GA4 error message received', event.data.error);
          error = event.data.error || 'Failed to connect to Google Analytics 4.';
          isConnecting = false;
          connectSuccess = false;
        }
      };

      window.addEventListener('message', handleMessage);
      return () => window.removeEventListener('message', handleMessage);
    }
  });

  async function initialize() {
    if (!projectId) return;
    try {
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
    if (!projectId) return;

    isLoadingProperties = true;
    error = null;

    const result = await fetchProjectGA4Properties(projectId);
    if (result.error) {
      error = result.error.message || 'Failed to load properties';
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

  async function connectGA4() {
    if (!projectId) {
      error = 'Project ID is required';
      return;
    }

    isConnecting = true;
    error = null;

    try {
      const result = await fetchProjectGA4Connect(projectId);
      if (result.error) {
        error = result.error.message || 'Failed to get authorization URL';
        isConnecting = false;
        return;
      }

      const { auth_url } = result.data;
      if (!auth_url) {
        error = 'No authorization URL returned';
        isConnecting = false;
        return;
      }

      // Open OAuth popup window
      const width = 600;
      const height = 700;
      const left = (window.screen.width - width) / 2;
      const top = (window.screen.height - height) / 2;
      
      const popup = window.open(
        auth_url,
        'ga4-oauth',
        `width=${width},height=${height},left=${left},top=${top},resizable=yes,scrollbars=yes`
      );

      if (!popup) {
        error = 'Popup blocked. Please allow popups for this site.';
        isConnecting = false;
        return;
      }

      // Store popup reference
      activePopup = popup;
      let messageReceived = false;

      // Listen for OAuth completion message from popup
      const messageHandler = async (event) => {
        // Only handle GA4 messages
        if (!event.data?.type || !event.data.type.startsWith('ga4_')) {
          return;
        }
        
        // Security: Verify origin matches expected API origin
        const expectedOrigin = window.location.origin;
        if (event.origin !== expectedOrigin && !event.origin.includes('localhost') && !event.origin.includes('127.0.0.1')) {
          console.warn('Message from unexpected origin:', event.origin);
          return;
        }
        
        console.log('OAuth popup message received', event.data);
        
        // If project_id is specified, it must match
        if (event.data?.project_id && event.data.project_id !== projectId) {
          console.log('Popup message project_id mismatch', { expected: projectId, received: event.data.project_id });
          return;
        }
        
        messageReceived = true;
        
        if (event.data?.type === 'ga4_connected') {
          console.log('GA4 connected via popup handler');
          // Clean up handlers immediately
          window.removeEventListener('message', messageHandler);
          if (checkClosed) {
            clearInterval(checkClosed);
          }
          
          // Close popup if still open
          if (popup && !popup.closed) {
            try {
              popup.close();
            } catch (e) {
              console.warn('Could not close popup:', e);
            }
          }
          activePopup = null;
          
          // Update UI
          await initialize();
          connectSuccess = true;
          isConnecting = false;
          // Clear success message after 5 seconds
          setTimeout(() => {
            connectSuccess = false;
          }, 5000);
        } else if (event.data?.type === 'ga4_error') {
          console.error('GA4 error via popup handler', event.data.error);
          // Clean up handlers immediately
          window.removeEventListener('message', messageHandler);
          if (checkClosed) {
            clearInterval(checkClosed);
          }
          
          // Close popup if still open
          if (popup && !popup.closed) {
            try {
              popup.close();
            } catch (e) {
              console.warn('Could not close popup:', e);
            }
          }
          activePopup = null;
          
          // Update UI
          error = event.data.error || 'Failed to connect Google Analytics 4';
          isConnecting = false;
          connectSuccess = false;
        }
      };

      // Set up message listener BEFORE popup navigates
      window.addEventListener('message', messageHandler);

      // Check if popup was closed manually (but give time for message to arrive)
      const checkClosed = setInterval(() => {
        if (popup.closed) {
          console.log('Popup closed', { messageReceived });
          clearInterval(checkClosed);
          window.removeEventListener('message', messageHandler);
          activePopup = null;
          
          // Only show error if we didn't receive a success/error message
          // The callback page should have sent a message before closing
          if (!messageReceived && isConnecting && !connectSuccess) {
            // Wait a bit more in case message is delayed
            setTimeout(() => {
              if (isConnecting && !connectSuccess) {
                isConnecting = false;
                error = 'Connection cancelled or popup was closed';
              }
            }, 500);
          }
        }
      }, 300);
    } catch (err) {
      error = err.message || 'Failed to connect Google Analytics 4';
      isConnecting = false;
    }
  }

  async function disconnectGA4() {
    if (!projectId) return;
    
    if (!confirm('Are you sure you want to disconnect Google Analytics 4? This will stop data synchronization.')) {
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
    isConnected = false;

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
    {#if connectSuccess}
      <div class="alert alert-success">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
        </svg>
        <span>Successfully connected to Google Analytics 4! Please wait while we load your properties...</span>
      </div>
    {:else}
      <div class="alert alert-info">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
        </svg>
        <div class="flex-1">
          <div class="font-semibold mb-1">Connect Google Analytics 4</div>
          <div class="text-sm">Connect your Google Analytics 4 account to enhance recommendations with real user behavior data.</div>
        </div>
      </div>
      <button
        class="btn btn-primary w-full"
        on:click={connectGA4}
        disabled={isConnecting || !projectId}
      >
        {#if isConnecting}
          <span class="loading loading-spinner loading-sm"></span>
          Connecting...
        {:else}
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5 mr-2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z" />
          </svg>
          Connect Google Analytics 4
        {/if}
      </button>
    {/if}
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
        Disconnect Google Analytics 4
      </button>
    </div>
  </div>
{/if}
