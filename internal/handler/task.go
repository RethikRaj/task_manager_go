package handler

import (
	"encoding/json"
	"net/http"
	"strings"

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

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.taskService.List(r.Context())

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
	newTask, err := h.taskService.Create(r.Context(), req.Title)
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
