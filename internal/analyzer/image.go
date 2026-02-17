package analyzer

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/dillonlara115/barracudaseo/internal/utils"
	"github.com/dillonlara115/barracudaseo/pkg/models"
)

const (
	// MaxImageSizeKB is the threshold for considering images as "large"
	MaxImageSizeKB = 100

	// imageAnalysisWorkers is the number of concurrent image size checks.
	// Keeps crawls fast without overwhelming target servers.
	imageAnalysisWorkers = 16
)

// ImageSizeInfo contains image size information
type ImageSizeInfo struct {
	URL    string
	SizeKB int64
	Size   int64
	Error  error
}

// CheckImageSize fetches image size using HEAD request
func CheckImageSize(imageURL string, timeout time.Duration) ImageSizeInfo {
	info := ImageSizeInfo{
		URL: imageURL,
	}

	client := &http.Client{
		Timeout: timeout,
	}

	// Try HEAD first (more efficient)
	req, err := http.NewRequest("HEAD", imageURL, nil)
	if err != nil {
		info.Error = err
		return info
	}

	resp, err := client.Do(req)
	if err != nil {
		info.Error = err
		return info
	}
	defer resp.Body.Close()

	// Check Content-Length header
	if resp.StatusCode == 200 {
		contentLength := resp.ContentLength
		if contentLength > 0 {
			info.Size = contentLength
			info.SizeKB = contentLength / 1024
			return info
		}
	}

	// If HEAD doesn't provide size (Content-Length is 0 or -1), try GET but limit body read
	if resp.StatusCode == 200 && resp.ContentLength <= 0 {
		getReq, err := http.NewRequest("GET", imageURL, nil)
		if err != nil {
			info.Error = err
			return info
		}

		// Only read first 1MB to get size
		getResp, err := client.Do(getReq)
		if err != nil {
			info.Error = err
			return info
		}
		defer getResp.Body.Close()

		if getResp.StatusCode == 200 {
			// Try to get size from Content-Length
			if getResp.ContentLength > 0 {
				info.Size = getResp.ContentLength
				info.SizeKB = getResp.ContentLength / 1024
			} else {
				// Read limited bytes to estimate
				limitedReader := io.LimitReader(getResp.Body, 1024*1024) // 1MB max
				bytesRead, _ := io.Copy(io.Discard, limitedReader)
				if bytesRead >= 1024*1024 {
					// If we hit the limit, it's at least 1MB
					info.Size = 1024 * 1024
					info.SizeKB = 1024
				} else {
					info.Size = bytesRead
					info.SizeKB = bytesRead / 1024
				}
			}
		} else {
			// GET request failed, return error
			info.Error = fmt.Errorf("GET request returned status %d", getResp.StatusCode)
		}
	} else if resp.StatusCode != 200 {
		// HEAD request failed
		info.Error = fmt.Errorf("HEAD request returned status %d", resp.StatusCode)
	}

	return info
}

// imageRef holds a page URL and image for building issues after parallel fetch
type imageRef struct {
	pageURL string
	img     models.Image
}

// AnalyzeImages analyzes images from page results and detects issues.
// Image size checks run in parallel (imageAnalysisWorkers) for faster analysis.
func AnalyzeImages(results []*models.PageResult, timeout time.Duration) []Issue {
	var issues []Issue
	var sizeCheckRefs []imageRef
	urlsToFetch := make(map[string]bool)
	totalImages := 0
	imagesWithoutAlt := 0

	// First pass: collect missing alt issues, and refs/URLs for size checks
	for _, result := range results {
		if utils.IsImageURL(result.URL) {
			continue
		}
		if result.StatusCode != 200 || result.Error != "" {
			continue
		}
		if len(result.Images) == 0 {
			continue
		}

		utils.Debug("Analyzing images from page",
			utils.NewField("url", result.URL),
			utils.NewField("image_count", len(result.Images)))

		for _, img := range result.Images {
			totalImages++

			if img.Alt == "" {
				imagesWithoutAlt++
				issues = append(issues, Issue{
					Type:           IssueMissingImageAlt,
					Severity:       "warning",
					URL:            result.URL,
					Message:        fmt.Sprintf("Image missing alt text: %s", img.URL),
					Value:          img.URL,
					Recommendation: "Add descriptive alt text for accessibility and SEO",
				})
			}

			sizeCheckRefs = append(sizeCheckRefs, imageRef{result.URL, img})
			urlsToFetch[img.URL] = true
		}
	}

	// Parallel fetch: populate cache for all unique image URLs
	imageSizeCache := fetchImageSizesInParallel(urlsToFetch, timeout)

	// Second pass: build large image issues from cache
	largeImages := 0
	for _, ref := range sizeCheckRefs {
		sizeInfo := imageSizeCache[ref.img.URL]
		if sizeInfo.Error == nil && sizeInfo.SizeKB > MaxImageSizeKB {
			largeImages++
			issues = append(issues, Issue{
				Type:           IssueLargeImage,
				Severity:       "warning",
				URL:            ref.pageURL,
				Message:        fmt.Sprintf("Large image detected: %s (%d KB)", ref.img.URL, sizeInfo.SizeKB),
				Value:          fmt.Sprintf("%s (%d KB)", ref.img.URL, sizeInfo.SizeKB),
				Recommendation: fmt.Sprintf("Optimize image to reduce size below %d KB", MaxImageSizeKB),
			})
		}
	}

	if totalImages > 0 {
		utils.Debug("Image analysis complete",
			utils.NewField("total_images", totalImages),
			utils.NewField("missing_alt", imagesWithoutAlt),
			utils.NewField("large_images", largeImages),
			utils.NewField("threshold_kb", MaxImageSizeKB))
	}

	return issues
}

// fetchImageSizesInParallel fetches sizes for the given URLs using a worker pool.
func fetchImageSizesInParallel(urls map[string]bool, timeout time.Duration) map[string]ImageSizeInfo {
	cache := make(map[string]ImageSizeInfo)
	var mu sync.Mutex

	work := make(chan string, len(urls))
	for url := range urls {
		work <- url
	}
	close(work)

	var wg sync.WaitGroup
	workers := imageAnalysisWorkers
	if len(urls) < workers {
		workers = len(urls)
	}
	if workers < 1 {
		return cache
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range work {
				info := CheckImageSize(url, timeout)
				mu.Lock()
				cache[url] = info
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	return cache
}
