<script lang="ts">
	import { X, Check } from '@lucide/svelte';

	const comparisons = [
		{ before: 'Spreadsheets and gut feel', after: 'Clear, explainable priorities' },
		{ before: 'Endless issue lists', after: 'Focus on what actually matters' },
		{ before: 'Tool overload', after: 'One place to decide, then act' },
		{ before: 'Audit anxiety', after: 'Confidence and clarity' }
	];
</script>

<section class="bg-[#1d1d1d] py-24 px-4">
	<div class="container mx-auto max-w-5xl">
		<div class="text-center mb-16">
			<p class="text-sm uppercase tracking-widest text-[#d79921] font-medium mb-4">Before & after</p>
			<h2 class="text-3xl md:text-4xl lg:text-5xl font-heading font-bold text-white mb-4">
				What Barracuda replaces
			</h2>
			<p class="text-lg text-white/50 max-w-2xl mx-auto">
				You're not replacing tools like SEMrush or Ahrefs. You're replacing the stress of interpreting them alone.
			</p>
		</div>

		<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
			{#each comparisons as item}
				<div class="flex items-stretch gap-4 bg-[#282828] border border-white/5 rounded-2xl p-5 overflow-hidden">
					<div class="flex-1 flex items-start gap-3">
						<div class="w-7 h-7 rounded-lg bg-[#cc241d]/10 flex items-center justify-center shrink-0 mt-0.5">
							<X class="w-4 h-4 text-[#cc241d]/60" />
						</div>
						<div>
							<p class="text-xs text-white/30 uppercase tracking-wider mb-1">Instead of</p>
							<p class="text-sm text-white/60">{item.before}</p>
						</div>
					</div>
					<div class="w-px bg-white/5 shrink-0"></div>
					<div class="flex-1 flex items-start gap-3">
						<div class="w-7 h-7 rounded-lg bg-[#8ec07c]/10 flex items-center justify-center shrink-0 mt-0.5">
							<Check class="w-4 h-4 text-[#8ec07c]" />
						</div>
						<div>
							<p class="text-xs text-[#8ec07c]/60 uppercase tracking-wider mb-1">You get</p>
							<p class="text-sm text-white font-medium">{item.after}</p>
						</div>
					</div>
				</div>
			{/each}
		</div>
	</div>
</section>
