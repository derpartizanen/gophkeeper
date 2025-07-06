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

func TestRegisterUser(t *testing.T) {
	tt := []struct {
		name     string
		userName string
	}{
		{
			name:     "Register user",
			userName: gophtest.Username,
		},
		{
			name:     "Register user with long name",
			userName: strings.Repeat("#", cgrpc.DefaultMaxUsernameLength),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := newServicesMock()
			m.Users.(*service.UsersServiceMock).On(
				"Register",
				mock.Anything,
				tc.userName,
				gophtest.SecurityKey,
			).
				Return(entity.AccessToken(gophtest.AccessToken), nil)

			conn := createTestServer(t, m)

			req := &proto.RegisterUserRequest{
				Username:    tc.userName,
				SecurityKey: gophtest.SecurityKey,
			}

			client := proto.NewUsersClient(conn)
			resp, err := client.Register(context.Background(), req)

			require.NoError(t, err)
			require.Equal(t, gophtest.AccessToken, resp.GetAccessToken())
			m.Users.(*service.UsersServiceMock).AssertExpectations(t)
		})
	}
}

func TestRegisterUserWithBadRequest(t *testing.T) {
	tt := []struct {
		name     string
		username string
		key      string
	}{
		{
			name:     "Register user fails if username is empty",
			username: "",
			key:      gophtest.SecurityKey,
		},
		{
			name:     "Register user fails if security key is empty",
			username: gophtest.Username,
			key:      "",
		},
		{
			name:     "Register user fails if username is too long",
			username: strings.Repeat("#", cgrpc.DefaultMaxUsernameLength+1),
			key:      gophtest.SecurityKey,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			conn := createTestServer(t, newServicesMock())

			req := &proto.RegisterUserRequest{
				Username:    tc.username,
				SecurityKey: tc.key,
			}

			client := proto.NewUsersClient(conn)
			_, err := client.Register(context.Background(), req)

			requireEqualCode(t, codes.InvalidArgument, err)
		})
	}
}

func TestRegisterUserOnServiceFailure(t *testing.T) {
	tt := []struct {
		name       string
		serviceErr error
		expected   codes.Code
	}{
		{
			name:       "Register user fails if user already exists",
			serviceErr: entity.ErrUserExists,
			expected:   codes.AlreadyExists,
		},
		{
			name:       "Register user fails if something bad happened",
			serviceErr: gophtest.ErrUnexpected,
			expected:   codes.Internal,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := newServicesMock()
			m.Users.(*service.UsersServiceMock).On(
				"Register",
				mock.Anything,
				gophtest.Username,
				gophtest.SecurityKey,
			).
				Return(entity.AccessToken(""), tc.serviceErr)

			conn := createTestServer(t, m)

			req := &proto.RegisterUserRequest{
				Username:    gophtest.Username,
				SecurityKey: gophtest.SecurityKey,
			}

			client := proto.NewUsersClient(conn)
			_, err := client.Register(context.Background(), req)

			requireEqualCode(t, tc.expected, err)
			m.Users.(*service.UsersServiceMock).AssertExpectations(t)
		})
	}
}
