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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// The below line convert tasks(a slice of go structs variable) to JSON and stream it to the HTTP response(w)
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTaskRequest
	// 1. Deserialization
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	// 2. Req validation
	req.Title = strings.TrimSpace(req.Title)

	// 3. Call service
	newTask, err := h.taskService.Create(r.Context(), req.Title)
	if err != nil {
		switch err {
		case errs.ErrTitleRequired, errs.ErrTitleTooLong:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	// 4. Success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(newTask); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
