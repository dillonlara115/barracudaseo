<script>
  import { link, location } from 'svelte-spa-router';
  import { 
    LayoutDashboard, 
    FileText, 
    AlertTriangle, 
    Lightbulb, 
    Network, 
    TrendingUp, 
    Binoculars, 
    Target,
    Search,
    Settings,
    BarChart,
    FileSearch
  } from 'lucide-svelte';

  export let projectId = null;
  export let gscStatus = null; // Optional: GSC status for conditional display

  // Determine active route based on current location
  $: currentPath = $location || '';
  $: isDashboard = currentPath === `/project/${projectId}` || currentPath === `/project/${projectId}/` || (currentPath.includes(`/project/${projectId}`) && currentPath.includes('?tab=dashboard'));
  $: isResults = currentPath.includes(`/project/${projectId}`) && currentPath.includes('?tab=results');
  $: isIssues = currentPath.includes(`/project/${projectId}`) && currentPath.includes('?tab=issues');
  $: isRecommendations = currentPath.includes(`/project/${projectId}`) && currentPath.includes('?tab=recommendations');
  $: isGraph = currentPath.includes(`/project/${projectId}`) && currentPath.includes('?tab=graph');
  $: isCrawls = currentPath.includes('/crawls') && !currentPath.includes('/crawl/');
  $: isRankTracker = currentPath.includes('/rank-tracker');
  $: isDiscoverKeywords = currentPath.includes('/discover-keywords');
  $: isImpactFirst = currentPath.includes('/impact-first');
  $: isGSCDashboard = currentPath.includes('/gsc') && !currentPath.includes('/keywords');
  $: isGSCKeywords = currentPath.includes('/gsc/keywords');
  $: isSettings = currentPath.includes('/settings');

  // Check if GSC is connected
  $: gscConnected = gscStatus?.integration?.property_url ? true : false;
</script>

<div class="flex flex-col lg:flex-row min-h-[calc(100vh-200px)] bg-base-100 border-t border-base-200">
  <!-- Sidebar Navigation -->
  <aside class="w-full lg:w-64 bg-base-100 lg:border-r border-base-200 flex-shrink-0">
    <ul class="menu menu-horizontal lg:menu-vertical p-2 lg:p-4 w-full overflow-x-auto lg:overflow-visible whitespace-nowrap lg:whitespace-normal space-x-2 lg:space-x-0 lg:space-y-1">
      
      <!-- Core SEO Section -->
      <li>
        <a
          href="/project/{projectId}"
          use:link
          class:active={isDashboard}
        >
          <LayoutDashboard class="w-5 h-5" />
          Dashboard
        </a>
      </li>
      <li>
        <a
          href="/project/{projectId}?tab=results"
          use:link
          class:active={isResults}
        >
          <FileText class="w-5 h-5" />
          Results
        </a>
      </li>
      <li>
        <a
          href="/project/{projectId}?tab=issues"
          use:link
          class:active={isIssues}
        >
          <AlertTriangle class="w-5 h-5" />
          Issues
        </a>
      </li>
      <li>
        <a
          href="/project/{projectId}?tab=recommendations"
          use:link
          class:active={isRecommendations}
        >
          <Lightbulb class="w-5 h-5" />
          Recommendations
        </a>
      </li>
      <li>
        <a
          href="/project/{projectId}?tab=graph"
          use:link
          class:active={isGraph}
        >
          <Network class="w-5 h-5" />
          Link Graph
        </a>
      </li>
      {#if projectId}
        <li>
          <a
            href="/project/{projectId}/crawls"
            use:link
            class:active={isCrawls}
          >
            <FileSearch class="w-5 h-5" />
            Crawls
          </a>
        </li>
      {/if}
      
      <!-- Rank Tracking Section -->
      {#if projectId}
        <li class="hidden lg:block border-b border-base-200 my-2 pointer-events-none"></li>
        <li>
          <a
            href="/project/{projectId}/rank-tracker"
            use:link
            class:active={isRankTracker}
          >
            <TrendingUp class="w-5 h-5" />
            Rank Tracker
          </a>
        </li>
        <li>
          <a
            href="/project/{projectId}/discover-keywords"
            use:link
            class:active={isDiscoverKeywords}
          >
            <Binoculars class="w-5 h-5" />
            Discover Keywords
          </a>
        </li>
        <li>
          <a
            href="/project/{projectId}/impact-first"
            use:link
            class:active={isImpactFirst}
          >
            <Target class="w-5 h-5" />
            Impact-First View
          </a>
        </li>
      {/if}
      
      <!-- Google Search Console Section (only if connected) -->
      {#if projectId && gscConnected}
        <li class="hidden lg:block border-b border-base-200 my-2 pointer-events-none"></li>
        <li>
          <a
            href="/project/{projectId}/gsc"
            use:link
            class:active={isGSCDashboard}
          >
            <BarChart class="w-5 h-5" />
            GSC Dashboard
          </a>
        </li>
        <li>
          <a
            href="/project/{projectId}/gsc/keywords"
            use:link
            class:active={isGSCKeywords}
          >
            <Search class="w-5 h-5" />
            GSC Keywords
          </a>
        </li>
      {/if}
      
      <!-- Project Settings -->
      {#if projectId}
        <li class="hidden lg:block border-b border-base-200 my-2 pointer-events-none"></li>
        <li>
          <a 
            href="/project/{projectId}/settings" 
            use:link
            class:active={isSettings}
          >
            <Settings class="w-5 h-5" />
            Settings
          </a>
        </li>
      {/if}
    </ul>
  </aside>

  <!-- Main Content -->
  <main class="flex-1 p-4 lg:p-8 overflow-y-auto">
    <slot />
  </main>
</div>

<style>
  :global(.menu a.active) {
    background-color: hsl(var(--p) / 0.1);
    color: hsl(var(--p));
  }
  
  :global(.menu a.active svg) {
    color: hsl(var(--p));
  }
</style>

