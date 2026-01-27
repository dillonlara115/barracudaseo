package cmd

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dillonlara115/barracuda/internal/utils"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate the CLI with your Barracuda account",
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Sign in and link this CLI to your account",
	RunE:  runAuthLogin,
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current CLI authentication status",
	RunE:  runAuthStatus,
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove stored CLI credentials",
	RunE:  runAuthLogout,
}

func init() {
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(authLogoutCmd)
	rootCmd.AddCommand(authCmd)
}

type cliAuthPayload struct {
	AccessToken     string `json:"access_token"`
	RefreshToken    string `json:"refresh_token"`
	ExpiresAt       int64  `json:"expires_at"`
	TokenType       string `json:"token_type"`
	SupabaseURL     string `json:"supabase_url"`
	SupabaseAnonKey string `json:"supabase_anon_key"`
	APIURL          string `json:"api_url"`
	UserID          string `json:"user_id"`
	UserEmail       string `json:"user_email"`
	State           string `json:"state"`
}

func runAuthLogin(cmd *cobra.Command, args []string) error {
	loadEnv()

	supabaseURL, supabaseAnonKey := resolveSupabaseConfig()
	if supabaseURL == "" || supabaseAnonKey == "" {
		return errors.New("missing Supabase configuration. Set PUBLIC_SUPABASE_URL and PUBLIC_SUPABASE_ANON_KEY")
	}

	appURL := resolveAppURL()
	apiURL := resolveAPIURL()

	state, err := generateState()
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("failed to start callback listener: %w", err)
	}
	defer listener.Close()

	callbackURL := fmt.Sprintf("http://%s/callback", listener.Addr().String())

	done := make(chan error, 1)
	mux := http.NewServeMux()
	server := &http.Server{Handler: mux}
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)
		var payload cliAuthPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			select {
			case done <- fmt.Errorf("invalid payload: %w", err):
			default:
			}
			return
		}

		if payload.State == "" || payload.State != state {
			http.Error(w, "Invalid state", http.StatusUnauthorized)
			select {
			case done <- errors.New("invalid state returned from auth callback"):
			default:
			}
			return
		}

		if payload.AccessToken == "" || payload.RefreshToken == "" {
			http.Error(w, "Missing tokens", http.StatusBadRequest)
			select {
			case done <- errors.New("missing tokens in auth callback"):
			default:
			}
			return
		}

		if payload.SupabaseURL == "" {
			payload.SupabaseURL = supabaseURL
		}
		if payload.SupabaseAnonKey == "" {
			payload.SupabaseAnonKey = supabaseAnonKey
		}
		if payload.APIURL == "" {
			payload.APIURL = apiURL
		}

		creds := &utils.Credentials{
			AccessToken:     payload.AccessToken,
			RefreshToken:    payload.RefreshToken,
			ExpiresAt:       payload.ExpiresAt,
			TokenType:       payload.TokenType,
			SupabaseURL:     payload.SupabaseURL,
			SupabaseAnonKey: payload.SupabaseAnonKey,
			APIURL:          payload.APIURL,
			UserID:          payload.UserID,
			UserEmail:       payload.UserEmail,
		}

		if err := utils.SaveCredentials(creds); err != nil {
			http.Error(w, "Failed to save credentials", http.StatusInternalServerError)
			select {
			case done <- err:
			default:
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "linked",
		})

		select {
		case done <- nil:
		default:
		}
	})

	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			select {
			case done <- err:
			default:
			}
		}
	}()

	loginURL := fmt.Sprintf("%s/#/cli-auth?callback=%s&state=%s", strings.TrimSuffix(appURL, "/"), url.QueryEscape(callbackURL), url.QueryEscape(state))

	fmt.Println("Opening browser to link your CLI...")
	if err := openBrowserURL(loginURL); err != nil {
		fmt.Printf("Open this URL in your browser to continue:\n%s\n", loginURL)
	}

	select {
	case err := <-done:
		_ = server.Shutdown(context.Background())
		if err != nil {
			return err
		}
		fmt.Println("✓ CLI successfully linked to your account.")
		return nil
	case <-time.After(10 * time.Minute):
		_ = server.Shutdown(context.Background())
		return errors.New("login timed out")
	}
}

func runAuthStatus(cmd *cobra.Command, args []string) error {
	loadEnv()
	creds, err := utils.LoadCredentials()
	if err != nil {
		return err
	}
	if creds == nil {
		fmt.Println("Not authenticated. Run `barracuda auth login`.")
		return nil
	}

	fmt.Println("Authenticated CLI session:")
	if creds.UserEmail != "" {
		fmt.Printf("- Email: %s\n", creds.UserEmail)
	}
	fmt.Printf("- API URL: %s\n", creds.APIURL)
	if creds.ExpiresAt > 0 {
		expiresAt := time.Unix(creds.ExpiresAt, 0)
		fmt.Printf("- Token expires: %s\n", expiresAt.Local().Format(time.RFC1123))
	} else {
		fmt.Println("- Token expires: unknown (will refresh on next request)")
	}
	return nil
}

func runAuthLogout(cmd *cobra.Command, args []string) error {
	if err := utils.ClearCredentials(); err != nil {
		return err
	}
	fmt.Println("Signed out. Stored credentials removed.")
	return nil
}

func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}
	return hex.EncodeToString(b), nil
}
