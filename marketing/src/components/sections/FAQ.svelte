<script lang="ts">
	import { CircleQuestionMark, ArrowRight } from '@lucide/svelte';

	let open1 = $state(false);
	let open2 = $state(false);
	let open3 = $state(false);
	let open4 = $state(false);
	let open5 = $state(false);
	let open6 = $state(false);

	function toggleItem(index: number) {
		switch (index) {
			case 1: open1 = !open1; break;
			case 2: open2 = !open2; break;
			case 3: open3 = !open3; break;
			case 4: open4 = !open4; break;
			case 5: open5 = !open5; break;
			case 6: open6 = !open6; break;
		}
	}
</script>

<section class="py-20 bg-[#3c3836]">
	<div class="container mx-auto px-4">
		<div class="max-w-6xl mx-auto">
			<div class="text-center mb-12">
				<div class="flex items-center justify-center gap-4 mb-6">
					<div class="p-3 bg-[#8ec07c]/10 rounded-lg">
						<CircleQuestionMark class="w-8 h-8 text-[#8ec07c]" />
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
							Will BarracudaSEO replace my existing SEO tools?
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
						<div class="px-6 py-4 text-white/80 text-sm space-y-4">
							<p>In many cases, yes — and in others, it complements them.</p>
							<p>BarracudaSEO includes its own crawler and issue detection, covering many of the same technical audit capabilities found in larger tools like SEMrush and Ahrefs.</p>
							<p>For many consultants, agencies, and teams, Barracuda can replace those tools entirely — especially when the goal is to run audits, identify issues, and understand what to fix next without paying for enterprise-level feature bloat.</p>
							<p>Some teams still choose to keep specialized tools for niche use cases. But Barracuda’s goal is simple:</p>
							<p class="font-medium text-white">Replace expensive, overwhelming SEO platforms with a focused system that finds issues, prioritizes what matters, and tells you exactly what to do next.</p>
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
							How is this different from other SEO tools?
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
						<div class="px-6 py-4 text-white/80 text-sm space-y-4">
							<p>Most SEO tools focus on surfacing issues. Barracuda focuses on helping you decide:</p>
							<ul class="list-disc list-inside space-y-1 ml-2">
								<li>What to fix first</li>
								<li>What can safely wait</li>
								<li>What doesn’t matter right now</li>
								<li>How to explain those decisions clearly</li>
							</ul>
							<p>If you’re looking for more metrics, Barracuda may not be the right fit. If you’re looking for confidence and clarity, it probably is.</p>
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
							Can I trust the prioritization and recommendations?
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
						<div class="px-6 py-4 text-white/80 text-sm space-y-4">
							<p>Barracuda doesn’t hide its reasoning. For every priority, you can see:</p>
							<ul class="list-disc list-inside space-y-1 ml-2">
								<li>Why it matters</li>
								<li>What data informed it</li>
								<li>What was deprioritized and why</li>
							</ul>
							<p>You’re never asked to blindly trust a score or a black box. Barracuda supports your judgment — it doesn’t replace it.</p>
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
							Do you offer refunds?
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
						<div class="px-6 py-4 text-white/80 text-sm space-y-4">
							<p>Yes.</p>
							<p>We offer a 30-day money-back guarantee. If Barracuda doesn’t reduce your audit stress or help you feel more confident in your decisions, we don’t want you stuck with it.</p>
							<p>No hard feelings. No hoops.</p>
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
							Is BarracudaSEO right for beginners?
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
						<div class="px-6 py-4 text-white/80 text-sm space-y-4">
							<p>Barracuda works best for people who already understand SEO basics but want help deciding what matters most.</p>
							<p>If you’re brand new to SEO, you may want a learning-focused tool first. If you’re experienced and tired of audit overwhelm, Barracuda was built for you.</p>
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
							Can I try BarracudaSEO before I buy?
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
						<div class="px-6 py-4 text-white/80 text-sm space-y-4">
							<p>Yes.</p>
							<p>The Free plan gives you access to the core crawling and analysis experience so you can see how Barracuda approaches prioritization and clarity before committing.</p>
							<p>We want you to understand how decisions are surfaced and explained, not just click around a demo. Upgrade only when you’re ready for cloud projects, collaboration, and ongoing decision support.</p>
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
