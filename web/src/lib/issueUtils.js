/**
 * Issue types that relate to a specific image/asset URL.
 * For these types, grouping by asset (value) rather than by page is preferred.
 */
export const IMAGE_ISSUE_TYPES = ['missing_image_alt', 'broken_image', 'large_image'];

/**
 * Check if an issue type is asset/image-centric.
 * @param {string} type - Issue type
 * @returns {boolean}
 */
export function isImageIssueType(type) {
	return ['missing_image_alt', 'broken_image', 'large_image'].includes(type || '');
}

/**
 * Heuristic: does the issue have a value that looks like an image URL?
 * Used when type might be missing or non-standard.
 * @param {object} issue - Issue object
 * @returns {boolean}
 */
export function hasImageUrlValue(issue) {
	const v = issue?.value || '';
	return (
		/\.(jpg|jpeg|png|gif|webp|svg|avif)(\?|$)/i.test(v) ||
		/googleusercontent\.com|wp-content\/uploads|cdn\.|cloudinary/i.test(v)
	);
}

/**
 * Collapse image-type issues by asset (value) URL.
 * Returns display items: one per unique (type, value) with aggregated affected page URLs.
 * Non-image issues are passed through as-is (one display item per issue).
 *
 * @param {Array<{type: string, url: string, value?: string, message?: string, recommendation?: string, severity?: string, [key: string]: any}>} issues
 * @returns {Array<{type: string, message: string, recommendation?: string, severity?: string, assetUrl?: string, affectedUrls: string[], isAssetGroup: boolean, [key: string]: any}>}
 */
export function collapseImageIssuesByAsset(issues) {
	if (!Array.isArray(issues) || issues.length === 0) return [];

	const assetMap = new Map(); // key: `${type}|${value}` -> { ..., affectedUrls: Set }
	const nonImageItems = [];

	for (const issue of issues) {
		const type = issue.type || 'unknown';
		const value = (issue.value || '').trim();

		if (isImageIssueType(type) && value) {
			// large_image value is "url (N KB)" - strip for linkable assetUrl
			const linkableUrl = value.replace(/\s*\(\d+\s*KB\)\s*$/i, '').trim() || value;
			const key = `${type}|${value}`;
			if (!assetMap.has(key)) {
				assetMap.set(key, {
					type,
					message: issue.message || '',
					recommendation: issue.recommendation || '',
					severity: issue.severity || 'info',
					assetUrl: value,
					assetLinkUrl: linkableUrl,
					affectedUrls: [],
					// For priority/enrichment: use first page as representative; keep all constituent ids
					url: issue.url || '',
					constituentIds: [],
					isAssetGroup: true
				});
			}
			const entry = assetMap.get(key);
			const pageUrl = issue.url || '';
			if (pageUrl && !entry.affectedUrls.includes(pageUrl)) {
				entry.affectedUrls.push(pageUrl);
			}
			entry.constituentIds.push(`${pageUrl}|${type}`);
			if (!entry.url) entry.url = pageUrl;
		} else {
			nonImageItems.push({
				...issue,
				affectedUrls: issue.url ? [issue.url] : [],
				isAssetGroup: false
			});
		}
	}

	return [...Array.from(assetMap.values()), ...nonImageItems];
}
