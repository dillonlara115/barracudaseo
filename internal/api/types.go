package api

import "github.com/dillonlara115/barracudaseo/pkg/models"

// CreateCrawlRequest represents a crawl ingestion request
type CreateCrawlRequest struct {
	ProjectID string               `json:"project_id"`
	Pages     []*models.PageResult `json:"pages"`
	Source    string               `json:"source,omitempty"` // "cli", "web", "schedule"
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

// TriggerCrawlRequest represents a request to trigger a new crawl
type TriggerCrawlRequest struct {
	URL              string `json:"url"`                // Starting URL to crawl
	MaxDepth         int    `json:"max_depth"`          // Maximum crawl depth (default: 3)
	MaxPages         int    `json:"max_pages"`          // Maximum pages to crawl (default: 1000)
	Workers          int    `json:"workers"`            // Number of concurrent workers (default: 10)
	RespectRobots    *bool  `json:"respect_robots"`     // Respect robots.txt (default: true)
	ParseSitemap     *bool  `json:"parse_sitemap"`      // Parse sitemap.xml (default: false)
	CrawlSitemapOnly *bool  `json:"crawl_sitemap_only"` // Crawl only sitemap URLs, no link discovery—like indexed pages (default: false, requires parse_sitemap)
}
