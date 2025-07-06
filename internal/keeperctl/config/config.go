package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/derpartizanen/gophkeeper/internal/libraries/creds"
)

// Config is configuration of keeperctl.
type Config struct {
	Username string
	Password creds.Password
	Address  string
	CAPath   string
	Verbose  bool
}

// New create application config from ENVs and cmd flags.
func New() *Config {
	viper.SetDefault("address", "127.0.0.1:9090")
	viper.SetDefault("verbose", false)

	viper.SetEnvPrefix("GOPH")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	cfg := &Config{
		Username: viper.GetString("username"),
		Password: creds.Password(viper.GetString("password")),
		Address:  viper.GetString("address"),
		CAPath:   viper.GetString("ca-path"),
		Verbose:  viper.GetBool("verbose"),
	}

	return cfg
}

func (c *Config) String() string {
	var sb strings.Builder

	sb.WriteString("Config:\n")
	sb.WriteString(fmt.Sprintf("\t\tUsername: %s\n", c.Username))
	sb.WriteString(fmt.Sprintf("\t\tPassword: %s\n", c.Password))
	sb.WriteString(fmt.Sprintf("\t\tAddress: %s\n", c.Address))
	sb.WriteString(fmt.Sprintf("\t\tCA path: %s\n", c.CAPath))
	sb.WriteString(fmt.Sprintf("\t\tVerbose: %t", c.Verbose))

	return sb.String()
}
