<script lang="ts">
	import '../app.css';
	import Header from '../components/layout/Header.svelte';
	import Footer from '../components/layout/Footer.svelte';
	import GoogleAnalytics from '../components/GoogleAnalytics.svelte';
	import PageTransitionLoader from '../components/PageTransitionLoader.svelte';
	import { getOrganizationSchema, getWebSiteSchema } from '$lib/meta';

	let { children } = $props();

	// Site-wide structured data only. Each page renders its own <MetaTags>,
	// so rendering it here too would duplicate title/description/canonical.
	const siteSchema = JSON.stringify([getOrganizationSchema(), getWebSiteSchema()]);
</script>

<svelte:head>
	{@html `<script type="application/ld+json">${siteSchema}</script>`}
</svelte:head>
<GoogleAnalytics />
<PageTransitionLoader />

<div class="flex min-h-screen flex-col" data-theme="barracuda">
	<Header />
	<main class="flex-grow">
		{@render children()}
	</main>
	<Footer />
</div>
