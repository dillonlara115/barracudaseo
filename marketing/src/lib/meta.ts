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
	/** Append " - Barracuda SEO" to the title. Defaults to true. */
	withSuffix?: boolean;
}

/** Google truncates title tags at roughly this many characters. */
export const MAX_TITLE_LENGTH = 60;

export function getMetaTags(meta: MetaTags = {}): MetaTagsConfig {
	let title = `${SITE_NAME} - AI-Powered SEO Crawler & Technical Audit Tool`;
	if (meta.title) {
		const suffixed = `${meta.title} - ${SITE_NAME}`;
		// Keep the brand suffix unless it would push the title past the SERP cutoff
		const withSuffix = meta.withSuffix ?? suffixed.length <= MAX_TITLE_LENGTH;
		title = withSuffix ? suffixed : meta.title;
	}

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
		description: SITE_DESCRIPTION,
		url: APP_URL
	};
}

export function getFAQPageSchema(faqs: Array<{ question: string; answer: string }>) {
	return {
		'@context': 'https://schema.org',
		'@type': 'FAQPage',
		mainEntity: faqs.map((faq) => ({
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

export function getArticleSchema(article: {
	title: string;
	description: string;
	author: string;
	publishDate: string;
	updatedDate?: string;
	image?: string;
	url: string;
}) {
	return {
		'@context': 'https://schema.org',
		'@type': 'Article',
		headline: article.title,
		description: article.description,
		author: {
			'@type': 'Person',
			name: article.author
		},
		datePublished: article.publishDate,
		dateModified: article.updatedDate ?? article.publishDate,
		image: article.image || `${SITE_URL}/mockups/barracuda-dashboard.png`,
		url: article.url.startsWith('http') ? article.url : `${SITE_URL}${article.url}`,
		publisher: {
			'@type': 'Organization',
			name: SITE_NAME,
			logo: {
				'@type': 'ImageObject',
				url: `${SITE_URL}/favicon.svg`
			}
		}
	};
}

export function getHowToSchema(howTo: {
	name: string;
	description: string;
	steps: Array<{ name: string; text: string; url?: string }>;
}) {
	return {
		'@context': 'https://schema.org',
		'@type': 'HowTo',
		name: howTo.name,
		description: howTo.description,
		step: howTo.steps.map((step, index) => ({
			'@type': 'HowToStep',
			position: index + 1,
			name: step.name,
			itemListElement: [
				{
					'@type': 'HowToDirection',
					text: step.text
				}
			],
			...(step.url && { url: step.url })
		}))
	};
}
