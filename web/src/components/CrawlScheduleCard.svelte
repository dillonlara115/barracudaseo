<script>
	import { fetchCrawlSchedule, saveCrawlSchedule, deleteCrawlSchedule } from '../lib/data.js';

	export let projectId = null;

	const FREQUENCY_LABELS = {
		daily: 'Every day',
		weekly: 'Every week',
		monthly: 'Every month'
	};

	let loading = false;
	let saving = false;
	let error = null;
	let saved = false;
	let schedule = null;
	let allowedFrequencies = [];
	let frequency = 'off';
	let loadedFor = null;

	$: if (projectId && projectId !== loadedFor) {
		loadSchedule(projectId);
	}

	$: canSchedule = allowedFrequencies.length > 0;
	$: currentFrequency = schedule?.enabled ? schedule.frequency : 'off';
	$: dirty = frequency !== currentFrequency;
	$: nextRun =
		schedule?.enabled && schedule.next_run_at
			? new Date(schedule.next_run_at).toLocaleString()
			: null;

	async function loadSchedule(id) {
		loadedFor = id;
		loading = true;
		error = null;
		const { data, error: fetchError } = await fetchCrawlSchedule(id);
		loading = false;

		if (fetchError) {
			error = fetchError.message;
			return;
		}
		schedule = data.schedule;
		allowedFrequencies = data.allowed_frequencies || [];
		frequency = schedule?.enabled ? schedule.frequency : 'off';
	}

	async function handleSave() {
		saving = true;
		saved = false;
		error = null;

		const result =
			frequency === 'off'
				? await deleteCrawlSchedule(projectId)
				: await saveCrawlSchedule(projectId, { frequency, enabled: true });
		saving = false;

		if (result.error) {
			error = result.error.message;
			return;
		}
		schedule = frequency === 'off' ? null : result.data.schedule;
		saved = true;
	}
</script>

<div class="card bg-base-200 shadow-lg">
	<div class="card-body">
		<h3 class="card-title">Scheduled crawls</h3>

		{#if loading}
			<div class="flex items-center gap-2 text-base-content/70">
				<span class="loading loading-spinner loading-sm"></span>
				<span>Loading schedule...</span>
			</div>
		{:else if !canSchedule && !error}
			<p class="text-base-content/70">
				Re-crawl this site automatically and see what changed each time. Scheduled crawls are
				available on paid plans.
			</p>
		{:else}
			<p class="text-sm text-base-content/70">
				Crawls the project's domain automatically, up to your plan's page limit.
			</p>

			<div class="flex flex-wrap items-end gap-4">
				<div class="form-control w-full max-w-xs">
					<label class="label" for="crawl-schedule-frequency">
						<span class="label-text">Frequency</span>
					</label>
					<select
						id="crawl-schedule-frequency"
						class="select select-bordered w-full"
						bind:value={frequency}
						disabled={saving}
					>
						<option value="off">Off</option>
						{#each allowedFrequencies as option}
							<option value={option}>{FREQUENCY_LABELS[option] || option}</option>
						{/each}
					</select>
				</div>
				<button class="btn btn-primary" on:click={handleSave} disabled={!dirty || saving}>
					{#if saving}
						<span class="loading loading-spinner loading-sm"></span>
					{/if}
					Save
				</button>
			</div>

			{#if nextRun}
				<p class="text-sm text-base-content/70">Next crawl: {nextRun}</p>
			{/if}
			{#if saved && !dirty}
				<p class="text-sm text-success">Schedule saved.</p>
			{/if}
		{/if}

		{#if error}
			<div class="alert alert-error">{error}</div>
		{/if}
	</div>
</div>
