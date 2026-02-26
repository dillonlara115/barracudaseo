<script lang="ts">
	import {
		Zap,
		Layers,
		TrendingUp,
		ChartBar,
		CircleCheck,
		CircleQuestionMark,
		ArrowRight
	} from '@lucide/svelte';
	import MetaTags from '../../../components/MetaTags.svelte';
	import OptimizedImage from '../../../components/OptimizedImage.svelte';
	import { getMetaTags, getBreadcrumbSchema, getFAQPageSchema } from '$lib/meta';

	const meta = getMetaTags({
		title: 'Programmatic SEO Tool',
		description:
			'Scale your SEO content strategy with Barracuda SEO. Audit thousands of programmatic pages, identify technical issues, and optimize at scale today.'
	});

	const faqs = [
		{
			question: 'How does Barracuda help with programmatic SEO?',
			answer:
				'Barracuda audits thousands of programmatic pages efficiently, identifying technical SEO issues at scale. It groups issues by URL pattern to help you fix template-level problems, connects with Google Search Console to prioritize fixes based on performance, and provides bulk export capabilities for analysis.'
		},
		{
			question: 'Can Barracuda handle large-scale programmatic sites?',
			answer:
				"Yes. Barracuda's Pro plan supports crawling 10,000+ pages, making it ideal for programmatic SEO sites. The fast, concurrent crawling engine handles large-scale audits efficiently, and you can export results for bulk analysis."
		},
		{
			question: 'How does Barracuda identify template-level issues?',
			answer:
				'Barracuda groups issues by URL pattern, helping you identify problems that affect multiple programmatic pages. This makes it easy to fix template-level issues like duplicate meta tags, missing schema markup, or inconsistent content structure across your programmatic pages.'
		},
		{
			question: 'Can I export programmatic SEO audit results?',
			answer:
				'Yes. Barracuda supports CSV and JSON export, making it easy to analyze programmatic page issues in bulk. You can also track improvements over time with historical crawl comparisons.'
		}
	];

	const structuredData = [
		getBreadcrumbSchema([
			{ name: 'Home', url: '/' },
			{ name: 'Use Cases', url: '/use-cases' },
			{ name: 'Programmatic SEO', url: '/use-cases/programmatic-seo' }
		]),
		getFAQPageSchema(faqs)
	];
</script>

<MetaTags config={{ ...meta, structuredData }} />

<!-- Hero Section -->
<section class="bg-gradient-to-b from-[#3c3836] to-[#2d2826] py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-4xl text-center">
			<h1 class="mb-6 font-heading text-5xl font-bold text-white md:text-6xl">
				Audit Programmatic SEO Pages at Scale
			</h1>
			<p class="mx-auto mb-10 max-w-3xl text-xl text-white/80 md:text-2xl">
				Barracuda SEO helps you discover and fix technical SEO issues across thousands of
				programmatic pages, ensuring your scaled content strategy performs in search results.
			</p>
			<div class="flex flex-col justify-center gap-4 sm:flex-row">
				<a
					href="https://app.barracudaseo.com"
					class="rounded-lg bg-[#8ec07c] px-8 py-4 text-lg font-medium text-[#3c3836] transition-colors hover:bg-[#a0d28c]"
					target="_blank"
					rel="noopener noreferrer"
				>
					Start Your Programmatic SEO Audit
				</a>
			</div>
		</div>
	</div>
</section>

<!-- Dashboard Preview -->
<section class="bg-[#2d2826] py-12">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-6xl">
			<OptimizedImage
				src="/mockups/barracuda-dashboard.png"
				alt="Barracuda SEO Dashboard showing programmatic SEO audit at scale"
				class="w-full rounded-lg border border-white/10 shadow-2xl"
				loading="lazy"
				aspectRatio="16/9"
			/>
		</div>
	</div>
</section>

<!-- Features Section -->
<section class="bg-[#3c3836] py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-6xl">
			<h2 class="mb-12 text-center font-heading text-4xl font-bold text-white md:text-5xl">
				Your Complete Programmatic SEO Audit Tool
			</h2>
			<div class="grid grid-cols-1 gap-8 md:grid-cols-2">
				<div class="rounded-lg border border-white/10 bg-[#2d2826] p-8">
					<div class="mb-4 flex items-center gap-4">
						<div class="rounded-lg bg-[#8ec07c]/10 p-3">
							<Zap class="h-6 w-6 text-[#8ec07c]" />
						</div>
						<h3 class="font-heading text-2xl font-bold text-white">Scale & Performance</h3>
					</div>
					<p class="mb-4 text-white/80">
						Crawl thousands of programmatic pages efficiently. Barracuda's optimized engine handles
						large-scale audits without breaking, identifying issues across your entire programmatic
						content library.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Crawl 10,000+ pages with Pro plan</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Fast, concurrent crawling engine</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Export results for bulk analysis</span>
						</li>
					</ul>
				</div>

				<div class="rounded-lg border border-white/10 bg-[#2d2826] p-8">
					<div class="mb-4 flex items-center gap-4">
						<div class="rounded-lg bg-[#8ec07c]/10 p-3">
							<Layers class="h-6 w-6 text-[#8ec07c]" />
						</div>
						<h3 class="font-heading text-2xl font-bold text-white">
							Template-Level Issue Detection
						</h3>
					</div>
					<p class="mb-4 text-white/80">
						Identify patterns across programmatic pages. Barracuda groups issues by URL structure,
						helping you fix template-level problems that affect hundreds or thousands of pages at
						once.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Group issues by URL pattern</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Identify duplicate content across templates</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Detect missing or inconsistent meta tags</span>
						</li>
					</ul>
				</div>

				<div class="rounded-lg border border-white/10 bg-[#2d2826] p-8">
					<div class="mb-4 flex items-center gap-4">
						<div class="rounded-lg bg-[#8ec07c]/10 p-3">
							<TrendingUp class="h-6 w-6 text-[#8ec07c]" />
						</div>
						<h3 class="font-heading text-2xl font-bold text-white">Data-Driven Prioritization</h3>
					</div>
					<p class="mb-4 text-white/80">
						Connect Google Search Console to see which programmatic pages drive traffic and
						conversions. Prioritize fixes on high-performing pages and identify opportunities on
						underperforming ones.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>GSC integration for performance data</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Identify high-traffic programmatic pages</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Focus fixes on pages that matter most</span>
						</li>
					</ul>
				</div>

				<div class="rounded-lg border border-white/10 bg-[#2d2826] p-8">
					<div class="mb-4 flex items-center gap-4">
						<div class="rounded-lg bg-[#8ec07c]/10 p-3">
							<ChartBar class="h-6 w-6 text-[#8ec07c]" />
						</div>
						<h3 class="font-heading text-2xl font-bold text-white">Bulk Issue Resolution</h3>
					</div>
					<p class="mb-4 text-white/80">
						Export crawl results to CSV or JSON for bulk analysis. Use Barracuda's issue grouping to
						fix template-level problems efficiently, then track improvements with historical
						comparisons.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>CSV and JSON export for analysis</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Historical crawl comparisons</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Track improvements over time</span>
						</li>
					</ul>
				</div>
			</div>
		</div>
	</div>
</section>

<!-- How It Works Section -->
<section class="bg-[#2d2826] py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-4xl">
			<h2 class="mb-4 text-center font-heading text-4xl font-bold text-white md:text-5xl">
				How It Works
			</h2>
			<p class="mx-auto mb-12 max-w-2xl text-center text-xl text-white/70">
				Audit your programmatic SEO pages in minutes
			</p>

			<div class="space-y-8">
				<div class="flex gap-6">
					<div class="flex-shrink-0">
						<div
							class="flex h-12 w-12 items-center justify-center rounded-full bg-[#8ec07c] text-xl font-bold text-[#3c3836]"
						>
							1
						</div>
					</div>
					<div class="flex-grow">
						<h3 class="mb-3 font-heading text-2xl font-bold text-white">
							Crawl Your Programmatic Pages
						</h3>
						<p class="mb-4 text-white/80">
							Enter your website URL and let Barracuda crawl your programmatic pages. Use sitemap
							seeding or custom depth limits to focus on your programmatic content structure.
						</p>
					</div>
				</div>

				<div class="flex gap-6">
					<div class="flex-shrink-0">
						<div
							class="flex h-12 w-12 items-center justify-center rounded-full bg-[#8ec07c] text-xl font-bold text-[#3c3836]"
						>
							2
						</div>
					</div>
					<div class="flex-grow">
						<h3 class="mb-3 font-heading text-2xl font-bold text-white">
							Identify Template-Level Issues
						</h3>
						<p class="mb-4 text-white/80">
							Barracuda groups issues by URL pattern, helping you identify template-level problems
							affecting multiple pages. See duplicate content, missing meta tags, and schema markup
							issues at scale.
						</p>
						<div class="mt-6">
							<OptimizedImage
								src="/mockups/barracuda-results.png"
								alt="Barracuda SEO Results showing programmatic pages grouped by pattern"
								class="w-full rounded-lg border border-white/10 shadow-lg"
								loading="lazy"
								aspectRatio="16/9"
							/>
						</div>
					</div>
				</div>

				<div class="flex gap-6">
					<div class="flex-shrink-0">
						<div
							class="flex h-12 w-12 items-center justify-center rounded-full bg-[#8ec07c] text-xl font-bold text-[#3c3836]"
						>
							3
						</div>
					</div>
					<div class="flex-grow">
						<h3 class="mb-3 font-heading text-2xl font-bold text-white">
							Connect Performance Data
						</h3>
						<p class="mb-4 text-white/80">
							Link Google Search Console to see which programmatic pages drive traffic. Prioritize
							fixes on high-performing pages and identify opportunities on underperforming ones.
						</p>
					</div>
				</div>

				<div class="flex gap-6">
					<div class="flex-shrink-0">
						<div
							class="flex h-12 w-12 items-center justify-center rounded-full bg-[#8ec07c] text-xl font-bold text-[#3c3836]"
						>
							4
						</div>
					</div>
					<div class="flex-grow">
						<h3 class="mb-3 font-heading text-2xl font-bold text-white">Fix & Scale</h3>
						<p class="mb-4 text-white/80">
							Export results for bulk analysis, fix template-level issues efficiently, and track
							improvements over time with historical crawl comparisons.
						</p>
						<div class="mt-6">
							<OptimizedImage
								src="/mockups/barracuda-link-graph.png"
								alt="Barracuda SEO Link Graph showing programmatic site structure"
								class="w-full rounded-lg border border-white/10 shadow-lg"
								loading="lazy"
								aspectRatio="16/9"
							/>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</section>

<!-- FAQs Section -->
<section class="bg-[#3c3836] py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-4xl">
			<div class="mb-8 flex items-center justify-center gap-4">
				<div class="rounded-lg bg-[#8ec07c]/10 p-3">
					<CircleQuestionMark class="h-8 w-8 text-[#8ec07c]" />
				</div>
				<h2 class="font-heading text-4xl font-bold text-white md:text-5xl">
					Programmatic SEO FAQs
				</h2>
			</div>

			<div class="space-y-4">
				<div class="overflow-hidden rounded-lg border border-white/10 bg-[#2d2826]">
					<div class="p-6">
						<h3 class="mb-3 font-heading text-xl font-bold text-white">
							How does Barracuda help with programmatic SEO?
						</h3>
						<p class="text-white/80">
							Barracuda audits thousands of programmatic pages efficiently, identifying <a
								href="/blog/complete-technical-seo-audit-guide"
								class="text-[#8ec07c] underline hover:text-[#a0d28c]">technical SEO issues</a
							> at scale. It groups issues by URL pattern to help you fix template-level problems, connects
							with Google Search Console to prioritize fixes based on performance, and provides bulk export
							capabilities for analysis.
						</p>
					</div>
				</div>

				<div class="overflow-hidden rounded-lg border border-white/10 bg-[#2d2826]">
					<div class="p-6">
						<h3 class="mb-3 font-heading text-xl font-bold text-white">
							Can Barracuda handle large-scale programmatic sites?
						</h3>
						<p class="text-white/80">
							Yes. Barracuda's Pro plan supports crawling <a
								href="/blog/audit-large-sites-10000-pages"
								class="text-[#8ec07c] underline hover:text-[#a0d28c]">10,000+ pages</a
							>, making it ideal for programmatic SEO sites. The fast, concurrent crawling engine
							handles large-scale audits efficiently, and you can export results for bulk analysis.
						</p>
					</div>
				</div>

				<div class="overflow-hidden rounded-lg border border-white/10 bg-[#2d2826]">
					<div class="p-6">
						<h3 class="mb-3 font-heading text-xl font-bold text-white">
							How does Barracuda identify template-level issues?
						</h3>
						<p class="text-white/80">
							Barracuda groups issues by URL pattern, helping you identify problems that affect
							multiple programmatic pages. This makes it easy to fix template-level issues like
							duplicate meta tags, missing schema markup, or inconsistent content structure across
							your programmatic pages.
						</p>
					</div>
				</div>

				<div class="overflow-hidden rounded-lg border border-white/10 bg-[#2d2826]">
					<div class="p-6">
						<h3 class="mb-3 font-heading text-xl font-bold text-white">
							Can I export programmatic SEO audit results?
						</h3>
						<p class="text-white/80">
							Yes. Barracuda supports CSV and JSON export, making it easy to analyze programmatic
							page issues in bulk. You can also track improvements over time with historical crawl
							comparisons.
						</p>
					</div>
				</div>
			</div>
		</div>
	</div>
</section>

<!-- CTA Section -->
<section class="bg-gradient-to-b from-[#2d2826] to-[#3c3836] py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-3xl text-center">
			<h2 class="mb-6 font-heading text-4xl font-bold text-white md:text-5xl">
				Ready to audit your programmatic SEO pages?
			</h2>
			<p class="mb-10 text-xl text-white/80">
				Start your free 100-page audit and discover the technical issues affecting your programmatic
				content strategy.
			</p>
			<a
				href="https://app.barracudaseo.com"
				class="inline-flex items-center gap-2 rounded-lg bg-[#8ec07c] px-8 py-4 text-lg font-medium text-[#3c3836] transition-colors hover:bg-[#a0d28c]"
				target="_blank"
				rel="noopener noreferrer"
			>
				Start Your Programmatic SEO Audit
				<ArrowRight class="h-5 w-5" />
			</a>
		</div>
	</div>
</section>
