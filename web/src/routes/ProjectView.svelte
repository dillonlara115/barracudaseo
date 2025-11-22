<script>
  import { onMount } from 'svelte';
  import { push, params } from 'svelte-spa-router';
  import { fetchProjects, fetchCrawls, fetchCrawl } from '../lib/data.js';
  import ProjectsView from '../components/ProjectsView.svelte';
  import CrawlSelector from '../components/CrawlSelector.svelte';
  import TriggerCrawlButton from '../components/TriggerCrawlButton.svelte';
  
  let projects = [];
  let project = null;
  let selectedProject = null;
  let crawls = [];
  let loading = true;
  let error = null;
  let currentProjectId = null;
  let showActiveCrawlNotification = false;
  let activeCrawl = null;

  $: projectId = $params?.id || null;

  onMount(async () => {
    // Wait a tick for params to be available
    await new Promise(resolve => setTimeout(resolve, 0));
    if ($params?.id) {
      await loadData();
      await checkActiveCrawl();
    }
  });
  
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

  $: if (projectId && projectId !== currentProjectId && $params?.id) {
    currentProjectId = projectId;
    loadData();
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

      // If crawls exist, redirect to latest crawl
      if (crawls.length > 0) {
        push(`/project/${projectId}/crawl/${crawls[0].id}`);
      }
    } catch (err) {
      error = err.message;
      loading = false;
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
    
    // Reload crawls and redirect to the new crawl
    const { data: crawlsData, error: crawlsError } = await fetchCrawls(projectId);
    if (!crawlsError && crawlsData && crawlsData.length > 0) {
      crawls = crawlsData;
      // Find the new crawl (should be first/latest)
      const newCrawl = crawlsData.find(c => c.id === e.detail.crawl_id) || crawlsData[0];
      push(`/project/${projectId}/crawl/${newCrawl.id}`);
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
  <ProjectsView {projects} {selectedProject} on:select={(e) => handleProjectSelect(e.detail)} />
  
  <div class="container mx-auto p-4">
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
      <h2 class="text-2xl font-bold text-base-content">Crawls</h2>
      <TriggerCrawlButton {projectId} project={project} on:created={handleCrawlCreated} />
    </div>
    
    {#if crawls.length === 0}
      <div class="alert alert-info">
        <span>No crawls found for this project. Start a crawl to get started.</span>
      </div>
    {:else}
      <CrawlSelector {crawls} {projectId} on:select={(e) => handleCrawlSelect(e.detail)} />
    {/if}
  </div>
{/if}

