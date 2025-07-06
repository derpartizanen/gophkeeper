package service

import (
	"context"
	"fmt"

	"github.com/derpartizanen/gophkeeper/internal/keeperctl/encryption"
	"github.com/derpartizanen/gophkeeper/internal/keeperctl/repo"
)

var _ Users = (*UsersService)(nil)

// UsersService contains business logic related to users management.
type UsersService struct {
	usersRepo repo.Users
}

// NewUsersService create and initializes new UsersService object.
func NewUsersService(users repo.Users) *UsersService {
	return &UsersService{users}
}

// Register creates a new user.
func (uc *UsersService) Register(
	ctx context.Context,
	username string,
	key encryption.Key,
) (string, error) {
	securityKey := key.Hash()

	accessToken, err := uc.usersRepo.Register(ctx, username, securityKey)
	if err != nil {
		return "", fmt.Errorf("UsersService - Register - uc.usersRepo.Register: %w", err)
	}

	return accessToken, nil
}
