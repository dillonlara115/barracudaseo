<script lang="ts">
	import MetaTags from '../../components/MetaTags.svelte';
	import { getMetaTags, getBreadcrumbSchema } from '$lib/meta';
	import { trackSignup } from '$lib/analytics';

	const meta = getMetaTags({
		title: 'About Barracuda SEO — Why we built this',
		description: 'SEO tools stopped helping people make decisions. Barracuda exists to change that. Learn why we built a decision-support layer for technical SEO.'
	});

	const structuredData = getBreadcrumbSchema([
		{ name: 'Home', url: '/' },
		{ name: 'About', url: '/about' }
	]);

	function handleSignupClick() {
		trackSignup({
			source: 'About page CTA',
			location: 'about-footer'
		});
	}
</script>

<MetaTags config={{ ...meta, structuredData }} />

<!-- Hero Section -->
<section class="py-20 bg-[#3c3836]">
	<div class="container mx-auto px-4 max-w-4xl text-center">
		<h1 class="text-5xl md:text-6xl font-heading font-bold mb-8 text-white">About BarracudaSEO</h1>
		
		<div class="text-xl md:text-2xl text-white/80 leading-relaxed max-w-3xl mx-auto">
			<p class="mb-8">BarracudaSEO exists for one simple reason:</p>
			<p class="font-bold text-white text-3xl md:text-4xl mb-8">SEO tools stopped helping people make decisions.</p>
			<p>
				They got very good at collecting data.
				<br />
				They got very bad at helping people decide what to do with it.
			</p>
		</div>
	</div>
</section>

<!-- The Problem Section -->
<section class="py-20 bg-[#282828]">
	<div class="container mx-auto px-4 max-w-4xl">
		<h2 class="text-3xl font-heading font-bold mb-12 text-white text-center">The problem we kept running into</h2>
		
		<div class="grid grid-cols-1 md:grid-cols-2 gap-12 items-center">
			<div class="text-lg text-white/80 space-y-6">
				<p>We’ve all been there.</p>
				<p>
					You run a site crawl and suddenly you’re staring at hundreds of issues:
					missing metadata, duplicate titles, redirects, slow pages, warnings, scores.
				</p>
				<p class="text-[#d79921] font-medium italic text-xl">
					That’s where confusion, stress, and audit paralysis creep in.
				</p>
			</div>
			
			<div class="bg-[#3c3836] p-8 rounded-xl border border-white/5">
				<p class="text-white font-medium mb-4">And then comes the hard part — the part no tool really helps with:</p>
				<ul class="space-y-4">
					<li class="flex items-start">
						<span class="text-[#cc241d] mr-3 font-bold">?</span>
						<span class="text-white/90">What actually matters right now?</span>
					</li>
					<li class="flex items-start">
						<span class="text-[#cc241d] mr-3 font-bold">?</span>
						<span class="text-white/90">What can safely wait?</span>
					</li>
					<li class="flex items-start">
						<span class="text-[#cc241d] mr-3 font-bold">?</span>
						<span class="text-white/90">What do I fix first?</span>
					</li>
					<li class="flex items-start">
						<span class="text-[#cc241d] mr-3 font-bold">?</span>
						<span class="text-white/90">How do I explain this to a client without second-guessing myself?</span>
					</li>
				</ul>
				<div class="mt-6 pt-6 border-t border-white/10 text-center text-white/60 text-sm">
					Most SEO tools stop at detection. The interpretation is left to you.
				</div>
			</div>
		</div>
	</div>
</section>

<!-- The Solution / Difference -->
<section class="py-20 bg-[#3c3836]">
	<div class="container mx-auto px-4 max-w-4xl text-center">
		<h2 class="text-3xl font-heading font-bold mb-8 text-[#8ec07c]">Why BarracudaSEO is different</h2>
		
		<div class="text-lg md:text-xl text-white/80 space-y-8 max-w-3xl mx-auto">
			<p>
				BarracudaSEO wasn’t built to collect more SEO data.
				<br />
				<span class="text-white font-bold text-2xl mt-2 block">It was built to help you decide.</span>
			</p>
			
			<div class="bg-[#282828] p-8 rounded-xl my-12 border border-[#8ec07c]/20">
				<p class="mb-4">
					Barracuda acts as a <strong>decision-support layer</strong> between raw crawl data and action.
				</p>
				<p>
					It helps you prioritize issues, understand impact, and explain <em>why</em> certain fixes matter more than others.
				</p>
			</div>

			<p>
				You still apply judgment.
				<br />
				You still make the final call.
				<br />
				<span class="text-white font-medium block mt-4">Barracuda just removes the chaos in between.</span>
			</p>
		</div>
	</div>
</section>

<!-- Opinionated Section -->
<section class="py-20 bg-[#282828]">
	<div class="container mx-auto px-4 max-w-5xl">
		<div class="flex flex-col md:flex-row gap-12 items-center">
			<div class="w-full md:w-1/2">
				<h2 class="text-3xl font-heading font-bold mb-6 text-white">Opinionated by design</h2>
				<div class="text-lg text-white/80 space-y-6">
					<p>We believe clarity requires focus.</p>
					<p>That’s why Barracuda is intentionally opinionated.</p>
					<p>
						We don’t surface every possible issue just because it exists.
						We highlight what matters <strong>now</strong>, explain why it matters, and make it defensible.
					</p>
				</div>
			</div>
			
			<div class="w-full md:w-1/2 bg-[#3c3836] p-8 rounded-xl border-l-4 border-[#8ec07c]">
				<h3 class="text-xl font-bold text-white mb-4">And when something doesn’t matter?</h3>
				<p class="text-white/90 text-lg mb-6">We say that too.</p>
				<p class="text-white/70 italic">
					That’s not hiding data — it’s respecting your time, your responsibility, and the reality of how SEO work actually gets done.
				</p>
			</div>
		</div>
	</div>
</section>

<!-- Built For Pressure -->
<section class="py-20 bg-[#3c3836]">
	<div class="container mx-auto px-4 max-w-6xl">
		<h2 class="text-3xl font-heading font-bold mb-12 text-white text-center">Built for real-world SEO pressure</h2>
		
		<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
			<div class="bg-[#282828] p-6 rounded-lg border border-white/5 hover:border-[#8ec07c]/30 transition-colors">
				<h3 class="text-xl font-bold text-white mb-3">Client calls</h3>
				<p class="text-white/70">Where you need clear answers immediately, not after digging through spreadsheets.</p>
			</div>
			
			<div class="bg-[#282828] p-6 rounded-lg border border-white/5 hover:border-[#8ec07c]/30 transition-colors">
				<h3 class="text-xl font-bold text-white mb-3">Time-limited audits</h3>
				<p class="text-white/70">When accountability is high and you can't afford to waste time on low-impact fixes.</p>
			</div>
			
			<div class="bg-[#282828] p-6 rounded-lg border border-white/5 hover:border-[#8ec07c]/30 transition-colors">
				<h3 class="text-xl font-bold text-white mb-3">High-stakes decisions</h3>
				<p class="text-white/70">Situations where confidence in your recommendation is as important as technical correctness.</p>
			</div>
		</div>
		
		<p class="text-center text-white/60 mt-12 text-lg">
			We design for clarity under pressure, not endless dashboards.
		</p>
	</div>
</section>

<!-- Not Replacement + Audience -->
<section class="py-20 bg-[#282828]">
	<div class="container mx-auto px-4 max-w-5xl">
		<div class="grid grid-cols-1 md:grid-cols-2 gap-16">
			<!-- Not a replacement -->
			<div>
				<h2 class="text-2xl font-heading font-bold mb-6 text-white">Not a replacement for expertise</h2>
				<p class="text-white/80 mb-6">
					BarracudaSEO isn’t here to replace your experience or instincts. Think of it as:
				</p>
				<ul class="space-y-3">
					<li class="flex items-center text-white/80">
						<span class="w-2 h-2 bg-[#8ec07c] rounded-full mr-3"></span>
						A second set of eyes
					</li>
					<li class="flex items-center text-white/80">
						<span class="w-2 h-2 bg-[#8ec07c] rounded-full mr-3"></span>
						A prioritization partner
					</li>
					<li class="flex items-center text-white/80">
						<span class="w-2 h-2 bg-[#8ec07c] rounded-full mr-3"></span>
						A way to sanity-check decisions
					</li>
					<li class="flex items-center text-white/80">
						<span class="w-2 h-2 bg-[#8ec07c] rounded-full mr-3"></span>
						A way to explain your thinking clearly
					</li>
				</ul>
				<p class="mt-6 text-white font-medium">
					You decide what to act on. Barracuda helps you stand behind those decisions.
				</p>
			</div>

			<!-- Who it is for -->
			<div>
				<h2 class="text-2xl font-heading font-bold mb-6 text-white">Who BarracudaSEO is for</h2>
				<ul class="space-y-4 mb-8">
					<li class="flex items-start">
						<svg class="w-6 h-6 text-[#8ec07c] mr-3 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
						</svg>
						<span class="text-white/80">Consultants who need to explain SEO decisions clearly</span>
					</li>
					<li class="flex items-start">
						<svg class="w-6 h-6 text-[#8ec07c] mr-3 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
						</svg>
						<span class="text-white/80">Agencies balancing speed, accuracy, and trust</span>
					</li>
					<li class="flex items-start">
						<svg class="w-6 h-6 text-[#8ec07c] mr-3 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
						</svg>
						<span class="text-white/80">Technical SEOs tired of audit overwhelm</span>
					</li>
					<li class="flex items-start">
						<svg class="w-6 h-6 text-[#8ec07c] mr-3 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
						</svg>
						<span class="text-white/80">Teams who value clarity over clutter</span>
					</li>
				</ul>
				<p class="text-white/60 italic text-sm border-t border-white/10 pt-4">
					If you love digging through raw CSVs and chasing every possible metric, Barracuda may not be for you — and that’s okay.
				</p>
			</div>
		</div>
	</div>
</section>

<!-- Closing Section -->
<section class="py-24 bg-gradient-to-t from-[#3c3836] to-[#282828]">
	<div class="container mx-auto px-4 max-w-4xl text-center">
		<h2 class="text-3xl font-heading font-bold mb-8 text-white">Why we’re building this</h2>
		
		<div class="text-xl text-white/80 space-y-6 max-w-3xl mx-auto mb-16">
			<p>SEO shouldn’t feel like guesswork.</p>
			<p>It shouldn’t feel like anxiety.</p>
			<p>And it shouldn’t require heroic mental effort just to decide what to do next.</p>
			
			<p class="text-white font-bold text-2xl py-4">
				BarracudaSEO exists to make the decision part of SEO easier.
			</p>
			
			<p>
				So when someone asks, “What’s wrong with the site?”
				<br />
				You don’t panic.
				<br />
				You answer — with clarity, confidence, and reasons you can explain.
			</p>
		</div>

		<div class="bg-[#282828] border border-[#8ec07c]/30 rounded-2xl p-8 md:p-12 inline-block shadow-xl">
			<h3 class="text-2xl font-bold text-white mb-4">Ready to stop second-guessing your audits?</h3>
			<p class="text-white/70 mb-8">Start using BarracudaSEO and turn crawl data into decisions you can stand behind.</p>
			<a
				href="https://app.barracudaseo.com"
				class="bg-[#8ec07c] hover:bg-[#a0d28c] text-[#3c3836] px-8 py-4 rounded-lg font-bold text-lg transition-colors inline-block"
				target="_blank"
				rel="noopener noreferrer"
				onclick={handleSignupClick}
			>
				Start your first audit
			</a>
		</div>
	</div>
</section>
