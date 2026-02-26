<script lang="ts">
	import { Check, X } from '@lucide/svelte';
</script>

<section class="bg-[#222] px-4 py-24">
	<div class="container mx-auto max-w-5xl">
		<div class="mb-16 text-center">
			<p class="mb-4 text-sm font-medium tracking-widest text-[#8ec07c] uppercase">Who it's for</p>
			<h2 class="font-heading text-3xl font-bold text-white md:text-4xl lg:text-5xl">
				Built for decision-makers, not data hoarders
			</h2>
		</div>

		<div class="grid grid-cols-1 gap-6 md:grid-cols-2">
			<!-- For You -->
			<div class="rounded-2xl border border-[#8ec07c]/20 bg-[#1d1d1d] p-8">
				<div
					class="mb-6 inline-flex items-center gap-2 rounded-full bg-[#8ec07c]/10 px-3 py-1 text-sm font-medium text-[#8ec07c]"
				>
					<Check class="h-3.5 w-3.5" />
					Great fit
				</div>
				<h3 class="mb-6 font-heading text-xl font-bold text-white">Barracuda is for you if…</h3>
				<ul class="space-y-4">
					<li class="flex items-start gap-3">
						<div
							class="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-[#8ec07c]/10"
						>
							<Check class="h-3 w-3 text-[#8ec07c]" />
						</div>
						<span class="text-sm leading-relaxed text-white/70"
							>You run audits and feel overwhelmed by issue volume</span
						>
					</li>
					<li class="flex items-start gap-3">
						<div
							class="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-[#8ec07c]/10"
						>
							<Check class="h-3 w-3 text-[#8ec07c]" />
						</div>
						<span class="text-sm leading-relaxed text-white/70"
							>You need to explain SEO decisions to clients or stakeholders</span
						>
					</li>
					<li class="flex items-start gap-3">
						<div
							class="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-[#8ec07c]/10"
						>
							<Check class="h-3 w-3 text-[#8ec07c]" />
						</div>
						<span class="text-sm leading-relaxed text-white/70"
							>You want fewer, better decisions — not more metrics</span
						>
					</li>
					<li class="flex items-start gap-3">
						<div
							class="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-[#8ec07c]/10"
						>
							<Check class="h-3 w-3 text-[#8ec07c]" />
						</div>
						<span class="text-sm leading-relaxed text-white/70"
							>You value transparency in how priorities are determined</span
						>
					</li>
				</ul>
			</div>

			<!-- Not For You -->
			<div class="rounded-2xl border border-white/5 bg-[#1d1d1d] p-8">
				<div
					class="mb-6 inline-flex items-center gap-2 rounded-full bg-white/5 px-3 py-1 text-sm font-medium text-white/40"
				>
					<X class="h-3.5 w-3.5" />
					Not the right fit
				</div>
				<h3 class="mb-6 font-heading text-xl font-bold text-white/50">Probably not for you if…</h3>
				<ul class="space-y-4">
					<li class="flex items-start gap-3">
						<div
							class="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-white/5"
						>
							<X class="h-3 w-3 text-white/30" />
						</div>
						<span class="text-sm leading-relaxed text-white/40"
							>You want every possible SEO metric on one screen</span
						>
					</li>
					<li class="flex items-start gap-3">
						<div
							class="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-white/5"
						>
							<X class="h-3 w-3 text-white/30" />
						</div>
						<span class="text-sm leading-relaxed text-white/40"
							>You enjoy digging through raw CSVs for hours</span
						>
					</li>
					<li class="flex items-start gap-3">
						<div
							class="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-white/5"
						>
							<X class="h-3 w-3 text-white/30" />
						</div>
						<span class="text-sm leading-relaxed text-white/40"
							>You believe more data always equals better decisions</span
						>
					</li>
				</ul>
			</div>
		</div>
	</div>
</section>
