<script>
  import { onMount } from 'svelte';
  import { initAuth, user } from './lib/auth.js';
  import { fetchProjects, fetchCrawls, fetchPages, fetchIssues } from './lib/data.js';
  import { supabase } from './lib/supabase.js';
  import Auth from './components/Auth.svelte';
  import Dashboard from './components/Dashboard.svelte';
  import ProjectsView from './components/ProjectsView.svelte';
  import ConfigError from './components/ConfigError.svelte';

  let projects = [];
  let selectedProject = null;
  let selectedCrawl = null;
  let summary = null;
  let results = [];
  let loading = true;
  let error = null;
  let configError = null;

  // Check Supabase configuration
  $: {
    const supabaseUrl = import.meta.env.PUBLIC_SUPABASE_URL || import.meta.env.VITE_PUBLIC_SUPABASE_URL;
    const supabaseAnonKey = import.meta.env.PUBLIC_SUPABASE_ANON_KEY || import.meta.env.VITE_PUBLIC_SUPABASE_ANON_KEY;
    
    if (!supabaseUrl || !supabaseAnonKey) {
      configError = 'Missing Supabase configuration. Please set PUBLIC_SUPABASE_URL and PUBLIC_SUPABASE_ANON_KEY environment variables.';
    } else {
      configError = null;
    }
  }

  onMount(async () => {
    // Check for auth callback (email confirmation, password reset, etc.)
    const hashParams = new URLSearchParams(window.location.hash.substring(1));
    const accessToken = hashParams.get('access_token');
    const type = hashParams.get('type');
    
    if (accessToken) {
      // Handle auth callback from email confirmation
      const { data, error } = await supabase.auth.setSession({
        access_token: accessToken,
        refresh_token: hashParams.get('refresh_token') || ''
      });
      
      if (error) {
        console.error('Auth callback error:', error);
      } else {
        // Clear hash from URL
        window.history.replaceState(null, '', window.location.pathname);
      }
    }

    // Initialize auth
    await initAuth();

    // React to auth state changes
    user.subscribe(async (currentUser) => {
      if (currentUser) {
        await loadProjects();
      } else {
        projects = [];
        selectedProject = null;
        summary = null;
        results = [];
      }
      loading = false;
    });
  });

  async function loadProjects() {
    try {
      const { data, error: fetchError } = await fetchProjects();
      if (fetchError) throw fetchError;
      projects = data || [];
      
      // Auto-select first project if available
      if (projects.length > 0 && !selectedProject) {
        selectedProject = projects[0];
        await loadProjectData(projects[0].id);
      }
    } catch (err) {
      error = err.message;
    }
  }

  async function loadProjectData(projectId) {
    if (!projectId) return;

    try {
      loading = true;
      
      // Fetch latest crawl for this project
      const { data: crawls, error: crawlsError } = await fetchCrawls(projectId);
      if (crawlsError) throw crawlsError;

      if (crawls && crawls.length > 0) {
        selectedCrawl = crawls[0];
        await loadCrawlData(crawls[0].id);
      } else {
        // No crawls yet
        summary = null;
        results = [];
        loading = false;
      }
    } catch (err) {
      error = err.message;
      loading = false;
    }
  }

  async function loadCrawlData(crawlId) {
    if (!crawlId) return;

    try {
      // Fetch pages and issues in parallel
      const [pagesResult, issuesResult] = await Promise.all([
        fetchPages(crawlId),
        fetchIssues(crawlId)
      ]);

      if (pagesResult.error) throw pagesResult.error;
      if (issuesResult.error) throw issuesResult.error;

      results = pagesResult.data || [];
      const issues = issuesResult.data || [];

      // Generate summary from data
      summary = {
        total_pages: results.length,
        total_issues: issues.length,
        issues_by_type: {},
        issues: issues.map(issue => ({
          type: issue.type,
          severity: issue.severity,
          url: issue.page_id ? results.find(p => p.id === issue.page_id)?.url || '' : '',
          message: issue.message,
          value: issue.value,
          recommendation: issue.recommendation
        })),
        average_response_time_ms: results.length > 0
          ? Math.round(results.reduce((sum, p) => sum + (p.response_time_ms || 0), 0) / results.length)
          : 0,
        pages_with_errors: results.filter(p => p.status_code >= 400).length
      };

      // Count issues by type
      issues.forEach(issue => {
        summary.issues_by_type[issue.type] = (summary.issues_by_type[issue.type] || 0) + 1;
      });

      loading = false;
    } catch (err) {
      error = err.message;
      loading = false;
    }
  }

  function handleProjectSelect(project) {
    selectedProject = project;
    selectedCrawl = null;
    summary = null;
    results = [];
    loadProjectData(project.id);
  }
</script>

<div class="min-h-screen bg-base-200">
  {#if configError}
    <!-- Show configuration error -->
    <div class="flex items-center justify-center min-h-screen p-4">
      <ConfigError error={configError} />
    </div>
  {:else if !$user}
    <!-- Show auth UI when not logged in -->
    <Auth />
  {:else if loading}
    <!-- Loading state -->
    <div class="flex items-center justify-center min-h-screen">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
  {:else if error}
    <!-- Error state -->
    <div class="flex items-center justify-center min-h-screen">
      <div class="alert alert-error max-w-md">
        <span>Error: {error}</span>
      </div>
    </div>
  {:else}
    <!-- Main app with projects and dashboard -->
    <ProjectsView
      {projects}
      {selectedProject}
      on:select={(e) => handleProjectSelect(e.detail)}
    />
    
    {#if selectedProject && selectedCrawl}
      <Dashboard {summary} {results} />
    {:else if selectedProject}
      <div class="container mx-auto p-4">
        <div class="alert alert-info">
          <span>No crawls found for this project. Upload a crawl to get started.</span>
        </div>
      </div>
    {:else}
      <div class="container mx-auto p-4">
        <div class="alert alert-info">
          <span>No projects yet. Create a project to get started.</span>
        </div>
      </div>
    {/if}
  {/if}
</div>
