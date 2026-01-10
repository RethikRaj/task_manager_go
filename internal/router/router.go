package router

import (
	"net/http"

	"github.com/RethikRaj/task_manager_go/internal/handler"
)

func NewRouter(healthHandler *handler.HealthHandler, authHandler *handler.AuthHandler, taskHandler *handler.TaskHandler, authMiddleware func(http.Handler) http.Handler) http.Handler {
	mux := http.NewServeMux()

	// ---------- Public routes ----------
	// health check
	mux.HandleFunc("GET /health", healthHandler.Check)
	mux.HandleFunc("POST /auth/signup", authHandler.SignUp)
	mux.HandleFunc("POST /auth/signin", authHandler.SignIn)

	// ---------- Private routes ----------
	// We wrap each specific method/handler with the middleware

	mux.Handle("GET /tasks", authMiddleware(http.HandlerFunc(taskHandler.ListAllTasksByUser)))
	mux.Handle("POST /tasks", authMiddleware(http.HandlerFunc(taskHandler.Create)))
	// Use {id} to define a path parameter
	mux.Handle("GET /tasks/{id}", authMiddleware(http.HandlerFunc(taskHandler.GetByID)))

	return mux
}
