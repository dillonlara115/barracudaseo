<script>
  import { onMount } from 'svelte';
  import { user } from '../lib/auth.js';
  import { updateEmail } from '../lib/auth.js';
  import { supabase } from '../lib/supabase.js';
  import { CreditCard, Check, X, Loader, ArrowLeft, Mail } from 'lucide-svelte';
  import { userProfile } from '../lib/subscription.js';
  import { push, link } from 'svelte-spa-router';
  import Logo from './Logo.svelte';
  import Auth from './Auth.svelte';
  import TeamManagement from './TeamManagement.svelte';
  import { fetchBillingSummary, createBillingCheckout, createBillingPortal, redeemPromoCode } from '../lib/data.js';

  let loading = true;
  let profile = null;
  let subscription = null;
  let teamInfo = null; // Team membership info
  let error = null;
  let creatingCheckout = false;
  let creatingPortal = false;
  
  const API_URL = import.meta.env.VITE_CLOUD_RUN_API_URL || 'http://localhost:8080';
  const STRIPE_PRICE_ID_PRO = import.meta.env.VITE_STRIPE_PRICE_ID_PRO || '';
  const STRIPE_PRICE_ID_PRO_ANNUAL = import.meta.env.VITE_STRIPE_PRICE_ID_PRO_ANNUAL || '';
  const STRIPE_PRICE_ID_TEAM_SEAT = import.meta.env.VITE_STRIPE_PRICE_ID_TEAM_SEAT || '';
  
  let selectedBillingPeriod = 'monthly'; // 'monthly' or 'annual'
  let teamSeatsQuantity = 0; // Number of additional team seats to add
  let hasLoaded = false; // Track if we've attempted to load
  
  let redeemCode = '';
  let redeemTeamSize = 1;
  let redeeming = false;
  let redeemError = null;

  // Email change state
  let newEmail = '';
  let passwordForEmailChange = '';
  let updatingEmail = false;
  let emailChangeError = null;
  let emailChangeSuccess = null;
  let showEmailChangeForm = false;

  // Load data when component mounts and user is available
  onMount(() => {
    // Check for success/cancel parameters from Stripe redirect
    const urlParams = new URLSearchParams(window.location.search);
    const success = urlParams.get('success');
    const canceled = urlParams.get('canceled');
    
    // Clean up URL parameters
    if (success || canceled) {
      window.history.replaceState({}, '', window.location.pathname);
    }
    
    // Subscribe to user store and load when user becomes available
    let unsubscribe;
    unsubscribe = user.subscribe(async (currentUser) => {
      if (currentUser) {
        // Check if we have a valid session before making API calls
        const { data: { session } } = await supabase.auth.getSession();
        if (!session) {
          // User exists but no session (email not confirmed?)
          loading = false;
          error = 'Please confirm your email address to access billing information.';
          return;
        }
        
        // Only load if we haven't loaded yet or if explicitly refreshing
        if (!hasLoaded) {
          hasLoaded = true;
          await loadBillingData();
        }
        
        // Show success message if returning from checkout
        if (success) {
          // Small delay to ensure data is loaded
          setTimeout(() => {
            // You could add a toast notification here if you have one
            console.log('Payment successful! Subscription updated.');
          }, 500);
        }
      } else if (!currentUser && !hasLoaded) {
        // No user yet, but don't keep loading state forever
        loading = false;
      }
    });
    
    return unsubscribe;
  });

  // Removed local getValidAccessToken - now using centralized auth wrapper from data.js

  async function loadBillingData() {
    if (!$user) {
      loading = false;
      return;
    }
    
    loading = true;
    error = null;
    
    try {
      const { data, error: fetchError } = await fetchBillingSummary();
      
      if (fetchError) {
        throw fetchError;
      }

      // Backend should always return a profile (it creates one if missing)
      // But handle the case where it might be null/undefined
      if (data?.profile) {
        profile = data.profile;
        // Update the subscription store
        userProfile.set(data.profile);
      } else {
        // Fallback: create a default profile object
        profile = {
          id: $user.id,
          subscription_tier: 'free',
          subscription_status: 'active',
          team_size: 1
        };
        userProfile.set(profile);
      }
      
      subscription = data?.subscription || null;
      teamInfo = data?.team_info || null; // Team membership info
    } catch (err) {
      error = err.message || 'Failed to load subscription data';
      console.error('Failed to load subscription data:', err);
      
      // Set a default profile on error so the UI can still render
      if ($user && !profile) {
        profile = {
          id: $user.id,
          subscription_tier: 'free',
          subscription_status: 'active',
          team_size: 1
        };
      }
    } finally {
      loading = false;
    }
  }

  async function createCheckoutSession(priceId) {
    if (!$user) return;
    
    creatingCheckout = true;
    error = null;
    
    try {
      const { data, error: fetchError } = await createBillingCheckout(priceId);
      
      if (fetchError) {
        throw fetchError;
      }
      
      if (data?.url) {
        window.location.href = data.url;
      } else {
        throw new Error('No checkout URL returned');
      }
    } catch (err) {
      error = err.message || 'Failed to create checkout session';
      console.error('Failed to create checkout session:', err);
    } finally {
      creatingCheckout = false;
    }
  }

  async function openCustomerPortal() {
    if (!$user) return;
    
    creatingPortal = true;
    error = null;
    
    try {
      const { data, error: fetchError } = await createBillingPortal();
      
      if (fetchError) {
        throw fetchError;
      }
      
      if (data?.url) {
        window.location.href = data.url;
      } else {
        throw new Error('No portal URL returned');
      }
    } catch (err) {
      error = err.message || 'Failed to create portal session';
      console.error('Failed to create portal session:', err);
    } finally {
      creatingPortal = false;
    }
  }

  // Alias for consistency
  const openBillingPortal = openCustomerPortal;

  async function handleRedeemCode() {
    if (!$user || !redeemCode.trim()) return;
    
    redeeming = true;
    redeemError = null;
    
    try {
      const { data, error: fetchError } = await redeemPromoCode(redeemCode.trim(), redeemTeamSize);
      
      if (fetchError) {
        throw fetchError;
      }
      
      // Success - reload billing data and clear form
      redeemCode = '';
      redeemTeamSize = 1;
      await loadBillingData();
      
      // Show success message (you could add a toast here)
      console.log('Promo code redeemed successfully');
    } catch (err) {
      redeemError = err.message || 'Failed to redeem code';
      console.error('Failed to redeem code:', err);
    } finally {
      redeeming = false;
    }
  }


  function getPlanFeatures(tier) {
    switch (tier) {
      case 'pro':
        return {
          pages: '10,000',
          users: profile?.team_size || 1,
          integrations: true,
          recommendations: true,
        };
      case 'team':
        // Team plans are custom - use same features as Pro but with custom limits
        return {
          pages: 'Custom',
          users: profile?.team_size || 5,
          integrations: true,
          recommendations: true,
        };
      default:
        return {
          pages: '100',
          users: 1,
          integrations: false,
          recommendations: false,
        };
    }
  }

  function formatDate(dateString) {
    if (!dateString) return 'N/A';
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  }

  $: planFeatures = getPlanFeatures(profile?.subscription_tier || 'free');
  $: isProOrTeam = profile?.subscription_tier === 'pro' || profile?.subscription_tier === 'team';

  async function redeemBetaCode() {
    if (!redeemCode.trim()) return;
    
    redeeming = true;
    redeemError = null;
    
    try {
      const { data, error: fetchError } = await redeemPromoCode(redeemCode.trim(), redeemTeamSize);
      
      if (fetchError) {
        throw fetchError;
      }
      
      // Success - reload billing data and clear form
      redeemCode = '';
      redeemTeamSize = 1;
      await loadBillingData();
      
      // Show success message (you could add a toast here)
      console.log('Promo code redeemed successfully');
    } catch (err) {
      redeemError = err.message || 'Failed to redeem code';
      console.error('Failed to redeem code:', err);
    } finally {
      redeeming = false;
    }
  }


  async function handleEmailChange() {
    if (!newEmail.trim()) {
      emailChangeError = 'Please enter a new email address';
      return;
    }

    // Basic email validation
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(newEmail)) {
      emailChangeError = 'Please enter a valid email address';
      return;
    }

    // Check if email is different from current
    if ($user && newEmail.toLowerCase() === $user.email?.toLowerCase()) {
      emailChangeError = 'New email must be different from your current email';
      return;
    }

    if (!passwordForEmailChange) {
      emailChangeError = 'Please enter your password to confirm this change';
      return;
    }

    updatingEmail = true;
    emailChangeError = null;
    emailChangeSuccess = null;

    try {
      const { data, error: updateError } = await updateEmail(newEmail, passwordForEmailChange);
      
      if (updateError) {
        throw updateError;
      }

      emailChangeSuccess = 'Email change request sent! Please check both your old and new email addresses to confirm the change.';
      newEmail = '';
      passwordForEmailChange = '';
      showEmailChangeForm = false;
      
      // Refresh user data after a short delay
      setTimeout(async () => {
        const { data: { user: updatedUser } } = await supabase.auth.getUser();
        if (updatedUser) {
          user.set(updatedUser);
        }
      }, 1000);
    } catch (err) {
      emailChangeError = err.message || 'Failed to update email address. Please try again.';
    } finally {
      updatingEmail = false;
    }
  }
</script>

<!-- Header Navigation -->
<div class="navbar bg-base-100 shadow-lg border-b border-base-300 gap-2">
  <div class="flex-1">
    <a href="#/" use:link class="btn btn-ghost">
      <Logo size="md" />
    </a>
  </div>
  <div class="flex gap-2">
    <Auth />
  </div>
</div>

<div class="container mx-auto p-6 max-w-4xl">
  <div class="mb-6">
    <button 
      class="btn btn-ghost btn-sm mb-4"
      on:click={() => push('#/')}
    >
      <ArrowLeft class="w-4 h-4 mr-2" />
      Back to Projects
    </button>
    <h1 class="text-3xl font-bold mb-2">Billing & Subscription</h1>
    <p class="text-base-content/70">
      Manage your subscription and billing information.
    </p>
  </div>

  {#if loading}
    <div class="flex items-center justify-center min-h-[400px]">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
  {:else}
    {#if error}
      <div class="alert alert-error mb-6">
        <X class="w-5 h-5" />
        <span>{error}</span>
      </div>
    {/if}
    
    {#if profile}
    <div class="space-y-6">
      <!-- Account Settings Card -->
      <div class="card bg-base-100 shadow-lg">
        <div class="card-body">
          <h2 class="card-title text-xl mb-4">
            <Mail class="w-5 h-5" />
            Account Settings
          </h2>
          
          <div class="space-y-4">
            <!-- Current Email Display -->
            <div>
              <div class="label">
                <span class="label-text font-semibold">Current Email</span>
              </div>
              <div class="flex items-center gap-2">
                <input 
                  type="email" 
                  value={$user?.email || ''} 
                  disabled
                  class="input input-bordered flex-1 bg-base-200"
                />
                <button 
                  class="btn btn-outline btn-sm"
                  on:click={() => {
                    showEmailChangeForm = !showEmailChangeForm;
                    emailChangeError = null;
                    emailChangeSuccess = null;
                    newEmail = '';
                    passwordForEmailChange = '';
                  }}
                >
                  {showEmailChangeForm ? 'Cancel' : 'Change Email'}
                </button>
              </div>
            </div>

            <!-- Email Change Form -->
            {#if showEmailChangeForm}
              <div class="border-t border-base-300 pt-4 mt-4">
                {#if emailChangeError}
                  <div class="alert alert-error mb-4">
                    <X class="w-5 h-5" />
                    <span>{emailChangeError}</span>
                  </div>
                {/if}

                {#if emailChangeSuccess}
                  <div class="alert alert-success mb-4">
                    <Check class="w-5 h-5" />
                    <span>{emailChangeSuccess}</span>
                  </div>
                {/if}

                <div class="space-y-4">
                  <div>
                    <label class="label" for="new-email">
                      <span class="label-text font-semibold">New Email Address</span>
                    </label>
                    <input 
                      id="new-email"
                      type="email" 
                      placeholder="Enter new email address"
                      class="input input-bordered w-full"
                      bind:value={newEmail}
                      disabled={updatingEmail}
                    />
                  </div>

                  <div>
                    <label class="label" for="password-for-email-change">
                      <span class="label-text font-semibold">Confirm Password</span>
                      <span class="label-text-alt">Required to change email</span>
                    </label>
                    <input 
                      id="password-for-email-change"
                      type="password" 
                      placeholder="Enter your password"
                      class="input input-bordered w-full"
                      bind:value={passwordForEmailChange}
                      disabled={updatingEmail}
                    />
                    <div class="label">
                      <span class="label-text-alt text-base-content/60">
                        For security, you'll need to confirm the change via email on both your old and new email addresses.
                      </span>
                    </div>
                  </div>

                  <button 
                    class="btn btn-primary w-full"
                    on:click={handleEmailChange}
                    disabled={updatingEmail || !newEmail.trim() || !passwordForEmailChange}
                  >
                    {#if updatingEmail}
                      <Loader class="w-4 h-4 animate-spin" />
                      Updating Email...
                    {:else}
                      Update Email Address
                    {/if}
                  </button>
                </div>
              </div>
            {/if}
          </div>
        </div>
      </div>

      <!-- Current Plan Card -->
      <div class="card bg-base-100 shadow-lg">
        <div class="card-body">
          <h2 class="card-title text-xl mb-4">Current Plan</h2>
          
          <div class="flex items-center justify-between mb-4">
            <div>
              <div class="badge badge-lg badge-primary badge-outline uppercase">
                {profile.subscription_tier || 'free'}
              </div>
              {#if subscription}
                <p class="text-sm text-base-content/70 mt-2">
                  Status: <span class="badge badge-sm badge-success">{subscription.status}</span>
                </p>
              {/if}
            </div>
            
            {#if isProOrTeam && (!teamInfo || teamInfo.is_owner)}
              {#if profile?.stripe_subscription_id}
                <!-- Paid User: Manage Billing Button -->
                <button 
                  class="btn btn-primary"
                  on:click={openBillingPortal}
                  disabled={creatingPortal}
                >
                  {#if creatingPortal}
                    <Loader class="w-4 h-4 animate-spin" />
                  {:else}
                    <CreditCard class="w-4 h-4" />
                  {/if}
                  Manage Subscription
                </button>
              {:else}
                <!-- Beta User: Manage Team Button -->
                <label for="beta-team-modal" class="btn btn-primary btn-outline">
                  Manage Team Size
                </label>
              {/if}
            {/if}
          </div>

          <div class="grid grid-cols-2 gap-4 mt-4">
            <div>
              <p class="text-sm text-base-content/70">Crawl Limit</p>
              <p class="text-lg font-semibold">{planFeatures.pages}{planFeatures.pages !== 'Custom' ? ' pages' : ''}</p>
            </div>
            <div>
              <p class="text-sm text-base-content/70">Team Members</p>
              <p class="text-lg font-semibold">{planFeatures.users}</p>
            </div>
          </div>

          {#if teamInfo}
            <div class="divider my-4"></div>
            <div class="bg-base-200 rounded-lg p-4">
              {#if teamInfo.is_owner}
                <!-- Account Owner View -->
                <div class="flex items-center justify-between">
                  <div>
                    <h3 class="font-semibold text-sm mb-1">Team Management</h3>
                    <p class="text-xs text-base-content/70">
                      {teamInfo.active_count} of {teamInfo.team_size_limit} seat{teamInfo.team_size_limit === 1 ? '' : 's'} used
                    </p>
                  </div>
                  {#if profile?.stripe_subscription_id}
                    <button 
                      class="btn btn-primary btn-sm"
                      on:click={openBillingPortal}
                      disabled={creatingPortal}
                    >
                      {#if creatingPortal}
                        <Loader class="w-4 h-4 animate-spin" />
                      {:else}
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                        </svg>
                      {/if}
                      Add Seats
                    </button>
                  {/if}
                </div>
              {:else}
                <!-- Team Member View -->
                <div>
                  <h3 class="font-semibold text-sm mb-1">Team Member</h3>
                  <p class="text-xs text-base-content/70">
                    You're part of a team with {teamInfo.active_count} of {teamInfo.team_size_limit} seat{teamInfo.team_size_limit === 1 ? '' : 's'} used. Contact your account owner to manage billing.
                  </p>
                </div>
              {/if}
            </div>
          {/if}

          {#if subscription}
            <div class="divider my-4"></div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <p class="text-sm text-base-content/70">Current Period</p>
                <p class="text-sm">
                  {formatDate(subscription.current_period_start)} - {formatDate(subscription.current_period_end)}
                </p>
              </div>
              {#if subscription.cancel_at_period_end}
                <div>
                  <p class="text-sm text-warning">Cancels on</p>
                  <p class="text-sm">{formatDate(subscription.current_period_end)}</p>
                </div>
              {/if}
            </div>
          {/if}
        </div>
      </div>

      <!-- Plan Features -->
      <div class="card bg-base-100 shadow">
        <div class="card-body">
          <h2 class="card-title text-xl mb-4">Plan Features</h2>
          <div class="space-y-2">
            <div class="flex items-center gap-2">
              {#if planFeatures.integrations}
                <Check class="w-5 h-5 text-success" />
              {:else}
                <X class="w-5 h-5 text-base-content/30" />
              {/if}
              <span>Google Search Console & Analytics integrations</span>
            </div>
            <div class="flex items-center gap-2">
              {#if planFeatures.recommendations}
                <Check class="w-5 h-5 text-success" />
              {:else}
                <X class="w-5 h-5 text-base-content/30" />
              {/if}
              <span>AI-powered recommendations</span>
            </div>
            <div class="flex items-center gap-2">
              {#if isProOrTeam}
                <Check class="w-5 h-5 text-success" />
              {:else}
                <X class="w-5 h-5 text-base-content/30" />
              {/if}
              <span>Team collaboration</span>
            </div>
            <div class="flex items-center gap-2">
              {#if isProOrTeam}
                <Check class="w-5 h-5 text-success" />
              {:else}
                <X class="w-5 h-5 text-base-content/30" />
              {/if}
              <span>Priority support</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Team Management -->
      {#if isProOrTeam && (!teamInfo || teamInfo.is_owner)}
        <TeamManagement />
      {/if}

      <!-- Upgrade Options -->
      {#if !isProOrTeam && (!teamInfo || teamInfo.is_owner)}
        <div class="card bg-base-100 shadow">
          <div class="card-body">
            <h2 class="card-title text-xl mb-4">Upgrade Plan</h2>
            <p class="text-base-content/70 mb-4">
              Unlock more features with a Pro subscription.
            </p>
            
            <!-- Billing Period Toggle -->
            <div class="flex justify-center mb-6">
              <div class="btn-group">
                <button 
                  class="btn btn-sm {selectedBillingPeriod === 'monthly' ? 'btn-primary' : 'btn-outline'}"
                  on:click={() => selectedBillingPeriod = 'monthly'}
                >
                  Monthly
                </button>
                <button 
                  class="btn btn-sm {selectedBillingPeriod === 'annual' ? 'btn-primary' : 'btn-outline'}"
                  on:click={() => selectedBillingPeriod = 'annual'}
                >
                  Annual
                  <span class="badge badge-success badge-sm ml-2">Save 20%</span>
                </button>
              </div>
            </div>
            
            <div class="bg-primary/10 rounded-lg p-4 mb-4">
              {#if selectedBillingPeriod === 'monthly'}
                <h3 class="font-semibold mb-2">Pro Plan - $29/month</h3>
              {:else}
                <h3 class="font-semibold mb-2">Pro Plan - Annual</h3>
                <p class="text-sm text-base-content/70 mb-2">Billed annually, save 20%</p>
              {/if}
              <ul class="text-sm space-y-1 mb-4">
                <li>✓ Crawl up to <strong>10,000 pages</strong></li>
                <li>✓ <strong>Team collaboration</strong> — invite teammates with role-based permissions (1 user included, +$5/user/month)</li>
                <li>✓ Integrations: Google Search Console, Analytics, Clarity, Slack</li>
                <li>✓ Full recommendation engine with contextual fixes</li>
                <li>✓ Historical comparisons and advanced exports</li>
                <li>✓ CLI automation and priority support</li>
              </ul>
              <p class="text-xs text-base-content/60 mt-2">
                Need 5+ team members or custom crawl limits? <a href="mailto:sales@barracudaseo.com" class="link link-primary">Contact Sales</a> for custom Team plans.
              </p>
            </div>

            <!-- Team Seats Selection -->
            {#if STRIPE_PRICE_ID_TEAM_SEAT}
              <div class="mb-4 p-4 border border-base-300 rounded-lg">
                <label class="label" for="team-seats-quantity">
                  <span class="label-text font-semibold">Additional Team Seats</span>
                  <span class="label-text-alt">$5/month each</span>
                </label>
                <div class="flex items-center gap-4 mt-2">
                  <input 
                    id="team-seats-quantity"
                    type="number" 
                    min="0" 
                    max="20" 
                    bind:value={teamSeatsQuantity}
                    class="input input-bordered w-24"
                    placeholder="0"
                  />
                  <span class="text-sm text-base-content/70">
                    {#if teamSeatsQuantity > 0}
                      {teamSeatsQuantity} seat{teamSeatsQuantity === 1 ? '' : 's'} = ${teamSeatsQuantity * 5}/month
                    {:else}
                      No additional seats
                    {/if}
                  </span>
                </div>
                <p class="text-xs text-base-content/60 mt-2">
                  Pro plan includes 1 user. Add more seats for your team members at $5/user/month.
                </p>
                <p class="text-xs text-base-content/50 mt-1">
                  Need 5+ team members or custom crawl limits? <a href="mailto:sales@barracudaseo.com" class="link link-primary">Contact Sales</a> for custom Team plans.
                </p>
              </div>
            {/if}

            <button 
              class="btn btn-primary w-full"
              on:click={() => {
                const priceId = selectedBillingPeriod === 'monthly' 
                  ? STRIPE_PRICE_ID_PRO 
                  : STRIPE_PRICE_ID_PRO_ANNUAL;
                createCheckoutSession(priceId);
              }}
              disabled={creatingCheckout || (!STRIPE_PRICE_ID_PRO && !STRIPE_PRICE_ID_PRO_ANNUAL)}
            >
              {#if creatingCheckout}
                <Loader class="w-4 h-4 animate-spin" />
                Processing...
              {:else}
                {#if teamSeatsQuantity > 0}
                  Upgrade to Pro {selectedBillingPeriod === 'annual' ? '(Annual)' : ''} + {teamSeatsQuantity} seat{teamSeatsQuantity === 1 ? '' : 's'}
                {:else}
                  Upgrade to Pro {selectedBillingPeriod === 'annual' ? '(Annual)' : ''}
                {/if}
              {/if}
            </button>

            {#if !STRIPE_PRICE_ID_PRO && !STRIPE_PRICE_ID_PRO_ANNUAL}
              <p class="text-sm text-warning mt-2">
                Stripe is not configured. Please set VITE_STRIPE_PRICE_ID_PRO and VITE_STRIPE_PRICE_ID_PRO_ANNUAL environment variables.
              </p>
            {/if}

            <!-- Beta Code Redemption -->
            <div class="divider mt-6">OR</div>
            
            <div class="collapse collapse-arrow border border-base-300 bg-base-100 rounded-box">
              <input type="checkbox" /> 
              <div class="collapse-title text-sm font-medium">
                Have a beta invite code?
              </div>
              <div class="collapse-content"> 
                <div class="form-control pt-2 gap-4">
                  <!-- Team Size Input -->
                  <div class="form-control w-full">
                    <label class="label" for="redeem-team-size">
                      <span class="label-text">Team Size (Optional)</span>
                      <span class="label-text-alt">Max 10 for Beta</span>
                    </label>
                    <div class="join">
                      <button class="btn join-item" on:click={() => redeemTeamSize = Math.max(1, redeemTeamSize - 1)}>-</button>
                      <input id="redeem-team-size" type="number" min="1" max="10" bind:value={redeemTeamSize} class="input input-bordered join-item w-full text-center" />
                      <button class="btn join-item" on:click={() => redeemTeamSize = Math.min(10, redeemTeamSize + 1)}>+</button>
                    </div>
                    <div class="label">
                      <span class="label-text-alt">Includes {redeemTeamSize} user{redeemTeamSize > 1 ? 's' : ''}</span>
                    </div>
                  </div>

                  <div class="join w-full">
                    <input 
                      type="text" 
                      placeholder="Enter code" 
                      class="input input-bordered join-item w-full"
                      bind:value={redeemCode}
                    />
                    <button 
                      class="btn btn-primary join-item" 
                      on:click={redeemBetaCode}
                      disabled={redeeming || !redeemCode}
                    >
                      {#if redeeming}
                        <Loader class="w-4 h-4 animate-spin" />
                      {:else}
                        Redeem
                      {/if}
                    </button>
                  </div>
                  {#if redeemError}
                    <div class="label">
                      <span class="label-text-alt text-error">{redeemError}</span>
                    </div>
                  {/if}
                </div>
              </div>
            </div>
          </div>
        </div>
      {/if}
    </div>

    <!-- Modal for Beta Users to Update Team Size -->
    <input 
      type="checkbox" 
      id="beta-team-modal" 
      class="modal-toggle"
      on:change={(e) => {
        // Initialize team size with current value when modal opens
        if (e.target.checked && profile?.team_size) {
          redeemTeamSize = profile.team_size || 1;
        }
      }}
    />
    <div class="modal" role="dialog">
      <div class="modal-box">
        <h3 class="font-bold text-lg">Update Team Size</h3>
        <p class="py-4">Enter your beta invite code again to update your team size.</p>
        
        <div class="form-control gap-4">
          <div>
            <label class="label" for="beta-team-size">
              <span class="label-text">New Team Size</span>
              {#if profile?.team_size}
                <span class="label-text-alt">Current: {profile.team_size}</span>
              {/if}
            </label>
            <div class="join w-full">
              <button class="btn join-item" on:click={() => redeemTeamSize = Math.max(1, redeemTeamSize - 1)}>-</button>
              <input id="beta-team-size" type="number" min="1" max="10" bind:value={redeemTeamSize} class="input input-bordered join-item w-full text-center" />
              <button class="btn join-item" on:click={() => redeemTeamSize = Math.min(10, redeemTeamSize + 1)}>+</button>
            </div>
          </div>

          <div>
            <label class="label" for="beta-invite-code">
              <span class="label-text">Beta Invite Code</span>
            </label>
            <input 
              id="beta-invite-code"
              type="text" 
              placeholder="Enter code to confirm" 
              class="input input-bordered w-full"
              bind:value={redeemCode}
            />
          </div>
          
          {#if redeemError}
            <div class="text-error text-sm">{redeemError}</div>
          {/if}
        </div>

        <div class="modal-action">
          <label for="beta-team-modal" class="btn">Cancel</label>
          <button 
            class="btn btn-primary" 
            on:click={redeemBetaCode}
            disabled={redeeming || !redeemCode}
          >
            {#if redeeming}
              <Loader class="w-4 h-4 animate-spin" />
            {/if}
            Update Team
          </button>
        </div>
      </div>
    </div>

    {:else}
      <!-- No profile loaded - show loading or error message -->
      <div class="alert alert-warning">
        <span>Unable to load subscription information. Please try refreshing the page.</span>
      </div>
    {/if}
  {/if}
</div>

<style>
  :global(.badge-success) {
    background-color: #8ec07c;
    color: white;
  }
</style>



