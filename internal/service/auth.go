package service

import (
	"context"
	"fmt"

	"github.com/RethikRaj/task_manager_go/internal/errs"
	"github.com/RethikRaj/task_manager_go/internal/model"
	"github.com/RethikRaj/task_manager_go/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// We create a interface so that handlers depend on interfaces instead of concrete implementations.
type AuthService interface {
	Ping(ctx context.Context) error
	SignUp(ctx context.Context, email string, password string) (model.User, error)
	Login(ctx context.Context, email string, password string) (string, error)
}

// Private implementation of authService
type authService struct {
	authRepo  repository.AuthRepository
	jwtSecret string
}

// We return interface so that handlers/callers don't know the concrete type.
func NewAuthService(authRepo repository.AuthRepository, jwtSecret string) AuthService {
	return &authService{
		authRepo:  authRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) Ping(ctx context.Context) error {
	return nil
}

func (s *authService) SignUp(ctx context.Context, email string, password string) (model.User, error) {

	if email == "" || len(password) < 6 || len(password) > 50 {
		return model.User{}, errs.ErrInvalidCredentials
	}

	// Check if email already exists -> DB layer unique check

	// Hash password
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return model.User{}, err
	}

	hashedPassword := string(bytes)

	user, err := s.authRepo.Create(ctx, email, hashedPassword)

	if err != nil {
		// TODO: DUPLICATE EMAIL ERROR HANDLING
		fmt.Println(err)
		return model.User{}, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, email string, password string) (string, error) {
	// Check if user exists
	user, err := s.authRepo.FindUserByEmail(ctx, email)

	if err != nil {
		return "", errs.ErrInvalidCredentials
	}

	// Check if password match
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errs.ErrInvalidCredentials
	}

	// Generate token
	token, err := GenerateToken(user.ID, s.jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}
