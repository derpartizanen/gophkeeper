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

func newLoginRequest() *proto.LoginRequest {
	return &proto.LoginRequest{
		Username:    gophtest.Username,
		SecurityKey: gophtest.SecurityKey,
	}
}

func TestLogin(t *testing.T) {
	resp := &proto.LoginResponse{
		AccessToken: gophtest.AccessToken,
	}

	m := &proto.AuthClientMock{}
	m.On(
		"Login",
		mock.Anything,
		newLoginRequest(),
		mock.Anything,
	).
		Return(resp, nil)

	sat := repo.NewAuthRepo(m)
	token, err := sat.Login(context.Background(), gophtest.Username, gophtest.SecurityKey)

	require.NoError(t, err)
	require.Equal(t, gophtest.AccessToken, token)
	m.AssertExpectations(t)
}

func TestLoginOnClientFailure(t *testing.T) {
	m := &proto.AuthClientMock{}
	m.On(
		"Login",
		mock.Anything,
		newLoginRequest(),
		mock.Anything,
	).
		Return(nil, gophtest.ErrUnexpected)

	sat := repo.NewAuthRepo(m)
	_, err := sat.Login(context.Background(), gophtest.Username, gophtest.SecurityKey)

	require.Error(t, err)
	m.AssertExpectations(t)
}
