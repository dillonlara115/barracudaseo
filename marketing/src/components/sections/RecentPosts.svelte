<script lang="ts">
	import { getAllBlogPosts, type BlogPost } from '$lib/blog';
	import { Calendar, Clock, ChevronRight } from '@lucide/svelte';

	let { count = 3, title = 'From the Blog', subtitle = '' }: { count?: number; title?: string; subtitle?: string } = $props();

	const posts: BlogPost[] = $derived(getAllBlogPosts().slice(0, count));

	const formatDate = (dateStr: string) => {
		return new Date(dateStr).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			year: 'numeric'
		});
	};
</script>

<section class="bg-[#222] px-4 py-24">
	<div class="container mx-auto max-w-6xl">
		<div class="mb-12 text-center">
			<p class="mb-4 text-sm font-medium tracking-widest text-[#8ec07c] uppercase">
				Guides & Insights
			</p>
			<h2 class="mb-4 font-heading text-3xl font-bold text-white md:text-4xl lg:text-5xl">
				{title}
			</h2>
			{#if subtitle}
				<p class="mx-auto max-w-2xl text-lg text-white/50">{subtitle}</p>
			{/if}
		</div>

		<div class="grid grid-cols-1 gap-6 md:grid-cols-3">
			{#each posts as post}
				<a
					href="/blog/{post.slug}"
					class="group flex flex-col rounded-2xl border border-white/5 bg-[#282828] p-6 transition-all duration-300 hover:border-[#8ec07c]/20"
				>
					<div class="mb-3">
						<span
							class="inline-block rounded-full bg-[#8ec07c]/10 px-3 py-1 text-xs font-medium text-[#8ec07c]"
						>
							{post.category}
						</span>
					</div>
					<h3
						class="mb-3 font-heading text-lg font-bold leading-snug text-white transition-colors group-hover:text-[#8ec07c]"
					>
						{post.title}
					</h3>
					<p class="mb-4 flex-1 text-sm leading-relaxed text-white/50">
						{post.description}
					</p>
					<div class="flex items-center justify-between text-xs text-white/30">
						<div class="flex items-center gap-3">
							<span class="flex items-center gap-1">
								<Calendar class="h-3 w-3" />
								{formatDate(post.publishDate)}
							</span>
							<span class="flex items-center gap-1">
								<Clock class="h-3 w-3" />
								{post.readTime} min
							</span>
						</div>
						<ChevronRight
							class="h-4 w-4 transition-transform group-hover:translate-x-0.5"
						/>
					</div>
				</a>
			{/each}
		</div>

		<div class="mt-10 text-center">
			<a
				href="/blog"
				class="inline-flex items-center gap-2 rounded-xl border border-white/15 px-6 py-3 text-sm font-medium text-white/70 transition-all hover:border-[#8ec07c]/30 hover:text-white"
			>
				View all posts
				<ChevronRight class="h-4 w-4" />
			</a>
		</div>
	</div>
</section>
