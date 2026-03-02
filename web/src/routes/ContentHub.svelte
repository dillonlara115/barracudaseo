<script>
  import { params, querystring } from 'svelte-spa-router';
  import ProjectPageLayout from '../components/ProjectPageLayout.svelte';
  import ContentBriefs from '../components/ContentBriefs.svelte';
  import ArticleWriter from '../components/ArticleWriter.svelte';

  let projectId = '';
  let activeTab = 'briefs';

  $: if ($params?.projectId) {
    projectId = $params.projectId;
  }

  $: initialKeyword = $querystring
    ? new URLSearchParams($querystring).get('keyword') || ''
    : '';
</script>

<ProjectPageLayout {projectId}>
  <div class="max-w-7xl mx-auto">
    <div class="mb-6">
      <h2 class="text-2xl font-bold">Content</h2>
      <p class="text-sm opacity-60">Generate data-backed briefs and AI-written articles for your project.</p>
    </div>

    <div class="tabs tabs-boxed mb-6">
      <button class="tab tab-lg" class:tab-active={activeTab === 'briefs'} on:click={() => activeTab = 'briefs'}>
        Content Briefs
      </button>
      <button class="tab tab-lg" class:tab-active={activeTab === 'articles'} on:click={() => activeTab = 'articles'}>
        Article Writer
      </button>
    </div>

    {#if projectId}
      {#if activeTab === 'briefs'}
        <ContentBriefs {projectId} {initialKeyword} />
      {:else}
        <ArticleWriter {projectId} />
      {/if}
    {:else}
      <div class="text-center py-12 opacity-60">Loading project...</div>
    {/if}
  </div>
</ProjectPageLayout>
