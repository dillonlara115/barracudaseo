<script>
  import { params } from 'svelte-spa-router';
  import ProjectPageLayout from '../components/ProjectPageLayout.svelte';
  import GSCIntelligence from '../components/GSCIntelligence.svelte';

  let projectId = '';

  $: if ($params?.projectId) {
    projectId = $params.projectId;
  }
</script>

<ProjectPageLayout {projectId}>
  <div class="max-w-7xl mx-auto">
    <div class="mb-6">
      <h2 class="text-2xl font-bold">GSC Intelligence</h2>
      <p class="text-sm opacity-60">AI-powered analysis of your Google Search Console data.</p>
    </div>

    {#if projectId}
      <GSCIntelligence {projectId} />
    {:else}
      <div class="text-center py-12 opacity-60">Loading project...</div>
    {/if}
  </div>
</ProjectPageLayout>
