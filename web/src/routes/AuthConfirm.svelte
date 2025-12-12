<script>
  import { onMount } from 'svelte';
  import { supabase } from '../lib/supabase.js';
  import { push } from 'svelte-spa-router';

  let loading = true;
  let error = null;
  let success = false;
  let verificationAttempted = false; // Guard to prevent multiple verifications
  
  // Temporary helper: decode JWT payload for debugging project/issuer mismatch
  function decodeJwtClaims(token) {
    try {
      const payload = token.split('.')[1];
      const padded = payload.padEnd(payload.length + (4 - (payload.length % 4)) % 4, '=');
      const json = atob(padded.replace(/-/g, '+').replace(/_/g, '/'));
      return JSON.parse(json);
    } catch (err) {
      console.warn('Failed to decode JWT payload:', err);
      return null;
    }
  }

  onMount(async () => {
    // Prevent multiple verification attempts
    if (verificationAttempted) {
      console.log('⚠️ Verification already attempted, skipping...');
      return;
    }

    console.log('🔍 Auth confirm page loaded');
    console.log('🔍 Full URL:', window.location.href);
    console.log('🔍 Hash:', window.location.hash);
    console.log('🔍 Search:', window.location.search);

    // Check if user is already logged in
    const { data: { session: existingSession } } = await supabase.auth.getSession();
    if (existingSession) {
      console.log('✅ User already logged in, redirecting...');
      push('/');
      return;
    }

    try {
      // Extract token_hash and type from URL
      // URL format: /#/auth/confirm?token_hash=XXXXX&type=email
      const hashParts = window.location.hash.split('?');
      const queryString = hashParts.length > 1 ? hashParts[1] : window.location.search.substring(1);
      const params = new URLSearchParams(queryString);
      
      const tokenHash = params.get('token_hash');
      const type = params.get('type');

      console.log('🔍 Token hash:', tokenHash ? 'present' : 'missing');
      console.log('🔍 Type:', type);

      if (!tokenHash || !type) {
        throw new Error('Missing token_hash or type parameter in URL');
      }

      // Mark verification as attempted
      verificationAttempted = true;

      // Clear URL immediately to prevent re-runs
      window.history.replaceState(null, '', window.location.pathname + '#/');

      console.log('🔄 Verifying OTP with token hash...');

      // Wait for SIGNED_IN event to ensure session is fully established
      const signedInPromise = new Promise((resolve, reject) => {
        const timeout = setTimeout(() => {
          reject(new Error('Timeout waiting for SIGNED_IN event'));
        }, 5000);

        const { data: { subscription } } = supabase.auth.onAuthStateChange((event, session) => {
          console.log('🔔 Auth state change:', event, session ? 'session present' : 'no session');
          if (event === 'SIGNED_IN' && session) {
            clearTimeout(timeout);
            subscription.unsubscribe();
            resolve(session);
          }
        });
      });

      // Exchange the token hash for a session (PKCE flow)
      // Per Supabase docs: https://supabase.com/docs/guides/auth/auth-email-passwordless
      // verifyOtp automatically sets the session and triggers SIGNED_IN event
      const { data, error: verifyError } = await supabase.auth.verifyOtp({
        token_hash: tokenHash,
        type: type
      });

      if (verifyError) {
        console.error('🔴 OTP verification failed:', verifyError);
        throw verifyError;
      }

      console.log('✅ OTP verified successfully');
      console.log('🔍 Response session:', data?.session ? 'present' : 'missing');
      console.log('🔍 Response user:', data?.user ? 'present' : 'missing');
      
      // Debug claims before any downstream calls/sign-outs can clear session
      if (data?.session?.access_token) {
        const initialClaims = decodeJwtClaims(data.session.access_token);
        if (initialClaims) {
          console.log('🔬 Access token claims (verifyOtp response):', {
            iss: initialClaims.iss,
            aud: initialClaims.aud,
            sub: initialClaims.sub,
            exp: initialClaims.exp,
            projectRef: (() => {
              try {
                return new URL(initialClaims.iss).hostname.split('.')[0];
              } catch {
                return initialClaims.iss;
              }
            })()
          });
        }
      }

      // Explicitly set the session if provided in response
      // This ensures the session is properly stored and tokens are valid
      if (data?.session) {
        console.log('🔄 Explicitly setting session from verifyOtp response...');
        console.log('🔍 Access token length:', data.session.access_token?.length || 0);
        console.log('🔍 Refresh token length:', data.session.refresh_token?.length || 0);
        
        const { data: sessionData, error: setSessionError } = await supabase.auth.setSession({
          access_token: data.session.access_token,
          refresh_token: data.session.refresh_token
        });
        
        if (setSessionError) {
          console.error('🔴 Failed to set session:', setSessionError);
          throw setSessionError;
        }
        
        if (!sessionData.session) {
          throw new Error('Session was not set after setSession call');
        }
        
        console.log('✅ Session explicitly set and verified');
      } else {
        console.warn('⚠️ No session in verifyOtp response, relying on automatic session setting');
      }

      // Wait for SIGNED_IN event to fire (ensures session is fully established)
      try {
        const confirmedSession = await signedInPromise;
        console.log('✅ SIGNED_IN event received, session confirmed:', confirmedSession ? 'present' : 'missing');
      } catch (err) {
        console.warn('⚠️ SIGNED_IN event timeout, checking session directly...');
        // Fallback: check session directly
        const { data: { session: fallbackSession } } = await supabase.auth.getSession();
        if (!fallbackSession) {
          throw new Error('Session was not established after verification');
        }
        console.log('✅ Session confirmed via getSession() fallback');
      }

      // Additional wait to ensure all stores are updated and API calls use new session
      await new Promise(resolve => setTimeout(resolve, 500));

      // Final verification that session is valid and persisted
      const { data: { session: finalSession }, error: sessionCheckError } = await supabase.auth.getSession();
      if (sessionCheckError) {
        console.error('🔴 Session check error:', sessionCheckError);
        throw sessionCheckError;
      }
      if (!finalSession) {
        throw new Error('Session was not persisted');
      }
      
      // Debug which Supabase project issued the token
      const accessClaims = decodeJwtClaims(finalSession.access_token || '');
      if (accessClaims) {
        console.log('🔬 Access token claims (temporary debug):', {
          iss: accessClaims.iss,
          aud: accessClaims.aud,
          sub: accessClaims.sub,
          exp: accessClaims.exp,
          projectRef: (() => {
            try {
              return new URL(accessClaims.iss).hostname.split('.')[0];
            } catch {
              return accessClaims.iss;
            }
          })()
        });
      }

      // Verify session is in localStorage (Supabase stores it there)
      try {
        const supabaseAuthKey = Object.keys(localStorage).find(
          (key) => key.startsWith('sb-') && key.endsWith('-auth-token')
        );
        const storedAuth = supabaseAuthKey ? localStorage.getItem(supabaseAuthKey) : null;
        if (!storedAuth) {
          console.warn('⚠️ Session not found in localStorage, but getSession() returned session');
          // This is okay - Supabase might use a different key format
        } else {
          console.log('✅ Session found in localStorage');
        }
      } catch (err) {
        console.warn('⚠️ Could not inspect localStorage for session:', err);
      }

      console.log('✅ Session fully established and ready');
      console.log('🔍 Final session user:', finalSession.user?.email || 'no email');
      console.log('🔍 Session expires at:', finalSession.expires_at ? new Date(finalSession.expires_at * 1000).toISOString() : 'no expiry');
      
      // Set success state
      success = true;
      loading = false;

      // Redirect using window.location.hash for reliable hash-based routing
      // This ensures the SPA router picks up the change
      console.log('🔄 Redirecting to dashboard...');
      window.location.hash = '#/';
      
      // Small delay to ensure hash change is processed, then verify redirect worked
      setTimeout(() => {
        if (window.location.hash === '#/') {
          console.log('✅ Redirect successful - hash is now #/');
        } else {
          console.warn('⚠️ Hash not updated, trying push() fallback');
          push('/');
        }
      }, 100);

    } catch (err) {
      console.error('🔴 Auth confirm error:', err);
      error = err.message || 'Failed to verify login link. Please try again.';
      loading = false;
      verificationAttempted = false; // Allow retry on error
    }
  });
</script>

<div class="min-h-screen flex items-center justify-center bg-base-200">
  <div class="card w-full max-w-md bg-base-100 shadow-xl">
    <div class="card-body space-y-4">
      {#if loading}
        <div class="text-center py-8">
          <div class="flex justify-center mb-4">
            <span class="loading loading-spinner loading-lg"></span>
          </div>
          <h2 class="text-2xl font-bold mb-2">Signing you in...</h2>
          <p class="text-base-content/70">Please wait while we verify your login link.</p>
        </div>
      {:else if success}
        <div class="text-center py-8">
          <svg class="w-16 h-16 mx-auto mb-4 text-success" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <h2 class="text-2xl font-bold mb-2 text-success">Successfully signed in!</h2>
          <p class="text-base-content/70">Redirecting you to the dashboard...</p>
        </div>
      {:else if error}
        <div class="text-center py-8">
          <svg class="w-16 h-16 mx-auto mb-4 text-error" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <h2 class="text-2xl font-bold mb-2 text-error">Verification Failed</h2>
          <div class="alert alert-error mb-4">
            <span>{error}</span>
          </div>
          <p class="text-sm text-base-content/70 mb-4">
            Your login link may have expired or already been used.
          </p>
          <button class="btn btn-primary w-full" on:click={() => push('/')}>
            Back to Login
          </button>
        </div>
      {/if}
    </div>
  </div>
</div>
