package cmd

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func loadEnv() {
	if envFiles := strings.TrimSpace(os.Getenv("BARRACUDA_ENV_FILE")); envFiles != "" {
		for _, path := range splitEnvFiles(envFiles) {
			if path == "" {
				continue
			}
			_ = godotenv.Overload(path)
		}
		return
	}

	if !isTruthy(os.Getenv("BARRACUDA_LOAD_ENV")) {
		return
	}

	_ = godotenv.Load()
	if os.Getenv("PORT") == "" {
		_ = godotenv.Overload(".env.local")
	}
}

func resolveSupabaseConfig() (string, string) {
	supabaseURL := firstEnv(
		"PUBLIC_SUPABASE_URL",
		"VITE_PUBLIC_SUPABASE_URL",
		"SUPABASE_URL",
	)
	supabaseAnonKey := firstEnv(
		"PUBLIC_SUPABASE_ANON_KEY",
		"VITE_PUBLIC_SUPABASE_ANON_KEY",
		"SUPABASE_ANON_KEY",
	)
	return supabaseURL, supabaseAnonKey
}

func resolveAPIURL() string {
	return defaultEnv("http://localhost:8080",
		"BARRACUDA_API_URL",
		"CLOUD_RUN_API_URL",
		"VITE_CLOUD_RUN_API_URL",
		"API_URL",
	)
}

func resolveAppURL() string {
	return defaultEnv("https://app.barracudaseo.com",
		"APP_URL",
		"BARRACUDA_APP_URL",
	)
}

func firstEnv(keys ...string) string {
	for _, key := range keys {
		if val := strings.TrimSpace(os.Getenv(key)); val != "" {
			return val
		}
	}
	return ""
}

func defaultEnv(fallback string, keys ...string) string {
	if val := firstEnv(keys...); val != "" {
		return val
	}
	return fallback
}

func isTruthy(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "y":
		return true
	default:
		return false
	}
}

func splitEnvFiles(value string) []string {
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == ';' || r == ':' || r == '|'
	})
	return parts
}
