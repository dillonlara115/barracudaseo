<script lang="ts">
	import TextSkeleton from './TextSkeleton.svelte';

	interface Props {
		class?: string;
		showImage?: boolean;
		lines?: number;
	}

	let {
		class: className = '',
		showImage = false,
		lines = 3
	}: Props = $props();
</script>

<div
	class="bg-[#3c3836] rounded-lg p-6 border border-white/10 animate-pulse {className}"
	role="status"
	aria-label="Loading card..."
>
	{#if showImage}
		<div class="w-full h-48 bg-[#2d2826] rounded-lg mb-4"></div>
	{/if}
	<div class="h-6 bg-[#2d2826] rounded w-3/4 mb-4"></div>
	<TextSkeleton {lines} />
</div>
