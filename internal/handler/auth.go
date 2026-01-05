package handler

import (
	"encoding/json"
	"net/http"
	"strings"

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
		errResp := ErrorResponse{
			Message: "failed to decode request body",
			Code:    "INVALID_JSON_DECODE_FAILED",
			Status:  http.StatusBadRequest,
			Success: false,
		}
		SendJSONResponse(w, errResp.Status, errResp)
		return
	}

	// 2. Req validation
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	// 3. Call service layer
	user, err := h.authService.SignUp(r.Context(), req.Email, req.Password)

	if err != nil {
		errResp := ErrorResponse{
			Message: err.Error(),
			Success: false,
		}

		switch err {
		case errs.ErrInvalidCredentials:
			errResp.Code = "INVALID_CREDENTIALS"
			errResp.Status = http.StatusBadRequest
		default:
			errResp.Code = "INTERNAL_SERVER_ERROR"
			errResp.Status = http.StatusInternalServerError
		}
		SendJSONResponse(w, errResp.Status, errResp)
		return
	}

	// 4. Construct and Serialize and send repsonse

	// 4.1 Construct Response
	data := dto.SignUpResponse{
		ID:    user.ID,
		Email: user.Email,
	}

	succesResp := SuccessResponse{
		Status:  http.StatusCreated,
		Message: "User created successfully",
		Data:    data,
		Success: true,
	}

	// 4.2 Serialize and send response
	SendJSONResponse(w, succesResp.Status, succesResp)

}
