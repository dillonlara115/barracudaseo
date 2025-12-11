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
  import { loadSubscriptionData } from './lib/subscription.js';

  let loading = true;
  let configError = null;
  let currentHash = typeof window !== 'undefined' ? window.location.hash : '';
  
  $: isLegalPage = currentHash === '#/privacy' || currentHash === '#/terms';
  $: isPublicPage = isLegalPage || currentHash.startsWith('#/team/accept') || currentHash.startsWith('#/reports/');

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
    
    // Fix double hash issue if present (e.g., #/#/billing -> #/billing)
    if (window.location.hash.startsWith('#/#/')) {
      const fixedHash = window.location.hash.replace('#/#/', '#/');
      window.location.hash = fixedHash;
      console.log('Auto-fixed hash from #/#/ to:', fixedHash);
    }
    
    // Check for auth callback (email confirmation, password reset, etc.)
    // Support auth hashes like "#access_token=..." and "#/reset#access_token=..."
    const hashParts = window.location.hash.split('#').filter(Boolean);
    const lastHash = hashParts.length ? hashParts[hashParts.length - 1] : '';
    const hashParams = new URLSearchParams(lastHash.startsWith('/') ? lastHash.substring(1) : lastHash);
    const accessToken = hashParams.get('access_token');
    
    if (accessToken) {
      // Handle auth callback from email confirmation
      const { data, error } = await supabase.auth.setSession({
        access_token: accessToken,
        refresh_token: hashParams.get('refresh_token') || ''
      });
      
      if (error) {
        console.error('Auth callback error:', error);
      } else {
        // Clear hash from URL
        window.history.replaceState(null, '', window.location.pathname);
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
