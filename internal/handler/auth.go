package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/RethikRaj/task_manager_go/internal/common"
	"github.com/RethikRaj/task_manager_go/internal/dto"
	"github.com/RethikRaj/task_manager_go/internal/errs"
	"github.com/RethikRaj/task_manager_go/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	// 1. Deserialization
	var req dto.SignUpRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteJSONError(w, "failed to decode body", "INVALID_JSON_DECODE_FAILED", http.StatusInternalServerError)
	}

	// 2. Req validation
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	// 3. Call service layer
	user, err := h.authService.SignUp(r.Context(), req.Email, req.Password)

	if err != nil {
		switch err {
		case errs.ErrInvalidCredentials:
			common.WriteJSONError(w, err.Error(), "INVALID_CREDENTIALS", http.StatusBadRequest)
		default:
			common.WriteJSONError(w, err.Error(), "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
		}
		return
	}

	// 4. Construct and Serialize and send repsonse

	// 4.1 Construct Response
	resp := dto.SignUpResponse{
		ID:    user.ID,
		Email: user.Email,
	}

	// 4.2 Serialize response
	buf := new(bytes.Buffer)

	if err := json.NewEncoder(buf).Encode(resp); err != nil {
		common.WriteJSONError(w, "failed to encode response", "JSON_ENCODE_FAILED", http.StatusInternalServerError)
		return
	}

	// 4.3 Send Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())

}
