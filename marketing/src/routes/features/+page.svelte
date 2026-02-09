<script lang="ts">
	import { Bolt, ChartBar, Users, Plug, Code, CircleCheck, ChevronRight, Github, Target } from '@lucide/svelte';
	import MetaTags from '../../components/MetaTags.svelte';
	import { getMetaTags, getBreadcrumbSchema } from '$lib/meta';

	const meta = getMetaTags({
		title: 'Features — Comprehensive Technical SEO Auditing & Analysis',
		description: 'Explore Barracuda SEO features: High-speed crawling, interactive link graphs, issue prioritization, and team collaboration tools for agencies and in-house teams.'
	});

	const structuredData = getBreadcrumbSchema([
		{ name: 'Home', url: '/' },
		{ name: 'Features', url: '/features' }
	]);
</script>

<MetaTags config={{ ...meta, structuredData }} />

<!-- Hero Section -->
<section class="py-20 bg-gradient-to-b from-[#3c3836] to-[#2d2826]">
	<div class="container mx-auto px-4">
		<div class="text-center max-w-4xl mx-auto">
			<h1 class="text-5xl md:text-6xl font-heading font-bold mb-6 text-white">
				Crawl Smarter. Fix Faster. Grow Visibility.
			</h1>
			<p class="text-xl md:text-2xl text-white/80 mb-10 max-w-2xl mx-auto">
				Discover and fix the SEO issues holding your website back. Barracuda delivers fast, reliable crawls with actionable insights—all in one sleek web dashboard.
			</p>
			<div class="flex flex-col sm:flex-row gap-4 justify-center">
				<a href="https://app.barracudaseo.com" class="bg-[#8ec07c] hover:bg-[#a0d28c] text-[#3c3836] px-8 py-4 rounded-lg font-medium text-lg transition-colors" target="_blank" rel="noopener noreferrer">
					Start Your Free 100-Page Audit
				</a>
				<a href="#dashboard" class="border-2 border-white/20 hover:border-[#8ec07c] text-white px-8 py-4 rounded-lg font-medium text-lg transition-colors">
					See Barracuda in Action
				</a>
			</div>
		</div>
	</div>
</section>

<!-- Section 1: Intelligent Crawling -->
<section id="crawling" class="py-20 bg-[#3c3836]">
	<div class="container mx-auto px-4">
		<div class="max-w-4xl mx-auto">
			<div class="flex items-center gap-4 mb-6">
				<div class="p-3 bg-[#8ec07c]/10 rounded-lg">
					<Bolt class="w-8 h-8 text-[#8ec07c]" />
				</div>
				<h2 class="text-4xl md:text-5xl font-heading font-bold text-white">Lightning-Fast Crawling for Every Site</h2>
			</div>
			<p class="text-xl text-white/80 mb-8 max-w-3xl">
				Scan your site in minutes, not hours. Barracuda's optimized crawl engine uncovers broken links, missing tags, redirect chains, and slow-loading pages with unmatched precision. Whether you're running a local crawl or managing multiple client sites, Barracuda adapts to your workflow.
			</p>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
				<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
					<h3 class="text-lg font-heading font-bold mb-3 text-[#8ec07c]">Scale & Performance</h3>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Crawl up to <strong class="text-white">10,000+ pages</strong> with the Pro plan</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span><strong class="text-white">Respect robots.txt</strong> automatically</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Supports <strong class="text-white">sitemap seeding</strong> and custom depth limits</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Export results to <strong class="text-white">CSV, JSON, or Supabase Cloud</strong></span>
						</li>
					</ul>
				</div>
				<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
					<h3 class="text-lg font-heading font-bold mb-3 text-[#8ec07c]">What We Detect</h3>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Broken links and redirect chains</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Missing or duplicate meta tags</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>Slow-loading pages</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span>SEO best practice violations</span>
						</li>
					</ul>
				</div>
			</div>
			<div class="text-center">
				<a href="https://app.barracudaseo.com" class="inline-flex items-center gap-2 bg-[#8ec07c] hover:bg-[#a0d28c] text-[#3c3836] px-6 py-3 rounded-lg font-medium transition-colors" target="_blank" rel="noopener noreferrer">
					Start Crawling Free
					<ChevronRight class="w-4 h-4" />
				</a>
			</div>
		</div>
	</div>
</section>

<!-- Section 2: Insightful Dashboard -->
<section id="dashboard" class="py-20 bg-[#2d2826]">
	<div class="container mx-auto px-4">
		<div class="max-w-4xl mx-auto">
			<div class="flex items-center gap-4 mb-6">
				<div class="p-3 bg-[#8ec07c]/10 rounded-lg">
					<ChartBar class="w-8 h-8 text-[#8ec07c]" />
				</div>
				<h2 class="text-4xl md:text-5xl font-heading font-bold text-white">Turn Data into Decisions</h2>
			</div>
			<p class="text-xl text-white/80 mb-8 max-w-3xl">
				No more spreadsheets or endless tabs. The Barracuda dashboard organizes crawl results into a clean, visual interface that surfaces the most critical issues first.
			</p>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
				<div class="bg-[#3c3836] rounded-lg p-6 border border-white/10">
					<h3 class="text-lg font-heading font-bold mb-3 text-[#8ec07c]">Priority Scoring</h3>
					<p class="text-white/70">
						Automatically ranks issues by impact, so you know exactly what to fix first.
					</p>
				</div>
				<div class="bg-[#3c3836] rounded-lg p-6 border border-white/10">
					<h3 class="text-lg font-heading font-bold mb-3 text-[#8ec07c]">Interactive Link Graph</h3>
					<p class="text-white/70">
						Visualize site structure and discover orphaned pages at a glance.
					</p>
				</div>
				<div class="bg-[#3c3836] rounded-lg p-6 border border-white/10">
					<h3 class="text-lg font-heading font-bold mb-3 text-[#8ec07c]">Issue Modals</h3>
					<p class="text-white/70">
						Get details, examples, and best-practice guidance for every problem.
					</p>
				</div>
				<div class="bg-[#3c3836] rounded-lg p-6 border border-[#8ec07c]/30 relative">
					<div class="absolute -top-2 -right-2 bg-[#8ec07c] text-[#3c3836] px-2 py-1 rounded text-xs font-bold">Pro</div>
					<h3 class="text-lg font-heading font-bold mb-3 text-[#8ec07c]">AI-Powered Recommendations</h3>
					<p class="text-white/70">
						See how to fix issues, not just what's broken. (Pro feature)
					</p>
				</div>
			</div>
			<div class="text-center">
				<a href="https://app.barracudaseo.com" class="inline-flex items-center gap-2 bg-[#8ec07c] hover:bg-[#a0d28c] text-[#3c3836] px-6 py-3 rounded-lg font-medium transition-colors" target="_blank" rel="noopener noreferrer">
					Try the Dashboard — Free for 100 Pages
					<ChevronRight class="w-4 h-4" />
				</a>
			</div>
		</div>
	</div>
</section>

<!-- Section 3: Collaboration Made Easy -->
<section id="collaboration" class="py-20 bg-[#3c3836]">
	<div class="container mx-auto px-4">
		<div class="max-w-4xl mx-auto">
			<div class="flex items-center gap-4 mb-6">
				<div class="p-3 bg-[#8ec07c]/10 rounded-lg">
					<Users class="w-8 h-8 text-[#8ec07c]" />
				</div>
				<h2 class="text-4xl md:text-5xl font-heading font-bold text-white">Built for Teams and Agencies</h2>
			</div>
			<p class="text-xl text-white/80 mb-8 max-w-3xl">
				Work together seamlessly. Invite your teammates, assign roles, and share crawl results in real-time. Whether you're an agency managing multiple client sites or an in-house SEO team, Barracuda makes collaboration effortless.
			</p>
			<div class="bg-[#2d2826] rounded-lg p-8 border border-white/10 mb-8">
				<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
					<div>
						<h3 class="text-lg font-heading font-bold mb-2 text-[#8ec07c]">Owner</h3>
						<p class="text-white/70 text-sm">
							Full permissions including billing, invitations, and project deletion.
						</p>
					</div>
					<div>
						<h3 class="text-lg font-heading font-bold mb-2 text-[#8ec07c]">Editor</h3>
						<p class="text-white/70 text-sm">
							Can trigger crawls, manage issues, and view recommendations.
						</p>
					</div>
					<div>
						<h3 class="text-lg font-heading font-bold mb-2 text-[#8ec07c]">Viewer</h3>
						<p class="text-white/70 text-sm">
							Read-only access to dashboard and reports.
						</p>
					</div>
				</div>
				<div class="mt-6 pt-6 border-t border-white/10">
					<p class="text-white/70 text-center">
						Team add-ons available starting at <strong class="text-white">$5 per user</strong>
					</p>
				</div>
			</div>
			<div class="text-center">
				<a href="/pricing" class="inline-flex items-center gap-2 bg-[#8ec07c] hover:bg-[#a0d28c] text-[#3c3836] px-6 py-3 rounded-lg font-medium transition-colors">
					Upgrade to a Team Account
					<ChevronRight class="w-4 h-4" />
				</a>
			</div>
		</div>
	</div>
</section>

<!-- Section 4: Advanced Integrations -->
<section id="integrations" class="py-20 bg-[#2d2826]">
	<div class="container mx-auto px-4">
		<div class="max-w-4xl mx-auto">
			<div class="flex items-center gap-4 mb-6">
				<div class="p-3 bg-[#8ec07c]/10 rounded-lg">
					<Plug class="w-8 h-8 text-[#8ec07c]" />
				</div>
				<h2 class="text-4xl md:text-5xl font-heading font-bold text-white">Connect Your Favorite SEO Tools</h2>
			</div>
			<p class="text-xl text-white/80 mb-8 max-w-3xl">
				Take your audits further with integrations that give you context beyond the crawl. Pro users can link Google Search Console, Google Analytics, and Microsoft Clarity to enrich recommendations and track performance directly from the dashboard.
			</p>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
				<div class="bg-[#3c3836] rounded-lg p-6 border border-white/10">
					<h3 class="text-lg font-heading font-bold mb-2 text-[#8ec07c]">Google Search Console</h3>
					<p class="text-white/70 text-sm">
						Import impressions, clicks, and CTR data to prioritize fixes that impact visibility.
					</p>
				</div>
				<div class="bg-[#3c3836] rounded-lg p-6 border border-white/10">
					<h3 class="text-lg font-heading font-bold mb-2 text-[#8ec07c]">Google Analytics</h3>
					<p class="text-white/70 text-sm">
						Measure engagement and conversions by page to focus on high-value fixes.
					</p>
				</div>
				<div class="bg-[#3c3836] rounded-lg p-6 border border-white/10">
					<h3 class="text-lg font-heading font-bold mb-2 text-[#8ec07c]">Microsoft Clarity</h3>
					<p class="text-white/70 text-sm">
						Visualize user behavior heatmaps alongside technical issues.
					</p>
				</div>
				<div class="bg-[#3c3836] rounded-lg p-6 border border-white/10">
					<h3 class="text-lg font-heading font-bold mb-2 text-[#8ec07c]">Slack</h3>
					<p class="text-white/70 text-sm">
						Get instant crawl reports in your workspace. <span class="text-white/50">(Coming Soon)</span>
					</p>
				</div>
			</div>
			<div class="text-center">
				<a href="/pricing" class="inline-flex items-center gap-2 bg-[#8ec07c] hover:bg-[#a0d28c] text-[#3c3836] px-6 py-3 rounded-lg font-medium transition-colors">
					Unlock Integrations — Go Pro
					<ChevronRight class="w-4 h-4" />
				</a>
			</div>
		</div>
	</div>
</section>

<!-- Section 5: CLI for Developers -->
<section id="cli" class="py-20 bg-[#3c3836]">
	<div class="container mx-auto px-4">
		<div class="max-w-4xl mx-auto">
			<div class="flex items-center gap-4 mb-6">
				<div class="p-3 bg-[#8ec07c]/10 rounded-lg">
					<Code class="w-8 h-8 text-[#8ec07c]" />
				</div>
				<h2 class="text-4xl md:text-5xl font-heading font-bold text-white">Automate Audits with the Barracuda CLI</h2>
			</div>
			<p class="text-xl text-white/80 mb-8 max-w-3xl">
				Developers, we've got you covered. The same powerful engine that powers the web app is also available as an open-source CLI. Integrate Barracuda into your CI/CD pipelines, trigger crawls via scripts, and export clean data for analysis.
			</p>
			<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10 mb-8">
				<ul class="space-y-3 text-white/70">
					<li class="flex items-start gap-3">
						<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
						<span>Run local or remote crawls</span>
					</li>
					<li class="flex items-start gap-3">
						<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
						<span>JSON and CSV export support</span>
					</li>
					<li class="flex items-start gap-3">
						<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
						<span>Upload directly to the Barracuda Cloud workspace</span>
					</li>
					<li class="flex items-start gap-3">
						<CircleCheck class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
						<span>Open-source under the MIT License</span>
					</li>
				</ul>
			</div>
			<div class="text-center">
				<a href="https://github.com/dillonlara115/barracudaseo" class="inline-flex items-center gap-2 border-2 border-white/20 hover:border-[#8ec07c] text-white px-6 py-3 rounded-lg font-medium transition-colors" target="_blank" rel="noopener noreferrer">
					<Github class="w-5 h-5" />
					Download the CLI on GitHub
				</a>
			</div>
		</div>
	</div>
</section>

<!-- Section 6: Why Barracuda -->
<section class="py-20 bg-[#2d2826]">
	<div class="container mx-auto px-4">
		<div class="max-w-4xl mx-auto">
			<h2 class="text-4xl md:text-5xl font-heading font-bold mb-12 text-center text-white">Why Agencies Choose Barracuda</h2>
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6 mb-12">
				<div class="text-center">
					<div class="p-3 bg-[#8ec07c]/10 rounded-lg inline-block mb-3">
						<Bolt class="w-6 h-6 text-[#8ec07c]" />
					</div>
					<h3 class="font-heading font-bold mb-2 text-white">Fast</h3>
					<p class="text-white/70 text-sm">Multi-threaded crawling engine built in Go</p>
				</div>
				<div class="text-center">
					<div class="p-3 bg-[#8ec07c]/10 rounded-lg inline-block mb-3">
						<Plug class="w-6 h-6 text-[#8ec07c]" />
					</div>
					<h3 class="font-heading font-bold mb-2 text-white">Flexible</h3>
					<p class="text-white/70 text-sm">Use via web or command line</p>
				</div>
				<div class="text-center">
					<div class="p-3 bg-[#8ec07c]/10 rounded-lg inline-block mb-3">
						<CircleCheck class="w-6 h-6 text-[#8ec07c]" />
					</div>
					<h3 class="font-heading font-bold mb-2 text-white">Transparent</h3>
					<p class="text-white/70 text-sm">Open-source foundation and honest pricing</p>
				</div>
				<div class="text-center">
					<div class="p-3 bg-[#8ec07c]/10 rounded-lg inline-block mb-3">
						<Target class="w-6 h-6 text-[#8ec07c]" />
					</div>
					<h3 class="font-heading font-bold mb-2 text-white">Actionable</h3>
					<p class="text-white/70 text-sm">Prioritized, fix-ready insights</p>
				</div>
				<div class="text-center">
					<div class="p-3 bg-[#8ec07c]/10 rounded-lg inline-block mb-3">
						<Users class="w-6 h-6 text-[#8ec07c]" />
					</div>
					<h3 class="font-heading font-bold mb-2 text-white">Collaborative</h3>
					<p class="text-white/70 text-sm">Built for teams, not just individuals</p>
				</div>
			</div>
			<div class="bg-[#3c3836] rounded-lg p-8 border border-white/10 mb-8">
				<blockquote class="text-xl text-white/90 italic mb-4">
					"Barracuda makes technical SEO feel effortless. It's fast, accurate, and the dashboard is genuinely useful."
				</blockquote>
				<p class="text-white/70 text-sm">
					— Early Access User, SEO Agency Owner
				</p>
			</div>
			<div class="text-center">
				<a href="https://app.barracudaseo.com" class="inline-flex items-center gap-2 bg-[#8ec07c] hover:bg-[#a0d28c] text-[#3c3836] px-8 py-4 rounded-lg font-medium text-lg transition-colors" target="_blank" rel="noopener noreferrer">
					Start Your Free 100-Page Crawl
					<ChevronRight class="w-5 h-5" />
				</a>
			</div>
		</div>
	</div>
</section>

<!-- Section 7: Pricing Preview -->
<section class="py-20 bg-[#3c3836]">
	<div class="container mx-auto px-4">
		<div class="max-w-4xl mx-auto">
			<h2 class="text-4xl md:text-5xl font-heading font-bold mb-12 text-center text-white">Simple, Fair Pricing</h2>
			<div class="bg-[#2d2826] rounded-lg overflow-hidden border border-white/10">
				<div class="overflow-x-auto">
					<table class="w-full">
						<thead class="bg-[#3c3836]">
							<tr>
								<th class="px-6 py-4 text-left text-sm font-heading font-bold text-white">Plan</th>
								<th class="px-6 py-4 text-left text-sm font-heading font-bold text-white">Price</th>
								<th class="px-6 py-4 text-left text-sm font-heading font-bold text-white">Crawl Limit</th>
								<th class="px-6 py-4 text-left text-sm font-heading font-bold text-white">Users</th>
								<th class="px-6 py-4 text-left text-sm font-heading font-bold text-white">Integrations</th>
								<th class="px-6 py-4 text-left text-sm font-heading font-bold text-white">Recommendations</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-white/10">
							<tr>
								<td class="px-6 py-4 text-white font-medium">Free</td>
								<td class="px-6 py-4 text-white/70">$0</td>
								<td class="px-6 py-4 text-white/70">100 pages</td>
								<td class="px-6 py-4 text-white/70">1</td>
								<td class="px-6 py-4 text-white/50">—</td>
								<td class="px-6 py-4 text-white/50">—</td>
							</tr>
							<tr class="bg-[#8ec07c]/10">
								<td class="px-6 py-4 text-white font-medium">Pro</td>
								<td class="px-6 py-4 text-white">$29/mo</td>
								<td class="px-6 py-4 text-white">10,000+ pages</td>
								<td class="px-6 py-4 text-white">1</td>
								<td class="px-6 py-4 text-[#8ec07c]">✅</td>
								<td class="px-6 py-4 text-[#8ec07c]">✅</td>
							</tr>
							<tr>
								<td class="px-6 py-4 text-white font-medium">Team</td>
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
			<div class="text-center mt-8">
				<a href="/pricing" class="inline-flex items-center gap-2 bg-[#8ec07c] hover:bg-[#a0d28c] text-[#3c3836] px-6 py-3 rounded-lg font-medium transition-colors">
					Explore Plans & Upgrade
					<ChevronRight class="w-4 h-4" />
				</a>
			</div>
		</div>
	</div>
</section>

<!-- Section 8: Final CTA -->
<section class="py-20 bg-gradient-to-b from-[#2d2826] to-[#3c3836]">
	<div class="container mx-auto px-4">
		<div class="max-w-3xl mx-auto text-center">
			<h2 class="text-4xl md:text-5xl font-heading font-bold mb-6 text-white">
				Your site deserves better SEO visibility.
			</h2>
			<p class="text-xl text-white/80 mb-10">
				Join hundreds of Agencies and developers who use Barracuda to uncover issues, prioritize fixes, and grow their traffic.
			</p>
			<div class="flex flex-col sm:flex-row gap-4 justify-center">
				<a href="https://app.barracudaseo.com" class="bg-[#8ec07c] hover:bg-[#a0d28c] text-[#3c3836] px-8 py-4 rounded-lg font-medium text-lg transition-colors" target="_blank" rel="noopener noreferrer">
					Start Free — No Credit Card Required
				</a>
				<a href="#dashboard" class="border-2 border-white/20 hover:border-[#8ec07c] text-white px-8 py-4 rounded-lg font-medium text-lg transition-colors">
					View Demo
				</a>
			</div>
		</div>
	</div>
</section>

