<script lang="ts">
	import MetaTags from '../../../components/MetaTags.svelte';
	import { getMetaTags, getBreadcrumbSchema, getArticleSchema, getHowToSchema } from '$lib/meta';
	import { Calendar, Clock, ArrowLeft, Tag } from '@lucide/svelte';
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

	if (post.slug === 'find-declining-pages-google-search-console') {
		structuredData.push(
			getHowToSchema({
				name: 'How to Find Declining Pages in Google Search Console',
				description:
					'Step-by-step guide to finding and analyzing declining pages using GSC data comparison.',
				steps: [
					{
						name: 'Open GSC and navigate to Search Results',
						text: 'Log into Search Console, select your property, and click "Search results" in the left sidebar.'
					},
					{
						name: 'Set a date comparison',
						text: 'Click the date filter at the top and switch to "Compare." Use "Last 3 months" compared to the "Previous period."'
					},
					{
						name: 'Switch to the Pages tab and sort by clicks change',
						text: 'In the table below the chart, click the "Pages" tab. Then click the "Clicks Difference" column header to sort ascending.'
					},
					{
						name: 'Cross-reference with impressions and position data',
						text: 'For each declining page, check whether impressions also dropped or held steady. Toggle on "Average position" to see if rankings shifted too.'
					},
					{
						name: 'Export and document',
						text: 'Export the data to a spreadsheet and flag the pages worth investigating.'
					},
					{
						name: 'Repeat the process weekly',
						text: 'A recurring process tells you where things are heading.'
					}
				]
			})
		);
	}

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
<section class="bg-[#3c3836] pt-8 pb-4">
	<div class="container mx-auto px-4">
		<a
			href="/blog"
			class="inline-flex items-center gap-2 text-white/70 transition-colors hover:text-[#8ec07c]"
		>
			<ArrowLeft class="h-4 w-4" />
			Back to Blog
		</a>
	</div>
</section>

<!-- Article Header -->
<section class="bg-gradient-to-b from-[#3c3836] to-[#2d2826] py-12">
	<div class="container mx-auto max-w-4xl px-4">
		<div class="mb-6">
			<a
				href="/blog/category/{post.category.toLowerCase()}"
				class="rounded-full bg-[#8ec07c]/20 px-3 py-1 text-sm font-medium text-[#8ec07c] transition-colors hover:bg-[#8ec07c]/30"
			>
				{post.category}
			</a>
		</div>
		<h1 class="mb-6 font-heading text-4xl font-bold text-white md:text-5xl">
			{post.title}
		</h1>
		<p class="mb-8 max-w-3xl text-xl text-white/80">
			{post.description}
		</p>
		<div class="flex flex-wrap items-center gap-6 text-white/60">
			<div class="flex items-center gap-2">
				<Calendar class="h-5 w-5" />
				{formatDate(post.publishDate)}
			</div>
			<div class="flex items-center gap-2">
				<Clock class="h-5 w-5" />
				{post.readTime} min read
			</div>
			<div class="flex items-center gap-2">
				<span>By {post.author}</span>
			</div>
		</div>
		{#if post.tags.length > 0}
			<div class="mt-6 flex flex-wrap gap-2">
				{#each post.tags as tag}
					<span class="flex items-center gap-1 rounded bg-white/5 px-2 py-1 text-xs text-white/70">
						<Tag class="h-3 w-3" />
						{tag}
					</span>
				{/each}
			</div>
		{/if}
	</div>
</section>

<!-- Article Content -->
<article class="bg-[#2d2826] py-12">
	<div class="container mx-auto max-w-4xl px-4">
		<div class="blog-content">
			{@html content}
		</div>
	</div>
</article>

<!-- Related Posts -->
{#if relatedPosts.length > 0}
	<section class="bg-[#3c3836] py-12">
		<div class="container mx-auto max-w-4xl px-4">
			<h2 class="mb-8 font-heading text-3xl font-bold text-white">Related Posts</h2>
			<div class="grid grid-cols-1 gap-6 md:grid-cols-3">
				{#each relatedPosts as relatedPost}
					<a
						href="/blog/{relatedPost.slug}"
						class="group block rounded-lg border border-white/10 bg-[#2d2826] p-6 transition-all hover:border-[#8ec07c]/50"
						onclick={() => handleRelatedPostClick(relatedPost.title)}
					>
						<h3
							class="mb-2 font-heading text-lg font-bold text-white transition-colors group-hover:text-[#8ec07c]"
						>
							{relatedPost.title}
						</h3>
						<p class="mb-4 line-clamp-2 text-sm text-white/70">
							{relatedPost.description}
						</p>
						<div class="flex items-center gap-2 text-sm text-white/50">
							<Clock class="h-4 w-4" />
							{relatedPost.readTime} min read
						</div>
					</a>
				{/each}
			</div>
		</div>
	</section>
{/if}

<!-- CTA Section -->
<section class="bg-[#2d2826] py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-3xl text-center">
			<h2 class="mb-6 font-heading text-4xl font-bold text-white md:text-5xl">
				Ready to audit your site?
			</h2>
			<p class="mb-10 text-xl text-white/80">
				Start your free 100-page audit and discover technical SEO issues in minutes.
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
