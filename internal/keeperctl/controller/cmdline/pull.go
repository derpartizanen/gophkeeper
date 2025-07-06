package cmdline

import (
	"github.com/cheynewallace/tabby"
	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/derpartizanen/gophkeeper/internal/keeperctl/app"
	"github.com/derpartizanen/gophkeeper/internal/keeperctl/errors"
	"github.com/derpartizanen/gophkeeper/proto"
)

var pullCmd = &cobra.Command{
	Use:   "pull [secret id] [flags]",
	Short: "Show the secret and stored data",
	Args:  cobra.MinimumNArgs(1),
	RunE:  doPull,
}

func init() {
	rootCmd.AddCommand(pullCmd)
}

func doPull(cmd *cobra.Command, args []string) error {
	id, err := uuid.Parse(args[0])
	if err != nil {
		return err
	}

	clientApp, err := app.FromContext(cmd.Context())
	if err != nil {
		return err
	}

	secret, data, err := clientApp.Services.Secrets.Get(cmd.Context(), clientApp.AccessToken, id)
	if err != nil {
		clientApp.Log.Debug().Err(err).Msg("")

		return errors.Unwrap(err)
	}

	header := []any{"ID", "Name", "Kind", "Description"}
	line := []any{
		secret.GetId(),
		secret.GetName(),
		secret.GetKind().String(),
		string(secret.GetMetadata()),
	}
	messages := make([]string, 0)

	switch d := data.(type) {
	case *proto.Binary:
		messages = append(messages, string(d.GetBinary()))

	case *proto.Card:
		header = append(header, "Number", "Expiration", "Holder", "CVV")
		line = append(line, d.GetNumber(), d.GetExpiration(), d.GetHolder(), d.GetCvv())

	case *proto.Credentials:
		header = append(header, "Login", "Password")
		line = append(line, d.GetLogin(), d.GetPassword())

	case *proto.Text:
		messages = append(messages, d.GetText())
	}

	t := tabby.New()
	t.AddHeader(header...)
	t.AddLine(line...)
	t.Print()

	for _, msg := range messages {
		clientApp.Log.Info().Msg(msg)
	}

	return nil
}
