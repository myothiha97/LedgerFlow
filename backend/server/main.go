// Command server runs the LedgerFlow API. This file only builds and wires the parts.
package main

import (
	"context"
	"errors"
	"fmt"
	"ledgerflow/backend/internal/config"
	"ledgerflow/backend/internal/handler"
	"ledgerflow/backend/internal/store"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(logger); err != nil {
		logger.Error("server stopped", slog.Any("error", err))
		os.Exit(1)
	}
}

// run returns errors instead of exiting, so deferred cleanup (closing the pool) always runs.
func run(logger *slog.Logger) error {
	// load env variables
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	// setting up prod env
	if cfg.Env == config.EnvProd {
		gin.SetMode(gin.ReleaseMode)
	}

	// ctx is cancelled on Ctrl+C (SIGINT) or SIGTERM, which Docker sends on stop.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// initiate connecting to db
	pool, err := store.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	// defining the go server using builtin http package
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.Port),
		Handler: handler.NewRouter(logger, pool),
		// Timeouts stop slow or stuck clients from holding connections forever.
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		// Send net/http's own errors (such as TLS handshake failures) through slog as JSON.
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	// create the pipe object called serveErr to store error infos and create goroutine function and then start listening the http request in new goroutine
	// if there are any error occur during listening, return put the errors infos from srv.ListenAndServe to serveErr
	serveErr := make(chan error, 1)
	go func() {
		logger.Info("server listening", slog.Int("port", cfg.Port), slog.String("env", string(cfg.Env)))
		serveErr <- srv.ListenAndServe()
	}()

	select {
	case err := <-serveErr:
		// ListenAndServe only returns early on a real failure, such as a port in use.
		return fmt.Errorf("listen: %w", err)
	case <-ctx.Done():
	}

	select {
	case err := <-serveErr:
		// ListenAndServe only returns early on a real failure, such as a port in use.
		return fmt.Errorf("listen: %w", err)
	case <-ctx.Done():
	}
	// Stop accepting new requests and give in-flight ones time to finish.
	logger.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	if err := <-serveErr; !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen: %w", err)
	}
	logger.Info("server stopped cleanly")
	return nil
}
