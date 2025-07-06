package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/derpartizanen/gophkeeper/internal/keeperd/entity"
	"github.com/derpartizanen/gophkeeper/internal/keeperd/service"
	"github.com/derpartizanen/gophkeeper/proto"
)

// UsersServer provides implementation of the Users API.
type UsersServer struct {
	proto.UnimplementedUsersServer

	usersService service.Users
}

// NewUsersServer initializes and creates new UsersServer.
func NewUsersServer(users service.Users) *UsersServer {
	return &UsersServer{usersService: users}
}

// Register creates new user.
func (s UsersServer) Register(
	ctx context.Context,
	req *proto.RegisterUserRequest,
) (*proto.RegisterUserResponse, error) {
	username := req.GetUsername()
	key := req.GetSecurityKey()

	if details, ok := validateCredentials(username, key); !ok {
		st := composeBadRequestError(details)

		return nil, st.Err()
	}

	accessToken, err := s.usersService.Register(ctx, username, key)
	if err != nil {
		if errors.Is(err, entity.ErrUserExists) {
			return nil, status.Errorf(codes.AlreadyExists, entity.ErrUserExists.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &proto.RegisterUserResponse{AccessToken: accessToken.String()}, nil
}
