package gsc

import (
	"fmt"
	"sort"
	"time"

	"github.com/dillonlara115/barracudaseo/pkg/models"
	"google.golang.org/api/searchconsole/v1"
)

// PerformanceRow represents aggregated metrics for a single dimension value.
type PerformanceRow struct {
	Value   string
	Metrics map[string]float64
}

// PerformanceReport aggregates Search Console metrics for multiple dimensions.
type PerformanceReport struct {
	Totals      map[string]float64
	Queries     []PerformanceRow
	Pages       []PerformanceRow
	Countries   []PerformanceRow
	Devices     []PerformanceRow
	Appearance  []PerformanceRow
	Dates       []PerformanceRow
	PageQueries map[string][]models.Query
}

// FetchPerformanceReport returns a comprehensive performance report across dimensions for the given property.
func FetchPerformanceReport(userID, siteURL string, startDate, endDate time.Time) (*PerformanceReport, error) {
	service, err := GetService(userID)
	if err != nil {
		return nil, err
	}

	report := &PerformanceReport{
		Totals:      map[string]float64{"clicks": 0, "impressions": 0, "ctr": 0, "position": 0},
		PageQueries: make(map[string][]models.Query),
	}

	// Totals (no dimensions)
	if totals, err := queryTotals(service, siteURL, startDate, endDate); err == nil {
		report.Totals = totals
	} else {
		return nil, fmt.Errorf("failed to fetch totals: %w", err)
	}

	// Queries
	if rows, err := queryDimension(service, siteURL, startDate, endDate, []string{"query"}, 25000, nil); err == nil {
		report.Queries = toPerformanceRows(rows, 0)
	} else {
		return nil, fmt.Errorf("failed to fetch queries: %w", err)
	}

	// Pages
	pageRows, err := queryDimension(service, siteURL, startDate, endDate, []string{"page"}, 25000, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pages: %w", err)
	}
	report.Pages = toPerformanceRows(pageRows, 1)

	// Countries
	if rows, err := queryDimension(service, siteURL, startDate, endDate, []string{"country"}, 200, nil); err == nil {
		report.Countries = toPerformanceRows(rows, 0)
	} else {
		return nil, fmt.Errorf("failed to fetch countries: %w", err)
	}

	// Devices
	if rows, err := queryDimension(service, siteURL, startDate, endDate, []string{"device"}, 10, nil); err == nil {
		report.Devices = toPerformanceRows(rows, 0)
	} else {
		return nil, fmt.Errorf("failed to fetch devices: %w", err)
	}

	// Search appearance
	if rows, err := queryDimension(service, siteURL, startDate, endDate, []string{"searchAppearance"}, 50, nil); err == nil {
		report.Appearance = toPerformanceRows(rows, 0)
	} else {
		return nil, fmt.Errorf("failed to fetch search appearance: %w", err)
	}

	// Dates (limited to 400 to keep response manageable)
	if rows, err := queryDimension(service, siteURL, startDate, endDate, []string{"date"}, 400, nil); err == nil {
		report.Dates = toPerformanceRows(rows, 0)
	} else {
		return nil, fmt.Errorf("failed to fetch daily metrics: %w", err)
	}

	// Top queries per page (limit to the top 50 pages by impressions)
	report.PageQueries = fetchTopQueriesForPages(service, siteURL, startDate, endDate, report.Pages, 50)

	return report, nil
}

func queryTotals(service *searchconsole.Service, siteURL string, startDate, endDate time.Time) (map[string]float64, error) {
	request := &searchconsole.SearchAnalyticsQueryRequest{
		StartDate: startDate.Format("2006-01-02"),
		EndDate:   endDate.Format("2006-01-02"),
		RowLimit:  1,
	}

	response, err := service.Searchanalytics.Query(siteURL, request).Do()
	if err != nil {
		return nil, err
	}

	totals := map[string]float64{"clicks": 0, "impressions": 0, "ctr": 0, "position": 0}
	if len(response.Rows) > 0 {
		row := response.Rows[0]
		totals["clicks"] = row.Clicks
		totals["impressions"] = row.Impressions
		totals["ctr"] = row.Ctr
		totals["position"] = row.Position
	}

	return totals, nil
}

func queryDimension(service *searchconsole.Service, siteURL string, startDate, endDate time.Time, dimensions []string, limit int64, filterGroups []*searchconsole.ApiDimensionFilterGroup) ([]*searchconsole.ApiDataRow, error) {
	request := &searchconsole.SearchAnalyticsQueryRequest{
		StartDate:             startDate.Format("2006-01-02"),
		EndDate:               endDate.Format("2006-01-02"),
		Dimensions:            dimensions,
		RowLimit:              limit,
		DimensionFilterGroups: filterGroups,
	}

	response, err := service.Searchanalytics.Query(siteURL, request).Do()
	if err != nil {
		return nil, err
	}

	return response.Rows, nil
}

func toPerformanceRows(rows []*searchconsole.ApiDataRow, normalizePageIndex int) []PerformanceRow {
	result := make([]PerformanceRow, 0, len(rows))
	for _, row := range rows {
		value := ""
		if len(row.Keys) > 0 {
			value = row.Keys[0]
		}
		if normalizePageIndex >= 0 && len(row.Keys) > normalizePageIndex {
			value = normalizeURL(row.Keys[normalizePageIndex])
		}

		result = append(result, PerformanceRow{
			Value: value,
			Metrics: map[string]float64{
				"clicks":      row.Clicks,
				"impressions": row.Impressions,
				"ctr":         row.Ctr,
				"position":    row.Position,
			},
		})
	}

	return result
}

func fetchTopQueriesForPages(service *searchconsole.Service, siteURL string, startDate, endDate time.Time, pages []PerformanceRow, limit int) map[string][]models.Query {
	if len(pages) == 0 || limit <= 0 {
		return map[string][]models.Query{}
	}

	// Sort by impressions descending to get top pages
	sorted := make([]PerformanceRow, len(pages))
	copy(sorted, pages)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Metrics["impressions"] > sorted[j].Metrics["impressions"]
	})

	if len(sorted) > limit {
		sorted = sorted[:limit]
	}

	result := make(map[string][]models.Query, len(sorted))
	for _, row := range sorted {
		pageURL := row.Value
		request := &searchconsole.SearchAnalyticsQueryRequest{
			StartDate:  startDate.Format("2006-01-02"),
			EndDate:    endDate.Format("2006-01-02"),
			Dimensions: []string{"query"},
			DimensionFilterGroups: []*searchconsole.ApiDimensionFilterGroup{
				{
					Filters: []*searchconsole.ApiDimensionFilter{
						{
							Dimension:  "page",
							Expression: pageURL,
							Operator:   "equals",
						},
					},
				},
			},
			RowLimit: 10,
		}

		response, err := service.Searchanalytics.Query(siteURL, request).Do()
		if err != nil || len(response.Rows) == 0 {
			continue
		}

		queries := make([]models.Query, 0, len(response.Rows))
		for _, row := range response.Rows {
			if len(row.Keys) == 0 {
				continue
			}
			queries = append(queries, models.Query{
				Query:       row.Keys[0],
				Impressions: int64(row.Impressions),
				Clicks:      int64(row.Clicks),
				CTR:         row.Ctr,
				Position:    row.Position,
			})
		}

		if len(queries) > 0 {
			result[pageURL] = queries
		}
	}

	return result
}
