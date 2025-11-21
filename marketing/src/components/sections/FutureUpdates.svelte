<script lang="ts">
	import { Rocket, Sparkles, Zap, ArrowRight } from 'lucide-svelte';
</script>

<section class="py-20 bg-[#2d2826]">
	<div class="container mx-auto px-4">
		<div class="max-w-6xl mx-auto">
			<div class="text-center mb-12">
				<div class="flex items-center justify-center gap-4 mb-6">
					<div class="p-3 bg-[#8ec07c]/10 rounded-lg">
						<Rocket class="w-8 h-8 text-[#8ec07c]" />
					</div>
					<h2 class="text-4xl md:text-5xl font-heading font-bold text-white">What's Coming Next</h2>
				</div>
				<p class="text-xl text-white/80 max-w-3xl mx-auto">
					We're constantly improving Barracuda with new features, AI-powered insights, and integrations. Here's a glimpse of what's on the horizon.
				</p>
			</div>

			<div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-10">
				<div class="bg-[#3c3836] rounded-lg p-6 border border-white/10 hover:border-[#8ec07c]/30 transition-colors">
					<div class="flex items-center gap-3 mb-4">
						<div class="p-2 bg-[#8ec07c]/10 rounded-lg">
							<Sparkles class="w-6 h-6 text-[#8ec07c]" />
						</div>
						<h3 class="text-xl font-heading font-bold text-[#8ec07c]">AI Integration</h3>
					</div>
					<p class="text-white/70 mb-4">
						AI-powered recommendations, crawl summaries, and intelligent priority scoring to help you fix issues faster.
					</p>
					<ul class="space-y-2 text-sm text-white/60">
						<li class="flex items-start gap-2">
							<span class="text-[#8ec07c]">•</span>
							<span>AI recommendations for issues</span>
						</li>
						<li class="flex items-start gap-2">
							<span class="text-[#8ec07c]">•</span>
							<span>Page-level rewrite suggestions</span>
						</li>
						<li class="flex items-start gap-2">
							<span class="text-[#8ec07c]">•</span>
							<span>Intelligent SEO assistant</span>
						</li>
					</ul>
				</div>

				<div class="bg-[#3c3836] rounded-lg p-6 border border-white/10 hover:border-[#8ec07c]/30 transition-colors">
					<div class="flex items-center gap-3 mb-4">
						<div class="p-2 bg-[#458588]/10 rounded-lg">
							<Zap class="w-6 h-6 text-[#458588]" />
						</div>
						<h3 class="text-xl font-heading font-bold text-[#8ec07c]">More Integrations</h3>
					</div>
					<p class="text-white/70 mb-4">
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

				<div class="bg-[#3c3836] rounded-lg p-6 border border-white/10 hover:border-[#8ec07c]/30 transition-colors">
					<div class="flex items-center gap-3 mb-4">
						<div class="p-2 bg-[#d79921]/10 rounded-lg">
							<Rocket class="w-6 h-6 text-[#d79921]" />
						</div>
						<h3 class="text-xl font-heading font-bold text-[#8ec07c]">Automation</h3>
					</div>
					<p class="text-white/70 mb-4">
						Set it and forget it. Scheduled crawls, automated reports, and change tracking to keep your SEO on track.
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
			</div>

			<div class="text-center">
				<a 
					href="/roadmap" 
					class="inline-flex items-center gap-2 bg-[#8ec07c] hover:bg-[#a0d28c] text-[#3c3836] px-6 py-3 rounded-lg font-medium transition-colors"
				>
					View Full Roadmap
					<ArrowRight class="w-4 h-4" />
				</a>
			</div>
		</div>
	</div>
</section>

