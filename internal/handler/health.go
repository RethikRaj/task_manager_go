package handler

import "net/http"

type HealthHandler struct {
	// no dependencies yet
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok from handler"))
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}
