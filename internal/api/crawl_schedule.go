package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/supabase-community/postgrest-go"
	"go.uber.org/zap"
)

// maxScheduledCrawlsPerRun caps how many crawls one cron invocation starts.
// Crawls run inside the API process, so starting many at once would starve it;
// schedules left over are picked up by the next invocation.
const maxScheduledCrawlsPerRun = 3

// UpdateCrawlScheduleRequest creates or replaces a project's crawl schedule.
type UpdateCrawlScheduleRequest struct {
	Frequency        string `json:"frequency"` // "daily", "weekly" or "monthly"
	Enabled          *bool  `json:"enabled,omitempty"`
	URL              string `json:"url,omitempty"` // defaults to the project's domain
	MaxPages         int    `json:"max_pages,omitempty"`
	MaxDepth         int    `json:"max_depth,omitempty"`
	RespectRobots    *bool  `json:"respect_robots,omitempty"`
	ParseSitemap     *bool  `json:"parse_sitemap,omitempty"`
	CrawlSitemapOnly *bool  `json:"crawl_sitemap_only,omitempty"`
}

// allowedScheduleFrequencies returns the schedule frequencies a subscription
// tier ("free", "pro", "team") may use. An empty result means the tier cannot
// schedule crawls at all.
//
// TODO(dillon): this is a pricing decision. The baseline gives free accounts no
// scheduling and paid accounts every frequency. Each scheduled crawl spends
// crawl compute with nobody watching, and a daily crawl at the Pro page limit
// is 30x the load of a monthly one, so you may want e.g. weekly+monthly on Pro
// and daily only on Team.
func allowedScheduleFrequencies(tier string) []string {
	switch tier {
	case "pro", "team":
		return []string{"daily", "weekly", "monthly"}
	default:
		return []string{}
	}
}

// frequencyAllowed reports whether a tier may use a schedule frequency.
func frequencyAllowed(tier, frequency string) bool {
	for _, allowed := range allowedScheduleFrequencies(tier) {
		if allowed == frequency {
			return true
		}
	}
	return false
}

// nextRunAfter returns when a schedule should next run, counting from a time.
func nextRunAfter(frequency string, from time.Time) (time.Time, error) {
	switch frequency {
	case "daily":
		return from.AddDate(0, 0, 1), nil
	case "weekly":
		return from.AddDate(0, 0, 7), nil
	case "monthly":
		return from.AddDate(0, 1, 0), nil
	default:
		return time.Time{}, fmt.Errorf("unknown frequency %q", frequency)
	}
}

// scheduleCrawlRequest turns a schedule's stored config into a crawl request,
// applying defaults and capping the page count at the plan limit.
func scheduleCrawlRequest(config map[string]interface{}, maxPagesLimit int) TriggerCrawlRequest {
	boolOr := func(key string, fallback bool) *bool {
		value := fallback
		if v, ok := config[key].(bool); ok {
			value = v
		}
		return &value
	}
	intOr := func(key string, fallback int) int {
		if v, ok := config[key].(float64); ok && v > 0 {
			return int(v)
		}
		return fallback
	}

	req := TriggerCrawlRequest{
		MaxDepth:         intOr("max_depth", 3),
		MaxPages:         intOr("max_pages", maxPagesLimit),
		Workers:          10,
		RespectRobots:    boolOr("respect_robots", true),
		ParseSitemap:     boolOr("parse_sitemap", false),
		CrawlSitemapOnly: boolOr("crawl_sitemap_only", false),
	}
	req.URL, _ = config["url"].(string)
	if req.MaxPages > maxPagesLimit {
		req.MaxPages = maxPagesLimit
	}
	return req
}

// fetchCrawlSchedule returns a project's schedule, or nil if it has none.
func (s *Server) fetchCrawlSchedule(projectID string) (map[string]interface{}, error) {
	data, _, err := s.serviceRole.From("crawl_schedules").
		Select("*", "", false).
		Eq("project_id", projectID).
		Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch crawl schedule: %w", err)
	}
	var schedules []map[string]interface{}
	if err := json.Unmarshal(data, &schedules); err != nil {
		return nil, fmt.Errorf("failed to parse crawl schedule: %w", err)
	}
	if len(schedules) == 0 {
		return nil, nil
	}
	return schedules[0], nil
}

// handleProjectCrawlSchedule handles /api/v1/projects/:id/crawl-schedule
// (GET to read, PUT to create or replace, DELETE to remove).
func (s *Server) handleProjectCrawlSchedule(w http.ResponseWriter, r *http.Request, projectID, userID string) {
	hasAccess, err := s.verifyProjectAccess(userID, projectID)
	if err != nil {
		s.logger.Error("Failed to verify project access", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to verify access")
		return
	}
	if !hasAccess {
		s.respondError(w, http.StatusForbidden, "You don't have access to this project")
		return
	}

	subscription, err := s.resolveSubscription(userID)
	if err != nil {
		s.logger.Error("Failed to resolve subscription", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to verify subscription")
		return
	}
	tier := subscription.EffectiveTier

	switch r.Method {
	case http.MethodGet:
		schedule, err := s.fetchCrawlSchedule(projectID)
		if err != nil {
			s.logger.Error("Failed to fetch crawl schedule", zap.Error(err), zap.String("project_id", projectID))
			s.respondError(w, http.StatusInternalServerError, "Failed to fetch crawl schedule")
			return
		}
		s.respondJSON(w, http.StatusOK, map[string]interface{}{
			"schedule":            schedule,
			"allowed_frequencies": allowedScheduleFrequencies(tier),
			"max_pages_limit":     getMaxPagesLimit(tier),
		})

	case http.MethodPut:
		s.handleUpdateCrawlSchedule(w, r, projectID, userID, tier)

	case http.MethodDelete:
		_, _, err := s.serviceRole.From("crawl_schedules").
			Delete("", "").
			Eq("project_id", projectID).
			Execute()
		if err != nil {
			s.logger.Error("Failed to delete crawl schedule", zap.Error(err), zap.String("project_id", projectID))
			s.respondError(w, http.StatusInternalServerError, "Failed to delete crawl schedule")
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// handleUpdateCrawlSchedule creates or replaces a project's crawl schedule.
func (s *Server) handleUpdateCrawlSchedule(w http.ResponseWriter, r *http.Request, projectID, userID, tier string) {
	var req UpdateCrawlScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
		return
	}

	now := time.Now().UTC()
	nextRun, err := nextRunAfter(req.Frequency, now)
	if err != nil {
		s.respondError(w, http.StatusBadRequest, "frequency must be one of: daily, weekly, monthly")
		return
	}
	if !frequencyAllowed(tier, req.Frequency) {
		s.respondError(w, http.StatusForbidden, fmt.Sprintf("Your %s plan does not include %s scheduled crawls. Please upgrade to use this schedule.", tier, req.Frequency))
		return
	}

	maxPagesLimit := getMaxPagesLimit(tier)
	if req.MaxPages > maxPagesLimit {
		s.respondError(w, http.StatusForbidden, fmt.Sprintf("Your %s plan allows a maximum of %d pages per crawl.", tier, maxPagesLimit))
		return
	}
	if req.MaxPages < 0 || req.MaxDepth < 0 {
		s.respondError(w, http.StatusBadRequest, "max_pages and max_depth must not be negative")
		return
	}

	url, err := s.resolveScheduleURL(projectID, req.URL)
	if err != nil {
		s.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	config := map[string]interface{}{"url": url}
	if req.MaxPages > 0 {
		config["max_pages"] = req.MaxPages
	}
	if req.MaxDepth > 0 {
		config["max_depth"] = req.MaxDepth
	}
	if req.RespectRobots != nil {
		config["respect_robots"] = *req.RespectRobots
	}
	if req.ParseSitemap != nil {
		config["parse_sitemap"] = *req.ParseSitemap
	}
	if req.CrawlSitemapOnly != nil {
		config["crawl_sitemap_only"] = *req.CrawlSitemapOnly
	}

	record := map[string]interface{}{
		"project_id":  projectID,
		"created_by":  userID,
		"frequency":   req.Frequency,
		"enabled":     enabled,
		"config":      config,
		"next_run_at": nextRun.Format(time.RFC3339),
	}

	data, _, err := s.serviceRole.From("crawl_schedules").
		Insert(record, true, "project_id", "representation", "").
		Execute()
	if err != nil {
		s.logger.Error("Failed to save crawl schedule", zap.Error(err), zap.String("project_id", projectID))
		s.respondError(w, http.StatusInternalServerError, "Failed to save crawl schedule")
		return
	}

	var saved []map[string]interface{}
	if err := json.Unmarshal(data, &saved); err != nil || len(saved) == 0 {
		s.logger.Error("Failed to parse saved crawl schedule", zap.Error(err), zap.String("project_id", projectID))
		s.respondError(w, http.StatusInternalServerError, "Failed to save crawl schedule")
		return
	}

	s.logger.Info("Saved crawl schedule",
		zap.String("project_id", projectID),
		zap.String("user_id", userID),
		zap.String("frequency", req.Frequency),
		zap.Bool("enabled", enabled))
	s.respondJSON(w, http.StatusOK, map[string]interface{}{"schedule": saved[0]})
}

// resolveScheduleURL returns the start URL for a schedule. The URL must be on
// the project's own domain; when none is given the domain itself is used.
func (s *Server) resolveScheduleURL(projectID, requested string) (string, error) {
	data, _, err := s.serviceRole.From("projects").
		Select("domain", "", false).
		Eq("id", projectID).
		Execute()
	var projects []map[string]interface{}
	if err != nil || json.Unmarshal(data, &projects) != nil || len(projects) == 0 {
		return "", fmt.Errorf("could not load the project's domain")
	}
	domain, _ := projects[0]["domain"].(string)
	domain = strings.TrimSpace(domain)
	if domain == "" {
		return "", fmt.Errorf("the project has no domain set")
	}

	host := strings.ToLower(bareHost(domain))
	requested = strings.TrimSpace(requested)
	if requested == "" {
		if strings.Contains(domain, "://") {
			return domain, nil
		}
		return "https://" + domain, nil
	}

	if !strings.HasPrefix(requested, "http://") && !strings.HasPrefix(requested, "https://") {
		return "", fmt.Errorf("url must start with http:// or https://")
	}
	requestedHost := strings.ToLower(bareHost(requested))
	if strings.TrimPrefix(requestedHost, "www.") != strings.TrimPrefix(host, "www.") {
		return "", fmt.Errorf("url must be on the project's domain (%s)", host)
	}
	return requested, nil
}

// bareHost strips the scheme, path and port from a URL or domain string.
func bareHost(value string) string {
	if i := strings.Index(value, "://"); i >= 0 {
		value = value[i+3:]
	}
	if i := strings.IndexAny(value, "/?#"); i >= 0 {
		value = value[:i]
	}
	if i := strings.LastIndex(value, ":"); i >= 0 {
		value = value[:i]
	}
	return value
}

// handleScheduledCrawls handles POST /api/internal/crawls/run-scheduled. A cron
// job calls it to start crawls for schedules that are due.
func (s *Server) handleScheduledCrawls(w http.ResponseWriter, r *http.Request) {
	secret := r.Header.Get("X-Cron-Secret")
	if secret == "" || secret != s.cronSecret {
		s.respondError(w, http.StatusUnauthorized, "Invalid or missing cron secret")
		return
	}
	if r.Method != http.MethodPost {
		s.respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	now := time.Now().UTC()
	data, _, err := s.serviceRole.From("crawl_schedules").
		Select("*", "", false).
		Eq("enabled", "true").
		Lte("next_run_at", now.Format(time.RFC3339)).
		Order("next_run_at", &postgrest.OrderOpts{Ascending: true}).
		Limit(maxScheduledCrawlsPerRun, "").
		Execute()
	if err != nil {
		s.logger.Error("Failed to query due crawl schedules", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to query due crawl schedules")
		return
	}
	var schedules []map[string]interface{}
	if err := json.Unmarshal(data, &schedules); err != nil {
		s.logger.Error("Failed to parse due crawl schedules", zap.Error(err))
		s.respondError(w, http.StatusInternalServerError, "Failed to parse due crawl schedules")
		return
	}

	started := []map[string]interface{}{}
	skipped := []map[string]interface{}{}
	for _, schedule := range schedules {
		projectID, _ := schedule["project_id"].(string)
		crawlID, reason := s.runScheduledCrawl(schedule, now)
		if crawlID == "" {
			s.logger.Warn("Skipped scheduled crawl", zap.String("project_id", projectID), zap.String("reason", reason))
			skipped = append(skipped, map[string]interface{}{"project_id": projectID, "reason": reason})
			continue
		}
		s.logger.Info("Started scheduled crawl", zap.String("project_id", projectID), zap.String("crawl_id", crawlID))
		started = append(started, map[string]interface{}{"project_id": projectID, "crawl_id": crawlID})
	}

	s.respondJSON(w, http.StatusOK, map[string]interface{}{
		"due":     len(schedules),
		"started": started,
		"skipped": skipped,
	})
}

// runScheduledCrawl starts the crawl for one due schedule. It returns the new
// crawl's ID, or an empty ID and the reason the crawl was not started.
func (s *Server) runScheduledCrawl(schedule map[string]interface{}, now time.Time) (string, string) {
	scheduleID, _ := schedule["id"].(string)
	projectID, _ := schedule["project_id"].(string)
	createdBy, _ := schedule["created_by"].(string)
	frequency, _ := schedule["frequency"].(string)
	config, _ := schedule["config"].(map[string]interface{})

	// Move the schedule forward before doing anything else, so a schedule that
	// keeps failing is retried at its next slot instead of on every cron run.
	nextRun, err := nextRunAfter(frequency, now)
	if err != nil {
		return "", err.Error()
	}
	_, _, err = s.serviceRole.From("crawl_schedules").
		Update(map[string]interface{}{"next_run_at": nextRun.Format(time.RFC3339)}, "", "").
		Eq("id", scheduleID).
		Execute()
	if err != nil {
		s.logger.Error("Failed to advance crawl schedule", zap.Error(err), zap.String("schedule_id", scheduleID))
		return "", "could not advance the schedule"
	}

	subscription, err := s.resolveSubscription(createdBy)
	if err != nil {
		s.logger.Error("Failed to resolve subscription for scheduled crawl", zap.Error(err), zap.String("user_id", createdBy))
		return "", "could not verify the subscription"
	}
	tier := subscription.EffectiveTier
	if !frequencyAllowed(tier, frequency) {
		return "", fmt.Sprintf("the %s plan does not include %s scheduled crawls", tier, frequency)
	}

	// Don't stack a scheduled crawl on top of one that is still running
	runningData, _, err := s.serviceRole.From("crawls").
		Select("id", "", false).
		Eq("project_id", projectID).
		Eq("status", "running").
		Limit(1, "").
		Execute()
	var running []map[string]interface{}
	if err != nil || json.Unmarshal(runningData, &running) != nil {
		s.logger.Error("Failed to check for running crawls", zap.Error(err), zap.String("project_id", projectID))
		return "", "could not check for running crawls"
	}
	if len(running) > 0 {
		return "", "a crawl is already running for this project"
	}

	req := scheduleCrawlRequest(config, getMaxPagesLimit(tier))
	if req.URL == "" {
		return "", "the schedule has no URL"
	}

	crawlID := uuid.New().String()
	crawl := map[string]interface{}{
		"id":           crawlID,
		"project_id":   projectID,
		"initiated_by": createdBy,
		"source":       "schedule",
		"status":       "running",
		"started_at":   now.Format(time.RFC3339),
		"total_pages":  0,
		"total_issues": 0,
		"meta": map[string]interface{}{
			"url":                req.URL,
			"max_depth":          req.MaxDepth,
			"max_pages":          req.MaxPages,
			"workers":            req.Workers,
			"respect_robots":     *req.RespectRobots,
			"parse_sitemap":      *req.ParseSitemap,
			"crawl_sitemap_only": *req.CrawlSitemapOnly,
			"schedule_id":        scheduleID,
		},
	}
	if _, _, err := s.serviceRole.From("crawls").Insert(crawl, false, "", "", "").Execute(); err != nil {
		s.logger.Error("Failed to insert scheduled crawl", zap.Error(err), zap.String("project_id", projectID))
		return "", "could not create the crawl"
	}

	_, _, err = s.serviceRole.From("crawl_schedules").
		Update(map[string]interface{}{
			"last_run_at":   now.Format(time.RFC3339),
			"last_crawl_id": crawlID,
		}, "", "").
		Eq("id", scheduleID).
		Execute()
	if err != nil {
		s.logger.Warn("Failed to record scheduled crawl on schedule", zap.Error(err), zap.String("schedule_id", scheduleID))
	}

	go s.runCrawlAsync(crawlID, projectID, req)
	return crawlID, ""
}
