// Package config loads runtime configuration from the environment (Twelve-Factor;
// Architecture Guidelines §1.2). Nothing host-specific lives in code.
package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config is the typed view of the process environment.
type Config struct {
	Port         string
	DatabaseURL  string
	Env          string // "dev" | "prod"
	CookieSecure bool   // set true behind HTTPS; off for local http so login works.
}

// Load reads a local .env (optional — real env vars always win) and builds a Config.
func Load() (Config, error) {
	_ = godotenv.Load() // ignore "file not found": env may be set by the shell/compose.

	cfg := Config{
		Port:        getenv("PORT", "8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Env:         getenv("APP_ENV", "dev"),
	}
	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	secure, err := strconv.ParseBool(getenv("COOKIE_SECURE", "false"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid COOKIE_SECURE: %w", err)
	}
	cfg.CookieSecure = secure

	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
