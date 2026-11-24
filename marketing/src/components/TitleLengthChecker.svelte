<script lang="ts">
	import { onMount } from 'svelte';

	// Google's desktop title links are set in 20px Arial and cut at about 600px.
	const TITLE_FONT = '20px Arial';
	const DESKTOP_LIMIT_PX = 600;

	type Verdict = { tone: 'good' | 'warn' | 'bad' | 'idle'; label: string; detail: string };

	let title = $state('');
	let context: CanvasRenderingContext2D | null = $state(null);

	const chars = $derived(title.trim().length);
	const pixels = $derived(measure(title.trim()));
	const result = $derived(verdict(pixels, chars));
	const fill = $derived(Math.min(100, Math.round((pixels / DESKTOP_LIMIT_PX) * 100)));

	// Width of the title as Google would draw it. Before the canvas is ready
	// (and on the server) this falls back to an average of 9.5px per character.
	function measure(text: string): number {
		if (!text) return 0;
		if (!context) return Math.round(text.length * 9.5);
		return Math.round(context.measureText(text).width);
	}

	// Thresholds behind the verdict. These are judgement calls: 30 characters is
	// where a title stops describing a page, and the 20px margin under the limit
	// allows for small differences between browsers' text rendering.
	function verdict(px: number, count: number): Verdict {
		if (count === 0) {
			return { tone: 'idle', label: 'Type a title', detail: 'Paste a title tag to measure it.' };
		}
		if (px > DESKTOP_LIMIT_PX) {
			return {
				tone: 'bad',
				label: 'Will be cut off',
				detail: `About ${px - DESKTOP_LIMIT_PX}px over the desktop limit. Move the keyword and the main promise into the first 50 characters, or shorten it.`
			};
		}
		if (px > DESKTOP_LIMIT_PX - 20) {
			return {
				tone: 'warn',
				label: 'Borderline',
				detail:
					'Within 20px of the limit, so it may be cut in some browsers. Trim a word to be safe.'
			};
		}
		if (count < 30) {
			return {
				tone: 'warn',
				label: 'Fits, but short',
				detail:
					'Very short titles describe little, and they are the ones Google rewrites most often.'
			};
		}
		return {
			tone: 'good',
			label: 'Fits',
			detail:
				'Shows in full on desktop. Check that the first 50 characters make sense on their own for mobile.'
		};
	}

	const toneClass: Record<Verdict['tone'], string> = {
		good: 'text-[#8ec07c]',
		warn: 'text-[#fabd2f]',
		bad: 'text-[#fb4934]',
		idle: 'text-white/60'
	};
	const barClass: Record<Verdict['tone'], string> = {
		good: 'bg-[#8ec07c]',
		warn: 'bg-[#fabd2f]',
		bad: 'bg-[#fb4934]',
		idle: 'bg-white/20'
	};

	onMount(() => {
		const ctx = document.createElement('canvas').getContext('2d');
		if (ctx) {
			ctx.font = TITLE_FONT;
			context = ctx;
		}
	});
</script>

<div class="my-8 rounded-lg border border-white/10 bg-[#3c3836] p-6">
	<label for="title-length-input" class="mb-2 block font-bold text-white">
		Title tag length checker
	</label>
	<input
		id="title-length-input"
		type="text"
		bind:value={title}
		placeholder="Paste or type a title tag"
		autocomplete="off"
		class="w-full rounded-lg border border-white/20 bg-[#282828] px-4 py-3 text-white placeholder-white/40 focus:border-[#8ec07c] focus:outline-none"
	/>

	<div class="mt-4 h-2 w-full overflow-hidden rounded bg-white/10" aria-hidden="true">
		<div class="h-full {barClass[result.tone]}" style="width: {fill}%"></div>
	</div>

	<div class="mt-3 flex flex-wrap items-baseline gap-x-6 gap-y-1 text-sm" aria-live="polite">
		<span class="text-white/80"><strong class="text-white">{chars}</strong> characters</span>
		<span class="text-white/80">
			<strong class="text-white">{pixels}</strong> / {DESKTOP_LIMIT_PX} px
		</span>
		<span class="font-bold {toneClass[result.tone]}">{result.label}</span>
	</div>
	<p class="mt-2 mb-0 text-sm text-white/70">{result.detail}</p>

	{#if chars > 0}
		<div class="mt-5 rounded-lg bg-white p-4">
			<div class="mb-1 text-xs text-[#4d5156]">Desktop preview</div>
			<div
				class="overflow-hidden text-ellipsis whitespace-nowrap text-[#1a0dab]"
				style="font: {TITLE_FONT}; max-width: {DESKTOP_LIMIT_PX}px; line-height: 1.3"
			>
				{title.trim()}
			</div>
		</div>
	{/if}

	<p class="mt-4 mb-0 text-xs text-white/50">
		Measured in your browser as 20px Arial against a 600px limit. Google's exact cut-off shifts
		slightly over time, so treat the last few pixels as a margin, not a target.
	</p>
</div>
