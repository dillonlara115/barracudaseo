<script>
	import {
		summarizeReading,
		rate,
		ratingClass,
		formatLCP,
		formatINP,
		formatCLS
	} from '../lib/webVitals.js';

	// The crawl record. Its meta holds the readings taken during the crawl.
	export let crawl = null;

	$: readings = (crawl?.meta?.web_vitals || []).map(summarizeReading);
	$: measured = readings.filter((r) => !r.error);
	$: failed = readings.length - measured.length;
	$: passing = measured.filter((r) => r.overall === 'good').length;
	$: rendering = crawl?.meta?.rendering || null;
	$: strategy = crawl?.meta?.web_vitals?.[0]?.vitals?.strategy || 'mobile';

	function pagePath(url) {
		try {
			const parsed = new URL(url);
			return parsed.pathname + parsed.search || '/';
		} catch {
			return url;
		}
	}

	function overallLabel(overall) {
		if (overall === 'good') return 'Passes';
		if (overall === 'needs-improvement') return 'Needs work';
		if (overall === 'poor') return 'Fails';
		return 'No data';
	}

	function overallBadge(overall) {
		if (overall === 'good') return 'badge-success';
		if (overall === 'needs-improvement') return 'badge-warning';
		if (overall === 'poor') return 'badge-error';
		return 'badge-ghost';
	}
</script>

{#if readings.length > 0 || rendering}
	<div class="card bg-base-200 shadow-lg">
		<div class="card-body">
			<h3 class="card-title">Core Web Vitals</h3>

			{#if measured.length > 0}
				<p class="text-sm text-base-content/70">
					{passing} of {measured.length} key pages pass, measured with Google PageSpeed Insights
					({strategy}).
				</p>

				<div class="overflow-x-auto">
					<table class="table table-sm">
						<thead>
							<tr>
								<th>Page</th>
								<th>Result</th>
								<th>LCP</th>
								<th>INP</th>
								<th>CLS</th>
								<th>Score</th>
								<th>Data</th>
							</tr>
						</thead>
						<tbody>
							{#each measured as reading}
								<tr>
									<td class="max-w-xs break-all" title={reading.url}>{pagePath(reading.url)}</td>
									<td>
										<span class="badge badge-sm {overallBadge(reading.overall)}">
											{overallLabel(reading.overall)}
										</span>
									</td>
									<td class="font-semibold {ratingClass(rate('lcp', reading.lcp))}">
										{formatLCP(reading.lcp)}
									</td>
									<td class="font-semibold {ratingClass(rate('inp', reading.inp))}">
										{formatINP(reading.inp)}
									</td>
									<td class="font-semibold {ratingClass(rate('cls', reading.cls))}">
										{formatCLS(reading.cls)}
									</td>
									<td>{reading.score ?? '—'}</td>
									<td class="text-base-content/70">{reading.source}</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>

				<p class="text-xs text-base-content/70">
					Good: LCP up to 2.5 s, INP up to 200 ms, CLS up to 0.1. A lab test cannot measure INP.
				</p>
			{:else if readings.length > 0}
				<p class="text-base-content/70">
					None of the {readings.length} key pages could be measured this time.
				</p>
			{/if}

			{#if failed > 0 && measured.length > 0}
				<p class="text-sm text-base-content/70">
					{failed} page{failed === 1 ? '' : 's'} could not be measured.
				</p>
			{/if}

			{#if rendering}
				<p class="text-sm text-base-content/70">
					JavaScript rendering: {rendering.rendered} of {rendering.selected} key pages were loaded in
					headless Chrome and compared with their raw HTML. Differences are listed under Issues.
				</p>
			{/if}
		</div>
	</div>
{/if}
