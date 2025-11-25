<script>
  import { onMount } from 'svelte';
  import { viewPublicReport } from '../lib/data.js';
  import { AlertCircle, Lock, Calendar, FileText, ExternalLink } from 'lucide-svelte';

  let loading = true;
  let error = null;
  let reportData = null;
  let password = '';
  let passwordRequired = false;
  let passwordError = null;
  let token = '';

  // Extract token from route params on mount
  onMount(() => {
    // Get token from hash route: #/reports/:token
    const hash = typeof window !== 'undefined' ? window.location.hash : '';
    const match = hash.match(/^#\/reports\/([^\/\?]+)/);
    if (match) {
      token = match[1];
      loadReport();
    } else {
      error = 'Invalid report URL';
      loading = false;
    }
  });

  async function loadReport() {
    if (!token) return;
    
    loading = true;
    error = null;
    passwordError = null;

    const { data, error: err } = await viewPublicReport(token, password || null);
    
    if (err) {
      // Check if password is required
      if (err.message?.includes('Password is required') || err.message?.includes('Invalid password')) {
        passwordRequired = true;
        passwordError = err.message;
      } else {
        error = err.message || 'Failed to load report';
      }
      loading = false;
      return;
    }

    reportData = data;
    passwordRequired = false;
    loading = false;
  }

  async function handlePasswordSubmit() {
    await loadReport();
  }


  // Group issues by type
  function groupIssuesByType(issues) {
    const grouped = {};
    issues.forEach(issue => {
      const type = issue.type || 'unknown';
      if (!grouped[type]) {
        grouped[type] = {
          type,
          severity: issue.severity,
          count: 0,
          issues: []
        };
      }
      grouped[type].count++;
      grouped[type].issues.push(issue);
    });
    return Object.values(grouped);
  }

  // Get severity color
  function getSeverityColor(severity) {
    switch (severity) {
      case 'error': return 'text-error';
      case 'warning': return 'text-warning';
      case 'info': return 'text-info';
      default: return 'text-base-content';
    }
  }

  // Get severity badge
  function getSeverityBadge(severity) {
    switch (severity) {
      case 'error': return 'badge-error';
      case 'warning': return 'badge-warning';
      case 'info': return 'badge-info';
      default: return 'badge';
    }
  }
</script>

<svelte:head>
  <title>{reportData?.report?.title || 'SEO Audit Report'} - Barracuda SEO</title>
</svelte:head>

<div class="min-h-screen bg-base-200">
  {#if loading}
    <div class="flex items-center justify-center min-h-screen">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
  {:else if error && !passwordRequired}
    <div class="container mx-auto px-4 py-16">
      <div class="max-w-2xl mx-auto">
        <div class="alert alert-error">
          <AlertCircle class="w-6 h-6" />
          <span>{error}</span>
        </div>
      </div>
    </div>
  {:else if passwordRequired}
    <div class="container mx-auto px-4 py-16">
      <div class="max-w-md mx-auto">
        <div class="card bg-base-100 shadow-xl">
          <div class="card-body">
            <div class="flex items-center gap-3 mb-4">
              <Lock class="w-8 h-8 text-primary" />
              <h2 class="card-title">Password Required</h2>
            </div>
            <p class="text-base-content/70 mb-4">
              This report is password protected. Please enter the password to view it.
            </p>
            {#if passwordError}
              <div class="alert alert-error mb-4">
                <span>{passwordError}</span>
              </div>
            {/if}
            <div class="form-control">
              <label class="label">
                <span class="label-text">Password</span>
              </label>
              <input
                type="password"
                class="input input-bordered"
                bind:value={password}
                on:keydown={(e) => e.key === 'Enter' && handlePasswordSubmit()}
                placeholder="Enter password"
              />
            </div>
            <div class="card-actions justify-end mt-4">
              <button class="btn btn-primary" on:click={handlePasswordSubmit}>
                View Report
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  {:else if reportData}
    <div class="container mx-auto px-4 py-8">
      <!-- Header -->
      <div class="bg-base-100 rounded-lg shadow-lg p-6 mb-6">
        <div class="flex items-start justify-between">
          <div class="flex-1">
            <!-- Project Info -->
            {#if reportData.project && (reportData.project.name || reportData.project.domain)}
              <div class="mb-4 pb-4 border-b border-base-300">
                {#if reportData.project.name}
                  <h2 class="text-xl font-semibold mb-1">{reportData.project.name}</h2>
                {/if}
                {#if reportData.project.domain}
                  <a 
                    href={reportData.project.domain.startsWith('http') ? reportData.project.domain : `https://${reportData.project.domain}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    class="text-primary hover:underline flex items-center gap-1"
                  >
                    <ExternalLink class="w-4 h-4" />
                    {reportData.project.domain}
                  </a>
                {/if}
              </div>
            {/if}
            
            <!-- Report Title -->
            <h1 class="text-3xl font-bold mb-2">{reportData.report.title || 'SEO Audit Report'}</h1>
            {#if reportData.report.description}
              <p class="text-base-content/70">{reportData.report.description}</p>
            {/if}
            <div class="flex items-center gap-4 mt-4 text-sm text-base-content/60">
              <span>Generated: {new Date(reportData.report.created_at).toLocaleDateString()}</span>
              {#if reportData.report.expires_at}
                <span class="flex items-center gap-1">
                  <Calendar class="w-4 h-4" />
                  Expires: {new Date(reportData.report.expires_at).toLocaleDateString()}
                </span>
              {/if}
            </div>
          </div>
          <div class="badge badge-primary badge-lg ml-4">
            <FileText class="w-4 h-4 mr-1" />
            Public Report
          </div>
        </div>
      </div>

      <!-- Summary Stats -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <div class="stat bg-base-100 rounded-lg shadow">
          <div class="stat-title">Total Pages</div>
          <div class="stat-value text-primary">{reportData.summary?.total_pages || 0}</div>
        </div>
        <div class="stat bg-base-100 rounded-lg shadow">
          <div class="stat-title">Total Issues</div>
          <div class="stat-value text-warning">{reportData.summary?.total_issues || 0}</div>
        </div>
        <div class="stat bg-base-100 rounded-lg shadow">
          <div class="stat-title">Crawl Date</div>
          <div class="stat-value text-sm">
            {reportData.crawl?.started_at 
              ? new Date(reportData.crawl.started_at).toLocaleDateString()
              : 'N/A'}
          </div>
        </div>
      </div>

      <!-- Issues by Type -->
      {#if reportData.issues && reportData.issues.length > 0}
        <div class="bg-base-100 rounded-lg shadow-lg p-6 mb-6">
          <h2 class="text-2xl font-bold mb-4">Issues Found</h2>
          
          {#each groupIssuesByType(reportData.issues) as group}
            <div class="mb-6">
              <div class="flex items-center justify-between mb-2">
                <h3 class="text-lg font-semibold flex items-center gap-2">
                  <span class={getSeverityBadge(group.severity)}>{group.severity}</span>
                  {group.type.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase())}
                </h3>
                <span class="badge badge-outline">{group.count} {group.count === 1 ? 'issue' : 'issues'}</span>
              </div>
              
              <div class="space-y-3">
                {#each group.issues.slice(0, 10) as issue}
                  <div class="border-l-4 border-base-300 pl-4 py-2 bg-base-50 rounded-r">
                    <div class="flex items-start justify-between">
                      <div class="flex-1">
                        {#if issue.url}
                          <a 
                            href={issue.url} 
                            target="_blank" 
                            rel="noopener noreferrer"
                            class="text-sm font-medium text-primary hover:underline flex items-center gap-1 mb-2"
                          >
                            <ExternalLink class="w-4 h-4" />
                            {issue.url}
                          </a>
                        {/if}
                        <p class="font-medium text-base-content">{issue.message}</p>
                        {#if issue.recommendation}
                          <p class="text-sm text-base-content/70 mt-2">{issue.recommendation}</p>
                        {/if}
                      </div>
                    </div>
                  </div>
                {/each}
                {#if group.issues.length > 10}
                  <p class="text-sm text-base-content/60 italic">
                    ... and {group.issues.length - 10} more {group.issues.length - 10 === 1 ? 'issue' : 'issues'} of this type
                  </p>
                {/if}
              </div>
            </div>
          {/each}
        </div>
      {:else}
        <div class="bg-base-100 rounded-lg shadow-lg p-6 mb-6">
          <div class="text-center py-8">
            <p class="text-lg text-base-content/70">No issues found in this crawl.</p>
          </div>
        </div>
      {/if}

      <!-- Footer -->
      <div class="text-center text-sm text-base-content/60 py-4">
        <p>Report generated by <a href="https://barracudaseo.com" class="link link-primary">Barracuda SEO</a></p>
      </div>
    </div>
  {/if}
</div>

<style>
  :global(body) {
    background-color: hsl(var(--b2));
  }
</style>

