<script>
  import { link } from 'svelte-spa-router';
</script>

<svelte:head>
  <title>Terms of Service - Barracuda SEO</title>
</svelte:head>

<div class="container mx-auto p-6 max-w-4xl">
  <div class="mb-6">
    <button 
      class="btn btn-ghost btn-sm mb-4"
      on:click={() => window.history.back()}
    >
      ← Back
    </button>
    <h1 class="text-3xl font-bold mb-2">Terms of Service</h1>
    <p class="text-base-content/70">
      Last updated: {new Date().toLocaleDateString()}
    </p>
  </div>

  <div class="prose prose-lg max-w-none">
    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">1. Acceptance of Terms</h2>
      <p class="mb-4">
        By accessing and using Barracuda SEO ("the Service"), you accept and agree to be bound by the terms and provision of this agreement. If you do not agree to these Terms of Service, you may not use the Service.
      </p>
    </section>

    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">2. Description of Service</h2>
      <p class="mb-4">
        Barracuda SEO is a web-based SEO crawler and auditing tool that helps you discover, organize, and fix technical SEO issues on websites. The Service includes:
      </p>
      <ul class="list-disc pl-6 mb-4 space-y-2">
        <li>Automated website crawling and analysis</li>
        <li>SEO issue detection and reporting</li>
        <li>Integration with third-party services (e.g., Google Search Console)</li>
        <li>Data visualization and reporting tools</li>
      </ul>
    </section>

    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">3. User Accounts</h2>
      
      <h3 class="text-xl font-semibold mb-3">3.1 Account Creation</h3>
      <p class="mb-4">
        To use certain features of the Service, you must create an account. You agree to:
      </p>
      <ul class="list-disc pl-6 mb-4 space-y-2">
        <li>Provide accurate, current, and complete information</li>
        <li>Maintain and update your account information</li>
        <li>Maintain the security of your password</li>
        <li>Accept responsibility for all activities under your account</li>
        <li>Notify us immediately of any unauthorized use</li>
      </ul>

      <h3 class="text-xl font-semibold mb-3">3.2 Account Termination</h3>
      <p class="mb-4">
        We reserve the right to suspend or terminate your account at any time for violation of these Terms or for any other reason we deem necessary.
      </p>
    </section>

    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">4. Acceptable Use</h2>
      <p class="mb-4">You agree not to use the Service to:</p>
      <ul class="list-disc pl-6 mb-4 space-y-2">
        <li>Crawl websites without proper authorization or in violation of robots.txt</li>
        <li>Violate any applicable laws or regulations</li>
        <li>Infringe upon intellectual property rights</li>
        <li>Transmit any malicious code, viruses, or harmful data</li>
        <li>Attempt to gain unauthorized access to the Service or related systems</li>
        <li>Interfere with or disrupt the Service or servers</li>
        <li>Use the Service for any illegal or unauthorized purpose</li>
        <li>Impersonate any person or entity</li>
        <li>Collect or store personal data about other users</li>
        <li>Use automated systems to access the Service without permission</li>
      </ul>
    </section>

    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">5. Subscription and Payment</h2>
      
      <h3 class="text-xl font-semibold mb-3">5.1 Subscription Plans</h3>
      <p class="mb-4">
        We offer various subscription plans with different features and usage limits. Subscription fees are billed in advance on a monthly or annual basis.
      </p>

      <h3 class="text-xl font-semibold mb-3">5.2 Payment Terms</h3>
      <ul class="list-disc pl-6 mb-4 space-y-2">
        <li>All fees are non-refundable except as required by law</li>
        <li>You authorize us to charge your payment method for all fees</li>
        <li>If payment fails, we may suspend or terminate your access</li>
        <li>Prices may change with 30 days' notice</li>
      </ul>

      <h3 class="text-xl font-semibold mb-3">5.3 Cancellation</h3>
      <p class="mb-4">
        You may cancel your subscription at any time. Cancellation takes effect at the end of your current billing period. You will continue to have access until the end of the period you've paid for.
      </p>
    </section>

    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">6. Intellectual Property</h2>
      
      <h3 class="text-xl font-semibold mb-3">6.1 Service Ownership</h3>
      <p class="mb-4">
        The Service and its original content, features, and functionality are owned by Barracuda SEO and are protected by international copyright, trademark, patent, trade secret, and other intellectual property laws.
      </p>

      <h3 class="text-xl font-semibold mb-3">6.2 Your Content</h3>
      <p class="mb-4">
        You retain ownership of any content you submit to the Service (including website URLs and crawl data). By using the Service, you grant us a license to use, store, and process your content solely for the purpose of providing the Service.
      </p>
    </section>

    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">7. Third-Party Services</h2>
      <p class="mb-4">
        The Service may integrate with third-party services (such as Google Search Console). Your use of these integrations is subject to the terms and conditions of those third-party services. We are not responsible for the availability, accuracy, or content of third-party services.
      </p>
    </section>

    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">8. Disclaimers</h2>
      <p class="mb-4">
        THE SERVICE IS PROVIDED "AS IS" AND "AS AVAILABLE" WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO:
      </p>
      <ul class="list-disc pl-6 mb-4 space-y-2">
        <li>Warranties of merchantability, fitness for a particular purpose, or non-infringement</li>
        <li>Accuracy, reliability, or completeness of crawl data or reports</li>
        <li>Uninterrupted or error-free operation</li>
        <li>Security or freedom from viruses or harmful components</li>
      </ul>
      <p class="mb-4">
        We do not guarantee that the Service will meet your requirements or that it will be available at all times.
      </p>
    </section>

    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">9. Limitation of Liability</h2>
      <p class="mb-4">
        TO THE MAXIMUM EXTENT PERMITTED BY LAW, IN NO EVENT SHALL BARRACUDA SEO, ITS AFFILIATES, OR THEIR RESPECTIVE OFFICERS, DIRECTORS, EMPLOYEES, OR AGENTS BE LIABLE FOR ANY INDIRECT, INCIDENTAL, SPECIAL, CONSEQUENTIAL, OR PUNITIVE DAMAGES, INCLUDING BUT NOT LIMITED TO:
      </p>
      <ul class="list-disc pl-6 mb-4 space-y-2">
        <li>Loss of profits, data, or business opportunities</li>
        <li>Costs of substitute services</li>
        <li>Damages resulting from use or inability to use the Service</li>
      </ul>
      <p class="mb-4">
        Our total liability for any claims arising from or related to the Service shall not exceed the amount you paid us in the 12 months preceding the claim.
      </p>
    </section>

    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">10. Indemnification</h2>
      <p class="mb-4">
        You agree to indemnify, defend, and hold harmless Barracuda SEO and its affiliates from any claims, damages, losses, liabilities, and expenses (including legal fees) arising from:
      </p>
      <ul class="list-disc pl-6 mb-4 space-y-2">
        <li>Your use of the Service</li>
        <li>Your violation of these Terms</li>
        <li>Your violation of any rights of another party</li>
        <li>Content you submit to the Service</li>
      </ul>
    </section>

    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">11. Data and Privacy</h2>
      <p class="mb-4">
        Your use of the Service is also governed by our Privacy Policy. By using the Service, you consent to the collection and use of your information as described in the Privacy Policy.
      </p>
    </section>

    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">12. Modifications to Terms</h2>
      <p class="mb-4">
        We reserve the right to modify these Terms at any time. We will notify you of material changes by posting the updated Terms on this page and updating the "Last updated" date. Your continued use of the Service after such modifications constitutes acceptance of the updated Terms.
      </p>
    </section>

    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">13. Termination</h2>
      <p class="mb-4">
        We may terminate or suspend your access to the Service immediately, without prior notice, for any reason, including breach of these Terms. Upon termination:
      </p>
      <ul class="list-disc pl-6 mb-4 space-y-2">
        <li>Your right to use the Service will immediately cease</li>
        <li>We may delete your account and data</li>
        <li>All provisions of these Terms that by their nature should survive termination shall survive</li>
      </ul>
    </section>

    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">14. Governing Law</h2>
      <p class="mb-4">
        These Terms shall be governed by and construed in accordance with the laws of [Your Jurisdiction], without regard to its conflict of law provisions. Any disputes arising from these Terms or the Service shall be resolved in the courts of [Your Jurisdiction].
      </p>
    </section>

    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">15. Severability</h2>
      <p class="mb-4">
        If any provision of these Terms is found to be unenforceable or invalid, that provision shall be limited or eliminated to the minimum extent necessary, and the remaining provisions shall remain in full force and effect.
      </p>
    </section>

    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">16. Entire Agreement</h2>
      <p class="mb-4">
        These Terms, together with our Privacy Policy, constitute the entire agreement between you and Barracuda SEO regarding the Service and supersede all prior agreements and understandings.
      </p>
    </section>

    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">17. Contact Information</h2>
      <p class="mb-4">
        If you have any questions about these Terms of Service, please contact us:
      </p>
      <ul class="list-none mb-4 space-y-2">
        <li><strong>Email:</strong> legal@barracudaseo.com</li>
        <li><strong>Website:</strong> <a href="https://barracudaseo.com" class="link link-primary" target="_blank" rel="noopener noreferrer">barracudaseo.com</a></li>
      </ul>
    </section>
  </div>
</div>

