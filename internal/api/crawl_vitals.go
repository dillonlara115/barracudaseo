package api

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/dillonlara115/barracudaseo/internal/psi"
	"github.com/dillonlara115/barracudaseo/pkg/models"
)

const (
	// defaultVitalsPages is how many pages per crawl are measured when
	// PSI_MAX_PAGES is not set. Each measurement takes 10 to 30 seconds.
	defaultVitalsPages = 5

	// vitalsTimeout bounds the whole measurement step so a slow PageSpeed API
	// cannot stall a crawl.
	vitalsTimeout = 3 * time.Minute
)

// measureWebVitals measures Core Web Vitals for a crawl's key pages and stores
// the readings on the crawl. It does nothing unless PSI_API_KEY is set, and it
// never fails the crawl: PageSpeed is a slow, rate-limited outside service.
func (s *Server) measureWebVitals(crawlID string, results []*models.PageResult) {
	apiKey := os.Getenv("PSI_API_KEY")
	if apiKey == "" {
		return
	}

	limit := defaultVitalsPages
	if v, err := strconv.Atoi(os.Getenv("PSI_MAX_PAGES")); err == nil && v >= 0 {
		limit = v
	}
	selected := psi.SelectPages(results, limit)
	if len(selected) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), vitalsTimeout)
	defer cancel()
	measured := psi.NewClient(apiKey).Enrich(ctx, selected, "mobile", 3)
	s.logger.Info("Measured Core Web Vitals",
		zap.String("crawl_id", crawlID),
		zap.Int("selected", len(selected)),
		zap.Int("measured", measured))

	// Pages are already stored by the time the crawl reaches this step, so the
	// readings are kept on the crawl record instead of on each page row.
	readings := make([]map[string]interface{}, 0, len(selected))
	for _, p := range selected {
		if p.WebVitals == nil {
			continue
		}
		if p.WebVitals.Error != "" {
			s.logger.Warn("Core Web Vitals measurement failed",
				zap.String("crawl_id", crawlID),
				zap.String("url", p.URL),
				zap.String("error", p.WebVitals.Error))
		}
		readings = append(readings, map[string]interface{}{"url": p.URL, "vitals": p.WebVitals})
	}
	if err := s.mergeCrawlMeta(crawlID, map[string]interface{}{"web_vitals": readings}); err != nil {
		s.logger.Warn("Failed to store Core Web Vitals on crawl", zap.String("crawl_id", crawlID), zap.Error(err))
	}
}

// mergeCrawlMeta adds keys to a crawl's meta without discarding what is there.
func (s *Server) mergeCrawlMeta(crawlID string, values map[string]interface{}) error {
	data, _, err := s.serviceRole.From("crawls").Select("meta", "", false).Eq("id", crawlID).Execute()
	if err != nil {
		return err
	}
	var crawls []map[string]interface{}
	if err := json.Unmarshal(data, &crawls); err != nil {
		return err
	}

	meta := map[string]interface{}{}
	if len(crawls) > 0 {
		if existing, ok := crawls[0]["meta"].(map[string]interface{}); ok {
			meta = existing
		}
	}
	for key, value := range values {
		meta[key] = value
	}

	_, _, err = s.serviceRole.From("crawls").
		Update(map[string]interface{}{"meta": meta}, "", "").
		Eq("id", crawlID).
		Execute()
	return err
}
