<script>
  import { onMount } from 'svelte';
  import { push } from 'svelte-spa-router';
  import { Search, BarChart3, Zap, Globe, Slack, FileText } from 'lucide-svelte';
  import { fetchProjects } from '../lib/data.js';
  import ProjectGSCSelector from '../components/ProjectGSCSelector.svelte';
  
  let summary = null; // Could be passed as prop or fetched if needed
  let projects = [];
  let selectedProjectId = null;
  let selectedProject = null;
  let loadingProjects = false;
  let loadError = null;
  
  // Integration statuses
  let integrations = {
    gsc: { connected: false, name: 'Google Search Console' },
    analytics: { connected: false, name: 'Google Analytics' },
    pagespeed: { connected: false, name: 'PageSpeed Insights' },
    bing: { connected: false, name: 'Bing Webmaster Tools' },
    slack: { connected: false, name: 'Slack' },
    jira: { connected: false, name: 'Jira' },
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
  });

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

    <!-- Jira -->
    <div class="card bg-base-100 shadow">
      <div class="card-body">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h2 class="card-title text-xl">
              <FileText class="w-6 h-6 mr-2" />
              Jira
            </h2>
            <p class="text-sm text-base-content/70 mt-1">
              Automatically create Jira tickets for SEO issues that need to be fixed by your development team.
            </p>
          </div>
          <div class="badge badge-ghost badge-lg">Coming Soon</div>
        </div>
        <div class="alert alert-info">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          </svg>
          <span>Streamline your workflow by automatically creating tickets for SEO issues, making it easy to track and assign fixes to your team.</span>
        </div>
      </div>
    </div>
  </div>
</div>
