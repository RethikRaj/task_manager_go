package repository

import (
	"context"

	"github.com/RethikRaj/task_manager_go/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository interface {
	Ping(ctx context.Context) error
	List(ctx context.Context) ([]model.Task, error)
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

func (r *taskRepository) List(ctx context.Context) ([]model.Task, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, title, created_at
		FROM tasks
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task

	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
