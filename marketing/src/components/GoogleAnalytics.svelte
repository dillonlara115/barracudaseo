<script lang="ts">
	import { page } from '$app/stores';
	import { browser } from '$app/environment';
	import { onMount } from 'svelte';
	import { trackPageView } from '$lib/analytics';

	// Track page views on navigation - using onMount for safer initialization
	if (browser) {
		onMount(() => {
			let retryCount = 0;
			const maxRetries = 10;

			const trackCurrentPage = () => {
				try {
					if (
						typeof window !== 'undefined' &&
						typeof document !== 'undefined' &&
						typeof window.gtag === 'function'
					) {
						const url = window.location.pathname + window.location.search;
						const title = document.title || '';
						trackPageView(url, title);
					} else if (retryCount < maxRetries) {
						retryCount++;
						setTimeout(trackCurrentPage, 100);
					}
				} catch (error) {
					// Silently fail - don't break the page
					console.warn('GA tracking error:', error);
				}
			};

			// Initial track
			trackCurrentPage();

			// Track on navigation
			const unsubscribe = page.subscribe(() => {
				setTimeout(() => {
					try {
						if (typeof window !== 'undefined' && typeof window.gtag === 'function') {
							const url = window.location.pathname + window.location.search;
							const title = document.title || '';
							trackPageView(url, title);
						}
					} catch (error) {
						console.warn('GA navigation tracking error:', error);
					}
				}, 100);
			});

			return () => {
				unsubscribe();
			};
		});
	}
</script>
