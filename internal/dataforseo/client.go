package dataforseo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
		
		// Decode the actual response
		if err := json.Unmarshal(bodyBytes, v); err != nil {
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
	var matchedItem *OrganicResultItem
	if targetURL != "" {
		for i := range items {
			// Simple URL matching - could be enhanced with normalization
			if items[i].URL == targetURL {
				matchedItem = &items[i]
				break
			}
		}
	}

	// If no target URL or no match, return first result
	if matchedItem == nil && len(items) > 0 {
		matchedItem = &items[0]
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

