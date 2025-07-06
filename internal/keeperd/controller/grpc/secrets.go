package grpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/derpartizanen/gophkeeper/internal/keeperd/entity"
	"github.com/derpartizanen/gophkeeper/internal/keeperd/service"
	"github.com/derpartizanen/gophkeeper/proto"
)

// SecretsServer provides implementation of the Secrets API.
type SecretsServer struct {
	proto.UnimplementedSecretsServer

	secretsService service.Secrets
}

// NewSecretsServer initializes and creates new SecretsServer.
func NewSecretsServer(secrets service.Secrets) *SecretsServer {
	return &SecretsServer{secretsService: secrets}
}

// Create creates new secret for a user.
func (s SecretsServer) Create(
	ctx context.Context,
	req *proto.CreateSecretRequest,
) (*proto.CreateSecretResponse, error) {
	owner := entity.UserFromContext(ctx)
	if owner == nil {
		return nil, status.Errorf(codes.Unauthenticated, entity.ErrInvalidCredentials.Error())
	}

	if details, ok := validateCreateSecretReq(req); !ok {
		st := composeBadRequestError(details)

		return nil, st.Err()
	}

	id, err := s.secretsService.Create(
		ctx,
		owner.ID,
		req.GetName(),
		req.GetKind(),
		req.GetMetadata(),
		req.GetData(),
	)
	if err != nil {
		if errors.Is(err, entity.ErrSecretExists) {
			return nil, status.Errorf(codes.AlreadyExists, entity.ErrSecretExists.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &proto.CreateSecretResponse{Id: id.String()}, nil
}

// List retrieves list of the secrets stored a user.
func (s SecretsServer) List(
	ctx context.Context,
	_ *proto.ListSecretsRequest,
) (*proto.ListSecretsResponse, error) {
	owner := entity.UserFromContext(ctx)
	if owner == nil {
		return nil, status.Errorf(codes.Unauthenticated, entity.ErrInvalidCredentials.Error())
	}

	data, err := s.secretsService.List(ctx, owner.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	rv := make([]*proto.Secret, 0, len(data))
	for _, val := range data {
		rv = append(rv, &proto.Secret{
			Id:       val.ID.String(),
			Name:     val.Name,
			Kind:     val.Kind,
			Metadata: val.Metadata,
		})
	}

	return &proto.ListSecretsResponse{Secrets: rv}, nil
}

// Get returns particular secret with data.
func (s SecretsServer) Get(
	ctx context.Context,
	req *proto.GetSecretRequest,
) (*proto.GetSecretResponse, error) {
	owner := entity.UserFromContext(ctx)
	if owner == nil {
		return nil, status.Errorf(codes.Unauthenticated, entity.ErrInvalidCredentials.Error())
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	secret, err := s.secretsService.Get(ctx, owner.ID, id)
	if err != nil {
		if errors.Is(err, entity.ErrSecretNotFound) {
			return nil, status.Errorf(codes.NotFound, entity.ErrSecretNotFound.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &proto.GetSecretResponse{
		Secret: &proto.Secret{
			Id:       secret.ID.String(),
			Name:     secret.Name,
			Kind:     secret.Kind,
			Metadata: secret.Metadata,
		},
		Data: secret.Data,
	}, nil
}

// Update updates particular secret stored by a user.
func (s SecretsServer) Update(
	ctx context.Context,
	req *proto.UpdateSecretRequest,
) (*proto.UpdateSecretResponse, error) {
	owner := entity.UserFromContext(ctx)
	if owner == nil {
		return nil, status.Errorf(codes.Unauthenticated, entity.ErrInvalidCredentials.Error())
	}

	id, details := validateUpdateSecretReq(req)
	if details != nil {
		st := composeBadRequestError(details)

		return nil, st.Err()
	}

	mask := req.GetUpdateMask()
	mask.Normalize()

	if err := s.secretsService.Update(
		ctx,
		owner.ID,
		id,
		mask.GetPaths(),
		req.GetName(),
		req.GetMetadata(),
		req.GetData(),
	); err != nil {
		if errors.Is(err, entity.ErrSecretNotFound) {
			return nil, status.Errorf(codes.NotFound, entity.ErrSecretNotFound.Error())
		}

		if errors.Is(err, entity.ErrSecretNameConflict) {
			return nil, status.Errorf(codes.AlreadyExists, entity.ErrSecretNameConflict.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &proto.UpdateSecretResponse{}, nil
}

// Delete removes particular secret stored by a user.
func (s SecretsServer) Delete(
	ctx context.Context,
	req *proto.DeleteSecretRequest,
) (*proto.DeleteSecretResponse, error) {
	owner := entity.UserFromContext(ctx)
	if owner == nil {
		return nil, status.Errorf(codes.Unauthenticated, entity.ErrInvalidCredentials.Error())
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := s.secretsService.Delete(ctx, owner.ID, id); err != nil {
		if errors.Is(err, entity.ErrSecretNotFound) {
			return nil, status.Errorf(codes.NotFound, entity.ErrSecretNotFound.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &proto.DeleteSecretResponse{}, nil
}
