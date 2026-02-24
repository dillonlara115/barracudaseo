package sitecrawl

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dillonlara115/barracudaseo/internal/ai/embeddings"
	"github.com/dillonlara115/barracudaseo/internal/crawler"
	"github.com/google/uuid"
	"github.com/supabase-community/supabase-go"
	"go.uber.org/zap"
)

// CrawledPage represents a page ready for storage in the crawled_pages table
type CrawledPage struct {
	ID              string      `json:"id"`
	ProjectID       string      `json:"project_id"`
	URL             string      `json:"url"`
	Title           string      `json:"title"`
	MetaDescription string      `json:"meta_description"`
	Headings        interface{} `json:"headings"`
	BodySummary     string      `json:"body_summary"`
	InternalLinks   interface{} `json:"internal_links"`
	WordCount       int         `json:"word_count"`
	Embedding       string      `json:"embedding,omitempty"`
	CrawledAt       string      `json:"crawled_at"`
	ManuallyEdited  bool        `json:"manually_edited"`
}

// Heading represents a heading extracted from a page
type Heading struct {
	Level int    `json:"level"`
	Text  string `json:"text"`
}

// Service orchestrates sitemap-based crawling with content storage and embeddings
type Service struct {
	serviceRole  *supabase.Client
	embedService *embeddings.Service
	logger       *zap.Logger
}

// NewService creates a new site crawl service
func NewService(serviceRoleClient *supabase.Client, embedService *embeddings.Service, logger *zap.Logger) *Service {
	return &Service{
		serviceRole:  serviceRoleClient,
		embedService: embedService,
		logger:       logger,
	}
}

// CrawlFromSitemap crawls all URLs from a site's sitemap and stores clean content with embeddings
func (s *Service) CrawlFromSitemap(ctx context.Context, projectID, siteURL string, maxPages int) (int, error) {
	fetcher := crawler.NewFetcher(10*time.Second, "BarracudaSEO/1.0")
	sitemapParser := crawler.NewSitemapParser(fetcher)

	sitemapURL := sitemapParser.DiscoverSitemapURL(siteURL)
	urls, err := sitemapParser.ParseSitemap(sitemapURL)
	if err != nil {
		return 0, fmt.Errorf("failed to parse sitemap: %w", err)
	}

	if maxPages > 0 && len(urls) > maxPages {
		urls = urls[:maxPages]
	}

	s.logger.Info("Starting site crawl",
		zap.String("project_id", projectID),
		zap.Int("urls_found", len(urls)),
	)

	crawled := 0
	for _, pageURL := range urls {
		select {
		case <-ctx.Done():
			return crawled, ctx.Err()
		default:
		}

		if err := s.crawlAndStorePage(ctx, projectID, pageURL, fetcher); err != nil {
			s.logger.Warn("Failed to crawl page",
				zap.String("url", pageURL),
				zap.Error(err),
			)
			continue
		}
		crawled++
	}

	s.logger.Info("Site crawl complete",
		zap.String("project_id", projectID),
		zap.Int("crawled", crawled),
	)

	return crawled, nil
}

func (s *Service) crawlAndStorePage(ctx context.Context, projectID, pageURL string, fetcher *crawler.Fetcher) error {
	result := fetcher.Fetch(pageURL)
	if result.Error != nil {
		return fmt.Errorf("fetch failed: %w", result.Error)
	}
	if result.PageResult.StatusCode != 200 {
		return fmt.Errorf("HTTP %d", result.PageResult.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(result.Body)))
	if err != nil {
		return fmt.Errorf("parse HTML failed: %w", err)
	}

	title := strings.TrimSpace(doc.Find("title").First().Text())
	metaDesc := ""
	doc.Find("meta[name='description']").Each(func(_ int, sel *goquery.Selection) {
		if content, exists := sel.Attr("content"); exists {
			metaDesc = strings.TrimSpace(content)
		}
	})

	var headings []Heading
	for level := 1; level <= 6; level++ {
		tag := fmt.Sprintf("h%d", level)
		lvl := level
		doc.Find(tag).Each(func(_ int, sel *goquery.Selection) {
			text := strings.TrimSpace(sel.Text())
			if text != "" {
				headings = append(headings, Heading{Level: lvl, Text: text})
			}
		})
	}

	// Extract body text (strip nav, header, footer, script, style)
	doc.Find("nav, header, footer, script, style, noscript").Remove()
	bodyText := strings.TrimSpace(doc.Find("body").Text())
	bodyText = collapseWhitespace(bodyText)

	// Truncate body summary to ~2000 chars for storage
	bodySummary := bodyText
	if len(bodySummary) > 2000 {
		bodySummary = bodySummary[:2000]
	}

	wordCount := len(strings.Fields(bodyText))

	// Extract internal links
	var internalLinks []string
	doc.Find("a[href]").Each(func(_ int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		if !exists {
			return
		}
		href = strings.TrimSpace(href)
		if strings.HasPrefix(href, "/") || strings.Contains(href, pageURL) {
			internalLinks = append(internalLinks, href)
		}
	})

	// Generate embedding from title + body summary
	embeddingText := title + " " + metaDesc + " " + bodySummary
	if len(embeddingText) > 5000 {
		embeddingText = embeddingText[:5000]
	}

	var embeddingStr string
	if s.embedService != nil {
		vec, err := s.embedService.Embed(ctx, embeddingText)
		if err != nil {
			s.logger.Warn("Failed to generate embedding", zap.String("url", pageURL), zap.Error(err))
		} else {
			embeddingStr = embeddings.FormatForPgvector(vec)
		}
	}

	headingsJSON, _ := json.Marshal(headings)
	linksJSON, _ := json.Marshal(internalLinks)

	page := CrawledPage{
		ID:              uuid.New().String(),
		ProjectID:       projectID,
		URL:             pageURL,
		Title:           title,
		MetaDescription: metaDesc,
		Headings:        json.RawMessage(headingsJSON),
		BodySummary:     bodySummary,
		InternalLinks:   json.RawMessage(linksJSON),
		WordCount:       wordCount,
		Embedding:       embeddingStr,
		CrawledAt:       time.Now().UTC().Format(time.RFC3339),
		ManuallyEdited:  false,
	}

	data, err := json.Marshal(page)
	if err != nil {
		return fmt.Errorf("failed to marshal page: %w", err)
	}

	// Upsert: insert or update on conflict (project_id, url), preserving manually_edited pages
	_, _, err = s.serviceRole.From("crawled_pages").
		Upsert(json.RawMessage(data), "project_id,url", "id,manually_edited", "").
		Execute()
	if err != nil {
		return fmt.Errorf("failed to upsert crawled page: %w", err)
	}

	return nil
}

func collapseWhitespace(s string) string {
	var b strings.Builder
	prevSpace := false
	for _, r := range s {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			if !prevSpace {
				b.WriteRune(' ')
			}
			prevSpace = true
		} else {
			b.WriteRune(r)
			prevSpace = false
		}
	}
	return b.String()
}
