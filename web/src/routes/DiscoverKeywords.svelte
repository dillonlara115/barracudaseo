<script>
  import { onMount } from 'svelte';
  import { params } from 'svelte-spa-router';
  import ProjectPageLayout from '../components/ProjectPageLayout.svelte';
  import KeywordDiscovery from '../components/KeywordDiscovery.svelte';
  import { fetchProjectGSCStatus } from '../lib/data.js';

  let projectId = null;
  let gscStatus = null;
  let gscLoading = false;

  $: projectId = $params?.projectId || null;

  onMount(async () => {
    if (projectId) {
      await loadGSCStatus();
    }
  });

  $: if (projectId) {
    loadGSCStatus();
  }

  async function loadGSCStatus() {
    if (!projectId) return;
    gscLoading = true;
    const result = await fetchProjectGSCStatus(projectId);
    if (!result.error && result.data) {
      gscStatus = result.data;
    }
    gscLoading = false;
  }
</script>

<ProjectPageLayout {projectId} {gscStatus} showCrawlSection={false}>
  <div class="max-w-7xl mx-auto">
    <div class="mb-6">
      <h1 class="text-3xl font-bold mb-2">Discover Keywords</h1>
      <p class="text-base-content/70">
        Find keywords that your domain or specific URLs are currently ranking for using our keyword discovery tools. 
        Discover new opportunities and add them to your tracking list.
      </p>
    </div>

    <KeywordDiscovery {projectId} />
  </div>
</ProjectPageLayout>

