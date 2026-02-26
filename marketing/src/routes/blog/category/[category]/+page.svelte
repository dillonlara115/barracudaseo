<script lang="ts">
	import MetaTags from '../../../../components/MetaTags.svelte';
	import { getMetaTags, getBreadcrumbSchema } from '$lib/meta';
	import { Calendar, Clock, ArrowRight, ArrowLeft, BookOpen } from '@lucide/svelte';
	import { trackCTA } from '$lib/analytics';

	let { data } = $props();

	const { category, posts } = data;

	const meta = getMetaTags({
		title: `${category} - SEO Blog`,
		description: `Read the latest articles, guides, and ${category.toLowerCase()} from the Barracuda SEO team.`,
		keywords: `SEO blog, ${category.toLowerCase()}, technical SEO, website crawling`
	});

	const structuredData = getBreadcrumbSchema([
		{ name: 'Home', url: '/' },
		{ name: 'Blog', url: '/blog' },
		{ name: category, url: `/blog/category/${category.toLowerCase()}` }
	]);

	function formatDate(dateString: string): string {
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' });
	}

	function handlePostClick(title: string) {
		trackCTA({
			source: title,
			location: `blog_category_${category.toLowerCase()}`
		});
	}
</script>

<MetaTags config={{ ...meta, structuredData }} />

<!-- Hero Section -->
<section class="bg-gradient-to-b from-[#3c3836] to-[#2d2826] py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-4xl text-center">
			<a
				href="/blog"
				class="mb-8 inline-flex items-center gap-2 text-white/70 transition-colors hover:text-[#8ec07c]"
			>
				<ArrowLeft class="h-4 w-4" />
				Back to Blog
			</a>
			<div class="mb-6 flex items-center justify-center gap-4">
				<div class="rounded-lg bg-[#8ec07c]/10 p-3">
					<BookOpen class="h-8 w-8 text-[#8ec07c]" />
				</div>
				<h1 class="font-heading text-5xl font-bold text-white md:text-6xl">{category}</h1>
			</div>
			<p class="mx-auto mb-10 max-w-3xl text-xl text-white/80 md:text-2xl">
				Explore our collection of articles and guides specifically about {category.toLowerCase()}.
			</p>
		</div>
	</div>
</section>

<!-- Category Posts -->
<section class="bg-[#2d2826] py-12">
	<div class="container mx-auto px-4">
		<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
			{#each posts as post (post.slug)}
				<div
					class="group relative overflow-hidden rounded-lg border border-white/10 bg-[#3c3836] transition-all hover:border-[#8ec07c]/50 hover:shadow-lg"
				>
					<div class="p-6">
						<a
							href="/blog/{post.slug}"
							class="after:absolute after:inset-0"
							onclick={() => handlePostClick(post.title)}
						>
							<h3
								class="mb-2 font-heading text-xl font-bold text-white transition-colors group-hover:text-[#8ec07c]"
							>
								{post.title}
							</h3>
						</a>
						<p class="mb-4 line-clamp-2 text-white/70">
							{post.description}
						</p>
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-4 text-sm text-white/50">
								<div class="flex items-center gap-1">
									<Calendar class="h-4 w-4" />
									{formatDate(post.publishDate)}
								</div>
								<div class="flex items-center gap-1">
									<Clock class="h-4 w-4" />
									{post.readTime} min
								</div>
							</div>
							<ArrowRight
								class="h-5 w-5 text-[#8ec07c] opacity-0 transition-opacity group-hover:opacity-100"
							/>
						</div>
					</div>
				</div>
			{/each}
		</div>
	</div>
</section>

<!-- CTA Section -->
<section class="bg-[#3c3836] py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-3xl text-center">
			<h2 class="mb-6 font-heading text-4xl font-bold text-white md:text-5xl">
				Ready to improve your SEO?
			</h2>
			<p class="mb-10 text-xl text-white/80">
				Start your free 100-page audit and discover technical issues holding your site back.
			</p>
			<a
				href="https://app.barracudaseo.com"
				class="inline-block rounded-lg bg-[#8ec07c] px-8 py-4 text-lg font-medium text-[#3c3836] transition-colors hover:bg-[#a0d28c]"
				target="_blank"
				rel="noopener noreferrer"
			>
				Start Your Free Audit
			</a>
		</div>
	</div>
</section>
