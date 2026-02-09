<script>
	import { onMount } from 'svelte';
	import {
		connectClarity,
		disconnectClarity,
		fetchProjectClarityStatus,
		triggerProjectClaritySync
	} from '../lib/data.js';

	export let project = null;
	export let projectId = null;

	let clarityProjectId = '';
	let apiToken = '';
	let isConnected = false;
	let isConnecting = false;
	let isSyncing = false;
	let error = null;
	let clarityStatus = null;
	let loading = false;
	let lastProjectId = null;

	const formatDateTime = (value) => {
		if (!value) return null;
		const date = new Date(value);
		if (Number.isNaN(date.getTime())) return null;
		return `${date.toLocaleDateString()} ${date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}`;
	};

	$: isConnected = Boolean(clarityStatus?.integration?.connected);
	$: lastSyncedDisplay = clarityStatus?.sync_state?.last_synced_at
		? formatDateTime(clarityStatus.sync_state.last_synced_at)
		: null;

	$: if (projectId && projectId !== lastProjectId && !loading) {
		lastProjectId = projectId;
		initialize();
	}

	onMount(() => {
		if (projectId) {
			lastProjectId = projectId;
		}
		initialize();
	});

	async function initialize() {
		if (!projectId) return;
		loading = true;
		error = null;

		try {
			const result = await fetchProjectClarityStatus(projectId);
			if (result.error) {
				error = result.error.message || 'Failed to load Clarity status';
			} else {
				clarityStatus = result.data;
			}
		} catch (err) {
			error = err.message || 'Failed to initialize Clarity integration.';
		} finally {
			loading = false;
		}
	}

	async function handleConnect() {
		if (!clarityProjectId.trim() || !apiToken.trim()) {
			error = 'Both Clarity Project ID and API Token are required.';
			return;
		}

		isConnecting = true;
		error = null;

		const result = await connectClarity(projectId, clarityProjectId.trim(), apiToken.trim());
		if (result.error) {
			error = result.error.message || 'Failed to connect Clarity';
			isConnecting = false;
			return;
		}

		apiToken = '';
		await initialize();
		isConnecting = false;
	}

	async function handleSync() {
		if (!projectId) return;
		isSyncing = true;
		error = null;

		const result = await triggerProjectClaritySync(projectId, { num_days: 3 });
		if (result.error) {
			error = result.error.message || 'Failed to sync Clarity data';
			isSyncing = false;
			return;
		}

		await initialize();
		isSyncing = false;
	}

	async function handleDisconnect() {
		if (!projectId) return;
		if (!confirm('Disconnect Microsoft Clarity from this project?')) return;

		loading = true;
		error = null;

		const result = await disconnectClarity(projectId);
		if (result.error) {
			error = result.error.message || 'Failed to disconnect Clarity';
			loading = false;
			return;
		}

		clarityStatus = null;
		clarityProjectId = '';
		apiToken = '';
		await initialize();
		loading = false;
	}
</script>

{#if loading}
	<div class="alert alert-info">
		<span>Loading Microsoft Clarity status...</span>
	</div>
{:else if !isConnected}
	<div class="space-y-4">
		<div class="alert alert-info">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
				class="stroke-current shrink-0 w-6 h-6"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
				></path>
			</svg>
			<div class="flex-1">
				<div class="font-semibold mb-1">Connect Microsoft Clarity</div>
				<div class="text-sm">
					Enter your Clarity Project ID and API Token to enable engagement analytics. Data
					covers the last 1-3 days with a limit of 10 API requests per project per day.
				</div>
			</div>
		</div>

		<div class="form-control w-full">
			<label class="label" for="clarity-project-id">
				<span class="label-text">Clarity Project ID</span>
			</label>
			<input
				id="clarity-project-id"
				type="text"
				placeholder="e.g. abc123def"
				class="input input-bordered w-full"
				bind:value={clarityProjectId}
				disabled={isConnecting}
			/>
		</div>

		<div class="form-control w-full">
			<label class="label" for="clarity-api-token">
				<span class="label-text">API Token</span>
			</label>
			<input
				id="clarity-api-token"
				type="password"
				placeholder="Your Clarity API token"
				class="input input-bordered w-full"
				bind:value={apiToken}
				disabled={isConnecting}
			/>
			<label class="label">
				<span class="label-text-alt text-base-content/70"
					>Generate an API token in Clarity Settings &gt; API Access</span
				>
			</label>
		</div>

		<button
			class="btn btn-primary w-full"
			on:click={handleConnect}
			disabled={isConnecting || !clarityProjectId.trim() || !apiToken.trim()}
		>
			{#if isConnecting}
				<span class="loading loading-spinner loading-sm"></span>
				Connecting...
			{:else}
				Connect Clarity
			{/if}
		</button>

		{#if error}
			<div class="alert alert-error">
				<span>{error}</span>
			</div>
		{/if}
	</div>
{:else}
	<div class="space-y-4">
		<div
			class="flex flex-col md:flex-row md:items-center md:justify-between gap-3 rounded-box border border-base-300 bg-base-100 p-4 shadow-sm"
		>
			<div>
				<div class="text-sm font-semibold text-base-content/80">Microsoft Clarity</div>
				<div class="text-sm">
					Connected to project <span class="font-semibold"
						>{clarityStatus?.integration?.clarity_project_id}</span
					>.
				</div>
				{#if lastSyncedDisplay}
					<div class="text-xs text-base-content/60">Last synced {lastSyncedDisplay}</div>
				{:else}
					<div class="text-xs text-base-content/60">
						No data yet. Sync to pull engagement metrics.
					</div>
				{/if}
				<div class="text-xs text-base-content/40 mt-1">
					Data covers last 1-3 days. 10 API calls/day limit.
				</div>
			</div>
			<div class="flex gap-2">
				<button
					class="btn btn-sm btn-outline"
					on:click={handleSync}
					disabled={isSyncing || loading}
				>
					{#if isSyncing}
						<span class="loading loading-spinner loading-xs"></span>
						Syncing...
					{:else}
						Sync Data
					{/if}
				</button>
			</div>
		</div>

		{#if error}
			<div class="alert alert-error">
				<span>{error}</span>
			</div>
		{/if}

		<!-- Disconnect -->
		<div class="pt-4 border-t border-base-200 mt-4">
			<h3 class="text-sm font-bold text-base-content/70 mb-2">Danger Zone</h3>
			<button class="btn btn-error btn-sm" on:click={handleDisconnect} disabled={loading}>
				Disconnect Clarity
			</button>
		</div>
	</div>
{/if}
