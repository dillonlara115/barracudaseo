<script>
  import { onMount } from 'svelte';
  import { params, push } from 'svelte-spa-router';
  import { fetchProjects, fetchCrawls } from '../lib/data.js';
  import ProjectPageLayout from '../components/ProjectPageLayout.svelte';

  let project = null;
  let loading = true;
  let error = null;

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
      project = projects?.find(p => p.id === projectId) || null;
      if (!project) {
        error = 'Project not found';
      }
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  }

  async function navigateBackToProject() {
    if (!projectId) return;
    try {
      const { data: crawlsData } = await fetchCrawls(projectId);
      if (crawlsData && crawlsData.length > 0) {
        push(`/project/${projectId}/crawl/${crawlsData[0].id}`);
      } else {
        push(`/project/${projectId}`);
      }
    } catch (err) {
      push(`/project/${projectId}`);
    }
  }
</script>

<ProjectPageLayout {projectId} showCrawlSection={false}>
  <div class="container mx-auto p-6 max-w-4xl">
    <div class="mb-6">
      <button class="btn btn-ghost btn-sm mb-4" on:click={navigateBackToProject}>
        ← Back to Project
      </button>
      <h1 class="text-3xl font-bold mb-2">CLI Setup</h1>
      <p class="text-base-content/70">
        Install the Barracuda CLI and connect it to this project for fast local crawls.
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
      <div class="card bg-base-100 shadow">
        <div class="card-body space-y-4">
          <div>
            <p class="text-sm font-semibold mb-1">1) Install (macOS/Linux)</p>
            <pre class="bg-base-200 rounded-lg p-3 text-sm overflow-x-auto"><code>curl -fsSL https://raw.githubusercontent.com/dillonlara115/barracuda/main/scripts/install-barracuda.sh | bash</code></pre>
          </div>

          <div>
            <p class="text-sm font-semibold mb-1">Homebrew (macOS/Linux)</p>
            <pre class="bg-base-200 rounded-lg p-3 text-sm overflow-x-auto"><code>brew install barracuda/tap/barracuda</code></pre>
          </div>

          <div>
            <p class="text-sm font-semibold mb-1">2) Link your account</p>
            <pre class="bg-base-200 rounded-lg p-3 text-sm overflow-x-auto"><code>barracuda auth login</code></pre>
          </div>

          <div>
            <p class="text-sm font-semibold mb-1">3) Run a crawl and upload to this project</p>
            <pre class="bg-base-200 rounded-lg p-3 text-sm overflow-x-auto"><code>barracuda crawl https://{project.domain || 'example.com'} --cloud --project-id {project.id}</code></pre>
          </div>

          <div class="text-sm text-base-content/70">
            <a
              class="link link-primary"
              href="https://github.com/dillonlara115/barracuda/releases/latest"
              target="_blank"
              rel="noopener noreferrer"
            >
              Download other binaries
            </a>
            <span> · </span>
            <a
              class="link link-primary"
              href="https://github.com/dillonlara115/barracuda"
              target="_blank"
              rel="noopener noreferrer"
            >
              View CLI docs
            </a>
          </div>

          <p class="text-xs text-base-content/60">
            Windows users: download the latest .exe from Releases and add it to your PATH.
          </p>
        </div>
      </div>
    {/if}
  </div>
</ProjectPageLayout>
