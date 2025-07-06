package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/derpartizanen/gophkeeper/internal/keeperd/entity"
	"github.com/derpartizanen/gophkeeper/internal/keeperd/repo"
	"github.com/derpartizanen/gophkeeper/internal/keeperd/service"
	"github.com/derpartizanen/gophkeeper/internal/libraries/gophtest"
)

func doRegisterUser(t *testing.T, repoErr error) (entity.AccessToken, error) {
	t.Helper()

	m := &repo.UsersRepoMock{}
	m.On(
		"Register",
		mock.Anything,
		gophtest.Username,
		gophtest.SecurityKey,
	).
		Return(uuid.New(), repoErr)

	sat := service.NewUsersService(gophtest.Secret, m)
	token, err := sat.Register(context.Background(), gophtest.Username, gophtest.SecurityKey)

	m.AssertExpectations(t)

	return token, err
}

func TestRegisterUser(t *testing.T) {
	token, err := doRegisterUser(t, nil)

	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestRegisterUserFailsIfUserExists(t *testing.T) {
	_, err := doRegisterUser(t, entity.ErrUserExists)

	require.Error(t, err)
}
