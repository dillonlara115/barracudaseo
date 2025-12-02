<script>
  import { createEventDispatcher } from 'svelte';
  import KeywordDiscoveryContent from './KeywordDiscoveryContent.svelte';

  export let projectId = null;
  export let defaultTarget = '';
  export let showAsModal = false;

  const dispatch = createEventDispatcher();

  let formData = {
    target: defaultTarget,
    location_name: 'United States',
    language_name: 'English',
    limit: 1000,
    min_position: 0,
    max_position: 0
  };

  function handleKeywordAdded(event) {
    dispatch('keyword-added', event.detail);
  }

  function handleKeywordsAdded(event) {
    dispatch('keywords-added', event.detail);
  }
</script>

{#if showAsModal}
<div class="modal modal-open">
  <div class="modal-box max-w-6xl max-h-[90vh] overflow-y-auto">
    <div class="flex items-center justify-between mb-4">
      <h3 class="font-bold text-lg">Discover Keywords</h3>
      <button class="btn btn-sm btn-circle btn-ghost" on:click={() => dispatch('close')}>
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
          <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>
    <KeywordDiscoveryContent
      {projectId}
      {defaultTarget}
      {formData}
      on:keyword-added={handleKeywordAdded}
      on:keywords-added={handleKeywordsAdded}
    />
  </div>
  <div 
    class="modal-backdrop" 
    role="button"
    tabindex="0"
    on:click={() => dispatch('close')}
    on:keydown={(e) => e.key === 'Enter' || e.key === ' ' ? dispatch('close') : null}
  ></div>
</div>
{:else}
<KeywordDiscoveryContent
  {projectId}
  {defaultTarget}
  {formData}
  on:keyword-added={handleKeywordAdded}
  on:keywords-added={handleKeywordsAdded}
/>
{/if}
