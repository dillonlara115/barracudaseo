import { SITE_URL } from '$lib/constants';

// Define all routes with their priority and change frequency
const routes = [
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
	{ path: '/use-cases/programmatic-seo', priority: '0.8', changefreq: 'monthly' }
];

export async function GET() {
	const sitemap = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
${routes
	.map(
		(route) => `  <url>
    <loc>${SITE_URL}${route.path}</loc>
    <lastmod>${new Date().toISOString().split('T')[0]}</lastmod>
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
