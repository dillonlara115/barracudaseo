package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/dillonlara115/barracudaseo/internal/api"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	apiPort               int
	apiSupabaseURL        string
	apiSupabaseServiceKey string
	apiSupabaseAnonKey    string
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the Cloud Run API server",
	Long: `Start the REST API server for Cloud Run deployment.
This server handles crawl ingestion, project management, and provides
authenticated endpoints for the dashboard.`,
	RunE: runAPI,
}

func init() {
	apiCmd.Flags().IntVar(&apiPort, "port", 8091, "Port to run the API server on")
	apiCmd.Flags().StringVar(&apiSupabaseURL, "supabase-url", "", "Supabase project URL (or set PUBLIC_SUPABASE_URL env var)")
	apiCmd.Flags().StringVar(&apiSupabaseServiceKey, "supabase-service-key", "", "Supabase service role key (or set SUPABASE_SERVICE_ROLE_KEY env var)")
	apiCmd.Flags().StringVar(&apiSupabaseAnonKey, "supabase-anon-key", "", "Supabase anon key (or set PUBLIC_SUPABASE_ANON_KEY env var)")

	rootCmd.AddCommand(apiCmd)
}

func runAPI(cmd *cobra.Command, args []string) error {
	loadAPIEnvFiles()

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer logger.Sync()

	// Check if PORT is set (Cloud Run sets this)
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		if p, err := strconv.Atoi(portEnv); err == nil {
			apiPort = p
		}
	}

	// Make the resolved API port available to downstream initialization code that
	// derives localhost callback URLs from PORT.
	_ = os.Setenv("PORT", strconv.Itoa(apiPort))

	server, err := newAPIServer(logger)
	if err != nil {
		return err
	}

	// Create HTTP server
	// Increased WriteTimeout to 60s to accommodate AI operations which can take longer
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", apiPort),
		Handler:      server.Router(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second, // Increased for AI operations
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info("Starting API server",
			zap.Int("port", apiPort),
		)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Crawls running in this process cannot outlive it; record that now, while
	// there is still time, instead of leaving them marked as running.
	server.InterruptActiveCrawls()

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
		return err
	}

	logger.Info("Server exited")
	return nil
}

// loadAPIEnvFiles loads .env, and .env.local in local development.
func loadAPIEnvFiles() {
	// Ignore errors - these files are optional
	_ = godotenv.Load() // Load .env first

	// Skip .env.local in production (Cloud Run sets PORT for services and
	// CLOUD_RUN_JOB for jobs). This prevents local dev values from overriding
	// production config.
	if os.Getenv("PORT") == "" && os.Getenv("CLOUD_RUN_JOB") == "" {
		_ = godotenv.Overload(".env.local") // Only override with .env.local in local dev
	}
}

// newAPIServer builds the API server from flags and environment variables.
// The api command serves HTTP with it; the crawl-job command runs one crawl.
func newAPIServer(logger *zap.Logger) (*api.Server, error) {
	supabaseURL := apiSupabaseURL
	if supabaseURL == "" {
		supabaseURL = os.Getenv("PUBLIC_SUPABASE_URL")
	}
	if supabaseURL == "" {
		return nil, fmt.Errorf("PUBLIC_SUPABASE_URL is required (flag or environment variable)")
	}

	supabaseServiceKey := apiSupabaseServiceKey
	if supabaseServiceKey == "" {
		supabaseServiceKey = os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
	}
	if supabaseServiceKey == "" {
		return nil, fmt.Errorf("SUPABASE_SERVICE_ROLE_KEY is required (flag or environment variable)")
	}

	supabaseAnonKey := apiSupabaseAnonKey
	if supabaseAnonKey == "" {
		supabaseAnonKey = os.Getenv("PUBLIC_SUPABASE_ANON_KEY")
	}
	if supabaseAnonKey == "" {
		return nil, fmt.Errorf("PUBLIC_SUPABASE_ANON_KEY is required (flag or environment variable)")
	}

	// Log configuration (without sensitive data)
	logger.Info("Initializing API server",
		zap.String("supabase_url", supabaseURL),
		zap.Bool("has_service_key", supabaseServiceKey != ""),
		zap.Bool("has_anon_key", supabaseAnonKey != ""))

	server, err := api.NewServer(api.Config{
		SupabaseURL:        supabaseURL,
		SupabaseServiceKey: supabaseServiceKey,
		SupabaseAnonKey:    supabaseAnonKey,
		SupabaseJWTSecret:  os.Getenv("SUPABASE_JWT_SECRET"),
		CronSyncSecret:     os.Getenv("GSC_SYNC_SECRET"),
		CrawlJob:           os.Getenv("CRAWL_JOB"),
		Logger:             logger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize API server: %w", err)
	}
	return server, nil
}
