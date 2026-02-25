<script lang="ts">
	import { AlertTriangle, List, HelpCircle } from '@lucide/svelte';
</script>

<section class="bg-[#1d1d1d] py-24 px-4">
	<div class="container mx-auto max-w-6xl">
		<div class="grid grid-cols-1 lg:grid-cols-2 gap-16 items-center">
			<!-- Left: Visual -->
			<div class="relative order-2 lg:order-1">
				<div class="bg-[#282828] rounded-2xl border border-white/5 p-8 space-y-5">
					<!-- Simulated overwhelming audit output -->
					<div class="flex items-center gap-3 text-[#cc241d]/80">
						<AlertTriangle class="w-5 h-5 shrink-0" />
						<div class="flex-1 h-3 bg-[#cc241d]/15 rounded-full">
							<div class="h-3 bg-[#cc241d]/40 rounded-full w-[85%]"></div>
						</div>
						<span class="text-xs font-mono text-white/40">142 issues</span>
					</div>
					<div class="flex items-center gap-3 text-[#d79921]/80">
						<AlertTriangle class="w-5 h-5 shrink-0" />
						<div class="flex-1 h-3 bg-[#d79921]/15 rounded-full">
							<div class="h-3 bg-[#d79921]/40 rounded-full w-[60%]"></div>
						</div>
						<span class="text-xs font-mono text-white/40">87 warnings</span>
					</div>
					<div class="flex items-center gap-3 text-[#458588]/80">
						<List class="w-5 h-5 shrink-0" />
						<div class="flex-1 h-3 bg-[#458588]/15 rounded-full">
							<div class="h-3 bg-[#458588]/40 rounded-full w-[40%]"></div>
						</div>
						<span class="text-xs font-mono text-white/40">53 notices</span>
					</div>

					<div class="border-t border-white/5 pt-5 mt-2">
						<div class="grid grid-cols-3 gap-3">
							{#each ['Missing meta', 'Dup titles', 'Redirect chains', 'Slow pages', 'Broken links', '404 errors'] as issue}
								<div class="bg-[#1d1d1d] rounded-lg px-3 py-2 text-xs text-white/40 text-center">{issue}</div>
							{/each}
						</div>
					</div>

					<div class="flex items-center gap-2 pt-2">
						<HelpCircle class="w-4 h-4 text-[#d79921]" />
						<p class="text-sm text-[#d79921] italic">Now what? Where do you even start?</p>
					</div>
				</div>
				<!-- Subtle glow -->
				<div class="absolute inset-0 -z-10 bg-[#cc241d]/5 rounded-2xl blur-2xl scale-105"></div>
			</div>

			<!-- Right: Copy -->
			<div class="order-1 lg:order-2">
				<p class="text-sm uppercase tracking-widest text-[#d79921] font-medium mb-4">The problem</p>
				<h2 class="text-3xl md:text-4xl lg:text-5xl font-heading font-bold text-white leading-tight mb-6">
					SEO tools give you data.
					<br />
					<span class="text-white/50">They don't give you clarity.</span>
				</h2>
				<div class="space-y-5 text-lg text-white/60 leading-relaxed">
					<p>
						You run a crawl and get hundreds of issues. Missing meta descriptions. Duplicate titles. Redirect chains. Slow pages.
					</p>
					<p class="text-white/80 font-medium">
						Now comes the hard part: What actually matters? What can wait? What do you fix first — and why?
					</p>
					<p>
						Most SEO tools stop at detection. The interpretation is left entirely to you.
					</p>
				</div>
			</div>
		</div>
	</div>
</section>
