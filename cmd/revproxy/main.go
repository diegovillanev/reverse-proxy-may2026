package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"diegovillanev/reverse-proxy-may2026/internal/config"
	"diegovillanev/reverse-proxy-may2026/internal/proxy"

	"github.com/lmittmann/tint"
	slogmulti "github.com/samber/slog-multi"
)

func main() {
	// CONFIG ##########################################################################################################
	cfg := config.Load()

	var level slog.Level
	if err := level.UnmarshalText([]byte(cfg.App.LogLevel)); err != nil {
		level = slog.LevelInfo
	}

	var logHandler slog.Handler
	if cfg.App.LogFormat == "json" {
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	} else {
		logHandler = tint.NewHandler(os.Stdout, &tint.Options{
			Level:   level,
			NoColor: os.Getenv("NO_COLOR") != "",
		})
	}

	if cfg.App.LogFile != "" {
		f, err := os.OpenFile(cfg.App.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			slog.Warn("Could not open log file, using console only", "path", cfg.App.LogFile, "error", err)
		} else {
			fileHandler := slog.NewJSONHandler(f, &slog.HandlerOptions{Level: level})
			logHandler = slogmulti.Fanout(logHandler, fileHandler)
		}
	}

	rootLogger := slog.New(logHandler)

	// CONFIG ##########################################################################################################

	u, err := url.Parse(cfg.Upstream.ProxyPass)
	if err != nil || u.Scheme == "" || u.Host == "" {
		rootLogger.Error("Bad ProxyPass", "url", cfg.Upstream.ProxyPass)
		os.Exit(1)
	}

	if cfg.Server.TTL <= 0 {
		rootLogger.Error("TTL must be greater than 0", "ttl", cfg.Server.TTL)
		os.Exit(1)
	}

	proxy := proxy.ReverseProxy{
		Upstream: u,
		Client: &http.Client{
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		Logger: rootLogger,
		Cache:  proxy.NewCache(cfg.Server.TTL),
	}

	// Setup signal context for Graceful Shutdown
	signalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Ensure in-flight requests aren't cancelled immediately on SIGTERM, for Graceful Shutdown
	ongoingCtx, stopOngoingGracefully := context.WithCancel(context.Background())

	addr := fmt.Sprintf(":%d", cfg.Server.Port)

	srv := &http.Server{
		Addr:    addr,
		Handler: &proxy,
		BaseContext: func(_ net.Listener) context.Context {
			return ongoingCtx
		},
	}

	// ListenAndServe() intentionally returns http.ErrServerClosed. This is a normal lifecycle event,
	// not an actual network failure, so we filter it out to avoid logging a false error.
	go func() {
		rootLogger.Info("Server listening", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			rootLogger.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for SIGINT/SIGTERM
	<-signalCtx.Done()
	stop()
	rootLogger.Info("Received shutdown signal, shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	err = srv.Shutdown(shutdownCtx)
	stopOngoingGracefully()
	if err != nil {
		rootLogger.Info("Failed to wait for ongoing requests to finish, waiting for forced cancellation...")
		time.Sleep(cfg.Server.ShutdownTimeout)
	}

	rootLogger.Info("Server shutdown gracefully")
}
