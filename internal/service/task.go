package service

import (
	"context"

	"github.com/RethikRaj/task_manager_go/internal/errs"
	"github.com/RethikRaj/task_manager_go/internal/model"
	"github.com/RethikRaj/task_manager_go/internal/repository"
)

type TaskService interface {
	Ping(ctx context.Context) error
	ListAllTasksByUser(ctx context.Context, userId int) ([]model.Task, error)
	Create(ctx context.Context, title string, userId int) (model.Task, error)
	GetByID(ctx context.Context, taskId int, userId int) (model.Task, error)
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

func (s *taskService) ListAllTasksByUser(ctx context.Context, userId int) ([]model.Task, error) {
	return s.taskRepo.ListAllTasksByUser(ctx, userId)
}

func (s *taskService) Create(ctx context.Context, title string, userId int) (model.Task, error) {
	// Validation
	if title == "" {
		return model.Task{}, errs.ErrTitleRequired
	}

	if len(title) > 200 {
		return model.Task{}, errs.ErrTitleTooLong
	}

	return s.taskRepo.Create(ctx, title, userId)
}

func (s *taskService) GetByID(ctx context.Context, taskId int, userId int) (model.Task, error) {
	return s.taskRepo.GetByID(ctx, taskId, userId)
}
