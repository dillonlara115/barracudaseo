<script lang="ts">
	import { trackSignup } from '$lib/analytics';
	import { ArrowRight, Zap } from '@lucide/svelte';

	function handleSignupClick() {
		trackSignup({
			source: 'Run a smarter site audit',
			location: 'hero'
		});
	}
</script>

<section class="relative overflow-hidden bg-[#1d1d1d] pt-24 pb-20 md:pt-32 md:pb-28">
	<!-- Subtle gradient orbs -->
	<div
		class="pointer-events-none absolute top-0 left-1/4 h-[600px] w-[600px] rounded-full bg-[#8ec07c]/5 blur-[120px]"
	></div>
	<div
		class="pointer-events-none absolute right-1/4 bottom-0 h-[400px] w-[400px] rounded-full bg-[#d79921]/5 blur-[100px]"
	></div>

	<div class="relative z-10 container mx-auto px-4">
		<div class="grid grid-cols-1 items-center gap-12 lg:grid-cols-2 lg:gap-16">
			<!-- Left: Copy -->
			<div class="max-w-xl">
				<div
					class="mb-8 inline-flex items-center gap-2 rounded-full border border-white/10 bg-white/5 px-4 py-1.5 text-sm font-medium text-white/80"
				>
					<Zap class="h-3.5 w-3.5 text-[#8ec07c]" />
					AI-powered SEO decision engine
				</div>

				<h1
					class="mb-6 font-heading text-4xl leading-[1.1] font-bold tracking-tight text-white md:text-5xl lg:text-6xl"
				>
					Stop guessing.
					<br />
					<span class="text-[#8ec07c]">Start deciding.</span>
				</h1>

				<p class="mb-10 max-w-lg text-lg leading-relaxed text-white/60 md:text-xl">
					Barracuda turns overwhelming crawl data into clear, defensible priorities — powered by AI
					that understands your site, your voice, and your data.
				</p>

				<div class="flex flex-col gap-4 sm:flex-row">
					<a
						href="https://app.barracudaseo.com"
						class="group inline-flex items-center justify-center gap-2 rounded-xl bg-[#8ec07c] px-7 py-3.5 text-base font-semibold text-[#1d1d1d] shadow-lg shadow-[#8ec07c]/10 transition-all hover:bg-[#a0d28c] hover:shadow-[#8ec07c]/20"
						target="_blank"
						rel="noopener noreferrer"
						onclick={handleSignupClick}
					>
						Start your first audit
						<ArrowRight class="h-4 w-4 transition-transform group-hover:translate-x-0.5" />
					</a>
					<a
						href="#how-it-works"
						class="inline-flex items-center justify-center gap-2 rounded-xl border border-white/15 px-7 py-3.5 text-base font-medium text-white/80 transition-all hover:border-white/30 hover:text-white"
					>
						See how it works
					</a>
				</div>

				<p class="mt-6 text-sm text-white/40">No credit card required. Free tier available.</p>
			</div>

			<!-- Right: Product mockup -->
			<div class="relative">
				<div
					class="relative overflow-hidden rounded-2xl border border-white/10 bg-[#282828] shadow-2xl shadow-black/40"
				>
					<!-- Window chrome -->
					<div class="flex items-center gap-2 border-b border-white/5 bg-[#222] px-4 py-3">
						<div class="h-3 w-3 rounded-full bg-[#cc241d]/60"></div>
						<div class="h-3 w-3 rounded-full bg-[#d79921]/60"></div>
						<div class="h-3 w-3 rounded-full bg-[#8ec07c]/60"></div>
						<span class="ml-3 font-mono text-xs text-white/30">app.barracudaseo.com</span>
					</div>
					<!-- Dashboard preview content -->
					<div class="space-y-4 p-6">
						<!-- Stats row -->
						<div class="grid grid-cols-3 gap-3">
							<div class="rounded-lg bg-[#1d1d1d] p-3">
								<p class="mb-1 text-xs text-white/40">Pages Crawled</p>
								<p class="font-mono text-xl font-bold text-white">2,847</p>
							</div>
							<div class="rounded-lg bg-[#1d1d1d] p-3">
								<p class="mb-1 text-xs text-white/40">Issues Found</p>
								<p class="font-mono text-xl font-bold text-[#d79921]">142</p>
							</div>
							<div class="rounded-lg bg-[#1d1d1d] p-3">
								<p class="mb-1 text-xs text-white/40">Priority Actions</p>
								<p class="font-mono text-xl font-bold text-[#8ec07c]">12</p>
							</div>
						</div>
						<!-- Issue list preview -->
						<div class="space-y-2">
							<div class="flex items-center gap-3 rounded-lg bg-[#1d1d1d] p-3">
								<div class="h-2 w-2 shrink-0 rounded-full bg-[#cc241d]"></div>
								<span class="truncate text-sm text-white/70">Missing H1 tags on 23 pages</span>
								<span class="ml-auto text-xs font-medium whitespace-nowrap text-[#cc241d]"
									>Critical</span
								>
							</div>
							<div class="flex items-center gap-3 rounded-lg bg-[#1d1d1d] p-3">
								<div class="h-2 w-2 shrink-0 rounded-full bg-[#d79921]"></div>
								<span class="truncate text-sm text-white/70">Duplicate meta descriptions</span>
								<span class="ml-auto text-xs font-medium whitespace-nowrap text-[#d79921]"
									>Warning</span
								>
							</div>
							<div class="flex items-center gap-3 rounded-lg bg-[#1d1d1d] p-3">
								<div class="h-2 w-2 shrink-0 rounded-full bg-[#458588]"></div>
								<span class="truncate text-sm text-white/70">Images without alt text</span>
								<span class="ml-auto text-xs font-medium whitespace-nowrap text-[#458588]"
									>Info</span
								>
							</div>
						</div>
						<!-- AI summary bar -->
						<div
							class="rounded-lg border border-[#8ec07c]/20 bg-gradient-to-r from-[#8ec07c]/10 to-transparent p-3"
						>
							<div class="mb-1 flex items-center gap-2">
								<Zap class="h-3.5 w-3.5 text-[#8ec07c]" />
								<span class="text-xs font-semibold text-[#8ec07c]">AI Analysis</span>
							</div>
							<p class="text-xs text-white/60">
								Focus on the 12 critical H1 and title tag issues first — they affect your
								highest-traffic pages and have the biggest indexing impact.
							</p>
						</div>
					</div>
				</div>
				<!-- Glow effect behind the card -->
				<div
					class="absolute inset-0 -z-10 scale-105 rounded-2xl bg-gradient-to-br from-[#8ec07c]/10 via-transparent to-[#d79921]/10 blur-xl"
				></div>
			</div>
		</div>
	</div>
</section>
