<script lang="ts">
	interface Props {
		lines?: number;
		width?: string;
		class?: string;
	}

	let {
		lines = 3,
		width = 'w-full',
		class: className = ''
	}: Props = $props();
</script>

<div class="space-y-2 {className}" role="status" aria-label="Loading content...">
	{#each Array(lines) as _, i}
		<div
			class="h-4 animate-pulse bg-[#2d2826] rounded {width}"
			style="width: {i === lines - 1 ? '75%' : '100%'}"
		></div>
	{/each}
</div>
