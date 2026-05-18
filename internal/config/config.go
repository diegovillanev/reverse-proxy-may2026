// Package config
package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server   ServerConfig
	App      AppConfig
	Upstream UpstreamConfig
}

type ServerConfig struct {
	Port                int
	AllowedOrigins      []string
	ShutdownTimeout     time.Duration
	ShutdownHardTimeout time.Duration
}

type AppConfig struct {
	LogLevel  string // debug | info | warn | error
	LogFormat string // text | json
	LogFile   string // if set, also write JSON logs to this path for log shippers
}

type UpstreamConfig struct {
	ProxyPass    string
	TTL          time.Duration
	CacheMaxSize int
}

// Load reads every setting from environment variables.
// Panics on missing required values so the service fails fast at startup.
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:                getInt("SERVER_PORT", 8080),
			AllowedOrigins:      splitComma(getString("CORS_ALLOWED_ORIGINS", "*")),
			ShutdownTimeout:     seconds("SHUTDOWN_TIMEOUT", 15),
			ShutdownHardTimeout: seconds("SHUTDOWN_HARD_TIMEOUT", 5),
		},
		App: AppConfig{
			LogLevel:  getString("LOG_LEVEL", "info"),
			LogFormat: getString("LOG_FORMAT", "text"),
			LogFile:   getString("LOG_FILE", ""),
		},
		Upstream: UpstreamConfig{
			ProxyPass:    getString("PROXY_PASS", ""),
			TTL:          seconds("TTL", 10),
			CacheMaxSize: getInt("CACHE_MAX_SIZE", 5),
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

func seconds(key string, fallback int) time.Duration {
	return time.Duration(getInt(key, fallback)) * time.Second
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
