<script>
	import { fetchCrawlDiff } from '../lib/data.js';

	export let crawlId = null;

	const LIST_LIMIT = 25;

	let loading = false;
	let error = null;
	let result = null;
	let loadedFor = null;

	$: if (crawlId && crawlId !== loadedFor) {
		loadDiff(crawlId);
	}

	$: diff = result?.available ? result.diff : null;
	$: baselineDate = result?.baseline?.started_at
		? new Date(result.baseline.started_at).toLocaleDateString()
		: null;

	async function loadDiff(id) {
		loadedFor = id;
		loading = true;
		error = null;
		const { data, error: fetchError } = await fetchCrawlDiff(id);
		if (id !== crawlId) return; // a newer crawl was selected while this loaded
		loading = false;

		if (fetchError) {
			error = fetchError.message;
			result = null;
			return;
		}
		result = data;
	}

	function label(type) {
		return type.replace(/_/g, ' ');
	}
</script>

<div class="card bg-base-200 shadow-lg">
	<div class="card-body">
		<h3 class="card-title">Changes since previous crawl</h3>

		{#if loading}
			<div class="flex items-center gap-2 text-base-content/70">
				<span class="loading loading-spinner loading-sm"></span>
				<span>Comparing crawls...</span>
			</div>
		{:else if error}
			<div class="alert alert-error">{error}</div>
		{:else if result && !result.available}
			<p class="text-base-content/70">{result.reason}</p>
		{:else if diff}
			{#if baselineDate}
				<p class="text-sm text-base-content/70">Compared with the crawl from {baselineDate}.</p>
			{/if}

			<div class="grid grid-cols-2 gap-4 md:grid-cols-4">
				<div>
					<div class="text-sm text-base-content/70">New issues</div>
					<div class="text-2xl font-bold" class:text-error={diff.new_count > 0}>
						{diff.new_count}
					</div>
				</div>
				<div>
					<div class="text-sm text-base-content/70">Fixed issues</div>
					<div class="text-2xl font-bold" class:text-success={diff.fixed_count > 0}>
						{diff.fixed_count}
					</div>
				</div>
				<div>
					<div class="text-sm text-base-content/70">Unchanged</div>
					<div class="text-2xl font-bold">{diff.unchanged_count}</div>
				</div>
				<div>
					<div class="text-sm text-base-content/70">Pages</div>
					<div class="text-2xl font-bold">
						+{diff.pages_added_count} / −{diff.pages_removed_count}
					</div>
				</div>
			</div>

			{#if diff.not_recrawled_count > 0}
				<p class="text-sm text-base-content/70">
					{diff.not_recrawled_count} earlier issue{diff.not_recrawled_count === 1 ? ' is' : 's are'} on
					pages this crawl did not reach, so they are not counted as fixed.
				</p>
			{/if}

			{#if diff.new_issues.length > 0}
				<div class="collapse collapse-arrow bg-base-100">
					<input type="checkbox" aria-label="Show new issues" />
					<div class="collapse-title font-semibold">New issues ({diff.new_count})</div>
					<div class="collapse-content space-y-2">
						{#each diff.new_issues.slice(0, LIST_LIMIT) as issue}
							<div class="text-sm">
								<span
									class="badge badge-sm"
									class:badge-error={issue.severity === 'error'}
									class:badge-warning={issue.severity === 'warning'}
									class:badge-info={issue.severity === 'info'}
								>
									{label(issue.type)}
								</span>
								<span class="break-all">{issue.url}</span>
							</div>
						{/each}
						{#if diff.new_count > LIST_LIMIT}
							<div class="text-sm text-base-content/70">
								... and {diff.new_count - LIST_LIMIT} more
							</div>
						{/if}
					</div>
				</div>
			{/if}

			{#if diff.fixed_issues.length > 0}
				<div class="collapse collapse-arrow bg-base-100">
					<input type="checkbox" aria-label="Show fixed issues" />
					<div class="collapse-title font-semibold">Fixed issues ({diff.fixed_count})</div>
					<div class="collapse-content space-y-2">
						{#each diff.fixed_issues.slice(0, LIST_LIMIT) as issue}
							<div class="text-sm">
								<span class="badge badge-sm badge-success">{label(issue.type)}</span>
								<span class="break-all">{issue.url}</span>
							</div>
						{/each}
						{#if diff.fixed_count > LIST_LIMIT}
							<div class="text-sm text-base-content/70">
								... and {diff.fixed_count - LIST_LIMIT} more
							</div>
						{/if}
					</div>
				</div>
			{/if}
		{/if}
	</div>
</div>
