<script lang="ts">
	import { Quote } from '@lucide/svelte';
</script>

<section class="bg-[#1d1d1d] px-4 py-24">
	<div class="container mx-auto max-w-6xl">
		<!-- Stats bar -->
		<div class="mb-20 grid grid-cols-2 gap-6 md:grid-cols-4">
			<div class="text-center">
				<p class="mb-2 font-heading text-4xl font-bold text-white md:text-5xl">100+</p>
				<p class="text-sm text-white/40">Sites audited</p>
			</div>
			<div class="text-center">
				<p class="mb-2 font-heading text-4xl font-bold text-[#8ec07c] md:text-5xl">50k+</p>
				<p class="text-sm text-white/40">Pages crawled</p>
			</div>
			<div class="text-center">
				<p class="mb-2 font-heading text-4xl font-bold text-white md:text-5xl">80%</p>
				<p class="text-sm text-white/40">Fewer audit hours</p>
			</div>
			<div class="text-center">
				<p class="mb-2 font-heading text-4xl font-bold text-[#8ec07c] md:text-5xl">12x</p>
				<p class="text-sm text-white/40">Faster to priorities</p>
			</div>
		</div>

		<!-- Testimonial -->
		<div class="mx-auto max-w-3xl">
			<div class="relative rounded-2xl border border-white/5 bg-[#282828] p-8 md:p-12">
				<Quote class="absolute top-8 left-8 h-10 w-10 text-[#8ec07c]/20 md:top-12 md:left-12" />
				<blockquote class="relative z-10">
					<p
						class="mb-6 pl-0 text-xl leading-relaxed font-medium text-white/90 md:pl-6 md:text-2xl"
					>
						"Barracuda helped me finally explain <span class="text-[#8ec07c]">why</span> we were fixing
						certain issues instead of just reacting to audit noise."
					</p>
					<footer class="pl-0 md:pl-6">
						<div class="flex items-center gap-3">
							<div class="flex h-10 w-10 items-center justify-center rounded-full bg-[#8ec07c]/10">
								<span class="text-sm font-bold text-[#8ec07c]">SC</span>
							</div>
							<div>
								<p class="text-sm font-medium text-white">SEO Consultant</p>
								<p class="text-xs text-white/40">Technical SEO Specialist</p>
							</div>
						</div>
					</footer>
				</blockquote>
			</div>
		</div>
	</div>
</section>
