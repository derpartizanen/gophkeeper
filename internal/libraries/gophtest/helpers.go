package gophtest

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// CreateUUID create new UUID from the provided string.
func CreateUUID(t *testing.T, val string) uuid.UUID {
	t.Helper()

	v, err := uuid.Parse(val)

	require.NoError(t, err)

	return v
}
