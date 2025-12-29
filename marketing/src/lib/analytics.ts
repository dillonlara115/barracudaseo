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
