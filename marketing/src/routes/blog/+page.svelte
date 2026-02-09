<script lang="ts">
	import MetaTags from '../../components/MetaTags.svelte';
	import { getMetaTags, getBreadcrumbSchema } from '$lib/meta';
	import { getAllBlogPosts, getFeaturedPosts } from '$lib/blog';
	import { Calendar, Clock, ArrowRight, BookOpen } from '@lucide/svelte';
	import { trackCTA } from '$lib/analytics';

	const meta = getMetaTags({
		title: 'Technical SEO Guides & Best Practices',
		description: 'Learn about technical SEO, website crawling, and SEO best practices. Tips, guides, and comparisons from the Barracuda SEO team.',
		keywords: 'SEO blog, technical SEO, website crawling, SEO guides, SEO tips'
	});

	const structuredData = getBreadcrumbSchema([
		{ name: 'Home', url: '/' },
		{ name: 'Blog', url: '/blog' }
	]);

	const allPosts = getAllBlogPosts();
	const featuredPosts = getFeaturedPosts();

	function formatDate(dateString: string): string {
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' });
	}

	function handlePostClick(title: string) {
		trackCTA({
			source: title,
			location: 'blog_archive'
		});
	}
</script>

<MetaTags config={{ ...meta, structuredData }} />

<!-- Hero Section -->
<section class="py-20 bg-gradient-to-b from-[#3c3836] to-[#2d2826]">
	<div class="container mx-auto px-4">
		<div class="text-center max-w-4xl mx-auto">
			<div class="flex items-center justify-center gap-4 mb-6">
				<div class="p-3 bg-[#8ec07c]/10 rounded-lg">
					<BookOpen class="w-8 h-8 text-[#8ec07c]" />
				</div>
				<h1 class="text-5xl md:text-6xl font-heading font-bold text-white">
					Barracuda Blog
				</h1>
			</div>
			<p class="text-xl md:text-2xl text-white/80 mb-10 max-w-3xl mx-auto">
				Learn about technical SEO, website crawling, and SEO best practices. Tips, guides, and comparisons from our team.
			</p>
		</div>
	</div>
</section>

<!-- Featured Posts -->
{#if featuredPosts.length > 0}
	<section class="py-12 bg-[#3c3836]">
		<div class="container mx-auto px-4">
			<h2 class="text-3xl font-heading font-bold mb-8 text-white">Featured Posts</h2>
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
				{#each featuredPosts as post}
					<a
						href="/blog/{post.slug}"
						class="block bg-[#2d2826] rounded-lg border border-white/10 hover:border-[#8ec07c]/50 transition-all hover:shadow-lg overflow-hidden group"
						onclick={() => handlePostClick(post.title)}
					>
						<div class="p-6">
							<div class="flex items-center gap-2 mb-3">
								<span class="px-3 py-1 bg-[#8ec07c]/20 text-[#8ec07c] text-xs font-medium rounded-full">
									{post.category}
								</span>
								<span class="text-white/50 text-sm">Featured</span>
							</div>
							<h3 class="text-xl font-heading font-bold mb-2 text-white group-hover:text-[#8ec07c] transition-colors">
								{post.title}
							</h3>
							<p class="text-white/70 mb-4 line-clamp-2">
								{post.description}
							</p>
							<div class="flex items-center gap-4 text-sm text-white/50">
								<div class="flex items-center gap-1">
									<Calendar class="w-4 h-4" />
									{formatDate(post.publishDate)}
								</div>
								<div class="flex items-center gap-1">
									<Clock class="w-4 h-4" />
									{post.readTime} min read
								</div>
							</div>
						</div>
					</a>
				{/each}
			</div>
		</div>
	</section>
{/if}

<!-- All Posts -->
<section class="py-12 bg-[#2d2826]">
	<div class="container mx-auto px-4">
		<h2 class="text-3xl font-heading font-bold mb-8 text-white">All Posts</h2>
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
			{#each allPosts as post}
				<a
					href="/blog/{post.slug}"
					class="block bg-[#3c3836] rounded-lg border border-white/10 hover:border-[#8ec07c]/50 transition-all hover:shadow-lg overflow-hidden group"
					onclick={() => handlePostClick(post.title)}
				>
					<div class="p-6">
						<div class="flex items-center gap-2 mb-3">
							<span class="px-3 py-1 bg-[#8ec07c]/20 text-[#8ec07c] text-xs font-medium rounded-full">
								{post.category}
							</span>
						</div>
						<h3 class="text-xl font-heading font-bold mb-2 text-white group-hover:text-[#8ec07c] transition-colors">
							{post.title}
						</h3>
						<p class="text-white/70 mb-4 line-clamp-2">
							{post.description}
						</p>
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-4 text-sm text-white/50">
								<div class="flex items-center gap-1">
									<Calendar class="w-4 h-4" />
									{formatDate(post.publishDate)}
								</div>
								<div class="flex items-center gap-1">
									<Clock class="w-4 h-4" />
									{post.readTime} min
								</div>
							</div>
							<ArrowRight class="w-5 h-5 text-[#8ec07c] opacity-0 group-hover:opacity-100 transition-opacity" />
						</div>
					</div>
				</a>
			{/each}
		</div>
	</section>

<!-- CTA Section -->
<section class="py-20 bg-[#3c3836]">
	<div class="container mx-auto px-4">
		<div class="max-w-3xl mx-auto text-center">
			<h2 class="text-4xl md:text-5xl font-heading font-bold mb-6 text-white">
				Ready to improve your SEO?
			</h2>
			<p class="text-xl text-white/80 mb-10">
				Start your free 100-page audit and discover technical issues holding your site back.
			</p>
			<a
				href="https://app.barracudaseo.com"
				class="inline-block bg-[#8ec07c] hover:bg-[#a0d28c] text-[#3c3836] px-8 py-4 rounded-lg font-medium text-lg transition-colors"
				target="_blank"
				rel="noopener noreferrer"
			>
				Start Your Free Audit
			</a>
		</div>
	</div>
</section>
