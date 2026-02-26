<script lang="ts">
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';
	import { Play } from '@lucide/svelte';

	let iframeLoaded = $state(false);
	let shouldLoad = $state(false);
	let containerRef: HTMLDivElement | undefined = $state();

	const LOOM_VIDEO_ID = '8c848f21dcb14be9af3a0cc1b070a6cb';
	const LOOM_THUMBNAIL = `https://cdn.loom.com/sessions/thumbnails/${LOOM_VIDEO_ID}-with-play.gif`;

	function loadIframe() {
		if (!iframeLoaded) {
			iframeLoaded = true;
		}
	}

	// Use Intersection Observer to load iframe only when visible
	if (browser) {
		onMount(() => {
			if (!containerRef) return;

			const observer = new IntersectionObserver(
				(entries) => {
					entries.forEach((entry) => {
						if (entry.isIntersecting && !shouldLoad) {
							shouldLoad = true;
							// Load iframe when scrolled into view (with margin)
							// This prevents loading heavy Loom scripts until needed
							setTimeout(() => {
								if (!iframeLoaded) {
									iframeLoaded = true;
								}
							}, 1000); // Delay to ensure page is interactive first
						}
					});
				},
				{
					rootMargin: '300px' // Start loading 300px before it comes into view
				}
			);

			observer.observe(containerRef);

			return () => {
				observer.disconnect();
			};
		});
	}
</script>

<section class="bg-[#3c3836] pt-12 pb-20">
	<div class="container mx-auto px-4">
		<div class="mb-8 text-center">
			<h2 class="mb-4 font-heading text-4xl font-bold text-white md:text-5xl">
				See Barracuda SEO in Action
			</h2>
			<p class="mx-auto max-w-2xl text-xl text-white/70">
				Watch how Barracuda helps you discover, analyze, and fix SEO issues quickly
			</p>
		</div>
		<div class="mx-auto max-w-5xl" bind:this={containerRef}>
			<div class="relative" style="padding-bottom: 53.57894736842105%; height: 0;">
				{#if !iframeLoaded}
					<!-- Thumbnail placeholder with play button -->
					<button
						type="button"
						class="group absolute inset-0 h-full w-full cursor-pointer overflow-hidden rounded-lg border-0 bg-[#2d2826] p-0 shadow-2xl focus:ring-2 focus:ring-[#8ec07c] focus:ring-offset-2 focus:outline-none"
						onclick={loadIframe}
						aria-label="Play Barracuda SEO Video Walkthrough"
					>
						<img
							src={LOOM_THUMBNAIL}
							alt="Barracuda SEO Video Walkthrough"
							class="h-full w-full object-cover"
							loading="lazy"
							decoding="async"
						/>
						<div
							class="absolute inset-0 flex items-center justify-center bg-black/20 transition-colors group-hover:bg-black/30"
						>
							<div
								class="flex h-20 w-20 items-center justify-center rounded-full bg-[#8ec07c] shadow-2xl transition-transform group-hover:scale-110"
							>
								<Play class="ml-1 h-10 w-10 text-[#3c3836]" fill="currentColor" />
							</div>
						</div>
					</button>
				{/if}

				{#if iframeLoaded}
					<!-- Load iframe only when user clicks or scrolls near -->
					<iframe
						src="https://www.loom.com/embed/{LOOM_VIDEO_ID}"
						frameborder="0"
						webkitallowfullscreen
						mozallowfullscreen
						allowfullscreen
						title="Barracuda SEO Video Walkthrough - See how Barracuda helps you discover, analyze, and fix SEO issues"
						loading="lazy"
						class="absolute top-0 left-0 h-full w-full rounded-lg shadow-2xl"
					></iframe>
				{/if}
			</div>
		</div>
	</div>
</section>
