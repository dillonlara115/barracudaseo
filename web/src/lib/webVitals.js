// Core Web Vitals thresholds, at the 75th percentile of page loads. These match
// the ones the analyzer uses to raise issues (internal/analyzer/rendering.go).
const THRESHOLDS = {
	lcp: { good: 2500, poor: 4000 },
	inp: { good: 200, poor: 500 },
	cls: { good: 0.1, poor: 0.25 }
};

// rate returns 'good', 'needs-improvement' or 'poor' for a metric value, or
// null when the value is missing.
export function rate(metric, value) {
	if (value == null) return null;
	const { good, poor } = THRESHOLDS[metric];
	if (value <= good) return 'good';
	if (value <= poor) return 'needs-improvement';
	return 'poor';
}

export function ratingClass(rating) {
	if (rating === 'good') return 'text-success';
	if (rating === 'needs-improvement') return 'text-warning';
	if (rating === 'poor') return 'text-error';
	return 'text-base-content/50';
}

export function formatLCP(ms) {
	return ms == null ? '—' : `${(ms / 1000).toFixed(1)} s`;
}

export function formatINP(ms) {
	return ms == null ? '—' : `${Math.round(ms)} ms`;
}

export function formatCLS(value) {
	return value == null ? '—' : value.toFixed(2);
}

// summarizeReading turns one stored reading ({ url, vitals }) into the values
// the dashboard shows. Real-user (field) data is used when Google has it;
// otherwise the lab test stands in, and INP is unavailable.
export function summarizeReading(reading) {
	const vitals = reading?.vitals || {};
	const hasField = Boolean(vitals.field_source);
	const lcp = hasField ? vitals.lcp_ms : vitals.lab_lcp_ms;
	const inp = hasField ? vitals.inp_ms : null;
	const cls = hasField ? vitals.cls : vitals.lab_cls;

	const ratings = [rate('lcp', lcp), rate('inp', inp), rate('cls', cls)].filter(Boolean);
	let overall = null;
	if (ratings.includes('poor')) overall = 'poor';
	else if (ratings.includes('needs-improvement')) overall = 'needs-improvement';
	else if (ratings.length > 0) overall = 'good';

	let source = 'Lab test';
	if (vitals.field_source === 'page') source = 'Real users';
	else if (vitals.field_source === 'origin') source = 'Real users, site-wide';

	return {
		url: reading?.url || '',
		error: vitals.error || '',
		source,
		lcp,
		inp,
		cls,
		score: vitals.performance_score ?? null,
		overall
	};
}
