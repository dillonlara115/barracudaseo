<script>
  import { onMount } from 'svelte';
  import Router, { link, push } from 'svelte-spa-router';
  import { initAuth, user } from './lib/auth.js';
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
      '/project/:projectId/impact-first': ImpactFirstView,
    '/team/accept': TeamAccept,
    '/reports/:token': PublicReportView,
    '/auth': Auth, // Auth route for when user is authenticated but needs to redirect
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

  onMount(async () => {
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
    const hashParams = new URLSearchParams(window.location.hash.substring(1));
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
      if (!currentUser) {
        // Don't redirect if on public pages (legal pages or invite acceptance)
        if (!isPublicPage) {
          push('/');
        }
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
