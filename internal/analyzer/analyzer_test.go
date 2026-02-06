package analyzer

import (
	"testing"

	"github.com/dillonlara115/barracudaseo/pkg/models"
)

func TestAnalyze(t *testing.T) {
	tests := []struct {
		name           string
		results        []*models.PageResult
		expectedIssues map[IssueType]int
		expectedPages  int
	}{
		{
			name: "Perfect Page",
			results: []*models.PageResult{
				{
					URL:                "https://example.com",
					StatusCode:         200,
					Title:              "Perfect Title for SEO Optimization",                                                                                                                  // ~34 chars
					MetaDesc:           "This is a perfect meta description that falls right within the recommended length range of 120 to 160 characters for optimal search engine display.", // ~145 chars
					H1:                 []string{"Main Heading"},
					Canonical:          "https://example.com",
					IndexabilityStatus: models.IndexabilityIndexable,
					ResponseTime:       200,
				},
			},
			expectedIssues: map[IssueType]int{}, // No issues expected
			expectedPages:  1,
		},
		{
			name: "Missing Title and H1",
			results: []*models.PageResult{
				{
					URL:                "https://example.com/bad",
					StatusCode:         200,
					Title:              "",
					H1:                 []string{},
					MetaDesc:           "Valid meta description for this test case to avoid noise.",
					Canonical:          "https://example.com/bad",
					IndexabilityStatus: models.IndexabilityIndexable,
				},
			},
			expectedIssues: map[IssueType]int{
				IssueMissingTitle:  1,
				IssueMissingH1:     1,
				IssueShortMetaDesc: 1, // "Valid..." is < 120 chars
			},
			expectedPages: 1,
		},
		{
			name: "Length Issues",
			results: []*models.PageResult{
				{
					URL:                "https://example.com/lengths",
					StatusCode:         200,
					Title:              "Hi",    // Too short
					MetaDesc:           "Short", // Too short
					H1:                 []string{"Heading"},
					Canonical:          "https://example.com/lengths",
					IndexabilityStatus: models.IndexabilityIndexable,
				},
				{
					URL:                "https://example.com/long",
					StatusCode:         200,
					Title:              "This is a very very very very very very very very very long title that exceeds the limit",                                                                                                                     // Too long
					MetaDesc:           "This description is way too long and just keeps going and going and going and going and going and going and going and going and going and going and going and going and going and going and going and going.", // Too long
					H1:                 []string{"Heading"},
					Canonical:          "https://example.com/long",
					IndexabilityStatus: models.IndexabilityIndexable,
				},
			},
			expectedIssues: map[IssueType]int{
				IssueShortTitle:    1,
				IssueShortMetaDesc: 1,
				IssueLongTitle:     1,
				IssueLongMetaDesc:  1,
			},
			expectedPages: 2,
		},
		{
			name: "Multiple H1 and No Canonical",
			results: []*models.PageResult{
				{
					URL:                "https://example.com/multi-h1",
					StatusCode:         200,
					Title:              "Valid Title Length For This Test Page",
					MetaDesc:           "This is a valid meta description that is long enough to pass the minimum length check but short enough to avoid the maximum length check. It is perfect.",
					H1:                 []string{"Heading 1", "Heading 2"},
					Canonical:          "", // Missing
					IndexabilityStatus: models.IndexabilityIndexable,
				},
			},
			expectedIssues: map[IssueType]int{
				IssueMultipleH1:  1,
				IssueNoCanonical: 1,
			},
			expectedPages: 1,
		},
		{
			name: "Broken Link (404)",
			results: []*models.PageResult{
				{
					URL:        "https://example.com/404",
					StatusCode: 404,
				},
			},
			expectedIssues: map[IssueType]int{
				IssueBrokenLink: 1,
			},
			expectedPages: 1,
		},
		{
			name: "Redirect Chain",
			results: []*models.PageResult{
				{
					URL:           "https://example.com/redirect",
					StatusCode:    200,
					RedirectChain: []string{"https://example.com/start", "https://example.com/end"},
					// Even if other fields are valid, redirect chain is an issue
					Title:              "Valid Title",
					IndexabilityStatus: models.IndexabilityIndexable,
				},
			},
			expectedIssues: map[IssueType]int{
				IssueRedirectChain: 1,
				// Note: Other issues might be skipped or present depending on logic,
				// but we primarily check for the redirect chain here.
				// Based on code: it continues to analyze SEO issues for redirects if status is 200.
				// "Valid Title" is short (<30), so IssueShortTitle might also trigger.
				IssueShortTitle:      1,
				IssueMissingH1:       1, // No H1 provided
				IssueMissingMetaDesc: 1, // No Meta provided
				IssueNoCanonical:     1, // No Canonical provided
			},
			expectedPages: 1,
		},
		{
			name: "Non-Indexable Page",
			results: []*models.PageResult{
				{
					URL:                "https://example.com/noindex",
					StatusCode:         200,
					Title:              "", // Missing, but should be ignored
					IndexabilityStatus: models.IndexabilityNoindex,
				},
			},
			expectedIssues: map[IssueType]int{
				// Should have NO SEO issues reported because it's noindex
			},
			expectedPages: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary := Analyze(tt.results)

			if summary.TotalPages != tt.expectedPages {
				t.Errorf("Analyze() TotalPages = %v, want %v", summary.TotalPages, tt.expectedPages)
			}

			// Check expected issue counts
			for issueType, expectedCount := range tt.expectedIssues {
				if gotCount := summary.IssuesByType[issueType]; gotCount != expectedCount {
					t.Errorf("Analyze() Issue %v count = %v, want %v", issueType, gotCount, expectedCount)
				}
			}

			// Check that we don't have unexpected issues
			// (Optional validation to ensure test cases are strict)
			for issueType, count := range summary.IssuesByType {
				if _, ok := tt.expectedIssues[issueType]; !ok && count > 0 {
					t.Errorf("Analyze() Unexpected issue %v found (count: %v)", issueType, count)
				}
			}
		})
	}
}
