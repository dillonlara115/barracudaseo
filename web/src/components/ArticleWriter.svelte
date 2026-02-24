<script>
  import { onMount } from 'svelte';
  import { fetchArticles, fetchBriefs, updateArticle, streamAIRequest } from '../lib/data.js';

  export let projectId;

  let articles = [];
  let approvedBriefs = [];
  let loading = false;
  let error = null;

  let generating = false;
  let streamText = '';
  let selectedBriefId = '';

  let activeArticle = null;
  let editContent = '';

  onMount(loadData);

  async function loadData() {
    loading = true;
    const [artRes, briefRes] = await Promise.all([
      fetchArticles(projectId),
      fetchBriefs(projectId)
    ]);
    loading = false;

    if (artRes.error) { error = artRes.error.message; return; }
    articles = artRes.data || [];

    if (!briefRes.error) {
      approvedBriefs = (briefRes.data || []).filter(b => b.status === 'approved');
    }
  }

  async function generateArticle() {
    if (!selectedBriefId) return;
    generating = true;
    streamText = '';
    error = null;

    await streamAIRequest('/api/v1/ai/articles/generate', {
      project_id: projectId,
      brief_id: selectedBriefId
    }, {
      onChunk: (text) => { streamText += text; },
      onDone: async () => {
        generating = false;
        selectedBriefId = '';
        await loadData();
      },
      onError: (err) => {
        generating = false;
        error = err.message;
      }
    });
  }

  function openArticle(article) {
    activeArticle = article;
    editContent = article.content || '';
  }

  async function saveArticle() {
    if (!activeArticle) return;
    await updateArticle({ id: activeArticle.id, project_id: projectId, content: editContent });
    await loadData();
  }

  async function changeStatus(article, status) {
    await updateArticle({ id: article.id, project_id: projectId, status });
    await loadData();
  }

  function copyContent(content, format) {
    if (format === 'html') {
      const html = content.replace(/^### (.*$)/gm, '<h3>$1</h3>')
        .replace(/^## (.*$)/gm, '<h2>$1</h2>')
        .replace(/^# (.*$)/gm, '<h1>$1</h1>')
        .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
        .replace(/\*(.*?)\*/g, '<em>$1</em>')
        .replace(/\[(.*?)\]\((.*?)\)/g, '<a href="$2">$1</a>')
        .replace(/\n/g, '<br>');
      navigator.clipboard.writeText(html);
    } else {
      navigator.clipboard.writeText(content);
    }
  }

  function wordCount(text) {
    return text ? text.split(/\s+/).filter(Boolean).length : 0;
  }

  function readTime(text) {
    const words = wordCount(text);
    return Math.max(1, Math.ceil(words / 250));
  }
</script>

<div class="space-y-6">
  <div class="flex items-center justify-between">
    <h3 class="text-lg font-bold">Article Writer</h3>
    <span class="text-sm opacity-60">{articles.length} article{articles.length !== 1 ? 's' : ''}</span>
  </div>

  {#if error}
    <div class="alert alert-error text-sm">{error}</div>
  {/if}

  {#if !activeArticle}
    {#if approvedBriefs.length > 0}
      <div class="card bg-base-200">
        <div class="card-body">
          <h4 class="card-title text-sm">Generate Article from Brief</h4>
          <div class="flex gap-2">
            <select class="select select-bordered select-sm flex-1" bind:value={selectedBriefId} disabled={generating}>
              <option value="">Select an approved brief...</option>
              {#each approvedBriefs as brief}
                <option value={brief.id}>{brief.keyword}</option>
              {/each}
            </select>
            <button class="btn btn-primary btn-sm" on:click={generateArticle} disabled={generating || !selectedBriefId}>
              {#if generating}
                <span class="loading loading-spinner loading-sm"></span>
              {:else}
                Write Article
              {/if}
            </button>
          </div>
        </div>
      </div>
    {/if}

    {#if generating}
      <div class="card bg-base-200">
        <div class="card-body">
          <p class="text-sm opacity-60 mb-2">Writing article...</p>
          <div class="prose prose-sm max-w-none whitespace-pre-wrap">{streamText}</div>
        </div>
      </div>
    {/if}

    {#if loading}
      <div class="flex justify-center py-12">
        <span class="loading loading-spinner loading-lg"></span>
      </div>
    {:else}
      <div class="space-y-3">
        {#each articles as article}
          <div class="card bg-base-200 cursor-pointer hover:bg-base-300 transition-colors" on:click={() => openArticle(article)}>
            <div class="card-body py-4">
              <div class="flex items-center justify-between">
                <div>
                  <p class="font-semibold">Article #{articles.indexOf(article) + 1}</p>
                  <p class="text-xs opacity-60">{article.word_count || wordCount(article.content)} words · {readTime(article.content)} min read</p>
                </div>
                <span class="badge badge-sm">{article.status}</span>
              </div>
            </div>
          </div>
        {:else}
          <div class="text-center py-8 opacity-60">No articles yet. Approve a brief first, then generate an article.</div>
        {/each}
      </div>
    {/if}
  {:else}
    <div class="flex items-center gap-2 mb-4">
      <button class="btn btn-ghost btn-sm" on:click={() => { activeArticle = null; }}>← Back</button>
      <div class="flex-1"></div>
      <span class="text-sm opacity-60">{wordCount(editContent)} words · {readTime(editContent)} min</span>
      <div class="dropdown dropdown-end">
        <label tabindex="0" class="btn btn-ghost btn-sm">Export</label>
        <ul tabindex="0" class="dropdown-content menu p-2 shadow bg-base-200 rounded-box w-40">
          <li><button on:click={() => copyContent(editContent, 'md')}>Copy Markdown</button></li>
          <li><button on:click={() => copyContent(editContent, 'html')}>Copy HTML</button></li>
        </ul>
      </div>
      {#if activeArticle.status === 'draft'}
        <button class="btn btn-success btn-sm" on:click={() => changeStatus(activeArticle, 'reviewed')}>Mark Reviewed</button>
      {:else if activeArticle.status === 'reviewed'}
        <button class="btn btn-primary btn-sm" on:click={() => changeStatus(activeArticle, 'approved')}>Approve</button>
      {/if}
      <button class="btn btn-primary btn-sm" on:click={saveArticle}>Save</button>
    </div>

    <textarea
      class="textarea textarea-bordered w-full h-[60vh] font-mono text-sm"
      bind:value={editContent}
    ></textarea>
  {/if}
</div>
