package grpc_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	cgrpc "github.com/derpartizanen/gophkeeper/internal/keeperd/controller/grpc"
	"github.com/derpartizanen/gophkeeper/internal/keeperd/entity"
	"github.com/derpartizanen/gophkeeper/internal/keeperd/service"
	"github.com/derpartizanen/gophkeeper/internal/libraries/gophtest"
	"github.com/derpartizanen/gophkeeper/proto"
)

func TestLoginUser(t *testing.T) {
	m := newServicesMock()
	m.Auth.(*service.AuthServiceMock).On(
		"Login",
		mock.Anything,
		gophtest.Username,
		gophtest.SecurityKey,
	).
		Return(entity.AccessToken(gophtest.AccessToken), nil)

	conn := createTestServer(t, m)

	req := &proto.LoginRequest{
		Username:    gophtest.Username,
		SecurityKey: gophtest.SecurityKey,
	}

	client := proto.NewAuthClient(conn)
	resp, err := client.Login(context.Background(), req)

	require.NoError(t, err)
	require.Equal(t, gophtest.AccessToken, resp.AccessToken)
	m.Auth.(*service.AuthServiceMock).AssertExpectations(t)
}

func TestLoginWithBadRequest(t *testing.T) {
	tt := []struct {
		name     string
		username string
		key      string
	}{
		{
			name:     "Login fails if username is empty",
			username: "",
			key:      gophtest.SecurityKey,
		},
		{
			name:     "Login fails if security key is empty",
			username: gophtest.Username,
			key:      "",
		},
		{
			name:     "Login fails if username is too long",
			username: strings.Repeat("#", cgrpc.DefaultMaxUsernameLength+1),
			key:      gophtest.SecurityKey,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			conn := createTestServer(t, newServicesMock())

			req := &proto.LoginRequest{
				Username:    tc.username,
				SecurityKey: tc.key,
			}

			client := proto.NewAuthClient(conn)
			_, err := client.Login(context.Background(), req)

			requireEqualCode(t, codes.InvalidArgument, err)
		})
	}
}

func TestLoginOnServiceFailure(t *testing.T) {
	tt := []struct {
		name       string
		serviceErr error
		expected   codes.Code
	}{
		{
			name:       "Login fails on invalid credentials",
			serviceErr: entity.ErrInvalidCredentials,
			expected:   codes.Unauthenticated,
		},
		{
			name:       "Login fails if something bad happened",
			serviceErr: gophtest.ErrUnexpected,
			expected:   codes.Internal,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := newServicesMock()
			m.Auth.(*service.AuthServiceMock).On(
				"Login",
				mock.Anything,
				gophtest.Username,
				gophtest.SecurityKey,
			).
				Return(entity.AccessToken(""), tc.serviceErr)

			conn := createTestServer(t, m)

			req := &proto.LoginRequest{
				Username:    gophtest.Username,
				SecurityKey: gophtest.SecurityKey,
			}

			client := proto.NewAuthClient(conn)
			_, err := client.Login(context.Background(), req)

			requireEqualCode(t, tc.expected, err)
			m.Auth.(*service.AuthServiceMock).AssertExpectations(t)
		})
	}
}
