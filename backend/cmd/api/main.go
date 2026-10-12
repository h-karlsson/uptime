package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/h-karlsson/uptime/internal/checker"
	"github.com/h-karlsson/uptime/internal/store"
)

// Matches compose.yaml. Override with DATABASE_URL outside local development.
const defaultDatabaseURL = "postgres://admin:password@localhost:5432/uptime?sslmode=disable"

const shutdownTimeout = 10 * time.Second

func main() {
	if err := run(); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	// Cancelled on Ctrl+C or SIGTERM (e.g. docker stop), which starts a graceful shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := store.Open(ctx, getenv("DATABASE_URL", defaultDatabaseURL))
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := store.Migrate(ctx, pool); err != nil {
		return err
	}

	server := &http.Server{
		Addr:              getenv("HTTP_ADDR", ":8080"),
		Handler:           RegisterRoutes(pool),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	checkerDone := make(chan struct{})
	go func() {
		defer close(checkerDone)
		checker.New(pool, checker.Config{
			Workers:      getenvInt("CHECKER_WORKERS", 10),
			PollInterval: time.Second,
		}).Run(ctx)
	}()

	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("api listening", "addr", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return err
	case <-ctx.Done():
	}

	slog.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return err
	}
	if err := <-serverErrors; !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	select {
	case <-checkerDone:
	case <-shutdownCtx.Done():
		slog.Warn("checks still running at shutdown; their results are lost")
	}
	return nil
}

func getenvInt(key string, fallback int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
