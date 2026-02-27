<script>
	import { link, location, querystring } from 'svelte-spa-router';
	import {
		LayoutDashboard,
		FileText,
		AlertTriangle,
		Lightbulb,
		Network,
		TrendingUp,
		ScanSearch,
		Target,
		Search,
		Settings,
		BarChart,
		FileSearch,
		Terminal,
		BrainCircuit,
		PenTool,
		Link2,
		Gauge,
		ChevronDown
	} from 'lucide-svelte';

	/** @type {string|null} */
	export let projectId = null;

	/**
	 * In "crawl" mode this is the active tab name (e.g. "dashboard").
	 * In "project" mode we derive active state from $location.
	 */
	export let activeTab = 'dashboard';

	/** @type {any} */
	export let gscStatus = null;
	/** @type {any} */
	export let ga4Status = null;
	/** @type {any} */
	export let clarityStatus = null;

	/**
	 * Callback for tab-based navigation in crawl mode.
	 * null in project mode.
	 * @type {((tab: string, filters?: any) => void)|null}
	 */
	export let navigateToTab = null;

	/** "crawl" = Dashboard.svelte tabs, "project" = route-based pages */
	export let mode = 'project';

	$: currentPath = $location || '';

	// In crawl mode, derive the active tab from querystring for accurate state
	$: computedActiveTab = mode === 'crawl' && $querystring
		? new URLSearchParams($querystring).get('tab') || activeTab
		: activeTab;

	// Integration connection flags
	$: gscConnected = !!(gscStatus && gscStatus.integration && gscStatus.integration.property_url);
	$: ga4Connected = !!(ga4Status && ga4Status.integration && ga4Status.integration.property_id);
	$: clarityConnected = !!(clarityStatus && clarityStatus.integration && clarityStatus.integration.connected);

	// Active-state helpers for project mode
	$: isSettingsPage = currentPath.includes(`/project/${projectId}/settings`) && !currentPath.includes('/settings/cli');
	$: isCLISettings = currentPath.includes('/settings/cli');

	// Section definitions
	const SECTION_KEYS = ['crawl-analysis', 'keywords-rankings', 'search-console', 'analytics', 'ai-tools', 'settings'];

	function loadSectionState(key) {
		try {
			const val = localStorage.getItem(`sidebar-${key}-open`);
			return val === null ? true : val === 'true';
		} catch {
			return true;
		}
	}

	let sections = {};
	SECTION_KEYS.forEach((key) => {
		sections[key] = loadSectionState(key);
	});

	function toggleSection(key) {
		sections[key] = !sections[key];
		try {
			localStorage.setItem(`sidebar-${key}-open`, String(sections[key]));
		} catch {
			// localStorage unavailable
		}
	}

	// Auto-expand the section that contains the active item
	$: {
		if (mode === 'crawl') {
			const crawlTabs = ['dashboard', 'results', 'issues', 'recommendations', 'graph'];
			const gscTabs = ['gsc-dashboard', 'gsc-keywords'];
			const analyticsTabs = ['ga4-dashboard', 'clarity-dashboard'];

			if (crawlTabs.includes(activeTab) || activeTab === 'crawls') {
				if (!sections['crawl-analysis']) {
					sections['crawl-analysis'] = true;
					try { localStorage.setItem('sidebar-crawl-analysis-open', 'true'); } catch {}
				}
			}
			if (gscTabs.includes(activeTab)) {
				if (!sections['search-console']) {
					sections['search-console'] = true;
					try { localStorage.setItem('sidebar-search-console-open', 'true'); } catch {}
				}
			}
			if (analyticsTabs.includes(activeTab)) {
				if (!sections['analytics']) {
					sections['analytics'] = true;
					try { localStorage.setItem('sidebar-analytics-open', 'true'); } catch {}
				}
			}
		} else {
			if (currentPath.includes('/rank-tracker') || currentPath.includes('/discover-keywords')) {
				if (!sections['keywords-rankings']) {
					sections['keywords-rankings'] = true;
					try { localStorage.setItem('sidebar-keywords-rankings-open', 'true'); } catch {}
				}
			}
			if (currentPath.includes('/gsc') || currentPath.includes('/gsc-intelligence')) {
				if (!sections['search-console']) {
					sections['search-console'] = true;
					try { localStorage.setItem('sidebar-search-console-open', 'true'); } catch {}
				}
			}
			if (currentPath.includes('tab=ga4') || currentPath.includes('tab=clarity')) {
				if (!sections['analytics']) {
					sections['analytics'] = true;
					try { localStorage.setItem('sidebar-analytics-open', 'true'); } catch {}
				}
			}
			if (currentPath.includes('/content') || currentPath.includes('/internal-links') || currentPath.includes('/impact-first') || currentPath.includes('/ai-usage')) {
				if (!sections['ai-tools']) {
					sections['ai-tools'] = true;
					try { localStorage.setItem('sidebar-ai-tools-open', 'true'); } catch {}
				}
			}
			if (currentPath.includes('/settings')) {
				if (!sections['settings']) {
					sections['settings'] = true;
					try { localStorage.setItem('sidebar-settings-open', 'true'); } catch {}
				}
			}
		}
	}

	// Helper: is a given item active?
	function isActive(tabOrPath) {
		if (mode === 'crawl') {
			// In crawl mode, use the computed activeTab from querystring
			return computedActiveTab === tabOrPath;
		}
		// Project mode: match by path segment
		switch (tabOrPath) {
			case 'dashboard':
				return (
					(currentPath === `/project/${projectId}` || currentPath === `/project/${projectId}/`) &&
					!currentPath.includes('?tab=')
				);
			case 'results':
				return currentPath.includes(`/project/${projectId}`) && currentPath.includes('?tab=results');
			case 'issues':
				return currentPath.includes(`/project/${projectId}`) && currentPath.includes('?tab=issues');
			case 'recommendations':
				return currentPath.includes(`/project/${projectId}`) && currentPath.includes('?tab=recommendations');
			case 'graph':
				return currentPath.includes(`/project/${projectId}`) && currentPath.includes('?tab=graph');
			case 'crawls':
				return currentPath.includes('/crawls') && !currentPath.includes('/crawl/');
			case 'rank-tracker':
				return currentPath.includes('/rank-tracker');
			case 'discover-keywords':
				return currentPath.includes('/discover-keywords');
			case 'gsc-dashboard':
				return currentPath.includes('/gsc') && !currentPath.includes('/keywords') && !currentPath.includes('/gsc-intelligence');
			case 'gsc-keywords':
				return currentPath.includes('/gsc/keywords');
			case 'gsc-intelligence':
				return currentPath.includes('/gsc-intelligence');
			case 'ga4-dashboard':
				return currentPath.includes('tab=ga4-dashboard');
			case 'clarity-dashboard':
				return currentPath.includes('tab=clarity-dashboard');
			case 'content':
				return currentPath.includes('/content');
			case 'internal-links':
				return currentPath.includes('/internal-links');
			case 'impact-first':
				return currentPath.includes('/impact-first');
			case 'ai-usage':
				return currentPath.includes('/ai-usage');
			case 'settings':
				return isSettingsPage;
			case 'settings-cli':
				return isCLISettings;
			default:
				return false;
		}
	}

	// Navigation: in crawl mode use navigateToTab, in project mode use links
	// The crawl-view tabs (dashboard, results, issues, recommendations, graph) are
	// buttons in crawl mode but links in project mode.
	// Everything else (crawls, rank-tracker, etc.) is always a link.
	function isCrawlTab(id) {
		return ['dashboard', 'results', 'issues', 'recommendations', 'graph', 'gsc-dashboard', 'gsc-keywords', 'ga4-dashboard', 'clarity-dashboard'].includes(id);
	}

	function getHref(id) {
		switch (id) {
			case 'dashboard':
				return `/project/${projectId}`;
			case 'results':
				return `/project/${projectId}?tab=results`;
			case 'issues':
				return `/project/${projectId}?tab=issues`;
			case 'recommendations':
				return `/project/${projectId}?tab=recommendations`;
			case 'graph':
				return `/project/${projectId}?tab=graph`;
			case 'crawls':
				return `/project/${projectId}/crawls`;
			case 'rank-tracker':
				return `/project/${projectId}/rank-tracker`;
			case 'discover-keywords':
				return `/project/${projectId}/discover-keywords`;
			case 'gsc-dashboard':
				return `/project/${projectId}/gsc`;
			case 'gsc-keywords':
				return `/project/${projectId}/gsc/keywords`;
			case 'gsc-intelligence':
				return `/project/${projectId}/gsc-intelligence`;
			case 'ga4-dashboard':
				return `/project/${projectId}?tab=ga4-dashboard`;
			case 'clarity-dashboard':
				return `/project/${projectId}?tab=clarity-dashboard`;
			case 'content':
				return `/project/${projectId}/content`;
			case 'internal-links':
				return `/project/${projectId}/internal-links`;
			case 'impact-first':
				return `/project/${projectId}/impact-first`;
			case 'ai-usage':
				return `/project/${projectId}/ai-usage`;
			case 'settings':
				return `/project/${projectId}/settings`;
			case 'settings-cli':
				return `/project/${projectId}/settings/cli`;
			default:
				return '#';
		}
	}
</script>

<aside class="sidebar-nav w-full lg:w-64 bg-base-100 lg:border-r border-base-200 flex-shrink-0 lg:overflow-y-auto lg:max-h-[calc(100vh-200px)]">
	<!-- Mobile: flat horizontal scroll -->
	<ul class="menu menu-horizontal lg:hidden p-2 w-full overflow-x-auto whitespace-nowrap space-x-2">
		{#if mode === 'crawl' && navigateToTab}
			<li><button type="button" class:active={isActive('dashboard')} on:click={() => navigateToTab('dashboard')}><LayoutDashboard class="w-5 h-5" />Dashboard</button></li>
			<li><button type="button" class:active={isActive('results')} on:click={() => navigateToTab('results')}><FileText class="w-5 h-5" />Results</button></li>
			<li><button type="button" class:active={isActive('issues')} on:click={() => navigateToTab('issues')}><AlertTriangle class="w-5 h-5" />Issues</button></li>
			<li><button type="button" class:active={isActive('recommendations')} on:click={() => navigateToTab('recommendations')}><Lightbulb class="w-5 h-5" />Recommendations</button></li>
			<li><button type="button" class:active={isActive('graph')} on:click={() => navigateToTab('graph')}><Network class="w-5 h-5" />Link Graph</button></li>
		{:else}
			<li><a href={getHref('dashboard')} use:link class:active={isActive('dashboard')}><LayoutDashboard class="w-5 h-5" />Dashboard</a></li>
			<li><a href={getHref('results')} use:link class:active={isActive('results')}><FileText class="w-5 h-5" />Results</a></li>
			<li><a href={getHref('issues')} use:link class:active={isActive('issues')}><AlertTriangle class="w-5 h-5" />Issues</a></li>
			<li><a href={getHref('recommendations')} use:link class:active={isActive('recommendations')}><Lightbulb class="w-5 h-5" />Recommendations</a></li>
			<li><a href={getHref('graph')} use:link class:active={isActive('graph')}><Network class="w-5 h-5" />Link Graph</a></li>
		{/if}
		{#if projectId}
			<li><a href={getHref('crawls')} use:link class:active={isActive('crawls')}><FileSearch class="w-5 h-5" />Crawls</a></li>
			<li><a href={getHref('rank-tracker')} use:link class:active={isActive('rank-tracker')}><TrendingUp class="w-5 h-5" />Rank Tracker</a></li>
			<li><a href={getHref('settings')} use:link class:active={isActive('settings')}><Settings class="w-5 h-5" />Settings</a></li>
		{/if}
	</ul>

	<!-- Desktop: collapsible sections -->
	<nav class="hidden lg:block p-4 space-y-1">
		<!-- Crawl Analysis -->
		{#if projectId}
			<button
				type="button"
				class="sidebar-section-header"
				on:click={() => toggleSection('crawl-analysis')}
			>
				<span>Crawl Analysis</span>
				<ChevronDown class="w-4 h-4 transition-transform {sections['crawl-analysis'] ? '' : '-rotate-90'}" />
			</button>
			{#if sections['crawl-analysis']}
				<ul class="menu menu-sm p-0 pl-1">
					{#if mode === 'crawl' && navigateToTab}
						<li><button type="button" class:active={isActive('dashboard')} on:click={() => navigateToTab('dashboard')}><LayoutDashboard class="w-4 h-4" />Dashboard</button></li>
						<li><button type="button" class:active={isActive('results')} on:click={() => navigateToTab('results')}><FileText class="w-4 h-4" />Results</button></li>
						<li><button type="button" class:active={isActive('issues')} on:click={() => navigateToTab('issues')}><AlertTriangle class="w-4 h-4" />Issues</button></li>
						<li><button type="button" class:active={isActive('recommendations')} on:click={() => navigateToTab('recommendations')}><Lightbulb class="w-4 h-4" />Recommendations</button></li>
						<li><button type="button" class:active={isActive('graph')} on:click={() => navigateToTab('graph')}><Network class="w-4 h-4" />Link Graph</button></li>
					{:else}
						<li><a href={getHref('dashboard')} use:link class:active={isActive('dashboard')}><LayoutDashboard class="w-4 h-4" />Dashboard</a></li>
						<li><a href={getHref('results')} use:link class:active={isActive('results')}><FileText class="w-4 h-4" />Results</a></li>
						<li><a href={getHref('issues')} use:link class:active={isActive('issues')}><AlertTriangle class="w-4 h-4" />Issues</a></li>
						<li><a href={getHref('recommendations')} use:link class:active={isActive('recommendations')}><Lightbulb class="w-4 h-4" />Recommendations</a></li>
						<li><a href={getHref('graph')} use:link class:active={isActive('graph')}><Network class="w-4 h-4" />Link Graph</a></li>
					{/if}
					<li><a href={getHref('crawls')} use:link class:active={isActive('crawls')}><FileSearch class="w-4 h-4" />Crawls</a></li>
				</ul>
			{/if}
		{/if}

		<!-- Keywords & Rankings -->
		{#if projectId}
			<button
				type="button"
				class="sidebar-section-header"
				on:click={() => toggleSection('keywords-rankings')}
			>
				<span>Keywords & Rankings</span>
				<ChevronDown class="w-4 h-4 transition-transform {sections['keywords-rankings'] ? '' : '-rotate-90'}" />
			</button>
			{#if sections['keywords-rankings']}
				<ul class="menu menu-sm p-0 pl-1">
					<li><a href={getHref('rank-tracker')} use:link class:active={isActive('rank-tracker')}><TrendingUp class="w-4 h-4" />Rank Tracker</a></li>
					<li><a href={getHref('discover-keywords')} use:link class:active={isActive('discover-keywords')}><ScanSearch class="w-4 h-4" />Discover Keywords</a></li>
				</ul>
			{/if}
		{/if}

		<!-- Search Console -->
		{#if projectId && gscConnected}
			<button
				type="button"
				class="sidebar-section-header"
				on:click={() => toggleSection('search-console')}
			>
				<span>Search Console</span>
				<ChevronDown class="w-4 h-4 transition-transform {sections['search-console'] ? '' : '-rotate-90'}" />
			</button>
			{#if sections['search-console']}
				<ul class="menu menu-sm p-0 pl-1">
					{#if mode === 'crawl' && navigateToTab}
						<li><button type="button" class:active={isActive('gsc-dashboard')} on:click={() => navigateToTab('gsc-dashboard')}><BarChart class="w-4 h-4" />GSC Dashboard</button></li>
						<li><button type="button" class:active={isActive('gsc-keywords')} on:click={() => navigateToTab('gsc-keywords')}><Search class="w-4 h-4" />GSC Keywords</button></li>
					{:else}
						<li><a href={getHref('gsc-dashboard')} use:link class:active={isActive('gsc-dashboard')}><BarChart class="w-4 h-4" />GSC Dashboard</a></li>
						<li><a href={getHref('gsc-keywords')} use:link class:active={isActive('gsc-keywords')}><Search class="w-4 h-4" />GSC Keywords</a></li>
					{/if}
					<li><a href={getHref('gsc-intelligence')} use:link class:active={isActive('gsc-intelligence')}><BrainCircuit class="w-4 h-4" />GSC Intelligence</a></li>
				</ul>
			{/if}
		{/if}

		<!-- Analytics -->
		{#if projectId && (ga4Connected || clarityConnected)}
			<button
				type="button"
				class="sidebar-section-header"
				on:click={() => toggleSection('analytics')}
			>
				<span>Analytics</span>
				<ChevronDown class="w-4 h-4 transition-transform {sections['analytics'] ? '' : '-rotate-90'}" />
			</button>
			{#if sections['analytics']}
				<ul class="menu menu-sm p-0 pl-1">
					{#if ga4Connected}
						{#if mode === 'crawl' && navigateToTab}
							<li><button type="button" class:active={isActive('ga4-dashboard')} on:click={() => navigateToTab('ga4-dashboard')}><BarChart class="w-4 h-4" />GA4 Dashboard</button></li>
						{:else}
							<li><a href={getHref('ga4-dashboard')} use:link class:active={isActive('ga4-dashboard')}><BarChart class="w-4 h-4" />GA4 Dashboard</a></li>
						{/if}
					{/if}
					{#if clarityConnected}
						{#if mode === 'crawl' && navigateToTab}
							<li><button type="button" class:active={isActive('clarity-dashboard')} on:click={() => navigateToTab('clarity-dashboard')}><AlertTriangle class="w-4 h-4" />Clarity</button></li>
						{:else}
							<li><a href={getHref('clarity-dashboard')} use:link class:active={isActive('clarity-dashboard')}><AlertTriangle class="w-4 h-4" />Clarity</a></li>
						{/if}
					{/if}
				</ul>
			{/if}
		{/if}

		<!-- AI Tools -->
		{#if projectId}
			<button
				type="button"
				class="sidebar-section-header"
				on:click={() => toggleSection('ai-tools')}
			>
				<span>AI Tools</span>
				<ChevronDown class="w-4 h-4 transition-transform {sections['ai-tools'] ? '' : '-rotate-90'}" />
			</button>
			{#if sections['ai-tools']}
				<ul class="menu menu-sm p-0 pl-1">
					<li><a href={getHref('content')} use:link class:active={isActive('content')}><PenTool class="w-4 h-4" />Content</a></li>
					<li><a href={getHref('internal-links')} use:link class:active={isActive('internal-links')}><Link2 class="w-4 h-4" />Internal Links</a></li>
					<li><a href={getHref('impact-first')} use:link class:active={isActive('impact-first')}><Target class="w-4 h-4" />Impact-First View</a></li>
					<li><a href={getHref('ai-usage')} use:link class:active={isActive('ai-usage')}><Gauge class="w-4 h-4" />AI Usage</a></li>
				</ul>
			{/if}
		{/if}

		<!-- Settings -->
		{#if projectId}
			<button
				type="button"
				class="sidebar-section-header"
				on:click={() => toggleSection('settings')}
			>
				<span>Settings</span>
				<ChevronDown class="w-4 h-4 transition-transform {sections['settings'] ? '' : '-rotate-90'}" />
			</button>
			{#if sections['settings']}
				<ul class="menu menu-sm p-0 pl-1">
					<li><a href={getHref('settings')} use:link class:active={isActive('settings')}><Settings class="w-4 h-4" />Settings</a></li>
					<li><a href={getHref('settings-cli')} use:link class:active={isActive('settings-cli')}><Terminal class="w-4 h-4" />CLI Setup</a></li>
				</ul>
			{/if}
		{/if}
	</nav>
</aside>

<style>
	.sidebar-section-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		width: 100%;
		padding: 0.4rem 0.75rem;
		margin-top: 0.25rem;
		font-size: 0.7rem;
		font-weight: 700;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		color: hsl(var(--bc) / 0.5);
		border-radius: 0.375rem;
		cursor: pointer;
		user-select: none;
		transition: color 0.15s;
	}

	.sidebar-section-header:hover {
		color: hsl(var(--bc) / 0.8);
	}

	.sidebar-section-header:first-child {
		margin-top: 0;
	}

	/* Sidebar link base */
	.sidebar-nav :global(.menu li button),
	.sidebar-nav :global(.menu li a) {
		transition: background-color 0.15s ease, border-color 0.15s ease;
		border-radius: 0.375rem;
	}

	/* Hover: solid darker background for clear visibility */
	.sidebar-nav :global(.menu li button:hover),
	.sidebar-nav :global(.menu li a:hover) {
		background-color: hsl(var(--b2)) !important;
	}

	/* Active: solid background + prominent left border */
	.sidebar-nav :global(.menu li button.active),
	.sidebar-nav :global(.menu li a.active) {
		background-color: hsl(var(--b3)) !important;
		color: hsl(var(--bc)) !important;
		font-weight: 600;
		border-left: 4px solid hsl(var(--p)) !important;
		position: relative;
	}

	.sidebar-nav :global(.menu li button.active:hover),
	.sidebar-nav :global(.menu li a.active:hover) {
		background-color: hsl(var(--b3)) !important;
	}

	.sidebar-nav :global(.menu li button.active svg),
	.sidebar-nav :global(.menu li a.active svg) {
		color: hsl(var(--p));
	}
</style>
