<script lang="ts">
	import { page } from '$app/stores';
	import Logo from '../Logo.svelte';

	let mobileMenuOpen = $state(false);

	function toggleMobileMenu() {
		mobileMenuOpen = !mobileMenuOpen;
	}

	function closeMobileMenu() {
		mobileMenuOpen = false;
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
