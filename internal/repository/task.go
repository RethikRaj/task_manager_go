package repository

import "context"

type TaskRepository interface {
	Ping(ctx context.Context) error
}

type taskRepository struct {
	// later: *pgxpool.Pool
}

func NewTaskRepository() TaskRepository {
	return &taskRepository{}
}

func (r *taskRepository) Ping(ctx context.Context) error {
	return nil
}
