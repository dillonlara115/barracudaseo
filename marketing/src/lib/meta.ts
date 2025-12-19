import { SITE_NAME, SITE_DESCRIPTION, SITE_URL, APP_URL } from './constants';

export interface MetaTagsConfig {
	title: string;
	description: string;
	ogTitle?: string;
	ogDescription?: string;
	ogImage?: string;
	ogType?: string;
	ogImageWidth?: number;
	ogImageHeight?: number;
	twitterTitle?: string;
	twitterDescription?: string;
	twitterSite?: string;
	twitterCreator?: string;
	keywords?: string;
	author?: string;
	robots?: string;
	structuredData?: Record<string, any> | Record<string, any>[];
}

export interface MetaTags {
	title?: string;
	description?: string;
	ogImage?: string;
	ogType?: string;
}

export function getMetaTags(meta: MetaTags = {}): MetaTagsConfig {
	const title = meta.title 
		? `${meta.title} - ${SITE_NAME}`
		: `${SITE_NAME} - Web-Based SEO Crawler & Auditing Tool`;
	
	const description = meta.description || SITE_DESCRIPTION;

	return {
		title,
		description,
		ogTitle: title,
		ogDescription: description,
		ogImage: meta.ogImage || '/mockups/barracuda-dashboard.png',
		ogType: meta.ogType || 'website'
	};
}

// Structured Data Helpers

export function getOrganizationSchema() {
	return {
		'@context': 'https://schema.org',
		'@type': 'Organization',
		name: SITE_NAME,
		url: SITE_URL,
		logo: `${SITE_URL}/favicon.svg`,
		sameAs: [
			// Add social media URLs here when available
		],
		description: SITE_DESCRIPTION
	};
}

export function getWebSiteSchema() {
	return {
		'@context': 'https://schema.org',
		'@type': 'WebSite',
		name: SITE_NAME,
		url: SITE_URL,
		description: SITE_DESCRIPTION,
		potentialAction: {
			'@type': 'SearchAction',
			target: {
				'@type': 'EntryPoint',
				urlTemplate: `${SITE_URL}/search?q={search_term_string}`
			},
			'query-input': 'required name=search_term_string'
		}
	};
}

export function getSoftwareApplicationSchema() {
	return {
		'@context': 'https://schema.org',
		'@type': 'SoftwareApplication',
		name: SITE_NAME,
		applicationCategory: 'SEO Tool',
		operatingSystem: 'Web, CLI',
		offers: {
			'@type': 'Offer',
			price: '0',
			priceCurrency: 'USD',
			description: 'Free tier available'
		},
		aggregateRating: {
			'@type': 'AggregateRating',
			ratingValue: '4.8',
			ratingCount: '1'
		},
		description: SITE_DESCRIPTION,
		url: APP_URL
	};
}

export function getFAQPageSchema(faqs: Array<{ question: string; answer: string }>) {
	return {
		'@context': 'https://schema.org',
		'@type': 'FAQPage',
		mainEntity: faqs.map(faq => ({
			'@type': 'Question',
			name: faq.question,
			acceptedAnswer: {
				'@type': 'Answer',
				text: faq.answer
			}
		}))
	};
}

export function getBreadcrumbSchema(items: Array<{ name: string; url: string }>) {
	return {
		'@context': 'https://schema.org',
		'@type': 'BreadcrumbList',
		itemListElement: items.map((item, index) => ({
			'@type': 'ListItem',
			position: index + 1,
			name: item.name,
			item: item.url.startsWith('http') ? item.url : `${SITE_URL}${item.url}`
		}))
	};
}
