<script>
  import { params } from 'svelte-spa-router';
  import ProjectPageLayout from '../components/ProjectPageLayout.svelte';
  import AIUsageDashboard from '../components/AIUsageDashboard.svelte';

  let projectId = '';

  $: if ($params?.projectId) {
    projectId = $params.projectId;
  }
</script>

<ProjectPageLayout {projectId}>
  <div class="max-w-7xl mx-auto">
    <div class="mb-6">
      <h2 class="text-2xl font-bold">AI Usage</h2>
      <p class="text-sm opacity-60">Track your AI feature usage against plan limits.</p>
    </div>

    {#if projectId}
      <AIUsageDashboard {projectId} />
    {:else}
      <div class="text-center py-12 opacity-60">Loading project...</div>
    {/if}
  </div>
</ProjectPageLayout>
