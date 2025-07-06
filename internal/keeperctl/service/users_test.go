package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/derpartizanen/gophkeeper/internal/keeperctl/repo"
	"github.com/derpartizanen/gophkeeper/internal/keeperctl/service"
	"github.com/derpartizanen/gophkeeper/internal/libraries/gophtest"
)

func TestRegister(t *testing.T) {
	key := newTestKey()

	m := &repo.UsersRepoMock{}
	m.On(
		"Register",
		mock.Anything,
		gophtest.Username,
		key.Hash(),
	).
		Return(gophtest.AccessToken, nil)

	sat := service.NewUsersService(m)
	token, err := sat.Register(context.Background(), gophtest.Username, key)

	require.NoError(t, err)
	require.Equal(t, gophtest.AccessToken, token)
	m.AssertExpectations(t)
}

func TestRegisterOnRepoFailure(t *testing.T) {
	key := newTestKey()

	m := &repo.UsersRepoMock{}
	m.On(
		"Register",
		mock.Anything,
		gophtest.Username,
		key.Hash(),
	).
		Return("", gophtest.ErrUnexpected)

	sat := service.NewUsersService(m)
	_, err := sat.Register(context.Background(), gophtest.Username, key)

	require.Error(t, err)
	m.AssertExpectations(t)
}
