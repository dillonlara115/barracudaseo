<script>
  import { onMount } from 'svelte';
  import { supabase } from '../lib/supabase.js';
  import { getApiUrl } from '../lib/data.js';

  let status = 'loading';
  let error = null;
  let userEmail = '';
  let callbackUrl = '';
  let state = '';

  const supabaseUrl = import.meta.env.PUBLIC_SUPABASE_URL || import.meta.env.VITE_PUBLIC_SUPABASE_URL || '';
  const supabaseAnonKey = import.meta.env.PUBLIC_SUPABASE_ANON_KEY || import.meta.env.VITE_PUBLIC_SUPABASE_ANON_KEY || '';

  function getQueryParams() {
    if (typeof window === 'undefined') return new URLSearchParams();
    const hash = window.location.hash || '';
    if (hash.includes('?')) {
      return new URLSearchParams(hash.split('?')[1]);
    }
    return new URLSearchParams(window.location.search);
  }

  function isLocalCallback(urlString) {
    try {
      const parsed = new URL(urlString);
      if (!['http:', 'https:'].includes(parsed.protocol)) {
        return false;
      }
      const host = parsed.hostname;
      return host === 'localhost' || host === '127.0.0.1' || host === '::1';
    } catch {
      return false;
    }
  }

  async function sendTokens() {
    const params = getQueryParams();
    callbackUrl = params.get('callback') || '';
    state = params.get('state') || '';

    if (!callbackUrl || !state) {
      error = 'Missing callback or state. Please restart CLI login.';
      status = 'error';
      return;
    }

    if (!isLocalCallback(callbackUrl)) {
      error = 'Invalid callback. Only localhost callbacks are allowed.';
      status = 'error';
      return;
    }

    if (!supabaseUrl || !supabaseAnonKey) {
      error = 'Supabase configuration is missing. Contact support.';
      status = 'error';
      return;
    }

    const { data: { session } } = await supabase.auth.getSession();
    if (!session || !session.access_token || !session.refresh_token) {
      error = 'You are not logged in yet. Please sign in first.';
      status = 'error';
      return;
    }

    userEmail = session.user?.email || '';

    const payload = {
      access_token: session.access_token,
      refresh_token: session.refresh_token,
      expires_at: session.expires_at,
      token_type: session.token_type,
      supabase_url: supabaseUrl,
      supabase_anon_key: supabaseAnonKey,
      api_url: getApiUrl(),
      user_id: session.user?.id || '',
      user_email: session.user?.email || '',
      state
    };

    try {
      const response = await fetch(callbackUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(payload)
      });

      if (!response.ok) {
        throw new Error(`Callback failed (${response.status})`);
      }

      status = 'success';
    } catch (err) {
      error = err.message || 'Failed to send credentials to CLI.';
      status = 'error';
    }
  }

  onMount(() => {
    sendTokens();
  });
</script>

<div class="min-h-screen flex items-center justify-center bg-base-200 px-4">
  <div class="card w-full max-w-md bg-base-100 shadow-xl">
    <div class="card-body space-y-4">
      <h1 class="text-2xl font-bold">Linking CLI</h1>

      {#if status === 'loading'}
        <div class="flex items-center gap-3">
          <span class="loading loading-spinner loading-sm"></span>
          <span>Connecting to your Barracuda CLI...</span>
        </div>
      {:else if status === 'success'}
        <div class="alert alert-success">
          <span>✅ CLI linked successfully. You can return to your terminal.</span>
        </div>
        {#if userEmail}
          <p class="text-sm text-base-content/70">Signed in as {userEmail}</p>
        {/if}
      {:else}
        <div class="alert alert-error">
          <span>{error}</span>
        </div>
        <p class="text-sm text-base-content/70">
          Close this tab and run <code>barracuda auth login</code> again.
        </p>
      {/if}
    </div>
  </div>
</div>
