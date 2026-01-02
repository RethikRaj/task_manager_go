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
	"github.com/RethikRaj/task_manager_go/internal/handler"
	"github.com/RethikRaj/task_manager_go/internal/router"
	"github.com/RethikRaj/task_manager_go/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Helloooooooooo... main entry point")

	// 1. Load .env for local development (optional, safe in prod)
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}

	// 2. Load configuration (single source of truth)
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 3. Root application context (controls entire app lifetime)
	// This context represents the entire lifetime of the application.
	// ctx → cancelled when signal arrives
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// cleanup signal listener
	defer stop()

	log.Println("task-manager starting...")

	// 4. Initialize services, handlers and router
	// services
	authService := service.NewAuthService()
	taskService := service.NewTaskService()

	// handlers
	healthHandler := handler.NewHealthHandler()
	authHandler := handler.NewAuthHandler(authService)
	taskHandler := handler.NewTaskHandler(taskService)

	// router
	router := router.NewRouter(healthHandler, authHandler, taskHandler)

	// 5. Create HTTP server with explicit configuration
	server := &http.Server{
		Addr:         cfg.HTTP.Addr,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
		Handler:      router,
	}

	// 6. Start the server in a goroutine
	// ListenAndServe blocks, so it must not run on the main goroutine.
	go func() {
		log.Printf("HTTP server started on %s\n", cfg.HTTP.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	// 7. Wait for shutdown signal
	// This blocks until ctx is cancelled by SIGINT or SIGTERM.
	<-ctx.Done()
	log.Println("shutdown signal received")

	// 8. Graceful shutdown with timeout
	// We allow in-flight requests to complete within this window.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown failed: %v", err)
	}

	log.Println("server exited properly")
}
