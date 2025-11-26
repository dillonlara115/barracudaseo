<script>
  import { createEventDispatcher } from 'svelte';
  import { createKeyword } from '../lib/data.js';
  import { X } from 'lucide-svelte';

  export let projectId = null;

  const dispatch = createEventDispatcher();

  let loading = false;
  let error = null;
  
  let formData = {
    keyword: '',
    target_url: '',
    location_name: 'United States',
    language_name: 'English',
    device: 'desktop',
    search_engine: 'google.com',
    check_frequency: 'manual',
    tags: []
  };

  let tagInput = '';

  const commonLocations = [
    'United States',
    'United Kingdom',
    'Canada',
    'Australia',
    'Germany',
    'France',
    'Spain',
    'Italy',
    'Netherlands',
    'Sweden',
    'Denmark',
    'Norway',
    'Japan',
    'South Korea',
    'India',
    'Brazil',
    'Mexico',
    'Argentina',
    'Chile',
    'New Zealand'
  ];

  async function handleSubmit() {
    if (!formData.keyword.trim()) {
      error = 'Keyword is required';
      return;
    }

    if (!formData.location_name) {
      error = 'Location is required';
      return;
    }

    loading = true;
    error = null;

    const keywordData = {
      project_id: projectId,
      keyword: formData.keyword.trim(),
      location_name: formData.location_name,
      language_name: formData.language_name,
      device: formData.device,
      search_engine: formData.search_engine,
      check_frequency: formData.check_frequency,
    };

    if (formData.target_url.trim()) {
      keywordData.target_url = formData.target_url.trim();
    }

    if (formData.tags.length > 0) {
      keywordData.tags = formData.tags;
    }

    const result = await createKeyword(keywordData);
    
    loading = false;

    if (result.error) {
      error = result.error.message || 'Failed to create keyword';
      return;
    }

    dispatch('created');
    dispatch('close');
  }

  function addTag() {
    const tag = tagInput.trim();
    if (tag && !formData.tags.includes(tag)) {
      formData.tags = [...formData.tags, tag];
      tagInput = '';
    }
  }

  function removeTag(tag) {
    formData.tags = formData.tags.filter(t => t !== tag);
  }

  function handleKeydown(event) {
    if (event.key === 'Enter') {
      event.preventDefault();
      addTag();
    }
  }
</script>

<div class="modal modal-open">
  <div class="modal-box max-w-2xl">
    <div class="flex items-center justify-between mb-4">
      <h3 class="font-bold text-lg">Add Keyword</h3>
      <button class="btn btn-sm btn-circle btn-ghost" on:click={() => dispatch('close')}>
        <X class="w-4 h-4" />
      </button>
    </div>

    {#if error}
      <div class="alert alert-error mb-4">
        <span>{error}</span>
      </div>
    {/if}

    <form on:submit|preventDefault={handleSubmit}>
      <div class="space-y-4">
        <div class="form-control">
          <label class="label">
            <span class="label-text">Keyword *</span>
          </label>
          <input
            type="text"
            class="input input-bordered w-full"
            placeholder="e.g., best singing bowls"
            bind:value={formData.keyword}
            required
          />
        </div>

        <div class="form-control">
          <label class="label">
            <span class="label-text">Target URL (optional)</span>
          </label>
          <input
            type="url"
            class="input input-bordered w-full"
            placeholder="https://example.com/page"
            bind:value={formData.target_url}
          />
          <label class="label">
            <span class="label-text-alt">The URL you want to rank for this keyword</span>
          </label>
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div class="form-control">
            <label class="label">
              <span class="label-text">Location *</span>
            </label>
            <select class="select select-bordered w-full" bind:value={formData.location_name} required>
              {#each commonLocations as location}
                <option value={location}>{location}</option>
              {/each}
            </select>
          </div>

          <div class="form-control">
            <label class="label">
              <span class="label-text">Language</span>
            </label>
            <select class="select select-bordered w-full" bind:value={formData.language_name}>
              <option value="English">English</option>
              <option value="Spanish">Spanish</option>
              <option value="French">French</option>
              <option value="German">German</option>
              <option value="Italian">Italian</option>
              <option value="Portuguese">Portuguese</option>
              <option value="Japanese">Japanese</option>
              <option value="Korean">Korean</option>
              <option value="Chinese">Chinese</option>
            </select>
          </div>
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div class="form-control">
            <label class="label">
              <span class="label-text">Device</span>
            </label>
            <select class="select select-bordered w-full" bind:value={formData.device}>
              <option value="desktop">Desktop</option>
              <option value="mobile">Mobile</option>
            </select>
          </div>

          <div class="form-control">
            <label class="label">
              <span class="label-text">Search Engine</span>
            </label>
            <select class="select select-bordered w-full" bind:value={formData.search_engine}>
              <option value="google.com">Google</option>
              <option value="google.co.uk">Google UK</option>
              <option value="google.ca">Google Canada</option>
              <option value="google.com.au">Google Australia</option>
            </select>
          </div>
        </div>

        <div class="form-control">
          <label class="label">
            <span class="label-text">Check Frequency</span>
          </label>
          <select class="select select-bordered w-full" bind:value={formData.check_frequency}>
            <option value="manual">Manual (check when you click "Check Now")</option>
            <option value="daily">Daily (automatic daily checks)</option>
            <option value="weekly">Weekly (automatic weekly checks)</option>
          </select>
          <label class="label">
            <span class="label-text-alt">Scheduled checks run automatically via cron job</span>
          </label>
        </div>

        <div class="form-control">
          <label class="label">
            <span class="label-text">Tags (optional)</span>
          </label>
          <div class="flex gap-2 mb-2">
            <input
              type="text"
              class="input input-bordered flex-1"
              placeholder="Add a tag and press Enter"
              bind:value={tagInput}
              on:keydown={handleKeydown}
            />
            <button type="button" class="btn btn-outline" on:click={addTag}>
              Add
            </button>
          </div>
          {#if formData.tags.length > 0}
            <div class="flex flex-wrap gap-2">
              {#each formData.tags as tag}
                <span class="badge badge-lg gap-2">
                  {tag}
                  <button type="button" class="btn btn-xs btn-circle btn-ghost" on:click={() => removeTag(tag)}>
                    <X class="w-3 h-3" />
                  </button>
                </span>
              {/each}
            </div>
          {/if}
        </div>
      </div>

      <div class="modal-action">
        <button type="button" class="btn btn-ghost" on:click={() => dispatch('close')} disabled={loading}>
          Cancel
        </button>
        <button type="submit" class="btn btn-primary" disabled={loading}>
          {#if loading}
            <span class="loading loading-spinner loading-sm"></span>
            Creating...
          {:else}
            Create Keyword
          {/if}
        </button>
      </div>
    </form>
  </div>
  <div class="modal-backdrop" on:click={() => dispatch('close')}></div>
</div>

