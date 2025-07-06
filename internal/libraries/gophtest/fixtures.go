package gophtest

import (
	"errors"

	"github.com/derpartizanen/gophkeeper/internal/libraries/creds"
)

const (
	Username                   = "admin"
	Password    creds.Password = "1q2w3e"
	SecurityKey                = "88bb5abaa61568b9f11ba091445d81772a3a264fb3f3054088f78baf7a091a9d"
	AccessToken                = "SomeLongTokenInJWT"
	Secret      creds.Password = "xxx"

	SecretName = "my-secret"
	Metadata   = "encrypted extra data"
	TextData   = "encrypted secret data"
)

var ErrUnexpected = errors.New("runtime error")
