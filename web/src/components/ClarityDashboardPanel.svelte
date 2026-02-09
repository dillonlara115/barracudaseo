<script>
	import { onMount } from 'svelte';
	import { Bar } from 'svelte-chartjs';
	import {
		Chart,
		CategoryScale,
		LinearScale,
		BarElement,
		Title,
		Tooltip,
		Legend
	} from 'chart.js';
	import {
		fetchProjectClarityDimensions,
		triggerProjectClaritySync
	} from '../lib/data.js';

	Chart.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend);

	export let projectId = null;
	export let clarityStatus = null;
	export let clarityLoading = false;
	export let clarityRefreshing = false;
	export let clarityError = null;
	export let onRefresh = null;

	let loading = false;
	let error = null;

	let urlRows = [];
	let deviceRows = [];
	let sourceRows = [];

	let frustrationChartData = null;
	let deviceChartData = null;

	onMount(() => {
		if (projectId && clarityStatus?.integration?.connected) {
			loadData();
		}
	});

	$: if (projectId && clarityStatus?.integration?.connected && !loading && !urlRows.length) {
		loadData();
	}

	async function loadData() {
		if (!projectId) return;

		loading = true;
		error = null;

		try {
			await Promise.all([
				loadDimension('url', (data) => {
					urlRows = data;
					prepareFrustrationChart();
				}),
				loadDimension('device', (data) => {
					deviceRows = data;
					prepareDeviceChart();
				}),
				loadDimension('source', (data) => {
					sourceRows = data;
				})
			]);

			loading = false;
		} catch (err) {
			error = err.message || 'Failed to load Clarity data';
			loading = false;
		}
	}

	async function loadDimension(type, callback) {
		try {
			const result = await fetchProjectClarityDimensions(projectId, type, { limit: 1000 });
			if (!result.error && result.data?.rows) {
				callback(result.data.rows);
			}
		} catch (err) {
			console.error(`Failed to load ${type} dimension:`, err);
		}
	}

	function prepareFrustrationChart() {
		if (!urlRows.length) return;

		// Sort by rage + dead clicks
		const sorted = [...urlRows]
			.map((r) => ({
				...r,
				frustrationScore:
					((r.metrics || {}).rage_click_count || 0) +
					((r.metrics || {}).dead_click_count || 0)
			}))
			.sort((a, b) => b.frustrationScore - a.frustrationScore)
			.slice(0, 10);

		frustrationChartData = {
			labels: sorted.map((r) => {
				const url = r.dimension_value || '';
				return url.length > 40 ? url.substring(0, 37) + '...' : url;
			}),
			datasets: [
				{
					label: 'Rage Clicks',
					data: sorted.map((r) => (r.metrics || {}).rage_click_count || 0),
					backgroundColor: 'rgba(239, 68, 68, 0.8)'
				},
				{
					label: 'Dead Clicks',
					data: sorted.map((r) => (r.metrics || {}).dead_click_count || 0),
					backgroundColor: 'rgba(245, 158, 11, 0.8)'
				}
			]
		};
	}

	function prepareDeviceChart() {
		if (!deviceRows.length) return;

		deviceChartData = {
			labels: deviceRows.map((r) => r.dimension_value || 'Unknown'),
			datasets: [
				{
					label: 'Sessions',
					data: deviceRows.map((r) => (r.metrics || {}).traffic || 0),
					backgroundColor: [
						'rgba(59, 130, 246, 0.8)',
						'rgba(16, 185, 129, 0.8)',
						'rgba(245, 158, 11, 0.8)'
					]
				}
			]
		};
	}

	async function refreshData() {
		if (!projectId || clarityRefreshing) return;
		if (onRefresh) {
			await onRefresh();
		}
		await loadData();
	}

	function formatNumber(num) {
		if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M';
		if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
		return num.toLocaleString();
	}

	function formatPercent(num) {
		return (num * 100).toFixed(1) + '%';
	}

	$: totals = clarityStatus?.summary?.totals || {};
	$: hasData = clarityStatus?.integration?.connected && (urlRows.length > 0 || deviceRows.length > 0);
	$: clarityProjectIdDisplay = clarityStatus?.integration?.clarity_project_id || '';
</script>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div>
			<h2 class="text-2xl font-bold mb-2">Microsoft Clarity Dashboard</h2>
			{#if clarityProjectIdDisplay}
				<p class="text-sm text-base-content/60">Project: {clarityProjectIdDisplay}</p>
			{/if}
		</div>
		<div class="flex gap-2">
			{#if clarityProjectIdDisplay}
				<a
					href="https://clarity.microsoft.com/projects/view/{clarityProjectIdDisplay}/dashboard"
					target="_blank"
					rel="noopener noreferrer"
					class="btn btn-ghost btn-sm"
				>
					Open Clarity
				</a>
			{/if}
			<button
				class="btn btn-primary"
				on:click={refreshData}
				disabled={clarityRefreshing || loading || clarityLoading}
			>
				{#if clarityRefreshing}
					<span class="loading loading-spinner loading-sm"></span>
					Syncing...
				{:else}
					Sync Data
				{/if}
			</button>
		</div>
	</div>

	<div class="alert alert-info text-sm">
		<svg
			xmlns="http://www.w3.org/2000/svg"
			fill="none"
			viewBox="0 0 24 24"
			class="stroke-current shrink-0 w-5 h-5"
		>
			<path
				stroke-linecap="round"
				stroke-linejoin="round"
				stroke-width="2"
				d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
			></path>
		</svg>
		<span>Clarity data covers the last 1-3 days only. API limit: 10 requests/project/day.</span>
	</div>

	{#if clarityLoading || loading}
		<div class="flex justify-center items-center py-20">
			<span class="loading loading-spinner loading-lg"></span>
		</div>
	{:else if clarityError || error}
		<div class="alert alert-error">
			<span>{clarityError || error}</span>
		</div>
	{:else if !clarityStatus?.integration?.connected}
		<div class="alert alert-warning">
			<span
				>Microsoft Clarity is not connected. Configure it in Project Settings with your Clarity
				Project ID and API Token.</span
			>
		</div>
	{:else if !hasData}
		<div class="alert alert-info">
			<span>No Clarity data available yet. Click "Sync Data" to fetch engagement metrics.</span>
		</div>
	{:else}
		<!-- Overview Cards -->
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title text-sm text-base-content/70">Total Sessions</h2>
					<p class="text-3xl font-bold">{formatNumber(totals.traffic || 0)}</p>
				</div>
			</div>
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title text-sm text-base-content/70">Avg Scroll Depth</h2>
					<p class="text-3xl font-bold">{formatPercent(totals.scroll_depth || 0)}</p>
				</div>
			</div>
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title text-sm text-base-content/70">Rage Clicks</h2>
					<p class="text-3xl font-bold text-error">
						{formatNumber(totals.rage_click_count || 0)}
					</p>
				</div>
			</div>
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title text-sm text-base-content/70">Dead Clicks</h2>
					<p class="text-3xl font-bold text-warning">
						{formatNumber(totals.dead_click_count || 0)}
					</p>
				</div>
			</div>
		</div>

		<!-- Frustration Signals Chart -->
		{#if frustrationChartData}
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title mb-4">Frustration Signals by URL</h2>
					<div class="h-72">
						<Bar
							data={frustrationChartData}
							options={{
								responsive: true,
								maintainAspectRatio: false,
								indexAxis: 'y',
								scales: {
									x: { stacked: true },
									y: { stacked: true }
								},
								plugins: {
									legend: { position: 'top' }
								}
							}}
						/>
					</div>
				</div>
			</div>
		{/if}

		<!-- Device Engagement -->
		{#if deviceChartData}
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title mb-4">Engagement by Device</h2>
					<div class="h-64">
						<Bar
							data={deviceChartData}
							options={{
								responsive: true,
								maintainAspectRatio: false,
								plugins: { legend: { display: false } }
							}}
						/>
					</div>
				</div>
			</div>
		{/if}

		<!-- Pages Table -->
		{#if urlRows.length > 0}
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title mb-4">Pages by Engagement Issues</h2>
					<div class="overflow-x-auto">
						<table class="table table-zebra">
							<thead>
								<tr>
									<th>URL</th>
									<th>Sessions</th>
									<th>Scroll Depth</th>
									<th>Rage Clicks</th>
									<th>Dead Clicks</th>
									<th>Quickbacks</th>
								</tr>
							</thead>
							<tbody>
								{#each urlRows.slice(0, 20) as row}
									{@const m = row.metrics || {}}
									<tr>
										<td>
											<span class="text-sm truncate max-w-xs inline-block">
												{row.dimension_value}
											</span>
										</td>
										<td>{formatNumber(m.traffic || 0)}</td>
										<td>{formatPercent(m.scroll_depth || 0)}</td>
										<td>
											{#if (m.rage_click_count || 0) > 0}
												<span class="text-error font-semibold"
													>{m.rage_click_count}</span
												>
											{:else}
												<span class="text-base-content/40">0</span>
											{/if}
										</td>
										<td>
											{#if (m.dead_click_count || 0) > 0}
												<span class="text-warning font-semibold"
													>{m.dead_click_count}</span
												>
											{:else}
												<span class="text-base-content/40">0</span>
											{/if}
										</td>
										<td>
											{#if (m.quickback_click || 0) > 0}
												<span class="text-info font-semibold"
													>{m.quickback_click}</span
												>
											{:else}
												<span class="text-base-content/40">0</span>
											{/if}
										</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				</div>
			</div>
		{/if}
	{/if}
</div>
