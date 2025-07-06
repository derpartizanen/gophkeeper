package repo

import (
	"context"
	"fmt"

	"github.com/derpartizanen/gophkeeper/internal/keeperctl/errors"
	"github.com/derpartizanen/gophkeeper/proto"
)

var _ Users = (*UsersRepo)(nil)

// UsersRepo is facade to operations regarding Keeper.
type UsersRepo struct {
	client proto.UsersClient
}

// NewUsersRepo creates and initializes UsersRepo object.
func NewUsersRepo(client proto.UsersClient) *UsersRepo {
	return &UsersRepo{client}
}

// Register creates a new user.
func (r *UsersRepo) Register(
	ctx context.Context,
	username, securityKey string,
) (string, error) {
	req := &proto.RegisterUserRequest{
		Username:    username,
		SecurityKey: securityKey,
	}

	resp, err := r.client.Register(ctx, req)
	if err != nil {
		return "", fmt.Errorf("UsersRepo - Register - r.client.Register: %w", errors.NewRequestError(err))
	}

	return resp.GetAccessToken(), nil
}
