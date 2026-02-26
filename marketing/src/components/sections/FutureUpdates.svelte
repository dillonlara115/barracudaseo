<script lang="ts">
	import { Rocket, Sparkles, Zap, ArrowRight, TrendingUp } from '@lucide/svelte';
</script>

<section class="bg-[#2d2826] py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-6xl">
			<div class="mb-12 text-center">
				<div class="mb-6 flex items-center justify-center gap-4">
					<div class="rounded-lg bg-[#8ec07c]/10 p-3">
						<Rocket class="h-8 w-8 text-[#8ec07c]" />
					</div>
					<h2 class="font-heading text-4xl font-bold text-white md:text-5xl">What's Coming Next</h2>
				</div>
				<p class="mx-auto max-w-3xl text-xl text-white/80">
					We're constantly improving Barracuda with new features, AI-powered insights, and
					integrations. Here's a glimpse of what's on the horizon.
				</p>
			</div>

			<div class="mb-8 rounded-lg border border-[#8ec07c]/30 bg-[#8ec07c]/10 p-6">
				<h3 class="mb-3 flex items-center gap-2 font-heading text-xl font-bold text-[#8ec07c]">
					<Sparkles class="h-5 w-5" />
					Now Available
				</h3>
				<div class="grid grid-cols-1 gap-4 text-sm text-white/80 md:grid-cols-2">
					<div>
						<p class="mb-2 font-semibold text-white">✓ AI-Powered Features</p>
						<ul class="ml-4 space-y-1 text-white/70">
							<li>• AI recommendations for issues</li>
							<li>• Intelligent crawl summaries</li>
							<li>• Traffic-aware priority scoring</li>
						</ul>
					</div>
					<div>
						<p class="mb-2 font-semibold text-white">✓ Google Search Console Integration</p>
						<ul class="ml-4 space-y-1 text-white/70">
							<li>• Connect GSC properties</li>
							<li>• Enrich issues with traffic data</li>
							<li>• Prioritize fixes by impressions</li>
						</ul>
					</div>
				</div>
			</div>

			<div class="mb-10 grid grid-cols-1 gap-6 md:grid-cols-3">
				<div
					class="rounded-lg border border-white/10 bg-[#3c3836] p-6 transition-colors hover:border-[#8ec07c]/30"
				>
					<div class="mb-4 flex items-center gap-3">
						<div class="rounded-lg bg-[#458588]/10 p-2">
							<Zap class="h-6 w-6 text-[#458588]" />
						</div>
						<h3 class="font-heading text-xl font-bold text-[#8ec07c]">More Integrations</h3>
					</div>
					<p class="mb-4 text-white/70">
						Connect Barracuda with your favorite tools and workflows for seamless SEO management.
					</p>
					<ul class="space-y-2 text-sm text-white/60">
						<li class="flex items-start gap-2">
							<span class="text-[#8ec07c]">•</span>
							<span>Google Analytics 4 (GA4)</span>
						</li>
						<li class="flex items-start gap-2">
							<span class="text-[#8ec07c]">•</span>
							<span>Slack alerts & notifications</span>
						</li>
						<li class="flex items-start gap-2">
							<span class="text-[#8ec07c]">•</span>
							<span>Zapier / Make.com workflows</span>
						</li>
					</ul>
				</div>

				<div
					class="rounded-lg border border-white/10 bg-[#3c3836] p-6 transition-colors hover:border-[#8ec07c]/30"
				>
					<div class="mb-4 flex items-center gap-3">
						<div class="rounded-lg bg-[#d79921]/10 p-2">
							<Rocket class="h-6 w-6 text-[#d79921]" />
						</div>
						<h3 class="font-heading text-xl font-bold text-[#8ec07c]">Automation</h3>
					</div>
					<p class="mb-4 text-white/70">
						Set it and forget it. Scheduled crawls, automated reports, and change tracking to keep
						your SEO on track.
					</p>
					<ul class="space-y-2 text-sm text-white/60">
						<li class="flex items-start gap-2">
							<span class="text-[#8ec07c]">•</span>
							<span>Scheduled daily/weekly crawls</span>
						</li>
						<li class="flex items-start gap-2">
							<span class="text-[#8ec07c]">•</span>
							<span>Email & Slack summaries</span>
						</li>
						<li class="flex items-start gap-2">
							<span class="text-[#8ec07c]">•</span>
							<span>Change tracking between crawls</span>
						</li>
					</ul>
				</div>

				<div
					class="rounded-lg border border-white/10 bg-[#3c3836] p-6 transition-colors hover:border-[#8ec07c]/30"
				>
					<div class="mb-4 flex items-center gap-3">
						<div class="rounded-lg bg-[#b16286]/10 p-2">
							<TrendingUp class="h-6 w-6 text-[#b16286]" />
						</div>
						<h3 class="font-heading text-xl font-bold text-[#8ec07c]">Advanced Analytics</h3>
					</div>
					<p class="mb-4 text-white/70">
						Deeper insights and advanced reporting capabilities to understand your SEO performance.
					</p>
					<ul class="space-y-2 text-sm text-white/60">
						<li class="flex items-start gap-2">
							<span class="text-[#8ec07c]">•</span>
							<span>Local Pack tracking</span>
						</li>
						<li class="flex items-start gap-2">
							<span class="text-[#8ec07c]">•</span>
							<span>GeoGrid visualization</span>
						</li>
						<li class="flex items-start gap-2">
							<span class="text-[#8ec07c]">•</span>
							<span>Advanced SERP feature tracking</span>
						</li>
					</ul>
				</div>
			</div>

			<div class="text-center">
				<a
					href="/roadmap"
					class="inline-flex items-center gap-2 rounded-lg bg-[#8ec07c] px-6 py-3 font-medium text-[#3c3836] transition-colors hover:bg-[#a0d28c]"
				>
					View Full Roadmap
					<ArrowRight class="h-4 w-4" />
				</a>
			</div>
		</div>
	</div>
</section>
