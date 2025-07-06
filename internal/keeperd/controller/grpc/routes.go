// Package grpc implements the gRPC API.
package grpc

import (
	"google.golang.org/grpc"

	"github.com/derpartizanen/gophkeeper/internal/keeperd/service"
	"github.com/derpartizanen/gophkeeper/proto"
)

// DefaultMaxMessageSize suggests limit for maximum length of gRPC message.
const DefaultMaxMessageSize = DefaultDataLimit + DefaultMetadataLimit + 2*DefaultMaxSecretNameLength

// RegisterRoutes injects new routes into the provided gRPC server.
func RegisterRoutes(server *grpc.Server, services *service.Services) {
	auth := NewAuthServer(services.Auth)
	proto.RegisterAuthServer(server, auth)

	secrets := NewSecretsServer(services.Secrets)
	proto.RegisterSecretsServer(server, secrets)

	users := NewUsersServer(services.Users)
	proto.RegisterUsersServer(server, users)
}
