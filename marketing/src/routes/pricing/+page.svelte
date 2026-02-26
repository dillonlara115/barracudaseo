<script lang="ts">
	import { onMount } from 'svelte';
	import Pricing from '../../components/sections/Pricing.svelte';
	import MetaTags from '../../components/MetaTags.svelte';
	import { getMetaTags, getBreadcrumbSchema, getFAQPageSchema } from '$lib/meta';
	import { trackPricingAction } from '$lib/analytics';
	import { CircleQuestionMark } from '@lucide/svelte';

	const meta = getMetaTags({
		title: 'Pricing — Pay for clarity, not complexity',
		description:
			'Most SEO tools charge more for complexity. Barracuda charges for clarity. Simple, transparent pricing for decision support.'
	});

	const faqs = [
		{
			question: 'What am I actually paying for when I upgrade?',
			answer:
				"You're paying for clarity at scale. As you move to paid plans, you're not unlocking more noise or extra dashboards. You're getting: deeper prioritization and recommendations, cloud-based projects and history, and the ability to collaborate and explain decisions across teams or clients. Every plan is designed to reduce the mental overhead of SEO audits — especially when accountability is high."
		},
		{
			question: 'How does team collaboration work?',
			answer:
				'Barracuda is built for shared decision-making, not just shared access. Pro includes one user, with the option to add teammates for $5/user/month. Each team member can see the same prioritized issues, reasoning, and context — so everyone stays aligned on what matters and why. This works especially well for agencies and in-house teams that need consistent explanations across people and projects.'
		},
		{
			question: 'What are the different team roles?',
			answer:
				'Roles are designed around responsibility and trust: Owner has full control including billing, invitations, and project management; Editor can run crawls, review priorities, manage issues, and see recommendations; Viewer has read-only access for stakeholders who need visibility without risk. The goal is clarity without chaos — everyone sees the same story, with the right level of control.'
		},
		{
			question: 'Do I need a Team plan for multiple users?',
			answer:
				"Not necessarily. The Pro plan includes one user and lets you add additional team members as needed. This keeps things flexible as your responsibility grows. If you're managing larger teams or have specific requirements, we offer custom Team plans to match how you actually work."
		}
	];

	const structuredData = [
		getBreadcrumbSchema([
			{ name: 'Home', url: '/' },
			{ name: 'Pricing', url: '/pricing' }
		]),
		getFAQPageSchema(faqs)
	];

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
			case 1:
				open1 = !open1;
				break;
			case 2:
				open2 = !open2;
				break;
			case 3:
				open3 = !open3;
				break;
			case 4:
				open4 = !open4;
				break;
		}
	}
</script>

<MetaTags config={{ ...meta, structuredData }} />

<!-- Pricing Hero -->
<section class="bg-[#3c3836] py-20">
	<div class="container mx-auto max-w-4xl px-4 text-center">
		<h1 class="mb-8 font-heading text-5xl font-bold text-white md:text-6xl">Pricing</h1>

		<h2 class="mb-8 text-3xl font-bold text-white md:text-4xl">Pay for clarity, not complexity.</h2>

		<div class="mx-auto max-w-3xl space-y-6 text-xl leading-relaxed text-white/80">
			<p>
				Most SEO tools charge more as they add more data, more dashboards, and more things you’re
				expected to interpret on your own.
			</p>
			<p>BarracudaSEO takes a different approach.</p>
			<p class="text-2xl font-medium text-white">
				We price based on the value we deliver at the moment it matters most: helping you decide
				what to fix, what to ignore, and how to explain it with confidence.
			</p>
		</div>
	</div>
</section>

<!-- Why Pricing Works This Way -->
<section class="bg-[#282828] py-16">
	<div class="container mx-auto max-w-4xl px-4">
		<div class="mb-12 text-center">
			<h2 class="mb-6 font-heading text-3xl font-bold text-[#8ec07c]">
				Why our pricing works this way
			</h2>
			<p class="mx-auto max-w-3xl text-xl text-white/80">
				Barracuda isn’t about unlocking more metrics.
				<br />
				It’s about reducing the mental and emotional load that comes with SEO audits.
			</p>
		</div>

		<div class="mb-12 grid grid-cols-1 gap-6 md:grid-cols-3">
			<div class="rounded-lg border border-white/5 bg-[#3c3836] p-6">
				<p class="font-medium text-white/90">Help you move from crawl data to clear priorities</p>
			</div>
			<div class="rounded-lg border border-white/5 bg-[#3c3836] p-6">
				<p class="font-medium text-white/90">Reduce second-guessing and audit paralysis</p>
			</div>
			<div class="rounded-lg border border-white/5 bg-[#3c3836] p-6">
				<p class="font-medium text-white/90">
					Support confident decisions under real-world pressure
				</p>
			</div>
		</div>

		<div class="rounded-xl border-l-4 border-[#8ec07c] bg-[#3c3836] p-8 text-center md:text-left">
			<p class="text-lg text-white/80">
				As you move up plans, you’re not buying “more noise.”
				<br />
				<span class="font-bold text-white"
					>You’re getting more support for better decisions — faster, clearer, and easier to defend.</span
				>
			</p>
		</div>
	</div>
</section>

<!-- Choose Clarity Level -->
<section class="bg-[#3c3836] py-16">
	<div class="container mx-auto max-w-4xl px-4 text-center">
		<h2 class="mb-6 font-heading text-3xl font-bold text-white">
			Choose the level of clarity you need
		</h2>
		<p class="mx-auto mb-8 max-w-3xl text-xl text-white/80">
			Whether you’re running the occasional audit or managing multiple sites and clients, there’s a
			plan built around how much <strong>decision support</strong> you need — not how many features you
			can unlock.
		</p>
		<p class="text-lg text-white/60 italic">
			Start where you are. Upgrade when your responsibility (and accountability) grows.
		</p>
	</div>
</section>

<!-- Pricing Component -->
<Pricing />

<!-- One Quiet Promise -->
<section class="bg-[#282828] py-20">
	<div class="container mx-auto max-w-3xl px-4 text-center">
		<h2 class="mb-8 font-heading text-3xl font-bold text-[#8ec07c]">One quiet promise</h2>

		<div class="mb-12 flex flex-col justify-center gap-8 text-lg text-white/70 md:flex-row">
			<div class="flex items-center justify-center">
				<span class="mr-2 font-bold text-[#cc241d]">×</span> No black-box scoring
			</div>
			<div class="flex items-center justify-center">
				<span class="mr-2 font-bold text-[#cc241d]">×</span> No feature bloat
			</div>
			<div class="flex items-center justify-center">
				<span class="mr-2 font-bold text-[#cc241d]">×</span> No paying for overwhelm
			</div>
		</div>

		<p class="mb-12 text-2xl font-medium text-white">
			Every plan is built around the same core belief:
			<br />
			<span class="mt-2 block font-bold text-white"
				>SEO should be easier to decide, not harder to understand.</span
			>
		</p>

		<p class="text-lg text-white/60">
			Below, you’ll find plans designed to match how you work — from focused audits to ongoing,
			multi-site decision support.
		</p>
	</div>
</section>

<section class="bg-[#3c3836] py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-3xl">
			<h2 class="mb-8 text-center font-heading text-3xl font-bold text-white">
				Frequently Asked Questions
			</h2>

			<div class="space-y-4">
				<!-- FAQ 1 -->
				<div class="overflow-hidden rounded-lg border border-white/10 bg-[#2d2826]">
					<button
						class="flex w-full items-center justify-between px-6 py-4 text-left transition-colors hover:bg-[#3c3836]"
						onclick={() => toggleItem(1)}
					>
						<h3 class="pr-4 font-heading text-lg font-bold text-white">
							What am I actually paying for when I upgrade?
						</h3>
						<svg
							class="h-5 w-5 flex-shrink-0 text-[#8ec07c] transition-transform {open1
								? 'rotate-180'
								: ''}"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M19 9l-7 7-7-7"
							/>
						</svg>
					</button>
					{#if open1}
						<div class="space-y-4 px-6 py-4 text-sm text-white/80">
							<p>You’re paying for clarity at scale.</p>
							<p>
								As you move to paid plans, you’re not unlocking more noise or extra dashboards.
								You’re getting:
							</p>
							<ul class="ml-2 list-inside list-disc space-y-1">
								<li>
									Deeper <a
										href="/blog/how-to-prioritize-seo-issues"
										class="text-[#8ec07c] underline hover:text-[#a0d28c]">prioritization</a
									> and recommendations
								</li>
								<li>Cloud-based projects and history</li>
								<li>The ability to collaborate and explain decisions across teams or clients</li>
							</ul>
							<p>
								Every plan is designed to reduce the mental overhead of <a
									href="/blog/complete-technical-seo-audit-guide"
									class="text-[#8ec07c] underline hover:text-[#a0d28c]">SEO audits</a
								> — especially when accountability is high.
							</p>
						</div>
					{/if}
				</div>

				<!-- FAQ 2 -->
				<div class="overflow-hidden rounded-lg border border-white/10 bg-[#2d2826]">
					<button
						class="flex w-full items-center justify-between px-6 py-4 text-left transition-colors hover:bg-[#3c3836]"
						onclick={() => toggleItem(2)}
					>
						<h3 class="pr-4 font-heading text-lg font-bold text-white">
							How does team collaboration work?
						</h3>
						<svg
							class="h-5 w-5 flex-shrink-0 text-[#8ec07c] transition-transform {open2
								? 'rotate-180'
								: ''}"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M19 9l-7 7-7-7"
							/>
						</svg>
					</button>
					{#if open2}
						<div class="space-y-4 px-6 py-4 text-sm text-white/80">
							<p>Barracuda is built for shared decision-making, not just shared access.</p>
							<p>
								Pro includes one user, with the option to add teammates for $5/user/month. Each team
								member can see the same prioritized issues, reasoning, and context — so everyone
								stays aligned on what matters and why.
							</p>
							<p>
								This works especially well for agencies and in-house teams that need consistent
								explanations across people and projects.
							</p>
						</div>
					{/if}
				</div>

				<!-- FAQ 3 -->
				<div class="overflow-hidden rounded-lg border border-white/10 bg-[#2d2826]">
					<button
						class="flex w-full items-center justify-between px-6 py-4 text-left transition-colors hover:bg-[#3c3836]"
						onclick={() => toggleItem(3)}
					>
						<h3 class="pr-4 font-heading text-lg font-bold text-white">
							What are the different team roles?
						</h3>
						<svg
							class="h-5 w-5 flex-shrink-0 text-[#8ec07c] transition-transform {open3
								? 'rotate-180'
								: ''}"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M19 9l-7 7-7-7"
							/>
						</svg>
					</button>
					{#if open3}
						<div class="space-y-4 px-6 py-4 text-sm text-white/80">
							<p>Roles are designed around responsibility and trust:</p>
							<ul class="ml-2 list-inside list-disc space-y-1">
								<li>
									<strong>Owner:</strong> Full control, including billing, invitations, and project management
								</li>
								<li>
									<strong>Editor:</strong> Can run crawls, review priorities, manage issues, and see recommendations
								</li>
								<li>
									<strong>Viewer:</strong> Read-only access for stakeholders who need visibility without
									risk
								</li>
							</ul>
							<p>
								The goal is clarity without chaos — everyone sees the same story, with the right
								level of control.
							</p>
						</div>
					{/if}
				</div>

				<!-- FAQ 4 -->
				<div class="overflow-hidden rounded-lg border border-white/10 bg-[#2d2826]">
					<button
						class="flex w-full items-center justify-between px-6 py-4 text-left transition-colors hover:bg-[#3c3836]"
						onclick={() => toggleItem(4)}
					>
						<h3 class="pr-4 font-heading text-lg font-bold text-white">
							Do I need a Team plan for multiple users?
						</h3>
						<svg
							class="h-5 w-5 flex-shrink-0 text-[#8ec07c] transition-transform {open4
								? 'rotate-180'
								: ''}"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M19 9l-7 7-7-7"
							/>
						</svg>
					</button>
					{#if open4}
						<div class="space-y-4 px-6 py-4 text-sm text-white/80">
							<p>Not necessarily.</p>
							<p>
								The Pro plan includes one user and lets you add additional team members as needed.
								This keeps things flexible as your responsibility grows.
							</p>
							<p>
								If you’re managing larger teams or have specific requirements, we offer custom Team
								plans to match how you actually work.
							</p>
						</div>
					{/if}
				</div>
			</div>
		</div>
	</div>
</section>
