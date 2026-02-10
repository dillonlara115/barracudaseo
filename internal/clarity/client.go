package clarity

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const baseURL = "https://www.clarity.ms/export-data/api/v1"

// RateLimitError indicates the Clarity API daily limit (10 req/day) was exceeded.
// RetryAfter is the UTC midnight when the limit resets.
var ErrRateLimited = errors.New("clarity API rate limit exceeded")

// RateLimitError wraps ErrRateLimited with RetryAfter timestamp.
type RateLimitError struct {
	RetryAfter time.Time
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("clarity API rate limit exceeded, retry after %s", e.RetryAfter.Format(time.RFC3339))
}

func (e *RateLimitError) Unwrap() error { return ErrRateLimited }

// nextMidnightUTC returns the next UTC midnight (when Clarity daily limit resets).
func nextMidnightUTC(now time.Time) time.Time {
	y, m, d := now.Date()
	return time.Date(y, m, d+1, 0, 0, 0, 0, time.UTC)
}

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

// InsightsResponse represents the legacy object-format API response (unused; API returns array)
type InsightsResponse struct {
	Rows []struct {
		Name  string         `json:"name"`
		Value InsightMetrics `json:"value"`
	} `json:"rows"`
	Summary InsightMetrics `json:"summary"`
}

// clarityAPIArrayEntry represents one metric block in the Microsoft array response.
// Response format: [{"metricName":"Traffic","information":[{...}]}, ...]
type clarityAPIArrayEntry struct {
	MetricName  string                   `json:"metricName"`
	Information []map[string]interface{} `json:"information"`
}

// dimensionKeysForDimension returns possible keys used in info records (API may vary casing).
func dimensionKeysForDimension(dimension string) []string {
	switch dimension {
	case "Url":
		return []string{"URL", "Url", "url"}
	default:
		return []string{dimension}
	}
}

// parseClarityArrayResponse parses the Microsoft array format into rows and summary.
func parseClarityArrayResponse(body []byte, dimension string) ([]InsightRow, *InsightMetrics, error) {
	var arr []clarityAPIArrayEntry
	if err := json.Unmarshal(body, &arr); err != nil {
		return nil, nil, fmt.Errorf("failed to parse array: %w", err)
	}

	dimKeys := dimensionKeysForDimension(dimension)
	rowMap := make(map[string]*InsightMetrics) // dimension value -> aggregated metrics
	var summaryTotals InsightMetrics

	for _, entry := range arr {
		metricName := entry.MetricName
		for _, info := range entry.Information {
			var dimVal string
			for _, k := range dimKeys {
				if v, ok := info[k].(string); ok && v != "" {
					dimVal = v
					break
				}
			}
			if dimVal == "" {
				// Try common dimension keys as fallback
				for _, k := range []string{"Device", "URL", "Url", "OS", "Source", "Browser"} {
					if v, ok := info[k].(string); ok && v != "" {
						dimVal = v
						break
					}
				}
			}
			if dimVal == "" {
				continue
			}

			if rowMap[dimVal] == nil {
				rowMap[dimVal] = &InsightMetrics{}
			}
			row := rowMap[dimVal]

			// Extract metric value based on metricName
			switch metricName {
			case "Traffic":
				if v, ok := parseNum(info["totalSessionCount"]); ok {
					row.Traffic += v
					summaryTotals.Traffic += v
				}
			case "Scroll Depth":
				if v, ok := parseFloat64(info["scrollDepth"]); ok {
					row.ScrollDepth = v
				}
				if v, ok := parseFloat64(info["averageScrollDepthPercentage"]); ok && row.ScrollDepth == 0 {
					row.ScrollDepth = v / 100
				}
			case "Engagement Time":
				if v, ok := parseFloat64(info["engagementTime"]); ok {
					row.EngagementTime = v
				}
				if v, ok := parseFloat64(info["averageEngagementTime"]); ok && row.EngagementTime == 0 {
					row.EngagementTime = v
				}
			case "Rage Click Count":
				if v, ok := parseNum(info["rageClickCount"]); ok {
					row.RageClickCount += v
					summaryTotals.RageClickCount += v
				}
			case "Dead Click Count":
				if v, ok := parseNum(info["deadClickCount"]); ok {
					row.DeadClickCount += v
					summaryTotals.DeadClickCount += v
				}
			case "Quickback Click":
				if v, ok := parseNum(info["quickbackClick"]); ok {
					row.QuickbackClick += v
					summaryTotals.QuickbackClick += v
				}
			case "Excessive Scroll":
				if v, ok := parseNum(info["excessiveScroll"]); ok {
					row.ExcessiveScroll += v
					summaryTotals.ExcessiveScroll += v
				}
			case "Error Click Count", "Script Error Count":
				if v, ok := parseNum(info["errorClickCount"]); ok {
					row.ErrorClickCount += v
					summaryTotals.ErrorClickCount += v
				}
				if v, ok := parseNum(info["scriptErrorCount"]); ok {
					row.ErrorClickCount += v
					summaryTotals.ErrorClickCount += v
				}
			}
		}
	}

	// Build rows slice
	var rows []InsightRow
	for dimVal, metrics := range rowMap {
		rows = append(rows, InsightRow{DimensionValue: dimVal, Metrics: *metrics})
	}

	return rows, &summaryTotals, nil
}

func parseNum(v interface{}) (int64, bool) {
	switch n := v.(type) {
	case float64:
		return int64(n), true
	case int:
		return int64(n), true
	case int64:
		return n, true
	case string:
		var i int64
		_, err := fmt.Sscanf(n, "%d", &i)
		return i, err == nil
	}
	return 0, false
}

func parseFloat64(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case string:
		var f float64
		_, err := fmt.Sscanf(n, "%f", &f)
		return f, err == nil
	}
	return 0, false
}

// FetchInsights fetches live insights from Microsoft Clarity for a given dimension.
// Uses the documented project-live-insights endpoint. The API token is project-scoped
// (generated per project in Clarity Settings → Data Export), so projectID is stored
// for display only and is not sent to the API.
// Dimension can be: "Url", "Device", "Browser", "Country", "Source", "Medium", "OS"
// numDays is limited to 1-3 for the Clarity API.
func FetchInsights(apiToken, projectID string, numDays int, dimension string) ([]InsightRow, *InsightMetrics, error) {
	if numDays < 1 {
		numDays = 1
	}
	if numDays > 3 {
		numDays = 3
	}

	params := url.Values{}
	params.Set("numOfDays", fmt.Sprintf("%d", numDays))
	params.Set("dimension1", dimension)

	// Documented endpoint: https://www.clarity.ms/export-data/api/v1/project-live-insights
	// Token is project-scoped, no project ID in URL
	reqURL := fmt.Sprintf("%s/project-live-insights?%s", baseURL, params.Encode())

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

	if resp.StatusCode == http.StatusTooManyRequests ||
		(resp.StatusCode == 429 || (resp.StatusCode >= 400 && strings.Contains(string(body), "Exceeded daily limit"))) {
		retryAfter := nextMidnightUTC(time.Now().UTC())
		return nil, nil, &RateLimitError{RetryAfter: retryAfter}
	}
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("clarity API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Microsoft API returns an array: [{"metricName":"Traffic","information":[...]}, ...]
	// Try array format first - attempt unmarshal to detect format
	var arr []clarityAPIArrayEntry
	if err := json.Unmarshal(body, &arr); err == nil {
		return parseClarityArrayResponse(body, dimension)
	}

	// Fallback: legacy object format {rows, summary}
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

// MultiDimensionResult holds rows and summary for a single dimension type.
type MultiDimensionResult struct {
	Rows    []InsightRow
	Summary *InsightMetrics
}

// FetchInsightsMulti fetches insights for up to 3 dimensions in a single API call.
// Reduces API usage from 3 calls to 1. dimensions should be e.g. []string{"Url", "Device", "Source"}.
func FetchInsightsMulti(apiToken, projectID string, numDays int, dimensions []string) (map[string]*MultiDimensionResult, error) {
	if numDays < 1 {
		numDays = 1
	}
	if numDays > 3 {
		numDays = 3
	}
	if len(dimensions) == 0 || len(dimensions) > 3 {
		return nil, fmt.Errorf("dimensions must have 1-3 elements")
	}

	params := url.Values{}
	params.Set("numOfDays", fmt.Sprintf("%d", numDays))
	for i, d := range dimensions {
		if i == 0 {
			params.Set("dimension1", d)
		} else if i == 1 {
			params.Set("dimension2", d)
		} else {
			params.Set("dimension3", d)
		}
	}

	reqURL := fmt.Sprintf("%s/project-live-insights?%s", baseURL, params.Encode())
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("clarity API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode == http.StatusTooManyRequests ||
		(resp.StatusCode == 429 || (resp.StatusCode >= 400 && strings.Contains(string(body), "Exceeded daily limit"))) {
		retryAfter := nextMidnightUTC(time.Now().UTC())
		return nil, &RateLimitError{RetryAfter: retryAfter}
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("clarity API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Map API dimension names to our row_type keys
	apiToRowType := map[string]string{
		"Url":    "url",
		"Device": "device",
		"Source": "source",
	}
	results := parseClarityArrayResponseMulti(body, dimensions, apiToRowType)
	return results, nil
}

// parseClarityArrayResponseMulti parses the array response when multiple dimensions are requested.
// Builds separate row maps per dimension and aggregates metrics.
func parseClarityArrayResponseMulti(body []byte, dimensions []string, apiToRowType map[string]string) map[string]*MultiDimensionResult {
	var arr []clarityAPIArrayEntry
	if err := json.Unmarshal(body, &arr); err != nil {
		return nil
	}

	// rowType -> dimensionValue -> metrics
	rowMaps := make(map[string]map[string]*InsightMetrics)
	var summaryTotals InsightMetrics

	for _, entry := range arr {
		metricName := entry.MetricName
		for _, info := range entry.Information {
			for _, apiDim := range dimensions {
				rowType, ok := apiToRowType[apiDim]
				if !ok {
					rowType = strings.ToLower(apiDim)
				}
				dimKeys := dimensionKeysForDimension(apiDim)
				var dimVal string
				for _, k := range dimKeys {
					if v, ok := info[k].(string); ok && v != "" {
						dimVal = v
						break
					}
				}
				if dimVal == "" {
					if v, ok := info[apiDim].(string); ok && v != "" {
						dimVal = v
					}
				}
				if dimVal == "" {
					continue
				}
				if rowMaps[rowType] == nil {
					rowMaps[rowType] = make(map[string]*InsightMetrics)
				}
				if rowMaps[rowType][dimVal] == nil {
					rowMaps[rowType][dimVal] = &InsightMetrics{}
				}
				row := rowMaps[rowType][dimVal]

				// Only add to summary once per info row (when on first dimension)
				addToSummary := apiDim == dimensions[0]

				switch metricName {
				case "Traffic":
					if v, ok := parseNum(info["totalSessionCount"]); ok {
						row.Traffic += v
						if addToSummary {
							summaryTotals.Traffic += v
						}
					}
				case "Scroll Depth":
					if v, ok := parseFloat64(info["scrollDepth"]); ok {
						row.ScrollDepth = v
					}
					if v, ok := parseFloat64(info["averageScrollDepthPercentage"]); ok && row.ScrollDepth == 0 {
						row.ScrollDepth = v / 100
					}
				case "Engagement Time":
					if v, ok := parseFloat64(info["engagementTime"]); ok {
						row.EngagementTime = v
					}
					if v, ok := parseFloat64(info["averageEngagementTime"]); ok && row.EngagementTime == 0 {
						row.EngagementTime = v
					}
				case "Rage Click Count":
					if v, ok := parseNum(info["rageClickCount"]); ok {
						row.RageClickCount += v
						if addToSummary {
							summaryTotals.RageClickCount += v
						}
					}
				case "Dead Click Count":
					if v, ok := parseNum(info["deadClickCount"]); ok {
						row.DeadClickCount += v
					}
				case "Quickback Click":
					if v, ok := parseNum(info["quickbackClick"]); ok {
						row.QuickbackClick += v
					}
				case "Excessive Scroll":
					if v, ok := parseNum(info["excessiveScroll"]); ok {
						row.ExcessiveScroll += v
					}
				case "Error Click Count", "Script Error Count":
					if v, ok := parseNum(info["errorClickCount"]); ok {
						row.ErrorClickCount += v
					}
					if v, ok := parseNum(info["scriptErrorCount"]); ok {
						row.ErrorClickCount += v
					}
				}
			}
		}
	}

	out := make(map[string]*MultiDimensionResult)
	for rowType, m := range rowMaps {
		var rows []InsightRow
		for dimVal, metrics := range m {
			rows = append(rows, InsightRow{DimensionValue: dimVal, Metrics: *metrics})
		}
		sum := summaryTotals
		out[rowType] = &MultiDimensionResult{Rows: rows, Summary: &sum}
	}
	return out
}

// ValidateToken tests if the API token is valid by making a minimal API request.
// projectID is stored for display; the token is project-scoped so it identifies the project.
func ValidateToken(apiToken, projectID string) error {
	if apiToken == "" {
		return fmt.Errorf("API token is required")
	}
	params := url.Values{}
	params.Set("numOfDays", "1")
	params.Set("dimension1", "Device")
	reqURL := fmt.Sprintf("%s/project-live-insights?%s", baseURL, params.Encode())

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("clarity API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("invalid or expired token (HTTP %d)", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("clarity API returned status %d: %s", resp.StatusCode, string(body))
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
