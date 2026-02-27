<script>
	import ProjectSidebar from './ProjectSidebar.svelte';

	export let projectId = null;
	/** @type {any} */
	export let gscStatus = null;
	/** @type {any} */
	export let ga4Status = null;
	/** @type {any} */
	export let clarityStatus = null;
</script>

<div class="flex flex-col lg:flex-row min-h-[calc(100vh-200px)] bg-base-100 border-t border-base-200">
	<ProjectSidebar
		{projectId}
		{gscStatus}
		{ga4Status}
		{clarityStatus}
		mode="project"
	/>

	<!-- Main Content -->
	<main class="flex-1 p-4 lg:p-8 overflow-y-auto">
		<slot />
	</main>
</div>
