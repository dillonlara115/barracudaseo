<script lang="ts">
	import { onMount } from 'svelte';
	import Pricing from '../../components/sections/Pricing.svelte';
	import MetaTags from '../../components/MetaTags.svelte';
	import { getMetaTags, getBreadcrumbSchema } from '$lib/meta';
	import { trackPricingAction } from '$lib/analytics';
	import { CircleQuestionMark } from '@lucide/svelte';

	const meta = getMetaTags({
		title: 'Pricing — Pay for clarity, not complexity',
		description: 'Most SEO tools charge more for complexity. Barracuda charges for clarity. Simple, transparent pricing for decision support.'
	});

	const structuredData = getBreadcrumbSchema([
		{ name: 'Home', url: '/' },
		{ name: 'Pricing', url: '/pricing' }
	]);

	// Track pricing page view
	onMount(() => {
		trackPricingAction('view');
	});

	let open1 = $state(false);
	let open2 = $state(false);
	let open3 = $state(false);
	let open4 = $state(false);

	function toggleItem(index: number) {
		switch (index) {
			case 1: open1 = !open1; break;
			case 2: open2 = !open2; break;
			case 3: open3 = !open3; break;
			case 4: open4 = !open4; break;
		}
	}
</script>

<MetaTags config={{ ...meta, structuredData }} />

<!-- Pricing Hero -->
<section class="py-20 bg-[#3c3836]">
	<div class="container mx-auto px-4 max-w-4xl text-center">
		<h1 class="text-5xl md:text-6xl font-heading font-bold mb-8 text-white">Pricing</h1>
		
		<h2 class="text-3xl md:text-4xl font-bold text-white mb-8">Pay for clarity, not complexity.</h2>
		
		<div class="text-xl text-white/80 leading-relaxed max-w-3xl mx-auto space-y-6">
			<p>
				Most SEO tools charge more as they add more data, more dashboards, and more things you’re expected to interpret on your own.
			</p>
			<p>
				BarracudaSEO takes a different approach.
			</p>
			<p class="text-white font-medium text-2xl">
				We price based on the value we deliver at the moment it matters most: helping you decide what to fix, what to ignore, and how to explain it with confidence.
			</p>
		</div>
	</div>
</section>

<!-- Why Pricing Works This Way -->
<section class="py-16 bg-[#282828]">
	<div class="container mx-auto px-4 max-w-4xl">
		<div class="text-center mb-12">
			<h2 class="text-3xl font-heading font-bold mb-6 text-[#8ec07c]">Why our pricing works this way</h2>
			<p class="text-xl text-white/80 max-w-3xl mx-auto">
				Barracuda isn’t about unlocking more metrics.
				<br />
				It’s about reducing the mental and emotional load that comes with SEO audits.
			</p>
		</div>

		<div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
			<div class="bg-[#3c3836] p-6 rounded-lg border border-white/5">
				<p class="text-white/90 font-medium">Help you move from crawl data to clear priorities</p>
			</div>
			<div class="bg-[#3c3836] p-6 rounded-lg border border-white/5">
				<p class="text-white/90 font-medium">Reduce second-guessing and audit paralysis</p>
			</div>
			<div class="bg-[#3c3836] p-6 rounded-lg border border-white/5">
				<p class="text-white/90 font-medium">Support confident decisions under real-world pressure</p>
			</div>
		</div>

		<div class="bg-[#3c3836] p-8 rounded-xl border-l-4 border-[#8ec07c] text-center md:text-left">
			<p class="text-white/80 text-lg">
				As you move up plans, you’re not buying “more noise.”
				<br />
				<span class="text-white font-bold">You’re getting more support for better decisions — faster, clearer, and easier to defend.</span>
			</p>
		</div>
	</div>
</section>

<!-- Choose Clarity Level -->
<section class="py-16 bg-[#3c3836]">
	<div class="container mx-auto px-4 max-w-4xl text-center">
		<h2 class="text-3xl font-heading font-bold mb-6 text-white">Choose the level of clarity you need</h2>
		<p class="text-xl text-white/80 max-w-3xl mx-auto mb-8">
			Whether you’re running the occasional audit or managing multiple sites and clients, there’s a plan built around how much <strong>decision support</strong> you need — not how many features you can unlock.
		</p>
		<p class="text-lg text-white/60 italic">
			Start where you are. Upgrade when your responsibility (and accountability) grows.
		</p>
	</div>
</section>

<!-- Pricing Component -->
<Pricing />

<!-- One Quiet Promise -->
<section class="py-20 bg-[#282828]">
	<div class="container mx-auto px-4 max-w-3xl text-center">
		<h2 class="text-3xl font-heading font-bold mb-8 text-[#8ec07c]">One quiet promise</h2>
		
		<div class="flex flex-col md:flex-row justify-center gap-8 mb-12 text-white/70 text-lg">
			<div class="flex items-center justify-center">
				<span class="text-[#cc241d] mr-2 font-bold">×</span> No black-box scoring
			</div>
			<div class="flex items-center justify-center">
				<span class="text-[#cc241d] mr-2 font-bold">×</span> No feature bloat
			</div>
			<div class="flex items-center justify-center">
				<span class="text-[#cc241d] mr-2 font-bold">×</span> No paying for overwhelm
			</div>
		</div>

		<p class="text-white font-medium text-2xl mb-12">
			Every plan is built around the same core belief:
			<br />
			<span class="text-white font-bold mt-2 block">SEO should be easier to decide, not harder to understand.</span>
		</p>

		<p class="text-white/60 text-lg">
			Below, you’ll find plans designed to match how you work — from focused audits to ongoing, multi-site decision support.
		</p>
	</div>
</section>

<section class="py-20 bg-[#3c3836]">
	<div class="container mx-auto px-4">
		<div class="max-w-3xl mx-auto">
			<h2 class="text-3xl font-heading font-bold mb-8 text-white text-center">Frequently Asked Questions</h2>
			
			<div class="space-y-4">
				<!-- FAQ 1 -->
				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<button
						class="w-full px-6 py-4 text-left flex items-center justify-between hover:bg-[#3c3836] transition-colors"
						onclick={() => toggleItem(1)}
					>
						<h3 class="text-lg font-heading font-bold text-white pr-4">
							What am I actually paying for when I upgrade?
						</h3>
						<svg
							class="w-5 h-5 text-[#8ec07c] flex-shrink-0 transition-transform {open1 ? 'rotate-180' : ''}"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
						</svg>
					</button>
					{#if open1}
						<div class="px-6 py-4 text-white/80 text-sm space-y-4">
							<p>You’re paying for clarity at scale.</p>
							<p>As you move to paid plans, you’re not unlocking more noise or extra dashboards. You’re getting:</p>
							<ul class="list-disc list-inside space-y-1 ml-2">
								<li>Deeper prioritization and recommendations</li>
								<li>Cloud-based projects and history</li>
								<li>The ability to collaborate and explain decisions across teams or clients</li>
							</ul>
							<p>Every plan is designed to reduce the mental overhead of SEO audits — especially when accountability is high.</p>
						</div>
					{/if}
				</div>

				<!-- FAQ 2 -->
				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<button
						class="w-full px-6 py-4 text-left flex items-center justify-between hover:bg-[#3c3836] transition-colors"
						onclick={() => toggleItem(2)}
					>
						<h3 class="text-lg font-heading font-bold text-white pr-4">
							How does team collaboration work?
						</h3>
						<svg
							class="w-5 h-5 text-[#8ec07c] flex-shrink-0 transition-transform {open2 ? 'rotate-180' : ''}"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
						</svg>
					</button>
					{#if open2}
						<div class="px-6 py-4 text-white/80 text-sm space-y-4">
							<p>Barracuda is built for shared decision-making, not just shared access.</p>
							<p>Pro includes one user, with the option to add teammates for $5/user/month. Each team member can see the same prioritized issues, reasoning, and context — so everyone stays aligned on what matters and why.</p>
							<p>This works especially well for agencies and in-house teams that need consistent explanations across people and projects.</p>
						</div>
					{/if}
				</div>

				<!-- FAQ 3 -->
				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<button
						class="w-full px-6 py-4 text-left flex items-center justify-between hover:bg-[#3c3836] transition-colors"
						onclick={() => toggleItem(3)}
					>
						<h3 class="text-lg font-heading font-bold text-white pr-4">
							What are the different team roles?
						</h3>
						<svg
							class="w-5 h-5 text-[#8ec07c] flex-shrink-0 transition-transform {open3 ? 'rotate-180' : ''}"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
						</svg>
					</button>
					{#if open3}
						<div class="px-6 py-4 text-white/80 text-sm space-y-4">
							<p>Roles are designed around responsibility and trust:</p>
							<ul class="list-disc list-inside space-y-1 ml-2">
								<li><strong>Owner:</strong> Full control, including billing, invitations, and project management</li>
								<li><strong>Editor:</strong> Can run crawls, review priorities, manage issues, and see recommendations</li>
								<li><strong>Viewer:</strong> Read-only access for stakeholders who need visibility without risk</li>
							</ul>
							<p>The goal is clarity without chaos — everyone sees the same story, with the right level of control.</p>
						</div>
					{/if}
				</div>

				<!-- FAQ 4 -->
				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<button
						class="w-full px-6 py-4 text-left flex items-center justify-between hover:bg-[#3c3836] transition-colors"
						onclick={() => toggleItem(4)}
					>
						<h3 class="text-lg font-heading font-bold text-white pr-4">
							Do I need a Team plan for multiple users?
						</h3>
						<svg
							class="w-5 h-5 text-[#8ec07c] flex-shrink-0 transition-transform {open4 ? 'rotate-180' : ''}"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
						</svg>
					</button>
					{#if open4}
						<div class="px-6 py-4 text-white/80 text-sm space-y-4">
							<p>Not necessarily.</p>
							<p>The Pro plan includes one user and lets you add additional team members as needed. This keeps things flexible as your responsibility grows.</p>
							<p>If you’re managing larger teams or have specific requirements, we offer custom Team plans to match how you actually work.</p>
						</div>
					{/if}
				</div>
			</div>
		</div>
	</div>
</section>
