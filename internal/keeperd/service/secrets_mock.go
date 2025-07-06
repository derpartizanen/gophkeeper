package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/derpartizanen/gophkeeper/internal/keeperd/entity"
	"github.com/derpartizanen/gophkeeper/proto"
)

var _ Secrets = (*SecretsServiceMock)(nil)

type SecretsServiceMock struct {
	mock.Mock
}

func (m *SecretsServiceMock) Create(
	ctx context.Context,
	owner uuid.UUID,
	name string,
	kind proto.DataKind,
	metadata, data []byte,
) (uuid.UUID, error) {
	args := m.Called(ctx, owner, name, kind, metadata, data)

	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *SecretsServiceMock) List(
	ctx context.Context,
	owner uuid.UUID,
) ([]entity.Secret, error) {
	args := m.Called(ctx, owner)

	if args.Get(0) == 1 {
		return nil, args.Error(1)
	}

	return args.Get(0).([]entity.Secret), args.Error(1)
}

func (m *SecretsServiceMock) Get(
	ctx context.Context,
	owner, id uuid.UUID,
) (*entity.Secret, error) {
	args := m.Called(ctx, owner, id)

	return args.Get(0).(*entity.Secret), args.Error(1)
}

func (m *SecretsServiceMock) Update(
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

func (m *SecretsServiceMock) Delete(
	ctx context.Context,
	owner, id uuid.UUID,
) error {
	args := m.Called(ctx, owner, id)

	return args.Error(0)
}
