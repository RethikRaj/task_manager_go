package handler

import "net/http"

type AuthHandler struct {
	// later: auth service
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// login logic here
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}
