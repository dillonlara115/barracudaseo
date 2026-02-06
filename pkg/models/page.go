package models

import (
	"strings"
	"time"
)

// IndexabilityStatus represents whether a page is indexable by search engines
type IndexabilityStatus string

const (
	IndexabilityIndexable IndexabilityStatus = "indexable" // Page can be indexed
	IndexabilityNoindex   IndexabilityStatus = "noindex"   // Page has noindex directive
	IndexabilityBlocked   IndexabilityStatus = "blocked"   // Page blocked by robots.txt
)

// PageResult represents the SEO data extracted from a crawled page
type PageResult struct {
	URL                string             `json:"url"`
	StatusCode         int                `json:"status_code"`
	ResponseTime       int64              `json:"response_time_ms"` // Duration in milliseconds
	Title              string             `json:"title"`
	MetaDesc           string             `json:"meta_description"`
	Canonical          string             `json:"canonical"`
	H1                 []string           `json:"h1"`
	H2                 []string           `json:"h2"`
	H3                 []string           `json:"h3"`
	H4                 []string           `json:"h4"`
	H5                 []string           `json:"h5"`
	H6                 []string           `json:"h6"`
	InternalLinks      []string           `json:"internal_links"`
	ExternalLinks      []string           `json:"external_links"`
	Images             []Image            `json:"images,omitempty"`
	RedirectChain      []string           `json:"redirect_chain,omitempty"`
	Error              string             `json:"error,omitempty"`
	XRobotsTag         string             `json:"x_robots_tag,omitempty"` // HTTP X-Robots-Tag header value
	MetaRobots         string             `json:"meta_robots,omitempty"`  // HTML meta robots tag value
	IndexabilityStatus IndexabilityStatus `json:"indexability_status,omitempty"`
	CrawledAt          time.Time          `json:"crawled_at"`
}

// Image represents an image found on a page
type Image struct {
	URL string `json:"url"`
	Alt string `json:"alt,omitempty"`
}

// DetermineIndexabilityStatus determines the indexability status based on x-robots-tag, meta robots, and robots.txt blocking
// This should be called after parsing both HTTP headers and HTML content
func (p *PageResult) DetermineIndexabilityStatus(isBlockedByRobots bool) {
	// Check if blocked by robots.txt first (highest priority)
	if isBlockedByRobots {
		p.IndexabilityStatus = IndexabilityBlocked
		return
	}

	// Check x-robots-tag header
	xRobotsLower := strings.ToLower(p.XRobotsTag)
	if strings.Contains(xRobotsLower, "noindex") {
		p.IndexabilityStatus = IndexabilityNoindex
		return
	}

	// Check meta robots tag
	metaRobotsLower := strings.ToLower(p.MetaRobots)
	if strings.Contains(metaRobotsLower, "noindex") {
		p.IndexabilityStatus = IndexabilityNoindex
		return
	}

	// Default to indexable if no noindex directives found
	p.IndexabilityStatus = IndexabilityIndexable
}
