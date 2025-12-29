<script lang="ts">
	import '../app.css';
	import Header from '../components/layout/Header.svelte';
	import Footer from '../components/layout/Footer.svelte';
	import MetaTags from '../components/MetaTags.svelte';
	import GoogleAnalytics from '../components/GoogleAnalytics.svelte';
	import PageTransitionLoader from '../components/PageTransitionLoader.svelte';
	import { getMetaTags, getOrganizationSchema, getWebSiteSchema } from '$lib/meta';

	let { children } = $props();

	const defaultMeta = getMetaTags();
	const structuredData = [getOrganizationSchema(), getWebSiteSchema()];
</script>

<MetaTags config={{ ...defaultMeta, structuredData }} />
<GoogleAnalytics />
<PageTransitionLoader />

<div class="min-h-screen flex flex-col" data-theme="barracuda">
	<Header />
	<main class="flex-grow">
		{@render children()}
	</main>
	<Footer />
</div>
