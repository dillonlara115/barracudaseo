package crawler

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dillonlara115/barracuda/internal/utils"
	"github.com/dillonlara115/barracuda/pkg/models"
)

// Parser extracts SEO data from HTML content
type Parser struct {
	baseURL string
	domain  string
}

// NewParser creates a new Parser instance
func NewParser(baseURL string) (*Parser, error) {
	domain, err := utils.ExtractDomain(baseURL)
	if err != nil {
		return nil, err
	}

	return &Parser{
		baseURL: baseURL,
		domain:  domain,
	}, nil
}

// Parse extracts SEO data from HTML content
func (p *Parser) Parse(htmlContent []byte) (*models.PageResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(htmlContent)))
	if err != nil {
		return nil, err
	}

	result := &models.PageResult{
		URL:           p.baseURL,
		H1:            make([]string, 0),
		H2:            make([]string, 0),
		H3:            make([]string, 0),
		H4:            make([]string, 0),
		H5:            make([]string, 0),
		H6:            make([]string, 0),
		InternalLinks: make([]string, 0),
		ExternalLinks: make([]string, 0),
		Images:        make([]models.Image, 0),
	}
	
	// Log HTML content size for debugging
	utils.Debug("Parsing HTML", 
		utils.NewField("url", p.baseURL),
		utils.NewField("html_size", len(htmlContent)),
		utils.NewField("h1_count_in_html", doc.Find("h1").Length()),
		utils.NewField("link_count_in_html", doc.Find("a[href]").Length()))

	// Extract title
	result.Title = strings.TrimSpace(doc.Find("title").First().Text())

	// Extract meta description
	doc.Find("meta[name='description']").Each(func(i int, s *goquery.Selection) {
		if content, exists := s.Attr("content"); exists {
			result.MetaDesc = strings.TrimSpace(content)
		}
	})

	// Extract canonical link
	doc.Find("link[rel='canonical']").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			result.Canonical = strings.TrimSpace(href)
		}
	})

	// Extract headings
	// Helper function to extract clean text from heading elements
	// Handles nested elements (spans, divs, etc.) and normalizes whitespace
	extractHeadingText := func(s *goquery.Selection) string {
		// Remove script and style elements that might contain text we don't want
		// Clone first to avoid modifying the original document
		clone := s.Clone()
		clone.Find("script, style").Remove()
		
		// Get text content - goquery's Text() method extracts text from all nested elements
		// This will get text from <span>, <div>, etc. nested inside the heading
		text := clone.Text()
		
		// Normalize whitespace: strings.Fields splits on any whitespace (spaces, tabs, newlines)
		// and strings.Join combines them with single spaces
		// This handles <br> tags, multiple spaces, tabs, etc.
		text = strings.Join(strings.Fields(text), " ")
		
		// Final trim to remove leading/trailing spaces
		return strings.TrimSpace(text)
	}
	
	doc.Find("h1").Each(func(i int, s *goquery.Selection) {
		// Get raw HTML for debugging
		rawHTML, _ := s.Html()
		
		// Try multiple extraction methods
		var text string
		
		// Method 1: Direct text extraction (should work for most cases)
		directText := s.Text()
		text = strings.TrimSpace(directText)
		
		// Method 2: If empty, try normalizing whitespace
		if text == "" && len(directText) > 0 {
			text = strings.Join(strings.Fields(directText), " ")
			text = strings.TrimSpace(text)
		}
		
		// Method 3: If still empty, try the helper function with cloning
		if text == "" {
			text = extractHeadingText(s)
		}
		
		// Method 4: Last resort - try getting text from all child elements
		if text == "" {
			var parts []string
			s.Contents().Each(func(j int, child *goquery.Selection) {
				childText := strings.TrimSpace(child.Text())
				if childText != "" {
					parts = append(parts, childText)
				}
			})
			if len(parts) > 0 {
				text = strings.Join(parts, " ")
				text = strings.Join(strings.Fields(text), " ")
				text = strings.TrimSpace(text)
			}
		}
		
		if text != "" {
			result.H1 = append(result.H1, text)
		} else {
			// Log when H1 tag exists but text extraction returns empty
			utils.Debug("H1 tag found but text extraction returned empty", 
				utils.NewField("url", p.baseURL),
				utils.NewField("h1_html", rawHTML),
				utils.NewField("h1_count", doc.Find("h1").Length()),
				utils.NewField("direct_text", directText),
				utils.NewField("direct_text_len", len(directText)),
				utils.NewField("direct_text_bytes", []byte(directText)))
		}
	})
	doc.Find("h2").Each(func(i int, s *goquery.Selection) {
		text := extractHeadingText(s)
		if text != "" {
			result.H2 = append(result.H2, text)
		}
	})
	doc.Find("h3").Each(func(i int, s *goquery.Selection) {
		text := extractHeadingText(s)
		if text != "" {
			result.H3 = append(result.H3, text)
		}
	})
	doc.Find("h4").Each(func(i int, s *goquery.Selection) {
		text := extractHeadingText(s)
		if text != "" {
			result.H4 = append(result.H4, text)
		}
	})
	doc.Find("h5").Each(func(i int, s *goquery.Selection) {
		text := extractHeadingText(s)
		if text != "" {
			result.H5 = append(result.H5, text)
		}
	})
	doc.Find("h6").Each(func(i int, s *goquery.Selection) {
		text := extractHeadingText(s)
		if text != "" {
			result.H6 = append(result.H6, text)
		}
	})

	// Extract links
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		// Resolve relative URLs
		resolvedURL, err := utils.ResolveURL(p.baseURL, href)
		if err != nil {
			return
		}

		// Normalize URL
		normalizedURL, err := utils.NormalizeURL(resolvedURL)
		if err != nil {
			return
		}

		// Skip fragments, javascript:, mailto:, etc.
		u, err := url.Parse(normalizedURL)
		if err != nil {
			return
		}

		if u.Scheme != "http" && u.Scheme != "https" {
			return
		}

		// Categorize as internal or external
		if utils.IsSameDomain(normalizedURL, p.baseURL) {
			// Avoid duplicates
			for _, existing := range result.InternalLinks {
				if existing == normalizedURL {
					return
				}
			}
			result.InternalLinks = append(result.InternalLinks, normalizedURL)
		} else {
			// Avoid duplicates
			for _, existing := range result.ExternalLinks {
				if existing == normalizedURL {
					return
				}
			}
			result.ExternalLinks = append(result.ExternalLinks, normalizedURL)
		}
	})

	// Extract images
	imageCount := 0
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if !exists {
			return
		}

		// Resolve relative URLs
		resolvedURL, err := utils.ResolveURL(p.baseURL, src)
		if err != nil {
			utils.Debug("Failed to resolve image URL", utils.NewField("src", src), utils.NewField("error", err.Error()))
			return
		}

		// Normalize URL
		normalizedURL, err := utils.NormalizeURL(resolvedURL)
		if err != nil {
			utils.Debug("Failed to normalize image URL", utils.NewField("resolved_url", resolvedURL), utils.NewField("error", err.Error()))
			return
		}

		// Skip data URIs and non-HTTP schemes
		u, err := url.Parse(normalizedURL)
		if err != nil {
			return
		}

		if u.Scheme != "http" && u.Scheme != "https" {
			return
		}

		// Get alt text
		alt := s.AttrOr("alt", "")

		// Avoid duplicates
		for _, existing := range result.Images {
			if existing.URL == normalizedURL {
				return
			}
		}

		result.Images = append(result.Images, models.Image{
			URL: normalizedURL,
			Alt: alt,
		})
		imageCount++
	})
	
	if imageCount > 0 {
		utils.Debug("Extracted images from page", 
			utils.NewField("url", p.baseURL),
			utils.NewField("image_count", imageCount))
	}

	return result, nil
}

// ExtractLinks extracts all links from HTML content and returns them as a slice
func (p *Parser) ExtractLinks(htmlContent []byte) ([]string, error) {
	result, err := p.Parse(htmlContent)
	if err != nil {
		return nil, err
	}

	links := make([]string, 0, len(result.InternalLinks)+len(result.ExternalLinks))
	links = append(links, result.InternalLinks...)
	links = append(links, result.ExternalLinks...)

	return links, nil
}

