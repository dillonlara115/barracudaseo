<script lang="ts">
	import MetaTags from '../../../components/MetaTags.svelte';
	import { getMetaTags, getBreadcrumbSchema, getArticleSchema } from '$lib/meta';
	import { SITE_URL } from '$lib/constants';
	import { Calendar, Clock, ArrowLeft, ArrowRight, Tag } from '@lucide/svelte';
	import { page } from '$app/stores';
	import { trackCTA } from '$lib/analytics';
	import { blogContent } from '$lib/blog-content';

	let { data } = $props();

	const { post, relatedPosts } = data;
	const content = blogContent[post.slug] || '<p>Content coming soon...</p>';

	const meta = getMetaTags({
		title: post.title,
		description: post.description,
		keywords: post.tags.join(', '),
		author: post.author,
		ogType: 'article',
		ogImage: '/mockups/barracuda-dashboard.png'
	});

	const structuredData = [
		getBreadcrumbSchema([
			{ name: 'Home', url: '/' },
			{ name: 'Blog', url: '/blog' },
			{ name: post.title, url: `/blog/${post.slug}` }
		]),
		getArticleSchema({
			title: post.title,
			description: post.description,
			author: post.author,
			publishDate: post.publishDate,
			url: `/blog/${post.slug}`
		})
	];

	function formatDate(dateString: string): string {
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' });
	}

	function handleRelatedPostClick(title: string) {
		trackCTA({
			source: title,
			location: 'blog_post_related'
		});
	}
</script>

<MetaTags config={{ ...meta, structuredData }} />

<!-- Back to Blog -->
<section class="pt-8 pb-4 bg-[#3c3836]">
	<div class="container mx-auto px-4">
		<a
			href="/blog"
			class="inline-flex items-center gap-2 text-white/70 hover:text-[#8ec07c] transition-colors"
		>
			<ArrowLeft class="w-4 h-4" />
			Back to Blog
		</a>
	</div>
</section>

<!-- Article Header -->
<section class="py-12 bg-gradient-to-b from-[#3c3836] to-[#2d2826]">
	<div class="container mx-auto px-4 max-w-4xl">
		<div class="mb-6">
			<span class="px-3 py-1 bg-[#8ec07c]/20 text-[#8ec07c] text-sm font-medium rounded-full">
				{post.category}
			</span>
		</div>
		<h1 class="text-4xl md:text-5xl font-heading font-bold mb-6 text-white">
			{post.title}
		</h1>
		<p class="text-xl text-white/80 mb-8 max-w-3xl">
			{post.description}
		</p>
		<div class="flex flex-wrap items-center gap-6 text-white/60">
			<div class="flex items-center gap-2">
				<Calendar class="w-5 h-5" />
				{formatDate(post.publishDate)}
			</div>
			<div class="flex items-center gap-2">
				<Clock class="w-5 h-5" />
				{post.readTime} min read
			</div>
			<div class="flex items-center gap-2">
				<span>By {post.author}</span>
			</div>
		</div>
		{#if post.tags.length > 0}
			<div class="flex flex-wrap gap-2 mt-6">
				{#each post.tags as tag}
					<span class="px-2 py-1 bg-white/5 text-white/70 text-xs rounded flex items-center gap-1">
						<Tag class="w-3 h-3" />
						{tag}
					</span>
				{/each}
			</div>
		{/if}
	</div>
</section>

<!-- Article Content -->
<article class="py-12 bg-[#2d2826]">
	<div class="container mx-auto px-4 max-w-4xl">
		<div class="blog-content">
			{@html content}
		</div>
	</div>
</article>

<style>
	.blog-content {
		color: rgba(255, 255, 255, 0.8);
		line-height: 1.75;
	}

	/* Heading styles are now in app.css @layer components */
	/* Paragraph spacing */
	.blog-content p {
		margin-top: 0 !important;
		margin-bottom: 1.25rem !important;
		color: rgba(255, 255, 255, 0.8);
		line-height: 1.75;
	}

	/* Ensure spacing after headings */
	.blog-content h1 + p,
	.blog-content h2 + p,
	.blog-content h3 + p,
	.blog-content h4 + p,
	.blog-content h5 + p,
	.blog-content h6 + p {
		margin-top: 0 !important;
	}

	/* List styles - these complement the app.css overrides */
	.blog-content ul,
	.blog-content ol {
		margin-top: 1rem !important;
		margin-bottom: 1.25rem !important;
		color: rgba(255, 255, 255, 0.8);
	}

	/* Spacing after headings before lists */
	.blog-content h2 + ul,
	.blog-content h2 + ol,
	.blog-content h3 + ul,
	.blog-content h3 + ol,
	.blog-content h4 + ul,
	.blog-content h4 + ol {
		margin-top: 0.75rem !important;
	}

	.blog-content li {
		line-height: 1.75;
		color: rgba(255, 255, 255, 0.8);
	}

	.blog-content ul li::marker {
		color: #8ec07c;
	}

	.blog-content ol li::marker {
		color: #8ec07c;
	}

	.blog-content li > ul,
	.blog-content li > ol {
		margin-top: 0.5rem;
		margin-bottom: 0.5rem;
	}

	.blog-content strong {
		color: white;
		font-weight: 600;
	}

	.blog-content a {
		color: #8ec07c;
		text-decoration: none;
		transition: color 0.2s;
	}

	.blog-content a:hover {
		color: #a0d28c;
		text-decoration: underline;
	}

	/* Code blocks */
	.blog-content pre {
		margin-top: 1rem !important;
		margin-bottom: 1.25rem !important;
		padding: 1rem;
		background-color: rgba(0, 0, 0, 0.3);
		border-radius: 0.5rem;
		overflow-x: auto;
	}

	.blog-content code {
		font-family: var(--font-mono);
		font-size: 0.875rem;
	}

	.blog-content pre code {
		background: transparent;
		padding: 0;
	}

	/* Tables */
	.blog-content table {
		width: 100%;
		border-collapse: collapse;
		margin-top: 1.5rem !important;
		margin-bottom: 1.5rem !important;
	}

	.blog-content table th,
	.blog-content table td {
		border: 1px solid rgba(255, 255, 255, 0.2);
		padding: 0.75rem;
		text-align: left;
	}

	.blog-content table th {
		background-color: rgba(60, 56, 54, 0.5);
		color: white;
		font-weight: 600;
	}

	.blog-content table td {
		color: rgba(255, 255, 255, 0.8);
	}

	/* Blockquotes */
	.blog-content blockquote {
		margin-top: 1.5rem !important;
		margin-bottom: 1.5rem !important;
		padding-left: 1.5rem;
		border-left: 4px solid #8ec07c;
		color: rgba(255, 255, 255, 0.7);
		font-style: italic;
	}

	/* Images */
	.blog-content img {
		margin-top: 1.5rem !important;
		margin-bottom: 1.5rem !important;
		max-width: 100%;
		height: auto;
		border-radius: 0.5rem;
	}

	/* Horizontal rules */
	.blog-content hr {
		margin-top: 2rem !important;
		margin-bottom: 2rem !important;
		border: none;
		border-top: 1px solid rgba(255, 255, 255, 0.2);
	}

	/* Ensure first element doesn't have top margin */
	.blog-content > *:first-child {
		margin-top: 0 !important;
	}

	/* Ensure last element doesn't have bottom margin */
	.blog-content > *:last-child {
		margin-bottom: 0 !important;
	}
</style>

<!-- Related Posts -->
{#if relatedPosts.length > 0}
	<section class="py-12 bg-[#3c3836]">
		<div class="container mx-auto px-4 max-w-4xl">
			<h2 class="text-3xl font-heading font-bold mb-8 text-white">Related Posts</h2>
			<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
				{#each relatedPosts as relatedPost}
					<a
						href="/blog/{relatedPost.slug}"
						class="block bg-[#2d2826] rounded-lg border border-white/10 hover:border-[#8ec07c]/50 transition-all p-6 group"
						onclick={() => handleRelatedPostClick(relatedPost.title)}
					>
						<h3 class="text-lg font-heading font-bold mb-2 text-white group-hover:text-[#8ec07c] transition-colors">
							{relatedPost.title}
						</h3>
						<p class="text-white/70 text-sm mb-4 line-clamp-2">
							{relatedPost.description}
						</p>
						<div class="flex items-center gap-2 text-white/50 text-sm">
							<Clock class="w-4 h-4" />
							{relatedPost.readTime} min read
						</div>
					</a>
				{/each}
			</div>
		</div>
	</section>
{/if}

<!-- CTA Section -->
<section class="py-20 bg-[#2d2826]">
	<div class="container mx-auto px-4">
		<div class="max-w-3xl mx-auto text-center">
			<h2 class="text-4xl md:text-5xl font-heading font-bold mb-6 text-white">
				Ready to audit your site?
			</h2>
			<p class="text-xl text-white/80 mb-10">
				Start your free 100-page audit and discover technical SEO issues in minutes.
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
