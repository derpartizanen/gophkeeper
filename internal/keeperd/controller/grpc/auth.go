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

// AuthServer provides implementation of the Auth API.
type AuthServer struct {
	proto.UnimplementedAuthServer

	authService service.Auth
}

// NewAuthServer initializes and creates new AuthServer.
func NewAuthServer(auth service.Auth) *AuthServer {
	return &AuthServer{authService: auth}
}

// Login authenticates a user in the service.
func (s AuthServer) Login(
	ctx context.Context,
	req *proto.LoginRequest,
) (*proto.LoginResponse, error) {
	username := req.GetUsername()
	key := req.GetSecurityKey()

	if details, ok := validateCredentials(username, key); !ok {
		st := composeBadRequestError(details)

		return nil, st.Err()
	}

	accessToken, err := s.authService.Login(ctx, username, key)
	if err != nil {
		if errors.Is(err, entity.ErrInvalidCredentials) {
			return nil, status.Errorf(codes.Unauthenticated, entity.ErrInvalidCredentials.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &proto.LoginResponse{AccessToken: accessToken.String()}, nil
}
