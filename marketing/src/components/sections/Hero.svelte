<script lang="ts">
	import { trackSignup } from '$lib/analytics';

	function handleSignupClick() {
		trackSignup({
			source: 'Run a smarter site audit',
			location: 'hero'
		});
	}
</script>

<section class="min-h-[60vh] flex items-center justify-center bg-gradient-to-b from-[#3c3836] via-[#8ec07c]/10 to-[#3c3836] py-16 px-4">
	<div class="text-center max-w-4xl mx-auto">
		<h1 class="mb-6 text-4xl md:text-6xl font-heading font-bold text-white leading-tight">
			When SEO data gets overwhelming, Barracuda tells you what actually matters.
		</h1>
		<p class="mb-10 text-xl md:text-2xl text-white/80 font-body max-w-3xl mx-auto leading-relaxed">
			BarracudaSEO turns complex site audits into clear, defensible priorities you can explain with confidence — without drowning in SEO noise.
		</p>
		<div class="flex flex-col sm:flex-row gap-4 justify-center mb-8">
			<a
				href="https://app.barracudaseo.com"
				class="bg-[#8ec07c] hover:bg-[#a0d28c] text-[#3c3836] px-8 py-4 rounded-lg font-medium text-lg transition-colors shadow-lg hover:shadow-[#8ec07c]/20"
				target="_blank"
				rel="noopener noreferrer"
				onclick={handleSignupClick}
			>
				Run a smarter site audit
			</a>
			<a href="#how-it-works" class="border-2 border-white/20 hover:border-[#8ec07c] text-white px-8 py-4 rounded-lg font-medium text-lg transition-colors">
				See how decisions are made
			</a>
		</div>
	</div>
</section>
