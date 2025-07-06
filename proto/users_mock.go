package proto

import (
	context "context"

	"github.com/stretchr/testify/mock"
	grpc "google.golang.org/grpc"
)

var _ UsersClient = (*UsersClientMock)(nil)

type UsersClientMock struct {
	mock.Mock
}

func (m *UsersClientMock) Register(
	ctx context.Context,
	in *RegisterUserRequest,
	opts ...grpc.CallOption,
) (*RegisterUserResponse, error) {
	args := m.Called(ctx, in, opts)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*RegisterUserResponse), args.Error(1)
}
