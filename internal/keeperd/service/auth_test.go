package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/derpartizanen/gophkeeper/internal/keeperd/entity"
	"github.com/derpartizanen/gophkeeper/internal/keeperd/service"
	"github.com/derpartizanen/gophkeeper/internal/libraries/gophtest"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/derpartizanen/gophkeeper/internal/keeperd/repo"
)

func doLogin(t *testing.T, repoErr error) (entity.AccessToken, error) {
	t.Helper()

	m := &repo.UsersRepoMock{}
	m.On(
		"Verify",
		mock.Anything,
		gophtest.Username,
		gophtest.SecurityKey,
	).
		Return(entity.User{ID: uuid.New(), Username: gophtest.Username}, repoErr)

	sat := service.NewAuthService(gophtest.Secret, m)
	accessToken, err := sat.Login(context.Background(), gophtest.Username, gophtest.SecurityKey)

	m.AssertExpectations(t)

	return accessToken, err
}

func TestLogin(t *testing.T) {
	token, err := doLogin(t, nil)

	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestLoginOnBadCredentials(t *testing.T) {
	_, err := doLogin(t, entity.ErrInvalidCredentials)

	require.ErrorIs(t, err, entity.ErrInvalidCredentials)
}
