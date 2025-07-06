package editcmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/derpartizanen/gophkeeper/internal/keeperctl/service"
)

var (
	text string

	textCmd = &cobra.Command{
		Use:     "text [secret id] [flags]",
		Short:   "Edit text secret",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: preRun,
		RunE:    doEditText,
	}
)

func init() {
	textCmd.Flags().StringVarP(
		&text,
		"text",
		"t",
		"",
		"New text",
	)
}

func doEditText(cmd *cobra.Command, _args []string) error {
	if secretName == "" && description == "" && !noDescription && text == "" {
		return errFlagsRequired
	}

	if err := clientApp.Services.Secrets.EditText(
		cmd.Context(),
		clientApp.AccessToken,
		secretID,
		secretName,
		description,
		noDescription,
		text,
	); err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		if errors.Is(err, service.ErrKindMismatch) {
			return service.ErrKindMismatch
		}

		return errors.Unwrap(err)
	}

	return nil
}
