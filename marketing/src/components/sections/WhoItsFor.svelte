<script lang="ts">
	import { Check, X } from '@lucide/svelte';
</script>

<section class="bg-[#222] py-24 px-4">
	<div class="container mx-auto max-w-5xl">
		<div class="text-center mb-16">
			<p class="text-sm uppercase tracking-widest text-[#8ec07c] font-medium mb-4">Who it's for</p>
			<h2 class="text-3xl md:text-4xl lg:text-5xl font-heading font-bold text-white">
				Built for decision-makers, not data hoarders
			</h2>
		</div>

		<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
			<!-- For You -->
			<div class="bg-[#1d1d1d] rounded-2xl border border-[#8ec07c]/20 p-8">
				<div class="inline-flex items-center gap-2 bg-[#8ec07c]/10 text-[#8ec07c] px-3 py-1 rounded-full text-sm font-medium mb-6">
					<Check class="w-3.5 h-3.5" />
					Great fit
				</div>
				<h3 class="text-xl font-heading font-bold text-white mb-6">Barracuda is for you if…</h3>
				<ul class="space-y-4">
					<li class="flex items-start gap-3">
						<div class="w-5 h-5 rounded-full bg-[#8ec07c]/10 flex items-center justify-center shrink-0 mt-0.5">
							<Check class="w-3 h-3 text-[#8ec07c]" />
						</div>
						<span class="text-white/70 text-sm leading-relaxed">You run audits and feel overwhelmed by issue volume</span>
					</li>
					<li class="flex items-start gap-3">
						<div class="w-5 h-5 rounded-full bg-[#8ec07c]/10 flex items-center justify-center shrink-0 mt-0.5">
							<Check class="w-3 h-3 text-[#8ec07c]" />
						</div>
						<span class="text-white/70 text-sm leading-relaxed">You need to explain SEO decisions to clients or stakeholders</span>
					</li>
					<li class="flex items-start gap-3">
						<div class="w-5 h-5 rounded-full bg-[#8ec07c]/10 flex items-center justify-center shrink-0 mt-0.5">
							<Check class="w-3 h-3 text-[#8ec07c]" />
						</div>
						<span class="text-white/70 text-sm leading-relaxed">You want fewer, better decisions — not more metrics</span>
					</li>
					<li class="flex items-start gap-3">
						<div class="w-5 h-5 rounded-full bg-[#8ec07c]/10 flex items-center justify-center shrink-0 mt-0.5">
							<Check class="w-3 h-3 text-[#8ec07c]" />
						</div>
						<span class="text-white/70 text-sm leading-relaxed">You value transparency in how priorities are determined</span>
					</li>
				</ul>
			</div>

			<!-- Not For You -->
			<div class="bg-[#1d1d1d] rounded-2xl border border-white/5 p-8">
				<div class="inline-flex items-center gap-2 bg-white/5 text-white/40 px-3 py-1 rounded-full text-sm font-medium mb-6">
					<X class="w-3.5 h-3.5" />
					Not the right fit
				</div>
				<h3 class="text-xl font-heading font-bold text-white/50 mb-6">Probably not for you if…</h3>
				<ul class="space-y-4">
					<li class="flex items-start gap-3">
						<div class="w-5 h-5 rounded-full bg-white/5 flex items-center justify-center shrink-0 mt-0.5">
							<X class="w-3 h-3 text-white/30" />
						</div>
						<span class="text-white/40 text-sm leading-relaxed">You want every possible SEO metric on one screen</span>
					</li>
					<li class="flex items-start gap-3">
						<div class="w-5 h-5 rounded-full bg-white/5 flex items-center justify-center shrink-0 mt-0.5">
							<X class="w-3 h-3 text-white/30" />
						</div>
						<span class="text-white/40 text-sm leading-relaxed">You enjoy digging through raw CSVs for hours</span>
					</li>
					<li class="flex items-start gap-3">
						<div class="w-5 h-5 rounded-full bg-white/5 flex items-center justify-center shrink-0 mt-0.5">
							<X class="w-3 h-3 text-white/30" />
						</div>
						<span class="text-white/40 text-sm leading-relaxed">You believe more data always equals better decisions</span>
					</li>
				</ul>
			</div>
		</div>
	</div>
</section>
