package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository interface {
	Ping(ctx context.Context) error
}

type taskRepository struct {
	pool *pgxpool.Pool
}

func NewTaskRepository(pool *pgxpool.Pool) TaskRepository {
	return &taskRepository{
		pool: pool,
	}
}

func (r *taskRepository) Ping(ctx context.Context) error {
	return nil
}
