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
			description: "Know what was intentionally set aside and why — so you can defend the choice.",
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
			description: 'Content briefs, article drafts, and internal link suggestions — all grounded in your site data.',
			span: 'md:col-span-1'
		},
		{
			icon: BarChart3,
			title: 'GSC-powered insights',
			description: 'Quick wins, declining pages, and content gaps surfaced automatically from Search Console.',
			span: 'md:col-span-1'
		}
	];
</script>

<section class="bg-[#222] py-24 px-4">
	<div class="container mx-auto max-w-6xl">
		<div class="text-center mb-16">
			<p class="text-sm uppercase tracking-widest text-[#8ec07c] font-medium mb-4">Why Barracuda</p>
			<h2 class="text-3xl md:text-4xl lg:text-5xl font-heading font-bold text-white mb-4">
				We don't hide the logic. We show it.
			</h2>
			<p class="text-lg text-white/50 max-w-2xl mx-auto">
				No mysterious scores. No black-box recommendations. Full transparency into every decision.
			</p>
		</div>

		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
			{#each features as feature}
				<div class="group bg-[#282828] border border-white/5 rounded-2xl p-6 hover:border-[#8ec07c]/20 transition-all duration-300 {feature.span}">
					<div class="w-10 h-10 rounded-xl bg-[#8ec07c]/10 flex items-center justify-center mb-4 group-hover:bg-[#8ec07c]/20 transition-colors">
						<svelte:component this={feature.icon} class="w-5 h-5 text-[#8ec07c]" />
					</div>
					<h3 class="text-lg font-heading font-bold text-white mb-2">{feature.title}</h3>
					<p class="text-sm text-white/50 leading-relaxed">{feature.description}</p>
				</div>
			{/each}
		</div>
	</div>
</section>
