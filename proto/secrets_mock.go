package proto

import (
	context "context"

	"github.com/stretchr/testify/mock"
	grpc "google.golang.org/grpc"
)

var _ SecretsClient = (*SecretsClientMock)(nil)

type SecretsClientMock struct {
	mock.Mock
}

func (m *SecretsClientMock) Create(
	ctx context.Context,
	in *CreateSecretRequest,
	opts ...grpc.CallOption,
) (*CreateSecretResponse, error) {
	args := m.Called(ctx, in, opts)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*CreateSecretResponse), args.Error(1)
}

func (m *SecretsClientMock) List(
	ctx context.Context,
	in *ListSecretsRequest,
	opts ...grpc.CallOption,
) (*ListSecretsResponse, error) {
	args := m.Called(ctx, in, opts)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*ListSecretsResponse), args.Error(1)
}

func (m *SecretsClientMock) Get(
	ctx context.Context,
	in *GetSecretRequest,
	opts ...grpc.CallOption,
) (*GetSecretResponse, error) {
	args := m.Called(ctx, in, opts)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*GetSecretResponse), args.Error(1)
}

func (m *SecretsClientMock) Update(
	ctx context.Context,
	in *UpdateSecretRequest,
	opts ...grpc.CallOption,
) (*UpdateSecretResponse, error) {
	args := m.Called(ctx, in, opts)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*UpdateSecretResponse), args.Error(1)
}

func (m *SecretsClientMock) Delete(
	ctx context.Context,
	in *DeleteSecretRequest,
	opts ...grpc.CallOption,
) (*DeleteSecretResponse, error) {
	args := m.Called(ctx, in, opts)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*DeleteSecretResponse), args.Error(1)
}
