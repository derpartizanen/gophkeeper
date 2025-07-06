package service

import (
	"context"
	"fmt"

	"github.com/derpartizanen/gophkeeper/internal/keeperctl/encryption"
	"github.com/derpartizanen/gophkeeper/internal/keeperctl/repo"
)

var _ Auth = (*AuthService)(nil)

// AuthService contains business logic related to authentication.
type AuthService struct {
	authRepo repo.Auth
}

// NewAuthService create and initializes new AuthService object.
func NewAuthService(auth repo.Auth) *AuthService {
	return &AuthService{auth}
}

// Login authenticates a user.
func (s *AuthService) Login(ctx context.Context, username string, key encryption.Key) (string, error) {
	securityKey := key.Hash()

	token, err := s.authRepo.Login(ctx, username, securityKey)
	if err != nil {
		return "", fmt.Errorf("login error: %w", err)
	}

	return token, nil
}
