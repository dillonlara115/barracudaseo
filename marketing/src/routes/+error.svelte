<script lang="ts">
	import { page } from '$app/stores';
	import { browser } from '$app/environment';
	import { onMount } from 'svelte';
	import MetaTags from '../components/MetaTags.svelte';
	import { getMetaTags } from '$lib/meta';
	import { trackEvent } from '$lib/analytics';
	import { Home, Search, ArrowLeft, AlertCircle } from '@lucide/svelte';

	let { error, status }: { error: Error; status: number } = $props();

	const is404 = $derived(status === 404);
	const pageTitle = $derived(is404 ? 'Page Not Found - Barracuda SEO' : `Error ${status} - Barracuda SEO`);
	const pageDescription = $derived(
		is404
			? "The page you're looking for doesn't exist. Return to Barracuda SEO homepage or explore our features."
			: 'An error occurred. Please try again or return to the homepage.'
	);

	const meta = $derived(
		getMetaTags({
			title: is404 ? 'Page Not Found' : `Error ${status}`,
			description: pageDescription,
			robots: 'noindex, nofollow' // Don't index error pages
		})
	);

	// Track 404 errors in analytics
	if (browser) {
		onMount(() => {
			if (is404) {
				trackEvent('404_error', {
					category: 'error',
					action: 'page_not_found',
					label: $page.url.pathname,
					page: $page.url.pathname
				});
			}
		});
	}
</script>

<MetaTags config={meta} />

<section class="min-h-[70vh] flex items-center justify-center bg-gradient-to-b from-[#3c3836] via-[#8ec07c]/20 to-[#3c3836] py-20 px-4">
	<div class="text-center max-w-2xl mx-auto">
		<div class="flex items-center justify-center gap-4 mb-6">
			<div class="p-4 bg-[#8ec07c]/10 rounded-full">
				{#if is404}
					<Search class="w-12 h-12 text-[#8ec07c]" />
				{:else}
					<AlertCircle class="w-12 h-12 text-[#d79921]" />
				{/if}
			</div>
		</div>

		<h1 class="text-6xl md:text-7xl font-heading font-bold mb-4 text-white">
			{status}
		</h1>

		{#if is404}
			<h2 class="text-3xl md:text-4xl font-heading font-bold mb-4 text-white">
				Page Not Found
			</h2>
			<p class="text-xl text-white/80 mb-8">
				The page you're looking for doesn't exist or has been moved.
			</p>
		{:else}
			<h2 class="text-3xl md:text-4xl font-heading font-bold mb-4 text-white">
				Something Went Wrong
			</h2>
			<p class="text-xl text-white/80 mb-8">
				{error?.message || 'An unexpected error occurred. Please try again.'}
			</p>
		{/if}

		<div class="flex flex-col sm:flex-row gap-4 justify-center mb-12">
			<a
				href="/"
				class="inline-flex items-center justify-center gap-2 bg-[#8ec07c] hover:bg-[#a0d28c] text-[#3c3836] px-8 py-4 rounded-lg font-medium text-lg transition-colors"
			>
				<Home class="w-5 h-5" />
				Go Home
			</a>
			<button
				onclick={() => window.history.back()}
				class="inline-flex items-center justify-center gap-2 border-2 border-white/20 hover:border-[#8ec07c] text-white px-8 py-4 rounded-lg font-medium text-lg transition-colors"
			>
				<ArrowLeft class="w-5 h-5" />
				Go Back
			</button>
		</div>

		{#if is404}
			<div class="bg-[#2d2826] rounded-lg p-6 border border-white/10 text-left">
				<h3 class="text-lg font-heading font-bold mb-4 text-white">Popular Pages</h3>
				<div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
					<a href="/features" class="text-white/70 hover:text-[#8ec07c] transition-colors flex items-center gap-2">
						<span>→</span>
						<span>Features</span>
					</a>
					<a href="/pricing" class="text-white/70 hover:text-[#8ec07c] transition-colors flex items-center gap-2">
						<span>→</span>
						<span>Pricing</span>
					</a>
					<a href="/about" class="text-white/70 hover:text-[#8ec07c] transition-colors flex items-center gap-2">
						<span>→</span>
						<span>About</span>
					</a>
					<a href="/faq" class="text-white/70 hover:text-[#8ec07c] transition-colors flex items-center gap-2">
						<span>→</span>
						<span>FAQ</span>
					</a>
					<a href="/use-cases/e-commerce" class="text-white/70 hover:text-[#8ec07c] transition-colors flex items-center gap-2">
						<span>→</span>
						<span>E-commerce SEO</span>
					</a>
					<a href="/use-cases/local-seo" class="text-white/70 hover:text-[#8ec07c] transition-colors flex items-center gap-2">
						<span>→</span>
						<span>Local SEO</span>
					</a>
				</div>
			</div>
		{/if}
	</div>
</section>
