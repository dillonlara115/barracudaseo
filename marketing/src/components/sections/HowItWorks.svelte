<script lang="ts">
	import { Upload, Search, Lightbulb, CheckCircle } from '@lucide/svelte';

	const steps = [
		{
			icon: Upload,
			number: '01',
			title: 'Ingest your crawl data',
			description: 'Run a site crawl and bring the results into Barracuda. Connect GSC for even deeper insights.'
		},
		{
			icon: Search,
			number: '02',
			title: 'See what matters now',
			description: 'AI analyzes issues, context, and impact to surface the problems that deserve attention right now.'
		},
		{
			icon: Lightbulb,
			number: '03',
			title: 'Understand the "why"',
			description: 'Every priority comes with clear reasoning — so you know why it matters and can explain it to others.'
		},
		{
			icon: CheckCircle,
			number: '04',
			title: 'Act with confidence',
			description: 'Fix what matters, ignore what doesn\'t, and walk into stakeholder conversations prepared.'
		}
	];
</script>

<section id="how-it-works" class="bg-[#1d1d1d] py-24 px-4">
	<div class="container mx-auto max-w-6xl">
		<div class="text-center mb-16">
			<p class="text-sm uppercase tracking-widest text-[#8ec07c] font-medium mb-4">How it works</p>
			<h2 class="text-3xl md:text-4xl lg:text-5xl font-heading font-bold text-white">
				From crawl to clarity in minutes
			</h2>
		</div>

		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
			{#each steps as step, i}
				<div class="group relative">
					<div class="bg-[#282828] border border-white/5 rounded-2xl p-6 h-full hover:border-[#8ec07c]/20 transition-all duration-300">
						<span class="text-5xl font-heading font-bold text-white/5 group-hover:text-[#8ec07c]/10 transition-colors">{step.number}</span>
						<div class="mt-4 mb-4 w-10 h-10 rounded-xl bg-[#8ec07c]/10 flex items-center justify-center">
							<svelte:component this={step.icon} class="w-5 h-5 text-[#8ec07c]" />
						</div>
						<h3 class="text-lg font-heading font-bold text-white mb-2">{step.title}</h3>
						<p class="text-sm text-white/50 leading-relaxed">{step.description}</p>
					</div>
					<!-- Connector line (hidden on last) -->
					{#if i < steps.length - 1}
						<div class="hidden lg:block absolute top-1/2 -right-3 w-6 border-t border-dashed border-white/10"></div>
					{/if}
				</div>
			{/each}
		</div>
	</div>
</section>
