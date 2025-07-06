package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/derpartizanen/gophkeeper/internal/keeperd/entity"
	"github.com/derpartizanen/gophkeeper/internal/keeperd/repo"
	"github.com/derpartizanen/gophkeeper/proto"
)

var _ Secrets = (*SecretsService)(nil)

// SecretsService contains business logic related to secrets management.
type SecretsService struct {
	secretsRepo repo.Secrets
}

// NewSecretsService create and initializes new SecretsService object.
func NewSecretsService(secrets repo.Secrets) *SecretsService {
	return &SecretsService{secrets}
}

// Create creates new secret.
func (uc *SecretsService) Create(
	ctx context.Context,
	owner uuid.UUID,
	name string,
	kind proto.DataKind,
	metadata, data []byte,
) (uuid.UUID, error) {
	id, err := uc.secretsRepo.Create(ctx, owner, name, kind, metadata, data)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("SecretsService - Create - uc.secretsRepo.Create: %w", err)
	}

	return id, nil
}

// List returns list of user's secrets.
func (uc *SecretsService) List(
	ctx context.Context,
	owner uuid.UUID,
) ([]entity.Secret, error) {
	secrets, err := uc.secretsRepo.List(ctx, owner)
	if err != nil {
		return nil, fmt.Errorf("SecretsService - List - uc.secretsRepo.List: %w", err)
	}

	return secrets, nil
}

// Get retrieves full secret info from database.
func (uc *SecretsService) Get(
	ctx context.Context,
	owner, id uuid.UUID,
) (*entity.Secret, error) {
	secret, err := uc.secretsRepo.Get(ctx, owner, id)
	if err != nil {
		return nil, fmt.Errorf("SecretsService - Get - uc.secretsRepo.Get: %w", err)
	}

	return secret, nil
}

// Update changes secret info and data.
func (uc *SecretsService) Update(
	ctx context.Context,
	owner, id uuid.UUID,
	changed []string,
	name string,
	metadata []byte,
	data []byte,
) error {
	if err := uc.secretsRepo.Update(ctx, owner, id, changed, name, metadata, data); err != nil {
		return fmt.Errorf("SecretsService - Update - uc.secretsRepo.Update: %w", err)
	}

	return nil
}

// Delete removes secret owned by user.
func (uc *SecretsService) Delete(
	ctx context.Context,
	owner, id uuid.UUID,
) error {
	if err := uc.secretsRepo.Delete(ctx, owner, id); err != nil {
		return fmt.Errorf("SecretsService - Delete - uc.secretsRepo.Delete: %w", err)
	}

	return nil
}
