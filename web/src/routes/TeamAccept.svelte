<script>
  import { onMount } from 'svelte';
  import { push } from 'svelte-spa-router';
  import { supabase } from '../lib/supabase.js';
  import { user } from '../lib/auth.js';
  import { Check, X, Loader, AlertCircle, Mail, UserPlus, LogIn } from 'lucide-svelte';

  let loading = true;
  let error = null;
  let success = false;
  let token = null;
  let inviteDetails = null;
  let needsAuth = false;

  const API_URL = import.meta.env.VITE_CLOUD_RUN_API_URL || 'http://localhost:8080';

  onMount(async () => {
    // Get token from URL - check both hash and search params
    // Hash-based routing puts query params in the hash
    let params;
    if (window.location.hash.includes('?')) {
      // Token is in hash: #/team/accept?token=...
      const hashPart = window.location.hash.split('?')[1];
      params = new URLSearchParams(hashPart);
    } else {
      // Fallback to search params (for non-hash URLs)
      params = new URLSearchParams(window.location.search);
    }
    token = params.get('token');

    if (!token) {
      error = 'No invite token provided';
      loading = false;
      return;
    }

    // Fetch invite details first
    await loadInviteDetails();

    // Check if user is authenticated
    if (!$user) {
      needsAuth = true;
      loading = false;
      return;
    }

    // User is signed in - accept the invite
    await acceptInvite();
  });

  async function loadInviteDetails() {
    try {
      const response = await fetch(`${API_URL}/api/v1/team/${encodeURIComponent(token)}/details`);
      
      if (!response.ok) {
        const errorData = await response.json().catch(() => null);
        throw new Error(errorData?.error || 'Failed to load invite details');
      }

      inviteDetails = await response.json();
      loading = false; // Set loading to false on success
    } catch (err) {
      error = err.message || 'Failed to load invite details';
      loading = false;
    }
  }

  async function acceptInvite() {
    if (!token) {
      error = 'No invite token provided';
      loading = false;
      return;
    }

    loading = true;
    error = null;

    try {
      const { data: sessionData, error: sessionError } = await supabase.auth.getSession();
      if (sessionError || !sessionData.session) {
        throw new Error('Not authenticated. Please sign in first.');
      }

      const accessToken = sessionData.session.access_token;
      const response = await fetch(`${API_URL}/api/v1/team/${encodeURIComponent(token)}/accept`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${accessToken}`,
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Failed to accept invite');
      }

      success = true;
      
      // Redirect to billing/team page after 2 seconds
      setTimeout(() => {
        push('#/billing');
      }, 2000);
    } catch (err) {
      error = err.message || 'Failed to accept invite';
    } finally {
      loading = false;
    }
  }

  // Watch for auth state changes - if user signs in, accept invite
  $: if ($user && token && !success && !loading && !error && !needsAuth) {
    acceptInvite();
  }

  function goToAuth() {
    push(`#/auth?invite_token=${token}`);
  }
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-50 px-4">
  <div class="max-w-md w-full bg-white rounded-lg shadow-lg p-8">
    {#if loading}
      <div class="flex flex-col items-center justify-center space-y-4">
        <Loader class="w-12 h-12 text-blue-600 animate-spin" />
        <p class="text-gray-600">Loading invite details...</p>
      </div>
    {:else if success}
      <div class="flex flex-col items-center justify-center space-y-4">
        <div class="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center">
          <Check class="w-8 h-8 text-green-600" />
        </div>
        <h2 class="text-2xl font-bold text-gray-900">Invite Accepted!</h2>
        <p class="text-gray-600 text-center">
          You've successfully joined the team. Redirecting you to the billing page...
        </p>
      </div>
    {:else if error}
      <div class="flex flex-col items-center justify-center space-y-4">
        <div class="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center">
          <AlertCircle class="w-8 h-8 text-red-600" />
        </div>
        <h2 class="text-2xl font-bold text-gray-900">Error</h2>
        <p class="text-gray-600 text-center">{error}</p>
        <div class="flex space-x-4 mt-4">
          <button
            on:click={() => push('#/auth')}
            class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
          >
            Sign In
          </button>
          <button
            on:click={() => { error = null; loading = true; loadInviteDetails(); }}
            class="px-4 py-2 bg-gray-200 text-gray-800 rounded-lg hover:bg-gray-300"
          >
            Try Again
          </button>
        </div>
      </div>
    {:else if needsAuth && inviteDetails}
      <!-- Show invite details and prompt for account creation/login -->
      <div class="flex flex-col items-center justify-center space-y-6">
        <div class="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center">
          <Mail class="w-8 h-8 text-blue-600" />
        </div>
        <div class="text-center">
          <h2 class="text-2xl font-bold text-gray-900 mb-2">You've Been Invited!</h2>
          <p class="text-gray-600 mb-4">
            You've been invited to join a team on Barracuda SEO.
          </p>
          {#if inviteDetails.email}
            <p class="text-sm text-gray-500 mb-6">
              Invited email: <span class="font-medium">{inviteDetails.email}</span>
            </p>
          {/if}
        </div>
        
        <div class="w-full bg-blue-50 border border-blue-200 rounded-lg p-4 mb-4">
          <div class="flex items-start space-x-3">
            <UserPlus class="w-5 h-5 text-blue-600 mt-0.5 flex-shrink-0" />
            <div>
              <p class="text-sm font-medium text-blue-900 mb-1">Account Activated</p>
              <p class="text-sm text-blue-700">
                Your account has been activated. Please complete your account setup to accept this invitation.
              </p>
            </div>
          </div>
        </div>

        <div class="w-full space-y-3">
          <button
            on:click={goToAuth}
            class="w-full px-4 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium flex items-center justify-center space-x-2"
          >
            <LogIn class="w-5 h-5" />
            <span>Complete Account Setup</span>
          </button>
          <p class="text-xs text-gray-500 text-center">
            You'll be able to accept the invitation after signing in or creating your account.
          </p>
        </div>
      </div>
    {:else}
      <div class="flex flex-col items-center justify-center space-y-4">
        <Loader class="w-12 h-12 text-blue-600 animate-spin" />
        <p class="text-gray-600">Loading...</p>
      </div>
    {/if}
  </div>
</div>

