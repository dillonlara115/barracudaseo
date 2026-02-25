<script lang="ts">
	import { ArrowRight, Layers, Brain, Shield } from '@lucide/svelte';
</script>

<section class="bg-[#222] py-24 px-4">
	<div class="container mx-auto max-w-6xl">
		<div class="grid grid-cols-1 lg:grid-cols-2 gap-16 items-center">
			<!-- Left: Copy -->
			<div>
				<p class="text-sm uppercase tracking-widest text-[#8ec07c] font-medium mb-4">The shift</p>
				<h2 class="text-3xl md:text-4xl lg:text-5xl font-heading font-bold text-white leading-tight mb-6">
					A decision layer
					<br />
					<span class="text-[#8ec07c]">between data and action.</span>
				</h2>
				<p class="text-lg text-white/60 leading-relaxed mb-8">
					Barracuda sits between raw crawl data and your next move — helping you prioritize issues, understand impact, and explain your decisions with confidence.
				</p>

				<div class="space-y-4">
					<div class="flex items-start gap-4">
						<div class="w-10 h-10 rounded-xl bg-[#8ec07c]/10 flex items-center justify-center shrink-0">
							<Layers class="w-5 h-5 text-[#8ec07c]" />
						</div>
						<div>
							<p class="text-white font-medium mb-1">You still run crawls</p>
							<p class="text-sm text-white/50">Bring in data from your existing tools or use our built-in crawler.</p>
						</div>
					</div>
					<div class="flex items-start gap-4">
						<div class="w-10 h-10 rounded-xl bg-[#8ec07c]/10 flex items-center justify-center shrink-0">
							<Brain class="w-5 h-5 text-[#8ec07c]" />
						</div>
						<div>
							<p class="text-white font-medium mb-1">You still apply judgment</p>
							<p class="text-sm text-white/50">AI assists your thinking — it doesn't replace it.</p>
						</div>
					</div>
					<div class="flex items-start gap-4">
						<div class="w-10 h-10 rounded-xl bg-[#8ec07c]/10 flex items-center justify-center shrink-0">
							<Shield class="w-5 h-5 text-[#8ec07c]" />
						</div>
						<div>
							<p class="text-white font-medium mb-1">Barracuda removes the chaos</p>
							<p class="text-sm text-white/50">Clear priorities, transparent reasoning, defensible decisions.</p>
						</div>
					</div>
				</div>
			</div>

			<!-- Right: Visual flow diagram -->
			<div class="relative">
				<div class="space-y-4">
					<!-- Raw Data -->
					<div class="bg-[#1d1d1d] border border-white/5 rounded-xl p-5">
						<p class="text-xs uppercase tracking-wider text-white/30 mb-3">Raw Crawl Data</p>
						<div class="flex flex-wrap gap-2">
							{#each ['2,847 pages', '142 errors', '87 warnings', '53 notices', '12 redirects'] as item}
								<span class="bg-[#282828] text-white/50 text-xs px-3 py-1.5 rounded-lg">{item}</span>
							{/each}
						</div>
					</div>

					<!-- Arrow -->
					<div class="flex justify-center">
						<div class="w-10 h-10 rounded-full bg-[#8ec07c]/10 border border-[#8ec07c]/20 flex items-center justify-center">
							<ArrowRight class="w-4 h-4 text-[#8ec07c] rotate-90" />
						</div>
					</div>

					<!-- Barracuda processing -->
					<div class="bg-gradient-to-r from-[#8ec07c]/10 to-[#8ec07c]/5 border border-[#8ec07c]/20 rounded-xl p-5">
						<p class="text-xs uppercase tracking-wider text-[#8ec07c] mb-3 font-medium">Barracuda Analysis</p>
						<p class="text-sm text-white/60">Contextual AI prioritization, impact scoring, and reasoning generation</p>
					</div>

					<!-- Arrow -->
					<div class="flex justify-center">
						<div class="w-10 h-10 rounded-full bg-[#8ec07c]/10 border border-[#8ec07c]/20 flex items-center justify-center">
							<ArrowRight class="w-4 h-4 text-[#8ec07c] rotate-90" />
						</div>
					</div>

					<!-- Clear Priorities -->
					<div class="bg-[#1d1d1d] border border-[#8ec07c]/20 rounded-xl p-5">
						<p class="text-xs uppercase tracking-wider text-white/30 mb-3">Clear Priorities</p>
						<div class="space-y-2">
							<div class="flex items-center gap-3">
								<span class="w-5 h-5 rounded-full bg-[#8ec07c] text-[#1d1d1d] text-xs font-bold flex items-center justify-center">1</span>
								<span class="text-sm text-white/80">Fix critical H1 tags on high-traffic pages</span>
							</div>
							<div class="flex items-center gap-3">
								<span class="w-5 h-5 rounded-full bg-[#8ec07c]/60 text-[#1d1d1d] text-xs font-bold flex items-center justify-center">2</span>
								<span class="text-sm text-white/60">Resolve redirect chains affecting 34 pages</span>
							</div>
							<div class="flex items-center gap-3">
								<span class="w-5 h-5 rounded-full bg-[#8ec07c]/30 text-white/60 text-xs font-bold flex items-center justify-center">3</span>
								<span class="text-sm text-white/40">Add meta descriptions to new content</span>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</section>
