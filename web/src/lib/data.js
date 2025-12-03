import { supabase } from './supabase.js';

export const getApiUrl = () => import.meta.env.VITE_CLOUD_RUN_API_URL || 'http://localhost:8080';

async function getValidAccessToken() {
  // Get current session
  const { data: { session }, error: sessionError } = await supabase.auth.getSession();
  
  // If there's an error or no session, try to refresh
  if (sessionError || !session) {
    const { data: refreshed, error: refreshError } = await supabase.auth.refreshSession();
    if (refreshError || !refreshed.session) {
      throw new Error('Not authenticated. Please sign in again.');
    }
    return refreshed.session.access_token;
  }

  // Check if token is expired or about to expire (within 60 seconds)
  const expiresAt = session.expires_at;
  const now = Date.now();
  const expiresAtMs = expiresAt ? expiresAt * 1000 : 0;
  
  // Refresh if expired or expiring within 60 seconds
  if (!expiresAt || expiresAtMs < now + 60000) {
    const { data: refreshed, error: refreshError } = await supabase.auth.refreshSession();
    if (refreshError || !refreshed.session) {
      // If refresh fails but we have a token, try using it anyway
      // The 401 retry logic will handle it if it's truly expired
      if (session.access_token) {
        return session.access_token;
      }
      throw new Error('Session expired. Please sign in again.');
    }
    return refreshed.session.access_token;
  }

  return session.access_token;
}

async function authorizedRequest(path, { method = 'GET', body, headers = {}, retryOn401 = true } = {}) {
  let token = await getValidAccessToken();

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
    try {
      // Always try to refresh on 401, even if we just checked
      const { data: refreshed, error: refreshError } = await supabase.auth.refreshSession();
      if (!refreshError && refreshed.session && refreshed.session.access_token) {
        // Retry with refreshed token (only once)
        const retryHeaders = new Headers(headers);
        retryHeaders.set('Authorization', `Bearer ${refreshed.session.access_token}`);
        if (body && !(body instanceof FormData) && typeof body === 'object' && !(body instanceof Blob)) {
          retryHeaders.set('Content-Type', 'application/json');
        }

        const retryResponse = await fetch(`${getApiUrl()}${path}`, {
          method,
          headers: retryHeaders,
          body: requestBody,
        });
        return retryResponse;
      } else {
        // Refresh failed - log for debugging
        console.warn('Token refresh failed on 401:', refreshError || 'No session returned');
      }
    } catch (refreshErr) {
      // If refresh fails, return original 401 response
      console.error('Failed to refresh token on 401:', refreshErr);
    }
  }

  return response;
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
      throw new Error(message);
    }

    if (response.status === 204) {
      return { data: null, error: null };
    }

    const data = await response.json();
    return { data, error: null };
  } catch (error) {
    console.error('API request failed:', error);
    return { data: null, error };
  }
}

// Fetch user's projects
export async function fetchProjects() {
  try {
    const { data, error } = await supabase
      .from('projects')
      .select('*')
      .order('created_at', { ascending: false });

    if (error) {
      console.error('Error fetching projects:', error);
      console.error('Error details:', JSON.stringify(error, null, 2));
      throw error;
    }
    
    console.log('Fetched projects:', data?.length || 0, 'projects');
    console.log('Projects data:', data);
    
    return { data, error: null };
  } catch (error) {
    console.error('Error fetching projects:', error);
    return { data: null, error };
  }
}

// Fetch crawls for a project
export async function fetchCrawls(projectId) {
  try {
    const { data, error } = await supabase
      .from('crawls')
      .select('*')
      .eq('project_id', projectId)
      .order('started_at', { ascending: false });

    if (error) throw error;
    return { data, error: null };
  } catch (error) {
    return { data: null, error };
  }
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

// Fetch pages for a crawl
export async function fetchPages(crawlId) {
  try {
    const { data, error } = await supabase
      .from('pages')
      .select('*')
      .eq('crawl_id', crawlId)
      .order('created_at', { ascending: false });

    if (error) throw error;
    
    // Flatten the data field - merge data.* fields into the top-level page object
    const flattenedData = (data || []).map(page => {
      const flattened = { ...page };
      if (page.data && typeof page.data === 'object') {
        // Merge data fields into top level
        Object.assign(flattened, page.data);
      }
      return flattened;
    });
    
    return { data: flattenedData, error: null };
  } catch (error) {
    return { data: null, error };
  }
}

// Fetch issues for a crawl
export async function fetchIssues(crawlId) {
  try {
    const { data, error } = await supabase
      .from('issues')
      .select('*')
      .eq('crawl_id', crawlId)
      .order('created_at', { ascending: false });

    if (error) throw error;
    return { data, error: null };
  } catch (error) {
    return { data: null, error };
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
export async function generateCrawlSummary(crawlId, forceRefresh = false) {
  return authorizedJSON('/api/v1/ai/crawl-summary', {
    method: 'POST',
    body: {
      crawl_id: String(crawlId), // Ensure crawlId is a string
      force_refresh: Boolean(forceRefresh) // Explicitly convert to boolean
    }
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
  try {
    const { data, error } = await authorizedJSON(`/api/v1/projects/${projectId}/discover-keywords`, {
      method: 'POST',
      body: discoveryData
    });
    return { data, error };
  } catch (error) {
    console.error('Error discovering keywords:', error);
    return { data: null, error };
  }
}


