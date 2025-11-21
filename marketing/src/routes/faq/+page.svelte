<script lang="ts">
	import { HelpCircle } from 'lucide-svelte';

	let openItems: Record<number, boolean> = $state({});

	function toggleItem(index: number) {
		const currentValue = openItems[index] ?? false;
		// Create a new object to ensure reactivity
		openItems = {
			...openItems,
			[index]: !currentValue
		};
	}
</script>

<svelte:head>
	<title>Frequently Asked Questions - Barracuda SEO</title>
	<meta name="description" content="Get answers to common questions about Barracuda SEO, including how it compares to other tools, pricing, features, and more." />
</svelte:head>

<!-- Hero Section -->
<section class="py-20 bg-gradient-to-b from-[#3c3836] to-[#2d2826]">
	<div class="container mx-auto px-4">
		<div class="text-center max-w-4xl mx-auto">
			<div class="flex items-center justify-center gap-4 mb-6">
				<div class="p-3 bg-[#8ec07c]/10 rounded-lg">
					<HelpCircle class="w-8 h-8 text-[#8ec07c]" />
				</div>
				<h1 class="text-5xl md:text-6xl font-heading font-bold text-white">
					Frequently Asked Questions
				</h1>
			</div>
			<p class="text-xl md:text-2xl text-white/80 mb-10 max-w-3xl mx-auto">
				Everything you need to know about Barracuda SEO, from how it compares to other tools to pricing and features.
			</p>
		</div>
	</div>
</section>

<!-- FAQ Content -->
<section class="py-20 bg-[#3c3836]">
	<div class="container mx-auto px-4">
		<div class="max-w-4xl mx-auto">
			<div class="space-y-4">
				<!-- FAQ 1 -->
				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<button
						class="w-full px-6 py-5 text-left flex items-center justify-between hover:bg-[#3c3836] transition-colors"
						onclick={() => toggleItem(1)}
					>
						<h2 class="text-xl font-heading font-bold text-white pr-4">
							What makes Barracuda different from tools like Screaming Frog?
						</h2>
						<svg
							class="w-5 h-5 text-[#8ec07c] flex-shrink-0 transition-transform {openItems[1] ? 'rotate-180' : ''}"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
						</svg>
					</button>
					{#if openItems[1]}
						<div class="px-6 pb-5 text-white/80">
							<p class="mb-4">
								Screaming Frog is a powerful desktop crawler—but it's single-device, manual, and lacks team collaboration unless you export CSVs and spreadsheets.
							</p>
							<p class="mb-4">Barracuda offers a hybrid model:</p>
							<ul class="list-disc list-inside space-y-2 mb-4 ml-4">
								<li><strong class="text-white">Fast local crawls (CLI — coming soon)</strong></li>
								<li><strong class="text-white">Cloud-based dashboard for teams</strong></li>
								<li><strong class="text-white">AI-powered issue explanations (coming soon)</strong></li>
								<li><strong class="text-white">Built-in integrations</strong> with Google Search Console, GA4, Microsoft Clarity, Slack, and Google Drive</li>
							</ul>
							<p class="italic">Think of Barracuda as <em>Screaming Frog + team dashboard + automation layer</em>.</p>
						</div>
					{/if}
				</div>

				<!-- FAQ 2 -->
				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<button
						class="w-full px-6 py-5 text-left flex items-center justify-between hover:bg-[#3c3836] transition-colors"
						onclick={() => toggleItem(2)}
					>
						<h2 class="text-xl font-heading font-bold text-white pr-4">
							How is Barracuda different from SEMrush and Ahrefs?
						</h2>
						<svg
							class="w-5 h-5 text-[#8ec07c] flex-shrink-0 transition-transform {openItems[2] ? 'rotate-180' : ''}"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
						</svg>
					</button>
					{#if openItems[2]}
						<div class="px-6 pb-5 text-white/80">
							<p class="mb-4">SEMrush and Ahrefs run crawls on their servers, which means:</p>
							<ul class="list-disc list-inside space-y-2 mb-4 ml-4">
								<li>You have limited control over crawl behavior</li>
								<li>Crawl speeds are throttled</li>
								<li>You don't get raw crawl data</li>
								<li>Pricing is tied to credits</li>
								<li>They may miss pages blocked to outside IPs</li>
							</ul>
							<p class="mb-4">Barracuda:</p>
							<ul class="list-disc list-inside space-y-2 mb-4 ml-4">
								<li>Runs <strong class="text-white">local crawls via CLI (coming soon)</strong> for full control</li>
								<li>Allows uploading data to the cloud dashboard</li>
								<li>Stores complete crawl datasets</li>
								<li>Has <strong class="text-white">no credits or crawl caps</strong></li>
								<li>Mirrors what your site serves to <em>real browsers</em></li>
							</ul>
							<p class="font-semibold text-white">Barracuda is <em>your own crawler, your own infrastructure</em>.</p>
						</div>
					{/if}
				</div>

				<!-- FAQ 3 -->
				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<button
						class="w-full px-6 py-5 text-left flex items-center justify-between hover:bg-[#3c3836] transition-colors"
						onclick={() => toggleItem(3)}
					>
						<h2 class="text-xl font-heading font-bold text-white pr-4">
							Why choose Barracuda over Looker Studio dashboards?
						</h2>
						<svg
							class="w-5 h-5 text-[#8ec07c] flex-shrink-0 transition-transform {openItems[3] ? 'rotate-180' : ''}"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
						</svg>
					</button>
					{#if openItems[3]}
						<div class="px-6 pb-5 text-white/80">
							<p class="mb-4">Looker Studio is excellent for reporting—but it can't crawl your site or diagnose SEO issues.</p>
							<p class="mb-4">Barracuda does:</p>
							<ul class="list-disc list-inside space-y-2 mb-4 ml-4">
								<li>Website crawling (CLI coming soon)</li>
								<li>SEO issue detection</li>
								<li>Data collection from GSC, GA4, and Clarity</li>
								<li>Dashboard generation automatically</li>
								<li>Full issue prioritization and insights</li>
							</ul>
							<p class="mb-2">Looker Studio = reporting</p>
							<p class="font-semibold text-white">Barracuda = crawling + diagnosing + reporting + recommendations</p>
						</div>
					{/if}
				</div>

				<!-- FAQ 4 -->
				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<button
						class="w-full px-6 py-5 text-left flex items-center justify-between hover:bg-[#3c3836] transition-colors"
						onclick={() => toggleItem(4)}
					>
						<h2 class="text-xl font-heading font-bold text-white pr-4">
							Does Barracuda replace my current SEO tools?
						</h2>
						<svg
							class="w-5 h-5 text-[#8ec07c] flex-shrink-0 transition-transform {openItems[4] ? 'rotate-180' : ''}"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
						</svg>
					</button>
					{#if openItems[4]}
						<div class="px-6 pb-5 text-white/80">
							<p class="mb-4">Barracuda enhances your SEO stack without replacing everything.</p>
							<p class="mb-2 font-semibold text-white">Replaces:</p>
							<ul class="list-disc list-inside space-y-2 mb-4 ml-4">
								<li>Screaming Frog for most audits (CLI coming soon)</li>
								<li>Manual CSV exports</li>
								<li>Custom Looker Studio builds</li>
								<li>GA4/GSC spreadsheet workflows</li>
							</ul>
							<p class="mb-2 font-semibold text-white">Complements:</p>
							<ul class="list-disc list-inside space-y-2 mb-4 ml-4">
								<li>SEMrush/Ahrefs for keyword research</li>
								<li>GSC for search insights</li>
								<li>GA4 for user behavior</li>
								<li>Clarity for UX signals</li>
							</ul>
							<p class="font-semibold text-white">Barracuda becomes your <em>technical SEO command center</em>.</p>
						</div>
					{/if}
				</div>

				<!-- FAQ 5 -->
				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<button
						class="w-full px-6 py-5 text-left flex items-center justify-between hover:bg-[#3c3836] transition-colors"
						onclick={() => toggleItem(5)}
					>
						<h2 class="text-xl font-heading font-bold text-white pr-4">
							How fast is the Barracuda crawler?
						</h2>
						<svg
							class="w-5 h-5 text-[#8ec07c] flex-shrink-0 transition-transform {openItems[5] ? 'rotate-180' : ''}"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
						</svg>
					</button>
					{#if openItems[5]}
						<div class="px-6 pb-5 text-white/80">
							<p class="mb-4">Once released, the CLI will be extremely fast due to:</p>
							<ul class="list-disc list-inside space-y-2 mb-4 ml-4">
								<li>A Go-based worker-pool engine</li>
								<li>Concurrent crawling</li>
								<li>Intelligent robots.txt + sitemap handling</li>
								<li>Lightweight architecture</li>
							</ul>
							<p class="font-semibold text-white">Expected speed: <strong>100–500 pages/min</strong>, depending on target site performance.</p>
						</div>
					{/if}
				</div>

				<!-- FAQ 6 -->
				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<button
						class="w-full px-6 py-5 text-left flex items-center justify-between hover:bg-[#3c3836] transition-colors"
						onclick={() => toggleItem(6)}
					>
						<h2 class="text-xl font-heading font-bold text-white pr-4">
							What makes Barracuda's issue detection unique?
						</h2>
						<svg
							class="w-5 h-5 text-[#8ec07c] flex-shrink-0 transition-transform {openItems[6] ? 'rotate-180' : ''}"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
						</svg>
					</button>
					{#if openItems[6]}
						<div class="px-6 pb-5 text-white/80">
							<p class="mb-4">Most crawlers simply report raw issues.</p>
							<p class="mb-4">Barracuda:</p>
							<ul class="list-disc list-inside space-y-2 mb-4 ml-4">
								<li>Detects issues</li>
								<li>Prioritizes them by severity + impact</li>
								<li>Groups them by URL structure</li>
								<li>Connects crawl data with GSC/GA4</li>
								<li>(Coming soon) Uses AI to explain root causes and exact fixes</li>
							</ul>
							<p>This intelligence layer goes beyond traditional crawlers.</p>
						</div>
					{/if}
				</div>

				<!-- FAQ 7 -->
				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<button
						class="w-full px-6 py-5 text-left flex items-center justify-between hover:bg-[#3c3836] transition-colors"
						onclick={() => toggleItem(7)}
					>
						<h2 class="text-xl font-heading font-bold text-white pr-4">
							What can I do in the Barracuda Cloud Workspace?
						</h2>
						<svg
							class="w-5 h-5 text-[#8ec07c] flex-shrink-0 transition-transform {openItems[7] ? 'rotate-180' : ''}"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
						</svg>
					</button>
					{#if openItems[7]}
						<div class="px-6 pb-5 text-white/80">
							<p class="mb-4">The dashboard lets you:</p>
							<ul class="list-disc list-inside space-y-2 mb-4 ml-4">
								<li>Save historical crawls</li>
								<li>Compare issues over time</li>
								<li>Invite team members</li>
								<li>View page-level insights</li>
								<li>Send issues to Slack</li>
								<li>Push exports to Google Sheets</li>
								<li>Connect GSC/GA4/Clarity for blended insights</li>
								<li>View recommendations tied to real metrics</li>
							</ul>
							<p>This makes Barracuda collaborative and data-driven.</p>
						</div>
					{/if}
				</div>

				<!-- FAQ 9 -->
				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<button
						class="w-full px-6 py-5 text-left flex items-center justify-between hover:bg-[#3c3836] transition-colors"
						onclick={() => toggleItem(9)}
					>
						<h2 class="text-xl font-heading font-bold text-white pr-4">
							How does billing work?
						</h2>
						<svg
							class="w-5 h-5 text-[#8ec07c] flex-shrink-0 transition-transform {openItems[9] ? 'rotate-180' : ''}"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
						</svg>
					</button>
					{#if openItems[9]}
						<div class="px-6 pb-5 text-white/80">
							<p class="mb-4 font-semibold text-white">No credits. No crawl caps.</p>
							<p class="mb-2 font-semibold text-white">Free tier:</p>
							<ul class="list-disc list-inside space-y-2 mb-4 ml-4">
								<li>100-page crawl limit</li>
								<li>Issues only (no recommendations)</li>
								<li>No integrations</li>
							</ul>
							<p class="mb-2 font-semibold text-white">Pro tier:</p>
							<ul class="list-disc list-inside space-y-2 mb-4 ml-4">
								<li>Unlimited crawls (via CLI—coming soon)</li>
								<li>Full recommendations</li>
								<li>All integrations</li>
								<li>Priority scoring</li>
							</ul>
							<p class="mb-2 font-semibold text-white">Team add-ons:</p>
							<ul class="list-disc list-inside space-y-2 mb-4 ml-4">
								<li>Multi-user support for agencies and in-house teams</li>
							</ul>
							<p>Billing uses secure Stripe subscriptions.</p>
						</div>
					{/if}
				</div>

				<!-- FAQ 10 -->
				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<button
						class="w-full px-6 py-5 text-left flex items-center justify-between hover:bg-[#3c3836] transition-colors"
						onclick={() => toggleItem(10)}
					>
						<h2 class="text-xl font-heading font-bold text-white pr-4">
							Who is Barracuda designed for?
						</h2>
						<svg
							class="w-5 h-5 text-[#8ec07c] flex-shrink-0 transition-transform {openItems[10] ? 'rotate-180' : ''}"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
						</svg>
					</button>
					{#if openItems[10]}
						<div class="px-6 pb-5 text-white/80">
							<p class="mb-4">Barracuda is ideal for:</p>
							<ul class="list-disc list-inside space-y-2 mb-4 ml-4">
								<li>SEO agencies</li>
								<li>Freelancers</li>
								<li>In-house SEO teams</li>
								<li>Developers</li>
								<li>Technical SEOs</li>
								<li>Anyone tired of managing SEO audits manually</li>
							</ul>
							<p>The dashboard is non-technical.</p>
							<p>The CLI (coming soon) is perfect for developers and power users.</p>
						</div>
					{/if}
				</div>

				<!-- FAQ 11 -->
				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<button
						class="w-full px-6 py-5 text-left flex items-center justify-between hover:bg-[#3c3836] transition-colors"
						onclick={() => toggleItem(11)}
					>
						<h2 class="text-xl font-heading font-bold text-white pr-4">
							Does Barracuda support team accounts?
						</h2>
						<svg
							class="w-5 h-5 text-[#8ec07c] flex-shrink-0 transition-transform {openItems[11] ? 'rotate-180' : ''}"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
						</svg>
					</button>
					{#if openItems[11]}
						<div class="px-6 pb-5 text-white/80">
							<p class="mb-4">Yes. You can:</p>
							<ul class="list-disc list-inside space-y-2 mb-4 ml-4">
								<li>Invite collaborators</li>
								<li>Assign roles</li>
								<li>Share projects</li>
								<li>Track crawl history together</li>
								<li>Maintain secure access controls</li>
							</ul>
							<p>Team collaboration is built into the platform.</p>
						</div>
					{/if}
				</div>

				<!-- FAQ 12 -->
				<div class="bg-[#2d2826] rounded-lg border border-white/10 overflow-hidden">
					<button
						class="w-full px-6 py-5 text-left flex items-center justify-between hover:bg-[#3c3836] transition-colors"
						onclick={() => toggleItem(12)}
					>
						<h2 class="text-xl font-heading font-bold text-white pr-4">
							What features are coming soon?
						</h2>
						<svg
							class="w-5 h-5 text-[#8ec07c] flex-shrink-0 transition-transform {openItems[12] ? 'rotate-180' : ''}"
							fill="none"
							stroke="currentColor"
							viewBox="0 0 24 24"
						>
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
						</svg>
					</button>
					{#if openItems[12]}
						<div class="px-6 pb-5 text-white/80">
							<p class="mb-4">Our roadmap includes:</p>
							<ul class="list-disc list-inside space-y-2 mb-4 ml-4">
								<li><strong class="text-white">Full CLI release</strong> for local crawls</li>
								<li><strong class="text-white">AI-powered issue explanations & fix recommendations</strong></li>
								<li><strong class="text-white">Slack + email audit summaries</strong></li>
								<li><strong class="text-white">Google Drive auto-sync</strong></li>
								<li><strong class="text-white">GA4/GSC blended insights</strong></li>
								<li><strong class="text-white">Template-level issue clustering</strong></li>
								<li><strong class="text-white">Automated priority scoring using traffic + crawl data</strong></li>
							</ul>
							<p>Barracuda is rapidly evolving to give you an end-to-end technical SEO workflow.</p>
						</div>
					{/if}
				</div>
			</div>

			<!-- CTA Section -->
			<div class="mt-12 text-center">
				<p class="text-lg text-white/80 mb-6">
					Still have questions? Check out our <a href="/roadmap" class="text-[#8ec07c] hover:text-[#a0d28c] underline">product roadmap</a> or <a href="https://github.com/dillonlara115/barracuda/discussions" class="text-[#8ec07c] hover:text-[#a0d28c] underline" target="_blank" rel="noopener noreferrer">join our community discussions</a>.
				</p>
				<a
					href="https://app.barracudaseo.com"
					class="inline-block bg-[#8ec07c] hover:bg-[#a0d28c] text-[#3c3836] px-8 py-4 rounded-lg font-medium text-lg transition-colors"
					target="_blank"
					rel="noopener noreferrer"
				>
					Start Your Free 100-Page Audit
				</a>
			</div>
		</div>
	</div>
</section>

