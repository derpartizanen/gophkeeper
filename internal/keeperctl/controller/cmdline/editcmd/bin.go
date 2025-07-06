package editcmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/derpartizanen/gophkeeper/internal/keeperctl/service"
)

var (
	data []byte

	binCmd = &cobra.Command{
		Use:     "bin [secret id] [flags]",
		Short:   "Edit stored binary secret",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: preRun,
		RunE:    doEditBin,
	}
)

func init() {
	binCmd.Flags().BytesHexVarP(
		&data,
		"binary-data",
		"b",
		nil,
		"New binary data in hex format",
	)
}

func doEditBin(cmd *cobra.Command, _args []string) error {
	if secretName == "" && description == "" && !noDescription && len(data) == 0 {
		return errFlagsRequired
	}

	if err := clientApp.Services.Secrets.EditBinary(
		cmd.Context(),
		clientApp.AccessToken,
		secretID,
		secretName,
		description,
		noDescription,
		data,
	); err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		if errors.Is(err, service.ErrKindMismatch) {
			return service.ErrKindMismatch
		}

		return errors.Unwrap(err)
	}

	return nil
}
