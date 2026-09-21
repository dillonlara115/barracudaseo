<script lang="ts">
	import {
		MapPin,
		Search,
		Target,
		TrendingUp,
		CircleCheck,
		CircleX,
		CircleQuestionMark,
		ArrowRight
	} from '@lucide/svelte';
	import MetaTags from '../../../components/MetaTags.svelte';
	import OptimizedImage from '../../../components/OptimizedImage.svelte';
	import { getMetaTags, getBreadcrumbSchema, getFAQPageSchema } from '$lib/meta';

	const meta = getMetaTags({
		title: 'Free Local SEO Audit Tool',
		withSuffix: true,
		description:
			'Audit the website side of local SEO: location and service pages, titles, broken links, redirects, speed, orphan pages, and real Search Console data — prioritized. Free for 100 pages, no card.'
	});

	const faqs = [
		{
			question: 'Is there a free local SEO audit tool?',
			answer:
				'Yes. Barracuda crawls up to 100 pages free with no card, which covers most single-location business sites, and shows every issue prioritized in a dashboard. For the listings side, Google Business Profile itself is free, and Google Search Console shows your real local impressions and clicks.'
		},
		{
			question: 'What does a local SEO audit include?',
			answer:
				'Two halves. The website half: location and service pages, titles and descriptions, headings, internal links, broken links and redirects, page speed, and whether Google can crawl and index each page. The listings half: Google Business Profile completeness, categories, reviews, photos, and citation consistency across directories. Barracuda covers the website half; your Google Business Profile dashboard covers the other.'
		},
		{
			question: 'What local SEO issues does Barracuda detect?',
			answer:
				'Missing, short, long or empty page titles and meta descriptions on location and service pages; missing or multiple H1s; broken internal and external links; redirect chains; missing canonical tags; slow-responding pages; oversized images and images without alt text; and location pages that no other page links to. Each is scored by impact and, with Search Console connected, weighted by real impressions.'
		},
		{
			question: 'Does Barracuda check my Google Business Profile, citations or NAP consistency?',
			answer:
				'No. Barracuda audits your website, not your listings. It does not read your Google Business Profile, check directory citations, compare name/address/phone details across the web, or validate LocalBusiness schema. Use your GBP dashboard and a citation tool for those; use Barracuda for the site they all link to.'
		},
		{
			question: 'How often should I audit a local business website?',
			answer:
				'Monthly for most businesses, and immediately after adding location pages, changing hours or addresses on the site, or a redesign. Barracuda keeps crawl history, so each audit shows what changed since the last one instead of the full list again.'
		},
		{
			question: 'Does it work for multi-location businesses?',
			answer:
				'Yes. The Pro plan crawls up to 10,000 pages, the link graph shows which location pages are buried or orphaned, and Search Console data shows which locations earn impressions and which do not. Rank tracking lets you follow "service + city" keywords per location.'
		}
	];

	const structuredData = [
		getBreadcrumbSchema([
			{ name: 'Home', url: '/' },
			{ name: 'Use Cases', url: '/use-cases' },
			{ name: 'Local SEO', url: '/use-cases/local-seo' }
		]),
		getFAQPageSchema(faqs)
	];

	const relatedGuides = [
		{
			href: '/blog/google-business-profile-optimization-the-2026-playbook',
			title: 'Google Business Profile Optimization: The 2026 Playbook'
		},
		{
			href: '/blog/local-pack-ranking-factors-what-the-data-shows-in-2026',
			title: 'Local Pack Ranking Factors: What the Data Shows in 2026'
		},
		{
			href: '/blog/how-to-build-local-citations-that-still-move-the-needle',
			title: 'How to Build Local Citations That Still Move the Needle'
		},
		{
			href: '/blog/multi-location-seo-strategy-managing-5-locations-without-duplicate-content',
			title: 'Multi-Location SEO Without Duplicate Content'
		},
		{
			href: '/blog/hyperlocal-content-strategy-writing-pages-that-rank-for-service-near-me',
			title: 'Hyperlocal Content: Pages That Rank for "Service Near Me"'
		},
		{
			href: '/blog/review-management-for-seo-getting-more-reviews-without-breaking-googles-guidelines',
			title: "Review Management for SEO Without Breaking Google's Guidelines"
		}
	];
</script>

<MetaTags config={{ ...meta, structuredData }} />

<!-- Hero Section -->
<section class="bg-gradient-to-b from-[#3c3836] to-[#2d2826] py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-4xl text-center">
			<h1 class="mb-6 font-heading text-5xl font-bold text-white md:text-6xl">
				Free Local SEO Audit Tool for the Half Most Tools Skip: Your Website
			</h1>
			<p class="mx-auto mb-10 max-w-3xl text-xl text-white/80 md:text-2xl">
				Most local SEO checkers score your Google Business Profile and citations. Barracuda audits
				the site they all point to — every location and service page, prioritized by what actually
				earns local impressions in Search Console. 100 pages free, no card.
			</p>
			<div class="flex flex-col justify-center gap-4 sm:flex-row">
				<a
					href="https://app.barracudaseo.com"
					class="rounded-lg bg-[#8ec07c] px-8 py-4 text-lg font-medium text-[#3c3836] transition-colors hover:bg-[#a0d28c]"
					target="_blank"
					rel="noopener noreferrer"
				>
					Start Your Free Local SEO Audit
				</a>
			</div>
		</div>
	</div>
</section>

<!-- Dashboard Preview -->
<section class="bg-[#2d2826] py-12">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-6xl">
			<OptimizedImage
				src="/mockups/barracuda-dashboard.png"
				alt="Barracuda SEO dashboard showing a prioritized audit of a local business website"
				class="w-full rounded-lg border border-white/10 shadow-2xl"
				loading="lazy"
				aspectRatio="16/9"
			/>
		</div>
	</div>
</section>

<!-- What a local SEO audit covers -->
<section class="bg-[#3c3836] py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-5xl">
			<h2 class="mb-4 text-center font-heading text-4xl font-bold text-white md:text-5xl">
				What a local SEO audit actually covers
			</h2>
			<p class="mx-auto mb-10 max-w-3xl text-center text-xl text-white/70">
				Local search has two halves. The free tools that rank for "local SEO audit" almost all check
				the listings half. Nobody ranks in the local pack with a website Google can't crawl, so here
				is the whole checklist — and which tool does which part.
			</p>
			<div class="overflow-x-auto rounded-lg border border-white/10 bg-[#2d2826]">
				<table class="w-full border-collapse text-left text-sm">
					<thead>
						<tr class="border-b border-white/10">
							<th class="p-4 font-heading text-[#8ec07c]">Check</th>
							<th class="p-4 font-heading text-[#8ec07c]">Why it matters locally</th>
							<th class="p-4 font-heading text-[#8ec07c]">Barracuda</th>
							<th class="p-4 font-heading text-[#8ec07c]">Use instead</th>
						</tr>
					</thead>
					<tbody class="text-white/80">
						<tr class="border-b border-white/10">
							<td class="p-4 font-bold text-white"
								>Location &amp; service page titles and descriptions</td
							>
							<td class="p-4">"Plumber in Boulder" only ranks if the page says so, once, clearly</td
							>
							<td class="p-4"
								><CircleCheck class="inline h-5 w-5 text-[#8ec07c]" /> missing, short, long, empty</td
							>
							<td class="p-4">—</td>
						</tr>
						<tr class="border-b border-white/10">
							<td class="p-4 font-bold text-white">H1 on every location page</td>
							<td class="p-4">Missing or duplicated H1s blur which city a page targets</td>
							<td class="p-4"
								><CircleCheck class="inline h-5 w-5 text-[#8ec07c]" /> missing, multiple, empty</td
							>
							<td class="p-4">—</td>
						</tr>
						<tr class="border-b border-white/10">
							<td class="p-4 font-bold text-white">Orphaned location pages</td>
							<td class="p-4">A city page nothing links to is invisible to Google and customers</td>
							<td class="p-4"
								><CircleCheck class="inline h-5 w-5 text-[#8ec07c]" /> link graph + orphan report</td
							>
							<td class="p-4">—</td>
						</tr>
						<tr class="border-b border-white/10">
							<td class="p-4 font-bold text-white">Broken links, redirect chains, canonicals</td>
							<td class="p-4"
								>Old location URLs and moved service pages leak crawl budget and clicks</td
							>
							<td class="p-4"><CircleCheck class="inline h-5 w-5 text-[#8ec07c]" /></td>
							<td class="p-4">—</td>
						</tr>
						<tr class="border-b border-white/10">
							<td class="p-4 font-bold text-white"
								>Slow pages, oversized images, missing alt text</td
							>
							<td class="p-4">Local searches are mostly mobile; a 4 MB hero photo costs calls</td>
							<td class="p-4"><CircleCheck class="inline h-5 w-5 text-[#8ec07c]" /></td>
							<td class="p-4">PageSpeed Insights for Core Web Vitals detail</td>
						</tr>
						<tr class="border-b border-white/10">
							<td class="p-4 font-bold text-white">Which local pages actually get impressions</td>
							<td class="p-4"
								>Fix the Denver page that gets 2,000 impressions before the one that gets 12</td
							>
							<td class="p-4"
								><CircleCheck class="inline h-5 w-5 text-[#8ec07c]" /> Search Console overlay</td
							>
							<td class="p-4">—</td>
						</tr>
						<tr class="border-b border-white/10">
							<td class="p-4 font-bold text-white">"Service + city" rank tracking</td>
							<td class="p-4">Know whether the fixes moved anything</td>
							<td class="p-4"
								><CircleCheck class="inline h-5 w-5 text-[#8ec07c]" /> scheduled checks (Pro)</td
							>
							<td class="p-4">—</td>
						</tr>
						<tr class="border-b border-white/10">
							<td class="p-4 font-bold text-white"
								>Google Business Profile completeness, categories, photos, Q&amp;A</td
							>
							<td class="p-4">The single biggest local pack factor</td>
							<td class="p-4"><CircleX class="inline h-5 w-5 text-white/40" /> not audited</td>
							<td class="p-4">Your GBP dashboard (free)</td>
						</tr>
						<tr class="border-b border-white/10">
							<td class="p-4 font-bold text-white"
								>Citations and NAP consistency across directories</td
							>
							<td class="p-4">Conflicting phone numbers erode trust signals</td>
							<td class="p-4"><CircleX class="inline h-5 w-5 text-white/40" /> not audited</td>
							<td class="p-4"
								>BrightLocal, Whitespark, or a manual check of the top 10 directories</td
							>
						</tr>
						<tr class="border-b border-white/10">
							<td class="p-4 font-bold text-white">LocalBusiness schema validation</td>
							<td class="p-4">Helps Google connect the page to the listing</td>
							<td class="p-4"><CircleX class="inline h-5 w-5 text-white/40" /> not validated</td>
							<td class="p-4">Google Rich Results Test (free)</td>
						</tr>
						<tr>
							<td class="p-4 font-bold text-white">Reviews volume, velocity and replies</td>
							<td class="p-4">Second biggest local pack factor</td>
							<td class="p-4"><CircleX class="inline h-5 w-5 text-white/40" /> not audited</td>
							<td class="p-4">GBP dashboard; see our review management guide below</td>
						</tr>
					</tbody>
				</table>
			</div>
			<p class="mt-6 text-center text-white/60">
				We would rather show you the gaps than have you find them on day two. For the listings half,
				our
				<a
					href="/blog/local-pack-ranking-factors-what-the-data-shows-in-2026"
					class="text-[#8ec07c] underline hover:text-[#a0d28c]">local pack ranking factors</a
				>
				guide covers what moves the map results.
			</p>
		</div>
	</div>
</section>

<!-- Features Section -->
<section class="bg-[#2d2826] py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-6xl">
			<h2 class="mb-12 text-center font-heading text-4xl font-bold text-white md:text-5xl">
				What Barracuda does for a local business site
			</h2>
			<div class="grid grid-cols-1 gap-8 md:grid-cols-2">
				<div class="rounded-lg border border-white/10 bg-[#3c3836] p-8">
					<div class="mb-4 flex items-center gap-4">
						<div class="rounded-lg bg-[#8ec07c]/10 p-3">
							<MapPin class="h-6 w-6 text-[#8ec07c]" />
						</div>
						<h3 class="font-heading text-2xl font-bold text-white">Every location page, checked</h3>
					</div>
					<p class="mb-4 text-white/80">
						Crawl the whole site and see each city, service and contact page with its title, meta
						description, H1, canonical, links and images checked — so a "Denver" page with a
						"Boulder" title gets caught before a customer does.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Missing, short or over-long titles and descriptions per page</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Missing or multiple H1s on location pages</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Missing canonical tags on pages that need them</span>
						</li>
					</ul>
				</div>

				<div class="rounded-lg border border-white/10 bg-[#3c3836] p-8">
					<div class="mb-4 flex items-center gap-4">
						<div class="rounded-lg bg-[#8ec07c]/10 p-3">
							<Search class="h-6 w-6 text-[#8ec07c]" />
						</div>
						<h3 class="font-heading text-2xl font-bold text-white">Find the buried locations</h3>
					</div>
					<p class="mb-4 text-white/80">
						The interactive link graph shows how your location pages connect. Pages nothing links to
						are flagged as orphans, and pages six clicks from the homepage stand out — the usual
						fate of the fifth location added last year.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Orphaned page report</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Broken internal and external links, redirect chains</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Internal link suggestions between related location and service pages</span>
						</li>
					</ul>
				</div>

				<div class="rounded-lg border border-white/10 bg-[#3c3836] p-8">
					<div class="mb-4 flex items-center gap-4">
						<div class="rounded-lg bg-[#8ec07c]/10 p-3">
							<Target class="h-6 w-6 text-[#8ec07c]" />
						</div>
						<h3 class="font-heading text-2xl font-bold text-white">
							Prioritized by real local traffic
						</h3>
					</div>
					<p class="mb-4 text-white/80">
						Connect Google Search Console and every issue is weighted by the impressions and clicks
						its page already earns. A missing title on the page that ranks for "emergency plumber
						Boulder" outranks the same issue on your privacy policy.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>See which local queries each page appears for</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Declining pages and quick-win keywords surfaced automatically</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Track "service + city" rankings on a schedule (Pro)</span>
						</li>
					</ul>
				</div>

				<div class="rounded-lg border border-white/10 bg-[#3c3836] p-8">
					<div class="mb-4 flex items-center gap-4">
						<div class="rounded-lg bg-[#8ec07c]/10 p-3">
							<TrendingUp class="h-6 w-6 text-[#8ec07c]" />
						</div>
						<h3 class="font-heading text-2xl font-bold text-white">
							Fast pages for mobile searchers
						</h3>
					</div>
					<p class="mb-4 text-white/80">
						Most "near me" searches happen on a phone in a parking lot. Barracuda flags
						slow-responding pages, oversized images and missing alt text — the three fixes that most
						often turn a slow location page into a fast one.
					</p>
					<ul class="space-y-2 text-white/70">
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Slow server response detection per page</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Large image and broken image reports</span>
						</li>
						<li class="flex items-start gap-2">
							<CircleCheck class="mt-0.5 h-5 w-5 flex-shrink-0 text-[#8ec07c]" />
							<span>Microsoft Clarity integration to see where mobile visitors give up</span>
						</li>
					</ul>
				</div>
			</div>
		</div>
	</div>
</section>

<!-- How It Works Section -->
<section class="bg-[#3c3836] py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-4xl">
			<h2 class="mb-4 text-center font-heading text-4xl font-bold text-white md:text-5xl">
				How to run a local SEO audit with Barracuda
			</h2>
			<p class="mx-auto mb-12 max-w-2xl text-center text-xl text-white/70">
				About ten minutes for a single-location site
			</p>

			<div class="space-y-8">
				<div class="flex gap-6">
					<div class="flex-shrink-0">
						<div
							class="flex h-12 w-12 items-center justify-center rounded-full bg-[#8ec07c] text-xl font-bold text-[#3c3836]"
						>
							1
						</div>
					</div>
					<div class="flex-grow">
						<h3 class="mb-3 font-heading text-2xl font-bold text-white">Crawl the site</h3>
						<p class="mb-4 text-white/80">
							Enter your URL. The free plan crawls 100 pages, which covers most single-location
							businesses; Pro crawls 10,000 for multi-location sites and agencies.
						</p>
					</div>
				</div>

				<div class="flex gap-6">
					<div class="flex-shrink-0">
						<div
							class="flex h-12 w-12 items-center justify-center rounded-full bg-[#8ec07c] text-xl font-bold text-[#3c3836]"
						>
							2
						</div>
					</div>
					<div class="flex-grow">
						<h3 class="mb-3 font-heading text-2xl font-bold text-white">
							Work the prioritized list
						</h3>
						<p class="mb-4 text-white/80">
							Issues arrive scored by impact and grouped by type — <a
								href="/blog/are-missing-meta-descriptions-important"
								class="text-[#8ec07c] underline hover:text-[#a0d28c]">missing descriptions</a
							>, broken links, redirect chains, slow pages, orphaned locations — with the reasoning
							shown, so you fix the five that matter this month rather than reading four hundred
							rows.
						</p>
						<div class="mt-6">
							<OptimizedImage
								src="/mockups/barracuda-issues.png"
								alt="Barracuda issues dashboard listing prioritized technical SEO issues on a local business site"
								class="w-full rounded-lg border border-white/10 shadow-lg"
								loading="lazy"
								aspectRatio="16/9"
							/>
						</div>
					</div>
				</div>

				<div class="flex gap-6">
					<div class="flex-shrink-0">
						<div
							class="flex h-12 w-12 items-center justify-center rounded-full bg-[#8ec07c] text-xl font-bold text-[#3c3836]"
						>
							3
						</div>
					</div>
					<div class="flex-grow">
						<h3 class="mb-3 font-heading text-2xl font-bold text-white">
							Connect Google Search Console
						</h3>
						<p class="mb-4 text-white/80">
							One click. From then on, every priority is weighted by the impressions and clicks its
							page earns, and the GSC dashboard shows which location pages are gaining, which are
							declining, and which keywords are one nudge from page one.
						</p>
					</div>
				</div>

				<div class="flex gap-6">
					<div class="flex-shrink-0">
						<div
							class="flex h-12 w-12 items-center justify-center rounded-full bg-[#8ec07c] text-xl font-bold text-[#3c3836]"
						>
							4
						</div>
					</div>
					<div class="flex-grow">
						<h3 class="mb-3 font-heading text-2xl font-bold text-white">
							Fix, then do the listings half
						</h3>
						<p class="mb-4 text-white/80">
							Ship the fixes, then open your Google Business Profile and work through categories,
							photos, hours and reviews — our
							<a
								href="/blog/google-business-profile-optimization-the-2026-playbook"
								class="text-[#8ec07c] underline hover:text-[#a0d28c]">GBP playbook</a
							> covers it. Next month's crawl shows what changed.
						</p>
					</div>
				</div>
			</div>
		</div>
	</div>
</section>

<!-- FAQs Section -->
<section class="bg-[#2d2826] py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-4xl">
			<div class="mb-8 flex items-center justify-center gap-4">
				<div class="rounded-lg bg-[#8ec07c]/10 p-3">
					<CircleQuestionMark class="h-8 w-8 text-[#8ec07c]" />
				</div>
				<h2 class="font-heading text-4xl font-bold text-white md:text-5xl">Local SEO audit FAQs</h2>
			</div>

			<div class="space-y-4">
				{#each faqs as faq}
					<div class="overflow-hidden rounded-lg border border-white/10 bg-[#3c3836]">
						<div class="p-6">
							<h3 class="mb-3 font-heading text-xl font-bold text-white">{faq.question}</h3>
							<p class="text-white/80">{faq.answer}</p>
						</div>
					</div>
				{/each}
			</div>
		</div>
	</div>
</section>

<!-- Related guides -->
<section class="bg-[#3c3836] py-16">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-4xl">
			<h2 class="mb-8 text-center font-heading text-3xl font-bold text-white">
				Local SEO guides from the Barracuda blog
			</h2>
			<ul class="grid grid-cols-1 gap-4 md:grid-cols-2">
				{#each relatedGuides as guide}
					<li>
						<a
							href={guide.href}
							class="block rounded-lg border border-white/10 bg-[#2d2826] p-5 text-white/90 transition-colors hover:border-[#8ec07c]/50 hover:text-white"
						>
							{guide.title}
						</a>
					</li>
				{/each}
			</ul>
		</div>
	</div>
</section>

<!-- CTA Section -->
<section class="bg-gradient-to-b from-[#2d2826] to-[#3c3836] py-20">
	<div class="container mx-auto px-4">
		<div class="mx-auto max-w-3xl text-center">
			<h2 class="mb-6 font-heading text-4xl font-bold text-white md:text-5xl">
				Audit the website half in ten minutes
			</h2>
			<p class="mb-10 text-xl text-white/80">
				100 pages free, no card. See every location page's issues, ranked by what they cost you in
				local search.
			</p>
			<a
				href="https://app.barracudaseo.com"
				class="inline-flex items-center gap-2 rounded-lg bg-[#8ec07c] px-8 py-4 text-lg font-medium text-[#3c3836] transition-colors hover:bg-[#a0d28c]"
				target="_blank"
				rel="noopener noreferrer"
			>
				Start Your Free Local SEO Audit
				<ArrowRight class="h-5 w-5" />
			</a>
		</div>
	</div>
</section>
