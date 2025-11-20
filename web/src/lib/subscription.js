import { writable } from 'svelte/store';
import { getApiUrl } from './data.js';
import { supabase } from './supabase.js';

// Store for user subscription/profile data
export const userProfile = writable(null);
export const userSubscription = writable(null);

// Prevent multiple simultaneous calls
let isLoading = false;
let lastLoadTime = 0;
const LOAD_DEBOUNCE_MS = 1000; // Don't load more than once per second

/**
 * Load user subscription data from the API
 */
export async function loadSubscriptionData() {
  // Prevent multiple simultaneous calls
  if (isLoading) {
    return;
  }
  
  // Debounce rapid calls
  const now = Date.now();
  if (now - lastLoadTime < LOAD_DEBOUNCE_MS) {
    return;
  }
  
  isLoading = true;
  lastLoadTime = now;
  
  try {
    const { data: { session }, error: sessionError } = await supabase.auth.getSession();
    
    if (sessionError) {
      console.error('Session error:', sessionError);
      userProfile.set(null);
      userSubscription.set(null);
      return;
    }
    
    if (!session) {
      // No session - user not authenticated, set defaults
      userProfile.set(null);
      userSubscription.set(null);
      return;
    }

    // Refresh session if needed
    let currentSession = session;
    const expiresAt = session.expires_at;
    if (expiresAt && expiresAt * 1000 < Date.now() + 60000) {
      const { data: refreshed, error: refreshError } = await supabase.auth.refreshSession();
      if (refreshError) {
        console.error('Failed to refresh session:', refreshError);
        userProfile.set(null);
        userSubscription.set(null);
        return;
      }
      if (refreshed.session) {
        currentSession = refreshed.session;
      }
    }

    const apiUrl = getApiUrl();
    const response = await fetch(`${apiUrl}/api/v1/billing/summary`, {
      headers: {
        'Authorization': `Bearer ${currentSession.access_token}`,
      },
    });

    if (!response.ok) {
      if (response.status === 401) {
        // 401 means unauthorized - try refreshing once, then give up
        const { data: refreshed, error: refreshError } = await supabase.auth.refreshSession();
        if (!refreshError && refreshed.session) {
          // Retry with refreshed token (only once)
          const retryResponse = await fetch(`${apiUrl}/api/v1/billing/summary`, {
            headers: {
              'Authorization': `Bearer ${refreshed.session.access_token}`,
            },
          });
          if (retryResponse.ok) {
            const data = await retryResponse.json();
            userProfile.set(data.profile || null);
            userSubscription.set(data.subscription || null);
            return;
          }
        }
        // If refresh failed or retry failed, user is not authenticated
        // Set defaults and don't retry again
        userProfile.set(null);
        userSubscription.set(null);
        return;
      }
      // For other errors, log but don't retry
      console.error('Failed to load subscription data:', response.status);
      return;
    }

    const data = await response.json();
    userProfile.set(data.profile || null);
    userSubscription.set(data.subscription || null);
  } catch (error) {
    console.error('Error loading subscription data:', error);
  } finally {
    isLoading = false;
  }
}

/**
 * Check if user has access to a feature based on subscription tier
 * @param {string} feature - Feature name ('integrations', 'recommendations', etc.)
 * @param {object} profile - Optional profile object, otherwise uses store value
 */
export function hasFeatureAccess(feature, profile = null) {
  let subscriptionTier = 'free';
  
  if (profile) {
    subscriptionTier = getSubscriptionTier(profile);
  } else {
    // Get from store synchronously (note: this requires store to be subscribed)
    let currentProfile = null;
    const unsubscribe = userProfile.subscribe(value => {
      currentProfile = value;
    });
    unsubscribe();
    subscriptionTier = getSubscriptionTier(currentProfile);
  }
  
  switch (feature) {
    case 'integrations':
      return subscriptionTier === 'pro' || subscriptionTier === 'team';
    case 'recommendations':
      return subscriptionTier === 'pro' || subscriptionTier === 'team';
    case 'team_collaboration':
      return subscriptionTier === 'pro' || subscriptionTier === 'team';
    case 'priority_support':
      return subscriptionTier === 'pro' || subscriptionTier === 'team';
    default:
      return false;
  }
}

/**
 * Get subscription tier from profile
 */
export function getSubscriptionTier(profile) {
  return profile?.subscription_tier || 'free';
}

/**
 * Check if user is on free plan
 */
export function isFreeUser(profile) {
  return getSubscriptionTier(profile) === 'free';
}

/**
 * Check if user is on pro or team plan
 */
export function isProOrTeam(profile) {
  const tier = getSubscriptionTier(profile);
  return tier === 'pro' || tier === 'team';
}

