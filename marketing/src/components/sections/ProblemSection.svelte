<script lang="ts">
	import { AlertTriangle, List, HelpCircle } from '@lucide/svelte';
</script>

<section class="bg-[#1d1d1d] px-4 py-24">
	<div class="container mx-auto max-w-6xl">
		<div class="grid grid-cols-1 items-center gap-16 lg:grid-cols-2">
			<!-- Left: Visual -->
			<div class="relative order-2 lg:order-1">
				<div class="space-y-5 rounded-2xl border border-white/5 bg-[#282828] p-8">
					<!-- Simulated overwhelming audit output -->
					<div class="flex items-center gap-3 text-[#cc241d]/80">
						<AlertTriangle class="h-5 w-5 shrink-0" />
						<div class="h-3 flex-1 rounded-full bg-[#cc241d]/15">
							<div class="h-3 w-[85%] rounded-full bg-[#cc241d]/40"></div>
						</div>
						<span class="font-mono text-xs text-white/40">142 issues</span>
					</div>
					<div class="flex items-center gap-3 text-[#d79921]/80">
						<AlertTriangle class="h-5 w-5 shrink-0" />
						<div class="h-3 flex-1 rounded-full bg-[#d79921]/15">
							<div class="h-3 w-[60%] rounded-full bg-[#d79921]/40"></div>
						</div>
						<span class="font-mono text-xs text-white/40">87 warnings</span>
					</div>
					<div class="flex items-center gap-3 text-[#458588]/80">
						<List class="h-5 w-5 shrink-0" />
						<div class="h-3 flex-1 rounded-full bg-[#458588]/15">
							<div class="h-3 w-[40%] rounded-full bg-[#458588]/40"></div>
						</div>
						<span class="font-mono text-xs text-white/40">53 notices</span>
					</div>

					<div class="mt-2 border-t border-white/5 pt-5">
						<div class="grid grid-cols-3 gap-3">
							{#each ['Missing meta', 'Dup titles', 'Redirect chains', 'Slow pages', 'Broken links', '404 errors'] as issue}
								<div class="rounded-lg bg-[#1d1d1d] px-3 py-2 text-center text-xs text-white/40">
									{issue}
								</div>
							{/each}
						</div>
					</div>

					<div class="flex items-center gap-2 pt-2">
						<HelpCircle class="h-4 w-4 text-[#d79921]" />
						<p class="text-sm text-[#d79921] italic">Now what? Where do you even start?</p>
					</div>
				</div>
				<!-- Subtle glow -->
				<div class="absolute inset-0 -z-10 scale-105 rounded-2xl bg-[#cc241d]/5 blur-2xl"></div>
			</div>

			<!-- Right: Copy -->
			<div class="order-1 lg:order-2">
				<p class="mb-4 text-sm font-medium tracking-widest text-[#d79921] uppercase">The problem</p>
				<h2
					class="mb-6 font-heading text-3xl leading-tight font-bold text-white md:text-4xl lg:text-5xl"
				>
					SEO tools give you data.
					<br />
					<span class="text-white/50">They don't give you clarity.</span>
				</h2>
				<div class="space-y-5 text-lg leading-relaxed text-white/60">
					<p>
						You run a crawl and get hundreds of issues. Missing meta descriptions. Duplicate titles.
						Redirect chains. Slow pages.
					</p>
					<p class="font-medium text-white/80">
						Now comes the hard part: What actually matters? What can wait? What do you fix first —
						and why?
					</p>
					<p>Most SEO tools stop at detection. The interpretation is left entirely to you.</p>
				</div>
			</div>
		</div>
	</div>
</section>
