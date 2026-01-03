package service

import (
	"context"

	"github.com/RethikRaj/task_manager_go/internal/errs"
	"github.com/RethikRaj/task_manager_go/internal/model"
	"github.com/RethikRaj/task_manager_go/internal/repository"
)

type TaskService interface {
	Ping(ctx context.Context) error
	List(ctx context.Context) ([]model.Task, error)
	Create(ctx context.Context, title string) (model.Task, error)
}

type taskService struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(taskRepo repository.TaskRepository) TaskService {
	return &taskService{
		taskRepo: taskRepo,
	}
}

func (s *taskService) Ping(ctx context.Context) error {
	return s.taskRepo.Ping(ctx)
}

func (s *taskService) List(ctx context.Context) ([]model.Task, error) {
	return s.taskRepo.List(ctx)
}

func (s *taskService) Create(ctx context.Context, title string) (model.Task, error) {
	// Validation
	if title == "" {
		return model.Task{}, errs.ErrTitleRequired
	}

	if len(title) > 200 {
		return model.Task{}, errs.ErrTitleTooLong
	}

	return s.taskRepo.Create(ctx, title)
}
