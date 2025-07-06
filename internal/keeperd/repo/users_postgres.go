package repo

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/derpartizanen/gophkeeper/internal/keeperd/entity"
	"github.com/derpartizanen/gophkeeper/internal/keeperd/postgres"
)

var _ Users = (*UsersRepo)(nil)

// UsersRepo is facade to users stored in Postgres.
type UsersRepo struct {
	pg *postgres.Postgres
}

// NewUsersRepo creates and initializes UsersRepo object.
func NewUsersRepo(
	pg *postgres.Postgres,
) *UsersRepo {
	return &UsersRepo{pg}
}

// Register creates a new user.
func (r *UsersRepo) Register(
	ctx context.Context,
	username, securityKey string,
) (uuid.UUID, error) {
	var id uuid.UUID

	fn := func(tx postgres.Transaction) error {
		err := tx.QueryRow(
			ctx,
			`INSERT INTO
           users (username, security_key)
       VALUES
           ($1, crypt($2, gen_salt('bf', 8)))
       RETURNING user_id`,
			username,
			securityKey,
		).Scan(&id)
		if err != nil {
			if postgres.IsEntityExists(err) {
				return entity.ErrUserExists
			}

			return fmt.Errorf("UsersRepo - Register - tx.QueryRow.Scan: %w", err)
		}

		return nil
	}

	if err := r.pg.RunAtomic(ctx, fn); err != nil {
		return id, fmt.Errorf("UsersRepo - Register - r.pg.RunAtomic: %w", err)
	}

	return id, nil
}

// Verify checks provided username and security key against data stored in database.
// Returns entity.User, if verification was successful.
func (r *UsersRepo) Verify(
	ctx context.Context,
	username, securityKey string,
) (entity.User, error) {
	var user entity.User

	err := r.pg.Pool.
		QueryRow(
			ctx,
			`SELECT
           user_id, username
       FROM
           users
       WHERE username=$1 AND security_key = crypt($2, security_key)`,
			username,
			securityKey,
		).
		Scan(&user.ID, &user.Username)
	if err != nil {
		if postgres.IsEmptyResponse(err) {
			return user, entity.ErrInvalidCredentials
		}

		return user, fmt.Errorf("UsersRepo - Verify - r.pg.Pool.QueryRow.Scan: %w", err)
	}

	return user, nil
}
