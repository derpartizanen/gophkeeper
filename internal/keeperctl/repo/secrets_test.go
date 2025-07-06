package repo_test

import (
	"context"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/derpartizanen/gophkeeper/internal/keeperctl/repo"
	"github.com/derpartizanen/gophkeeper/internal/libraries/gophtest"
	"github.com/derpartizanen/gophkeeper/proto"
)

func doCreateSecret(
	t *testing.T,
	mockRV *proto.CreateSecretResponse,
	mockErr error,
) (uuid.UUID, error) {
	t.Helper()

	req := &proto.CreateSecretRequest{
		Name:     gophtest.SecretName,
		Metadata: []byte(gophtest.Metadata),
		Kind:     proto.DataKind_TEXT,
		Data:     []byte(gophtest.TextData),
	}

	m := &proto.SecretsClientMock{}
	m.On(
		"Create",
		mock.Anything,
		req,
		mock.Anything,
	).
		Return(mockRV, mockErr)

	sat := repo.NewSecretsRepo(m)
	rv, err := sat.Push(
		context.Background(),
		gophtest.AccessToken,
		gophtest.SecretName,
		proto.DataKind_TEXT,
		[]byte(gophtest.Metadata),
		[]byte(gophtest.TextData),
	)

	m.AssertExpectations(t)

	return rv, err
}

func doListSecrets(
	t *testing.T,
	mockRV *proto.ListSecretsResponse,
	mockErr error,
) ([]*proto.Secret, error) {
	t.Helper()

	req := &proto.ListSecretsRequest{}

	m := &proto.SecretsClientMock{}
	m.On(
		"List",
		mock.Anything,
		req,
		mock.Anything,
	).
		Return(mockRV, mockErr)

	sat := repo.NewSecretsRepo(m)
	rv, err := sat.List(context.Background(), gophtest.AccessToken)

	m.AssertExpectations(t)

	return rv, err
}

func doGetSecret(
	t *testing.T,
	mockRV *proto.GetSecretResponse,
	mockErr error,
) (*proto.Secret, []byte, error) {
	t.Helper()

	id := uuid.New()
	req := &proto.GetSecretRequest{Id: id.String()}

	m := &proto.SecretsClientMock{}
	m.On(
		"Get",
		mock.Anything,
		req,
		mock.Anything,
	).
		Return(mockRV, mockErr)

	sat := repo.NewSecretsRepo(m)
	secret, data, err := sat.Get(context.Background(), gophtest.AccessToken, id)

	m.AssertExpectations(t)

	return secret, data, err
}

func doUpdateSecret(
	t *testing.T,
	name string,
	description []byte,
	noDescription bool,
	data []byte,
	changed []string,
	clientErr error,
) error {
	t.Helper()

	id := uuid.New()
	req := &proto.UpdateSecretRequest{
		Id:       id.String(),
		Name:     name,
		Metadata: description,
		Data:     data,
	}

	mask, err := fieldmaskpb.New(req, changed...)
	require.NoError(t, err)

	req.UpdateMask = mask

	m := &proto.SecretsClientMock{}
	m.On(
		"Update",
		mock.Anything,
		req,
		mock.Anything,
	).
		Return(&proto.UpdateSecretResponse{}, clientErr)

	sat := repo.NewSecretsRepo(m)
	err = sat.Update(
		context.Background(),
		gophtest.AccessToken,
		id,
		name,
		description,
		noDescription,
		data,
	)

	m.AssertExpectations(t)

	return err
}

func doDeleteSecret(t *testing.T, mockErr error) error {
	t.Helper()

	id := uuid.New()
	req := &proto.DeleteSecretRequest{Id: id.String()}

	m := &proto.SecretsClientMock{}
	m.On(
		"Delete",
		mock.Anything,
		req,
		mock.Anything,
	).
		Return(&proto.DeleteSecretResponse{}, mockErr)

	sat := repo.NewSecretsRepo(m)
	err := sat.Delete(context.Background(), gophtest.AccessToken, id)

	m.AssertExpectations(t)

	return err
}

func TestCreateSecret(t *testing.T) {
	expected := uuid.New()
	resp := &proto.CreateSecretResponse{
		Id: expected.String(),
	}

	id, err := doCreateSecret(t, resp, nil)

	require.NoError(t, err)
	require.Equal(t, expected, id)
}

func TestCreateSecretOnClientFailure(t *testing.T) {
	_, err := doCreateSecret(t, nil, gophtest.ErrUnexpected)

	require.Error(t, err)
}

func TestListSecrets(t *testing.T) {
	tt := []struct {
		name    string
		secrets []*proto.Secret
	}{
		{
			name: "List secrets of a user",
			secrets: []*proto.Secret{
				{
					Id:       gophtest.CreateUUID(t, "df566e25-43a5-4c34-9123-3931fb809b45").String(),
					Name:     gophtest.SecretName,
					Kind:     proto.DataKind_TEXT,
					Metadata: []byte(gophtest.Metadata),
				},
			},
		},
		{
			name:    "List secrets of a user who has no secrets",
			secrets: []*proto.Secret{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			resp := &proto.ListSecretsResponse{
				Secrets: tc.secrets,
			}

			rv, err := doListSecrets(t, resp, nil)

			require.NoError(t, err)
			snaps.MatchSnapshot(t, rv)
		})
	}
}

func TestListSecretsOnClientFailure(t *testing.T) {
	_, err := doListSecrets(t, nil, gophtest.ErrUnexpected)

	require.Error(t, err)
}

func TestGetSecret(t *testing.T) {
	expSecret := &proto.Secret{
		Id:       uuid.New().String(),
		Name:     gophtest.SecretName,
		Kind:     proto.DataKind_TEXT,
		Metadata: []byte(gophtest.Metadata),
	}
	expData := []byte(gophtest.TextData)

	mockRV := &proto.GetSecretResponse{
		Secret: expSecret,
		Data:   expData,
	}

	secret, data, err := doGetSecret(t, mockRV, nil)

	require.NoError(t, err)
	require.Equal(t, expSecret, secret)
	require.Equal(t, expData, data)
}

func TestGetSecretOnClientFailure(t *testing.T) {
	_, _, err := doGetSecret(t, nil, gophtest.ErrUnexpected)

	require.Error(t, err)
}

func TestUpdateSecret(t *testing.T) {
	tt := []struct {
		name          string
		secretName    string
		description   []byte
		noDescription bool
		data          []byte
		changed       []string
	}{
		{
			name:        "Update all fields of a secret",
			secretName:  gophtest.SecretName,
			description: []byte(gophtest.Metadata),
			data:        []byte(gophtest.TextData),
			changed:     []string{"name", "metadata", "data"},
		},
		{
			name:       "Update secret's name",
			secretName: gophtest.SecretName,
			changed:    []string{"name"},
		},
		{
			name:        "Update secret's description",
			description: []byte(gophtest.Metadata),
			changed:     []string{"metadata"},
		},
		{
			name:          "Reset secret's description",
			noDescription: true,
			changed:       []string{"metadata"},
		},
		{
			name:    "Update secret's data",
			data:    []byte(gophtest.TextData),
			changed: []string{"data"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := doUpdateSecret(
				t,
				tc.secretName,
				tc.description,
				tc.noDescription,
				tc.data,
				tc.changed,
				nil,
			)

			require.NoError(t, err)
		})
	}
}

func TestUpdateSecretOnClientFailure(t *testing.T) {
	err := doUpdateSecret(t, "", nil, false, nil, nil, gophtest.ErrUnexpected)

	require.Error(t, err)
}

func TestDeleteSecret(t *testing.T) {
	err := doDeleteSecret(t, nil)

	require.NoError(t, err)
}

func TestDeleteSecretOnClientFailure(t *testing.T) {
	err := doDeleteSecret(t, gophtest.ErrUnexpected)

	require.Error(t, err)
}
