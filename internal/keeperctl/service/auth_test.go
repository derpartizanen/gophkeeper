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

func TestLogin(t *testing.T) {
	key := newTestKey()

	m := &repo.AuthRepoMock{}
	m.On(
		"Login",
		mock.Anything,
		gophtest.Username,
		key.Hash(),
	).
		Return(gophtest.AccessToken, nil)

	sat := service.NewAuthService(m)
	token, err := sat.Login(context.Background(), gophtest.Username, key)

	require.NoError(t, err)
	require.Equal(t, gophtest.AccessToken, token)
	m.AssertExpectations(t)
}

func TestLoginOnRepoFailure(t *testing.T) {
	key := newTestKey()

	m := &repo.AuthRepoMock{}
	m.On(
		"Login",
		mock.Anything,
		gophtest.Username,
		key.Hash(),
	).
		Return("", gophtest.ErrUnexpected)

	sat := service.NewAuthService(m)
	_, err := sat.Login(context.Background(), gophtest.Username, key)

	require.Error(t, err)
	m.AssertExpectations(t)
}
