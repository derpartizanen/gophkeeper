package postgres

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// IsEmptyResponse returns true if error means that response contains no data.
func IsEmptyResponse(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

// IsEntityExists returns true if the provided error related to
// cases when entity already exists in the database.
func IsEntityExists(err error) bool {
	var pgErr *pgconn.PgError

	return errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation
}
