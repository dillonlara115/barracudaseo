package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Credentials stores CLI auth details for Supabase-backed API access.
type Credentials struct {
	AccessToken     string `json:"access_token"`
	RefreshToken    string `json:"refresh_token"`
	ExpiresAt       int64  `json:"expires_at"` // Unix timestamp (seconds)
	TokenType       string `json:"token_type"`
	SupabaseURL     string `json:"supabase_url"`
	SupabaseAnonKey string `json:"supabase_anon_key"`
	APIURL          string `json:"api_url"`
	UserID          string `json:"user_id,omitempty"`
	UserEmail       string `json:"user_email,omitempty"`
}

func credentialsPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to locate config dir: %w", err)
	}
	return filepath.Join(configDir, "barracuda", "credentials.json"), nil
}

// LoadCredentials reads stored CLI credentials.
func LoadCredentials() (*Credentials, error) {
	path, err := credentialsPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read credentials: %w", err)
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("failed to parse credentials: %w", err)
	}
	return &creds, nil
}

// SaveCredentials persists CLI credentials to disk.
func SaveCredentials(creds *Credentials) error {
	if creds == nil {
		return errors.New("credentials are nil")
	}
	path, err := credentialsPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}
	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode credentials: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write credentials: %w", err)
	}
	return nil
}

// ClearCredentials removes stored credentials.
func ClearCredentials() error {
	path, err := credentialsPath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove credentials: %w", err)
	}
	return nil
}

// NeedsRefresh returns true if the token is expired or about to expire.
func (c *Credentials) NeedsRefresh() bool {
	if c == nil {
		return true
	}
	if c.ExpiresAt == 0 {
		return true
	}
	return time.Now().Unix() >= c.ExpiresAt-60
}

// EnsureValidAccessToken refreshes the access token when needed.
func EnsureValidAccessToken(ctx context.Context, creds *Credentials) (string, *Credentials, error) {
	if creds == nil {
		return "", nil, errors.New("missing credentials")
	}
	if creds.AccessToken == "" || creds.RefreshToken == "" {
		return "", nil, errors.New("credentials are incomplete - run `barracuda auth login`")
	}
	if !creds.NeedsRefresh() {
		return creds.AccessToken, creds, nil
	}

	refreshed, err := refreshSupabaseToken(ctx, creds)
	if err != nil {
		return "", nil, err
	}
	if err := SaveCredentials(refreshed); err != nil {
		return "", nil, err
	}
	return refreshed.AccessToken, refreshed, nil
}

type refreshResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	ExpiresAt    int64  `json:"expires_at"`
	RefreshToken string `json:"refresh_token"`
}

func refreshSupabaseToken(ctx context.Context, creds *Credentials) (*Credentials, error) {
	if creds.SupabaseURL == "" || creds.SupabaseAnonKey == "" {
		return nil, errors.New("missing Supabase configuration in credentials")
	}

	baseURL := strings.TrimSuffix(creds.SupabaseURL, "/")
	endpoint := fmt.Sprintf("%s/auth/v1/token?grant_type=refresh_token", baseURL)

	payload := map[string]string{"refresh_token": creds.RefreshToken}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to build refresh request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", creds.SupabaseAnonKey)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("refresh failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(msg)))
	}

	var parsed refreshResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("failed to parse refresh response: %w", err)
	}

	expiresAt := parsed.ExpiresAt
	if expiresAt == 0 && parsed.ExpiresIn > 0 {
		expiresAt = time.Now().Unix() + parsed.ExpiresIn
	}

	updated := *creds
	if parsed.AccessToken != "" {
		updated.AccessToken = parsed.AccessToken
	}
	if parsed.RefreshToken != "" {
		updated.RefreshToken = parsed.RefreshToken
	}
	if parsed.TokenType != "" {
		updated.TokenType = parsed.TokenType
	}
	if expiresAt != 0 {
		updated.ExpiresAt = expiresAt
	}

	return &updated, nil
}
