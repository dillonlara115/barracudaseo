package ga4

import (
	"fmt"
	"strings"
	"time"

	analyticsdata "google.golang.org/api/analyticsdata/v1beta"

	"github.com/dillonlara115/barracudaseo/internal/analyzer"
	"github.com/dillonlara115/barracudaseo/pkg/models"
)

// EnrichedIssue extends analyzer.Issue with GA4 performance data
type EnrichedIssue struct {
	Issue                analyzer.Issue         `json:"issue"`
	GA4Performance       *models.GA4Performance `json:"ga4_performance,omitempty"`
	EnrichedPriority     float64                `json:"enriched_priority"`
	RecommendationReason string                 `json:"recommendation_reason"`
}

// DimensionRow represents a single row of GA4 data with a dimension and metrics
type DimensionRow struct {
	RowType        string             `json:"row_type"`
	DimensionValue string             `json:"dimension_value"`
	Metrics        map[string]float64 `json:"metrics"`
}

// FetchPerformanceData fetches GA4 performance data for pages
func FetchPerformanceData(userID string, propertyID string, startDate, endDate time.Time) (map[string]*models.GA4Performance, error) {
	service, err := GetService(userID)
	if err != nil {
		return nil, err
	}

	// Build the request for page-level metrics
	request := &analyticsdata.RunReportRequest{
		DateRanges: []*analyticsdata.DateRange{
			{
				StartDate: startDate.Format("2006-01-02"),
				EndDate:   endDate.Format("2006-01-02"),
			},
		},
		Dimensions: []*analyticsdata.Dimension{
			{Name: "pagePath"},
		},
		Metrics: []*analyticsdata.Metric{
			{Name: "sessions"},
			{Name: "totalUsers"},
			{Name: "screenPageViews"},
			{Name: "bounceRate"},
			{Name: "averageSessionDuration"},
			{Name: "conversions"},
			{Name: "totalRevenue"},
		},
		Limit: 25000, // Max allowed by API
	}

	// Run the report
	response, err := service.Properties.RunReport(fmt.Sprintf("properties/%s", propertyID), request).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to run GA4 report: %w", err)
	}

	// Convert to our model
	performanceMap := make(map[string]*models.GA4Performance)

	for _, row := range response.Rows {
		if len(row.DimensionValues) == 0 || len(row.MetricValues) == 0 {
			continue
		}

		pagePath := row.DimensionValues[0].Value
		// Normalize URL to match crawl results
		normalizedURL := normalizeURL(pagePath)

		// Parse metrics
		sessions := parseMetricValue(row.MetricValues[0].Value)
		users := parseMetricValue(row.MetricValues[1].Value)
		pageViews := parseMetricValue(row.MetricValues[2].Value)
		bounceRate := parseFloatValue(row.MetricValues[3].Value)
		avgDuration := parseFloatValue(row.MetricValues[4].Value)
		conversions := parseMetricValue(row.MetricValues[5].Value)
		revenue := parseFloatValue(row.MetricValues[6].Value)

		performanceMap[normalizedURL] = &models.GA4Performance{
			URL:                normalizedURL,
			Sessions:           sessions,
			Users:              users,
			PageViews:          pageViews,
			BounceRate:         bounceRate,
			AvgSessionDuration: avgDuration,
			Conversions:        conversions,
			Revenue:            revenue,
			LastUpdated:        time.Now(),
		}
	}

	return performanceMap, nil
}

// FetchMultiDimensionData fetches GA4 data across multiple dimensions (source, medium, device, country, date)
func FetchMultiDimensionData(userID string, propertyID string, startDate, endDate time.Time) (map[string][]DimensionRow, error) {
	service, err := GetService(userID)
	if err != nil {
		return nil, err
	}

	dateRange := &analyticsdata.DateRange{
		StartDate: startDate.Format("2006-01-02"),
		EndDate:   endDate.Format("2006-01-02"),
	}

	metrics := []*analyticsdata.Metric{
		{Name: "sessions"},
		{Name: "totalUsers"},
		{Name: "screenPageViews"},
		{Name: "bounceRate"},
		{Name: "averageSessionDuration"},
		{Name: "conversions"},
	}

	dimensions := map[string]string{
		"source":  "sessionSource",
		"medium":  "sessionMedium",
		"device":  "deviceCategory",
		"country": "country",
		"date":    "date",
	}

	result := make(map[string][]DimensionRow)

	for rowType, dimName := range dimensions {
		request := &analyticsdata.RunReportRequest{
			DateRanges: []*analyticsdata.DateRange{dateRange},
			Dimensions: []*analyticsdata.Dimension{{Name: dimName}},
			Metrics:    metrics,
			Limit:      10000,
		}

		response, err := service.Properties.RunReport(fmt.Sprintf("properties/%s", propertyID), request).Do()
		if err != nil {
			return nil, fmt.Errorf("failed to run GA4 %s report: %w", rowType, err)
		}

		var rows []DimensionRow
		for _, row := range response.Rows {
			if len(row.DimensionValues) == 0 || len(row.MetricValues) == 0 {
				continue
			}

			dr := DimensionRow{
				RowType:        rowType,
				DimensionValue: row.DimensionValues[0].Value,
				Metrics: map[string]float64{
					"sessions":             float64(parseMetricValue(row.MetricValues[0].Value)),
					"users":                float64(parseMetricValue(row.MetricValues[1].Value)),
					"page_views":           float64(parseMetricValue(row.MetricValues[2].Value)),
					"bounce_rate":          parseFloatValue(row.MetricValues[3].Value),
					"avg_session_duration": parseFloatValue(row.MetricValues[4].Value),
					"conversions":          float64(parseMetricValue(row.MetricValues[5].Value)),
				},
			}
			rows = append(rows, dr)
		}
		result[rowType] = rows
	}

	return result, nil
}

// normalizeURL normalizes URLs to match crawl results
func normalizeURL(url string) string {
	// Remove leading slash if present
	url = strings.TrimPrefix(url, "/")
	// Remove trailing slash
	url = strings.TrimSuffix(url, "/")
	// Ensure lowercase
	url = strings.ToLower(url)
	return url
}

// parseMetricValue parses a metric value string to int64
func parseMetricValue(value string) int64 {
	var result int64
	fmt.Sscanf(value, "%d", &result)
	return result
}

// parseFloatValue parses a float value string
func parseFloatValue(value string) float64 {
	var result float64
	fmt.Sscanf(value, "%f", &result)
	return result
}

// EnrichIssues merges GA4 performance data with issues
func EnrichIssues(issues []analyzer.Issue, performanceMap map[string]*models.GA4Performance) []EnrichedIssue {
	enriched := make([]EnrichedIssue, 0, len(issues))

	for _, issue := range issues {
		enrichedIssue := EnrichedIssue{
			Issue: issue,
		}

		// Normalize issue URL to match GA4 data
		normalizedURL := normalizeURL(issue.URL)

		// Find matching performance data
		if perf, exists := performanceMap[normalizedURL]; exists {
			enrichedIssue.GA4Performance = perf
			enrichedIssue.EnrichedPriority = calculateEnrichedPriority(issue, perf)
			enrichedIssue.RecommendationReason = generateRecommendationReason(issue, perf)
		} else {
			// No GA4 data available - use base priority
			enrichedIssue.EnrichedPriority = float64(getSeverityWeight(issue.Severity))
		}

		enriched = append(enriched, enrichedIssue)
	}

	return enriched
}

// calculateEnrichedPriority calculates priority with GA4 data
func calculateEnrichedPriority(issue analyzer.Issue, perf *models.GA4Performance) float64 {
	basePriority := float64(getSeverityWeight(issue.Severity))

	// Traffic multiplier based on sessions
	trafficMultiplier := 1.0
	if perf.Sessions > 10000 {
		trafficMultiplier = 3.0 // High traffic = 3x priority
	} else if perf.Sessions > 1000 {
		trafficMultiplier = 2.0 // Medium traffic = 2x priority
	} else if perf.Sessions < 100 {
		trafficMultiplier = 0.5 // Low traffic = 0.5x priority
	}

	// Bounce rate multiplier (high bounce = opportunity)
	bounceMultiplier := 1.0
	if perf.BounceRate > 70 && perf.Sessions > 500 {
		bounceMultiplier = 1.5 // High bounce with decent traffic = opportunity
	}

	// Conversion multiplier (pages with conversions are critical)
	conversionMultiplier := 1.0
	if perf.Conversions > 0 {
		conversionMultiplier = 2.0 // Converting pages = 2x priority
	}

	return basePriority * trafficMultiplier * bounceMultiplier * conversionMultiplier
}

// generateRecommendationReason creates contextual recommendation based on GA4 data
func generateRecommendationReason(issue analyzer.Issue, perf *models.GA4Performance) string {
	if perf.Sessions > 10000 {
		return fmt.Sprintf("This page has high traffic (%d sessions/month). Fixing this issue could significantly impact user experience and conversions.", perf.Sessions)
	} else if perf.Sessions > 1000 {
		if perf.BounceRate > 70 {
			return fmt.Sprintf("This page has moderate traffic (%d sessions/month) but high bounce rate (%.1f%%). Optimizing this could improve engagement.", perf.Sessions, perf.BounceRate)
		}
		return fmt.Sprintf("This page has moderate traffic (%d sessions/month).", perf.Sessions)
	} else if perf.Conversions > 0 {
		return fmt.Sprintf("This page generates conversions (%d conversions). Fixing this issue is critical for revenue.", perf.Conversions)
	} else if perf.Sessions < 100 {
		return "This page has minimal traffic. Consider fixing as part of broader technical SEO improvements."
	}
	return ""
}

// getSeverityWeight returns weight for severity level
func getSeverityWeight(severity string) int {
	switch severity {
	case "error":
		return 10
	case "warning":
		return 5
	case "info":
		return 1
	default:
		return 1
	}
}
