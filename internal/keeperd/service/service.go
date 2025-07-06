package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/derpartizanen/gophkeeper/internal/keeperd/config"
	"github.com/derpartizanen/gophkeeper/internal/keeperd/entity"
	"github.com/derpartizanen/gophkeeper/internal/keeperd/repo"
	"github.com/derpartizanen/gophkeeper/proto"
)

type Auth interface {
	Login(ctx context.Context, username, securityKey string) (entity.AccessToken, error)
}

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
	Register(ctx context.Context, username, securityKey string) (entity.AccessToken, error)
}

// Services is a collection of business logic.
type Services struct {
	Auth    Auth
	Secrets Secrets
	Users   Users
}

// New creates and initializes collection of business logic.
func New(cfg *config.Config, repos *repo.Repositories) *Services {
	return &Services{
		Auth:    NewAuthService(cfg.Secret, repos.Users),
		Secrets: NewSecretsService(repos.Secrets),
		Users:   NewUsersService(cfg.Secret, repos.Users),
	}
}
