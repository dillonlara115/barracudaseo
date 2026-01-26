import { supabase } from './supabase.js';

export const getApiUrl = () => import.meta.env.VITE_CLOUD_RUN_API_URL || 'http://localhost:8080';

// Track ongoing refresh to prevent concurrent refresh attempts
let refreshPromise = null;
let lastRefreshTime = 0;
const REFRESH_COOLDOWN = 5000; // 5 seconds cooldown between refresh attempts

let isRefreshing = false;
let pendingRequests = [];

async function onTokenRefreshed(token) {
  const requests = [...pendingRequests];
  pendingRequests = [];
  requests.forEach(({ resolve }) => resolve(token));
}

async function onTokenRefreshFailed(error) {
  const requests = [...pendingRequests];
  pendingRequests = [];
  requests.forEach(({ reject }) => reject(error));
}

async function getValidAccessToken() {
  // If refresh is already in progress, wait for it (check this FIRST to prevent race conditions)
  if (isRefreshing && refreshPromise) {
    try {
      const { data: refreshed, error: refreshError } = await refreshPromise;
      if (!refreshError && refreshed.session && refreshed.session.access_token) {
        // Verify the refreshed token is not expired
        const expiresAt = refreshed.session.expires_at;
        const expiresAtMs = expiresAt ? expiresAt * 1000 : 0;
        if (expiresAt && expiresAtMs > Date.now()) {
          return refreshed.session.access_token;
        }
        // Token is still expired, fall through to refresh again
      }
      // Refresh failed, fall through to try again
    } catch (err) {
      // Refresh promise failed, fall through to try again
    }
  }
  
  // Get current session
  const { data: { session }, error: sessionError } = await supabase.auth.getSession();
  
  // If there's an error or no session, try to refresh
  if (sessionError || !session) {
    // If refresh is already in progress, wait for it
    if (isRefreshing && refreshPromise) {
      return new Promise((resolve, reject) => {
        pendingRequests.push({ resolve, reject });
      });
    }

    // Atomic check-and-set: set flag BEFORE creating promise to prevent race conditions
    if (isRefreshing) {
      // Another request beat us, wait for it
      return new Promise((resolve, reject) => {
        pendingRequests.push({ resolve, reject });
      });
    }
    
    isRefreshing = true;
    
    try {
      // Start new refresh - create promise BEFORE any await to prevent race conditions
      refreshPromise = supabase.auth.refreshSession();
      lastRefreshTime = Date.now();
      const { data: refreshed, error: refreshError } = await refreshPromise;
      refreshPromise = null;
      
      if (refreshError || !refreshed.session) {
        throw new Error('Not authenticated. Please sign in again.');
      }
      
      const token = refreshed.session.access_token;
      onTokenRefreshed(token);
      return token;
    } catch (err) {
      onTokenRefreshFailed(err);
      throw err;
    } finally {
      isRefreshing = false;
    }
  }

  // Check if token is expired or about to expire (within 60 seconds)
  const expiresAt = session.expires_at;
  const now = Date.now();
  const expiresAtMs = expiresAt ? expiresAt * 1000 : 0;
  
  // Check if token is already expired (more aggressive check)
  const isExpired = !expiresAt || expiresAtMs < now;
  
  // Refresh if expired or expiring within 60 seconds
  // If already expired, refresh immediately (don't use expired token)
  if (isExpired || expiresAtMs < now + 60000) {
    // If refresh is already in progress, wait for it
    if (isRefreshing && refreshPromise) {
      return new Promise((resolve, reject) => {
        pendingRequests.push({ resolve, reject });
      });
    }

    // Atomic check-and-set: set flag BEFORE creating promise to prevent race conditions
    if (isRefreshing) {
      // Another request beat us, wait for it
      return new Promise((resolve, reject) => {
        pendingRequests.push({ resolve, reject });
      });
    }

    isRefreshing = true;

    try {
      // Start new refresh - create promise BEFORE any await to prevent race conditions
      refreshPromise = supabase.auth.refreshSession();
      lastRefreshTime = Date.now();
      const { data: refreshed, error: refreshError } = await refreshPromise;
      refreshPromise = null;
      
      if (refreshError || !refreshed.session) {
        // Don't use expired token - refresh failed means session is invalid
        // Check if token is actually expired before throwing
        const expiresAt = session.expires_at;
        const expiresAtMs = expiresAt ? expiresAt * 1000 : 0;
        const isExpired = !expiresAt || expiresAtMs < Date.now();
        
        if (isExpired) {
          // Token is expired and refresh failed - user needs to sign in again
          onTokenRefreshFailed(new Error('Session expired. Please sign in again.'));
          throw new Error('Session expired. Please sign in again.');
        }
        
        // Token not expired yet but refresh failed - try using current token
        // This handles transient refresh failures
        const token = session.access_token;
        if (token) {
          // Still notify pending requests about the token (even if refresh failed)
          onTokenRefreshed(token);
          return token;
        }
        
        onTokenRefreshFailed(new Error('Session expired. Please sign in again.'));
        throw new Error('Session expired. Please sign in again.');
      }
      
      const token = refreshed.session.access_token;
      onTokenRefreshed(token);
      return token;
    } catch (err) {
      onTokenRefreshFailed(err);
      throw err;
    } finally {
      isRefreshing = false;
    }
  }
  
  // Token is valid - return it
  return session.access_token;
}

async function authorizedRequest(path, { method = 'GET', body, headers = {}, retryOn401 = true } = {}) {
  try {
    // Ensure we have a valid token before making request
    // This will refresh if needed and wait for any in-progress refresh
    let token = await getValidAccessToken();
    
    // Double-check token is still valid (race condition protection)
    if (!token) {
      throw new Error('No valid access token available');
    }

    const requestHeaders = new Headers(headers);
    requestHeaders.set('Authorization', `Bearer ${token}`);

    let requestBody = body;
    if (body && !(body instanceof FormData) && typeof body === 'object' && !(body instanceof Blob)) {
      requestHeaders.set('Content-Type', 'application/json');
      requestBody = JSON.stringify(body);
    }

    const response = await fetch(`${getApiUrl()}${path}`, {
      method,
      headers: requestHeaders,
      body: requestBody,
    });

    // Handle 401 errors by refreshing token and retrying once
    if (response.status === 401 && retryOn401) {
      // Use centralized refresh function to coordinate with other requests
      try {
        const newToken = await getValidAccessToken();
        
        if (newToken) {
          // Retry request with new token
          requestHeaders.set('Authorization', `Bearer ${newToken}`);
          const retryResponse = await fetch(`${getApiUrl()}${path}`, {
            method,
            headers: requestHeaders,
            body: requestBody,
          });
          
          // If retry still fails with 401, the refresh token itself may be expired
          if (retryResponse.status === 401) {
            console.warn('Retry with refreshed token still returned 401 - refresh token may be expired');
            // Sign out user to force re-authentication
            try {
              await supabase.auth.signOut();
            } catch (signOutErr) {
              console.error('Failed to sign out after token refresh failure:', signOutErr);
            }
          }
          
          return retryResponse;
        } else {
          // No token available - refresh must have failed
          console.warn('Token refresh failed on 401 - no token returned');
          // Sign out user to force re-authentication
          try {
            await supabase.auth.signOut();
          } catch (signOutErr) {
            console.error('Failed to sign out after refresh failure:', signOutErr);
          }
        }
      } catch (refreshErr) {
        console.error('Failed to refresh token on 401:', refreshErr);
        
        // If error indicates expired refresh token, sign out
        if (refreshErr?.message?.includes('refresh_token') || refreshErr?.message?.includes('expired') || refreshErr?.message?.includes('Session expired') || refreshErr?.message?.includes('Not authenticated')) {
          try {
            await supabase.auth.signOut();
          } catch (signOutErr) {
            console.error('Failed to sign out after refresh token expiration:', signOutErr);
          }
        }
      }
    }

    return response;
  } catch (err) {
    console.error('Request failed:', err);
    throw err;
  }
}

async function authorizedJSON(path, options = {}) {
  try {
    const response = await authorizedRequest(path, options);

    if (!response.ok) {
      let message = `Request failed with status ${response.status}`;
      try {
        const errorPayload = await response.json();
        message = errorPayload.error || errorPayload.message || message;
      } catch (_) {
        // Ignore JSON parse errors
      }
      
      // If it's a 401 and we've already retried, the session is likely expired
      if (response.status === 401) {
        message = 'Session expired. Please sign in again.';
        // Sign out to clear invalid session
        try {
          await supabase.auth.signOut();
        } catch (signOutErr) {
          console.error('Failed to sign out after 401:', signOutErr);
        }
      }
      
      throw new Error(message);
    }

    if (response.status === 204) {
      return { data: null, error: null };
    }

    const data = await response.json();
    return { data, error: null };
  } catch (error) {
    // If error indicates session expired, sign out
    if (error.message?.includes('expired') || error.message?.includes('Not authenticated') || error.message?.includes('Session expired')) {
      try {
        await supabase.auth.signOut();
      } catch (signOutErr) {
        console.error('Failed to sign out after auth error:', signOutErr);
      }
    }
    
    console.error('API request failed:', error);
    return { data: null, error };
  }
}

const PROJECTS_CACHE_TTL_MS = 60000;
const CRAWLS_CACHE_TTL_MS = 60000;

const projectsCache = {
  data: null,
  error: null,
  fetchedAt: 0,
  inFlight: null
};

const crawlsCache = new Map();

function isCacheFresh(fetchedAt, ttlMs) {
  return fetchedAt > 0 && Date.now() - fetchedAt < ttlMs;
}

// Fetch user's projects
export async function fetchProjects() {
  try {
    if (projectsCache.inFlight) {
      return await projectsCache.inFlight;
    }

    if (projectsCache.data && isCacheFresh(projectsCache.fetchedAt, PROJECTS_CACHE_TTL_MS)) {
      return { data: projectsCache.data, error: null };
    }

    projectsCache.inFlight = (async () => {
      const { data, error } = await supabase
        .from('projects')
        .select('*')
        .order('created_at', { ascending: false });

      if (error) {
        console.error('Error fetching projects:', error);
        console.error('Error details:', JSON.stringify(error, null, 2));
        return { data: null, error };
      }

      projectsCache.data = data;
      projectsCache.error = null;
      projectsCache.fetchedAt = Date.now();
      return { data, error: null };
    })();

    return await projectsCache.inFlight;
  } catch (error) {
    console.error('Error fetching projects:', error);
    return { data: null, error };
  } finally {
    projectsCache.inFlight = null;
  }
}

// Fetch crawls for a project - use backend API for reliable access
export async function fetchCrawls(projectId) {
  const cacheEntry = crawlsCache.get(projectId);
  if (cacheEntry?.inFlight) {
    return await cacheEntry.inFlight;
  }
  if (cacheEntry?.data && isCacheFresh(cacheEntry.fetchedAt, CRAWLS_CACHE_TTL_MS)) {
    return { data: cacheEntry.data, error: null };
  }

  const inFlight = (async () => {
    try {
      const response = await authorizedRequest(`/api/v1/projects/${projectId}/crawls`);
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ error: `HTTP ${response.status}` }));
        const errorMessage = errorData.error || `Failed to fetch crawls: ${response.status}`;
        const error = new Error(errorMessage);
        error.status = response.status;
        throw error;
      }
      const responseData = await response.json();
      const data = responseData.crawls || (Array.isArray(responseData) ? responseData : []);
      crawlsCache.set(projectId, {
        data,
        error: null,
        fetchedAt: Date.now(),
        inFlight: null
      });
      return { data, error: null };
    } catch (error) {
      console.error('Error fetching crawls:', error);
      return {
        data: null,
        error: error instanceof Error ? error : new Error(String(error))
      };
    }
  })();

  crawlsCache.set(projectId, {
    data: cacheEntry?.data || null,
    error: null,
    fetchedAt: cacheEntry?.fetchedAt || 0,
    inFlight
  });

  const result = await inFlight;
  const updatedEntry = crawlsCache.get(projectId);
  if (updatedEntry) {
    updatedEntry.inFlight = null;
  }
  return result;
}

// Fetch a single crawl by ID - use backend API for real-time updates
export async function fetchCrawl(crawlId) {
  try {
    const response = await authorizedRequest(`/api/v1/crawls/${crawlId}`);
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: `HTTP ${response.status}` }));
      const errorMessage = errorData.error || `Failed to fetch crawl: ${response.status}`;
      const error = new Error(errorMessage);
      error.status = response.status;
      throw error;
    }
    const data = await response.json();
    return { data, error: null };
  } catch (error) {
    console.error('Error fetching crawl:', error);
    // Return error object with message property for consistent handling
    return { 
      data: null, 
      error: error instanceof Error ? error : new Error(String(error))
    };
  }
}

// Fetch page count for a crawl (for progress tracking)
export async function fetchCrawlPageCount(crawlId) {
  try {
    const { count, error } = await supabase
      .from('pages')
      .select('*', { count: 'exact', head: true })
      .eq('crawl_id', crawlId);

    if (error) throw error;
    return { count: count || 0, error: null };
  } catch (error) {
    return { count: 0, error };
  }
}

// Fetch pages for a crawl - use backend API for reliable access
export async function fetchPages(crawlId) {
  try {
    const response = await authorizedRequest(`/api/v1/crawls/${crawlId}/pages`);
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: `HTTP ${response.status}` }));
      const errorMessage = errorData.error || `Failed to fetch pages: ${response.status}`;
      const error = new Error(errorMessage);
      error.status = response.status;
      throw error;
    }
    const data = await response.json();
    return { data, error: null };
  } catch (error) {
    console.error('Error fetching pages:', error);
    return { 
      data: null, 
      error: error instanceof Error ? error : new Error(String(error))
    };
  }
}

// Fetch issues for a crawl - use backend API for reliable access
export async function fetchIssues(crawlId) {
  try {
    const response = await authorizedRequest(`/api/v1/crawls/${crawlId}/issues`);
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: `HTTP ${response.status}` }));
      const errorMessage = errorData.error || `Failed to fetch issues: ${response.status}`;
      const error = new Error(errorMessage);
      error.status = response.status;
      throw error;
    }
    const data = await response.json();
    return { data, error: null };
  } catch (error) {
    console.error('Error fetching issues:', error);
    return { 
      data: null, 
      error: error instanceof Error ? error : new Error(String(error))
    };
  }
}

// Fetch project issue summary (using the view)
export async function fetchProjectIssueSummary(projectId) {
  try {
    const { data, error } = await supabase
      .from('project_issue_summary')
      .select('*')
      .eq('project_id', projectId)
      .single();

    if (error) throw error;
    return { data, error: null };
  } catch (error) {
    return { data: null, error };
  }
}

// Create a new project
export async function createProject(name, domain, settings = {}) {
  try {
    const { data: { user } } = await supabase.auth.getUser();
    if (!user) throw new Error('Not authenticated');

    const { data, error } = await supabase
      .from('projects')
      .insert({
        name,
        domain,
        owner_id: user.id,
        settings
      })
      .select()
      .single();

    if (error) throw error;

    // Also add the owner as a project member
    await supabase
      .from('project_members')
      .insert({
        project_id: data.id,
        user_id: user.id,
        role: 'owner'
      });

    return { data, error: null };
  } catch (error) {
    return { data: null, error };
  }
}

// Update a project
export async function updateProject(projectId, updates) {
  if (!projectId) return { data: null, error: new Error('projectId is required') };
  
  try {
    const response = await authorizedRequest(`/api/v1/projects/${projectId}`, {
      method: 'PUT',
      body: updates
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: `HTTP ${response.status}` }));
      const errorMessage = errorData.error || `Failed to update project: ${response.status}`;
      const error = new Error(errorMessage);
      error.status = response.status;
      throw error;
    }

    const data = await response.json();
    return { data, error: null };
  } catch (error) {
    console.error('Error updating project:', error);
    return { 
      data: null, 
      error: error instanceof Error ? error : new Error(String(error))
    };
  }
}

// Delete a project
export async function deleteProject(projectId) {
  if (!projectId) return { data: null, error: new Error('projectId is required') };
  
  try {
    const response = await authorizedRequest(`/api/v1/projects/${projectId}`, {
      method: 'DELETE'
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: `HTTP ${response.status}` }));
      const errorMessage = errorData.error || `Failed to delete project: ${response.status}`;
      const error = new Error(errorMessage);
      error.status = response.status;
      throw error;
    }

    const data = await response.json();
    return { data, error: null };
  } catch (error) {
    console.error('Error deleting project:', error);
    return { 
      data: null, 
      error: error instanceof Error ? error : new Error(String(error))
    };
  }
}

// Trigger a new crawl for a project
export async function triggerCrawl(projectId, crawlConfig) {
  try {
    const token = await getValidAccessToken();
    const apiUrl = getApiUrl();
    
    const response = await fetch(`${apiUrl}/api/v1/projects/${projectId}/crawl`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(crawlConfig)
    });

    // Handle 401 by refreshing and retrying once
    if (response.status === 401) {
      const { data: refreshed, error: refreshError } = await supabase.auth.refreshSession();
      if (!refreshError && refreshed.session) {
        const retryResponse = await fetch(`${apiUrl}/api/v1/projects/${projectId}/crawl`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${refreshed.session.access_token}`
          },
          body: JSON.stringify(crawlConfig)
        });
        
        if (!retryResponse.ok) {
          const error = await retryResponse.json();
          throw new Error(error.error || 'Failed to trigger crawl');
        }
        const data = await retryResponse.json();
        return { data, error: null };
      }
    }

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to trigger crawl');
    }

    const data = await response.json();
    return { data, error: null };
  } catch (error) {
    console.error('Error triggering crawl:', error);
    return { data: null, error };
  }
}

// Update project settings
export async function updateProjectSettings(projectId, settings) {
  try {
    const { data: { user } } = await supabase.auth.getUser();
    if (!user) throw new Error('Not authenticated');

    // Get current project to merge settings
    const { data: project, error: fetchError } = await supabase
      .from('projects')
      .select('settings')
      .eq('id', projectId)
      .single();

    if (fetchError) throw fetchError;

    // Merge with existing settings
    const updatedSettings = {
      ...(project.settings || {}),
      ...settings
    };

    const { data, error } = await supabase
      .from('projects')
      .update({ settings: updatedSettings })
      .eq('id', projectId)
      .select()
      .single();

    if (error) throw error;
    return { data, error: null };
  } catch (error) {
    return { data: null, error };
  }
}

// Update issue status
export async function updateIssueStatus(issueId, status, notes = null) {
  try {
    const { data: { user } } = await supabase.auth.getUser();
    if (!user) throw new Error('Not authenticated');

    const { data, error } = await supabase
      .from('issues')
      .update({
        status,
        status_updated_at: new Date().toISOString()
      })
      .eq('id', issueId)
      .select()
      .single();

    if (error) throw error;

    // Log status change
    if (data) {
      await supabase
        .from('issue_status_history')
        .insert({
          issue_id: issueId,
          old_status: data.status, // This won't work perfectly - we'd need the old value
          new_status: status,
          changed_by: user.id,
          notes
        });
    }

    return { data, error: null };
  } catch (error) {
    return { data: null, error };
  }
}

export async function fetchProjectGSCConnect(projectId) {
  if (!projectId) return { data: null, error: new Error('projectId is required') };
  return authorizedJSON(`/api/v1/projects/${projectId}/gsc/connect`);
}

export async function fetchProjectGSCStatus(projectId) {
  if (!projectId) return { data: null, error: new Error('projectId is required') };
  return authorizedJSON(`/api/v1/projects/${projectId}/gsc/status`);
}

export async function fetchProjectGSCProperties(projectId) {
  if (!projectId) return { data: null, error: new Error('projectId is required') };
  return authorizedJSON(`/api/v1/projects/${projectId}/gsc/properties`);
}

export async function updateProjectGSCProperty(projectId, propertyUrl, propertyType = null) {
  if (!projectId) return { data: null, error: new Error('projectId is required') };
  if (!propertyUrl) return { data: null, error: new Error('propertyUrl is required') };
  return authorizedJSON(`/api/v1/projects/${projectId}/gsc/property`, {
    method: 'POST',
    body: {
      property_url: propertyUrl,
      property_type: propertyType,
    },
  });
}

export async function triggerProjectGSCSync(projectId, options = {}) {
  if (!projectId) return { data: null, error: new Error('projectId is required') };
  return authorizedJSON(`/api/v1/projects/${projectId}/gsc/trigger-sync`, {
    method: 'POST',
    body: options,
  });
}

export async function fetchProjectGSCDimensions(projectId, type, params = {}) {
  if (!projectId) return { data: null, error: new Error('projectId is required') };
  if (!type) return { data: null, error: new Error('type is required') };

  const searchParams = new URLSearchParams({ type });
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') {
      searchParams.set(key, value.toString());
    }
  });

  return authorizedJSON(`/api/v1/projects/${projectId}/gsc/dimensions?${searchParams.toString()}`);
}

// Fetch link graph for a crawl
export async function fetchCrawlGraph(crawlId) {
  if (!crawlId) return { data: null, error: new Error('crawlId is required') };
  return authorizedJSON(`/api/v1/crawls/${crawlId}/graph`);
}

// Delete a crawl
export async function deleteCrawl(crawlId) {
  if (!crawlId) return { data: null, error: new Error('crawlId is required') };
  try {
    const response = await authorizedRequest(`/api/v1/crawls/${crawlId}`, {
      method: 'DELETE'
    });
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: `HTTP ${response.status}` }));
      const errorMessage = errorData.error || `Failed to delete crawl: ${response.status}`;
      const error = new Error(errorMessage);
      error.status = response.status;
      throw error;
    }
    const data = await response.json();
    return { data, error: null };
  } catch (error) {
    console.error('Error deleting crawl:', error);
    return { 
      data: null, 
      error: error instanceof Error ? error : new Error(String(error))
    };
  }
}

export async function disconnectProjectGSC(projectId) {
  if (!projectId) return { data: null, error: new Error('projectId is required') };
  return authorizedJSON(`/api/v1/projects/${projectId}/gsc/disconnect`, {
    method: 'POST'
  });
}

// GA4 API functions
export async function fetchProjectGA4Connect(projectId) {
  if (!projectId) return { data: null, error: new Error('projectId is required') };
  return authorizedJSON(`/api/v1/projects/${projectId}/ga4/connect`);
}

export async function fetchProjectGA4Status(projectId) {
  if (!projectId) return { data: null, error: new Error('projectId is required') };
  return authorizedJSON(`/api/v1/projects/${projectId}/ga4/status`);
}

export async function fetchProjectGA4Properties(projectId) {
  if (!projectId) return { data: null, error: new Error('projectId is required') };
  return authorizedJSON(`/api/v1/projects/${projectId}/ga4/properties`);
}

export async function updateProjectGA4Property(projectId, propertyId, propertyName = null) {
  if (!projectId) return { data: null, error: new Error('projectId is required') };
  if (!propertyId) return { data: null, error: new Error('propertyId is required') };
  return authorizedJSON(`/api/v1/projects/${projectId}/ga4/property`, {
    method: 'POST',
    body: {
      property_id: propertyId,
      property_name: propertyName,
    },
  });
}

export async function triggerProjectGA4Sync(projectId, options = {}) {
  if (!projectId) return { data: null, error: new Error('projectId is required') };
  return authorizedJSON(`/api/v1/projects/${projectId}/ga4/trigger-sync`, {
    method: 'POST',
    body: options,
  });
}

export async function disconnectProjectGA4(projectId) {
  if (!projectId) return { data: null, error: new Error('projectId is required') };
  return authorizedJSON(`/api/v1/projects/${projectId}/ga4/disconnect`, {
    method: 'POST'
  });
}

// AI-related API functions

// Save OpenAI API key
export async function saveOpenAIKey(openaiApiKey) {
  return authorizedJSON('/api/v1/integrations/openai-key', {
    method: 'POST',
    body: { openai_api_key: openaiApiKey }
  });
}

// Get OpenAI API key status (doesn't return the actual key)
export async function getOpenAIKeyStatus() {
  return authorizedJSON('/api/v1/integrations/openai-key', {
    method: 'GET'
  });
}

// Disconnect/remove OpenAI API key
export async function disconnectOpenAIKey() {
  return authorizedJSON('/api/v1/integrations/openai-key', {
    method: 'DELETE'
  });
}

// Generate AI insight for an issue
export async function generateIssueInsight(issueId, crawlId) {
  return authorizedJSON('/api/v1/ai/issue-insight', {
    method: 'POST',
    body: {
      issue_id: String(issueId), // Convert to string as API expects string type
      crawl_id: String(crawlId || '') // Ensure crawlId is also a string
    }
  });
}

// Generate AI summary for a crawl
// Get existing crawl summary
export async function getCrawlSummary(crawlId) {
  return authorizedJSON(`/api/v1/ai/crawl-summary?crawl_id=${encodeURIComponent(String(crawlId))}`, {
    method: 'GET'
  });
}

// Generate or regenerate crawl summary
export async function generateCrawlSummary(crawlId, forceRefresh = false) {
  return authorizedJSON('/api/v1/ai/crawl-summary', {
    method: 'POST',
    body: {
      crawl_id: String(crawlId), // Ensure crawlId is a string
      force_refresh: Boolean(forceRefresh) // Explicitly convert to boolean
    }
  });
}

// Delete crawl summary
export async function deleteCrawlSummary(crawlId) {
  return authorizedJSON(`/api/v1/ai/crawl-summary?crawl_id=${encodeURIComponent(String(crawlId))}`, {
    method: 'DELETE'
  });
}

// Create a public report for a crawl
export async function createPublicReport(crawlId, options = {}) {
  return authorizedJSON('/api/v1/reports/public', {
    method: 'POST',
    body: {
      crawl_id: String(crawlId),
      title: options.title || '',
      description: options.description || '',
      password: options.password || '',
      expires_in_days: options.expiresInDays || null,
      settings: options.settings || {}
    }
  });
}

// List all public reports for the current user
export async function listPublicReports(projectId = null) {
  const url = projectId 
    ? `/api/v1/reports/public?project_id=${projectId}`
    : '/api/v1/reports/public';
  return authorizedJSON(url, {
    method: 'GET'
  });
}

// Delete a public report
export async function deletePublicReport(reportId) {
  return authorizedJSON(`/api/v1/reports/public/${reportId}`, {
    method: 'DELETE'
  });
}

// View a public report (no auth required)
export async function viewPublicReport(accessToken, password = null) {
  const apiUrl = getApiUrl();
  const url = `${apiUrl}/api/public/reports/${accessToken}`;
  
  const options = {
    method: password ? 'POST' : 'GET',
    headers: {
      'Content-Type': 'application/json'
    }
  };
  
  if (password) {
    options.body = JSON.stringify({ password });
  }
  
  try {
    const response = await fetch(url, options);
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: `HTTP ${response.status}` }));
      throw new Error(errorData.error || `Failed to fetch report: ${response.status}`);
    }
    
    const data = await response.json();
    return { data, error: null };
  } catch (error) {
    console.error('Error viewing public report:', error);
    return { data: null, error };
  }
}

// Keyword management functions

// Create a keyword
export async function createKeyword(projectId, keywordData) {
  if (!projectId) return { data: null, error: new Error('projectId is required') };
  if (!keywordData) return { data: null, error: new Error('keywordData is required') };
  
  try {
    // Create body with project_id from parameter, and spread keywordData (which may or may not have project_id)
    const { project_id, ...restData } = keywordData;
    const body = {
      project_id: projectId,
      ...restData
    };
    
    const { data, error } = await authorizedJSON('/api/v1/keywords', {
      method: 'POST',
      body
    });
    return { data, error };
  } catch (error) {
    console.error('Error creating keyword:', error);
    return { data: null, error };
  }
}

// List keywords for a project
export async function listKeywords(projectId, filters = {}) {
  try {
    const params = new URLSearchParams({ project_id: projectId });
    if (filters.device) params.append('device', filters.device);
    if (filters.location) params.append('location', filters.location);
    if (filters.tag) params.append('tag', filters.tag);
    
    const { data, error } = await authorizedJSON(`/api/v1/keywords?${params.toString()}`);
    return { data, error };
  } catch (error) {
    console.error('Error listing keywords:', error);
    return { data: null, error };
  }
}

// Get a single keyword
export async function getKeyword(keywordId) {
  try {
    const { data, error } = await authorizedJSON(`/api/v1/keywords/${keywordId}`);
    return { data, error };
  } catch (error) {
    console.error('Error fetching keyword:', error);
    return { data: null, error };
  }
}

// Update a keyword
export async function updateKeyword(keywordId, updates) {
  try {
    const { data, error } = await authorizedJSON(`/api/v1/keywords/${keywordId}`, {
      method: 'PUT',
      body: updates
    });
    return { data, error };
  } catch (error) {
    console.error('Error updating keyword:', error);
    return { data: null, error };
  }
}

// Delete a keyword
export async function deleteKeyword(keywordId) {
  try {
    const response = await authorizedRequest(`/api/v1/keywords/${keywordId}`, {
      method: 'DELETE'
    });
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: `HTTP ${response.status}` }));
      throw new Error(errorData.error || `Failed to delete keyword: ${response.status}`);
    }
    
    return { data: null, error: null };
  } catch (error) {
    console.error('Error deleting keyword:', error);
    return { data: null, error };
  }
}

// Check keyword ranking (trigger rank check)
export async function checkKeyword(keywordId) {
  try {
    const { data, error } = await authorizedJSON(`/api/v1/keywords/${keywordId}/check`, {
      method: 'POST'
    });
    return { data, error };
  } catch (error) {
    console.error('Error checking keyword:', error);
    return { data: null, error };
  }
}

// Get keyword snapshots (historical rank data)
export async function getKeywordSnapshots(keywordId, limit = 30) {
  try {
    const { data, error } = await authorizedJSON(`/api/v1/keywords/${keywordId}/snapshots?limit=${limit}`);
    return { data, error };
  } catch (error) {
    console.error('Error fetching keyword snapshots:', error);
    return { data: null, error };
  }
}

// Get project keyword metrics
export async function getProjectKeywordMetrics(projectId) {
  try {
    const { data, error } = await authorizedJSON(`/api/v1/projects/${projectId}/keyword-metrics`);
    return { data, error };
  } catch (error) {
    console.error('Error fetching project keyword metrics:', error);
    return { data: null, error };
  }
}

export async function fetchProjectKeywordUsage(projectId) {
  try {
    const { data, error } = await authorizedJSON(`/api/v1/projects/${projectId}/keyword-usage`);
    return { data, error };
  } catch (error) {
    console.error('Error fetching project keyword usage:', error);
    return { data: null, error };
  }
}

export async function fetchProjectImpactFirst(projectId) {
  try {
    const { data, error } = await authorizedJSON(`/api/v1/projects/${projectId}/impact-first`);
    return { data, error };
  } catch (error) {
    console.error('Error fetching impact-first view:', error);
    return { data: null, error };
  }
}

// Discover keywords for a domain/URL
export async function discoverKeywords(projectId, discoveryData) {
  if (!projectId) return { data: null, error: new Error('projectId is required') };
  if (!discoveryData) return { data: null, error: new Error('discoveryData is required') };
  
  try {
    const { data, error } = await authorizedJSON(`/api/v1/projects/${projectId}/discover-keywords`, {
      method: 'POST',
      body: {
        ...discoveryData
      }
    });
    
    return { data, error };
  } catch (err) {
    return { data: null, error: err };
  }
}

// Billing functions - use authorizedRequest wrapper for coordinated auth

/**
 * Fetch billing summary (profile, subscription, team info)
 */
export async function fetchBillingSummary() {
  try {
    const { data, error } = await authorizedJSON(`/api/v1/billing/summary`);
    return { data, error };
  } catch (err) {
    return { data: null, error: err };
  }
}

/**
 * Create Stripe checkout session
 */
export async function createBillingCheckout(priceId) {
  if (!priceId) return { data: null, error: new Error('priceId is required') };
  
  try {
    const { data, error } = await authorizedJSON(`/api/v1/billing/checkout`, {
      method: 'POST',
      body: { price_id: priceId }
    });
    return { data, error };
  } catch (err) {
    return { data: null, error: err };
  }
}

/**
 * Create Stripe customer portal session
 */
export async function createBillingPortal() {
  try {
    const { data, error } = await authorizedJSON(`/api/v1/billing/portal`, {
      method: 'POST'
    });
    return { data, error };
  } catch (err) {
    return { data: null, error: err };
  }
}

/**
 * Redeem a promo code
 */
export async function redeemPromoCode(code, teamSize = 1) {
  if (!code) return { data: null, error: new Error('code is required') };
  
  try {
    const { data, error } = await authorizedJSON(`/api/v1/billing/redeem`, {
      method: 'POST',
      body: { code, team_size: teamSize }
    });
    return { data, error };
  } catch (err) {
    return { data: null, error: err };
  }
}
