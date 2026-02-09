package clarity

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const baseURL = "https://www.clarity.ms/export-data/api/v1"

// InsightMetrics represents metrics returned by the Clarity API
type InsightMetrics struct {
	ScrollDepth     float64 `json:"scrollDepth"`
	EngagementTime  float64 `json:"engagementTime"`
	Traffic         int64   `json:"traffic"`
	DeadClickCount  int64   `json:"deadClickCount"`
	RageClickCount  int64   `json:"rageClickCount"`
	QuickbackClick  int64   `json:"quickbackClick"`
	ExcessiveScroll int64   `json:"excessiveScroll"`
	ErrorClickCount int64   `json:"errorClickCount"`
}

// InsightRow represents a single row of insight data from Clarity
type InsightRow struct {
	DimensionValue string         `json:"dimensionValue"`
	Metrics        InsightMetrics `json:"metrics"`
}

// InsightsResponse represents the API response from Clarity
type InsightsResponse struct {
	Rows []struct {
		Name  string         `json:"name"`
		Value InsightMetrics `json:"value"`
	} `json:"rows"`
	Summary InsightMetrics `json:"summary"`
}

// FetchInsights fetches live insights from Microsoft Clarity for a given dimension.
// Dimension can be: "Url", "Device", "Browser", "Country", "Source", "Medium"
// numDays is limited to 1-3 for the Clarity API.
func FetchInsights(apiToken, projectID string, numDays int, dimension string) ([]InsightRow, *InsightMetrics, error) {
	if numDays < 1 {
		numDays = 1
	}
	if numDays > 3 {
		numDays = 3
	}

	params := url.Values{}
	params.Set("numDays", fmt.Sprintf("%d", numDays))
	params.Set("dimension", dimension)

	reqURL := fmt.Sprintf("%s/project/%s/live-insights?%s", baseURL, projectID, params.Encode())

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("clarity API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("clarity API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp InsightsResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, nil, fmt.Errorf("failed to parse clarity response: %w", err)
	}

	var rows []InsightRow
	for _, r := range apiResp.Rows {
		rows = append(rows, InsightRow{
			DimensionValue: r.Name,
			Metrics:        r.Value,
		})
	}

	return rows, &apiResp.Summary, nil
}

// ValidateToken tests if the API token and project ID are valid
func ValidateToken(apiToken, projectID string) error {
	_, _, err := FetchInsights(apiToken, projectID, 1, "Device")
	if err != nil {
		return fmt.Errorf("invalid Clarity credentials: %w", err)
	}
	return nil
}

// MetricsToMap converts InsightMetrics to a map for storage
func MetricsToMap(m InsightMetrics) map[string]interface{} {
	return map[string]interface{}{
		"scroll_depth":      m.ScrollDepth,
		"engagement_time":   m.EngagementTime,
		"traffic":           m.Traffic,
		"dead_click_count":  m.DeadClickCount,
		"rage_click_count":  m.RageClickCount,
		"quickback_click":   m.QuickbackClick,
		"excessive_scroll":  m.ExcessiveScroll,
		"error_click_count": m.ErrorClickCount,
	}
}
