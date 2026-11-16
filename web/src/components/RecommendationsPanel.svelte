<script>
  export let issues = [];
  export let navigateToTab = null;
  export let enrichedIssues = {}; // Map of enriched issue data

  // Recommendation data for each issue type
  const recommendations = {
    missing_h1: {
      title: "Add H1 Heading",
      impact: "High",
      description: "Each page should have exactly one H1 heading that clearly describes the page content.",
      codeSnippet: `<h1>Your Page Title Here</h1>`,
      explanation: "Place the H1 near the top of your main content, ideally within the first 100 words.",
      resources: [
        { name: "H1 Tag Best Practices", url: "https://moz.com/learn/seo/h1-tag" },
        { name: "Google H1 Guidelines", url: "https://developers.google.com/search/docs/appearance/heading-titles" }
      ]
    },
    missing_title: {
      title: "Add Page Title",
      impact: "Critical",
      description: "Every page must have a unique, descriptive title tag for SEO and user experience.",
      codeSnippet: `<title>Your Page Title - Your Brand Name</title>`,
      explanation: "Place in the <head> section. Keep it under 60 characters for optimal display in search results.",
      resources: [
        { name: "Title Tag Best Practices", url: "https://moz.com/learn/seo/title-tag" },
        { name: "Google Title Guidelines", url: "https://developers.google.com/search/docs/appearance/title-link" }
      ]
    },
    missing_meta_description: {
      title: "Add Meta Description",
      impact: "Medium",
      description: "Meta descriptions improve click-through rates from search results.",
      codeSnippet: `<meta name="description" content="A compelling description of your page content (150-160 characters).">`,
      explanation: "Place in the <head> section. Write compelling copy that encourages clicks while accurately describing the content.",
      resources: [
        { name: "Meta Description Guide", url: "https://moz.com/learn/seo/meta-description" },
        { name: "Writing Effective Descriptions", url: "https://developers.google.com/search/docs/appearance/snippet" }
      ]
    },
    long_title: {
      title: "Shorten Page Title",
      impact: "Medium",
      description: "Titles over 60 characters may be truncated in search results.",
      codeSnippet: `<title>Shorter, More Focused Title</title>`,
      explanation: "Keep the most important keywords at the beginning. Aim for 50-60 characters.",
      resources: [
        { name: "Title Length Guidelines", url: "https://moz.com/learn/seo/title-tag" }
      ]
    },
    short_title: {
      title: "Expand Page Title",
      impact: "Low",
      description: "Very short titles may not provide enough context for search engines and users.",
      codeSnippet: `<title>Descriptive Title That Provides Context</title>`,
      explanation: "Include primary keywords and brand name. Aim for at least 30 characters.",
      resources: [
        { name: "Title Tag Best Practices", url: "https://moz.com/learn/seo/title-tag" }
      ]
    },
    long_meta_description: {
      title: "Shorten Meta Description",
      impact: "Low",
      description: "Meta descriptions over 160 characters may be truncated in search results.",
      codeSnippet: `<meta name="description" content="Concise description within 150-160 characters.">`,
      explanation: "Keep the most compelling information at the beginning. Aim for 150-160 characters.",
      resources: [
        { name: "Meta Description Best Practices", url: "https://moz.com/learn/seo/meta-description" }
      ]
    },
    short_meta_description: {
      title: "Expand Meta Description",
      impact: "Low",
      description: "Very short meta descriptions may not be compelling enough to encourage clicks.",
      codeSnippet: `<meta name="description" content="A more detailed description that provides value and encourages clicks.">`,
      explanation: "Include key benefits and a call-to-action. Aim for at least 120 characters.",
      resources: [
        { name: "Meta Description Guide", url: "https://moz.com/learn/seo/meta-description" }
      ]
    },
    multiple_h1: {
      title: "Use Single H1 Per Page",
      impact: "Medium",
      description: "Having multiple H1 tags can confuse search engines about your page's main topic.",
      codeSnippet: `<h1>Main Page Heading</h1>
<h2>Subsection Heading</h2>
<h2>Another Subsection</h2>`,
      explanation: "Use only one H1 for the main page heading. Use H2-H6 for subsections.",
      resources: [
        { name: "H1 Tag Best Practices", url: "https://moz.com/learn/seo/h1-tag" }
      ]
    },
    empty_h1: {
      title: "Add Content to H1 Tag",
      impact: "High",
      description: "Empty H1 tags provide no value to users or search engines.",
      codeSnippet: `<h1>Meaningful Heading Text</h1>`,
      explanation: "Ensure your H1 contains descriptive text that summarizes the page content.",
      resources: [
        { name: "H1 Tag Guidelines", url: "https://moz.com/learn/seo/h1-tag" }
      ]
    },
    missing_image_alt: {
      title: "Add Alt Text to Images",
      impact: "Medium",
      description: "Alt text improves accessibility and helps search engines understand images.",
      codeSnippet: `<img src="image.jpg" alt="Descriptive text explaining what the image shows">`,
      explanation: "Write concise, descriptive alt text that explains the image's content and purpose.",
      resources: [
        { name: "Alt Text Best Practices", url: "https://moz.com/learn/seo/alt-text" },
        { name: "Accessibility Guidelines", url: "https://www.w3.org/WAI/tutorials/images/" }
      ]
    },
    broken_image: {
      title: "Fix Broken Images",
      impact: "High",
      description: "Images that return 404 or other errors hurt user experience and may indicate missing or moved assets.",
      codeSnippet: `<!-- Update image URL to valid path -->
<img src="/correct-image-path.jpg" alt="Image description">

<!-- Or remove broken image if no longer needed -->
<!-- Or use placeholder while fixing -->`,
      explanation: "Update image src to point to a valid URL, remove the image if obsolete, or replace with a working alternative.",
      resources: [
        { name: "Image Best Practices", url: "https://web.dev/fast/#optimize-your-images" }
      ]
    },
    large_image: {
      title: "Optimize Image Size",
      impact: "Medium",
      description: "Large images slow down page load times, impacting user experience and SEO.",
      codeSnippet: `<!-- Optimize images before uploading -->
<img src="optimized-image.jpg" 
     srcset="image-400w.jpg 400w, image-800w.jpg 800w"
     sizes="(max-width: 400px) 400px, 800px"
     alt="Image description">`,
      explanation: "Compress images, use modern formats (WebP), and implement responsive images with srcset.",
      resources: [
        { name: "Image Optimization Guide", url: "https://web.dev/fast/#optimize-your-images" },
        { name: "PageSpeed Insights", url: "https://pagespeed.web.dev/" }
      ]
    },
    slow_response: {
      title: "Improve Page Speed",
      impact: "High",
      description: "Slow-loading pages hurt user experience and search rankings.",
      codeSnippet: `<!-- Enable compression -->
<!-- Use CDN for static assets -->
<!-- Minimize HTTP requests -->
<!-- Optimize server response time -->`,
      explanation: "Optimize images, enable caching, use a CDN, minimize JavaScript/CSS, and improve server response time.",
      resources: [
        { name: "PageSpeed Optimization", url: "https://web.dev/fast/" },
        { name: "Core Web Vitals", url: "https://web.dev/vitals/" },
        { name: "Google PageSpeed Insights", url: "https://pagespeed.web.dev/" }
      ]
    },
    redirect_chain: {
      title: "Simplify Redirect Chains",
      impact: "Medium",
      description: "Redirect chains waste crawl budget and slow down page loads. The issue URL is the source (the one that redirects)—test that URL with curl -I -L to verify; the destination URL will not redirect.",
      codeSnippet: `<!-- Update links to point directly to final destination -->
<a href="https://example.com/final-page">Link Text</a>

<!-- Or configure server to redirect directly -->
<!-- Apache: Redirect 301 /old-page /final-page -->
<!-- Nginx: return 301 /final-page; -->`,
      explanation: "Update internal links to point directly to the final destination. Use 301 redirects for permanent moves.",
      resources: [
        { name: "Redirect Best Practices", url: "https://moz.com/learn/seo/redirection" },
        { name: "HTTP Status Codes", url: "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status" }
      ]
    },
    no_canonical: {
      title: "Add Canonical Tag",
      impact: "Low",
      description: "Canonical tags help prevent duplicate content issues.",
      codeSnippet: `<link rel="canonical" href="https://example.com/canonical-page-url">`,
      explanation: "Place in the <head> section. Point to the preferred version of the page (usually self-referential).",
      resources: [
        { name: "Canonical Tag Guide", url: "https://moz.com/learn/seo/canonicalization" },
        { name: "Google Canonical Guidelines", url: "https://developers.google.com/search/docs/crawling-indexing/consolidate-duplicate-urls" }
      ]
    },
    broken_link: {
      title: "Fix Broken Links",
      impact: "Medium",
      description: "Broken links hurt user experience and may indicate crawlability issues.",
      codeSnippet: `<!-- Update link to correct URL -->
<a href="https://example.com/correct-page">Link Text</a>

<!-- Or remove if page no longer exists -->
<!-- Or redirect to relevant alternative page -->`,
      explanation: "Update links to point to valid pages, remove dead links, or redirect to relevant alternatives.",
      resources: [
        { name: "Link Building Best Practices", url: "https://moz.com/learn/seo/internal-link" },
        { name: "HTTP Status Codes", url: "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status" }
      ]
    },
    duplicate_title: {
      title: "Write Unique Page Titles",
      impact: "High",
      description: "Several pages share the same title, so search engines cannot tell which one to rank for a query.",
      codeSnippet: `<!-- Before: the same title on every page -->
<title>Our Services</title>

<!-- After: each title describes its own page -->
<title>Roof Repair in Austin | Brand</title>
<title>Roof Replacement in Austin | Brand</title>`,
      explanation: "Give every indexable page a title that names its specific topic. If two pages really are the same content, consolidate them or set a canonical.",
      resources: [
        { name: "Google Title Link Guidelines", url: "https://developers.google.com/search/docs/appearance/title-link" },
        { name: "Consolidate Duplicate URLs", url: "https://developers.google.com/search/docs/crawling-indexing/consolidate-duplicate-urls" }
      ]
    },
    duplicate_meta_description: {
      title: "Write Unique Meta Descriptions",
      impact: "Medium",
      description: "Several pages share the same meta description, which makes their search snippets indistinguishable.",
      codeSnippet: `<meta name="description" content="A summary specific to this page, 120-160 characters.">`,
      explanation: "Describe what is unique about each page. Templates that fill in the page's own product, location or topic are fine.",
      resources: [
        { name: "Google Snippet Guidelines", url: "https://developers.google.com/search/docs/appearance/snippet" }
      ]
    },
    duplicate_h1: {
      title: "Use Distinct H1 Headings",
      impact: "Low",
      description: "Several pages share the same H1, which weakens the signal of what each page is about.",
      codeSnippet: `<h1>A heading that names this page's topic</h1>`,
      explanation: "Check that the H1 is page content rather than a site-wide element such as the logo or site name.",
      resources: [
        { name: "H1 Tag Best Practices", url: "https://moz.com/learn/seo/h1-tag" }
      ]
    },
    orphan_page: {
      title: "Link to Orphan Pages",
      impact: "High",
      description: "No crawled page links to these pages. They were found through the sitemap only, so users cannot reach them and search engines give them little weight.",
      codeSnippet: `<!-- Add a contextual link from a related page -->
<a href="/orphan-page">Descriptive anchor text</a>`,
      explanation: "Link to each page from a relevant page, hub or navigation. If the page is no longer needed, remove it from the sitemap and redirect or delete it. If the crawl stopped at its page limit, re-crawl with a higher limit to confirm.",
      resources: [
        { name: "Internal Linking Guide", url: "https://moz.com/learn/seo/internal-link" },
        { name: "Google Link Best Practices", url: "https://developers.google.com/search/docs/crawling-indexing/links-crawlable" }
      ]
    },
    broken_internal_link: {
      title: "Fix Links to Broken Pages",
      impact: "High",
      description: "These pages contain links to URLs on your site that return an error.",
      codeSnippet: `<!-- Point the link at a working URL -->
<a href="/working-page">Link Text</a>`,
      explanation: "Update each link to the correct URL, or remove it. If the target moved, also add a redirect so external links keep working.",
      resources: [
        { name: "HTTP Status Codes", url: "https://developer.mozilla.org/en-US/docs/Web/HTTP/Status" }
      ]
    },
    internal_link_to_redirect: {
      title: "Link Directly to Final URLs",
      impact: "Low",
      description: "These pages link to URLs that redirect. Each one adds an extra hop for users and crawlers.",
      codeSnippet: `<!-- Before -->
<a href="/old-url">Link Text</a>

<!-- After -->
<a href="/final-url/">Link Text</a>`,
      explanation: "Update the link to the redirect's destination. This is often a navigation or footer link, so one template fix clears many pages.",
      resources: [
        { name: "Redirect Best Practices", url: "https://moz.com/learn/seo/redirection" }
      ]
    },
    invalid_structured_data: {
      title: "Fix Unreadable Structured Data",
      impact: "High",
      description: "These pages contain JSON-LD that is not valid JSON or has no @type. Search engines ignore a block they cannot parse, so the page gets no rich results from it.",
      codeSnippet: `<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@type": "Organization",
  "name": "Example"
}
<\/script>`,
      explanation: "The usual causes are a trailing comma, an unescaped quote in a text value, or a template that printed nothing. Validate the block after fixing it.",
      resources: [
        { name: "Rich Results Test", url: "https://search.google.com/test/rich-results" },
        { name: "Schema Markup Validator", url: "https://validator.schema.org/" }
      ]
    },
    structured_data_missing_fields: {
      title: "Complete Structured Data",
      impact: "Medium",
      description: "These pages declare a Schema.org type but leave out properties that are required before the markup can produce a rich result.",
      codeSnippet: `{
  "@context": "https://schema.org",
  "@type": "Product",
  "name": "Widget",
  "offers": { "@type": "Offer", "price": "9.99", "priceCurrency": "USD" }
}`,
      explanation: "Add the listed properties using values that are visible on the page. Markup that describes content users cannot see violates Google's guidelines.",
      resources: [
        { name: "Structured Data Search Gallery", url: "https://developers.google.com/search/docs/appearance/structured-data/search-gallery" },
        { name: "Rich Results Test", url: "https://search.google.com/test/rich-results" }
      ]
    },
    hreflang_invalid_language: {
      title: "Fix Invalid hreflang Values",
      impact: "High",
      description: "These pages use hreflang values that are not valid language or region codes, so search engines ignore the annotation.",
      codeSnippet: `<!-- Wrong -->
<link rel="alternate" hreflang="en_US" href="https://example.com/us/">
<link rel="alternate" hreflang="en-UK" href="https://example.com/uk/">

<!-- Right -->
<link rel="alternate" hreflang="en-US" href="https://example.com/us/">
<link rel="alternate" hreflang="en-GB" href="https://example.com/uk/">`,
      explanation: "Use an ISO 639-1 language code, optionally followed by a hyphen and an ISO 3166-1 region. Use a hyphen, not an underscore, and GB rather than UK.",
      resources: [
        { name: "Google hreflang Guide", url: "https://developers.google.com/search/docs/specialty/international/localized-versions" }
      ]
    },
    hreflang_missing_return_link: {
      title: "Add hreflang Return Links",
      impact: "High",
      description: "These pages point to an alternate language version that does not point back. hreflang only works when both pages reference each other.",
      codeSnippet: `<!-- On the English page AND on the German page -->
<link rel="alternate" hreflang="en" href="https://example.com/en/page">
<link rel="alternate" hreflang="de" href="https://example.com/de/page">`,
      explanation: "Every page in a language set should carry the same full list of alternates, including itself. Generating the list from one source for all versions prevents this.",
      resources: [
        { name: "Google hreflang Guide", url: "https://developers.google.com/search/docs/specialty/international/localized-versions" }
      ]
    },
    hreflang_to_non_indexable: {
      title: "Point hreflang at Indexable URLs",
      impact: "High",
      description: "These pages list an alternate language URL that redirects, is broken, is noindexed or is canonicalised elsewhere.",
      codeSnippet: `<link rel="alternate" hreflang="de" href="https://example.com/de/final-indexable-url">`,
      explanation: "Each hreflang URL must be the final, indexable address of that language version. Update the annotation, or fix the target page.",
      resources: [
        { name: "Google hreflang Guide", url: "https://developers.google.com/search/docs/specialty/international/localized-versions" }
      ]
    },
    hreflang_missing_self_reference: {
      title: "Add a Self-Referencing hreflang",
      impact: "Medium",
      description: "These pages list their alternate versions but leave themselves out of the set.",
      codeSnippet: `<!-- On https://example.com/en/page -->
<link rel="alternate" hreflang="en" href="https://example.com/en/page">
<link rel="alternate" hreflang="de" href="https://example.com/de/page">`,
      explanation: "Include an entry for the page's own language pointing at its own URL, alongside the alternates.",
      resources: [
        { name: "Google hreflang Guide", url: "https://developers.google.com/search/docs/specialty/international/localized-versions" }
      ]
    },
    hreflang_missing_x_default: {
      title: "Add an x-default hreflang",
      impact: "Low",
      description: "These language sets have no x-default, the version shown to users whose language matches none of the alternates.",
      codeSnippet: `<link rel="alternate" hreflang="x-default" href="https://example.com/">`,
      explanation: "Point x-default at your language selector or your most general version. It is optional but recommended.",
      resources: [
        { name: "Google hreflang Guide", url: "https://developers.google.com/search/docs/specialty/international/localized-versions" }
      ]
    },
    mixed_content: {
      title: "Remove Mixed Content",
      impact: "Medium",
      description: "These HTTPS pages load images, scripts, styles or frames over plain HTTP. Browsers block or warn about them, which can break the page and removes the secure padlock.",
      codeSnippet: `<!-- Before -->
<img src="http://example.com/photo.jpg">

<!-- After -->
<img src="https://example.com/photo.jpg">`,
      explanation: "Change each resource URL to HTTPS. If a third party does not offer HTTPS, host the file yourself or replace it.",
      resources: [
        { name: "Fixing Mixed Content", url: "https://web.dev/articles/fixing-mixed-content" }
      ]
    },
    insecure_page: {
      title: "Serve Pages over HTTPS",
      impact: "High",
      description: "These pages are served over plain HTTP. Browsers mark them as not secure and HTTPS is a ranking signal.",
      codeSnippet: `# Redirect all HTTP requests to HTTPS (nginx)
server {
  listen 80;
  return 301 https://$host$request_uri;
}`,
      explanation: "Install a certificate, redirect every HTTP URL to its HTTPS equivalent with a 301, and update internal links and canonicals to HTTPS.",
      resources: [
        { name: "Why HTTPS Matters", url: "https://web.dev/articles/why-https-matters" }
      ]
    },
    missing_security_header: {
      title: "Add Security Headers",
      impact: "Low",
      description: "The site does not send one or more recommended security headers. These protect visitors; they are not a direct ranking factor.",
      codeSnippet: `Strict-Transport-Security: max-age=31536000; includeSubDomains
X-Content-Type-Options: nosniff
X-Frame-Options: SAMEORIGIN
Referrer-Policy: strict-origin-when-cross-origin
Content-Security-Policy: default-src 'self'`,
      explanation: "Set these once at the server, CDN or hosting platform and they apply to every page. Test a Content-Security-Policy in report-only mode first, since a strict one can block legitimate scripts.",
      resources: [
        { name: "MDN HTTP Security Headers", url: "https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers#security" }
      ]
    },
    thin_content: {
      title: "Expand Thin Pages",
      impact: "Medium",
      description: "These pages have very little main text once navigation, header and footer are set aside. Thin pages rarely rank and can drag down how the whole site is assessed.",
      codeSnippet: `<main>
  <h1>Topic of the page</h1>
  <p>Enough original content to fully answer what a visitor came for...</p>
</main>`,
      explanation: "Add content that covers the topic properly, merge the page into a stronger one, or noindex it if it is not meant to rank. Word count is a signal, not a target: a short page that fully answers its query is fine.",
      resources: [
        { name: "Creating Helpful Content", url: "https://developers.google.com/search/docs/fundamentals/creating-helpful-content" }
      ]
    },
    duplicate_content: {
      title: "Consolidate Duplicate Pages",
      impact: "High",
      description: "These pages have exactly the same main content as another page on the site, so they compete with each other in search.",
      codeSnippet: `<!-- On each duplicate, point to the version you want indexed -->
<link rel="canonical" href="https://example.com/preferred-page">`,
      explanation: "Pick one version to keep. Redirect the others to it with a 301, or set their canonical to it if they must stay reachable (for example URL-parameter variants).",
      resources: [
        { name: "Consolidate Duplicate URLs", url: "https://developers.google.com/search/docs/crawling-indexing/consolidate-duplicate-urls" }
      ]
    },
    near_duplicate_content: {
      title: "Differentiate Near-Duplicate Pages",
      impact: "Medium",
      description: "These pages share almost all of their main content with another page. This is common with location or product-variant pages built from one template.",
      codeSnippet: `<!-- Each page needs content the others do not have -->
<h1>Roof Repair in Austin</h1>
<p>Details, examples and FAQs specific to Austin...</p>`,
      explanation: "Give each page substantial content of its own, or consolidate the set into one page. Swapping a single word such as a city name is not enough for search engines to treat pages as distinct.",
      resources: [
        { name: "Spam Policies: Doorway and Scaled Content", url: "https://developers.google.com/search/docs/essentials/spam-policies" }
      ]
    },
    noindex_page: {
      title: "Review Noindex Pages",
      impact: "Medium",
      description: "These pages carry a noindex directive, so search engines will not show them. This is often intentional, but an accidental noindex silently removes a page from search.",
      codeSnippet: `<!-- In the HTML head -->
<meta name="robots" content="noindex">

<!-- Or as an HTTP response header -->
X-Robots-Tag: noindex`,
      explanation: "Check each page against what you expect to rank. To restore a page, remove the directive from the HTML and from any server or CDN header rule.",
      resources: [
        { name: "Block Indexing with noindex", url: "https://developers.google.com/search/docs/crawling-indexing/block-indexing" }
      ]
    },
    blocked_by_robots: {
      title: "Review Pages Blocked by robots.txt",
      impact: "Medium",
      description: "These pages are disallowed in robots.txt. Search engines cannot crawl them, but can still index the bare URL if other pages link to it.",
      codeSnippet: `# robots.txt
User-agent: *
Disallow: /private/`,
      explanation: "Keep the rule if the section should not be crawled. To keep a page out of search results, allow crawling and use noindex instead; a blocked page's noindex is never seen.",
      resources: [
        { name: "Google robots.txt Introduction", url: "https://developers.google.com/search/docs/crawling-indexing/robots/intro" }
      ]
    },
    canonical_points_elsewhere: {
      title: "Confirm Canonicalised Pages",
      impact: "Low",
      description: "These pages declare another URL as canonical, so they will not be indexed themselves.",
      codeSnippet: `<link rel="canonical" href="https://example.com/preferred-url">`,
      explanation: "This is correct for true duplicates such as tracking-parameter or filtered URLs. If a page has its own content and should rank, make its canonical self-referencing.",
      resources: [
        { name: "Consolidate Duplicate URLs", url: "https://developers.google.com/search/docs/crawling-indexing/consolidate-duplicate-urls" }
      ]
    },
    canonical_to_non_indexable: {
      title: "Fix Canonicals Pointing to Non-Indexable URLs",
      impact: "High",
      description: "These pages name a canonical URL that is broken, redirects, is noindexed, or is itself canonicalised elsewhere. Search engines may ignore the canonical or drop both pages.",
      codeSnippet: `<!-- The canonical target must return 200 and be indexable -->
<link rel="canonical" href="https://example.com/working-indexable-url">`,
      explanation: "Point the canonical at the final, indexable URL. This is usually a template or plugin setting, so one fix often clears many pages.",
      resources: [
        { name: "Canonicalization Troubleshooting", url: "https://developers.google.com/search/docs/crawling-indexing/canonicalization-troubleshooting" }
      ]
    },
    non_indexable_in_sitemap: {
      title: "Clean Up the Sitemap",
      impact: "Medium",
      description: "The sitemap lists URLs that cannot be indexed: redirects, errors, noindexed or canonicalised pages.",
      codeSnippet: `<!-- List only final, indexable, canonical URLs -->
<url>
  <loc>https://example.com/indexable-page</loc>
</url>`,
      explanation: "Either fix the page so it is indexable, or remove it from the sitemap. Mixed signals waste crawl budget and make sitemap reports harder to read.",
      resources: [
        { name: "Build and Submit a Sitemap", url: "https://developers.google.com/search/docs/crawling-indexing/sitemaps/build-sitemap" }
      ]
    },
    indexable_not_in_sitemap: {
      title: "Add Indexable Pages to the Sitemap",
      impact: "Low",
      description: "These pages are indexable and linked on the site but missing from the sitemap.",
      codeSnippet: `<url>
  <loc>https://example.com/missing-page</loc>
</url>`,
      explanation: "Add them so search engines discover changes faster. If a page should not be in search, noindex it instead.",
      resources: [
        { name: "Build and Submit a Sitemap", url: "https://developers.google.com/search/docs/crawling-indexing/sitemaps/build-sitemap" }
      ]
    }
  };

  // Group issues by type and get unique recommendations
  $: issueTypes = [...new Set(issues.map(i => i.type))];
  
  $: recommendationsToShow = issueTypes.map(type => {
    const rec = recommendations[type];
    if (!rec) return null;
    
    // Count affected pages for this issue type
    const affectedPages = issues.filter(i => i.type === type).length;
    
    return {
      issueType: type,
      ...rec,
      affectedPages,
      severity: issues.find(i => i.type === type)?.severity || 'info'
    };
  }).filter(r => r !== null);

  // Sort by impact (High > Medium > Low) and affected pages
  // If enriched data available, use enriched priority for sorting
  $: sortedRecommendations = [...recommendationsToShow].sort((a, b) => {
    // First, check if we have enriched priority for these issue types
    const enrichedA = issues.find(i => i.type === a.issueType && enrichedIssues[`${i.url}|${i.type}`]);
    const enrichedB = issues.find(i => i.type === b.issueType && enrichedIssues[`${i.url}|${i.type}`]);
    
    if (enrichedA && enrichedB) {
      const priorityA = enrichedIssues[`${enrichedA.url}|${enrichedA.type}`]?.enriched_priority || 0;
      const priorityB = enrichedIssues[`${enrichedB.url}|${enrichedB.type}`]?.enriched_priority || 0;
      if (priorityA !== priorityB) {
        return priorityB - priorityA; // Descending
      }
    }
    
    // Fallback to impact and affected pages
    const impactOrder = { 'Critical': 4, 'High': 3, 'Medium': 2, 'Low': 1 };
    const impactDiff = impactOrder[b.impact] - impactOrder[a.impact];
    if (impactDiff !== 0) return impactDiff;
    return b.affectedPages - a.affectedPages;
  });

  const copyToClipboard = async (text) => {
    try {
      await navigator.clipboard.writeText(text);
      // Show temporary feedback
      return true;
    } catch (err) {
      console.error('Failed to copy:', err);
      return false;
    }
  };

  const getImpactColor = (impact) => {
    switch (impact) {
      case 'Critical': return 'badge-error';
      case 'High': return 'badge-warning';
      case 'Medium': return 'badge-info';
      case 'Low': return 'badge-ghost';
      default: return 'badge-ghost';
    }
  };
</script>

<div class="card bg-base-100 shadow">
  <div class="card-body">
    <h2 class="card-title mb-4">Actionable Recommendations</h2>
    
    {#if sortedRecommendations.length === 0}
      <div class="alert alert-success">
        <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <span>No issues found! Your site is looking great.</span>
      </div>
    {:else}
      <div class="space-y-6">
        {#each sortedRecommendations as rec}
          <div class="card bg-base-200 shadow-md">
            <div class="card-body">
              <div class="flex items-start justify-between mb-2">
                <div class="flex-1">
                  <h3 class="card-title text-lg">{rec.title}</h3>
                  <div class="flex gap-2 mt-1">
                    <span class="badge {getImpactColor(rec.impact)}">{rec.impact} Impact</span>
                    <span class="badge badge-ghost">{rec.affectedPages} page{rec.affectedPages !== 1 ? 's' : ''} affected</span>
                  </div>
                </div>
              </div>
              
              <p class="text-base-content/80 mb-4">{rec.description}</p>
              
              <!-- Code Snippet -->
              <div class="mb-4">
                <div class="flex items-center justify-between mb-2">
                  <span class="text-sm font-semibold">Code Example:</span>
                  <button 
                    class="btn btn-xs btn-ghost"
                    on:click={() => copyToClipboard(rec.codeSnippet)}
                    title="Copy to clipboard"
                  >
                    📋 Copy
                  </button>
                </div>
                <pre class="bg-base-300 p-4 rounded-lg overflow-x-auto text-sm"><code>{rec.codeSnippet}</code></pre>
              </div>
              
              <p class="text-sm text-base-content/70 mb-4">{rec.explanation}</p>
              
              <!-- Resources -->
              {#if rec.resources && rec.resources.length > 0}
                <div class="divider"></div>
                <div>
                  <span class="text-sm font-semibold mb-2 block">Learn More:</span>
                  <div class="flex flex-wrap gap-2">
                    {#each rec.resources as resource}
                      <a 
                        href={resource.url} 
                        target="_blank" 
                        rel="noopener noreferrer"
                        class="btn btn-xs btn-outline"
                      >
                        {resource.name} ↗
                      </a>
                    {/each}
                  </div>
                </div>
              {/if}
              
              <!-- Link to view issues -->
              <div class="card-actions justify-end mt-4">
                {#if navigateToTab}
                  <button 
                    class="btn btn-sm btn-primary"
                    on:click={() => navigateToTab('issues', { type: rec.issueType })}
                  >
                    View {rec.affectedPages} Issue{rec.affectedPages !== 1 ? 's' : ''}
                  </button>
                {:else}
                  <a 
                    href="#issues" 
                    class="btn btn-sm btn-primary"
                  >
                    View {rec.affectedPages} Issue{rec.affectedPages !== 1 ? 's' : ''}
                  </a>
                {/if}
              </div>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  pre code {
    white-space: pre;
    font-family: 'Courier New', monospace;
  }
</style>

