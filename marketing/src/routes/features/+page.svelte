<script lang="ts">
	import {
		Bolt,
		ChartBar,
		Users,
		Plug,
		Code,
		CircleCheck,
		ChevronRight,
		Target,
		BrainCircuit,
		PenTool,
		FileText,
		Sparkles,
		Link2
	} from '@lucide/svelte';
	import MetaTags from '../../components/MetaTags.svelte';
	import RecentPosts from '../../components/sections/RecentPosts.svelte';
	import { getMetaTags, getBreadcrumbSchema } from '$lib/meta';

	const meta = getMetaTags({
		title: 'Decision-Support for SEO Audits',
		description:
			'Barracuda SEO features: High-speed crawling, interactive link graphs, issue prioritization, and team collaboration tools for agencies and in-house teams.'
	});

	const structuredData = getBreadcrumbSchema([
		{ name: 'Home', url: '/' },
		{ name: 'Features', url: '/features' }
	]);
</script>

<MetaTags config={{ ...meta, structuredData }} />

<!-- Hero Section -->
<section class="bg-gradient-to-b from-[#3c3836] to-[#2d2826] py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-4xl text-center">
			<h1 class="mb-6 font-heading text-5xl font-bold text-white md:text-6xl">
				Crawl Smarter. Fix Faster. Grow Visibility.
			</h1>
			<p class="mx-auto mb-10 max-w-2xl text-xl text-white/80 md:text-2xl">
				Discover and fix the SEO issues holding your website back. Barracuda delivers fast, reliable
				crawls with actionable insights—all in one sleek web dashboard.
			</p>
			<div class="flex flex-col justify-center gap-4 sm:flex-row">
				<a
					href="https://app.barracudaseo.com"
					class="rounded-lg bg-[#8ec07c] px-8 py-4 text-lg font-medium text-[#3c3836] transition-colors hover:bg-[#a0d28c]"
					target="_blank"
					rel="noopener noreferrer"
				>
					Start Your Free 100-Page Audit
				</a>
				<a
					href="#dashboard"
					class="rounded-lg border-2 border-white/20 px-8 py-4 text-lg font-medium text-white transition-colors hover:border-[#8ec07c]"
				>
					See Barracuda in Action
				</a>
			</div>
		</div>
	</div>
</section>

<!-- Section 1: Intelligent Crawling -->
<section id="crawling" class="bg-[#3c3836] py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-4xl">
			<div class="mb-6 flex items-center gap-4">
				<div class="rounded-lg bg-[#8ec07c]/10 p-3">
					<Bolt class="h-8 w-8 text-[#8ec07c]" />
				</div>
				<h2 class="font-heading text-4xl font-bold text-white md:text-5xl">
					Lightning-Fast Crawling for Every Site
				</h2>
			</div>
			<p class="mb-8 max-w-3xl text-xl text-white/80">
				Scan your site in minutes, not hours. Barracuda's optimized crawl engine uncovers <a
					href="/blog/find-fix-broken-links"
					class="text-[#8ec07c] underline hover:text-[#a0d28c]">broken links</a
				>, missing tags,
				<a
					href="/blog/redirect-chains-seo-killer"
					class="text-[#8ec07c] underline hover:text-[#a0d28c]">redirect chains</a
				>, and slow-loading pages with unmatched precision. Whether you're running a local crawl or
				managing multiple client sites, Barracuda adapts to your workflow.
			</p>
			<div class="mb-8 grid grid-cols-1 gap-6 md:grid-cols-2">
				<div class="rounded-lg border border-white/10 bg-[#2d2826] p-6">
					<h3 class="mb-3 font-heading text-lg font-bold text-[#8ec07c]">Scale & Performance</h3>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span
								>Crawl up to <strong class="text-white">10,000+ pages</strong> with the Pro plan</span
							>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span><strong class="text-white">Respect robots.txt</strong> automatically</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span
								>Supports <strong class="text-white">sitemap seeding</strong> and custom depth limits</span
							>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span
								>Export results to <strong class="text-white">CSV, JSON, or Supabase Cloud</strong
								></span
							>
						</li>
					</ul>
				</div>
				<div class="rounded-lg border border-white/10 bg-[#2d2826] p-6">
					<h3 class="mb-3 font-heading text-lg font-bold text-[#8ec07c]">What We Detect</h3>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Broken links and redirect chains</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Missing or <a
								href="/blog/duplicate-meta-tags-fix"
								class="text-[#8ec07c] underline hover:text-[#a0d28c]"
								>duplicate meta tags</a
							></span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Slow-loading pages</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span><a
								href="/blog/crawled-not-indexed"
								class="text-[#8ec07c] underline hover:text-[#a0d28c]"
								>Indexability issues</a
							> and best practice violations</span>
						</li>
					</ul>
				</div>
			</div>
			<div class="text-center">
				<a
					href="https://app.barracudaseo.com"
					class="inline-flex items-center gap-2 rounded-lg bg-[#8ec07c] px-6 py-3 font-medium text-[#3c3836] transition-colors hover:bg-[#a0d28c]"
					target="_blank"
					rel="noopener noreferrer"
				>
					Start Crawling Free
					<ChevronRight class="h-4 w-4" />
				</a>
			</div>
		</div>
	</div>
</section>

<!-- Section 2: Insightful Dashboard -->
<section id="dashboard" class="bg-[#2d2826] py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-4xl">
			<div class="mb-6 flex items-center gap-4">
				<div class="rounded-lg bg-[#8ec07c]/10 p-3">
					<ChartBar class="h-8 w-8 text-[#8ec07c]" />
				</div>
				<h2 class="font-heading text-4xl font-bold text-white md:text-5xl">
					Turn Data into Decisions
				</h2>
			</div>
			<p class="mb-8 max-w-3xl text-xl text-white/80">
				No more spreadsheets or endless tabs. The Barracuda dashboard organizes crawl results into a
				clean, visual interface that surfaces the most critical issues first.
			</p>
			<div class="mb-8 grid grid-cols-1 gap-6 md:grid-cols-2">
				<div class="rounded-lg border border-white/10 bg-[#3c3836] p-6">
					<h3 class="mb-3 font-heading text-lg font-bold text-[#8ec07c]">Priority Scoring</h3>
					<p class="text-white/70">
						Automatically <a
							href="/blog/how-to-prioritize-seo-issues"
							class="text-[#8ec07c] underline hover:text-[#a0d28c]"
							>ranks issues by impact</a
						>, so you know exactly what to fix first.
					</p>
				</div>
				<div class="rounded-lg border border-white/10 bg-[#3c3836] p-6">
					<h3 class="mb-3 font-heading text-lg font-bold text-[#8ec07c]">Interactive Link Graph</h3>
					<p class="text-white/70">
						<a
							href="/blog/visualize-site-structure-link-graph"
							class="text-[#8ec07c] underline hover:text-[#a0d28c]"
							>Visualize site structure</a
						> and discover orphaned pages at a glance.
					</p>
				</div>
				<div class="rounded-lg border border-white/10 bg-[#3c3836] p-6">
					<h3 class="mb-3 font-heading text-lg font-bold text-[#8ec07c]">Issue Modals</h3>
					<p class="text-white/70">
						Get details, examples, and best-practice guidance for every problem.
					</p>
				</div>
				<div class="relative rounded-lg border border-[#8ec07c]/30 bg-[#3c3836] p-6">
					<div
						class="absolute -top-2 -right-2 rounded bg-[#8ec07c] px-2 py-1 text-xs font-bold text-[#3c3836]"
					>
						Pro
					</div>
					<h3 class="mb-3 font-heading text-lg font-bold text-[#8ec07c]">
						AI-Powered Recommendations
					</h3>
					<p class="text-white/70">See how to fix issues, not just what's broken. (Pro feature)</p>
				</div>
			</div>
			<div class="text-center">
				<a
					href="https://app.barracudaseo.com"
					class="inline-flex items-center gap-2 rounded-lg bg-[#8ec07c] px-6 py-3 font-medium text-[#3c3836] transition-colors hover:bg-[#a0d28c]"
					target="_blank"
					rel="noopener noreferrer"
				>
					Try the Dashboard — Free for 100 Pages
					<ChevronRight class="h-4 w-4" />
				</a>
			</div>
		</div>
	</div>
</section>

<!-- Section 3: Collaboration Made Easy -->
<section id="collaboration" class="bg-[#3c3836] py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-4xl">
			<div class="mb-6 flex items-center gap-4">
				<div class="rounded-lg bg-[#8ec07c]/10 p-3">
					<Users class="h-8 w-8 text-[#8ec07c]" />
				</div>
				<h2 class="font-heading text-4xl font-bold text-white md:text-5xl">
					Built for Teams and Agencies
				</h2>
			</div>
			<p class="mb-8 max-w-3xl text-xl text-white/80">
				Work together seamlessly. Invite your teammates, assign roles, and share crawl results in
				real-time. Whether you're an agency managing multiple client sites or an in-house SEO team,
				Barracuda makes collaboration effortless.
			</p>
			<div class="mb-8 rounded-lg border border-white/10 bg-[#2d2826] p-8">
				<div class="grid grid-cols-1 gap-6 md:grid-cols-3">
					<div>
						<h3 class="mb-2 font-heading text-lg font-bold text-[#8ec07c]">Owner</h3>
						<p class="text-sm text-white/70">
							Full permissions including billing, invitations, and project deletion.
						</p>
					</div>
					<div>
						<h3 class="mb-2 font-heading text-lg font-bold text-[#8ec07c]">Editor</h3>
						<p class="text-sm text-white/70">
							Can trigger crawls, manage issues, and view recommendations.
						</p>
					</div>
					<div>
						<h3 class="mb-2 font-heading text-lg font-bold text-[#8ec07c]">Viewer</h3>
						<p class="text-sm text-white/70">Read-only access to dashboard and reports.</p>
					</div>
				</div>
				<div class="mt-6 border-t border-white/10 pt-6">
					<p class="text-center text-white/70">
						Team add-ons available starting at <strong class="text-white">$5 per user</strong>
					</p>
				</div>
			</div>
			<div class="text-center">
				<a
					href="/pricing"
					class="inline-flex items-center gap-2 rounded-lg bg-[#8ec07c] px-6 py-3 font-medium text-[#3c3836] transition-colors hover:bg-[#a0d28c]"
				>
					Upgrade to a Team Account
					<ChevronRight class="h-4 w-4" />
				</a>
			</div>
		</div>
	</div>
</section>

<!-- Section 4: Advanced Integrations -->
<section id="integrations" class="bg-[#2d2826] py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-4xl">
			<div class="mb-6 flex items-center gap-4">
				<div class="rounded-lg bg-[#8ec07c]/10 p-3">
					<Plug class="h-8 w-8 text-[#8ec07c]" />
				</div>
				<h2 class="font-heading text-4xl font-bold text-white md:text-5xl">
					Connect Your Favorite SEO Tools
				</h2>
			</div>
			<p class="mb-8 max-w-3xl text-xl text-white/80">
				Take your audits further with integrations that give you context beyond the crawl. Pro users
				can link <a
					href="/blog/find-declining-pages-google-search-console"
					class="text-[#8ec07c] underline hover:text-[#a0d28c]">Google Search Console</a
				>, Google Analytics, and Microsoft Clarity to enrich recommendations and track performance
				directly from the dashboard.
			</p>
			<div class="mb-8 grid grid-cols-1 gap-6 md:grid-cols-2">
				<div class="rounded-lg border border-white/10 bg-[#3c3836] p-6">
					<h3 class="mb-2 font-heading text-lg font-bold text-[#8ec07c]">Google Search Console</h3>
					<p class="text-sm text-white/70">
						Import impressions, clicks, and CTR data to prioritize fixes that impact visibility.
					</p>
				</div>
				<div class="rounded-lg border border-white/10 bg-[#3c3836] p-6">
					<h3 class="mb-2 font-heading text-lg font-bold text-[#8ec07c]">Google Analytics</h3>
					<p class="text-sm text-white/70">
						Measure engagement and conversions by page to focus on high-value fixes.
					</p>
				</div>
				<div class="rounded-lg border border-white/10 bg-[#3c3836] p-6">
					<h3 class="mb-2 font-heading text-lg font-bold text-[#8ec07c]">Microsoft Clarity</h3>
					<p class="text-sm text-white/70">
						Visualize user behavior heatmaps alongside technical issues.
					</p>
				</div>
				<div class="rounded-lg border border-white/10 bg-[#3c3836] p-6">
					<h3 class="mb-2 font-heading text-lg font-bold text-[#8ec07c]">Slack</h3>
					<p class="text-sm text-white/70">
						Get instant crawl reports in your workspace. <span class="text-white/50"
							>(Coming Soon)</span
						>
					</p>
				</div>
			</div>
			<div class="text-center">
				<a
					href="/pricing"
					class="inline-flex items-center gap-2 rounded-lg bg-[#8ec07c] px-6 py-3 font-medium text-[#3c3836] transition-colors hover:bg-[#a0d28c]"
				>
					Unlock Integrations — Go Pro
					<ChevronRight class="h-4 w-4" />
				</a>
			</div>
		</div>
	</div>
</section>

<!-- Section 5: AI-Powered Features Coming Soon -->
<section id="ai" class="bg-[#3c3836] py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-4xl">
			<div class="mb-6 flex items-center gap-4">
				<div class="rounded-lg bg-[#d79921]/10 p-3">
					<BrainCircuit class="h-8 w-8 text-[#d79921]" />
				</div>
				<div>
					<h2 class="font-heading text-4xl font-bold text-white md:text-5xl">
						AI-Powered Intelligence
					</h2>
					<span
						class="mt-2 inline-block rounded-full bg-[#d79921]/10 px-3 py-1 text-sm font-semibold text-[#d79921]"
						>Coming Soon</span
					>
				</div>
			</div>
			<p class="mb-8 max-w-3xl text-xl text-white/80">
				AI baked in, not bolted on. Every AI feature is contextual and action-triggered —
				intelligence at the moment you need it, grounded in your real site data and GSC metrics.
			</p>
			<div class="mb-8 grid grid-cols-1 gap-6 md:grid-cols-2">
				<div class="relative rounded-lg border border-[#d79921]/20 bg-[#2d2826] p-6">
					<div class="mb-3 flex items-center gap-3">
						<PenTool class="h-5 w-5 text-[#d79921]" />
						<h3 class="font-heading text-lg font-bold text-white">Writing Voice Analysis</h3>
					</div>
					<p class="text-sm text-white/70">
						Generate a brand voice profile from your existing content. Every AI output matches your
						tone, structure, and style — not generic AI filler.
					</p>
				</div>
				<div class="relative rounded-lg border border-[#d79921]/20 bg-[#2d2826] p-6">
					<div class="mb-3 flex items-center gap-3">
						<ChartBar class="h-5 w-5 text-[#d79921]" />
						<h3 class="font-heading text-lg font-bold text-white">GSC Intelligence Dashboard</h3>
					</div>
					<p class="text-sm text-white/70">
						Quick wins, declining pages, content gaps — all surfaced from your GSC data with
						one-click AI explanations and weekly digest summaries.
					</p>
				</div>
				<div class="relative rounded-lg border border-[#d79921]/20 bg-[#2d2826] p-6">
					<div class="mb-3 flex items-center gap-3">
						<FileText class="h-5 w-5 text-[#d79921]" />
						<h3 class="font-heading text-lg font-bold text-white">Content Brief Generator</h3>
					</div>
					<p class="text-sm text-white/70">
						Data-backed content briefs with target keywords, outlines, internal link opportunities,
						and intent classification — generated from real GSC data.
					</p>
				</div>
				<div class="relative rounded-lg border border-[#d79921]/20 bg-[#2d2826] p-6">
					<div class="mb-3 flex items-center gap-3">
						<Sparkles class="h-5 w-5 text-[#d79921]" />
						<h3 class="font-heading text-lg font-bold text-white">Article Writer</h3>
					</div>
					<p class="text-sm text-white/70">
						Full article drafts from approved briefs, written in your brand voice. Markdown editor
						with live preview, internal link highlights, and one-click export.
					</p>
				</div>
				<div class="relative rounded-lg border border-[#d79921]/20 bg-[#2d2826] p-6 md:col-span-2">
					<div class="mb-3 flex items-center gap-3">
						<Link2 class="h-5 w-5 text-[#d79921]" />
						<h3 class="font-heading text-lg font-bold text-white">Internal Link Suggester</h3>
					</div>
					<p class="text-sm text-white/70">
						Semantic similarity across crawled pages surfaces missing internal links. Page-level
						suggestions with AI anchor text, plus a site-wide audit for orphaned pages and topic
						clusters.
					</p>
				</div>
			</div>
			<div class="text-center">
				<a
					href="/roadmap"
					class="inline-flex items-center gap-2 rounded-lg border-2 border-[#d79921]/30 px-6 py-3 font-medium text-white transition-colors hover:border-[#d79921]"
				>
					See Full Roadmap
					<ChevronRight class="h-4 w-4" />
				</a>
			</div>
		</div>
	</div>
</section>

<!-- Section 6: CLI for Developers -->
<section id="cli" class="bg-[#2d2826] py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-4xl">
			<div class="mb-6 flex items-center gap-4">
				<div class="rounded-lg bg-[#8ec07c]/10 p-3">
					<Code class="h-8 w-8 text-[#8ec07c]" />
				</div>
				<h2 class="font-heading text-4xl font-bold text-white md:text-5xl">
					Automate Audits with the Barracuda CLI
				</h2>
			</div>
			<p class="mb-8 max-w-3xl text-xl text-white/80">
				Developers, we've got you covered. The same powerful engine that powers the web app is also
				available as an open-source CLI. Integrate Barracuda into your CI/CD pipelines, trigger
				crawls via scripts, and export clean data for analysis.
			</p>
			<div class="mb-8 rounded-lg border border-white/10 bg-[#2d2826] p-6">
				<ul class="space-y-3 text-white/70">
					<li class="flex items-start gap-3">
						<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
						<span>Run local or remote crawls</span>
					</li>
					<li class="flex items-start gap-3">
						<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
						<span>JSON and CSV export support</span>
					</li>
					<li class="flex items-start gap-3">
						<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
						<span>Upload directly to the Barracuda Cloud workspace</span>
					</li>
					<li class="flex items-start gap-3">
						<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
						<span>Open-source under the MIT License</span>
					</li>
				</ul>
			</div>
			<div class="text-center">
				<a
					href="https://app.barracudaseo.com"
					class="inline-flex items-center gap-2 rounded-lg border-2 border-white/20 px-6 py-3 font-medium text-white transition-colors hover:border-[#8ec07c]"
					target="_blank"
					rel="noopener noreferrer"
				>
					Get Started
				</a>
			</div>
		</div>
	</div>
</section>

<!-- Section 7: Why Barracuda -->
<section class="bg-[#3c3836] py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-4xl">
			<h2 class="mb-12 text-center font-heading text-4xl font-bold text-white md:text-5xl">
				Why Agencies Choose Barracuda
			</h2>
			<div class="mb-12 grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-5">
				<div class="text-center">
					<div class="mb-3 inline-block rounded-lg bg-[#8ec07c]/10 p-3">
						<Bolt class="h-6 w-6 text-[#8ec07c]" />
					</div>
					<h3 class="mb-2 font-heading font-bold text-white">Fast</h3>
					<p class="text-sm text-white/70">Multi-threaded crawling engine built in Go</p>
				</div>
				<div class="text-center">
					<div class="mb-3 inline-block rounded-lg bg-[#8ec07c]/10 p-3">
						<Plug class="h-6 w-6 text-[#8ec07c]" />
					</div>
					<h3 class="mb-2 font-heading font-bold text-white">Flexible</h3>
					<p class="text-sm text-white/70">Use via web or command line</p>
				</div>
				<div class="text-center">
					<div class="mb-3 inline-block rounded-lg bg-[#8ec07c]/10 p-3">
						<CircleCheck class="h-6 w-6 text-[#8ec07c]" />
					</div>
					<h3 class="mb-2 font-heading font-bold text-white">Transparent</h3>
					<p class="text-sm text-white/70">Open-source foundation and honest pricing</p>
				</div>
				<div class="text-center">
					<div class="mb-3 inline-block rounded-lg bg-[#8ec07c]/10 p-3">
						<Target class="h-6 w-6 text-[#8ec07c]" />
					</div>
					<h3 class="mb-2 font-heading font-bold text-white">Actionable</h3>
					<p class="text-sm text-white/70">Prioritized, fix-ready insights</p>
				</div>
				<div class="text-center">
					<div class="mb-3 inline-block rounded-lg bg-[#8ec07c]/10 p-3">
						<Users class="h-6 w-6 text-[#8ec07c]" />
					</div>
					<h3 class="mb-2 font-heading font-bold text-white">Collaborative</h3>
					<p class="text-sm text-white/70">Built for teams, not just individuals</p>
				</div>
			</div>
			<div class="mb-8 rounded-lg border border-white/10 bg-[#3c3836] p-8">
				<blockquote class="mb-4 text-xl text-white/90 italic">
					"Barracuda makes technical SEO feel effortless. It's fast, accurate, and the dashboard is
					genuinely useful."
				</blockquote>
				<p class="text-sm text-white/70">— Early Access User, SEO Agency Owner</p>
			</div>
			<div class="text-center">
				<a
					href="https://app.barracudaseo.com"
					class="inline-flex items-center gap-2 rounded-lg bg-[#8ec07c] px-8 py-4 text-lg font-medium text-[#3c3836] transition-colors hover:bg-[#a0d28c]"
					target="_blank"
					rel="noopener noreferrer"
				>
					Start Your Free 100-Page Crawl
					<ChevronRight class="h-5 w-5" />
				</a>
			</div>
		</div>
	</div>
</section>

<!-- Section 8: Pricing Preview -->
<section class="bg-[#2d2826] py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-4xl">
			<h2 class="mb-12 text-center font-heading text-4xl font-bold text-white md:text-5xl">
				Simple, Fair Pricing
			</h2>
			<div class="overflow-hidden rounded-lg border border-white/10 bg-[#2d2826]">
				<div class="overflow-x-auto">
					<table class="w-full">
						<thead class="bg-[#3c3836]">
							<tr>
								<th class="px-6 py-4 text-left font-heading text-sm font-bold text-white">Plan</th>
								<th class="px-6 py-4 text-left font-heading text-sm font-bold text-white">Price</th>
								<th class="px-6 py-4 text-left font-heading text-sm font-bold text-white"
									>Crawl Limit</th
								>
								<th class="px-6 py-4 text-left font-heading text-sm font-bold text-white">Users</th>
								<th class="px-6 py-4 text-left font-heading text-sm font-bold text-white"
									>Integrations</th
								>
								<th class="px-6 py-4 text-left font-heading text-sm font-bold text-white"
									>Recommendations</th
								>
							</tr>
						</thead>
						<tbody class="divide-y divide-white/10">
							<tr>
								<td class="px-6 py-4 font-medium text-white">Free</td>
								<td class="px-6 py-4 text-white/70">$0</td>
								<td class="px-6 py-4 text-white/70">100 pages</td>
								<td class="px-6 py-4 text-white/70">1</td>
								<td class="px-6 py-4 text-white/50">—</td>
								<td class="px-6 py-4 text-white/50">—</td>
							</tr>
							<tr class="bg-[#8ec07c]/10">
								<td class="px-6 py-4 font-medium text-white">Pro</td>
								<td class="px-6 py-4 text-white">$29/mo</td>
								<td class="px-6 py-4 text-white">10,000+ pages</td>
								<td class="px-6 py-4 text-white">1</td>
								<td class="px-6 py-4 text-[#8ec07c]">✅</td>
								<td class="px-6 py-4 text-[#8ec07c]">✅</td>
							</tr>
							<tr>
								<td class="px-6 py-4 font-medium text-white">Team</td>
								<td class="px-6 py-4 text-white/70">+$5/user</td>
								<td class="px-6 py-4 text-white/70">25,000+ pages</td>
								<td class="px-6 py-4 text-white/70">Unlimited</td>
								<td class="px-6 py-4 text-[#8ec07c]">✅</td>
								<td class="px-6 py-4 text-[#8ec07c]">✅</td>
							</tr>
						</tbody>
					</table>
				</div>
			</div>
			<div class="mt-8 text-center">
				<a
					href="/pricing"
					class="inline-flex items-center gap-2 rounded-lg bg-[#8ec07c] px-6 py-3 font-medium text-[#3c3836] transition-colors hover:bg-[#a0d28c]"
				>
					Explore Plans & Upgrade
					<ChevronRight class="h-4 w-4" />
				</a>
			</div>
		</div>
	</div>
</section>

<RecentPosts
	title="SEO Guides"
	subtitle="Learn how to find, prioritize, and fix the issues that matter most."
/>

<!-- Section 8: Final CTA -->
<section class="bg-gradient-to-b from-[#2d2826] to-[#3c3836] py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-3xl text-center">
			<h2 class="mb-6 font-heading text-4xl font-bold text-white md:text-5xl">
				Your site deserves better SEO visibility.
			</h2>
			<p class="mb-10 text-xl text-white/80">
				Join hundreds of Agencies and developers who use Barracuda to uncover issues, prioritize
				fixes, and grow their traffic.
			</p>
			<div class="flex flex-col justify-center gap-4 sm:flex-row">
				<a
					href="https://app.barracudaseo.com"
					class="rounded-lg bg-[#8ec07c] px-8 py-4 text-lg font-medium text-[#3c3836] transition-colors hover:bg-[#a0d28c]"
					target="_blank"
					rel="noopener noreferrer"
				>
					Start Free — No Credit Card Required
				</a>
				<a
					href="#dashboard"
					class="rounded-lg border-2 border-white/20 px-8 py-4 text-lg font-medium text-white transition-colors hover:border-[#8ec07c]"
				>
					View Demo
				</a>
			</div>
		</div>
	</div>
</section>
