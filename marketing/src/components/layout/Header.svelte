<script lang="ts">
	import { page } from '$app/stores';
	import Logo from '../Logo.svelte';
	import { ChevronDown } from 'lucide-svelte';

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
</script>

<header class="bg-[#3c3836] border-b border-white/10 sticky top-0 z-50 backdrop-blur-sm bg-[#3c3836]/80">
	<div class="container mx-auto px-4">
		<nav class="flex items-center justify-between h-16">
			<div class="flex-1">
				<a href="/" class="hover:opacity-80 transition-opacity" onclick={closeMobileMenu}>
					<Logo size="md" />
				</a>
			</div>
			
			<!-- Desktop Navigation -->
			<div class="hidden md:flex flex-none">
				<ul class="flex items-center gap-6">
					<li>
						<a href="/" class={`hover:text-[#8ec07c] transition-colors ${$page.url.pathname === '/' ? 'text-[#8ec07c]' : 'text-white'}`}>
							Home
						</a>
					</li>
					<li>
						<a href="/features" class={`hover:text-[#8ec07c] transition-colors ${$page.url.pathname === '/features' ? 'text-[#8ec07c]' : 'text-white'}`}>
							Features
						</a>
					</li>
					<li class="relative" onmouseenter={() => useCasesDropdownOpen = true} onmouseleave={() => useCasesDropdownOpen = false}>
						<button 
							class={`flex items-center gap-1 hover:text-[#8ec07c] transition-colors ${$page.url.pathname.startsWith('/use-cases') ? 'text-[#8ec07c]' : 'text-white'}`}
							onclick={toggleUseCasesDropdown}
						>
							Use Cases
							<ChevronDown class={`w-4 h-4 transition-transform ${useCasesDropdownOpen ? 'rotate-180' : ''}`} />
						</button>
						{#if useCasesDropdownOpen}
							<div class="absolute top-full left-0 mt-2 w-48 bg-[#2d2826] border border-white/10 rounded-lg shadow-lg py-2 z-50">
								<a href="/use-cases/local-seo" class="block px-4 py-2 text-white/70 hover:text-[#8ec07c] hover:bg-[#3c3836] transition-colors">
									Local SEO
								</a>
								<a href="/use-cases/programmatic-seo" class="block px-4 py-2 text-white/70 hover:text-[#8ec07c] hover:bg-[#3c3836] transition-colors">
									Programmatic SEO
								</a>
								<a href="/use-cases/e-commerce" class="block px-4 py-2 text-white/70 hover:text-[#8ec07c] hover:bg-[#3c3836] transition-colors">
									E-commerce
								</a>
							</div>
						{/if}
					</li>
					<li>
						<a href="/pricing" class={`hover:text-[#8ec07c] transition-colors ${$page.url.pathname === '/pricing' ? 'text-[#8ec07c]' : 'text-white'}`}>
							Pricing
						</a>
					</li>
					<li>
						<a href="/about" class={`hover:text-[#8ec07c] transition-colors ${$page.url.pathname === '/about' ? 'text-[#8ec07c]' : 'text-white'}`}>
							About
						</a>
					</li>
					<li>
						<a href="https://app.barracudaseo.com" class="bg-[#8ec07c] hover:bg-[#a0d28c] text-[#3c3836] px-4 py-2 rounded-lg font-medium transition-colors" target="_blank" rel="noopener noreferrer">
							Get Started
						</a>
					</li>
				</ul>
			</div>

			<!-- Mobile Hamburger Button -->
			<button
				class="md:hidden flex flex-col gap-1.5 p-2 text-white hover:text-[#8ec07c] transition-colors"
				onclick={toggleMobileMenu}
				aria-label="Toggle menu"
				aria-expanded={mobileMenuOpen}
			>
				<span class="block w-6 h-0.5 bg-current transition-all duration-300 {mobileMenuOpen ? 'rotate-45 translate-y-2' : ''}"></span>
				<span class="block w-6 h-0.5 bg-current transition-all duration-300 {mobileMenuOpen ? 'opacity-0' : ''}"></span>
				<span class="block w-6 h-0.5 bg-current transition-all duration-300 {mobileMenuOpen ? '-rotate-45 -translate-y-2' : ''}"></span>
			</button>
		</nav>

		<!-- Mobile Navigation Menu -->
		<div class="md:hidden overflow-hidden transition-all duration-300 ease-in-out {mobileMenuOpen ? 'max-h-96 opacity-100' : 'max-h-0 opacity-0'}">
			<ul class="flex flex-col gap-4 py-4 border-t border-white/10">
				<li>
					<a 
						href="/" 
						class={`block py-2 hover:text-[#8ec07c] transition-colors ${$page.url.pathname === '/' ? 'text-[#8ec07c]' : 'text-white'}`}
						onclick={closeMobileMenu}
					>
						Home
					</a>
				</li>
				<li>
					<a 
						href="/features" 
						class={`block py-2 hover:text-[#8ec07c] transition-colors ${$page.url.pathname === '/features' ? 'text-[#8ec07c]' : 'text-white'}`}
						onclick={closeMobileMenu}
					>
						Features
					</a>
				</li>
				<li>
					<button 
						class={`w-full flex items-center justify-between py-2 hover:text-[#8ec07c] transition-colors ${$page.url.pathname.startsWith('/use-cases') ? 'text-[#8ec07c]' : 'text-white'}`}
						onclick={toggleMobileUseCases}
					>
						<span>Use Cases</span>
						<ChevronDown class={`w-4 h-4 transition-transform ${mobileUseCasesOpen ? 'rotate-180' : ''}`} />
					</button>
					{#if mobileUseCasesOpen}
						<ul class="pl-4 mt-2 space-y-2">
							<li>
								<a 
									href="/use-cases/local-seo" 
									class={`block py-2 hover:text-[#8ec07c] transition-colors ${$page.url.pathname === '/use-cases/local-seo' ? 'text-[#8ec07c]' : 'text-white/70'}`}
									onclick={closeMobileMenu}
								>
									Local SEO
								</a>
							</li>
							<li>
								<a 
									href="/use-cases/programmatic-seo" 
									class={`block py-2 hover:text-[#8ec07c] transition-colors ${$page.url.pathname === '/use-cases/programmatic-seo' ? 'text-[#8ec07c]' : 'text-white/70'}`}
									onclick={closeMobileMenu}
								>
									Programmatic SEO
								</a>
							</li>
							<li>
								<a 
									href="/use-cases/e-commerce" 
									class={`block py-2 hover:text-[#8ec07c] transition-colors ${$page.url.pathname === '/use-cases/e-commerce' ? 'text-[#8ec07c]' : 'text-white/70'}`}
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
						class={`block py-2 hover:text-[#8ec07c] transition-colors ${$page.url.pathname === '/pricing' ? 'text-[#8ec07c]' : 'text-white'}`}
						onclick={closeMobileMenu}
					>
						Pricing
					</a>
				</li>
				<li>
					<a 
						href="/about" 
						class={`block py-2 hover:text-[#8ec07c] transition-colors ${$page.url.pathname === '/about' ? 'text-[#8ec07c]' : 'text-white'}`}
						onclick={closeMobileMenu}
					>
						About
					</a>
				</li>
				<li class="pt-2">
					<a 
						href="https://app.barracudaseo.com" 
						class="block bg-[#8ec07c] hover:bg-[#a0d28c] text-[#3c3836] px-4 py-2 rounded-lg font-medium transition-colors text-center" 
						target="_blank" 
						rel="noopener noreferrer"
						onclick={closeMobileMenu}
					>
						Get Started
					</a>
				</li>
			</ul>
		</div>
	</div>
</header>
