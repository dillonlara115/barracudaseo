<script lang="ts">
	import {
		MapPin,
		Search,
		Target,
		TrendingUp,
		CircleCheck,
		CircleQuestionMark,
		ArrowRight
	} from '@lucide/svelte';
	import MetaTags from '../../../components/MetaTags.svelte';
	import OptimizedImage from '../../../components/OptimizedImage.svelte';
	import { getMetaTags, getBreadcrumbSchema, getFAQPageSchema } from '$lib/meta';

	const meta = getMetaTags({
		title: 'Local SEO Auditing Tool',
		description:
			'Optimize your local business website with Barracuda SEO. Discover technical issues, check schema markup, and improve local search visibility today.'
	});

	const faqs = [
		{
			question: 'How does Barracuda help with local SEO specifically?',
			answer:
				"Barracuda identifies technical SEO issues that impact local search rankings, such as missing location schema, inconsistent NAP information, mobile usability problems, and slow page speeds. By fixing these technical issues, you improve your site's ability to rank for local searches and support your Google Business Profile visibility."
		},
		{
			question: 'Can Barracuda check my Google Business Profile?',
			answer:
				"Barracuda focuses on technical SEO audits of your website. While we don't directly audit your Google Business Profile, fixing technical issues on your site (like schema markup, NAP consistency, and mobile performance) helps improve your overall local search visibility and supports your GBP performance."
		},
		{
			question: 'What local SEO issues does Barracuda detect?',
			answer:
				'Barracuda detects missing or incorrect local business schema, inconsistent name/address/phone information across pages, duplicate location pages, mobile usability issues, slow-loading pages that hurt mobile search, and broken links that impact user experience for local searchers.'
		},
		{
			question: 'How often should I audit my local business website?',
			answer:
				'We recommend running a technical SEO audit monthly or whenever you make significant changes to your website, add new location pages, or update your business information. Regular audits help you catch issues before they impact your local search rankings.'
		}
	];

	const structuredData = [
		getBreadcrumbSchema([
			{ name: 'Home', url: '/' },
			{ name: 'Use Cases', url: '/use-cases' },
			{ name: 'Local SEO', url: '/use-cases/local-seo' }
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
				Dominate Local Search with Technical SEO Audits
			</h1>
			<p class="mx-auto mb-10 max-w-3xl text-xl text-white/80 md:text-2xl">
				Barracuda SEO helps local businesses discover and fix technical SEO issues that impact local
				search visibility, Google Business Profile performance, and "near me" search rankings.
			</p>
			<div class="flex flex-col justify-center gap-4 sm:flex-row">
				<a
					href="https://app.barracudaseo.com"
					class="rounded-lg bg-[#8ec07c] px-8 py-4 text-lg font-medium text-[#3c3836] transition-colors hover:bg-[#a0d28c]"
					target="_blank"
					rel="noopener noreferrer"
				>
					Start Your Free Local SEO Audit
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
				alt="Barracuda SEO Dashboard showing local SEO audit results"
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
				Your Complete Local SEO Technical Audit Tool
			</h2>
			<div class="grid grid-cols-1 gap-8 md:grid-cols-2">
				<div class="rounded-lg border border-white/10 bg-[#2d2826] p-8">
					<div class="mb-4 flex items-center gap-4">
						<div class="rounded-lg bg-[#8ec07c]/10 p-3">
							<MapPin class="h-6 w-6 text-[#8ec07c]" />
						</div>
						<h3 class="font-heading text-2xl font-bold text-white">Location Page Optimization</h3>
					</div>
					<p class="mb-4 text-white/80">
						Ensure your location pages, service area pages, and contact pages are properly
						structured with schema markup, consistent NAP (Name, Address, Phone) information, and
						optimized meta tags for local search queries.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Detect missing or inconsistent location data</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Identify duplicate location pages</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Verify schema markup implementation</span>
						</li>
					</ul>
				</div>

				<div class="rounded-lg border border-white/10 bg-[#2d2826] p-8">
					<div class="mb-4 flex items-center gap-4">
						<div class="rounded-lg bg-[#8ec07c]/10 p-3">
							<Search class="h-6 w-6 text-[#8ec07c]" />
						</div>
						<h3 class="font-heading text-2xl font-bold text-white">Mobile-First Performance</h3>
					</div>
					<p class="mb-4 text-white/80">
						Local searches happen on mobile devices. Barracuda identifies slow-loading pages, mobile
						usability issues, and technical problems that hurt your local search rankings.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Page speed optimization for mobile users</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Mobile-friendly markup validation</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Click-to-call and map integration checks</span>
						</li>
					</ul>
				</div>

				<div class="rounded-lg border border-white/10 bg-[#2d2826] p-8">
					<div class="mb-4 flex items-center gap-4">
						<div class="rounded-lg bg-[#8ec07c]/10 p-3">
							<Target class="h-6 w-6 text-[#8ec07c]" />
						</div>
						<h3 class="font-heading text-2xl font-bold text-white">Local Keyword Targeting</h3>
					</div>
					<p class="mb-4 text-white/80">
						Connect Google Search Console to see which local keywords drive traffic, identify pages
						with ranking potential, and prioritize fixes that impact "near me" and location-based
						searches.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>GSC integration for local search insights</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Track local keyword performance</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Identify high-opportunity local pages</span>
						</li>
					</ul>
				</div>

				<div class="rounded-lg border border-white/10 bg-[#2d2826] p-8">
					<div class="mb-4 flex items-center gap-4">
						<div class="rounded-lg bg-[#8ec07c]/10 p-3">
							<TrendingUp class="h-6 w-6 text-[#8ec07c]" />
						</div>
						<h3 class="font-heading text-2xl font-bold text-white">
							Citation & Directory Consistency
						</h3>
					</div>
					<p class="mb-4 text-white/80">
						Maintain consistent business information across your website to improve local SEO
						signals. Barracuda helps you spot inconsistencies that confuse search engines and hurt
						local rankings.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>NAP consistency across all pages</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Local business schema validation</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Hours of operation and service area checks</span>
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
				Get your local SEO audit running in minutes
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
						<h3 class="mb-3 font-heading text-2xl font-bold text-white">Start Your Crawl</h3>
						<p class="mb-4 text-white/80">
							Enter your local business website URL and let Barracuda crawl your site. Focus on
							location pages, service pages, and key landing pages that drive local traffic.
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
						<h3 class="mb-3 font-heading text-2xl font-bold text-white">Review Technical Issues</h3>
						<p class="mb-4 text-white/80">
							Barracuda identifies broken links, <a
								href="/blog/are-missing-meta-descriptions-important"
								class="text-[#8ec07c] underline hover:text-[#a0d28c]">missing meta tags</a
							>, duplicate content, slow pages, and schema markup issues that impact your local
							search visibility.
						</p>
						<div class="mt-6">
							<OptimizedImage
								src="/mockups/barracuda-issues.png"
								alt="Barracuda SEO Issues Dashboard showing local SEO issues"
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
							Connect Google Search Console
						</h3>
						<p class="mb-4 text-white/80">
							Link your GSC account to see which local keywords drive traffic, identify pages with
							ranking potential, and prioritize fixes based on actual search performance.
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
						<h3 class="mb-3 font-heading text-2xl font-bold text-white">Fix & Monitor</h3>
						<p class="mb-4 text-white/80">
							Use Barracuda's prioritized issue list to fix technical problems, then track
							improvements over time with historical crawl comparisons and GSC performance data.
						</p>
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
				<h2 class="font-heading text-4xl font-bold text-white md:text-5xl">Local SEO FAQs</h2>
			</div>

			<div class="space-y-4">
				<div class="overflow-hidden rounded-lg border border-white/10 bg-[#2d2826]">
					<div class="p-6">
						<h3 class="mb-3 font-heading text-xl font-bold text-white">
							How does Barracuda help with local SEO specifically?
						</h3>
						<p class="text-white/80">
							Barracuda identifies technical SEO issues that impact local search rankings, such as
							missing location schema, inconsistent NAP information, mobile usability problems, and
							slow page speeds. By fixing these technical issues, you improve your site's ability to
							rank for local searches and support your Google Business Profile visibility.
						</p>
					</div>
				</div>

				<div class="overflow-hidden rounded-lg border border-white/10 bg-[#2d2826]">
					<div class="p-6">
						<h3 class="mb-3 font-heading text-xl font-bold text-white">
							Can Barracuda check my Google Business Profile?
						</h3>
						<p class="text-white/80">
							Barracuda focuses on technical SEO audits of your website. While we don't directly
							audit your Google Business Profile, fixing technical issues on your site (like schema
							markup, NAP consistency, and mobile performance) helps improve your overall local
							search visibility and supports your GBP performance.
						</p>
					</div>
				</div>

				<div class="overflow-hidden rounded-lg border border-white/10 bg-[#2d2826]">
					<div class="p-6">
						<h3 class="mb-3 font-heading text-xl font-bold text-white">
							What local SEO issues does Barracuda detect?
						</h3>
						<p class="text-white/80">
							Barracuda detects missing or incorrect local business schema, inconsistent
							name/address/phone information across pages, duplicate location pages, mobile
							usability issues, slow-loading pages that hurt mobile search, and broken links that
							impact user experience for local searchers.
						</p>
					</div>
				</div>

				<div class="overflow-hidden rounded-lg border border-white/10 bg-[#2d2826]">
					<div class="p-6">
						<h3 class="mb-3 font-heading text-xl font-bold text-white">
							How often should I audit my local business website?
						</h3>
						<p class="text-white/80">
							We recommend running a <a
								href="/blog/complete-technical-seo-audit-guide"
								class="text-[#8ec07c] underline hover:text-[#a0d28c]">technical SEO audit</a
							> monthly or whenever you make significant changes to your website, add new location pages,
							or update your business information. Regular audits help you catch issues before they impact
							your local search rankings.
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
				Ready to improve your local search visibility?
			</h2>
			<p class="mb-10 text-xl text-white/80">
				Start your free 100-page audit and discover the technical issues holding back your local SEO
				performance.
			</p>
			<a
				href="https://app.barracudaseo.com"
				class="inline-flex items-center gap-2 rounded-lg bg-[#8ec07c] px-8 py-4 text-lg font-medium text-[#3c3836] transition-colors hover:bg-[#a0d28c]"
				target="_blank"
				rel="noopener noreferrer"
			>
				Start Your Free Local SEO Audit
				<ArrowRight class="h-5 w-5" />
			</a>
		</div>
	</div>
</section>
