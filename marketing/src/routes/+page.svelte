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
	import { BrainCircuit, ChevronRight } from '@lucide/svelte';

	const meta = getMetaTags({
		title: 'Turn crawl data into decisions',
		description: 'When SEO data gets overwhelming, Barracuda tells you what actually matters. AI-powered decision support for technical SEO audits.'
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

<!-- AI Coming Soon Teaser -->
<section class="bg-[#3c3836] py-16 px-4 border-t border-b border-[#d79921]/20">
	<div class="max-w-4xl mx-auto text-center">
		<div class="inline-flex items-center gap-2 bg-[#d79921]/10 text-[#d79921] px-4 py-2 rounded-full text-sm font-semibold mb-6">
			<BrainCircuit class="w-4 h-4" />
			Coming Soon
		</div>
		<h2 class="text-3xl md:text-4xl font-heading font-bold text-white mb-6">
			AI that knows your site, your voice, and your data.
		</h2>
		<p class="text-lg text-white/70 max-w-2xl mx-auto mb-8">
			Content briefs from GSC data. Articles in your brand voice. Internal link suggestions powered by semantic search. Intelligence at the moment you need it — not a chatbot.
		</p>
		<a href="/roadmap" class="inline-flex items-center gap-2 border-2 border-[#d79921]/30 hover:border-[#d79921] text-white px-6 py-3 rounded-lg font-medium transition-colors">
			See the full AI roadmap
			<ChevronRight class="w-4 h-4" />
		</a>
	</div>
</section>

<FinalCTA />
