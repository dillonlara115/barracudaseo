<script lang="ts">
	import {
		Search,
		ChartBar,
		Target,
		Plug,
		Users,
		Sparkles,
		TrendingUp,
		Share2
	} from '@lucide/svelte';
</script>

<section class="bg-[#3c3836] py-20">
	<div class="container mx-auto px-4">
		<div class="mb-16 text-center">
			<h2 class="mb-4 font-heading text-4xl font-bold text-white md:text-5xl">Powerful Features</h2>
			<p class="mx-auto max-w-2xl text-xl text-white/70">
				Everything you need to crawl, analyze, and optimize your website's SEO
			</p>
		</div>
		<div class="grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
			<div class="rounded-lg border border-white/10 bg-[#2d2826] p-6 shadow-xl">
				<div class="mb-3 flex items-center gap-3">
					<div class="rounded-lg bg-[#8ec07c]/10 p-2">
						<Search class="h-5 w-5 text-[#8ec07c]" />
					</div>
					<h3 class="font-heading text-xl font-bold text-[#8ec07c]">Discover</h3>
				</div>
				<p class="text-white/70">
					Crawl up to tens of thousands of URLs and identify broken links, duplicate tags, missing
					meta data, slow pages, and more — fast.
				</p>
			</div>
			<div class="rounded-lg border border-white/10 bg-[#2d2826] p-6 shadow-xl">
				<div class="mb-3 flex items-center gap-3">
					<div class="rounded-lg bg-[#8ec07c]/10 p-2">
						<ChartBar class="h-5 w-5 text-[#8ec07c]" />
					</div>
					<h3 class="font-heading text-xl font-bold text-[#8ec07c]">Analyze</h3>
				</div>
				<p class="text-white/70">
					Visualize your site architecture, link structure, and SEO health inside a clean, intuitive
					web dashboard.
				</p>
			</div>
			<div class="rounded-lg border border-white/10 bg-[#2d2826] p-6 shadow-xl">
				<div class="mb-3 flex items-center gap-3">
					<div class="rounded-lg bg-[#8ec07c]/10 p-2">
						<Target class="h-5 w-5 text-[#8ec07c]" />
					</div>
					<h3 class="font-heading text-xl font-bold text-[#8ec07c]">Prioritize</h3>
				</div>
				<p class="text-white/70">
					Get clear issue summaries with built-in severity and impact scoring to focus on what truly
					matters.
				</p>
			</div>
			<div class="relative rounded-lg border border-[#8ec07c]/30 bg-[#2d2826] p-6 shadow-xl">
				<div
					class="absolute -top-2 -right-2 rounded bg-[#8ec07c] px-2 py-1 text-xs font-bold text-[#3c3836]"
				>
					Pro
				</div>
				<div class="mb-3 flex items-center gap-3">
					<div class="rounded-lg bg-[#8ec07c]/10 p-2">
						<Sparkles class="h-5 w-5 text-[#8ec07c]" />
					</div>
					<h3 class="font-heading text-xl font-bold text-[#8ec07c]">AI Insights</h3>
				</div>
				<p class="text-white/70">
					AI-powered issue explanations, crawl summaries, and intelligent recommendations to help
					you fix issues faster and understand root causes.
				</p>
			</div>
			<div class="relative rounded-lg border border-[#8ec07c]/30 bg-[#2d2826] p-6 shadow-xl">
				<div
					class="absolute -top-2 -right-2 rounded bg-[#8ec07c] px-2 py-1 text-xs font-bold text-[#3c3836]"
				>
					Pro
				</div>
				<div class="mb-3 flex items-center gap-3">
					<div class="rounded-lg bg-[#8ec07c]/10 p-2">
						<TrendingUp class="h-5 w-5 text-[#8ec07c]" />
					</div>
					<h3 class="font-heading text-xl font-bold text-[#8ec07c]">Rank Tracking</h3>
				</div>
				<p class="text-white/70">
					Track keyword rankings over time with automatic daily or weekly checks. Monitor position
					changes, view historical trends, and connect rankings to your crawl data.
				</p>
			</div>
			<div class="relative rounded-lg border border-[#8ec07c]/30 bg-[#2d2826] p-6 shadow-xl">
				<div
					class="absolute -top-2 -right-2 rounded bg-[#8ec07c] px-2 py-1 text-xs font-bold text-[#3c3836]"
				>
					Pro
				</div>
				<div class="mb-3 flex items-center gap-3">
					<div class="rounded-lg bg-[#8ec07c]/10 p-2">
						<Users class="h-5 w-5 text-[#8ec07c]" />
					</div>
					<h3 class="font-heading text-xl font-bold text-[#8ec07c]">Collaborate</h3>
				</div>
				<p class="text-white/70">
					Invite teammates with role-based permissions. Perfect for agencies and teams who need
					shared project access and seamless collaboration.
				</p>
			</div>
			<div class="relative rounded-lg border border-[#8ec07c]/30 bg-[#2d2826] p-6 shadow-xl">
				<div
					class="absolute -top-2 -right-2 rounded bg-[#8ec07c] px-2 py-1 text-xs font-bold text-[#3c3836]"
				>
					Pro
				</div>
				<div class="mb-3 flex items-center gap-3">
					<div class="rounded-lg bg-[#8ec07c]/10 p-2">
						<Plug class="h-5 w-5 text-[#8ec07c]" />
					</div>
					<h3 class="font-heading text-xl font-bold text-[#8ec07c]">Integrate</h3>
				</div>
				<p class="text-white/70">
					Connect Google Search Console, Analytics, and Microsoft Clarity to enrich your audit data
					and track performance over time.
				</p>
			</div>
			<div class="relative rounded-lg border border-[#8ec07c]/30 bg-[#2d2826] p-6 shadow-xl">
				<div
					class="absolute -top-2 -right-2 rounded bg-[#8ec07c] px-2 py-1 text-xs font-bold text-[#3c3836]"
				>
					Pro
				</div>
				<div class="mb-3 flex items-center gap-3">
					<div class="rounded-lg bg-[#8ec07c]/10 p-2">
						<Share2 class="h-5 w-5 text-[#8ec07c]" />
					</div>
					<h3 class="font-heading text-xl font-bold text-[#8ec07c]">Share Reports</h3>
				</div>
				<p class="text-white/70">
					Create public, shareable reports for clients or stakeholders. Share crawl results and
					issue summaries without requiring login.
				</p>
			</div>
		</div>
	</div>
</section>
