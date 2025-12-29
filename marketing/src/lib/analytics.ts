/**
 * Google Analytics utility functions
 */

declare global {
	interface Window {
		dataLayer: any[];
		gtag: (...args: any[]) => void;
	}
}

const GA_MEASUREMENT_ID = 'G-B9442HLM46';

/**
 * Track a page view in Google Analytics
 */
export function trackPageView(url: string, title?: string) {
	if (typeof window === 'undefined' || !window.gtag) {
		return;
	}

	window.gtag('config', GA_MEASUREMENT_ID, {
		page_path: url,
		page_title: title
	});
}

/**
 * Track a custom event in Google Analytics
 */
export function trackEvent(
	eventName: string,
	eventParams?: {
		action?: string;
		category?: string;
		label?: string;
		value?: number;
		[key: string]: any;
	}
) {
	if (typeof window === 'undefined' || !window.gtag) {
		return;
	}

	window.gtag('event', eventName, eventParams);
}

/**
 * Track a CTA click (conversion event)
 */
export function trackCTA(
	ctaName: string,
	options?: {
		location?: string;
		page?: string;
		plan?: string;
	}
) {
	trackEvent('cta_click', {
		category: 'engagement',
		action: 'click',
		label: ctaName,
		cta_name: ctaName,
		location: options?.location || 'unknown',
		page: options?.page || (typeof window !== 'undefined' ? window.location.pathname : 'unknown'),
		plan: options?.plan
	});
}

/**
 * Track a signup link click (primary conversion)
 */
export function trackSignup(
	options?: {
		source?: string;
		location?: string;
		plan?: string;
	}
) {
	trackEvent('signup_click', {
		category: 'conversion',
		action: 'signup_click',
		label: options?.source || 'unknown',
		source: options?.source || 'unknown',
		location: options?.location || 'unknown',
		page: typeof window !== 'undefined' ? window.location.pathname : 'unknown',
		plan: options?.plan
	});
}

/**
 * Track pricing page interaction
 */
export function trackPricingAction(
	action: 'view' | 'cta_click' | 'plan_select',
	options?: {
		plan?: string;
		cta_name?: string;
	}
) {
	trackEvent('pricing_interaction', {
		category: 'engagement',
		action,
		label: options?.plan || options?.cta_name || action,
		plan: options?.plan,
		cta_name: options?.cta_name
	});
}

/**
 * Track outbound link click
 */
export function trackOutboundLink(url: string, linkText?: string) {
	trackEvent('outbound_click', {
		category: 'engagement',
		action: 'click',
		label: linkText || url,
		link_url: url,
		link_text: linkText
	});
}
