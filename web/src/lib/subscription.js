import { writable } from 'svelte/store';
import { fetchBillingSummary } from './data.js';

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
    // Use centralized billing API function that handles auth coordination
    // This will automatically handle session checks and token refresh
    const { data, error } = await fetchBillingSummary();
    
    if (error) {
      // If 401, user is not authenticated - set defaults
      if (error.message?.includes('401') || error.message?.includes('Unauthorized')) {
        userProfile.set(null);
        userSubscription.set(null);
        return;
      }
      // For other errors, log but don't retry
      console.error('Failed to load subscription data:', error);
      return;
    }

    if (data) {
      userProfile.set(data.profile || null);
      userSubscription.set(data.subscription || null);
    }
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
    case 'ai_insights':
    case 'rank_tracker':
    case 'discover_keywords':
    case 'public_reports':
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

