package router

import (
	"net/http"

	"github.com/RethikRaj/task_manager_go/internal/handler"
)

func NewRouter(healthHandler *handler.HealthHandler) http.Handler {
	mux := http.NewServeMux()

	// health check
	mux.HandleFunc("/health", healthHandler.Check)

	return mux
}
