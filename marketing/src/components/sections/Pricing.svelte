<script lang="ts">
	import { trackSignup, trackPricingAction } from '$lib/analytics';

	function handleFreeSignup() {
		trackSignup({
			source: 'Get Started',
			location: 'pricing_free',
			plan: 'free'
		});
	}

	function handleProSignup() {
		trackSignup({
			source: 'Start Free Trial',
			location: 'pricing_pro',
			plan: 'pro'
		});
	}

	function handleContactSales() {
		trackPricingAction('cta_click', {
			cta_name: 'Contact Sales'
		});
	}
</script>

<section class="bg-[#2d2826] py-20">
	<div class="container mx-auto px-4">
		<div class="mb-16 text-center">
			<h2 class="mb-4 font-heading text-4xl font-bold text-white md:text-5xl">Simple Pricing</h2>
			<p class="mx-auto max-w-2xl text-xl text-white/70">Start free, scale with your team</p>
		</div>
		<div class="mx-auto grid max-w-4xl grid-cols-1 gap-8 md:grid-cols-2">
			<!-- Free Tier -->
			<div class="rounded-lg border border-white/10 bg-[#3c3836] p-8 shadow-xl">
				<h3 class="mb-2 font-heading text-2xl font-bold text-white">Free</h3>
				<div class="mb-4 text-4xl font-bold text-white">
					$0
					<span class="text-lg font-normal text-white/70">/month</span>
				</div>
				<ul class="mb-6 space-y-3">
					<li class="flex items-center gap-2 text-white/70">
						<span class="text-[#8ec07c]">✓</span>
						<span>Crawl up to <strong class="text-white">100 pages</strong> per project</span>
					</li>
					<li class="flex items-center gap-2 text-white/70">
						<span class="text-[#8ec07c]">✓</span>
						<span>View detected issues and opportunities</span>
					</li>
					<li class="flex items-center gap-2 text-white/70">
						<span class="text-[#8ec07c]">✓</span>
						<span>Web dashboard access</span>
					</li>
					<li class="flex items-center gap-2 text-white/50">
						<span class="text-white/30">✗</span>
						<span class="text-white/50">No team access</span>
					</li>
					<li class="flex items-center gap-2 text-white/50">
						<span class="text-white/30">✗</span>
						<span class="text-white/50">No recommendations or integrations</span>
					</li>
				</ul>
				<a
					href="https://app.barracudaseo.com"
					class="block w-full rounded-lg border-2 border-white/20 px-6 py-3 text-center font-medium text-white transition-colors hover:border-[#8ec07c]"
					target="_blank"
					rel="noopener noreferrer"
					onclick={handleFreeSignup}
				>
					Get Started
				</a>
				<p class="mt-4 text-center text-sm text-white/50">
					Great for quick technical audits or demos
				</p>
			</div>
			<!-- Pro Tier -->
			<div class="relative rounded-lg border-2 border-[#8ec07c] bg-[#8ec07c] p-8 shadow-xl">
				<div
					class="absolute top-4 right-4 rounded-full bg-white px-3 py-1 text-xs font-bold text-[#3c3836]"
				>
					Popular
				</div>
				<h3 class="mb-2 font-heading text-2xl font-bold text-[#3c3836]">Pro</h3>
				<div class="mb-2">
					<div class="text-4xl font-bold text-[#3c3836]">
						$29
						<span class="text-lg font-normal text-[#3c3836]/80">/month</span>
					</div>
					<p class="mt-1 text-sm text-[#3c3836]/80">+ $5 per additional team member</p>
				</div>
				<ul class="mb-6 space-y-3">
					<li class="flex items-center gap-2 text-[#3c3836]">
						<span>✓</span>
						<span>Crawl up to <strong>10,000 pages</strong></span>
					</li>
					<li class="flex items-center gap-2 text-[#3c3836]">
						<span>✓</span>
						<span
							><strong>Team collaboration</strong> — invite teammates with role-based permissions</span
						>
					</li>
					<li class="flex items-center gap-2 text-[#3c3836]">
						<span>✓</span>
						<span
							><strong>AI-powered recommendations</strong> with contextual fixes and crawl summaries</span
						>
					</li>
					<li class="flex items-center gap-2 text-[#3c3836]">
						<span>✓</span>
						<span><strong>Keyword rank tracking</strong> with automatic daily/weekly checks</span>
					</li>
					<li class="flex items-center gap-2 text-[#3c3836]">
						<span>✓</span>
						<span>Integrations: Google Search Console, Analytics, Clarity, Slack</span>
					</li>
					<li class="flex items-center gap-2 text-[#3c3836]">
						<span>✓</span>
						<span>Public report sharing for clients and stakeholders</span>
					</li>
					<li class="flex items-center gap-2 text-[#3c3836]">
						<span>✓</span>
						<span>Historical comparisons and advanced exports</span>
					</li>
					<li class="flex items-center gap-2 text-[#3c3836]">
						<span>✓</span>
						<span>CLI automation and priority support</span>
					</li>
				</ul>
				<a
					href="https://app.barracudaseo.com"
					class="block w-full rounded-lg bg-white px-6 py-3 text-center font-medium text-[#3c3836] transition-colors hover:bg-white/90"
					target="_blank"
					rel="noopener noreferrer"
					onclick={handleProSignup}
				>
					Start Free Trial
				</a>
				<p class="mt-4 text-center text-sm text-[#3c3836]/80">Perfect for agencies and teams</p>
			</div>
		</div>
		<div class="mx-auto mt-12 max-w-3xl">
			<div class="rounded-lg border border-white/10 bg-[#3c3836] p-6">
				<h3 class="mb-3 text-center font-heading text-xl font-bold text-white">Need more?</h3>
				<p class="mb-4 text-center text-white/70">
					For larger organizations with 5+ team members or custom crawl limits, we offer custom Team
					plans with dedicated support.
				</p>
				<a
					href="https://app.barracudaseo.com"
					class="block w-full rounded-lg border-2 border-[#8ec07c] px-6 py-3 text-center font-medium text-[#8ec07c] transition-colors hover:bg-[#8ec07c]/10"
					target="_blank"
					rel="noopener noreferrer"
					onclick={handleContactSales}
				>
					Contact Sales
				</a>
			</div>
		</div>
	</div>
</section>
