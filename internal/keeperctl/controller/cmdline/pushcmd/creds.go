package pushcmd

import (
	"github.com/spf13/cobra"

	"github.com/derpartizanen/gophkeeper/internal/keeperctl/errors"
)

var (
	login    string
	password string

	credsCmd = &cobra.Command{
		Use:     "creds [flags]",
		Short:   "Save credentials",
		PreRunE: preRun,
		RunE:    doPushCreds,
	}
)

func init() {
	credsCmd.Flags().StringVarP(
		&login,
		"login",
		"l",
		"",
		"Login or username to save",
	)
	credsCmd.Flags().StringVarP(
		&password,
		"password",
		"p",
		"",
		"Password to save",
	)

	credsCmd.MarkFlagRequired("login")
	credsCmd.MarkFlagRequired("password")
}

func doPushCreds(cmd *cobra.Command, _args []string) error {
	id, err := clientApp.Services.Secrets.PushCreds(
		cmd.Context(),
		clientApp.AccessToken,
		secretName,
		description,
		login,
		password,
	)
	if err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		return errors.Unwrap(err)
	}

	clientApp.Log.Debug().Str("secret-id", id.String()).Msg("Secret saved successfully")

	return nil
}
