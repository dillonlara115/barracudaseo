<script>
  import { onMount } from 'svelte';
  import Router, { link, push } from 'svelte-spa-router';
  import { initAuth, user, authEvent, session } from './lib/auth.js';
  import { supabase } from './lib/supabase.js';
  import Auth from './components/Auth.svelte';
  import ConfigError from './components/ConfigError.svelte';
  import ProjectsList from './routes/ProjectsList.svelte';
  import ProjectView from './routes/ProjectView.svelte';
  import CrawlView from './routes/CrawlView.svelte';
  import IntegrationsProtected from './routes/IntegrationsProtected.svelte';
  import Settings from './routes/Settings.svelte';
  import Billing from './routes/Billing.svelte';
  import PrivacyPolicy from './routes/PrivacyPolicy.svelte';
  import TermsOfService from './routes/TermsOfService.svelte';
  import GSCDashboard from './routes/GSCDashboard.svelte';
  import GSCKeywords from './routes/GSCKeywords.svelte';
  import TeamAccept from './routes/TeamAccept.svelte';
  import PublicReportView from './routes/PublicReportView.svelte';
  import RankTracker from './routes/RankTracker.svelte';
  import ImpactFirstView from './routes/ImpactFirstView.svelte';
  import DiscoverKeywords from './routes/DiscoverKeywords.svelte';
  import Crawls from './routes/Crawls.svelte';
  import ResetPassword from './routes/ResetPassword.svelte';
  import AuthConfirm from './routes/AuthConfirm.svelte';
  import { loadSubscriptionData } from './lib/subscription.js';

  let loading = true;
  let configError = null;
  let currentHash = typeof window !== 'undefined' ? window.location.hash : '';
  
  $: isLegalPage = currentHash === '#/privacy' || currentHash === '#/terms';
  $: isPublicPage = isLegalPage || currentHash.startsWith('#/team/accept') || currentHash.startsWith('#/reports/') || currentHash.startsWith('#/auth/confirm');

  // Route definitions
  const routes = {
    '/': ProjectsList,
    '/project/:id': ProjectView,
    '/project/:projectId/crawl/:crawlId': CrawlView,
    '/project/:projectId/settings': Settings,
    '/integrations': IntegrationsProtected,
    '/billing': Billing,
    '/settings': Billing, // Alias for billing
    '/privacy': PrivacyPolicy,
    '/terms': TermsOfService,
    '/project/:projectId/gsc': GSCDashboard,
    '/project/:projectId/gsc/keywords': GSCKeywords,
    '/project/:projectId/rank-tracker': RankTracker,
    '/project/:projectId/discover-keywords': DiscoverKeywords,
    '/project/:projectId/impact-first': ImpactFirstView,
    '/project/:projectId/crawls': Crawls,
    '/team/accept': TeamAccept,
    '/reports/:token': PublicReportView,
    '/auth': Auth, // Auth route for when user is authenticated but needs to redirect
    '/auth/confirm': AuthConfirm, // PKCE magic link verification endpoint
    '/reset': ResetPassword, // Password reset flow after Supabase recovery email
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
    
    // Handle PKCE magic link redirect (Supabase redirects to /auth/confirm?token_hash=...)
    // Convert regular path to hash route for SPA routing
    if (window.location.pathname === '/auth/confirm' && window.location.search) {
      const queryParams = window.location.search;
      console.log('🔍 PKCE redirect detected - converting to hash route');
      console.log('🔍 Query params:', queryParams);
      // Convert /auth/confirm?token_hash=... to /#/auth/confirm?token_hash=...
      window.location.replace(`${window.location.origin}/#/auth/confirm${queryParams}`);
      return; // Exit early, let the redirect happen
    }
    
    // Fix double hash issue if present (e.g., #/#/billing -> #/billing)
    if (window.location.hash.startsWith('#/#/')) {
      const fixedHash = window.location.hash.replace('#/#/', '#/');
      window.location.hash = fixedHash;
      console.log('Auto-fixed hash from #/#/ to:', fixedHash);
    }
    
    // Check for auth callback (magic link, email confirmation, password reset)
    // Magic links use URL fragments: #access_token=...&refresh_token=...&type=magiclink
    // They can appear in different formats:
    // 1. Direct: https://app.barracudaseo.com#access_token=...
    // 2. With hash route: https://app.barracudaseo.com/#/#access_token=...
    // 3. With path: https://app.barracudaseo.com/#/some-path#access_token=...
    
    const fullHash = window.location.hash;
    console.log('Full URL hash:', fullHash || '<empty string>');
    console.log('Full URL pathname:', window.location.pathname);
    console.log('Full URL search:', window.location.search);
    
    // Extract auth parameters from hash (handle multiple # symbols)
    let authParams = new URLSearchParams();
    let hasAuthToken = false;
    
    // Check if hash contains access_token
    if (fullHash.includes('access_token=')) {
      // Find the auth token portion (everything after the last # or first access_token)
      const accessTokenIndex = fullHash.indexOf('access_token=');
      const authFragment = fullHash.substring(accessTokenIndex);
      authParams = new URLSearchParams(authFragment);
      hasAuthToken = true;
      console.log('Found auth token in URL');
    }
    
    if (hasAuthToken) {
      const accessToken = authParams.get('access_token');
      const refreshToken = authParams.get('refresh_token');
      const tokenType = authParams.get('type');
      
      console.log('Auth callback detected:', { type: tokenType, hasAccessToken: !!accessToken, hasRefreshToken: !!refreshToken });
      
      if (accessToken && refreshToken) {
        try {
          // Set the session using tokens from magic link
          const { data, error } = await supabase.auth.setSession({
            access_token: accessToken,
            refresh_token: refreshToken
          });
          
          if (error) {
            console.error('Auth callback error:', error);
          } else {
            console.log('Magic link session set successfully:', data);
            // Clean up URL and redirect to dashboard
            window.history.replaceState(null, '', window.location.pathname + '#/');
            // Force reload to ensure app state is fresh
            window.location.hash = '#/';
          }
        } catch (err) {
          console.error('Failed to set session:', err);
        }
      }
    }

    // Initialize auth
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
