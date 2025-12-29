import { error } from '@sveltejs/kit';
import { getBlogPost, getAllBlogPosts } from '$lib/blog';

export async function load({ params }) {
	const post = getBlogPost(params.slug);

	if (!post) {
		throw error(404, 'Blog post not found');
	}

	// Get related posts (excluding current post)
	const allPosts = getAllBlogPosts();
	const relatedPosts = allPosts
		.filter(p => p.slug !== params.slug && p.category === post.category)
		.slice(0, 3);

	return {
		post,
		relatedPosts
	};
}

// Generate static paths for all blog posts at build time
export async function entries() {
	const posts = getAllBlogPosts();
	return posts.map(post => ({ slug: post.slug }));
}
