<script lang="ts">
	import ImageSkeleton from './skeletons/ImageSkeleton.svelte';
	import { browser } from '$app/environment';

	interface Props {
		src: string;
		alt: string;
		width?: number;
		height?: number;
		class?: string;
		loading?: 'lazy' | 'eager';
		priority?: boolean;
		aspectRatio?: string;
	}

	let {
		src,
		alt,
		width,
		height,
		class: className = '',
		loading = 'lazy',
		priority = false,
		aspectRatio
	}: Props = $props();

	let imageLoaded = $state(false);
	let imageError = $state(false);

	// Use eager loading for priority images (above the fold)
	const loadingAttr = priority ? 'eager' : loading;

	function handleLoad() {
		imageLoaded = true;
	}

	function handleError() {
		imageError = true;
		imageLoaded = true; // Stop showing skeleton on error
	}
</script>

<div class="relative {className}">
	{#if browser && !imageLoaded && !imageError}
		<ImageSkeleton
			width={width ? `${width}px` : 'w-full'}
			height={height ? `${height}px` : undefined}
			aspectRatio={aspectRatio}
			class="absolute inset-0"
		/>
	{/if}
	<img
		{src}
		{alt}
		{width}
		{height}
		class="{imageLoaded ? 'opacity-100' : 'opacity-0'} transition-opacity duration-300"
		loading={loadingAttr}
		decoding="async"
		onload={handleLoad}
		onerror={handleError}
	/>
</div>
