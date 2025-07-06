package grpc_test

import (
	"context"
	"strings"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	cgrpc "github.com/derpartizanen/gophkeeper/internal/keeperd/controller/grpc"
	"github.com/derpartizanen/gophkeeper/internal/keeperd/entity"
	"github.com/derpartizanen/gophkeeper/internal/keeperd/service"
	"github.com/derpartizanen/gophkeeper/internal/libraries/gophtest"
	"github.com/derpartizanen/gophkeeper/proto"
)

func doListSecrets(
	t *testing.T,
	mockRV []entity.Secret,
	mockErr error,
) (*proto.ListSecretsResponse, error) {
	t.Helper()

	m := newServicesMock()
	m.Secrets.(*service.SecretsServiceMock).On(
		"List",
		mock.Anything,
		mock.AnythingOfType("uuid.UUID"),
	).
		Return(mockRV, mockErr)

	conn := createTestServerWithFakeAuth(t, m)
	req := &proto.ListSecretsRequest{}

	client := proto.NewSecretsClient(conn)
	rv, err := client.List(context.Background(), req)

	m.Secrets.(*service.SecretsServiceMock).AssertExpectations(t)

	return rv, err
}

func doGetSecret(
	t *testing.T,
	mockRV *entity.Secret,
	mockErr error,
) (*proto.GetSecretResponse, error) {
	t.Helper()

	m := newServicesMock()
	m.Secrets.(*service.SecretsServiceMock).On(
		"Get",
		mock.Anything,
		mock.AnythingOfType("uuid.UUID"),
		mock.AnythingOfType("uuid.UUID"),
	).
		Return(mockRV, mockErr)

	conn := createTestServerWithFakeAuth(t, m)
	req := &proto.GetSecretRequest{Id: uuid.New().String()}

	client := proto.NewSecretsClient(conn)
	rv, err := client.Get(context.Background(), req)

	m.Secrets.(*service.SecretsServiceMock).AssertExpectations(t)

	return rv, err
}

func doDeleteSecret(
	t *testing.T,
	mockErr error,
) (*proto.DeleteSecretResponse, error) {
	t.Helper()

	m := newServicesMock()
	m.Secrets.(*service.SecretsServiceMock).On(
		"Delete",
		mock.Anything,
		mock.AnythingOfType("uuid.UUID"),
		mock.AnythingOfType("uuid.UUID"),
	).
		Return(mockErr)

	conn := createTestServerWithFakeAuth(t, m)
	req := &proto.DeleteSecretRequest{Id: uuid.New().String()}

	client := proto.NewSecretsClient(conn)
	rv, err := client.Delete(context.Background(), req)

	m.Secrets.(*service.SecretsServiceMock).AssertExpectations(t)

	return rv, err
}

func TestCreateSecret(t *testing.T) {
	tt := []struct {
		name       string
		secretName string
		metadata   []byte
		data       []byte
	}{
		{
			name:       "Create secret",
			secretName: gophtest.SecretName,
			metadata:   []byte(gophtest.Metadata),
			data:       []byte(gophtest.TextData),
		},
		{
			name:       "Create secret without metadata",
			secretName: gophtest.Username,
			metadata:   nil,
			data:       []byte(gophtest.TextData),
		},
		{
			name:       "Create secret of maximum size",
			secretName: strings.Repeat("#", cgrpc.DefaultMaxSecretNameLength),
			metadata:   []byte(strings.Repeat("#", cgrpc.DefaultMetadataLimit)),
			data:       []byte(strings.Repeat("#", cgrpc.DefaultDataLimit)),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			expected := uuid.New()

			m := newServicesMock()
			m.Secrets.(*service.SecretsServiceMock).On(
				"Create",
				mock.Anything,
				mock.AnythingOfType("uuid.UUID"),
				tc.secretName,
				proto.DataKind_BINARY,
				tc.metadata,
				tc.data,
			).
				Return(expected, nil)

			conn := createTestServerWithFakeAuth(t, m)

			req := &proto.CreateSecretRequest{
				Name:     tc.secretName,
				Kind:     proto.DataKind_BINARY,
				Metadata: tc.metadata,
				Data:     tc.data,
			}

			client := proto.NewSecretsClient(conn)
			resp, err := client.Create(context.Background(), req)

			require.NoError(t, err)
			require.Equal(t, expected.String(), resp.GetId())
			m.Secrets.(*service.SecretsServiceMock).AssertExpectations(t)
		})
	}
}

func TestCreateSecretWithBadRequest(t *testing.T) {
	tt := []struct {
		name       string
		secretName string
		metadata   []byte
		data       []byte
	}{
		{
			name:       "Create secret fails if secret name is empty",
			secretName: "",
			metadata:   []byte(gophtest.Metadata),
			data:       []byte(gophtest.TextData),
		},
		{
			name:       "Create secret fails if secret name is too long",
			secretName: strings.Repeat("#", cgrpc.DefaultMaxSecretNameLength+1),
			metadata:   []byte(gophtest.Metadata),
			data:       []byte(gophtest.TextData),
		},
		{
			name:       "Create secret fails if metadata is too long",
			secretName: gophtest.Username,
			metadata:   []byte(strings.Repeat("#", cgrpc.DefaultMetadataLimit+1)),
			data:       make([]byte, 0),
		},
		{
			name:       "Create secret fails if data is empty",
			secretName: gophtest.Username,
			metadata:   []byte(gophtest.Metadata),
			data:       make([]byte, 0),
		},
		{
			name:       "Create secret fails if data is too long",
			secretName: gophtest.Username,
			metadata:   []byte(gophtest.Metadata),
			data:       []byte(strings.Repeat("#", cgrpc.DefaultDataLimit+1)),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			conn := createTestServerWithFakeAuth(t, newServicesMock())

			req := &proto.CreateSecretRequest{
				Name:     tc.secretName,
				Kind:     proto.DataKind_BINARY,
				Metadata: tc.metadata,
				Data:     tc.data,
			}

			client := proto.NewSecretsClient(conn)
			_, err := client.Create(context.Background(), req)

			requireEqualCode(t, codes.InvalidArgument, err)
		})
	}
}

func TestCreateSecretFailsIfNoUserInfo(t *testing.T) {
	conn := createTestServer(t, newServicesMock())

	client := proto.NewSecretsClient(conn)
	_, err := client.Create(context.Background(), &proto.CreateSecretRequest{})

	requireEqualCode(t, codes.Unauthenticated, err)
}

func TestCreateServerOnServiceFailure(t *testing.T) {
	tt := []struct {
		name     string
		err      error
		expected codes.Code
	}{
		{
			name:     "Create secret fails if secret already exists",
			err:      entity.ErrSecretExists,
			expected: codes.AlreadyExists,
		},
		{
			name:     "Create secret fails if use case fails unexpectedly",
			err:      gophtest.ErrUnexpected,
			expected: codes.Internal,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := newServicesMock()
			m.Secrets.(*service.SecretsServiceMock).On(
				"Create",
				mock.Anything,
				mock.AnythingOfType("uuid.UUID"),
				gophtest.SecretName,
				proto.DataKind_BINARY,
				[]byte(gophtest.Metadata),
				[]byte(gophtest.TextData),
			).
				Return(uuid.UUID{}, tc.err)

			conn := createTestServerWithFakeAuth(t, m)

			req := &proto.CreateSecretRequest{
				Name:     gophtest.SecretName,
				Kind:     proto.DataKind_BINARY,
				Metadata: []byte(gophtest.Metadata),
				Data:     []byte(gophtest.TextData),
			}

			client := proto.NewSecretsClient(conn)
			_, err := client.Create(context.Background(), req)

			requireEqualCode(t, tc.expected, err)
			m.Secrets.(*service.SecretsServiceMock).AssertExpectations(t)
		})
	}
}

func TestListSecrets(t *testing.T) {
	tt := []struct {
		name    string
		secrets []entity.Secret
	}{
		{
			name: "List secrets of a user",
			secrets: []entity.Secret{
				{
					ID:   gophtest.CreateUUID(t, "7728154c-9400-4f1b-a2a3-01deb83ece05"),
					Name: gophtest.SecretName,
					Kind: proto.DataKind_BINARY,
				},
				{
					ID:       gophtest.CreateUUID(t, "df566e25-43a5-4c34-9123-3931fb809b45"),
					Name:     gophtest.SecretName + "ex",
					Kind:     proto.DataKind_TEXT,
					Metadata: []byte(gophtest.Metadata),
				},
			},
		},
		{
			name:    "List secrets when user has no secrets",
			secrets: []entity.Secret{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rv, err := doListSecrets(t, tc.secrets, nil)

			require.NoError(t, err)
			snaps.MatchSnapshot(t, rv.GetSecrets())
		})
	}
}

func TestListSecretsFailsIfNoUserInfo(t *testing.T) {
	conn := createTestServer(t, newServicesMock())

	client := proto.NewSecretsClient(conn)
	_, err := client.List(context.Background(), &proto.ListSecretsRequest{})

	requireEqualCode(t, codes.Unauthenticated, err)
}

func TestListSecretsFailsOnServiceFailure(t *testing.T) {
	_, err := doListSecrets(t, nil, gophtest.ErrUnexpected)

	requireEqualCode(t, codes.Internal, err)
}

func TestGetSecret(t *testing.T) {
	tt := []struct {
		name   string
		secret *entity.Secret
	}{
		{
			name: "Get secret",
			secret: &entity.Secret{
				ID:       gophtest.CreateUUID(t, "df566e25-43a5-4c34-9123-3931fb809b45"),
				Name:     gophtest.SecretName,
				Kind:     proto.DataKind_TEXT,
				Metadata: []byte(gophtest.Metadata),
				Data:     []byte(gophtest.TextData),
			},
		},
		{
			name: "Get secret without metadata",
			secret: &entity.Secret{
				ID:       gophtest.CreateUUID(t, "df566e25-43a5-4c34-9123-3931fb809b45"),
				Name:     gophtest.SecretName,
				Kind:     proto.DataKind_TEXT,
				Metadata: []byte(gophtest.Metadata),
				Data:     []byte(gophtest.TextData),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := doGetSecret(t, tc.secret, nil)

			require.NoError(t, err)
			snaps.MatchSnapshot(t, resp.GetSecret())
			require.Equal(t, tc.secret.Data, resp.GetData())
		})
	}
}

func TestGetSecretOnBadRequest(t *testing.T) {
	conn := createTestServerWithFakeAuth(t, newServicesMock())

	req := &proto.GetSecretRequest{Id: "xxx"}

	client := proto.NewSecretsClient(conn)
	_, err := client.Get(context.Background(), req)

	requireEqualCode(t, codes.InvalidArgument, err)
}

func TestGetSecretFailsIfNoUserInfo(t *testing.T) {
	conn := createTestServer(t, newServicesMock())

	client := proto.NewSecretsClient(conn)
	_, err := client.Get(context.Background(), &proto.GetSecretRequest{})

	requireEqualCode(t, codes.Unauthenticated, err)
}

func TestGetSecretOnServiceFailure(t *testing.T) {
	tt := []struct {
		name     string
		ucErr    error
		expected codes.Code
	}{
		{
			name:     "Get secret fails if secret not found",
			ucErr:    entity.ErrSecretNotFound,
			expected: codes.NotFound,
		},
		{
			name:     "Get secret fails on expected error",
			ucErr:    gophtest.ErrUnexpected,
			expected: codes.Internal,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := doGetSecret(t, nil, tc.ucErr)

			requireEqualCode(t, tc.expected, err)
		})
	}
}

func TestUpdateSecret(t *testing.T) {
	tt := []struct {
		name    string
		req     *proto.UpdateSecretRequest
		changed []string
	}{
		{
			name: "Update all fields of a secret",
			req: &proto.UpdateSecretRequest{
				Name:     gophtest.SecretName,
				Metadata: []byte(gophtest.Metadata),
				Data:     []byte(gophtest.TextData),
			},
			changed: []string{"data", "metadata", "name"},
		},
		{
			name: "Update secret's name",
			req: &proto.UpdateSecretRequest{
				Name: gophtest.SecretName,
			},
			changed: []string{"name"},
		},
		{
			name: "Update secret's metadata",
			req: &proto.UpdateSecretRequest{
				Metadata: []byte(gophtest.Metadata),
			},
			changed: []string{"metadata"},
		},
		{
			name: "Reset secret's metadata",
			req: &proto.UpdateSecretRequest{
				Metadata: []byte(nil),
			},
			changed: []string{"metadata"},
		},
		{
			name: "Update secret's data",
			req: &proto.UpdateSecretRequest{
				Data: []byte(gophtest.TextData),
			},
			changed: []string{"data"},
		},
		{
			name: "Update secret with maximum fields limits",
			req: &proto.UpdateSecretRequest{
				Name:     strings.Repeat("#", cgrpc.DefaultMaxSecretNameLength),
				Metadata: []byte(strings.Repeat("#", cgrpc.DefaultMetadataLimit)),
				Data:     []byte(strings.Repeat("#", cgrpc.DefaultDataLimit)),
			},
			changed: []string{"data", "metadata", "name"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			id := uuid.New()

			mask, err := fieldmaskpb.New(tc.req, tc.changed...)
			require.NoError(t, err)

			tc.req.Id = id.String()
			tc.req.UpdateMask = mask

			m := newServicesMock()
			m.Secrets.(*service.SecretsServiceMock).On(
				"Update",
				mock.Anything,
				mock.AnythingOfType("uuid.UUID"),
				id,
				tc.changed,
				tc.req.Name,
				tc.req.Metadata,
				tc.req.Data,
			).
				Return(nil)

			conn := createTestServerWithFakeAuth(t, m)

			client := proto.NewSecretsClient(conn)
			_, err = client.Update(context.Background(), tc.req)

			m.Secrets.(*service.SecretsServiceMock).AssertExpectations(t)

			require.NoError(t, err)
		})
	}
}

func TestUpdateSecretOnBadRequest(t *testing.T) {
	tt := []struct {
		name    string
		req     *proto.UpdateSecretRequest
		changed []string
	}{
		{
			name: "Update fails if no mask specified",
			req: &proto.UpdateSecretRequest{
				Id:   uuid.New().String(),
				Name: gophtest.SecretName,
			},
			changed: nil,
		},
		{
			name: "Update fails if bad secret id provided",
			req: &proto.UpdateSecretRequest{
				Id:   "xxx",
				Name: gophtest.SecretName,
			},
			changed: []string{"name"},
		},
		{
			name: "Update fails if empty name provided",
			req: &proto.UpdateSecretRequest{
				Id:   uuid.New().String(),
				Name: "",
			},
			changed: []string{"name"},
		},
		{
			name: "Update fails if too long name provided",
			req: &proto.UpdateSecretRequest{
				Id:   uuid.New().String(),
				Name: strings.Repeat("#", cgrpc.DefaultMaxSecretNameLength+1),
			},
			changed: []string{"name"},
		},
		{
			name: "Update fails if too long metadata provided",
			req: &proto.UpdateSecretRequest{
				Id:       uuid.New().String(),
				Metadata: []byte(strings.Repeat("#", cgrpc.DefaultMetadataLimit+1)),
			},
			changed: []string{"metadata"},
		},
		{
			name: "Update fails if empty data provided",
			req: &proto.UpdateSecretRequest{
				Id: uuid.New().String(),
			},
			changed: []string{"data"},
		},
		{
			name: "Update fails if too long data provided",
			req: &proto.UpdateSecretRequest{
				Id:   uuid.New().String(),
				Data: []byte(strings.Repeat("#", cgrpc.DefaultDataLimit+1)),
			},
			changed: []string{"data"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			conn := createTestServerWithFakeAuth(t, newServicesMock())

			mask, err := fieldmaskpb.New(tc.req, tc.changed...)
			require.NoError(t, err)

			tc.req.UpdateMask = mask

			client := proto.NewSecretsClient(conn)
			_, err = client.Update(context.Background(), tc.req)

			requireEqualCode(t, codes.InvalidArgument, err)
		})
	}
}

func TestUpdateSecretFailsIfNoUserInfo(t *testing.T) {
	conn := createTestServer(t, newServicesMock())

	client := proto.NewSecretsClient(conn)
	_, err := client.Update(context.Background(), &proto.UpdateSecretRequest{})

	requireEqualCode(t, codes.Unauthenticated, err)
}

func TestUpdateSecretOnServiceFailure(t *testing.T) {
	tt := []struct {
		name     string
		ucErr    error
		expected codes.Code
	}{
		{
			name:     "Update secret fails if secret not found",
			ucErr:    entity.ErrSecretNotFound,
			expected: codes.NotFound,
		},
		{
			name:     "Update secret fails if secret with provided name already exists",
			ucErr:    entity.ErrSecretNameConflict,
			expected: codes.AlreadyExists,
		},
		{
			name:     "Update secret fails on expected error",
			ucErr:    gophtest.ErrUnexpected,
			expected: codes.Internal,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			id := uuid.New()

			m := newServicesMock()
			m.Secrets.(*service.SecretsServiceMock).On(
				"Update",
				mock.Anything,
				mock.AnythingOfType("uuid.UUID"),
				id,
				[]string{"name"},
				gophtest.SecretName,
				[]byte(nil),
				[]byte(nil),
			).
				Return(tc.ucErr)

			conn := createTestServerWithFakeAuth(t, m)
			req := &proto.UpdateSecretRequest{
				Id:   id.String(),
				Name: gophtest.SecretName,
			}

			mask, err := fieldmaskpb.New(req, "name")
			require.NoError(t, err)

			req.UpdateMask = mask

			client := proto.NewSecretsClient(conn)
			_, err = client.Update(context.Background(), req)

			m.Secrets.(*service.SecretsServiceMock).AssertExpectations(t)
			requireEqualCode(t, tc.expected, err)
		})
	}
}

func TestDeleteSecret(t *testing.T) {
	_, err := doDeleteSecret(t, nil)

	require.NoError(t, err)
}

func TestDeleteSecretOnBadRequest(t *testing.T) {
	conn := createTestServerWithFakeAuth(t, newServicesMock())

	req := &proto.DeleteSecretRequest{Id: "xxx"}

	client := proto.NewSecretsClient(conn)
	_, err := client.Delete(context.Background(), req)

	requireEqualCode(t, codes.InvalidArgument, err)
}

func TestDeleteSecretFailsIfNoUserInfo(t *testing.T) {
	conn := createTestServer(t, newServicesMock())

	client := proto.NewSecretsClient(conn)
	_, err := client.Delete(context.Background(), &proto.DeleteSecretRequest{})

	requireEqualCode(t, codes.Unauthenticated, err)
}

func TestDeleteSecretOnServiceFailure(t *testing.T) {
	tt := []struct {
		name     string
		ucErr    error
		expected codes.Code
	}{
		{
			name:     "Delete secret fails if secret not found",
			ucErr:    entity.ErrSecretNotFound,
			expected: codes.NotFound,
		},
		{
			name:     "Delete secret fails on expected error",
			ucErr:    gophtest.ErrUnexpected,
			expected: codes.Internal,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := doDeleteSecret(t, tc.ucErr)

			requireEqualCode(t, tc.expected, err)
		})
	}
}
