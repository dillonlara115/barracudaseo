package api

import (
	"encoding/json"
	"os"
	"strconv"

	"go.uber.org/zap"

	"github.com/dillonlara115/barracudaseo/internal/psi"
	"github.com/dillonlara115/barracudaseo/internal/render"
	"github.com/dillonlara115/barracudaseo/internal/utils"
	"github.com/dillonlara115/barracudaseo/pkg/models"
)

const (
	// phasePageExperience is the crawl phase shown while Core Web Vitals are
	// measured and pages are rendered in headless Chrome.
	phasePageExperience = "page_experience"

	// renderConcurrency is how many tabs render at once. Each tab is a full
	// page load, so this is kept low to stay inside the job's memory.
	renderConcurrency = 2

	// crawlUserAgent identifies cloud crawls, both when fetching and rendering.
	crawlUserAgent = "barracuda/1.0.0"
)

// renderPageLimit is how many pages per crawl are rendered in headless Chrome.
// Rendering is off unless RENDER_MAX_PAGES is set to a positive number: it
// needs Chrome in the image and far more memory than the API instance has, so
// only the crawl job turns it on.
func renderPageLimit() int {
	if v, err := strconv.Atoi(os.Getenv("RENDER_MAX_PAGES")); err == nil && v > 0 {
		return v
	}
	return 0
}

// pageExperienceEnabled reports whether a crawl will measure Core Web Vitals
// or render pages, so the caller knows whether to show that phase.
func pageExperienceEnabled() bool {
	return os.Getenv("PSI_API_KEY") != "" || renderPageLimit() > 0
}

// renderKeyPages renders a crawl's key pages in headless Chrome and records
// how each differs from the HTML the server sent. The comparison is set on the
// results (so the analyzer can report issues) and stored on each page row. It
// does nothing unless RENDER_MAX_PAGES is set, and it never fails the crawl.
func (s *Server) renderKeyPages(crawlID string, results []*models.PageResult) {
	limit := renderPageLimit()
	if limit == 0 {
		return
	}
	selected := render.SelectPages(results, psi.SelectPages(results, limit), limit)
	if len(selected) == 0 {
		return
	}

	renderer, err := render.NewRenderer(crawlUserAgent)
	if err != nil {
		s.logger.Warn("Skipping JavaScript rendering", zap.String("crawl_id", crawlID), zap.Error(err))
		return
	}
	defer renderer.Close()

	rendered := renderer.Enrich(selected, renderConcurrency)
	s.logger.Info("Rendered pages in headless Chrome",
		zap.String("crawl_id", crawlID),
		zap.Int("selected", len(selected)),
		zap.Int("rendered", rendered))

	for _, p := range selected {
		if p.Rendered == nil {
			continue
		}
		if p.Rendered.Error != "" {
			s.logger.Warn("JavaScript rendering failed",
				zap.String("crawl_id", crawlID),
				zap.String("url", p.URL),
				zap.String("error", p.Rendered.Error))
			continue
		}
		s.mergePageData(crawlID, p.URL, map[string]interface{}{"rendered": p.Rendered})
	}

	summary := map[string]interface{}{"selected": len(selected), "rendered": rendered}
	if err := s.mergeCrawlMeta(crawlID, map[string]interface{}{"rendering": summary}); err != nil {
		s.logger.Warn("Failed to store rendering summary on crawl", zap.String("crawl_id", crawlID), zap.Error(err))
	}
}

// mergePageData adds keys to the data of a page that is already stored. Pages
// are written while the crawl runs, so measurements taken afterwards are added
// this way. Failures are logged and otherwise ignored.
func (s *Server) mergePageData(crawlID, pageURL string, values map[string]interface{}) {
	storedURL := storedPageURL(pageURL)

	data, _, err := s.serviceRole.From("pages").
		Select("id,data", "", false).
		Eq("crawl_id", crawlID).
		Eq("url", storedURL).
		Limit(1, "").
		Execute()
	if err != nil {
		s.logger.Warn("Failed to load page for extra data", zap.Error(err), zap.String("url", storedURL))
		return
	}
	var rows []map[string]interface{}
	if err := json.Unmarshal(data, &rows); err != nil || len(rows) == 0 {
		s.logger.Warn("Page not found for extra data", zap.String("crawl_id", crawlID), zap.String("url", storedURL))
		return
	}
	pageID, ok := parsePageID(rows[0]["id"])
	if !ok {
		return
	}

	existing, _ := rows[0]["data"].(map[string]interface{})
	_, _, err = s.serviceRole.From("pages").
		Update(map[string]interface{}{"data": mergedData(existing, values)}, "", "").
		Eq("id", strconv.FormatInt(pageID, 10)).
		Execute()
	if err != nil {
		s.logger.Warn("Failed to store extra page data", zap.Error(err), zap.String("url", storedURL))
	}
}

// mergedData returns a copy of existing with values added or replaced.
func mergedData(existing, values map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(existing)+len(values))
	for key, value := range existing {
		out[key] = value
	}
	for key, value := range values {
		out[key] = value
	}
	return out
}

// storedPageURL is the form a page's URL takes in the pages table.
func storedPageURL(pageURL string) string {
	if normalized, err := utils.NormalizeURL(pageURL); err == nil {
		return normalized
	}
	return pageURL
}
