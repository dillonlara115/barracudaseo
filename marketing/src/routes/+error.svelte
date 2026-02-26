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
	const pageTitle = $derived(
		is404 ? 'Page Not Found - Barracuda SEO' : `Error ${status} - Barracuda SEO`
	);
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

<section
	class="flex min-h-[70vh] items-center justify-center bg-gradient-to-b from-[#3c3836] via-[#8ec07c]/20 to-[#3c3836] px-4 py-20"
>
	<div class="mx-auto max-w-2xl text-center">
		<div class="mb-6 flex items-center justify-center gap-4">
			<div class="rounded-full bg-[#8ec07c]/10 p-4">
				{#if is404}
					<Search class="h-12 w-12 text-[#8ec07c]" />
				{:else}
					<AlertCircle class="h-12 w-12 text-[#d79921]" />
				{/if}
			</div>
		</div>

		<h1 class="mb-4 font-heading text-6xl font-bold text-white md:text-7xl">
			{status}
		</h1>

		{#if is404}
			<h2 class="mb-4 font-heading text-3xl font-bold text-white md:text-4xl">Page Not Found</h2>
			<p class="mb-8 text-xl text-white/80">
				The page you're looking for doesn't exist or has been moved.
			</p>
		{:else}
			<h2 class="mb-4 font-heading text-3xl font-bold text-white md:text-4xl">
				Something Went Wrong
			</h2>
			<p class="mb-8 text-xl text-white/80">
				{error?.message || 'An unexpected error occurred. Please try again.'}
			</p>
		{/if}

		<div class="mb-12 flex flex-col justify-center gap-4 sm:flex-row">
			<a
				href="/"
				class="inline-flex items-center justify-center gap-2 rounded-lg bg-[#8ec07c] px-8 py-4 text-lg font-medium text-[#3c3836] transition-colors hover:bg-[#a0d28c]"
			>
				<Home class="h-5 w-5" />
				Go Home
			</a>
			<button
				onclick={() => window.history.back()}
				class="inline-flex items-center justify-center gap-2 rounded-lg border-2 border-white/20 px-8 py-4 text-lg font-medium text-white transition-colors hover:border-[#8ec07c]"
			>
				<ArrowLeft class="h-5 w-5" />
				Go Back
			</button>
		</div>

		{#if is404}
			<div class="rounded-lg border border-white/10 bg-[#2d2826] p-6 text-left">
				<h3 class="mb-4 font-heading text-lg font-bold text-white">Popular Pages</h3>
				<div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
					<a
						href="/features"
						class="flex items-center gap-2 text-white/70 transition-colors hover:text-[#8ec07c]"
					>
						<span>→</span>
						<span>Features</span>
					</a>
					<a
						href="/pricing"
						class="flex items-center gap-2 text-white/70 transition-colors hover:text-[#8ec07c]"
					>
						<span>→</span>
						<span>Pricing</span>
					</a>
					<a
						href="/about"
						class="flex items-center gap-2 text-white/70 transition-colors hover:text-[#8ec07c]"
					>
						<span>→</span>
						<span>About</span>
					</a>
					<a
						href="/faq"
						class="flex items-center gap-2 text-white/70 transition-colors hover:text-[#8ec07c]"
					>
						<span>→</span>
						<span>FAQ</span>
					</a>
					<a
						href="/use-cases/e-commerce"
						class="flex items-center gap-2 text-white/70 transition-colors hover:text-[#8ec07c]"
					>
						<span>→</span>
						<span>E-commerce SEO</span>
					</a>
					<a
						href="/use-cases/local-seo"
						class="flex items-center gap-2 text-white/70 transition-colors hover:text-[#8ec07c]"
					>
						<span>→</span>
						<span>Local SEO</span>
					</a>
				</div>
			</div>
		{/if}
	</div>
</section>
