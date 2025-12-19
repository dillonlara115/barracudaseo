<script lang="ts">
	import Hero from '../components/sections/Hero.svelte';
	import VideoWalkthrough from '../components/sections/VideoWalkthrough.svelte';
	import Features from '../components/sections/Features.svelte';
	import Pricing from '../components/sections/Pricing.svelte';
	import FutureUpdates from '../components/sections/FutureUpdates.svelte';
	import FAQ from '../components/sections/FAQ.svelte';
	import MetaTags from '../components/MetaTags.svelte';
	import { getMetaTags, getSoftwareApplicationSchema } from '$lib/meta';

	const meta = getMetaTags({
		title: 'Barracuda SEO',
		description: 'Barracuda is a web-based SEO crawler and auditing tool that helps you discover technical issues, prioritize fixes, and improve your site\'s performance.'
	});

	const structuredData = getSoftwareApplicationSchema();
</script>

<MetaTags config={{ ...meta, structuredData }} />

<Hero />
<VideoWalkthrough />
<Features />
<Pricing />
<FutureUpdates />
<FAQ />
