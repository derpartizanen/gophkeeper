package cmdline

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/derpartizanen/gophkeeper/internal/keeperctl/app"
	"github.com/derpartizanen/gophkeeper/internal/keeperctl/config"
	"github.com/derpartizanen/gophkeeper/internal/keeperctl/controller/cmdline/editcmd"
	"github.com/derpartizanen/gophkeeper/internal/keeperctl/controller/cmdline/pushcmd"
)

var (
	cfg *config.Config

	verbose  bool
	address  string
	caPath   string
	username string
	password string

	rootCmd = &cobra.Command{
		Use:               "keeperctl",
		Short:             "keeperctl is an interactive commandline client for the keeperd service",
		PersistentPreRunE: initApp,
		PersistentPostRun: finalizeApp,
	}
)

// Execute executes the root command.
func Execute(buildVersion, buildDate string) error {
	rootCmd.Version = fmt.Sprintf("%s (%s)", buildVersion, buildDate)

	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initializeConfig)

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().StringVar(
		&address,
		"address",
		"",
		"Address and port of the keeperd service",
	)
	rootCmd.PersistentFlags().StringVar(
		&caPath,
		"ca-path",
		"",
		"Path to certificate authority to verify server certificate",
	)
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "Name of a user")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Master password")

	rootCmd.MarkFlagRequired("username")
	rootCmd.MarkFlagRequired("password")

	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("address", rootCmd.PersistentFlags().Lookup("address"))
	viper.BindPFlag("ca-path", rootCmd.PersistentFlags().Lookup("ca-path"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.AddCommand(pushcmd.PushCmd)
	rootCmd.AddCommand(editcmd.EditCmd)
}

// initializeConfig does initialization routine before reading commandline flags.
func initializeConfig() {
	cfg = config.New()
}

// initApp does initialization routine before reading commandline flags.
func initApp(cmd *cobra.Command, args []string) error {
	// Prerun is executed for EVERY command, even for help.
	if cmd.Name() == "help" {
		return nil
	}

	cfg := config.New()

	clientApp, err := app.New(cfg)
	if err != nil {
		return err
	}

	cmd.SetContext(clientApp.WithContext(cmd.Context()))

	if cmd.Name() == "register" {
		return nil
	}

	return login(cmd, args)
}

// finalizeApp does cleanup at the end of commandline application.
func finalizeApp(cmd *cobra.Command, _ []string) {
	clientApp, err := app.FromContext(cmd.Context())
	if err == nil {
		clientApp.Shutdown()
	}
}
