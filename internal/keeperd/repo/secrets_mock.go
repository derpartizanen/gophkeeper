package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/derpartizanen/gophkeeper/internal/keeperd/entity"
	"github.com/derpartizanen/gophkeeper/proto"
)

var _ Secrets = (*SecretsRepoMock)(nil)

type SecretsRepoMock struct {
	mock.Mock
}

func (m *SecretsRepoMock) Create(
	ctx context.Context,
	owner uuid.UUID,
	name string,
	kind proto.DataKind,
	metadata, data []byte,
) (uuid.UUID, error) {
	args := m.Called(ctx, owner, name, kind, metadata, data)

	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *SecretsRepoMock) List(
	ctx context.Context,
	owner uuid.UUID,
) ([]entity.Secret, error) {
	args := m.Called(ctx, owner)

	return args.Get(0).([]entity.Secret), args.Error(1)
}

func (m *SecretsRepoMock) Get(
	ctx context.Context,
	owner, id uuid.UUID,
) (*entity.Secret, error) {
	args := m.Called(ctx, owner, id)

	return args.Get(0).(*entity.Secret), args.Error(1)
}

func (m *SecretsRepoMock) Update(
	ctx context.Context,
	owner, id uuid.UUID,
	changed []string,
	name string,
	metadata []byte,
	data []byte,
) error {
	args := m.Called(ctx, owner, id, changed, name, metadata, data)

	return args.Error(0)
}

func (m *SecretsRepoMock) Delete(
	ctx context.Context,
	owner, id uuid.UUID,
) error {
	args := m.Called(ctx, owner, id)

	return args.Error(0)
}
