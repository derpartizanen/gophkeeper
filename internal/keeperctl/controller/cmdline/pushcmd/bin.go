package pushcmd

import (
	"github.com/spf13/cobra"

	"github.com/derpartizanen/gophkeeper/internal/keeperctl/errors"
)

var (
	data []byte

	binCmd = &cobra.Command{
		Use:     "bin [flags]",
		Short:   "Save arbitrary binary data",
		PreRunE: preRun,
		RunE:    doPushBinary,
	}
)

func init() {
	binCmd.Flags().BytesHexVarP(
		&data,
		"binary-data",
		"b",
		nil,
		"Binary data in hex format",
	)

	binCmd.MarkFlagRequired("data")
}

func doPushBinary(cmd *cobra.Command, _args []string) error {
	id, err := clientApp.Services.Secrets.PushBinary(
		cmd.Context(),
		clientApp.AccessToken,
		secretName,
		description,
		data,
	)
	if err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		return errors.Unwrap(err)
	}

	clientApp.Log.Debug().Str("secret-id", id.String()).Msg("Secret saved successfully")

	return nil
}
