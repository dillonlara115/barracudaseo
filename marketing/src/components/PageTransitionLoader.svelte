<script lang="ts">
	import { navigating } from '$app/stores';
	import { browser } from '$app/environment';
</script>

{#if browser && $navigating}
	<div
		class="fixed top-0 right-0 left-0 z-50 h-1 bg-[#2d2826]"
		role="progressbar"
		aria-label="Loading page..."
	>
		<div class="animate-progress h-full bg-[#8ec07c]" style="width: 100%"></div>
	</div>
{/if}

<style>
	@keyframes progress {
		0% {
			transform: translateX(-100%);
		}
		50% {
			transform: translateX(0%);
		}
		100% {
			transform: translateX(100%);
		}
	}

	.animate-progress {
		animation: progress 1.5s ease-in-out infinite;
	}
</style>
