package cmdline

import (
	"github.com/spf13/cobra"

	"github.com/derpartizanen/gophkeeper/internal/keeperctl/app"
	"github.com/derpartizanen/gophkeeper/internal/keeperctl/errors"
)

func login(cmd *cobra.Command, _ []string) error {
	clientApp, err := app.FromContext(cmd.Context())
	if err != nil {
		return err
	}

	token, err := clientApp.Services.Auth.Login(
		cmd.Context(),
		cfg.Username,
		clientApp.EncryptionKey,
	)
	if err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		return errors.Unwrap(err)
	}

	clientApp.AccessToken = token
	clientApp.Log.Debug().
		Str("access-token", token).
		Msg("Login successful")

	return nil
}
