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

    // STEP 1: Check if user is already logged in (App.svelte may have already handled tokens)
    const { data: { session: existingSession } } = await supabase.auth.getSession();
    if (existingSession) {
      console.log('✅ User already logged in (session set by App.svelte), redirecting...');
      success = true;
      loading = false;
      setTimeout(() => {
        push('/');
      }, 500);
      return;
    }

    // STEP 2: Check if tokens are in hash (implicit flow - should be handled by App.svelte)
    const fullHash = window.location.hash;
    if (fullHash.includes('access_token=')) {
      console.log('⚠️ Tokens found in hash - App.svelte should handle this');
      console.log('⏳ Waiting for App.svelte to process tokens...');
      
      // Wait a bit for App.svelte to process
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      // Check again if session was set
      const { data: { session: sessionAfterWait } } = await supabase.auth.getSession();
      if (sessionAfterWait) {
        console.log('✅ Session set by App.svelte, redirecting...');
        success = true;
        loading = false;
        setTimeout(() => {
          push('/');
        }, 500);
        return;
      } else {
        console.log('⚠️ Session not set yet, will handle tokens here');
        // Fall through to handle tokens below
      }
    }

    try {
      // Extract parameters from URL
      // PKCE flow can use either:
      // - token_hash + type (older format)
      // - code (newer format)
      const hashParts = window.location.hash.split('?');
      const queryString = hashParts.length > 1 ? hashParts[1] : window.location.search.substring(1);
      const params = new URLSearchParams(queryString);
      
      const tokenHash = params.get('token_hash');
      const code = params.get('code');
      const type = params.get('type');

      console.log('🔍 Token hash:', tokenHash ? 'present' : 'missing');
      console.log('🔍 Code:', code ? 'present' : 'missing');
      console.log('🔍 Type:', type);

      // Mark verification as attempted
      verificationAttempted = true;

      // Clear URL immediately to prevent re-runs
      window.history.replaceState(null, '', window.location.pathname + '#/');

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

      let data, verifyError;

      // Handle implicit flow tokens in hash (if App.svelte didn't handle them)
      const urlHash = window.location.hash;
      if (urlHash.includes('access_token=')) {
        console.log('🔄 Processing access_token from URL hash (implicit flow)...');
        // Extract tokens from URL hash
        // Handle case where hash might be: #/auth/confirm#access_token=... or #access_token=...
        const hashPart = urlHash.includes('#access_token=') 
          ? urlHash.substring(urlHash.indexOf('#access_token=') + 1) // Remove leading #
          : urlHash.substring(urlHash.indexOf('access_token='));
        
        const hashParams = new URLSearchParams(hashPart);
        const accessToken = hashParams.get('access_token');
        const refreshToken = hashParams.get('refresh_token');
        
        if (accessToken && refreshToken) {
          console.log('🔄 Setting session from URL tokens...');
          const { data: sessionData, error: sessionError } = await supabase.auth.setSession({
            access_token: accessToken,
            refresh_token: refreshToken
          });
          
          if (sessionError) {
            console.error('🔴 Failed to set session:', sessionError);
            throw sessionError;
          }
          
          if (sessionData?.session) {
            console.log('✅ Session set from URL tokens');
            data = sessionData;
            verifyError = null;
          } else {
            throw new Error('Session was not created from URL tokens');
          }
        } else {
          throw new Error('Missing access_token or refresh_token in URL');
        }
      }
      // Handle code parameter (PKCE flow - shouldn't happen with implicit flow)
      else if (code) {
        console.log('⚠️ Code parameter found - this shouldn\'t happen with implicit flow');
        console.log('🔄 Attempting code exchange (may fail)...');
        const { data: exchangeData, error: exchangeError } = await supabase.auth.exchangeCodeForSession(code);
        data = exchangeData;
        verifyError = exchangeError;
        
        if (exchangeError) {
          throw new Error('Code exchange failed. Please request a new magic link.');
        }
      }
      // Handle token_hash + type (older PKCE format)
      else if (tokenHash && type) {
        console.log('🔄 Verifying OTP with token hash...');
        // Exchange the token hash for a session (PKCE flow)
        // Per Supabase docs: https://supabase.com/docs/guides/auth/auth-email-passwordless
        // verifyOtp automatically sets the session and triggers SIGNED_IN event
        const { data: verifyData, error: verifyErr } = await supabase.auth.verifyOtp({
          token_hash: tokenHash,
          type: type
        });
        data = verifyData;
        verifyError = verifyErr;
      }
      else {
        throw new Error('Missing code or token_hash+type parameters in URL');
      }

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

