<script lang="ts">
	import { page } from '$app/stores';
	import { browser } from '$app/environment';
	import { trackPageView } from '$lib/analytics';

	// Track page views on navigation
	$effect(() => {
		if (!browser) return;

		// Get current page info
		const url = $page.url.pathname + $page.url.search;
		const title = document.title;

		// Wait for gtag to be available (it loads asynchronously)
		const checkAndTrack = () => {
			if (typeof window !== 'undefined' && typeof window.gtag === 'function') {
				trackPageView(url, title);
			} else {
				// Retry after a short delay if gtag isn't ready yet
				setTimeout(checkAndTrack, 100);
			}
		};

		checkAndTrack();
	});
</script>
