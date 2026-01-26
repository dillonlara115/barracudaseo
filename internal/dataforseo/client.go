package dataforseo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Client represents a DataForSEO API client
type Client struct {
	httpClient *http.Client
	baseURL    string
	login      string
	password   string
}

// NewClient creates a new DataForSEO client
func NewClient() (*Client, error) {
	baseURL := os.Getenv("DATAFORSEO_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.dataforseo.com"
	}

	login := os.Getenv("DATAFORSEO_LOGIN")
	password := os.Getenv("DATAFORSEO_PASSWORD")

	if login == "" || password == "" {
		return nil, fmt.Errorf("DATAFORSEO_LOGIN and DATAFORSEO_PASSWORD must be set")
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:  baseURL,
		login:    login,
		password: password,
	}, nil
}

// do performs an HTTP request to DataForSEO API
func (c *Client) do(ctx context.Context, method, path string, body interface{}, v interface{}) error {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return fmt.Errorf("encode body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, &buf)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	req.SetBasicAuth(c.login, c.password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// DataForSEO returns HTTP 200 even for errors, with error info in JSON body
	// Read the body first to check for JSON status codes
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode >= 300 {
		// HTTP-level error
		var errorResp struct {
			StatusCode    int    `json:"status_code"`
			StatusMessage string `json:"status_message"`
		}
		if err := json.Unmarshal(bodyBytes, &errorResp); err == nil {
			return fmt.Errorf("dataforseo API error: status %d, message: %s", errorResp.StatusCode, errorResp.StatusMessage)
		}
		return fmt.Errorf("dataforseo API error: HTTP status %d", resp.StatusCode)
	}

	if v != nil {
		// Check for DataForSEO error codes in JSON response (even with HTTP 200)
		var checkResp struct {
			StatusCode    int    `json:"status_code"`
			StatusMessage string `json:"status_message"`
		}
		if err := json.Unmarshal(bodyBytes, &checkResp); err == nil {
			// 40400 = Not Found, 40000+ = errors
			if checkResp.StatusCode >= 40000 {
				return fmt.Errorf("dataforseo API error: status %d, message: %s", checkResp.StatusCode, checkResp.StatusMessage)
			}
		}

		// For ranked keywords API, log raw response for debugging
		if strings.Contains(path, "ranked_keywords") && len(bodyBytes) > 0 && len(bodyBytes) < 50000 {
			// Log first 1000 chars of response for debugging
			preview := string(bodyBytes)
			if len(preview) > 1000 {
				preview = preview[:1000] + "..."
			}
			// This will help us see the actual response structure
		}

		// Decode the actual response
		if err := json.Unmarshal(bodyBytes, v); err != nil {
			// Log response body for debugging if decode fails
			if len(bodyBytes) > 0 && len(bodyBytes) < 10000 {
				return fmt.Errorf("decode response: %w, response body: %s", err, string(bodyBytes))
			}
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

// CreateOrganicTask creates a new organic SERP task
func (c *Client) CreateOrganicTask(ctx context.Context, task OrganicTaskPost) (*OrganicTaskPostResponse, error) {
	// DataForSEO expects tasks as a map with numeric string keys
	body := OrganicTaskPostRequest{
		"0": task,
	}

	var resp OrganicTaskPostResponse
	if err := c.do(ctx, http.MethodPost, "/v3/serp/google/organic/task_post", body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetOrganicTask retrieves the result of an organic SERP task
func (c *Client) GetOrganicTask(ctx context.Context, taskID string) (*OrganicTaskGetResponse, error) {
	var resp OrganicTaskGetResponse
	// DataForSEO API format: GET /v3/serp/google/organic/task_get/{id}
	path := fmt.Sprintf("/v3/serp/google/organic/task_get/%s", taskID)
	if err := c.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetOrganicTasksReady retrieves a list of completed organic SERP tasks that are ready for collection
func (c *Client) GetOrganicTasksReady(ctx context.Context) (*OrganicTasksReadyResponse, error) {
	var resp OrganicTasksReadyResponse
	path := "/v3/serp/google/organic/tasks_ready"
	if err := c.do(ctx, http.MethodGet, path, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// CreateOrganicTaskLive creates a new organic SERP task using the Live API
// Live API returns results immediately in a single request (no polling needed)
// More expensive but instant - perfect for "check now" functionality
func (c *Client) CreateOrganicTaskLive(ctx context.Context, task OrganicTaskPost) (*OrganicTaskGetResponse, error) {
	// DataForSEO expects tasks as a map with numeric string keys
	body := OrganicLiveRequest{
		"0": task,
	}

	var resp OrganicTaskGetResponse
	if err := c.do(ctx, http.MethodPost, "/v3/serp/google/organic/live/regular", body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// IsTaskReady checks if a task is ready (status_code 20000 means task is ready)
func IsTaskReady(resp *OrganicTaskGetResponse) bool {
	if len(resp.Tasks) == 0 {
		return false
	}
	// Status code 20000 means task is ready
	return resp.Tasks[0].StatusCode == 20000
}

// ExtractRanking extracts ranking information from task result for a target URL
func ExtractRanking(result *OrganicTaskGetResponse, targetURL string) (*RankingData, error) {
	if len(result.Tasks) == 0 || len(result.Tasks[0].Result) == 0 {
		return nil, fmt.Errorf("no result data available")
	}

	items := result.Tasks[0].Result[0].Items

	// Find the item matching target URL (if provided)
	normalize := func(u string) string {
		u = strings.TrimSpace(strings.ToLower(u))
		u = strings.TrimSuffix(u, "/")
		return u
	}

	var matchedItem *OrganicResultItem
	if targetURL != "" {
		targetNorm := normalize(targetURL)
		for i := range items {
			if normalize(items[i].URL) == targetNorm {
				matchedItem = &items[i]
				break
			}
		}
		
		// If target URL was provided but no match found, the site is not ranking
		if matchedItem == nil {
			return nil, fmt.Errorf("target URL %s is not ranking in the search results", targetURL)
		}
	} else {
		// No target URL provided - use first organic result as reference
		for i := range items {
			if strings.ToLower(items[i].Type) == "organic" {
				matchedItem = &items[i]
				break
			}
		}
		
		// Fallback to first result if no organic result found
		if matchedItem == nil && len(items) > 0 {
			matchedItem = &items[0]
		}
	}

	if matchedItem == nil {
		return nil, fmt.Errorf("no ranking data found")
	}

	return &RankingData{
		PositionAbsolute: matchedItem.RankAbsolute,
		PositionOrganic:  matchedItem.RankGroup,
		URL:              matchedItem.URL,
		Title:            matchedItem.Title,
		Snippet:          matchedItem.Description,
		SERPFeatures:     matchedItem.SERPFeatures,
	}, nil
}

// RankingData represents extracted ranking information
type RankingData struct {
	PositionAbsolute int
	PositionOrganic  int
	URL              string
	Title            string
	Snippet          string
	SERPFeatures     []string
}

// GetRankedKeywordsLive discovers keywords that a domain/URL is currently ranking for
// Uses the DataForSEO Labs Ranked Keywords API (Live endpoint)
// Returns results immediately - no polling needed
// Note: This API uses an array format for the request body, not a map like SERP API
func (c *Client) GetRankedKeywordsLive(ctx context.Context, task RankedKeywordsTask) (*RankedKeywordsResponse, error) {
	// DataForSEO Ranked Keywords API expects an array format: [{...}]
	// Unlike SERP API which uses map format: {"0": {...}}
	body := RankedKeywordsRequest{task}

	var resp RankedKeywordsResponse
	// Endpoint: /v3/dataforseo_labs/google/ranked_keywords/live
	// Based on docs: https://docs.dataforseo.com/v3/dataforseo_labs/google/ranked_keywords/live/
	path := "/v3/dataforseo_labs/google/ranked_keywords/live"
	if err := c.do(ctx, http.MethodPost, path, body, &resp); err != nil {
		// Provide helpful error message
		if strings.Contains(err.Error(), "40400") || strings.Contains(err.Error(), "Not Found") {
			return nil, fmt.Errorf("Ranked Keywords API endpoint not found (40400). This feature may require a DataForSEO Labs subscription or the endpoint may not be available in your account tier. Error: %w", err)
		}
		return nil, err
	}

	return &resp, nil
}
