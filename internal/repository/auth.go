package repository

import (
	"context"

	"github.com/RethikRaj/task_manager_go/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	Create(ctx context.Context, email string, password string) (model.User, error)
	FindUserByEmail(ctx context.Context, email string) (model.User, error)
}

type authRepository struct {
	pool *pgxpool.Pool
}

func NewAuthRepository(pool *pgxpool.Pool) AuthRepository {
	return &authRepository{
		pool: pool,
	}
}

func (r *authRepository) Create(ctx context.Context, email string, password string) (model.User, error) {
	var u model.User
	query := `INSERT INTO USERS (email, password) VALUES ($1, $2) RETURNING id, email`

	err := r.pool.QueryRow(ctx, query, email, password).Scan(&u.ID, &u.Email)

	if err != nil {
		return model.User{}, err
	}

	return u, nil
}

func (r *authRepository) FindUserByEmail(ctx context.Context, email string) (model.User, error) {
	var u model.User
	query := `SELECT id, email, password FROM users WHERE email = $1`

	err := r.pool.QueryRow(ctx, query, email).Scan(&u.ID, &u.Email, &u.Password)

	if err != nil {
		return model.User{}, err
	}

	return u, nil
}
