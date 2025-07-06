package pushcmd

import (
	"github.com/spf13/cobra"

	"github.com/derpartizanen/gophkeeper/internal/keeperctl/errors"
)

var (
	number     string
	expiration string
	holder     string
	cvv        int32

	cardCmd = &cobra.Command{
		Use:     "card [flags]",
		Short:   "Save bank card info",
		PreRunE: preRun,
		RunE:    doPushCard,
	}
)

func init() {
	cardCmd.Flags().StringVar(
		&number,
		"number",
		"",
		"Card number",
	)
	cardCmd.Flags().StringVar(
		&expiration,
		"expiration",
		"",
		"Card expiration date",
	)
	cardCmd.Flags().StringVar(
		&holder,
		"holder",
		"",
		"Card holder name and surname",
	)
	cardCmd.Flags().Int32Var(
		&cvv,
		"cvv",
		0,
		"Card verification value",
	)

	cardCmd.MarkFlagRequired("number")
	cardCmd.MarkFlagRequired("expiration")
	cardCmd.MarkFlagRequired("holder")
	cardCmd.MarkFlagRequired("cvv")
}

func doPushCard(cmd *cobra.Command, _args []string) error {
	id, err := clientApp.Services.Secrets.PushCard(
		cmd.Context(),
		clientApp.AccessToken,
		secretName,
		description,
		number,
		expiration,
		holder,
		cvv,
	)
	if err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		return errors.Unwrap(err)
	}

	clientApp.Log.Debug().Str("secret-id", id.String()).Msg("Secret saved successfully")

	return nil
}
