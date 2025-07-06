package grpc_test

import (
	"context"
	"net"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	cgrpc "github.com/derpartizanen/gophkeeper/internal/keeperd/controller/grpc"
	"github.com/derpartizanen/gophkeeper/internal/keeperd/entity"
	"github.com/derpartizanen/gophkeeper/internal/keeperd/service"
	"github.com/derpartizanen/gophkeeper/internal/libraries/gophtest"
	"github.com/derpartizanen/gophkeeper/internal/logger"
)

func requireEqualCode(t *testing.T, expected codes.Code, err error) {
	t.Helper()

	rv, ok := status.FromError(err)

	require.True(t, ok)
	require.Equal(t, expected, rv.Code())
}

func newServicesMock() service.Services {
	return service.Services{
		Auth:    &service.AuthServiceMock{},
		Secrets: &service.SecretsServiceMock{},
		Users:   &service.UsersServiceMock{},
	}
}

func fakeAuthInterceptor(
	ctx context.Context,
	req any,
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	user := entity.User{
		ID:       uuid.New(),
		Username: gophtest.Username,
	}

	return handler(user.WithContext(ctx), req)
}

func createTestServer(
	t *testing.T,
	services service.Services,
	opts ...grpc.ServerOption,
) *grpc.ClientConn {
	t.Helper()
	require := require.New(t)

	srvOpts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(cgrpc.DefaultMaxMessageSize),
	}
	srvOpts = append(srvOpts, opts...)

	srv := grpc.NewServer(srvOpts...)
	cgrpc.RegisterRoutes(srv, &services)

	lis := bufconn.Listen(1024 * 1024)
	go func() {
		require.NoError(srv.Serve(lis))
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, err := grpc.Dial(
		"",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(err)

	t.Cleanup(func() {
		require.NoError(conn.Close())
		srv.Stop()
		require.NoError(lis.Close())
	})

	return conn
}

func createTestServerWithFakeAuth(
	t *testing.T,
	services service.Services,
) *grpc.ClientConn {
	t.Helper()

	log, err := logger.New("info")
	require.NoError(t, err)

	return createTestServer(
		t,
		services,
		grpc.ChainUnaryInterceptor(
			cgrpc.LoggingUnaryInterceptor(log),
			fakeAuthInterceptor,
		),
	)
}
