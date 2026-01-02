package service

import "context"

// We create a interface so that handlers depend on interfaces instead of concrete implementations.
type AuthService interface {
	// later: Login, Register, Refresh
	Ping(ctx context.Context) error
}

// Private implementation of authService
type authService struct {
}

// We return interface so that handlers/callers don't know the concrete type.
func NewAuthService() AuthService {
	return &authService{}
}

func (s *authService) Ping(ctx context.Context) error {
	return nil
}
