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
)

func TestRegisterUser(t *testing.T) {
	expected := uuid.New()

	rows := pgxmock.NewRows([]string{"id"}).
		AddRow(expected.String())

	m := newPoolMock(t)
	m.ExpectBeginTx(postgres.DefaultTxOptions)
	m.ExpectQuery("INSERT INTO users").
		WithArgs(gophtest.Username, gophtest.SecurityKey).
		WillReturnRows(rows)
	m.ExpectCommit()

	sat := newTestRepos(t, m).Users
	id, err := sat.Register(context.Background(), gophtest.Username, gophtest.SecurityKey)

	require.NoError(t, err)
	require.Equal(t, expected, id)
	require.NoError(t, m.ExpectationsWereMet())
}

func TestRegisterUserOnDBFailure(t *testing.T) {
	tt := []struct {
		name     string
		err      error
		expected error
	}{
		{
			name:     "Register user fails if user exists",
			err:      errUniqueViolation,
			expected: entity.ErrUserExists,
		},
		{
			name:     "Register user fails on unexpected error",
			err:      gophtest.ErrUnexpected,
			expected: gophtest.ErrUnexpected,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := newPoolMock(t)
			m.ExpectBeginTx(postgres.DefaultTxOptions)
			m.ExpectQuery("INSERT").
				WithArgs(gophtest.Username, gophtest.SecurityKey).
				WillReturnError(tc.err)
			m.ExpectRollback()

			sat := newTestRepos(t, m).Users
			_, err := sat.Register(context.Background(), gophtest.Username, gophtest.SecurityKey)

			require.ErrorIs(t, err, tc.expected)
			require.NoError(t, m.ExpectationsWereMet())
		})
	}
}

func TestVerifyUser(t *testing.T) {
	expected := entity.User{
		ID:       uuid.New(),
		Username: gophtest.Username,
	}

	rows := pgxmock.NewRows([]string{"user_id", "username"}).
		AddRow(expected.ID.String(), expected.Username)

	m := newPoolMock(t)
	m.ExpectQuery("SELECT user_id, username FROM users").
		WithArgs(gophtest.Username, gophtest.SecurityKey).
		WillReturnRows(rows)

	sat := newTestRepos(t, m).Users
	rv, err := sat.Verify(context.Background(), gophtest.Username, gophtest.SecurityKey)

	require.NoError(t, err)
	require.Equal(t, expected, rv)
	require.NoError(t, m.ExpectationsWereMet())
}

func TestVerifyFailsOnBadCredentials(t *testing.T) {
	rows := pgxmock.NewRows([]string{"user_id", "username"})

	m := newPoolMock(t)
	m.ExpectQuery("SELECT").
		WithArgs(gophtest.Username, gophtest.SecurityKey).
		WillReturnRows(rows)

	sat := newTestRepos(t, m).Users
	_, err := sat.Verify(context.Background(), gophtest.Username, gophtest.SecurityKey)

	require.ErrorIs(t, err, entity.ErrInvalidCredentials)
	require.NoError(t, m.ExpectationsWereMet())
}

func TestVerifyFailsOnDBFailure(t *testing.T) {
	m := newPoolMock(t)
	m.ExpectQuery("SELECT").
		WithArgs(gophtest.Username, gophtest.SecurityKey).
		WillReturnError(gophtest.ErrUnexpected)

	sat := newTestRepos(t, m).Users
	_, err := sat.Verify(context.Background(), gophtest.Username, gophtest.SecurityKey)

	require.ErrorIs(t, err, gophtest.ErrUnexpected)
	require.NoError(t, m.ExpectationsWereMet())
}
