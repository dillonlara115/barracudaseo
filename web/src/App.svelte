<script>
  import { onMount } from 'svelte';
  import Router, { link, push } from 'svelte-spa-router';
  import { wrap } from 'svelte-spa-router/wrap';
  import { initAuth, user, authEvent, session } from './lib/auth.js';
  import { supabase } from './lib/supabase.js';
  import Auth from './components/Auth.svelte';
  import ConfigError from './components/ConfigError.svelte';
  import { loadSubscriptionData } from './lib/subscription.js';

  let loading = true;
  let configError = null;
  let currentHash = typeof window !== 'undefined' ? window.location.hash : '';
  
  $: isLegalPage = currentHash === '#/privacy' || currentHash === '#/terms';
  $: isPublicPage = isLegalPage || currentHash.startsWith('#/team/accept') || currentHash.startsWith('#/reports/') || currentHash.startsWith('#/auth/confirm');

  // Route definitions with dynamic imports (code splitting)
  // Using wrap() for proper async component loading
  const routes = {
    '/': wrap({
      asyncComponent: () => import('./routes/ProjectsList.svelte')
    }),
    '/project/:id': wrap({
      asyncComponent: () => import('./routes/ProjectView.svelte')
    }),
    '/project/:projectId/crawl/:crawlId': wrap({
      asyncComponent: () => import('./routes/CrawlView.svelte')
    }),
    '/project/:projectId/settings': wrap({
      asyncComponent: () => import('./routes/Settings.svelte')
    }),
    '/integrations': wrap({
      asyncComponent: () => import('./routes/IntegrationsProtected.svelte')
    }),
    '/billing': wrap({
      asyncComponent: () => import('./routes/Billing.svelte')
    }),
    '/settings': wrap({
      asyncComponent: () => import('./routes/Billing.svelte') // Alias for billing
    }),
    '/privacy': wrap({
      asyncComponent: () => import('./routes/PrivacyPolicy.svelte')
    }),
    '/terms': wrap({
      asyncComponent: () => import('./routes/TermsOfService.svelte')
    }),
    '/project/:projectId/gsc': wrap({
      asyncComponent: () => import('./routes/GSCDashboard.svelte')
    }),
    '/project/:projectId/gsc/keywords': wrap({
      asyncComponent: () => import('./routes/GSCKeywords.svelte')
    }),
    '/project/:projectId/rank-tracker': wrap({
      asyncComponent: () => import('./routes/RankTracker.svelte')
    }),
    '/project/:projectId/discover-keywords': wrap({
      asyncComponent: () => import('./routes/DiscoverKeywords.svelte')
    }),
    '/project/:projectId/impact-first': wrap({
      asyncComponent: () => import('./routes/ImpactFirstView.svelte')
    }),
    '/project/:projectId/crawls': wrap({
      asyncComponent: () => import('./routes/Crawls.svelte')
    }),
    '/team/accept': wrap({
      asyncComponent: () => import('./routes/TeamAccept.svelte')
    }),
    '/reports/:token': wrap({
      asyncComponent: () => import('./routes/PublicReportView.svelte')
    }),
    '/auth': Auth, // Auth route - keep static as it's needed immediately
    '/auth/confirm': wrap({
      asyncComponent: () => import('./routes/AuthConfirm.svelte') // PKCE magic link verification endpoint
    }),
    '/reset': wrap({
      asyncComponent: () => import('./routes/ResetPassword.svelte') // Password reset flow after Supabase recovery email
    }),
  };

  // Check Supabase configuration
  $: {
    const supabaseUrl = import.meta.env.PUBLIC_SUPABASE_URL || import.meta.env.VITE_PUBLIC_SUPABASE_URL;
    const supabaseAnonKey = import.meta.env.PUBLIC_SUPABASE_ANON_KEY || import.meta.env.VITE_PUBLIC_SUPABASE_ANON_KEY;
    
    if (!supabaseUrl || !supabaseAnonKey) {
      configError = 'Missing Supabase configuration. Please set PUBLIC_SUPABASE_URL and PUBLIC_SUPABASE_ANON_KEY environment variables.';
    } else {
      configError = null;
    }
  }

  // Track redirect timeout to prevent multiple redirects
  let redirectTimeout = null;
  let isHandlingRedirect = false;
  let lastAuthEvent = null;

  onMount(async () => {
    // Track auth events to differentiate real sign-outs from refresh failures
    authEvent.subscribe((event) => {
      lastAuthEvent = event;
    });
    // Watch hash changes for legal pages
    if (typeof window !== 'undefined') {
      currentHash = window.location.hash;
      const updateHash = () => {
        currentHash = window.location.hash;
      };
      window.addEventListener('hashchange', updateHash);
      
      // Also check periodically in case hashchange doesn't fire
      setInterval(updateHash, 100);
    }
    
    // STEP 1: Handle /auth/confirm path redirect (legacy PKCE support)
    // Convert regular path to hash route for SPA routing
    // This is only needed if someone uses PKCE flow with token_hash
    if (window.location.pathname === '/auth/confirm') {
      const queryParams = window.location.search;
      const hashParams = window.location.hash;
      
      if (queryParams) {
        console.log('🔍 Converting /auth/confirm path with query params to hash route');
        // Convert /auth/confirm?token_hash=... to /#/auth/confirm?token_hash=...
        window.location.replace(`${window.location.origin}/#/auth/confirm${queryParams}`);
        return; // Exit early, let the redirect happen
      } else if (hashParams && hashParams.includes('access_token=')) {
        // If tokens are in hash, extract them and redirect to root
        console.log('🔍 Converting /auth/confirm path - extracting tokens to root');
        // Extract tokens and redirect to root where App.svelte will handle them
        window.location.replace(`${window.location.origin}/#/${hashParams.substring(1)}`);
        return;
      } else {
        // No params, redirect to root
        window.location.replace(`${window.location.origin}/#/`);
        return;
      }
    }
    
    // Fix double hash issue if present (e.g., #/#/billing -> #/billing)
    if (window.location.hash.startsWith('#/#/')) {
      const fixedHash = window.location.hash.replace('#/#/', '#/');
      window.location.hash = fixedHash;
      console.log('Auto-fixed hash from #/#/ to:', fixedHash);
    }
    
    // STEP 2: Check for auth callback tokens in hash (implicit flow)
    // Magic links use URL fragments: #access_token=...&refresh_token=...&type=magiclink
    // Supabase redirects to: origin#access_token=... (no hash route initially)
    // They can appear as: #access_token=... or #/#access_token=... or #/route#access_token=...
    const fullHash = window.location.hash;
    console.log('Full URL hash:', fullHash || '<empty string>');
    console.log('Full URL pathname:', window.location.pathname);
    console.log('Full URL search:', window.location.search);
    
    // Extract auth parameters from hash (handle multiple # symbols)
    let authParams = new URLSearchParams();
    let hasAuthToken = false;
    
    // Check if hash contains access_token (implicit flow)
    if (fullHash.includes('access_token=')) {
      // Find the auth token portion
      // Handle cases like: 
      // - #access_token=... (direct from Supabase)
      // - #/#access_token=... (double hash)
      // - #/route#access_token=... (route + tokens)
      const accessTokenIndex = fullHash.indexOf('access_token=');
      
      // Extract everything from access_token onwards
      let authFragment;
      if (accessTokenIndex > 0 && fullHash[accessTokenIndex - 1] === '#') {
        // Case: #/route#access_token=... - extract from the second #
        authFragment = fullHash.substring(accessTokenIndex - 1);
        // Remove the leading # to make it a proper fragment
        authFragment = authFragment.substring(1);
      } else {
        // Case: #access_token=... - extract from access_token
        authFragment = fullHash.substring(accessTokenIndex);
      }
      
      authParams = new URLSearchParams(authFragment);
      hasAuthToken = true;
      console.log('✅ Found auth token in URL');
    }
    
    if (hasAuthToken) {
      const accessToken = authParams.get('access_token');
      const refreshToken = authParams.get('refresh_token');
      const tokenType = authParams.get('type');
      
      console.log('Auth callback detected:', { type: tokenType, hasAccessToken: !!accessToken, hasRefreshToken: !!refreshToken });
      
      if (accessToken && refreshToken) {
        try {
          // Set the session using tokens from magic link
          console.log('🔄 Setting session from magic link tokens...');
          const { data, error } = await supabase.auth.setSession({
            access_token: accessToken,
            refresh_token: refreshToken
          });
          
          if (error) {
            console.error('Auth callback error:', error);
          } else {
            console.log('✅ Magic link session set successfully:', data);
            
            // STEP 3: Wait for session to be fully established and propagated
            // Give Supabase time to persist the session and update all stores
            await new Promise(resolve => setTimeout(resolve, 300));
            
            // Verify session is actually set
            const { data: { session: verifySession } } = await supabase.auth.getSession();
            if (!verifySession) {
              console.error('⚠️ Session was not persisted after setSession');
            } else {
              console.log('✅ Session verified and persisted');
            }
            
            // Clean up URL - remove tokens, redirect to dashboard
            // After setting session, redirect to app root with clean hash
            window.history.replaceState(null, '', window.location.pathname);
            window.location.hash = '#/';
            
            // Don't return here - let initAuth() run to ensure everything is synced
          }
        } catch (err) {
          console.error('Failed to set session:', err);
        }
      }
    }

    // STEP 3: Initialize auth (this will sync the session state)
    await initAuth();

    // React to auth state changes
    user.subscribe(async (currentUser) => {
      // Clear any pending redirect
      if (redirectTimeout) {
        clearTimeout(redirectTimeout);
        redirectTimeout = null;
      }

      if (!currentUser) {
        // Don't redirect if we're already handling one
        if (isHandlingRedirect) {
          return;
        }

        // Check if this is a real sign-out or just a token refresh failure
        // Only redirect on explicit SIGNED_OUT events, not on TOKEN_REFRESHED failures
        const eventType = lastAuthEvent;
        
        // If the event is SIGNED_OUT, it's a real sign-out
        // If it's TOKEN_REFRESHED or null, it might be a temporary failure
        if (eventType === 'SIGNED_OUT') {
          // Real sign-out - redirect immediately (but still check for public pages)
          if (!isPublicPage) {
            isHandlingRedirect = true;
            push('/');
            isHandlingRedirect = false;
          }
          return;
        }

        // For other cases (TOKEN_REFRESHED failure, null event, etc.), use a delay
        // and verify there's truly no session before redirecting
        redirectTimeout = setTimeout(async () => {
          // Double-check session still doesn't exist
          const { data: { session: storedSession } } = await supabase.auth.getSession();
          
          // Also check the session store
          const currentSession = $session;
          
          // Only redirect if there's truly no session (real sign-out)
          // Don't redirect if on public pages (legal pages or invite acceptance)
          // Don't redirect if session was restored
          if (!isPublicPage && !storedSession && !currentSession && !isHandlingRedirect) {
            isHandlingRedirect = true;
            push('/');
            // Reset flag after a delay to allow navigation
            setTimeout(() => {
              isHandlingRedirect = false;
            }, 2000);
          } else {
            // Session was restored or we're on a public page - don't redirect
            redirectTimeout = null;
          }
        }, 2000); // Increased delay to 2 seconds to allow token refresh to complete
      } else {
        // Load subscription data when user is authenticated
        // Add a small delay to ensure session is fully propagated before making API calls
        await new Promise(resolve => setTimeout(resolve, 200));
        await loadSubscriptionData();
        
        // If user just authenticated and is on /auth route, check for invite token and redirect
        if (currentHash.startsWith('#/auth')) {
          let params;
          if (currentHash.includes('?')) {
            const hashPart = currentHash.split('?')[1];
            params = new URLSearchParams(hashPart);
          } else {
            params = new URLSearchParams(window.location.search);
          }
          const inviteToken = params.get('invite_token');
          if (inviteToken) {
            // Redirect to accept invite page
            setTimeout(() => {
              push(`#/team/accept?token=${inviteToken}`);
            }, 500);
          } else {
            // Redirect to home
            setTimeout(() => {
              push('#/');
            }, 500);
          }
        }
      }
      loading = false;
    });
  });
</script>

<div class="min-h-screen bg-base-100">
  {#if configError}
    <!-- Show configuration error -->
    <div class="flex items-center justify-center min-h-screen p-4">
      <ConfigError error={configError} />
    </div>
  {:else if !$user}
    <!-- Show auth UI when not logged in, or public pages -->
    {#if isPublicPage}
      <!-- Public pages - accessible without login (legal pages, invite acceptance) -->
      <Router {routes} />
    {:else}
      <Auth />
    {/if}
  {:else if loading}
    <!-- Loading state -->
    <div class="flex items-center justify-center min-h-screen">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
  {:else}
    <!-- Router handles all routes (including public legal pages) -->
    <!-- Using hash mode for routing -->
    <Router {routes} />
  {/if}
</div>
