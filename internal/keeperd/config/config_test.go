package config_test

import (
	"os"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/require"

	"github.com/derpartizanen/gophkeeper/internal/keeperd/config"
)

func TestNewConfig(t *testing.T) {
	initialArgs := os.Args

	defer t.Cleanup(func() {
		os.Args = initialArgs
	})

	os.Args = []string{
		"",
		"--secret=xxx",
		"--crt-path=../../ssl/ca/keeper.crt",
		"--key-path=../../ssl/ca/keeper.key",
	}

	sat, err := config.New()

	require.NoError(t, err)
	snaps.MatchSnapshot(t, sat.String())
}

func TestNewConfigFailsIfSecretNotSet(t *testing.T) {
	_, err := config.New()

	require.ErrorIs(t, err, config.ErrSecretNotSet)
}

func TestNewConfigFailsIfCrtPathNotSet(t *testing.T) {
	initialArgs := os.Args

	defer t.Cleanup(func() {
		os.Args = initialArgs
	})

	os.Args = []string{
		"",
		"--secret=xxx",
	}

	_, err := config.New()

	require.ErrorIs(t, err, config.ErrCrtNotSet)
}

func TestNewConfigFailsIfCrtKeyPathNotSet(t *testing.T) {
	initialArgs := os.Args

	defer t.Cleanup(func() {
		os.Args = initialArgs
	})

	os.Args = []string{
		"",
		"--secret=xxx",
		"--crt-path=../../ssl/ca/keeper.crt",
	}

	_, err := config.New()

	require.ErrorIs(t, err, config.ErrCrtKeyNotSet)
}
