import { writable } from 'svelte/store';
import { supabase } from './supabase.js';

// Auth state store
export const user = writable(null);
export const session = writable(null);
export const loading = writable(true);

// Initialize auth state
export async function initAuth() {
  try {
    // Get initial session
    const { data: { session: initialSession } } = await supabase.auth.getSession();
    session.set(initialSession);
    user.set(initialSession?.user ?? null);

    // Listen for auth changes
    supabase.auth.onAuthStateChange((_event, newSession) => {
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

// Sign up
export async function signUp(email, password, displayName = null) {
  try {
    // Get the current origin (works for both localhost and production)
    const redirectTo = typeof window !== 'undefined' 
      ? `${window.location.origin}/auth/callback`
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

// Sign in
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

