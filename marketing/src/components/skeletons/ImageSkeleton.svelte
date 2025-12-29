<script lang="ts">
	interface Props {
		width?: string;
		height?: string;
		aspectRatio?: string;
		class?: string;
		rounded?: boolean;
	}

	let {
		width = 'w-full',
		height,
		aspectRatio,
		class: className = '',
		rounded = true
	}: Props = $props();

	const style = aspectRatio ? `aspect-ratio: ${aspectRatio};` : '';
</script>

<div
	class="bg-[#2d2826] {width} {height} {rounded ? 'rounded-lg' : ''} {className}"
	style={style}
	role="status"
	aria-label="Loading image..."
>
	<div class="w-full h-full bg-gradient-to-r from-[#2d2826] via-[#3c3836] to-[#2d2826] bg-[length:200%_100%] animate-shimmer {rounded ? 'rounded-lg' : ''}"></div>
</div>

<style>
	@keyframes shimmer {
		0% {
			background-position: -200% 0;
		}
		100% {
			background-position: 200% 0;
		}
	}

	.animate-shimmer {
		animation: shimmer 2s infinite;
	}
</style>
