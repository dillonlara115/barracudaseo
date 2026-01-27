package ga4

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dillonlara115/barracudaseo/pkg/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	analyticsadmin "google.golang.org/api/analyticsadmin/v1beta"
	analyticsdata "google.golang.org/api/analyticsdata/v1beta"
	"google.golang.org/api/option"
)

var (
	// OAuth2 config - will be initialized with credentials
	oauthConfig *oauth2.Config
	// In-memory token storage (in production, use database)
	tokenStore = make(map[string]*oauth2.Token)
	tokenMu    sync.RWMutex
	// State storage for OAuth flow
	stateStore = make(map[string]oauthState)
	stateMu    sync.RWMutex
)

type oauthState struct {
	ProjectID string
	Expires   time.Time
}

// InitializeOAuth sets up OAuth2 configuration for GA4
// Credentials can be provided via environment variables
// Users authorize Barracuda to access their Analytics data
// Separate from GSC to allow connecting different Google accounts
func InitializeOAuth(redirectURL string) error {
	// Get credentials from environment variables (required)
	clientID := os.Getenv("GA4_CLIENT_ID")
	clientSecret := os.Getenv("GA4_CLIENT_SECRET")

	// If not set, try credentials JSON
	if clientID == "" || clientSecret == "" {
		credentialsJSON := os.Getenv("GA4_CREDENTIALS_JSON")
		if credentialsJSON != "" {
			config, err := google.ConfigFromJSON([]byte(credentialsJSON), analyticsdata.AnalyticsReadonlyScope)
			if err != nil {
				return fmt.Errorf("failed to parse credentials JSON: %w", err)
			}
			config.RedirectURL = redirectURL
			// Add Admin API scope for property listing
			config.Scopes = []string{
				analyticsdata.AnalyticsReadonlyScope,
				"https://www.googleapis.com/auth/analytics.readonly",
			}
			oauthConfig = config
			return nil
		}
	}

	// Final check - if still empty, return error with helpful message
	if clientID == "" || clientSecret == "" {
		return fmt.Errorf("GA4 OAuth credentials not configured. Set environment variables:\n" +
			"\n" +
			"export GA4_CLIENT_ID='your-client-id'\n" +
			"export GA4_CLIENT_SECRET='your-client-secret'\n" +
			"\n" +
			"Or set GA4_CREDENTIALS_JSON with your full credentials JSON.\n" +
			"\n" +
			"For setup instructions, see: docs/GA4_SETUP_CHECKLIST.md")
	}

	// Use both Analytics Data API and Admin API scopes
	// Admin API is needed to list properties
	oauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			analyticsdata.AnalyticsReadonlyScope,
			"https://www.googleapis.com/auth/analytics.readonly", // Admin API scope
		},
		Endpoint: google.Endpoint,
	}

	return nil
}

// GenerateAuthURL creates an OAuth2 authorization URL and binds it to a project
func GenerateAuthURL(projectID string) (string, string, error) {
	if oauthConfig == nil {
		return "", "", fmt.Errorf("GA4 OAuth not initialized")
	}

	// Generate secure random state
	stateBytes := make([]byte, 32)
	if _, err := rand.Read(stateBytes); err != nil {
		return "", "", fmt.Errorf("failed to generate state: %w", err)
	}
	state := base64.URLEncoding.EncodeToString(stateBytes)

	// Store state with project ID
	stateMu.Lock()
	stateStore[state] = oauthState{
		ProjectID: projectID,
		Expires:   time.Now().Add(10 * time.Minute),
	}
	stateMu.Unlock()

	// Start cleanup goroutine if not already running
	go cleanupExpiredStates()

	authURL := oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	return authURL, state, nil
}

// ConsumeState validates and consumes an OAuth state, returning the project ID
func ConsumeState(state string) (string, bool) {
	stateMu.Lock()
	defer stateMu.Unlock()

	entry, exists := stateStore[state]
	if !exists {
		return "", false
	}

	if time.Now().After(entry.Expires) {
		delete(stateStore, state)
		return "", false
	}

	// Consume state (one-time use)
	delete(stateStore, state)
	return entry.ProjectID, true
}

// ExchangeCode exchanges an authorization code for a token
func ExchangeCode(code string) (*oauth2.Token, error) {
	if oauthConfig == nil {
		return nil, fmt.Errorf("GA4 OAuth not initialized")
	}

	ctx := context.Background()
	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	return token, nil
}

// StoreToken stores an OAuth token for a project
func StoreToken(projectID string, token *oauth2.Token) {
	tokenMu.Lock()
	defer tokenMu.Unlock()
	tokenStore[projectID] = token
}

// GetToken retrieves a stored token for a project
func GetToken(projectID string) (*oauth2.Token, bool) {
	tokenMu.RLock()
	token, exists := tokenStore[projectID]
	tokenMu.RUnlock()

	if !exists {
		return nil, false
	}

	// Check if token needs refresh
	if !token.Valid() {
		// Attempt to refresh
		if token.RefreshToken != "" {
			ctx := context.Background()
			ts := oauthConfig.TokenSource(ctx, token)
			newToken, err := ts.Token()
			if err == nil {
				tokenMu.Lock()
				tokenStore[projectID] = newToken
				tokenMu.Unlock()
				return newToken, true
			}
		}
		return nil, false
	}

	return token, true
}

// cleanupExpiredStates removes expired OAuth states
func cleanupExpiredStates() {
	now := time.Now()
	stateMu.Lock()
	defer stateMu.Unlock()
	for state, entry := range stateStore {
		if now.After(entry.Expires) {
			delete(stateStore, state)
		}
	}
}

// GetClient creates an authenticated HTTP client
func GetClient(projectID string) (*http.Client, error) {
	token, exists := GetToken(projectID)
	if !exists {
		return nil, fmt.Errorf("no valid token for project")
	}

	ctx := context.Background()
	client := oauthConfig.Client(ctx, token)
	return client, nil
}

// GetService creates an Analytics Data API service client
func GetService(projectID string) (*analyticsdata.Service, error) {
	client, err := GetClient(projectID)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	service, err := analyticsdata.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create analytics data service: %w", err)
	}

	return service, nil
}

// GetAdminService creates an Analytics Admin API service client
func GetAdminService(projectID string) (*analyticsadmin.Service, error) {
	client, err := GetClient(projectID)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	service, err := analyticsadmin.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create analytics admin service: %w", err)
	}

	return service, nil
}

// GetProperties lists all GA4 properties for authenticated user
func GetProperties(projectID string) ([]*models.GA4Property, error) {
	adminService, err := GetAdminService(projectID)
	if err != nil {
		return nil, err
	}

	// List all accounts first
	accountsResp, err := adminService.Accounts.List().Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}

	properties := make([]*models.GA4Property, 0)

	// For each account, list properties
	for _, account := range accountsResp.Accounts {
		propertiesResp, err := adminService.Properties.List().
			Filter(fmt.Sprintf("parent:accounts/%s", strings.TrimPrefix(account.Name, "accounts/"))).
			Do()
		if err != nil {
			// Log error but continue with other accounts
			continue
		}

		for _, property := range propertiesResp.Properties {
			properties = append(properties, &models.GA4Property{
				PropertyID:   strings.TrimPrefix(property.Name, "properties/"),
				PropertyName: property.DisplayName,
				DisplayName:  property.DisplayName,
			})
		}
	}

	return properties, nil
}
