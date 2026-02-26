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
<section class="relative overflow-hidden bg-[#222] px-4 py-24">
	<div
		class="pointer-events-none absolute top-0 right-0 h-[500px] w-[500px] rounded-full bg-[#d79921]/5 blur-[120px]"
	></div>
	<div class="relative z-10 container mx-auto max-w-6xl">
		<div class="grid grid-cols-1 items-center gap-16 lg:grid-cols-2">
			<div>
				<div
					class="mb-6 inline-flex items-center gap-2 rounded-full border border-[#d79921]/20 bg-[#d79921]/10 px-4 py-1.5 text-sm font-medium text-[#d79921]"
				>
					<BrainCircuit class="h-3.5 w-3.5" />
					Coming Soon
				</div>
				<h2
					class="mb-6 font-heading text-3xl leading-tight font-bold text-white md:text-4xl lg:text-5xl"
				>
					AI that knows your site,
					<br />
					<span class="text-[#d79921]">your voice, and your data.</span>
				</h2>
				<p class="mb-8 text-lg leading-relaxed text-white/50">
					Content briefs from GSC data. Articles in your brand voice. Internal link suggestions
					powered by semantic search. Intelligence at the moment you need it.
				</p>
				<a
					href="/roadmap"
					class="inline-flex items-center gap-2 rounded-xl border border-[#d79921]/30 px-6 py-3 font-medium text-white/80 transition-all hover:border-[#d79921]/60 hover:text-white"
				>
					See the full AI roadmap
					<ChevronRight class="h-4 w-4" />
				</a>
			</div>

			<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
				<div
					class="rounded-2xl border border-[#d79921]/10 bg-[#1d1d1d] p-5 transition-colors hover:border-[#d79921]/25"
				>
					<div class="mb-3 flex h-9 w-9 items-center justify-center rounded-xl bg-[#d79921]/10">
						<PenTool class="h-4 w-4 text-[#d79921]" />
					</div>
					<h3 class="mb-1 font-heading text-sm font-bold text-white">Writing Voice</h3>
					<p class="text-xs leading-relaxed text-white/40">
						Brand voice profiles from your existing content
					</p>
				</div>
				<div
					class="rounded-2xl border border-[#d79921]/10 bg-[#1d1d1d] p-5 transition-colors hover:border-[#d79921]/25"
				>
					<div class="mb-3 flex h-9 w-9 items-center justify-center rounded-xl bg-[#d79921]/10">
						<BarChart3 class="h-4 w-4 text-[#d79921]" />
					</div>
					<h3 class="mb-1 font-heading text-sm font-bold text-white">GSC Intelligence</h3>
					<p class="text-xs leading-relaxed text-white/40">
						Quick wins and declining pages from Search Console
					</p>
				</div>
				<div
					class="rounded-2xl border border-[#d79921]/10 bg-[#1d1d1d] p-5 transition-colors hover:border-[#d79921]/25"
				>
					<div class="mb-3 flex h-9 w-9 items-center justify-center rounded-xl bg-[#d79921]/10">
						<BrainCircuit class="h-4 w-4 text-[#d79921]" />
					</div>
					<h3 class="mb-1 font-heading text-sm font-bold text-white">Content Briefs</h3>
					<p class="text-xs leading-relaxed text-white/40">
						Data-backed briefs with outlines and keywords
					</p>
				</div>
				<div
					class="rounded-2xl border border-[#d79921]/10 bg-[#1d1d1d] p-5 transition-colors hover:border-[#d79921]/25"
				>
					<div class="mb-3 flex h-9 w-9 items-center justify-center rounded-xl bg-[#d79921]/10">
						<Link2 class="h-4 w-4 text-[#d79921]" />
					</div>
					<h3 class="mb-1 font-heading text-sm font-bold text-white">Internal Links</h3>
					<p class="text-xs leading-relaxed text-white/40">
						Semantic similarity for missing link opportunities
					</p>
				</div>
			</div>
		</div>
	</div>
</section>

<WhoItsFor />
<SocialProof />
<FinalCTA />
