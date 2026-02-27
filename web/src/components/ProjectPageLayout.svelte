<script>
	import { onMount } from 'svelte';
	import { push } from 'svelte-spa-router';
	import { fetchProjects, fetchProjectGSCStatus, fetchProjectGA4Status, fetchProjectClarityStatus } from '../lib/data.js';
	import ProjectsView from './ProjectsView.svelte';
	import ProjectLayout from './ProjectLayout.svelte';

	export let projectId = null;
	export let gscStatus = null;

	let projects = [];
	let project = null;
	let selectedProject = null;
	let loading = true;
	let error = null;
	let gscStatusLoaded = null;
	let ga4StatusLoaded = null;
	let clarityStatusLoaded = null;

	$: currentProjectId = projectId;

	onMount(async () => {
		if (projectId) {
			await loadData();
			await loadIntegrationStatuses();
		}
	});

	$: if (projectId && projectId !== currentProjectId) {
		loadData();
		loadIntegrationStatuses();
	}

	async function loadIntegrationStatuses() {
		if (!projectId) return;

		if (gscStatus === null) {
			try {
				const result = await fetchProjectGSCStatus(projectId);
				if (!result.error && result.data) {
					gscStatusLoaded = result.data;
				}
			} catch (err) {
				console.error('Failed to load GSC status:', err);
			}
		}

		try {
			const result = await fetchProjectGA4Status(projectId);
			if (!result.error && result.data) {
				ga4StatusLoaded = result.data;
			}
		} catch (err) {
			console.error('Failed to load GA4 status:', err);
		}

		try {
			const result = await fetchProjectClarityStatus(projectId);
			if (!result.error && result.data) {
				clarityStatusLoaded = result.data;
			}
		} catch (err) {
			console.error('Failed to load Clarity status:', err);
		}
	}

	$: finalGSCStatus = gscStatus !== null ? gscStatus : gscStatusLoaded;

	async function loadData() {
		if (!projectId) return;

		loading = true;
		try {
			const { data: projectsData, error: projectsError } = await fetchProjects();
			if (projectsError) throw projectsError;
			projects = projectsData || [];

			project = projects.find(p => p.id === projectId);
			selectedProject = project;
			if (!project) {
				error = 'Project not found';
				loading = false;
				return;
			}
		} catch (err) {
			error = err.message;
		} finally {
			loading = false;
		}
	}

	function handleProjectSelect(selectedProject) {
		push(`/project/${selectedProject.id}`);
	}
</script>

{#if loading}
	<div class="flex items-center justify-center min-h-screen">
		<span class="loading loading-spinner loading-lg"></span>
	</div>
{:else if error}
	<div class="flex items-center justify-center min-h-screen">
		<div class="alert alert-error max-w-md">
			<span>Error: {error}</span>
		</div>
	</div>
{:else if project}
	<ProjectsView {projects} {selectedProject} on:select={(e) => handleProjectSelect(e.detail)} />

	<ProjectLayout {projectId} gscStatus={finalGSCStatus} ga4Status={ga4StatusLoaded} clarityStatus={clarityStatusLoaded}>
		<slot></slot>
	</ProjectLayout>
{/if}
