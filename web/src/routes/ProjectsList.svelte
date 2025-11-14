<script>
  import { onMount } from 'svelte';
  import { push } from 'svelte-spa-router';
  import { initAuth, user } from '../lib/auth.js';
  import { fetchProjects } from '../lib/data.js';
  import ProjectsView from '../components/ProjectsView.svelte';
  import CreateProjectButton from '../components/CreateProjectButton.svelte';
  
  let projects = [];
  let loading = true;
  let error = null;

  onMount(async () => {
    await initAuth();
    
    user.subscribe(async (currentUser) => {
      if (currentUser) {
        await loadProjects();
      } else {
        projects = [];
        loading = false;
      }
    });
  });

  async function loadProjects() {
    try {
      const { data, error: fetchError } = await fetchProjects();
      if (fetchError) throw fetchError;
      projects = data || [];
      
      // Sort projects: prefer "barracudaseo" or alphabetically by name
      projects.sort((a, b) => {
        const aName = (a.name || '').toLowerCase();
        const bName = (b.name || '').toLowerCase();
        
        // Put "barracudaseo" first if it exists
        if (aName.includes('barracuda') && !bName.includes('barracuda')) return -1;
        if (bName.includes('barracuda') && !aName.includes('barracuda')) return 1;
        
        // Otherwise sort alphabetically
        return aName.localeCompare(bName);
      });
      
      // Auto-redirect to first project if available (now sorted with barracudaseo first)
      if (projects.length > 0) {
        push(`/project/${projects[0].id}`);
      }
    } catch (err) {
      error = err.message;
      loading = false;
    } finally {
      loading = false;
    }
  }

  function handleProjectSelect(project) {
    push(`/project/${project.id}`);
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
{:else if projects.length === 0}
  <ProjectsView {projects} selectedProject={null} on:select={(e) => handleProjectSelect(e.detail)} />
  <div class="container mx-auto p-4">
    <div class="alert alert-info">
      <span>No projects yet.</span>
      <CreateProjectButton 
        className="ml-2 inline" 
        on:created={async (e) => {
          await loadProjects();
          handleProjectSelect(e.detail);
        }} 
      />
    </div>
  </div>
{:else}
  <ProjectsView {projects} selectedProject={null} on:select={(e) => handleProjectSelect(e.detail)} />
{/if}

