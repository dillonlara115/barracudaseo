package utils

import (
	"errors"
	"net/url"
	"strings"
)

var (
	ErrInvalidURL      = errors.New("invalid URL")
	ErrEmptyStartURL   = errors.New("start URL cannot be empty")
	ErrInvalidMaxDepth = errors.New("max depth must be non-negative")
	ErrInvalidMaxPages = errors.New("max pages must be at least 1")
	ErrInvalidWorkers  = errors.New("workers must be at least 1")
	ErrInvalidExportFormat = errors.New("export format must be 'csv' or 'json'")
)

// NormalizeURL normalizes a URL by removing fragments and trailing slashes
func NormalizeURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", ErrInvalidURL
	}

	// Remove fragment
	u.Fragment = ""
	// Remove trailing slash unless it's root
	normalized := u.String()
	if normalized != u.Scheme+"://"+u.Host+"/" && strings.HasSuffix(normalized, "/") {
		normalized = normalized[:len(normalized)-1]
	}

	return normalized, nil
}

// ExtractDomain extracts the domain from a URL
func ExtractDomain(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", ErrInvalidURL
	}
	return u.Host, nil
}

// IsSameDomain checks if two URLs belong to the same domain
// It handles www vs non-www by treating them as the same domain
func IsSameDomain(url1, url2 string) bool {
	domain1, err1 := ExtractDomain(url1)
	domain2, err2 := ExtractDomain(url2)
	if err1 != nil || err2 != nil {
		return false
	}
	
	// Exact match
	if domain1 == domain2 {
		return true
	}
	
	// Handle www vs non-www: treat as same domain
	// Remove www. prefix for comparison
	normalizeDomain := func(domain string) string {
		if strings.HasPrefix(domain, "www.") {
			return domain[4:]
		}
		return domain
	}
	
	normalized1 := normalizeDomain(domain1)
	normalized2 := normalizeDomain(domain2)
	
	return normalized1 == normalized2
}

// ResolveURL resolves a relative URL against a base URL
func ResolveURL(baseURL, relativeURL string) (string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	rel, err := url.Parse(relativeURL)
	if err != nil {
		return "", err
	}

	resolved := base.ResolveReference(rel)
	normalized, err := NormalizeURL(resolved.String())
	if err != nil {
		return "", err
	}

	return normalized, nil
}

// IsValidURL checks if a string is a valid URL
func IsValidURL(rawURL string) bool {
	_, err := url.Parse(rawURL)
	return err == nil
}

// IsImageURL checks if a URL points to an image file based on its extension
// This function checks both the path extension and common image URL patterns
func IsImageURL(rawURL string) bool {
	if rawURL == "" {
		return false
	}
	
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	
	// Get the path (without query parameters or fragments)
	path := u.Path
	if path == "" {
		return false
	}
	
	// Remove trailing slash
	path = strings.TrimSuffix(path, "/")
	if path == "" {
		return false
	}
	
	// Find the last dot in the path
	lastDot := strings.LastIndex(path, ".")
	if lastDot == -1 {
		// No extension found
		return false
	}
	
	// Make sure the dot is not at the end
	if lastDot >= len(path)-1 {
		return false
	}
	
	// Get file extension (lowercase) - everything after the last dot
	ext := strings.ToLower(path[lastDot:])
	if ext == "" {
		return false
	}
	
	// Common image extensions (must match exactly)
	imageExtensions := []string{
		".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg",
		".bmp", ".ico", ".tiff", ".tif", ".avif", ".heic", ".heif",
		".jp2", ".j2k", ".jpx", ".jpf", ".jpm", ".mj2", ".mjp2",
	}
	
	for _, imgExt := range imageExtensions {
		if ext == imgExt {
			return true
		}
	}
	
	return false
}

