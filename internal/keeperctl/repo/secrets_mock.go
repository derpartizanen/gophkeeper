package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/derpartizanen/gophkeeper/proto"
)

var _ Secrets = (*SecretsRepoMock)(nil)

type SecretsRepoMock struct {
	mock.Mock
}

func (m *SecretsRepoMock) Push(
	ctx context.Context,
	token, name string,
	kind proto.DataKind,
	description, payload []byte,
) (uuid.UUID, error) {
	args := m.Called(ctx, token, name, kind, description, payload)

	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *SecretsRepoMock) List(
	ctx context.Context,
	token string,
) ([]*proto.Secret, error) {
	args := m.Called(ctx, token)

	return args.Get(0).([]*proto.Secret), args.Error(1)
}

func (m *SecretsRepoMock) Get(
	ctx context.Context,
	token string,
	id uuid.UUID,
) (*proto.Secret, []byte, error) {
	args := m.Called(ctx, token, id)

	return args.Get(0).(*proto.Secret), args.Get(1).([]byte), args.Error(2)
}

func (m *SecretsRepoMock) Update(
	ctx context.Context,
	token string,
	id uuid.UUID,
	name string,
	description []byte,
	noDescription bool,
	data []byte,
) error {
	args := m.Called(ctx, token, id, name, description, noDescription, data)

	return args.Error(0)
}

func (m *SecretsRepoMock) Delete(
	ctx context.Context,
	token string,
	id uuid.UUID,
) error {
	args := m.Called(ctx, token, id)

	return args.Error(0)
}
