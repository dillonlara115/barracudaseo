<script lang="ts">
	import { page } from '$app/stores';
	import Logo from '../Logo.svelte';
	import { ChevronDown } from '@lucide/svelte';
	import { trackSignup } from '$lib/analytics';

	let mobileMenuOpen = $state(false);
	let useCasesDropdownOpen = $state(false);
	let mobileUseCasesOpen = $state(false);

	function toggleMobileMenu() {
		mobileMenuOpen = !mobileMenuOpen;
	}

	function closeMobileMenu() {
		mobileMenuOpen = false;
		mobileUseCasesOpen = false;
	}

	function toggleUseCasesDropdown() {
		useCasesDropdownOpen = !useCasesDropdownOpen;
	}

	function toggleMobileUseCases() {
		mobileUseCasesOpen = !mobileUseCasesOpen;
	}

	function handleHeaderSignupClick() {
		trackSignup({
			source: 'Get Started',
			location: 'header'
		});
	}
</script>

<header class="sticky top-0 z-50 border-b border-white/5 bg-[#1d1d1d]/90 backdrop-blur-md">
	<div class="container mx-auto px-4">
		<nav class="flex h-16 items-center justify-between">
			<div class="flex-1">
				<a href="/" class="transition-opacity hover:opacity-80" onclick={closeMobileMenu}>
					<Logo size="md" />
				</a>
			</div>

			<!-- Desktop Navigation -->
			<div class="hidden flex-none md:flex">
				<ul class="flex items-center gap-6">
					<li>
						<a
							href="/"
							class={`transition-colors hover:text-[#8ec07c] ${$page.url.pathname === '/' ? 'text-[#8ec07c]' : 'text-white'}`}
						>
							Home
						</a>
					</li>
					<li>
						<a
							href="/features"
							class={`transition-colors hover:text-[#8ec07c] ${$page.url.pathname === '/features' ? 'text-[#8ec07c]' : 'text-white'}`}
						>
							Features
						</a>
					</li>
					<li
						class="relative"
						onmouseenter={() => (useCasesDropdownOpen = true)}
						onmouseleave={() => (useCasesDropdownOpen = false)}
					>
						<button
							class={`-mx-1 flex items-center gap-1 px-1 py-2 transition-colors hover:text-[#8ec07c] ${$page.url.pathname.startsWith('/use-cases') ? 'text-[#8ec07c]' : 'text-white'}`}
							onclick={toggleUseCasesDropdown}
						>
							Use Cases
							<ChevronDown
								class={`h-4 w-4 transition-transform ${useCasesDropdownOpen ? 'rotate-180' : ''}`}
							/>
						</button>
						{#if useCasesDropdownOpen}
							<div class="absolute top-full left-0 z-50 w-48 pt-1">
								<div
									class="rounded-xl border border-white/10 bg-[#282828] py-2 shadow-2xl shadow-black/40"
								>
									<a
										href="/use-cases/local-seo"
										class="block px-4 py-2 text-white/60 transition-colors hover:bg-white/5 hover:text-[#8ec07c]"
									>
										Local SEO
									</a>
									<a
										href="/use-cases/programmatic-seo"
										class="block px-4 py-2 text-white/60 transition-colors hover:bg-white/5 hover:text-[#8ec07c]"
									>
										Programmatic SEO
									</a>
									<a
										href="/use-cases/e-commerce"
										class="block px-4 py-2 text-white/60 transition-colors hover:bg-white/5 hover:text-[#8ec07c]"
									>
										E-commerce
									</a>
								</div>
							</div>
						{/if}
					</li>
					<li>
						<a
							href="/pricing"
							class={`transition-colors hover:text-[#8ec07c] ${$page.url.pathname === '/pricing' ? 'text-[#8ec07c]' : 'text-white'}`}
						>
							Pricing
						</a>
					</li>
					<li>
						<a
							href="/about"
							class={`transition-colors hover:text-[#8ec07c] ${$page.url.pathname === '/about' ? 'text-[#8ec07c]' : 'text-white'}`}
						>
							About
						</a>
					</li>
					<li>
						<a
							href="/blog"
							class={`transition-colors hover:text-[#8ec07c] ${$page.url.pathname.startsWith('/blog') ? 'text-[#8ec07c]' : 'text-white'}`}
						>
							Blog
						</a>
					</li>
					<li>
						<a
							href="https://app.barracudaseo.com"
							class="rounded-xl bg-[#8ec07c] px-4 py-2 font-medium text-[#1d1d1d] transition-colors hover:bg-[#a0d28c]"
							target="_blank"
							rel="noopener noreferrer"
							onclick={handleHeaderSignupClick}
						>
							Get Started
						</a>
					</li>
				</ul>
			</div>

			<!-- Mobile Hamburger Button -->
			<button
				class="flex flex-col gap-1.5 p-2 text-white transition-colors hover:text-[#8ec07c] md:hidden"
				onclick={toggleMobileMenu}
				aria-label="Toggle menu"
				aria-expanded={mobileMenuOpen}
			>
				<span
					class="block h-0.5 w-6 bg-current transition-all duration-300 {mobileMenuOpen
						? 'translate-y-2 rotate-45'
						: ''}"
				></span>
				<span
					class="block h-0.5 w-6 bg-current transition-all duration-300 {mobileMenuOpen
						? 'opacity-0'
						: ''}"
				></span>
				<span
					class="block h-0.5 w-6 bg-current transition-all duration-300 {mobileMenuOpen
						? '-translate-y-2 -rotate-45'
						: ''}"
				></span>
			</button>
		</nav>

		<!-- Mobile Navigation Menu -->
		<div
			class="overflow-hidden transition-all duration-300 ease-in-out md:hidden {mobileMenuOpen
				? 'max-h-96 opacity-100'
				: 'max-h-0 opacity-0'}"
		>
			<ul class="flex flex-col gap-4 border-t border-white/10 py-4">
				<li>
					<a
						href="/"
						class={`block py-2 transition-colors hover:text-[#8ec07c] ${$page.url.pathname === '/' ? 'text-[#8ec07c]' : 'text-white'}`}
						onclick={closeMobileMenu}
					>
						Home
					</a>
				</li>
				<li>
					<a
						href="/features"
						class={`block py-2 transition-colors hover:text-[#8ec07c] ${$page.url.pathname === '/features' ? 'text-[#8ec07c]' : 'text-white'}`}
						onclick={closeMobileMenu}
					>
						Features
					</a>
				</li>
				<li>
					<button
						class={`flex w-full items-center justify-between py-2 transition-colors hover:text-[#8ec07c] ${$page.url.pathname.startsWith('/use-cases') ? 'text-[#8ec07c]' : 'text-white'}`}
						onclick={toggleMobileUseCases}
					>
						<span>Use Cases</span>
						<ChevronDown
							class={`h-4 w-4 transition-transform ${mobileUseCasesOpen ? 'rotate-180' : ''}`}
						/>
					</button>
					{#if mobileUseCasesOpen}
						<ul class="mt-2 space-y-2 pl-4">
							<li>
								<a
									href="/use-cases/local-seo"
									class={`block py-2 transition-colors hover:text-[#8ec07c] ${$page.url.pathname === '/use-cases/local-seo' ? 'text-[#8ec07c]' : 'text-white/70'}`}
									onclick={closeMobileMenu}
								>
									Local SEO
								</a>
							</li>
							<li>
								<a
									href="/use-cases/programmatic-seo"
									class={`block py-2 transition-colors hover:text-[#8ec07c] ${$page.url.pathname === '/use-cases/programmatic-seo' ? 'text-[#8ec07c]' : 'text-white/70'}`}
									onclick={closeMobileMenu}
								>
									Programmatic SEO
								</a>
							</li>
							<li>
								<a
									href="/use-cases/e-commerce"
									class={`block py-2 transition-colors hover:text-[#8ec07c] ${$page.url.pathname === '/use-cases/e-commerce' ? 'text-[#8ec07c]' : 'text-white/70'}`}
									onclick={closeMobileMenu}
								>
									E-commerce
								</a>
							</li>
						</ul>
					{/if}
				</li>
				<li>
					<a
						href="/pricing"
						class={`block py-2 transition-colors hover:text-[#8ec07c] ${$page.url.pathname === '/pricing' ? 'text-[#8ec07c]' : 'text-white'}`}
						onclick={closeMobileMenu}
					>
						Pricing
					</a>
				</li>
				<li>
					<a
						href="/about"
						class={`block py-2 transition-colors hover:text-[#8ec07c] ${$page.url.pathname === '/about' ? 'text-[#8ec07c]' : 'text-white'}`}
						onclick={closeMobileMenu}
					>
						About
					</a>
				</li>
				<li>
					<a
						href="/blog"
						class={`block py-2 transition-colors hover:text-[#8ec07c] ${$page.url.pathname.startsWith('/blog') ? 'text-[#8ec07c]' : 'text-white'}`}
						onclick={closeMobileMenu}
					>
						Blog
					</a>
				</li>
				<li class="pt-2">
					<a
						href="https://app.barracudaseo.com"
						class="block rounded-lg bg-[#8ec07c] px-4 py-2 text-center font-medium text-[#3c3836] transition-colors hover:bg-[#a0d28c]"
						target="_blank"
						rel="noopener noreferrer"
						onclick={(e) => {
							handleHeaderSignupClick();
							closeMobileMenu();
						}}
					>
						Get Started
					</a>
				</li>
			</ul>
		</div>
	</div>
</header>
