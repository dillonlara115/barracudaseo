// Package psi reads Core Web Vitals from Google's PageSpeed Insights API.
package psi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/dillonlara115/barracudaseo/pkg/models"
)

// defaultEndpoint is the PageSpeed Insights v5 API.
const defaultEndpoint = "https://www.googleapis.com/pagespeedonline/v5/runPagespeed"

// Client calls the PageSpeed Insights API. An API key is optional, but the
// keyless quota is shared and small, so anything beyond a trial needs one.
type Client struct {
	APIKey   string
	Endpoint string
	HTTP     *http.Client
}

// NewClient returns a client. A PageSpeed run loads the page in a simulated
// browser, so requests routinely take 10 to 30 seconds.
func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:   apiKey,
		Endpoint: defaultEndpoint,
		HTTP:     &http.Client{Timeout: 90 * time.Second},
	}
}

// apiResponse is the part of the PageSpeed response this package uses.
type apiResponse struct {
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	LoadingExperience *loadingExperience `json:"loadingExperience"`
	LighthouseResult  *struct {
		Audits map[string]struct {
			NumericValue *float64 `json:"numericValue"`
		} `json:"audits"`
		Categories map[string]struct {
			Score *float64 `json:"score"`
		} `json:"categories"`
	} `json:"lighthouseResult"`
}

// loadingExperience is field data from the Chrome UX Report.
type loadingExperience struct {
	// OriginFallback is true when Google had too little data for this URL and
	// returned site-wide figures instead.
	OriginFallback bool `json:"origin_fallback"`
	Metrics        map[string]struct {
		Percentile *float64 `json:"percentile"`
	} `json:"metrics"`
}

// Analyze runs PageSpeed Insights for one URL. strategy is "mobile" or
// "desktop".
func (c *Client) Analyze(ctx context.Context, pageURL, strategy string) (*models.WebVitals, error) {
	params := url.Values{}
	params.Set("url", pageURL)
	params.Set("strategy", strategy)
	params.Set("category", "performance")
	if c.APIKey != "" {
		params.Set("key", c.APIKey)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.Endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build PageSpeed request: %w", err)
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("PageSpeed request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read PageSpeed response: %w", err)
	}

	var parsed apiResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse PageSpeed response (HTTP %d): %w", resp.StatusCode, err)
	}
	if parsed.Error != nil {
		return nil, fmt.Errorf("PageSpeed API error %d: %s", parsed.Error.Code, firstLine(parsed.Error.Message))
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("PageSpeed API returned HTTP %d", resp.StatusCode)
	}

	return vitalsFromResponse(&parsed, strategy), nil
}

// vitalsFromResponse converts an API response into WebVitals.
func vitalsFromResponse(parsed *apiResponse, strategy string) *models.WebVitals {
	vitals := &models.WebVitals{Strategy: strategy}

	if field := parsed.LoadingExperience; field != nil && len(field.Metrics) > 0 {
		vitals.LCPMs = field.Metrics["LARGEST_CONTENTFUL_PAINT_MS"].Percentile
		vitals.INPMs = field.Metrics["INTERACTION_TO_NEXT_PAINT"].Percentile
		// The API reports CLS multiplied by 100 (a percentile of 12 means 0.12)
		if cls := field.Metrics["CUMULATIVE_LAYOUT_SHIFT_SCORE"].Percentile; cls != nil {
			scaled := *cls / 100
			vitals.CLS = &scaled
		}
		if vitals.LCPMs != nil || vitals.INPMs != nil || vitals.CLS != nil {
			vitals.FieldSource = "page"
			if field.OriginFallback {
				vitals.FieldSource = "origin"
			}
		}
	}

	if lab := parsed.LighthouseResult; lab != nil {
		vitals.LabLCPMs = lab.Audits["largest-contentful-paint"].NumericValue
		vitals.LabCLS = lab.Audits["cumulative-layout-shift"].NumericValue
		vitals.LabTBTMs = lab.Audits["total-blocking-time"].NumericValue
		if score := lab.Categories["performance"].Score; score != nil {
			rounded := int(*score*100 + 0.5)
			vitals.PerformanceScore = &rounded
		}
	}

	return vitals
}

// firstLine returns the first line of a possibly multi-line message.
func firstLine(message string) string {
	if i := strings.IndexByte(message, '\n'); i >= 0 {
		return message[:i]
	}
	return message
}

// SelectPages picks which crawled pages to test. PageSpeed runs are slow and
// rate limited, so only a handful of pages per crawl can be measured.
//
// TODO(dillon): the baseline takes the pages most other pages link to, which
// favours the homepage and navigation targets. That says "important" but not
// "representative": ten blog posts share one template and one result. Decide
// whether to pick one page per template (by first path segment, say), or to
// prefer pages with the most Search Console clicks when that data exists.
func SelectPages(pages []*models.PageResult, limit int) []*models.PageResult {
	if limit <= 0 {
		return nil
	}

	inlinks := make(map[string]int)
	for _, p := range pages {
		for _, link := range p.InternalLinks {
			inlinks[strings.TrimRight(link, "/")]++
		}
	}

	var candidates []*models.PageResult
	for _, p := range pages {
		indexable := p.IndexabilityStatus == models.IndexabilityIndexable || p.IndexabilityStatus == ""
		if p.StatusCode == 200 && p.Error == "" && indexable && len(p.RedirectChain) == 0 {
			candidates = append(candidates, p)
		}
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		a, b := candidates[i], candidates[j]
		linksA, linksB := inlinks[strings.TrimRight(a.URL, "/")], inlinks[strings.TrimRight(b.URL, "/")]
		if linksA != linksB {
			return linksA > linksB
		}
		if len(a.URL) != len(b.URL) {
			return len(a.URL) < len(b.URL)
		}
		return a.URL < b.URL
	})

	if len(candidates) > limit {
		candidates = candidates[:limit]
	}
	return candidates
}

// Enrich measures the given pages and stores the result on each. A page whose
// test fails gets a WebVitals carrying the error, so one failure does not hide
// the others. It returns how many pages were measured successfully.
func (c *Client) Enrich(ctx context.Context, pages []*models.PageResult, strategy string, concurrency int) int {
	if concurrency < 1 {
		concurrency = 1
	}

	var (
		wg        sync.WaitGroup
		mu        sync.Mutex
		succeeded int
	)
	slots := make(chan struct{}, concurrency)

	for _, p := range pages {
		wg.Add(1)
		go func(p *models.PageResult) {
			defer wg.Done()
			slots <- struct{}{}
			defer func() { <-slots }()

			vitals, err := c.Analyze(ctx, p.URL, strategy)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				p.WebVitals = &models.WebVitals{Strategy: strategy, Error: err.Error()}
				return
			}
			p.WebVitals = vitals
			succeeded++
		}(p)
	}

	wg.Wait()
	return succeeded
}
