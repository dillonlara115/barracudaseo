import { writable } from 'svelte/store';
import { supabase } from './supabase.js';

// Auth state store
export const user = writable(null);
export const session = writable(null);
export const loading = writable(true);
export const authEvent = writable(null); // Track the last auth event type

// Initialize auth state
export async function initAuth() {
  try {
    // Get initial session
    const { data: { session: initialSession } } = await supabase.auth.getSession();
    session.set(initialSession);
    user.set(initialSession?.user ?? null);

    // Listen for auth changes
    supabase.auth.onAuthStateChange((event, newSession) => {
      authEvent.set(event); // Track the event type
      session.set(newSession);
      user.set(newSession?.user ?? null);
      loading.set(false);
    });

    loading.set(false);
  } catch (error) {
    console.error('Error initializing auth:', error);
    loading.set(false);
  }
}

// Sign up with magic link (passwordless)
export async function signUpWithMagicLink(email, displayName = null) {
  try {
    // For PKCE flow, redirect to /auth/confirm endpoint (not hash route)
    // App.svelte will convert /auth/confirm to /#/auth/confirm for SPA routing
    const redirectTo = typeof window !== 'undefined' 
      ? `${window.location.origin}/auth/confirm`
      : undefined;

    console.log('🔍 Requesting magic link signup for:', email);
    console.log('🔍 Redirect URL:', redirectTo);
    console.log('🔍 Display name:', displayName);

    const { data, error } = await supabase.auth.signInWithOtp({
      email,
      options: {
        data: {
          display_name: displayName
        },
        emailRedirectTo: redirectTo,
        shouldCreateUser: true
      }
    });

    if (error) {
      console.error('🔴 Magic link signup error:', error);
      throw error;
    }

    console.log('✅ Magic link signup sent successfully');
    return { data, error: null };
  } catch (error) {
    console.error('🔴 Failed to send magic link signup:', error);
    return { data: null, error };
  }
}

// Legacy sign up with password (kept for existing users)
export async function signUp(email, password, displayName = null) {
  try {
    // Don't include hash in redirect URL
    const redirectTo = typeof window !== 'undefined' 
      ? window.location.origin
      : undefined;

    const { data, error } = await supabase.auth.signUp({
      email,
      password,
      options: {
        data: {
          display_name: displayName
        },
        emailRedirectTo: redirectTo
      }
    });

    if (error) throw error;

    // Create profile if user was created
    if (data.user) {
      const { error: profileError } = await supabase
        .from('profiles')
        .insert({
          id: data.user.id,
          display_name: displayName || email.split('@')[0]
        });

      if (profileError) {
        console.error('Error creating profile:', profileError);
        // Don't throw - profile might already exist
      }
    }

    return { data, error: null };
  } catch (error) {
    return { data: null, error };
  }
}

// Sign in with magic link (primary method)
export async function signInWithMagicLink(email) {
  try {
    // For PKCE flow, redirect to /auth/confirm endpoint (not hash route)
    // App.svelte will convert /auth/confirm to /#/auth/confirm for SPA routing
    const redirectTo = typeof window !== 'undefined' 
      ? `${window.location.origin}/auth/confirm`
      : undefined;

    console.log('🔍 Requesting magic link for:', email);
    console.log('🔍 Redirect URL:', redirectTo);
    console.log('🔍 Current origin:', window?.location?.origin);

    const { data, error } = await supabase.auth.signInWithOtp({
      email,
      options: {
        emailRedirectTo: redirectTo,
        shouldCreateUser: false // Don't create user on sign-in, only on sign-up
      }
    });

    if (error) {
      console.error('🔴 Magic link error:', error);
      throw error;
    }
    
    console.log('✅ Magic link sent successfully');
    return { data, error: null };
  } catch (error) {
    console.error('🔴 Failed to send magic link:', error);
    return { data: null, error };
  }
}

// Sign in with password (legacy/optional method)
export async function signIn(email, password) {
  try {
    const { data, error } = await supabase.auth.signInWithPassword({
      email,
      password
    });

    if (error) throw error;
    return { data, error: null };
  } catch (error) {
    return { data: null, error };
  }
}

// Sign out
export async function signOut() {
  try {
    const { error } = await supabase.auth.signOut();
    if (error) throw error;
    return { error: null };
  } catch (error) {
    return { error };
  }
}

// Get current user
export async function getCurrentUser() {
  const { data: { user } } = await supabase.auth.getUser();
  return user;
}

// Update user password (optional for users who want to set one)
export async function updatePassword(newPassword) {
  try {
    const { data, error } = await supabase.auth.updateUser({
      password: newPassword
    });

    if (error) throw error;
    return { data, error: null };
  } catch (error) {
    return { data: null, error };
  }
}

// Update user email
// Note: Supabase will send confirmation emails to both old and new email addresses
// (if double_confirm_changes is enabled in Supabase config)
export async function updateEmail(newEmail, password = null) {
  try {
    // Verify password by attempting to sign in (for security)
    // This ensures the user knows their password before changing email
    if (password) {
      const { data: { user: currentUser } } = await supabase.auth.getUser();
      if (!currentUser || !currentUser.email) {
        throw new Error('No user found');
      }

      // Verify password by attempting to sign in
      const { error: signInError } = await supabase.auth.signInWithPassword({
        email: currentUser.email,
        password: password
      });

      if (signInError) {
        throw new Error('Password incorrect. Please verify your password.');
      }
    }

    // Update email - Supabase will send confirmation emails to both addresses
    // The email change will only be applied after both confirmations
    const { data, error } = await supabase.auth.updateUser({
      email: newEmail
    });

    if (error) throw error;
    return { data, error: null };
  } catch (error) {
    return { data: null, error };
  }
}

