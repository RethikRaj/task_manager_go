package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/RethikRaj/task_manager_go/internal/common"
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
		common.WriteJSONError(w, err.Error(), "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// The below line convert tasks(a slice of go structs variable) to JSON and stream it to the HTTP response(w)
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		common.WriteJSONError(w, "failed to encode response", "FAILED_TO_ENCODE_RESPONSE", http.StatusInternalServerError)
	}
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTaskRequest
	// 1. Deserialization
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteJSONError(w, "failed to decode request body", "INVALID_JSON_DECODE_FAILED", http.StatusBadRequest)
		return
	}

	// 2. Req validation
	req.Title = strings.TrimSpace(req.Title)

	// 3. Call service
	newTask, err := h.taskService.Create(r.Context(), req.Title)
	if err != nil {
		switch err {
		case errs.ErrTitleRequired:
			common.WriteJSONError(w, err.Error(), "TITLE_REQUIRED", http.StatusBadRequest)
		case errs.ErrTitleTooLong:
			common.WriteJSONError(w, err.Error(), "TITLE_TOO_LONG", http.StatusBadRequest)
		default:
			common.WriteJSONError(w, err.Error(), "INTERNAL_SERVER_ERROR", http.StatusInternalServerError)
		}
		return
	}

	// 4. Success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(newTask); err != nil {
		common.WriteJSONError(w, "failed to encode response", "FAILED_TO_ENCODE_RESPONSE", http.StatusInternalServerError)
	}
}
