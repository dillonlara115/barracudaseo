<script lang="ts">
	import {
		ShoppingCart,
		Package,
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
		title: 'E-commerce SEO Auditing Tool',
		description:
			'Optimize your e-commerce site for search with Barracuda SEO. Discover technical issues, improve product page performance, and boost organic traffic today.'
	});

	const faqs = [
		{
			question: 'How does Barracuda help with e-commerce SEO?',
			answer:
				'Barracuda identifies technical SEO issues specific to e-commerce sites, including missing product schema, duplicate content, broken category links, slow-loading product pages, and navigation problems. By fixing these issues, you improve product page rankings, organic traffic, and conversion rates.'
		},
		{
			question: 'Can Barracuda handle large e-commerce catalogs?',
			answer:
				"Yes. Barracuda's Pro plan supports crawling 10,000+ pages, making it ideal for large e-commerce sites. The fast crawling engine efficiently audits product catalogs, and you can use sitemap seeding to ensure comprehensive coverage."
		},
		{
			question: 'What e-commerce SEO issues does Barracuda detect?',
			answer:
				'Barracuda detects missing product schema markup, duplicate product content, thin product pages, broken category and filter links, orphaned products, slow-loading pages, mobile usability issues, and internal linking problems that impact e-commerce SEO performance.'
		},
		{
			question: 'How can I prioritize fixes for my product pages?',
			answer:
				"Connect Google Search Console and Analytics to see which product pages drive traffic and conversions. Barracuda's integration helps you prioritize fixes on high-value pages that impact revenue, ensuring you focus your SEO efforts where they matter most."
		}
	];

	const structuredData = [
		getBreadcrumbSchema([
			{ name: 'Home', url: '/' },
			{ name: 'Use Cases', url: '/use-cases' },
			{ name: 'E-commerce', url: '/use-cases/e-commerce' }
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
				Optimize E-commerce Sites for Search Success
			</h1>
			<p class="mx-auto mb-10 max-w-3xl text-xl text-white/80 md:text-2xl">
				Barracuda SEO helps e-commerce businesses discover and fix technical SEO issues that impact
				product page rankings, organic traffic, and conversion rates.
			</p>
			<div class="flex flex-col justify-center gap-4 sm:flex-row">
				<a
					href="https://app.barracudaseo.com"
					class="rounded-lg bg-[#8ec07c] px-8 py-4 text-lg font-medium text-[#3c3836] transition-colors hover:bg-[#a0d28c]"
					target="_blank"
					rel="noopener noreferrer"
				>
					Start Your E-commerce SEO Audit
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
				alt="Barracuda SEO Dashboard showing e-commerce SEO audit results"
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
				Your Complete E-commerce SEO Audit Tool
			</h2>
			<div class="grid grid-cols-1 gap-8 md:grid-cols-2">
				<div class="rounded-lg border border-white/10 bg-[#2d2826] p-8">
					<div class="mb-4 flex items-center gap-4">
						<div class="rounded-lg bg-[#8ec07c]/10 p-3">
							<ShoppingCart class="h-6 w-6 text-[#8ec07c]" />
						</div>
						<h3 class="font-heading text-2xl font-bold text-white">Product Page Optimization</h3>
					</div>
					<p class="mb-4 text-white/80">
						Ensure your product pages are properly optimized for search. Barracuda identifies
						missing product schema, <a
							href="/blog/duplicate-meta-tags-fix"
							class="text-[#8ec07c] underline hover:text-[#a0d28c]">duplicate content</a
						>, thin content issues, and optimization opportunities across your catalog.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Product schema markup validation</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Duplicate product content detection</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Missing meta descriptions and titles</span>
						</li>
					</ul>
				</div>

				<div class="rounded-lg border border-white/10 bg-[#2d2826] p-8">
					<div class="mb-4 flex items-center gap-4">
						<div class="rounded-lg bg-[#8ec07c]/10 p-3">
							<Package class="h-6 w-6 text-[#8ec07c]" />
						</div>
						<h3 class="font-heading text-2xl font-bold text-white">
							Category & Navigation Structure
						</h3>
					</div>
					<p class="mb-4 text-white/80">
						Optimize your site architecture for both users and search engines. Barracuda identifies
						<a
							href="/blog/find-fix-broken-links"
							class="text-[#8ec07c] underline hover:text-[#a0d28c]">broken category links</a
						>, orphaned product pages, navigation issues, and internal linking problems that hurt
						SEO performance.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Broken category and filter links</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Orphaned product pages</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Internal linking structure analysis</span>
						</li>
					</ul>
				</div>

				<div class="rounded-lg border border-white/10 bg-[#2d2826] p-8">
					<div class="mb-4 flex items-center gap-4">
						<div class="rounded-lg bg-[#8ec07c]/10 p-3">
							<TrendingUp class="h-6 w-6 text-[#8ec07c]" />
						</div>
						<h3 class="font-heading text-2xl font-bold text-white">
							Performance-Driven Prioritization
						</h3>
					</div>
					<p class="mb-4 text-white/80">
						Connect Google Search Console and Analytics to see which product pages drive traffic and
						conversions. Prioritize fixes on high-value pages and identify opportunities on
						underperforming products.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>GSC and GA4 integration</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Identify high-traffic product pages</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Conversion-focused issue prioritization</span>
						</li>
					</ul>
				</div>

				<div class="rounded-lg border border-white/10 bg-[#2d2826] p-8">
					<div class="mb-4 flex items-center gap-4">
						<div class="rounded-lg bg-[#8ec07c]/10 p-3">
							<ChartBar class="h-6 w-6 text-[#8ec07c]" />
						</div>
						<h3 class="font-heading text-2xl font-bold text-white">
							Site Speed & Mobile Performance
						</h3>
					</div>
					<p class="mb-4 text-white/80">
						Slow-loading pages hurt both SEO rankings and conversion rates. Barracuda identifies
						performance bottlenecks, mobile usability issues, and technical problems that impact
						user experience and search visibility.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Page speed issue detection</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Mobile usability validation</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Core Web Vitals optimization opportunities</span>
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
				Audit your e-commerce site in minutes
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
							Crawl Your E-commerce Site
						</h3>
						<p class="mb-4 text-white/80">
							Enter your store URL and let Barracuda crawl your product pages, category pages, and
							key landing pages. Use sitemap seeding to ensure comprehensive coverage of your
							catalog.
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
							Identify Product Page Issues
						</h3>
						<p class="mb-4 text-white/80">
							Barracuda detects missing product schema, duplicate content, thin pages, broken links,
							slow-loading products, and navigation issues that impact your e-commerce SEO
							performance.
						</p>
						<div class="mt-6">
							<OptimizedImage
								src="/mockups/barracuda-issues.png"
								alt="Barracuda SEO Issues Dashboard showing e-commerce product page issues"
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
							Link Google Search Console and Analytics to see which product pages drive traffic and
							conversions. Prioritize fixes on high-value pages that impact revenue.
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
						<h3 class="mb-3 font-heading text-2xl font-bold text-white">Fix & Convert</h3>
						<p class="mb-4 text-white/80">
							Use Barracuda's prioritized issue list to fix technical problems affecting your
							product pages. Track improvements over time with historical crawl comparisons and
							monitor organic traffic growth.
						</p>
						<div class="mt-6 grid grid-cols-1 gap-4 md:grid-cols-2">
							<OptimizedImage
								src="/mockups/barracuda-recommendations.png"
								alt="Barracuda SEO Recommendations for e-commerce pages"
								class="w-full rounded-lg border border-white/10 shadow-lg"
								loading="lazy"
								aspectRatio="16/9"
							/>
							<OptimizedImage
								src="/mockups/barracuda-results.png"
								alt="Barracuda SEO Results showing product page audit results"
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
				<h2 class="font-heading text-4xl font-bold text-white md:text-5xl">E-commerce SEO FAQs</h2>
			</div>

			<div class="space-y-4">
				<div class="overflow-hidden rounded-lg border border-white/10 bg-[#2d2826]">
					<div class="p-6">
						<h3 class="mb-3 font-heading text-xl font-bold text-white">
							How does Barracuda help with e-commerce SEO?
						</h3>
						<p class="text-white/80">
							Barracuda identifies technical SEO issues specific to e-commerce sites (see our <a
								href="/blog/ecommerce-seo-audit"
								class="text-[#8ec07c] underline hover:text-[#a0d28c]">e-commerce SEO audit guide</a
							>), including missing product schema, duplicate content, broken category links,
							slow-loading product pages, and navigation problems. By fixing these issues, you
							improve product page rankings, organic traffic, and conversion rates.
						</p>
					</div>
				</div>

				<div class="overflow-hidden rounded-lg border border-white/10 bg-[#2d2826]">
					<div class="p-6">
						<h3 class="mb-3 font-heading text-xl font-bold text-white">
							Can Barracuda handle large e-commerce catalogs?
						</h3>
						<p class="text-white/80">
							Yes. Barracuda's Pro plan supports crawling 10,000+ pages, making it ideal for large
							e-commerce sites. The fast crawling engine efficiently audits product catalogs, and
							you can use sitemap seeding to ensure comprehensive coverage.
						</p>
					</div>
				</div>

				<div class="overflow-hidden rounded-lg border border-white/10 bg-[#2d2826]">
					<div class="p-6">
						<h3 class="mb-3 font-heading text-xl font-bold text-white">
							What e-commerce SEO issues does Barracuda detect?
						</h3>
						<p class="text-white/80">
							Barracuda detects missing product schema markup, duplicate product content, thin
							product pages, broken category and filter links, orphaned products, slow-loading
							pages, mobile usability issues, and internal linking problems that impact e-commerce
							SEO performance.
						</p>
					</div>
				</div>

				<div class="overflow-hidden rounded-lg border border-white/10 bg-[#2d2826]">
					<div class="p-6">
						<h3 class="mb-3 font-heading text-xl font-bold text-white">
							How can I prioritize fixes for my product pages?
						</h3>
						<p class="text-white/80">
							Connect Google Search Console and Analytics to see which product pages drive traffic
							and conversions. Barracuda's integration helps you prioritize fixes on high-value
							pages that impact revenue, ensuring you focus your SEO efforts where they matter most.
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
				Ready to optimize your e-commerce site for search?
			</h2>
			<p class="mb-10 text-xl text-white/80">
				Start your free 100-page audit and discover the technical issues holding back your product
				page rankings and organic traffic.
			</p>
			<a
				href="https://app.barracudaseo.com"
				class="inline-flex items-center gap-2 rounded-lg bg-[#8ec07c] px-8 py-4 text-lg font-medium text-[#3c3836] transition-colors hover:bg-[#a0d28c]"
				target="_blank"
				rel="noopener noreferrer"
			>
				Start Your E-commerce SEO Audit
				<ArrowRight class="h-5 w-5" />
			</a>
		</div>
	</div>
</section>
