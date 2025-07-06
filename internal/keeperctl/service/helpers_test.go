package service_test

import (
	"github.com/derpartizanen/gophkeeper/internal/keeperctl/encryption"
	"github.com/derpartizanen/gophkeeper/internal/libraries/gophtest"
)

func newTestKey() encryption.Key {
	return encryption.NewKey(gophtest.Username, gophtest.Password)
}
