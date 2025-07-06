package entity

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type userKey string

const (
	userKeyName userKey = "owner"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or security key")
	ErrUserExists         = errors.New("user already exists")
)

// User represents basic user of the system.
type User struct {
	ID       uuid.UUID `db:"user_id"`
	Username string
}

// WithContext injects user info into context.
func (u User) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, userKeyName, &u)
}

// UserFromContext extracts User object from context.
// Returns nil, if there is no info regarding user in the context.
func UserFromContext(ctx context.Context) *User {
	if val := ctx.Value(userKeyName); val != nil {
		return val.(*User)
	}

	return nil
}
