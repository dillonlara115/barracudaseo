<script lang="ts">
	import Hero from '../components/sections/Hero.svelte';
	import TrustSignal from '../components/sections/TrustSignal.svelte';
	import ProblemSection from '../components/sections/ProblemSection.svelte';
	import TheShift from '../components/sections/TheShift.svelte';
	import HowItWorks from '../components/sections/HowItWorks.svelte';
	import WhyThisIsDifferent from '../components/sections/WhyThisIsDifferent.svelte';
	import Comparison from '../components/sections/Comparison.svelte';
	import WhoItsFor from '../components/sections/WhoItsFor.svelte';
	import SocialProof from '../components/sections/SocialProof.svelte';
	import FinalCTA from '../components/sections/FinalCTA.svelte';
	import MetaTags from '../components/MetaTags.svelte';
	import { getMetaTags, getSoftwareApplicationSchema } from '$lib/meta';

	const meta = getMetaTags({
		title: 'Barracuda SEO — Turn crawl data into decisions',
		description: 'When SEO data gets overwhelming, Barracuda tells you what actually matters. The decision-support layer for technical SEO audits.'
	});

	const structuredData = getSoftwareApplicationSchema();
</script>

<MetaTags config={{ ...meta, structuredData }} />

<Hero />
<TrustSignal />
<ProblemSection />
<TheShift />
<HowItWorks />
<WhyThisIsDifferent />
<Comparison />
<WhoItsFor />
<SocialProof />
<FinalCTA />
