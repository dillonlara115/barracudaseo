// Helpers for showing the progress of a running crawl.
//
// A crawl goes through these phases in order. Scanning is by far the longest,
// so it gets most of the progress bar; the later phases share the rest.
export const PHASES = [
	{ key: 'scanning', label: 'Scanning pages' },
	{ key: 'page_experience', label: 'Measuring page speed and rendering' },
	{ key: 'metadata_review', label: 'Reviewing metadata' },
	{ key: 'image_analysis', label: 'Analyzing images' },
	{ key: 'storing', label: 'Storing results' }
];

// Share of the progress bar (in percent) that scanning fills
const SCAN_SHARE = 80;

export function phaseLabel(phase) {
	return PHASES.find((p) => p.key === phase)?.label || '';
}

// phaseStep returns the 1-based position of a phase, or 0 if it is unknown.
export function phaseStep(phase) {
	return PHASES.findIndex((p) => p.key === phase) + 1;
}

export function isScanning(phase) {
	return !phase || phase === 'scanning';
}

// estimateTotal returns how many pages the crawl is expected to end with, or
// null when nothing is known yet. `discovered` is the crawler's count of pages
// found so far; the site may turn out bigger while links are still being
// followed, and it never exceeds the crawl's page limit.
export function estimateTotal({ crawled = 0, discovered = 0, maxPages = 0 }) {
	let total = Math.max(crawled, discovered || 0);
	if (maxPages > 0 && total > maxPages) total = maxPages;
	return total > 0 ? total : null;
}

// progressPercent maps the crawl's position to 0-100 across all phases.
export function progressPercent({ phase, crawled = 0, total = null }) {
	if (isScanning(phase)) {
		if (!total) return 0;
		return Math.min(crawled / total, 1) * SCAN_SHARE;
	}
	const step = phaseStep(phase);
	if (step === 0) return SCAN_SHARE;
	// Later phases: the middle of each one's slice of the remaining bar
	const later = PHASES.length - 1;
	return SCAN_SHARE + ((step - 1 - 0.5) / later) * (100 - SCAN_SHARE);
}

// rollingRate returns pages per second over the most recent window, from
// samples of { time (ms), count }. It returns null until there is enough data,
// so a slow start does not drag the rate down for the rest of the crawl.
export function rollingRate(samples, windowMs = 20000) {
	if (!samples || samples.length < 2) return null;
	const latest = samples[samples.length - 1];
	const earliest = samples.find((s) => latest.time - s.time <= windowMs) || samples[0];
	const seconds = (latest.time - earliest.time) / 1000;
	const pages = latest.count - earliest.count;
	if (seconds < 3 || pages <= 0) return null;
	return pages / seconds;
}

// secondsRemaining estimates the time left in the scanning phase, or null
// when it cannot be estimated.
export function secondsRemaining({ crawled = 0, total = null, rate = null }) {
	if (!rate || !total || total <= crawled) return null;
	return (total - crawled) / rate;
}

export function formatDuration(seconds) {
	if (seconds == null || seconds < 0) return null;
	if (seconds < 60) return `${Math.max(1, Math.round(seconds))}s`;
	if (seconds < 3600) return `${Math.round(seconds / 60)}m`;
	const hours = Math.floor(seconds / 3600);
	const minutes = Math.round((seconds % 3600) / 60);
	return minutes > 0 ? `${hours}h ${minutes}m` : `${hours}h`;
}
