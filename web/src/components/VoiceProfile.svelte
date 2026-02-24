<script>
  import { onMount } from 'svelte';
  import { fetchVoiceProfile, updateVoiceProfile, streamAIRequest } from '../lib/data.js';

  export let projectId;

  let profile = null;
  let loading = false;
  let generating = false;
  let streamText = '';
  let error = null;
  let saveStatus = '';

  const fields = [
    { key: 'tone', label: 'Tone', desc: 'Formality level, personality traits, humor, emotional register' },
    { key: 'structure', label: 'Structure', desc: 'How articles open, heading patterns, paragraph norms, CTA style' },
    { key: 'sentence_style', label: 'Sentence Style', desc: 'Length, person, punctuation tendencies' },
    { key: 'brand_context', label: 'Brand Context', desc: 'Audience, expertise areas, products/services, positioning' },
    { key: 'avoid_list', label: 'Avoid List', desc: 'Overused phrases, off-brand language, generic AI filler' }
  ];

  onMount(loadProfile);

  async function loadProfile() {
    loading = true;
    const { data, error: err } = await fetchVoiceProfile(projectId);
    loading = false;
    if (err) {
      error = err.message;
      return;
    }
    profile = data;
  }

  async function generate() {
    generating = true;
    streamText = '';
    error = null;

    await streamAIRequest('/api/v1/ai/voice/generate', { project_id: projectId }, {
      onChunk: (text) => { streamText += text; },
      onDone: async () => {
        generating = false;
        await loadProfile();
      },
      onError: (err) => {
        generating = false;
        error = err.message;
      }
    });
  }

  let saveTimeout;
  function handleFieldChange(key, value) {
    if (!profile) return;
    profile[key] = value;
    saveStatus = 'Saving...';
    clearTimeout(saveTimeout);
    saveTimeout = setTimeout(async () => {
      await updateVoiceProfile({ project_id: projectId, [key]: value });
      saveStatus = 'Saved';
      setTimeout(() => { saveStatus = ''; }, 2000);
    }, 500);
  }
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <h3 class="text-lg font-bold">Brand Voice Profile</h3>
      <p class="text-sm opacity-60">AI-generated writing voice that matches your brand across all content.</p>
    </div>
    <div class="flex items-center gap-3">
      {#if saveStatus}
        <span class="text-sm opacity-60">{saveStatus}</span>
      {/if}
      <button class="btn btn-primary btn-sm" on:click={generate} disabled={generating}>
        {#if generating}
          <span class="loading loading-spinner loading-sm"></span>
          Analyzing...
        {:else}
          {profile ? 'Regenerate' : 'Analyze Brand Voice'}
        {/if}
      </button>
    </div>
  </div>

  {#if error}
    <div class="alert alert-error text-sm">{error}</div>
  {/if}

  {#if generating}
    <div class="card bg-base-200">
      <div class="card-body">
        <p class="text-sm opacity-60 mb-2">Analyzing your site content...</p>
        <div class="prose prose-sm max-w-none whitespace-pre-wrap">{streamText}</div>
      </div>
    </div>
  {:else if loading}
    <div class="flex justify-center py-12">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
  {:else if profile}
    {#each fields as { key, label, desc }}
      <div class="form-control">
        <label class="label" for="voice-{key}">
          <span class="label-text font-semibold">{label}</span>
          <span class="label-text-alt opacity-50">{desc}</span>
        </label>
        <textarea
          id="voice-{key}"
          class="textarea textarea-bordered h-32 text-sm"
          on:input={(e) => handleFieldChange(key, /** @type {HTMLTextAreaElement} */ (e.target).value)}
        >{profile[key] || ''}</textarea>
      </div>
    {/each}
  {:else}
    <div class="card bg-base-200">
      <div class="card-body text-center">
        <p class="opacity-60">No voice profile generated yet. Click "Analyze Brand Voice" to create one from your site content.</p>
      </div>
    </div>
  {/if}
</div>
