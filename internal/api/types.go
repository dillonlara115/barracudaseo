package api

import "github.com/dillonlara115/barracuda/pkg/models"

// CreateCrawlRequest represents a crawl ingestion request
type CreateCrawlRequest struct {
	ProjectID string              `json:"project_id"`
	Pages     []*models.PageResult `json:"pages"`
	Source    string              `json:"source,omitempty"` // "cli", "web", "schedule"
}

// CreateCrawlResponse represents the response after creating a crawl
type CreateCrawlResponse struct {
	CrawlID     string `json:"crawl_id"`
	ProjectID   string `json:"project_id"`
	TotalPages  int    `json:"total_pages"`
	TotalIssues int    `json:"total_issues"`
	Status      string `json:"status"`
}

// CreateProjectRequest represents a project creation request
type CreateProjectRequest struct {
	Name     string                 `json:"name"`
	Domain   string                 `json:"domain"`
	Settings map[string]interface{} `json:"settings,omitempty"`
}

