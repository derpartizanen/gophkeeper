package pushcmd

import (
	"github.com/spf13/cobra"

	"github.com/derpartizanen/gophkeeper/internal/keeperctl/app"
)

var (
	clientApp *app.App

	secretName  string
	description string

	PushCmd = &cobra.Command{
		Use:   "push",
		Short: "Push secret to the Keeperd service",
	}
)

func init() {
	PushCmd.PersistentFlags().StringVarP(
		&secretName,
		"name",
		"n",
		"",
		"Name of the stored secret",
	)
	PushCmd.PersistentFlags().StringVarP(
		&description,
		"description",
		"d",
		"",
		"Additional description of stored data (activation codes, names of banks etc)",
	)

	PushCmd.MarkPersistentFlagRequired("name")

	PushCmd.AddCommand(binCmd)
	PushCmd.AddCommand(cardCmd)
	PushCmd.AddCommand(credsCmd)
	PushCmd.AddCommand(textCmd)
}

// preRun executes preparational operations common for all sub commands.
func preRun(cmd *cobra.Command, _args []string) error {
	var err error

	clientApp, err = app.FromContext(cmd.Context())

	return err
}
