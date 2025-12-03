package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/dillonlara115/barracuda/internal/dataforseo"
	"github.com/dillonlara115/barracuda/internal/gsc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/supabase-community/supabase-go"
	"go.uber.org/zap"
)

// Config holds API server configuration
type Config struct {
	SupabaseURL        string
	SupabaseServiceKey string
	SupabaseAnonKey    string
	CronSyncSecret     string
	Logger             *zap.Logger
}

// Server represents the API server
type Server struct {
	config       Config
	supabase     *supabase.Client
	serviceRole  *supabase.Client
	logger       *zap.Logger
	cronSecret   string
	emailService EmailService
	tokenCache   map[string]TokenCacheEntry
	tokenCacheMu sync.RWMutex
	// Inflight requests to prevent thundering herd
	tokenInflight   map[string]chan struct{}
	tokenInflightMu sync.Mutex
	// JWKS for local JWT validation (avoids hitting Supabase for every request)
	jwks    *keyfunc.JWKS
	jwksURL string
}

// TokenCacheEntry represents a cached token validation result
type TokenCacheEntry struct {
	User      *User
	ExpiresAt time.Time
}

// NewServer creates a new API server instance
func NewServer(cfg Config) (*Server, error) {
	// Create Supabase client with anon key (for RLS-protected queries)
	supabaseClient, err := supabase.NewClient(cfg.SupabaseURL, cfg.SupabaseAnonKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Supabase client: %w", err)
	}

	// Create service role client (bypasses RLS for admin operations)
	serviceRoleClient, err := supabase.NewClient(cfg.SupabaseURL, cfg.SupabaseServiceKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Supabase service role client: %w", err)
	}

	// Initialize email service
	emailService := NewEmailService(cfg.SupabaseURL, cfg.SupabaseServiceKey, cfg.Logger)

	server := &Server{
		config:        cfg,
		supabase:      supabaseClient,
		serviceRole:   serviceRoleClient,
		logger:        cfg.Logger,
		cronSecret:    cfg.CronSyncSecret,
		emailService:  emailService,
		tokenCache:    make(map[string]TokenCacheEntry),
		tokenInflight: make(map[string]chan struct{}),
		jwksURL:       strings.TrimSuffix(cfg.SupabaseURL, "/") + "/auth/v1/keys",
	}

	// Load JWKS for local token validation; fall back to Supabase API if it fails
	jwks, err := keyfunc.Get(server.jwksURL, keyfunc.Options{
		RefreshErrorHandler: func(err error) {
			cfg.Logger.Warn("JWKS refresh failed", zap.Error(err))
		},
		RefreshInterval:   time.Hour,
		RefreshRateLimit:  5 * time.Minute,
		RefreshTimeout:    5 * time.Second,
		RefreshUnknownKID: true,
	})
	if err != nil {
		cfg.Logger.Warn("Failed to load JWKS for local auth validation - falling back to Supabase auth API", zap.Error(err))
	} else {
		server.jwks = jwks
		cfg.Logger.Info("JWKS loaded for local auth validation", zap.String("jwks_url", server.jwksURL))
	}

	return server, nil
}

// Router returns the HTTP router with all routes configured
func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()

	// Health check (no auth required)
	mux.HandleFunc("/health", s.handleHealth)

	// Initialize GSC OAuth (non-blocking - will fail gracefully if credentials not set)
	// Determine redirect URL - use GSC_REDIRECT_URL if set, otherwise use APP_URL or localhost
	var gscRedirectURL string
	if redirectURL := os.Getenv("GSC_REDIRECT_URL"); redirectURL != "" {
		// Explicit redirect URL (e.g., https://barracuda-api-xxx.run.app/api/gsc/callback)
		gscRedirectURL = redirectURL
	} else if appURL := os.Getenv("APP_URL"); appURL != "" {
		// Use APP_URL for production (e.g., https://app.barracudaseo.com)
		// Note: This assumes the API is accessible at the same domain
		// For Cloud Run, you should set GSC_REDIRECT_URL explicitly
		gscRedirectURL = fmt.Sprintf("%s/api/gsc/callback", strings.TrimSuffix(appURL, "/"))
	} else {
		// Fallback to localhost for local development
		apiPort := os.Getenv("PORT")
		if apiPort == "" {
			apiPort = "8080"
		}
		gscRedirectURL = fmt.Sprintf("http://localhost:%s/api/gsc/callback", apiPort)
	}
	if err := gsc.InitializeOAuth(gscRedirectURL); err != nil {
		s.logger.Warn("GSC integration disabled", zap.Error(err))
		s.logger.Info("Set GSC_CLIENT_ID, GSC_CLIENT_SECRET, or GSC_CREDENTIALS_JSON to enable")
	} else {
		s.logger.Info("GSC OAuth initialized", zap.String("redirect_url", gscRedirectURL))
	}

	// Initialize Stripe (non-blocking - will fail gracefully if credentials not set)
	stripeConfig := GetStripeConfig()
	if stripeConfig.SecretKey != "" {
		InitializeStripe(stripeConfig.SecretKey)
		s.logger.Info("Stripe initialized")
	} else {
		s.logger.Warn("Stripe integration disabled - set STRIPE_SECRET_KEY to enable")
	}

	// Initialize DataForSEO (non-blocking - will fail gracefully if credentials not set)
	if _, err := dataforseo.NewClient(); err != nil {
		s.logger.Warn("DataForSEO integration disabled", zap.Error(err))
		s.logger.Info("Set DATAFORSEO_LOGIN and DATAFORSEO_PASSWORD to enable")
	} else {
		s.logger.Info("DataForSEO integration initialized")
	}

	// GSC OAuth callback (OAuth handles its own security)
	mux.HandleFunc("/api/gsc/callback", s.handleGSCCallback)
	// Internal cron endpoint for background sync (protected via shared secret)
	mux.HandleFunc("/api/internal/gsc/sync", s.handleGSCGlobalSync)

	// Stripe webhook (no auth required - verified by signature)
	mux.HandleFunc("/api/stripe/webhook", s.handleStripeWebhook)

	// Internal cron endpoint for keyword task polling (protected via shared secret)
	mux.HandleFunc("/api/internal/keywords/poll", s.handleKeywordTaskPoll)
	// Internal cron endpoint for scheduled keyword checks (protected via shared secret)
	mux.HandleFunc("/api/internal/keywords/check-scheduled", s.handleScheduledKeywordChecks)

	// API v1 routes
	v1 := http.NewServeMux()
	// Register /crawls/ FIRST (more specific) - Go's ServeMux matches longest prefix first
	// This ensures /crawls/:id matches /crawls/ before /crawls
	v1.HandleFunc("/crawls/", s.handleCrawlByID)
	// Register /crawls second for collection operations
	v1.HandleFunc("/crawls", s.handleCrawls)
	v1.HandleFunc("/projects", s.handleProjects)
	v1.HandleFunc("/projects/", s.handleProjectByID)
	v1.HandleFunc("/exports", s.handleExports)
	v1.HandleFunc("/billing/", s.handleBilling)
	v1.HandleFunc("/team/", s.handleTeam)
	v1.HandleFunc("/team", s.handleTeam)
	// AI routes
	v1.HandleFunc("/ai/issue-insight", s.handleIssueInsight)
	v1.HandleFunc("/ai/crawl-summary", s.handleCrawlSummary)
	// Integrations routes
	v1.HandleFunc("/integrations/openai-key", s.handleOpenAIKey)
	// Public report routes (authenticated)
	v1.HandleFunc("/reports/public/", s.handlePublicReportByID) // Handles DELETE for specific report
	v1.HandleFunc("/reports/public", s.handlePublicReports)     // Handles GET (list) and POST (create)
	// Keyword routes
	v1.HandleFunc("/keywords/", s.handleKeywordByID) // Handles GET, PUT, DELETE, and sub-resources
	v1.HandleFunc("/keywords", s.handleKeywords)     // Handles GET (list) and POST (create)

	// Wrap v1 routes with authentication middleware
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", s.authMiddleware(v1)))

	// Public report viewing (no auth required)
	mux.HandleFunc("/api/public/reports/", s.handleViewPublicReport)

	// Start background keyword task poller if DataForSEO is configured
	if _, err := dataforseo.NewClient(); err == nil {
		ctx := context.Background()
		s.StartKeywordTaskPoller(ctx, 15*time.Second) // Poll every 15 seconds for faster processing
		s.logger.Info("Started background keyword task poller", zap.Duration("interval", 15*time.Second))
	}

	return s.corsMiddleware(s.loggingMiddleware(mux))
}

// authMiddleware validates Supabase JWT tokens
func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow public access to team invite details endpoint (no auth required)
		if r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/details") {
			// Extract token from path (e.g., "/team/abc123/details" -> "abc123")
			path := strings.TrimPrefix(r.URL.Path, "/team/")
			path = strings.Trim(path, "/")
			parts := strings.Split(path, "/")
			if len(parts) == 2 && parts[1] == "details" {
				s.handleGetInviteDetailsPublic(w, r, parts[0])
				return
			}
		}

		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			s.respondError(w, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		// Parse Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			s.respondError(w, http.StatusUnauthorized, "Invalid Authorization header format")
			return
		}

		token := parts[1]

		// Validate token with Supabase
		// Note: Supabase Go client doesn't have built-in JWT validation
		// We'll use the Supabase REST API to verify the token
		user, err := s.validateToken(token)
		if err != nil {
			s.logger.Debug("Token validation failed", zap.Error(err))
			s.respondError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// Add user info to request context
		ctx := r.Context()
		ctx = contextWithUserID(ctx, user.ID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// validateToken validates a Supabase JWT token and returns user info
func (s *Server) validateToken(token string) (*User, error) {
	// 1. Check cache first
	if user, ok := s.checkTokenCache(token); ok {
		return user, nil
	}

	// 2. Check if another request is already validating this token (Singleflight)
	s.tokenInflightMu.Lock()
	if waitCh, ok := s.tokenInflight[token]; ok {
		// Another request is validating, wait for it
		s.tokenInflightMu.Unlock()
		<-waitCh
		// Re-check cache after waiting
		if user, ok := s.checkTokenCache(token); ok {
			return user, nil
		}
		// If still not in cache (validation failed), try one more time or return error
		// Ideally we shouldn't hammer the API if it just failed, but for robustness we'll let it fall through
		// to try again (which will likely fail again if invalid)
	} else {
		// We are the leader, claim the spot
		waitCh := make(chan struct{})
		s.tokenInflight[token] = waitCh
		s.tokenInflightMu.Unlock()

		// Ensure we clean up
		defer func() {
			s.tokenInflightMu.Lock()
			delete(s.tokenInflight, token)
			s.tokenInflightMu.Unlock()
			close(waitCh) // Wake up waiters
		}()
	}

	// 3. Validate token via Supabase Auth API
	var user *User
	var err error

	// Prefer local JWT validation to avoid hammering Supabase (reduces 401s when rate limited)
	if s.jwks != nil {
		if user, err = s.validateTokenLocally(token); err != nil {
			s.logger.Debug("Local token validation failed, falling back to Supabase", zap.Error(err))
		}
	}

	// Fall back to Supabase API validation if local validation was unavailable or failed
	if user == nil {
		user, err = s.validateTokenViaAPI(token)
		if err != nil {
			return nil, err
		}
	}

	// 4. Cache successful validation for 1 minute
	s.tokenCacheMu.Lock()
	s.tokenCache[token] = TokenCacheEntry{
		User:      user,
		ExpiresAt: time.Now().Add(1 * time.Minute),
	}
	s.tokenCacheMu.Unlock()

	return user, nil
}

// checkTokenCache helper
func (s *Server) checkTokenCache(token string) (*User, bool) {
	s.tokenCacheMu.RLock()
	defer s.tokenCacheMu.RUnlock()

	if entry, ok := s.tokenCache[token]; ok {
		if time.Now().Before(entry.ExpiresAt) {
			return entry.User, true
		}
	}
	return nil, false
}

// validateTokenLocally validates JWT signatures using cached JWKS keys to avoid rate limits
func (s *Server) validateTokenLocally(token string) (*User, error) {
	if s.jwks == nil {
		return nil, errors.New("jwks not initialized")
	}

	parsed, err := jwt.ParseWithClaims(token, &supabaseClaims{}, s.jwks.Keyfunc)
	if err != nil {
		return nil, err
	}

	claims, ok := parsed.Claims.(*supabaseClaims)
	if !ok || !parsed.Valid {
		return nil, errors.New("invalid token claims")
	}

	if claims.Subject == "" {
		return nil, errors.New("missing subject in token")
	}

	return &User{
		ID:    claims.Subject,
		Email: claims.Email,
	}, nil
}

// validateTokenViaAPI validates token by making a request to Supabase Auth API
func (s *Server) validateTokenViaAPI(token string) (*User, error) {
	// Make request to Supabase Auth API to get user
	authURL := s.config.SupabaseURL + "/auth/v1/user"
	req, err := http.NewRequest("GET", authURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("apikey", s.config.SupabaseAnonKey)

	client := &http.Client{Timeout: 10 * 1000000000} // 10 seconds
	resp, err := client.Do(req)
	if err != nil {
		s.logger.Debug("Token validation request failed",
			zap.String("url", authURL),
			zap.Error(err))
		return nil, fmt.Errorf("token validation request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read response body for error details
		bodyBytes, _ := io.ReadAll(resp.Body)
		s.logger.Debug("Token validation failed",
			zap.String("url", authURL),
			zap.Int("status", resp.StatusCode),
			zap.String("response", string(bodyBytes)),
			zap.String("supabase_url", s.config.SupabaseURL),
			zap.String("anon_key_prefix", keyPrefix))
		return nil, fmt.Errorf("token validation failed: status %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// corsMiddleware adds CORS headers
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs HTTP requests
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		s.logger.Info("HTTP request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", wrapped.statusCode),
			zap.Duration("duration", time.Since(start)),
			zap.String("remote_addr", r.RemoteAddr),
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// respondError sends a JSON error response
func (s *Server) respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

// respondJSON sends a JSON response
func (s *Server) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.logger.Error("Failed to encode JSON response", zap.Error(err))
	}
}

// User represents a Supabase user
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// supabaseClaims mirrors the JWT claims Supabase includes for auth
type supabaseClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
}
