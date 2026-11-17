package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dillonlara115/barracudaseo/internal/api"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var crawlJobID string

var crawlJobCmd = &cobra.Command{
	Use:   "crawl-job",
	Short: "Run one cloud crawl to completion (used by the Cloud Run crawl job)",
	Long: `Run a crawl that the API has already created, then exit.

The API starts one Cloud Run Job execution per crawl and passes the crawl's ID
in the CRAWL_ID environment variable. The job runs this command, which crawls
the site, stores pages and issues, and records the result on the crawl.

It needs the same Supabase configuration as the api command.`,
	RunE: runCrawlJob,
}

func init() {
	crawlJobCmd.Flags().StringVar(&crawlJobID, "crawl-id", "", "ID of the crawl to run (or set CRAWL_ID env var)")

	rootCmd.AddCommand(crawlJobCmd)
}

func runCrawlJob(cmd *cobra.Command, args []string) error {
	loadAPIEnvFiles()

	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer logger.Sync()

	crawlID := crawlJobID
	if crawlID == "" {
		crawlID = os.Getenv(api.CrawlIDEnvVar)
	}
	if crawlID == "" {
		return fmt.Errorf("a crawl ID is required (--crawl-id or %s)", api.CrawlIDEnvVar)
	}

	// The job runs crawls itself; it must never hand them to another job
	_ = os.Unsetenv("CRAWL_JOB")

	server, err := newAPIServer(logger)
	if err != nil {
		return err
	}

	// Cloud Run sends SIGTERM when the execution is cancelled or times out
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Info("Running crawl job", zap.String("crawl_id", crawlID))
	if err := server.RunCrawlJob(ctx, crawlID); err != nil {
		logger.Error("Crawl job failed", zap.Error(err), zap.String("crawl_id", crawlID))
		return err
	}
	logger.Info("Crawl job finished", zap.String("crawl_id", crawlID))
	return nil
}
