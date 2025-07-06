package repo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/derpartizanen/gophkeeper/internal/keeperctl/repo"
	"github.com/derpartizanen/gophkeeper/internal/libraries/gophtest"
	"github.com/derpartizanen/gophkeeper/proto"
)

func newRegisterUserRequest() *proto.RegisterUserRequest {
	return &proto.RegisterUserRequest{
		Username:    gophtest.Username,
		SecurityKey: gophtest.SecurityKey,
	}
}

func TestRegister(t *testing.T) {
	resp := &proto.RegisterUserResponse{
		AccessToken: gophtest.AccessToken,
	}

	m := &proto.UsersClientMock{}
	m.On(
		"Register",
		mock.Anything,
		newRegisterUserRequest(),
		mock.Anything,
	).
		Return(resp, nil)

	sat := repo.NewUsersRepo(m)
	token, err := sat.Register(context.Background(), gophtest.Username, gophtest.SecurityKey)

	require.NoError(t, err)
	require.Equal(t, gophtest.AccessToken, token)
	m.AssertExpectations(t)
}

func TestRegisterOnClientFailure(t *testing.T) {
	m := &proto.UsersClientMock{}
	m.On(
		"Register",
		mock.Anything,
		newRegisterUserRequest(),
		mock.Anything,
	).
		Return(nil, gophtest.ErrUnexpected)

	sat := repo.NewUsersRepo(m)
	_, err := sat.Register(context.Background(), gophtest.Username, gophtest.SecurityKey)

	require.Error(t, err)
	m.AssertExpectations(t)
}
