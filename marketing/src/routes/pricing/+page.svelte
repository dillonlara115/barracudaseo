<script lang="ts">
	import Pricing from '../../components/sections/Pricing.svelte';
	import MetaTags from '../../components/MetaTags.svelte';
	import { getMetaTags, getBreadcrumbSchema } from '$lib/meta';

	const meta = getMetaTags({
		title: 'Pricing',
		description: 'Start with a free 100-page audit. Upgrade to Pro ($29/mo) for 10,000+ pages, team collaboration with role-based permissions, recommendations, and integrations.'
	});

	const structuredData = getBreadcrumbSchema([
		{ name: 'Home', url: '/' },
		{ name: 'Pricing', url: '/pricing' }
	]);
</script>

<MetaTags config={{ ...meta, structuredData }} />

<div class="py-20 bg-[#3c3836]">
	<div class="container mx-auto px-4">
		<div class="text-center mb-16">
			<h1 class="text-5xl md:text-6xl font-heading font-bold mb-4 text-white">Pricing</h1>
			<p class="text-xl text-white/70 max-w-2xl mx-auto">
				Simple, transparent pricing for teams of all sizes
			</p>
		</div>
	</div>
</div>

<Pricing />

<section class="py-20 bg-[#3c3836]">
	<div class="container mx-auto px-4">
		<div class="max-w-3xl mx-auto text-center">
			<h2 class="text-3xl font-heading font-bold mb-4 text-white">Frequently Asked Questions</h2>
			<div class="space-y-6 mt-8 text-left">
				<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
					<h3 class="text-xl font-heading font-bold mb-3 text-[#8ec07c]">Can I try before I buy?</h3>
					<p class="text-white/70">
						Yes! The Free plan includes all core crawling features. Upgrade to Pro or Team when you need cloud workspace and team collaboration.
					</p>
				</div>
				<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
					<h3 class="text-xl font-heading font-bold mb-3 text-[#8ec07c]">How does team collaboration work?</h3>
					<p class="text-white/70">
						Pro includes 1 user. Add team members for $5/user/month. Each member gets role-based access (Owner, Editor, or Viewer) to shared projects. Perfect for agencies managing multiple clients or in-house SEO teams collaborating on site audits.
					</p>
				</div>
				<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
					<h3 class="text-xl font-heading font-bold mb-3 text-[#8ec07c]">What are the different team roles?</h3>
					<p class="text-white/70">
						<strong>Owner:</strong> Full permissions including billing, invitations, and project deletion. <strong>Editor:</strong> Can trigger crawls, manage issues, and view recommendations. <strong>Viewer:</strong> Read-only access to dashboard and reports.
					</p>
				</div>
				<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
					<h3 class="text-xl font-heading font-bold mb-3 text-[#8ec07c]">Do I need a Team plan for multiple users?</h3>
					<p class="text-white/70">
						No! Pro includes flexible team add-ons. Start with Pro ($29/mo) for 1 user, then add teammates at $5/user/month. Custom Team plans are available for organizations with 5+ members or specific requirements.
					</p>
				</div>
				<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10">
					<h3 class="text-xl font-heading font-bold mb-3 text-[#8ec07c]">Do you offer refunds?</h3>
					<p class="text-white/70">
						Yes, we offer a 30-day money-back guarantee. If you're not satisfied, contact us for a full refund.
					</p>
				</div>
			</div>
		</div>
	</div>
</section>

