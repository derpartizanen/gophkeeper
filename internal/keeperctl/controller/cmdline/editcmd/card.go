package editcmd

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/derpartizanen/gophkeeper/internal/keeperctl/service"
)

var (
	number     string
	expiration string
	holder     string
	cvv        int32

	cardCmd = &cobra.Command{
		Use:     "card [secret id] [flags]",
		Short:   "Edit stored bank card secret",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: preRun,
		RunE:    doEditCard,
	}
)

func init() {
	cardCmd.Flags().StringVar(
		&number,
		"number",
		"",
		"New card number",
	)
	cardCmd.Flags().StringVar(
		&expiration,
		"expiration",
		"",
		"New card expiration date",
	)
	cardCmd.Flags().StringVar(
		&holder,
		"holder",
		"",
		"New card holder name and surname",
	)
	cardCmd.Flags().Int32Var(
		&cvv,
		"cvv",
		0,
		"New card verification value",
	)
}

func doEditCard(cmd *cobra.Command, _args []string) error {
	if secretName == "" && description == "" && !noDescription &&
		number == "" && expiration == "" && holder == "" && cvv == 0 {
		return errFlagsRequired
	}

	if err := clientApp.Services.Secrets.EditCard(
		cmd.Context(),
		clientApp.AccessToken,
		secretID,
		secretName,
		description,
		noDescription,
		number,
		expiration,
		holder,
		cvv,
	); err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		if errors.Is(err, service.ErrKindMismatch) {
			return service.ErrKindMismatch
		}

		return errors.Unwrap(err)
	}

	return nil
}
