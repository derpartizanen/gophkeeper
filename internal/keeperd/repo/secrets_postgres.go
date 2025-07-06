package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/derpartizanen/gophkeeper/internal/keeperd/entity"
	"github.com/derpartizanen/gophkeeper/internal/keeperd/postgres"
	"github.com/derpartizanen/gophkeeper/proto"
)

var _ Secrets = (*SecretsRepo)(nil)

var ErrNoValuesToUpdate = errors.New("no values to update")

// SecretsRepo is facade to secrets stored in Postgres.
type SecretsRepo struct {
	pg *postgres.Postgres
}

// NewSecretsRepo creates and initializes SecretsRepo object.
func NewSecretsRepo(pg *postgres.Postgres) *SecretsRepo {
	return &SecretsRepo{pg}
}

// Create stores new secret in database.
func (r *SecretsRepo) Create(
	ctx context.Context,
	owner uuid.UUID,
	name string,
	kind proto.DataKind,
	metadata, data []byte,
) (id uuid.UUID, err error) {
	fn := func(tx postgres.Transaction) error {
		err := tx.QueryRow(
			ctx,
			`INSERT INTO
           secrets (owner_id, name, kind, metadata, data)
       VALUES
           ($1, $2, $3, $4, $5)
       RETURNING secret_id`,
			owner,
			name,
			kind,
			metadata,
			data,
		).Scan(&id)
		if err != nil {
			if postgres.IsEntityExists(err) {
				return entity.ErrSecretExists
			}

			return fmt.Errorf("SecretsRepo - Create - tx.QueryRow.Scan: %w", err)
		}

		return nil
	}

	if err := r.pg.RunAtomic(ctx, fn); err != nil {
		return id, fmt.Errorf("SecretsRepo - Create - r.pg.RunAtomic: %w", err)
	}

	return id, nil
}

// List returns all secrets of the provided user.
// Data is not filled in this case to reduce load on service.
func (r *SecretsRepo) List(
	ctx context.Context,
	owner uuid.UUID,
) ([]entity.Secret, error) {
	rv := make([]entity.Secret, 0)
	if err := r.pg.Select(
		ctx,
		&rv,
		`SELECT
         secret_id, name, kind, metadata
     FROM
         secrets
     WHERE owner_id = $1`,
		owner,
	); err != nil {
		return nil, fmt.Errorf("SecretsRepo - List - r.Select: %w", err)
	}

	return rv, nil
}

// Get returns full secret info and data.
func (r *SecretsRepo) Get(
	ctx context.Context,
	owner, id uuid.UUID,
) (*entity.Secret, error) {
	var secret entity.Secret

	err := r.pg.Pool.
		QueryRow(
			ctx,
			`SELECT
           secret_id, name, kind, metadata, data
       FROM
           secrets
       WHERE secret_id=$1 AND owner_id = $2`,
			id,
			owner,
		).
		Scan(&secret.ID, &secret.Name, &secret.Kind, &secret.Metadata, &secret.Data)
	if err != nil {
		if postgres.IsEmptyResponse(err) {
			return nil, entity.ErrSecretNotFound
		}

		return nil, fmt.Errorf("SecretsRepo - Get - r.pg.Pool.QueryRow.Scan: %w", err)
	}

	return &secret, nil
}

// Update changes secret info and data.
func (r *SecretsRepo) Update(
	ctx context.Context,
	owner, id uuid.UUID,
	changed []string,
	name string,
	metadata, data []byte,
) error {
	fn := func(tx postgres.Transaction) error {
		qb := newQueryBuilder("UPDATE secrets").Set()

		for _, field := range changed {
			switch field {
			case "name":
				qb.Append("name", "=", name)

			case "metadata":
				qb.Append("metadata", "=", metadata)

			case "data":
				qb.Append("data", "=", data)
			}
		}

		if len(qb.Values()) == 0 {
			return fmt.Errorf("SecretsRepo - Update: %w", ErrNoValuesToUpdate)
		}

		qb.Where().
			Append("secret_id", "=", id).
			And().
			Append("owner_id", "=", owner)

		tag, err := tx.Exec(ctx, qb.Query(), qb.Values()...)
		if err != nil {
			if postgres.IsEntityExists(err) {
				return entity.ErrSecretNameConflict
			}

			return fmt.Errorf("SecretsRepo - Update - tx.Exec: %w", err)
		}

		if tag.RowsAffected() == 0 {
			return entity.ErrSecretNotFound
		}

		return nil
	}

	if err := r.pg.RunAtomic(ctx, fn); err != nil {
		return fmt.Errorf("SecretsRepo - Update - r.pg.RunAtomic: %w", err)
	}

	return nil
}

// Delete removes secret from database.
func (r *SecretsRepo) Delete(
	ctx context.Context,
	owner, id uuid.UUID,
) (err error) {
	fn := func(tx postgres.Transaction) error {
		tag, err := tx.Exec(
			ctx,
			`DELETE FROM
           secrets
       WHERE secret_id = $1 AND owner_id = $2`,
			id,
			owner,
		)
		if err != nil {
			return fmt.Errorf("SecretsRepo - Delete - tx.Exec: %w", err)
		}

		if tag.RowsAffected() == 0 {
			return entity.ErrSecretNotFound
		}

		return nil
	}

	if err := r.pg.RunAtomic(ctx, fn); err != nil {
		return fmt.Errorf("SecretsRepo - Delete - r.pg.RunAtomic: %w", err)
	}

	return nil
}
