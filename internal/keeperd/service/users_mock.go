package service

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/derpartizanen/gophkeeper/internal/keeperd/entity"
)

var _ Users = (*UsersServiceMock)(nil)

type UsersServiceMock struct {
	mock.Mock
}

func (m *UsersServiceMock) Register(
	ctx context.Context,
	username, securityKey string,
) (entity.AccessToken, error) {
	args := m.Called(ctx, username, securityKey)

	return args.Get(0).(entity.AccessToken), args.Error(1)
}
