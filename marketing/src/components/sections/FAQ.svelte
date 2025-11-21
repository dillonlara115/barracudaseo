<script lang="ts">
	import { HelpCircle, ArrowRight } from 'lucide-svelte';

	let open1 = $state(false);
	let open2 = $state(false);
	let open3 = $state(false);
	let open4 = $state(false);
	let open5 = $state(false);
	let open6 = $state(false);
	let open7 = $state(false);

	function toggleItem(index: number) {
		switch (index) {
			case 1: open1 = !open1; break;
			case 2: open2 = !open2; break;
			case 3: open3 = !open3; break;
			case 4: open4 = !open4; break;
			case 5: open5 = !open5; break;
			case 6: open6 = !open6; break;
			case 7: open7 = !open7; break;
		}
	}
</script>

<section class="py-20 bg-[#3c3836]">
	<div class="container mx-auto px-4">
		<div class="max-w-6xl mx-auto">
			<div class="text-center mb-12">
				<div class="flex items-center justify-center gap-4 mb-6">
					<div class="p-3 bg-[#8ec07c]/10 rounded-lg">
						<HelpCircle class="w-8 h-8 text-[#8ec07c]" />
					</div>
					<h2 class="text-4xl md:text-5xl font-heading font-bold text-white">Frequently Asked Questions</h2>
				</div>
				<p class="text-xl text-white/80 max-w-3xl mx-auto">
					Get quick answers to common questions about Barracuda SEO.
				</p>
			</div>

			<div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-8 items-start">
				<!-- FAQ 1 -->
				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<button
						class="w-full px-6 py-4 text-left flex items-center justify-between hover:bg-[#3c3836] transition-colors"
						onclick={() => toggleItem(1)}
					>
						<h3 class="text-lg font-heading font-bold text-white pr-4">
							What makes Barracuda different from Screaming Frog?
						</h3>
						<svg
							class="w-5 h-5 text-[#8ec07c] flex-shrink-0 transition-transform {open1 ? 'rotate-180' : ''}"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
						</svg>
					</button>
					{#if open1}
						<div class="px-6 py-4 text-white/80 text-sm">
							<p class="mb-2">Screaming Frog is powerful but single-device and manual. Barracuda offers:</p>
							<ul class="list-disc list-inside space-y-1 ml-2">
								<li>Fast local crawls (CLI coming soon)</li>
								<li>Cloud-based dashboard for teams</li>
								<li>AI-powered issue explanations</li>
								<li>Built-in integrations</li>
							</ul>
						</div>
					{/if}
				</div>

				<!-- FAQ 2 -->
				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<button
						class="w-full px-6 py-4 text-left flex items-center justify-between hover:bg-[#3c3836] transition-colors"
						onclick={() => toggleItem(2)}
					>
						<h3 class="text-lg font-heading font-bold text-white pr-4">
							How is Barracuda different from SEMrush/Ahrefs?
						</h3>
						<svg
							class="w-5 h-5 text-[#8ec07c] flex-shrink-0 transition-transform {open2 ? 'rotate-180' : ''}"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
						</svg>
					</button>
					{#if open2}
						<div class="px-6 py-4 text-white/80 text-sm">
							<p class="mb-2">SEMrush and Ahrefs run crawls on their servers with limited control. Barracuda:</p>
							<ul class="list-disc list-inside space-y-1 ml-2">
								<li>Runs local crawls via CLI for full control</li>
								<li>Has no credits or crawl caps</li>
								<li>Stores complete crawl datasets</li>
								<li>Mirrors what your site serves to real browsers</li>
							</ul>
						</div>
					{/if}
				</div>

				<!-- FAQ 3 -->
				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<button
						class="w-full px-6 py-4 text-left flex items-center justify-between hover:bg-[#3c3836] transition-colors"
						onclick={() => toggleItem(3)}
					>
						<h3 class="text-lg font-heading font-bold text-white pr-4">
							Why choose Barracuda over Looker Studio?
						</h3>
						<svg
							class="w-5 h-5 text-[#8ec07c] flex-shrink-0 transition-transform {open3 ? 'rotate-180' : ''}"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
						</svg>
					</button>
					{#if open3}
						<div class="px-6 py-4 text-white/80 text-sm">
							<p class="mb-2">Looker Studio is great for reporting but can't crawl or diagnose SEO issues. Barracuda does:</p>
							<ul class="list-disc list-inside space-y-1 ml-2">
								<li>Website crawling</li>
								<li>SEO issue detection</li>
								<li>Data collection from GSC, GA4, Clarity</li>
								<li>Automatic dashboard generation</li>
								<li>Full issue prioritization</li>
							</ul>
						</div>
					{/if}
				</div>

				<!-- FAQ 4 -->
				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<button
						class="w-full px-6 py-4 text-left flex items-center justify-between hover:bg-[#3c3836] transition-colors"
						onclick={() => toggleItem(4)}
					>
						<h3 class="text-lg font-heading font-bold text-white pr-4">
							Does Barracuda replace my SEO tools?
						</h3>
						<svg
							class="w-5 h-5 text-[#8ec07c] flex-shrink-0 transition-transform {open4 ? 'rotate-180' : ''}"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
						</svg>
					</button>
					{#if open4}
						<div class="px-6 py-4 text-white/80 text-sm">
							<p class="mb-2">Barracuda enhances your SEO stack:</p>
							<p class="mb-1 font-semibold text-white">Replaces:</p>
							<ul class="list-disc list-inside space-y-1 ml-2 mb-2">
								<li>Screaming Frog for audits</li>
								<li>Manual CSV exports</li>
								<li>Custom Looker Studio builds</li>
							</ul>
							<p class="mb-1 font-semibold text-white">Complements:</p>
							<ul class="list-disc list-inside space-y-1 ml-2">
								<li>SEMrush/Ahrefs for keyword research</li>
								<li>GSC, GA4, Clarity for insights</li>
							</ul>
						</div>
					{/if}
				</div>

				<!-- FAQ 5 -->
				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<button
						class="w-full px-6 py-4 text-left flex items-center justify-between hover:bg-[#3c3836] transition-colors"
						onclick={() => toggleItem(5)}
					>
						<h3 class="text-lg font-heading font-bold text-white pr-4">
							How fast is the Barracuda crawler?
						</h3>
						<svg
							class="w-5 h-5 text-[#8ec07c] flex-shrink-0 transition-transform {open5 ? 'rotate-180' : ''}"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
						</svg>
					</button>
					{#if open5}
						<div class="px-6 py-4 text-white/80 text-sm">
							<p class="mb-2">The CLI will be extremely fast due to:</p>
							<ul class="list-disc list-inside space-y-1 ml-2">
								<li>Go-based worker-pool engine</li>
								<li>Concurrent crawling</li>
								<li>Intelligent robots.txt + sitemap handling</li>
								<li>Lightweight architecture</li>
							</ul>
							<p class="mt-2 font-semibold text-white">Expected speed: 100–500 pages/min</p>
						</div>
					{/if}
				</div>

				<!-- FAQ 6 -->
				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<button
						class="w-full px-6 py-4 text-left flex items-center justify-between hover:bg-[#3c3836] transition-colors"
						onclick={() => toggleItem(6)}
					>
						<h3 class="text-lg font-heading font-bold text-white pr-4">
							What makes issue detection unique?
						</h3>
						<svg
							class="w-5 h-5 text-[#8ec07c] flex-shrink-0 transition-transform {open6 ? 'rotate-180' : ''}"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
						</svg>
					</button>
					{#if open6}
						<div class="px-6 py-4 text-white/80 text-sm">
							<p class="mb-2">Most crawlers just report raw issues. Barracuda:</p>
							<ul class="list-disc list-inside space-y-1 ml-2">
								<li>Prioritizes by severity + impact</li>
								<li>Groups by URL structure</li>
								<li>Connects crawl data with GSC/GA4</li>
								<li>Uses AI to explain root causes (coming soon)</li>
							</ul>
						</div>
					{/if}
				</div>

			</div>

			<div class="text-center mt-8">
				<a
					href="/faq"
					class="inline-flex items-center gap-2 text-[#8ec07c] hover:text-[#a0d28c] transition-colors font-medium"
				>
					View All FAQs
					<ArrowRight class="w-4 h-4" />
				</a>
			</div>
		</div>
	</div>
</section>

