<script>
  import { onMount } from 'svelte';
  import { push } from 'svelte-spa-router';
  import { Search, BarChart3, Zap, Globe, Slack, Sparkles, MousePointerClick } from 'lucide-svelte';
  import {
    saveOpenAIKey,
    getOpenAIKeyStatus,
    disconnectOpenAIKey,
    fetchGSCStatus,
    fetchGSCConnect,
    disconnectGSCIntegration,
    fetchGA4Status,
    fetchGA4Connect,
    disconnectGA4Integration
  } from '../lib/data.js';
  
  // OpenAI API Key state
  let openaiApiKey = '';
  let hasOpenAIKey = false;
  let loadingOpenAIStatus = false;
  let savingOpenAIKey = false;
  let disconnectingOpenAIKey = false;
  let openaiError = null;
  let openaiSuccess = false;
  
  // Global integrations state
  let gscConnected = false;
  let ga4Connected = false;
  let gscLoading = false;
  let ga4Loading = false;
  let gscError = null;
  let ga4Error = null;
  let gscConnecting = false;
  let ga4Connecting = false;
  let gscSuccess = false;
  let ga4Success = false;

  onMount(async () => {
    // Load OpenAI key status
    await loadOpenAIKeyStatus();
    await loadGSCStatus();
    await loadGA4Status();
  });

  async function loadOpenAIKeyStatus() {
    loadingOpenAIStatus = true;
    console.log('Loading OpenAI key status...');
    try {
      const { data, error } = await getOpenAIKeyStatus();
      console.log('OpenAI key status response:', { data, error });
      console.log('Full data object:', JSON.stringify(data, null, 2));
      if (error) {
        console.error('Failed to load OpenAI key status:', error);
        // Don't update hasOpenAIKey on error - keep current state
      } else if (data) {
        // Handle both direct boolean and object with has_key property
        const keyStatus = typeof data === 'boolean' ? data : (data.has_key === true);
        hasOpenAIKey = keyStatus;
        console.log('OpenAI key status loaded:', { hasOpenAIKey, rawData: data });
      } else {
        console.warn('OpenAI key status response has no data');
        hasOpenAIKey = false;
      }
    } catch (err) {
      console.error('Exception loading OpenAI key status:', err);
      // Don't update hasOpenAIKey on exception - keep current state
    } finally {
      loadingOpenAIStatus = false;
    }
  }

  async function handleSaveOpenAIKey() {
    savingOpenAIKey = true;
    openaiError = null;
    openaiSuccess = false;
    
    console.log('Saving OpenAI API key...', { hasKey: !!openaiApiKey, keyLength: openaiApiKey.length });
    
    try {
      const { data, error } = await saveOpenAIKey(openaiApiKey);
      console.log('Save OpenAI key response:', { data, error });
      console.log('Save response data:', JSON.stringify(data, null, 2));
      
      if (error) {
        console.error('Failed to save OpenAI key:', error);
        openaiError = error.message || 'Failed to save OpenAI API key';
        savingOpenAIKey = false;
        return;
      }
      
      // Success - check if response includes has_key
      console.log('OpenAI key saved successfully', data);
      openaiSuccess = true;
      openaiApiKey = ''; // Clear input after saving
      
      // Use has_key from response if available, otherwise optimistically set to true
      if (data && typeof data.has_key === 'boolean') {
        hasOpenAIKey = data.has_key;
        console.log('Status from save response:', hasOpenAIKey);
      } else {
        // Optimistically set status to true (will be verified by reload)
        hasOpenAIKey = true;
        console.log('No has_key in response, optimistically setting to true');
      }
      
      // Reload status with retry logic to ensure it's persisted
      let retries = 3;
      let statusLoaded = false;
      
      while (retries > 0 && !statusLoaded) {
        // Wait a bit longer for database write to be visible
        await new Promise(resolve => setTimeout(resolve, 800));
        
        await loadOpenAIKeyStatus();
        
        // If status confirms key exists, we're done
        if (hasOpenAIKey) {
          statusLoaded = true;
          console.log('Status confirmed, hasOpenAIKey:', hasOpenAIKey);
        } else {
          retries--;
          console.log(`Status check failed, retries remaining: ${retries}`);
          if (retries > 0) {
            console.log('Retrying status check...');
          }
        }
      }
      
      if (!statusLoaded) {
        console.warn('Status check failed after retries, but save was successful');
        // Keep hasOpenAIKey as true since save succeeded
      }
      
      setTimeout(() => {
        openaiSuccess = false;
      }, 3000);
    } catch (err) {
      console.error('Exception saving OpenAI key:', err);
      openaiError = err.message || 'Failed to save OpenAI API key';
    } finally {
      savingOpenAIKey = false;
    }
  }

  async function handleDisconnectOpenAIKey() {
    if (!confirm('Are you sure you want to disconnect your OpenAI API key? You can reconnect it later.')) {
      return;
    }
    
    disconnectingOpenAIKey = true;
    openaiError = null;
    openaiSuccess = false;
    
    try {
      const { data, error } = await disconnectOpenAIKey();
      if (error) {
        openaiError = error.message || 'Failed to disconnect OpenAI API key';
        disconnectingOpenAIKey = false;
        return;
      }
      
      // Success - optimistically update status
      openaiSuccess = true;
      openaiApiKey = ''; // Clear input
      hasOpenAIKey = false; // Set immediately
      
      // Reload status from server to confirm
      await new Promise(resolve => setTimeout(resolve, 500));
      await loadOpenAIKeyStatus();
      
      setTimeout(() => {
        openaiSuccess = false;
      }, 3000);
    } catch (err) {
      console.error('Exception disconnecting OpenAI key:', err);
      openaiError = err.message || 'Failed to disconnect OpenAI API key';
    } finally {
      disconnectingOpenAIKey = false;
    }
  }

  async function loadGSCStatus() {
    gscLoading = true;
    gscError = null;
    const { data, error } = await fetchGSCStatus();
    if (error) {
      gscError = error.message || 'Failed to load Google Search Console status';
      gscConnected = false;
    } else {
      gscConnected = Boolean(data?.connected);
    }
    gscLoading = false;
  }

  async function loadGA4Status() {
    ga4Loading = true;
    ga4Error = null;
    const { data, error } = await fetchGA4Status();
    if (error) {
      ga4Error = error.message || 'Failed to load Google Analytics 4 status';
      ga4Connected = false;
    } else {
      ga4Connected = Boolean(data?.connected);
    }
    ga4Loading = false;
  }

  function openOAuthPopup(authUrl, type) {
    const width = 600;
    const height = 700;
    const left = (window.screen.width - width) / 2;
    const top = (window.screen.height - height) / 2;

    const popup = window.open(
      authUrl,
      `${type}-oauth`,
      `width=${width},height=${height},left=${left},top=${top},resizable=yes,scrollbars=yes`
    );

    if (!popup) {
      return { error: 'Popup blocked. Please allow popups for this site.' };
    }

    const expectedOrigin = window.location.origin;

    const handler = (event) => {
      if (!event.data?.type || !event.data.type.startsWith(type)) {
        return;
      }
      if (event.origin !== expectedOrigin && !event.origin.includes('localhost') && !event.origin.includes('127.0.0.1')) {
        return;
      }

      if (event.data.type === `${type}_connected`) {
        if (type === 'gsc') {
          gscSuccess = true;
          gscConnecting = false;
          loadGSCStatus();
          setTimeout(() => (gscSuccess = false), 5000);
        } else if (type === 'ga4') {
          ga4Success = true;
          ga4Connecting = false;
          loadGA4Status();
          setTimeout(() => (ga4Success = false), 5000);
        }
      } else if (event.data.type === `${type}_error`) {
        if (type === 'gsc') {
          gscError = event.data.error || 'Failed to connect Google Search Console';
          gscConnecting = false;
        } else if (type === 'ga4') {
          ga4Error = event.data.error || 'Failed to connect Google Analytics 4';
          ga4Connecting = false;
        }
      }

      window.removeEventListener('message', handler);
    };

    window.addEventListener('message', handler);
    return { popup };
  }

  async function connectGSC() {
    gscConnecting = true;
    gscError = null;
    const result = await fetchGSCConnect();
    if (result.error) {
      gscError = result.error.message || 'Failed to get authorization URL';
      gscConnecting = false;
      return;
    }
    const { auth_url } = result.data || {};
    if (!auth_url) {
      gscError = 'No authorization URL returned';
      gscConnecting = false;
      return;
    }
    const popupResult = openOAuthPopup(auth_url, 'gsc');
    if (popupResult.error) {
      gscError = popupResult.error;
      gscConnecting = false;
    }
  }

  async function connectGA4() {
    ga4Connecting = true;
    ga4Error = null;
    const result = await fetchGA4Connect();
    if (result.error) {
      ga4Error = result.error.message || 'Failed to get authorization URL';
      ga4Connecting = false;
      return;
    }
    const { auth_url } = result.data || {};
    if (!auth_url) {
      ga4Error = 'No authorization URL returned';
      ga4Connecting = false;
      return;
    }
    const popupResult = openOAuthPopup(auth_url, 'ga4');
    if (popupResult.error) {
      ga4Error = popupResult.error;
      ga4Connecting = false;
    }
  }

  async function disconnectGSC() {
    if (!confirm('Disconnect Google Search Console for your account?')) return;
    const result = await disconnectGSCIntegration();
    if (result.error) {
      gscError = result.error.message || 'Failed to disconnect Google Search Console';
      return;
    }
    gscConnected = false;
    await loadGSCStatus();
  }

  async function disconnectGA4() {
    if (!confirm('Disconnect Google Analytics 4 for your account?')) return;
    const result = await disconnectGA4Integration();
    if (result.error) {
      ga4Error = result.error.message || 'Failed to disconnect Google Analytics 4';
      return;
    }
    ga4Connected = false;
    await loadGA4Status();
  }
</script>

<div class="container mx-auto p-6 max-w-4xl">
  <div class="mb-6">
    <button 
      class="btn btn-ghost btn-sm mb-4"
      on:click={() => push('/')}
    >
      ← Back to Projects
    </button>
    <h1 class="text-3xl font-bold mb-2">Integrations</h1>
    <p class="text-base-content/70">
      Connect external services to enhance your SEO workflow and get more insights.
    </p>
  </div>

  <div class="space-y-6">
    <!-- Google Search Console -->
    <div class="card bg-base-100 shadow">
      <div class="card-body">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h2 class="card-title text-xl">
              <Search class="w-6 h-6 mr-2" />
              Google Search Console
            </h2>
            <p class="text-sm text-base-content/70 mt-1">
              Connect your Google account once. Then select the appropriate property per project in Project Settings.
            </p>
          </div>
          <div class="badge badge-lg whitespace-nowrap" class:badge-success={gscConnected} class:badge-ghost={!gscConnected}>
            {gscConnected ? 'Connected' : 'Not Connected'}
          </div>
        </div>

        {#if gscLoading}
          <div class="alert alert-info">
            <span class="loading loading-spinner loading-sm"></span>
            <span>Loading status...</span>
          </div>
        {:else}
          <div class="space-y-4">
            {#if gscSuccess}
              <div class="alert alert-success">
                <span>Google Search Console connected successfully.</span>
              </div>
            {/if}
            {#if gscConnected}
              <div class="alert alert-success">
                <span>Your Google Search Console account is connected.</span>
              </div>
              <button class="btn btn-error" on:click={disconnectGSC}>
                Disconnect Google Search Console
              </button>
            {:else}
              <div class="alert alert-info">
                <span>Connect your Google Search Console account to enable property selection per project.</span>
              </div>
              <button class="btn btn-primary" on:click={connectGSC} disabled={gscConnecting}>
                {#if gscConnecting}
                  <span class="loading loading-spinner loading-sm"></span>
                  Connecting...
                {:else}
                  Connect Google Search Console
                {/if}
              </button>
            {/if}

            {#if gscError}
              <div class="alert alert-error">
                <span>{gscError}</span>
              </div>
            {/if}
          </div>
        {/if}
      </div>
    </div>

    <!-- OpenAI API Key -->
    <div class="card bg-base-100 shadow">
      <div class="card-body">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h2 class="card-title text-xl">
              <Sparkles class="w-6 h-6 mr-2" />
              OpenAI API Key
            </h2>
            <p class="text-sm text-base-content/70 mt-1">
              Connect your own OpenAI API key to use AI features. If provided, Barracuda will use YOUR OpenAI key for AI features to reduce cost and increase privacy.
            </p>
          </div>
          <div class="badge badge-lg whitespace-nowrap" class:badge-success={hasOpenAIKey} class:badge-ghost={!hasOpenAIKey}>
            {hasOpenAIKey ? 'Connected' : 'Not Connected'}
          </div>
        </div>
        
        {#if loadingOpenAIStatus}
          <div class="alert alert-info">
            <span class="loading loading-spinner loading-sm"></span>
            <span>Loading status...</span>
          </div>
        {:else}
          <div class="space-y-4">
            <div class="form-control">
              <label class="label" for="openai-api-key">
                <span class="label-text">Your OpenAI API Key</span>
              </label>
              <input
                id="openai-api-key"
                type="password"
                placeholder="sk-..."
                class="input input-bordered"
                bind:value={openaiApiKey}
                disabled={savingOpenAIKey}
              />
              <div class="label">
                <span class="label-text-alt">Your key is encrypted and stored securely. Leave empty to use the app-wide key.</span>
              </div>
            </div>
            
            {#if openaiError}
              <div class="alert alert-error">
                <span>{openaiError}</span>
              </div>
            {/if}
            
            {#if openaiSuccess}
              <div class="alert alert-success">
                <span>OpenAI API key saved successfully!</span>
              </div>
            {/if}
            
            <div class="alert alert-info">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
              <span>If provided, Barracuda will use YOUR OpenAI key for AI features to reduce cost and increase privacy. If not provided, the app-wide key will be used.</span>
            </div>
            
            <div class="flex gap-2">
              {#if hasOpenAIKey}
                <button
                  class="btn btn-error"
                  on:click={handleDisconnectOpenAIKey}
                  disabled={disconnectingOpenAIKey}
                >
                  {#if disconnectingOpenAIKey}
                    <span class="loading loading-spinner loading-sm"></span>
                    Disconnecting...
                  {:else}
                    Disconnect OpenAI API Key
                  {/if}
                </button>
              {/if}
              
              <button
                class="btn btn-primary"
                on:click={handleSaveOpenAIKey}
                disabled={savingOpenAIKey || !openaiApiKey.trim()}
              >
                {#if savingOpenAIKey}
                  <span class="loading loading-spinner loading-sm"></span>
                  Saving...
                {:else}
                  {hasOpenAIKey ? 'Update OpenAI API Key' : 'Save OpenAI API Key'}
                {/if}
              </button>
            </div>
          </div>
        {/if}
      </div>
    </div>

    <!-- Google Analytics -->
    <div class="card bg-base-100 shadow">
      <div class="card-body">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h2 class="card-title text-xl">
              <BarChart3 class="w-6 h-6 mr-2" />
              Google Analytics
            </h2>
            <p class="text-sm text-base-content/70 mt-1">
              Connect your Google account once. Then select the appropriate GA4 property per project in Project Settings.
            </p>
          </div>
          <div class="badge badge-lg whitespace-nowrap" class:badge-success={ga4Connected} class:badge-ghost={!ga4Connected}>
            {ga4Connected ? 'Connected' : 'Not Connected'}
          </div>
        </div>
        {#if ga4Loading}
          <div class="alert alert-info">
            <span class="loading loading-spinner loading-sm"></span>
            <span>Loading status...</span>
          </div>
        {:else}
          <div class="space-y-4">
            {#if ga4Success}
              <div class="alert alert-success">
                <span>Google Analytics 4 connected successfully.</span>
              </div>
            {/if}
            {#if ga4Connected}
              <div class="alert alert-success">
                <span>Your Google Analytics 4 account is connected.</span>
              </div>
              <button class="btn btn-error" on:click={disconnectGA4}>
                Disconnect Google Analytics 4
              </button>
            {:else}
              <div class="alert alert-info">
                <span>Connect your Google Analytics 4 account to enable property selection per project.</span>
              </div>
              <button class="btn btn-primary" on:click={connectGA4} disabled={ga4Connecting}>
                {#if ga4Connecting}
                  <span class="loading loading-spinner loading-sm"></span>
                  Connecting...
                {:else}
                  Connect Google Analytics 4
                {/if}
              </button>
            {/if}

            {#if ga4Error}
              <div class="alert alert-error">
                <span>{ga4Error}</span>
              </div>
            {/if}
          </div>
        {/if}
      </div>
    </div>

    <!-- Microsoft Clarity -->
    <div class="card bg-base-100 shadow">
      <div class="card-body">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h2 class="card-title text-xl">
              <MousePointerClick class="w-6 h-6 mr-2" />
              Microsoft Clarity
            </h2>
            <p class="text-sm text-base-content/70 mt-1">
              Connect Clarity for UX engagement metrics—rage clicks, dead clicks, scroll depth—to prioritize SEO fixes. Each project can connect to a different Clarity project or account.
            </p>
          </div>
          <div class="badge badge-ghost badge-lg">Per Project</div>
        </div>
        <div class="alert alert-info">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          </svg>
          <div>
            <p class="font-medium mb-1">Configure in Project Settings</p>
            <p class="text-sm">
              Clarity uses per-project API tokens (no OAuth). In each project's Settings, add your Clarity Project ID and API Token from <strong>Settings → Data Export</strong> in the Clarity dashboard. Data covers the last 1–3 days with a 10 requests/project/day limit.
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- PageSpeed Insights -->
    <div class="card bg-base-100 shadow">
      <div class="card-body">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h2 class="card-title text-xl">
              <Zap class="w-6 h-6 mr-2" />
              PageSpeed Insights
            </h2>
            <p class="text-sm text-base-content/70 mt-1">
              Automatically test page performance and get Core Web Vitals scores for crawled pages.
            </p>
          </div>
          <div class="badge badge-ghost badge-lg">Coming Soon</div>
        </div>
        <div class="alert alert-info">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          </svg>
          <span>Automatically test page speed and Core Web Vitals for all crawled pages to identify performance issues.</span>
        </div>
      </div>
    </div>

    <!-- Bing Webmaster Tools -->
    <div class="card bg-base-100 shadow">
      <div class="card-body">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h2 class="card-title text-xl">
              <Globe class="w-6 h-6 mr-2" />
              Bing Webmaster Tools
            </h2>
            <p class="text-sm text-base-content/70 mt-1">
              Import search performance data from Bing to get a complete picture of your SEO performance.
            </p>
          </div>
          <div class="badge badge-ghost badge-lg">Coming Soon</div>
        </div>
        <div class="alert alert-info">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          </svg>
          <span>Combine data from both Google and Bing to understand your search visibility across major search engines.</span>
        </div>
      </div>
    </div>

    <!-- Slack -->
    <div class="card bg-base-100 shadow">
      <div class="card-body">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h2 class="card-title text-xl">
              <Slack class="w-6 h-6 mr-2" />
              Slack
            </h2>
            <p class="text-sm text-base-content/70 mt-1">
              Get notified in Slack when critical SEO issues are found or when crawls complete.
            </p>
          </div>
          <div class="badge badge-ghost badge-lg">Coming Soon</div>
        </div>
        <div class="alert alert-info">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          </svg>
          <span>Receive real-time notifications about critical SEO issues and crawl status updates directly in your Slack workspace.</span>
        </div>
      </div>
    </div>
  </div>
</div>
