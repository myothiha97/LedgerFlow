// Package config loads and validates the app settings from environment variables.
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

// Env names the environment the app runs in.
type Env string

const (
	EnvDev  Env = "dev"
	EnvProd Env = "prod"
)

// Config holds every setting the app reads at startup.
type Config struct {
	Port        int
	Env         Env
	DatabaseURL string
}

// Load reads the settings from environment variables and validates them.
// It returns every problem at once, so a bad deploy shows all fixes in one run.
func Load() (Config, error) {
	cfg := Config{
		Port: 4000,
		// Fail closed: a deploy that forgets APP_ENV gets the safe prod settings.
		Env:         EnvProd,
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}

	var errs []error

	if v := os.Getenv("PORT"); v != "" {
		port, err := strconv.Atoi(v)
		if err != nil || port < 1 || port > 65535 {
			errs = append(errs, fmt.Errorf("PORT must be a number from 1 to 65535, got %q", v))
		} else {
			cfg.Port = port
		}
	}

	if v := os.Getenv("APP_ENV"); v != "" {
		cfg.Env = Env(v)
	}
	if cfg.Env != EnvDev && cfg.Env != EnvProd {
		errs = append(errs, fmt.Errorf("APP_ENV must be %q or %q, got %q", EnvDev, EnvProd, cfg.Env))
	}

	// The value is never printed: it can hold the database password.
	if cfg.DatabaseURL == "" {
		errs = append(errs, errors.New("DATABASE_URL is required"))
	}

	if len(errs) > 0 {
		return Config{}, fmt.Errorf("invalid config: %w", errors.Join(errs...))
	}
	return cfg, nil
}
