package service

import "context"

type TaskService interface {
	Ping(ctx context.Context) error
}

type taskService struct {
}

func NewTaskService() TaskService {
	return &taskService{}
}

func (s *taskService) Ping(ctx context.Context) error {
	return nil
}
