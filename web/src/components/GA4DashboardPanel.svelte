<script>
	import { onMount } from 'svelte';
	import { Bar, Line } from 'svelte-chartjs';
	import {
		Chart,
		CategoryScale,
		LinearScale,
		BarElement,
		LineElement,
		PointElement,
		Title,
		Tooltip,
		Legend
	} from 'chart.js';
	import { fetchProjectGA4Dimensions, triggerProjectGA4Sync } from '../lib/data.js';

	Chart.register(
		CategoryScale,
		LinearScale,
		BarElement,
		LineElement,
		PointElement,
		Title,
		Tooltip,
		Legend
	);

	export let projectId = null;
	export let ga4Status = null;
	export let ga4Loading = false;
	export let ga4Refreshing = false;
	export let ga4Error = null;
	export let onRefresh = null;

	let loading = false;
	let error = null;

	let dateRows = [];
	let pageRows = [];
	let sourceRows = [];
	let deviceRows = [];
	let countryRows = [];

	let dateChartData = null;
	let deviceChartData = null;
	let sourceChartData = null;
	let countryChartData = null;

	onMount(() => {
		if (projectId && ga4Status?.integration?.property_id) {
			loadData();
		}
	});

	$: if (projectId && ga4Status?.integration?.property_id && !loading && !dateRows.length) {
		loadData();
	}

	async function loadData() {
		if (!projectId) return;

		loading = true;
		error = null;

		try {
			await Promise.all([
				loadDimension('date', (data) => {
					dateRows = data;
					prepareDateChart();
				}),
				loadDimension('page', (data) => {
					pageRows = data;
				}),
				loadDimension('source', (data) => {
					sourceRows = data;
					prepareSourceChart();
				}),
				loadDimension('device', (data) => {
					deviceRows = data;
					prepareDeviceChart();
				}),
				loadDimension('country', (data) => {
					countryRows = data;
					prepareCountryChart();
				})
			]);

			loading = false;
		} catch (err) {
			error = err.message || 'Failed to load GA4 data';
			loading = false;
		}
	}

	async function loadDimension(type, callback) {
		try {
			const result = await fetchProjectGA4Dimensions(projectId, type, { limit: 1000 });
			if (!result.error && result.data?.rows) {
				let processedRows = deduplicateDimension(result.data.rows);
				callback(processedRows);
			}
		} catch (err) {
			console.error(`Failed to load ${type} dimension:`, err);
		}
	}

	function deduplicateDimension(rows) {
		const dimensionMap = new Map();

		for (const row of rows) {
			const value = row.dimension_value;
			if (!value) continue;

			const metrics = row.metrics || {};
			const existing = dimensionMap.get(value);

			if (existing) {
				const existingMetrics = existing.metrics || {};
				existingMetrics.sessions =
					(existingMetrics.sessions || 0) + (metrics.sessions || 0);
				existingMetrics.users = (existingMetrics.users || 0) + (metrics.users || 0);
				existingMetrics.page_views =
					(existingMetrics.page_views || 0) + (metrics.page_views || 0);
				existingMetrics.conversions =
					(existingMetrics.conversions || 0) + (metrics.conversions || 0);

				const totalSessions =
					(existingMetrics.sessions || 0) + (metrics.sessions || 0);
				if (totalSessions > 0) {
					existingMetrics.bounce_rate =
						((existingMetrics.bounce_rate || 0) *
							((existingMetrics.sessions || 0) - (metrics.sessions || 0)) +
							(metrics.bounce_rate || 0) * (metrics.sessions || 0)) /
						totalSessions;
					existingMetrics.avg_session_duration =
						((existingMetrics.avg_session_duration || 0) *
							((existingMetrics.sessions || 0) - (metrics.sessions || 0)) +
							(metrics.avg_session_duration || 0) * (metrics.sessions || 0)) /
						totalSessions;
				}
			} else {
				dimensionMap.set(value, {
					...row,
					metrics: { ...metrics }
				});
			}
		}

		const deduplicated = Array.from(dimensionMap.values());
		deduplicated.sort((a, b) => {
			return (b.metrics?.sessions || 0) - (a.metrics?.sessions || 0);
		});

		return deduplicated;
	}

	function prepareDateChart() {
		if (!dateRows.length) return;

		const sorted = [...dateRows].sort((a, b) => {
			return (a.dimension_value || '').localeCompare(b.dimension_value || '');
		});

		const labels = sorted.map((r) => {
			const date = r.dimension_value;
			if (!date) return '';
			const parts = date.split('-');
			if (parts.length === 3) return `${parts[1]}/${parts[2]}`;
			return date;
		});

		dateChartData = {
			labels,
			datasets: [
				{
					label: 'Sessions',
					data: sorted.map((r) => Math.round((r.metrics || {}).sessions || 0)),
					borderColor: 'rgb(59, 130, 246)',
					backgroundColor: 'rgba(59, 130, 246, 0.1)',
					yAxisID: 'y'
				},
				{
					label: 'Users',
					data: sorted.map((r) => Math.round((r.metrics || {}).users || 0)),
					borderColor: 'rgb(16, 185, 129)',
					backgroundColor: 'rgba(16, 185, 129, 0.1)',
					yAxisID: 'y'
				},
				{
					label: 'Bounce Rate (%)',
					data: sorted.map((r) => ((r.metrics || {}).bounce_rate || 0) * 100),
					borderColor: 'rgb(245, 158, 11)',
					backgroundColor: 'rgba(245, 158, 11, 0.1)',
					yAxisID: 'y1'
				}
			]
		};
	}

	function prepareDeviceChart() {
		if (!deviceRows.length) return;

		const top = deviceRows.slice(0, 10);
		deviceChartData = {
			labels: top.map((r) => r.dimension_value || 'Unknown'),
			datasets: [
				{
					label: 'Sessions',
					data: top.map((r) => Math.round((r.metrics || {}).sessions || 0)),
					backgroundColor: [
						'rgba(59, 130, 246, 0.8)',
						'rgba(16, 185, 129, 0.8)',
						'rgba(245, 158, 11, 0.8)',
						'rgba(239, 68, 68, 0.8)',
						'rgba(139, 92, 246, 0.8)'
					]
				}
			]
		};
	}

	function prepareSourceChart() {
		if (!sourceRows.length) return;

		const top = sourceRows.slice(0, 10);
		sourceChartData = {
			labels: top.map((r) => r.dimension_value || 'Unknown'),
			datasets: [
				{
					label: 'Sessions',
					data: top.map((r) => Math.round((r.metrics || {}).sessions || 0)),
					backgroundColor: 'rgba(59, 130, 246, 0.8)'
				}
			]
		};
	}

	function prepareCountryChart() {
		if (!countryRows.length) return;

		const top = countryRows.slice(0, 10);
		countryChartData = {
			labels: top.map((r) => r.dimension_value || 'Unknown'),
			datasets: [
				{
					label: 'Sessions',
					data: top.map((r) => Math.round((r.metrics || {}).sessions || 0)),
					backgroundColor: 'rgba(16, 185, 129, 0.8)'
				}
			]
		};
	}

	async function refreshData() {
		if (!projectId || ga4Refreshing) return;
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
		return (num * 100).toFixed(2) + '%';
	}

	function formatDuration(seconds) {
		if (seconds < 60) return Math.round(seconds) + 's';
		const mins = Math.floor(seconds / 60);
		const secs = Math.round(seconds % 60);
		return `${mins}m ${secs}s`;
	}

	$: totals = (() => {
		if (!dateRows.length) {
			return ga4Status?.summary?.totals || {};
		}

		let totalSessions = 0;
		let totalUsers = 0;
		let totalConversions = 0;
		let weightedBounceSum = 0;

		for (const row of dateRows) {
			const m = row.metrics || {};
			const sessions = m.sessions || 0;
			totalSessions += sessions;
			totalUsers += m.users || 0;
			totalConversions += m.conversions || 0;
			weightedBounceSum += (m.bounce_rate || 0) * sessions;
		}

		const avgBounceRate = totalSessions > 0 ? weightedBounceSum / totalSessions : 0;

		return {
			total_sessions: totalSessions,
			total_users: totalUsers,
			total_conversions: totalConversions,
			avg_bounce_rate: avgBounceRate
		};
	})();

	$: hasData =
		ga4Status?.integration?.property_id && (dateRows.length > 0 || pageRows.length > 0);
</script>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div>
			<h2 class="text-2xl font-bold mb-2">Google Analytics 4 Dashboard</h2>
			{#if ga4Status?.integration?.property_name}
				<p class="text-sm text-base-content/60">
					Property: {ga4Status.integration.property_name}
				</p>
			{/if}
		</div>
		<button
			class="btn btn-primary"
			on:click={refreshData}
			disabled={ga4Refreshing || loading || ga4Loading}
		>
			{#if ga4Refreshing}
				<span class="loading loading-spinner loading-sm"></span>
				Refreshing...
			{:else}
				Refresh Data
			{/if}
		</button>
	</div>

	{#if ga4Loading || loading}
		<div class="flex justify-center items-center py-20">
			<span class="loading loading-spinner loading-lg"></span>
		</div>
	{:else if ga4Error || error}
		<div class="alert alert-error">
			<span>{ga4Error || error}</span>
		</div>
	{:else if !ga4Status?.integration?.property_id}
		<div class="alert alert-warning">
			<span
				>Google Analytics 4 is not connected for this project. Connect it in Integrations, then
				select a property in Project Settings.</span
			>
		</div>
	{:else if !hasData}
		<div class="alert alert-info">
			<span
				>No GA4 data available yet. Click "Refresh Data" to sync data from Google Analytics
				4.</span
			>
		</div>
	{:else}
		<!-- Overview Cards -->
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title text-sm text-base-content/70">Total Sessions</h2>
					<p class="text-3xl font-bold">{formatNumber(totals.total_sessions || 0)}</p>
				</div>
			</div>
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title text-sm text-base-content/70">Total Users</h2>
					<p class="text-3xl font-bold">{formatNumber(totals.total_users || 0)}</p>
				</div>
			</div>
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title text-sm text-base-content/70">Avg Bounce Rate</h2>
					<p class="text-3xl font-bold">{formatPercent(totals.avg_bounce_rate || 0)}</p>
				</div>
			</div>
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title text-sm text-base-content/70">Total Conversions</h2>
					<p class="text-3xl font-bold">{formatNumber(totals.total_conversions || 0)}</p>
				</div>
			</div>
		</div>

		<!-- Traffic Over Time -->
		{#if dateChartData}
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title mb-4">Traffic Over Time</h2>
					<div class="h-64">
						<Line
							data={dateChartData}
							options={{
								responsive: true,
								maintainAspectRatio: false,
								scales: {
									y: {
										type: 'linear',
										position: 'left',
										title: { display: true, text: 'Sessions / Users' }
									},
									y1: {
										type: 'linear',
										position: 'right',
										title: { display: true, text: 'Bounce Rate (%)' },
										grid: { drawOnChartArea: false }
									}
								},
								plugins: {
									legend: { position: 'top' },
									tooltip: { mode: 'index', intersect: false }
								}
							}}
						/>
					</div>
				</div>
			</div>
		{/if}

		<!-- Charts Row -->
		<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
			{#if deviceChartData}
				<div class="card bg-base-100 shadow">
					<div class="card-body">
						<h2 class="card-title mb-4">Sessions by Device</h2>
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

			{#if sourceChartData}
				<div class="card bg-base-100 shadow">
					<div class="card-body">
						<h2 class="card-title mb-4">Top Sources</h2>
						<div class="h-64">
							<Bar
								data={sourceChartData}
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
		</div>

		<!-- Top Countries -->
		{#if countryChartData}
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title mb-4">Top Countries</h2>
					<div class="h-64">
						<Bar
							data={countryChartData}
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

		<!-- Top Pages Table -->
		{#if pageRows.length > 0}
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title mb-4">Top Pages</h2>
					<div class="overflow-x-auto">
						<table class="table table-zebra">
							<thead>
								<tr>
									<th>Page</th>
									<th>Sessions</th>
									<th>Users</th>
									<th>Bounce Rate</th>
									<th>Avg Duration</th>
									<th>Conversions</th>
								</tr>
							</thead>
							<tbody>
								{#each pageRows.slice(0, 20) as row}
									{@const metrics = row.metrics || {}}
									<tr>
										<td>
											<span class="text-sm truncate max-w-xs inline-block">
												/{row.dimension_value}
											</span>
										</td>
										<td>{formatNumber(Math.round(metrics.sessions || 0))}</td>
										<td>{formatNumber(Math.round(metrics.users || 0))}</td>
										<td>{formatPercent(metrics.bounce_rate || 0)}</td>
										<td>{formatDuration(metrics.avg_session_duration || 0)}</td>
										<td>{Math.round(metrics.conversions || 0)}</td>
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
