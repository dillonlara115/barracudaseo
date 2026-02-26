<script lang="ts">
	import { Upload, Search, Lightbulb, CheckCircle } from '@lucide/svelte';

	const steps = [
		{
			icon: Upload,
			number: '01',
			title: 'Ingest your crawl data',
			description:
				'Run a site crawl and bring the results into Barracuda. Connect GSC for even deeper insights.'
		},
		{
			icon: Search,
			number: '02',
			title: 'See what matters now',
			description:
				'AI analyzes issues, context, and impact to surface the problems that deserve attention right now.'
		},
		{
			icon: Lightbulb,
			number: '03',
			title: 'Understand the "why"',
			description:
				'Every priority comes with clear reasoning — so you know why it matters and can explain it to others.'
		},
		{
			icon: CheckCircle,
			number: '04',
			title: 'Act with confidence',
			description:
				"Fix what matters, ignore what doesn't, and walk into stakeholder conversations prepared."
		}
	];
</script>

<section id="how-it-works" class="bg-[#1d1d1d] px-4 py-24">
	<div class="container mx-auto max-w-6xl">
		<div class="mb-16 text-center">
			<p class="mb-4 text-sm font-medium tracking-widest text-[#8ec07c] uppercase">How it works</p>
			<h2 class="font-heading text-3xl font-bold text-white md:text-4xl lg:text-5xl">
				From crawl to clarity in minutes
			</h2>
		</div>

		<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-4">
			{#each steps as step, i}
				<div class="group relative">
					<div
						class="h-full rounded-2xl border border-white/5 bg-[#282828] p-6 transition-all duration-300 hover:border-[#8ec07c]/20"
					>
						<span
							class="font-heading text-5xl font-bold text-white/5 transition-colors group-hover:text-[#8ec07c]/10"
							>{step.number}</span
						>
						<div
							class="mt-4 mb-4 flex h-10 w-10 items-center justify-center rounded-xl bg-[#8ec07c]/10"
						>
							<svelte:component this={step.icon} class="h-5 w-5 text-[#8ec07c]" />
						</div>
						<h3 class="mb-2 font-heading text-lg font-bold text-white">{step.title}</h3>
						<p class="text-sm leading-relaxed text-white/50">{step.description}</p>
					</div>
					<!-- Connector line (hidden on last) -->
					{#if i < steps.length - 1}
						<div
							class="absolute top-1/2 -right-3 hidden w-6 border-t border-dashed border-white/10 lg:block"
						></div>
					{/if}
				</div>
			{/each}
		</div>
	</div>
</section>
