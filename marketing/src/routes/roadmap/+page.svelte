<script lang="ts">
	// Import the correct icon components from lucide-svelte
	import CheckCircle2 from 'lucide-svelte/CheckCircle2.svelte';
	import Clock from 'lucide-svelte/icons/Clock.svelte';
	import Rocket from 'lucide-svelte/icons/Rocket.svelte';
	import Target from 'lucide-svelte/icons/Target.svelte';
	import Sparkles from 'lucide-svelte/icons/Sparkles.svelte';
	import Zap from 'lucide-svelte/icons/Zap.svelte';
	import Github from 'lucide-svelte/icons/Github.svelte';
	import MessageSquare from 'lucide-svelte/icons/MessageSquare.svelte';
	import ExternalLink from 'lucide-svelte/icons/ExternalLink.svelte';
	'@lucide/svelte/icons/circle-alert';
</script>

<svelte:head>
	<title>Product Roadmap - Barracuda SEO</title>
	<meta name="description" content="See what's coming next for Barracuda SEO. Our roadmap includes AI-powered insights, scheduled crawls, rank tracking, and more." />
</svelte:head>

<!-- Hero Section -->
<section class="py-20 bg-gradient-to-b from-[#3c3836] to-[#2d2826]">
	<div class="container mx-auto px-4">
		<div class="text-center max-w-4xl mx-auto">
			<h1 class="text-5xl md:text-6xl font-heading font-bold mb-6 text-white">
				Barracuda Product Roadmap
			</h1>
			<p class="text-xl md:text-2xl text-white/80 mb-10 max-w-3xl mx-auto">
				We're building Barracuda to be the fastest, most intelligent cloud-native SEO crawler and insights platform. Our roadmap is organized around three core pillars: <strong class="text-[#8ec07c]">Performance</strong>, <strong class="text-[#8ec07c]">Collaboration</strong>, and <strong class="text-[#8ec07c]">Intelligence (AI)</strong>.
			</p>
		</div>
	</div>
</section>

<!-- Section 1: Completed Milestones -->
<section class="py-20 bg-[#3c3836]">
	<div class="container mx-auto px-4">
		<div class="max-w-6xl mx-auto">
			<div class="flex items-center gap-4 mb-8">
				<div class="p-3 bg-[#8ec07c]/10 rounded-lg">
					<CheckCircle2 class="w-8 h-8 text-[#8ec07c]" />
				</div>
				<h2 class="text-4xl md:text-5xl font-heading font-bold text-white">Completed Milestones</h2>
			</div>

			<!-- 1.1 Core Infrastructure -->
			<div class="mb-12">
				<h3 class="text-2xl md:text-3xl font-heading font-bold mb-6 text-[#8ec07c]">1.1 Core Infrastructure</h3>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Supabase backend with RLS</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Cloud Run API server</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">SvelteKit app deployed to Vercel</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Production environment variables & secrets configured</span>
						</div>
					</div>
				</div>
			</div>

			<!-- 1.2 Core Features -->
			<div class="mb-12">
				<h3 class="text-2xl md:text-3xl font-heading font-bold mb-6 text-[#8ec07c]">1.2 Core Features</h3>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Website crawler</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Google Search Console integration</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Project workspace system</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Team accounts + invites</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Billing + Stripe integration</span>
						</div>
					</div>
				</div>
			</div>

			<!-- 1.3 CLI Functionality -->
			<div class="mb-12">
				<h3 class="text-2xl md:text-3xl font-heading font-bold mb-6 text-[#8ec07c]">1.3 CLI Functionality</h3>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Local crawling</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Exporting results (CSV, JSON)</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Optional cloud upload</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Embedded dashboard (serve command)</span>
						</div>
					</div>
				</div>
			</div>

			<!-- 1.4 Cloud Dashboard -->
			<div class="mb-12">
				<h3 class="text-2xl md:text-3xl font-heading font-bold mb-6 text-[#8ec07c]">1.4 Cloud Dashboard</h3>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Project overview pages</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Crawl history</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Issue management interface</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Page-level detail views</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Link graph visualization</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">GSC overlays inside UI</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">AI insights with GSC data integration</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Public-facing client reports</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">HTML report generator</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">AI executive summaries</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Dedicated Crawls & Audits management page</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Improved navigation and UX with page descriptions</span>
						</div>
					</div>
				</div>
			</div>

			<!-- 1.6 Rank Tracking -->
			<div class="mb-12">
				<h3 class="text-2xl md:text-3xl font-heading font-bold mb-6 text-[#8ec07c]">1.6 Rank Tracking</h3>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">SERP data integration for rank tracking</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Keyword rank tracking with historical snapshots</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Ranked keywords discovery (reverse lookup)</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Impact-First View (pages ranking with issues)</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Scheduled keyword checks</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Position trend tracking and visualization</span>
						</div>
					</div>
				</div>
			</div>

			<!-- 1.5 Marketing Site -->
			<div class="mb-12">
				<h3 class="text-2xl md:text-3xl font-heading font-bold mb-6 text-[#8ec07c]">1.5 Marketing Site</h3>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Updated feature descriptions</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">New visual direction to match the app</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">SEO-optimized landing pages</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Pricing, features, and signup funnels</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<CheckCircle2 class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Product roadmap page</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</section>

<!-- Section 2: In Progress -->
<section class="py-20 bg-[#2d2826]">
	<div class="container mx-auto px-4">
		<div class="max-w-6xl mx-auto">
			<div class="flex items-center gap-4 mb-8">
				<div class="p-3 bg-[#d79921]/10 rounded-lg">
					<Clock class="w-8 h-8 text-[#d79921]" />
				</div>
				<h2 class="text-4xl md:text-5xl font-heading font-bold text-white">In Progress</h2>
			</div>

			<!-- Note: Marketing Site Enhancements moved to Completed -->
			<div class="text-center py-8 text-white/60">
				<p>More features coming soon...</p>
			</div>
		</div>
	</div>
</section>

<!-- Section 3: Upcoming Features -->
<section class="py-20 bg-[#3c3836]">
	<div class="container mx-auto px-4">
		<div class="max-w-6xl mx-auto">
			<div class="flex items-center gap-4 mb-8">
				<div class="p-3 bg-[#458588]/10 rounded-lg">
					<Rocket class="w-8 h-8 text-[#458588]" />
				</div>
				<h2 class="text-4xl md:text-5xl font-heading font-bold text-white">Upcoming Features</h2>
				<span class="text-lg text-white/70">(Estimated next 1–3 months)</span>
			</div>

			<!-- 3.1 AI Integration -->
			<div class="mb-12">
				<h3 class="text-2xl md:text-3xl font-heading font-bold mb-6 text-[#8ec07c]">3.1 AI Integration</h3>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div class="bg-[#2d2826] rounded-lg p-6 border border-[#8ec07c]/30">
						<div class="flex items-start gap-3">
							<Sparkles class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">AI recommendations for issues</span>
							<span class="ml-auto text-xs text-[#8ec07c] font-semibold">✓ COMPLETE</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-[#8ec07c]/30">
						<div class="flex items-start gap-3">
							<Sparkles class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">AI crawl summary reports</span>
							<span class="ml-auto text-xs text-[#8ec07c] font-semibold">✓ COMPLETE</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Sparkles class="w-5 h-5 text-[#458588] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Page-level rewrite suggestions</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Sparkles class="w-5 h-5 text-[#458588] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">AI-powered priority scoring</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-[#8ec07c]/30">
						<div class="flex items-start gap-3">
							<Sparkles class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">AI insights combining crawl + GSC</span>
							<span class="ml-auto text-xs text-[#8ec07c] font-semibold">✓ COMPLETE</span>
						</div>
					</div>
				</div>
			</div>

			<!-- 3.2 Additional Integrations -->
			<div class="mb-12">
				<h3 class="text-2xl md:text-3xl font-heading font-bold mb-6 text-[#8ec07c]">3.2 Additional Integrations</h3>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Zap class="w-5 h-5 text-[#458588] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Google Analytics 4 (GA4)</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Zap class="w-5 h-5 text-[#458588] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Google Drive / Google Sheets export</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Zap class="w-5 h-5 text-[#458588] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Slack alerts</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Zap class="w-5 h-5 text-[#458588] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Microsoft Clarity integration</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Zap class="w-5 h-5 text-[#458588] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Zapier / Make.com workflows</span>
						</div>
					</div>
				</div>
			</div>

			<!-- 3.3 CLI → Cloud Sync Improvements -->
			<div class="mb-12">
				<h3 class="text-2xl md:text-3xl font-heading font-bold mb-6 text-[#8ec07c]">3.3 CLI → Cloud Sync Improvements</h3>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Zap class="w-5 h-5 text-[#458588] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">barracuda auth login</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Zap class="w-5 h-5 text-[#458588] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">barracuda projects list</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Zap class="w-5 h-5 text-[#458588] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">--cloud ingestion polish</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Zap class="w-5 h-5 text-[#458588] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Auto-associate crawls to default project</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</section>

<!-- Section 4: Mid-Term Goals -->
<section class="py-20 bg-[#2d2826]">
	<div class="container mx-auto px-4">
		<div class="max-w-6xl mx-auto">
			<div class="flex items-center gap-4 mb-8">
				<div class="p-3 bg-[#458588]/10 rounded-lg">
					<Target class="w-8 h-8 text-[#458588]" />
				</div>
				<h2 class="text-4xl md:text-5xl font-heading font-bold text-white">Mid-Term Goals</h2>
				<span class="text-lg text-white/70">(3–6 month window)</span>
			</div>

			<!-- 4.1 Scheduled Crawls -->
			<div class="mb-12">
				<h3 class="text-2xl md:text-3xl font-heading font-bold mb-6 text-[#8ec07c]">4.1 Scheduled Crawls</h3>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div class="bg-[#3c3836] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Target class="w-5 h-5 text-[#458588] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Daily/weekly automation</span>
						</div>
					</div>
					<div class="bg-[#3c3836] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Target class="w-5 h-5 text-[#458588] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Cloud Run cron triggers</span>
						</div>
					</div>
					<div class="bg-[#3c3836] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Target class="w-5 h-5 text-[#458588] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Email & Slack summaries</span>
						</div>
					</div>
					<div class="bg-[#3c3836] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Target class="w-5 h-5 text-[#458588] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Automatic comparisons between past crawls</span>
						</div>
					</div>
				</div>
			</div>

			<!-- 4.2 Reporting & White-Labeling -->
			<div class="mb-12">
				<h3 class="text-2xl md:text-3xl font-heading font-bold mb-6 text-[#8ec07c]">4.2 Reporting & White-Labeling</h3>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div class="bg-[#3c3836] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Target class="w-5 h-5 text-[#458588] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">PDF report generator</span>
						</div>
					</div>
					<div class="bg-[#3c3836] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Target class="w-5 h-5 text-[#458588] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Change tracking between crawls</span>
						</div>
					</div>
					<div class="bg-[#3c3836] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Target class="w-5 h-5 text-[#458588] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Shareable links for client reporting</span>
							<span class="ml-auto text-xs text-[#8ec07c] font-semibold">✓ COMPLETE</span>
						</div>
					</div>
					<div class="bg-[#3c3836] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Target class="w-5 h-5 text-[#458588] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Public-facing client reports (no login required)</span>
							<span class="ml-auto text-xs text-[#8ec07c] font-semibold">✓ COMPLETE</span>
						</div>
					</div>
				</div>
			</div>

			<!-- 4.3 Agency-Focused Enhancements -->
			<div class="mb-12">
				<h3 class="text-2xl md:text-3xl font-heading font-bold mb-6 text-[#8ec07c]">4.3 Agency-Focused Enhancements</h3>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div class="bg-[#3c3836] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Target class="w-5 h-5 text-[#458588] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Multi-project GSC/GA dashboards</span>
						</div>
					</div>
					<div class="bg-[#3c3836] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Target class="w-5 h-5 text-[#458588] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Project grouping/folders</span>
						</div>
					</div>
					<div class="bg-[#3c3836] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Target class="w-5 h-5 text-[#458588] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Integrations with task managers (Jira, Asana, ClickUp)</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</section>

<!-- Section 5: Long-Term Vision -->
<section class="py-20 bg-[#3c3836]">
	<div class="container mx-auto px-4">
		<div class="max-w-6xl mx-auto">
			<div class="flex items-center gap-4 mb-8">
				<div class="p-3 bg-[#8ec07c]/10 rounded-lg">
					<Sparkles class="w-8 h-8 text-[#8ec07c]" />
				</div>
				<h2 class="text-4xl md:text-5xl font-heading font-bold text-white">Long-Term Vision</h2>
				<span class="text-lg text-white/70">(6–12+ months)</span>
			</div>

			<!-- 5.1 Rank Tracking Enhancements -->
			<div class="mb-12">
				<h3 class="text-2xl md:text-3xl font-heading font-bold mb-6 text-[#8ec07c]">5.1 Rank Tracking Enhancements</h3>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Sparkles class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Daily SERP position monitoring</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Sparkles class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Competitor comparisons</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Sparkles class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Keyword grouping & insights</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Sparkles class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Alerts for ranking drops</span>
						</div>
					</div>
				</div>
			</div>

			<!-- 5.2 Intelligent SEO Assistant -->
			<div class="mb-12">
				<h3 class="text-2xl md:text-3xl font-heading font-bold mb-6 text-[#8ec07c]">5.2 Intelligent SEO Assistant</h3>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Sparkles class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Full conversational assistant inside dashboard</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Sparkles class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Auto-generated optimization briefs</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Sparkles class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Content audits and SEO scoring</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Sparkles class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Multi-domain insights for agencies</span>
						</div>
					</div>
				</div>
			</div>

			<!-- 5.3 Plugin Ecosystem -->
			<div class="mb-12">
				<h3 class="text-2xl md:text-3xl font-heading font-bold mb-6 text-[#8ec07c]">5.3 Plugin Ecosystem</h3>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Sparkles class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Community-driven custom issue detectors</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Sparkles class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Webhooks & API extensions</span>
						</div>
					</div>
					<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
						<div class="flex items-start gap-3">
							<Sparkles class="w-5 h-5 text-[#8ec07c] flex-shrink-0 mt-0.5" />
							<span class="text-white/90">Private plugins for enterprise clients</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</section>

<!-- Section 6: Transparency & Feedback -->
<section class="py-20 bg-[#2d2826]">
	<div class="container mx-auto px-4">
		<div class="max-w-4xl mx-auto">
			<div class="flex items-center gap-4 mb-8">
				<div class="p-3 bg-[#458588]/10 rounded-lg">
					<MessageSquare class="w-8 h-8 text-[#458588]" />
				</div>
				<h2 class="text-4xl md:text-5xl font-heading font-bold text-white">Transparency & Feedback</h2>
			</div>
			<div class="bg-[#3c3836] rounded-lg p-8 border border-white/10 mb-8">
				<p class="text-lg text-white/90 mb-6">
					We believe in building Barracuda in the open. Your feedback shapes our roadmap, and we're always listening to what matters most to you.
				</p>
				<ul class="space-y-4 text-white/80 mb-6">
					
					<li class="flex items-start gap-3">
						<MessageSquare class="w-5 h-5 text-[#458588] flex-shrink-0 mt-0.5" />
						<span>This roadmap is updated regularly as priorities shift and new opportunities emerge.</span>
					</li>
				</ul>
			</div>
		</div>
	</div>
</section>

<!-- Call to Action -->
<section class="py-20 bg-gradient-to-b from-[#3c3836] to-[#2d2826]">
	<div class="container mx-auto px-4">
		<div class="max-w-4xl mx-auto text-center">
			<h2 class="text-4xl md:text-5xl font-heading font-bold mb-6 text-white">
				Ready to Get Started?
			</h2>
			<p class="text-xl text-white/80 mb-10">
				Start crawling your site today and see what Barracuda can do for your SEO workflow.
			</p>
			<div class="flex flex-col sm:flex-row gap-4 justify-center mb-8">
				<a 
					href="https://app.barracudaseo.com" 
					class="inline-flex items-center justify-center gap-2 bg-[#8ec07c] hover:bg-[#a0d28c] text-[#3c3836] px-8 py-4 rounded-lg font-medium text-lg transition-colors" 
					target="_blank" 
					rel="noopener noreferrer"
				>
					Start a Free Crawl
					<ExternalLink class="w-5 h-5" />
				</a>
				
			</div>
			
		</div>
	</div>
</section>

