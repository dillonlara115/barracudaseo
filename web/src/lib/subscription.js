import { writable } from 'svelte/store';
import { getApiUrl } from './data.js';
import { supabase } from './supabase.js';

// Store for user subscription/profile data
export const userProfile = writable(null);
export const userSubscription = writable(null);

/**
 * Load user subscription data from the API
 */
export async function loadSubscriptionData() {
  try {
    const { data: { session }, error: sessionError } = await supabase.auth.getSession();
    
    if (sessionError) {
      console.error('Session error:', sessionError);
      userProfile.set(null);
      userSubscription.set(null);
      return;
    }
    
    if (!session) {
      console.log('No session available');
      userProfile.set(null);
      userSubscription.set(null);
      return;
    }

    // Refresh session if needed
    let currentSession = session;
    const expiresAt = session.expires_at;
    if (expiresAt && expiresAt * 1000 < Date.now() + 60000) {
      console.log('Session expiring soon, refreshing...');
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
        console.error('Unauthorized - token may be expired. Attempting to refresh session...');
        // Try refreshing the session once more
        const { data: refreshed, error: refreshError } = await supabase.auth.refreshSession();
        if (!refreshError && refreshed.session) {
          // Retry with refreshed token
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
      }
      console.error('Failed to load subscription data:', response.status);
      return;
    }

    const data = await response.json();
    userProfile.set(data.profile || null);
    userSubscription.set(data.subscription || null);
  } catch (error) {
    console.error('Error loading subscription data:', error);
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

