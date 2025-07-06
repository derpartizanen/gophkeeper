package service

import (
	"context"
	"fmt"

	"github.com/derpartizanen/gophkeeper/internal/keeperd/entity"
	"github.com/derpartizanen/gophkeeper/internal/keeperd/repo"
	"github.com/derpartizanen/gophkeeper/internal/libraries/creds"
)

var _ Users = (*UsersService)(nil)

// UsersService contains business logic related to users management.
type UsersService struct {
	secret    creds.Password
	usersRepo repo.Users
}

// NewUsersService create and initializes new UsersService object.
func NewUsersService(secret creds.Password, users repo.Users) *UsersService {
	return &UsersService{secret, users}
}

// Register creates a new user.
func (uc UsersService) Register(
	ctx context.Context,
	username, securityKey string,
) (entity.AccessToken, error) {
	id, err := uc.usersRepo.Register(ctx, username, securityKey)
	if err != nil {
		return "", fmt.Errorf("UsersService - Register - uc.usersRepo.Register: %w", err)
	}

	user := entity.User{
		ID:       id,
		Username: username,
	}

	accessToken, err := entity.NewAccessToken(user, uc.secret)
	if err != nil {
		return "", fmt.Errorf("AuthService - Login - entity.NewAccessToken: %w", err)
	}

	return accessToken, nil
}
