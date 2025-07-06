package entity

import (
	"errors"

	"github.com/google/uuid"

	"github.com/derpartizanen/gophkeeper/proto"
)

var (
	ErrSecretNotFound     = errors.New("secret not found")
	ErrSecretExists       = errors.New("secret already exists")
	ErrSecretNameConflict = errors.New("secret with such name already exists")
)

// Secret represents full secret info stored in the service.
type Secret struct {
	ID       uuid.UUID `db:"secret_id"`
	Name     string
	Kind     proto.DataKind
	Metadata []byte
	Data     []byte
}
