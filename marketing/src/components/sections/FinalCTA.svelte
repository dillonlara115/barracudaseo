<script lang="ts">
	import { trackSignup } from '$lib/analytics';

	function handleSignupClick() {
		trackSignup({
			source: 'Start your first audit',
			location: 'final-cta'
		});
	}
</script>

<section class="py-24 px-4 bg-gradient-to-t from-[#282828] to-[#3c3836]">
	<div class="max-w-4xl mx-auto text-center">
		<h2 class="text-4xl md:text-5xl font-heading font-bold text-white mb-6">
			Stop second-guessing your audits.
		</h2>
		<p class="text-xl md:text-2xl text-white/80 font-body mb-10 max-w-2xl mx-auto">
			BarracudaSEO helps you move from data to decisions — with clarity you can stand behind.
		</p>
		
		<div class="flex flex-col items-center">
			<a
				href="https://app.barracudaseo.com"
				class="bg-[#8ec07c] hover:bg-[#a0d28c] text-[#3c3836] px-10 py-5 rounded-lg font-bold text-xl transition-all shadow-lg hover:shadow-[#8ec07c]/20 hover:-translate-y-1 mb-4"
				target="_blank"
				rel="noopener noreferrer"
				onclick={handleSignupClick}
			>
				Start your first audit
			</a>
			<p class="text-white/50 text-sm">
				No credit card. No black box. Just clearer decisions.
			</p>
		</div>
	</div>
</section>
