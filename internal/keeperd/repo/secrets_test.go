package repo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/require"

	"github.com/derpartizanen/gophkeeper/internal/keeperd/entity"
	"github.com/derpartizanen/gophkeeper/internal/keeperd/postgres"
	"github.com/derpartizanen/gophkeeper/internal/libraries/gophtest"
	"github.com/derpartizanen/gophkeeper/proto"
)

func doGetSecret(
	t *testing.T,
	owner, id uuid.UUID,
	m pgxmock.PgxPoolIface,
) (*entity.Secret, error) {
	t.Helper()

	sat := newTestRepos(t, m).Secrets
	secret, err := sat.Get(context.Background(), owner, id)

	require.NoError(t, m.ExpectationsWereMet())

	return secret, err
}

func doUpdateSecret(
	t *testing.T,
	owner, id uuid.UUID,
	changed []string,
	name string,
	metadata, data []byte,
	m pgxmock.PgxPoolIface,
) error {
	t.Helper()

	sat := newTestRepos(t, m).Secrets
	err := sat.Update(
		context.Background(),
		owner,
		id,
		changed,
		name,
		metadata,
		data,
	)

	require.NoError(t, m.ExpectationsWereMet())

	return err
}

func doDeleteSecret(t *testing.T, owner, id uuid.UUID, m pgxmock.PgxPoolIface) error {
	t.Helper()

	sat := newTestRepos(t, m).Secrets
	err := sat.Delete(context.Background(), owner, id)

	require.NoError(t, m.ExpectationsWereMet())

	return err
}

func TestCreateSecret(t *testing.T) {
	owner := uuid.New()
	expected := uuid.New()

	rows := pgxmock.NewRows([]string{"id"}).
		AddRow(expected.String())

	m := newPoolMock(t)
	m.ExpectBeginTx(postgres.DefaultTxOptions)
	m.ExpectQuery("INSERT INTO secrets").
		WithArgs(
			owner,
			gophtest.SecretName,
			proto.DataKind_TEXT,
			[]byte(gophtest.Metadata),
			[]byte(gophtest.TextData),
		).
		WillReturnRows(rows)
	m.ExpectCommit()

	sat := newTestRepos(t, m).Secrets
	id, err := sat.Create(
		context.Background(),
		owner,
		gophtest.SecretName,
		proto.DataKind_TEXT,
		[]byte(gophtest.Metadata),
		[]byte(gophtest.TextData),
	)

	require.NoError(t, err)
	require.Equal(t, expected, id)
	require.NoError(t, m.ExpectationsWereMet())
}

func TestCreateSecretOnDBFailure(t *testing.T) {
	tt := []struct {
		name     string
		err      error
		expected error
	}{
		{
			name:     "Create secret fails if secret exists",
			err:      errUniqueViolation,
			expected: entity.ErrSecretExists,
		},
		{
			name:     "Create secret fails on unexpected error",
			err:      gophtest.ErrUnexpected,
			expected: gophtest.ErrUnexpected,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			owner := uuid.New()

			m := newPoolMock(t)
			m.ExpectBeginTx(postgres.DefaultTxOptions)
			m.ExpectQuery("INSERT").
				WithArgs(
					owner,
					gophtest.SecretName,
					proto.DataKind_TEXT,
					[]byte(gophtest.Metadata),
					[]byte(gophtest.TextData),
				).
				WillReturnError(tc.err)
			m.ExpectRollback()

			sat := newTestRepos(t, m).Secrets
			_, err := sat.Create(
				context.Background(),
				owner,
				gophtest.SecretName,
				proto.DataKind_TEXT,
				[]byte(gophtest.Metadata),
				[]byte(gophtest.TextData),
			)

			require.ErrorIs(t, err, tc.expected)
			require.NoError(t, m.ExpectationsWereMet())
		})
	}
}

func TestListSecrets(t *testing.T) {
	tt := []struct {
		name string
		rows [][]any
	}{
		{
			name: "List secrets of a user",
			rows: [][]any{
				{uuid.New().String(), gophtest.SecretName, proto.DataKind_TEXT, []byte("xxx")},
				{uuid.New().String(), gophtest.SecretName + "ex", proto.DataKind_BINARY, []byte{}},
			},
		},
		{
			name: "List secrets returns empty list",
			rows: [][]any{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			owner := uuid.New()
			rows := pgxmock.NewRows([]string{"secret_id", "name", "kind", "metadata"})

			for _, row := range tc.rows {
				rows.AddRow(row...)
			}

			m := newPoolMock(t)
			m.ExpectQuery("SELECT secret_id, name, kind, metadata FROM secrets").
				WithArgs(owner).
				WillReturnRows(rows)

			sat := newTestRepos(t, m).Secrets
			secrets, err := sat.List(context.Background(), owner)

			require.NoError(t, err)
			require.Len(t, secrets, len(tc.rows))
			require.NoError(t, m.ExpectationsWereMet())
		})
	}
}

func TestListSecretsOnDBFailure(t *testing.T) {
	owner := uuid.New()

	m := newPoolMock(t)
	m.ExpectQuery("SELECT").
		WithArgs(owner).
		WillReturnError(gophtest.ErrUnexpected)

	sat := newTestRepos(t, m).Secrets
	_, err := sat.List(context.Background(), owner)

	require.Error(t, err)
	require.NoError(t, m.ExpectationsWereMet())
}

func TestGetSecret(t *testing.T) {
	owner := uuid.New()

	expected := &entity.Secret{
		ID:       uuid.New(),
		Name:     gophtest.SecretName,
		Kind:     proto.DataKind_TEXT,
		Metadata: []byte(gophtest.Metadata),
		Data:     []byte(gophtest.TextData),
	}

	rows := pgxmock.NewRows([]string{"secret_id", "name", "kind", "metadata", "data"}).
		AddRow(expected.ID.String(), expected.Name, expected.Kind, expected.Metadata, expected.Data)

	m := newPoolMock(t)
	m.ExpectQuery("SELECT secret_id, name, kind, metadata, data FROM secrets").
		WithArgs(expected.ID, owner).
		WillReturnRows(rows)

	secret, err := doGetSecret(t, owner, expected.ID, m)

	require.NoError(t, err)
	require.Equal(t, expected, secret)
}

func TestGetUnexistingSecret(t *testing.T) {
	rows := pgxmock.NewRows([]string{"secret_id", "name", "kind", "metadata", "data"})

	owner := uuid.New()
	id := uuid.New()

	m := newPoolMock(t)
	m.ExpectQuery("SELECT").
		WithArgs(id, owner).
		WillReturnRows(rows)

	_, err := doGetSecret(t, owner, id, m)

	require.ErrorIs(t, err, entity.ErrSecretNotFound)
}

func TestGetSecretOnDBFailure(t *testing.T) {
	owner := uuid.New()
	id := uuid.New()

	m := newPoolMock(t)
	m.ExpectQuery("SELECT").
		WithArgs(id, owner).
		WillReturnError(gophtest.ErrUnexpected)

	_, err := doGetSecret(t, owner, id, m)

	require.ErrorIs(t, err, gophtest.ErrUnexpected)
}

func TestUpdateSecret(t *testing.T) {
	owner := uuid.New()
	id := uuid.New()

	type expected struct {
		query string
		args  []any
	}

	tt := []struct {
		name       string
		secretName string
		changed    []string
		metadata   []byte
		data       []byte
		expected   expected
	}{
		{
			name:       "Update all fields",
			changed:    []string{"name", "metadata", "data"},
			secretName: gophtest.SecretName,
			metadata:   []byte(gophtest.Metadata),
			data:       []byte(gophtest.TextData),
			expected: expected{
				query: "UPDATE secrets SET name = \\$1, metadata = \\$2, data = \\$3",
				args: []any{
					gophtest.SecretName,
					[]byte(gophtest.Metadata),
					[]byte(gophtest.TextData),
					id,
					owner,
				},
			},
		},
		{
			name:       "Update name",
			changed:    []string{"name"},
			secretName: gophtest.SecretName,
			expected: expected{
				query: "UPDATE secrets SET name = \\$1",
				args:  []any{gophtest.SecretName, id, owner},
			},
		},
		{
			name:     "Update metadata",
			changed:  []string{"metadata"},
			metadata: []byte(gophtest.Metadata),
			expected: expected{
				query: "UPDATE secrets SET metadata = \\$1",
				args:  []any{[]byte(gophtest.Metadata), id, owner},
			},
		},
		{
			name:    "Update data",
			changed: []string{"data"},
			data:    []byte(gophtest.TextData),
			expected: expected{
				query: "UPDATE secrets SET data = \\$1",
				args:  []any{[]byte(gophtest.TextData), id, owner},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := newPoolMock(t)
			m.ExpectBeginTx(postgres.DefaultTxOptions)
			m.ExpectExec(tc.expected.query).
				WithArgs(tc.expected.args...).
				WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			m.ExpectCommit()

			err := doUpdateSecret(
				t,
				owner,
				id,
				tc.changed,
				tc.secretName,
				tc.metadata,
				tc.data,
				m,
			)

			require.NoError(t, err)
		})
	}
}

func TestUpdateUnexistingSecret(t *testing.T) {
	owner := uuid.New()
	id := uuid.New()

	m := newPoolMock(t)
	m.ExpectBeginTx(postgres.DefaultTxOptions)
	m.ExpectExec("UPDATE").
		WithArgs(gophtest.SecretName, id, owner).
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))
	m.ExpectRollback()

	err := doUpdateSecret(
		t,
		owner,
		id,
		[]string{"name"},
		gophtest.SecretName,
		nil,
		nil,
		m,
	)

	require.ErrorIs(t, err, entity.ErrSecretNotFound)
}

func TestUpdateSecretOnDBFailure(t *testing.T) {
	tt := []struct {
		name     string
		err      error
		expected error
	}{
		{
			name:     "Update secret fails if new name conflicts with other secret",
			err:      errUniqueViolation,
			expected: entity.ErrSecretNameConflict,
		},
		{
			name:     "Update secret fails on unexpected error",
			err:      gophtest.ErrUnexpected,
			expected: gophtest.ErrUnexpected,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			owner := uuid.New()
			id := uuid.New()

			m := newPoolMock(t)
			m.ExpectBeginTx(postgres.DefaultTxOptions)
			m.ExpectExec("UPDATE").
				WithArgs(gophtest.SecretName, id, owner).
				WillReturnError(tc.err)
			m.ExpectRollback()

			err := doUpdateSecret(
				t,
				owner,
				id,
				[]string{"name"},
				gophtest.SecretName,
				nil,
				nil,
				m,
			)

			require.ErrorIs(t, err, tc.expected)
			require.NoError(t, m.ExpectationsWereMet())
		})
	}
}

func TestDeleteSecret(t *testing.T) {
	owner := uuid.New()
	id := uuid.New()

	m := newPoolMock(t)
	m.ExpectBeginTx(postgres.DefaultTxOptions)
	m.ExpectExec("DELETE FROM secrets").
		WithArgs(id, owner).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	m.ExpectCommit()

	err := doDeleteSecret(t, owner, id, m)

	require.NoError(t, err)
}

func TestDeleteUnexistingSecret(t *testing.T) {
	owner := uuid.New()
	id := uuid.New()

	m := newPoolMock(t)
	m.ExpectBeginTx(postgres.DefaultTxOptions)
	m.ExpectExec("DELETE").
		WithArgs(id, owner).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))
	m.ExpectRollback()

	err := doDeleteSecret(t, owner, id, m)

	require.ErrorIs(t, err, entity.ErrSecretNotFound)
}

func TestDeleteSecretOnDBFailure(t *testing.T) {
	owner := uuid.New()
	id := uuid.New()

	m := newPoolMock(t)
	m.ExpectBeginTx(postgres.DefaultTxOptions)
	m.ExpectExec("DELETE").
		WithArgs(id, owner).
		WillReturnError(gophtest.ErrUnexpected)
	m.ExpectRollback()

	err := doDeleteSecret(t, owner, id, m)

	require.ErrorIs(t, err, gophtest.ErrUnexpected)
}
