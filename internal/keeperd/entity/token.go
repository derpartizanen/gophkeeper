package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/derpartizanen/gophkeeper/internal/libraries/creds"
)

const TokenLifeTime = 15 * time.Minute

var SigningMethod = jwt.SigningMethodHS256

// AccessToken is JWT token used for authentication.
type AccessToken string

// Claims contain token's payload with various info about token itself and authenticated user.
type Claims struct {
	jwt.RegisteredClaims
	Username string
}

// NewAccessToken issues new access token valid for limited period of time.
func NewAccessToken(user User, secret creds.Password) (AccessToken, error) {
	now := time.Now()

	claims := jwt.MapClaims{}

	claims["iss"] = "Goph"
	claims["jti"] = uuid.New()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()
	claims["exp"] = now.Add(TokenLifeTime).Unix()

	// User info
	claims["sub"] = user.ID
	claims["username"] = user.Username

	rawToken := jwt.NewWithClaims(SigningMethod, claims)

	signedToken, err := rawToken.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("AccessToken - NewAccessToken - rawToken.SignedString: %w", err)
	}

	return AccessToken(signedToken), nil
}

// TokenFromString create new AccessToken from provided string.
// Supports string in form of "Bearer xxx" or raw tokens.
func TokenFromString(src string) AccessToken {
	rawToken := strings.TrimSpace(src)

	return AccessToken(strings.TrimPrefix(rawToken, "Bearer "))
}

// String converts AccessToken to string.
func (t AccessToken) String() string {
	return string(t)
}

// Decode decodes token, verifies its signature and return claims if the token is valid.
func (t AccessToken) Decode(secret creds.Password) (*Claims, error) {
	claims := new(Claims)

	if _, err := jwt.ParseWithClaims(
		t.String(),
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	); err != nil {
		return nil, err
	}

	return claims, nil
}
