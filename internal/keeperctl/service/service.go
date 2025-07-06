package service

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"github.com/derpartizanen/gophkeeper/internal/keeperctl/encryption"
	"github.com/derpartizanen/gophkeeper/internal/keeperctl/repo"
	p "github.com/derpartizanen/gophkeeper/proto"
)

type Auth interface {
	Login(ctx context.Context, username string, key encryption.Key) (string, error)
}

type Secrets interface {
	//todo: split to multiple interfaces
	PushBinary(ctx context.Context, token, name, description string, binary []byte) (uuid.UUID, error)
	PushCard(ctx context.Context, token, name, description string, number, expiration, holder string, cvv int32) (uuid.UUID, error)
	PushCreds(ctx context.Context, token, name, description, login, password string) (uuid.UUID, error)
	PushText(ctx context.Context, token, name, description, text string) (uuid.UUID, error)
	List(ctx context.Context, token string) ([]*p.Secret, error)
	Get(ctx context.Context, token string, id uuid.UUID) (*p.Secret, proto.Message, error)
	EditBinary(ctx context.Context, token string, id uuid.UUID, name, description string, noDescription bool, binary []byte) error
	EditCard(ctx context.Context, token string, id uuid.UUID, name, description string, noDescription bool, number, expiration, holder string, cvv int32) error
	EditCreds(ctx context.Context, token string, id uuid.UUID, name, description string, noDescription bool, login, password string) error
	EditText(ctx context.Context, token string, id uuid.UUID, name, description string, noDescription bool, text string) error
	Delete(ctx context.Context, token string, id uuid.UUID) error
}

type Users interface {
	Register(ctx context.Context, username string, key encryption.Key) (string, error)
}

// Services is a collection of business logic.
type Services struct {
	Auth    Auth
	Secrets Secrets
	Users   Users
}

// New creates and initializes collection of services.
func New(key encryption.Key, repos *repo.Repositories) *Services {
	return &Services{
		Auth:    NewAuthService(repos.Auth),
		Secrets: NewSecretsService(key, repos.Secrets),
		Users:   NewUsersService(repos.Users),
	}
}
