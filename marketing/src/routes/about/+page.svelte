<script lang="ts">
	import MetaTags from '../../components/MetaTags.svelte';
	import RecentPosts from '../../components/sections/RecentPosts.svelte';
	import { getMetaTags, getBreadcrumbSchema } from '$lib/meta';
	import { trackSignup } from '$lib/analytics';

	const meta = getMetaTags({
		title: 'About Barracuda SEO — Why we built this',
		description:
			'SEO tools stopped helping people make decisions. Barracuda exists to change that. Learn why we built a decision-support layer for technical SEO.'
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
<section class="bg-[#3c3836] py-20">
	<div class="container mx-auto max-w-4xl px-4 text-center">
		<h1 class="mb-8 font-heading text-5xl font-bold text-white md:text-6xl">About BarracudaSEO</h1>

		<div class="mx-auto max-w-3xl text-xl leading-relaxed text-white/80 md:text-2xl">
			<p class="mb-8">BarracudaSEO exists for one simple reason:</p>
			<p class="mb-8 text-3xl font-bold text-white md:text-4xl">
				SEO tools stopped helping people make decisions.
			</p>
			<p>
				They got very good at collecting data.
				<br />
				They got very bad at helping people decide what to do with it.
			</p>
		</div>
	</div>
</section>

<!-- The Problem Section -->
<section class="bg-[#282828] py-20">
	<div class="container mx-auto max-w-4xl px-4">
		<h2 class="mb-12 text-center font-heading text-3xl font-bold text-white">
			The problem we kept running into
		</h2>

		<div class="grid grid-cols-1 items-center gap-12 md:grid-cols-2">
			<div class="space-y-6 text-lg text-white/80">
				<p>We’ve all been there.</p>
				<p>
					You run a <a
						href="/blog/complete-technical-seo-audit-guide"
						class="text-[#8ec07c] underline hover:text-[#a0d28c]">site crawl</a
					> and suddenly you’re staring at hundreds of issues: missing metadata, duplicate titles, redirects,
					slow pages, warnings, scores.
				</p>
				<p class="text-xl font-medium text-[#d79921] italic">
					That’s where confusion, stress, and audit paralysis creep in.
				</p>
			</div>

			<div class="rounded-xl border border-white/5 bg-[#3c3836] p-8">
				<p class="mb-4 font-medium text-white">
					And then comes the hard part — the part no tool really helps with:
				</p>
				<ul class="space-y-4">
					<li class="flex items-start">
						<span class="mr-3 font-bold text-[#cc241d]">?</span>
						<span class="text-white/90">What actually matters right now?</span>
					</li>
					<li class="flex items-start">
						<span class="mr-3 font-bold text-[#cc241d]">?</span>
						<span class="text-white/90">What can safely wait?</span>
					</li>
					<li class="flex items-start">
						<span class="mr-3 font-bold text-[#cc241d]">?</span>
						<span class="text-white/90">What do I fix first?</span>
					</li>
					<li class="flex items-start">
						<span class="mr-3 font-bold text-[#cc241d]">?</span>
						<span class="text-white/90"
							>How do I explain this to a client without second-guessing myself?</span
						>
					</li>
				</ul>
				<div class="mt-6 border-t border-white/10 pt-6 text-center text-sm text-white/60">
					Most SEO tools stop at detection. The interpretation is left to you.
				</div>
			</div>
		</div>
	</div>
</section>

<!-- The Solution / Difference -->
<section class="bg-[#3c3836] py-20">
	<div class="container mx-auto max-w-4xl px-4 text-center">
		<h2 class="mb-8 font-heading text-3xl font-bold text-[#8ec07c]">
			Why BarracudaSEO is different
		</h2>

		<div class="mx-auto max-w-3xl space-y-8 text-lg text-white/80 md:text-xl">
			<p>
				BarracudaSEO wasn’t built to collect more SEO data.
				<br />
				<span class="mt-2 block text-2xl font-bold text-white"
					>It was built to help you decide.</span
				>
			</p>

			<div class="my-12 rounded-xl border border-[#8ec07c]/20 bg-[#282828] p-8">
				<p class="mb-4">
					Barracuda acts as a <strong>decision-support layer</strong> between raw crawl data and action.
				</p>
				<p>
					It helps you <a
						href="/blog/how-to-prioritize-seo-issues"
						class="text-[#8ec07c] underline hover:text-[#a0d28c]">prioritize issues</a
					>, understand impact, and explain <em>why</em> certain fixes matter more than others.
				</p>
			</div>

			<p>
				You still apply judgment.
				<br />
				You still make the final call.
				<br />
				<span class="mt-4 block font-medium text-white"
					>Barracuda just removes the chaos in between. <a
						href="/features"
						class="text-[#8ec07c] underline hover:text-[#a0d28c]"
						>See how it works</a
					>.</span
				>
			</p>
		</div>
	</div>
</section>

<!-- Opinionated Section -->
<section class="bg-[#282828] py-20">
	<div class="container mx-auto max-w-5xl px-4">
		<div class="flex flex-col items-center gap-12 md:flex-row">
			<div class="w-full md:w-1/2">
				<h2 class="mb-6 font-heading text-3xl font-bold text-white">Opinionated by design</h2>
				<div class="space-y-6 text-lg text-white/80">
					<p>We believe clarity requires focus.</p>
					<p>That’s why Barracuda is intentionally opinionated.</p>
					<p>
						We don’t surface every possible issue just because it exists. We highlight what matters <strong
							>now</strong
						>, explain why it matters, and make it defensible.
					</p>
				</div>
			</div>

			<div class="w-full rounded-xl border-l-4 border-[#8ec07c] bg-[#3c3836] p-8 md:w-1/2">
				<h3 class="mb-4 text-xl font-bold text-white">And when something doesn’t matter?</h3>
				<p class="mb-6 text-lg text-white/90">We say that too.</p>
				<p class="text-white/70 italic">
					That’s not hiding data — it’s respecting your time, your responsibility, and the reality
					of how SEO work actually gets done.
				</p>
			</div>
		</div>
	</div>
</section>

<!-- Built For Pressure -->
<section class="bg-[#3c3836] py-20">
	<div class="container mx-auto max-w-6xl px-4">
		<h2 class="mb-12 text-center font-heading text-3xl font-bold text-white">
			Built for real-world SEO pressure
		</h2>

		<div class="grid grid-cols-1 gap-6 md:grid-cols-3">
			<div
				class="rounded-lg border border-white/5 bg-[#282828] p-6 transition-colors hover:border-[#8ec07c]/30"
			>
				<h3 class="mb-3 text-xl font-bold text-white">Client calls</h3>
				<p class="text-white/70">
					Where you need clear answers immediately — like explaining why a page is <a
						href="/blog/crawled-not-indexed"
						class="text-[#8ec07c] underline hover:text-[#a0d28c]"
						>crawled but not indexed</a
					> — not after digging through spreadsheets.
				</p>
			</div>

			<div
				class="rounded-lg border border-white/5 bg-[#282828] p-6 transition-colors hover:border-[#8ec07c]/30"
			>
				<h3 class="mb-3 text-xl font-bold text-white">Time-limited audits</h3>
				<p class="text-white/70">
					When accountability is high and you can't afford to waste time chasing <a
						href="/blog/duplicate-meta-tags-fix"
						class="text-[#8ec07c] underline hover:text-[#a0d28c]"
						>duplicate meta tags</a
					> on low-traffic pages.
				</p>
			</div>

			<div
				class="rounded-lg border border-white/5 bg-[#282828] p-6 transition-colors hover:border-[#8ec07c]/30"
			>
				<h3 class="mb-3 text-xl font-bold text-white">High-stakes decisions</h3>
				<p class="text-white/70">
					Situations where confidence in your recommendation is as important as technical
					correctness.
				</p>
			</div>
		</div>

		<p class="mt-12 text-center text-lg text-white/60">
			We design for clarity under pressure, not endless dashboards.
		</p>
	</div>
</section>

<!-- Not Replacement + Audience -->
<section class="bg-[#282828] py-20">
	<div class="container mx-auto max-w-5xl px-4">
		<div class="grid grid-cols-1 gap-16 md:grid-cols-2">
			<!-- Not a replacement -->
			<div>
				<h2 class="mb-6 font-heading text-2xl font-bold text-white">
					Not a replacement for expertise
				</h2>
				<p class="mb-6 text-white/80">
					BarracudaSEO isn’t here to replace your experience or instincts. Think of it as:
				</p>
				<ul class="space-y-3">
					<li class="flex items-center text-white/80">
						<span class="mr-3 h-2 w-2 rounded-full bg-[#8ec07c]"></span>
						A second set of eyes
					</li>
					<li class="flex items-center text-white/80">
						<span class="mr-3 h-2 w-2 rounded-full bg-[#8ec07c]"></span>
						A prioritization partner
					</li>
					<li class="flex items-center text-white/80">
						<span class="mr-3 h-2 w-2 rounded-full bg-[#8ec07c]"></span>
						A way to sanity-check decisions
					</li>
					<li class="flex items-center text-white/80">
						<span class="mr-3 h-2 w-2 rounded-full bg-[#8ec07c]"></span>
						A way to explain your thinking clearly
					</li>
				</ul>
				<p class="mt-6 font-medium text-white">
					You decide what to act on. Barracuda helps you stand behind those decisions.
				</p>
			</div>

			<!-- Who it is for -->
			<div>
				<h2 class="mb-6 font-heading text-2xl font-bold text-white">Who BarracudaSEO is for</h2>
				<ul class="mb-8 space-y-4">
					<li class="flex items-start">
						<svg
							class="mr-3 h-6 w-6 shrink-0 text-[#8ec07c]"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M5 13l4 4L19 7"
							></path>
						</svg>
						<span class="text-white/80">Consultants who need to explain SEO decisions clearly</span>
					</li>
					<li class="flex items-start">
						<svg
							class="mr-3 h-6 w-6 shrink-0 text-[#8ec07c]"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M5 13l4 4L19 7"
							></path>
						</svg>
						<span class="text-white/80">Agencies balancing speed, accuracy, and trust</span>
					</li>
					<li class="flex items-start">
						<svg
							class="mr-3 h-6 w-6 shrink-0 text-[#8ec07c]"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M5 13l4 4L19 7"
							></path>
						</svg>
						<span class="text-white/80"
							>Technical SEOs tired of <a
								href="/blog/why-seo-audits-feel-overwhelming"
								class="text-[#8ec07c] underline hover:text-[#a0d28c]">audit overwhelm</a
							></span
						>
					</li>
					<li class="flex items-start">
						<svg
							class="mr-3 h-6 w-6 shrink-0 text-[#8ec07c]"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M5 13l4 4L19 7"
							></path>
						</svg>
						<span class="text-white/80">Teams who value clarity over clutter</span>
					</li>
				</ul>
				<p class="border-t border-white/10 pt-4 text-sm text-white/60 italic">
					If you love digging through raw CSVs and chasing every possible metric, Barracuda may not
					be for you — and that’s okay.
				</p>
			</div>
		</div>
	</div>
</section>

<RecentPosts
	title="Learn More"
	subtitle="Guides on the technical SEO problems Barracuda helps you solve."
/>

<!-- Closing Section -->
<section class="bg-gradient-to-t from-[#3c3836] to-[#282828] py-24">
	<div class="container mx-auto max-w-4xl px-4 text-center">
		<h2 class="mb-8 font-heading text-3xl font-bold text-white">Why we’re building this</h2>

		<div class="mx-auto mb-16 max-w-3xl space-y-6 text-xl text-white/80">
			<p>SEO shouldn’t feel like guesswork.</p>
			<p>It shouldn’t feel like anxiety.</p>
			<p>And it shouldn’t require heroic mental effort just to decide what to do next.</p>

			<p class="py-4 text-2xl font-bold text-white">
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

		<div
			class="inline-block rounded-2xl border border-[#8ec07c]/30 bg-[#282828] p-8 shadow-xl md:p-12"
		>
			<h3 class="mb-4 text-2xl font-bold text-white">Ready to stop second-guessing your audits?</h3>
			<p class="mb-8 text-white/70">
				Start using BarracudaSEO and turn crawl data into decisions you can stand behind.
			</p>
			<a
				href="https://app.barracudaseo.com"
				class="inline-block rounded-lg bg-[#8ec07c] px-8 py-4 text-lg font-bold text-[#3c3836] transition-colors hover:bg-[#a0d28c]"
				target="_blank"
				rel="noopener noreferrer"
				onclick={handleSignupClick}
			>
				Start your first audit
			</a>
		</div>
	</div>
</section>
