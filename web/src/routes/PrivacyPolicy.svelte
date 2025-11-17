<script>
  import { link } from 'svelte-spa-router';
</script>

<svelte:head>
  <title>Privacy Policy - Barracuda SEO</title>
</svelte:head>

<div class="container mx-auto p-6 max-w-4xl">
  <div class="mb-6">
    <button 
      class="btn btn-ghost btn-sm mb-4"
      on:click={() => window.history.back()}
    >
      ← Back
    </button>
    <h1 class="text-3xl font-bold mb-2">Privacy Policy</h1>
    <p class="text-base-content/70">
      Last updated: {new Date().toLocaleDateString()}
    </p>
  </div>

  <div class="prose prose-lg max-w-none">
    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">1. Introduction</h2>
      <p class="mb-4">
        Barracuda SEO ("we," "our," or "us") is committed to protecting your privacy. This Privacy Policy explains how we collect, use, disclose, and safeguard your information when you use our web-based SEO crawler and auditing service (the "Service").
      </p>
      <p class="mb-4">
        By using our Service, you agree to the collection and use of information in accordance with this policy.
      </p>
    </section>

    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">2. Information We Collect</h2>
      
      <h3 class="text-xl font-semibold mb-3">2.1 Information You Provide</h3>
      <ul class="list-disc pl-6 mb-4 space-y-2">
        <li><strong>Account Information:</strong> Email address, name, and password when you create an account</li>
        <li><strong>Project Information:</strong> Project names, domains, and URLs you submit for crawling</li>
        <li><strong>Payment Information:</strong> Billing details processed through our third-party payment processor (Stripe)</li>
        <li><strong>Integration Data:</strong> OAuth tokens and credentials for connected services (e.g., Google Search Console)</li>
      </ul>

      <h3 class="text-xl font-semibold mb-3">2.2 Automatically Collected Information</h3>
      <ul class="list-disc pl-6 mb-4 space-y-2">
        <li><strong>Usage Data:</strong> How you interact with our Service, including pages visited and features used</li>
        <li><strong>Technical Data:</strong> IP address, browser type, device information, and operating system</li>
        <li><strong>Cookies and Tracking:</strong> We use cookies and similar tracking technologies to track activity on our Service</li>
      </ul>

      <h3 class="text-xl font-semibold mb-3">2.3 Crawled Website Data</h3>
      <p class="mb-4">
        When you use our Service to crawl websites, we collect and store:
      </p>
      <ul class="list-disc pl-6 mb-4 space-y-2">
        <li>HTML content and structure of crawled pages</li>
        <li>SEO-related data (titles, meta descriptions, headings, links)</li>
        <li>Technical issues and recommendations</li>
        <li>Performance metrics and analytics data</li>
      </ul>
    </section>

    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">3. How We Use Your Information</h2>
      <p class="mb-4">We use the collected information for various purposes:</p>
      <ul class="list-disc pl-6 mb-4 space-y-2">
        <li>To provide, maintain, and improve our Service</li>
        <li>To process your transactions and manage your account</li>
        <li>To perform website crawls and generate SEO reports</li>
        <li>To send you notifications, updates, and support communications</li>
        <li>To detect, prevent, and address technical issues</li>
        <li>To comply with legal obligations</li>
        <li>To analyze usage patterns and improve user experience</li>
      </ul>
    </section>

    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">4. Data Storage and Security</h2>
      <p class="mb-4">
        We use Supabase (a PostgreSQL database hosted on secure cloud infrastructure) to store your data. We implement appropriate technical and organizational measures to protect your personal information, including:
      </p>
      <ul class="list-disc pl-6 mb-4 space-y-2">
        <li>Encryption of data in transit and at rest</li>
        <li>Access controls and authentication mechanisms</li>
        <li>Regular security assessments and updates</li>
        <li>Secure API endpoints with authentication</li>
      </ul>
      <p class="mb-4">
        However, no method of transmission over the Internet or electronic storage is 100% secure. While we strive to use commercially acceptable means to protect your data, we cannot guarantee absolute security.
      </p>
    </section>

    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">5. Data Sharing and Disclosure</h2>
      <p class="mb-4">We do not sell your personal information. We may share your information only in the following circumstances:</p>
      
      <h3 class="text-xl font-semibold mb-3">5.1 Service Providers</h3>
      <p class="mb-4">
        We may share information with third-party service providers who perform services on our behalf, including:
      </p>
      <ul class="list-disc pl-6 mb-4 space-y-2">
        <li><strong>Supabase:</strong> Database and authentication services</li>
        <li><strong>Stripe:</strong> Payment processing</li>
        <li><strong>Google:</strong> OAuth authentication for Google Search Console integration</li>
        <li><strong>Hosting Providers:</strong> Cloud infrastructure services</li>
      </ul>

      <h3 class="text-xl font-semibold mb-3">5.2 Legal Requirements</h3>
      <p class="mb-4">
        We may disclose your information if required by law or in response to valid requests by public authorities.
      </p>

      <h3 class="text-xl font-semibold mb-3">5.3 Business Transfers</h3>
      <p class="mb-4">
        In the event of a merger, acquisition, or sale of assets, your information may be transferred as part of that transaction.
      </p>
    </section>

    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">6. Your Rights and Choices</h2>
      <p class="mb-4">Depending on your location, you may have certain rights regarding your personal information:</p>
      <ul class="list-disc pl-6 mb-4 space-y-2">
        <li><strong>Access:</strong> Request access to your personal data</li>
        <li><strong>Correction:</strong> Request correction of inaccurate data</li>
        <li><strong>Deletion:</strong> Request deletion of your data</li>
        <li><strong>Portability:</strong> Request transfer of your data</li>
        <li><strong>Objection:</strong> Object to processing of your data</li>
        <li><strong>Withdrawal:</strong> Withdraw consent where processing is based on consent</li>
      </ul>
      <p class="mb-4">
        To exercise these rights, please contact us at the email address provided below. We will respond to your request within a reasonable timeframe.
      </p>
    </section>

    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">7. Cookies and Tracking Technologies</h2>
      <p class="mb-4">
        We use cookies and similar tracking technologies to track activity on our Service and hold certain information. You can instruct your browser to refuse all cookies or to indicate when a cookie is being sent. However, if you do not accept cookies, you may not be able to use some portions of our Service.
      </p>
    </section>

    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">8. Data Retention</h2>
      <p class="mb-4">
        We retain your personal information for as long as necessary to provide our Service and fulfill the purposes outlined in this Privacy Policy, unless a longer retention period is required or permitted by law. When you delete your account, we will delete or anonymize your personal information, except where we are required to retain it for legal purposes.
      </p>
    </section>

    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">9. Children's Privacy</h2>
      <p class="mb-4">
        Our Service is not intended for individuals under the age of 18. We do not knowingly collect personal information from children. If you become aware that a child has provided us with personal information, please contact us, and we will take steps to delete such information.
      </p>
    </section>

    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">10. International Data Transfers</h2>
      <p class="mb-4">
        Your information may be transferred to and maintained on computers located outside of your state, province, country, or other governmental jurisdiction where data protection laws may differ. By using our Service, you consent to the transfer of your information to our facilities and those third parties with whom we share it as described in this Privacy Policy.
      </p>
    </section>

    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">11. Changes to This Privacy Policy</h2>
      <p class="mb-4">
        We may update our Privacy Policy from time to time. We will notify you of any changes by posting the new Privacy Policy on this page and updating the "Last updated" date. You are advised to review this Privacy Policy periodically for any changes.
      </p>
    </section>

    <section class="mb-8">
      <h2 class="text-2xl font-bold mb-4">12. Contact Us</h2>
      <p class="mb-4">
        If you have any questions about this Privacy Policy, please contact us:
      </p>
      <ul class="list-none mb-4 space-y-2">
        <li><strong>Email:</strong> privacy@barracudaseo.com</li>
        <li><strong>Website:</strong> <a href="https://barracudaseo.com" class="link link-primary" target="_blank" rel="noopener noreferrer">barracudaseo.com</a></li>
      </ul>
    </section>
  </div>
</div>

