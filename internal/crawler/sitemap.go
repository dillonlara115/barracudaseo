package crawler

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"

	"github.com/dillonlara115/barracudaseo/internal/utils"
)

// SitemapIndex represents a sitemap index file (no namespace)
type SitemapIndex struct {
	XMLName  xml.Name  `xml:"sitemapindex"`
	Sitemaps []Sitemap `xml:"sitemap"`
}

// SitemapIndexNS is sitemap index with namespace
type SitemapIndexNS struct {
	XMLName  xml.Name    `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 sitemapindex"`
	Sitemaps []SitemapNS `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 sitemap"`
}

// Sitemap represents a single sitemap entry
type Sitemap struct {
	Loc string `xml:"loc"`
}

// SitemapNS is Sitemap with namespace
type SitemapNS struct {
	Loc string `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 loc"`
}

// URLSet represents a sitemap URL set (no namespace)
type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	URLs    []URL    `xml:"url"`
}

// URLSetNS is URLSet with namespace - for sitemaps that declare xmlns
type URLSetNS struct {
	XMLName xml.Name `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 urlset"`
	URLs    []URLNS  `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 url"`
}

// URL represents a single URL in a sitemap
type URL struct {
	Loc string `xml:"loc"`
}

// URLNS is URL with namespace
type URLNS struct {
	Loc string `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 loc"`
}

// SitemapParser parses sitemap.xml files
type SitemapParser struct {
	fetcher *Fetcher
}

// NewSitemapParser creates a new SitemapParser instance
func NewSitemapParser(fetcher *Fetcher) *SitemapParser {
	return &SitemapParser{
		fetcher: fetcher,
	}
}

// ParseSitemap fetches and parses a sitemap URL, returning all URLs found
func (s *SitemapParser) ParseSitemap(sitemapURL string) ([]string, error) {
	body, statusCode, err := s.fetcher.FetchRaw(sitemapURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sitemap: %w", err)
	}

	if statusCode != 200 {
		return nil, fmt.Errorf("sitemap returned HTTP %d", statusCode)
	}

	// Try parsing as sitemap index first (with namespace - most sitemaps use xmlns)
	var indexNS SitemapIndexNS
	if err := xml.Unmarshal(body, &indexNS); err == nil && len(indexNS.Sitemaps) > 0 {
		return s.collectFromIndexNS(indexNS.Sitemaps)
	}

	// Try sitemap index without namespace
	var index SitemapIndex
	if err := xml.Unmarshal(body, &index); err == nil && len(index.Sitemaps) > 0 {
		return s.collectFromIndex(index.Sitemaps)
	}

	// Try parsing as URL set with namespace first (barracudaseo.com and most sites use this)
	var urlSetNS URLSetNS
	if err := xml.Unmarshal(body, &urlSetNS); err == nil && len(urlSetNS.URLs) > 0 {
		return s.extractURLsNS(urlSetNS.URLs)
	}

	// Fallback: URL set without namespace
	var urlSet URLSet
	if err := xml.Unmarshal(body, &urlSet); err != nil {
		return nil, fmt.Errorf("failed to parse sitemap XML: %w", err)
	}
	return s.extractURLs(urlSet.URLs)
}

func (s *SitemapParser) collectFromIndex(sitemaps []Sitemap) ([]string, error) {
	seen := make(map[string]bool)
	urls := make([]string, 0)
	for _, sm := range sitemaps {
		subURLs, err := s.ParseSitemap(strings.TrimSpace(sm.Loc))
		if err != nil {
			utils.Debug("Failed to parse sub-sitemap", utils.NewField("url", sm.Loc), utils.NewField("error", err.Error()))
			continue
		}
		for _, u := range subURLs {
			if !seen[u] {
				seen[u] = true
				urls = append(urls, u)
			}
		}
	}
	return urls, nil
}

func (s *SitemapParser) collectFromIndexNS(sitemaps []SitemapNS) ([]string, error) {
	seen := make(map[string]bool)
	urls := make([]string, 0)
	for _, sm := range sitemaps {
		subURLs, err := s.ParseSitemap(strings.TrimSpace(sm.Loc))
		if err != nil {
			utils.Debug("Failed to parse sub-sitemap", utils.NewField("url", sm.Loc), utils.NewField("error", err.Error()))
			continue
		}
		for _, u := range subURLs {
			if !seen[u] {
				seen[u] = true
				urls = append(urls, u)
			}
		}
	}
	return urls, nil
}

func (s *SitemapParser) extractURLs(urlList []URL) ([]string, error) {
	seen := make(map[string]bool)
	urls := make([]string, 0, len(urlList))
	for _, u := range urlList {
		raw := strings.TrimSpace(u.Loc)
		normalized, err := utils.NormalizeURL(raw)
		if err != nil {
			utils.Debug("Invalid URL in sitemap", utils.NewField("url", raw), utils.NewField("error", err.Error()))
			continue
		}
		if seen[normalized] {
			utils.Debug("Skipping duplicate URL in sitemap", utils.NewField("original", raw), utils.NewField("normalized", normalized))
			continue
		}
		seen[normalized] = true
		urls = append(urls, normalized)
	}
	return urls, nil
}

func (s *SitemapParser) extractURLsNS(urlList []URLNS) ([]string, error) {
	seen := make(map[string]bool)
	urls := make([]string, 0, len(urlList))
	for _, u := range urlList {
		raw := strings.TrimSpace(u.Loc)
		normalized, err := utils.NormalizeURL(raw)
		if err != nil {
			utils.Debug("Invalid URL in sitemap", utils.NewField("url", raw), utils.NewField("error", err.Error()))
			continue
		}
		if seen[normalized] {
			utils.Debug("Skipping duplicate URL in sitemap", utils.NewField("original", raw), utils.NewField("normalized", normalized))
			continue
		}
		seen[normalized] = true
		urls = append(urls, normalized)
	}
	return urls, nil
}

// DiscoverSitemapURL attempts to discover sitemap.xml URL from a base URL
func (s *SitemapParser) DiscoverSitemapURL(baseURL string) string {
	u, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}

	defaultSitemapURL := fmt.Sprintf("%s://%s/sitemap.xml", u.Scheme, u.Host)

	candidates := s.discoverSitemapCandidates(u, defaultSitemapURL)
	for _, candidate := range candidates {
		body, statusCode, err := s.fetcher.FetchRaw(candidate)
		if err != nil || statusCode != 200 {
			continue
		}
		if looksLikeSitemapXML(body) {
			return candidate
		}
	}

	return defaultSitemapURL
}

func (s *SitemapParser) discoverSitemapCandidates(base *url.URL, defaultSitemapURL string) []string {
	candidates := make([]string, 0, 8)
	seen := make(map[string]bool)

	addCandidate := func(raw string) {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			return
		}

		parsed, err := base.Parse(raw)
		if err != nil {
			return
		}

		candidate := parsed.String()
		if seen[candidate] {
			return
		}
		seen[candidate] = true
		candidates = append(candidates, candidate)
	}

	robotsURL := fmt.Sprintf("%s://%s/robots.txt", base.Scheme, base.Host)
	if body, statusCode, err := s.fetcher.FetchRaw(robotsURL); err == nil && statusCode == 200 {
		for _, line := range strings.Split(string(body), "\n") {
			line = strings.TrimSpace(line)
			if len(line) < len("sitemap:") || !strings.EqualFold(line[:len("sitemap:")], "sitemap:") {
				continue
			}
			addCandidate(strings.TrimSpace(line[len("sitemap:"):]))
		}
	}

	addCandidate(defaultSitemapURL)
	addCandidate(fmt.Sprintf("%s://%s/sitemap_index.xml", base.Scheme, base.Host))
	addCandidate(fmt.Sprintf("%s://%s/wp-sitemap.xml", base.Scheme, base.Host))

	return candidates
}

func looksLikeSitemapXML(body []byte) bool {
	content := strings.ToLower(strings.TrimSpace(string(body)))
	return strings.Contains(content, "<urlset") || strings.Contains(content, "<sitemapindex")
}
