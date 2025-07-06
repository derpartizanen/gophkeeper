package repo_test

import (
	"testing"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/require"

	"github.com/derpartizanen/gophkeeper/internal/keeperd/postgres"
	"github.com/derpartizanen/gophkeeper/internal/keeperd/repo"
)

var errUniqueViolation = error(&pgconn.PgError{Code: pgerrcode.UniqueViolation})

func newPoolMock(t *testing.T) pgxmock.PgxPoolIface {
	t.Helper()

	m, err := pgxmock.NewPool()
	require.NoError(t, err)

	t.Cleanup(m.Close)

	return m
}

func newTestRepos(t *testing.T, m pgxmock.PgxPoolIface) *repo.Repositories {
	t.Helper()

	pg := &postgres.Postgres{
		Pool: m,
	}

	return repo.New(pg)
}
