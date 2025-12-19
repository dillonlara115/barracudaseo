<script lang="ts">
	import { MapPin, Search, Target, TrendingUp, CircleCheck, CircleQuestionMark, ArrowRight } from '@lucide/svelte';
	import MetaTags from '../../../components/MetaTags.svelte';
	import { getMetaTags, getBreadcrumbSchema } from '$lib/meta';

	const meta = getMetaTags({
		title: 'Local SEO Auditing Tool',
		description: 'Optimize your local business website for local search with Barracuda SEO. Discover technical issues, improve Google Business Profile visibility, and dominate local search results.'
	});

	const structuredData = getBreadcrumbSchema([
		{ name: 'Home', url: '/' },
		{ name: 'Use Cases', url: '/use-cases' },
		{ name: 'Local SEO', url: '/use-cases/local-seo' }
	]);
</script>

<MetaTags config={{ ...meta, structuredData }} />

<!-- Hero Section -->
<section class="py-20 bg-gradient-to-b from-[#3c3836] to-[#2d2826]">
	<div class="container mx-auto px-4">
		<div class="text-center max-w-4xl mx-auto">
			<h1 class="text-5xl md:text-6xl font-heading font-bold mb-6 text-white">
				Dominate Local Search with Technical SEO Audits
			</h1>
			<p class="text-xl md:text-2xl text-white/80 mb-10 max-w-3xl mx-auto">
				Barracuda SEO helps local businesses discover and fix technical SEO issues that impact local search visibility, Google Business Profile performance, and "near me" search rankings.
			</p>
			<div class="flex flex-col sm:flex-row gap-4 justify-center">
				<a href="https://app.barracudaseo.com" class="bg-[#8ec07c] hover:bg-[#a0d28c] text-[#3c3836] px-8 py-4 rounded-lg font-medium text-lg transition-colors" target="_blank" rel="noopener noreferrer">
					Start Your Free Local SEO Audit
				</a>
			</div>
		</div>
	</div>
</section>

<!-- Dashboard Preview -->
<section class="py-12 bg-[#2d2826]">
	<div class="container mx-auto px-4">
		<div class="max-w-6xl mx-auto">
			<img src="/mockups/barracuda-dashboard.png" alt="Barracuda SEO Dashboard showing local SEO audit results" class="w-full rounded-lg shadow-2xl border border-white/10" />
		</div>
	</div>
</section>

<!-- Features Section -->
<section class="py-20 bg-[#3c3836]">
	<div class="container mx-auto px-4">
		<div class="max-w-6xl mx-auto">
			<h2 class="text-4xl md:text-5xl font-heading font-bold mb-12 text-center text-white">
				Your Complete Local SEO Technical Audit Tool
			</h2>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-8">
				<div class="bg-[#2d2826] rounded-lg p-8 border border-white/10">
					<div class="flex items-center gap-4 mb-4">
						<div class="p-3 bg-[#8ec07c]/10 rounded-lg">
							<MapPin class="w-6 h-6 text-[#8ec07c]" />
						</div>
						<h3 class="text-2xl font-heading font-bold text-white">Location Page Optimization</h3>
					</div>
					<p class="text-white/80 mb-4">
						Ensure your location pages, service area pages, and contact pages are properly structured with schema markup, consistent NAP (Name, Address, Phone) information, and optimized meta tags for local search queries.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Detect missing or inconsistent location data</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Identify duplicate location pages</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Verify schema markup implementation</span>
						</li>
					</ul>
				</div>

				<div class="bg-[#2d2826] rounded-lg p-8 border border-white/10">
					<div class="flex items-center gap-4 mb-4">
						<div class="p-3 bg-[#8ec07c]/10 rounded-lg">
							<Search class="w-6 h-6 text-[#8ec07c]" />
						</div>
						<h3 class="text-2xl font-heading font-bold text-white">Mobile-First Performance</h3>
					</div>
					<p class="text-white/80 mb-4">
						Local searches happen on mobile devices. Barracuda identifies slow-loading pages, mobile usability issues, and technical problems that hurt your local search rankings.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Page speed optimization for mobile users</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Mobile-friendly markup validation</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Click-to-call and map integration checks</span>
						</li>
					</ul>
				</div>

				<div class="bg-[#2d2826] rounded-lg p-8 border border-white/10">
					<div class="flex items-center gap-4 mb-4">
						<div class="p-3 bg-[#8ec07c]/10 rounded-lg">
							<Target class="w-6 h-6 text-[#8ec07c]" />
						</div>
						<h3 class="text-2xl font-heading font-bold text-white">Local Keyword Targeting</h3>
					</div>
					<p class="text-white/80 mb-4">
						Connect Google Search Console to see which local keywords drive traffic, identify pages with ranking potential, and prioritize fixes that impact "near me" and location-based searches.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>GSC integration for local search insights</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Track local keyword performance</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Identify high-opportunity local pages</span>
						</li>
					</ul>
				</div>

				<div class="bg-[#2d2826] rounded-lg p-8 border border-white/10">
					<div class="flex items-center gap-4 mb-4">
						<div class="p-3 bg-[#8ec07c]/10 rounded-lg">
							<TrendingUp class="w-6 h-6 text-[#8ec07c]" />
						</div>
						<h3 class="text-2xl font-heading font-bold text-white">Citation & Directory Consistency</h3>
					</div>
					<p class="text-white/80 mb-4">
						Maintain consistent business information across your website to improve local SEO signals. Barracuda helps you spot inconsistencies that confuse search engines and hurt local rankings.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>NAP consistency across all pages</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Local business schema validation</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Hours of operation and service area checks</span>
						</li>
					</ul>
				</div>
			</div>
		</div>
	</div>
</section>

<!-- How It Works Section -->
<section class="py-20 bg-[#2d2826]">
	<div class="container mx-auto px-4">
		<div class="max-w-4xl mx-auto">
			<h2 class="text-4xl md:text-5xl font-heading font-bold mb-4 text-center text-white">
				How It Works
			</h2>
			<p class="text-xl text-white/70 mb-12 text-center max-w-2xl mx-auto">
				Get your local SEO audit running in minutes
			</p>

			<div class="space-y-8">
				<div class="flex gap-6">
					<div class="flex-shrink-0">
						<div class="w-12 h-12 bg-[#8ec07c] rounded-full flex items-center justify-center text-[#3c3836] font-bold text-xl">
							1
						</div>
					</div>
					<div class="flex-grow">
						<h3 class="text-2xl font-heading font-bold mb-3 text-white">Start Your Crawl</h3>
						<p class="text-white/80 mb-4">
							Enter your local business website URL and let Barracuda crawl your site. Focus on location pages, service pages, and key landing pages that drive local traffic.
						</p>
					</div>
				</div>

				<div class="flex gap-6">
					<div class="flex-shrink-0">
						<div class="w-12 h-12 bg-[#8ec07c] rounded-full flex items-center justify-center text-[#3c3836] font-bold text-xl">
							2
						</div>
					</div>
					<div class="flex-grow">
						<h3 class="text-2xl font-heading font-bold mb-3 text-white">Review Technical Issues</h3>
						<p class="text-white/80 mb-4">
							Barracuda identifies broken links, missing meta tags, duplicate content, slow pages, and schema markup issues that impact your local search visibility.
						</p>
						<div class="mt-6">
							<img src="/mockups/barracuda-issues.png" alt="Barracuda SEO Issues Dashboard showing local SEO issues" class="w-full rounded-lg shadow-lg border border-white/10" />
						</div>
					</div>
				</div>

				<div class="flex gap-6">
					<div class="flex-shrink-0">
						<div class="w-12 h-12 bg-[#8ec07c] rounded-full flex items-center justify-center text-[#3c3836] font-bold text-xl">
							3
						</div>
					</div>
					<div class="flex-grow">
						<h3 class="text-2xl font-heading font-bold mb-3 text-white">Connect Google Search Console</h3>
						<p class="text-white/80 mb-4">
							Link your GSC account to see which local keywords drive traffic, identify pages with ranking potential, and prioritize fixes based on actual search performance.
						</p>
					</div>
				</div>

				<div class="flex gap-6">
					<div class="flex-shrink-0">
						<div class="w-12 h-12 bg-[#8ec07c] rounded-full flex items-center justify-center text-[#3c3836] font-bold text-xl">
							4
						</div>
					</div>
					<div class="flex-grow">
						<h3 class="text-2xl font-heading font-bold mb-3 text-white">Fix & Monitor</h3>
						<p class="text-white/80 mb-4">
							Use Barracuda's prioritized issue list to fix technical problems, then track improvements over time with historical crawl comparisons and GSC performance data.
						</p>
					</div>
				</div>
			</div>
		</div>
	</div>
</section>

<!-- FAQs Section -->
<section class="py-20 bg-[#3c3836]">
	<div class="container mx-auto px-4">
		<div class="max-w-4xl mx-auto">
			<div class="flex items-center justify-center gap-4 mb-8">
				<div class="p-3 bg-[#8ec07c]/10 rounded-lg">
					<CircleQuestionMark class="w-8 h-8 text-[#8ec07c]" />
				</div>
				<h2 class="text-4xl md:text-5xl font-heading font-bold text-white">
					Local SEO FAQs
				</h2>
			</div>

			<div class="space-y-4">
				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<div class="p-6">
						<h3 class="text-xl font-heading font-bold mb-3 text-white">
							How does Barracuda help with local SEO specifically?
						</h3>
						<p class="text-white/80">
							Barracuda identifies technical SEO issues that impact local search rankings, such as missing location schema, inconsistent NAP information, mobile usability problems, and slow page speeds. By fixing these technical issues, you improve your site's ability to rank for local searches and support your Google Business Profile visibility.
						</p>
					</div>
				</div>

				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<div class="p-6">
						<h3 class="text-xl font-heading font-bold mb-3 text-white">
							Can Barracuda check my Google Business Profile?
						</h3>
						<p class="text-white/80">
							Barracuda focuses on technical SEO audits of your website. While we don't directly audit your Google Business Profile, fixing technical issues on your site (like schema markup, NAP consistency, and mobile performance) helps improve your overall local search visibility and supports your GBP performance.
						</p>
					</div>
				</div>

				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<div class="p-6">
						<h3 class="text-xl font-heading font-bold mb-3 text-white">
							What local SEO issues does Barracuda detect?
						</h3>
						<p class="text-white/80">
							Barracuda detects missing or incorrect local business schema, inconsistent name/address/phone information across pages, duplicate location pages, mobile usability issues, slow-loading pages that hurt mobile search, and broken links that impact user experience for local searchers.
						</p>
					</div>
				</div>

				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<div class="p-6">
						<h3 class="text-xl font-heading font-bold mb-3 text-white">
							How often should I audit my local business website?
						</h3>
						<p class="text-white/80">
							We recommend running a technical SEO audit monthly or whenever you make significant changes to your website, add new location pages, or update your business information. Regular audits help you catch issues before they impact your local search rankings.
						</p>
					</div>
				</div>
			</div>
		</div>
	</div>
</section>

<!-- CTA Section -->
<section class="py-20 bg-gradient-to-b from-[#2d2826] to-[#3c3836]">
	<div class="container mx-auto px-4">
		<div class="max-w-3xl mx-auto text-center">
			<h2 class="text-4xl md:text-5xl font-heading font-bold mb-6 text-white">
				Ready to improve your local search visibility?
			</h2>
			<p class="text-xl text-white/80 mb-10">
				Start your free 100-page audit and discover the technical issues holding back your local SEO performance.
			</p>
			<a href="https://app.barracudaseo.com" class="inline-flex items-center gap-2 bg-[#8ec07c] hover:bg-[#a0d28c] text-[#3c3836] px-8 py-4 rounded-lg font-medium text-lg transition-colors" target="_blank" rel="noopener noreferrer">
				Start Your Free Local SEO Audit
				<ArrowRight class="w-5 h-5" />
			</a>
		</div>
	</div>
</section>

