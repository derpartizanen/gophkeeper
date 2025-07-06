// Package repo provides facade to data stored in external sources.
package repo

import (
	"context"

	"github.com/derpartizanen/gophkeeper/internal/keeperd/entity"
	"github.com/derpartizanen/gophkeeper/internal/keeperd/postgres"
	"github.com/derpartizanen/gophkeeper/proto"

	"github.com/google/uuid"
)

type Secrets interface {
	Create(
		ctx context.Context,
		owner uuid.UUID,
		name string,
		kind proto.DataKind,
		metadata, data []byte,
	) (uuid.UUID, error)

	List(ctx context.Context, owner uuid.UUID) ([]entity.Secret, error)
	Get(ctx context.Context, owner, id uuid.UUID) (*entity.Secret, error)

	Update(
		ctx context.Context,
		owner, id uuid.UUID,
		changed []string,
		name string,
		metadata []byte,
		data []byte,
	) error

	Delete(ctx context.Context, owner, id uuid.UUID) error
}

type Users interface {
	Register(ctx context.Context, username, securityKey string) (uuid.UUID, error)
	Verify(ctx context.Context, username, securityKey string) (entity.User, error)
}

// Repositories is a collection of data repositories.
type Repositories struct {
	Secrets Secrets
	Users   Users
}

// New creates and initializes collection of data repositories.
func New(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		Secrets: NewSecretsRepo(pg),
		Users:   NewUsersRepo(pg),
	}
}
