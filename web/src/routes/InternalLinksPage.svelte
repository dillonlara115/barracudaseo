<script>
  import { params } from 'svelte-spa-router';
  import ProjectPageLayout from '../components/ProjectPageLayout.svelte';
  import InternalLinks from '../components/InternalLinks.svelte';

  let projectId = '';

  $: if ($params?.projectId) {
    projectId = $params.projectId;
  }
</script>

<ProjectPageLayout {projectId}>
  <div class="max-w-7xl mx-auto">
    <div class="mb-6">
      <h2 class="text-2xl font-bold">Internal Links</h2>
      <p class="text-sm opacity-60">Semantic link suggestions and orphaned page detection across your site.</p>
    </div>

    {#if projectId}
      <InternalLinks {projectId} />
    {:else}
      <div class="text-center py-12 opacity-60">Loading project...</div>
    {/if}
  </div>
</ProjectPageLayout>
