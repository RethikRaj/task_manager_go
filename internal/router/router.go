package router

import (
	"net/http"

	"github.com/RethikRaj/task_manager_go/internal/handler"
)

func NewRouter(healthHandler *handler.HealthHandler, authHandler *handler.AuthHandler, taskHandler *handler.TaskHandler) http.Handler {
	mux := http.NewServeMux()

	// health check
	mux.HandleFunc("/health", healthHandler.Check)

	mux.HandleFunc("/tasks", taskHandler.List)
	return mux
}
