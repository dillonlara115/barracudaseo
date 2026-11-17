// Package render loads pages in headless Chrome so the DOM after JavaScript has
// run can be compared with the HTML the server sent.
package render

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"

	"github.com/dillonlara115/barracudaseo/internal/crawler"
	"github.com/dillonlara115/barracudaseo/pkg/models"
)

const (
	// pageTimeout bounds one page load, including the settle delay.
	pageTimeout = 30 * time.Second

	// settleDelay gives client-side code time to fetch data and render after
	// the document is ready.
	settleDelay = 2 * time.Second
)

// Renderer drives one headless Chrome process, opening a tab per page. It
// needs Chrome or Chromium installed on the machine.
type Renderer struct {
	browserCtx    context.Context
	cancelBrowser context.CancelFunc
	cancelAlloc   context.CancelFunc
}

// NewRenderer starts a headless browser. It fails if Chrome cannot be started,
// for example because it is not installed. Close must be called to stop the
// browser process.
func NewRenderer(userAgent string) (*Renderer, error) {
	options := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoSandbox,
		chromedp.DisableGPU,
		// Containers give /dev/shm very little space; without this Chrome
		// crashes on larger pages.
		chromedp.Flag("disable-dev-shm-usage", true),
	)
	if userAgent != "" {
		options = append(options, chromedp.UserAgent(userAgent))
	}
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), options...)
	browserCtx, cancelBrowser := chromedp.NewContext(allocCtx)

	// The first call on a context starts the browser and ties its lifetime to
	// that context, so it must not carry a timeout.
	if err := chromedp.Do(browserCtx); err != nil {
		cancelBrowser()
		cancelAlloc()
		return nil, fmt.Errorf("failed to start headless Chrome: %w", err)
	}

	return &Renderer{browserCtx: browserCtx, cancelBrowser: cancelBrowser, cancelAlloc: cancelAlloc}, nil
}

// Close stops the browser.
func (r *Renderer) Close() {
	r.cancelBrowser()
	r.cancelAlloc()
}

// RenderHTML loads a URL in a new tab, waits for JavaScript to settle, and
// returns the resulting document's HTML.
func (r *Renderer) RenderHTML(pageURL string) (string, error) {
	// A context derived from one that holds a browser opens a tab on its first
	// call. Open the tab without a timeout, then bound the page load itself.
	tabCtx, closeTab := chromedp.NewContext(r.browserCtx)
	defer closeTab()
	if err := chromedp.Do(tabCtx); err != nil {
		return "", fmt.Errorf("failed to open a browser tab: %w", err)
	}

	ctx, cancel := context.WithTimeout(tabCtx, pageTimeout)
	defer cancel()

	err := chromedp.Do(ctx,
		chromedp.Navigate(pageURL),
		chromedp.WaitReady("body"),
		chromedp.Sleep(settleDelay),
	)
	if err != nil {
		return "", fmt.Errorf("failed to load %s: %w", pageURL, err)
	}

	html, err := chromedp.Run(ctx, chromedp.OuterHTML("html"))
	if err != nil {
		return "", fmt.Errorf("failed to read rendered HTML of %s: %w", pageURL, err)
	}
	return html, nil
}

// Compare renders a page and records how it differs from the raw HTML the
// crawler already parsed.
func (r *Renderer) Compare(raw *models.PageResult) *models.RenderComparison {
	comparison := &models.RenderComparison{
		RawWordCount:  raw.WordCount,
		RawTitle:      raw.Title,
		RawCanonical:  raw.Canonical,
		RawMetaRobots: raw.MetaRobots,
	}

	html, err := r.RenderHTML(raw.URL)
	if err != nil {
		comparison.Error = err.Error()
		return comparison
	}

	parser, err := crawler.NewParser(raw.URL)
	if err != nil {
		comparison.Error = err.Error()
		return comparison
	}
	rendered, err := parser.Parse([]byte(html))
	if err != nil {
		comparison.Error = err.Error()
		return comparison
	}

	comparison.RenderedWordCount = rendered.WordCount
	comparison.RenderedTitle = rendered.Title
	comparison.RenderedCanonical = rendered.Canonical
	comparison.RenderedMetaRobots = rendered.MetaRobots

	inRaw := make(map[string]bool, len(raw.InternalLinks))
	for _, link := range raw.InternalLinks {
		inRaw[strings.TrimRight(link, "/")] = true
	}
	for _, link := range rendered.InternalLinks {
		if !inRaw[strings.TrimRight(link, "/")] {
			comparison.LinksOnlyRendered = append(comparison.LinksOnlyRendered, link)
		}
	}

	return comparison
}

// Enrich renders the given pages and stores the comparison on each. It returns
// how many pages rendered successfully.
func (r *Renderer) Enrich(pages []*models.PageResult, concurrency int) int {
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

			comparison := r.Compare(p)
			mu.Lock()
			defer mu.Unlock()
			p.Rendered = comparison
			if comparison.Error == "" {
				succeeded++
			}
		}(p)
	}

	wg.Wait()
	return succeeded
}

// SelectPages picks which pages to render: every page that looks client
// rendered (where rendering changes the most), then the given priority pages,
// up to limit.
func SelectPages(all, priority []*models.PageResult, limit int) []*models.PageResult {
	if limit <= 0 {
		return nil
	}

	seen := make(map[*models.PageResult]bool)
	var selected []*models.PageResult
	add := func(p *models.PageResult) {
		if len(selected) < limit && !seen[p] && p.StatusCode == 200 && p.Error == "" {
			seen[p] = true
			selected = append(selected, p)
		}
	}

	for _, p := range all {
		if p.JS != nil && p.JS.ClientRendered {
			add(p)
		}
	}
	for _, p := range priority {
		add(p)
	}
	return selected
}
