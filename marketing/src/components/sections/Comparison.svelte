<script lang="ts">
	import { X, Check } from '@lucide/svelte';

	const comparisons = [
		{ before: 'Spreadsheets and gut feel', after: 'Clear, explainable priorities' },
		{ before: 'Endless issue lists', after: 'Focus on what actually matters' },
		{ before: 'Tool overload', after: 'One place to decide, then act' },
		{ before: 'Audit anxiety', after: 'Confidence and clarity' }
	];
</script>

<section class="bg-[#1d1d1d] px-4 py-24">
	<div class="container mx-auto max-w-5xl">
		<div class="mb-16 text-center">
			<p class="mb-4 text-sm font-medium tracking-widest text-[#d79921] uppercase">
				Before & after
			</p>
			<h2 class="mb-4 font-heading text-3xl font-bold text-white md:text-4xl lg:text-5xl">
				What Barracuda replaces
			</h2>
			<p class="mx-auto max-w-2xl text-lg text-white/50">
				You're not replacing tools like SEMrush or Ahrefs. You're replacing the stress of
				interpreting them alone.
			</p>
			<p class="mx-auto mt-3 max-w-2xl text-sm text-white/40">
				Weighing your options? See how Barracuda compares to
				<a
					href="/blog/screaming-frog-vs-barracuda"
					class="text-[#8ec07c]/80 hover:text-[#8ec07c] hover:underline">Screaming Frog</a
				>
				and
				<a
					href="/blog/semrush-vs-barracuda"
					class="text-[#8ec07c]/80 hover:text-[#8ec07c] hover:underline">Semrush</a
				>, or browse the
				<a
					href="/blog/alternatives-to-ahrefs"
					class="text-[#8ec07c]/80 hover:text-[#8ec07c] hover:underline">Ahrefs alternatives</a
				>
				and
				<a
					href="/blog/best-seo-crawler-tools"
					class="text-[#8ec07c]/80 hover:text-[#8ec07c] hover:underline">SEO crawler comparison</a
				>.
			</p>
		</div>

		<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
			{#each comparisons as item}
				<div
					class="flex items-stretch gap-4 overflow-hidden rounded-2xl border border-white/5 bg-[#282828] p-5"
				>
					<div class="flex flex-1 items-start gap-3">
						<div
							class="mt-0.5 flex h-7 w-7 shrink-0 items-center justify-center rounded-lg bg-[#cc241d]/10"
						>
							<X class="h-4 w-4 text-[#cc241d]/60" />
						</div>
						<div>
							<p class="mb-1 text-xs tracking-wider text-white/30 uppercase">Instead of</p>
							<p class="text-sm text-white/60">{item.before}</p>
						</div>
					</div>
					<div class="w-px shrink-0 bg-white/5"></div>
					<div class="flex flex-1 items-start gap-3">
						<div
							class="mt-0.5 flex h-7 w-7 shrink-0 items-center justify-center rounded-lg bg-[#8ec07c]/10"
						>
							<Check class="h-4 w-4 text-[#8ec07c]" />
						</div>
						<div>
							<p class="mb-1 text-xs tracking-wider text-[#8ec07c]/60 uppercase">You get</p>
							<p class="text-sm font-medium text-white">{item.after}</p>
						</div>
					</div>
				</div>
			{/each}
		</div>
	</div>
</section>
