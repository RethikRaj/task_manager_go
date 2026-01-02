package handler

import (
	"net/http"

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

	// We will add JSON later
	for _, t := range tasks {
		w.Write([]byte(t.Title + "\n"))
	}
}
