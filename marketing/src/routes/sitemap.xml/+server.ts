import { SITE_URL } from '$lib/constants';

import { getAllBlogPosts } from '$lib/blog';

// The blog index changes whenever a post is published
function latestBlogDate(): string {
	const dates = getAllBlogPosts().map((p) => p.updatedDate ?? p.publishDate);
	return dates.sort().at(-1) ?? new Date().toISOString();
}

// Define all routes with their priority and change frequency.
// lastmod must be the date the page content actually changed — a build-time
// "today" on every page teaches Google to ignore the field entirely.
// Update the date when you meaningfully edit a page.
const staticRoutes = [
	{ path: '', priority: '1.0', changefreq: 'weekly', lastmod: '2026-09-14' }, // Home
	{ path: '/about', priority: '0.8', changefreq: 'monthly', lastmod: '2026-03-14' },
	{ path: '/features', priority: '0.9', changefreq: 'monthly', lastmod: '2026-09-21' },
	{ path: '/pricing', priority: '0.9', changefreq: 'monthly', lastmod: '2026-09-21' },
	{ path: '/faq', priority: '0.8', changefreq: 'monthly', lastmod: '2026-02-26' },
	{ path: '/roadmap', priority: '0.7', changefreq: 'monthly', lastmod: '2026-02-26' },
	{ path: '/privacy', priority: '0.5', changefreq: 'yearly', lastmod: '2026-02-26' },
	{ path: '/terms', priority: '0.5', changefreq: 'yearly', lastmod: '2026-02-26' },
	{ path: '/use-cases/e-commerce', priority: '0.8', changefreq: 'monthly', lastmod: '2026-09-14' },
	{ path: '/use-cases/local-seo', priority: '0.8', changefreq: 'monthly', lastmod: '2026-09-21' },
	{
		path: '/use-cases/programmatic-seo',
		priority: '0.8',
		changefreq: 'monthly',
		lastmod: '2026-02-26'
	},
	{ path: '/blog', priority: '0.9', changefreq: 'weekly', lastmod: latestBlogDate() }
];

// Get blog posts dynamically
const blogPosts = getAllBlogPosts();
const blogRoutes = blogPosts.map((post) => ({
	path: `/blog/${post.slug}`,
	priority: '0.8',
	changefreq: 'monthly',
	lastmod: post.updatedDate ?? post.publishDate
}));

const routes = [...staticRoutes, ...blogRoutes];

// Normalize path: remove trailing slash except for root
const normalizePath = (path: string): string => {
	if (path !== '/' && path.endsWith('/')) {
		return path.slice(0, -1);
	}
	return path;
};

export async function GET() {
	const sitemap = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
${routes
	.map(
		(route) => `  <url>
    <loc>${SITE_URL}${normalizePath(route.path)}</loc>
    <lastmod>${route.lastmod.split('T')[0]}</lastmod>
    <changefreq>${route.changefreq}</changefreq>
    <priority>${route.priority}</priority>
  </url>`
	)
	.join('\n')}
</urlset>`;

	return new Response(sitemap, {
		headers: {
			'Content-Type': 'application/xml',
			'Cache-Control': 'public, max-age=3600' // Cache for 1 hour
		}
	});
}
