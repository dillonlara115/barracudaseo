<script>
  import { onMount } from 'svelte';
  import { params, push, link } from 'svelte-spa-router';
  import { fetchProjects, updateProject, deleteProject, fetchCrawls } from '../lib/data.js';
  import ProjectGSCSelector from '../components/ProjectGSCSelector.svelte';
  import ProjectGA4Selector from '../components/ProjectGA4Selector.svelte';
  import ProjectClaritySelector from '../components/ProjectClaritySelector.svelte';
  import CrawlManagement from '../components/CrawlManagement.svelte';
  import ProjectPageLayout from '../components/ProjectPageLayout.svelte';
  import { Loader, X, Trash2, Edit2, Check, AlertTriangle } from 'lucide-svelte';
  
  let project = null;
  let summary = null; // For enriching issues if needed
  let loading = true;
  let error = null;

  // Project editing state
  let editingProject = false;
  let projectName = '';
  let projectDomain = '';
  let updatingProject = false;
  let updateError = null;
  let updateSuccess = null;

  // Delete state
  let showDeleteConfirm = false;
  let deletingProject = false;
  let deleteError = null;

  $: projectId = $params?.projectId || null;

  onMount(async () => {
    if (projectId) {
      await loadProject();
    }
  });

  $: if (projectId) {
    loadProject();
  }

  async function loadProject() {
    if (!projectId) return;
    
    loading = true;
    try {
      const { data: projects, error: projectsError } = await fetchProjects();
      if (projectsError) throw projectsError;
      
      project = projects?.find(p => p.id === projectId);
      if (!project) {
        error = 'Project not found';
      } else {
        // Initialize form fields
        projectName = project.name || '';
        projectDomain = project.domain || '';
      }
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  }

  function startEditing() {
    editingProject = true;
    projectName = project?.name || '';
    projectDomain = project?.domain || '';
    updateError = null;
    updateSuccess = null;
  }

  function cancelEditing() {
    editingProject = false;
    projectName = project?.name || '';
    projectDomain = project?.domain || '';
    updateError = null;
    updateSuccess = null;
  }

  async function saveProject() {
    if (!projectId || !projectName.trim() || !projectDomain.trim()) {
      updateError = 'Name and domain are required';
      return;
    }

    updatingProject = true;
    updateError = null;
    updateSuccess = null;

    try {
      const { data, error: updateErr } = await updateProject(projectId, {
        name: projectName.trim(),
        domain: projectDomain.trim()
      });

      if (updateErr) throw updateErr;

      // Update local project object
      project = { ...project, ...data };
      updateSuccess = 'Project updated successfully';
      editingProject = false;

      // Reload projects list to ensure consistency
      await loadProject();
    } catch (err) {
      updateError = err.message || 'Failed to update project';
    } finally {
      updatingProject = false;
    }
  }

  function confirmDelete() {
    showDeleteConfirm = true;
    deleteError = null;
  }

  function cancelDelete() {
    showDeleteConfirm = false;
    deleteError = null;
  }

  async function handleDelete() {
    if (!projectId) return;

    deletingProject = true;
    deleteError = null;

    try {
      const { error: deleteErr } = await deleteProject(projectId);
      if (deleteErr) throw deleteErr;

      // Redirect to projects list after successful deletion
      push('/');
    } catch (err) {
      deleteError = err.message || 'Failed to delete project';
    } finally {
      deletingProject = false;
    }
  }

  function handleEnriched(e) {
    // Navigate to issues tab to see enriched data
    push(`/project/${projectId}/crawl/${$params.crawlId || ''}?tab=issues`);
  }

  function handleCrawlDeleted(e) {
    // Crawl was deleted, could reload or show notification
    // The CrawlManagement component handles its own reload
  }

  async function navigateBackToProject() {
    if (!projectId) return;
    
    // Check if there are crawls - if so, navigate to the latest crawl
    // Otherwise navigate to project view
    try {
      const { data: crawlsData } = await fetchCrawls(projectId);
      if (crawlsData && crawlsData.length > 0) {
        // Navigate to latest crawl
        push(`/project/${projectId}/crawl/${crawlsData[0].id}`);
      } else {
        // Navigate to project view (no crawls)
        push(`/project/${projectId}`);
      }
    } catch (err) {
      // On error, just navigate to project view
      push(`/project/${projectId}`);
    }
  }
</script>

<ProjectPageLayout {projectId} showCrawlSection={false}>
<div class="container mx-auto p-6 max-w-4xl">
  <div class="mb-6">
    <button 
      class="btn btn-ghost btn-sm mb-4"
      on:click={navigateBackToProject}
    >
      ← Back to Project
    </button>
    <h1 class="text-3xl font-bold mb-2">Project Settings</h1>
    <p class="text-base-content/70">
      Configure project information, Google Search Console integration, crawl management, and other settings for this project.
    </p>
  </div>

  {#if loading}
    <div class="flex items-center justify-center min-h-[400px]">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
  {:else if error}
    <div class="alert alert-error">
      <span>{error}</span>
    </div>
  {:else if project}
    <div class="space-y-6">
      <!-- Project Information Card -->
      <div class="card bg-base-100 shadow">
        <div class="card-body">
          <div class="flex items-center justify-between mb-4">
            <h2 class="card-title text-xl">Project Information</h2>
            {#if !editingProject}
              <button 
                class="btn btn-sm btn-outline"
                on:click={startEditing}
              >
                <Edit2 class="w-4 h-4 mr-2" />
                Edit
              </button>
            {/if}
          </div>

          {#if editingProject}
            {#if updateError}
              <div class="alert alert-error mb-4">
                <X class="w-5 h-5" />
                <span>{updateError}</span>
              </div>
            {/if}

            {#if updateSuccess}
              <div class="alert alert-success mb-4">
                <Check class="w-5 h-5" />
                <span>{updateSuccess}</span>
              </div>
            {/if}

            <div class="space-y-4">
              <div class="form-control">
                <label class="label">
                  <span class="label-text font-semibold">Project Name</span>
                </label>
                <input 
                  type="text" 
                  class="input input-bordered"
                  bind:value={projectName}
                  placeholder="My Website"
                  disabled={updatingProject}
                />
              </div>

              <div class="form-control">
                <label class="label">
                  <span class="label-text font-semibold">Domain</span>
                </label>
                <input 
                  type="text" 
                  class="input input-bordered"
                  bind:value={projectDomain}
                  placeholder="example.com"
                  disabled={updatingProject}
                />
              </div>

              <div class="flex gap-2">
                <button 
                  class="btn btn-primary"
                  on:click={saveProject}
                  disabled={updatingProject || !projectName.trim() || !projectDomain.trim()}
                >
                  {#if updatingProject}
                    <Loader class="w-4 h-4 animate-spin" />
                  {:else}
                    <Check class="w-4 h-4" />
                  {/if}
                  Save Changes
                </button>
                <button 
                  class="btn btn-ghost"
                  on:click={cancelEditing}
                  disabled={updatingProject}
                >
                  Cancel
                </button>
              </div>
            </div>
          {:else}
            <div class="space-y-2">
              <div>
                <span class="text-sm text-base-content/70">Project Name</span>
                <p class="text-lg font-semibold">{project.name || 'N/A'}</p>
              </div>
              <div>
                <span class="text-sm text-base-content/70">Domain</span>
                <p class="text-lg font-semibold">{project.domain || 'N/A'}</p>
              </div>
            </div>
          {/if}
        </div>
      </div>

      <!-- Crawl Management -->
      {#if project.id}
        <CrawlManagement projectId={project.id} on:deleted={handleCrawlDeleted} />
      {:else}
        <div class="alert alert-warning">
          <span>Project ID not available.</span>
        </div>
      {/if}

      <!-- Google Search Console Integration -->
      <div class="card bg-base-100 shadow">
        <div class="card-body">
          <h2 class="card-title text-xl mb-4">Google Search Console Integration</h2>
          <ProjectGSCSelector {project} projectId={project.id} {summary} on:enriched={handleEnriched} />
        </div>
      </div>

      <!-- Google Analytics 4 Integration -->
      <div class="card bg-base-100 shadow">
        <div class="card-body">
          <h2 class="card-title text-xl mb-4">Google Analytics 4 Integration</h2>
          <ProjectGA4Selector {project} projectId={project.id} />
        </div>
      </div>

      <!-- Microsoft Clarity Integration -->
      <div class="card bg-base-100 shadow">
        <div class="card-body">
          <h2 class="card-title text-xl mb-4">Microsoft Clarity Integration</h2>
          <ProjectClaritySelector {project} projectId={project.id} />
        </div>
      </div>

      <!-- Danger Zone -->
      <div class="card bg-base-100 shadow">
        <div class="card-body">
          <p class="text-sm text-base-content/70 mb-4">
            Deleting a project will permanently remove it and all associated crawls, pages, and issues. This action cannot be undone.
          </p>
          <button 
            class="btn btn-error"
            on:click={confirmDelete}
            disabled={deletingProject}
          >
            <Trash2 class="w-4 h-4 mr-2" />
            Delete Project
          </button>
        </div>
      </div>
    </div>
  {/if}
</div>

<!-- Delete Confirmation Modal -->
{#if showDeleteConfirm}
  <div class="modal modal-open">
    <div class="modal-box">
      <h3 class="font-bold text-lg mb-4 flex items-center gap-2">
        <AlertTriangle class="w-6 h-6 text-error" />
        Delete Project
      </h3>
      
      <p class="py-4">
        Are you sure you want to delete <strong>{project?.name}</strong>? This will permanently delete:
      </p>
      
      <ul class="list-disc list-inside mb-4 text-sm text-base-content/70">
        <li>All crawls associated with this project</li>
        <li>All pages and issues from those crawls</li>
        <li>All project settings and integrations</li>
      </ul>

      <p class="text-error font-semibold mb-4">
        This action cannot be undone.
      </p>

      {#if deleteError}
        <div class="alert alert-error mb-4">
          <X class="w-5 h-5" />
          <span>{deleteError}</span>
        </div>
      {/if}

      <div class="modal-action">
        <button 
          class="btn btn-ghost"
          on:click={cancelDelete}
          disabled={deletingProject}
        >
          Cancel
        </button>
        <button 
          class="btn btn-error"
          on:click={handleDelete}
          disabled={deletingProject}
        >
          {#if deletingProject}
            <Loader class="w-4 h-4 animate-spin" />
            Deleting...
          {:else}
            <Trash2 class="w-4 h-4" />
            Delete Project
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}
</ProjectPageLayout>
