<script lang="ts">
	import { trackSignup, trackCTA, trackOutboundLink, trackPricingAction } from '$lib/analytics';

	interface Props {
		href?: string;
		trackAs?: 'signup' | 'cta' | 'pricing' | 'outbound';
		ctaName?: string;
		location?: string;
		plan?: string;
		source?: string;
		class?: string;
		target?: string;
		rel?: string;
		children?: any;
		[key: string]: any;
	}

	let {
		href,
		trackAs = 'outbound',
		ctaName,
		location,
		plan,
		source,
		class: className,
		target,
		rel,
		children,
		...restProps
	}: Props = $props();

	function handleClick(e: MouseEvent) {
		// Track based on type
		if (trackAs === 'signup') {
			trackSignup({
				source: source || ctaName || 'unknown',
				location: location || 'unknown',
				plan
			});
		} else if (trackAs === 'cta') {
			trackCTA(ctaName || href || 'unknown', {
				location,
				plan
			});
		} else if (trackAs === 'pricing') {
			trackPricingAction('cta_click', {
				plan,
				cta_name: ctaName
			});
		} else if (trackAs === 'outbound' && href) {
			trackOutboundLink(href, ctaName);
		}
	}
</script>

<a
	{...restProps}
	{href}
	class={className}
	{target}
	{rel}
	onclick={handleClick}
>
	{@render children()}
</a>
