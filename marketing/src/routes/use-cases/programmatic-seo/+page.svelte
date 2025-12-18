<script lang="ts">
	import { Zap, Layers, TrendingUp, ChartBar, CircleCheck, CircleQuestionMark, ArrowRight } from '@lucide/svelte';
</script>

<svelte:head>
	<title>Programmatic SEO Tool - Barracuda SEO</title>
	<meta name="description" content="Scale your SEO content strategy with Barracuda SEO. Audit thousands of programmatic pages, identify technical issues, and optimize at scale for better search rankings." />
</svelte:head>

<!-- Hero Section -->
<section class="py-20 bg-gradient-to-b from-[#3c3836] to-[#2d2826]">
	<div class="container mx-auto px-4">
		<div class="text-center max-w-4xl mx-auto">
			<h1 class="text-5xl md:text-6xl font-heading font-bold mb-6 text-white">
				Audit Programmatic SEO Pages at Scale
			</h1>
			<p class="text-xl md:text-2xl text-white/80 mb-10 max-w-3xl mx-auto">
				Barracuda SEO helps you discover and fix technical SEO issues across thousands of programmatic pages, ensuring your scaled content strategy performs in search results.
			</p>
			<div class="flex flex-col sm:flex-row gap-4 justify-center">
				<a href="https://app.barracudaseo.com" class="bg-[#8ec07c] hover:bg-[#a0d28c] text-[#3c3836] px-8 py-4 rounded-lg font-medium text-lg transition-colors" target="_blank" rel="noopener noreferrer">
					Start Your Programmatic SEO Audit
				</a>
			</div>
		</div>
	</div>
</section>

<!-- Dashboard Preview -->
<section class="py-12 bg-[#2d2826]">
	<div class="container mx-auto px-4">
		<div class="max-w-6xl mx-auto">
			<img src="/mockups/barracuda-dashboard.png" alt="Barracuda SEO Dashboard showing programmatic SEO audit at scale" class="w-full rounded-lg shadow-2xl border border-white/10" />
		</div>
	</div>
</section>

<!-- Features Section -->
<section class="py-20 bg-[#3c3836]">
	<div class="container mx-auto px-4">
		<div class="max-w-6xl mx-auto">
			<h2 class="text-4xl md:text-5xl font-heading font-bold mb-12 text-center text-white">
				Your Complete Programmatic SEO Audit Tool
			</h2>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-8">
				<div class="bg-[#2d2826] rounded-lg p-8 border border-white/10">
					<div class="flex items-center gap-4 mb-4">
						<div class="p-3 bg-[#8ec07c]/10 rounded-lg">
							<Zap class="w-6 h-6 text-[#8ec07c]" />
						</div>
						<h3 class="text-2xl font-heading font-bold text-white">Scale & Performance</h3>
					</div>
					<p class="text-white/80 mb-4">
						Crawl thousands of programmatic pages efficiently. Barracuda's optimized engine handles large-scale audits without breaking, identifying issues across your entire programmatic content library.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Crawl 10,000+ pages with Pro plan</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Fast, concurrent crawling engine</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Export results for bulk analysis</span>
						</li>
					</ul>
				</div>

				<div class="bg-[#2d2826] rounded-lg p-8 border border-white/10">
					<div class="flex items-center gap-4 mb-4">
						<div class="p-3 bg-[#8ec07c]/10 rounded-lg">
							<Layers class="w-6 h-6 text-[#8ec07c]" />
						</div>
						<h3 class="text-2xl font-heading font-bold text-white">Template-Level Issue Detection</h3>
					</div>
					<p class="text-white/80 mb-4">
						Identify patterns across programmatic pages. Barracuda groups issues by URL structure, helping you fix template-level problems that affect hundreds or thousands of pages at once.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Group issues by URL pattern</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Identify duplicate content across templates</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Detect missing or inconsistent meta tags</span>
						</li>
					</ul>
				</div>

				<div class="bg-[#2d2826] rounded-lg p-8 border border-white/10">
					<div class="flex items-center gap-4 mb-4">
						<div class="p-3 bg-[#8ec07c]/10 rounded-lg">
							<TrendingUp class="w-6 h-6 text-[#8ec07c]" />
						</div>
						<h3 class="text-2xl font-heading font-bold text-white">Data-Driven Prioritization</h3>
					</div>
					<p class="text-white/80 mb-4">
						Connect Google Search Console to see which programmatic pages drive traffic and conversions. Prioritize fixes on high-performing pages and identify opportunities on underperforming ones.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>GSC integration for performance data</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Identify high-traffic programmatic pages</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Focus fixes on pages that matter most</span>
						</li>
					</ul>
				</div>

				<div class="bg-[#2d2826] rounded-lg p-8 border border-white/10">
					<div class="flex items-center gap-4 mb-4">
						<div class="p-3 bg-[#8ec07c]/10 rounded-lg">
							<ChartBar class="w-6 h-6 text-[#8ec07c]" />
						</div>
						<h3 class="text-2xl font-heading font-bold text-white">Bulk Issue Resolution</h3>
					</div>
					<p class="text-white/80 mb-4">
						Export crawl results to CSV or JSON for bulk analysis. Use Barracuda's issue grouping to fix template-level problems efficiently, then track improvements with historical comparisons.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>CSV and JSON export for analysis</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Historical crawl comparisons</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Track improvements over time</span>
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
				Audit your programmatic SEO pages in minutes
			</p>

			<div class="space-y-8">
				<div class="flex gap-6">
					<div class="flex-shrink-0">
						<div class="w-12 h-12 bg-[#8ec07c] rounded-full flex items-center justify-center text-[#3c3836] font-bold text-xl">
							1
						</div>
					</div>
					<div class="flex-grow">
						<h3 class="text-2xl font-heading font-bold mb-3 text-white">Crawl Your Programmatic Pages</h3>
						<p class="text-white/80 mb-4">
							Enter your website URL and let Barracuda crawl your programmatic pages. Use sitemap seeding or custom depth limits to focus on your programmatic content structure.
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
						<h3 class="text-2xl font-heading font-bold mb-3 text-white">Identify Template-Level Issues</h3>
						<p class="text-white/80 mb-4">
							Barracuda groups issues by URL pattern, helping you identify template-level problems affecting multiple pages. See duplicate content, missing meta tags, and schema markup issues at scale.
						</p>
						<div class="mt-6">
							<img src="/mockups/barracuda-results.png" alt="Barracuda SEO Results showing programmatic pages grouped by pattern" class="w-full rounded-lg shadow-lg border border-white/10" />
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
							Link Google Search Console to see which programmatic pages drive traffic. Prioritize fixes on high-performing pages and identify opportunities on underperforming ones.
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
						<h3 class="text-2xl font-heading font-bold mb-3 text-white">Fix & Scale</h3>
						<p class="text-white/80 mb-4">
							Export results for bulk analysis, fix template-level issues efficiently, and track improvements over time with historical crawl comparisons.
						</p>
						<div class="mt-6">
							<img src="/mockups/barracuda-link-graph.png" alt="Barracuda SEO Link Graph showing programmatic site structure" class="w-full rounded-lg shadow-lg border border-white/10" />
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
					<CircleQuestionMark class="w-8 h-8 text-[#8ec07c]" />
				</div>
				<h2 class="text-4xl md:text-5xl font-heading font-bold text-white">
					Programmatic SEO FAQs
				</h2>
			</div>

			<div class="space-y-4">
				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<div class="p-6">
						<h3 class="text-xl font-heading font-bold mb-3 text-white">
							How does Barracuda help with programmatic SEO?
						</h3>
						<p class="text-white/80">
							Barracuda audits thousands of programmatic pages efficiently, identifying technical SEO issues at scale. It groups issues by URL pattern to help you fix template-level problems, connects with Google Search Console to prioritize fixes based on performance, and provides bulk export capabilities for analysis.
						</p>
					</div>
				</div>

				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<div class="p-6">
						<h3 class="text-xl font-heading font-bold mb-3 text-white">
							Can Barracuda handle large-scale programmatic sites?
						</h3>
						<p class="text-white/80">
							Yes. Barracuda's Pro plan supports crawling 10,000+ pages, making it ideal for programmatic SEO sites. The fast, concurrent crawling engine handles large-scale audits efficiently, and you can export results for bulk analysis.
						</p>
					</div>
				</div>

				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<div class="p-6">
						<h3 class="text-xl font-heading font-bold mb-3 text-white">
							How does Barracuda identify template-level issues?
						</h3>
						<p class="text-white/80">
							Barracuda groups issues by URL pattern, helping you identify problems that affect multiple programmatic pages. This makes it easy to fix template-level issues like duplicate meta tags, missing schema markup, or inconsistent content structure across your programmatic pages.
						</p>
					</div>
				</div>

				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<div class="p-6">
						<h3 class="text-xl font-heading font-bold mb-3 text-white">
							Can I export programmatic SEO audit results?
						</h3>
						<p class="text-white/80">
							Yes. Barracuda supports CSV and JSON export, making it easy to analyze programmatic page issues in bulk. You can also track improvements over time with historical crawl comparisons.
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
				Ready to audit your programmatic SEO pages?
			</h2>
			<p class="text-xl text-white/80 mb-10">
				Start your free 100-page audit and discover the technical issues affecting your programmatic content strategy.
			</p>
			<a href="https://app.barracudaseo.com" class="inline-flex items-center gap-2 bg-[#8ec07c] hover:bg-[#a0d28c] text-[#3c3836] px-8 py-4 rounded-lg font-medium text-lg transition-colors" target="_blank" rel="noopener noreferrer">
				Start Your Programmatic SEO Audit
				<ArrowRight class="w-5 h-5" />
			</a>
		</div>
	</div>
</section>

