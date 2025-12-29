<script lang="ts">
	interface Props {
		src: string;
		alt: string;
		width?: number;
		height?: number;
		class?: string;
		loading?: 'lazy' | 'eager';
		priority?: boolean;
	}

	let {
		src,
		alt,
		width,
		height,
		class: className = '',
		loading = 'lazy',
		priority = false
	}: Props = $props();

	// Use eager loading for priority images (above the fold)
	const loadingAttr = priority ? 'eager' : loading;
</script>

<img
	{src}
	{alt}
	{width}
	{height}
	class={className}
	loading={loadingAttr}
	decoding="async"
/>
