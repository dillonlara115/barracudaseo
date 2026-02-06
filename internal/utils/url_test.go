package utils

import (
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "Simple URL",
			input:    "https://example.com",
			expected: "https://example.com",
			wantErr:  false,
		},
		{
			name:     "URL with trailing slash",
			input:    "https://example.com/",
			expected: "https://example.com",
			wantErr:  false,
		},
		{
			name:     "URL with fragment",
			input:    "https://example.com#section",
			expected: "https://example.com",
			wantErr:  false,
		},
		{
			name:     "URL with path and trailing slash",
			input:    "https://example.com/blog/",
			expected: "https://example.com/blog",
			wantErr:  false,
		},
		{
			name:     "URL with query params",
			input:    "https://example.com?q=test",
			expected: "https://example.com?q=test",
			wantErr:  false,
		},
		{
			name:     "Invalid URL",
			input:    "://invalid",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("NormalizeURL() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsSameDomain(t *testing.T) {
	tests := []struct {
		name string
		url1 string
		url2 string
		want bool
	}{
		{
			name: "Same domain exact match",
			url1: "https://example.com",
			url2: "https://example.com",
			want: true,
		},
		{
			name: "Different protocols",
			url1: "https://example.com",
			url2: "http://example.com",
			want: true, // ExtractDomain only checks Host
		},
		{
			name: "www and non-www",
			url1: "https://www.example.com",
			url2: "https://example.com",
			want: true,
		},
		{
			name: "Different subdomains",
			url1: "https://blog.example.com",
			url2: "https://shop.example.com",
			want: false,
		},
		{
			name: "Different domains",
			url1: "https://example.com",
			url2: "https://google.com",
			want: false,
		},
		{
			name: "Domain and subdomain",
			url1: "https://example.com",
			url2: "https://sub.example.com",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSameDomain(tt.url1, tt.url2); got != tt.want {
				t.Errorf("IsSameDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsImageURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "JPG file",
			input: "https://example.com/image.jpg",
			want:  true,
		},
		{
			name:  "PNG file",
			input: "https://example.com/assets/logo.png",
			want:  true,
		},
		{
			name:  "WebP file",
			input: "https://example.com/photo.webp",
			want:  true,
		},
		{
			name:  "Uppercase extension",
			input: "https://example.com/PHOTO.JPG",
			want:  true, // Extension check is case-insensitive
		},
		{
			name:  "No extension",
			input: "https://example.com/image",
			want:  false,
		},
		{
			name:  "HTML file",
			input: "https://example.com/index.html",
			want:  false,
		},
		{
			name:  "Directory",
			input: "https://example.com/images/",
			want:  false,
		},
		{
			name:  "Empty string",
			input: "",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsImageURL(tt.input); got != tt.want {
				t.Errorf("IsImageURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
