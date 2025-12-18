<script>
  import { createEventDispatcher } from 'svelte';
  import { createProject } from '../lib/data.js';
  
  const dispatch = createEventDispatcher();
  
  export let className = '';
  
  let showCreateModal = false;
  let newProjectName = '';
  let newProjectUrl = '';
  let creating = false;
  let error = null;

  // Extract domain from URL
  function extractDomain(url) {
    if (!url) return '';
    try {
      const urlObj = new URL(url.startsWith('http') ? url : `https://${url}`);
      return urlObj.hostname.replace(/^www\./, ''); // Remove www. prefix
    } catch (e) {
      // If URL parsing fails, try to extract domain manually
      const cleaned = url.replace(/^https?:\/\//, '').replace(/^www\./, '').split('/')[0];
      return cleaned;
    }
  }

  async function handleCreateProject() {
    if (!newProjectName || !newProjectUrl) {
      error = 'Name and starting URL are required';
      return;
    }

    // Validate URL format
    let validatedUrl = newProjectUrl.trim();
    if (!validatedUrl.startsWith('http://') && !validatedUrl.startsWith('https://')) {
      validatedUrl = `https://${validatedUrl}`;
    }

    try {
      new URL(validatedUrl);
    } catch (e) {
      error = 'Invalid URL format';
      return;
    }

    // Extract domain from URL
    const domain = extractDomain(validatedUrl);
    if (!domain) {
      error = 'Could not extract domain from URL';
      return;
    }

    creating = true;
    error = null;

    try {
      const { data, error: createError } = await createProject(
        newProjectName,
        domain,
        { url: validatedUrl }
      );

      if (createError) throw createError;

      // Emit the created project to parent
      dispatch('created', data);
      
      // Reset form and close modal
      newProjectName = '';
      newProjectUrl = '';
      showCreateModal = false;
    } catch (err) {
      error = err.message || 'Failed to create project';
    } finally {
      creating = false;
    }
  }
</script>

<button 
  class="btn btn-primary {className}"
  on:click={() => showCreateModal = true}
>
  Create Project
</button>

<!-- Create Project Modal -->
{#if showCreateModal}
  <div class="modal modal-open">
    <div class="modal-box bg-base-100">
      <h3 class="font-bold text-lg mb-4 text-base-content">Create New Project</h3>

      {#if error}
        <div class="alert alert-error mb-4">
          <span>{error}</span>
        </div>
      {/if}

      <div class="form-control w-full mb-4">
        <label class="label" for="project-name">
          <span class="label-text text-base-content">Project Name</span>
        </label>
        <input
          id="project-name"
          type="text"
          placeholder="My Website"
          class="input input-bordered w-full bg-base-200 text-base-content placeholder-gray-500 border-base-300 focus:border-primary"
          bind:value={newProjectName}
        />
      </div>

      <div class="form-control w-full mb-4">
        <label class="label" for="starting-url">
          <span class="label-text text-base-content">Starting URL</span>
        </label>
        <input
          id="starting-url"
          type="url"
          placeholder="https://example.com"
          class="input input-bordered w-full bg-base-200 text-base-content placeholder-gray-500 border-base-300 focus:border-primary"
          bind:value={newProjectUrl}
        />
        <div class="label">
          <span class="label-text-alt text-base-content opacity-70">The domain will be automatically extracted from this URL</span>
        </div>
      </div>

      <div class="modal-action">
        <button
          class="btn btn-ghost text-base-content hover:bg-base-200"
          on:click={() => {
            showCreateModal = false;
            error = null;
          }}
        >
          Cancel
        </button>
        <button
          class="btn btn-primary text-primary-content"
          on:click={handleCreateProject}
          disabled={creating || !newProjectName || !newProjectUrl}
        >
          {#if creating}
            <span class="loading loading-spinner loading-sm"></span>
          {:else}
            Create
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}

