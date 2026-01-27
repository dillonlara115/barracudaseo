package crawler

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/dillonlara115/barracudaseo/internal/graph"
	"github.com/dillonlara115/barracudaseo/internal/utils"
	"github.com/dillonlara115/barracudaseo/pkg/models"
)

// ProgressCallback is called when a page is crawled to allow real-time updates
type ProgressCallback func(page *models.PageResult, totalPages int)

// Manager orchestrates the crawling process
type Manager struct {
	config           *utils.Config
	fetcher          *Fetcher
	robotsChecker    *RobotsChecker
	sitemapParser    *SitemapParser
	linkGraph        *graph.Graph
	visited          sync.Map // map[string]bool for visited URLs
	queue            chan crawlTask
	results          []*models.PageResult
	resultsMu        sync.Mutex
	wg               sync.WaitGroup
	ctx              context.Context
	cancel           context.CancelFunc
	pending          int32 // Track pending tasks (atomic)
	queueClosed      int32 // Atomic flag to track if queue is closed
	progressCallback ProgressCallback // Optional callback for progress updates
	normalizedStartURL string // Store normalized start URL for domain comparison
}

// crawlTask represents a URL to be crawled with its depth
type crawlTask struct {
	URL   string
	Depth int
}

// NewManager creates a new Manager instance
func NewManager(config *utils.Config) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	manager := &Manager{
		config:  config,
		fetcher: NewFetcher(config.Timeout, config.UserAgent),
		queue:   make(chan crawlTask, config.MaxPages*2), // Buffer for queue
		results: make([]*models.PageResult, 0, config.MaxPages),
		ctx:     ctx,
		cancel:  cancel,
	}

	// Initialize robots checker
	manager.robotsChecker = NewRobotsChecker(manager.fetcher, config.UserAgent, config.RespectRobots)

	// Initialize sitemap parser
	manager.sitemapParser = NewSitemapParser(manager.fetcher)

	// Initialize link graph
	manager.linkGraph = graph.NewGraph()

	// Setup graceful shutdown
	go manager.handleSignals()

	return manager
}

// SetProgressCallback sets a callback function that will be called when each page is crawled
func (m *Manager) SetProgressCallback(callback ProgressCallback) {
	m.progressCallback = callback
}

// Crawl starts the crawling process
func (m *Manager) Crawl() ([]*models.PageResult, error) {
	// Normalize start URL
	startURL, err := utils.NormalizeURL(m.config.StartURL)
	if err != nil {
		return nil, fmt.Errorf("invalid start URL: %w", err)
	}
	
	// Store normalized start URL for domain comparison
	m.normalizedStartURL = startURL

		// Parse sitemap if enabled
		var seedURLs []string
		if m.config.ParseSitemap {
			sitemapURL := m.sitemapParser.DiscoverSitemapURL(startURL)
			utils.Info("Parsing sitemap", utils.NewField("url", sitemapURL))
			
			urls, err := m.sitemapParser.ParseSitemap(sitemapURL)
			if err != nil {
				utils.Debug("Failed to parse sitemap", utils.NewField("url", sitemapURL), utils.NewField("error", err.Error()))
			} else {
				// Filter out image URLs from sitemap
				for _, url := range urls {
					if !utils.IsImageURL(url) {
						seedURLs = append(seedURLs, url)
					} else {
						utils.Debug("Skipping image URL from sitemap", utils.NewField("url", url))
					}
				}
				utils.Info("Found URLs in sitemap", utils.NewField("count", len(seedURLs)), utils.NewField("filtered_images", len(urls)-len(seedURLs)))
			}
		}

	// If no sitemap URLs found, use start URL
	if len(seedURLs) == 0 {
		seedURLs = []string{startURL}
	}

	// Start worker pool
	for i := 0; i < m.config.Workers; i++ {
		m.wg.Add(1)
		go m.worker(i)
	}

	// Enqueue initial tasks (don't mark as visited yet - workers will do that)
	enqueueDone := make(chan bool)
	go func() {
		defer close(enqueueDone)
		for _, url := range seedURLs {
			// Normalize URL
			normalized, err := utils.NormalizeURL(url)
			if err != nil {
				utils.Debug("Failed to normalize seed URL", utils.NewField("url", url), utils.NewField("error", err.Error()))
				continue
			}
			
			atomic.AddInt32(&m.pending, 1)
			m.queue <- crawlTask{
				URL:   normalized,
				Depth: 0,
			}
		}
	}()

	// Wait for initial enqueueing to complete
	<-enqueueDone
	utils.Debug("Initial tasks enqueued", utils.NewField("count", len(seedURLs)))

	// Monitor queue and close when done
	go m.monitorQueue()

	// Wait for all workers to finish
	m.wg.Wait()

	// Return results - don't treat cancellation as error if we got results
	// (cancellation might be due to reaching max-pages, which is success)
	if m.ctx.Err() != nil && len(m.results) == 0 {
		return m.results, fmt.Errorf("crawl cancelled: %w", m.ctx.Err())
	}

	return m.results, nil
}

// GetLinkGraph returns the link graph
func (m *Manager) GetLinkGraph() *graph.Graph {
	return m.linkGraph
}

// worker processes crawl tasks from the queue
func (m *Manager) worker(id int) {
	defer m.wg.Done()

	for {
		select {
		case <-m.ctx.Done():
			utils.Debug("Worker stopping", utils.NewField("worker_id", id))
			return
		case task, ok := <-m.queue:
			if !ok {
				utils.Debug("Worker queue closed", utils.NewField("worker_id", id))
				return
			}

			// Decrement pending counter (but don't let it go negative)
			currentPending := atomic.AddInt32(&m.pending, -1)
			if currentPending < 0 {
				// Reset if it went negative (shouldn't happen, but safety check)
				atomic.StoreInt32(&m.pending, 0)
			}

			// Check if we've reached max pages BEFORE processing
			m.resultsMu.Lock()
			if len(m.results) >= m.config.MaxPages {
				m.resultsMu.Unlock()
				// Cancel to signal other workers to stop
				m.cancel()
				return
			}
			m.resultsMu.Unlock()

			// Check depth limit - pages at max depth should still be crawled,
			// but we won't discover links from them (handled later)
			// Only skip if depth exceeds max depth
			if task.Depth > m.config.MaxDepth {
				utils.Debug("Skipping task - depth exceeds max", utils.NewField("url", task.URL), utils.NewField("depth", task.Depth), utils.NewField("max_depth", m.config.MaxDepth))
				continue
			}

			// Check if already visited (before marking to avoid race condition)
			if _, visited := m.visited.LoadOrStore(task.URL, true); visited {
				continue
			}

			// Skip image URLs - they should never be crawled as pages
			if utils.IsImageURL(task.URL) {
				utils.Debug("Skipping image URL - not a crawlable page", utils.NewField("url", task.URL))
				continue
			}

			// Check robots.txt before fetching
			if allowed, err := m.robotsChecker.IsAllowed(task.URL); err != nil {
				utils.Debug("Robots check error", utils.NewField("url", task.URL), utils.NewField("error", err.Error()))
			} else if !allowed {
				utils.Debug("URL disallowed by robots.txt", utils.NewField("url", task.URL))
				continue
			}

			// Apply delay if configured
			if m.config.Delay > 0 {
				select {
				case <-m.ctx.Done():
					return
				case <-time.After(m.config.Delay):
				}
			}

			// Fetch the URL with retry logic
			result := m.fetcher.FetchWithRetry(task.URL, 3)

			// Skip non-HTML content (images, PDFs, etc.) - don't add to results
			if result.Error != nil && strings.Contains(result.Error.Error(), "skipped non-HTML") {
				utils.Debug("Skipping non-HTML content", utils.NewField("url", task.URL), utils.NewField("error", result.Error.Error()))
				continue
			}

			// Store result (check limit again before storing)
			m.resultsMu.Lock()
			resultCount := len(m.results)
			if resultCount >= m.config.MaxPages {
				m.resultsMu.Unlock()
				m.cancel()
				return
			}
			m.results = append(m.results, result.PageResult)
			resultCount = len(m.results)
			m.resultsMu.Unlock()

			utils.Info("Crawled page",
				utils.NewField("url", task.URL),
				utils.NewField("status", result.PageResult.StatusCode),
				utils.NewField("depth", task.Depth),
				utils.NewField("total", resultCount),
			)

			// If fetch failed or not HTML, call progress callback and continue
			if result.Error != nil || result.PageResult.StatusCode != 200 {
				utils.Info("Skipping link discovery - fetch failed or non-200", 
					utils.NewField("url", task.URL),
					utils.NewField("error", result.Error),
					utils.NewField("status", result.PageResult.StatusCode))
				// Call progress callback even for failed pages
				if m.progressCallback != nil {
					m.progressCallback(result.PageResult, resultCount)
				}
				continue
			}

			// Check if we have body content
			if len(result.Body) == 0 {
				utils.Warn("No body content to parse", utils.NewField("url", task.URL))
				// Call progress callback even if no body
				if m.progressCallback != nil {
					m.progressCallback(result.PageResult, resultCount)
				}
				continue
			}

			// Parse HTML and discover links
			parser, err := NewParser(task.URL)
			if err != nil {
				utils.Error("Failed to create parser", utils.NewField("url", task.URL), utils.NewField("error", err.Error()))
				// Call progress callback even if parser creation failed
				if m.progressCallback != nil {
					m.progressCallback(result.PageResult, resultCount)
				}
				continue
			}

			// Merge parsed SEO data into result
			parsedData, err := parser.Parse(result.Body)
			if err != nil {
				utils.Error("Failed to parse HTML", utils.NewField("url", task.URL), utils.NewField("error", err.Error()))
				// Call progress callback even if parsing failed
				if m.progressCallback != nil {
					m.progressCallback(result.PageResult, resultCount)
				}
				continue
			}
			
			utils.Info("Parsed page", 
				utils.NewField("url", task.URL), 
				utils.NewField("depth", task.Depth),
				utils.NewField("h1_count", len(parsedData.H1)),
				utils.NewField("h1_values", parsedData.H1),
				utils.NewField("internal_links", len(parsedData.InternalLinks)),
				utils.NewField("external_links", len(parsedData.ExternalLinks)),
				utils.NewField("images", len(parsedData.Images)),
				utils.NewField("body_size", len(result.Body)))

			// Merge parsed data into page result
			result.PageResult.Title = parsedData.Title
			result.PageResult.MetaDesc = parsedData.MetaDesc
			result.PageResult.Canonical = parsedData.Canonical
			result.PageResult.H1 = parsedData.H1
			result.PageResult.H2 = parsedData.H2
			result.PageResult.H3 = parsedData.H3
			result.PageResult.H4 = parsedData.H4
			result.PageResult.H5 = parsedData.H5
			result.PageResult.H6 = parsedData.H6
			result.PageResult.InternalLinks = parsedData.InternalLinks
			result.PageResult.ExternalLinks = parsedData.ExternalLinks
			result.PageResult.Images = parsedData.Images

			// Call progress callback AFTER parsing and merging data
			// This ensures the stored page has all the parsed SEO data (H1, links, etc.)
			if m.progressCallback != nil {
				m.progressCallback(result.PageResult, resultCount)
			}

			// Check if we've reached max pages after storing
			if resultCount >= m.config.MaxPages {
				m.cancel()
				return
			}

			// Add edges to link graph
			m.linkGraph.AddEdges(task.URL, parsedData.InternalLinks)
			m.linkGraph.AddEdges(task.URL, parsedData.ExternalLinks)

			// Enqueue discovered internal links for crawling
			// Only discover links if we haven't reached max depth yet
			if task.Depth < m.config.MaxDepth {
				enqueuedCount := 0
				skippedCount := 0
				domainSkippedCount := 0
				visitedSkippedCount := 0
				
				utils.Info("Discovering links", 
					utils.NewField("url", task.URL),
					utils.NewField("depth", task.Depth),
					utils.NewField("max_depth", m.config.MaxDepth),
					utils.NewField("total_internal_links", len(parsedData.InternalLinks)))
				
				for _, linkURL := range parsedData.InternalLinks {
					// Skip image URLs - they should not be crawled as pages
					if utils.IsImageURL(linkURL) {
						skippedCount++
						utils.Debug("Skipping image URL", utils.NewField("link", linkURL))
						continue
					}
					
					// Check domain filter (use normalized start URL for comparison)
					if m.config.DomainFilter == "same" && !utils.IsSameDomain(linkURL, m.normalizedStartURL) {
						domainSkippedCount++
						utils.Info("Skipping link - different domain", 
							utils.NewField("link", linkURL), 
							utils.NewField("start_url", m.normalizedStartURL))
						continue
					}

					// Check if already visited
					if _, visited := m.visited.Load(linkURL); visited {
						visitedSkippedCount++
						utils.Info("Skipping link - already visited", utils.NewField("link", linkURL))
						continue
					}

					// Enqueue new task (check if queue is still open)
					// Check if queue is closed before attempting to send
					if atomic.LoadInt32(&m.queueClosed) == 1 {
						utils.Warn("Queue closed, stopping link discovery", utils.NewField("url", task.URL))
						return
					}
					
					select {
					case <-m.ctx.Done():
						utils.Info("Context cancelled, stopping link discovery")
						return
					case m.queue <- crawlTask{URL: linkURL, Depth: task.Depth + 1}:
						// Successfully enqueued
						atomic.AddInt32(&m.pending, 1)
						enqueuedCount++
						utils.Info("Enqueued link", utils.NewField("link", linkURL), utils.NewField("new_depth", task.Depth+1))
					default:
						// Queue full, skip (but don't panic)
						utils.Warn("Queue full, skipping link", utils.NewField("url", linkURL))
						skippedCount++
					}
				}
				utils.Info("Link discovery complete", 
					utils.NewField("url", task.URL),
					utils.NewField("enqueued", enqueuedCount),
					utils.NewField("skipped_domain", domainSkippedCount),
					utils.NewField("skipped_visited", visitedSkippedCount),
					utils.NewField("skipped_queue_full", skippedCount),
					utils.NewField("total_internal", len(parsedData.InternalLinks)))
			} else {
				utils.Info("Max depth reached, not discovering links", 
					utils.NewField("url", task.URL),
					utils.NewField("depth", task.Depth),
					utils.NewField("max_depth", m.config.MaxDepth))
			}

			// Check if we've reached max pages
			if resultCount >= m.config.MaxPages {
				m.cancel()
				return
			}
		}
	}
}

// monitorQueue closes the queue when all tasks are processed
func (m *Manager) monitorQueue() {
	ticker := time.NewTicker(2 * time.Second) // Check every 2 seconds (less frequent)
	defer ticker.Stop()
	
	emptyCount := 0 // Count consecutive empty checks
	const maxEmptyChecks = 5 // Close after 5 consecutive empty checks (10 seconds total)
	// Increased timeout to allow workers time to discover and enqueue links

	for {
		select {
		case <-m.ctx.Done():
			utils.Info("Monitor queue: context cancelled, closing queue")
			atomic.StoreInt32(&m.queueClosed, 1)
			close(m.queue)
			return
		case <-ticker.C:
			// Check if queue is empty and no pending tasks
			pending := atomic.LoadInt32(&m.pending)
			queueLen := len(m.queue)
			resultCount := len(m.results)
			
			utils.Info("Monitor queue check", 
				utils.NewField("pending", pending),
				utils.NewField("queue_len", queueLen),
				utils.NewField("results", resultCount),
				utils.NewField("empty_count", emptyCount),
				utils.NewField("max_pages", m.config.MaxPages))
			
			// Only close if:
			// 1. Queue is empty AND no pending tasks
			// 2. We haven't reached max pages (if we have, workers will cancel)
			// 3. We've had multiple consecutive empty checks
			if pending <= 0 && queueLen == 0 && resultCount < m.config.MaxPages {
				emptyCount++
				if emptyCount >= maxEmptyChecks {
					utils.Info("Closing queue - no pending tasks after multiple checks", 
						utils.NewField("empty_checks", emptyCount),
						utils.NewField("total_results", resultCount),
						utils.NewField("max_pages", m.config.MaxPages))
					atomic.StoreInt32(&m.queueClosed, 1)
					close(m.queue)
					return
				}
			} else {
				// Reset counter if we have pending work or reached max pages
				emptyCount = 0
			}
		}
	}
}

// handleSignals sets up graceful shutdown on interrupt signals
func (m *Manager) handleSignals() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	utils.Info("Received interrupt signal, shutting down gracefully...")
	m.cancel()
}

