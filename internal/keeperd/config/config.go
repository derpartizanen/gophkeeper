// Package config provides configuration for the keeperd service.
package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/derpartizanen/gophkeeper/internal/libraries/creds"
)

var (
	ErrSecretNotSet = errors.New("secret key required")
	ErrCrtKeyNotSet = errors.New("certificate key required")
	ErrCrtNotSet    = errors.New("service certificate required")
)

type Config struct {
	Address     string
	DatabaseURI creds.ConnURI
	Secret      creds.Password
	CrtPath     string
	KeyPath     string
	LogLevel    string
}

// Validate verifies values stored in resulting config.
func validate(cfg *Config) error {
	if cfg.Secret == "" {
		return ErrSecretNotSet
	}

	if cfg.CrtPath == "" {
		return ErrCrtNotSet
	}

	if cfg.KeyPath == "" {
		return ErrCrtKeyNotSet
	}

	return nil
}

// New create application config by reading environment variables and
// commandline flags. The environment variables has first priority.
func New() (*Config, error) {
	viper.SetDefault("address", "0.0.0.0:9090")
	viper.SetDefault(
		"database-uri",
		"postgres://postgres:postgres@127.0.0.1:5432/goph?sslmode=disable",
	)

	flag := pflag.FlagSet{}
	flag.String("address", "", "address:port the service listens on")
	flag.String("database-uri", "", "full Postgres database connection URI")
	flag.String("secret", "", "secret key to sign JWT tokens")
	flag.String("crt-path", "", "path to server certificate")
	flag.String("key-path", "", "path to server key certificate")
	flag.String("log-level", "info", "log level of the service (info, warn, error, debug)")

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	if err := flag.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	if err := viper.BindPFlags(&flag); err != nil {
		return nil, err
	}

	cfg := &Config{
		Address:     viper.GetString("address"),
		DatabaseURI: creds.ConnURI(viper.GetString("database-uri")),
		Secret:      creds.Password(viper.GetString("secret")),
		CrtPath:     viper.GetString("crt-path"),
		KeyPath:     viper.GetString("key-path"),
		LogLevel:    viper.GetString("log-level"),
	}

	if err := validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) String() string {
	var sb strings.Builder

	sb.WriteString("Configuration:\n")
	sb.WriteString(fmt.Sprintf("\t\tListening address: %s\n", c.Address))
	sb.WriteString(fmt.Sprintf("\t\tDatabase URI: %s\n", c.DatabaseURI))
	sb.WriteString(fmt.Sprintf("\t\tSecret: %s\n", c.Secret))
	sb.WriteString(fmt.Sprintf("\t\tCertificate path: %s\n", c.CrtPath))
	sb.WriteString(fmt.Sprintf("\t\tCertificate key path: %s\n", c.KeyPath))
	sb.WriteString(fmt.Sprintf("\t\tLog level: %s", c.LogLevel))

	return sb.String()
}
