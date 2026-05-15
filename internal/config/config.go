// Package config
package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server ServerConfig
	App    AppConfig
}

type ServerConfig struct {
	Port            int
	AllowedOrigins  []string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type AppConfig struct {
	LogLevel  string // debug | info | warn | error
	LogFormat string // text | json
	LogFile   string // if set, also write JSON logs to this path for log shippers
}

// Load reads every setting from environment variables.
// Panics on missing required values so the service fails fast at startup.
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:           getInt("SERVER_PORT", 8080),
			AllowedOrigins: splitComma(getString("CORS_ALLOWED_ORIGINS", "*")),
		},
		App: AppConfig{
			LogLevel:  getString("LOG_LEVEL", "info"),
			LogFormat: getString("LOG_FORMAT", "text"),
			LogFile:   getString("LOG_FILE", ""),
		},
	}
}

func getString(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func splitComma(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}
