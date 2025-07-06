// Package app implements keeperctl service.
package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/derpartizanen/gophkeeper/internal/keeperctl/config"
	"github.com/derpartizanen/gophkeeper/internal/keeperctl/encryption"
	"github.com/derpartizanen/gophkeeper/internal/keeperctl/grpcconn"
	"github.com/derpartizanen/gophkeeper/internal/keeperctl/repo"
	"github.com/derpartizanen/gophkeeper/internal/keeperctl/service"
	"github.com/derpartizanen/gophkeeper/internal/logger"
)

type appKey string

var (
	appKeyName appKey = "gophkeeper"

	ErrNotInitialized = errors.New("application is not initialized")
)

// App implements keeperctl service.
type App struct {
	AccessToken   string
	conn          *grpcconn.Connection
	EncryptionKey encryption.Key
	Log           *logger.Logger
	Services      *service.Services
}

// New creates new App object.
func New(cfg *config.Config) (*App, error) {
	var loglevel string
	if cfg.Verbose {
		loglevel = zerolog.DebugLevel.String()
	} else {
		loglevel = zerolog.InfoLevel.String()
	}

	log, err := logger.New(loglevel)
	if err != nil {
		return nil, fmt.Errorf("logger error: %w", err)
	}

	log.Debug().Msg(cfg.String())

	conn, err := grpcconn.New(cfg.Address, cfg.CAPath)
	if err != nil {
		return nil, fmt.Errorf("grpc connection error: %w", err)
	}

	key := encryption.NewKey(cfg.Username, cfg.Password)
	repos := repo.New(conn)
	services := service.New(key, repos)

	return &App{
		Log:      log,
		conn:     conn,
		Services: services,
	}, nil
}

// WithContext injects App into provided context.
func (a *App) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, appKeyName, a)
}

// FromContext extracts App from provided context.
func FromContext(ctx context.Context) (*App, error) {
	if val := ctx.Value(appKeyName); val != nil {
		return val.(*App), nil
	}

	return nil, ErrNotInitialized
}

// Shutdown gracefully stops App.
func (a *App) Shutdown() {
	if err := a.conn.Close(); err != nil {
		a.Log.Warn().Err(err).Msg("grpc connection close error")
	}
}
