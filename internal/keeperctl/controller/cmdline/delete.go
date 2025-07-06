package cmdline

import (
	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/derpartizanen/gophkeeper/internal/keeperctl/app"
	"github.com/derpartizanen/gophkeeper/internal/keeperctl/errors"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [secret id] [flags]",
	Short: "Delete the secret",
	Args:  cobra.MinimumNArgs(1),
	RunE:  doDelete,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func doDelete(cmd *cobra.Command, args []string) error {
	id, err := uuid.Parse(args[0])
	if err != nil {
		return err
	}

	clientApp, err := app.FromContext(cmd.Context())
	if err != nil {
		return err
	}

	if err := clientApp.Services.Secrets.Delete(cmd.Context(), clientApp.AccessToken, id); err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		return errors.Unwrap(err)
	}

	return nil
}
