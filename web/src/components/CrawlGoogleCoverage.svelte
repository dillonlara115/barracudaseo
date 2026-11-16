<script>
	import { fetchCrawlGSCCoverage } from '../lib/data.js';

	export let projectId = null;
	export let crawlId = null;

	const LIST_LIMIT = 25;

	let loading = false;
	let error = null;
	let result = null;
	let loadedFor = null;

	$: requestKey = projectId && crawlId ? `${projectId}:${crawlId}` : null;
	$: if (requestKey && requestKey !== loadedFor) {
		loadCoverage(requestKey);
	}

	$: coverage = result?.available ? result.coverage : null;
	$: capturedOn = result?.gsc?.captured_on
		? new Date(result.gsc.captured_on).toLocaleDateString()
		: null;

	async function loadCoverage(key) {
		loadedFor = key;
		loading = true;
		error = null;
		const { data, error: fetchError } = await fetchCrawlGSCCoverage(projectId, crawlId);
		if (key !== requestKey) return; // a newer crawl was selected while this loaded
		loading = false;

		if (fetchError) {
			error = fetchError.message;
			result = null;
			return;
		}
		result = data;
	}

	function formatNumber(value) {
		return Math.round(value || 0).toLocaleString();
	}
</script>

<div class="card bg-base-200 shadow-lg">
	<div class="card-body">
		<h3 class="card-title">Crawl vs Google</h3>

		{#if loading}
			<div class="flex items-center gap-2 text-base-content/70">
				<span class="loading loading-spinner loading-sm"></span>
				<span>Comparing with Search Console...</span>
			</div>
		{:else if error}
			<div class="alert alert-error">{error}</div>
		{:else if result && !result.available}
			<p class="text-base-content/70">{result.reason}</p>
		{:else if coverage}
			<p class="text-sm text-base-content/70">
				What this crawl found, compared with the pages Google showed in search{capturedOn
					? ` (Search Console data from ${capturedOn})`
					: ''}.
			</p>

			<div class="grid grid-cols-2 gap-4 md:grid-cols-4">
				<div>
					<div class="text-sm text-base-content/70">Indexable pages shown by Google</div>
					<div class="text-2xl font-bold">
						{coverage.visible_indexable_pages} / {coverage.indexable_pages}
					</div>
				</div>
				<div>
					<div class="text-sm text-base-content/70">Search traffic at risk</div>
					<div
						class="text-2xl font-bold"
						class:text-error={coverage.search_traffic_at_risk_count > 0}
					>
						{coverage.search_traffic_at_risk_count}
					</div>
				</div>
				<div>
					<div class="text-sm text-base-content/70">In Google, not in crawl</div>
					<div class="text-2xl font-bold">{coverage.not_in_crawl_count}</div>
				</div>
				<div>
					<div class="text-sm text-base-content/70">Indexable, never shown</div>
					<div class="text-2xl font-bold">{coverage.no_search_visibility_count}</div>
				</div>
			</div>

			{#if coverage.search_traffic_at_risk.length > 0}
				<div class="collapse collapse-arrow bg-base-100">
					<input type="checkbox" aria-label="Show pages with search traffic at risk" />
					<div class="collapse-title font-semibold">
						Search traffic at risk ({coverage.search_traffic_at_risk_count})
					</div>
					<div class="collapse-content space-y-2">
						<p class="text-sm text-base-content/70">
							Google is showing these pages, but the crawl found they cannot be indexed. Fix them
							or expect them to drop out of search.
						</p>
						{#each coverage.search_traffic_at_risk.slice(0, LIST_LIMIT) as page}
							<div class="text-sm">
								<span class="badge badge-sm badge-error">{page.class}</span>
								<span class="break-all">{page.url}</span>
								<span class="text-base-content/70">
									{page.reason} · {formatNumber(page.clicks)} clicks,
									{formatNumber(page.impressions)} impressions
								</span>
							</div>
						{/each}
						{#if coverage.search_traffic_at_risk_count > LIST_LIMIT}
							<div class="text-sm text-base-content/70">
								... and {coverage.search_traffic_at_risk_count - LIST_LIMIT} more
							</div>
						{/if}
					</div>
				</div>
			{/if}

			{#if coverage.not_in_crawl.length > 0}
				<div class="collapse collapse-arrow bg-base-100">
					<input type="checkbox" aria-label="Show pages in Google but not in the crawl" />
					<div class="collapse-title font-semibold">
						In Google, not in crawl ({coverage.not_in_crawl_count})
					</div>
					<div class="collapse-content space-y-2">
						<p class="text-sm text-base-content/70">
							Google is showing these pages, but the crawl never reached them. They may be
							orphaned, removed, or URL variants.
							{#if result.crawl_partial}
								This crawl stopped at its page limit, so some may simply not have been reached.
							{/if}
						</p>
						{#each coverage.not_in_crawl.slice(0, LIST_LIMIT) as page}
							<div class="text-sm">
								<span class="break-all">{page.url}</span>
								<span class="text-base-content/70">
									{formatNumber(page.clicks)} clicks, {formatNumber(page.impressions)} impressions
								</span>
							</div>
						{/each}
						{#if coverage.not_in_crawl_count > LIST_LIMIT}
							<div class="text-sm text-base-content/70">
								... and {coverage.not_in_crawl_count - LIST_LIMIT} more
							</div>
						{/if}
					</div>
				</div>
			{/if}

			{#if coverage.no_search_visibility.length > 0}
				<div class="collapse collapse-arrow bg-base-100">
					<input type="checkbox" aria-label="Show indexable pages never shown in search" />
					<div class="collapse-title font-semibold">
						Indexable, never shown ({coverage.no_search_visibility_count})
					</div>
					<div class="collapse-content space-y-2">
						<p class="text-sm text-base-content/70">
							The crawl found these pages indexable, but Google did not show them for any search
							in the period. They may not be indexed, or may not rank for anything yet.
						</p>
						{#each coverage.no_search_visibility.slice(0, LIST_LIMIT) as page}
							<div class="text-sm break-all">{page.url}</div>
						{/each}
						{#if coverage.no_search_visibility_count > LIST_LIMIT}
							<div class="text-sm text-base-content/70">
								... and {coverage.no_search_visibility_count - LIST_LIMIT} more
							</div>
						{/if}
					</div>
				</div>
			{/if}
		{/if}
	</div>
</div>
