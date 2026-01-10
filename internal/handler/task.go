package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/RethikRaj/task_manager_go/internal/ctx"
	"github.com/RethikRaj/task_manager_go/internal/dto"
	"github.com/RethikRaj/task_manager_go/internal/errs"
	"github.com/RethikRaj/task_manager_go/internal/service"
)

type TaskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

func (h *TaskHandler) ListAllTasksByUser(w http.ResponseWriter, r *http.Request) {
	user, ok := ctx.GetUserFromContext(r.Context())

	if !ok {
		// This should technically never happen if the middleware is working
		errResp := ErrorResponse{
			Status:  http.StatusUnauthorized,
			Message: "User not found in context",
			Code:    "UNAUTHORIZED",
			Success: false,
		}
		SendJSONResponse(w, errResp.Status, errResp)
		return
	}

	tasks, err := h.taskService.ListAllTasksByUser(r.Context(), user.ID)

	if err != nil {
		errResp := ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Code:    "INTERNAL_SERVER_ERROR",
			Success: false,
		}
		SendJSONResponse(w, errResp.Status, errResp)
		return
	}

	// Success Response

	successResp := SuccessResponse{
		Status:  http.StatusOK,
		Message: "Fetched all task successfully",
		Data:    tasks,
		Success: true,
	}

	SendJSONResponse(w, successResp.Status, successResp)
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	// 0. Get the user context
	user, ok := ctx.GetUserFromContext(r.Context())

	if !ok {
		// This should technically never happen if the middleware is working
		errResp := ErrorResponse{
			Status:  http.StatusUnauthorized,
			Message: "User not found in context",
			Code:    "UNAUTHORIZED",
			Success: false,
		}
		SendJSONResponse(w, errResp.Status, errResp)
		return
	}

	var req dto.CreateTaskRequest
	// 1. Deserialization
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
	req.Title = strings.TrimSpace(req.Title)

	// 3. Call service
	newTask, err := h.taskService.Create(r.Context(), req.Title, user.ID)
	if err != nil {
		// Create a base error response
		errResp := ErrorResponse{
			Message: err.Error(),
			Success: false,
		}

		switch err {
		case errs.ErrTitleRequired:
			errResp.Status = http.StatusBadRequest
			errResp.Code = "TITLE_REQUIRED"

		case errs.ErrTitleTooLong:
			errResp.Status = http.StatusBadRequest
			errResp.Code = "TITLE_TOO_LONG"

		default:
			errResp.Status = http.StatusInternalServerError
			errResp.Code = "INTERNAL_SERVER_ERROR"
		}

		// Send the response using your new buffered helper
		SendJSONResponse(w, errResp.Status, errResp)
		return
	}

	// 4. Construct , Serialize , and send suceess response
	successResp := SuccessResponse{
		Status:  http.StatusCreated,
		Message: "Created task successfully",
		Data:    newTask,
		Success: true,
	}

	SendJSONResponse(w, successResp.Status, successResp)
}

func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// 0. Get the user context
	user, ok := ctx.GetUserFromContext(r.Context())

	if !ok {
		// This should technically never happen if the middleware is working
		errResp := ErrorResponse{
			Status:  http.StatusUnauthorized,
			Message: "User not found in context",
			Code:    "UNAUTHORIZED",
			Success: false,
		}
		SendJSONResponse(w, errResp.Status, errResp)
		return
	}

	// Retrieve ID from path
	taskIdStr := r.PathValue("id")
	taskId, err := strconv.Atoi(taskIdStr)

	if err != nil {
		errResp := ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Code:    "CONVERSION_ERROR",
			Success: false,
		}
		SendJSONResponse(w, errResp.Status, errResp)
	}

	// Call Service
	task, err := h.taskService.GetByID(r.Context(), taskId, user.ID)

	if err != nil {
		errResp := ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Code:    "INTERNAL_SERVER_ERROR",
			Success: false,
		}
		SendJSONResponse(w, errResp.Status, errResp)
		return
	}

	sucessResp := SuccessResponse{
		Success: true,
		Message: "Fetched task succesfully",
		Data:    task,
		Status:  http.StatusOK,
	}

	SendJSONResponse(w, sucessResp.Status, sucessResp)

}
