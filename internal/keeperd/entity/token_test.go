package entity_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/derpartizanen/gophkeeper/internal/keeperd/entity"
	"github.com/derpartizanen/gophkeeper/internal/libraries/gophtest"
)

func TestAccessTokenEncodeDecode(t *testing.T) {
	user := entity.User{
		ID:       uuid.New(),
		Username: gophtest.Username,
	}

	token, err := entity.NewAccessToken(user, gophtest.Secret)
	require.NoError(t, err)

	claims, err := token.Decode(gophtest.Secret)
	require.NoError(t, err)

	require.Equal(t, user.ID.String(), claims.Subject)
	require.Equal(t, user.Username, claims.Username)
}

func TestAccessTokenDecodeWithWrongSecret(t *testing.T) {
	user := entity.User{
		ID:       uuid.New(),
		Username: gophtest.Username,
	}

	token, err := entity.NewAccessToken(user, gophtest.Secret)
	require.NoError(t, err)

	_, err = token.Decode("yyy")
	require.Error(t, err)
}

func TestAccesTokenFromString(t *testing.T) {
	tt := []struct {
		name     string
		src      string
		expected string
	}{
		{
			name:     "Bearer token",
			src:      "Bearer JWT-token-value",
			expected: "JWT-token-value",
		},
		{
			name:     "Plain token",
			src:      "JWT-token-value",
			expected: "JWT-token-value",
		},
		{
			name:     "Bearer token with trailing spaces",
			src:      "  Bearer JWT-token-value   ",
			expected: "JWT-token-value",
		},
		{
			name:     "Empty string",
			src:      "",
			expected: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sat := entity.TokenFromString(tc.src)

			require.Equal(t, tc.expected, sat.String())
		})
	}
}
