package service

import (
	"context"
	"fmt"

	"github.com/derpartizanen/gophkeeper/internal/keeperd/entity"
	"github.com/derpartizanen/gophkeeper/internal/keeperd/repo"
	"github.com/derpartizanen/gophkeeper/internal/libraries/creds"
)

var _ Auth = (*AuthService)(nil)

// AuthService contains business logic related to authentication.
type AuthService struct {
	secret    creds.Password
	usersRepo repo.Users
}

// NewAuthService create and initializes new AuthService object.
func NewAuthService(
	secret creds.Password,
	users repo.Users,
) *AuthService {
	return &AuthService{secret, users}
}

func (uc *AuthService) Login(
	ctx context.Context,
	username, securityKey string,
) (entity.AccessToken, error) {
	user, err := uc.usersRepo.Verify(ctx, username, securityKey)
	if err != nil {
		return "", fmt.Errorf("AuthService - Login - uc.usersRepo.Verify: %w", err)
	}

	accessToken, err := entity.NewAccessToken(user, uc.secret)
	if err != nil {
		return "", fmt.Errorf("AuthService - Login - entity.NewAccessToken: %w", err)
	}

	return accessToken, nil
}
