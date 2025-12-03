<script>
  import { onMount } from 'svelte';
  import { push } from 'svelte-spa-router';
  import { Search, BarChart3, Zap, Globe, Slack, Sparkles } from 'lucide-svelte';
  import { fetchProjects, saveOpenAIKey, getOpenAIKeyStatus, disconnectOpenAIKey } from '../lib/data.js';
  import ProjectGSCSelector from '../components/ProjectGSCSelector.svelte';
  
  let summary = null; // Could be passed as prop or fetched if needed
  let projects = [];
  let selectedProjectId = null;
  let selectedProject = null;
  let loadingProjects = false;
  let loadError = null;
  
  // OpenAI API Key state
  let openaiApiKey = '';
  let hasOpenAIKey = false;
  let loadingOpenAIStatus = false;
  let savingOpenAIKey = false;
  let disconnectingOpenAIKey = false;
  let openaiError = null;
  let openaiSuccess = false;
  
  // Integration statuses
  let integrations = {
    gsc: { connected: false, name: 'Google Search Console' },
    analytics: { connected: false, name: 'Google Analytics' },
    pagespeed: { connected: false, name: 'PageSpeed Insights' },
    bing: { connected: false, name: 'Bing Webmaster Tools' },
    slack: { connected: false, name: 'Slack' },
  };

  onMount(async () => {
    loadingProjects = true;
    const { data, error } = await fetchProjects();
    if (error) {
      loadError = error.message || 'Failed to load projects';
    } else {
      projects = data || [];
      if (projects.length > 0) {
        selectedProject = projects[0];
        selectedProjectId = selectedProject.id;
      }
    }
    loadingProjects = false;
    
    // Load OpenAI key status
    await loadOpenAIKeyStatus();
  });

  async function loadOpenAIKeyStatus() {
    loadingOpenAIStatus = true;
    console.log('Loading OpenAI key status...');
    const { data, error } = await getOpenAIKeyStatus();
    console.log('OpenAI key status response:', { data, error });
    console.log('Full data object:', JSON.stringify(data, null, 2));
    if (!error && data) {
      hasOpenAIKey = data.has_key || false;
      console.log('OpenAI key status loaded:', { hasOpenAIKey, has_key: data.has_key, dataKeys: Object.keys(data) });
    } else if (error) {
      console.error('Failed to load OpenAI key status:', error);
    }
    loadingOpenAIStatus = false;
  }

  async function handleSaveOpenAIKey() {
    savingOpenAIKey = true;
    openaiError = null;
    openaiSuccess = false;
    
    console.log('Saving OpenAI API key...', { hasKey: !!openaiApiKey, keyLength: openaiApiKey.length });
    
    const { data, error } = await saveOpenAIKey(openaiApiKey);
    console.log('Save OpenAI key response:', { data, error });
    console.log('Save response data:', JSON.stringify(data, null, 2));
    
    if (error) {
      console.error('Failed to save OpenAI key:', error);
      openaiError = error.message || 'Failed to save OpenAI API key';
    } else {
      console.log('OpenAI key saved successfully, waiting 500ms before reloading status...');
      openaiSuccess = true;
      openaiApiKey = ''; // Clear input after saving
      // Small delay to ensure database write is complete
      await new Promise(resolve => setTimeout(resolve, 500));
      // Reload status from server to ensure it's persisted
      await loadOpenAIKeyStatus();
      console.log('Status reloaded, hasOpenAIKey:', hasOpenAIKey);
      setTimeout(() => {
        openaiSuccess = false;
      }, 3000);
    }
    savingOpenAIKey = false;
  }

  async function handleDisconnectOpenAIKey() {
    if (!confirm('Are you sure you want to disconnect your OpenAI API key? You can reconnect it later.')) {
      return;
    }
    
    disconnectingOpenAIKey = true;
    openaiError = null;
    openaiSuccess = false;
    
    const { data, error } = await disconnectOpenAIKey();
    if (error) {
      openaiError = error.message || 'Failed to disconnect OpenAI API key';
    } else {
      openaiSuccess = true;
      openaiApiKey = ''; // Clear input
      // Reload status from server to ensure it's updated
      await loadOpenAIKeyStatus();
      setTimeout(() => {
        openaiSuccess = false;
      }, 3000);
    }
    disconnectingOpenAIKey = false;
  }

  $: if (selectedProjectId) {
    selectedProject = projects.find((p) => p.id === selectedProjectId) || selectedProject;
  }

  $: integrations.gsc.connected = Boolean(selectedProject?.settings?.gsc_property_url);
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
              Enhance recommendations with real search performance data. Prioritize fixes based on actual traffic.
            </p>
        </div>
        <div class="badge badge-success badge-lg" class:badge-success={integrations.gsc.connected} class:badge-ghost={!integrations.gsc.connected}>
          {integrations.gsc.connected ? 'Connected' : 'Available'}
        </div>
      </div>
        {#if loadError}
          <div class="alert alert-error mt-4">
            <span>{loadError}</span>
          </div>
        {:else if loadingProjects}
          <div class="alert alert-info mt-4">
            <span>Loading projects...</span>
          </div>
        {:else if projects.length === 0}
          <div class="alert alert-warning mt-4">
            <span>Create a project to connect Google Search Console.</span>
          </div>
        {:else}
          {#if projects.length > 1}
            <div class="form-control w-full mb-4">
              <label class="label" for="gsc-project-select">
                <span class="label-text">Project</span>
              </label>
              <select
                id="gsc-project-select"
                class="select select-bordered"
                bind:value={selectedProjectId}
              >
                {#each projects as projectOption}
                  <option value={projectOption.id}>{projectOption.name}</option>
                {/each}
              </select>
            </div>
          {/if}
          <ProjectGSCSelector summary={summary} project={selectedProject} projectId={selectedProjectId} />
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
              <label class="label">
                <span class="label-text">Your OpenAI API Key</span>
              </label>
              <input
                type="password"
                placeholder="sk-..."
                class="input input-bordered"
                bind:value={openaiApiKey}
                disabled={savingOpenAIKey}
              />
              <label class="label">
                <span class="label-text-alt">Your key is encrypted and stored securely. Leave empty to use the app-wide key.</span>
              </label>
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
              Connect your Google Analytics account to correlate SEO issues with user behavior and conversion data.
            </p>
          </div>
          <div class="badge badge-ghost badge-lg">Coming Soon</div>
        </div>
        <div class="alert alert-info">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          </svg>
          <span>This integration will allow you to see which SEO issues are affecting your most valuable pages based on conversion data.</span>
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
