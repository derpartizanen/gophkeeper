package editcmd

import (
	"errors"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/derpartizanen/gophkeeper/internal/keeperctl/app"
)

var (
	errFlagsRequired = errors.New("at least one flag required")

	clientApp *app.App

	secretID      uuid.UUID
	secretName    string
	description   string
	noDescription bool
)

var EditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit secret data stored in the keeperd service",
}

func init() {
	EditCmd.PersistentFlags().StringVarP(
		&secretName,
		"name",
		"n",
		"",
		"New name of the secret",
	)
	EditCmd.PersistentFlags().StringVarP(
		&description,
		"description",
		"d",
		"",
		"New description of secret (activation codes, names of banks etc)",
	)
	EditCmd.PersistentFlags().BoolVar(
		&noDescription,
		"no-description",
		false,
		"Remove description from the secret",
	)

	EditCmd.MarkFlagsMutuallyExclusive("description", "no-description")

	EditCmd.AddCommand(binCmd)
	EditCmd.AddCommand(cardCmd)
	EditCmd.AddCommand(credsCmd)
	EditCmd.AddCommand(textCmd)
}

// preRun executes preparation operations common for all sub commands.
func preRun(cmd *cobra.Command, args []string) error {
	var err error

	secretID, err = uuid.Parse(args[0])
	if err != nil {
		return err
	}

	clientApp, err = app.FromContext(cmd.Context())

	return err
}
