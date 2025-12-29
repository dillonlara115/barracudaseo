export interface BlogPost {
	slug: string;
	title: string;
	description: string;
	author: string;
	publishDate: string; // ISO date string
	readTime: number; // minutes
	category: string;
	tags: string[];
	featured?: boolean;
}

export const blogPosts: BlogPost[] = [
	{
		slug: 'complete-technical-seo-audit-guide',
		title: 'How to Perform a Complete Technical SEO Audit: A Step-by-Step Guide',
		description: 'Learn how to conduct a comprehensive technical SEO audit. From crawling to fixing issues, this guide covers everything you need to improve your site\'s search visibility.',
		author: 'Barracuda Team',
		publishDate: '2025-01-20',
		readTime: 12,
		category: 'Guides',
		tags: ['technical SEO', 'SEO audit', 'website audit', 'SEO guide', 'crawling'],
		featured: true
	},
	{
		slug: 'find-fix-broken-links',
		title: 'Broken Links: How to Find and Fix Them (Before They Hurt Your Rankings)',
		description: 'Broken links damage user experience and SEO. Learn how to identify, prioritize, and fix broken links at scale using modern crawling tools.',
		author: 'Barracuda Team',
		publishDate: '2025-01-18',
		readTime: 10,
		category: 'Guides',
		tags: ['broken links', '404 errors', 'link building', 'technical SEO', 'SEO fixes'],
		featured: true
	},
	{
		slug: 'screaming-frog-vs-barracuda',
		title: 'Screaming Frog vs Barracuda: Which SEO Crawler Should You Choose?',
		description: 'A detailed comparison of Screaming Frog and Barracuda SEO crawlers. Discover which tool fits your workflow, budget, and team needs.',
		author: 'Barracuda Team',
		publishDate: '2025-01-15',
		readTime: 8,
		category: 'Comparisons',
		tags: ['screaming frog', 'SEO crawler', 'comparison', 'technical SEO'],
		featured: true
	},
	{
		slug: 'semrush-vs-barracuda',
		title: 'SEMrush vs Barracuda: When to Use a Dedicated SEO Crawler',
		description: 'SEMrush is great for keyword research, but when do you need a dedicated crawler? Compare SEMrush\'s crawl features with Barracuda\'s specialized approach.',
		author: 'Barracuda Team',
		publishDate: '2025-01-12',
		readTime: 9,
		category: 'Comparisons',
		tags: ['SEMrush', 'SEO tools', 'comparison', 'crawling', 'technical SEO'],
		featured: false
	},
	{
		slug: 'automated-seo-audits-cicd',
		title: 'How to Set Up Automated SEO Audits with CI/CD Pipelines',
		description: 'Automate your technical SEO audits by integrating crawlers into your CI/CD workflow. Catch issues before they go live and maintain SEO quality at scale.',
		author: 'Barracuda Team',
		publishDate: '2025-01-10',
		readTime: 11,
		category: 'Automation',
		tags: ['CI/CD', 'automation', 'devops', 'SEO automation', 'technical SEO'],
		featured: false
	},
	{
		slug: 'duplicate-meta-tags-fix',
		title: 'Duplicate Meta Tags: Why They Matter and How to Fix Them at Scale',
		description: 'Duplicate meta tags confuse search engines and hurt rankings. Learn how to identify and fix duplicate title tags and meta descriptions across your entire site.',
		author: 'Barracuda Team',
		publishDate: '2025-01-08',
		readTime: 7,
		category: 'Guides',
		tags: ['duplicate content', 'meta tags', 'title tags', 'SEO fixes', 'on-page SEO'],
		featured: false
	},
	{
		slug: 'redirect-chains-seo-killer',
		title: 'Redirect Chains: The Hidden SEO Killer (And How to Fix Them)',
		description: 'Redirect chains slow down pages and waste crawl budget. Discover how to identify redirect chains and consolidate them into single redirects for better SEO performance.',
		author: 'Barracuda Team',
		publishDate: '2025-01-05',
		readTime: 8,
		category: 'Guides',
		tags: ['redirects', '301 redirects', 'redirect chains', 'technical SEO', 'site speed'],
		featured: false
	},
	{
		slug: 'prioritizing-seo-fixes',
		title: 'Prioritizing SEO Fixes: A Data-Driven Framework',
		description: 'Not all SEO issues are created equal. Learn how to prioritize fixes based on impact, effort, and data to maximize your SEO ROI.',
		author: 'Barracuda Team',
		publishDate: '2025-01-03',
		readTime: 9,
		category: 'Guides',
		tags: ['SEO prioritization', 'SEO strategy', 'technical SEO', 'data-driven SEO'],
		featured: false
	},
	{
		slug: 'audit-large-sites-10000-pages',
		title: 'How to Audit 10,000+ Pages: A Guide for Enterprise SEO',
		description: 'Auditing large websites requires different strategies than small sites. Learn how to crawl, analyze, and fix issues at scale for enterprise-level SEO.',
		author: 'Barracuda Team',
		publishDate: '2025-01-01',
		readTime: 10,
		category: 'Guides',
		tags: ['enterprise SEO', 'large site audit', 'scalable SEO', 'technical SEO'],
		featured: false
	},
	{
		slug: 'visualize-site-structure-link-graph',
		title: 'How to Visualize Your Site Structure: Link Graph Analysis Guide',
		description: 'Understanding your site\'s internal linking structure helps identify orphaned pages, improve crawlability, and optimize information architecture.',
		author: 'Barracuda Team',
		publishDate: '2024-12-29',
		readTime: 8,
		category: 'Guides',
		tags: ['site structure', 'internal linking', 'link graph', 'information architecture', 'crawling'],
		featured: false
	},
	{
		slug: 'seo-audit-checklist',
		title: 'The SEO Audit Checklist Every Agency Should Use',
		description: 'A comprehensive checklist covering all aspects of technical SEO audits. Use this framework to ensure nothing falls through the cracks.',
		author: 'Barracuda Team',
		publishDate: '2024-12-27',
		readTime: 7,
		category: 'Guides',
		tags: ['SEO checklist', 'SEO audit', 'agency SEO', 'technical SEO'],
		featured: false
	},
	{
		slug: 'ecommerce-seo-audit',
		title: 'How to Audit E-commerce Sites: Common Issues and Fixes',
		description: 'E-commerce sites have unique SEO challenges. Learn how to audit product pages, category structures, and technical issues specific to online stores.',
		author: 'Barracuda Team',
		publishDate: '2024-12-25',
		readTime: 11,
		category: 'Guides',
		tags: ['ecommerce SEO', 'product pages', 'category pages', 'technical SEO'],
		featured: false
	}
];

export function getBlogPost(slug: string): BlogPost | undefined {
	return blogPosts.find(post => post.slug === slug);
}

export function getAllBlogPosts(): BlogPost[] {
	return blogPosts.sort((a, b) => 
		new Date(b.publishDate).getTime() - new Date(a.publishDate).getTime()
	);
}

export function getFeaturedPosts(): BlogPost[] {
	return blogPosts.filter(post => post.featured);
}
