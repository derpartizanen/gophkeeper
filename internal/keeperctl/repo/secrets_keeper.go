package repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/derpartizanen/gophkeeper/internal/keeperctl/errors"
	"github.com/derpartizanen/gophkeeper/proto"
)

var _ Secrets = (*SecretsRepo)(nil)

// SecretsRepo is facade to secrets stored in Keeper.
type SecretsRepo struct {
	client proto.SecretsClient
}

// NewSecretsRepo creates and initializes SecretsRepo object.
func NewSecretsRepo(client proto.SecretsClient) *SecretsRepo {
	return &SecretsRepo{client}
}

// Push send new secret data to the server.
func (r *SecretsRepo) Push(
	ctx context.Context,
	token, name string,
	kind proto.DataKind,
	description, payload []byte,
) (uuid.UUID, error) {
	var id uuid.UUID

	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	req := &proto.CreateSecretRequest{
		Name:     name,
		Metadata: description,
		Kind:     kind,
		Data:     payload,
	}

	resp, err := r.client.Create(ctx, req)
	if err != nil {
		return id, fmt.Errorf("SecretsRepo - Push - r.client.Create: %w", errors.NewRequestError(err))
	}

	id, err = uuid.Parse(resp.GetId())
	if err != nil {
		return id, fmt.Errorf("SecretsRepo - Push - uuid.Parse: %w", err)
	}

	return id, nil
}

// List returns list of user's secrets without data.
func (r *SecretsRepo) List(
	ctx context.Context,
	token string,
) ([]*proto.Secret, error) {
	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	req := &proto.ListSecretsRequest{}

	resp, err := r.client.List(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("SecretsRepo - List - r.client.List: %w", errors.NewRequestError(err))
	}

	return resp.Secrets, nil
}

// Get downloads full user's secret.
func (r *SecretsRepo) Get(
	ctx context.Context,
	token string,
	id uuid.UUID,
) (*proto.Secret, []byte, error) {
	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	req := &proto.GetSecretRequest{Id: id.String()}

	resp, err := r.client.Get(ctx, req)
	if err != nil {
		return nil, nil, fmt.Errorf("SecretsRepo - Get - r.client.Get: %w", errors.NewRequestError(err))
	}

	return resp.GetSecret(), resp.GetData(), nil
}

// Update changes parameters of stored secret.
func (r *SecretsRepo) Update(
	ctx context.Context,
	token string,
	id uuid.UUID,
	name string,
	description []byte,
	noDescription bool,
	data []byte,
) error {
	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	req := &proto.UpdateSecretRequest{Id: id.String()}

	mask, err := fieldmaskpb.New(req)
	if err != nil {
		return fmt.Errorf("SecretsRepo - Update - fieldmaskpb.New: %w", err)
	}

	if name != "" {
		if err := mask.Append(req, "name"); err != nil {
			return fmt.Errorf("SecretsRepo - Update - mask.Append: %w", err)
		}

		req.Name = name
	}

	if len(description) != 0 || noDescription {
		if err := mask.Append(req, "metadata"); err != nil {
			return fmt.Errorf("SecretsRepo - Update - mask.Append: %w", err)
		}

		req.Metadata = description
	}

	if len(data) != 0 {
		if err := mask.Append(req, "data"); err != nil {
			return fmt.Errorf("SecretsRepo - Update - mask.Append: %w", err)
		}

		req.Data = data
	}

	req.UpdateMask = mask

	if _, err := r.client.Update(ctx, req); err != nil {
		return fmt.Errorf("SecretsRepo - Update - r.client.Update: %w", errors.NewRequestError(err))
	}

	return nil
}

// Delete removes user's secret.
func (r *SecretsRepo) Delete(
	ctx context.Context,
	token string,
	id uuid.UUID,
) error {
	md := metadata.New(map[string]string{"authorization": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	req := &proto.DeleteSecretRequest{Id: id.String()}

	if _, err := r.client.Delete(ctx, req); err != nil {
		return fmt.Errorf("SecretsRepo - Delete - r.client.Delete: %w", errors.NewRequestError(err))
	}

	return nil
}
