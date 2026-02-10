<script>
	import { onMount } from 'svelte';
	import { ExternalLink } from 'lucide-svelte';
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
	let clarityProjectLabel = '';
	let isConnected = false;
	let isConnecting = false;
	let isSyncing = false;
	let error = null;
	let clarityStatus = null;
	let loading = false;
	let lastProjectId = null;

	const CLARITY_SETUP_URL = 'https://learn.microsoft.com/en-us/clarity/setup-and-installation/clarity-data-export-api#obtaining-access-tokens';
	const CLARITY_DASHBOARD_URL = 'https://clarity.microsoft.com';

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

	// Display name: label if set, otherwise project ID
	$: connectedDisplayName = clarityStatus?.integration?.clarity_project_label
		|| clarityStatus?.integration?.clarity_project_id
		|| '';

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
			error = 'Clarity Project ID and API Token are required.';
			return;
		}

		isConnecting = true;
		error = null;

		const result = await connectClarity(projectId, clarityProjectId.trim(), apiToken.trim(), clarityProjectLabel.trim() || null);
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
		if (!confirm(`Disconnect Microsoft Clarity from ${project?.name || 'this project'}? This will remove the stored credentials.`)) return;

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
		clarityProjectLabel = '';
		await initialize();
		loading = false;
	}
</script>

{#if loading}
	<div class="alert alert-info">
		<span class="loading loading-spinner loading-sm"></span>
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
				<p class="text-sm mb-2">
					Clarity provides UX metrics like rage clicks, dead clicks, and scroll depth to help prioritize SEO fixes. Each Barracuda project can connect to a different Clarity project—useful when sites are in different Clarity accounts.
				</p>
				<p class="text-sm font-medium text-base-content/80">How to get your credentials:</p>
				<ol class="text-sm list-decimal list-inside mt-1 space-y-0.5 text-base-content/70">
					<li>
						<a
							href={CLARITY_DASHBOARD_URL}
							target="_blank"
							rel="noopener noreferrer"
							class="link link-hover inline-flex items-center gap-0.5"
						>
							Sign in to Clarity <ExternalLink class="w-3 h-3" />
						</a>
						and open your project
					</li>
					<li>Go to <strong>Settings → Data Export → Generate new API token</strong></li>
					<li>Copy the Project ID (from the project URL) and the generated token</li>
				</ol>
				<a
					href={CLARITY_SETUP_URL}
					target="_blank"
					rel="noopener noreferrer"
					class="link link-hover text-xs mt-2 inline-flex items-center gap-1"
				>
					View Microsoft's setup guide <ExternalLink class="w-3 h-3" />
				</a>
			</div>
		</div>

		<div class="form-control w-full">
			<label class="label" for="clarity-project-id">
				<span class="label-text">Clarity Project ID</span>
			</label>
			<input
				id="clarity-project-id"
				type="text"
				placeholder="e.g. abc123xyz (from the project URL in Clarity)"
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
				placeholder="Paste the token from Settings → Data Export"
				class="input input-bordered w-full"
				bind:value={apiToken}
				disabled={isConnecting}
			/>
			<div class="label">
				<span class="label-text-alt text-base-content/70">
					Only project admins can generate tokens. Store securely—tokens cannot be viewed again after creation.
				</span>
			</div>
		</div>

		<div class="form-control w-full">
			<label class="label" for="clarity-project-label">
				<span class="label-text">Label <span class="text-base-content/50 font-normal">(optional)</span></span>
			</label>
			<input
				id="clarity-project-label"
				type="text"
				placeholder="e.g. blog.example.com or Marketing site"
				class="input input-bordered w-full"
				bind:value={clarityProjectLabel}
				disabled={isConnecting}
			/>
			<div class="label">
				<span class="label-text-alt text-base-content/70">
					Helpful when you have multiple Clarity projects. Shown in the dashboard instead of the Project ID.
				</span>
			</div>
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

		<div class="text-xs text-base-content/50">
			Data covers the last 1–3 days. Clarity limits Data Export API to 10 requests per project per day.
		</div>

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
					Connected to <span class="font-semibold">{connectedDisplayName}</span>
					{#if clarityStatus?.integration?.clarity_project_label && clarityStatus?.integration?.clarity_project_id}
						<span class="text-base-content/50">({clarityStatus.integration.clarity_project_id})</span>
					{/if}
				</div>
				{#if lastSyncedDisplay}
					<div class="text-xs text-base-content/60">Last synced {lastSyncedDisplay}</div>
				{:else}
					<div class="text-xs text-base-content/60">No data yet. Sync to pull engagement metrics.</div>
				{/if}
				<div class="text-xs text-base-content/40 mt-1">
					Data covers last 1–3 days. 10 API calls/day limit.
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

		<div class="pt-4 border-t border-base-200 mt-4">
			<h3 class="text-sm font-bold text-base-content/70 mb-2">Danger Zone</h3>
			<button class="btn btn-error btn-sm" on:click={handleDisconnect} disabled={loading}>
				Disconnect Clarity
			</button>
		</div>
	</div>
{/if}
