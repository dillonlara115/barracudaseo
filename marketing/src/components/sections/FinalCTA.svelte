<script lang="ts">
	import { trackSignup } from '$lib/analytics';
	import { ArrowRight } from '@lucide/svelte';

	function handleSignupClick() {
		trackSignup({
			source: 'Start your first audit',
			location: 'final-cta'
		});
	}
</script>

<section class="relative overflow-hidden bg-[#1d1d1d] py-28 px-4">
	<!-- Gradient background -->
	<div class="absolute inset-0 bg-gradient-to-br from-[#8ec07c]/10 via-transparent to-[#d79921]/5 pointer-events-none"></div>
	<div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[400px] bg-[#8ec07c]/5 rounded-full blur-[120px] pointer-events-none"></div>

	<div class="relative z-10 max-w-3xl mx-auto text-center">
		<h2 class="text-3xl md:text-5xl lg:text-6xl font-heading font-bold text-white leading-tight mb-6">
			Stop second-guessing
			<br />
			your audits.
		</h2>
		<p class="text-lg md:text-xl text-white/50 mb-10 max-w-xl mx-auto">
			Move from data to decisions — with clarity you can stand behind.
		</p>

		<div class="flex flex-col sm:flex-row gap-4 justify-center">
			<a
				href="https://app.barracudaseo.com"
				class="group inline-flex items-center justify-center gap-2 bg-[#8ec07c] hover:bg-[#a0d28c] text-[#1d1d1d] px-8 py-4 rounded-xl font-semibold text-lg transition-all shadow-lg shadow-[#8ec07c]/10 hover:shadow-[#8ec07c]/25 hover:-translate-y-0.5"
				target="_blank"
				rel="noopener noreferrer"
				onclick={handleSignupClick}
			>
				Start your first audit
				<ArrowRight class="w-5 h-5 transition-transform group-hover:translate-x-0.5" />
			</a>
		</div>

		<p class="mt-6 text-sm text-white/30">No credit card required. No black box. Just clearer decisions.</p>
	</div>
</section>
