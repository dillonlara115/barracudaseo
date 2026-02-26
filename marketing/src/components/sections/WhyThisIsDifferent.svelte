<script lang="ts">
	import { Eye, Database, EyeOff, ShieldAlert, PenTool, BarChart3 } from '@lucide/svelte';

	const features = [
		{
			icon: Eye,
			title: 'Transparent reasoning',
			description: 'See exactly why each issue was prioritized. No black-box scores.',
			span: 'md:col-span-1'
		},
		{
			icon: Database,
			title: 'Data-informed decisions',
			description: 'Every recommendation is grounded in your real crawl data and GSC metrics.',
			span: 'md:col-span-1'
		},
		{
			icon: EyeOff,
			title: 'What was deprioritized',
			description: 'Know what was intentionally set aside and why — so you can defend the choice.',
			span: 'md:col-span-1'
		},
		{
			icon: ShieldAlert,
			title: 'Risk assessment',
			description: 'Understand the cost of inaction for every issue in your audit.',
			span: 'md:col-span-1'
		},
		{
			icon: PenTool,
			title: 'AI content intelligence',
			description:
				'Content briefs, article drafts, and internal link suggestions — all grounded in your site data.',
			span: 'md:col-span-1'
		},
		{
			icon: BarChart3,
			title: 'GSC-powered insights',
			description:
				'Quick wins, declining pages, and content gaps surfaced automatically from Search Console.',
			span: 'md:col-span-1'
		}
	];
</script>

<section class="bg-[#222] px-4 py-24">
	<div class="container mx-auto max-w-6xl">
		<div class="mb-16 text-center">
			<p class="mb-4 text-sm font-medium tracking-widest text-[#8ec07c] uppercase">Why Barracuda</p>
			<h2 class="mb-4 font-heading text-3xl font-bold text-white md:text-4xl lg:text-5xl">
				We don't hide the logic. We show it.
			</h2>
			<p class="mx-auto max-w-2xl text-lg text-white/50">
				No mysterious scores. No black-box recommendations. Full transparency into every decision.
			</p>
		</div>

		<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
			{#each features as feature}
				<div
					class="group rounded-2xl border border-white/5 bg-[#282828] p-6 transition-all duration-300 hover:border-[#8ec07c]/20 {feature.span}"
				>
					<div
						class="mb-4 flex h-10 w-10 items-center justify-center rounded-xl bg-[#8ec07c]/10 transition-colors group-hover:bg-[#8ec07c]/20"
					>
						<svelte:component this={feature.icon} class="h-5 w-5 text-[#8ec07c]" />
					</div>
					<h3 class="mb-2 font-heading text-lg font-bold text-white">{feature.title}</h3>
					<p class="text-sm leading-relaxed text-white/50">{feature.description}</p>
				</div>
			{/each}
		</div>
	</div>
</section>
