package api

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// How a cloud crawl stays durable:
//
//   - In production each crawl runs in its own Cloud Run Job execution
//     (startCrawl -> crawlJob.Run -> `barracuda crawl-job` -> RunCrawlJob), so
//     it is not tied to the API instance that received the request.
//   - Without a configured job (local development) the crawl runs in this
//     process, as it always did.
//   - Either way the runner writes a heartbeat, and sweepStaleCrawls marks a
//     crawl as failed when the heartbeat stops, so no crawl stays "running"
//     forever.

const (
	// How often a live runner refreshes crawls.heartbeat_at
	crawlHeartbeatInterval = 30 * time.Second
	// A running crawl with no heartbeat for this long is treated as dead. It
	// must comfortably exceed the heartbeat interval plus a job's cold start.
	crawlStaleAfter = 10 * time.Minute
	// How long the API waits for Cloud Run to accept a job execution
	crawlJobLaunchTimeout = 20 * time.Second
	// After a shutdown signal, how long RunCrawlJob lets the runner record the
	// interruption itself before recording it on the runner's behalf
	crawlInterruptGrace = 3 * time.Second

	crawlInterruptedMessage  = "The crawl was interrupted before it finished. Start it again to get a full result."
	crawlStaleMessage        = "The crawl stopped responding and was marked as failed. Start it again to get a full result."
	crawlLaunchFailedMessage = "The crawl could not be started. Please try again."
)

// crawlJobRunner starts a crawl in its own job execution, outside this process.
type crawlJobRunner interface {
	Run(ctx context.Context, crawlID string) error
}

// runInProcessWhenJobFails decides what happens when the Cloud Run Job cannot
// be launched (bad configuration, missing permission, Cloud Run API outage).
//
// TODO(dillon): choose the trade-off.
//   - true: run the crawl inside the API process instead. Small, watched crawls
//     will usually still finish, but a long or scheduled crawl can be killed
//     when the instance goes idle, and a broken job setup is only visible in
//     the logs.
//   - false: fail the crawl straight away with crawlLaunchFailedMessage. The
//     problem is obvious to the user and to you, but every crawl fails until
//     the job is fixed.
func runInProcessWhenJobFails(err error) bool {
	return true
}

// startCrawl hands a crawl whose row already exists to a runner.
func (s *Server) startCrawl(crawlID, projectID string, req TriggerCrawlRequest) {
	if s.crawlJob != nil {
		ctx, cancel := context.WithTimeout(context.Background(), crawlJobLaunchTimeout)
		err := s.crawlJob.Run(ctx, crawlID)
		cancel()
		if err == nil {
			s.logger.Info("Started crawl job", zap.String("crawl_id", crawlID))
			return
		}
		s.logger.Error("Failed to start crawl job", zap.Error(err), zap.String("crawl_id", crawlID))
		if !runInProcessWhenJobFails(err) {
			s.updateCrawlStatus(crawlID, "failed", crawlLaunchFailedMessage)
			return
		}
		s.logger.Warn("Running crawl in the API process instead", zap.String("crawl_id", crawlID))
	}
	s.runCrawlInProcess(crawlID, projectID, req)
}

// runCrawlInProcess runs a crawl in a goroutine of this process and tracks it
// so InterruptActiveCrawls can account for it at shutdown.
func (s *Server) runCrawlInProcess(crawlID, projectID string, req TriggerCrawlRequest) {
	s.activeCrawls.Store(crawlID, struct{}{})
	go func() {
		defer s.activeCrawls.Delete(crawlID)
		// A panic in a crawl must not take the whole API down with it
		defer func() {
			if r := recover(); r != nil {
				s.logger.Error("Crawl panicked", zap.Any("panic", r), zap.String("crawl_id", crawlID))
				s.failCrawlIfRunning(crawlID, crawlInterruptedMessage)
			}
		}()
		s.runCrawlAsync(crawlID, projectID, req)
	}()
}

// InterruptActiveCrawls marks every crawl still running in this process as
// failed. Call it when the process is about to exit: those crawls cannot
// finish, and without this they would sit in "running" until the sweeper
// noticed.
func (s *Server) InterruptActiveCrawls() {
	s.activeCrawls.Range(func(key, _ interface{}) bool {
		if crawlID, ok := key.(string); ok {
			s.failCrawlIfRunning(crawlID, crawlInterruptedMessage)
		}
		return true
	})
}

// RunCrawlJob runs one crawl to completion in this process. It is the body of
// the `barracuda crawl-job` command, which a Cloud Run Job executes.
//
// It returns an error only when the crawl could not be loaded, so that the job
// execution fails and Cloud Run retries it. A crawl that ran and failed is
// recorded on the crawl row and is not an error here: retrying would not help.
func (s *Server) RunCrawlJob(ctx context.Context, crawlID string) error {
	data, _, err := s.serviceRole.From("crawls").
		Select("id,project_id,status,meta", "", false).
		Eq("id", crawlID).
		Execute()
	if err != nil {
		return fmt.Errorf("load crawl %s: %w", crawlID, err)
	}
	var rows []map[string]interface{}
	if err := json.Unmarshal(data, &rows); err != nil {
		return fmt.Errorf("parse crawl %s: %w", crawlID, err)
	}
	if len(rows) == 0 {
		return fmt.Errorf("crawl %s not found", crawlID)
	}

	crawl := rows[0]
	status, _ := crawl["status"].(string)
	if status != "running" {
		// Already finished, failed or swept: nothing to do, and not worth a retry
		s.logger.Warn("Crawl is not running, skipping", zap.String("crawl_id", crawlID), zap.String("status", status))
		return nil
	}

	projectID, _ := crawl["project_id"].(string)
	meta, _ := crawl["meta"].(map[string]interface{})
	req := crawlRequestFromMeta(meta)
	if req.URL == "" {
		s.updateCrawlStatus(crawlID, "failed", "The crawl has no URL to start from.")
		return nil
	}

	// A retried execution starts the crawl again from the beginning, so drop
	// whatever the previous attempt stored.
	s.clearCrawlResults(crawlID)

	done := make(chan struct{})
	go func() {
		defer close(done)
		s.runCrawlAsync(crawlID, projectID, req)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		// Shutdown signal (job cancelled or timed out). The runner notices the
		// same signal while it is still fetching pages; give it a moment, then
		// make sure the crawl is not left "running".
		select {
		case <-done:
		case <-time.After(crawlInterruptGrace):
		}
		s.failCrawlIfRunning(crawlID, crawlInterruptedMessage)
		return nil
	}
}

// crawlRequestFromMeta rebuilds the crawl options stored in crawls.meta when
// the crawl was created.
func crawlRequestFromMeta(meta map[string]interface{}) TriggerCrawlRequest {
	// The page limit was enforced when the crawl was created; the ceiling here
	// only guards against a corrupt value.
	req := scheduleCrawlRequest(meta, getMaxPagesLimit("team"))
	if workers, ok := meta["workers"].(float64); ok && workers > 0 {
		req.Workers = int(workers)
	}
	return req
}

// clearCrawlResults removes pages and issues stored for a crawl.
func (s *Server) clearCrawlResults(crawlID string) {
	if _, _, err := s.serviceRole.From("issues").Delete("", "").Eq("crawl_id", crawlID).Execute(); err != nil {
		s.logger.Warn("Failed to clear crawl issues", zap.Error(err), zap.String("crawl_id", crawlID))
	}
	if _, _, err := s.serviceRole.From("pages").Delete("", "").Eq("crawl_id", crawlID).Execute(); err != nil {
		s.logger.Warn("Failed to clear crawl pages", zap.Error(err), zap.String("crawl_id", crawlID))
	}
	reset := map[string]interface{}{"total_pages": 0, "total_issues": 0, "pages_discovered": 0}
	if _, _, err := s.serviceRole.From("crawls").Update(reset, "", "").Eq("id", crawlID).Execute(); err != nil {
		s.logger.Warn("Failed to reset crawl counters", zap.Error(err), zap.String("crawl_id", crawlID))
	}
}

// startCrawlHeartbeat writes crawls.heartbeat_at now and then every
// crawlHeartbeatInterval until the returned stop function is called.
func (s *Server) startCrawlHeartbeat(crawlID string) (stop func()) {
	beat := func() {
		update := map[string]interface{}{"heartbeat_at": time.Now().UTC().Format(time.RFC3339)}
		if _, _, err := s.serviceRole.From("crawls").Update(update, "", "").Eq("id", crawlID).Execute(); err != nil {
			s.logger.Warn("Failed to write crawl heartbeat", zap.Error(err), zap.String("crawl_id", crawlID))
		}
	}
	beat()

	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(crawlHeartbeatInterval)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				beat()
			}
		}
	}()

	var once sync.Once
	return func() { once.Do(func() { close(done) }) }
}

// staleCrawlFilter is the PostgREST `or` filter matching crawls whose last
// sign of life is older than cutoff: the heartbeat, or the start time for a
// crawl whose runner never wrote one.
func staleCrawlFilter(cutoff time.Time) string {
	ts := cutoff.UTC().Format(time.RFC3339)
	return fmt.Sprintf(`heartbeat_at.lt."%s",and(heartbeat_at.is.null,started_at.lt."%s")`, ts, ts)
}

// sweepStaleCrawls marks running crawls whose runner has died as failed and
// returns how many it marked.
func (s *Server) sweepStaleCrawls(now time.Time) int {
	data, _, err := s.serviceRole.From("crawls").
		Select("id", "", false).
		Eq("status", "running").
		Or(staleCrawlFilter(now.Add(-crawlStaleAfter)), "").
		Execute()
	if err != nil {
		s.logger.Error("Failed to query stale crawls", zap.Error(err))
		return 0
	}
	var stale []map[string]interface{}
	if err := json.Unmarshal(data, &stale); err != nil {
		s.logger.Error("Failed to parse stale crawls", zap.Error(err))
		return 0
	}

	swept := 0
	for _, crawl := range stale {
		crawlID, _ := crawl["id"].(string)
		if crawlID == "" {
			continue
		}
		if s.failCrawlIfRunning(crawlID, crawlStaleMessage) {
			s.logger.Warn("Marked stale crawl as failed", zap.String("crawl_id", crawlID))
			swept++
		}
	}
	return swept
}

// failCrawlIfRunning marks a crawl as failed unless it has already finished.
// It reports whether the update was sent without error.
func (s *Server) failCrawlIfRunning(crawlID, message string) bool {
	update := map[string]interface{}{
		"status":       "failed",
		"completed_at": time.Now().UTC().Format(time.RFC3339),
		"phase":        "",
		"meta":         metaWithError(s.fetchCrawlMeta(crawlID), message),
	}
	_, _, err := s.serviceRole.From("crawls").
		Update(update, "", "").
		Eq("id", crawlID).
		Eq("status", "running").
		Execute()
	if err != nil {
		s.logger.Error("Failed to mark crawl as failed", zap.Error(err), zap.String("crawl_id", crawlID))
		return false
	}
	return true
}

// fetchCrawlMeta returns a crawl's meta object, or nil if it cannot be read.
func (s *Server) fetchCrawlMeta(crawlID string) map[string]interface{} {
	data, _, err := s.serviceRole.From("crawls").Select("meta", "", false).Eq("id", crawlID).Execute()
	if err != nil {
		s.logger.Warn("Failed to read crawl meta", zap.Error(err), zap.String("crawl_id", crawlID))
		return nil
	}
	var rows []map[string]interface{}
	if err := json.Unmarshal(data, &rows); err != nil || len(rows) == 0 {
		return nil
	}
	meta, _ := rows[0]["meta"].(map[string]interface{})
	return meta
}

// metaWithError returns a copy of meta with the error message added.
func metaWithError(meta map[string]interface{}, message string) map[string]interface{} {
	out := make(map[string]interface{}, len(meta)+1)
	for key, value := range meta {
		out[key] = value
	}
	out["error"] = message
	return out
}

// discoveredPages is the progress estimate stored while a crawl scans: the
// crawler's count of URLs found, never less than the pages already crawled.
func discoveredPages(crawled, discovered int) int {
	if discovered < crawled {
		return crawled
	}
	return discovered
}

// setCrawlPagesDiscovered records the crawl's page estimate.
func (s *Server) setCrawlPagesDiscovered(crawlID string, pages int) {
	update := map[string]interface{}{"pages_discovered": pages}
	if _, _, err := s.serviceRole.From("crawls").Update(update, "", "").Eq("id", crawlID).Execute(); err != nil {
		s.logger.Warn("Failed to update crawl page estimate", zap.Error(err), zap.String("crawl_id", crawlID))
	}
}
