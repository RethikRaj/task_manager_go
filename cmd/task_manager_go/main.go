package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("Helloooooooooo... main entry point")

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
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
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
