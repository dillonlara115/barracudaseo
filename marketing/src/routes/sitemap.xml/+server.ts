import { SITE_URL } from '$lib/constants';

import { getAllBlogPosts } from '$lib/blog';

// Define all routes with their priority and change frequency
const staticRoutes = [
	{ path: '', priority: '1.0', changefreq: 'weekly' }, // Home
	{ path: '/about', priority: '0.8', changefreq: 'monthly' },
	{ path: '/features', priority: '0.9', changefreq: 'monthly' },
	{ path: '/pricing', priority: '0.9', changefreq: 'monthly' },
	{ path: '/faq', priority: '0.8', changefreq: 'monthly' },
	{ path: '/roadmap', priority: '0.7', changefreq: 'monthly' },
	{ path: '/privacy', priority: '0.5', changefreq: 'yearly' },
	{ path: '/terms', priority: '0.5', changefreq: 'yearly' },
	{ path: '/use-cases/e-commerce', priority: '0.8', changefreq: 'monthly' },
	{ path: '/use-cases/local-seo', priority: '0.8', changefreq: 'monthly' },
	{ path: '/use-cases/programmatic-seo', priority: '0.8', changefreq: 'monthly' },
	{ path: '/blog', priority: '0.9', changefreq: 'weekly' }
];

// Get blog posts dynamically
const blogPosts = getAllBlogPosts();
const blogRoutes = blogPosts.map(post => ({
	path: `/blog/${post.slug}`,
	priority: '0.8',
	changefreq: 'monthly',
	lastmod: post.publishDate
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
    <lastmod>${route.lastmod ? route.lastmod.split('T')[0] : new Date().toISOString().split('T')[0]}</lastmod>
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
