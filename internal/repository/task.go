package repository

import (
	"context"

	"github.com/RethikRaj/task_manager_go/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository interface {
	Ping(ctx context.Context) error
	ListAllTasksByUser(ctx context.Context, userId int) ([]model.Task, error)
	Create(ctx context.Context, title string, userId int) (model.Task, error)
	GetByID(ctx context.Context, taskId int, userId int) (model.Task, error)
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

func (r *taskRepository) ListAllTasksByUser(ctx context.Context, userId int) ([]model.Task, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, title, created_at
		FROM tasks
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userId)
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

func (r *taskRepository) Create(ctx context.Context, title string, userId int) (model.Task, error) {
	var t model.Task

	err := r.pool.QueryRow(ctx, `
		INSERT INTO tasks (title, user_id)
		VALUES ($1, $2)
		RETURNING id, title, created_at
	`, title, userId).Scan(&t.ID, &t.Title, &t.CreatedAt)

	if err != nil {
		return model.Task{}, err
	}

	return t, nil
}

func (r *taskRepository) GetByID(ctx context.Context, taskId int, userId int) (model.Task, error) {
	query := `SELECT id, title FROM TASKS WHERE id=$1 AND user_id=$2`

	var t model.Task

	err := r.pool.QueryRow(ctx, query, taskId, userId).Scan(&t.ID, &t.Title)

	return t, err
}
