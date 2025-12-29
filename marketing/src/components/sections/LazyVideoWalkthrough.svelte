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

<section class="pt-12 pb-20 bg-[#3c3836]">
	<div class="container mx-auto px-4">
		<div class="text-center mb-8">
			<h2 class="text-4xl md:text-5xl font-heading font-bold mb-4 text-white">
				See Barracuda SEO in Action
			</h2>
			<p class="text-xl text-white/70 max-w-2xl mx-auto">
				Watch how Barracuda helps you discover, analyze, and fix SEO issues quickly
			</p>
		</div>
		<div class="max-w-5xl mx-auto" bind:this={containerRef}>
			<div class="relative" style="padding-bottom: 53.57894736842105%; height: 0;">
				{#if !iframeLoaded}
					<!-- Thumbnail placeholder with play button -->
					<button
						type="button"
						class="absolute inset-0 rounded-lg shadow-2xl overflow-hidden bg-[#2d2826] cursor-pointer group border-0 p-0 w-full h-full focus:outline-none focus:ring-2 focus:ring-[#8ec07c] focus:ring-offset-2"
						onclick={loadIframe}
						aria-label="Play Barracuda SEO Video Walkthrough"
					>
						<img
							src={LOOM_THUMBNAIL}
							alt="Barracuda SEO Video Walkthrough"
							class="w-full h-full object-cover"
							loading="lazy"
							decoding="async"
						/>
						<div class="absolute inset-0 flex items-center justify-center bg-black/20 group-hover:bg-black/30 transition-colors">
							<div class="w-20 h-20 bg-[#8ec07c] rounded-full flex items-center justify-center shadow-2xl group-hover:scale-110 transition-transform">
								<Play class="w-10 h-10 text-[#3c3836] ml-1" fill="currentColor" />
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
						class="absolute top-0 left-0 w-full h-full rounded-lg shadow-2xl"
					></iframe>
				{/if}
			</div>
		</div>
	</div>
</section>
