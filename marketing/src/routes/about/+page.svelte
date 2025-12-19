<script lang="ts">
	import MetaTags from '../../components/MetaTags.svelte';
	import { getMetaTags, getBreadcrumbSchema } from '$lib/meta';

	const meta = getMetaTags({
		title: 'About',
		description: 'Learn about Barracuda SEO, a web-based SEO crawler and auditing tool that helps you discover, organize, and fix technical SEO issues.'
	});

	const structuredData = getBreadcrumbSchema([
		{ name: 'Home', url: '/' },
		{ name: 'About', url: '/about' }
	]);
</script>

<MetaTags config={{ ...meta, structuredData }} />

<div class="py-20 bg-[#3c3836]">
	<div class="container mx-auto px-4 max-w-4xl">
		<div class="text-center mb-16">
			<h1 class="text-5xl md:text-6xl font-heading font-bold mb-4 text-white">About Barracuda SEO</h1>
			<p class="text-xl text-white/70">
				Fast, transparent, and empowering
			</p>
		</div>

		<div class="max-w-none">
			<section class="mb-12">
				<h2 class="text-3xl font-heading font-bold mb-4 text-[#8ec07c]">Our Mission</h2>
				<p class="text-lg text-white/80 mb-4">
					Barracuda SEO was created to provide technical Agencies, agencies, and developers with a fast, web-based alternative to traditional website crawlers. Unlike bulky, expensive, or hard-to-scale solutions, Barracuda gives you the flexibility of a web app with optional CLI support for automation.
				</p>
				<p class="text-lg text-white/80">
					Our goal is to make SEO analysis accessible to everyone — from individual developers to large teams — through a transparent, modern platform focused on actionable fixes, not just data.
				</p>
			</section>

			<section class="mb-12">
				<h2 class="text-3xl font-heading font-bold mb-4 text-[#8ec07c]">What Makes Us Different</h2>
				<ul class="list-disc list-inside space-y-2 text-lg text-white/80">
					<li><strong>Web-First Experience:</strong> Log in, start a crawl, and view issues directly from your browser — no installation required</li>
					<li><strong>Team Collaboration:</strong> Built-in role-based permissions and shared project access — perfect for agencies and teams who need seamless collaboration</li>
					<li><strong>Transparent:</strong> Clear limits, no dark patterns, and straightforward pricing with flexible team add-ons</li>
					<li><strong>Actionable:</strong> Built-in severity and impact scoring helps you prioritize what truly matters</li>
					<li><strong>Modern:</strong> Fast, clean UI with cloud-based performance that scales with your needs</li>
					<li><strong>Empowering:</strong> Built by Agencies and engineers for Agencies and engineers</li>
				</ul>
			</section>

			<section class="mb-12">
				<h2 class="text-3xl font-heading font-bold mb-4 text-[#8ec07c]">Who We Serve</h2>
				<ul class="list-disc list-inside space-y-2 text-lg text-white/80">
					<li><strong>Technical Agencies</strong> who need scalable crawl data and clean reporting</li>
					<li><strong>Agencies</strong> managing multiple clients and domains</li>
					<li><strong>Developers</strong> who prefer API/CLI automation over GUIs</li>
					<li><strong>Website owners</strong> who want to understand their SEO issues without enterprise pricing</li>
				</ul>
			</section>

			<section>
				<h2 class="text-3xl font-heading font-bold mb-4 text-[#8ec07c]">Get Started</h2>
				<p class="text-lg text-white/80 mb-6">
					Ready to start crawling? Sign up for free and begin analyzing your website's SEO in minutes — no credit card required.
				</p>
				<a href="https://app.barracudaseo.com" class="inline-block bg-[#8ec07c] hover:bg-[#a0d28c] text-[#3c3836] px-8 py-4 rounded-lg font-medium text-lg transition-colors" target="_blank" rel="noopener noreferrer">
					Start your free 100-page audit
				</a>
			</section>
		</div>
	</div>
</div>

