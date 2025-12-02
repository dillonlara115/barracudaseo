<script>
  import { onMount } from 'svelte';
  import { push, params } from 'svelte-spa-router';
  import { fetchProjects, fetchCrawls, fetchCrawl, fetchProjectGSCStatus } from '../lib/data.js';
  import ProjectsView from './ProjectsView.svelte';
  import CrawlSelector from './CrawlSelector.svelte';
  import TriggerCrawlButton from './TriggerCrawlButton.svelte';
  import ProjectLayout from './ProjectLayout.svelte';

  export let projectId = null;
  export let gscStatus = null; // Optional: can be passed in, otherwise will be loaded
  export let showCrawlSection = false; // Default to false - crawl section removed from most pages

  let projects = [];
  let project = null;
  let selectedProject = null;
  let crawls = [];
  let loading = true;
  let error = null;
  let showActiveCrawlNotification = false;
  let activeCrawl = null;
  let gscStatusLoaded = null;

  $: currentProjectId = projectId;

  onMount(async () => {
    if (projectId) {
      await loadData();
      await checkActiveCrawl();
      await loadGSCStatus();
    }
  });

  $: if (projectId && projectId !== currentProjectId) {
    loadData();
    checkActiveCrawl();
    loadGSCStatus();
  }

  async function loadGSCStatus() {
    if (!projectId || gscStatus !== null) return; // Don't load if already provided
    try {
      const result = await fetchProjectGSCStatus(projectId);
      if (!result.error && result.data) {
        gscStatusLoaded = result.data;
      }
    } catch (err) {
      console.error('Failed to load GSC status:', err);
    }
  }

  // Use provided gscStatus or loaded one
  $: finalGSCStatus = gscStatus !== null ? gscStatus : gscStatusLoaded;

  async function checkActiveCrawl() {
    if (!projectId) return;
    
    // Check localStorage for active crawl
    const activeCrawlId = localStorage.getItem(`activeCrawl_${projectId}`);
    if (activeCrawlId) {
      try {
        const { data: crawl, error: crawlError } = await fetchCrawl(activeCrawlId);
        
        if (!crawlError && crawl) {
          if (crawl.status === 'running' || crawl.status === 'pending') {
            activeCrawl = crawl;
            showActiveCrawlNotification = true;
          } else {
            // Clean up if crawl is done
            localStorage.removeItem(`activeCrawl_${projectId}`);
          }
        } else {
          // Crawl not found, clean up
          localStorage.removeItem(`activeCrawl_${projectId}`);
        }
      } catch (err) {
        console.error('Error checking active crawl:', err);
        // Clean up on error
        localStorage.removeItem(`activeCrawl_${projectId}`);
      }
    }
  }

  async function loadData() {
    if (!projectId) return;
    
    loading = true;
    try {
      // Load projects
      const { data: projectsData, error: projectsError } = await fetchProjects();
      if (projectsError) throw projectsError;
      projects = projectsData || [];
      
      // Find current project
      project = projects.find(p => p.id === projectId);
      selectedProject = project;
      if (!project) {
        error = 'Project not found';
        loading = false;
        return;
      }

      // Load crawls for this project
      const { data: crawlsData, error: crawlsError } = await fetchCrawls(projectId);
      if (crawlsError) throw crawlsError;
      crawls = crawlsData || [];
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  }

  function handleProjectSelect(selectedProject) {
    push(`/project/${selectedProject.id}`);
  }

  function handleCrawlSelect(crawl) {
    push(`/project/${projectId}/crawl/${crawl.id}`);
  }

  async function handleCrawlCreated(e) {
    // Hide notification when new crawl is created
    showActiveCrawlNotification = false;
    activeCrawl = null;
    
    // Reload crawls
    const { data: crawlsData, error: crawlsError } = await fetchCrawls(projectId);
    if (!crawlsError && crawlsData) {
      crawls = crawlsData;
    }
  }
  
  function handleViewActiveCrawl() {
    if (activeCrawl?.id) {
      showActiveCrawlNotification = false;
      push(`/project/${projectId}/crawl/${activeCrawl.id}`);
    }
  }
  
  function handleDismissActiveCrawl() {
    showActiveCrawlNotification = false;
  }
</script>

{#if loading}
  <div class="flex items-center justify-center min-h-screen">
    <span class="loading loading-spinner loading-lg"></span>
  </div>
{:else if error}
  <div class="flex items-center justify-center min-h-screen">
    <div class="alert alert-error max-w-md">
      <span>Error: {error}</span>
    </div>
  </div>
{:else if project}
  <!-- Header -->
  <ProjectsView {projects} {selectedProject} on:select={(e) => handleProjectSelect(e.detail)} />
  
  <!-- Crawl Section -->
  {#if showCrawlSection}
    <div class="container mx-auto p-4 border-b border-base-200">
      {#if showActiveCrawlNotification && activeCrawl}
        <div class="alert alert-info mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="stroke-current shrink-0 w-6 h-6">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          </svg>
          <div class="flex-1">
            <h3 class="font-bold">Crawl in Progress</h3>
            <div class="text-sm">A crawl is currently running for this project.</div>
          </div>
          <div class="flex gap-2">
            <button class="btn btn-sm btn-primary" on:click={handleViewActiveCrawl}>
              View Progress
            </button>
            <button class="btn btn-sm btn-ghost" on:click={handleDismissActiveCrawl}>
              Dismiss
            </button>
          </div>
        </div>
      {/if}
      
      <div class="flex justify-between items-center mb-4">
        <CrawlSelector 
          {crawls} 
          {projectId} 
          on:select={(e) => handleCrawlSelect(e.detail)} 
        />
        <TriggerCrawlButton {projectId} project={project} on:created={handleCrawlCreated} />
      </div>
    </div>
  {/if}

  <!-- Project Layout with Sidebar -->
  <ProjectLayout {projectId} gscStatus={finalGSCStatus}>
    <slot></slot>
  </ProjectLayout>
{/if}

