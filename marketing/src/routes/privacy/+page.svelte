<script lang="ts">
	import MetaTags from '../../components/MetaTags.svelte';
	import { getMetaTags, getBreadcrumbSchema } from '$lib/meta';

	const meta = getMetaTags({
		title: 'Privacy Policy & Data Protection',
		description: 'Barracuda SEO Privacy Policy. Learn how we collect, protect, and use your data when you use our SEO auditing platform. We prioritize your data security and privacy.'
	});

	const structuredData = getBreadcrumbSchema([
		{ name: 'Home', url: '/' },
		{ name: 'Privacy Policy', url: '/privacy' }
	]);
</script>

<MetaTags config={{ ...meta, structuredData }} />

<!-- Hero Section -->
<section class="py-20 bg-gradient-to-b from-[#3c3836] to-[#2d2826]">
	<div class="container mx-auto px-4">
		<div class="text-center max-w-4xl mx-auto">
			<h1 class="text-5xl md:text-6xl font-heading font-bold mb-4 text-white">
				Privacy Policy
			</h1>
			<p class="text-xl text-white/70">
				Last updated: {new Date().toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })}
			</p>
		</div>
	</div>
</section>

<!-- Privacy Policy Content -->
<section class="py-20 bg-[#3c3836]">
	<div class="container mx-auto px-4">
		<div class="max-w-4xl mx-auto prose prose-invert prose-lg">
			<div class="text-white/80 space-y-8">
				<section>
					<h2 class="text-3xl font-heading font-bold mb-4 text-[#8ec07c]">1. Introduction</h2>
					<p class="mb-4">
						Barracuda SEO ("we," "our," or "us") is committed to protecting your privacy. This Privacy Policy explains how we collect, use, disclose, and safeguard your information when you use our web-based SEO crawler and auditing tool (the "Service").
					</p>
					<p>
						By using our Service, you agree to the collection and use of information in accordance with this policy.
					</p>
				</section>

				<section>
					<h2 class="text-3xl font-heading font-bold mb-4 text-[#8ec07c]">2. Information We Collect</h2>
					
					<h3 class="text-2xl font-heading font-bold mb-3 text-white mt-6">2.1 Information You Provide</h3>
					<ul class="list-disc list-inside space-y-2 mb-4 ml-4">
						<li><strong class="text-white">Account Information:</strong> Name, email address, and password when you create an account</li>
						<li><strong class="text-white">Project Data:</strong> Website URLs, crawl configurations, and project settings</li>
						<li><strong class="text-white">Payment Information:</strong> Billing details processed securely through Stripe (we do not store credit card information)</li>
						<li><strong class="text-white">Communication:</strong> Messages you send to us via email or support channels</li>
					</ul>

					<h3 class="text-2xl font-heading font-bold mb-3 text-white mt-6">2.2 Information Collected Automatically</h3>
					<ul class="list-disc list-inside space-y-2 mb-4 ml-4">
						<li><strong class="text-white">Usage Data:</strong> How you interact with the Service, including pages visited, features used, and time spent</li>
						<li><strong class="text-white">Technical Data:</strong> IP address, browser type, device information, and operating system</li>
						<li><strong class="text-white">Crawl Data:</strong> Website crawl results, SEO issues detected, and performance metrics</li>
						<li><strong class="text-white">Integration Data:</strong> Data from connected services (Google Search Console, Google Analytics, Microsoft Clarity) when you authorize integrations</li>
					</ul>
				</section>

				<section>
					<h2 class="text-3xl font-heading font-bold mb-4 text-[#8ec07c]">3. How We Use Your Information</h2>
					<p class="mb-4">We use the information we collect to:</p>
					<ul class="list-disc list-inside space-y-2 mb-4 ml-4">
						<li>Provide, maintain, and improve the Service</li>
						<li>Process transactions and manage your account</li>
						<li>Send you technical notices, updates, and support messages</li>
						<li>Respond to your comments, questions, and requests</li>
						<li>Monitor and analyze usage patterns to improve user experience</li>
						<li>Detect, prevent, and address technical issues and security threats</li>
						<li>Comply with legal obligations and enforce our terms</li>
					</ul>
				</section>

				<section>
					<h2 class="text-3xl font-heading font-bold mb-4 text-[#8ec07c]">4. Data Storage and Security</h2>
					<p class="mb-4">
						We take data security seriously and implement industry-standard measures to protect your information:
					</p>
					<ul class="list-disc list-inside space-y-2 mb-4 ml-4">
						<li><strong class="text-white">Encryption:</strong> Data is encrypted in transit using TLS/SSL and at rest</li>
						<li><strong class="text-white">Access Controls:</strong> Row-Level Security (RLS) policies ensure users can only access their own project data</li>
						<li><strong class="text-white">Secure Infrastructure:</strong> Data is stored in Supabase with enterprise-grade security</li>
						<li><strong class="text-white">API Security:</strong> API keys and tokens are stored securely using Cloud Run secrets</li>
						<li><strong class="text-white">Regular Audits:</strong> We regularly review and update our security practices</li>
					</ul>
					<p class="mb-4">
						However, no method of transmission over the Internet or electronic storage is 100% secure. While we strive to use commercially acceptable means to protect your data, we cannot guarantee absolute security.
					</p>
				</section>

				<section>
					<h2 class="text-3xl font-heading font-bold mb-4 text-[#8ec07c]">5. Data Sharing and Disclosure</h2>
					<p class="mb-4">We do not sell, trade, or rent your personal information to third parties. We may share your information only in the following circumstances:</p>
					<ul class="list-disc list-inside space-y-2 mb-4 ml-4">
						<li><strong class="text-white">Service Providers:</strong> With trusted third-party services that help us operate the Service (e.g., Supabase for data storage, Stripe for payments, Vercel for hosting)</li>
						<li><strong class="text-white">Team Collaboration:</strong> With team members you explicitly invite to shared projects</li>
						<li><strong class="text-white">Legal Requirements:</strong> When required by law, court order, or government regulation</li>
						<li><strong class="text-white">Business Transfers:</strong> In connection with a merger, acquisition, or sale of assets (with notice to users)</li>
						<li><strong class="text-white">With Your Consent:</strong> When you explicitly authorize us to share information</li>
					</ul>
				</section>

				<section>
					<h2 class="text-3xl font-heading font-bold mb-4 text-[#8ec07c]">6. Third-Party Integrations</h2>
					<p class="mb-4">
						Our Service integrates with third-party services (Google Search Console, Google Analytics, Microsoft Clarity, Slack, Google Drive). When you connect these services:
					</p>
					<ul class="list-disc list-inside space-y-2 mb-4 ml-4">
						<li>We access only the data you explicitly authorize</li>
						<li>We use OAuth tokens that you can revoke at any time</li>
						<li>Your data remains subject to the privacy policies of those third-party services</li>
						<li>We do not share your Barracuda data with these third parties beyond what's necessary for the integration</li>
					</ul>
				</section>

				<section>
					<h2 class="text-3xl font-heading font-bold mb-4 text-[#8ec07c]">7. Your Rights and Choices</h2>
					<p class="mb-4">You have the right to:</p>
					<ul class="list-disc list-inside space-y-2 mb-4 ml-4">
						<li><strong class="text-white">Access:</strong> Request a copy of your personal data</li>
						<li><strong class="text-white">Correction:</strong> Update or correct inaccurate information</li>
						<li><strong class="text-white">Deletion:</strong> Request deletion of your account and associated data</li>
						<li><strong class="text-white">Data Portability:</strong> Export your project data</li>
						<li><strong class="text-white">Opt-Out:</strong> Unsubscribe from marketing communications (account-related emails will still be sent)</li>
						<li><strong class="text-white">Revoke Access:</strong> Disconnect third-party integrations at any time</li>
					</ul>
					<p>
						To exercise these rights, please contact us at <a href="mailto:privacy@barracudaseo.com" class="text-[#8ec07c] hover:text-[#a0d28c] underline">privacy@barracudaseo.com</a>.
					</p>
				</section>

				<section>
					<h2 class="text-3xl font-heading font-bold mb-4 text-[#8ec07c]">8. Data Retention</h2>
					<p class="mb-4">
						We retain your information for as long as your account is active or as needed to provide the Service. When you delete your account:
					</p>
					<ul class="list-disc list-inside space-y-2 mb-4 ml-4">
						<li>We will delete your personal information and project data</li>
						<li>Some information may remain in backups for a limited time for disaster recovery purposes</li>
						<li>We may retain anonymized, aggregated data that cannot identify you</li>
					</ul>
				</section>

				<section>
					<h2 class="text-3xl font-heading font-bold mb-4 text-[#8ec07c]">9. Children's Privacy</h2>
					<p>
						Our Service is not intended for individuals under the age of 18. We do not knowingly collect personal information from children. If you become aware that a child has provided us with personal information, please contact us immediately.
					</p>
				</section>

				<section>
					<h2 class="text-3xl font-heading font-bold mb-4 text-[#8ec07c]">10. International Data Transfers</h2>
					<p>
						Your information may be transferred to and processed in countries other than your country of residence. These countries may have data protection laws that differ from those in your country. By using the Service, you consent to the transfer of your information to these countries.
					</p>
				</section>

				<section>
					<h2 class="text-3xl font-heading font-bold mb-4 text-[#8ec07c]">11. Changes to This Privacy Policy</h2>
					<p class="mb-4">
						We may update this Privacy Policy from time to time. We will notify you of any material changes by:
					</p>
					<ul class="list-disc list-inside space-y-2 mb-4 ml-4">
						<li>Posting the new Privacy Policy on this page</li>
						<li>Updating the "Last updated" date</li>
						<li>Sending you an email notification (for significant changes)</li>
					</ul>
					<p>
						Your continued use of the Service after changes become effective constitutes acceptance of the updated Privacy Policy.
					</p>
				</section>

				<section>
					<h2 class="text-3xl font-heading font-bold mb-4 text-[#8ec07c]">12. Contact Us</h2>
					<p class="mb-4">
						If you have questions about this Privacy Policy or our data practices, please contact us:
					</p>
					<ul class="list-none space-y-2 mb-4">
						<li><strong class="text-white">Email:</strong> <a href="mailto:privacy@barracudaseo.com" class="text-[#8ec07c] hover:text-[#a0d28c] underline">privacy@barracudaseo.com</a></li>
						<li><strong class="text-white">Website:</strong> <a href="https://barracudaseo.com" class="text-[#8ec07c] hover:text-[#a0d28c] underline">barracudaseo.com</a></li>
					</ul>
				</section>
			</div>
		</div>
	</div>
</section>



