package proto

import (
	context "context"

	"github.com/stretchr/testify/mock"
	grpc "google.golang.org/grpc"
)

var _ AuthClient = (*AuthClientMock)(nil)

type AuthClientMock struct {
	mock.Mock
}

func (m *AuthClientMock) Login(
	ctx context.Context,
	in *LoginRequest,
	opts ...grpc.CallOption,
) (*LoginResponse, error) {
	args := m.Called(ctx, in, opts)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*LoginResponse), args.Error(1)
}
