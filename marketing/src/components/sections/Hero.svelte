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
	<div class="absolute top-0 left-1/4 w-[600px] h-[600px] bg-[#8ec07c]/5 rounded-full blur-[120px] pointer-events-none"></div>
	<div class="absolute bottom-0 right-1/4 w-[400px] h-[400px] bg-[#d79921]/5 rounded-full blur-[100px] pointer-events-none"></div>

	<div class="container mx-auto px-4 relative z-10">
		<div class="grid grid-cols-1 lg:grid-cols-2 gap-12 lg:gap-16 items-center">
			<!-- Left: Copy -->
			<div class="max-w-xl">
				<div class="inline-flex items-center gap-2 bg-white/5 border border-white/10 text-white/80 px-4 py-1.5 rounded-full text-sm font-medium mb-8">
					<Zap class="w-3.5 h-3.5 text-[#8ec07c]" />
					AI-powered SEO decision engine
				</div>

				<h1 class="text-4xl md:text-5xl lg:text-6xl font-heading font-bold text-white leading-[1.1] tracking-tight mb-6">
					Stop guessing.
					<br />
					<span class="text-[#8ec07c]">Start deciding.</span>
				</h1>

				<p class="text-lg md:text-xl text-white/60 leading-relaxed mb-10 max-w-lg">
					Barracuda turns overwhelming crawl data into clear, defensible priorities — powered by AI that understands your site, your voice, and your data.
				</p>

				<div class="flex flex-col sm:flex-row gap-4">
					<a
						href="https://app.barracudaseo.com"
						class="group inline-flex items-center justify-center gap-2 bg-[#8ec07c] hover:bg-[#a0d28c] text-[#1d1d1d] px-7 py-3.5 rounded-xl font-semibold text-base transition-all shadow-lg shadow-[#8ec07c]/10 hover:shadow-[#8ec07c]/20"
						target="_blank"
						rel="noopener noreferrer"
						onclick={handleSignupClick}
					>
						Start your first audit
						<ArrowRight class="w-4 h-4 transition-transform group-hover:translate-x-0.5" />
					</a>
					<a
						href="#how-it-works"
						class="inline-flex items-center justify-center gap-2 border border-white/15 hover:border-white/30 text-white/80 hover:text-white px-7 py-3.5 rounded-xl font-medium text-base transition-all"
					>
						See how it works
					</a>
				</div>

				<p class="mt-6 text-sm text-white/40">No credit card required. Free tier available.</p>
			</div>

			<!-- Right: Product mockup -->
			<div class="relative">
				<div class="relative bg-[#282828] rounded-2xl border border-white/10 shadow-2xl shadow-black/40 overflow-hidden">
					<!-- Window chrome -->
					<div class="flex items-center gap-2 px-4 py-3 border-b border-white/5 bg-[#222]">
						<div class="w-3 h-3 rounded-full bg-[#cc241d]/60"></div>
						<div class="w-3 h-3 rounded-full bg-[#d79921]/60"></div>
						<div class="w-3 h-3 rounded-full bg-[#8ec07c]/60"></div>
						<span class="ml-3 text-xs text-white/30 font-mono">app.barracudaseo.com</span>
					</div>
					<!-- Dashboard preview content -->
					<div class="p-6 space-y-4">
						<!-- Stats row -->
						<div class="grid grid-cols-3 gap-3">
							<div class="bg-[#1d1d1d] rounded-lg p-3">
								<p class="text-xs text-white/40 mb-1">Pages Crawled</p>
								<p class="text-xl font-bold text-white font-mono">2,847</p>
							</div>
							<div class="bg-[#1d1d1d] rounded-lg p-3">
								<p class="text-xs text-white/40 mb-1">Issues Found</p>
								<p class="text-xl font-bold text-[#d79921] font-mono">142</p>
							</div>
							<div class="bg-[#1d1d1d] rounded-lg p-3">
								<p class="text-xs text-white/40 mb-1">Priority Actions</p>
								<p class="text-xl font-bold text-[#8ec07c] font-mono">12</p>
							</div>
						</div>
						<!-- Issue list preview -->
						<div class="space-y-2">
							<div class="flex items-center gap-3 bg-[#1d1d1d] rounded-lg p-3">
								<div class="w-2 h-2 rounded-full bg-[#cc241d] shrink-0"></div>
								<span class="text-sm text-white/70 truncate">Missing H1 tags on 23 pages</span>
								<span class="ml-auto text-xs text-[#cc241d] font-medium whitespace-nowrap">Critical</span>
							</div>
							<div class="flex items-center gap-3 bg-[#1d1d1d] rounded-lg p-3">
								<div class="w-2 h-2 rounded-full bg-[#d79921] shrink-0"></div>
								<span class="text-sm text-white/70 truncate">Duplicate meta descriptions</span>
								<span class="ml-auto text-xs text-[#d79921] font-medium whitespace-nowrap">Warning</span>
							</div>
							<div class="flex items-center gap-3 bg-[#1d1d1d] rounded-lg p-3">
								<div class="w-2 h-2 rounded-full bg-[#458588] shrink-0"></div>
								<span class="text-sm text-white/70 truncate">Images without alt text</span>
								<span class="ml-auto text-xs text-[#458588] font-medium whitespace-nowrap">Info</span>
							</div>
						</div>
						<!-- AI summary bar -->
						<div class="bg-gradient-to-r from-[#8ec07c]/10 to-transparent rounded-lg p-3 border border-[#8ec07c]/20">
							<div class="flex items-center gap-2 mb-1">
								<Zap class="w-3.5 h-3.5 text-[#8ec07c]" />
								<span class="text-xs font-semibold text-[#8ec07c]">AI Analysis</span>
							</div>
							<p class="text-xs text-white/60">Focus on the 12 critical H1 and title tag issues first — they affect your highest-traffic pages and have the biggest indexing impact.</p>
						</div>
					</div>
				</div>
				<!-- Glow effect behind the card -->
				<div class="absolute inset-0 -z-10 bg-gradient-to-br from-[#8ec07c]/10 via-transparent to-[#d79921]/10 rounded-2xl blur-xl scale-105"></div>
			</div>
		</div>
	</div>
</section>
