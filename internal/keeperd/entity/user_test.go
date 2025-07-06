package entity_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/derpartizanen/gophkeeper/internal/keeperd/entity"
	"github.com/derpartizanen/gophkeeper/internal/libraries/gophtest"
)

func TestUserWithFromContext(t *testing.T) {
	expected := entity.User{
		ID:       uuid.New(),
		Username: gophtest.Username,
	}

	ctx := expected.WithContext(context.Background())

	require.Equal(t, expected, *entity.UserFromContext(ctx))
}

func TestUserFromCleanContext(t *testing.T) {
	require.Nil(t, entity.UserFromContext(context.Background()))
}
