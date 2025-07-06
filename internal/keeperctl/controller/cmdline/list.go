package cmdline

import (
	"github.com/cheynewallace/tabby"
	"github.com/spf13/cobra"

	"github.com/derpartizanen/gophkeeper/internal/keeperctl/app"
	"github.com/derpartizanen/gophkeeper/internal/keeperctl/errors"
)

var listCmd = &cobra.Command{
	Use:   "list [flags]",
	Short: "List secrets of current user (without data)",
	RunE:  doList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func doList(cmd *cobra.Command, args []string) error {
	clientApp, err := app.FromContext(cmd.Context())
	if err != nil {
		return err
	}

	data, err := clientApp.Services.Secrets.List(cmd.Context(), clientApp.AccessToken)
	if err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		return errors.Unwrap(err)
	}

	t := tabby.New()
	t.AddHeader("ID", "Name", "Kind", "Description")

	for _, secret := range data {
		t.AddLine(secret.GetId(), secret.GetName(), secret.Kind.String(), string(secret.GetMetadata()))
	}

	t.Print()

	return nil
}
