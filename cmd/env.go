package cmd

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func loadEnv() {
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
