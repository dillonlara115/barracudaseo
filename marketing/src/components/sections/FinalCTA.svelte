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

<section class="relative overflow-hidden bg-[#1d1d1d] px-4 py-28">
	<!-- Gradient background -->
	<div
		class="pointer-events-none absolute inset-0 bg-gradient-to-br from-[#8ec07c]/10 via-transparent to-[#d79921]/5"
	></div>
	<div
		class="pointer-events-none absolute top-1/2 left-1/2 h-[400px] w-[800px] -translate-x-1/2 -translate-y-1/2 rounded-full bg-[#8ec07c]/5 blur-[120px]"
	></div>

	<div class="relative z-10 mx-auto max-w-3xl text-center">
		<h2
			class="mb-6 font-heading text-3xl leading-tight font-bold text-white md:text-5xl lg:text-6xl"
		>
			Stop second-guessing
			<br />
			your audits.
		</h2>
		<p class="mx-auto mb-10 max-w-xl text-lg text-white/50 md:text-xl">
			Move from data to decisions — with clarity you can stand behind.
		</p>

		<div class="flex flex-col justify-center gap-4 sm:flex-row">
			<a
				href="https://app.barracudaseo.com"
				class="group inline-flex items-center justify-center gap-2 rounded-xl bg-[#8ec07c] px-8 py-4 text-lg font-semibold text-[#1d1d1d] shadow-lg shadow-[#8ec07c]/10 transition-all hover:-translate-y-0.5 hover:bg-[#a0d28c] hover:shadow-[#8ec07c]/25"
				target="_blank"
				rel="noopener noreferrer"
				onclick={handleSignupClick}
			>
				Start your first audit
				<ArrowRight class="h-5 w-5 transition-transform group-hover:translate-x-0.5" />
			</a>
		</div>

		<p class="mt-6 text-sm text-white/30">
			No credit card required. No black box. Just clearer decisions.
		</p>
	</div>
</section>
