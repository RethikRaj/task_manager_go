package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RethikRaj/task_manager_go/internal/config"
)

func main() {
	log.Println("Helloooooooooo... main entry point")

	// 1. Load config
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Let's build a server
	// This context represents the entire lifetime of the application.
	// ctx → cancelled when signal arrives
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// cleanup signal listener
	defer stop()

	log.Println("task-manager starting...")

	// 2. HTTP server configuration
	// We use http.Server instead of http.ListenAndServe
	// so we can control timeouts and perform graceful shutdown.
	server := &http.Server{
		Addr:         cfg.HTTP.Addr,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	// 3. Start the server in a goroutine
	// ListenAndServe blocks, so it must not run on the main goroutine.
	go func() {
		log.Println("HTTP server started on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	// 4. Wait for shutdown signal
	// This blocks until ctx is cancelled by SIGINT or SIGTERM.
	<-ctx.Done()
	log.Println("shutdown signal received")

	// 5. Graceful shutdown with timeout
	// We allow in-flight requests to complete within this window.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown failed: %v", err)
	}

	log.Println("server exited properly")
}
