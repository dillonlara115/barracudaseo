<script lang="ts">
	import { page } from '$app/stores';
	import { SITE_NAME, SITE_URL } from '$lib/constants';
	import type { MetaTagsConfig } from '$lib/meta';

	let { config }: { config: MetaTagsConfig } = $props();

	// Normalize URL: remove trailing slash (except for root), ensure non-www
	const normalizePath = (pathname: string): string => {
		// Remove trailing slash except for root path
		if (pathname !== '/' && pathname.endsWith('/')) {
			return pathname.slice(0, -1);
		}
		return pathname;
	};
	
	// Build full URL for canonical and OG tags (normalized)
	const fullUrl = $derived(`${SITE_URL}${normalizePath($page.url.pathname)}`);
	const ogImage = $derived(
		config.ogImage 
			? (config.ogImage.startsWith('http') ? config.ogImage : `${SITE_URL}${config.ogImage}`)
			: `${SITE_URL}/mockups/barracuda-dashboard.png`
	);
	
	// Handle structured data (can be single object or array)
	const structuredDataArray = $derived(
		Array.isArray(config.structuredData) 
			? config.structuredData 
			: config.structuredData ? [config.structuredData] : []
	);
</script>

<svelte:head>
	<!-- Primary Meta Tags -->
	<title>{config.title}</title>
	<meta name="title" content={config.title} />
	<meta name="description" content={config.description} />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
	
	<!-- Canonical URL -->
	<link rel="canonical" href={fullUrl} />

	<!-- Open Graph / Facebook -->
	<meta property="og:type" content={config.ogType || 'website'} />
	<meta property="og:url" content={fullUrl} />
	<meta property="og:title" content={config.ogTitle || config.title} />
	<meta property="og:description" content={config.ogDescription || config.description} />
	<meta property="og:image" content={ogImage} />
	<meta property="og:site_name" content={SITE_NAME} />
	{#if config.ogImageWidth}
		<meta property="og:image:width" content={config.ogImageWidth.toString()} />
	{/if}
	{#if config.ogImageHeight}
		<meta property="og:image:height" content={config.ogImageHeight.toString()} />
	{/if}

	<!-- Twitter Card -->
	<meta name="twitter:card" content="summary_large_image" />
	<meta name="twitter:url" content={fullUrl} />
	<meta name="twitter:title" content={config.twitterTitle || config.title} />
	<meta name="twitter:description" content={config.twitterDescription || config.description} />
	<meta name="twitter:image" content={ogImage} />
	{#if config.twitterSite}
		<meta name="twitter:site" content={config.twitterSite} />
	{/if}
	{#if config.twitterCreator}
		<meta name="twitter:creator" content={config.twitterCreator} />
	{/if}

	<!-- Additional SEO Meta Tags -->
	{#if config.keywords}
		<meta name="keywords" content={config.keywords} />
	{/if}
	{#if config.author}
		<meta name="author" content={config.author} />
	{/if}
	{#if config.robots}
		<meta name="robots" content={config.robots} />
	{/if}

	<!-- Structured Data (JSON-LD) -->
	{#each structuredDataArray as data}
		{@html `<script type="application/ld+json">${JSON.stringify(data)}</script>`}
	{/each}
</svelte:head>
