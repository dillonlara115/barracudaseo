<script lang="ts">
	import { Quote } from '@lucide/svelte';
</script>

<section class="bg-[#1d1d1d] py-24 px-4">
	<div class="container mx-auto max-w-6xl">
		<!-- Stats bar -->
		<div class="grid grid-cols-2 md:grid-cols-4 gap-6 mb-20">
			<div class="text-center">
				<p class="text-4xl md:text-5xl font-heading font-bold text-white mb-2">100+</p>
				<p class="text-sm text-white/40">Sites audited</p>
			</div>
			<div class="text-center">
				<p class="text-4xl md:text-5xl font-heading font-bold text-[#8ec07c] mb-2">50k+</p>
				<p class="text-sm text-white/40">Pages crawled</p>
			</div>
			<div class="text-center">
				<p class="text-4xl md:text-5xl font-heading font-bold text-white mb-2">80%</p>
				<p class="text-sm text-white/40">Fewer audit hours</p>
			</div>
			<div class="text-center">
				<p class="text-4xl md:text-5xl font-heading font-bold text-[#8ec07c] mb-2">12x</p>
				<p class="text-sm text-white/40">Faster to priorities</p>
			</div>
		</div>

		<!-- Testimonial -->
		<div class="max-w-3xl mx-auto">
			<div class="relative bg-[#282828] rounded-2xl border border-white/5 p-8 md:p-12">
				<Quote class="w-10 h-10 text-[#8ec07c]/20 absolute top-8 left-8 md:top-12 md:left-12" />
				<blockquote class="relative z-10">
					<p class="text-xl md:text-2xl text-white/90 leading-relaxed font-medium mb-6 pl-0 md:pl-6">
						"Barracuda helped me finally explain <span class="text-[#8ec07c]">why</span> we were fixing certain issues instead of just reacting to audit noise."
					</p>
					<footer class="pl-0 md:pl-6">
						<div class="flex items-center gap-3">
							<div class="w-10 h-10 rounded-full bg-[#8ec07c]/10 flex items-center justify-center">
								<span class="text-[#8ec07c] font-bold text-sm">SC</span>
							</div>
							<div>
								<p class="text-white font-medium text-sm">SEO Consultant</p>
								<p class="text-white/40 text-xs">Technical SEO Specialist</p>
							</div>
						</div>
					</footer>
				</blockquote>
			</div>
		</div>
	</div>
</section>
