package repo

import (
	"context"

	"github.com/stretchr/testify/mock"
)

var _ Users = (*UsersRepoMock)(nil)

type UsersRepoMock struct {
	mock.Mock
}

func (m *UsersRepoMock) Register(
	ctx context.Context,
	username, securityKey string,
) (string, error) {
	args := m.Called(ctx, username, securityKey)

	return args.String(0), args.Error(1)
}
