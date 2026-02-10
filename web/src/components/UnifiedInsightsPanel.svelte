<script>
	import { onMount } from 'svelte';
	import { fetchProjectUnifiedInsights } from '../lib/data.js';
	import { HelpCircle } from 'lucide-svelte';

	export let projectId = null;
	export let onNavigateToIssues = null;

	let loading = false;
	let error = null;
	let insights = [];
	let summary = {};
	let connectedSources = [];
	let filterCategory = 'all';
	let sortBy = 'priority';
	let expandedPages = {};
	let copyFeedback = null;

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

	// Normalize raw score to 0-100 for display (percentile within current dataset)
	function normalizeScore(rawScore, allScores) {
		if (!allScores?.length || rawScore <= 0) return 0;
		const maxScore = Math.max(...allScores);
		if (maxScore <= 0) return 0;
		return Math.round(Math.min(100, (rawScore / maxScore) * 100));
	}

	function getPriorityBadge(score) {
		if (score >= 80) return { class: 'badge-error', label: 'Critical' };
		if (score >= 50) return { class: 'badge-warning', label: 'High' };
		if (score >= 25) return { class: 'badge-info', label: 'Medium' };
		return { class: 'badge-ghost', label: 'Low' };
	}

	function formatNumber(num) {
		if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M';
		if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
		return Math.round(num).toLocaleString();
	}

	function getScoreBreakdown(insight) {
		const parts = [];
		if (insight.issues?.length) {
			const e = insight.issue_severity_counts?.error || 0;
			const w = insight.issue_severity_counts?.warning || 0;
			const i = insight.issue_severity_counts?.info || 0;
			const desc = [e && `${e} error${e > 1 ? 's' : ''}`, w && `${w} warning${w > 1 ? 's' : ''}`, i && `${i} info`].filter(Boolean).join(', ');
			parts.push({ label: 'Crawl issues', value: desc });
		}
		if (insight.gsc_metrics) {
			const imp = insight.gsc_metrics.impressions || 0;
			const pos = insight.gsc_metrics.position;
			let desc = formatNumber(imp) + ' impressions';
			if (pos >= 5 && pos <= 20 && imp > 500) desc += ' (ranking opportunity)';
			parts.push({ label: 'GSC', value: desc });
		}
		if (insight.ga4_metrics) {
			const s = insight.ga4_metrics.sessions || 0;
			const br = insight.ga4_metrics.bounce_rate;
			let desc = formatNumber(s) + ' sessions';
			if (br > 0.7 && s > 500) desc += ' (high bounce)';
			parts.push({ label: 'GA4', value: desc });
		}
		if (insight.clarity_metrics) {
			const r = insight.clarity_metrics.rage_click_count || 0;
			const d = insight.clarity_metrics.dead_click_count || 0;
			if (r > 0 || d > 0) {
				parts.push({ label: 'Clarity', value: `${r} rage + ${d} dead clicks` });
			}
		}
		return parts;
	}

	async function copyUrl(insightUrl) {
		const fullUrl = insightUrl.startsWith('http') ? insightUrl : `https://${insightUrl}`;
		try {
			await navigator.clipboard.writeText(fullUrl);
			copyFeedback = insightUrl;
			setTimeout(() => (copyFeedback = null), 1500);
		} catch {
			copyFeedback = 'error';
		}
	}

	async function copyRationale(insight) {
		const r = insight.rationale;
		if (!r) return;
		const lines = [
			`Why this page matters: ${r.why_this_matters}`,
			r.what_informed_priority?.length
				? `Data informed by: ${r.what_informed_priority.join('; ')}`
				: null,
			r.deprioritized_context || null,
			`Risk of not fixing: ${r.risk_of_not_fixing}`
		].filter(Boolean);
		try {
			await navigator.clipboard.writeText(lines.join('\n\n'));
			copyFeedback = `rationale-${insight.url}`;
			setTimeout(() => (copyFeedback = null), 1500);
		} catch {
			copyFeedback = 'error';
		}
	}

	function handleViewIssues(url) {
		const fullUrl = url.startsWith('http') ? url : `https://${url}`;
		if (onNavigateToIssues) onNavigateToIssues(fullUrl);
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

	$: rawScores = filteredInsights.map((i) => i.priority_score || 0);
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
			<span>
				No insights available yet.
				{#if connectedSources.length === 0}
					Connect GSC, GA4, or Clarity in Project Settings, then sync and run a crawl to generate
					cross-referenced recommendations.
				{:else}
					Sync your connected integrations (GSC/GA4 dashboards) and run a crawl to pull data.
					Insights combine crawl issues with search traffic and UX metrics.
				{/if}
			</span>
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
					<div class="flex items-center gap-1.5">
						<h2 class="card-title text-sm text-base-content/70">Total Priority</h2>
						<div
							class="tooltip tooltip-right"
							data-tip="Sum of all page priority scores. Higher = more combined impact across your site."
						>
							<HelpCircle class="w-4 h-4 text-base-content/50 cursor-help" />
						</div>
					</div>
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
				{@const rawScore = insight.priority_score || 0}
				{@const displayScore = normalizeScore(rawScore, rawScores)}
				{@const priority = getPriorityBadge(displayScore)}
				<div class="card bg-base-100 shadow">
					<div class="card-body p-4">
						<!-- Header Row -->
						<div class="flex items-start justify-between gap-3">
							<div class="flex-1 min-w-0">
								<div class="flex items-center gap-2 flex-wrap">
									<span class="badge {priority.class} badge-sm">{priority.label}</span>
									<div class="flex items-center gap-1">
										<span class="text-xs text-base-content/50">
											{displayScore}/100
										</span>
										<div
											class="tooltip tooltip-right"
											data-tip="Higher = fix sooner. Based on crawl issues, traffic (GSC/GA4), and UX signals (Clarity)."
										>
											<HelpCircle class="w-3.5 h-3.5 text-base-content/40 cursor-help" />
										</div>
									</div>
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

						<!-- Expanded Details -->
						{#if expandedPages[insight.url]}
							{@const breakdown = getScoreBreakdown(insight)}
							{@const r = insight.rationale}
							<div class="mt-3 pt-3 border-t border-base-200 space-y-4">
								<!-- Decision Rationale (JTBD: explainable priorities) -->
								{#if r}
									<div class="p-3 rounded-lg bg-primary/5 border border-primary/20">
										<div class="flex items-center justify-between gap-2 mb-2">
											<h4 class="text-sm font-semibold text-primary">
												Why we're showing this
											</h4>
											<button
												class="btn btn-ghost btn-xs"
												on:click={() => copyRationale(insight)}
											>
												{#if copyFeedback === `rationale-${insight.url}`}
													Copied!
												{:else}
													Copy for client
												{/if}
											</button>
										</div>
										<div class="space-y-2 text-sm text-base-content/80">
											<p><strong>Why this matters:</strong> {r.why_this_matters}</p>
											{#if r.what_informed_priority?.length > 0}
												<p>
													<strong>Data that informed priority:</strong>
													{r.what_informed_priority.join('; ')}
												</p>
											{/if}
											{#if r.deprioritized_context}
												<p class="text-base-content/70 italic">
													{r.deprioritized_context}
												</p>
											{/if}
											<p><strong>Risk of not fixing:</strong> {r.risk_of_not_fixing}</p>
										</div>
									</div>
								{/if}

								<!-- Issues with message + recommendation (lead with actionable) -->
								{#if insight.issues?.length > 0}
									<div>
										<h4 class="text-sm font-semibold mb-2">
											{insight.issues.length} issue{insight.issues.length > 1 ? 's' : ''} to fix
										</h4>
										<div class="space-y-2">
											{#each insight.issues as issue}
												<div class="text-sm p-2 rounded bg-base-200/50">
													<span
														class="badge badge-xs mb-1"
														class:badge-error={issue.severity === 'error'}
														class:badge-warning={issue.severity === 'warning'}
														class:badge-info={issue.severity === 'info'}
													>
														{issue.type?.replace(/_/g, ' ') || 'issue'}
													</span>
													{#if issue.message}
														<p class="text-base-content/80">{issue.message}</p>
													{/if}
													{#if issue.recommendation}
														<p class="text-primary font-medium mt-0.5">
															→ {issue.recommendation}
														</p>
													{/if}
												</div>
											{/each}
										</div>
									</div>
								{:else if insight.recommendations?.length > 0}
									<!-- When no crawl issues but we have metric-based recs (GA4, Clarity, GSC) -->
									<div>
										<h4 class="text-sm font-semibold mb-2 text-primary">How to fix</h4>
										<ul class="space-y-2">
											{#each insight.recommendations as rec}
												<li class="flex gap-2 text-sm">
													<span class="text-success shrink-0">✓</span>
													<span class="text-base-content/90">{rec}</span>
												</li>
											{/each}
										</ul>
									</div>
								{:else}
									<div class="text-sm text-base-content/60">
										<p>No specific fixes available yet. Run a crawl to detect SEO issues on this page.</p>
									</div>
								{/if}

								<!-- Next steps: View in Issues + Copy URL -->
								<div>
									<h4 class="text-sm font-semibold mb-2">Next steps</h4>
									<div class="flex flex-wrap gap-2">
										{#if onNavigateToIssues && insight.issues?.length > 0}
											<button
												class="btn btn-sm btn-primary"
												on:click={() => handleViewIssues(insight.url)}
											>
												View full details in Issues tab
											</button>
										{/if}
										<button
											class="btn btn-sm btn-ghost"
											on:click={() => copyUrl(insight.url)}
										>
											{#if copyFeedback === insight.url}
												Copied!
											{:else}
												Copy page URL
											{/if}
										</button>
									</div>
									{#if insight.issues?.length > 0 && onNavigateToIssues}
										<p class="text-xs text-base-content/60 mt-1">
											The Issues tab shows each fix with context and links to the page.
										</p>
									{/if}
								</div>

								<!-- Score breakdown (collapsible context) -->
								{#if breakdown.length > 0}
									<details class="text-xs">
										<summary class="cursor-pointer text-base-content/60">What drove this score</summary>
										<ul class="mt-1 space-y-0.5 text-base-content/70">
											{#each breakdown as part}
												<li><span class="font-medium">{part.label}:</span> {part.value}</li>
											{/each}
										</ul>
									</details>
								{/if}
							</div>
						{/if}
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
