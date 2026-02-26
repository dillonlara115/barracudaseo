<script lang="ts">
	import TextSkeleton from './TextSkeleton.svelte';

	interface Props {
		class?: string;
		showImage?: boolean;
		lines?: number;
	}

	let { class: className = '', showImage = false, lines = 3 }: Props = $props();
</script>

<div
	class="animate-pulse rounded-lg border border-white/10 bg-[#3c3836] p-6 {className}"
	role="status"
	aria-label="Loading card..."
>
	{#if showImage}
		<div class="mb-4 h-48 w-full rounded-lg bg-[#2d2826]"></div>
	{/if}
	<div class="mb-4 h-6 w-3/4 rounded bg-[#2d2826]"></div>
	<TextSkeleton {lines} />
</div>
