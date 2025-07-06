package repo

import (
	"context"

	"github.com/stretchr/testify/mock"
)

var _ Auth = (*AuthRepoMock)(nil)

type AuthRepoMock struct {
	mock.Mock
}

func (m *AuthRepoMock) Login(
	ctx context.Context,
	username, securityKey string,
) (string, error) {
	args := m.Called(ctx, username, securityKey)

	return args.String(0), args.Error(1)
}
