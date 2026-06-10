// Command server is the LedgerFlow API entry point. It is wiring only (Architecture
// Guidelines §3.6): load config, open the DB pool, build the router, start the server,
// shut down gracefully. No business logic lives here.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/myothiha97/ledgerflow/backend/internal/config"
	"github.com/myothiha97/ledgerflow/backend/internal/handler"
	"github.com/myothiha97/ledgerflow/backend/internal/service"
	"github.com/myothiha97/ledgerflow/backend/internal/store"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx := context.Background()
	pool, err := store.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	// Compose the layers: concrete store → service → HTTP router.
	st := store.New(pool)
	authService := service.NewAuthService(st)
	router := handler.NewRouter(pool, authService, cfg.CookieSecure)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("listening on :%s (env=%s)", cfg.Port, cfg.Env)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	// Block until an interrupt, then drain in-flight requests.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
