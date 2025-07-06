package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/derpartizanen/gophkeeper/internal/keeperd/entity"
)

var _ Users = (*UsersRepoMock)(nil)

type UsersRepoMock struct {
	mock.Mock
}

func (m *UsersRepoMock) Register(
	ctx context.Context,
	username, securityKey string,
) (uuid.UUID, error) {
	args := m.Called(ctx, username, securityKey)

	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *UsersRepoMock) Verify(
	ctx context.Context,
	username, securityKey string,
) (entity.User, error) {
	args := m.Called(ctx, username, securityKey)

	return args.Get(0).(entity.User), args.Error(1)
}
