import { error } from '@sveltejs/kit';
import { getPostsByCategory, getAllCategories } from '$lib/blog';

export async function load({ params }) {
	const category = params.category;
	const posts = getPostsByCategory(category);

	if (posts.length === 0) {
		throw error(404, `No posts found in category: ${category}`);
	}

	// Normalize category name for display
	const categoryName = posts[0].category;

	return {
		category: categoryName,
		posts
	};
}

export async function entries() {
	const categories = getAllCategories();
	return categories.map((category) => ({ category: category.toLowerCase() }));
}
