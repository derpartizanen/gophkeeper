package repo

import (
	"context"
	"fmt"

	"github.com/derpartizanen/gophkeeper/internal/keeperctl/errors"
	"github.com/derpartizanen/gophkeeper/proto"
)

var _ Auth = (*AuthRepo)(nil)

// AuthRepo is facade to operations regarding authentication in Keeper.
type AuthRepo struct {
	client proto.AuthClient
}

// NewAuthRepo creates and initializes AuthRepo object.
func NewAuthRepo(client proto.AuthClient) *AuthRepo {
	return &AuthRepo{client}
}

// Login authenticates user in the Keeperd service.
func (r *AuthRepo) Login(ctx context.Context, username, securityKey string) (string, error) {
	req := &proto.LoginRequest{
		Username:    username,
		SecurityKey: securityKey,
	}

	resp, err := r.client.Login(ctx, req)
	if err != nil {
		return "", fmt.Errorf("AuthRepo - Login - r.client.Login: %w", errors.NewRequestError(err))
	}

	return resp.GetAccessToken(), nil
}
