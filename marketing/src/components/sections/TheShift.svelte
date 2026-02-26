<script lang="ts">
	import { ArrowRight, Layers, Brain, Shield } from '@lucide/svelte';
</script>

<section class="bg-[#222] px-4 py-24">
	<div class="container mx-auto max-w-6xl">
		<div class="grid grid-cols-1 items-center gap-16 lg:grid-cols-2">
			<!-- Left: Copy -->
			<div>
				<p class="mb-4 text-sm font-medium tracking-widest text-[#8ec07c] uppercase">The shift</p>
				<h2
					class="mb-6 font-heading text-3xl leading-tight font-bold text-white md:text-4xl lg:text-5xl"
				>
					A decision layer
					<br />
					<span class="text-[#8ec07c]">between data and action.</span>
				</h2>
				<p class="mb-8 text-lg leading-relaxed text-white/60">
					Barracuda sits between raw crawl data and your next move — helping you prioritize issues,
					understand impact, and explain your decisions with confidence.
				</p>

				<div class="space-y-4">
					<div class="flex items-start gap-4">
						<div
							class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-[#8ec07c]/10"
						>
							<Layers class="h-5 w-5 text-[#8ec07c]" />
						</div>
						<div>
							<p class="mb-1 font-medium text-white">You still run crawls</p>
							<p class="text-sm text-white/50">
								Bring in data from your existing tools or use our built-in crawler.
							</p>
						</div>
					</div>
					<div class="flex items-start gap-4">
						<div
							class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-[#8ec07c]/10"
						>
							<Brain class="h-5 w-5 text-[#8ec07c]" />
						</div>
						<div>
							<p class="mb-1 font-medium text-white">You still apply judgment</p>
							<p class="text-sm text-white/50">AI assists your thinking — it doesn't replace it.</p>
						</div>
					</div>
					<div class="flex items-start gap-4">
						<div
							class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-[#8ec07c]/10"
						>
							<Shield class="h-5 w-5 text-[#8ec07c]" />
						</div>
						<div>
							<p class="mb-1 font-medium text-white">Barracuda removes the chaos</p>
							<p class="text-sm text-white/50">
								Clear priorities, transparent reasoning, defensible decisions.
							</p>
						</div>
					</div>
				</div>
			</div>

			<!-- Right: Visual flow diagram -->
			<div class="relative">
				<div class="space-y-4">
					<!-- Raw Data -->
					<div class="rounded-xl border border-white/5 bg-[#1d1d1d] p-5">
						<p class="mb-3 text-xs tracking-wider text-white/30 uppercase">Raw Crawl Data</p>
						<div class="flex flex-wrap gap-2">
							{#each ['2,847 pages', '142 errors', '87 warnings', '53 notices', '12 redirects'] as item}
								<span class="rounded-lg bg-[#282828] px-3 py-1.5 text-xs text-white/50">{item}</span
								>
							{/each}
						</div>
					</div>

					<!-- Arrow -->
					<div class="flex justify-center">
						<div
							class="flex h-10 w-10 items-center justify-center rounded-full border border-[#8ec07c]/20 bg-[#8ec07c]/10"
						>
							<ArrowRight class="h-4 w-4 rotate-90 text-[#8ec07c]" />
						</div>
					</div>

					<!-- Barracuda processing -->
					<div
						class="rounded-xl border border-[#8ec07c]/20 bg-gradient-to-r from-[#8ec07c]/10 to-[#8ec07c]/5 p-5"
					>
						<p class="mb-3 text-xs font-medium tracking-wider text-[#8ec07c] uppercase">
							Barracuda Analysis
						</p>
						<p class="text-sm text-white/60">
							Contextual AI prioritization, impact scoring, and reasoning generation
						</p>
					</div>

					<!-- Arrow -->
					<div class="flex justify-center">
						<div
							class="flex h-10 w-10 items-center justify-center rounded-full border border-[#8ec07c]/20 bg-[#8ec07c]/10"
						>
							<ArrowRight class="h-4 w-4 rotate-90 text-[#8ec07c]" />
						</div>
					</div>

					<!-- Clear Priorities -->
					<div class="rounded-xl border border-[#8ec07c]/20 bg-[#1d1d1d] p-5">
						<p class="mb-3 text-xs tracking-wider text-white/30 uppercase">Clear Priorities</p>
						<div class="space-y-2">
							<div class="flex items-center gap-3">
								<span
									class="flex h-5 w-5 items-center justify-center rounded-full bg-[#8ec07c] text-xs font-bold text-[#1d1d1d]"
									>1</span
								>
								<span class="text-sm text-white/80">Fix critical H1 tags on high-traffic pages</span
								>
							</div>
							<div class="flex items-center gap-3">
								<span
									class="flex h-5 w-5 items-center justify-center rounded-full bg-[#8ec07c]/60 text-xs font-bold text-[#1d1d1d]"
									>2</span
								>
								<span class="text-sm text-white/60">Resolve redirect chains affecting 34 pages</span
								>
							</div>
							<div class="flex items-center gap-3">
								<span
									class="flex h-5 w-5 items-center justify-center rounded-full bg-[#8ec07c]/30 text-xs font-bold text-white/60"
									>3</span
								>
								<span class="text-sm text-white/40">Add meta descriptions to new content</span>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</section>
