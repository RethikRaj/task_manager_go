package router

import (
	"net/http"

	"github.com/RethikRaj/task_manager_go/internal/handler"
)

func NewRouter(healthHandler *handler.HealthHandler, authHandler *handler.AuthHandler, taskHandler *handler.TaskHandler, authMiddleware func(http.Handler) http.Handler) http.Handler {
	mux := http.NewServeMux()

	// ---------- Public routes ----------
	// health check
	mux.HandleFunc("/health", healthHandler.Check)
	mux.HandleFunc("/auth/signup", authHandler.SignUp)
	mux.HandleFunc("/auth/signin", authHandler.SignIn)

	// ---------- Private routes ----------

	// 1. Define the handler logic
	taskHandlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.List(w, r)
		case http.MethodPost:
			taskHandler.Create(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// 2. Wrap the handler with the middleware and register it
	// mux.Handle accepts an http.Handler, which authMiddleware returns
	mux.Handle("/tasks", authMiddleware(taskHandlerFunc))

	return mux
}
