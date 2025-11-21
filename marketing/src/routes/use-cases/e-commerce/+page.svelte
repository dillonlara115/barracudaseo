<script lang="ts">
	import { ShoppingCart, Package, TrendingUp, BarChart3, CheckCircle2, HelpCircle, ArrowRight } from 'lucide-svelte';
</script>

<svelte:head>
	<title>E-commerce SEO Auditing Tool - Barracuda SEO</title>
	<meta name="description" content="Optimize your e-commerce website for search with Barracuda SEO. Discover technical issues, improve product page performance, and boost organic traffic and conversions." />
</svelte:head>

<!-- Hero Section -->
<section class="py-20 bg-gradient-to-b from-[#3c3836] to-[#2d2826]">
	<div class="container mx-auto px-4">
		<div class="text-center max-w-4xl mx-auto">
			<h1 class="text-5xl md:text-6xl font-heading font-bold mb-6 text-white">
				Optimize E-commerce Sites for Search Success
			</h1>
			<p class="text-xl md:text-2xl text-white/80 mb-10 max-w-3xl mx-auto">
				Barracuda SEO helps e-commerce businesses discover and fix technical SEO issues that impact product page rankings, organic traffic, and conversion rates.
			</p>
			<div class="flex flex-col sm:flex-row gap-4 justify-center">
				<a href="https://app.barracudaseo.com" class="bg-[#8ec07c] hover:bg-[#a0d28c] text-[#3c3836] px-8 py-4 rounded-lg font-medium text-lg transition-colors" target="_blank" rel="noopener noreferrer">
					Start Your E-commerce SEO Audit
				</a>
			</div>
		</div>
	</div>
</section>

<!-- Dashboard Preview -->
<section class="py-12 bg-[#2d2826]">
	<div class="container mx-auto px-4">
		<div class="max-w-6xl mx-auto">
			<img src="/mockups/barracuda-dashboard.png" alt="Barracuda SEO Dashboard showing e-commerce SEO audit results" class="w-full rounded-lg shadow-2xl border border-white/10" />
		</div>
	</div>
</section>

<!-- Features Section -->
<section class="py-20 bg-[#3c3836]">
	<div class="container mx-auto px-4">
		<div class="max-w-6xl mx-auto">
			<h2 class="text-4xl md:text-5xl font-heading font-bold mb-12 text-center text-white">
				Your Complete E-commerce SEO Audit Tool
			</h2>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-8">
				<div class="bg-[#2d2826] rounded-lg p-8 border border-white/10">
					<div class="flex items-center gap-4 mb-4">
						<div class="p-3 bg-[#8ec07c]/10 rounded-lg">
							<ShoppingCart class="w-6 h-6 text-[#8ec07c]" />
						</div>
						<h3 class="text-2xl font-heading font-bold text-white">Product Page Optimization</h3>
					</div>
					<p class="text-white/80 mb-4">
						Ensure your product pages are properly optimized for search. Barracuda identifies missing product schema, duplicate content, thin content issues, and optimization opportunities across your catalog.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Product schema markup validation</span>
						</li>
						<li class="flex items-start gap-2">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Duplicate product content detection</span>
						</li>
						<li class="flex items-start gap-2">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Missing meta descriptions and titles</span>
						</li>
					</ul>
				</div>

				<div class="bg-[#2d2826] rounded-lg p-8 border border-white/10">
					<div class="flex items-center gap-4 mb-4">
						<div class="p-3 bg-[#8ec07c]/10 rounded-lg">
							<Package class="w-6 h-6 text-[#8ec07c]" />
						</div>
						<h3 class="text-2xl font-heading font-bold text-white">Category & Navigation Structure</h3>
					</div>
					<p class="text-white/80 mb-4">
						Optimize your site architecture for both users and search engines. Barracuda identifies broken category links, orphaned product pages, navigation issues, and internal linking problems that hurt SEO performance.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Broken category and filter links</span>
						</li>
						<li class="flex items-start gap-2">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Orphaned product pages</span>
						</li>
						<li class="flex items-start gap-2">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Internal linking structure analysis</span>
						</li>
					</ul>
				</div>

				<div class="bg-[#2d2826] rounded-lg p-8 border border-white/10">
					<div class="flex items-center gap-4 mb-4">
						<div class="p-3 bg-[#8ec07c]/10 rounded-lg">
							<TrendingUp class="w-6 h-6 text-[#8ec07c]" />
						</div>
						<h3 class="text-2xl font-heading font-bold text-white">Performance-Driven Prioritization</h3>
					</div>
					<p class="text-white/80 mb-4">
						Connect Google Search Console and Analytics to see which product pages drive traffic and conversions. Prioritize fixes on high-value pages and identify opportunities on underperforming products.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>GSC and GA4 integration</span>
						</li>
						<li class="flex items-start gap-2">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Identify high-traffic product pages</span>
						</li>
						<li class="flex items-start gap-2">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Conversion-focused issue prioritization</span>
						</li>
					</ul>
				</div>

				<div class="bg-[#2d2826] rounded-lg p-8 border border-white/10">
					<div class="flex items-center gap-4 mb-4">
						<div class="p-3 bg-[#8ec07c]/10 rounded-lg">
							<BarChart3 class="w-6 h-6 text-[#8ec07c]" />
						</div>
						<h3 class="text-2xl font-heading font-bold text-white">Site Speed & Mobile Performance</h3>
					</div>
					<p class="text-white/80 mb-4">
						Slow-loading pages hurt both SEO rankings and conversion rates. Barracuda identifies performance bottlenecks, mobile usability issues, and technical problems that impact user experience and search visibility.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Page speed issue detection</span>
						</li>
						<li class="flex items-start gap-2">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Mobile usability validation</span>
						</li>
						<li class="flex items-start gap-2">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Core Web Vitals optimization opportunities</span>
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
				Audit your e-commerce site in minutes
			</p>

			<div class="space-y-8">
				<div class="flex gap-6">
					<div class="flex-shrink-0">
						<div class="w-12 h-12 bg-[#8ec07c] rounded-full flex items-center justify-center text-[#3c3836] font-bold text-xl">
							1
						</div>
					</div>
					<div class="flex-grow">
						<h3 class="text-2xl font-heading font-bold mb-3 text-white">Crawl Your E-commerce Site</h3>
						<p class="text-white/80 mb-4">
							Enter your store URL and let Barracuda crawl your product pages, category pages, and key landing pages. Use sitemap seeding to ensure comprehensive coverage of your catalog.
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
						<h3 class="text-2xl font-heading font-bold mb-3 text-white">Identify Product Page Issues</h3>
						<p class="text-white/80 mb-4">
							Barracuda detects missing product schema, duplicate content, thin pages, broken links, slow-loading products, and navigation issues that impact your e-commerce SEO performance.
						</p>
						<div class="mt-6">
							<img src="/mockups/barracuda-issues.png" alt="Barracuda SEO Issues Dashboard showing e-commerce product page issues" class="w-full rounded-lg shadow-lg border border-white/10" />
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
						<h3 class="text-2xl font-heading font-bold mb-3 text-white">Connect Performance Data</h3>
						<p class="text-white/80 mb-4">
							Link Google Search Console and Analytics to see which product pages drive traffic and conversions. Prioritize fixes on high-value pages that impact revenue.
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
						<h3 class="text-2xl font-heading font-bold mb-3 text-white">Fix & Convert</h3>
						<p class="text-white/80 mb-4">
							Use Barracuda's prioritized issue list to fix technical problems affecting your product pages. Track improvements over time with historical crawl comparisons and monitor organic traffic growth.
						</p>
						<div class="mt-6 grid grid-cols-1 md:grid-cols-2 gap-4">
							<img src="/mockups/barracuda-recommendations.png" alt="Barracuda SEO Recommendations for e-commerce pages" class="w-full rounded-lg shadow-lg border border-white/10" />
							<img src="/mockups/barracuda-results.png" alt="Barracuda SEO Results showing product page audit results" class="w-full rounded-lg shadow-lg border border-white/10" />
						</div>
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
					<HelpCircle class="w-8 h-8 text-[#8ec07c]" />
				</div>
				<h2 class="text-4xl md:text-5xl font-heading font-bold text-white">
					E-commerce SEO FAQs
				</h2>
			</div>

			<div class="space-y-4">
				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<div class="p-6">
						<h3 class="text-xl font-heading font-bold mb-3 text-white">
							How does Barracuda help with e-commerce SEO?
						</h3>
						<p class="text-white/80">
							Barracuda identifies technical SEO issues specific to e-commerce sites, including missing product schema, duplicate content, broken category links, slow-loading product pages, and navigation problems. By fixing these issues, you improve product page rankings, organic traffic, and conversion rates.
						</p>
					</div>
				</div>

				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<div class="p-6">
						<h3 class="text-xl font-heading font-bold mb-3 text-white">
							Can Barracuda handle large e-commerce catalogs?
						</h3>
						<p class="text-white/80">
							Yes. Barracuda's Pro plan supports crawling 10,000+ pages, making it ideal for large e-commerce sites. The fast crawling engine efficiently audits product catalogs, and you can use sitemap seeding to ensure comprehensive coverage.
						</p>
					</div>
				</div>

				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<div class="p-6">
						<h3 class="text-xl font-heading font-bold mb-3 text-white">
							What e-commerce SEO issues does Barracuda detect?
						</h3>
						<p class="text-white/80">
							Barracuda detects missing product schema markup, duplicate product content, thin product pages, broken category and filter links, orphaned products, slow-loading pages, mobile usability issues, and internal linking problems that impact e-commerce SEO performance.
						</p>
					</div>
				</div>

				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<div class="p-6">
						<h3 class="text-xl font-heading font-bold mb-3 text-white">
							How can I prioritize fixes for my product pages?
						</h3>
						<p class="text-white/80">
							Connect Google Search Console and Analytics to see which product pages drive traffic and conversions. Barracuda's integration helps you prioritize fixes on high-value pages that impact revenue, ensuring you focus your SEO efforts where they matter most.
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
				Ready to optimize your e-commerce site for search?
			</h2>
			<p class="text-xl text-white/80 mb-10">
				Start your free 100-page audit and discover the technical issues holding back your product page rankings and organic traffic.
			</p>
			<a href="https://app.barracudaseo.com" class="inline-flex items-center gap-2 bg-[#8ec07c] hover:bg-[#a0d28c] text-[#3c3836] px-8 py-4 rounded-lg font-medium text-lg transition-colors" target="_blank" rel="noopener noreferrer">
				Start Your E-commerce SEO Audit
				<ArrowRight class="w-5 h-5" />
			</a>
		</div>
	</div>
</section>

