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
	import { BrainCircuit, ChevronRight, PenTool, BarChart3, Link2 } from '@lucide/svelte';

	const meta = getMetaTags({
		title: 'Turn crawl data into decisions',
		description:
			'When SEO data gets overwhelming, Barracuda tells you what actually matters. AI-powered decision support for technical SEO audits.'
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

<!-- AI Coming Soon — integrated naturally -->
<section class="bg-[#222] py-24 px-4 relative overflow-hidden">
	<div class="absolute top-0 right-0 w-[500px] h-[500px] bg-[#d79921]/5 rounded-full blur-[120px] pointer-events-none"></div>
	<div class="container mx-auto max-w-6xl relative z-10">
		<div class="grid grid-cols-1 lg:grid-cols-2 gap-16 items-center">
			<div>
				<div class="inline-flex items-center gap-2 bg-[#d79921]/10 border border-[#d79921]/20 text-[#d79921] px-4 py-1.5 rounded-full text-sm font-medium mb-6">
					<BrainCircuit class="w-3.5 h-3.5" />
					Coming Soon
				</div>
				<h2 class="text-3xl md:text-4xl lg:text-5xl font-heading font-bold text-white leading-tight mb-6">
					AI that knows your site,
					<br />
					<span class="text-[#d79921]">your voice, and your data.</span>
				</h2>
				<p class="text-lg text-white/50 leading-relaxed mb-8">
					Content briefs from GSC data. Articles in your brand voice. Internal link suggestions powered by semantic search. Intelligence at the moment you need it.
				</p>
				<a
					href="/roadmap"
					class="inline-flex items-center gap-2 border border-[#d79921]/30 hover:border-[#d79921]/60 text-white/80 hover:text-white px-6 py-3 rounded-xl font-medium transition-all"
				>
					See the full AI roadmap
					<ChevronRight class="w-4 h-4" />
				</a>
			</div>

			<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
				<div class="bg-[#1d1d1d] border border-[#d79921]/10 rounded-2xl p-5 hover:border-[#d79921]/25 transition-colors">
					<div class="w-9 h-9 rounded-xl bg-[#d79921]/10 flex items-center justify-center mb-3">
						<PenTool class="w-4 h-4 text-[#d79921]" />
					</div>
					<h3 class="text-sm font-heading font-bold text-white mb-1">Writing Voice</h3>
					<p class="text-xs text-white/40 leading-relaxed">Brand voice profiles from your existing content</p>
				</div>
				<div class="bg-[#1d1d1d] border border-[#d79921]/10 rounded-2xl p-5 hover:border-[#d79921]/25 transition-colors">
					<div class="w-9 h-9 rounded-xl bg-[#d79921]/10 flex items-center justify-center mb-3">
						<BarChart3 class="w-4 h-4 text-[#d79921]" />
					</div>
					<h3 class="text-sm font-heading font-bold text-white mb-1">GSC Intelligence</h3>
					<p class="text-xs text-white/40 leading-relaxed">Quick wins and declining pages from Search Console</p>
				</div>
				<div class="bg-[#1d1d1d] border border-[#d79921]/10 rounded-2xl p-5 hover:border-[#d79921]/25 transition-colors">
					<div class="w-9 h-9 rounded-xl bg-[#d79921]/10 flex items-center justify-center mb-3">
						<BrainCircuit class="w-4 h-4 text-[#d79921]" />
					</div>
					<h3 class="text-sm font-heading font-bold text-white mb-1">Content Briefs</h3>
					<p class="text-xs text-white/40 leading-relaxed">Data-backed briefs with outlines and keywords</p>
				</div>
				<div class="bg-[#1d1d1d] border border-[#d79921]/10 rounded-2xl p-5 hover:border-[#d79921]/25 transition-colors">
					<div class="w-9 h-9 rounded-xl bg-[#d79921]/10 flex items-center justify-center mb-3">
						<Link2 class="w-4 h-4 text-[#d79921]" />
					</div>
					<h3 class="text-sm font-heading font-bold text-white mb-1">Internal Links</h3>
					<p class="text-xs text-white/40 leading-relaxed">Semantic similarity for missing link opportunities</p>
				</div>
			</div>
		</div>
	</div>
</section>

<WhoItsFor />
<SocialProof />
<FinalCTA />
