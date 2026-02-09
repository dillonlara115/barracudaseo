<script>
	import { onMount } from 'svelte';
	import { fetchProjectUnifiedInsights } from '../lib/data.js';

	export let projectId = null;

	let loading = false;
	let error = null;
	let insights = [];
	let summary = {};
	let connectedSources = [];
	let filterCategory = 'all';
	let sortBy = 'priority';
	let expandedPages = {};

	onMount(() => {
		if (projectId) {
			loadInsights();
		}
	});

	async function loadInsights() {
		if (!projectId) return;
		loading = true;
		error = null;

		const result = await fetchProjectUnifiedInsights(projectId);
		if (result.error) {
			error = result.error.message || 'Failed to load insights';
			loading = false;
			return;
		}

		insights = result.data?.insights || [];
		summary = result.data?.summary || {};
		connectedSources = result.data?.connected_sources || [];
		loading = false;
	}

	function toggleExpanded(url) {
		expandedPages = { ...expandedPages, [url]: !expandedPages[url] };
	}

	function getPriorityBadge(score) {
		if (score >= 100) return { class: 'badge-error', label: 'Critical' };
		if (score >= 50) return { class: 'badge-warning', label: 'High' };
		if (score >= 20) return { class: 'badge-info', label: 'Medium' };
		return { class: 'badge-ghost', label: 'Low' };
	}

	function formatNumber(num) {
		if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M';
		if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
		return Math.round(num).toLocaleString();
	}

	$: filteredInsights = (() => {
		let filtered = insights;

		if (filterCategory === 'high-priority') {
			filtered = filtered.filter((i) => i.priority_score >= 50);
		} else if (filterCategory === 'frustration') {
			filtered = filtered.filter(
				(i) =>
					i.clarity_metrics &&
					((i.clarity_metrics.rage_click_count || 0) > 0 ||
						(i.clarity_metrics.dead_click_count || 0) > 0)
			);
		} else if (filterCategory === 'traffic') {
			filtered = filtered.filter(
				(i) =>
					(i.ga4_metrics && (i.ga4_metrics.sessions || 0) > 100) ||
					(i.gsc_metrics && (i.gsc_metrics.impressions || 0) > 1000)
			);
		}

		if (sortBy === 'priority') {
			filtered.sort((a, b) => b.priority_score - a.priority_score);
		} else if (sortBy === 'traffic') {
			filtered.sort((a, b) => {
				const aTraffic = (a.ga4_metrics?.sessions || 0) + (a.gsc_metrics?.clicks || 0);
				const bTraffic = (b.ga4_metrics?.sessions || 0) + (b.gsc_metrics?.clicks || 0);
				return bTraffic - aTraffic;
			});
		} else if (sortBy === 'frustration') {
			filtered.sort((a, b) => {
				const aFrustration =
					(a.clarity_metrics?.rage_click_count || 0) +
					(a.clarity_metrics?.dead_click_count || 0);
				const bFrustration =
					(b.clarity_metrics?.rage_click_count || 0) +
					(b.clarity_metrics?.dead_click_count || 0);
				return bFrustration - aFrustration;
			});
		}

		return filtered;
	})();
</script>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div>
			<h2 class="text-2xl font-bold mb-2">Unified Insights</h2>
			<p class="text-sm text-base-content/60">
				Cross-referenced recommendations from crawl data, search console, analytics, and UX
				signals.
			</p>
		</div>
		<button class="btn btn-primary" on:click={loadInsights} disabled={loading}>
			{#if loading}
				<span class="loading loading-spinner loading-sm"></span>
				Loading...
			{:else}
				Refresh
			{/if}
		</button>
	</div>

	<!-- Data Source Badges -->
	<div class="flex gap-2 flex-wrap">
		<span class="text-sm text-base-content/60 self-center">Data sources:</span>
		<span
			class="badge badge-sm"
			class:badge-success={connectedSources.includes('gsc')}
			class:badge-ghost={!connectedSources.includes('gsc')}
		>
			GSC {connectedSources.includes('gsc') ? '' : '(not connected)'}
		</span>
		<span
			class="badge badge-sm"
			class:badge-success={connectedSources.includes('ga4')}
			class:badge-ghost={!connectedSources.includes('ga4')}
		>
			GA4 {connectedSources.includes('ga4') ? '' : '(not connected)'}
		</span>
		<span
			class="badge badge-sm"
			class:badge-success={connectedSources.includes('clarity')}
			class:badge-ghost={!connectedSources.includes('clarity')}
		>
			Clarity {connectedSources.includes('clarity') ? '' : '(not connected)'}
		</span>
	</div>

	{#if loading}
		<div class="flex justify-center items-center py-20">
			<span class="loading loading-spinner loading-lg"></span>
		</div>
	{:else if error}
		<div class="alert alert-error">
			<span>{error}</span>
		</div>
	{:else if insights.length === 0}
		<div class="alert alert-info">
			<span
				>No insights available. Run a crawl and connect integrations to generate cross-referenced
				recommendations.</span
			>
		</div>
	{:else}
		<!-- Summary Cards -->
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title text-sm text-base-content/70">Pages with Issues</h2>
					<p class="text-3xl font-bold">{summary.pages_with_issues || 0}</p>
				</div>
			</div>
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title text-sm text-base-content/70">High-Priority Fixes</h2>
					<p class="text-3xl font-bold text-error">{summary.high_priority_fixes || 0}</p>
				</div>
			</div>
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title text-sm text-base-content/70">Frustration Signals</h2>
					<p class="text-3xl font-bold text-warning">
						{formatNumber(summary.total_frustration_signals || 0)}
					</p>
				</div>
			</div>
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title text-sm text-base-content/70">Opportunity Score</h2>
					<p class="text-3xl font-bold text-primary">
						{formatNumber(summary.opportunity_score || 0)}
					</p>
				</div>
			</div>
		</div>

		<!-- Filters -->
		<div class="flex flex-wrap gap-3">
			<select class="select select-bordered select-sm" bind:value={filterCategory}>
				<option value="all">All Pages</option>
				<option value="high-priority">High Priority Only</option>
				<option value="frustration">UX Frustration Issues</option>
				<option value="traffic">High Traffic Pages</option>
			</select>
			<select class="select select-bordered select-sm" bind:value={sortBy}>
				<option value="priority">Sort by Priority</option>
				<option value="traffic">Sort by Traffic</option>
				<option value="frustration">Sort by Frustration</option>
			</select>
			<span class="text-sm text-base-content/60 self-center">
				{filteredInsights.length} pages
			</span>
		</div>

		<!-- Insight Cards -->
		<div class="space-y-3">
			{#each filteredInsights.slice(0, 50) as insight}
				{@const priority = getPriorityBadge(insight.priority_score)}
				<div class="card bg-base-100 shadow">
					<div class="card-body p-4">
						<!-- Header Row -->
						<div class="flex items-start justify-between gap-3">
							<div class="flex-1 min-w-0">
								<div class="flex items-center gap-2 flex-wrap">
									<span class="badge {priority.class} badge-sm">{priority.label}</span>
									<span class="text-xs text-base-content/50">
										Score: {Math.round(insight.priority_score)}
									</span>
									{#if insight.issues?.length > 0}
										<span class="badge badge-error badge-xs badge-outline">
											{insight.issues.length} issue{insight.issues.length > 1
												? 's'
												: ''}
										</span>
									{/if}

									<!-- Data source indicators -->
									{#each insight.data_sources || [] as source}
										<span
											class="badge badge-xs"
											class:badge-primary={source === 'gsc'}
											class:badge-secondary={source === 'ga4'}
											class:badge-accent={source === 'clarity'}
										>
											{source.toUpperCase()}
										</span>
									{/each}
								</div>
								<p class="text-sm font-medium mt-1 truncate">/{insight.url}</p>

								<!-- Key Metrics Inline -->
								<div class="flex gap-4 mt-1 text-xs text-base-content/60 flex-wrap">
									{#if insight.gsc_metrics}
										<span
											>{formatNumber(
												insight.gsc_metrics.impressions || 0
											)} impressions</span
										>
										<span
											>pos {(insight.gsc_metrics.position || 0).toFixed(1)}</span
										>
									{/if}
									{#if insight.ga4_metrics}
										<span
											>{formatNumber(
												insight.ga4_metrics.sessions || 0
											)} sessions</span
										>
									{/if}
									{#if insight.clarity_metrics}
										{#if (insight.clarity_metrics.rage_click_count || 0) > 0}
											<span class="text-error"
												>{insight.clarity_metrics.rage_click_count} rage
												clicks</span
											>
										{/if}
										{#if (insight.clarity_metrics.dead_click_count || 0) > 0}
											<span class="text-warning"
												>{insight.clarity_metrics.dead_click_count} dead
												clicks</span
											>
										{/if}
									{/if}
								</div>
							</div>

							<button
								class="btn btn-ghost btn-xs"
								on:click={() => toggleExpanded(insight.url)}
							>
								{expandedPages[insight.url] ? 'Hide' : 'Details'}
							</button>
						</div>

						<!-- Expanded Recommendations -->
						{#if expandedPages[insight.url]}
							<div class="mt-3 pt-3 border-t border-base-200 space-y-2">
								{#if insight.recommendations?.length > 0}
									<h4 class="text-sm font-semibold">Recommendations</h4>
									<ul class="space-y-1">
										{#each insight.recommendations as rec}
											<li class="text-sm text-base-content/80 pl-4 relative">
												<span class="absolute left-0">-</span>
												{rec}
											</li>
										{/each}
									</ul>
								{/if}

								{#if insight.issues?.length > 0}
									<h4 class="text-sm font-semibold mt-2">Issues</h4>
									<div class="flex flex-wrap gap-1">
										{#each insight.issues as issue}
											<span
												class="badge badge-xs"
												class:badge-error={issue.severity === 'error'}
												class:badge-warning={issue.severity === 'warning'}
												class:badge-info={issue.severity === 'info'}
											>
												{issue.type}
											</span>
										{/each}
									</div>
								{/if}
							</div>
						{/if}
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
